package node

import "github.com/berquerant/ndql/pkg/util"

func (v *Op) Add(other *Op) (*Op, error) {
	switch d := v.data.(type) {
	case Float:
		switch e := other.data.(type) {
		case Float:
			return Float(util.Add(d, e)).AsOp(), nil
		case Int:
			return v.Add(util.Must(other.AsFloat()).AsOp())
		}
	case Int:
		switch e := other.data.(type) {
		case Float:
			return other.Add(v)
		case Int:
			return Int(util.Add(d, e)).AsOp(), nil
		}
	case String:
		if e, ok := other.String(); ok {
			return String(d.Raw() + e.Raw()).AsOp(), nil
		}
	case Bool:
		if e, ok := other.Bool(); ok {
			return Bool(d.Raw() || e.Raw()).AsOp(), nil
		}
	case Time:
		if e, ok := other.Duration(); ok {
			return Time(d.Raw().Add(e.Raw())).AsOp(), nil
		}
	case Duration:
		switch e := other.data.(type) {
		case Time:
			return other.Add(v)
		case Duration:
			return Duration(d.Raw() + e.Raw()).AsOp(), nil
		}
	}
	return nil, unavailable("Add", v, other)
}
