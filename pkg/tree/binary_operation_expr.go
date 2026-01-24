package tree

import (
	"errors"
	"fmt"

	"github.com/berquerant/ndql/pkg/iterx"
	"github.com/berquerant/ndql/pkg/node"
	"github.com/berquerant/ndql/pkg/util"
	. "github.com/pingcap/tidb/pkg/parser/ast"
	"github.com/pingcap/tidb/pkg/parser/opcode"
)

func (v TreeVisitor) VisitBinaryOperationExpr(n *BinaryOperationExpr) (NFunction, error) {
	f, err := v.VisitExpr(n.L)
	if err != nil {
		return nil, v.newErr(err, n, "BinaryOperationExpr: left")
	}
	g, err := v.VisitExpr(n.R)
	if err != nil {
		return nil, v.newErr(err, n, "BinaryOperationExpr: right")
	}
	if !(f.RetArity() == iterx.Unary && g.RetArity() == iterx.Unary) {
		return nil, v.invalidTree(n, "BinaryOperationExpr: left and right should be unary")
	}
	_, op, err := v.newBinaryOperation(n)
	if err != nil {
		return nil, v.newErr(err, n, "BinaryOperationExpr: invalid operation")
	}
	return iterx.NewMapFunction(func(x *N) (*N, error) {
		l, err := f.CallAny(x)
		if err != nil {
			return nil, v.newErr(err, n, "BinaryOperationExpr: eval left")
		}
		r, err := g.CallAny(x)
		if err != nil {
			return nil, v.newErr(err, n, "BinaryOperationExpr: eval right")
		}
		left, right := l[0], r[0]
		_, lv, ok := AsValueContainer(left).GetFirstValue()
		if !ok {
			return nil, v.invalidValue(n, "BinaryOperationExpr: no left value")
		}
		_, rv, ok := AsValueContainer(right).GetFirstValue()
		if !ok {
			return nil, v.invalidValue(n, "BinaryOperationExpr: no right value")
		}
		got, err := op(lv.AsOp(), rv.AsOp())
		if err != nil {
			return nil, v.newErr(err, n, "BinaryOperationExpr: exec")
		}
		res := AsValueContainer(node.New())
		res.SetContainerValue(got.AsData())
		return res.N, nil
	}), nil
}

func (v TreeVisitor) newBinaryOperation(n *BinaryOperationExpr) (string, func(*OP, *OP) (*OP, error), error) {
	switch n.Op {
	case opcode.LogicAnd: // AND
		return "LogicAnd", func(x, y *OP) (*OP, error) { return x.LogicalAnd(y) }, nil
	case opcode.LogicOr: // OR
		return "LogicOr", func(x, y *OP) (*OP, error) { return x.LogicalOr(y) }, nil
	case opcode.LogicXor: // XOR
		return "LogicXor", func(x, y *OP) (*OP, error) { return x.LogicalXor(y) }, nil
	case opcode.Plus: // +
		return "Plus", func(x, y *OP) (*OP, error) { return x.Add(y) }, nil
	case opcode.Minus: // -
		return "Minus", func(x, y *OP) (*OP, error) { return x.Subtract(y) }, nil
	case opcode.Mul: // *
		return "Mul", func(x, y *OP) (*OP, error) { return x.Multiply(y) }, nil
	case opcode.Div: // /
		return "Div", func(x, y *OP) (*OP, error) { return x.Divide(y) }, nil
	case opcode.Mod: // %
		return "Mod", func(x, y *OP) (*OP, error) { return x.Mod(y) }, nil
	case opcode.LeftShift: // <<
		return "LeftShift", func(x, y *OP) (*OP, error) { return x.LeftShift(y) }, nil
	case opcode.RightShift: // >>
		return "RightShift", func(x, y *OP) (*OP, error) { return x.RightShift(y) }, nil
	default:
		cmpFunc, err := v.compareBinaryOperation(n)
		if err != nil {
			return "", nil, fmt.Errorf("%w: op=%d", errors.Join(ErrNotImplmented, err), n.Op)
		}
		return "Compare", func(x, y *OP) (*OP, error) {
			return node.Bool(cmpFunc(x.Compare(y))).AsOp(), nil
		}, nil
	}
}

func (v TreeVisitor) compareBinaryOperation(n *BinaryOperationExpr) (func(node.CompareResult) bool, error) {
	switch n.Op {
	case opcode.GE: // >=
		return func(r node.CompareResult) bool { return r == util.CompareGreater || r == util.CompareEqual }, nil
	case opcode.LE: // <=
		return func(r node.CompareResult) bool { return r == util.CompareLess || r == util.CompareEqual }, nil
	case opcode.EQ: // =
		return func(r node.CompareResult) bool { return r == util.CompareEqual }, nil
	case opcode.NE: // <>, !=
		return func(r node.CompareResult) bool { return r != util.CompareEqual && r != util.CompareUnknown }, nil
	case opcode.LT: // <
		return func(r node.CompareResult) bool { return r == util.CompareLess }, nil
	case opcode.GT: // >
		return func(r node.CompareResult) bool { return r == util.CompareGreater }, nil
	default:
		return nil, v.invalidTree(n, "BinaryOperationExpr unknown compare code")
	}
}
