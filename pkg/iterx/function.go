package iterx

import (
	"errors"
	"fmt"
)

// Composite function; f(g(x)).
func Pipe[A, B, C any](f func(A) (B, error), g func(B) (C, error)) func(A) (C, error) {
	return func(a A) (C, error) {
		b, err := f(a)
		if err != nil {
			var c C
			return c, err
		}
		return g(b)
	}
}

type Arity int

const (
	Unary Arity = iota
	Variadic
)

var (
	ErrInvalidCall = errors.New("InvalidCall")
)

type Function[T any] interface {
	AsMapFunction() (*MapFunction[T], error)
	AsReduceFunction() (*ReduceFunction[T], error)
	AsFanoutFunction() (*FanoutFunction[T], error)
	AsMultiMapFunction() (*MultiMapFunction[T], error)
	ArgArity() Arity
	RetArity() Arity
	CallAny(v ...T) ([]T, error)
}

type MapFunction[T any] struct {
	f func(T) (T, error)
}

var _ Function[int] = &MapFunction[int]{}

func NewMapFunction[T any](f func(T) (T, error)) *MapFunction[T] { return &MapFunction[T]{f} }

func (f *MapFunction[T]) Call(v T) (T, error) { return f.f(v) }
func (f *MapFunction[T]) CallAny(v ...T) ([]T, error) {
	if len(v) == 1 {
		x, err := f.Call(v[0])
		if err != nil {
			return nil, err
		}
		return []T{x}, nil
	}
	return nil, fmt.Errorf("%w: failed to call MapFunction with %d arguments", ErrInvalidCall, len(v))
}
func (*MapFunction[T]) ArgArity() Arity { return Unary }
func (*MapFunction[T]) RetArity() Arity { return Unary }

type ReduceFunction[T any] struct {
	f func([]T) (T, error)
}

var _ Function[int] = &ReduceFunction[int]{}

func NewReduceFunction[T any](f func([]T) (T, error)) *ReduceFunction[T] {
	return &ReduceFunction[T]{f}
}

func (f *ReduceFunction[T]) Call(v []T) (T, error) { return f.f(v) }
func (f *ReduceFunction[T]) CallAny(v ...T) ([]T, error) {
	x, err := f.Call(v)
	if err != nil {
		return nil, err
	}
	return []T{x}, nil
}
func (*ReduceFunction[T]) ArgArity() Arity { return Variadic }
func (*ReduceFunction[T]) RetArity() Arity { return Unary }

type FanoutFunction[T any] struct {
	f func(T) ([]T, error)
}

var _ Function[int] = &FanoutFunction[int]{}

func NewFanoutFunction[T any](f func(T) ([]T, error)) *FanoutFunction[T] {
	return &FanoutFunction[T]{f}
}

func (f *FanoutFunction[T]) Call(v T) ([]T, error) { return f.f(v) }
func (f *FanoutFunction[T]) CallAny(v ...T) ([]T, error) {
	if len(v) == 1 {
		x, err := f.Call(v[0])
		if err != nil {
			return nil, err
		}
		return x, nil
	}
	return nil, fmt.Errorf("%w: failed to call FanoutFunction with %d arguments", ErrInvalidCall, len(v))
}
func (*FanoutFunction[T]) ArgArity() Arity { return Unary }
func (*FanoutFunction[T]) RetArity() Arity { return Variadic }

type MultiMapFunction[T any] struct {
	f func([]T) ([]T, error)
}

var _ Function[int] = &MultiMapFunction[int]{}

func NewMultiMapFunction[T any](f func([]T) ([]T, error)) *MultiMapFunction[T] {
	return &MultiMapFunction[T]{f}
}

func (f *MultiMapFunction[T]) Call(v []T) ([]T, error)     { return f.f(v) }
func (f *MultiMapFunction[T]) CallAny(v ...T) ([]T, error) { return f.Call(v) }
func (*MultiMapFunction[T]) ArgArity() Arity               { return Variadic }
func (*MultiMapFunction[T]) RetArity() Arity               { return Variadic }

//
// function conversions
//

var (
	ErrInvalidFunctionConversion = errors.New("InvalidFunctionConversion")
)

