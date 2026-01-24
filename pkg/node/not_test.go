package node_test

import (
	"testing"
	"time"

	"github.com/berquerant/ndql/pkg/node"
)

func TestNot(t *testing.T) {
	s := defaultFailedTestcaseSeed()
	runUnaryOpTest(t, func(v node.Data) (node.Data, error) {
		x, err := v.AsOp().Not()
		if err != nil {
			return nil, err
		}
		return x.AsData(), nil
	}, newFailedUnaryOpTestcases(s.except(s.f(), s.i(), s.b(), s.d())...), []*unaryOpTestcase{
		{
			v:    node.Float(1),
			want: node.Float(-1),
		},
		{
			v:    node.Int(1),
			want: node.Int(-1),
		},
		{
			v:    node.Bool(false),
			want: node.Bool(true),
		},
		{
			v:    node.Bool(true),
			want: node.Bool(false),
		},
		{
			v:    node.Duration(time.Second),
			want: node.Duration(-time.Second),
		},
	})
}
