package util

import "time"

func NewDuration[T Number](n T, d time.Duration) time.Duration {
	f := float64(n)
	return time.Duration(int(f * float64(d)))
}
