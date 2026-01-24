package node_test

import (
	"testing"
	"time"

	"github.com/berquerant/ndql/pkg/node"
	"github.com/berquerant/ndql/pkg/util"
)

func TestMultiply(t *testing.T) {
	bb := defaultFailedBinaryOpTestcaseBuilder()
	bb.
		nullPerm().
		pairPermExcept(bb.f(), bb.f(), bb.i(), bb.d()).
		pairPermExcept(bb.i(), bb.f(), bb.i(), bb.s(), bb.d()).
		pairPermExcept(bb.b(), bb.b()).
		pairPermExcept(bb.s(), bb.i()).
		pairPermExcept(bb.t()).
		pairPermExcept(bb.d(), bb.f(), bb.i())
	runBinaryOpTest(t, func(left, right node.Data) (node.Data, error) {
		x, err := left.AsOp().Multiply(right.AsOp())
		if err != nil {
			return nil, err
		}
		return x.AsData(), nil
	}, bb.build(), []*binaryOpTestacase{
		{
			left:  node.Float(1.2),
			right: node.Float(2),
			want:  node.Float(2.4),
		},
		{
			left:  node.Float(1.2),
			right: node.Int(2),
			want:  node.Float(2.4),
		},
		{
			left:  node.Float(1.2),
			right: node.Duration(time.Second),
			want:  node.Duration(util.NewDuration(1.2, time.Second)),
		},
		{
			left:  node.Int(2),
			right: node.Float(1.2),
			want:  node.Float(2.4),
		},
		{
			left:  node.Int(2),
			right: node.Int(3),
			want:  node.Int(6),
		},
		{
			left:  node.Int(2),
			right: node.String("s"),
			want:  node.String("ss"),
		},
		{
			left:  node.Int(2),
			right: node.Duration(time.Second),
			want:  node.Duration(2 * time.Second),
		},
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
		{
			left:  node.String("s"),
			right: node.Int(2),
			want:  node.String("ss"),
		},
		{
			left:  node.Duration(time.Second),
			right: node.Float(2.4),
			want:  node.Duration(util.NewDuration(2.4, time.Second)),
		},
		{
			left:  node.Duration(time.Second),
			right: node.Int(2),
			want:  node.Duration(2 * time.Second),
		},
	})
}
