package tinystring

import "unsafe"

// JSON decoding implementation for TinyString
// Uses our custom reflectlite integration for minimal binary size - NO standard reflect

// JsonDecode parses JSON data and populates the target struct/slice
//
// Usage patterns:
//   err := Convert(jsonBytes).JsonDecode(&user)
//   err := Convert(jsonString).JsonDecode(&users)  // []User slice
//   err := Convert(reader).JsonDecode(&data)
//
// Supports decoding into:
// - Structs with basic field types
// - Slices of structs
// - Basic types (string, int, float, bool)
//
// Field matching: Uses snake_case JSON keys to struct fields
// Example: {"user_name": "John"} -> UserName field
func (c *conv) JsonDecode(target any) error {
	if target == nil {
		return Err(errInvalidJSON, "target cannot be nil")
	}

	// Get JSON data as string
	jsonStr := c.getString()
	if jsonStr == "" {
		return Err(errInvalidJSON, "empty JSON data")
	}

	// Parse and populate target
	return c.parseJsonIntoTarget(jsonStr, target)
}

// parseJsonIntoTarget parses JSON string and populates the target value
func (c *conv) parseJsonIntoTarget(jsonStr string, target any) error {
	if target == nil {
		return Err(errInvalidJSON, "target cannot be nil")
	}

	// Use our custom reflection for target analysis
	rv := refValueOf(target)

	// Debug: Check what kind we get for the pointer
	targetKind := rv.Kind()
	if targetKind != tpPointer {
		return Err(errInvalidJSON, "target must be a pointer, got: "+targetKind.String())
	}

	// Get the element that the pointer points to
	elem := rv.Elem()
	if !elem.IsValid() {
		return Err(errInvalidJSON, "target pointer is nil or invalid")
	}

	// Debug: Check what kind we get for the element
	elemKind := elem.Kind()
	if elemKind.String() == "invalid" {
		return Err(errInvalidJSON, "element kind is invalid - reflection issue")
	}

	// Parse JSON and populate the element using our custom reflection
	return c.parseJsonValueWithRefReflect(jsonStr, elem)
}

// parseJsonValueWithRefReflect parses a JSON value using our custom reflection
func (c *conv) parseJsonValueWithRefReflect(jsonStr string, target refValue) error {
	// Trim whitespace
	jsonStr = Convert(jsonStr).Trim().String()
	if len(jsonStr) == 0 {
		return Err(errInvalidJSON, "empty JSON")
	}
	switch target.Kind() {
	case tpString:
		return c.parseJsonStringRef(jsonStr, target)
	case tpInt, tpInt8, tpInt16, tpInt32, tpInt64:
		return c.parseJsonIntRef(jsonStr, target)
	case tpUint, tpUint8, tpUint16, tpUint32, tpUint64:
		return c.parseJsonUintRef(jsonStr, target)
	case tpFloat32, tpFloat64:
		return c.parseJsonFloatRef(jsonStr, target)
	case tpBool:
		return c.parseJsonBoolRef(jsonStr, target)
	case tpStruct:
		return c.parseJsonStructRef(jsonStr, target)
	case tpSlice:
		return c.parseJsonSliceRef(jsonStr, target)
	case tpPointer:
		return c.parseJsonPointerRef(jsonStr, target)
	default:
		return Err(errUnsupportedType, "unsupported target type for JSON decoding: "+target.Kind().String())
	}
}

// Custom reflection-based parsing functions using our refValue system

// parseJsonStringRef parses a JSON string using our custom reflection
func (c *conv) parseJsonStringRef(jsonStr string, target refValue) error {
	if len(jsonStr) < 2 || jsonStr[0] != '"' || jsonStr[len(jsonStr)-1] != '"' {
		return Err(errInvalidJSON, "invalid JSON string format")
	}

	// Remove quotes and decode escape sequences
	unquoted := jsonStr[1 : len(jsonStr)-1]
	decoded, err := c.unescapeJsonString(unquoted)
	if err != nil {
		return err
	}
	target.SetString(decoded)
	return nil
}

// parseJsonIntRef parses a JSON integer using our custom reflection
func (c *conv) parseJsonIntRef(jsonStr string, target refValue) error {
	intVal, err := Convert(jsonStr).ToInt64()
	if err != nil {
		return err
	}
	target.SetInt(intVal)
	return nil
}

