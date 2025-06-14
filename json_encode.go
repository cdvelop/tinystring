package tinystring

// JSON encoding implementation for TinyString
// Uses our custom reflectlite integration for minimal binary size

// writer interface for JSON output - private interface compatible with io.Writer
// This allows writing JSON directly to any output that implements Write method
// without importing io package to maintain minimal binary size
type writer interface {
	Write(p []byte) (n int, err error)
}

// JsonEncode converts the current value to JSON format
//
// Usage patterns:
//   bytes, err := Convert(&user).JsonEncode()           // Returns JSON as []byte
//   err := Convert(&user).JsonEncode(writer)           // Writes JSON to writer, returns nil bytes
//   err := Convert(&user).JsonEncode(httpResponseWriter) // Direct HTTP response
//   err := Convert(&user).JsonEncode(buffer)           // To buffer/file
//
// The method accepts optional writer implementing Write([]byte) (int, error):
// - Without writer: Returns ([]byte, error) with JSON content
// - With writer: Writes to writer and returns (nil, error)
//
// Supported types for JSON encoding:
// - Basic types: string, int64, uint64, float64, bool
// - Slices: []string, []int, []float64, []bool
// - Structs: with basic field types and nested structs (max 8 levels)
// - Struct slices: []User, []Address, etc.
//
// Field naming: Automatically converts to snake_case (UserName -> "user_name")
// No JSON tags required - uses reflection for field inspection
func (c *conv) JsonEncode(w ...writer) ([]byte, error) {
	// Check if writer is provided
	if len(w) > 0 && w[0] != nil {
		// Write to provided writer
		jsonBytes, err := c.generateJsonBytes()
		if err != nil {
			return nil, err
		}

		_, writeErr := w[0].Write(jsonBytes)
		return nil, writeErr
	}

	// No writer provided, return bytes directly
	return c.generateJsonBytes()
}

// generateJsonBytes creates JSON representation of the current value
func (c *conv) generateJsonBytes() ([]byte, error) {
	switch c.vTpe {
	case tpString:
		return c.encodeJsonString()
	case tpInt, tpInt8, tpInt16, tpInt32, tpInt64:
		return c.encodeJsonInt()
	case tpUint, tpUint8, tpUint16, tpUint32, tpUint64:
		return c.encodeJsonUint()
	case tpFloat32, tpFloat64:
		return c.encodeJsonFloat()
	case tpBool:
		return c.encodeJsonBool()
	case tpStrSlice:
		return c.encodeJsonStringSlice()
	case tpStruct:
		return c.encodeJsonStruct()
	case tpSlice:
		return c.encodeJsonSlice()
	case tpPointer:
		return c.encodeJsonPointer()
	default:
		return nil, Err(errUnsupportedType, "for JSON encoding")
	}
}

// encodeJsonString encodes a string value to JSON
func (c *conv) encodeJsonString() ([]byte, error) {
	str := c.getString()
	return c.quoteJsonString(str), nil
}

// encodeJsonInt encodes an integer value to JSON
func (c *conv) encodeJsonInt() ([]byte, error) {
	// Use existing tinystring int formatting
	c.fmtInt(10)
	return []byte(c.tmpStr), nil
}

// encodeJsonUint encodes an unsigned integer value to JSON
func (c *conv) encodeJsonUint() ([]byte, error) {
	// Use existing tinystring uint formatting
	c.fmtUint(10)
	return []byte(c.tmpStr), nil
}

// encodeJsonFloat encodes a float value to JSON
func (c *conv) encodeJsonFloat() ([]byte, error) {
	// Use existing tinystring float formatting
	c.f2s()
	return []byte(c.tmpStr), nil
}

// encodeJsonBool encodes a boolean value to JSON
func (c *conv) encodeJsonBool() ([]byte, error) {
	if c.getBool() {
		return []byte("true"), nil
	}
	return []byte("false"), nil
}

