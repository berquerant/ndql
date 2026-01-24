package node

import "github.com/berquerant/ndql/pkg/util"

func (v *Op) BitNot() (*Op, error) {
	if d, ok := v.Int(); ok {
		return Int(^d.Raw()).AsOp(), nil
	}
	return nil, unavailable("BitNot", v)
}

func (v *Op) BitAnd(other *Op) (*Op, error) {
	if d, ok := v.Int(); ok {
		if e, ok := other.Int(); ok {
			return Int(d.Raw() & e.Raw()).AsOp(), nil
		}
	}
	return nil, unavailable("BitAnd", v, other)
}

func (v *Op) BitOr(other *Op) (*Op, error) {
	if d, ok := v.Int(); ok {
		if e, ok := other.Int(); ok {
			return Int(d.Raw() | e.Raw()).AsOp(), nil
		}
	}
	return nil, unavailable("BitOr", v, other)
}

func (v *Op) BitXor(other *Op) (*Op, error) {
	if d, ok := v.Int(); ok {
		if e, ok := other.Int(); ok {
			return Int(d.Raw() ^ e.Raw()).AsOp(), nil
		}
	}
	return nil, unavailable("BitXor", v, other)
}

func (v *Op) LeftShift(other *Op) (*Op, error) {
	if d, ok := v.Int(); ok {
		if e, ok := other.Int(); ok {
			x := e.Raw()
			if x < 0 {
				return v.RightShift(util.Must(other.Not()))
			}
			return Int(d.Raw() << x).AsOp(), nil
		}
	}
	return nil, unavailable("LeftShift", v, other)
}

func (v *Op) RightShift(other *Op) (*Op, error) {
	if d, ok := v.Int(); ok {
		if e, ok := other.Int(); ok {
			x := e.Raw()
			if x < 0 {
				return v.LeftShift(util.Must(other.Not()))
			}
			return Int(d.Raw() >> x).AsOp(), nil
		}
	}
	return nil, unavailable("RightShift", v, other)
}
