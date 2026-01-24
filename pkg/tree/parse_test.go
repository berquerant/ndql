package tree_test

import (
	"testing"
	"time"

	"github.com/berquerant/ndql/pkg/errorx"
	"github.com/berquerant/ndql/pkg/node"
	"github.com/berquerant/ndql/pkg/tree"
	"github.com/stretchr/testify/assert"
)

func TestParseGenResult(t *testing.T) {
	for _, tc := range []struct {
		title string
		b     []byte
		want  []*tree.N
		err   error
	}{
		{
			title: "equalpair",
			b:     []byte(`x=1`),
			want: []*tree.N{
				node.FromMap(map[string]node.Data{
					"x": node.String("1"),
				}),
			},
		},
		{
			title: "equalpairs",
			b: []byte(`x=1,y=2
x=3`),
			want: []*tree.N{
				node.FromMap(map[string]node.Data{
					"x": node.String("1"),
					"y": node.String("2"),
				}),
				node.FromMap(map[string]node.Data{
					"x": node.String("3"),
				}),
			},
		},
		{
			title: "equalpair missingvalue",
			b:     []byte(`x=,y=1`),
			want: []*tree.N{
				node.FromMap(map[string]node.Data{
					"x": node.String(""),
					"y": node.String("1"),
				}),
			},
		},
		{
			title: "empty",
			b:     []byte(``),
			want:  []*tree.N{},
		},
		{
			title: "equalpair many eq",
			b:     []byte(`x=y=2,y=1`),
			want: []*tree.N{
				node.FromMap(map[string]node.Data{
					"x": node.String("y=2"),
					"y": node.String("1"),
				}),
			},
		},
		{
			title: "single node",
			b:     []byte(`{"x":1}`),
			want: []*tree.N{
				node.FromMap(map[string]node.Data{
					"x": node.Int(1),
				}),
			},
		},
		{
			title: "invalid json",
			b:     []byte(`{"x":1`),
			err:   tree.ErrParseGenResult,
		},
		{
			title: "multiple nodes",
			b:     []byte(`[{"x":1},{"x":true,"y":"2s"}]`),
			want: []*tree.N{
				node.FromMap(map[string]node.Data{
					"x": node.Int(1),
				}),
				node.FromMap(map[string]node.Data{
					"x": node.Bool(true),
					"y": node.Duration(2 * time.Second),
				}),
			},
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			got, err := tree.ParseGenResult(tc.b)
			if tc.err != nil {
				assert.ErrorIs(t, err, tc.err)
				return
			}
			if !assert.Nil(t, err, errorx.AsString(err)) {
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}
