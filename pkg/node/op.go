package node

import (
	"errors"
	"fmt"

	"github.com/berquerant/ndql/pkg/util"
)

// Op defines some operations on Data.
type Op struct {
	data Data
}

var (
	ErrUnknownData = errors.New("UnknownData")
	ErrUnavailable = errors.New("Unavailable")
)

func withUnavailableErr(err error, op string, args []*Op, format string, v ...any) error {
	return fmt.Errorf("%w:%s(%v):%s",
		errors.Join(ErrUnavailable, err),
		op,
		opListToAny(args...),
		fmt.Sprintf(format, v...),
	)
}
func withUnavailable(op string, args []*Op, format string, v ...any) error {
	return fmt.Errorf("%w:%s(%v):%s",
		ErrUnavailable,
		op,
		opListToAny(args...),
		fmt.Sprintf(format, v...),
	)
}
func unavailableErr(err error, op string, args ...*Op) error {
	return fmt.Errorf("%w:%s(%v)",
		errors.Join(ErrUnavailable, err),
		op,
		opListToAny(args...),
	)
}
func unavailable(op string, args ...*Op) error {
	return fmt.Errorf("%w:%s(%v)",
		ErrUnavailable,
		op,
		opListToAny(args...),
	)
}

func opListToAny(v ...*Op) any {
	xs := make([]any, len(v))
	for i, x := range v {
		xs[i] = x.AsData().Any()
	}
	return xs
}

func (v *Op) AsData() Data { return v.data }

func (v *Op) Null() (Null, bool) {
	x, ok := v.data.(Null)
	return x, ok
}
func (v *Op) Float() (Float, bool) {
	x, ok := v.data.(Float)
	return x, ok
}
func (v *Op) Int() (Int, bool) {
	x, ok := v.data.(Int)
	return x, ok
}
func (v *Op) Bool() (Bool, bool) {
	x, ok := v.data.(Bool)
	return x, ok
}
func (v *Op) String() (String, bool) {
	x, ok := v.data.(String)
	return x, ok
}
func (v *Op) Time() (Time, bool) {
	x, ok := v.data.(Time)
	return x, ok
}
func (v *Op) Duration() (Duration, bool) {
	x, ok := v.data.(Duration)
	return x, ok
}

func (v *Op) EqualType(other *Op) bool {
	switch {
	case util.OK(v.Null()) && util.OK(other.Null()):
		return true
	case util.OK(v.Float()) && util.OK(other.Float()):
		return true
	case util.OK(v.Int()) && util.OK(other.Int()):
		return true
	case util.OK(v.Bool()) && util.OK(other.Bool()):
		return true
	case util.OK(v.String()) && util.OK(other.String()):
		return true
	case util.OK(v.Time()) && util.OK(other.Time()):
		return true
	case util.OK(v.Duration()) && util.OK(other.Duration()):
		return true
	default:
		return false
	}
}

func (v *Op) MarshalJSON() ([]byte, error) {
	switch v := v.data.(type) {
	case Null:
		return v.MarshalJSON()
	case Float:
		return v.MarshalJSON()
	case Int:
		return v.MarshalJSON()
	case Bool:
		return v.MarshalJSON()
	case String:
		return v.MarshalJSON()
	case Time:
		return v.MarshalJSON()
	case Duration:
		return v.MarshalJSON()
	default:
		return nil, ErrUnknownData
	}
}

func (v *Op) UnmarshalJSON(data []byte) error {
	def := Default()
	n := def.Null()
	if err := n.UnmarshalJSON(data); err == nil {
		*v = *n.AsOp()
		return nil
	}
	b := def.Bool()
	if err := b.UnmarshalJSON(data); err == nil {
		*v = *b.AsOp()
		return nil
	}
	t := def.Time()
	if err := t.UnmarshalJSON(data); err == nil {
		*v = *t.AsOp()
		return nil
	}
	d := def.Duration()
	if err := d.UnmarshalJSON(data); err == nil {
		*v = *d.AsOp()
		return nil
	}
	i := def.Int()
	if err := i.UnmarshalJSON(data); err == nil {
		*v = *i.AsOp()
		return nil
	}
	f := def.Float()
	if err := f.UnmarshalJSON(data); err == nil {
		*v = *f.AsOp()
		return nil
	}
	s := def.String()
	if err := s.UnmarshalJSON(data); err == nil {
		*v = *s.AsOp()
		return nil
	}
	return ErrUnknownData
}
