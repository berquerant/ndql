package iterx

import "iter"

// ClonableIter shares the source of the iterator.
type ClonableIter[T any] struct {
	cache *clonableIterCache[T]
	index int
}

type clonableIterCache[T any] struct {
	next  func() (T, bool)
	stop  func()
	cache []T
}

func (c *clonableIterCache[T]) read(index int) (T, bool) {
	if index < len(c.cache) {
		return c.cache[index], true
	}
	v, ok := c.next()
	if !ok {
		var t T
		return t, false
	}
	c.cache = append(c.cache, v)
	return v, true
}

func NewClonableIter[T any](it Iter[T]) *ClonableIter[T] {
	next, stop := iter.Pull(it)
	return &ClonableIter[T]{
		cache: &clonableIterCache[T]{
			next:  next,
			stop:  stop,
			cache: []T{},
		},
		index: 0,
	}
}

// Clone returns a cloned iterator that iterates from the beginning of the source.
func (it *ClonableIter[T]) Clone() *ClonableIter[T] {
	return &ClonableIter[T]{
		cache: it.cache,
		index: 0,
	}
}

func (it *ClonableIter[T]) read() (T, bool) {
	defer func() {
		it.index++
	}()
	return it.cache.read(it.index)
}

func (it *ClonableIter[T]) Values() Iter[T] {
	return func(yield func(T) bool) {
		for {
			v, ok := it.read()
			if !ok {
				return
			}
			if !yield(v) {
				return
			}
		}
	}
}

func (it *ClonableIter[T]) Close() { it.cache.stop() }
