package util

import "cmp"

type Comparable interface {
	cmp.Ordered
}

type CompareResult int

const (
	CompareLess CompareResult = iota - 1
	CompareEqual
	CompareGreater
	CompareUnknown
)

func Compare[T Comparable](a, b T) CompareResult { return CompareResult(cmp.Compare(a, b)) }