// encodeJsonStringSlice encodes a string slice to JSON
func (c *conv) encodeJsonStringSlice() ([]byte, error) {
	if len(c.stringSliceVal) == 0 {
		return []byte("[]"), nil
	}

	result := make([]byte, 0, len(c.stringSliceVal)*20) // Estimate capacity
	result = append(result, '[')

	for i, str := range c.stringSliceVal {
		if i > 0 {
			result = append(result, ',')
		}
		quoted := c.quoteJsonString(str)
		result = append(result, quoted...)
	}

	result = append(result, ']')
	return result, nil
}

// encodeJsonStruct encodes a struct to JSON using reflection
func (c *conv) encodeJsonStruct() ([]byte, error) {
	if !c.refIsValid() {
		return nil, Err(errInvalidJSON, "struct value is nil")
	}

	// Use our custom reflection to encode the struct directly
	return c.encodeStructValueWithConvReflect()
}

// encodeJsonSlice encodes a slice to JSON using reflection
func (c *conv) encodeJsonSlice() ([]byte, error) {
	if !c.refIsValid() {
		return []byte("[]"), nil
	}

	if c.refKind() != tpSlice {
		return []byte("[]"), nil
	}

	length := c.refLen()
	if length == 0 {
		return []byte("[]"), nil
	}

	result := make([]byte, 0, 256)
	result = append(result, '[')

	for i := range length {
		if i > 0 {
			result = append(result, ',')
		}

		// Get element at index i
		elem := c.refIndex(i)
		if !elem.refIsValid() {
			result = append(result, []byte("null")...)
			continue
		}

		// Encode the element recursively
		var elemBytes []byte
		var err error

		switch elem.refKind() {
		case tpString:
			strVal := elem.refString()
			elemBytes = c.quoteJsonString(strVal)
		case tpInt, tpInt8, tpInt16, tpInt32, tpInt64:
			intVal := elem.refInt()
			tempConv := newConv(nil)
			if tempConv.intToJsonString(intVal) {
				elemBytes = []byte(tempConv.tmpStr)
			} else {
				elemBytes = []byte("0")
			}
		case tpUint, tpUint8, tpUint16, tpUint32, tpUint64:
			uintVal := elem.refUint()
			tempConv := newConv(nil)
			if tempConv.uintToJsonString(uintVal) {
				elemBytes = []byte(tempConv.tmpStr)
			} else {
				elemBytes = []byte("0")
			}
		case tpFloat32, tpFloat64:
			floatVal := elem.refFloat()
			tempConv := newConv(nil)
			if tempConv.floatToJsonString(floatVal) {
				elemBytes = []byte(tempConv.tmpStr)
			} else {
				elemBytes = []byte("0")
			}
		case tpBool:
			boolVal := elem.refBool()
			if boolVal {
				elemBytes = []byte("true")
			} else {
				elemBytes = []byte("false")
			}
		case tpStruct:
			// Handle struct elements recursively
			elemBytes, err = elem.encodeStructValueWithConvReflect()
			if err != nil {
				elemBytes = []byte("{}")
			}
		case tpSlice:
			// Handle nested slices recursively
			elemBytes, err = elem.encodeJsonSlice()
			if err != nil {
				elemBytes = []byte("[]")
			}
		case tpPointer:
			// Handle pointers by dereferencing
			elemPtr := elem.refElem()
			if !elemPtr.refIsValid() {
				elemBytes = []byte("null")
			} else {
				// Recursively call slice encoding with the dereferenced element
				switch elemPtr.refKind() {
				case tpStruct:
					elemBytes, err = elemPtr.encodeStructValueWithConvReflect()
					if err != nil {
						elemBytes = []byte("{}")
					}
				case tpSlice:
					elemBytes, err = elemPtr.encodeJsonSlice()
					if err != nil {
						elemBytes = []byte("[]")
					}
				default:
					// For basic types, encode directly
					tempConv := newConv(nil)
					if tempConv.encodeFieldValueToJson(elemPtr) {
						elemBytes = []byte(tempConv.tmpStr)
					} else {
						elemBytes = []byte("null")
					}
				}
			}
		default:
			elemBytes = []byte("null")
		}

		result = append(result, elemBytes...)
	}

	result = append(result, ']')
	return result, nil
}

