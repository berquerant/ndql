package node_test

import (
	"testing"

	"github.com/berquerant/ndql/pkg/node"
)

func TestLen(t *testing.T) {
	s := defaultFailedTestcaseSeed()
	runUnaryOpTest(t, func(v node.Data) (node.Data, error) {
		x, err := v.AsOp().Len()
		if err != nil {
			return nil, err
		}
		return x.AsData(), nil
	}, newFailedUnaryOpTestcases(s.except(s.s())...), []*unaryOpTestcase{
		{
			v:    node.String(""),
			want: node.Int(0),
		},
		{
			v:    node.String("abc"),
			want: node.Int(3),
		},
		{
			v:    node.String("あいう"),
			want: node.Int(3),
		},
	})
}

func TestSize(t *testing.T) {
	s := defaultFailedTestcaseSeed()
	runUnaryOpTest(t, func(v node.Data) (node.Data, error) {
		x, err := v.AsOp().Size()
		if err != nil {
			return nil, err
		}
		return x.AsData(), nil
	}, newFailedUnaryOpTestcases(s.except(s.s())...), []*unaryOpTestcase{
		{
			v:    node.String(""),
			want: node.Int(0),
		},
		{
			v:    node.String("abc"),
			want: node.Int(3),
		},
		{
			v:    node.String("あいう"),
			want: node.Int(9),
		},
	})
}

func TestLike(t *testing.T) {
	bb := defaultFailedBinaryOpTestcaseBuilder()
	bb.
		perm(bb.except(bb.s())...)
	runBinaryOpTest(t, func(left, right node.Data) (node.Data, error) {
		x, err := left.AsOp().Like(right.AsOp())
		if err != nil {
			return nil, err
		}
		return x.AsData(), nil
	}, bb.build(), []*binaryOpTestacase{
		{
			left:  node.String("a"),
			right: node.String("a"),
			want:  node.Bool(true),
		},
		{
			left:  node.String("b"),
			right: node.String("a"),
			want:  node.Bool(false),
		},
		{
			left:  node.String("ab"),
			right: node.String("a_"),
			want:  node.Bool(true),
		},
		{
			left:  node.String("abc"),
			right: node.String("a%"),
			want:  node.Bool(true),
		},
	})
}

func TestRegexp(t *testing.T) {
	bb := defaultFailedBinaryOpTestcaseBuilder()
	bb.
		perm(bb.except(bb.s())...)
	runBinaryOpTest(t, func(left, right node.Data) (node.Data, error) {
		x, err := left.AsOp().Regexp(right.AsOp())
		if err != nil {
			return nil, err
		}
		return x.AsData(), nil
	}, bb.build(), []*binaryOpTestacase{
		{
			left:  node.String("a"),
			right: node.String("a"),
			want:  node.Bool(true),
		},
		{
			left:  node.String("b"),
			right: node.String("a"),
			want:  node.Bool(false),
		},
		{
			left:  node.String("ab"),
			right: node.String("a."),
			want:  node.Bool(true),
		},
		{
			left:  node.String("abc"),
			right: node.String("a.*"),
			want:  node.Bool(true),
		},
	})
}

func TestRegexpInstr(t *testing.T) {
	bb := defaultFailedBinaryOpTestcaseBuilder()
	bb.
		perm(bb.except(bb.s())...)
	runBinaryOpTest(t, func(left, right node.Data) (node.Data, error) {
		x, err := left.AsOp().RegexpInstr(right.AsOp())
		if err != nil {
			return nil, err
		}
		return x.AsData(), nil
	}, bb.build(), []*binaryOpTestacase{
		{
			left:  node.String("abc"),
			right: node.String("a"),
			want:  node.Int(1),
		},
		{
			left:  node.String("abc"),
			right: node.String("b"),
			want:  node.Int(2),
		},
		{
			left:  node.String("abc"),
			right: node.String("b."),
			want:  node.Int(2),
		},
		{
			left:  node.String("abc"),
			right: node.String("x.*"),
			want:  node.Int(0),
		},
	})
}

func TestRegexpSubstr(t *testing.T) {
	bb := defaultFailedBinaryOpTestcaseBuilder()
	bb.
		perm(bb.except(bb.s())...)
	runBinaryOpTest(t, func(left, right node.Data) (node.Data, error) {
		x, err := left.AsOp().RegexpSubstr(right.AsOp())
		if err != nil {
			return nil, err
		}
		return x.AsData(), nil
	}, bb.build(), []*binaryOpTestacase{
		{
			left:  node.String("abc"),
			right: node.String("a"),
			want:  node.String("a"),
		},
		{
			left:  node.String("abc"),
			right: node.String("b"),
			want:  node.String("b"),
		},
		{
			left:  node.String("abc"),
			right: node.String("b."),
			want:  node.String("bc"),
		},
		{
			left:  node.String("abc"),
			right: node.String("x.*"),
			want:  node.String(""),
		},
	})
}

