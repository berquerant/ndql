package node

import "path/filepath"

func (v *Op) Dir() (*Op, error) {
	switch d := v.data.(type) {
	case String:
		return String(filepath.Dir(d.Raw())).AsOp(), nil
	default:
		return nil, unavailable("Dir", v)
	}
}

func (v *Op) Basename() (*Op, error) {
	switch d := v.data.(type) {
	case String:
		return String(filepath.Base(d.Raw())).AsOp(), nil
	default:
		return nil, unavailable("Basename", v)
	}
}

func (v *Op) Extension() (*Op, error) {
	switch d := v.data.(type) {
	case String:
		return String(filepath.Ext(d.Raw())).AsOp(), nil
	default:
		return nil, unavailable("Extension", v)
	}
}

func (v *Op) AbsPath() (*Op, error) {
	switch d := v.data.(type) {
	case String:
		r, err := filepath.Abs(d.Raw())
		if err != nil {
			return nil, unavailableErr(err, "AbsPath", v)
		}
		return String(r).AsOp(), nil
	default:
		return nil, unavailable("AbsPath", v)
	}
}

func (v *Op) RelPath(basepath *Op) (*Op, error) {
	switch d := v.data.(type) {
	case String:
		e, ok := basepath.String()
		if !ok {
			return nil, withUnavailable("RelPath", []*Op{v, basepath}, "args[1] should be String")
		}
		r, err := filepath.Rel(e.Raw(), d.Raw())
		if err != nil {
			return nil, unavailableErr(err, "RelPath", v, basepath)
		}
		return String(r).AsOp(), nil
	}
	return nil, unavailable("RelPath", v, basepath)
}
