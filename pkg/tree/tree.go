package tree

import (
	"context"
	"errors"
	"log/slog"

	"github.com/berquerant/ndql/pkg/logx"
	. "github.com/pingcap/tidb/pkg/parser/ast"
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
