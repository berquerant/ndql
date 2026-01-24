package node

import (
	"slices"

	"github.com/berquerant/ndql/pkg/util"
)

type CompareResult = util.CompareResult

const (
	CmpUnknown = util.CompareUnknown
	CmpLess    = util.CompareLess
	CmpEqual   = util.CompareEqual
	CmpGreater = util.CompareGreater
)

func (v *Op) Compare(other *Op) CompareResult {
	switch d := v.data.(type) {
	case Null:
		if util.OK(other.Null()) {
			return CmpEqual
		}
	case Float:
		switch e := other.data.(type) {
		case Float:
			return util.Compare(d.Raw(), e.Raw())
		case Int:
			return util.Compare(d.Raw(), float64(e.Raw()))
		}
	case Int:
		switch e := other.data.(type) {
		case Float:
			return util.Compare(float64(d.Raw()), e.Raw())
		case Int:
			return util.Compare(d.Raw(), e.Raw())
		}
	case Bool:
		if util.OK(other.Bool()) {
			return util.Must(v.AsInt()).AsOp().Compare(util.Must(other.AsInt()).AsOp())
		}
	case String:
		if util.OK(other.String()) {
			return util.Compare(d.Raw(), util.MustOK(other.String()).Raw())
		}
	case Time:
		if util.OK(other.Time()) {
			return util.Must(v.AsInt()).AsOp().Compare(util.Must(other.AsInt()).AsOp())
		}
	case Duration:
		if util.OK(other.Duration()) {
			return util.Must(v.AsInt()).AsOp().Compare(util.Must(other.AsInt()).AsOp())
		}
	}
	return CmpUnknown
}

func (v *Op) IsNull() bool {
	_, err := v.AsNull()
	return err == nil
}

func (v *Op) IsTrue() bool {
	x, err := v.AsBool()
	return err == nil && x.Raw()
}

func (v *Op) IsFalse() bool {
	x, err := v.AsBool()
	return err == nil && !x.Raw()
}

func (v *Op) In(other ...*Op) bool {
	return slices.ContainsFunc(other, func(x *Op) bool {
		return v.Compare(x) == util.CompareEqual
	})
}

func (v *Op) Between(left, right *Op) (bool, error) {
	if !(v.EqualType(left) && left.EqualType(right)) {
		return false, withUnavailable("Between", []*Op{v, left, right}, "type conflict")
	}
	return slices.Contains(
		[]CompareResult{util.CompareEqual, util.CompareGreater},
		v.Compare(left),
	) && slices.Contains(
		[]CompareResult{util.CompareEqual, util.CompareLess},
		v.Compare(right),
	), nil
}

func (v *Op) Least(other ...*Op) *Op {
	r := v
	for _, x := range other {
		switch x.Compare(r) {
		case CmpLess:
			r = x
		case CmpUnknown:
			return NewNull().AsOp()
		}
	}
	return r
}

func (v *Op) Greatest(other ...*Op) *Op {
	r := v
	for _, x := range other {
		switch x.Compare(r) {
		case CmpGreater:
			r = x
		case CmpUnknown:
			return NewNull().AsOp()
		}
	}
	return r
}

func (v *Op) Coalesce(other ...*Op) *Op {
	if !v.IsNull() {
		return v
	}
	for _, x := range other {
		if !x.IsNull() {
			return x
		}
	}
	return NewNull().AsOp()
}
