package node_test

import (
	"testing"

	"github.com/berquerant/ndql/pkg/node"
)

func TestInverse(t *testing.T) {
	s := defaultFailedTestcaseSeed()
	runUnaryOpTest(t, func(v node.Data) (node.Data, error) {
		x, err := v.AsOp().Inverse()
		if err != nil {
			return nil, err
		}
		return x.AsData(), nil
	}, newFailedUnaryOpTestcases(s.except(s.f(), s.i(), s.s())...), []*unaryOpTestcase{
		{
			v:    node.Float(2),
			want: node.Float(0.5),
		},
		{
			v:    node.Int(2),
			want: node.Float(0.5),
		},
		{
			v: node.Float(0),
		},
		{
			v: node.Int(0),
		},
		{
			v:    node.String("str"),
			want: node.String("rts"),
		},
		{
			v:    node.String(""),
			want: node.String(""),
		},
	})
}
