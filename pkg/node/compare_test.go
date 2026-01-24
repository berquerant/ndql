package node_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/berquerant/ndql/pkg/node"
	"github.com/berquerant/ndql/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestCompare(t *testing.T) {
	for _, tc := range []struct {
		left, right node.Data
		want        node.CompareResult
	}{
		{
			left:  node.NewNull(),
			right: node.NewNull(),
			want:  node.CmpEqual,
		},
		{
			left:  node.NewNull(),
			right: node.Bool(false),
			want:  node.CmpUnknown,
		},
		{
			left:  node.Bool(false),
			right: node.Bool(true),
			want:  node.CmpLess,
		},
		{
			left:  node.Bool(true),
			right: node.Bool(false),
			want:  node.CmpGreater,
		},
		{
			left:  node.Bool(true),
			right: node.Bool(true),
			want:  node.CmpEqual,
		},
		{
			left:  node.Float(2),
			right: node.Bool(false),
			want:  node.CmpUnknown,
		},
		{
			left:  node.Float(2),
			right: node.Float(3),
			want:  node.CmpLess,
		},
		{
			left:  node.Float(3),
			right: node.Float(3),
			want:  node.CmpEqual,
		},
		{
			left:  node.Float(2),
			right: node.Float(1),
			want:  node.CmpGreater,
		},
		{
			left:  node.Int(2),
			right: node.Float(3),
			want:  node.CmpLess,
		},
		{
			left:  node.Float(3),
			right: node.Int(2),
			want:  node.CmpGreater,
		},
		{
			left:  node.Int(3),
			right: node.Int(4),
			want:  node.CmpLess,
		},
		{
			left:  node.Int(4),
			right: node.Int(4),
			want:  node.CmpEqual,
		},
		{
			left:  node.Int(5),
			right: node.Int(4),
			want:  node.CmpGreater,
		},
		{
			left:  node.Int(3),
			right: node.String(""),
			want:  node.CmpUnknown,
		},
		{
			left:  node.String("a"),
			right: node.String("b"),
			want:  node.CmpLess,
		},
		{
			left:  node.String("b"),
			right: node.String("b"),
			want:  node.CmpEqual,
		},
		{
			left:  node.String("c"),
			right: node.String("b"),
			want:  node.CmpGreater,
		},
		{
			left:  node.String("a"),
			right: node.Time(util.Must(time.Parse(time.DateTime, "2026-01-02 10:00:00"))),
			want:  node.CmpUnknown,
		},
		{
			left:  node.Time(util.Must(time.Parse(time.DateTime, "2026-01-02 10:00:00"))),
			right: node.Time(util.Must(time.Parse(time.DateTime, "2026-01-02 11:00:00"))),
			want:  node.CmpLess,
		},
		{
			left:  node.Time(util.Must(time.Parse(time.DateTime, "2026-01-02 11:00:00"))),
			right: node.Time(util.Must(time.Parse(time.DateTime, "2026-01-02 11:00:00"))),
			want:  node.CmpEqual,
		},
		{
			left:  node.Time(util.Must(time.Parse(time.DateTime, "2026-01-02 12:00:00"))),
			right: node.Time(util.Must(time.Parse(time.DateTime, "2026-01-02 11:00:00"))),
			want:  node.CmpGreater,
		},
		{
			left:  node.Time(util.Must(time.Parse(time.DateTime, "2026-01-02 10:00:00"))),
			right: node.Duration(time.Minute),
			want:  node.CmpUnknown,
		},
		{
			left:  node.Duration(time.Minute),
			right: node.Duration(time.Minute * 2),
			want:  node.CmpLess,
		},
		{
			left:  node.Duration(time.Minute),
			right: node.Duration(time.Minute),
			want:  node.CmpEqual,
		},
		{
			left:  node.Duration(time.Minute * 3),
			right: node.Duration(time.Minute * 2),
			want:  node.CmpGreater,
		},
	} {
		title := fmt.Sprintf("%s_%s", tc.left.Display(), tc.right.Display())
		t.Run(title, func(t *testing.T) {
			got := tc.left.AsOp().Compare(tc.right.AsOp())
			assert.Equal(t, tc.want, got)
		})
	}
}
