package node_test

import (
	"testing"

	"github.com/berquerant/ndql/pkg/node"
)

func TestDir(t *testing.T) {
	s := defaultFailedTestcaseSeed()
	runUnaryOpTest(t, func(v node.Data) (node.Data, error) {
		x, err := v.AsOp().Dir()
		if err != nil {
			return nil, err
		}
		return x.AsData(), nil
	}, newFailedUnaryOpTestcases(s.except(s.s())...), []*unaryOpTestcase{
		{
			v:    node.String("/"),
			want: node.String("/"),
		},
		{
			v:    node.String("/dir1"),
			want: node.String("/"),
		},
		{
			v:    node.String("/dir1/dir2"),
			want: node.String("/dir1"),
		},
	})
}

func TestBasename(t *testing.T) {
	s := defaultFailedTestcaseSeed()
	runUnaryOpTest(t, func(v node.Data) (node.Data, error) {
		x, err := v.AsOp().Basename()
		if err != nil {
			return nil, err
		}
		return x.AsData(), nil
	}, newFailedUnaryOpTestcases(s.except(s.s())...), []*unaryOpTestcase{
		{
			v:    node.String("/"),
			want: node.String("/"),
		},
		{
			v:    node.String("/dir1"),
			want: node.String("dir1"),
		},
		{
			v:    node.String("/dir1/dir2"),
			want: node.String("dir2"),
		},
		{
			v:    node.String("/dir1/dir2/some.txt"),
			want: node.String("some.txt"),
		},
	})
}

func TestExtension(t *testing.T) {
	s := defaultFailedTestcaseSeed()
	runUnaryOpTest(t, func(v node.Data) (node.Data, error) {
		x, err := v.AsOp().Extension()
		if err != nil {
			return nil, err
		}
		return x.AsData(), nil
	}, newFailedUnaryOpTestcases(s.except(s.s())...), []*unaryOpTestcase{
		{
			v:    node.String("/"),
			want: node.String(""),
		},
		{
			v:    node.String("/dir1/dir2/some.txt"),
			want: node.String(".txt"),
		},
		{
			v:    node.String("some.tar.gz"),
			want: node.String(".gz"),
		},
	})
}

func TestAbsPath(t *testing.T) {
	s := defaultFailedTestcaseSeed()
	runUnaryOpTest(t, func(v node.Data) (node.Data, error) {
		x, err := v.AsOp().AbsPath()
		if err != nil {
			return nil, err
		}
		return x.AsData(), nil
	}, newFailedUnaryOpTestcases(s.except(s.s())...))
}

func TestRelPath(t *testing.T) {
	bb := defaultFailedBinaryOpTestcaseBuilder()
	bb.
		perm(bb.except(bb.s())...)
	runBinaryOpTest(t, func(left, right node.Data) (node.Data, error) {
		x, err := left.AsOp().RelPath(right.AsOp())
		if err != nil {
			return nil, err
		}
		return x.AsData(), nil
	}, bb.build(), []*binaryOpTestacase{
		{
			left:  node.String("/root"),
			right: node.String("/root"),
			want:  node.String("."),
		},
		{
			left:  node.String("/root/sub"),
			right: node.String("/root"),
			want:  node.String("sub"),
		},
		{
			left:  node.String("/another"),
			right: node.String("/root"),
			want:  node.String("../another"),
		},
	})
}
