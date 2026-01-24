package tree

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"text/template"

	"github.com/berquerant/ndql/pkg/cachex"
	"github.com/berquerant/ndql/pkg/util"
)

//
// string gen template (text/template)
//
// tmpl(template)
//
// Data object is a node.
// Additional functions are defined.
//
// ## Functions
// ### env
// os.Getenv.

// ### envor
// os.Getenv with default value, like 'envor "KEY" "default_value"'.
type StringGenTemplate struct {
	text string
}

func NewStringGenTemplate(text string) *StringGenTemplate {
	return &StringGenTemplate{
		text: text,
	}
}

const (
	genTemplateFuncEnv   = "env"
	genTemplateFuncEnvOr = "envor"
	// genTemplateFuncGet   = "get"
	// genTemplateFuncGetOr = "getor"
)

var _ GenTemplate = &StringGenTemplate{}

func (g StringGenTemplate) Generate(_ context.Context, n *N) ([]byte, error) {
	t, err := stringGenTemplateCache.Get(g.text)
	if err != nil {
		return nil, fmt.Errorf("%w: cannot get text template", errors.Join(ErrGenTemplate, err))
	}
	var out bytes.Buffer
	if err := t.Execute(&out, NodeAsStructuredMap(n)); err != nil {
		return nil, fmt.Errorf("%w: failed to execute", errors.Join(ErrGenTemplate, err))
	}
	return out.Bytes(), nil
}

var stringGenTemplateCache = util.Must(cachex.NewTextTemplateCache(template.FuncMap{
	genTemplateFuncEnv:   os.Getenv,
	genTemplateFuncEnvOr: genTemplateEnvOr,
	// genTemplateFuncGetOr: genTemplateGetOr,
	// genTemplateFuncGet:   genTemplateGet,
}))
