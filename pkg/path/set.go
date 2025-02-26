package path

import (
	"fmt"
	"strconv"
	"strings"
)

// SetPathValue sets the given value at the given path in the data.
//
// Currently only supports []any and map[string]any data types.
// It returns an error if the path is invalid or the value cannot be set.
func SetPathValue[T any](data T, path string, value any) error {
	parts, err := ParseDottedPath(path)
	if err != nil {
		return fmt.Errorf("parsing dotted path: %w", err)
	}

	if _, err := setNonReflectPathValueRecursive(data, parts, value); err != nil {
		return fmt.Errorf("setting path value: %w", err)
	}

	return nil
}

// getNonReflectPathValueRecursive returns the value at the given path in the data.
//
// It is used for non-reflectable data types like map[string]any and []any.
// It returns an error if the path is invalid or the value is not found.
func setNonReflectPathValueRecursive(data any, parts []string, value any) (any, error) {
	var err error
	// If path parts is empty, return the new value
	if len(parts) == 0 {
		return value, nil
	}

	// If the data is nil, create the appropriate data type for the rest of the path parts
	if data == nil {
		if isSlicePart(parts[0]) {
			data = []any{}
		} else {
			data = map[string]any{}
		}
	}

	part := parts[0]
	parts = parts[1:]

	switch val := data.(type) {
	case map[string]any:
		val[part], err = setNonReflectPathValueRecursive(val[part], parts, value)
		if err != nil {
			return nil, fmt.Errorf("setting key %s in map-value %v: %w", part, val, err)
		}
		return val, nil

	case []any:
		i, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid index %s for slice-value %v: %w", part, val, err)
		}
		if i == len(val) {
			val = append(val, nil)
		} else if i < 0 || i >= len(val) {
			return nil, fmt.Errorf("index %d out of bounds for slice-value %v", i, val)
		}
		val[i], err = setNonReflectPathValueRecursive(val[i], parts, value)
		if err != nil {
			return nil, fmt.Errorf("setting index %d in slice-value %v: %w", i, val, err)
		}
		return val, nil
	default:
		return nil, fmt.Errorf("cannot set value at path %s from non-map non-slice value %v", strings.Join(parts, "."), data)
	}
}

// isSlicePart returns true if the given part is a slice index.
func isSlicePart(part string) bool {
	if _, err := strconv.Atoi(part); err == nil {
		return true
	}
	return false
}
