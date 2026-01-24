package run

import (
	"context"
	"fmt"

	"github.com/berquerant/ndql/pkg/iterx"
	"github.com/berquerant/ndql/pkg/parse"
	"github.com/berquerant/ndql/pkg/tree"
)

func (r *runner) query(ctx context.Context) error {
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
	it, err := r.Sources.ReadInput()
	if err != nil {
		return err
	}
	cit := iterx.NewClonableIter(it)
	defer cit.Close()
	rit := make([]tree.NIter, len(p.Nodes))
	for i, n := range p.Nodes {
		vit := cit.Clone()
		xit, err := tree.AsIter(ctx, vit.Values(), n)
		if err != nil {
			return fmt.Errorf("%w: node[%d]", err, i)
		}
		rit[i] = xit
	}
	for _, it := range rit {
		for n := range it {
			r.WriteNode(n)
		}
	}
	return nil
}
