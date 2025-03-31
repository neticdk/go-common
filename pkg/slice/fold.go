package slice

import "github.com/neticdk/go-stdlib/xslices"

// Fold applies a function to each element of the slice.
// storing the result in an accumulator.
// It applies the function from left to right.
// Deprecated: Fold is deprecated - use github.com/neticdk/go-stdlib/xslices.Fold
func Fold[T, S any](acc T, data []S, f func(T, S) T) T {
	return xslices.Fold(acc, data, f)
}

// FoldR applies a function to each element of the slice.
// storing the result in an accumulator.
// It applies the function from right to left.
// Deprecated: FoldR is deprecated - use github.com/neticdk/go-stdlib/xslices.FoldR
func FoldR[T, S any](acc T, data []S, f func(T, S) T) T {
	return xslices.FoldR(acc, data, f)
}