func newInvalidFunctionConversion(from, to string) error {
	return fmt.Errorf("%w: %s to %s", ErrInvalidFunctionConversion, from, to)
}

func (f *MapFunction[T]) AsMapFunction() (*MapFunction[T], error) {
	return f, nil
}
func (f *ReduceFunction[T]) AsMapFunction() (*MapFunction[T], error) {
	return NewMapFunction(func(v T) (T, error) {
		return f.Call([]T{v})
	}), nil
}
func (f *FanoutFunction[T]) AsMapFunction() (*MapFunction[T], error) {
	return nil, newInvalidFunctionConversion("fanout", "map")
}
func (f *MultiMapFunction[T]) AsMapFunction() (*MapFunction[T], error) {
	return nil, newInvalidFunctionConversion("multimap", "map")
}

func (f *MapFunction[T]) AsReduceFunction() (*ReduceFunction[T], error) {
	return nil, newInvalidFunctionConversion("map", "reduce")
}
func (f *ReduceFunction[T]) AsReduceFunction() (*ReduceFunction[T], error) {
	return f, nil
}
func (f *FanoutFunction[T]) AsReduceFunction() (*ReduceFunction[T], error) {
	return nil, newInvalidFunctionConversion("fanout", "reduce")
}
func (f *MultiMapFunction[T]) AsReduceFunction() (*ReduceFunction[T], error) {
	return nil, newInvalidFunctionConversion("multimap", "reduce")
}

func (f *MapFunction[T]) AsFanoutFunction() (*FanoutFunction[T], error) {
	return NewFanoutFunction(func(v T) ([]T, error) {
		r, err := f.Call(v)
		if err != nil {
			return nil, err
		}
		return []T{r}, nil
	}), nil
}
func (f *ReduceFunction[T]) AsFanoutFunction() (*FanoutFunction[T], error) {
	return NewFanoutFunction(func(v T) ([]T, error) {
		r, err := f.Call([]T{v})
		if err != nil {
			return nil, err
		}
		return []T{r}, nil
	}), nil
}
func (f *FanoutFunction[T]) AsFanoutFunction() (*FanoutFunction[T], error) {
	return f, nil
}
func (f *MultiMapFunction[T]) AsFanoutFunction() (*FanoutFunction[T], error) {
	return NewFanoutFunction(func(v T) ([]T, error) {
		return f.Call([]T{v})
	}), nil
}

func (f *MapFunction[T]) AsMultiMapFunction() (*MultiMapFunction[T], error) {
	return NewMultiMapFunction(func(v []T) ([]T, error) {
		r := []T{}
		for _, x := range v {
			a, err := f.Call(x)
			if err != nil {
				return nil, err
			}
			r = append(r, a)
		}
		return r, nil
	}), nil
}
func (f *ReduceFunction[T]) AsMultiMapFunction() (*MultiMapFunction[T], error) {
	return NewMultiMapFunction(func(v []T) ([]T, error) {
		r, err := f.Call(v)
		if err != nil {
			return nil, err
		}
		return []T{r}, nil
	}), nil
}
func (f *FanoutFunction[T]) AsMultiMapFunction() (*MultiMapFunction[T], error) {
	return NewMultiMapFunction(func(v []T) ([]T, error) {
		r := []T{}
		for _, x := range v {
			a, err := f.Call(x)
			if err != nil {
				return nil, err
			}
			r = append(r, a...)
		}
		return r, nil
	}), nil
}
func (f *MultiMapFunction[T]) AsMultiMapFunction() (*MultiMapFunction[T], error) {
	return f, nil
}

//
// function combination
//

var (
	ErrInvalidFunctionCombination = errors.New("InvalidFunctionCombination")
)

func newInvalidFunctionCombination(a, b string) error {
	return fmt.Errorf("%w: %s and %s", ErrInvalidFunctionCombination, a, b)
}