// encodeJsonPointer encodes a pointer value to JSON
func (c *conv) encodeJsonPointer() ([]byte, error) {
	// Handle nil pointer
	if c.ptr == nil {
		return []byte("null"), nil // Case 1: ptr is nil
	}

	// The current conv already represents the pointer, we need to get the element it points to
	if c.refKind() != tpPointer {
		return []byte("null"), nil // Case 2: not a pointer kind
	}

	// Get the element that the pointer points to using existing reflection
	elem := c.refElem()
	if !elem.refIsValid() {
		return []byte("null"), nil // Case 3: element not valid
	}

	// Create a new conv for the pointed-to value and encode it
	elemValue := elem.Interface()
	if elemValue == nil {
		return []byte("null"), nil // Case 4: element interface is nil
	}

	elemConv := Convert(elemValue)
	return elemConv.generateJsonBytes() // Case 5: should work
}

// quoteJsonString quotes a string for JSON output with proper escaping
func (c *conv) quoteJsonString(s string) []byte {
	// Add safety check for string length
	sLen := len(s)
	if sLen < 0 || sLen > 1<<20 { // 1MB limit for safety
		return []byte(`""`)
	}

	// Estimate capacity: original length + quotes + some escape characters
	result := make([]byte, 0, sLen+16)
	result = append(result, '"')

	for _, r := range s {
		switch r {
		case '"':
			result = append(result, '\\', '"')
		case '\\':
			result = append(result, '\\', '\\')
		case '\b':
			result = append(result, '\\', 'b')
		case '\f':
			result = append(result, '\\', 'f')
		case '\n':
			result = append(result, '\\', 'n')
		case '\r':
			result = append(result, '\\', 'r')
		case '\t':
			result = append(result, '\\', 't')
		default:
			if r < 32 {
				// Control characters need unicode escaping
				result = append(result, '\\', 'u', '0', '0')
				if r < 16 {
					result = append(result, '0')
				} else {
					result = append(result, '1')
					r -= 16
				}
				if r < 10 {
					result = append(result, byte('0'+r))
				} else {
					result = append(result, byte('a'+r-10))
				}
			} else {
				// Add the rune as UTF-8
				var buf [4]byte
				n := len(string(r))
				copy(buf[:], string(r))
				result = append(result, buf[:n]...)
			}
		}
	}

	result = append(result, '"')
	return result
}

// escapeAndQuoteJsonString escapes and quotes a string for JSON without heap allocation
// Stores result directly in c.tmpStr
func (c *conv) escapeAndQuoteJsonString(s string) {
	// Use fixed buffer to avoid heap allocation
	var buf [512]byte // Fixed size buffer for most strings
	idx := 0

	// Add opening quote
	if idx < len(buf) {
		buf[idx] = '"'
		idx++
	}

	// Escape and copy characters
	for _, r := range s {
		if idx >= len(buf)-6 { // Reserve space for closing quote and escape sequences
			break
		}

		switch r {
		case '"':
			buf[idx] = '\\'
			buf[idx+1] = '"'
			idx += 2
		case '\\':
			buf[idx] = '\\'
			buf[idx+1] = '\\'
			idx += 2
		case '\b':
			buf[idx] = '\\'
			buf[idx+1] = 'b'
			idx += 2
		case '\f':
			buf[idx] = '\\'
			buf[idx+1] = 'f'
			idx += 2
		case '\n':
			buf[idx] = '\\'
			buf[idx+1] = 'n'
			idx += 2
		case '\r':
			buf[idx] = '\\'
			buf[idx+1] = 'r'
			idx += 2
		case '\t':
			buf[idx] = '\\'
			buf[idx+1] = 't'
			idx += 2
		default:
			if r < 32 {
				// Control characters need unicode escaping \u00XX
				buf[idx] = '\\'
				buf[idx+1] = 'u'
				buf[idx+2] = '0'
				buf[idx+3] = '0'
				if r < 16 {
					buf[idx+4] = '0'
				} else {
					buf[idx+4] = '1'
					r -= 16
				}
				if r < 10 {
					buf[idx+5] = byte('0' + r)
				} else {
					buf[idx+5] = byte('A' + r - 10)
				}
				idx += 6
			} else {
				// Regular character - convert rune to UTF-8 bytes
				if r < 128 {
					buf[idx] = byte(r)
					idx++
				} else {
					// For non-ASCII, simplified handling
					buf[idx] = '?' // Placeholder for complex UTF-8
					idx++
				}
			}
		}
	}

	// Add closing quote
	if idx < len(buf) {
		buf[idx] = '"'
		idx++
	}

	// Copy to tmpStr
	c.tmpStr = string(buf[:idx])
}

