package node

import (
	"fmt"
	"strconv"
	"time"

	"github.com/berquerant/ndql/pkg/util"
)

func (v *Op) AsNull() (Null, error) {
	if x, ok := v.data.(Null); ok {
		return x, nil
	}
	return NewNull(), unavailable("AsNull", v)
}

func (v *Op) AsFloat() (Float, error) {
	switch d := v.data.(type) {
	case Float:
		return d, nil
	case Int:
		return Float(float64(d.Raw())), nil
	case Bool:
		return util.Must(v.AsInt()).AsOp().AsFloat()
	case String:
		x, err := strconv.ParseFloat(d.Raw(), 64)
		if err != nil {
			return Default().Float(), unavailableErr(err, "AsFloat", v)
		}
		return Float(x), nil
	case Time:
		return util.Must(v.AsInt()).AsOp().AsFloat()
	case Duration:
		return util.Must(v.AsInt()).AsOp().AsFloat()
	default:
		return Default().Float(), unavailable("AsFloat", v)
	}
}

func (v *Op) AsInt() (Int, error) {
	switch d := v.data.(type) {
	case Float:
		return Int(int64(d.Raw())), nil
	case Int:
		return d, nil
	case Bool:
		if d.Raw() {
			return Int(1), nil
		}
		return Int(0), nil
	case String:
		x, err := strconv.ParseInt(d.Raw(), 10, 64)
		if err != nil {
			return Default().Int(), unavailableErr(err, "AsInt", v)
		}
		return Int(x), nil
	case Time:
		return Int(d.Raw().Unix()), nil
	case Duration:
		return Int(d.Raw()), nil
	default:
		return Default().Int(), unavailable("AsInt", v)
	}
}

func (v *Op) AsBool() (Bool, error) {
	switch d := v.data.(type) {
	case Float:
		return util.Must(v.AsInt()).AsOp().AsBool()
	case Int:
		return Bool(d.Raw() != 0), nil
	case Bool:
		return d, nil
	case String:
		return Bool(d.Raw() != ""), nil
	default:
		return Default().Bool(), unavailable("AsBool", v)
	}
}

func (v *Op) AsString() (String, error) {
	switch d := v.data.(type) {
	case Null:
		return "null", nil
	case Float:
		return String(fmt.Sprint(d.Raw())), nil
	case Int:
		return String(fmt.Sprint(d.Raw())), nil
	case Bool:
		return String(fmt.Sprint(d.Raw())), nil
	case String:
		return d, nil
	case Time:
		return String(d.Raw().Format(time.DateTime)), nil
	case Duration:
		return String(d.Raw().String()), nil
	default:
		return Default().String(), unavailable("AsString", v)
	}
}

func (v *Op) AsTime() (Time, error) {
	switch d := v.data.(type) {
	case Float:
		return util.Must(v.AsInt()).AsOp().AsTime()
	case Int:
		return Time(time.Unix(d.Raw(), 0).In(time.UTC)), nil
	case String:
		x, err := time.Parse(time.DateTime, d.Raw())
		if err != nil {
			return Default().Time(), unavailableErr(err, "AsTime", v)
		}
		return Time(x), nil
	case Time:
		return d, nil
	default:
		return Default().Time(), unavailable("AsTime", v)
	}
}

func (v *Op) AsDuration() (Duration, error) {
	switch d := v.data.(type) {
	case Float:
		return util.Must(v.AsInt()).AsOp().AsDuration()
	case Int:
		return Duration(time.Duration(d.Raw())), nil
	case String:
		x, err := time.ParseDuration(d.Raw())
		if err != nil {
			return Default().Duration(), unavailableErr(err, "AsDuration", v)
		}
		return Duration(x), nil
	case Duration:
		return d, nil
	default:
		return Default().Duration(), unavailable("AsDuration", v)
	}
}
