package path

import (
	"reflect"
	"testing"
)

func TestParseDottedPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
		wantErr  bool
	}{
		{
			name:     "simple path",
			input:    "a.b.c",
			expected: []string{"a", "b", "c"},
			wantErr:  false,
		},
		{
			name:     "simple path with slices",
			input:    "a.b[1].c",
			expected: []string{"a", "b[1]", "c"},
			wantErr:  false,
		},
		{
			name:     "path with quoted space",
			input:    `a."b c".d`,
			expected: []string{"a", "b c", "d"},
			wantErr:  false,
		},
		{
			name:     "path with quoted dot",
			input:    `a."b.c".d`,
			expected: []string{"a", "b.c", "d"},
			wantErr:  false,
		},
		{
			name:     "path with multiple quoted sections",
			input:    `"a b"."c.d"."e f"`,
			expected: []string{"a b", "c.d", "e f"},
			wantErr:  false,
		},
		{
			name:    "unclosed quotes",
			input:   `a."b c.d`,
			wantErr: true,
		},
		{
			name:     "empty segments",
			input:    "a..b",
			expected: []string{"a", "b"},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDottedPath(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDottedPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ParseDottedPath() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestParseJSONPointer(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
		wantErr  bool
	}{
		{
			name:     "simple path",
			input:    "/a/b/c",
			expected: []string{"a", "b", "c"},
			wantErr:  false,
		},
		{
			name:     "simple path with slices",
			input:    "/a/b/1/c",
			expected: []string{"a", "b", "1", "c"},
			wantErr:  false,
		},
		{
			name:     "path with encoded '~'",
			input:    "/a/b~0c/d",
			expected: []string{"a", "b~c", "d"},
			wantErr:  false,
		},
		{
			name:     "path with encoded '/'",
			input:    "/a/b~1c/d",
			expected: []string{"a", "b/c", "d"},
			wantErr:  false,
		},
		{
			name:     "path with quotes",
			input:    `/a/\"b\"/c/d`,
			expected: []string{"a", "\"b\"", "c", "d"},
			wantErr:  false,
		},
		{
			name:     "path with escaped control character",
			input:    `/a/\\u0000/b/c`,
			expected: []string{"a", "\\u0000", "b", "c"},
			wantErr:  false,
		},
		{
			name:    "unused encoding",
			input:   "/a/b~2c/d",
			wantErr: true,
		},
		{
			name:    "unused enscaping",
			input:   "/a/\\b/d",
			wantErr: true,
		},
		{
			name:     "empty path segment",
			input:    `/a//b`,
			expected: []string{"a", "b"},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseJSONPointer(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseJSONPointer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ParseJSONPointer() = %v, want %v", got, tt.expected)
			}
		})
	}
}
