package tree

import (
	"fmt"
	"log/slog"

	"github.com/berquerant/ndql/pkg/iterx"
	"github.com/berquerant/ndql/pkg/logx"
	"github.com/berquerant/ndql/pkg/node"
	. "github.com/pingcap/tidb/pkg/parser/ast"
)

func (v TreeVisitor) VisitCaseExpr(n *CaseExpr) (NFunction, error) {
	cases, err := v.evalCaseExprWhenCaluses(n)
	if err != nil {
		return nil, v.newErr(err, n, "CaseExpr.Cases")
	}
	elseVal, err := v.evalCaseExprElse(n)
	if err != nil {
		return nil, v.newErr(err, n, "CaseExpr.Else")
	}
	if n.Value != nil { // CASE value WHEN ...
		return v.visitCaseExprWithValue(n, cases, elseVal)
	}
	// CASE WHEN ...
	return v.visitCaseExprWithoutValue(n, cases, elseVal)
}

func (v TreeVisitor) validateCaseExprFunction(n NFunction) error {
	if n.RetArity() != iterx.Unary {
		return fmt.Errorf("%w: CaseExpr requires unary ret", ErrInvalidFunctionArity)
	}
	return nil
}

func (v TreeVisitor) evalAndValidateCaseExpr(n ExprNode) (NFunction, error) {
	f, err := v.VisitExpr(n)
	if err != nil {
		return nil, err
	}
	if err := v.validateCaseExprFunction(f); err != nil {
		return nil, err
	}
	return f, nil
}

func (v TreeVisitor) evalCaseExprElse(n *CaseExpr) (NFunction, error) {
	if x := n.ElseClause; x != nil {
		return v.evalAndValidateCaseExpr(x)
	}
	return ReturnContainerValue("CaseExprElseClause", func(_ ND) (ND, error) {
		return node.NewNull(), nil
	}), nil
}

type caseExprWhenCaluse struct {
	expr, result NFunction
}

func (v TreeVisitor) evalCaseExprWhenCaluses(n *CaseExpr) ([]*caseExprWhenCaluse, error) {
	r := make([]*caseExprWhenCaluse, len(n.WhenClauses))
	for i, x := range n.WhenClauses {
		expr, err := v.evalAndValidateCaseExpr(x.Expr)
		if err != nil {
			return nil, v.newErr(err, n, "CaseExpr WhenCaluse[%d].Expr", i)
		}
		result, err := v.evalAndValidateCaseExpr(x.Result)
		if err != nil {
			return nil, v.newErr(err, n, "CaseExpr WhenCaluse[%d].Result", i)
		}
		r[i] = &caseExprWhenCaluse{
			expr:   expr,
			result: result,
		}
	}
	return r, nil
}

func (v TreeVisitor) visitCaseExprWithValue(n *CaseExpr, cases []*caseExprWhenCaluse, elseVal NFunction) (NFunction, error) {
	val, err := v.evalAndValidateCaseExpr(n.Value)
	if err != nil {
		return nil, v.newErr(err, n, "CaseExpr value")
	}
	return iterx.NewMapFunction(func(x *N) (*N, error) {
		values, err := val.CallAny(x)
		if err != nil {
			return nil, v.newErr(err, n, "CaseExpr value eval")
		}
		_, value, ok := AsValueContainer(values[0]).GetFirstValue()
		if !ok {
			return nil, v.newErr(ErrInvalidValue, n, "CaseExpr value eval no value")
		}
		for i, c := range cases {
			logger := slog.With(logx.Verbose("clause", c), slog.Int("index", i))
			cvals, err := c.expr.CallAny(x)
			if err != nil {
				logger.Warn("CaseExpr WhenCaluse failed to eval Expr", logx.Err(err))
				continue
			}
			_, cval, ok := AsValueContainer(cvals[0]).GetFirstValue()
			if !ok {
				logger.Warn("CaseExpr WhenCaluse failed to eval Expr no value")
				continue
			}
			switch value.AsOp().Compare(cval.AsOp()) {
			case node.CmpUnknown:
				logger.Warn("CaseExpr WhenClause failed to compare values", logx.Verbose("value", value), logx.Verbose("clause", cval))
				continue
			case node.CmpEqual:
				rvals, err := c.result.CallAny(x)
				if err != nil {
					return nil, v.newErr(err, n, "CaseExpr WhenClause[%d] failed to eval Result", i)
				}
				return rvals[0], nil
			}
		}
		evals, err := elseVal.CallAny(x)
		if err != nil {
			return nil, v.newErr(err, n, "CaseExpr Else failed to eval")
		}
		return evals[0], nil
	}), nil
}

func (v TreeVisitor) visitCaseExprWithoutValue(n *CaseExpr, cases []*caseExprWhenCaluse, elseVal NFunction) (NFunction, error) {
	return iterx.NewMapFunction(func(x *N) (*N, error) {
		for i, c := range cases {
			logger := slog.With(logx.Verbose("clause", c), slog.Int("index", i))
			cvals, err := c.expr.CallAny(x)
			if err != nil {
				logger.Warn("CaseExpr WhenCaluse failed to eval Expr", logx.Err(err))
				continue
			}
			_, cval, ok := AsValueContainer(cvals[0]).GetFirstValue()
			if !ok {
				logger.Warn("CaseExpr WhenCaluse failed to eval Expr no value")
				continue
			}
			b, ok := cval.AsOp().Bool()
			if !ok {
				logger.Warn("CaseExpr WhenClause failed to eval Expr not Bool", logx.Verbose("expr", cval))
				continue
			}
			if b.Raw() {
				rvals, err := c.result.CallAny(x)
				if err != nil {
					return nil, v.newErr(err, n, "CaseExpr WhenClause[%d] failed to eval Result", i)
				}
				return rvals[0], nil
			}
		}
		evals, err := elseVal.CallAny(x)
		if err != nil {
			return nil, v.newErr(err, n, "CaseExpr Else failed to eval")
		}
		return evals[0], nil
	}), nil
}
