package iox

import (
	"bufio"
	"io"
	"io/fs"
	"iter"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/berquerant/ndql/pkg/logx"
)

type Walker interface {
	Walk() iter.Seq[*WalkerEntry]
}

type WalkerEntry struct {
	Path    string // rel path
	Size    int64
	Mode    fs.FileMode
	ModTime time.Time
	IsDir   bool
}

type ReaderWalker struct {
	r io.Reader
}

// NewReaderWalker returns a Walker that walks the paths from the io.Reader.
func NewReaderWalker(r io.Reader) *ReaderWalker {
	return &ReaderWalker{
		r: r,
	}
}

var _ Walker = &ReaderWalker{}

func (w *ReaderWalker) Walk() iter.Seq[*WalkerEntry] {
	scanner := bufio.NewScanner(w.r)
	return func(yield func(*WalkerEntry) bool) {
		for scanner.Scan() {
			path := scanner.Text()
			e := &WalkerEntry{
				Path: path,
			}
			if stat, err := os.Stat(path); err == nil {
				e.IsDir = stat.IsDir()
				e.ModTime = stat.ModTime()
				e.Mode = stat.Mode()
				e.Size = stat.Size()
			}
			logx.Trace("ReaderWalker", slog.String("path", e.Path))
			if !yield(e) {
				return
			}
		}
	}
}

type PathWalker struct {
	root string
}

// NewPathWalker returns a Walker that walks the file tree.
func NewPathWalker(root string) *PathWalker {
	return &PathWalker{
		root: root,
	}
}

var _ Walker = &PathWalker{}

func (w PathWalker) Walk() iter.Seq[*WalkerEntry] {
	return func(yield func(*WalkerEntry) bool) {
		_ = filepath.Walk(w.root, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			e := &WalkerEntry{
				Path:    path,
				Size:    info.Size(),
				Mode:    info.Mode(),
				ModTime: info.ModTime(),
				IsDir:   info.IsDir(),
			}
			logx.Trace("PathWalker", slog.String("path", path))
			if !yield(e) {
				return filepath.SkipAll
			}
			return nil
		})
	}
}
