package node

func (v *Op) Not() (*Op, error) {
	switch d := v.data.(type) {
	case Float:
		return Float(-d.Raw()).AsOp(), nil
	case Int:
		return Int(-d.Raw()).AsOp(), nil
	case Bool:
		return Bool(!d.Raw()).AsOp(), nil
	case Duration:
		return Duration(-d.Raw()).AsOp(), nil
	}
	return nil, unavailable("Not", v)
}
