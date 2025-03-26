package path

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

// ParseDottedPath parses a dotted path into its segments.
//
// It returns an error if the path is invalid.
func ParseDottedPath(path string) ([]string, error) {
	var segments []string
	var current strings.Builder
	var inQuotes bool
	var inBraces bool

	runes := []rune(path)
	for i := 0; i < len(runes); i++ {
		ch := runes[i]

		switch {
		case ch == '"':
			if inQuotes {
				inQuotes = false
				segments = append(segments, current.String())
				current.Reset()
			} else {
				// Start of quoted section
				inQuotes = true
				// If we have accumulated any characters, add them as a segment
				if current.Len() > 0 {
					segments = append(segments, current.String())
					current.Reset()
				}
			}

		case ch == '.' && !inQuotes:
			// Dot outside quotes is a separator
			if current.Len() > 0 {
				segments = append(segments, current.String())
				current.Reset()
			}

		case ch == '[' && !inQuotes:
			// Bracket outside quotes is a slice separator
			if current.Len() > 0 {
				segments = append(segments, current.String())
				current.Reset()
			}
			inBraces = true

		case ch == ']' && inBraces:
			// End of slice separator
			if current.Len() == 0 {
				return nil, fmt.Errorf("empty index in path")
			}
			segments = append(segments, current.String())
			current.Reset()
			inBraces = false

		case unicode.IsSpace(ch) && !inQuotes:
			// Ignore spaces outside quotes
			continue

		default:
			current.WriteRune(ch)
		}
	}

	// Add final segment if there is one
	if current.Len() > 0 {
		segments = append(segments, current.String())
	}

	// Validate that we're not ending with unclosed quotes
	if inQuotes {
		return nil, fmt.Errorf("unclosed quotes in path")
	}

	return segments, nil
}

// ParseJSONPointer parses a JSON Pointer path into its segments.
// The JSON Pointer path is a string that represents a path to a value in a JSON document.
// The JSON Pointer path follows the specification defined in RFC 6901.
//
// It returns an error if the path is invalid.
func ParseJSONPointer(path string) ([]string, error) {
	var segments []string
	var current strings.Builder
	var inEncode bool
	var inEscape bool
	var isStarted bool

	runes := []rune(path)
	for i := 0; i < len(runes); i++ {
		ch := runes[i]

		switch {
		// If the character is ~ it encodes a special character ('/' or '~')
		case ch == '~':
			if inEncode {
				return nil, fmt.Errorf("invalid escape sequence in path")
			}
			inEncode = true

		// Check if the character '0' is part of an encoded character '~'
		case ch == '0':
			if inEncode {
				current.WriteRune('~')
				inEncode = false
			} else {
				current.WriteRune(ch)
			}

		// Check if the character '1' is part of an encoded character '/'
		case ch == '1':
			if inEncode {
				current.WriteRune('/')
				inEncode = false
			} else {
				current.WriteRune(ch)
			}

		// If the character is a control character it must be escaped
		case unicode.IsControl(ch):
			if !inEscape {
				return nil, fmt.Errorf("control character must be escaped in path")
			}
			current.WriteRune('\\')
			current.WriteRune(ch)
			inEscape = false

		// If the character is a backslash, it is the start of an escape sequence
		case ch == '\\':
			if inEscape {
				current.WriteRune(ch)
				inEscape = false
			} else {
				inEscape = true
			}

		// If the character is a quote, it is the start or end of a quoted section
		case ch == '"':
			if !inEscape {
				return nil, fmt.Errorf("quotes must be escaped in path")
			}
			current.WriteRune(ch)
			inEscape = false

		case ch == '/':
			// Check if trying to encode or escape slash
			if inEncode || inEscape {
				return nil, fmt.Errorf("invalid escape sequence in path")
			}
			if !isStarted {
				isStarted = true
			} else if current.Len() > 0 {
				segments = append(segments, current.String())
				current.Reset()
			}

		default:
			if inEncode || inEscape {
				return nil, fmt.Errorf("invalid escape sequence in path")
			}
			current.WriteRune(ch)
		}
	}

	// Add final segment if there is one
	if current.Len() > 0 {
		segments = append(segments, current.String())
	}

	return segments, nil
}

// ParseAnyToDottedPath parses any to dotted path.
//
// It returns a list of paths found in the given data, sorted by "path length" and lexical order.
func ParseAnyToDottedPath(data any, prefix string) []string {
	allPaths := extractPathsRecursively(data, []string{prefix})

	if len(allPaths) == 0 {
		return []string{} // Handle empty paths gracefully
	}

	parsedPaths := make([]string, 0, len(allPaths))

	for _, parts := range allPaths {
		dottedPath := PartsToDottedPath(parts)
		parsedPaths = append(parsedPaths, dottedPath)
	}

	// Sort paths by length and lexical order
	sort.Slice(parsedPaths, func(i, j int) bool {
		partsI, partsJ := strings.Split(parsedPaths[i], "."), strings.Split(parsedPaths[j], ".")
		lenI, lenJ := len(partsI), len(partsJ)
		if lenI != lenJ {
			return lenI < lenJ
		}
		for k := 0; k < lenI; k++ {
			if partsI[k] != partsJ[k] {
				return partsI[k] < partsJ[k]
			}
		}
		return parsedPaths[i] < parsedPaths[j]
	})

	return parsedPaths
}

func extractPathsRecursively(data any, path []string) [][]string {
	var allPaths [][]string
	extractPaths(data, path, &allPaths)
	return allPaths
}

func extractPaths(data any, path []string, allPaths *[][]string) {
	switch value := data.(type) {
	case map[string]any:
		for k, v := range value {
			newPath := make([]string, len(path), len(path)+1)
			copy(newPath, path)
			newPath = append(newPath, k)
			extractPaths(v, newPath, allPaths)
		}
	case []any:
		for i, item := range value {
			newPath := make([]string, len(path), len(path)+1)
			copy(newPath, path)
			newPath = append(newPath, strconv.Itoa(i))
			extractPaths(item, newPath, allPaths)
		}
	default:
		*allPaths = append(*allPaths, path)
	}
}
