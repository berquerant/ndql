package tree_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/berquerant/ndql/pkg/errorx"
	"github.com/berquerant/ndql/pkg/node"
	"github.com/berquerant/ndql/pkg/tree"
	"github.com/stretchr/testify/assert"
)

type genTemplateTestcase struct {
	title string
	g     tree.GenTemplate
	n     *tree.N
	e     map[string]string
	want  []byte
	err   error
}

func runGenTemplateTest(t *testing.T, cases []*genTemplateTestcase) {
	for _, tc := range cases {
		t.Run(tc.title, func(t *testing.T) {
			for k, v := range tc.e {
				os.Setenv(k, v)
			}
			defer func() {
				for k := range tc.e {
					os.Unsetenv(k)
				}
			}()

			got, err := tc.g.Generate(context.TODO(), tc.n)
			if tc.err != nil {
				assert.ErrorIs(t, err, tc.err)
				return
			}
			if !assert.Nil(t, err, errorx.AsString(err)) {
				return
			}
			assert.Equal(t, string(tc.want), string(got))
		})
	}
}

func TestGenTemplate(t *testing.T) {
	var (
		regexpTestfile = filepath.Join(t.TempDir(), "regexp.txt")
	)
	const (
		regexpTestContent = `a_key1=av1
a_key2=av2
b_key1=bv1
b_key2=bv2`
	)
	if !assert.Nil(t, os.WriteFile(regexpTestfile, []byte(regexpTestContent), 0644)) {
		return
	}

	runGenTemplateTest(t, []*genTemplateTestcase{
		{
			title: "expr const",
			g:     tree.NewExprGenTemplate(`"const"`),
			n:     node.New(),
			want:  []byte(`const`),
		},
		{
			title: "expr env",
			g:     tree.NewExprGenTemplate(`"k1=" + e.KEY`),
			n:     node.New(),
			e: map[string]string{
				"KEY": "VALUE",
			},
			want: []byte(`k1=VALUE`),
		},
		{
			title: "expr get",
			g:     tree.NewExprGenTemplate(`"k1=" + string(n.k1) + ",k2=" + string(n.a.k2)`),
			n: node.FromMap(map[string]node.Data{
				"k1":                              node.Int(1),
				tree.KeyFromName("a.k2").String(): node.Int(2),
			}),
			want: []byte(`k1=1,k2=2`),
		},
		{
			title: "lua const",
			g: tree.NewLuaGenTemplate(`function f(n)
  return "const"
end`, "f"),
			n:    node.New(),
			want: []byte(`const`),
		},
		{
			title: "lua env",
			g: tree.NewLuaGenTemplate(`function f(n)
  return "k1=" .. E.get("k1") .. ",k2=" .. E.get("k2") .. ",k3=" .. E.get("k2", "missing")
end`, "f"),
			n: node.New(),
			e: map[string]string{
				"k1": "v1",
			},
			want: []byte(`k1=v1,k2=,k3=missing`),
		},
		{
			title: "lua get",
			g: tree.NewLuaGenTemplate(`function f(n)
  local k1 = tostring(n.k1 or "")
  local k2 = tostring(n.k2 or "")
  local k3 = tostring(n.k3 or "missing")
  return string.format("k1=%s,k2=%s,k3=%s", k1, k2, k3)
end`, "f"),
			n: node.FromMap(map[string]node.Data{
				"k1": node.Int(1),
			}),
			want: []byte(`k1=1,k2=,k3=missing`),
		},
		{
			title: "regexp const",
			g:     tree.NewRegexpGenTemplate(`a_key1`, `const`),
			n: node.FromMap(map[string]node.Data{
				node.KeyPath: node.String(regexpTestfile),
			}),
			want: []byte(`const`),
		},
		{
			title: "regexp match",
			g:     tree.NewRegexpGenTemplate(`a_key1=(?P<a1>.+)`, `a1=$a1`),
			n: node.FromMap(map[string]node.Data{
				node.KeyPath: node.String(regexpTestfile),
			}),
			want: []byte(`a1=av1`),
		},
		{
			title: "regexp multiple match",
			g:     tree.NewRegexpGenTemplate(`(?P<k>b_[^=]+)=(?P<v>.+)`, `$k=$v`),
			n: node.FromMap(map[string]node.Data{
				node.KeyPath: node.String(regexpTestfile),
			}),
			want: []byte(`b_key1=bv1
b_key2=bv2`),
		},
		{
			title: "string const",
			g:     tree.NewStringGenTemplate(`const`),
			n:     node.New(),
			want:  []byte(`const`),
		},
		{
			title: "string env",
			g:     tree.NewStringGenTemplate(`k1={{ env "KEY" }},k2={{ env "missing" }}`),
			n:     node.New(),
			e: map[string]string{
				"KEY": "VALUE",
			},
			want: []byte(`k1=VALUE,k2=`),
		},
		{
			title: "string envor",
			g:     tree.NewStringGenTemplate(`k1={{ envor "KEY" "default" }},k2={{ envor "missing" "default" }}`),
			n:     node.New(),
			e: map[string]string{
				"KEY": "VALUE",
			},
			want: []byte(`k1=VALUE,k2=default`),
		},
		{
			title: "string get",
			g:     tree.NewStringGenTemplate(`k1={{ .key1 }},k2={{ or .missing "" }},k3={{ .a.key2 }}`),
			n: node.FromMap(map[string]node.Data{
				"key1":                              node.Int(1),
				tree.KeyFromName("a.key2").String(): node.Int(2),
			}),
			want: []byte(`k1=1,k2=,k3=2`),
		},
		{
			title: "string getor",
			g:     tree.NewStringGenTemplate(`k1={{ or .key1 "default" }},k2={{ or .missing "default" }},k3={{ or .a.key2 "default" }}`),
			n: node.FromMap(map[string]node.Data{
				"key1":                              node.Int(1),
				tree.KeyFromName("a.key2").String(): node.Int(2),
			}),
			want: []byte(`k1=1,k2=default,k3=2`),
		},
		{
			title: "shell const",
			g:     tree.NewShellGenTemplate(`echo const`),
			n:     node.New(),
			want:  []byte(`const`),
		},
		{
			title: "shell env",
			g:     tree.NewShellGenTemplate(`echo $KEY`),
			n:     node.New(),
			e: map[string]string{
				"KEY": "VALUE",
			},
			want: []byte(`VALUE`),
		},
		{
			title: "shell get",
			g:     tree.NewShellGenTemplate(`echo "k1=$(get k1),k2=$(get k2),k3=$(get k3),k4=$(get a.k4)"`),
			n: node.FromMap(map[string]node.Data{
				"k1":                              node.Int(1),
				tree.KeyFromName("a.k3").String(): node.Int(2),
				tree.KeyFromName("a.k4").String(): node.Int(3),
				"k4":                              node.Int(4),
			}),
			want: []byte(`k1=1,k2=,k3=2,k4=3`),
		},
		{
			title: "shell get_or",
			g:     tree.NewShellGenTemplate(`echo "k1=$(get_or k1 2),k2=$(get_or k2 2)"`),
			n: node.FromMap(map[string]node.Data{
				"k1": node.Int(1),
			}),
			want: []byte(`k1=1,k2=2`),
		},
	})
}