// encodeStructValueWithConvReflect encodes a struct using conv directly
func (c *conv) encodeStructValueWithConvReflect() ([]byte, error) {
	// Handle pointer to struct
	if c.refKind() == tpPointer {
		elem := c.refElem()
		if !elem.refIsValid() {
			return []byte("null"), nil
		}
		c = elem
	}

	if c.refKind() != tpStruct {
		return nil, Err(errUnsupportedType, "not a struct")
	}

	result := make([]byte, 0, 256)
	result = append(result, '{')
	fieldCount := 0
	numFields := c.refNumField()

	for i := range numFields {
		field := c.refField(i)

		// Skip invalid fields
		if !field.refIsValid() {
			continue
		}

		// Get field name from struct info - use original field name
		var structInfo refStructInfo
		getStructInfo(c.Type(), &structInfo)
		if structInfo.refType == nil || i >= len(structInfo.fields) {
			continue
		}

		jsonKey := structInfo.fields[i].name

		// Add comma separator for subsequent fields
		if fieldCount > 0 {
			result = append(result, ',')
		}

		// Add field name as quoted JSON key
		quotedKey := c.quoteJsonString(jsonKey)
		result = append(result, quotedKey...)
		result = append(result, ':') // Encode field value using our custom reflection
		if !c.encodeFieldValueToJson(field) {
			return nil, c
		}
		fieldJson := c.tmpStr
		result = append(result, fieldJson...)
		fieldCount++
	}

	result = append(result, '}')
	return result, nil
}

// encodeFieldValueToJson encodes a field value to JSON without heap allocation
// Stores result in c.tmpStr and returns success status
func (c *conv) encodeFieldValueToJson(fieldValue *conv) bool {
	if fieldValue == nil || !fieldValue.refIsValid() {
		c.tmpStr = "null"
		return true
	}

	switch fieldValue.refKind() {
	case tpString:
		strVal := fieldValue.refString() // Quote the string and store in tmpStr without heap allocation
		c.escapeAndQuoteJsonString(strVal)
		return true

	case tpInt, tpInt8, tpInt16, tpInt32, tpInt64:
		intVal := fieldValue.refInt()
		return c.intToJsonString(intVal)

	case tpUint, tpUint8, tpUint16, tpUint32, tpUint64:
		uintVal := fieldValue.refUint()
		return c.uintToJsonString(uintVal)

	case tpFloat32, tpFloat64:
		floatVal := fieldValue.refFloat()
		return c.floatToJsonString(floatVal)

	case tpBool:
		boolVal := fieldValue.refBool()
		if boolVal {
			c.tmpStr = "true"
		} else {
			c.tmpStr = "false"
		}
		return true
	case tpSlice:
		// Handle slices recursively by using reflection
		// Create temporary result and call existing slice encoding
		tempResult, err := fieldValue.encodeJsonSlice()
		if err != nil {
			c.tmpStr = "[]"
		} else {
			c.tmpStr = string(tempResult)
		}
		return true

	case tpStruct:
		// Handle nested structs recursively
		tempResult, err := fieldValue.encodeStructValueWithConvReflect()
		if err != nil {
			c.tmpStr = "{}"
		} else {
			c.tmpStr = string(tempResult)
		}
		return true

	case tpPointer:
		// Handle pointers by dereferencing
		elem := fieldValue.refElem()
		if !elem.refIsValid() {
			c.tmpStr = "null"
			return true
		}
		return c.encodeFieldValueToJson(elem)
	default:
		c.err = errUnsupportedType
		c.tmpStr = "null"
		return false
	}
}
