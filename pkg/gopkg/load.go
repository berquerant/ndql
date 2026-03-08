package gopkg

import (
	"context"
	"iter"
	"log/slog"

	"golang.org/x/tools/go/packages"
)

func NewLoader() *Loader {
	return &Loader{}
}

type Loader struct {
	pkgs []*packages.Package
}

func (r *Loader) Get() []*packages.Package { return r.pkgs }

func (r *Loader) Load(ctx context.Context, pattern ...string) error {
	slog.Info("Load: start", slog.Any("pattern", pattern))
	pkgs, err := packages.Load(&packages.Config{
		Context: ctx,
		Mode:    packages.NeedSyntax | packages.NeedFiles | packages.NeedTypes | packages.NeedName,
	}, pattern...)
	if err != nil {
		return err
	}
	r.pkgs = append(r.pkgs, pkgs...)
	slog.Info("Load: end", slog.Any("pattern", pattern))
	return nil
}

type Comment struct {
	Text string `json:"text"`
	Path string `json:"path"`
	Line int    `json:"line"`
}

func (r *Loader) Comments() iter.Seq[*Comment] {
	return func(yield func(*Comment) bool) {
		for _, pkg := range r.pkgs {
			for _, file := range pkg.Syntax {
				for _, comments := range file.Comments {
					pos := pkg.Fset.Position(comments.Pos())
					x := &Comment{
						Text: comments.Text(),
						Path: pos.Filename,
						Line: pos.Line,
					}
					if !yield(x) {
						return
					}
				}
			}
		}
	}
}
