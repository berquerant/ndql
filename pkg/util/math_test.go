package util_test

import (
	"testing"

	"github.com/berquerant/ndql/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestPerm2(t *testing.T) {
	for _, tc := range []struct {
		title string
		v     []int
		want  [][]int
	}{
		{
			title: "empty",
			v:     []int{},
			want:  [][]int{},
		},
		{
			title: "1",
			v:     []int{1},
			want:  [][]int{},
		},
		{
			title: "2",
			v:     []int{1, 2},
			want: [][]int{
				{1, 2},
				{2, 1},
			},
		},
		{
			title: "3",
			v:     []int{1, 2, 3},
			want: [][]int{
				{1, 2},
				{2, 1},
				{1, 3},
				{3, 1},
				{2, 3},
				{3, 2},
			},
		},
		{
			title: "4",
			v:     []int{1, 2, 3, 4},
			want: [][]int{
				{1, 2},
				{2, 1},
				{1, 3},
				{3, 1},
				{1, 4},
				{4, 1},
				{2, 3},
				{3, 2},
				{2, 4},
				{4, 2},
				{3, 4},
				{4, 3},
			},
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			assert.Equal(t, tc.want, util.Perm2(tc.v...))
		})
	}
}
