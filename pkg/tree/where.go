package tree

import (
	"log/slog"

	"github.com/berquerant/ndql/pkg/iterx"
	"github.com/berquerant/ndql/pkg/logx"
	. "github.com/pingcap/tidb/pkg/parser/ast"
)

func (v TreeVisitor) VisitWhere(n ExprNode) (NFunction, error) {
	f, err := v.VisitExpr(n)
	if err != nil {
		return nil, v.newErr(err, n, "Where failed to eval")
	}
	if f.RetArity() != iterx.Unary {
		return nil, v.newErr(ErrInvalidFunctionArity, n, "Where ret should be unary")
	}
	return iterx.NewMapFunction(func(x *N) (*N, error) {
		r, err := f.CallAny(x)
		if err != nil {
			return nil, v.newErr(err, n, "Where failed to eval")
		}
		_, rv, ok := AsValueContainer(r[0]).GetFirstValue()
		if !ok {
			slog.Warn("Where got no value", logx.Verbose("node", n), logx.Verbose("input", x))
			return nil, ErrIgnore
		}
		b, err := rv.AsOp().AsBool()
		if err != nil {
			slog.Warn("Where got not Bool value", logx.Verbose("node", n), logx.Verbose("input", x), logx.Verbose("value", rv), logx.Err(err))
			return nil, ErrIgnore
		}
		if !b.Raw() { // ignore the node because expr was evaluated as false
			return nil, ErrIgnore
		}
		return x, nil
	}), nil
}
