package tinystring

import "unsafe"

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
	case tpInt:
		return c.encodeJsonInt()
	case tpUint:
		return c.encodeJsonUint()
	case tpFloat64:
		return c.encodeJsonFloat()
	case tpBool:
		return c.encodeJsonBool()
	case tpStrSlice:
		return c.encodeJsonStringSlice()
	case tpStruct:
		return c.encodeJsonStruct()
	case tpSlice:
		return c.encodeJsonSlice()
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
	if !c.refVal.IsValid() {
		return nil, Err(errInvalidJSON, "struct value is nil")
	}

	// Use our custom reflection to encode the struct
	return c.encodeStructValueWithRefReflect(c.refVal)
}

// encodeJsonSlice encodes a slice to JSON using reflection
func (c *conv) encodeJsonSlice() ([]byte, error) {
	if !c.refVal.IsValid() {
		return []byte("[]"), nil
	}

	// Use our custom reflection to encode the slice
	return c.encodeJsonSliceWithRefReflect(c.refVal)
}

// quoteJsonString quotes a string for JSON output with proper escaping
func (c *conv) quoteJsonString(s string) []byte {
	// Estimate capacity: original length + quotes + some escape characters
	result := make([]byte, 0, len(s)+16)
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

// encodeStructValueWithRefReflect encodes a struct using refValue directly
func (c *conv) encodeStructValueWithRefReflect(rv refValue) ([]byte, error) {
	// Handle pointer to struct
	if rv.Kind() == tpPointer {
		elem := rv.Elem()
		if !elem.IsValid() {
			return []byte("null"), nil
		}
		rv = elem
	}

	if rv.Kind() != tpStruct {
		return nil, Err(errUnsupportedType, "not a struct")
	}

	result := make([]byte, 0, 256)
	result = append(result, '{')

	fieldCount := 0
	numFields := rv.NumField()

	for i := range numFields {
		field := rv.Field(i)

		// Skip invalid fields
		if !field.IsValid() {
			continue
		} // Get field name from struct info - use original field name
		var structInfo refStructInfo
		getStructInfo(rv.Type(), &structInfo)
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
		result = append(result, ':')

		// Encode field value using our custom reflection
		fieldJson, err := c.encodeFieldValueWithRefReflect(field)
		if err != nil {
			return nil, err
		}
		result = append(result, fieldJson...)
		fieldCount++
	}

	result = append(result, '}')
	return result, nil
}

// encodeFieldValueWithRefReflect encodes a field value using our custom reflection
func (c *conv) encodeFieldValueWithRefReflect(v refValue) ([]byte, error) {
	switch v.Kind() {
	case tpString:
		return c.quoteJsonString(v.String()), nil
	case tpInt, tpInt8, tpInt16, tpInt32, tpInt64:
		tempConv := Convert(v.Int())
		return tempConv.encodeJsonInt()
	case tpUint, tpUint8, tpUint16, tpUint32, tpUint64:
		tempConv := Convert(v.Uint())
		return tempConv.encodeJsonUint()
	case tpFloat32, tpFloat64:
		tempConv := Convert(v.Float())
		return tempConv.encodeJsonFloat()
	case tpBool:
		tempConv := Convert(v.Bool())
		return tempConv.encodeJsonBool()
	case tpStruct:
		return c.encodeStructValueWithRefReflect(v)
	case tpSlice:
		return c.encodeJsonSliceWithRefReflect(v)
	case tpMap, tpChan, tpFunc, tpUnsafePointer:
		return nil, Err(errUnsupportedType, "unsupported type for JSON encoding: "+v.Kind().String())
	case tpPointer:
		elem := v.Elem()
		if !elem.IsValid() {
			return []byte("null"), nil
		}
		return c.encodeFieldValueWithRefReflect(elem)
	default:
		return []byte("null"), nil
	}
}

// encodeJsonSliceWithRefReflect encodes a slice using our custom reflection
func (c *conv) encodeJsonSliceWithRefReflect(v refValue) ([]byte, error) {
	if v.Kind() != tpSlice {
		return nil, Err(errUnsupportedType, "not a slice")
	}

	length := v.Len()
	if length == 0 {
		return []byte("[]"), nil
	}

	result := make([]byte, 0, 256)
	result = append(result, '[')

	for i := 0; i < length; i++ {
		if i > 0 {
			result = append(result, ',')
		}

		// Get the element at index i
		elem := v.Index(i)

		// Encode the element
		elemJson, err := c.encodeFieldValueWithRefReflect(elem)
		if err != nil {
			return nil, err
		}

		result = append(result, elemJson...)
	}

	result = append(result, ']')
	return result, nil
}

// Len returns the length of v if v is a slice, array, map, string, or channel
func (v refValue) Len() int {
	switch v.Kind() {
	case tpSlice:
		return (*refSliceHeader)(v.ptr).Len
	case tpString:
		ptr := v.ptr
		if v.flag&flagIndir != 0 {
			ptr = *(*unsafe.Pointer)(ptr)
		}
		return len(*(*string)(ptr))
	default:
		panic("reflect: call of reflect.Value.Len on " + v.Kind().String() + " value")
	}
}

// Index returns v's i'th element. It panics if v's Kind is not Array, Slice, or String or i is out of range.
func (v refValue) Index(i int) refValue {
	switch k := v.Kind(); k {
	case tpSlice:
		s := (*refSliceHeader)(v.ptr)
		if i < 0 || i >= s.Len {
			panic("reflect: slice index out of range")
		}

		// Get element type from slice type
		elemType := v.typ.Elem()
		if elemType == nil {
			panic("reflect: slice element type is nil")
		}

		// Calculate element address
		elemPtr := add(s.Data, uintptr(i)*elemType.Size(), "same as &s[i]")
		// Create flags for the element
		// Elements in slices are accessed directly, no flagIndir needed
		fl := v.flag&flagRO | flagAddr | refFlag(elemType.Kind())

		return refValue{elemType, elemPtr, fl}
	default:
		panic("reflect: call of reflect.Value.Index on " + k.String() + " value")
	}
}
