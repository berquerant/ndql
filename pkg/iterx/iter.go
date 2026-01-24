package iterx

import (
	"errors"
	"iter"
	"log/slog"

	"github.com/berquerant/ndql/pkg/logx"
)

type Iter[T any] = iter.Seq[T]

func Filter[T any](it Iter[T], f func(T) bool) Iter[T] {
	return func(yield func(T) bool) {
		for x := range it {
			if !f(x) {
				continue
			}
			if !yield(x) {
				return
			}
		}
	}
}

func Always[T any](_ T) bool { return true }
func Never[T any](_ T) bool  { return false }

var (
	ErrIgnore = errors.New("Ignore")
)

func Map[T, U any](it Iter[T], f func(T) (U, error)) Iter[U] {
	return func(yield func(U) bool) {
		for x := range it {
			v, err := f(x)
			switch {
			case errors.Is(err, ErrIgnore):
				continue
			case err != nil:
				slog.Warn("iterx.Map", logx.Value("in", x), logx.Err(err))
				continue
			default:
				if !yield(v) {
					return
				}
			}
		}
	}
}

func Identity[T any](x T) (T, error) { return x, nil }
