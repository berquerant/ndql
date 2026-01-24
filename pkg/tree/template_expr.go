package tree

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/berquerant/ndql/pkg/cachex"
	"github.com/berquerant/ndql/pkg/util"
	"github.com/expr-lang/expr"
)

//
// expr gen template
//
// expr(expression)
//
// ## Environment
//
// - e: environment variables from os.Environ
// - n: node

type ExprGenTemplate struct {
	expr string
}

func NewExprGenTemplate(expr string) *ExprGenTemplate {
	return &ExprGenTemplate{
		expr: expr,
	}
}

var _ GenTemplate = &ExprGenTemplate{}

var exprGenTemplateCache = util.Must(cachex.NewExprCache())

const (
	exprGenTemplateEnvKey  = "e"
	exprGenTemplateNodeKey = "n"
)

func (g ExprGenTemplate) Generate(_ context.Context, n *N) ([]byte, error) {
	e, err := exprGenTemplateCache.Get(g.expr)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to compile expr", errors.Join(ErrGenTemplate, err))
	}

	r, err := expr.Run(e, g.newEnv(n))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to run expr", errors.Join(ErrGenTemplate, err))
	}
	return []byte(fmt.Sprint(r)), nil
}

func (g ExprGenTemplate) newEnv(n *N) map[string]any {
	return map[string]any{
		exprGenTemplateEnvKey:  g.environAsMap(),
		exprGenTemplateNodeKey: NodeAsStructuredMap(n),
	}
}

func (ExprGenTemplate) environAsMap() map[string]any {
	d := make(map[string]any)
	for _, x := range os.Environ() {
		xs := strings.SplitN(x, "=", 2)
		d[xs[0]] = xs[1]
	}
	return d
}
