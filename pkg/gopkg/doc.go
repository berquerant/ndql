package gopkg

import (
	"errors"
	"fmt"
	"strings"
)

type Annotation struct {
	Key   string
	Value string
	Linum int // 0 based
}

func (c *Comment) GetAnnotation(key string) (*Annotation, bool) {
	for linum, line := range strings.Split(c.Text, "\n") {
		xs := strings.Split(line, " ")
		for i, x := range xs {
			if i > 0 {
				break
			}
			if x == "@"+key {
				return &Annotation{
					Key:   key,
					Value: strings.Join(xs[1:], " "),
					Linum: linum,
				}, true
			}
		}
	}
	return nil, false
}

type Document struct {
	Path  string `json:"path"`  // @path
	Text  string `json:"text"`  // @document
	Title string `json:"title"` // @title
	File  string `json:"file"`
	Line  int    `json:"line"`
}

const (
	AnnotationPath     = "path"     // jsonpath
	AnnotationDocument = "document" // begin document
	AnnotationTitle    = "title"    // for markdown heading: #
)

var (
	ErrDocument = errors.New("DocumentError")
)

func (c *Comment) GetDocument() (*Document, error) {
	err := fmt.Errorf("%w: file %s line %d", ErrDocument, c.Path, c.Line)
	p, ok := c.GetAnnotation(AnnotationPath)
	if !ok || p.Value == "" {
		return nil, fmt.Errorf("%w: no path", err)
	}
	t, ok := c.GetAnnotation(AnnotationTitle)
	if !ok || t.Value == "" {
		return nil, fmt.Errorf("%w: no title", err)
	}
	d, ok := c.GetAnnotation(AnnotationDocument)
	if !ok {
		return nil, fmt.Errorf("%w: no document", err)
	}
	if d.Linum < p.Linum {
		return nil, fmt.Errorf("%w: document should preceed path", err)
	}
	if d.Linum < t.Linum {
		return nil, fmt.Errorf("%w: document should preceed title", err)
	}
	content := strings.Join(strings.Split(c.Text, "\n")[d.Linum+1:], "\n")
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, fmt.Errorf("%w: no content", err)
	}
	return &Document{
		Path:  p.Value,
		Text:  content,
		Title: t.Value,
		File:  c.Path,
		Line:  c.Line,
	}, nil
}

func (d *Document) String() string {
	return fmt.Sprintf(`# %s

%s`, d.Title, d.Text)
}
