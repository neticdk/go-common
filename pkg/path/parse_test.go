package path

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
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
			expected: []string{"a", "b", "1", "c"},
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

func TestParseAnyToDottedPath(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected []string
	}{
		{
			name: "simple map",
			input: map[string]interface{}{
				"a": map[string]interface{}{
					"b": map[string]interface{}{
						"c": "value",
					},
				},
			},
			expected: []string{"a.b.c"},
		},
		{
			name: "simple map with multiple maps",
			input: map[string]interface{}{
				"a": map[string]interface{}{
					"b": map[string]interface{}{
						"c": "value",
					},
				},
				"c": map[string]interface{}{
					"b": map[string]interface{}{
						"a": "value",
					},
				},
			},
			expected: []string{"a.b.c", "c.b.a"},
		},
		{
			name: "map with simple slice",
			input: map[string]interface{}{
				"a": map[string]interface{}{
					"b": map[string]interface{}{
						"c": []interface{}{
							"value1",
							"value2",
							"value2",
						},
					},
				},
			},
			expected: []string{"a.b.c[0]", "a.b.c[1]", "a.b.c[2]"},
		},
		{
			name: "map with slice of maps",
			input: map[string]interface{}{
				"a": map[string]interface{}{
					"b": map[string]interface{}{
						"c": []interface{}{
							map[string]interface{}{
								"a": map[string]interface{}{
									"c": "value",
								},
								"b": map[string]interface{}{
									"c": "value",
								},
							},
							map[string]interface{}{
								"a": map[string]interface{}{
									"c": "value",
								},
								"b": map[string]interface{}{
									"c": "value",
								},
							},
							map[string]interface{}{
								"a": map[string]interface{}{
									"c": "value",
								},
								"b": map[string]interface{}{
									"c": "value",
								},
							},
						},
					},
				},
			},
			expected: []string{"a.b.c[0].a.c", "a.b.c[0].b.c", "a.b.c[1].a.c", "a.b.c[1].b.c", "a.b.c[2].a.c", "a.b.c[2].b.c"},
		},
		{
			name:     "empty string",
			input:    "",
			expected: []string{""},
		},
		{
			name:     "nil input",
			input:    nil,
			expected: []string{""},
		},
		{
			name: "sort order",
			input: map[string]interface{}{
				"some": map[string]interface{}{
					"test": map[string]interface{}{
						"key": "value",
					},
				},
				"aSome": map[string]interface{}{
					"test": map[string]interface{}{
						"key": "value",
					},
				},
				"bSome": map[string]interface{}{
					"key": "value",
				},
				"c": map[string]interface{}{
					"key": "value",
				},
			},
			expected: []string{"bSome.key", "c.key", "aSome.test.key", "some.test.key"},
		},
		{
			name: "slice sorting order",
			input: map[string]interface{}{
				"some": []interface{}{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k"},
			},
			expected: []string{"some[0]", "some[1]", "some[2]", "some[3]", "some[4]", "some[5]", "some[6]", "some[7]", "some[8]", "some[9]", "some[10]"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paths := ParseAnyToDottedPath(tt.input)
			assert.Equal(t, tt.expected, paths)
		})
	}
}
