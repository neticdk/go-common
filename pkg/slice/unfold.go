package slice

type UnfoldConfig struct {
	Max  int
	Step int
}

type UnfoldOption func(*UnfoldConfig)

// WithStep sets the step of the unfold.
// This does not include the first step.
// If max is 5, the length of the result is <= 6.
func WithMax(max int) UnfoldOption {
	return func(c *UnfoldConfig) {
		c.Max = max
	}
}

// Unfold generates a slice by repeatedly applying a function to an accumulator.
// It includes the accumulator in the result as the first value.
// If the predicate is always false, it returns nil.
// It stops when the predicate returns false.
func Unfold[T any](acc T, f func(T) T, p func(T) bool, opts ...UnfoldOption) []T {
	config := &UnfoldConfig{
		Max:  100000,
		Step: 1,
	}

	for _, opt := range opts {
		opt(config)
	}

	current := 0
	var res []T
	for a := acc; p(a); a, current = f(a), current+1 {
		if current%config.Step == 0 {
			res = append(res, a)
		}
		if current >= config.Max {
			return res
		}
	}

	return res
}

// UnfoldI generates a slice by repeatedly applying a function to an accumulator.
// It includes the accumulator in the result as the first value.
// The length of the result is equal to i + 1.
// If i is negative, it returns nil.
// It stops after i iterations.
func UnfoldI[T any](acc T, f func(T) T, n int, opts ...UnfoldOption) []T {
	config := &UnfoldConfig{
		Max:  100000,
		Step: 1,
	}

	for _, opt := range opts {
		opt(config)
	}

	if n < 0 {
		return nil
	}

	n = min(n, config.Max)

	res := make([]T, 0, n)
	for i := 0; i <= n; i = i + config.Step {
		res = append(res, acc)
		acc = f(acc)
	}

	return res
}
