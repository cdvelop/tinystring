package tinystring

import "strings"

// Join concatenates the elements of a string slice to create a single string.
// If no separator is provided, it uses a space as default.
// Can be called with varargs to specify a custom separator.
// eg: Convert([]string{"Hello", "World"}).Join() => "Hello World"
// eg: Convert([]string{"Hello", "World"}).Join("-") => "Hello-World"
func (t *Text) Join(sep ...string) *Text {
	separator := " " // default separator is space
	if len(sep) > 0 {
		separator = sep[0]
	}

	// Handle case when we've received a string slice directly
	if t.contentSlice != nil {
		if len(t.contentSlice) == 0 {
			t.content = ""
		} else {
			t.content = strings.Join(t.contentSlice, separator)
		}
		return t
	}

	// If content is already a string, we split it and join it again with the new separator
	if t.content != "" {
		t.content = strings.Join(strings.Fields(t.content), separator)
	}

	return t
}

// JoinWithSpace concatenates the elements of a string slice to create a single string with the elements
// separated by spaces
// eg: JoinWithSpace([]string{"Hello", "World"}) => "Hello World"
// Deprecated: Use Convert(elements).Join() instead
func JoinWithSpace(elements []string) string {
	if len(elements) == 0 {
		return ""
	}

	return strings.Join(elements, " ")
}
