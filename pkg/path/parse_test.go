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
		name               string
		input              any
		expectedWitInc     []string
		expectedWithoutInc []string
	}{
		{
			name: "simple map",
			input: map[string]any{
				"a": map[string]any{
					"b": map[string]any{
						"c": "value",
					},
				},
			},
			expectedWitInc:     []string{"a.b.c"},
			expectedWithoutInc: []string{"a.b.c"},
		},
		{
			name: "simple map with multiple maps",
			input: map[string]any{
				"a": map[string]any{
					"b": map[string]any{
						"c": "value",
					},
				},
				"c": map[string]any{
					"b": map[string]any{
						"a": "value",
					},
				},
			},
			expectedWitInc:     []string{"a.b.c", "c.b.a"},
			expectedWithoutInc: []string{"a.b.c", "c.b.a"},
		},
		{
			name: "map with simple slice",
			input: map[string]any{
				"a": map[string]any{
					"b": map[string]any{
						"c": []any{
							"value1",
							"value2",
							"value2",
						},
					},
				},
			},
			expectedWitInc:     []string{"a.b.c[0]", "a.b.c[1]", "a.b.c[2]"},
			expectedWithoutInc: []string{"a.b.c[0]", "a.b.c[1]", "a.b.c[2]"},
		},
		{
			name: "map with slice of maps",
			input: map[string]any{
				"a": map[string]any{
					"b": map[string]any{
						"c": []any{
							map[string]any{
								"a": map[string]any{
									"c": "value",
								},
								"b": map[string]any{
									"c": "value",
								},
							},
							map[string]any{
								"a": map[string]any{
									"c": "value",
								},
								"b": map[string]any{
									"c": "value",
								},
							},
							map[string]any{
								"a": map[string]any{
									"c": "value",
								},
								"b": map[string]any{
									"c": "value",
								},
							},
						},
					},
				},
			},
			expectedWitInc:     []string{"a.b.c[0].a.c", "a.b.c[0].b.c", "a.b.c[1].a.c", "a.b.c[1].b.c", "a.b.c[2].a.c", "a.b.c[2].b.c"},
			expectedWithoutInc: []string{"a.b.c[0].a.c", "a.b.c[0].b.c", "a.b.c[1].a.c", "a.b.c[1].b.c", "a.b.c[2].a.c", "a.b.c[2].b.c"},
		},
		{
			name:               "empty string",
			input:              "",
			expectedWitInc:     []string{""},
			expectedWithoutInc: []string{""},
		},
		{
			name:               "nil input",
			input:              nil,
			expectedWitInc:     []string{""},
			expectedWithoutInc: []string{""},
		},
		{
			name: "sort order",
			input: map[string]any{
				"some": map[string]any{
					"test": map[string]any{
						"key": "value",
					},
				},
				"aSome": map[string]any{
					"test": map[string]any{
						"key": "value",
					},
				},
				"bSome": map[string]any{
					"key": "value",
				},
				"c": map[string]any{
					"key": "value",
				},
			},
			expectedWitInc:     []string{"bSome.key", "c.key", "aSome.test.key", "some.test.key"},
			expectedWithoutInc: []string{"bSome.key", "c.key", "aSome.test.key", "some.test.key"},
		},
		{
			name: "slice sorting order",
			input: map[string]any{
				"some": []any{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k"},
			},
			expectedWitInc:     []string{"some[0]", "some[1]", "some[2]", "some[3]", "some[4]", "some[5]", "some[6]", "some[7]", "some[8]", "some[9]", "some[10]"},
			expectedWithoutInc: []string{"some[0]", "some[1]", "some[2]", "some[3]", "some[4]", "some[5]", "some[6]", "some[7]", "some[8]", "some[9]", "some[10]"},
		},
		{
			name: "slice elements with nil values",
			input: map[string]any{
				"some": []any{"a", "b", nil, "d", "e", "f", nil, "h", "i", "j", "k"},
			},
			expectedWitInc:     []string{"some[0]", "some[1]", "some[3]", "some[4]", "some[5]", "some[7]", "some[8]", "some[9]", "some[10]"},
			expectedWithoutInc: []string{"some[0]", "some[1]", "some[2]", "some[3]", "some[4]", "some[5]", "some[6]", "some[7]", "some[8]", "some[9]", "some[10]"},
		},
		{
			name: "slice of maps with nil values",
			input: map[string]any{
				"some": []any{
					map[string]any{
						"key": "value",
					},
					map[string]any{
						"key":  "value",
						"list": []any{"a", "b", nil, "d", "e", "f", nil, "h", "i", "j", "k"},
					},
					nil,
					map[string]any{
						"key": map[string]any{
							"key": "value",
						},
					},
				},
			},
			expectedWitInc:     []string{"some[0].key", "some[1].key", "some[1].list[0]", "some[1].list[1]", "some[1].list[3]", "some[1].list[4]", "some[1].list[5]", "some[1].list[7]", "some[1].list[8]", "some[1].list[9]", "some[1].list[10]", "some[3].key.key"},
			expectedWithoutInc: []string{"some[2]", "some[0].key", "some[1].key", "some[1].list[0]", "some[1].list[1]", "some[1].list[2]", "some[1].list[3]", "some[1].list[4]", "some[1].list[5]", "some[1].list[6]", "some[1].list[7]", "some[1].list[8]", "some[1].list[9]", "some[1].list[10]", "some[3].key.key"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paths := ParseAnyToDottedPath(tt.input, func(key, value any) bool {
				// include everything except nil values when key is type int
				if _, ok := key.(int); ok {
					if value == nil {
						return false
					}
				}
				return true
			})
			assert.Equal(t, tt.expectedWitInc, paths)
			paths = ParseAnyToDottedPath(tt.input, nil)
			assert.Equal(t, tt.expectedWithoutInc, paths)
		})
	}
}
