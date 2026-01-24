package node

import (
	"math"

	"github.com/berquerant/ndql/pkg/util"
)

func (v *Op) Mod(other *Op) (*Op, error) {
	switch d := v.data.(type) {
	case Float:
		x := d.Raw()
		switch e := other.data.(type) {
		case Float:
			if y := e.Raw(); y != 0 {
				return Float(math.Mod(x, y)).AsOp(), nil
			}
			return nil, withUnavailable("Mod", []*Op{v, other}, "div by zero")
		case Int:
			return v.Mod(util.Must(other.AsFloat()).AsOp())
		}
	case Int:
		x := d.Raw()
		switch e := other.data.(type) {
		case Float:
			return util.Must(v.AsFloat()).AsOp().Mod(other)
		case Int:
			if y := e.Raw(); y != 0 {
				return Int(x % y).AsOp(), nil
			}
			return nil, withUnavailable("Mod", []*Op{v, other}, "div by zero")
		}
	}
	return nil, unavailable("Mod", v, other)
}
