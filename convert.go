package tinystring

// Buffer destination selection for anyToBuff universal conversion function
type buffDest int

const (
	buffOut  buffDest = iota // Primary output buffer
	buffWork                 // Working/temporary buffer
	buffErr                  // Error message buffer
)

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

	// ✅ UNIFIED BUFFER ARCHITECTURE - Only essential fields remain
	pointerVal any // ✅ Universal pointer for complex types ([]string, map[string]any, etc.)

	// TEMPORARY: Keep these fields until full migration to anyToBuff() completed
	intVal         int64    // TODO: eliminate after anyToBuff() migration
	uintVal        uint64   // TODO: eliminate after anyToBuff() migration
	floatVal       float64  // TODO: eliminate after anyToBuff() migration
	boolVal        bool     // TODO: eliminate after anyToBuff() migration
	stringSliceVal []string // TODO: eliminate after anyToBuff() migration
}

// Convert initializes a new conv struct with optional value for string,bool and number manipulation.
// REFACTORED: Now accepts variadic parameters - Convert() or Convert(value)
// Phase 7: Uses object pool internally for memory optimization (transparent to user)
func Convert(v ...any) *conv {
	c := getConv()
	// Validation: Only accept 0 or 1 parameter
	if len(v) > 1 {
		return c.wrErr(D.Invalid, D.Number, D.Of, D.Argument) // Consistent error handling pattern
	}
	// Initialize with value if provided, empty otherwise
	if len(v) == 1 {
		// Inlined withValue logic for performance
		val := v[0]
		if val == nil {
			return c.wrErr(D.String, D.Empty)
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
				c.pointerVal = typedVal
				c.kind = KPointer
			case bool:
				// Use anyToBuff() for immediate conversion instead of storing in field
				anyToBuff(c, buffOut, typedVal)
				c.kind = KBool
			case error:
				return c.wrErr(typedVal.Error())
			default:
				// Handle numeric types using generics
				c.handleAnyType(typedVal)
			}
		}
	}
	// If no value provided, conv is ready for builder pattern

	return c
}

// =============================================================================
// CENTRALIZED CONVERSION HELPER METHODS
// Implementation of type-specific conversions using centralized buffer ops
// =============================================================================

// fmtIntToOut converts integer to string and writes to out buffer
// Replaces fmtIntGeneric() with centralized buffer management
func (c *conv) fmtIntToOut(val int64, base int, signed bool) {
	if val == 0 {
		c.wrStringToOut("0")
		return
	}
	// Handle negative numbers for signed integers
	if signed && val < 0 {
		val = -val
		c.wrStringToOut("-")
	}

	// Convert using existing manual implementation logic
	// Use work buffer for intermediate operations
	c.rstWork()

	// Build digits in reverse order in work buffer
	for val > 0 {
		digit := byte(val%int64(base)) + '0'
		if digit > '9' {
			digit += 'a' - '9' - 1
		}
		c.work = append(c.work, digit)
		c.workLen++
		val /= int64(base)
	}

	// Reverse and write to out buffer
	for i := c.workLen - 1; i >= 0; i-- {
		c.writeByte(c.work[i])
	}
}

// floatToOut converts float64 to string and writes to out buffer
// Replaces floatToBufTmp() with centralized buffer management
func (c *conv) floatToOut() {
	val := c.floatVal

	// Handle special cases
	if val != val { // NaN
		c.wrStringToOut("NaN")
		return
	}
	if val == 0 {
		c.wrStringToOut("0")
		return
	}

	// Handle infinity
	if val > 1.7976931348623157e+308 {
		c.wrStringToOut("+Inf")
		return
	}
	if val < -1.7976931348623157e+308 {
		c.wrStringToOut("-Inf")
		return
	}

	// Handle negative numbers
	if val < 0 {
		c.wrStringToOut("-")
		val = -val
	}

	// Use existing floatToStringOptimized logic but write to out buffer
	// For now, use a simplified version - can be optimized later
	c.rstWork()

	// Integer part
	intPart := int64(val)
	c.fmtIntToWork(intPart, 10, false)

	// Copy integer part from work to out
	for i := 0; i < c.workLen; i++ {
		c.writeByte(c.work[i])
	}

	// Fractional part
	fracPart := val - float64(intPart)
	if fracPart > 0 {
		c.wrStringToOut(".")
		// Simple fractional conversion (can be optimized later)
		for i := 0; i < 6 && fracPart > 0; i++ {
			fracPart *= 10
			digit := int(fracPart)
			c.writeByte(byte(digit) + '0')
			fracPart -= float64(digit)
		}
	}
}

