package gopkg

import (
	"bytes"
	"encoding/json"
	"maps"
	"slices"

	"github.com/berquerant/ndql/pkg/jsonmap"
)

type DocumentMap = jsonmap.Map[*Document]

type DocumentSet struct {
	m *DocumentMap
}

func NewDocumentSet(doc ...*Document) *DocumentSet {
	d := jsonmap.NewMap[*Document]()
	for _, x := range doc {
		d.Set(x.Path, x)
	}
	return &DocumentSet{
		m: d,
	}
}

func (s *DocumentSet) MarshalJSON() ([]byte, error) { return json.Marshal(s.m) }
func (s *DocumentSet) UnmarshalJSON(b []byte) error {
	var m DocumentMap
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	s.m = &m
	return nil
}

func (s *DocumentSet) Keys() []string {
	c := s.m.Children()
	keys := slices.Collect(maps.Keys(c))
	slices.Sort(keys)
	return keys
}

func (s *DocumentSet) Get(key string) (string, bool) {
	m, ok := s.m.Get(key)
	if !ok {
		return "", false
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
		return b.String(), true
	}
	keys := slices.Collect(maps.Keys(c))
	slices.Sort(keys)
	w("")
	w("# Children")
	w("")
	for _, k := range keys {
		w("- " + k)
	}
	return b.String(), true
}
