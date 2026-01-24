package tree

import (
	"github.com/berquerant/ndql/pkg/iterx"
	"github.com/berquerant/ndql/pkg/node"
	. "github.com/pingcap/tidb/pkg/parser/ast"
)

func (v TreeVisitor) VisitExpr(n ExprNode) (NFunction, error) {
	switch n := n.(type) {
	case *CaseExpr:
		return v.VisitCaseExpr(n)
	case *FuncCallExpr:
		return v.VisitFuncCallExpr(n)
	case ValueExpr:
		return v.VisitValueExpr(n)
	case *BetweenExpr:
		return v.VisitBetweenExpr(n)
	case *IsTruthExpr:
		return v.VisitIsTruthExpr(n)
	case *IsNullExpr:
		return v.VisitIsNullExpr(n)
	case *ParenthesesExpr:
		return v.VisitParenthesesExpr(n)
	case *PatternInExpr:
		return v.VisitPatternInExpr(n)
	case *PatternLikeOrIlikeExpr:
		return v.VisitParrernLikeExpr(n)
	case *PatternRegexpExpr:
		return v.VisitPatternRegexpExpr(n)
	case *BinaryOperationExpr:
		return v.VisitBinaryOperationExpr(n)
	case *UnaryOperationExpr:
		return v.VisitUnaryOperationExpr(n)
	case *ColumnNameExpr:
		return v.VisitColumnNameExpr(n)
	default:
		return nil, v.notImplemented(n, "unknown ExprNode")
	}
}

func (v TreeVisitor) VisitColumnNameExpr(n *ColumnNameExpr) (NFunction, error) {
	return v.VisitColumnName(n.Name).NFunction(), nil
}

func (TreeVisitor) VisitColumnName(n *ColumnName) *Key {
	return NewKey(n.Table.O, n.Name.O)
}

// (EXPR)
func (v TreeVisitor) VisitParenthesesExpr(n *ParenthesesExpr) (NFunction, error) {
	return v.VisitExpr(n.Expr)
}

// IS NULL, IS NOT NULL
func (v TreeVisitor) VisitIsNullExpr(n *IsNullExpr) (NFunction, error) {
	f, err := v.VisitExpr(n.Expr)
	if err != nil {
		return nil, v.newErr(err, n, "IsNullExpr expr")
	}
	return iterx.CombineFunction(f, ReturnContainerValue("IsNullExpr", func(x ND) (ND, error) {
		return node.Bool(x.AsOp().IsNull() != n.Not), nil
	}))
}

// IS TRUE, IS FALSE
func (v TreeVisitor) VisitIsTruthExpr(n *IsTruthExpr) (NFunction, error) {
	f, err := v.VisitExpr(n.Expr)
	if err != nil {
		return nil, v.newErr(err, n, "IsTruthExpr expr")
	}
	return iterx.CombineFunction(f, ReturnContainerValue("IsTruthExpr", func(x ND) (ND, error) {
		if n.True > 0 {
			return node.Bool(x.AsOp().IsTrue() != n.Not), nil
		}
		return node.Bool(x.AsOp().IsFalse() != n.Not), nil
	}))
}
