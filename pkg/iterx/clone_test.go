package iterx_test

import (
	"iter"
	"slices"
	"testing"

	"github.com/berquerant/ndql/pkg/iterx"
	"github.com/stretchr/testify/assert"
)

func TestClonableIter(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		cit := iterx.NewClonableIter(slices.Values([]int{}))
		defer cit.Close()
		cit2 := cit.Clone()
		assert.Empty(t, slices.Collect(cit.Values()))
		assert.Empty(t, slices.Collect(cit2.Values()))
	})
	t.Run("values", func(t *testing.T) {
		cit := iterx.NewClonableIter(slices.Values([]int{1, 2, 3}))
		defer cit.Close()
		cit2 := cit.Clone()
		assert.Equal(t, []int{1, 2, 3}, slices.Collect(cit.Values()))
		assert.Equal(t, []int{1, 2, 3}, slices.Collect(cit2.Values()))
	})
	t.Run("clone after collect", func(t *testing.T) {
		cit := iterx.NewClonableIter(slices.Values([]int{1, 2, 3}))
		defer cit.Close()
		assert.Equal(t, []int{1, 2, 3}, slices.Collect(cit.Values()))
		cit2 := cit.Clone()
		assert.Equal(t, []int{1, 2, 3}, slices.Collect(cit2.Values()))
	})
	t.Run("origin stop early", func(t *testing.T) {
		cit := iterx.NewClonableIter(slices.Values([]int{1, 2, 3}))
		defer cit.Close()
		cit2 := cit.Clone()

		_, stop := iter.Pull(cit.Values())
		stop()
		assert.Equal(t, []int{1, 2, 3}, slices.Collect(cit2.Values()))
	})
	t.Run("origin stop early2", func(t *testing.T) {
		cit := iterx.NewClonableIter(slices.Values([]int{1, 2, 3}))
		defer cit.Close()
		cit2 := cit.Clone()

		next, stop := iter.Pull(cit.Values())
		v, ok := next()
		assert.True(t, ok)
		assert.Equal(t, 1, v)
		stop()
		assert.Equal(t, []int{1, 2, 3}, slices.Collect(cit2.Values()))
	})
	t.Run("origin and clone stop early", func(t *testing.T) {
		cit := iterx.NewClonableIter(slices.Values([]int{1, 2, 3}))
		defer cit.Close()
		cit2 := cit.Clone()

		next, stop := iter.Pull(cit.Values())
		next2, stop2 := iter.Pull(cit2.Values())
		v, ok := next()
		assert.True(t, ok)
		assert.Equal(t, 1, v)
		v, ok = next2()
		assert.True(t, ok)
		assert.Equal(t, 1, v)
		stop()
		v, ok = next2()
		assert.True(t, ok)
		assert.Equal(t, 2, v)
		stop2()

		cit3 := cit.Clone()
		assert.Equal(t, []int{1, 2, 3}, slices.Collect(cit3.Values()))
	})
}
