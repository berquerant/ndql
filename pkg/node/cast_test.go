package node_test

import (
	"testing"
	"time"

	"github.com/berquerant/ndql/pkg/node"
	"github.com/berquerant/ndql/pkg/util"
)

func TestAsNull(t *testing.T) {
	s := defaultFailedTestcaseSeed()
	runUnaryOpTest(t, func(v node.Data) (node.Data, error) {
		return v.AsOp().AsNull()
	}, newFailedUnaryOpTestcases(s.except(s.n())...), []*unaryOpTestcase{
		{
			v:    node.NewNull(),
			want: node.NewNull(),
		},
	})
}

func TestAsFloat(t *testing.T) {
	s := defaultFailedTestcaseSeed()
	runUnaryOpTest(t, func(v node.Data) (node.Data, error) {
		return v.AsOp().AsFloat()
	}, newFailedUnaryOpTestcases(s.n()), []*unaryOpTestcase{
		{
			v:    node.Float(1),
			want: node.Float(1),
		},
		{
			v:    node.Int(2),
			want: node.Float(2),
		},
		{
			v:    node.Bool(true),
			want: node.Float(1),
		},
		{
			v:    node.String("1.2"),
			want: node.Float(1.2),
		},
		{
			v: node.String("s"),
		},
		{
			v:    node.Time(time.Unix(1767348000, 0)),
			want: node.Float(1767348000),
		},
		{
			v:    node.Duration(time.Second),
			want: node.Float(1000000000),
		},
	})
}

func TestAsInt(t *testing.T) {
	s := defaultFailedTestcaseSeed()
	runUnaryOpTest(t, func(v node.Data) (node.Data, error) {
		return v.AsOp().AsInt()
	}, newFailedUnaryOpTestcases(s.n()), []*unaryOpTestcase{
		{
			v:    node.Float(1),
			want: node.Int(1),
		},
		{
			v:    node.Int(2),
			want: node.Int(2),
		},
		{
			v:    node.Bool(false),
			want: node.Int(0),
		},
		{
			v:    node.Bool(true),
			want: node.Int(1),
		},
		{
			v:    node.String("1"),
			want: node.Int(1),
		},
		{
			v: node.String("s"),
		},
		{
			v:    node.Time(time.Unix(1767348000, 0)),
			want: node.Int(1767348000),
		},
		{
			v:    node.Duration(time.Second),
			want: node.Int(1000000000),
		},
	})
}

func TestAsBool(t *testing.T) {
	s := defaultFailedTestcaseSeed()
	runUnaryOpTest(t, func(v node.Data) (node.Data, error) {
		return v.AsOp().AsBool()
	}, newFailedUnaryOpTestcases(s.except(s.f(), s.i(), s.b(), s.s())...), []*unaryOpTestcase{
		{
			v:    node.Float(0),
			want: node.Bool(false),
		},
		{
			v:    node.Float(1),
			want: node.Bool(true),
		},
		{
			v:    node.Int(0),
			want: node.Bool(false),
		},
		{
			v:    node.Int(1),
			want: node.Bool(true),
		},
		{
			v:    node.Bool(false),
			want: node.Bool(false),
		},
		{
			v:    node.Bool(true),
			want: node.Bool(true),
		},
		{
			v:    node.String(""),
			want: node.Bool(false),
		},
		{
			v:    node.String("s"),
			want: node.Bool(true),
		},
	})
}

func TestAsString(t *testing.T) {
	runUnaryOpTest(t, func(v node.Data) (node.Data, error) {
		return v.AsOp().AsString()
	}, []*unaryOpTestcase{
		{
			v:    node.NewNull(),
			want: node.String("null"),
		},
		{
			v:    node.Float(1.2),
			want: node.String("1.2"),
		},
		{
			v:    node.Int(2),
			want: node.String("2"),
		},
		{
			v:    node.Bool(false),
			want: node.String("false"),
		},
		{
			v:    node.Bool(true),
			want: node.String("true"),
		},
		{
			v:    node.String("str"),
			want: node.String("str"),
		},
		{
			v:    node.Time(util.Must(time.Parse(time.DateTime, "2026-01-02 10:00:00"))),
			want: node.String("2026-01-02 10:00:00"),
		},
		{
			v:    node.Duration(time.Second),
			want: node.String("1s"),
		},
	})
}

func TestAsTime(t *testing.T) {
	s := defaultFailedTestcaseSeed()
	runUnaryOpTest(t, func(v node.Data) (node.Data, error) {
		return v.AsOp().AsTime()
	}, newFailedUnaryOpTestcases(s.except(s.f(), s.i(), s.s(), s.t())...), []*unaryOpTestcase{
		{
			v:    node.Float(1767348000),
			want: node.Time(time.Unix(1767348000, 0).In(time.UTC)),
		},
		{
			v:    node.Int(1767348000),
			want: node.Time(time.Unix(1767348000, 0).In(time.UTC)),
		},
		{
			v:    node.String("2026-01-02 10:00:00"),
			want: node.Time(time.Unix(1767348000, 0).In(time.UTC)),
		},
		{
			v:    node.Time(util.Must(time.Parse(time.DateTime, "2026-01-02 10:00:00"))),
			want: node.Time(util.Must(time.Parse(time.DateTime, "2026-01-02 10:00:00"))),
		},
	})
}

func TestAsDuration(t *testing.T) {
	s := defaultFailedTestcaseSeed()
	runUnaryOpTest(t, func(v node.Data) (node.Data, error) {
		return v.AsOp().AsDuration()
	}, newFailedUnaryOpTestcases(s.except(s.f(), s.i(), s.s(), s.d())...), []*unaryOpTestcase{
		{
			v:    node.Float(1000000000),
			want: node.Duration(time.Second),
		},
		{
			v:    node.Int(1000000000),
			want: node.Duration(time.Second),
		},
		{
			v:    node.String("1s"),
			want: node.Duration(time.Second),
		},
		{
			v:    node.Duration(time.Second),
			want: node.Duration(time.Second),
		},
	})
}
