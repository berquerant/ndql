package tree

import (
	"github.com/berquerant/ndql/pkg/iterx"
	. "github.com/pingcap/tidb/pkg/parser/ast"
	"github.com/pingcap/tidb/pkg/parser/opcode"
)

func (v TreeVisitor) VisitUnaryOperationExpr(n *UnaryOperationExpr) (NFunction, error) {
	f, err := v.VisitExpr(n.V)
	if err != nil {
		return nil, v.newErr(err, n, "UnaryOperationExpr expr")
	}
	name, g, err := v.newUnaryOperationNodeDataMapper(n)
	if err != nil {
		return nil, err
	}
	return iterx.CombineFunction(f, ReturnContainerValue(name, func(x ND) (ND, error) {
		r, err := g(x.AsOp())
		if err != nil {
			return nil, err
		}
		return r.AsData(), nil
	}))
}

func (v TreeVisitor) newUnaryOperationNodeDataMapper(n *UnaryOperationExpr) (string, func(*OP) (*OP, error), error) {
	switch n.Op {
	case opcode.Minus, opcode.Not: // -
		return "Not", func(v *OP) (*OP, error) {
			return v.Not()
		}, nil
	case opcode.BitNeg: // ~
		return "BitNot", func(v *OP) (*OP, error) {
			return v.BitNot()
		}, nil
	default:
		return "", nil, v.notImplemented(n, "unknown UnaryOperationExpr opcode")
	}
}
