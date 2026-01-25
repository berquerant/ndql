package node_test

import (
	"testing"
	"time"

	"github.com/berquerant/ndql/pkg/node"
	"github.com/berquerant/ndql/pkg/util"
)

func TestStrToTime(t *testing.T) {
	bb := defaultFailedBinaryOpTestcaseBuilder()
	bb.
		nullPerm().
		perm(bb.except(bb.s())...).
		pairPermExcept(bb.s(), bb.s())
	runBinaryOpTest(t, func(left, right node.Data) (node.Data, error) {
		x, err := left.AsOp().StrToTime(right.AsOp())
		if err != nil {
			return nil, err
		}
		return x.AsData(), nil
	}, bb.build(), []*binaryOpTestacase{
		{
			left:  node.String("2026-01-02 10:00:00"),
			right: node.String(time.DateTime),
			want:  node.Time(util.Must(time.Parse(time.DateTime, "2026-01-02 10:00:00"))),
		},
		{
			left:  node.String("2026-01-02T10:00:00Z"),
			right: node.String(time.RFC3339),
			want:  node.Time(util.Must(time.Parse(time.DateTime, "2026-01-02 10:00:00"))),
		},
	})
}

func TestTimeFormat(t *testing.T) {
	bb := defaultFailedBinaryOpTestcaseBuilder()
	bb.
		nullPerm().
		perm(bb.except(bb.t(), bb.s())...).
		pairPermExcept(bb.t(), bb.s())
	runBinaryOpTest(t, func(left, right node.Data) (node.Data, error) {
		x, err := left.AsOp().TimeFormat(right.AsOp())
		if err != nil {
			return nil, err
		}
		return x.AsData(), nil
	}, bb.build(), []*binaryOpTestacase{
		{
			left:  node.Time(util.Must(time.Parse(time.DateTime, "2026-01-02 10:00:00"))),
			right: node.String(time.DateTime),
			want:  node.String("2026-01-02 10:00:00"),
		},
		{
			left:  node.Time(util.Must(time.Parse(time.DateTime, "2026-01-02 10:00:00"))),
			right: node.String("2006/01/02"),
			want:  node.String("2026/01/02"),
		},
	})
}

func TestYear(t *testing.T) {
	s := defaultFailedTestcaseSeed()
	runUnaryOpTest(t, func(v node.Data) (node.Data, error) {
		x, err := v.AsOp().Year()
		if err != nil {
			return nil, err
		}
		return x.AsData(), nil
	}, newFailedUnaryOpTestcases(s.except(s.t())...), []*unaryOpTestcase{
		{
			v:    node.Time(util.Must(time.Parse(time.DateTime, "2026-01-02 10:11:12"))),
			want: node.Int(2026),
		},
	})
}

func TestMonth(t *testing.T) {
	s := defaultFailedTestcaseSeed()
	runUnaryOpTest(t, func(v node.Data) (node.Data, error) {
		x, err := v.AsOp().Month()
		if err != nil {
			return nil, err
		}
		return x.AsData(), nil
	}, newFailedUnaryOpTestcases(s.except(s.t())...), []*unaryOpTestcase{
		{
			v:    node.Time(util.Must(time.Parse(time.DateTime, "2026-01-02 10:11:12"))),
			want: node.Int(1),
		},
	})
}

func TestDay(t *testing.T) {
	s := defaultFailedTestcaseSeed()
	runUnaryOpTest(t, func(v node.Data) (node.Data, error) {
		x, err := v.AsOp().Day()
		if err != nil {
			return nil, err
		}
		return x.AsData(), nil
	}, newFailedUnaryOpTestcases(s.except(s.t())...), []*unaryOpTestcase{
		{
			v:    node.Time(util.Must(time.Parse(time.DateTime, "2026-01-02 10:11:12"))),
			want: node.Int(2),
		},
	})
}

func TestHour(t *testing.T) {
	s := defaultFailedTestcaseSeed()
	runUnaryOpTest(t, func(v node.Data) (node.Data, error) {
		x, err := v.AsOp().Hour()
		if err != nil {
			return nil, err
		}
		return x.AsData(), nil
	}, newFailedUnaryOpTestcases(s.except(s.t())...), []*unaryOpTestcase{
		{
			v:    node.Time(util.Must(time.Parse(time.DateTime, "2026-01-02 10:11:12"))),
			want: node.Int(10),
		},
	})
}

func TestMinute(t *testing.T) {
	s := defaultFailedTestcaseSeed()
	runUnaryOpTest(t, func(v node.Data) (node.Data, error) {
		x, err := v.AsOp().Minute()
		if err != nil {
			return nil, err
		}
		return x.AsData(), nil
	}, newFailedUnaryOpTestcases(s.except(s.t())...), []*unaryOpTestcase{
		{
			v:    node.Time(util.Must(time.Parse(time.DateTime, "2026-01-02 10:11:12"))),
			want: node.Int(11),
		},
	})
}

