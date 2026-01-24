package node

func (v *Op) Subtract(other *Op) (*Op, error) {
	e, err := other.Not()
	if err != nil {
		return nil, unavailableErr(err, "Subtract", v, other)
	}
	r, err := v.Add(e)
	if err != nil {
		return nil, unavailableErr(err, "Subtract", v, other)
	}
	return r, nil
}
