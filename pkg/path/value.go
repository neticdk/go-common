package path

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// GetPathValue returns the value at the given path in the data.
//
// It returns an error if the path is invalid or the value is not found.
func GetPathValue[T any](data T, path string) (any, error) {
	parts, err := ParseDottedPath(path)
	if err != nil {
		return nil, err
	}

	return getPathValueRecursive(data, parts)
}

// GetMapPathValue returns the value at the given path in the map.
//
// It performs faster than the generic GetPathValue function.
// It returns an error if the path is invalid or the value is not found.
func GetMapPathValue(data map[string]any, path string) (any, error) {
	parts, err := ParseDottedPath(path)
	if err != nil {
		return nil, err
	}

	return getMapPathValueRecursive(data, parts)
}

// GetSlicePathValue returns the value at the given path in the slice.
//
// It performs faster than the generic GetPathValue function.
// It returns an error if the path is invalid or the value is not found.
func GetSlicePathValue(data []any, path string) (any, error) {
	parts, err := ParseDottedPath(path)
	if err != nil {
		return nil, err
	}

	return getSlicePathValueRecursive(data, parts)
}

func getPathValueRecursive(data any, parts []string) (any, error) {
	// If path parts is empty, return the data
	if len(parts) == 0 {
		return data, nil
	}

	if isNil(data) {
		if len(parts) > 0 {
			return nil, fmt.Errorf("nil value for non-empty path %s", strings.Join(parts, "."))
		}
		return data, nil
	}

	part := parts[0]
	parts = parts[1:]
	switch val := reflect.ValueOf(data); val.Kind() {
	case reflect.Struct:
		if f := reflect.Indirect(val).FieldByName(part); f.IsValid() {
			return getPathValueRecursive(f.Interface(), parts)
		} else {
			return nil, fmt.Errorf("cannot find field %s in struct-value %v", part, val)
		}

	case reflect.Map:
		if v, ok := data.(map[string]any)[part]; ok {
			return getPathValueRecursive(v, parts)
		} else {
			return nil, fmt.Errorf("cannot find key %s in map-value %v", part, val)
		}

	case reflect.Slice:
		v, ok := data.([]any)
		if !ok {
			return nil, fmt.Errorf("cannot get path %s from non-slice value %v", strings.Join(parts, "."), data)
		}
		if i, err := strconv.Atoi(part); err == nil {
			if i < 0 || i >= len(v) {
				return nil, fmt.Errorf("index %d out of bounds for slice-value %v", i, val)
			}
			return getPathValueRecursive(v[i], parts)
		} else {
			return nil, fmt.Errorf("invalid index %s for slice-value %v", part, val)
		}
	default:
		return nil, fmt.Errorf("cannot get path %s from non-map non-slice value %v", strings.Join(parts, "."), data)
	}
}

// getNonReflectPathValueRecursive returns the value at the given path in the data.
//
// It is used for non-reflectable data types like map[string]any and []any.
// This is increases performance as opposed to the reflect-version
// but is limited to only map[string]any and []any data types.
// It returns an error if the path is invalid or the value is not found.
func getNonReflectPathValueRecursive(data any, parts []string) (any, error) {
	// If path parts is empty, return the data
	if len(parts) == 0 {
		return data, nil
	}

	if data == nil {
		if len(parts) > 0 {
			return nil, fmt.Errorf("nil value for non-empty path %s", strings.Join(parts, "."))
		}
		return data, nil
	}

	switch val := data.(type) {
	case map[string]any:
		return getMapPathValueRecursive(val, parts)

	case []any:
		return getSlicePathValueRecursive(val, parts)
	default:
		return nil, fmt.Errorf("cannot get path %s from non-map non-slice value %v", strings.Join(parts, "."), data)
	}
}

// getMapPathValueRecursive returns the value at the given path in the map.
//
// It is used for map[string]any data type.
// It performs faster than the generic GetPathValue function.
// It returns an error if the path is invalid or the value is not found.
func getMapPathValueRecursive(data map[string]any, parts []string) (any, error) {
	// If path parts is empty, return the data
	if len(parts) == 0 {
		return data, nil
	}

	if v, ok := data[parts[0]]; ok {
		return getNonReflectPathValueRecursive(v, parts[1:])
	}

	return nil, fmt.Errorf("cannot find key %s in map-value %v", parts[0], data)
}

// getSlicePathValueRecursive returns the value at the given path in the slice.
//
// It is used for []any data type.
// It performs faster than the generic GetPathValue function.
// It returns an error if the path is invalid or the value is not found.
func getSlicePathValueRecursive(data []any, parts []string) (any, error) {
	// If path parts is empty, return the data
	if len(parts) == 0 {
		return data, nil
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("empty slice-value %v", data)
	}

	if i, err := strconv.Atoi(parts[0]); err == nil {
		if i < 0 || i >= len(data) {
			return nil, fmt.Errorf("index %d out of bounds for slice-value %v", i, data)
		}
		return getNonReflectPathValueRecursive(data[i], parts[1:])
	}

	return nil, fmt.Errorf("invalid index %s for slice-value %v", parts[0], data)
}

// isNil checks if generic data is nil
//
// It uses reflect, so it is slower than a direct comparison.
// It returns true if the data is nil, false otherwise.
func isNil[T any](data T) bool {
	v := reflect.ValueOf(data)
	return (v.Kind() == reflect.Ptr ||
		v.Kind() == reflect.Interface ||
		v.Kind() == reflect.Slice ||
		v.Kind() == reflect.Map ||
		v.Kind() == reflect.Chan ||
		v.Kind() == reflect.Func) && v.IsNil()
}
