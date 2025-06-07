// templui util templui.go - version: main installed by templui v0.73.1
package utils

import (
	"bytes"
	"fmt"
	"strings"

	"crypto/rand"

	"github.com/a-h/templ"
)

// TwMerge is not the original merge function provided by templui but i wanted to remove the
// external dependency, its by no means perfect but good enough
func TwMerge(classes ...string) string {
	return strings.Join(classes, " ")

	var classList []string
	classMap := make(map[string]string)

	for _, class := range strings.Fields(strings.Join(classes, " ")) {
		baseClass := strings.Split(class, "-")[0]
		if _, found := classMap[baseClass]; found {
			for i, v := range classList {
				if v != baseClass {
					continue
				}

				classList = append(classList[:i], classList[i+1:]...)
				break
			}
		}

		classMap[baseClass] = class
		classList = append(classList, baseClass)
	}

	var buf bytes.Buffer
	for _, baseClass := range classList {
		buf.WriteString(classMap[baseClass])
		buf.WriteByte(' ')
	}

	return buf.String()
}

// TwIf returns value if condition is true, otherwise an empty value of type T.
// Example: true, "bg-red-500" → "bg-red-500"
func If[T comparable](condition bool, value T) T {
	var empty T
	if condition {
		return value
	}
	return empty
}

// TwIfElse returns trueValue if condition is true, otherwise falseValue.
// Example: true, "bg-red-500", "bg-gray-300" → "bg-red-500"
func IfElse[T any](condition bool, trueValue T, falseValue T) T {
	if condition {
		return trueValue
	}
	return falseValue
}

// MergeAttributes combines multiple Attributes into one.
// Example: MergeAttributes(attr1, attr2) → combined attributes
func MergeAttributes(attrs ...templ.Attributes) templ.Attributes {
	merged := templ.Attributes{}
	for _, attr := range attrs {
		for k, v := range attr {
			merged[k] = v
		}
	}
	return merged
}

// RandomID generates a random ID string.
// Example: RandomID() → "id-1a2b3c"
func RandomID() string {
	return fmt.Sprintf("id-%s", rand.Text())
}
