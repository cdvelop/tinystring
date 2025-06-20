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
//   c.Write("hello").Write(" ").Write("world")  // Strings
//   c.Write(42).Write(" items")                 // Numbers
//   c.Write('A').Write(" grade")                // Runes
func (c *conv) Write(v any) *conv {
	if c.err != "" {
		return c // Error chain interruption
	}
	// BUILDER INTEGRATION: If buffer is empty but we have initial value, transfer it first
	if len(c.buf) == 0 && ((c.vTpe == typeStr && c.stringVal != "") ||
		(c.vTpe == typeStrPtr && c.stringPtrVal != nil && *c.stringPtrVal != "") ||
		(c.vTpe == typeStrSlice && len(c.stringSliceVal) > 0) ||
		(c.vTpe == typeInt || c.vTpe == typeUint || c.vTpe == typeFloat || c.vTpe == typeBool)) {
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
	c.stringVal = "" // CRITICAL: Clear string value
	c.intVal = 0
	c.uintVal = 0
	c.floatVal = 0
	c.boolVal = false
	c.stringSliceVal = nil
	c.stringPtrVal = nil
	c.vTpe = typeStr
	c.tmpStr = ""
	c.err = ""
	c.buf = c.buf[:0] // Reset buffer length, keep capacity - SINGLE SOURCE OF TRUTH
	return c
}

// setVal is the unified type handling method that consolidates ALL type switches
func (c *conv) setVal(v any, mode cm) {
	if c.err != "" && mode != mi {
		return // Skip operations if error exists (except initial mode)
	}

	if v == nil {
		if mode == mi {
			c.err = T(D.String, D.Empty)
		}
		return
	}
	switch val := v.(type) {
	case string:
		switch mode {
		case mi: // Initial mode
			c.stringVal = val
			c.vTpe = typeStr
		case mb: // Buffer mode
			c.buf = append(c.buf, val...)
		case ma: // Any mode
			c.tmpStr = val
		}
	case []string:
		switch mode {
		case mi: // Initial mode
			c.stringSliceVal = val
			c.vTpe = typeStrSlice
		case mb: // Buffer mode
			// Join slice with space and append to buffer
			for i, s := range val {
				if i > 0 {
					c.buf = append(c.buf, ' ')
				}
				c.buf = append(c.buf, s...)
			}
		case ma: // Any mode
			if len(val) == 0 {
				c.tmpStr = ""
			} else if len(val) == 1 {
				c.tmpStr = val[0]
			} else {
				tmp := Convert()
				for i, s := range val {
					if i > 0 {
						tmp.Write(" ")
					}
					tmp.Write(s)
				}
				c.tmpStr = tmp.String()
			}
		}
	case *string:
		if val == nil {
			if mode == mi {
				c.err = T(D.String, D.Empty)
			}
			return
		}

		switch mode {
		case mi: // Initial mode
			c.stringVal = *val
			c.stringPtrVal = val
			c.vTpe = typeStrPtr
		case mb: // Buffer mode
			c.buf = append(c.buf, *val...)
		case ma: // Any mode
			c.tmpStr = *val
		}
	case bool:
		switch mode {
		case mi: // Initial mode
			c.boolVal = val
			c.vTpe = typeBool
		case mb: // Buffer mode
			if val {
				c.buf = append(c.buf, trueStr...)
			} else {
				c.buf = append(c.buf, falseStr...)
			}
		case ma: // Any mode
			if val {
				c.tmpStr = trueStr
			} else {
				c.tmpStr = falseStr
			}
		}
	case error:
		errStr := val.Error()
		switch mode {
		case mi: // Initial mode
			c.err = errStr
			c.vTpe = typeErr
		case mb: // Buffer mode
			c.buf = append(c.buf, errStr...)
		case ma: // Any mode
			c.tmpStr = errStr
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
				c.vTpe = typeInt
			} else if isUint {
				c.uintVal = vUint
				c.vTpe = typeUint
			} else if isFloat {
				c.floatVal = vFloat
				c.vTpe = typeFloat
			}
		case mb: // Buffer mode
			oldTmpStr := c.tmpStr
			if isInt {
				c.fmtIntGeneric(vInt, 10, true)
			} else if isUint {
				c.fmtIntGeneric(int64(vUint), 10, false)
			} else if isFloat {
				c.floatVal = vFloat
				c.f2s()
			}
			c.buf = append(c.buf, c.tmpStr...)
			c.tmpStr = oldTmpStr // Restore original tmpStr
		case ma: // Any mode
			if isInt {
				c.intVal = vInt
				c.fmtIntGeneric(vInt, 10, true)
			} else if isUint {
				c.uintVal = vUint
				c.fmtIntGeneric(int64(vUint), 10, false)
			} else if isFloat {
				c.floatVal = vFloat
				c.f2s()
			}
		}
	default:
		// Unsupported type
		if mode == mi || mode == mb {
			c.err = T(D.Type, D.Not, D.Supported)
		}
	}
}

// val2Buf converts current value directly to buffer for maximum efficiency
func (c *conv) val2Buf() {
	c.buf = c.buf[:0] // Reset buffer for all cases
	switch c.vTpe {
	case typeStr:
		c.buf = append(c.buf, c.stringVal...)
	case typeStrPtr:
		if c.stringPtrVal != nil {
			c.buf = append(c.buf, *c.stringPtrVal...)
		}
	case typeStrSlice:
		for i, s := range c.stringSliceVal {
			if i > 0 {
				c.buf = append(c.buf, ' ')
			}
			c.buf = append(c.buf, s...)
		}
	case typeInt, typeUint, typeFloat:
		oldTmpStr := c.tmpStr
		switch c.vTpe {
		case typeInt:
			c.fmtIntGeneric(c.intVal, 10, true)
		case typeUint:
			c.fmtIntGeneric(int64(c.uintVal), 10, false)
		case typeFloat:
			c.f2s()
		}
		c.buf = append(c.buf, c.tmpStr...)
		c.tmpStr = oldTmpStr // Restore original tmpStr
	case typeBool:
		if c.boolVal {
			c.buf = append(c.buf, trueStr...)
		} else {
			c.buf = append(c.buf, falseStr...)
		}
	}
}
