package path

import (
	"reflect"
	"testing"
)

func TestSetPathValue(t *testing.T) {
	tests := []struct {
		name     string
		data     any
		path     string
		value    any
		expected any
		wantErr  bool
	}{
		{
			name: "simple map",
			data: map[string]any{
				"a": 1,
			},
			path:  "a",
			value: 2,
			expected: map[string]any{
				"a": 2,
			},
			wantErr: false,
		},
		{
			name: "map with slices and maps",
			data: map[string]any{
				"a": []any{
					map[string]any{"b": 1},
				},
			},
			path:  "a[0].b",
			value: 2,
			expected: map[string]any{
				"a": []any{
					map[string]any{"b": 2},
				},
			},
			wantErr: false,
		},
		{
			name: "new map path",
			data: map[string]any{
				"a": 1,
			},
			path:  "b.c",
			value: 2,
			expected: map[string]any{
				"a": 1,
				"b": map[string]any{
					"c": 2,
				},
			},
			wantErr: false,
		},
		{
			name: "new slice path",
			data: map[string]any{
				"a": 1,
			},
			path:  "b[0].c",
			value: 2,
			expected: map[string]any{
				"a": 1,
				"b": []any{
					map[string]any{
						"c": 2,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "new map path with slice value",
			data: map[string]any{
				"a": 1,
			},
			path:  "b.c",
			value: []any{1, 2, 3},
			expected: map[string]any{
				"a": 1,
				"b": map[string]any{
					"c": []any{1, 2, 3},
				},
			},
			wantErr: false,
		},
		{
			name: "new slice path with map value",
			data: map[string]any{
				"a": 1,
			},
			path: "b[0].c",
			value: map[string]any{
				"d": 2,
			},
			expected: map[string]any{
				"a": 1,
				"b": []any{
					map[string]any{
						"c": map[string]any{
							"d": 2,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "nil data",
			data: nil,
			path: "b[0].c",
			value: map[string]any{
				"d": 2,
			},
			expected: nil,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SetPathValue(tt.data, tt.path, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetPathValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(tt.data, tt.expected) {
				t.Errorf("SetPathValue() = %v, want %v", tt.data, tt.expected)
			}
		})
	}
}
