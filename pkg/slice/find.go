package slice

// FindFunc returns the first element in the slice that satisfies the
// predicate.
func FindFunc[T any](data []T, f func(T) bool) (T, bool) {
	var zero T
	for _, e := range data {
		if f(e) {
			return e, true
		}
	}

	return zero, false
}

// FindIFunc returns the index of the first element in the slice that satisfies
// the predicate.
func FindIFunc[T any](data []T, f func(T) bool) (int, bool) {
	for i, e := range data {
		if f(e) {
			return i, true
		}
	}

	return -1, false
}
