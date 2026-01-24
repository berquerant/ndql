package node

func (v *Op) If(left, right *Op) *Op {
	if v.IsTrue() {
		return left
	}
	return right
}

func (v *Op) IfNull(other *Op) *Op {
	if v.IsNull() {
		return other
	}
	return v
}

func (v *Op) NullIf(other *Op) *Op {
	switch v.Compare(other) {
	case CmpEqual:
		return NewNull().AsOp()
	default:
		return v
	}
}
