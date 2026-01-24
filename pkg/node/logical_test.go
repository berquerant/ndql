package node_test

import (
	"testing"

	"github.com/berquerant/ndql/pkg/node"
)

func TestLogicalNot(t *testing.T) {
	s := defaultFailedTestcaseSeed()
	runUnaryOpTest(t, func(v node.Data) (node.Data, error) {
		x, err := v.AsOp().LogicalNot()
		if err != nil {
			return nil, err
		}
		return x.AsData(), nil
	}, newFailedUnaryOpTestcases(s.except(s.b())...), []*unaryOpTestcase{
		{
			v:    node.Bool(false),
			want: node.Bool(true),
		},
		{
			v:    node.Bool(true),
			want: node.Bool(false),
		},
	})
}

func TestLogicalAnd(t *testing.T) {
	bb := defaultFailedBinaryOpTestcaseBuilder()
	bb.
		nullPerm().
		perm(bb.except(bb.b())...).
		pairPermExcept(bb.b(), bb.b())
	runBinaryOpTest(t, func(left, right node.Data) (node.Data, error) {
		x, err := left.AsOp().LogicalAnd(right.AsOp())
		if err != nil {
			return nil, err
		}
		return x.AsData(), nil
	}, bb.build(), []*binaryOpTestacase{
		{
			left:  node.Bool(false),
			right: node.Bool(false),
			want:  node.Bool(false),
		},
		{
			left:  node.Bool(false),
			right: node.Bool(true),
			want:  node.Bool(false),
		},
		{
			left:  node.Bool(true),
			right: node.Bool(false),
			want:  node.Bool(false),
		},
		{
			left:  node.Bool(true),
			right: node.Bool(true),
			want:  node.Bool(true),
		},
	})
}

func TestLogicalOr(t *testing.T) {
	bb := defaultFailedBinaryOpTestcaseBuilder()
	bb.
		nullPerm().
		perm(bb.except(bb.b())...).
		pairPermExcept(bb.b(), bb.b())
	runBinaryOpTest(t, func(left, right node.Data) (node.Data, error) {
		x, err := left.AsOp().LogicalOr(right.AsOp())
		if err != nil {
			return nil, err
		}
		return x.AsData(), nil
	}, bb.build(), []*binaryOpTestacase{
		{
			left:  node.Bool(false),
			right: node.Bool(false),
			want:  node.Bool(false),
		},
		{
			left:  node.Bool(false),
			right: node.Bool(true),
			want:  node.Bool(true),
		},
		{
			left:  node.Bool(true),
			right: node.Bool(false),
			want:  node.Bool(true),
		},
		{
			left:  node.Bool(true),
			right: node.Bool(true),
			want:  node.Bool(true),
		},
	})
}

func TestLogicalXor(t *testing.T) {
	bb := defaultFailedBinaryOpTestcaseBuilder()
	bb.
		nullPerm().
		perm(bb.except(bb.b())...).
		pairPermExcept(bb.b(), bb.b())
	runBinaryOpTest(t, func(left, right node.Data) (node.Data, error) {
		x, err := left.AsOp().LogicalXor(right.AsOp())
		if err != nil {
			return nil, err
		}
		return x.AsData(), nil
	}, bb.build(), []*binaryOpTestacase{
		{
			left:  node.Bool(false),
			right: node.Bool(false),
			want:  node.Bool(false),
		},
		{
			left:  node.Bool(false),
			right: node.Bool(true),
			want:  node.Bool(true),
		},
		{
			left:  node.Bool(true),
			right: node.Bool(false),
			want:  node.Bool(true),
		},
		{
			left:  node.Bool(true),
			right: node.Bool(true),
			want:  node.Bool(false),
		},
	})
}
