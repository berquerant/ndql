package run

import (
	"context"
	"fmt"

	"github.com/berquerant/ndql/pkg/iterx"
	"github.com/berquerant/ndql/pkg/parse"
	"github.com/berquerant/ndql/pkg/tree"
	"golang.org/x/sync/errgroup"
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

	var (
		eg, eCtx = errgroup.WithContext(ctx)
		recvCs   = make([]chan *tree.N, len(p.Nodes))
	)
	for i := range p.Nodes {
		recvC := make(chan *tree.N, 100)
		recvCs[i] = recvC
		eg.Go(func() error {
			for n := range recvC {
				r.WriteNode(n)
			}
			return nil
		})
	}

	concurrency := int(r.Concurrency)
	for i, n := range p.Nodes {
		vit := cit.Clone()
		eg.Go(func() error {
			if err := tree.AsChan(eCtx, vit.Values(), n, concurrency, recvCs[i]); err != nil {
				return fmt.Errorf("%w: node[%d]", err, i)
			}
			return nil
		})
	}

	return eg.Wait()
}
