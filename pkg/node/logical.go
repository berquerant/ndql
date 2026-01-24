package node

import "github.com/berquerant/ndql/pkg/util"

func (v *Op) LogicalNot() (*Op, error) {
	if util.OK(v.Bool()) {
		return v.Not()
	}
	return nil, unavailable("LogicalNot", v)
}

func (v *Op) LogicalAnd(other *Op) (*Op, error) {
	if d, ok := v.Bool(); ok {
		if e, ok := other.Bool(); ok {
			return Bool(d.Raw() && e.Raw()).AsOp(), nil
		}
	}
	return nil, unavailable("LogialAnd", v, other)
}

func (v *Op) LogicalOr(other *Op) (*Op, error) {
	if d, ok := v.Bool(); ok {
		if e, ok := other.Bool(); ok {
			return Bool(d.Raw() || e.Raw()).AsOp(), nil
		}
	}
	return nil, unavailable("LogicalOr", v, other)
}

func (v *Op) LogicalXor(other *Op) (*Op, error) {
	if d, ok := v.Bool(); ok {
		if e, ok := other.Bool(); ok {
			return Bool(d.Raw() != e.Raw()).AsOp(), nil
		}
	}
	return nil, unavailable("LogicalXor", v, other)
}
