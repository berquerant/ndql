package iox

import (
	"bytes"
	"io"
	"os"
	"strings"
)

type Source interface {
	ReadAll() ([]byte, error)
	AsReadCloser() io.ReadCloser
}

type psuedoReadCloser struct {
	io.Reader
}

func (psuedoReadCloser) Close() error { return nil }

var _ io.ReadCloser = &psuedoReadCloser{}

func asReadCloser(r io.Reader) *psuedoReadCloser { return &psuedoReadCloser{r} }

type ReaderSource struct {
	io.ReadCloser
}

func NewReaderSourceFromReader(r io.Reader) *ReaderSource { return &ReaderSource{asReadCloser(r)} }
func NewReaderSource(r io.ReadCloser) *ReaderSource       { return &ReaderSource{r} }

var _ Source = &ReaderSource{}

func (r *ReaderSource) ReadAll() ([]byte, error)    { return io.ReadAll(r) }
func (r *ReaderSource) AsReadCloser() io.ReadCloser { return r }

type StdinSource struct {
	*ReaderSource
}

func NewStdinSource() *StdinSource {
	return &StdinSource{NewReaderSourceFromReader(os.Stdin)}
}

var _ Source = &StdinSource{}

type FileSource struct {
	*ReaderSource
}

func NewFileSource(filename string) (*FileSource, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return &FileSource{NewReaderSource(f)}, nil
}

var _ Source = &FileSource{}

type StringSource struct {
	*ReaderSource
}

func NewStringSource(s string) *StringSource {
	return &StringSource{NewReaderSourceFromReader(bytes.NewBufferString(s))}
}

var _ Source = &StringSource{}

func NewFileOrStringSource(s string) (Source, error) {
	if strings.HasPrefix(s, "@") {
		return NewFileSource(s[1:])
	}
	return NewStringSource(s), nil
}
