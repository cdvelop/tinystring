package tinystring

// Convert mode constants for type-safe mode selection
type cm uint8

const (
	mi cm = iota // inicial - initial mode (like current withValue)
	mb           // buffer - buffer mode (for Write operations)
	ma           // any - any mode (for general conversions)
)

// Write appends any value to the buffer using unified type handling
// This is the core builder method that enables fluent chaining
//
// Usage:
//
//	c.Write("hello").Write(" ").Write("world")  // Strings
//	c.Write(42).Write(" items")                 // Numbers
//	c.Write('A').Write(" grade")                // Runes
func (c *conv) Write(v any) *conv {
	if len(c.err) > 0 {
		return c // Error chain interruption
	}
	// BUILDER INTEGRATION: If buffer is empty but we have initial value, transfer it first
	if len(c.out) == 0 && ((c.kind == KString && c.outLen > 0) ||
		(c.kind == KPointer && c.stringPtrVal != nil && *c.stringPtrVal != "") ||
		(c.kind == KSliceStr && len(c.stringSliceVal) > 0) ||
		(c.kind == KInt || c.kind == KUint || c.kind == KFloat64 || c.kind == KBool)) {
		c.val2Buf() // Transfer current value to buffer
	}

	// Use unified type handler with buffer mode
	c.setVal(v, mb)
	return c
}

// Reset clears all conv fields and resets the buffer
// Useful for reusing the same conv object for multiple operations
func (c *conv) Reset() *conv {
	// Reset all conv fields to default state
	c.out = c.out[:0] // Clear main buffer
	c.outLen = 0
	c.work = c.work[:0] // Clear temp buffer
	c.workLen = 0
	c.err = c.err[:0] // Clear error buffer
	c.intVal = 0
	c.uintVal = 0
	c.floatVal = 0
	c.boolVal = false
	c.stringSliceVal = nil
	c.stringPtrVal = nil
	c.kind = KString
	return c
}

