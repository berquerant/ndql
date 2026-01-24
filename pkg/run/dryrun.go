package run

import "github.com/berquerant/ndql/pkg/parse"

func (r *runner) dryrun() error {
	if err := r.SetupSources(); err != nil {
		return err
	}
	if err := r.SetupQuery(); err != nil {
		return err
	}
	p, err := parse.NewSQLParser().Parse(r.Query)
	if err != nil {
		return err
	}
	mode := r.dumpMode()
	for _, n := range p.Nodes {
		parse.Dump(r.Stdout, n, "|", mode)
	}
	return nil
}

func (r *runner) dumpMode() parse.DumpMode {
	if r.Verbose {
		return parse.DumpModeVerbose
	}
	return parse.DumpModeText
}
