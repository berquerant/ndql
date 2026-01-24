package run

import "github.com/berquerant/ndql/version"

func (r *runner) version() error {
	version.Write(r.Stdout)
	return nil
}
