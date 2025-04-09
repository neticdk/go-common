package slice

import "github.com/neticdk/go-stdlib/xslices"

// Unfold generates a slice by repeatedly applying a function to an accumulator.
// It includes the accumulator in the result as the first value.
// If the predicate is always false, it returns nil.
// It stops when the predicate returns false.
// Deprecated: Unfold is deprecated - use github.com/neticdk/go-stdlib/xslices.Unfold
func Unfold[T any](acc T, f func(T) T, p func(T) bool, opts ...xslices.UnfoldOption) []T {
	return xslices.Unfold(acc, f, p, opts...)
}

// UnfoldI generates a slice by repeatedly applying a function to an accumulator.
// It includes the accumulator in the result as the first value.
// The length of the result is equal to i + 1.
// If i is negative, it returns nil.
// It stops after i iterations.
// Deprecated: UnfoldI is deprecated - use github.com/neticdk/go-stdlib/xslices.UnfoldI
func UnfoldI[T any](acc T, f func(T) T, n int, opts ...xslices.UnfoldOption) []T {
	return xslices.UnfoldI(acc, f, n, opts...)
}
