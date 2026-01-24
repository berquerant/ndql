package tree_test

import (
	"testing"

	"github.com/berquerant/ndql/pkg/node"
	"github.com/berquerant/ndql/pkg/tree"
	"github.com/stretchr/testify/assert"
)

func TestKey(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		for _, tc := range []struct {
			title string
			key   *tree.Key
			n     *tree.N
			want  *tree.N
			ok    bool
		}{
			{
				title: "empty key",
				key:   tree.NewKey("", ""),
				n: node.FromMap(map[string]node.Data{
					"k1": node.Int(1),
				}),
				ok: false,
			},
			{
				title: "missing value",
				key:   tree.NewKey("", "missing"),
				n: node.FromMap(map[string]node.Data{
					"k1": node.Int(1),
				}),
				ok: false,
			},
			{
				title: "get from default table",
				key:   tree.NewKey("", "k1"),
				n: node.FromMap(map[string]node.Data{
					"k1": node.Int(1),
				}),
				want: node.FromMap(map[string]node.Data{
					"k1": node.Int(1),
				}),
				ok: true,
			},
			{
				title: "get from other table",
				key:   tree.NewKey("other", "k1"),
				n: node.FromMap(map[string]node.Data{
					tree.NewKey("other", "k1").String(): node.Int(1),
				}),
				want: node.FromMap(map[string]node.Data{
					tree.NewKey("other", "k1").String(): node.Int(1),
				}),
				ok: true,
			},
			{
				title: "get from other table ignoring default",
				key:   tree.NewKey("other", "k1"),
				n: node.FromMap(map[string]node.Data{
					"k1":                                node.Int(100),
					tree.NewKey("other", "k1").String(): node.Int(1),
				}),
				want: node.FromMap(map[string]node.Data{
					tree.NewKey("other", "k1").String(): node.Int(1),
				}),
				ok: true,
			},
			{
				title: "missing other.k1",
				key:   tree.NewKey("other", "k1"),
				n: node.FromMap(map[string]node.Data{
					"k1": node.Int(1),
				}),
				ok: false,
			},
		} {
			t.Run(tc.title, func(t *testing.T) {
				got, ok := tc.key.Get(tc.n)
				if !assert.Equal(t, tc.ok, ok) {
					return
				}
				assert.Equal(t, tc.want, got)
			})
		}
	})

	t.Run("NameAndString", func(t *testing.T) {
		for _, tc := range []struct {
			name string
			str  string
			key  *tree.Key
		}{
			{
				name: "",
				str:  "",
				key:  tree.NewKey("", ""),
			},
			{
				name: "a",
				str:  "a",
				key:  tree.NewKey("", "a"),
			},
			{
				name: "a.b",
				str:  "a___b",
				key:  tree.NewKey("a", "b"),
			},
			{
				name: "a.b.c",
				str:  "a___b.c",
				key:  tree.NewKey("a", "b.c"),
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				t.Run("KeyFromString", func(t *testing.T) {
					assert.Equal(t, tc.key, tree.KeyFromString(tc.str))
				})
				t.Run("KeyFromName", func(t *testing.T) {
					assert.Equal(t, tc.key, tree.KeyFromName(tc.name))
				})
				t.Run("String", func(t *testing.T) {
					assert.Equal(t, tc.str, tc.key.String())
				})
				t.Run("Name", func(t *testing.T) {
					assert.Equal(t, tc.name, tc.key.Name())
				})
			})
		}
	})
}
