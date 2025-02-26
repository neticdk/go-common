package path

import (
	"reflect"
	"testing"
)

func TestGetPathValue(t *testing.T) {
	tests := []struct {
		name     string
		data     any
		path     string
		expected any
		wantErr  bool
	}{
		{
			name: "simple map",
			data: map[string]any{
				"a": 1,
			},
			path:     "a",
			expected: 1,
			wantErr:  false,
		},
		{
			name: "map with nested slices and maps",
			data: map[string]any{
				"c": map[string]any{
					"d": []any{
						1,
						map[string]any{"e": 2},
					},
				},
			},
			path:     "c.d[1].e",
			expected: 2,
			wantErr:  false,
		},
		{
			name:     "simple slice",
			data:     []any{1, 2, 3},
			path:     "[1]",
			expected: 2,
			wantErr:  false,
		},
		{
			name: "slice with nested slices and maps",
			data: []any{
				map[string]any{
					"d": []any{
						1,
						map[string]any{"e": 2},
					},
				},
			},
			path:     "[0].d[1].e",
			expected: 2,
			wantErr:  false,
		},
		{
			name:     "simple struct",
			data:     struct{ A int }{A: 1},
			path:     "A",
			expected: 1,
			wantErr:  false,
		},
		{
			name: "struct with nested structs",
			data: struct {
				A struct {
					B []any
				}
			}{
				A: struct {
					B []any
				}{
					B: []any{
						1,
					},
				},
			},
			path:     "A.B[0]",
			expected: 1,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPathValue(tt.data, tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPathValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("GetPathValue() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetMapPathValue(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string]any
		path     string
		expected any
		wantErr  bool
	}{
		{
			name: "simple map",
			data: map[string]any{
				"a": 1,
			},
			path:     "a",
			expected: 1,
			wantErr:  false,
		},
		{
			name: "map with nested slices and maps",
			data: map[string]any{
				"c": map[string]any{
					"d": []any{
						1,
						map[string]any{"e": 2},
					},
				},
			},
			path:     "c.d[1].e",
			expected: 2,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetMapPathValue(tt.data, tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMapPathValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("GetMapPathValue() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetSlicePathValue(t *testing.T) {
	tests := []struct {
		name     string
		data     []any
		path     string
		expected any
		wantErr  bool
	}{
		{
			name:     "simple slice",
			data:     []any{1, 2, 3},
			path:     "[1]",
			expected: 2,
			wantErr:  false,
		},
		{
			name: "slice with nested slices and maps",
			data: []any{
				map[string]any{
					"d": []any{
						1,
						map[string]any{"e": 2},
					},
				},
			},
			path:     "[0].d[1].e",
			expected: 2,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSlicePathValue(tt.data, tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSlicePathValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("GetSlicePathValue() = %v, want %v", got, tt.expected)
			}
		})
	}
}
