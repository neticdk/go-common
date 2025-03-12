package string

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlugify(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"hello world", "hello-world"},
		{"hello  world", "hello-world"},
		{"hello---world", "hello-world"},
		{"hello- -world", "hello-world"},
		{"hello - - world", "hello-world"},
		{" hello world ", "hello-world"},
		{"!@#$%^&*()_+=-`~[]\\{}|;':\",./<>?", ""},
		{"HélLø Wörld", "hll-wrld"},
		{"你好世界", ""}, // Note: Non-latin characters are removed.
		{"with-hyphen", "with-hyphen"},
		{"with.dots.everywhere", "withdotseverywhere"},
		{"12345", "12345"},
		{"MixedCase123", "mixedcase123"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			actual := Slugify(tc.input)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
