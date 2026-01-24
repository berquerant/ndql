package node

import (
	"slices"
)

func (v *Op) Inverse() (*Op, error) {
	switch d := v.data.(type) {
	case Float:
		x := d.Raw()
		if x == 0 {
			return nil, withUnavailable("Inverse", []*Op{v}, "div by zeoo")
		}
		return Float(1 / x).AsOp(), nil
	case Int:
		x := d.Raw()
		if x == 0 {
			return nil, withUnavailable("Inverse", []*Op{v}, "div by zeoo")
		}
		return Float(1 / float64(x)).AsOp(), nil
	case String:
		x := []rune(d.Raw())
		slices.Reverse(x)
		return String(string(x)).AsOp(), nil
	}
	return nil, unavailable("Inverse", v)
}
