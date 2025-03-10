package slice

// Unfold generates a slice by repeatedly applying a function to an accumulator.
// It includes the accumulator in the result as the first value.
// If the predicate is always false, it returns nil.
// It stops when the predicate returns false.
func Unfold[T any](acc T, f func(T) T, p func(T) bool) []T {
	var res []T
	for a := acc; p(a); a = f(a) {
		res = append(res, a)
	}

	return res
}

// UnfoldI generates a slice by repeatedly applying a function to an accumulator.
// It includes the accumulator in the result as the first value.
// The length of the result is equal to i.
// If i is negative, it returns nil.
// It stops after i iterations.
func UnfoldI[T any](acc T, f func(T) T, i int) []T {
	if i < 0 {
		return nil
	}

	res := make([]T, 0, i)
	for range i {
		res = append(res, acc)
		acc = f(acc)
	}

	return res
}