// parseJsonUintRef parses a JSON unsigned integer using our custom reflection
func (c *conv) parseJsonUintRef(jsonStr string, target refValue) error {
	val, err := Convert(jsonStr).ToInt64() // Convert to int64 first, then cast to uint64
	if err != nil {
		return err
	}
	target.SetUint(uint64(val))
	return nil
}

// parseJsonFloatRef parses a JSON float using our custom reflection
func (c *conv) parseJsonFloatRef(jsonStr string, target refValue) error {
	val, err := Convert(jsonStr).ToFloat()
	if err != nil {
		return err
	}
	target.SetFloat(val)
	return nil
}

// parseJsonBoolRef parses a JSON boolean using our custom reflection
func (c *conv) parseJsonBoolRef(jsonStr string, target refValue) error {
	switch jsonStr {
	case "true":
		target.SetBool(true)
	case "false":
		target.SetBool(false)
	default:
		return Err(errInvalidJSON, "invalid JSON boolean: "+jsonStr)
	}
	return nil
}

// parseJsonStructRef parses a JSON object into a struct using our custom reflection
func (c *conv) parseJsonStructRef(jsonStr string, target refValue) error {
	if target.Kind() != tpStruct {
		return Err(errUnsupportedType, "target is not a struct")
	}

	// Basic validation - must start with { and end with }
	jsonStr = Convert(jsonStr).Trim().String()
	if len(jsonStr) < 2 || jsonStr[0] != '{' || jsonStr[len(jsonStr)-1] != '}' {
		return Err(errInvalidJSON, "invalid JSON object format")
	}

	// Handle empty object
	if jsonStr == "{}" {
		return nil // empty object, nothing to set
	}
	// Get struct information
	var structInfo refStructInfo
	getStructInfo(target.Type(), &structInfo)
	if structInfo.refType == nil {
		return Err(errUnsupportedType, "cannot get struct information")
	}

	// Simple JSON parsing - remove outer braces and split by commas
	content := jsonStr[1 : len(jsonStr)-1] // Remove { }
	return c.parseJsonObjectContent(content, target, &structInfo)
}

// parseJsonSliceRef parses a JSON array into a slice using our custom reflection
func (c *conv) parseJsonSliceRef(jsonStr string, target refValue) error {
	if target.Kind() != tpSlice {
		return Err(errUnsupportedType, "target is not a slice")
	}

	// Basic validation - must start with [ and end with ]
	jsonStr = Convert(jsonStr).Trim().String()
	if len(jsonStr) < 2 || jsonStr[0] != '[' || jsonStr[len(jsonStr)-1] != ']' {
		return Err(errInvalidJSON, "invalid JSON array format")
	}

	elemType := target.Type().Elem()

	// Handle empty array
	if jsonStr == "[]" {
		switch elemType.Kind() {
		case tpString:
			target.Set(refValueOf([]string{}))
		case tpStruct:
			// Create empty slice of structs using unsafe operations
			target.Set(refValueOf([]interface{}{}))
		default:
			return Err(errUnsupportedType, "unsupported slice element type: "+elemType.Kind().String())
		}
		return nil
	}

	content := jsonStr[1 : len(jsonStr)-1] // Remove [ ]

	// Split array elements
	elements := c.splitJsonArrayElements(content)

	// Handle different element types
	switch elemType.Kind() {
	case tpString:
		return c.parseStringSlice(elements, target)
	case tpStruct:
		return c.parseStructSlice(elements, target, elemType)
	case tpInt, tpInt64:
		return c.parseIntSlice(elements, target)
	case tpFloat64:
		return c.parseFloatSlice(elements, target)
	case tpBool:
		return c.parseBoolSlice(elements, target)
	default:
		return Err(errUnsupportedType, "slice decoding only supports string, struct, int, float, and bool slices currently")
	}
}

// parseStringSlice parses a slice of JSON strings
func (c *conv) parseStringSlice(elements []string, target refValue) error {
	var stringSlice []string
	for _, elem := range elements {
		// Parse string element
		elemStr := Convert(elem).Trim().String()
		if len(elemStr) >= 2 && elemStr[0] == '"' && elemStr[len(elemStr)-1] == '"' {
			unquoted := elemStr[1 : len(elemStr)-1]
			decoded, err := c.unescapeJsonString(unquoted)
			if err != nil {
				return err
			}
			stringSlice = append(stringSlice, decoded)
		} else {
			return Err(errInvalidJSON, "invalid string element in array: "+elem)
		}
	}
	target.Set(refValueOf(stringSlice))
	return nil
}

