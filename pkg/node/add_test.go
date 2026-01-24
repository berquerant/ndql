package node_test

import (
	"testing"
	"time"

	"github.com/berquerant/ndql/pkg/node"
)

func TestAdd(t *testing.T) {
	bb := defaultFailedBinaryOpTestcaseBuilder()
	bb.
		nullPerm().
		pairPermExcept(bb.f(), bb.f(), bb.i()).
		pairPermExcept(bb.i(), bb.f(), bb.i()).
		pairPermExcept(bb.s(), bb.s()).
		pairPermExcept(bb.b(), bb.b()).
		pairPermExcept(bb.t(), bb.d()).
		pairPermExcept(bb.d(), bb.t(), bb.d())
	runBinaryOpTest(t, func(left, right node.Data) (node.Data, error) {
		x, err := left.AsOp().Add(right.AsOp())
		if err != nil {
			return nil, err
		}
		return x.AsData(), nil
	}, bb.build(), []*binaryOpTestacase{
		{
			left:  node.Float(1.2),
			right: node.Float(2),
			want:  node.Float(3.2),
		},
		{
			left:  node.Float(1.2),
			right: node.Int(2),
			want:  node.Float(3.2),
		},
		{
			left:  node.Int(1),
			right: node.Float(2.2),
			want:  node.Float(3.2),
		},
		{
			left:  node.Int(1),
			right: node.Int(2),
			want:  node.Int(3),
		},
		{
			left:  node.String("str"),
			right: node.String("ing"),
			want:  node.String("string"),
		},
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
		{
			left:  node.Time(time.Unix(1767348000, 0)),
			right: node.Duration(time.Second),
			want:  node.Time(time.Unix(1767348001, 0)),
		},
		{
			left:  node.Duration(time.Second),
			right: node.Time(time.Unix(1767348000, 0)),
			want:  node.Time(time.Unix(1767348001, 0)),
		},
		{
			left:  node.Duration(time.Second),
			right: node.Duration(time.Second * 2),
			want:  node.Duration(time.Second * 3),
		},
	})
}
