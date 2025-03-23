package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// StringToInt64 converts a string to int64
func StringToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// Default returns the first non-zero value passed to it.
func Default[T comparable](val T, defaultVal T) T {
	var zero T // Zero value of type T
	if val != zero {
		return val
	}
	return defaultVal
}

// StringToFloat64 converts a string to float64
func StringToFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// StringToBool converts a string to bool
func StringToBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}

// TruncateString truncates a string to the specified length and adds ellipsis if needed
func TruncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength-3] + "..."
}

// RemoveWhitespace removes all whitespace characters from a string
func RemoveWhitespace(s string) string {
	return strings.Join(strings.Fields(s), "")
}

// IsEmpty checks if a string is empty or contains only whitespace
func IsEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

// ToSnakeCase converts a string to snake_case
func ToSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// ToCamelCase converts a string to camelCase
func ToCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i := 1; i < len(parts); i++ {
		if len(parts[i]) > 0 {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	return strings.Join(parts, "")
}

// FormatError formats an error message with optional parameters
func FormatError(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}
