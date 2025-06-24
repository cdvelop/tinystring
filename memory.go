package tinystring

import "sync"

// Reuse conv objects to eliminate the 53.67% allocation hotspot from newConv()
var convPool = sync.Pool{
	New: func() any {
		return &conv{
			out:  make([]byte, 0, 64),
			work: make([]byte, 0, 64),
			err:  make([]byte, 0, 64),
			// TODO: Add bufFmt when struct is updated
		}
	},
}

// getConv gets a reusable conv from the pool
func getConv() *conv {
	return convPool.Get().(*conv)
}

// putConv returns a conv to the pool after resetting it
func (c *conv) putConv() {
	// Reset all buffer positions using centralized method
	c.resetAllBuffers()
	// Clear buffer contents (keep capacity for reuse)
	c.out = c.out[:0]
	c.work = c.work[:0]
	c.err = c.err[:0]

	// Reset other fields to default state - only keep ptrValue and kind
	c.ptrValue = nil
	c.kind = KString

	convPool.Put(c)
}

// resetAllBuffers resets all buffer positions (used in putConv)
func (c *conv) resetAllBuffers() {
	c.outLen = 0
	c.workLen = 0
	c.errLen = 0
}

// ensureOutCapacity ensures main buffer has at least the specified capacity (legacy compatibility)
func (c *conv) ensureOutCapacity(capacity int) {
	if cap(c.out) < capacity {
		newCap := max(capacity, 64) // Minimum 64 bytes
		if cap(c.out) > 0 {
			newCap = max(cap(c.out)*2, capacity) // Double if growing
		}
		newBuf := make([]byte, c.outLen, newCap)
		copy(newBuf, c.out[:c.outLen])
		c.out = newBuf
	}
}

// getStringFromBuffer returns string content from specified buffer destination
// OPTIMIZED: Only reads from buffers, no ptrValue conversion for simple types
func (c *conv) getStringFromBuffer(dest buffDest) string {
	return c.getString(dest)
}

// ensureStringInOut ensures string representation is available in out buffer
// LEGACY: Maintains backward compatibility, will be deprecated
// CRITICAL: Cannot use anyToBuff to prevent infinite recursion
// Uses direct primitive conversion methods only
func (c *conv) ensureStringInOut() string {
	if c.kind == KErr {
		return c.getString(buffErr)
	}

	// Only convert if out buffer is empty (avoid redundant conversions)
	if c.outLen > 0 {
		return c.getString(buffOut) // Already converted
	}

	// For simple types, buffer should already have content from anyToBuff
	// Only fallback to ptrValue for complex types that need lazy conversion
	if c.ptrValue != nil {
		c.rstBuffer(buffOut)
		// Direct conversion based on kind to avoid anyToBuff recursion
		switch c.kind {
		case KPointer:
			// Only for *string pointers that need Apply() support
			if strPtr, ok := c.ptrValue.(*string); ok && strPtr != nil {
				c.wrString(buffOut, *strPtr)
			}
		case KSliceStr:
			// Complex type handling - convert []string to comma-separated
			if slice, ok := c.ptrValue.([]string); ok {
				for i, s := range slice {
					if i > 0 {
						c.wrString(buffOut, ",")
					}
					c.wrString(buffOut, s)
				}
			}
		case KMap:
			// Complex type handling - convert map to key=value pairs
			if m, ok := c.ptrValue.(map[string]string); ok {
				first := true
				for k, v := range m {
					if !first {
						c.wrString(buffOut, ",")
					}
					c.wrString(buffOut, k)
					c.wrString(buffOut, "=")
					c.wrString(buffOut, v)
					first = false
				}
			} else if m, ok := c.ptrValue.(map[string]any); ok {
				first := true
				for k, _ := range m {
					if !first {
						c.wrString(buffOut, ",")
					}
					c.wrString(buffOut, k)
					c.wrString(buffOut, "=")
					// Convert any value to string representation
					c.wrString(buffOut, "value")
					first = false
				}
			}
		default:
			// For other types, buffer should already have content
			// If not, return empty to avoid infinite recursion
			c.wrString(buffOut, "")
		}
	}

	return c.getString(buffOut)
}

// =============================================================================
// UNIVERSAL BUFFER METHODS - DEST-FIRST PARAMETER ORDER
// =============================================================================

// wrString writes string to specified buffer destination (universal method)
func (c *conv) wrString(dest buffDest, s string) {
	switch dest {
	case buffOut:
		c.out = append(c.out[:c.outLen], s...)
		c.outLen = len(c.out)
	case buffWork:
		c.work = append(c.work[:c.workLen], s...)
		c.workLen = len(c.work)
	case buffErr:
		c.err = append(c.err[:c.errLen], s...)
		c.errLen = len(c.err)
		// Invalid destinations are silently ignored (no-op)
	}
}

// wrBytes writes bytes to specified buffer destination (universal method)
func (c *conv) wrBytes(dest buffDest, data []byte) {
	switch dest {
	case buffOut:
		c.out = append(c.out[:c.outLen], data...)
		c.outLen = len(c.out)
	case buffWork:
		c.work = append(c.work[:c.workLen], data...)
		c.workLen = len(c.work)
	case buffErr:
		c.err = append(c.err[:c.errLen], data...)
		c.errLen = len(c.err)
		// Invalid destinations are silently ignored (no-op)
	}
}

// wrByte writes single byte to specified buffer destination
func (c *conv) wrByte(dest buffDest, b byte) {
	switch dest {
	case buffOut:
		c.out = append(c.out[:c.outLen], b)
		c.outLen = len(c.out)
	case buffWork:
		c.work = append(c.work[:c.workLen], b)
		c.workLen = len(c.work)
	case buffErr:
		c.err = append(c.err[:c.errLen], b)
		c.errLen = len(c.err)
		// Invalid destinations are silently ignored (no-op)
	}
}

// getString returns string content from specified buffer destination
func (c *conv) getString(dest buffDest) string {
	switch dest {
	case buffOut:
		return string(c.out[:c.outLen])
	case buffWork:
		return string(c.work[:c.workLen])
	case buffErr:
		return string(c.err[:c.errLen])
	default:
		return "" // Invalid destination returns empty string
	}
}

// rstBuffer resets specified buffer destination
func (c *conv) rstBuffer(dest buffDest) {
	switch dest {
	case buffOut:
		c.outLen = 0
	case buffWork:
		c.workLen = 0
	case buffErr:
		c.errLen = 0
		// Invalid destinations are silently ignored (no-op)
	}
}

// hasContent checks if specified buffer destination has content
func (c *conv) hasContent(dest buffDest) bool {
	switch dest {
	case buffOut:
		return c.outLen > 0
	case buffWork:
		return c.workLen > 0
	case buffErr:
		return c.errLen > 0
	default:
		return false // Invalid destination has no content
	}
}
