package tree

import (
	"context"
	"errors"
	"log/slog"

	"github.com/berquerant/ndql/pkg/logx"
	. "github.com/pingcap/tidb/pkg/parser/ast"
	"golang.org/x/sync/errgroup"
)

// AsIter converts AST into a function and applies to the iterator.
func AsIter(ctx context.Context, it NIter, n Node) (NIter, error) {
	function, err := NewTreeVisitor(ctx).Visit(n)
	if err != nil {
		return nil, err
	}
	return func(yield func(*N) bool) {
		for x := range it {
			rs, err := function.CallAny(x)
			if errors.Is(err, ErrIgnore) {
				logx.Trace("ignore node", logx.Err(err))
				continue
			}
			if err != nil {
				slog.Debug("failed to yield node", logx.JSON("node", x), logx.Err(err))
				continue
			}
			for _, r := range rs {
				if !yield(r) {
					return
				}
			}
		}
	}, nil
}

func AsChan(ctx context.Context, it NIter, n Node, concurrency int, recvC chan<- *N) error {
	defer close(recvC)

	if concurrency < 1 {
		concurrency = 1
	}
	function, err := NewTreeVisitor(ctx).Visit(n)
	if err != nil {
		return err
	}

	var (
		sendC    = make(chan *N, 100)
		eg, eCtx = errgroup.WithContext(ctx)
	)
	for range concurrency {
		eg.Go(func() error {
			for x := range sendC {
				select {
				case <-eCtx.Done():
					return eCtx.Err()
				default:
					rs, err := function.CallAny(x)
					if errors.Is(err, ErrIgnore) {
						logx.Trace("ignore node", logx.Err(err))
						continue
					}
					if err != nil {
						slog.Debug("failed to yield node", logx.JSON("node", x), logx.Err(err))
						continue
					}
					for _, r := range rs {
						select {
						case <-eCtx.Done():
							return eCtx.Err()
						default:
							recvC <- r
						}
					}
				}
			}
			return nil
		})
	}

	for x := range it {
		select {
		case <-eCtx.Done():
			return eCtx.Err()
		default:
			sendC <- x
		}
	}
	close(sendC)

	return eg.Wait()
}
