package tinystring

import (
	"sync"
	"unsafe"
)

// Reuse Conv objects to eliminate the 53.67% allocation hotspot from newConv()
var convPool = sync.Pool{
	New: func() any {
		return &Conv{
			out:  make([]byte, 0, 64),
			work: make([]byte, 0, 64),
			err:  make([]byte, 0, 64),
			// TODO: Add bufFmt when struct is updated
		}
	},
}

// =============================================================================
// ZERO-ALLOCATION CONVERSIONS (FastHTTP Optimizations)
// =============================================================================

// unsafeBytes converts string to []byte without memory allocation
// WARNING: Do not modify the returned []byte if the source string is still in use
// SAFE FOR: Immediate use in append operations where []byte is copied
func unsafeBytes(s string) []byte {
	if len(s) == 0 {
		return nil
	}
	// #nosec G103 - Intentional unsafe operation for performance
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// getConv gets a reusable Conv from the pool
// FIXED: Ensures object is completely clean to prevent race conditions
func getConv() *Conv {
	c := convPool.Get().(*Conv)
	// Defensive cleanup: ensure object is completely clean
	c.resetAllBuffers()
	c.out = c.out[:0]
	c.work = c.work[:0]
	c.err = c.err[:0]
	c.dataPtr = nil
	c.kind = K.String
	return c
}

// putConv returns a Conv to the pool after resetting it
func (c *Conv) putConv() {
	// Reset all buffer positions using centralized method
	c.resetAllBuffers()
	// Clear buffer contents (keep capacity for reuse)
	c.out = c.out[:0]
	c.work = c.work[:0]
	c.err = c.err[:0]

	// Reset other fields to default state - only keep dataPtr and Kind
	c.dataPtr = nil
	c.kind = K.String

	convPool.Put(c)
}

// resetAllBuffers resets all buffer positions (used in putConv)
func (c *Conv) resetAllBuffers() {
	c.outLen = 0
	c.workLen = 0
	c.errLen = 0
}

// =============================================================================
// UNIVERSAL BUFFER METHODS - DEST-FIRST PARAMETER ORDER
// =============================================================================

// wrString writes string to specified buffer destination (universal method)
// OPTIMIZED: Uses zero-allocation unsafe conversion for performance
func (c *Conv) wrString(dest buffDest, s string) {
	if len(s) == 0 {
		return // No-op for empty strings
	}

	// Convert string to []byte without allocation and reuse wrBytes logic
	data := unsafeBytes(s)
	c.wrBytes(dest, data)
}

// wrBytes writes bytes to specified buffer destination (universal method)
func (c *Conv) wrBytes(dest buffDest, data []byte) {
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
func (c *Conv) wrByte(dest buffDest, b byte) {
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
// SAFE: Uses standard conversion to avoid memory corruption in concurrent access
// NOTE: unsafeString() cannot be used here because returned strings outlive Conv lifecycle
func (c *Conv) getString(dest buffDest) string {
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

// getBytes returns []byte content from specified buffer destination
// OPTIMIZED: Returns slice directly without string conversion for io.Writer compatibility
func (c *Conv) getBytes(dest buffDest) []byte {
	switch dest {
	case buffOut:
		return c.out[:c.outLen]
	case buffWork:
		return c.work[:c.workLen]
	case buffErr:
		return c.err[:c.errLen]
	default:
		return nil // Invalid destination returns nil slice
	}
}

// rstBuffer resets specified buffer destination
// FIXED: Also resets slice length to prevent data contamination
func (c *Conv) rstBuffer(dest buffDest) {
	switch dest {
	case buffOut:
		c.outLen = 0
		c.out = c.out[:0]
	case buffWork:
		c.workLen = 0
		c.work = c.work[:0]
	case buffErr:
		c.errLen = 0
		c.err = c.err[:0]
		// Invalid destinations are silently ignored (no-op)
	}
}

// hasContent checks if specified buffer destination has content
func (c *Conv) hasContent(dest buffDest) bool {
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

// swapBuff safely copies content from source buffer to destination buffer
func (c *Conv) swapBuff(src, dest buffDest) {
	// Get source slice directly (no string allocation)
	var srcData []byte
	var srcLen int

	switch src {
	case buffOut:
		srcData, srcLen = c.out[:c.outLen], c.outLen
	case buffWork:
		srcData, srcLen = c.work[:c.workLen], c.workLen
	case buffErr:
		srcData, srcLen = c.err[:c.errLen], c.errLen
	}

	// Copy directly without string conversion
	c.rstBuffer(dest)
	c.wrBytes(dest, srcData[:srcLen])
	c.rstBuffer(src)
}

// addRuneToWork encodes rune to UTF-8 and appends to work buffer
func (c *Conv) addRuneToWork(r rune) {
	// Manually encode rune to UTF-8 in work buffer
	if r < 0x80 {
		// Single byte ASCII
		c.work = append(c.work, byte(r))
	} else if r < 0x800 {
		// Two bytes
		c.work = append(c.work, 0xC0|byte(r>>6), 0x80|byte(r&0x3F))
	} else if r < 0x10000 {
		// Three bytes
		c.work = append(c.work, 0xE0|byte(r>>12), 0x80|byte((r>>6)&0x3F), 0x80|byte(r&0x3F))
	} else {
		// Four bytes
		c.work = append(c.work, 0xF0|byte(r>>18), 0x80|byte((r>>12)&0x3F), 0x80|byte((r>>6)&0x3F), 0x80|byte(r&0x3F))
	}
	c.workLen = len(c.work)
}

// bytesEqual compares buffer content with given bytes slice for optimization
// This helper eliminates getString() allocations in boolean/comparison operations
func (c *Conv) bytesEqual(dest buffDest, target []byte) bool {
	var bufData []byte
	var bufLen int

	switch dest {
	case buffOut:
		bufData, bufLen = c.out, c.outLen
	case buffWork:
		bufData, bufLen = c.work, c.workLen
	case buffErr:
		bufData, bufLen = c.err, c.errLen
	default:
		return false
	}

	// Quick length check
	if bufLen != len(target) {
		return false
	}

	// Byte-by-byte comparison
	for i := 0; i < bufLen; i++ {
		if bufData[i] != target[i] {
			return false
		}
	}
	return true
}

// bufferContainsPattern checks if any pattern is present in the buffer (no allocations)
func (c *Conv) bufferContainsPattern(dest buffDest, patterns [][]byte) bool {
	bufData := c.getBytes(dest)
	for _, pattern := range patterns {
		if bytesContain(bufData, pattern) {
			return true
		}
	}
	return false
}

// bytesContain checks if needle is present in haystack (simple byte search)
func bytesContain(haystack, needle []byte) bool {
	n := len(needle)
	h := len(haystack)
	if n == 0 || h < n {
		return false
	}
	for i := 0; i <= h-n; i++ {
		match := true
		for j := 0; j < n; j++ {
			if haystack[i+j] != needle[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}
