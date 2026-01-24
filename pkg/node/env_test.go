package node_test

import (
	"os"
	"testing"

	"github.com/berquerant/ndql/pkg/node"
)

func TestEnvOr(t *testing.T) {
	os.Setenv("Test1", "testval")
	defer os.Unsetenv("Test1")
	bb := defaultFailedBinaryOpTestcaseBuilder()
	bb.
		perm(bb.except(bb.s())...)
	runBinaryOpTest(t, func(left, right node.Data) (node.Data, error) {
		x, err := left.AsOp().EnvOr(right.AsOp())
		if err != nil {
			return nil, err
		}
		return x.AsData(), nil
	}, bb.build(), []*binaryOpTestacase{
		{
			left:  node.String("Test1"),
			right: node.Int(1),
			want:  node.String("testval"),
		},
		{
			left:  node.String("unmatched"),
			right: node.Int(1),
			want:  node.Int(1),
		},
	})
}

func TestEnv(t *testing.T) {
	os.Setenv("Test1", "testval")
	defer os.Unsetenv("Test1")
	s := defaultFailedTestcaseSeed()
	runUnaryOpTest(t, func(v node.Data) (node.Data, error) {
		x, err := v.AsOp().Env()
		if err != nil {
			return nil, err
		}
		return x.AsData(), nil
	}, newFailedUnaryOpTestcases(s.except(s.s())...), []*unaryOpTestcase{
		{
			v:    node.String("Test1"),
			want: node.String("testval"),
		},
		{
			v:    node.String("missing"),
			want: node.NewNull(),
		},
	})
}
