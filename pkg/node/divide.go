package node

func (v *Op) Divide(other *Op) (*Op, error) {
	e, err := other.Inverse()
	if err != nil {
		return nil, unavailableErr(err, "Divide", v, other)
	}
	r, err := v.Multiply(e)
	if err != nil {
		return nil, unavailableErr(err, "Divide", v, other)
	}
	return r, nil
}