func CombineFunction[T any](a, b Function[T]) (Function[T], error) {
	switch a := any(a).(type) {
	case *MapFunction[T]:
		switch b := any(b).(type) {
		case *MapFunction[T]:
			// T -> T, T-> T => T -> T
			return NewMapFunction(Pipe(a.Call, b.Call)), nil
		case *ReduceFunction[T]:
			// T -> T, []T -> T => T -> T
			c, err := b.AsMapFunction()
			if err != nil {
				return nil, errors.Join(newInvalidFunctionCombination("map", "reduce"), err)
			}
			return NewMapFunction(Pipe(a.Call, c.Call)), nil
		case *FanoutFunction[T]:
			// T -> T, T -> []T => T -> []T
			return NewFanoutFunction(Pipe(a.Call, b.Call)), nil
		case *MultiMapFunction[T]:
			// T -> T, []T -> []T, T -> []T
			c, err := a.AsFanoutFunction()
			if err != nil {
				return nil, errors.Join(newInvalidFunctionCombination("map", "multimap"), err)
			}
			return NewFanoutFunction(Pipe(c.Call, b.Call)), nil
		}
	case *ReduceFunction[T]:
		switch b := any(b).(type) {
		case *MapFunction[T]:
			// []T -> T, T -> T => []T -> T
			return NewReduceFunction(Pipe(a.Call, b.Call)), nil
		case *ReduceFunction[T]:
			// []T -> T, []T -> T => []T -> T
			c, err := a.AsMultiMapFunction()
			if err != nil {
				return nil, errors.Join(newInvalidFunctionCombination("reduce", "reduce"), err)
			}
			return NewReduceFunction(Pipe(c.Call, b.Call)), nil
		case *FanoutFunction[T]:
			// []T -> T, T -> []T => []T -> []T
			return NewMultiMapFunction(Pipe(a.Call, b.Call)), nil
		case *MultiMapFunction[T]:
			// []T -> T, []T -> []T => []T -> []T
			c, err := a.AsMultiMapFunction()
			if err != nil {
				return nil, errors.Join(newInvalidFunctionCombination("reduce", "multimap"), err)
			}
			return NewMultiMapFunction(Pipe(c.Call, b.Call)), nil
		}
	case *FanoutFunction[T]:
		switch b := any(b).(type) {
		case *MapFunction[T]:
			// T -> []T, T -> T => T -> []T
			c, err := b.AsMultiMapFunction()
			if err != nil {
				return nil, errors.Join(newInvalidFunctionCombination("fanout", "map"), err)
			}
			return NewFanoutFunction(Pipe(a.Call, c.Call)), nil
		case *ReduceFunction[T]:
			// T -> []T, []T -> T => T -> T
			return NewMapFunction(Pipe(a.Call, b.Call)), nil
		case *FanoutFunction[T]:
			// T -> []T, T -> []T => T -> []T
			c, err := b.AsMultiMapFunction()
			if err != nil {
				return nil, errors.Join(newInvalidFunctionCombination("fanout", "fanout"), err)
			}
			return NewFanoutFunction(Pipe(a.Call, c.Call)), nil
		case *MultiMapFunction[T]:
			// T -> []T, []T -> []T => T -> []T
			return NewFanoutFunction(Pipe(a.Call, b.Call)), nil
		}
	case *MultiMapFunction[T]:
		switch b := any(b).(type) {
		case *MapFunction[T]:
			// []T -> []T, T -> T => []T -> []T
			c, err := b.AsMultiMapFunction()
			if err != nil {
				return nil, errors.Join(newInvalidFunctionCombination("multimap", "map"), err)
			}
			return NewMultiMapFunction(Pipe(a.Call, c.Call)), nil
		case *ReduceFunction[T]:
			// []T -> []T, []T -> T => []T -> T
			return NewReduceFunction(Pipe(a.Call, b.Call)), nil
		case *FanoutFunction[T]:
			// []T -> []T, T -> []T => []T -> []T
			c, err := b.AsMultiMapFunction()
			if err != nil {
				return nil, errors.Join(newInvalidFunctionCombination("multimap", "fanout"), err)
			}
			return NewMultiMapFunction(Pipe(a.Call, c.Call)), nil
		case *MultiMapFunction[T]:
			// []T -> []T, []T -> []T => []T -> []T
			return NewMultiMapFunction(Pipe(a.Call, b.Call)), nil
		}
	}

	return nil, fmt.Errorf("%w: unknown functions %#v and %#v", ErrInvalidFunctionCombination, a, b)
}
