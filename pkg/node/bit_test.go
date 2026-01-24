package node_test

import (
	"testing"

	"github.com/berquerant/ndql/pkg/node"
)

func TestBitNot(t *testing.T) {
	s := defaultFailedTestcaseSeed()
	runUnaryOpTest(t, func(v node.Data) (node.Data, error) {
		x, err := v.AsOp().BitNot()
		if err != nil {
			return nil, err
		}
		return x.AsData(), nil
	}, newFailedUnaryOpTestcases(s.except(s.i())...), []*unaryOpTestcase{
		{
			v:    node.Int(127),
			want: node.Int(-128),
		},
	})
}

func TestBitAnd(t *testing.T) {
	bb := defaultFailedBinaryOpTestcaseBuilder()
	bb.
		nullPerm().
		perm(bb.except(bb.i())...).
		pairPermExcept(bb.i(), bb.i())
	runBinaryOpTest(t, func(left, right node.Data) (node.Data, error) {
		x, err := left.AsOp().BitAnd(right.AsOp())
		if err != nil {
			return nil, err
		}
		return x.AsData(), nil
	}, bb.build(), []*binaryOpTestacase{
		{
			left:  node.Int(13),
			right: node.Int(127),
			want:  node.Int(13),
		},
	})
}

func TestBinOr(t *testing.T) {
	bb := defaultFailedBinaryOpTestcaseBuilder()
	bb.
		nullPerm().
		perm(bb.except(bb.i())...).
		pairPermExcept(bb.i(), bb.i())
	runBinaryOpTest(t, func(left, right node.Data) (node.Data, error) {
		x, err := left.AsOp().BitOr(right.AsOp())
		if err != nil {
			return nil, err
		}
		return x.AsData(), nil
	}, bb.build(), []*binaryOpTestacase{
		{
			left:  node.Int(17),
			right: node.Int(30),
			want:  node.Int(31),
		},
	})
}

func TestBitXor(t *testing.T) {
	bb := defaultFailedBinaryOpTestcaseBuilder()
	bb.
		nullPerm().
		perm(bb.except(bb.i())...).
		pairPermExcept(bb.i(), bb.i())
	runBinaryOpTest(t, func(left, right node.Data) (node.Data, error) {
		x, err := left.AsOp().BitXor(right.AsOp())
		if err != nil {
			return nil, err
		}
		return x.AsData(), nil
	}, bb.build(), []*binaryOpTestacase{
		{
			left:  node.Int(127),
			right: node.Int(10),
			want:  node.Int(117),
		},
	})

}

func TestLeftShift(t *testing.T) {
	bb := defaultFailedBinaryOpTestcaseBuilder()
	bb.
		nullPerm().
		perm(bb.except(bb.i())...).
		pairPermExcept(bb.i(), bb.i())
	runBinaryOpTest(t, func(left, right node.Data) (node.Data, error) {
		x, err := left.AsOp().LeftShift(right.AsOp())
		if err != nil {
			return nil, err
		}
		return x.AsData(), nil
	}, bb.build(), []*binaryOpTestacase{
		{
			left:  node.Int(1),
			right: node.Int(2),
			want:  node.Int(4),
		},
		{
			left:  node.Int(4),
			right: node.Int(-2),
			want:  node.Int(1),
		},
	})
}

func TestRightShift(t *testing.T) {
	bb := defaultFailedBinaryOpTestcaseBuilder()
	bb.
		nullPerm().
		perm(bb.except(bb.i())...).
		pairPermExcept(bb.i(), bb.i())
	runBinaryOpTest(t, func(left, right node.Data) (node.Data, error) {
		x, err := left.AsOp().RightShift(right.AsOp())
		if err != nil {
			return nil, err
		}
		return x.AsData(), nil
	}, bb.build(), []*binaryOpTestacase{
		{
			left:  node.Int(1),
			right: node.Int(-2),
			want:  node.Int(4),
		},
		{
			left:  node.Int(4),
			right: node.Int(2),
			want:  node.Int(1),
		},
		{
			left:  node.Int(4),
			right: node.Int(10),
			want:  node.Int(0),
		},
	})
}
