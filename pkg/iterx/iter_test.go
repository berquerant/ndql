package iterx_test

import (
	"slices"
	"strconv"
	"testing"

	"github.com/berquerant/ndql/pkg/iterx"
	"github.com/stretchr/testify/assert"
)

type Iter = iterx.Iter[int]

type testcase struct {
	title  string
	filter func(int) bool
	mapper func(int) (int, error)
	src    []int
	want   []int
}

func runTestcases(t *testing.T, f func(Iter, *testcase) Iter, cases []*testcase) {
	for _, tc := range cases {
		t.Run(tc.title, func(t *testing.T) {
			got := slices.Collect(f(slices.Values(tc.src), tc))
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestMapper(t *testing.T) {
	t.Run("toInt", func(t *testing.T) {
		var (
			want = []int{1, 2, 3}
			src  = []string{"1", "2", "3"}
		)
		got := slices.Collect(iterx.Map(slices.Values(src), strconv.Atoi))
		assert.Equal(t, want, got)
	})

	runTestcases(t, func(it Iter, tc *testcase) Iter {
		return iterx.Map(it, tc.mapper)
	}, []*testcase{
		{
			title:  "empty",
			mapper: func(x int) (int, error) { return x, nil },
			src:    []int{},
		},
		{
			title:  "identity",
			mapper: func(x int) (int, error) { return x, nil },
			src:    []int{1, 2, 3},
			want:   []int{1, 2, 3},
		},
		{
			title:  "double",
			mapper: func(x int) (int, error) { return x * 2, nil },
			src:    []int{1, 2, 3},
			want:   []int{2, 4, 6},
		},
		{
			title: "ignore",
			mapper: func(x int) (int, error) {
				if x%2 == 0 {
					return 0, iterx.ErrIgnore
				}
				return x * 2, nil
			},
			src:  []int{1, 2, 3},
			want: []int{2, 6},
		},
	})
}

func TestFilter(t *testing.T) {
	runTestcases(t, func(it Iter, tc *testcase) Iter {
		return iterx.Filter(it, tc.filter)
	}, []*testcase{
		{
			title:  "empty",
			filter: func(_ int) bool { return true },
			src:    []int{},
		},
		{
			title:  "identity",
			filter: func(_ int) bool { return true },
			src:    []int{1, 2, 3},
			want:   []int{1, 2, 3},
		},
		{
			title:  "even only",
			filter: func(x int) bool { return x%2 == 0 },
			src:    []int{1, 2, 3},
			want:   []int{2},
		},
	})
}
