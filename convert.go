package tinystring

// Generic type interfaces for consolidating repetitive type switches
type anyInt interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type anyUint interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type anyFloat interface {
	~float32 | ~float64
}

type conv struct {
	// Phase 13.3: Optimized memory allocation - all buffers are dynamic []byte
	// Buffers with initial capacity 64, grow as needed (no truncation)
	out     []byte // Buffer principal - make([]byte, 0, 64)
	outLen  int    // Longitud actual en out
	work    []byte // Buffer temporal - make([]byte, 0, 64)
	workLen int    // Longitud actual en work
	err     []byte // Buffer de errores - make([]byte, 0, 64)
	errLen  int    // Longitud actual en err

	// Type indicator - most frequently accessed
	kind kind // Hot path: type checking

	// PHASE 13.3: Variables ELIMINADAS completamente:
	// err       string ❌ DEPRECATED → Reemplazado por err []byte
	// stringVal string ❌ DEPRECATED → Reemplazado por out []byte
	// tmpStr    string ❌ DEPRECATED → Reemplazado por work []byte

	// Numeric values grouped together
	intVal   int64
	uintVal  uint64
	floatVal float64

	// Less frequently used fields
	stringSliceVal []string
	stringPtrVal   *string
	boolVal        bool
}

// Convert initializes a new conv struct with optional value for string,bool and number manipulation.
// REFACTORED: Now accepts variadic parameters - Convert() or Convert(value)
// Phase 7: Uses object pool internally for memory optimization (transparent to user)
func Convert(v ...any) *conv {
	c := getConv()
	// Validation: Only accept 0 or 1 parameter
	if len(v) > 1 {
		return c.setErr(D.Invalid, D.Number, D.Of, D.Argument) // Consistent error handling pattern
	}
	// Initialize with value if provided, empty otherwise
	if len(v) == 1 {
		// Inlined withValue logic for performance
		val := v[0]
		if val == nil {
			return c.setErr(D.String, D.Empty)
		} else {
			switch typedVal := val.(type) {
			case string:
				// Store string content directly in out
				c.out = append(c.out[:0], typedVal...)
				c.outLen = len(typedVal)
				c.kind = KString
			case []string:
				c.stringSliceVal = typedVal
				c.kind = KSliceStr
			case *string:
				// Store string content directly in out
				c.out = append(c.out[:0], (*typedVal)...)
				c.outLen = len(*typedVal)
				c.stringPtrVal = typedVal
				c.kind = KPointer
			case bool:
				c.boolVal = typedVal
				c.kind = KBool
			case error:
				return c.setErr(typedVal.Error())
			default:
				// Handle numeric types using generics
				c.handleAnyType(typedVal)
			}
		}
	}
	// If no value provided, conv is ready for builder pattern

	return c
}

// Consolidated generic functions - single function per type with operation parameter
func genInt[T anyInt](c *conv, v T, op int) {
	c.intVal = int64(v)
	switch op {
	case 0:
		c.kind = KInt // setValue
	case 1:
		c.fmtIntGeneric(c.intVal, 10, true) // any2s
	case 2:
		c.intToBufTmp() // format
	}
}

func genUint[T anyUint](c *conv, v T, op int) {
	c.uintVal = uint64(v)
	switch op {
	case 0:
		c.kind = KUint // setValue
	case 1:
		c.fmtIntGeneric(int64(c.uintVal), 10, false) // any2s
	case 2:
		c.uint64ToBufTmp() // format
	}
}

func genFloat[T anyFloat](c *conv, v T, op int) {
	c.floatVal = float64(v)
	switch op {
	case 0:
		c.kind = KFloat64 // setValue
	case 1:
		c.floatToBufTmp() // any2s
	case 2:
		c.f2sMan(-1) // format
	}
}

// Apply updates the original string pointer with the current content and auto-releases to pool.
// This method should be used when you want to modify the original string directly
// without additional allocations.
func (t *conv) Apply() {
	if t.kind == KPointer && t.stringPtrVal != nil {
		*t.stringPtrVal = t.getString()
	}
	// Auto-release back to pool for memory efficiency
	t.putConv()
}

// String method to return the content of the conv and automatically returns object to pool
// Phase 7: Auto-release makes pool usage completely transparent to user
func (t *conv) String() string {
	out := t.getString()
	// Auto-release back to pool for memory efficiency
	t.putConv()
	return out
}

