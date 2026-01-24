package config

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"iter"
	"log/slog"
	"strings"

	"github.com/berquerant/ndql/pkg/iox"
	"github.com/berquerant/ndql/pkg/logx"
	"github.com/berquerant/ndql/pkg/node"
	"github.com/berquerant/ndql/pkg/tree"
)

const (
	stdinNameLong  = "@stdin"
	stdinNameShort = "@-"
)

func DescribeSourceUsage() string {
	return fmt.Sprintf(`- %s or %s: from stdin
- @FILENAME: from file
- otherwise: as it is`, stdinNameLong, stdinNameShort)
}

func isStdinName(v string) bool {
	return v == stdinNameLong || v == stdinNameShort
}

var (
	ErrInvalidSource = errors.New("InvalidSource")
	ErrInvalidQuery  = errors.New("InvalidQuery")
	ErrStdinConflict = errors.New("StdinConflict")
	ErrInvalidInput  = errors.New("InvalidInput")
	ErrNoSources     = errors.New("NoSources")
)

type Sources struct {
	Query iox.Source
	Path  iox.Walker
	path  iox.Source // source underlying Path
	Index iox.Source
}

func NewSource(v string) (iox.Source, error) {
	switch {
	case isStdinName(v):
		return iox.NewStdinSource(), nil
	case strings.HasPrefix(v, "@") && len(v) > 1:
		filename := v[1:]
		return iox.NewFileSource(filename)
	default:
		return nil, ErrInvalidSource
	}
}

func NewSourceOrRaw(v string) (iox.Source, error) {
	s, err := NewSource(v)
	if errors.Is(err, ErrInvalidSource) {
		return iox.NewStringSource(v), nil
	}
	return s, err
}

func NewQuerySource(v string) (iox.Source, error) {
	slog.Debug("Use query", slog.String("source", v))
	r, err := NewSourceOrRaw(v)
	if err != nil {
		return nil, errors.Join(ErrInvalidQuery, err)
	}
	return r, nil
}

func NewPathSource(v string) (iox.Walker, iox.Source, error) {
	logger := slog.With(slog.String("source", v))
	logger.Debug("Use path")
	r, err := NewSource(v)
	switch {
	case errors.Is(err, ErrInvalidSource):
		logger.Debug("Use path as raw")
		return iox.NewPathWalker(v), nil, nil
	case err != nil:
		return nil, nil, fmt.Errorf("%w: invalid path source", err)
	default:
		return iox.NewReaderWalker(r.AsReadCloser()), r, nil
	}
}

func NewIndexSource(v string) (iox.Source, error) {
	slog.Debug("Use index", slog.String("source", v))
	r, err := NewSource(v)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid index source", err)
	}
	return r, nil
}

func IsStdinSource(s iox.Source) bool {
	_, ok := s.(*iox.StdinSource)
	return ok
}

func (s *Sources) HasStdinConflict() bool {
	var isStdin bool
	for _, x := range []iox.Source{s.Query, s.Index, s.path} {
		if IsStdinSource(x) {
			if isStdin {
				return true
			}
			isStdin = true
		}
	}
	return false
}

func ReadInputFromWalker(w iox.Walker) iter.Seq[*node.Node] {
	return func(yield func(*node.Node) bool) {
		for x := range w.Walk() {
			if !yield(node.FromWalkerEntry(x)) {
				return
			}
		}
	}
}

func ReadInputFromSource(s iox.Source, noValidate bool) iter.Seq[*node.Node] {
	return func(yield func(*node.Node) bool) {
		scanner := bufio.NewScanner(s.AsReadCloser())
		for scanner.Scan() {
			var n node.Node
			b := scanner.Bytes()
			if err := json.Unmarshal(b, &n); err != nil {
				slog.Warn("Failed to unmarshal input", logx.Err(err), logx.String("line", b))
				continue
			}
			if !noValidate {
				if err := n.Validate(); err != nil {
					slog.Warn("Invalid input", logx.Err(err), logx.String("line", b))
					continue
				}
			}
			if !yield(&n) {
				return
			}
		}
	}
}

func (c *Config) fixNodeKeys(n *node.Node) *node.Node {
	if c.RawOutput {
		return n
	}
	r := node.New()
	for k, v := range n.Unwrap() {
		r.Set(tree.KeyFromString(k).Name(), v)
	}
	return r
}

func (c *Config) WriteNode(n *node.Node) {
	b, err := json.Marshal(c.fixNodeKeys(n))
	if err != nil {
		logx.Error(err, "Failed to marshal output")
		return
	}
	_, _ = fmt.Fprintf(c.Stdout, "%s\n", b)
}

func (s *Sources) ReadInput() (iter.Seq[*node.Node], error) {
	switch {
	case s.Path != nil:
		return ReadInputFromWalker(s.Path), nil
	case s.Index != nil:
		return ReadInputFromSource(s.Index, true), nil
	default:
		return nil, ErrInvalidInput
	}
}

func (s *Sources) Close() error {
	var errs []error
	if x := s.Query; x != nil {
		errs = append(errs, x.AsReadCloser().Close())
	}
	if x := s.path; x != nil {
		errs = append(errs, x.AsReadCloser().Close())
	}
	if x := s.Index; x != nil {
		errs = append(errs, x.AsReadCloser().Close())
	}
	return errors.Join(errs...)
}