func TestRegexpReplace(t *testing.T) {
	bb := defaultFailedBinaryOpTestcaseBuilder()
	bb.
		perm(bb.except(bb.s())...)
	runBinaryOpTest(t, func(left, right node.Data) (node.Data, error) {
		repl := node.String("REPL").AsOp()
		x, err := left.AsOp().RegexpReplace(right.AsOp(), repl)
		if err != nil {
			return nil, err
		}
		return x.AsData(), nil
	}, bb.build(), []*binaryOpTestacase{
		{
			left:  node.String("abc"),
			right: node.String("a"),
			want:  node.String("REPLbc"),
		},
		{
			left:  node.String("abc"),
			right: node.String("b"),
			want:  node.String("aREPLc"),
		},
		{
			left:  node.String("abc"),
			right: node.String("b."),
			want:  node.String("aREPL"),
		},
		{
			left:  node.String("abc"),
			right: node.String("x.*"),
			want:  node.String("abc"),
		},
	})
}

func TestFormat(t *testing.T) {
	runVariadicOpTest(t, func(v ...node.Data) (node.Data, error) {
		a := DataListToOpList(v[1:]...)
		x, err := v[0].AsOp().Format(a...)
		if err != nil {
			return nil, err
		}
		return x.AsData(), nil
	}, []*variadicOpTestacase{
		{
			v:    []node.Data{node.String("str")},
			want: node.String("str"),
		},
		{
			v:    []node.Data{node.String("str=%s"), node.String("v")},
			want: node.String("str=v"),
		},
		{
			v:    []node.Data{node.String("str=%s=%d"), node.String("v"), node.Int(1)},
			want: node.String("str=v=1"),
		},
	})
}

func TestLower(t *testing.T) {
	s := defaultFailedTestcaseSeed()
	runUnaryOpTest(t, func(v node.Data) (node.Data, error) {
		x, err := v.AsOp().Lower()
		if err != nil {
			return nil, err
		}
		return x.AsData(), nil
	}, newFailedUnaryOpTestcases(s.except(s.s())...), []*unaryOpTestcase{
		{
			v:    node.String(""),
			want: node.String(""),
		},
		{
			v:    node.String("str"),
			want: node.String("str"),
		},
		{
			v:    node.String("STR"),
			want: node.String("str"),
		},
	})
}

func TestUpper(t *testing.T) {
	s := defaultFailedTestcaseSeed()
	runUnaryOpTest(t, func(v node.Data) (node.Data, error) {
		x, err := v.AsOp().Upper()
		if err != nil {
			return nil, err
		}
		return x.AsData(), nil
	}, newFailedUnaryOpTestcases(s.except(s.s())...), []*unaryOpTestcase{
		{
			v:    node.String(""),
			want: node.String(""),
		},
		{
			v:    node.String("str"),
			want: node.String("STR"),
		},
		{
			v:    node.String("STR"),
			want: node.String("STR"),
		},
	})
}

func TestSha2(t *testing.T) {
	s := defaultFailedTestcaseSeed()
	runUnaryOpTest(t, func(v node.Data) (node.Data, error) {
		x, err := v.AsOp().Sha2()
		if err != nil {
			return nil, err
		}
		return x.AsData(), nil
	}, newFailedUnaryOpTestcases(s.except(s.s())...), []*unaryOpTestcase{
		{
			v:    node.String(""),
			want: node.String("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"),
		},
		{
			v:    node.String("str"),
			want: node.String("8c25cb3686462e9a86d2883c5688a22fe738b0bbc85f458d2d2b5f3f667c6d5a"),
		},
	})
}

func TestConcatWs(t *testing.T) {
	runVariadicOpTest(t, func(v ...node.Data) (node.Data, error) {
		xs := DataListToOpList(v...)
		r, err := xs[0].ConcatWs(xs[1:]...)
		if err != nil {
			return nil, err
		}
		return r.AsData(), nil
	}, []*variadicOpTestacase{
		{
			v: []node.Data{node.Int(1)},
		},
		{
			v:    []node.Data{node.String(",")},
			want: node.String(""),
		},
		{
			v: []node.Data{
				node.String(","),
				node.String("a"),
			},
			want: node.String("a"),
		},
		{
			v: []node.Data{
				node.String(","),
				node.String("a"),
				node.String("b"),
			},
			want: node.String("a,b"),
		},
		{
			v: []node.Data{
				node.String(","),
				node.String("a"),
				node.String("b"),
				node.String("c"),
			},
			want: node.String("a,b,c"),
		},
		{
			v: []node.Data{
				node.String(","),
				node.String("a"),
				node.String("b"),
				node.String("c"),
				node.NewNull(),
			},
		},
	})
}