// fmtIntToWork converts integer to work buffer (helper for floatToOut)
func (c *conv) fmtIntToWork(val int64, base int, signed bool) {
	if val == 0 {
		c.work = append(c.work, '0')
		c.workLen++
		return
	}

	// Handle negative numbers for signed integers
	if signed && val < 0 {
		c.work = append(c.work, '-')
		c.workLen++
		val = -val
	}

	// Store starting position for reversal
	start := c.workLen

	// Build digits in reverse order
	for val > 0 {
		digit := byte(val%int64(base)) + '0'
		if digit > '9' {
			digit += 'a' - '9' - 1
		}
		c.work = append(c.work, digit)
		c.workLen++
		val /= int64(base)
	}

	// Reverse the digits portion
	end := c.workLen - 1
	for start < end {
		c.work[start], c.work[end] = c.work[end], c.work[start]
		start++
		end--
	}
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
	if t.kind == KPointer && t.pointerVal != nil {
		// Type assert to *string for Apply() functionality
		if strPtr, ok := t.pointerVal.(*string); ok {
			*strPtr = t.ensureStringInOut()
		}
	}
	// Auto-release back to pool for memory efficiency
	t.putConv()
}

// String method to return the content of the conv and automatically returns object to pool
// Phase 7: Auto-release makes pool usage completely transparent to user
func (t *conv) String() string {
	out := t.ensureStringInOut()
	// Auto-release back to pool for memory efficiency
	t.putConv()
	return out
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
func (t *conv) setString(s string) { // Store string content directly in out using API
	t.rstOut()         // Clear buffer using API
	t.wrStringToOut(s) // Write using API
	// If working with string pointer, update the original string
	if t.kind == KPointer && t.pointerVal != nil {
		// Type assert to *string for pointer functionality
		if strPtr, ok := t.pointerVal.(*string); ok {
			*strPtr = s
		}
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
		t.wrErr(val.Error())
	case string: // Store string content directly in out using API
		t.rstOut()           // Clear buffer using API
		t.wrStringToOut(val) // Write using API
	case bool:
		var out string
		if val {
			out = trueStr
		} else {
			out = falseStr
		} // Store boolean out in out using API
		t.rstOut()           // Clear buffer using API
		t.wrStringToOut(out) // Write using API
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

// =============================================================================
// UNIVERSAL CONVERSION FUNCTION - REUSES EXISTING IMPLEMENTATIONS
// =============================================================================

// anyToBuff converts any supported type to buffer using existing conversion logic
// REUSES: fmtIntToOut, floatToOut, wrStringToOut, writeStringToErr
// Supports: string, int variants, uint variants, float variants, bool, []byte, LocStr
func anyToBuff(c *conv, dest buffDest, value any) {
	switch v := value.(type) {
	// IMMEDIATE CONVERSION - Simple Types (REUSE existing implementations)
	case string:
		writeStringToDest(c, dest, v)

	case int:
		writeIntToDest(c, dest, int64(v))
	case int8:
		writeIntToDest(c, dest, int64(v))
	case int16:
		writeIntToDest(c, dest, int64(v))
	case int32:
		writeIntToDest(c, dest, int64(v))
	case int64:
		writeIntToDest(c, dest, v)

	case uint:
		writeUintToDest(c, dest, uint64(v))
	case uint8:
		writeUintToDest(c, dest, uint64(v))
	case uint16:
		writeUintToDest(c, dest, uint64(v))
	case uint32:
		writeUintToDest(c, dest, uint64(v))
	case uint64:
		writeUintToDest(c, dest, v)

	case float32:
		writeFloatToDest(c, dest, float64(v))
	case float64:
		writeFloatToDest(c, dest, v)

	case bool:
		if v {
			writeStringToDest(c, dest, "true")
		} else {
			writeStringToDest(c, dest, "false")
		}

	case []byte:
		writeBytesToDest(c, dest, v)
	case LocStr:
		// LocStr needs translation - for now use first language (English)
		writeStringToDest(c, dest, v[0]) // v[0] is English translation

	// LAZY CONVERSION - Complex Types (store pointer, convert on demand)
	case []string:
		c.pointerVal = v
		c.kind = KSliceStr
		// No immediate conversion - wait for operation

	case map[string]string:
		c.pointerVal = v
		c.kind = KMap
		// No immediate conversion - wait for operation

	case map[string]any:
		c.pointerVal = v
		c.kind = KMap
		// No immediate conversion - wait for operation

	default:
		// Unknown type - write error using DICTIONARY (REUSE existing wrErr)
		c.wrErr(D.Type, D.Not, D.Supported)
	}
}

// =============================================================================
// BUFFER DESTINATION HELPERS - REUSE EXISTING BUFFER OPERATIONS
// =============================================================================

// writeStringToDest writes string to specified buffer destination
// REUSES: wrStringToOut, wrStringToWork, writeStringToErr
func writeStringToDest(c *conv, dest buffDest, s string) {
	switch dest {
	case buffOut:
		c.wrStringToOut(s)
	case buffWork:
		c.wrStringToWork(s)
	case buffErr:
		c.writeStringToErr(s)
	}
}

// writeBytesToDest writes bytes to specified buffer destination
// REUSES: wrToOut, wrToWork, wrToErr
func writeBytesToDest(c *conv, dest buffDest, data []byte) {
	switch dest {
	case buffOut:
		c.wrToOut(data)
	case buffWork:
		c.wrToWork(data)
	case buffErr:
		c.wrToErr(data)
	}
}

// writeIntToDest converts int64 to string and writes to destination buffer
// REUSES: fmtIntToOut logic adapted for destination selection
func writeIntToDest(c *conv, dest buffDest, val int64) {
	// Store current out state
	var tempOut []byte
	var tempOutLen int
	if dest != buffOut {
		tempOut = make([]byte, len(c.out))
		copy(tempOut, c.out)
		tempOutLen = c.outLen
		c.rstOut()
	}

	// REUSE existing fmtIntToOut implementation
	c.fmtIntToOut(val, 10, true)

	// Move result to correct destination if not buffOut
	if dest != buffOut {
		result := string(c.out[:c.outLen])
		// Restore original out state
		c.out = tempOut
		c.outLen = tempOutLen
		// Write to correct destination
		writeStringToDest(c, dest, result)
	}
}

// writeUintToDest converts uint64 to string and writes to destination buffer
// REUSES: fmtIntToOut logic with unsigned flag
func writeUintToDest(c *conv, dest buffDest, val uint64) {
	// Store current out state
	var tempOut []byte
	var tempOutLen int
	if dest != buffOut {
		tempOut = make([]byte, len(c.out))
		copy(tempOut, c.out)
		tempOutLen = c.outLen
		c.rstOut()
	}

	// REUSE existing fmtIntToOut implementation for unsigned
	c.fmtIntToOut(int64(val), 10, false)

	// Move result to correct destination if not buffOut
	if dest != buffOut {
		result := string(c.out[:c.outLen])
		// Restore original out state
		c.out = tempOut
		c.outLen = tempOutLen
		// Write to correct destination
		writeStringToDest(c, dest, result)
	}
}

// writeFloatToDest converts float64 to string and writes to destination buffer
// REUSES: floatToOut logic adapted for destination selection
func writeFloatToDest(c *conv, dest buffDest, val float64) {
	// Store current out state
	var tempOut []byte
	var tempOutLen int
	if dest != buffOut {
		tempOut = make([]byte, len(c.out))
		copy(tempOut, c.out)
		tempOutLen = c.outLen
		c.rstOut()
	}

	// Store current floatVal and restore after
	oldFloatVal := c.floatVal
	c.floatVal = val

	// REUSE existing floatToOut implementation
	c.floatToOut()

	// Restore floatVal
	c.floatVal = oldFloatVal

	// Move result to correct destination if not buffOut
	if dest != buffOut {
		result := string(c.out[:c.outLen])
		// Restore original out state
		c.out = tempOut
		c.outLen = tempOutLen
		// Write to correct destination
		writeStringToDest(c, dest, result)
	}
}