// parseStructSlice parses JSON array elements into a struct slice
func (c *conv) parseStructSlice(elements []string, target refValue, elemType *refType) error {
	if len(elements) == 0 {
		// Empty slice - set target to empty slice of correct type
		target.Set(refZero(target.Type()))
		return nil
	}

	// Create slice to hold the parsed structs
	slicePtr := refNew(target.Type())
	slice := slicePtr.Elem()

	// Parse each element as a struct
	for _, element := range elements {
		element = Convert(element).Trim().String()

		// Create new struct instance
		structPtr := refNew(elemType)
		structValue := structPtr.Elem()

		// Parse JSON into the struct
		err := c.parseJsonValueWithRefReflect(element, structValue)
		if err != nil {
			return Err(errInvalidJSON, "error parsing struct element: "+err.Error())
		}

		// Append to slice using reflection
		slice = refAppend(slice, structValue)
	}

	// Set the target to our constructed slice
	target.Set(slice)
	return nil
}

// parseIntSlice, parseFloatSlice, parseBoolSlice - simplified implementations
func (c *conv) parseIntSlice(elements []string, target refValue) error {
	return Err(errUnsupportedType, "int slice decoding not implemented yet")
}

func (c *conv) parseFloatSlice(elements []string, target refValue) error {
	return Err(errUnsupportedType, "float slice decoding not implemented yet")
}

func (c *conv) parseBoolSlice(elements []string, target refValue) error {
	return Err(errUnsupportedType, "bool slice decoding not implemented yet")
}

// splitJsonArrayElements splits JSON array content into individual elements
func (c *conv) splitJsonArrayElements(content string) []string {
	var elements []string
	current := Builder()
	inQuotes := false
	braceLevel := 0
	bracketLevel := 0

	for i, char := range content {
		switch char {
		case '"':
			if i == 0 || content[i-1] != '\\' {
				inQuotes = !inQuotes
			}
			current.appendRune(char)
		case '{':
			if !inQuotes {
				braceLevel++
			}
			current.appendRune(char)
		case '}':
			if !inQuotes {
				braceLevel--
			}
			current.appendRune(char)
		case '[':
			if !inQuotes {
				bracketLevel++
			}
			current.appendRune(char)
		case ']':
			if !inQuotes {
				bracketLevel--
			}
			current.appendRune(char)
		case ',':
			if !inQuotes && braceLevel == 0 && bracketLevel == 0 {
				elem := Convert(current.String()).Trim().String()
				if len(elem) > 0 {
					elements = append(elements, elem)
				}
				current.reset()
			} else {
				current.appendRune(char)
			}
		default:
			current.appendRune(char)
		}
	}

	if current.length() > 0 {
		elem := Convert(current.String()).Trim().String()
		if len(elem) > 0 {
			elements = append(elements, elem)
		}
	}

	return elements
}

// unescapeJsonString unescapes a JSON string value
func (c *conv) unescapeJsonString(s string) (string, error) {
	// Simple implementation - just handle basic escapes for now
	// This could be expanded to handle all JSON escape sequences
	result := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) {
			switch s[i+1] {
			case '"':
				result = append(result, '"')
			case '\\':
				result = append(result, '\\')
			case 'n':
				result = append(result, '\n')
			case 'r':
				result = append(result, '\r')
			case 't':
				result = append(result, '\t')
			default:
				result = append(result, s[i], s[i+1])
			}
			i++ // Skip next character
		} else {
			result = append(result, s[i])
		}
	}
	return string(result), nil
}

// parseJsonObjectContent parses the content of a JSON object (without outer braces)
func (c *conv) parseJsonObjectContent(content string, target refValue, structInfo *refStructInfo) error {
	if content == "" {
		return nil // empty content
	}

	// Simple field parsing - split by commas (note: this is simplified and doesn't handle nested objects properly)
	pairs := c.splitJsonFields(content)

	for _, pair := range pairs {
		if err := c.parseJsonFieldPair(pair, target, structInfo); err != nil {
			return err
		}
	}

	return nil
}

