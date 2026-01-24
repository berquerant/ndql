package util

import (
	"math"

	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Float | constraints.Integer
}

func Add[T Number](a, b T) T      { return a + b }
func Subtract[T Number](a, b T) T { return a - b }
func Multiply[T Number](a, b T) T { return a * b }
func Divide[T Number](a, b T) T   { return a / b }

func Perm2[T any](v ...T) [][]T {
	r := [][]T{}
	for i, x := range v {
		for j := i + 1; j < len(v); j++ {
			y := v[j]
			r = append(r, []T{x, y}, []T{y, x})
		}
	}
	return r
}

func IsFinite(v float64) bool   { return !(IsNaN(v) || IsInfinite(v)) }
func IsNaN(v float64) bool      { return math.IsNaN(v) }
func IsInfinite(v float64) bool { return math.IsInf(v, 1) || math.IsInf(v, -1) }

func ToDegrees(radian float64) float64  { return radian * 180 / math.Pi }
func ToRadians(degrees float64) float64 { return degrees * math.Pi / 180 }