func TestSecond(t *testing.T) {
	s := defaultFailedTestcaseSeed()
	runUnaryOpTest(t, func(v node.Data) (node.Data, error) {
		x, err := v.AsOp().Second()
		if err != nil {
			return nil, err
		}
		return x.AsData(), nil
	}, newFailedUnaryOpTestcases(s.except(s.t())...), []*unaryOpTestcase{
		{
			v:    node.Time(util.Must(time.Parse(time.DateTime, "2026-01-02 10:11:12"))),
			want: node.Int(12),
		},
	})
}

func TestDayOfWeek(t *testing.T) {
	s := defaultFailedTestcaseSeed()
	runUnaryOpTest(t, func(v node.Data) (node.Data, error) {
		x, err := v.AsOp().DayOfWeek()
		if err != nil {
			return nil, err
		}
		return x.AsData(), nil
	}, newFailedUnaryOpTestcases(s.except(s.t())...), []*unaryOpTestcase{
		{
			v:    node.Time(util.Must(time.Parse(time.DateTime, "2026-01-02 10:11:12"))),
			want: node.Int(6),
		},
	})
}

func TestDayOfYear(t *testing.T) {
	s := defaultFailedTestcaseSeed()
	runUnaryOpTest(t, func(v node.Data) (node.Data, error) {
		x, err := v.AsOp().DayOfYear()
		if err != nil {
			return nil, err
		}
		return x.AsData(), nil
	}, newFailedUnaryOpTestcases(s.except(s.t())...), []*unaryOpTestcase{
		{
			v:    node.Time(util.Must(time.Parse(time.DateTime, "2026-01-02 10:11:12"))),
			want: node.Int(2),
		},
		{
			v:    node.Time(util.Must(time.Parse(time.DateTime, "2026-12-31 10:11:12"))),
			want: node.Int(365),
		},
	})
}

func TestNewTime(t *testing.T) {
	runVariadicOpTest(t, func(v ...node.Data) (node.Data, error) {
		a := DataListToOpList(v...)
		x, err := a[0].NewTime(a[1:]...)
		if err != nil {
			return nil, err
		}
		return x.AsData(), nil
	}, []*variadicOpTestacase{
		{
			v: []node.Data{
				node.Int(2026),
			},
			want: node.Time(time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)),
		},
		{
			v: []node.Data{
				node.Int(2026),
				node.Int(2),
			},
			want: node.Time(time.Date(2026, time.February, 1, 0, 0, 0, 0, time.UTC)),
		},
		{
			v: []node.Data{
				node.Int(2026),
				node.Int(2),
				node.Int(3),
			},
			want: node.Time(time.Date(2026, time.February, 3, 0, 0, 0, 0, time.UTC)),
		},
		{
			v: []node.Data{
				node.Int(2026),
				node.Int(2),
				node.Int(3),
				node.Int(4),
			},
			want: node.Time(time.Date(2026, time.February, 3, 4, 0, 0, 0, time.UTC)),
		},
		{
			v: []node.Data{
				node.Int(2026),
				node.Int(2),
				node.Int(3),
				node.Int(4),
				node.Int(5),
			},
			want: node.Time(time.Date(2026, time.February, 3, 4, 5, 0, 0, time.UTC)),
		},
		{
			v: []node.Data{
				node.Int(2026),
				node.Int(2),
				node.Int(3),
				node.Int(4),
				node.Int(5),
				node.Int(6),
			},
			want: node.Time(time.Date(2026, time.February, 3, 4, 5, 6, 0, time.UTC)),
		},
		{
			v: []node.Data{
				node.Int(2026),
				node.Int(2),
				node.Int(3),
				node.Int(4),
				node.Int(5),
				node.Int(6),
				node.Int(7),
			},
			want: node.Time(time.Date(2026, time.February, 3, 4, 5, 6, 0, time.UTC)),
		},
	})
}

func TestSleep(t *testing.T) {
	s := defaultFailedTestcaseSeed()
	runUnaryOpTest(t, func(x node.Data) (node.Data, error) {
		r, err := x.AsOp().Sleep()
		if err != nil {
			return nil, err
		}
		return r.AsData(), nil
	}, newFailedUnaryOpTestcases(s.except(s.f(), s.i(), s.d())...), []*unaryOpTestcase{
		{
			v:    node.Float(0.1),
			want: node.Int(0),
		},
		{
			v:    node.Int(1),
			want: node.Int(0),
		},
		{
			v:    node.Duration(100 * time.Millisecond),
			want: node.Int(0),
		},
	},
	)
}
