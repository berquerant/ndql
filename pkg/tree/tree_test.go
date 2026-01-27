package tree_test

import (
	"context"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"

	_ "github.com/pingcap/tidb/pkg/types/parser_driver"

	"github.com/berquerant/ndql/pkg/errorx"
	"github.com/berquerant/ndql/pkg/node"
	"github.com/berquerant/ndql/pkg/parse"
	"github.com/berquerant/ndql/pkg/tree"
	"github.com/berquerant/ndql/pkg/util"
	"github.com/stretchr/testify/assert"
)

func newNodes(v []map[string]node.Data) []*tree.N {
	r := make([]*tree.N, len(v))
	for i, x := range v {
		r[i] = node.FromMap(x)
	}
	return r
}

func TestAsIter(t *testing.T) {
	const (
		testEnvKey   = "TestEnvKey1"
		testEnvValue = "TestEnvValue1"
	)
	os.Setenv(testEnvKey, testEnvValue)
	defer os.Unsetenv(testEnvKey)
	var (
		tmplFile       = filepath.Join(t.TempDir(), "tmpl.txt")
		shFile         = filepath.Join(t.TempDir(), "sh.txt")
		luaFile        = filepath.Join(t.TempDir(), "lua.txt")
		exprFile       = filepath.Join(t.TempDir(), "expr.txt")
		exprJSONFile   = filepath.Join(t.TempDir(), "expr_json.txt")
		grepTargetFile = filepath.Join(t.TempDir(), "grep_target.txt")
	)
	const (
		tmplFileContent = `k2={{ .k1 }}`
		shFileContent   = `echo k2=$(get k1)`
		luaFileContent  = `function f(n)
  return "k2=" .. n.k1
end`
		exprFileContent       = `"k2=" + n.k1`
		exprJSONFileContent   = `toJSON([{"k2": n.k1},{"k2": n.k1 + "2"}])`
		grepTargetFileContent = `a_key1=av1
a_key2=av2
b_key1=bv1
b_key2=bv2`
	)
	for _, x := range []struct {
		name    string
		content string
	}{
		{name: tmplFile, content: tmplFileContent},
		{name: shFile, content: shFileContent},
		{name: luaFile, content: luaFileContent},
		{name: exprFile, content: exprFileContent},
		{name: exprJSONFile, content: exprJSONFileContent},
		{name: grepTargetFile, content: grepTargetFileContent},
	} {
		if !assert.Nil(t, os.WriteFile(x.name, []byte(x.content), 0644)) {
			return
		}
	}

	for _, tc := range []struct {
		title string
		data  []*tree.N
		query string
		want  []*tree.N
		err   error
	}{
		{
			title: "where",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(1),
				},
				{
					"k1": node.Int(0),
				},
			}),
			query: `select * where k1 > 0`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(1),
				},
			}),
		},
		{
			title: "from as select as",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(1),
				},
			}),
			query: `select t1.k1 as k10 from (select *) as t1`,
			want: newNodes([]map[string]node.Data{
				{
					tree.KeyFromName("t1.k10").String(): node.Int(1),
				},
			}),
		},
		{
			title: "from as with builtin column",
			data: newNodes([]map[string]node.Data{
				{
					node.KeySize: node.Int(100),
					"k1":         node.Int(1),
				},
				{
					"k1": node.Int(2),
				},
			}),
			query: `select * from (select *) as t1`,
			want: newNodes([]map[string]node.Data{
				{
					node.KeySize:                       node.Int(100),
					tree.KeyFromName("t1.k1").String(): node.Int(1),
				},
				{
					tree.KeyFromName("t1.k1").String(): node.Int(2),
				},
			}),
		},
		{
			title: "from as",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(1),
				},
			}),
			query: `select * from (select k1) as t1`,
			want: newNodes([]map[string]node.Data{
				{
					tree.KeyFromName("t1.k1").String(): node.Int(1),
				},
			}),
		},
		{
			title: "from",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(1),
				},
			}),
			query: `select * from (select k1)`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(1),
				},
			}),
		},
		{
			title: "case multiple cases",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(1),
				},
				{
					"k1": node.Int(0),
				},
				{
					"k1": node.Int(-1),
				},
			}),
			query: `select case
  when k1 = 1 then "one"
  when k1 = 0 then "zero"
  else false
end as k1`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.String("one"),
				},
				{
					"k1": node.String("zero"),
				},
				{
					"k1": node.Bool(false),
				},
			}),
		},
		{
			title: "case value multiple cases",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(1),
				},
				{
					"k1": node.Int(0),
				},
				{
					"k1": node.Int(-1),
				},
			}),
			query: `select case k1
  when 1 then "one"
  when 0 then "zero"
  else false
end as k1`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.String("one"),
				},
				{
					"k1": node.String("zero"),
				},
				{
					"k1": node.Bool(false),
				},
			}),
		},
		{
			title: "case match",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(1),
				},
			}),
			query: `select case when k1 > 0 then 100 end as k1`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(100),
				},
			}),
		},
		{
			title: "case no match no else",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(-1),
				},
			}),
			query: `select case when k1 > 0 then 100 end as k1`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.NewNull(),
				},
			}),
		},
		{
			title: "case no match",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(-1),
				},
			}),
			query: `select case when k1 > 0 then 100 else 1000 end as k1`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(1000),
				},
			}),
		},
		{
			title: "case value match",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(1),
				},
			}),
			query: `select case k1 when 1 then "one" end as k1`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.String("one"),
				},
			}),
		},
		{
			title: "case value no match no else",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(2),
				},
			}),
			query: `select case k1 when 1 then "one" end as k1`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.NewNull(),
				},
			}),
		},
		{
			title: "case value no match",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(2),
				},
			}),
			query: `select case k1 when 1 then "one" else "ELSE" end as k1`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.String("ELSE"),
				},
			}),
		},
		{
			title: "grep multiple results",
			data: newNodes([]map[string]node.Data{
				{
					node.KeyPath: node.String(grepTargetFile),
				},
			}),
			query: `select grep("(?P<k>b_[^=]+)=(?P<v>.+)", "$k=$v")`,
			want: newNodes([]map[string]node.Data{
				{
					node.KeyPath: node.String(grepTargetFile),
					"b_key1":     node.String("bv1"),
				},
				{
					node.KeyPath: node.String(grepTargetFile),
					"b_key2":     node.String("bv2"),
				},
			}),
		},
		{
			title: "grep",
			data: newNodes([]map[string]node.Data{
				{
					node.KeyPath: node.String(grepTargetFile),
				},
			}),
			query: `select grep("a_key1=(?P<a1>.+)", "k1=$a1")`,
			want: newNodes([]map[string]node.Data{
				{
					node.KeyPath: node.String(grepTargetFile),
					"k1":         node.String("av1"),
				},
			}),
		},
		{
			title: "grep no hit",
			data: newNodes([]map[string]node.Data{
				{
					node.KeyPath: node.String(grepTargetFile),
				},
			}),
			query: `select grep("unmatched", "k2=matched")`,
		},
		{
			title: "expr",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.String("v1"),
				},
			}),
			query: `select expr("'k2=' + n.k1")`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.String("v1"),
					"k2": node.String("v1"),
				},
			}),
		},
		{
			title: "expr file",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.String("v1"),
				},
			}),
			query: fmt.Sprintf(`select expr("@%s")`, exprFile),
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.String("v1"),
					"k2": node.String("v1"),
				},
			}),
		},
		{
			title: "expr file multiple results",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.String("v1"),
				},
			}),
			query: fmt.Sprintf(`select expr("@%s")`, exprJSONFile),
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.String("v1"),
					"k2": node.String("v1"),
				},
				{
					"k1": node.String("v1"),
					"k2": node.String("v12"),
				},
			}),
		},
		{
			title: "expr multiple results",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.String("v1"),
				},
			}),
			query: `select expr("'k2=' + n.k1 + '\\nk2=' + n.k1 + '2'")`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.String("v1"),
					"k2": node.String("v1"),
				},
				{
					"k1": node.String("v1"),
					"k2": node.String("v12"),
				},
			}),
		},
		{
			title: "lua",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.String("v1"),
				},
			}),
			query: `select lua("function f(n)
  return 'k2=' .. n.k1
end", "f")`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.String("v1"),
					"k2": node.String("v1"),
				},
			}),
		},
		{
			title: "lua file",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.String("v1"),
				},
			}),
			query: fmt.Sprintf(`select lua("@%s", "f")`, luaFile),
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.String("v1"),
					"k2": node.String("v1"),
				},
			}),
		},
		{
			title: "lua multiple results",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.String("v1"),
				},
			}),
			query: `select lua("function f(n)
  return 'k2=' .. n.k1 .. '\\nk2=' .. n.k1 .. '2'
end", "f")`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.String("v1"),
					"k2": node.String("v1"),
				},
				{
					"k1": node.String("v1"),
					"k2": node.String("v12"),
				},
			}),
		},
		{
			title: "sh",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.String("v1"),
				},
			}),
			query: `select sh("echo k2=$(get k1)")`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.String("v1"),
					"k2": node.String("v1"),
				},
			}),
		},
		{
			title: "sh file",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.String("v1"),
				},
			}),
			query: fmt.Sprintf(`select sh("@%s")`, shFile),
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.String("v1"),
					"k2": node.String("v1"),
				},
			}),
		},
		{
			title: "sh multiple results",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.String("v1"),
				},
			}),
			query: `select sh("echo k2=$(get k1); echo k2=$(get k1)1")`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.String("v1"),
					"k2": node.String("v1"),
				},
				{
					"k1": node.String("v1"),
					"k2": node.String("v11"),
				},
			}),
		},
		{
			title: "tmpl",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.String("v1"),
				},
			}),
			query: `select tmpl("k2={{ .k1 }}")`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.String("v1"),
					"k2": node.String("v1"),
				},
			}),
		},
		{
			title: "tmpl file",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.String("v1"),
				},
			}),
			query: fmt.Sprintf(`select tmpl("@%s")`, tmplFile),
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.String("v1"),
					"k2": node.String("v1"),
				},
			}),
		},
		{
			title: "tmpl multiple results",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.String("v1"),
				},
			}),
			query: `select tmpl("k2={{ .k1 }}
k2={{ .k1 }}2")`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.String("v1"),
					"k2": node.String("v1"),
				},
				{
					"k1": node.String("v1"),
					"k2": node.String("v12"),
				},
			}),
		},
		{
			title: "select all",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.String("v1"),
					"k2": node.String("v2"),
				},
			}),
			query: `select *`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.String("v1"),
					"k2": node.String("v2"),
				},
			}),
		},
		{
			title: "select k1",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.String("v1"),
					"k2": node.String("v2"),
				},
			}),
			query: `select k1`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.String("v1"),
				},
			}),
		},
		{
			title: "select as",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.String("v1"),
					"k2": node.String("v2"),
				},
			}),
			query: `select k1 as a, k2 as b`,
			want: newNodes([]map[string]node.Data{
				{
					"a": node.String("v1"),
					"b": node.String("v2"),
				},
			}),
		},
		{
			title: "select minus",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(1),
				},
			}),
			query: `select -k1 as k1`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(-1),
				},
			}),
		},
		{
			title: "select not",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.Bool(false),
				},
			}),
			query: `select not k1 as k1`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.Bool(true),
				},
			}),
		},
		{
			title: "select bitneg",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(250),
				},
			}),
			query: `select ~k1 as k1`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(-251),
				},
			}),
		},
		{
			title: "select regexp",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.String("foo"),
				},
			}),
			query: `select k1 regexp 'foo', k1 regexp 'bar' as k2, k1 not regexp 'bar' as k3`,
			want: newNodes([]map[string]node.Data{
				{
					`k1 regexp 'foo'`: node.Bool(true),
					"k2":              node.Bool(false),
					"k3":              node.Bool(true),
				},
			}),
		},
		{
			title: "select like",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.String("foo"),
				},
			}),
			query: `select k1 like 'fo.', k1 like 'bar' as k2, k1 not like 'f%' as k3`,
			want: newNodes([]map[string]node.Data{
				{
					`k1 like 'fo.'`: node.Bool(true),
					"k2":            node.Bool(false),
					"k3":            node.Bool(false),
				},
			}),
		},
		{
			title: "parentheses",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.String("foo"),
				},
			}),
			query: `select (k1)`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.String("foo"),
				},
			}),
		},
		{
			title: "is null",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.String("foo"),
					"k2": node.NewNull(),
				},
			}),
			query: `select k1 is null as k1, k2 is null as k2, k1 is not null as k3, k2 is not null as k4`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.Bool(false),
					"k2": node.Bool(true),
					"k3": node.Bool(true),
					"k4": node.Bool(false),
				},
			}),
		},
		{
			title: "is truth",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.Bool(true),
					"k2": node.NewNull(),
					"k3": node.Int(0),
				},
			}),
			query: `select k1 is true as k1, k2 is true as k2, k3 is false as k3, k2 is not true as k4`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.Bool(true),
					"k2": node.Bool(false),
					"k3": node.Bool(true),
					"k4": node.Bool(true),
				},
			}),
		},
		{
			title: "in",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(1),
					"k2": node.Bool(true),
				},
			}),
			query: `select k1 in (1) as k1, k1 in (10) as k2, k2 not in (100) as k3, k1 not in (1) as k4`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.Bool(true),
					"k2": node.Bool(false),
					"k3": node.Bool(true),
					"k4": node.Bool(false),
				},
			}),
		},
		{
			title: "between",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(1),
				},
			}),
			query: `select k1 between 0 and 2 as k1,
k1 between 2 and 3 as k2,
k1 between 3 and 0 as k3`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.Bool(true),
					"k2": node.Bool(false),
					"k3": node.Bool(false),
				},
			}),
		},
		{
			title: "logical and",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.Bool(false),
					"k2": node.Bool(true),
				},
			}),
			query: `select k1 and k2 as k1,
k1 or k2 as k2,
k1 xor k2 as k3`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.Bool(false),
					"k2": node.Bool(true),
					"k3": node.Bool(true),
				},
			}),
		},
		{
			title: "arithmetic",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(10),
					"k2": node.Int(3),
					"k3": node.Int(2),
				},
			}),
			query: `select k1 + k2 as k1,
k1 - k2 as k2,
k1 * k2 as k3,
k1 / k3 as k4,
k1 % k2 as k5`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(13),
					"k2": node.Int(7),
					"k3": node.Int(30),
					"k4": node.Float(5),
					"k5": node.Int(1),
				},
			}),
		},
		{
			title: "shift",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(4),
				},
			}),
			query: `select k1 << 2 as k1,
k1 >> 2 as k2`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(16),
					"k2": node.Int(1),
				},
			}),
		},
		{
			title: "compare",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(1),
					"k2": node.Int(2),
				},
			}),
			query: `select k1 < k2 as k1,
k1 = k2 as k2,
k1 > k2 as k3,
k1 <= k1 as k4,
k2 >= k1 as k5,
k1 <> k2 as k6`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.Bool(true),
					"k2": node.Bool(false),
					"k3": node.Bool(false),
					"k4": node.Bool(true),
					"k5": node.Bool(true),
					"k6": node.Bool(true),
				},
			}),
		},
		{
			title: "cast",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.String("1"),
					"k2": node.String("10s"),
					"k3": node.String("2026-01-02 10:00:00"),
					"k4": node.Float(1.2),
				},
			}),
			query: `select to_int(k1) as k1,
to_float(k1) as k2,
to_bool(k1) as k3,
to_string(k4) as k4,
to_time(k3) as k5,
to_duration(k2) as k6`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(1),
					"k2": node.Float(1),
					"k3": node.Bool(true),
					"k4": node.String("1.2"),
					"k5": node.Time(util.Must(time.Parse(time.DateTime, "2026-01-02 10:00:00"))),
					"k6": node.Duration(time.Second * 10),
				},
			}),
		},
		{
			title: "common",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(1),
				},
			}),
			query: `select least(k1) as k1, greatest(k1) as k2, coalesce(k1) as k3,
least(k1, 0) as k4, greatest(k1, 10) as k5, coalesce(NULL, k1) as k6,
least(k1, true) as k7, greatest(k1, true) as k8, coalesce(NULL, NULL) as k9`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(1),
					"k2": node.Int(1),
					"k3": node.Int(1),
					"k4": node.Int(0),
					"k5": node.Int(10),
					"k6": node.Int(1),
					"k7": node.NewNull(),
					"k8": node.NewNull(),
					"k9": node.NewNull(),
				},
			}),
		},
		{
			title: "control",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(1),
					"k2": node.Int(2),
				},
				{
					"k1": node.Int(0),
					"k2": node.NewNull(),
				},
			}),
			query: `select if(k1 > 0, 100, 10) as k1,
ifnull(k2, 100) as k2,
nullif(k1, 1) as k3`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(100),
					"k2": node.Int(2),
					"k3": node.NewNull(),
				},
				{
					"k1": node.Int(10),
					"k2": node.Int(100),
					"k3": node.Int(0),
				},
			}),
		},
		{
			title: "inverse",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(2),
					"k2": node.Float(0.5),
					"k3": node.String("str"),
				},
			}),
			query: `select inverse(k1) as k1,
inverse(k2) as k2,
inverse(k3) as k3`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.Float(0.5),
					"k2": node.Float(2),
					"k3": node.String("rts"),
				},
			}),
		},
		{
			title: "abs",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(1),
				},
				{
					"k1": node.Int(-1),
				},
			}),
			query: `select abs(k1) as k1`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.Float(1),
				},
				{
					"k1": node.Float(1),
				},
			}),
		},
		{
			title: "sqrt",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(1),
				},
				{
					"k1": node.Int(4),
				},
			}),
			query: `select sqrt(k1) as k1`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.Float(1),
				},
				{
					"k1": node.Float(2),
				},
			}),
		},
		{
			title: "degrees",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(0),
				},
				{
					"k1": node.Float(math.Pi),
				},
			}),
			query: `select degrees(k1) as k1`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.Float(0),
				},
				{
					"k1": node.Float(180),
				},
			}),
		},
		{
			title: "radians",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(0),
				},
				{
					"k1": node.Int(180),
				},
			}),
			query: `select radians(k1) as k1`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.Float(0),
				},
				{
					"k1": node.Float(math.Pi),
				},
			}),
		},
		{
			title: "arctri",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.Float(0.5),
					"k2": node.Int(1),
				},
			}),
			query: `select acos(k2) as k1,
asin(k1) as k2,
atan(k2) as k3,
atan2(k2, 1) as k4`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.Float(0),
					"k2": node.Float(math.Pi / 6),
					"k3": node.Float(math.Pi / 4),
					"k4": node.Float(math.Pi / 4),
				},
			}),
		},
		{
			title: "tri",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.Float(0),
					"k2": node.Float(math.Pi / 6),
				},
			}),
			query: `select cos(k1) as k1,
sin(k2) as k2,
tan(k1) as k3,
cot(k1) as k4`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.Float(1),
					"k2": node.Float(0.5),
					"k3": node.Float(0),
					"k4": node.NewNull(),
				},
			}),
		},
		{
			title: "log",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.Float(math.E),
					"k2": node.Int(2),
					"k3": node.Int(100),
				},
			}),
			query: `select ln(k1) as k1,
log2(k2) as k2,
log10(k3) as k3`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.Float(1),
					"k2": node.Float(1),
					"k3": node.Float(2),
				},
			}),
		},
		{
			title: "exp",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.Float(1),
				},
			}),
			query: `select exp(k1) as k1`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.Float(math.E),
				},
			}),
		},
		{
			title: "rounds",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.Float(0.5),
				},
				{
					"k1": node.Float(1.2),
				},
				{
					"k1": node.Float(1.7),
				},
			}),
			query: `select round(k1) as k1,
ceil(k1) as k2,
floor(k1) as k3`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.Float(1),
					"k2": node.Float(1),
					"k3": node.Float(0),
				},
				{
					"k1": node.Float(1),
					"k2": node.Float(2),
					"k3": node.Float(1),
				},
				{
					"k1": node.Float(2),
					"k2": node.Float(2),
					"k3": node.Float(1),
				},
			}),
		},
		{
			title: "pow",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.Float(2),
					"k2": node.Float(1),
				},
				{
					"k1": node.Float(2),
					"k2": node.Float(3),
				},
				{
					"k1": node.Float(2),
					"k2": node.Float(-1),
				},
			}),
			query: `select pow(k1, k2) as k1`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.Float(2),
				},
				{
					"k1": node.Float(8),
				},
				{
					"k1": node.Float(0.5),
				},
			}),
		},
		{
			title: "const",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(1),
				},
			}),
			query: `select e() as k1, pi() as k2`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.Float(math.E),
					"k2": node.Float(math.Pi),
				},
			}),
		},
		{
			title: "rand",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(1),
				},
			}),
			query: `select (rand() < 10) as k1`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.Bool(true),
				},
			}),
		},
		{
			title: "len and size",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.String("Hello, 世界"),
				},
			}),
			query: `select len(k1) as k1, size(k1) as k2`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(9),
					"k2": node.Int(13),
				},
			}),
		},
		{
			title: "regexp",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.String("Hello, 世界"),
				},
				{
					"k1": node.String("unmatched"),
				},
			}),
			query: `select regexp_like(k1, "llo") as k1,
regexp_instr(k1, "llo") as k2,
regexp_substr(k1, "ll.") as k3,
regexp_replace(k1, "llo", "y") as k4,
regexp_count(k1, "llo") as k5`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.Bool(true),
					"k2": node.Int(3),
					"k3": node.String("llo"),
					"k4": node.String("Hey, 世界"),
					"k5": node.Int(1),
				},
				{
					"k1": node.Bool(false),
					"k2": node.Int(0),
					"k3": node.String(""),
					"k4": node.String("unmatched"),
					"k5": node.Int(0),
				},
			}),
		},
		{
			title: "format",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(10),
				},
			}),
			query: `select format("i=%d", k1) as k1`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.String("i=10"),
				},
			}),
		},
		{
			title: "lower and upper",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.String("aBc"),
				},
			}),
			query: `select lower(k1) as k1, upper(k1) as k2`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.String("abc"),
					"k2": node.String("ABC"),
				},
			}),
		},
		{
			title: "sha2",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.String("str"),
				},
			}),
			query: `select sha2(k1) as k1`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.String("8c25cb3686462e9a86d2883c5688a22fe738b0bbc85f458d2d2b5f3f667c6d5a"),
				},
			}),
		},
		{
			title: "concat_ws",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.String("str"),
					"k2": node.String("ing"),
				},
			}),
			query: `select concat_ws(".", k1, k2) as k1`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.String("str.ing"),
				},
			}),
		},
		{
			title: "instr",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.String("str"),
				},
			}),
			query: `select instr(k1, "t") as k1`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(2),
				},
			}),
		},
		{
			title: "substr",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.String("str.i.ng"),
				},
			}),
			query: `select substr(k1, 3, 3) as k1, substr_index(k1, ".", 2) as k2`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.String("r.i"),
					"k2": node.String("str.i"),
				},
			}),
		},
		{
			title: "replace",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.String("str"),
				},
			}),
			query: `select replace(k1, "t", "s") as k1`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.String("ssr"),
				},
			}),
		},
		{
			title: "trim",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.String("  str "),
				},
			}),
			query: `select trim(k1) as k1`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.String("str"),
				},
			}),
		},
		{
			title: "strtotime",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.String("2026-01-02T10:00:00"),
				},
			}),
			query: `select strtotime(k1, "2006-01-02T15:04:05") as k1`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.Time(util.Must(time.Parse(time.DateTime, "2026-01-02 10:00:00"))),
				},
			}),
		},
		{
			title: "timeformat",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.Time(util.Must(time.Parse(time.DateTime, "2026-01-02 03:04:05"))),
					"k2": node.String("2006-01-02T15:04:05"),
				},
			}),
			query: `select timeformat(k1, k2) as k1`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.String("2026-01-02T03:04:05"),
				},
			}),
		},
		{
			title: "time extraction",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.Time(util.Must(time.Parse(time.DateTime, "2026-01-02 03:04:05"))),
				},
			}),
			query: `select year(k1) as k1,
month(k1) as k2,
day(k1) as k3,
hour(k1) as k4,
minute(k1) as k5,
second(k1) as k6`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(2026),
					"k2": node.Int(1),
					"k3": node.Int(2),
					"k4": node.Int(3),
					"k5": node.Int(4),
					"k6": node.Int(5),
				},
			}),
		},
		{
			title: "dayof",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.Time(util.Must(time.Parse(time.DateTime, "2026-02-01 03:04:05"))),
				},
			}),
			query: `select dayofweek(k1) as k1, dayofyear(k1) as k2`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(1),
					"k2": node.Int(32),
				},
			}),
		},
		{
			title: "newtime",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(2026),
				},
			}),
			query: `select newtime(k1, 1, 2) as k1`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.Time(util.Must(time.Parse(time.DateTime, "2026-01-02 00:00:00")).UTC()),
				},
			}),
		},
		{
			title: "sleep",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.Duration(100 * time.Millisecond),
				},
			}),
			query: `select sleep(k1) as k1`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.Int(0),
				},
			}),
		},
		{
			title: "env",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.String(testEnvKey),
				},
				{
					"k1": node.String(testEnvKey + "_missing"),
				},
			}),
			query: `select envor(k1, 1) as k1, env(k1) as k2`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.String(testEnvValue),
					"k2": node.String(testEnvValue),
				},
				{
					"k1": node.Int(1),
					"k2": node.NewNull(),
				},
			}),
		},
		{
			title: "path",
			data: newNodes([]map[string]node.Data{
				{
					"k1": node.String("/root/sub/some.zst"),
				},
			}),
			query: `select dir(k1) as k1,
basename(k1) as k2,
extension(k1) as k3,
abspath(k1) as k4,
relpath(k1, "/root") as k5`,
			want: newNodes([]map[string]node.Data{
				{
					"k1": node.String("/root/sub"),
					"k2": node.String("some.zst"),
					"k3": node.String(".zst"),
					"k4": node.String("/root/sub/some.zst"),
					"k5": node.String("sub/some.zst"),
				},
			}),
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			r, err := parse.NewSQLParser().Parse(tc.query)
			if !assert.Nil(t, err, "query syntax: %s", errorx.AsString(err)) {
				return
			}
			it, err := tree.AsIter(context.TODO(), slices.Values(tc.data), r.Nodes[0])
			if tc.err != nil {
				assert.ErrorIs(t, err, tc.err, errorx.AsString(err))
				return
			}
			if !assert.Nil(t, err, errorx.AsString(err)) {
				return
			}
			got := slices.Collect(it)
			assert.Equal(t, tc.want, got)
		})
	}
}
