package gopkg

import (
	"bytes"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
)

func (s *DocumentSet) IntoFiles(root string) error { return s.intoFiles(s.m, root) }

func (s *DocumentSet) intoFiles(m *DocumentMap, root string) error {
	if err := os.MkdirAll(root, 0755); err != nil {
		return err
	}
	var (
		b bytes.Buffer
		w = func(s string) { b.WriteString(s + "\n") }
	)
	if v, ok := m.Value(); ok {
		w(v.String())
	}

	c := m.Children()
	if len(c) == 0 {
		return nil
	}
	keys := slices.Collect(maps.Keys(c))
	slices.Sort(keys)
	w("")
	w("# Children")
	w("")
	for _, k := range keys {
		w(fmt.Sprintf("- [%s](./%s)", k, k))
	}
	if err := os.WriteFile(root, b.Bytes(), 0644); err != nil {
		return err
	}

	for _, k := range keys {
		v := c[k]
		if err := s.intoFiles(v, filepath.Join(root, k)); err != nil {
			return err
		}
	}
	return nil
}
