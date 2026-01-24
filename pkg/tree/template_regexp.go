package tree

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/berquerant/ndql/pkg/iox"
	"github.com/berquerant/ndql/pkg/regexpx"
)

//
// regexp gen template
//
// grep(pattern, template)

type RegexpGenTemplate struct {
	expr string
	tmpl string
}

func NewRegexpGenTemplate(expr, tmpl string) *RegexpGenTemplate {
	return &RegexpGenTemplate{
		expr: expr,
		tmpl: tmpl,
	}
}

var _ GenTemplate = &RegexpGenTemplate{}

func (g RegexpGenTemplate) Generate(_ context.Context, n *N) ([]byte, error) {
	re, err := regexpx.Compile(g.expr)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to compile expr", errors.Join(ErrGenTemplate, err))
	}

	p := n.GetPath()
	c, err := regexpGenTemplateCache.Get(p.Raw())
	if err != nil {
		return nil, fmt.Errorf("%w: failed to read content to grep, file=%s", errors.Join(ErrGenTemplate, err), p.Raw())
	}

	matches := re.FindAllSubmatchIndex(c, -1)
	result := make([][]byte, len(matches))
	for i, sm := range matches {
		r := []byte{}
		r = re.Expand(r, []byte(g.tmpl), c, sm)
		if len(r) > 0 {
			result[i] = r
		}
	}

	return bytes.Join(result, []byte("\n")), nil
}

var regexpGenTemplateCache = iox.NewFileContentCache()
