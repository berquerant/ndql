package node_test

import (
	"testing"

	"github.com/berquerant/ndql/pkg/node"
)

func TestMod(t *testing.T) {
	bb := defaultFailedBinaryOpTestcaseBuilder()
	bb.
		nullPerm().
		perm(bb.except(bb.s(), bb.b(), bb.t(), bb.d())...).
		pairPermExcept(bb.f(), bb.f(), bb.i()).
		pairPermExcept(bb.i(), bb.f(), bb.i())
	runBinaryOpTest(t, func(left, right node.Data) (node.Data, error) {
		x, err := left.AsOp().Mod(right.AsOp())
		if err != nil {
			return nil, err
		}
		return x.AsData(), nil
	}, bb.build(), []*binaryOpTestacase{
		{
			left:  node.Float(3.4),
			right: node.Float(1.2),
			want:  node.Float(1),
		},
		{
			left:  node.Float(3.4),
			right: node.Float(0),
		},
		{
			left:  node.Float(5),
			right: node.Int(2),
			want:  node.Float(1),
		},
		{
			left:  node.Int(11),
			right: node.Float(3),
			want:  node.Float(2),
		},
		{
			left:  node.Int(11),
			right: node.Int(3),
			want:  node.Int(2),
		},
		{
			left:  node.Int(1),
			right: node.Int(0),
		},
	})
}