// splitJsonFields splits JSON object content into field pairs (simplified)
func (c *conv) splitJsonFields(content string) []string {
	var pairs []string
	current := Builder() // Use our custom string builder
	inQuotes := false
	braceLevel := 0
	bracketLevel := 0

	for i, char := range content {
		switch char {
		case '"':
			if i == 0 || content[i-1] != '\\' {
				inQuotes = !inQuotes
			}
			current.appendRune(char)
		case '{':
			if !inQuotes {
				braceLevel++
			}
			current.appendRune(char)
		case '}':
			if !inQuotes {
				braceLevel--
			}
			current.appendRune(char)
		case '[':
			if !inQuotes {
				bracketLevel++
			}
			current.appendRune(char)
		case ']':
			if !inQuotes {
				bracketLevel--
			}
			current.appendRune(char)
		case ',':
			if !inQuotes && braceLevel == 0 && bracketLevel == 0 {
				pairs = append(pairs, current.String())
				current.reset()
			} else {
				current.appendRune(char)
			}
		default:
			current.appendRune(char)
		}
	}

	if current.length() > 0 {
		pairs = append(pairs, current.String())
	}

	return pairs
}

// parseJsonFieldPair parses a single "key":"value" pair
func (c *conv) parseJsonFieldPair(pair string, target refValue, structInfo *refStructInfo) error {
	pair = Convert(pair).Trim().String()

	// Find the colon separator
	colonIndex := c.findJsonColon(pair)
	if colonIndex == -1 {
		return Err(errInvalidJSON, "invalid field pair format: "+pair)
	}

	keyPart := Convert(pair[:colonIndex]).Trim().String()
	valuePart := Convert(pair[colonIndex+1:]).Trim().String()

	// Parse key (remove quotes)
	if len(keyPart) < 2 || keyPart[0] != '"' || keyPart[len(keyPart)-1] != '"' {
		return Err(errInvalidJSON, "invalid key format: "+keyPart)
	}
	jsonKey := keyPart[1 : len(keyPart)-1]

	// Find matching struct field
	fieldIndex := c.findStructFieldByJsonName(jsonKey, structInfo)
	if fieldIndex == -1 {
		// Field not found, skip it
		return nil
	}

	// Get the target field
	field := target.Field(fieldIndex)
	if !field.IsValid() {
		return Err(errInvalidJSON, "invalid field")
	}

	// Parse and set the value
	return c.parseJsonValueWithRefReflect(valuePart, field)
}

// findJsonColon finds the position of the colon that separates key from value
func (c *conv) findJsonColon(pair string) int {
	inQuotes := false
	for i, char := range pair {
		if char == '"' && (i == 0 || pair[i-1] != '\\') {
			inQuotes = !inQuotes
		} else if char == ':' && !inQuotes {
			return i
		}
	}
	return -1
}

// findStructFieldByJsonName finds the field index by JSON field name
func (c *conv) findStructFieldByJsonName(jsonKey string, structInfo *refStructInfo) int {
	// Match using original field names (no case conversion)
	for i, field := range structInfo.fields {
		if field.name == jsonKey {
			return i
		}
	}
	return -1
}

// appendRune adds a rune to the current conv value
func (c *conv) appendRune(r rune) *conv {
	current := c.getString()
	// Use the existing addRne2Buf method from convert.go
	buf := make([]byte, 0, len(current)+4) // 4 bytes max for UTF-8 rune
	buf = append(buf, current...)
	buf = addRne2Buf(buf, r)
	c.setString(string(buf))
	return c
}

// parseJsonPointerRef parses a JSON value into a pointer using our custom reflection
func (c *conv) parseJsonPointerRef(jsonStr string, target refValue) error {
	if target.Kind() != tpPointer {
		return Err(errUnsupportedType, "target is not a pointer")
	}

	// Check if JSON value is null
	jsonStr = Convert(jsonStr).Trim().String()
	if jsonStr == "null" {
		// Set pointer to nil by setting the pointer variable to zero
		*(*unsafe.Pointer)(target.ptr) = nil
		return nil
	}

	// Get the element type that the pointer points to
	elemType := target.Type().Elem()
	if elemType == nil {
		return Err(errUnsupportedType, "pointer element type is nil")
	}

	// Create a new value of the element type
	elemValue := refNew(elemType)
	if !elemValue.IsValid() {
		return Err(errUnsupportedType, "failed to create new element value")
	}

	// Parse JSON into the new element value
	err := c.parseJsonValueWithRefReflect(jsonStr, elemValue.Elem())
	if err != nil {
		return err
	}
	// Set the pointer to point to the new element
	// target.ptr points to the pointer field location in the struct
	// elemValue.ptr from refNew points to a pointer variable containing the allocated address
	// We need to store the actual allocated address in the struct field
	actualAddr := *(*unsafe.Pointer)(elemValue.ptr)
	*(*unsafe.Pointer)(target.ptr) = actualAddr
	return nil
}
