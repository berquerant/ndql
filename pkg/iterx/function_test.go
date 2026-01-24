package iterx_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/berquerant/ndql/pkg/iterx"
	"github.com/stretchr/testify/assert"
)

func TestPipe(t *testing.T) {
	var (
		errSample1 = errors.New("Sample1")
		errSample2 = errors.New("Sample2")
	)
	for _, tc := range []struct {
		title string
		f     func(int) (string, error)
		g     func(string) (int, error)
		input int
		want  int
		err   error
	}{
		{
			title: "success",
			f: func(x int) (string, error) {
				return strconv.Itoa(x + 1), nil
			},
			g:     strconv.Atoi,
			input: 10,
			want:  11,
		},
		{
			title: "failed f",
			f: func(_ int) (string, error) {
				return "", errSample1
			},
			g:     strconv.Atoi,
			input: 10,
			err:   errSample1,
		},
		{
			title: "failed g",
			f: func(x int) (string, error) {
				return strconv.Itoa(x + 1), nil
			},
			g: func(_ string) (int, error) {
				return 0, errSample2
			},
			input: 10,
			err:   errSample2,
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			got, err := iterx.Pipe(tc.f, tc.g)(tc.input)
			if tc.err != nil {
				assert.ErrorIs(t, err, tc.err)
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func mapperFunc(v int) (int, error) { return v + 1, nil }
func reducerFunc(v []int) (int, error) {
	var s int
	for _, x := range v {
		s += x
	}
	return s, nil
}
func fanoutFunc(v int) ([]int, error) {
	return []int{v, v}, nil
}
func multimapFunc(v []int) ([]int, error) {
	if len(v) == 0 {
		return nil, nil
	}
	x, y := make([]int, len(v)), make([]int, len(v))
	copy(x, v)
	copy(y, v)
	return append(v, append(x, y...)...), nil
}

func TestCombineFunction(t *testing.T) {
	var (
		mapper   = iterx.NewMapFunction(mapperFunc)
		reducer  = iterx.NewReduceFunction(reducerFunc)
		fanout   = iterx.NewFanoutFunction(fanoutFunc)
		multimap = iterx.NewMultiMapFunction(multimapFunc)
	)

	for i, tc := range []struct {
		f, g  iterx.Function[int]
		input any
		want  any
		err   error
	}{
		{
			f:     mapper,
			g:     mapper,
			input: 1,
			want:  3,
		},
		{
			f:     mapper,
			g:     reducer,
			input: 1,
			want:  2,
		},
		{
			f:     mapper,
			g:     fanout,
			input: 1,
			want:  []int{2, 2},
		},
		{
			f:     mapper,
			g:     multimap,
			input: 1,
			want:  []int{2, 2, 2},
		},
		{
			f:     reducer,
			g:     mapper,
			input: []int{1, 2},
			want:  4,
		},
		{
			f:     reducer,
			g:     reducer,
			input: []int{1, 2},
			want:  3,
		},
		{
			f:     reducer,
			g:     fanout,
			input: []int{1, 2},
			want:  []int{3, 3},
		},
		{
			f:     reducer,
			g:     multimap,
			input: []int{1, 2},
			want:  []int{3, 3, 3},
		},
		{
			f:     fanout,
			g:     mapper,
			input: 1,
			want:  []int{2, 2},
		},
		{
			f:     fanout,
			g:     reducer,
			input: 1,
			want:  2,
		},
		{
			f:     fanout,
			g:     fanout,
			input: 1,
			want:  []int{1, 1, 1, 1},
		},
		{
			f:     fanout,
			g:     multimap,
			input: 1,
			want:  []int{1, 1, 1, 1, 1, 1},
		},
		{
			f:     multimap,
			g:     mapper,
			input: []int{1, 2},
			want:  []int{2, 3, 2, 3, 2, 3},
		},
		{
			f:     multimap,
			g:     reducer,
			input: []int{1, 2},
			want:  9,
		},
		{
			f:     multimap,
			g:     fanout,
			input: []int{1, 2},
			want:  []int{1, 1, 2, 2, 1, 1, 2, 2, 1, 1, 2, 2},
		},
		{
			f:     multimap,
			g:     multimap,
			input: []int{1, 2},
			want:  []int{1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			h, err := iterx.CombineFunction(tc.f, tc.g)
			if tc.err != nil {
				assert.ErrorIs(t, err, tc.err)
				return
			}
			if !assert.Nil(t, err) {
				return
			}

			var got any
			switch h := any(h).(type) {
			case *iterx.MapFunction[int]:
				got, err = h.Call(tc.input.(int))
			case *iterx.ReduceFunction[int]:
				got, err = h.Call(tc.input.([]int))
			case *iterx.FanoutFunction[int]:
				got, err = h.Call(tc.input.(int))
			case *iterx.MultiMapFunction[int]:
				got, err = h.Call(tc.input.([]int))
			default:
				t.Fatal("unknown function type!")
			}
			if tc.err != nil {
				assert.ErrorIs(t, err, tc.err)
				return
			}
			if !assert.Nil(t, err) {
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}