// setVal is the unified type handling method that consolidates ALL type switches
func (c *conv) setVal(v any, mode cm) {
	if len(c.err) > 0 && mode != mi {
		return // Skip operations if error exists (except initial mode)
	}

	if v == nil {
		if mode == mi {
			c.setErr(D.String, D.Empty)
		}
		return
	}
	switch val := v.(type) {
	case string:
		switch mode {
		case mi: // Initial mode
			c.out = append(c.out[:0], val...)
			c.outLen = len(val)
			c.kind = KString
		case mb: // Buffer mode
			c.out = append(c.out, val...)
		case ma: // Any mode
			c.work = append(c.work[:0], val...)
			c.workLen = len(val)
		}
	case []string:
		switch mode {
		case mi: // Initial mode
			c.stringSliceVal = val
			c.kind = KSliceStr
		case mb: // Buffer mode
			// Join slice with space and append to buffer
			for i, s := range val {
				if i > 0 {
					c.out = append(c.out, ' ')
				}
				c.out = append(c.out, s...)
			}
		case ma: // Any mode
			if len(val) == 0 {
				c.work = c.work[:0]
				c.workLen = 0
			} else if len(val) == 1 {
				c.work = append(c.work[:0], val[0]...)
				c.workLen = len(val[0])
			} else {
				tmp := Convert()
				for i, s := range val {
					if i > 0 {
						tmp.Write(" ")
					}
					tmp.Write(s)
				}
				out := tmp.String()
				c.work = append(c.work[:0], out...)
				c.workLen = len(out)
			}
		}
	case *string:
		if val == nil {
			if mode == mi {
				c.setErr(D.String, D.Empty)
			}
			return
		}

		switch mode {
		case mi: // Initial mode
			c.out = append(c.out[:0], *val...)
			c.outLen = len(*val)
			c.stringPtrVal = val
			c.kind = KPointer
		case mb: // Buffer mode
			c.out = append(c.out, *val...)
		case ma: // Any mode
			c.work = append(c.work[:0], *val...)
			c.workLen = len(*val)
		}
	case bool:
		switch mode {
		case mi: // Initial mode
			c.boolVal = val
			c.kind = KBool
		case mb: // Buffer mode
			if val {
				c.out = append(c.out, trueStr...)
			} else {
				c.out = append(c.out, falseStr...)
			}
		case ma: // Any mode
			if val {
				c.work = append(c.work[:0], trueStr...)
				c.workLen = len(trueStr)
			} else {
				c.work = append(c.work[:0], falseStr...)
				c.workLen = len(falseStr)
			}
		}
	case error:
		errStr := val.Error()
		switch mode {
		case mi: // Initial mode
			c.setErr(errStr)
		case mb: // Buffer mode
			c.out = append(c.out, errStr...)
		case ma: // Any mode
			c.work = append(c.work[:0], errStr...)
			c.workLen = len(errStr)
		}
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		// Phase 3G.6: Consolidated numeric handling
		var vInt int64
		var vUint uint64
		var vFloat float64
		var isInt, isUint, isFloat bool

		switch v := val.(type) {
		case int:
			vInt, isInt = int64(v), true
		case int8:
			vInt, isInt = int64(v), true
		case int16:
			vInt, isInt = int64(v), true
		case int32:
			vInt, isInt = int64(v), true
		case int64:
			vInt, isInt = v, true
		case uint:
			vUint, isUint = uint64(v), true
		case uint8:
			vUint, isUint = uint64(v), true
		case uint16:
			vUint, isUint = uint64(v), true
		case uint32:
			vUint, isUint = uint64(v), true
		case uint64:
			vUint, isUint = v, true
		case float32:
			vFloat, isFloat = float64(v), true
		case float64:
			vFloat, isFloat = v, true
		}

		switch mode {
		case mi: // Initial mode
			if isInt {
				c.intVal = vInt
				c.kind = KInt
			} else if isUint {
				c.uintVal = vUint
				c.kind = KUint
			} else if isFloat {
				c.floatVal = vFloat
				c.kind = KFloat32
			}
		case mb: // Buffer mode
			if isInt {
				c.fmtIntGeneric(vInt, 10, true)
			} else if isUint {
				c.fmtIntGeneric(int64(vUint), 10, false)
			} else if isFloat {
				c.floatVal = vFloat
				c.floatToBufTmp()
			}
			c.out = append(c.out, c.work[:c.workLen]...)
			// tmpStr restore not needed with new buffer system
		case ma: // Any mode
			if isInt {
				c.intVal = vInt
				c.fmtIntGeneric(vInt, 10, true)
			} else if isUint {
				c.uintVal = vUint
				c.fmtIntGeneric(int64(vUint), 10, false)
			} else if isFloat {
				c.floatVal = vFloat
				c.floatToBufTmp()
			}
		}
	default:
		// Unsupported type
		if mode == mi || mode == mb {
			c.setErr(D.Type, D.Not, D.Supported)
		}
	}
}

// val2Buf converts current value directly to buffer for maximum efficiency
func (c *conv) val2Buf() {
	c.out = c.out[:0] // Reset buffer for all cases
	switch c.kind {
	case KString:
		c.out = append(c.out, c.out[:c.outLen]...)
	case KPointer:
		if c.stringPtrVal != nil {
			c.out = append(c.out, *c.stringPtrVal...)
		}
	case KSliceStr:
		for i, s := range c.stringSliceVal {
			if i > 0 {
				c.out = append(c.out, ' ')
			}
			c.out = append(c.out, s...)
		}
	case KInt, KUint, KFloat64:
		switch c.kind {
		case KInt:
			c.fmtIntGeneric(c.intVal, 10, true)
		case KUint:
			c.fmtIntGeneric(int64(c.uintVal), 10, false)
		case KFloat64:
			c.floatToBufTmp()
		}
		c.out = append(c.out, c.work[:c.workLen]...)
		// tmpStr restore not needed with new buffer system
	case KBool:
		if c.boolVal {
			c.out = append(c.out, trueStr...)
		} else {
			c.out = append(c.out, falseStr...)
		}
	}
}
