package run

func (r *runner) list() error {
	if err := r.SetupSources(); err != nil {
		return err
	}

	it, err := r.Sources.ReadInput()
	if err != nil {
		return err
	}
	for n := range it {
		r.WriteNode(n)
	}
	return nil
}
