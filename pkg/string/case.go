package string

import (
	"github.com/neticdk/go-stdlib/xstrings"
)

// ToKebabCase converts a string to kebab case.
// Deprecated: ToKebabCase is deprecated and has been moved to github.com/neticdk/go-stdlib/xstrings.ToKebabCase
func ToKebabCase(s string) string {
	return xstrings.ToKebabCase(s)
}

// ToSnakeCase converts a string to snake case.
// Deprecated: ToSnakeCase is deprecated and has been moved to github.com/neticdk/go-stdlib/xstrings.ToSnakeCase
func ToSnakeCase(s string) string {
	return xstrings.ToSnakeCase(s)
}

// ToDotCase converts a string to dot case.
// Deprecated: ToDotCase is deprecated and has been moved to github.com/neticdk/go-stdlib/xstrings.ToDotCase
func ToDotCase(s string) string {
	return xstrings.ToDotCase(s)
}

// ToCamelCase converts a string to camelCase.
// Deprecated: ToCamelCase is deprecated and has been moved to github.com/neticdk/go-stdlib/xstrings.ToCamelCase
func ToCamelCase(s string) string {
	return xstrings.ToCamelCase(s)
}

// ToPascalCase converts a string to pascal case.
// Deprecated: ToPascalCase is deprecated and has been moved to github.com/neticdk/go-stdlib/xstrings.ToPascalCase
func ToPascalCase(s string) string {
	return xstrings.ToPascalCase(s)
}

// ToDelimited converts a string to a delimited format using the specified
// delimiter
// Deprecated: ToDelimited is deprecated and has been moved to github.com/neticdk/go-stdlib/xstrings.ToDelimited
func ToDelimited(s string, delimiter string) string {
	return xstrings.ToDelimited(s, delimiter)
}
