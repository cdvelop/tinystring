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

// Rune Buffer Pool for memory optimization reuse rune buffers to eliminate allocation hotspot
var runePool = sync.Pool{
	New: func() any {
		// Start with a reasonable default capacity
		return make([]rune, 0, defaultBufCap)
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

	// Reset other fields to default state
	c.intVal = 0
	c.uintVal = 0
	c.floatVal = 0
	c.boolVal = false
	c.stringSliceVal = nil
	c.stringPtrVal = nil
	c.kind = KString

	convPool.Put(c)
}

// Phase 6.2: Buffer reuse methods for memory optimization
// ensureCapacity ensures the buffer has at least the specified capacity
func (c *conv) ensureCapacity(capacity int) {
	if cap(c.out) < capacity {
		newCap := max(capacity, 32)
		// Double the capacity if we need significant growth
		if newCap > cap(c.out)*2 {
			newCap = capacity
		} else if cap(c.out) > 0 {
			newCap = max(cap(c.out)*2, capacity)
		}
		newBuf := make([]byte, len(c.out), newCap)
		copy(newBuf, c.out)
		c.out = newBuf
	}
}

// getReusableBuffer returns a buffer with specified capacity, reusing existing if possible
func (c *conv) getReusableBuffer(capacity int) []byte {
	c.ensureCapacity(capacity)
	c.out = c.out[:0] // Inline resetBuffer
	return c.out
}

// Phase 13.2: Highly optimized buffer management with minimal allocations
func (c *conv) setStringFromBuffer() {
	// In the new system, the buffer already contains the string content
	// We just need to update the length and ensure proper type
	c.outLen = len(c.out)

	// If working with string pointer, update the original string
	if c.kind == KPointer && c.stringPtrVal != nil {
		*c.stringPtrVal = string(c.out)
		// Keep the kind as stringPtr to maintain the pointer relationship
	} else {
		c.kind = KString
	}

	// Note: We don't reset the buffer here as it contains our data
	// The buffer will be managed by the calling code as needed
}

// =============================================================================
// CENTRALIZED BUFFER MANAGEMENT - Phase 1 Implementation
// All buffer operations consolidated here for memory optimization
// =============================================================================

// writeString appends string to main buffer using length-controlled writing
func (c *conv) writeString(s string) {
	c.out = append(c.out[:c.outLen], s...)
	c.outLen = len(c.out)
}

// writeToBuffer appends byte slice to main buffer using length-controlled writing
func (c *conv) writeToBuffer(data []byte) {
	c.out = append(c.out[:c.outLen], data...)
	c.outLen = len(c.out)
}

// writeByte appends single byte to main buffer
func (c *conv) writeByte(b byte) {
	c.out = append(c.out[:c.outLen], b)
	c.outLen = len(c.out)
}

// resetBuffer resets write position to 0 (logical reset, keeps capacity)
func (c *conv) resetBuffer() {
	c.outLen = 0 // Previous data ignored, will be overwritten
}

// readBuffer returns only valid data from main buffer (up to outLen)
func (c *conv) readBuffer() []byte {
	return c.out[:c.outLen]
}

// getMainString returns main buffer content as string (length-controlled)
// Note: getString() exists in convert.go, this is the centralized version
func (c *conv) getMainString() string {
	return string(c.out[:c.outLen]) // Only valid data
}

// =============================================================================
// ERROR BUFFER OPERATIONS (temporary implementation using len() directly)
// TODO: Add bufErrLen field to conv struct for optimal performance
// =============================================================================

// writeToErrBuffer appends data to error buffer
func (c *conv) writeToErrBuffer(data []byte) {
	c.err = append(c.err, data...)
}

// writeStringToErr appends string to error buffer
func (c *conv) writeStringToErr(s string) {
	c.writeToErrBuffer([]byte(s))
}

// getErrorString returns error buffer content as string
func (c *conv) getErrorString() string {
	return string(c.err)
}

// resetErrBuffer resets error buffer
func (c *conv) resetErrBuffer() {
	c.err = c.err[:0]
}

// =============================================================================
// TEMPORARY BUFFER OPERATIONS
// =============================================================================

// writeToTmpBuffer appends data to temporary buffer with length control
func (c *conv) writeToTmpBuffer(data []byte) {
	c.work = append(c.work[:c.workLen], data...)
	c.workLen = len(c.work)
}

// writeStringToTmp appends string to temporary buffer
func (c *conv) writeStringToTmp(s string) {
	c.writeToTmpBuffer([]byte(s))
}

// getTmpString returns temporary buffer content as string
func (c *conv) getTmpString() string {
	return string(c.work[:c.workLen])
}

// resetTmpBuffer resets temporary buffer write position
func (c *conv) resetTmpBuffer() {
	c.workLen = 0
}

// =============================================================================
// FORMAT BUFFER OPERATIONS (temporary implementation without bufFmt field)
// TODO: Add bufFmt and bufFmtLen fields to conv struct for optimal performance
// =============================================================================

// cacheFormat stores format string in temporary buffer for now (shares work)
func (c *conv) cacheFormat(format string) {
	// For now, we'll implement format caching later when struct is updated
	// This is a placeholder to avoid compilation errors
}

// getCachedFormat returns cached format string
func (c *conv) getCachedFormat() string {
	return "" // Placeholder implementation
}

// hasCachedFormat checks if format matches cached format
func (c *conv) hasCachedFormat(format string) bool {
	return false // Always reparse for now, optimize later
}

// =============================================================================
// UNIFIED BUFFER STATE MANAGEMENT
// =============================================================================

// resetAllBuffers resets all buffer positions (used in putConv)
func (c *conv) resetAllBuffers() {
	c.outLen = 0
	c.workLen = 0
	// Note: bufErrLen and bufFmtLen will be added when struct is updated
}

// ensureMainCapacity ensures main buffer has at least the specified capacity
func (c *conv) ensureMainCapacity(capacity int) {
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

// bufferStats returns current buffer usage statistics (for debugging/monitoring)
func (c *conv) bufferStats() (mainLen, mainCap, tmpLen, tmpCap, errLen, errCap int) {
	return c.outLen, cap(c.out),
		c.workLen, cap(c.work),
		len(c.err), cap(c.err) // Using len() until bufErrLen field added
}
