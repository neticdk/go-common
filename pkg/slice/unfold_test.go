package slice

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnfold(t *testing.T) {
	tests := []struct {
		name     string
		acc      int
		f        func(int) int
		p        func(int) bool
		expected []int
	}{
		{
			name:     "double",
			acc:      1,
			f:        func(acc int) int { return acc * 2 },
			p:        func(acc int) bool { return acc < 100 },
			expected: []int{1, 2, 4, 8, 16, 32, 64},
		},
		{
			name:     "always false p",
			acc:      1,
			f:        func(acc int) int { return acc * 2 },
			p:        func(acc int) bool { return false },
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := Unfold(tt.acc, tt.f, tt.p)
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("expected %v but got %v", tt.expected, actual)
			}
		})
	}
}

func TestUnfoldI(t *testing.T) {
	tests := []struct {
		name     string
		acc      any
		f        func(any) any
		i        int
		expected []any
	}{
		{
			name:     "integers",
			acc:      1,
			f:        func(acc any) any { return acc.(int) * 2 },
			i:        7,
			expected: []any{1, 2, 4, 8, 16, 32, 64},
		},
		{
			name:     "negative i",
			acc:      1,
			f:        func(acc any) any { return acc.(int) * 2 },
			i:        -7,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := UnfoldI(tt.acc, tt.f, tt.i)
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("expected %v but got %v", tt.expected, actual)
			}
			assert.Equal(t, max(tt.i, 0), len(actual))
		})
	}
}
