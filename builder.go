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
	if len(c.buf) == 0 && c.hasInitialValue() {
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
	c.roundDown = false
	c.separator = "_"
	c.tmpStr = ""
	c.lastConvType = typeStr
	c.err = ""
	c.buf = c.buf[:0] // Reset buffer length, keep capacity - SINGLE SOURCE OF TRUTH
	return c
}

// setVal is the unified type handling method that consolidates ALL type switches
// This replaces the individual withValue, Write, and any2s type switches
//
// Modes:
//   mi (inicial): Initialize conv with value (like current withValue)
//   mb (buffer):  Append value to buffer (for Write operations)
//   ma (any):     General conversion (for any2s operations)
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
		c.handleString(val, mode)
	case []string:
		c.handleStringSlice(val, mode)
	case *string:
		c.handleStringPtr(val, mode)
	case bool:
		c.handleBool(val, mode)
	case error:
		c.handleError(val, mode)
	case int:
		c.handleInt(int64(val), mode)
	case int8:
		c.handleInt(int64(val), mode)
	case int16:
		c.handleInt(int64(val), mode)
	case int32: // This handles both int32 and rune
		c.handleInt(int64(val), mode)
	case int64:
		c.handleInt(val, mode)
	case uint:
		c.handleUint(uint64(val), mode)
	case uint8: // This handles both uint8 and byte
		c.handleUint(uint64(val), mode)
	case uint16:
		c.handleUint(uint64(val), mode)
	case uint32:
		c.handleUint(uint64(val), mode)
	case uint64:
		c.handleUint(val, mode)
	case float32:
		c.handleFloat(float64(val), mode)
	case float64:
		c.handleFloat(val, mode)
	default:
		// Unsupported type
		if mode == mi {
			c.err = T(D.Unsupported, D.Type)
		} else if mode == mb {
			c.err = T(D.Unsupported, D.Type)
		}
	}
}

// Type-specific handlers for unified setVal method

func (c *conv) handleString(val string, mode cm) {
	switch mode {
	case mi: // Initial mode
		c.stringVal = val
		c.vTpe = typeStr
	case mb: // Buffer mode
		c.buf = append(c.buf, val...)
	case ma: // Any mode
		c.tmpStr = val
	}
}

func (c *conv) handleStringSlice(val []string, mode cm) {
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
		c.tmpStr = c.joinSlice(" ")
	}
}

func (c *conv) handleStringPtr(val *string, mode cm) {
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
}

func (c *conv) handleBool(val bool, mode cm) {
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
}

func (c *conv) handleError(val error, mode cm) {
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
}

func (c *conv) handleInt(val int64, mode cm) {
	switch mode {
	case mi: // Initial mode
		c.intVal = val
		c.vTpe = typeInt
	case mb: // Buffer mode
		c.intVal = val
		c.appendIntToBuf(val)
	case ma: // Any mode
		c.intVal = val
		c.fmtInt(10)
	}
}

func (c *conv) handleUint(val uint64, mode cm) {
	switch mode {
	case mi: // Initial mode
		c.uintVal = val
		c.vTpe = typeUint
	case mb: // Buffer mode
		c.uintVal = val
		c.appendUintToBuf(val)
	case ma: // Any mode
		c.uintVal = val
		c.fmtUint(10)
	}
}

func (c *conv) handleFloat(val float64, mode cm) {
	switch mode {
	case mi: // Initial mode
		c.floatVal = val
		c.vTpe = typeFloat
	case mb: // Buffer mode
		c.floatVal = val
		c.appendFloatToBuf(val)
	case ma: // Any mode
		c.floatVal = val
		c.f2s()
	}
}

// grow ensures the buffer has sufficient capacity (renamed from ensureCapacity)
func (c *conv) grow(capacity int) {
	if cap(c.buf) < capacity {
		newCap := capacity
		if newCap < 32 {
			newCap = 32
		}
		// Double the capacity if we need significant growth
		if newCap > cap(c.buf)*2 {
			newCap = capacity
		} else if cap(c.buf) > 0 {
			newCap = cap(c.buf) * 2
			if newCap < capacity {
				newCap = capacity
			}
		}
		newBuf := make([]byte, len(c.buf), newCap)
		copy(newBuf, c.buf)
		c.buf = newBuf
	}
}

// getBuf returns the current buffer content as string (replaces getString for buffer operations)
func (c *conv) getBuf() string {
	if len(c.buf) == 0 {
		return ""
	}
	return string(c.buf)
}

// val2Buf converts current value directly to buffer for maximum efficiency
func (c *conv) val2Buf() {
	switch c.vTpe {
	case typeStr:
		c.buf = append(c.buf[:0], c.stringVal...)
	case typeStrPtr:
		if c.stringPtrVal != nil {
			c.buf = append(c.buf[:0], *c.stringPtrVal...)
		}
	case typeStrSlice:
		c.buf = c.buf[:0]
		for i, s := range c.stringSliceVal {
			if i > 0 {
				c.buf = append(c.buf, ' ')
			}
			c.buf = append(c.buf, s...)
		}
	case typeInt:
		c.buf = c.buf[:0]
		c.appendIntToBuf(c.intVal)
	case typeUint:
		c.buf = c.buf[:0]
		c.appendUintToBuf(c.uintVal)
	case typeFloat:
		c.buf = c.buf[:0]
		c.appendFloatToBuf(c.floatVal)
	case typeBool:
		c.buf = c.buf[:0]
		if c.boolVal {
			c.buf = append(c.buf, trueStr...)
		} else {
			c.buf = append(c.buf, falseStr...)
		}
	default:
		c.buf = c.buf[:0]
	}
}

// appendIntToBuf appends integer directly to buffer
func (c *conv) appendIntToBuf(val int64) {
	// Use existing fmtInt logic but append to buffer instead of tmpStr
	oldTmpStr := c.tmpStr
	c.intVal = val
	c.fmtInt(10)
	c.buf = append(c.buf, c.tmpStr...)
	c.tmpStr = oldTmpStr // Restore original tmpStr
}

// appendUintToBuf appends unsigned integer directly to buffer
func (c *conv) appendUintToBuf(val uint64) {
	// Use existing fmtUint logic but append to buffer instead of tmpStr
	oldTmpStr := c.tmpStr
	c.uintVal = val
	c.fmtUint(10)
	c.buf = append(c.buf, c.tmpStr...)
	c.tmpStr = oldTmpStr // Restore original tmpStr
}

// appendFloatToBuf appends float directly to buffer
func (c *conv) appendFloatToBuf(val float64) {
	// Use existing f2s logic but append to buffer instead of tmpStr
	oldTmpStr := c.tmpStr
	c.floatVal = val
	c.f2s()
	c.buf = append(c.buf, c.tmpStr...)
	c.tmpStr = oldTmpStr // Restore original tmpStr
}

// hasInitialValue checks if conv has an initial value that should be transferred to buffer
func (c *conv) hasInitialValue() bool {
	switch c.vTpe {
	case typeStr:
		return c.stringVal != ""
	case typeStrPtr:
		return c.stringPtrVal != nil && *c.stringPtrVal != ""
	case typeStrSlice:
		return len(c.stringSliceVal) > 0
	case typeInt, typeUint, typeFloat, typeBool:
		return true // Numeric and boolean types always have a value
	default:
		return false
	}
}
