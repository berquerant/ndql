package node_test

import (
	"testing"

	"github.com/berquerant/ndql/pkg/node"
)

func TestMath(t *testing.T) {
	t.Run("unary", func(t *testing.T) {
		t.Run("sqrt", func(t *testing.T) {
			s := defaultFailedTestcaseSeed()
			runUnaryOpTest(t, func(v node.Data) (node.Data, error) {
				x, err := v.AsOp().Sqrt()
				if err != nil {
					return nil, err
				}
				return x.AsData(), nil
			}, newFailedUnaryOpTestcases(s.except(s.f(), s.i())...), []*unaryOpTestcase{
				{
					v:    node.Int(4),
					want: node.Float(2),
				},
				{
					v:    node.Float(16),
					want: node.Float(4),
				},
				{
					v:    node.Float(-1),
					want: node.NewNull(),
				},
			})
		})
	})

	t.Run("binary", func(t *testing.T) {
		t.Run("pow", func(t *testing.T) {
			bb := defaultFailedBinaryOpTestcaseBuilder()
			bb.
				nullPerm().
				perm(bb.except(bb.f(), bb.i())...).
				pairPermExcept(bb.f(), bb.f(), bb.i()).
				pairPermExcept(bb.i(), bb.f(), bb.i())
			runBinaryOpTest(t, func(x, y node.Data) (node.Data, error) {
				z, err := x.AsOp().Pow(y.AsOp())
				if err != nil {
					return nil, err
				}
				return z.AsData(), nil
			}, bb.build(), []*binaryOpTestacase{
				{
					left:  node.Int(2),
					right: node.Int(3),
					want:  node.Float(8),
				},
				{
					left:  node.Int(4),
					right: node.Float(0.5),
					want:  node.Float(2),
				},
				{
					left:  node.Float(1.2),
					right: node.Int(2),
					want:  node.Float(1.44),
				},
				{
					left:  node.Int(0),
					right: node.Int(-1),
					want:  node.NewNull(),
				},
			})
		})
	})
}
