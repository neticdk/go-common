package path

import "testing"

func TestJSONPointerPathToDottedPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "simple path",
			input:    "/a/b/c",
			expected: "a.b.c",
			wantErr:  false,
		},
		{
			name:     "simple path with slices",
			input:    "/a/b/1/c",
			expected: "a.b[1].c",
			wantErr:  false,
		},
		{
			name:     "path with encoded '~'",
			input:    "/a/b~0c/d",
			expected: "a.b~c.d",
			wantErr:  false,
		},
		{
			name:     "path with encoded '/'",
			input:    "/a/b~1c/d",
			expected: "a.b/c.d",
			wantErr:  false,
		},
		{
			name:     "path with quotes",
			input:    `/a/\"b\"/c/d`,
			expected: `a."b".c.d`,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := JSONPointerPathToDottedPath(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONPointerPathToDottedPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if actual != tt.expected {
				t.Errorf("JSONPointerPathToDottedPath() = %v, want %v", actual, tt.expected)
			}
		})
	}
}

func TestDottedPathToJSONPointerPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "simple path",
			input:    "a.b.c",
			expected: "/a/b/c",
			wantErr:  false,
		},
		{
			name:     "simple path with slices",
			input:    "a.b[1].c",
			expected: "/a/b/1/c",
			wantErr:  false,
		},
		{
			name:     "path with encoded '~'",
			input:    "a.b~c.d",
			expected: "/a/b~0c/d",
			wantErr:  false,
		},
		{
			name:     "path with encoded '/'",
			input:    "a.b/c.d",
			expected: "/a/b~1c/d",
			wantErr:  false,
		},
		{
			name:     "path with quotes",
			input:    `a."b.c".d`,
			expected: `/a/b.c/d`,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := DottedPathToJSONPointerPath(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("DottedPathToJSONPointerPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if actual != tt.expected {
				t.Errorf("DottedPathToJSONPointerPath() = %v, want %v", actual, tt.expected)
			}
		})
	}
}
