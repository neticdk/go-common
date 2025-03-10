package slice

// Fold applies a function to each element of the slice,
// storing the result in an accumulator.
// It applies the function from left to right.
func Fold[T any](acc T, data []T, f func(T, T) T) T {
	res := acc
	for _, e := range data {
		res = f(res, e)
	}

	return res
}

// Fold applies a function to each element of the slice,
// storing the result in an accumulator.
// It applies the function from right to left.
func FoldR[T any](acc T, data []T, f func(T, T) T) T {
	res := acc
	for i := len(data) - 1; i >= 0; i-- {
		res = f(res, data[i])
	}

	return res
}
