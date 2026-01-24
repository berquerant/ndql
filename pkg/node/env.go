package node

import "os"

func (v *Op) EnvOr(other *Op) (*Op, error) {
	switch d := v.data.(type) {
	case String:
		s := os.Getenv(d.Raw())
		if s != "" {
			return String(s).AsOp(), nil
		}
		return other, nil
	}
	return nil, unavailable("EnvOr", v, other)
}

func (v *Op) Env() (*Op, error) {
	r, err := v.EnvOr(NewNull().AsOp())
	if err != nil {
		return nil, unavailableErr(err, "Env", v)
	}
	return r, nil
}
