package tinystring

// =============================================================================
// REFLECTION FUNCTIONALITY - CORE LOGIC MOVED FROM TINYREFLECT
// =============================================================================

// GetValue returns the interface{} value stored in dataPtr
// This allows external packages to access complex types managed by tinystring
func (c *conv) GetValue() any {
	if c.dataPtr == nil {
		return nil
	}

	// Convert unsafe.Pointer back to interface{} based on Kind
	switch c.Kind {
	case KStruct:
		return *(*any)(c.dataPtr)
	case KMap:
		switch c.Kind {
		// Will be handled by specific GetMap* methods
		}
	case KSlice, KSliceStr, KByte:
		// Will be handled by specific GetSlice* methods
	case KPointer:
		return c.dataPtr // Return the pointer itself
	}

	return nil
}

// GetFieldCount returns the number of fields in a struct
// Only works with KStruct types
func (c *conv) GetFieldCount() int {
	if c.Kind != KStruct || c.dataPtr == nil {
		return 0
	}

	// Use reflection-like interface detection for struct field counting
	// This is a minimal implementation - structs are complex and need careful handling

	// For now, return 0 as we need to implement proper struct introspection
	// TODO: Implement struct field enumeration
	return 0
}

// GetField returns the value of a struct field by index
// Only works with KStruct types
func (c *conv) GetField(index int) any {
	if c.Kind != KStruct || c.dataPtr == nil {
		return nil
	}

	// TODO: Implement struct field access by index
	return nil
}

// GetTag returns the tag of a struct field by index
// Only works with KStruct types
func (c *conv) GetTag(index int, tagName string) string {
	if c.Kind != KStruct || c.dataPtr == nil {
		return ""
	}

	// TODO: Implement struct tag reading
	return ""
}

// GetMapKeys returns all keys from a map
// Only works with KMap types
func (c *conv) GetMapKeys() []any {
	if c.Kind != KMap || c.dataPtr == nil {
		return nil
	}

	var keys []any

	// Access map data through unsafe pointer - need to determine map type
	// For now, we'll implement a basic version
	// TODO: Implement proper unsafe map access based on Kind detection

	return keys
}

// GetMapValue returns a value from a map by key
// Only works with KMap types
func (c *conv) GetMapValue(key any) any {
	if c.Kind != KMap || c.dataPtr == nil {
		return nil
	}

	// TODO: Implement proper unsafe map value access
	return nil
}

// GetSliceLen returns the length of a slice
// Works with KSlice, KSliceStr, and KByte types
func (c *conv) GetSliceLen() int {
	if c.dataPtr == nil {
		return 0
	}

	switch c.Kind {
	case KSliceStr:
		// dataPtr points to &[]string
		slice := *(*[]string)(c.dataPtr)
		return len(slice)
	case KByte:
		// dataPtr points to &[]byte
		slice := *(*[]byte)(c.dataPtr)
		return len(slice)
	case KSlice:
		// For generic slice types, we need to use unsafe slice header
		// TODO: Implement if needed for other slice types
		return 0
	default:
		return 0
	}
}

// GetSliceElement returns an element from a slice by index
// Works with KSlice, KSliceStr, and KByte types
func (c *conv) GetSliceElement(index int) any {
	if c.dataPtr == nil {
		return nil
	}

	switch c.Kind {
	case KSliceStr:
		// dataPtr points to &[]string
		slice := *(*[]string)(c.dataPtr)
		if index >= 0 && index < len(slice) {
			return slice[index]
		}
		return nil
	case KByte:
		// dataPtr points to &[]byte
		slice := *(*[]byte)(c.dataPtr)
		if index >= 0 && index < len(slice) {
			return slice[index]
		}
		return nil
	case KSlice:
		// For generic slice types, we need to use unsafe slice header
		// TODO: Implement if needed for other slice types
		return nil
	default:
		return nil
	}
}

// GetSliceElementKind returns the Kind of elements in a slice
// Uses fallback logic for TinyGo/WebAssembly compatibility
func (c *conv) GetSliceElementKind() Kind {
	if c.dataPtr == nil {
		return KInvalid
	}

	switch c.Kind {
	case KSliceStr:
		return KString
	case KByte:
		return KUint8
	case KSlice:
		// TODO: Implement proper unsafe slice element kind detection
		// For now, return KInvalid until we implement proper unsafe handling
		return KInvalid
	}

	return KInvalid
}

// GetPointerElement returns the element pointed to by a pointer
// Only works with KPointer types
func (c *conv) GetPointerElement() any {
	if c.Kind != KPointer || c.dataPtr == nil {
		return nil
	}

	// TODO: Implement proper unsafe pointer dereferencing
	// For now, return nil until we implement proper unsafe pointer handling
	return nil
}

// GetPointerElementKind returns the Kind of the element pointed to by a pointer
func (c *conv) GetPointerElementKind() Kind {
	if c.Kind != KPointer || c.dataPtr == nil {
		return KInvalid
	}

	// TODO: Implement proper unsafe pointer kind detection
	// For now, return KInvalid until we implement proper unsafe pointer handling
	return KInvalid
}

// looksLikeStruct is a helper function to identify struct types
// This replaces the previous isStruct method with better logic
func (c *conv) looksLikeStruct(value any) bool {
	if value == nil {
		return false
	}

	// Explicitly reject known non-struct types
	switch value.(type) {
	// Basic types - explicitly not structs
	case string, bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return false
	// Slice types - explicitly not structs
	case []string, []bool, []int, []int8, []int16, []int32, []int64, []uint, []uint16, []uint32, []uint64, []float32, []float64, []byte:
		return false
	// Map types - explicitly not structs
	case map[string]string, map[string]any, map[string]int, map[int]string:
		return false
	// Pointer types - explicitly not structs
	case *string, *int, *bool, *float64, **string:
		return false
	// Error type - explicitly not struct
	case error:
		return false
	default:
		// For anything else that's not explicitly known, we conservatively assume it could be a struct
		// This allows structs to be processed, while unknown types will still get an error during reflection
		return true
	}
}
