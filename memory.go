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

	// Reset other fields to default state - only keep pointerVal and kind
	c.pointerVal = nil
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
	if c.kind == KPointer && c.pointerVal != nil {
		if strPtr, ok := c.pointerVal.(*string); ok {
			*strPtr = string(c.out)
		}
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

// wrStringToOut appends string to main buffer using length-controlled writing
func (c *conv) wrStringToOut(s string) {
	c.out = append(c.out[:c.outLen], s...)
	c.outLen = len(c.out)
}

// wrToOut appends byte slice to main buffer using length-controlled writing
func (c *conv) wrToOut(data []byte) {
	c.out = append(c.out[:c.outLen], data...)
	c.outLen = len(c.out)
}

// writeByte appends single byte to main buffer
func (c *conv) writeByte(b byte) {
	c.out = append(c.out[:c.outLen], b)
	c.outLen = len(c.out)
}

// rstOut resets write position to 0 (logical reset, keeps capacity)
func (c *conv) rstOut() {
	c.outLen = 0 // Previous data ignored, will be overwritten
}

// readOut returns only valid data from main buffer (up to outLen)
func (c *conv) readOut() []byte {
	return c.out[:c.outLen]
}

// getOutString returns main buffer content as string (length-controlled)
// Note: ensureStringInOut() exists in convert.go, this is the centralized version
func (c *conv) getOutString() string {
	return string(c.out[:c.outLen]) // Only valid data
}

// =============================================================================
// ERROR BUFFER OPERATIONS (using errLen for length control)
// =============================================================================

// wrToErr appends data to error buffer with length control
func (c *conv) wrToErr(data []byte) {
	c.err = append(c.err[:c.errLen], data...)
	c.errLen = len(c.err)
}

// writeStringToErr appends string to error buffer with length control
func (c *conv) writeStringToErr(s string) {
	c.wrToErr([]byte(s))
}

// getErrorString returns error buffer content as string using errLen
func (c *conv) getErrorString() string {
	return string(c.err[:c.errLen])
}

// resetErr resets error buffer write position
func (c *conv) resetErr() {
	c.errLen = 0
}

// =============================================================================
// TEMPORARY BUFFER OPERATIONS
// =============================================================================

// wrToWork appends data to temporary buffer with length control
func (c *conv) wrToWork(data []byte) {
	c.work = append(c.work[:c.workLen], data...)
	c.workLen = len(c.work)
}

// wrStringToWork appends string to temporary buffer
func (c *conv) wrStringToWork(s string) {
	c.wrToWork([]byte(s))
}

// getWorkString returns temporary buffer content as string
func (c *conv) getWorkString() string {
	return string(c.work[:c.workLen])
}

// rstWork resets temporary buffer write position
func (c *conv) rstWork() {
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
	c.errLen = 0
}

// ensureOutCapacity ensures main buffer has at least the specified capacity
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

// bufferStats returns current buffer usage statistics (for debugging/monitoring)
func (c *conv) bufferStats() (mainLen, mainCap, tmpLen, tmpCap, errLen, errCap int) {
	return c.outLen, cap(c.out),
		c.workLen, cap(c.work),
		len(c.err), cap(c.err) // Using len() until bufErrLen field added
}

// =============================================================================
// BUFFER STATE CHECKING METHODS
// Use these instead of direct len() checks for buffer length control
// =============================================================================

// hasError checks if there's an error using errLen field
func (c *conv) hasError() bool {
	return c.errLen > 0
}

// hasWorkContent checks if work buffer has content using workLen field
func (c *conv) hasWorkContent() bool {
	return c.workLen > 0
}

// hasOutContent checks if out buffer has content using outLen field
func (c *conv) hasOutContent() bool {
	return c.outLen > 0
}

// isEmpty checks if all buffers are empty
func (c *conv) isEmpty() bool {
	return c.outLen == 0 && c.workLen == 0 && c.errLen == 0
}

// clearError resets the error buffer length
func (c *conv) clearError() {
	c.errLen = 0
}

// =============================================================================
// SIMPLIFIED CONVERSION - Uses only anyToBuff universal function
// =============================================================================

// ensureStringInOut ensures string representation is available in out buffer using anyToBuff
// Simplified to use only the universal anyToBuff function instead of custom conversion logic
func (c *conv) ensureStringInOut() string {
	if c.kind == KErr {
		return c.getErrorString()
	}

	// Only convert if out buffer is empty (avoid redundant conversions)
	if c.outLen > 0 {
		return c.getOutString() // Already converted
	}

	// Use anyToBuff for all conversions - universal approach
	if c.pointerVal != nil {
		c.rstOut()
		anyToBuff(c, buffOut, c.pointerVal)
	}

	return c.getOutString()
}
