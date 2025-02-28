package path

import (
	"strconv"
	"strings"
)

// JSONPointerPathToDottedPath converts a JSON Pointer path to a dotted path.
//
// It returns an error if the path is invalid.
func JSONPointerPathToDottedPath(path string) (string, error) {
	parts, err := ParseJSONPointer(path)
	if err != nil {
		return "", err
	}

	return PartsToDottedPath(parts)
}

// DottedPathToJSONPointerPath converts a dotted path to a JSON Pointer path.
//
// It returns an error if the path is invalid.
func DottedPathToJSONPointerPath(path string) (string, error) {
	parts, err := ParseDottedPath(path)
	if err != nil {
		return "", err
	}

	return PartsToJSONPointerPath(parts)
}

// PartsToDottedPath converts path parts to a dotted path.
//
// It returns an error if the parts constitute an invalid path.
func PartsToDottedPath(parts []string) (string, error) {
	var res strings.Builder
	var nextIsIndex bool

	for i, part := range parts {
		if strings.Contains(part, ".") || strings.Contains(part, "[") {
			res.WriteString(`"`)
			res.WriteString(part)
			res.WriteString(`"`)
		} else {
			if nextIsIndex {
				res.WriteString("[")
				res.WriteString(part)
				res.WriteString("]")
			} else {
				res.WriteString(part)
			}
		}

		nextIsIndex = false
		if len(parts) > i+1 {
			if _, err := strconv.Atoi(parts[i+1]); err == nil {
				nextIsIndex = true
			}
		}

		if i < len(parts)-1 && !nextIsIndex {
			res.WriteString(".")
		}
	}
	return res.String(), nil
}

// PartsToJSONPointerPath converts path parts to a JSON Pointer path.
//
// It returns an error if the parts constitute an invalid path.
func PartsToJSONPointerPath(parts []string) (string, error) {
	var res strings.Builder

	for _, part := range parts {

		res.WriteString("/")
		switch {
		case strings.Contains(part, "/"):
			res.WriteString(strings.ReplaceAll(part, "/", "~1"))
		case strings.Contains(part, "~"):
			res.WriteString(strings.ReplaceAll(part, "~", "~0"))
		default:
			res.WriteString(part)
		}
	}

	return res.String(), nil
}
