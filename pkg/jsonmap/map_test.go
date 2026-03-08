package jsonmap_test

import (
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"testing"

	"github.com/berquerant/ndql/pkg/jsonmap"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	m := jsonmap.NewMap[string]()
	for _, tc := range []*testcase{
		{
			name: "empty",
			key:  "key",
		},
		{
			name:  "top",
			set:   true,
			key:   "k1",
			value: "v1",
		},
		{
			name:     "top",
			key:      "k1",
			value:    "v1",
			exist:    true,
			hasValue: true,
		},
		{
			name:  "depth 2",
			set:   true,
			key:   "k2.k21",
			value: "v2.21",
		},
		{
			name:     "top",
			key:      "k2",
			exist:    true,
			children: []string{"k21"},
		},
		{
			name:     "depth 2",
			key:      "k2.k21",
			value:    "v2.21",
			exist:    true,
			hasValue: true,
		},
		{
			name:  "depth 2 overwrite",
			set:   true,
			key:   "k2.k21",
			value: "v2.21!",
		},
		{
			name:     "depth 2 overwrite",
			key:      "k2.k21",
			value:    "v2.21!",
			exist:    true,
			hasValue: true,
		},
		{
			name:  "top add",
			set:   true,
			key:   "k2",
			value: "v2",
		},
		{
			name:     "top add",
			key:      "k2",
			value:    "v2",
			exist:    true,
			hasValue: true,
			children: []string{"k21"},
		},
	} {
		tc.test(t, m)
	}

	t.Run("Marshal", func(t *testing.T) {
		b, err := json.Marshal(m)
		assert.Nil(t, err)
		var x jsonmap.Map[string]
		assert.Nil(t, json.Unmarshal(b, &x))
		assert.Equal(t, m, &x)
	})
}

type testcase struct {
	name     string
	set      bool
	key      string
	value    string
	hasValue bool
	exist    bool
	children []string
}

func (tc *testcase) title() string {
	if tc.set {
		return fmt.Sprintf("%s set %s", tc.name, tc.key)
	}
	return fmt.Sprintf("%s get %s", tc.name, tc.key)
}

func (tc *testcase) test(t *testing.T, m *jsonmap.Map[string]) {
	t.Run(tc.title(), func(t *testing.T) {
		b, _ := json.Marshal(m)
		t.Logf("m=%s", b)

		if tc.set {
			m.Set(tc.key, tc.value)
			return
		}

		got, ok := m.Get(tc.key)
		if !tc.exist {
			assert.False(t, ok, "want not exist")
			return
		}
		if !assert.True(t, ok, "want exist") {
			return
		}
		gotValue, gotHasValue := got.Value()
		assert.Equal(t, tc.hasValue, gotHasValue, "hasValue")
		assert.Equal(t, tc.value, gotValue, "value")

		slices.Sort(tc.children)
		gotChildren := slices.Collect(maps.Keys(got.Children()))
		slices.Sort(gotChildren)
		assert.Equal(t, tc.children, gotChildren, "children keys")
	})
}
