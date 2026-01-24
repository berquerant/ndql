package node

import (
	"log/slog"
	"time"
)

func Now() *Op { return Time(time.Now()).AsOp() }

func (v *Op) TimeFormat(other *Op) (*Op, error) {
	switch d := v.data.(type) {
	case Time:
		switch e := other.data.(type) {
		case String:
			return String(d.Raw().Format(e.Raw())).AsOp(), nil
		}
	}
	return nil, unavailable("TimeFormat", v, other)
}

func (v *Op) Year() (*Op, error) {
	switch d := v.data.(type) {
	case Time:
		return Int(int64(d.Raw().Year())).AsOp(), nil
	default:
		return nil, unavailable("Year", v)
	}
}

func (v *Op) Month() (*Op, error) {
	switch d := v.data.(type) {
	case Time:
		return Int(int64(d.Raw().Month())).AsOp(), nil
	default:
		return nil, unavailable("Month", v)
	}
}

func (v *Op) Day() (*Op, error) {
	switch d := v.data.(type) {
	case Time:
		return Int(int64(d.Raw().Day())).AsOp(), nil
	default:
		return nil, unavailable("Day", v)
	}
}

func (v *Op) Hour() (*Op, error) {
	switch d := v.data.(type) {
	case Time:
		return Int(int64(d.Raw().Hour())).AsOp(), nil
	default:
		return nil, unavailable("Hour", v)
	}
}

func (v *Op) Minute() (*Op, error) {
	switch d := v.data.(type) {
	case Time:
		return Int(int64(d.Raw().Minute())).AsOp(), nil
	default:
		return nil, unavailable("Minute", v)
	}
}

func (v *Op) Second() (*Op, error) {
	switch d := v.data.(type) {
	case Time:
		return Int(int64(d.Raw().Second())).AsOp(), nil
	default:
		return nil, unavailable("Second", v)
	}
}

func (v *Op) DayOfWeek() (*Op, error) {
	switch d := v.data.(type) {
	case Time:
		return Int(int64(d.Raw().Weekday()) + 1).AsOp(), nil
	default:
		return nil, unavailable("DayOfWeek", v)
	}
}

func (v *Op) DayOfYear() (*Op, error) {
	switch d := v.data.(type) {
	case Time:
		return Int(int64(d.Raw().YearDay())).AsOp(), nil
	default:
		return nil, unavailable("DayOfYear", v)
	}
}

func (v *Op) NewTime(param ...*Op) (*Op, error) {
	xs := []int{
		0, // year
		1, // month
		1, // day
		0, // hour
		0, // minute
		0, // second
	}
	pp := append([]*Op{v}, param...)
	get := func() (int, bool, bool) {
		if len(pp) == 0 {
			return 0, false, true
		}
		x, ok := pp[0].Int()
		pp = pp[1:]
		return int(x.Raw()), ok, false
	}
	for i := 0; i < len(xs); i++ {
		x, ok, end := get()
		if end {
			break
		}
		if !ok {
			return nil, withUnavailable("NewTime", append([]*Op{v}, param...), "invalid arg[%d]", i)
		}
		xs[i] = x
	}
	return Time(time.Date(xs[0], time.Month(xs[1]), xs[2], xs[3], xs[4], xs[5], 0, time.UTC)).AsOp(), nil
}

func (v *Op) Sleep() (*Op, error) {
	var x time.Duration
	switch d := v.data.(type) {
	case Float:
		x = time.Duration(int64(d.Raw() * float64(time.Second)))
	case Int:
		x = time.Duration(d.Raw()) * time.Second
	case Duration:
		x = d.Raw()
	default:
		return nil, unavailable("Sleep", v)
	}
	slog.Debug("Sleep", slog.Duration("duration", x))
	time.Sleep(x)
	return Int(0).AsOp(), nil
}
