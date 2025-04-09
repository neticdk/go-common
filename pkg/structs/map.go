// Deprecated: structs is deprecated - use github.com/neticdk/go-stdlib/xstructs
package structs

import (
	"github.com/neticdk/go-stdlib/xstructs"
)

// Deprecated: ToMap is deprecated - use github.com/neticdk/go-stdlib/xstructs.ToMap
// ToMap converts a struct or map to a map[string]any.
// It handles nested structs, maps, and slices.
// It uses the "json" and "yaml" tags to determine the key names.
// It respects the "omitempty" tag for fields.
// It respects the "inline" tag for nested structs.
// It respects the "-" tag to omit fields.
//
// If the input is nil, it returns nil.
// If the input is not a struct or map, it returns an error.
func ToMap(obj any) (map[string]any, error) {
	return xstructs.ToMap(obj)
}