func TestInstr(t *testing.T) {
	bb := defaultFailedBinaryOpTestcaseBuilder()
	bb.
		perm(bb.except(bb.s())...)
	runBinaryOpTest(t, func(left, right node.Data) (node.Data, error) {
		x, err := left.AsOp().Instr(right.AsOp())
		if err != nil {
			return nil, err
		}
		return x.AsData(), nil
	}, bb.build(), []*binaryOpTestacase{
		{
			left:  node.String("a"),
			right: node.String("a"),
			want:  node.Int(1),
		},
		{
			left:  node.String("b"),
			right: node.String("a"),
			want:  node.Int(0),
		},
		{
			left:  node.String("abc"),
			right: node.String("b"),
			want:  node.Int(2),
		},
	})
}

func TestSubstr(t *testing.T) {
	runVariadicOpTest(t, func(v ...node.Data) (node.Data, error) {
		xs := DataListToOpList(v...)
		r, err := xs[0].Substr(xs[1:]...)
		if err != nil {
			return nil, err
		}
		return r.AsData(), nil
	}, []*variadicOpTestacase{
		{
			v: []node.Data{
				node.String(""),
				node.Int(1),
			},
			want: node.String(""),
		},
		{
			v: []node.Data{
				node.String("abcd"),
				node.Int(1),
			},
			want: node.String("abcd"),
		},
		{
			v: []node.Data{
				node.String("abcd"),
				node.Int(2),
			},
			want: node.String("bcd"),
		},
		{
			v: []node.Data{
				node.String("abcd"),
				node.Int(2),
				node.Int(2),
			},
			want: node.String("bc"),
		},
		{
			v: []node.Data{
				node.String("abcd"),
				node.Int(-1),
			},
			want: node.String("d"),
		},
		{
			v: []node.Data{
				node.String("abcd"),
				node.Int(-2),
			},
			want: node.String("cd"),
		},
		{
			v: []node.Data{
				node.String("abcd"),
				node.Int(-3),
			},
			want: node.String("bcd"),
		},
		{
			v: []node.Data{
				node.String("abcd"),
				node.Int(-4),
			},
			want: node.String("abcd"),
		},
	})
}

func TestSubstrIndex(t *testing.T) {
	runVariadicOpTest(t, func(v ...node.Data) (node.Data, error) {
		xs := DataListToOpList(v...)
		r, err := xs[0].SubstrIndex(xs[1], xs[2])
		if err != nil {
			return nil, err
		}
		return r.AsData(), nil
	}, []*variadicOpTestacase{
		{
			v: []node.Data{
				node.String("a.b.c"),
				node.String("."),
				node.Int(1),
			},
			want: node.String("a"),
		},
		{
			v: []node.Data{
				node.String("a.b.c"),
				node.String("."),
				node.Int(2),
			},
			want: node.String("a.b"),
		},
		{
			v: []node.Data{
				node.String("a.b.c"),
				node.String("."),
				node.Int(3),
			},
			want: node.String("a.b.c"),
		},
		{
			v: []node.Data{
				node.String("a.b.c"),
				node.String("."),
				node.Int(4),
			},
			want: node.String("a.b.c"),
		},
		{
			v: []node.Data{
				node.String("a.b.c"),
				node.String("."),
				node.Int(0),
			},
			want: node.String(""),
		},
		{
			v: []node.Data{
				node.String("a.b.c"),
				node.String("."),
				node.Int(-1),
			},
			want: node.String("c"),
		},
		{
			v: []node.Data{
				node.String("a.b.c"),
				node.String("."),
				node.Int(-2),
			},
			want: node.String("b.c"),
		},
	})
}

func TestReplace(t *testing.T) {
	runVariadicOpTest(t, func(v ...node.Data) (node.Data, error) {
		xs := DataListToOpList(v...)
		r, err := xs[0].Replace(xs[1], xs[2])
		if err != nil {
			return nil, err
		}
		return r.AsData(), nil
	}, []*variadicOpTestacase{
		{
			v: []node.Data{
				node.String("a.b.c"),
				node.String("."),
				node.String("|"),
			},
			want: node.String("a|b|c"),
		},
		{
			v: []node.Data{
				node.String(""),
				node.String("."),
				node.String("|"),
			},
			want: node.String(""),
		},
	})
}

func TestTrim(t *testing.T) {
	runVariadicOpTest(t, func(v ...node.Data) (node.Data, error) {
		xs := DataListToOpList(v...)
		r, err := xs[0].Trim(xs[1:]...)
		if err != nil {
			return nil, err
		}
		return r.AsData(), nil
	}, []*variadicOpTestacase{
		{
			v: []node.Data{
				node.String("  space "),
			},
			want: node.String("space"),
		},
		{
			v: []node.Data{
				node.String("xxxXXXyyxx"),
				node.String("xy"),
			},
			want: node.String("XXX"),
		},
	})
}
