package node

import (
	"strings"

	"github.com/berquerant/ndql/pkg/util"
)

func (v *Op) Multiply(other *Op) (*Op, error) {
	switch d := v.data.(type) {
	case Float:
		switch e := other.data.(type) {
		case Float:
			return Float(util.Multiply(d, e)).AsOp(), nil
		case Int:
			return v.Multiply(util.Must(other.AsFloat()).AsOp())
		case Duration:
			return Duration(util.NewDuration(d.Raw(), e.Raw())).AsOp(), nil
		}
	case Int:
		switch e := other.data.(type) {
		case Float:
			return util.Must(v.AsFloat()).AsOp().Multiply(other)
		case Int:
			return Int(util.Multiply(d, e)).AsOp(), nil
		case String:
			if c := int(d.Raw()); c > 0 {
				return String(strings.Repeat(e.Raw(), c)).AsOp(), nil
			}
			return Default().String().AsOp(), nil
		case Duration:
			return util.Must(v.AsFloat()).AsOp().Multiply(other)
		}
	case Bool:
		if e, ok := other.Bool(); ok {
			return Bool(d.Raw() && e.Raw()).AsOp(), nil
		}
	case String:
		if util.OK(other.Int()) {
			return other.Multiply(v)
		}
	case Duration:
		switch other.data.(type) {
		case Float:
			return other.Multiply(v)
		case Int:
			return other.Multiply(v)
		}
	}
	return nil, unavailable("Multiply", v, other)
}