// getString converts the current value to string only when needed
// BUILDER INTEGRATION: Prioritizes buffer content when available
// Phase 13.2: Optimized to eliminate repeated string(t.out) allocations
func (t *conv) getString() string {
	if t.kind == KErr {
		return ""
	}

	// BUILDER PRIORITY: If buffer has content, use it as source of truth
	// Phase 13.2: Cache the buffer conversion using work to avoid repeated allocations
	if t.outLen > 0 {
		// Check if work needs to be updated (comparing lengths for efficiency)
		if t.workLen != t.outLen {
			// Copy out content to work and update string representation
			t.work = append(t.work[:0], t.out[:t.outLen]...)
			t.workLen = t.outLen
		}
		return string(t.work[:t.workLen])
	}

	// For string pointers, always return the current value (don't use cache)
	if t.kind == KPointer && t.stringPtrVal != nil {
		return *t.stringPtrVal
	}

	// Use cached string in work if available (simplified cache logic)
	if t.workLen > 0 && t.kind != KPointer {
		return string(t.work[:t.workLen])
	}

	// Convert to string using internal methods to avoid allocations
	switch t.kind {
	case KString:
		// String content should already be in out from Convert()
		if t.outLen > 0 {
			return string(t.out[:t.outLen])
		}
		return ""
	case KPointer: // now receive any pointer check please
		// For string pointers, always get the current value from the pointer
		if t.stringPtrVal != nil {
			value := *t.stringPtrVal
			// Store in work for caching
			t.work = append(t.work[:0], value...)
			t.workLen = len(value)
			return value
		}
		return ""
	case KSliceStr: // now receive any alice check please
		if len(t.stringSliceVal) == 0 {
			t.work = t.work[:0]
			t.workLen = 0
		} else if len(t.stringSliceVal) == 1 {
			value := t.stringSliceVal[0]
			t.work = append(t.work[:0], value...)
			t.workLen = len(value)
		} else {
			// Use builder API for zero-allocation string construction
			c := Convert() // Empty initialization for builder pattern
			for i, s := range t.stringSliceVal {
				if i > 0 {
					c.Write(" ")
				}
				c.Write(s)
			}
			out := c.String() // Auto-release to pool
			t.work = append(t.work[:0], out...)
			t.workLen = len(out)
		}
		return string(t.work[:t.workLen])
	case KInt:
		// Use internal method instead of external function
		t.fmtIntGeneric(t.intVal, 10, true)
	case KUint:
		// Use internal method instead of external function
		t.fmtIntGeneric(int64(t.uintVal), 10, false)
	case KFloat64:
		// Use internal method instead of external function
		t.floatToBufTmp()
	case KBool:
		var value string
		if t.boolVal {
			value = trueStr
		} else {
			value = falseStr
		}
		// Store in work for caching
		t.work = append(t.work[:0], value...)
		t.workLen = len(value)
		return value
	default:
		t.work = t.work[:0]
		t.workLen = 0
	}
	// Cache is now updated (simplified without type tracking)
	return string(t.work[:t.workLen])
}

// addRne2Buf manually encodes a rune to UTF-8 and appends it to the byte slice.
// This avoids importing the unicode/utf8 package for size optimization.
func addRne2Buf(out []byte, r rune) []byte {
	if r < 0x80 {
		return append(out, byte(r))
	} else if r < 0x800 {
		return append(out, byte(0xC0|(r>>6)), byte(0x80|(r&0x3F)))
	} else if r < 0x10000 {
		return append(out, byte(0xE0|(r>>12)), byte(0x80|((r>>6)&0x3F)), byte(0x80|(r&0x3F)))
	} else {
		return append(out, byte(0xF0|(r>>18)), byte(0x80|((r>>12)&0x3F)), byte(0x80|((r>>6)&0x3F)), byte(0x80|(r&0x3F)))
	}
}

// setString converts to string type and stores the value
func (t *conv) setString(s string) {
	// Store string content directly in out
	t.out = append(t.out[:0], s...)
	t.outLen = len(s)

	// If working with string pointer, update the original string
	if t.kind == KPointer && t.stringPtrVal != nil {
		*t.stringPtrVal = s
		// Keep the kind as stringPtr to maintain the pointer relationship
	} else {
		t.kind = KString
	}

	// Clear other values to save memory
	t.intVal = 0
	t.uintVal = 0
	t.floatVal = 0
	t.boolVal = false
	t.stringSliceVal = nil
	// Invalidate cache since we changed the string
	t.work = t.work[:0]
	t.workLen = 0
}

// Internal conversion methods - centralized in conv to minimize allocations
// These methods modify the conv struct directly instead of returning values

// any2s converts any type to string and stores in tmpStr
// default set to "" if no conversion is possible
// supports int, uint, float, bool, string and error types
func (t *conv) any2s(v any) {
	switch val := v.(type) {
	case error:
		t.setErr(val.Error())
	case string:
		// Store string content directly in out
		t.out = append(t.out[:0], val...)
		t.outLen = len(val)
	case bool:
		var out string
		if val {
			out = trueStr
		} else {
			out = falseStr
		}
		// Store boolean out in out
		t.out = append(t.out[:0], out...)
		t.outLen = len(out)
	default:
		// Inline handleAnyTypeForAny2s logic for numeric types
		switch v := v.(type) {
		case int:
			genInt(t, v, 1)
		case int8:
			genInt(t, v, 1)
		case int16:
			genInt(t, v, 1)
		case int32:
			genInt(t, v, 1)
		case int64:
			genInt(t, v, 1)
		case uint:
			genUint(t, v, 1)
		case uint8:
			genUint(t, v, 1)
		case uint16:
			genUint(t, v, 1)
		case uint32:
			genUint(t, v, 1)
		case uint64:
			genUint(t, v, 1)
		case float32:
			genFloat(t, v, 1)
		case float64:
			genFloat(t, v, 1)
		default:
			// Clear buffer for unknown types
			t.out = t.out[:0]
			t.outLen = 0
		}
	}
}

// handleIntType processes integer types using generics
// handleAnyType processes any numeric type using generics
func (c *conv) handleAnyType(val any) {
	switch v := val.(type) {
	case int:
		genInt(c, v, 0)
	case int8:
		genInt(c, v, 0)
	case int16:
		genInt(c, v, 0)
	case int32:
		genInt(c, v, 0)
	case int64:
		genInt(c, v, 0)
	case uint:
		genUint(c, v, 0)
	case uint8:
		genUint(c, v, 0)
	case uint16:
		genUint(c, v, 0)
	case uint32:
		genUint(c, v, 0)
	case uint64:
		genUint(c, v, 0)
	case float32:
		genFloat(c, v, 0)
	case float64:
		genFloat(c, v, 0)
	}
}
