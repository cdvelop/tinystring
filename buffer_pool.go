package tinystring

import "unsafe"

// Optimized string builder using unsafe for minimal memory allocations
// Based on Go's strings.Builder but with TinyGo compatibility and zero dependencies

// unsafeString converts byte slice to string without allocation
// This uses the same pattern as Go's strings.Builder.String() method
func unsafeString(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	// Convert []byte to string without copying by creating a string header
	// that points to the same underlying data
	return unsafe.String(&data[0], len(data))
}

// getBuilder returns a pooled string builder to minimize allocations
func getBuilder() *tinyStringBuilder {
	// Create new builder instead of using pool to avoid concurrency issues
	// TODO: Implement proper thread-safe pooling
	return &tinyStringBuilder{
		buf: make([]byte, 0, 64), // Start with reasonable capacity
	}
}

// putBuilder returns the builder to the pool for reuse
func putBuilder(builder *tinyStringBuilder) {
	// Temporarily disabled pooling to avoid concurrency issues
	// TODO: Implement proper thread-safe pooling
}

// Optimized string builder with unsafe operations
type tinyStringBuilder struct {
	buf []byte
}

// newTinyStringBuilder creates a new string builder with initial capacity
func newTinyStringBuilder(capacity int) *tinyStringBuilder {
	return &tinyStringBuilder{
		buf: make([]byte, 0, capacity),
	}
}

// writeString appends a string to the builder
func (tsb *tinyStringBuilder) writeString(s string) {
	tsb.buf = append(tsb.buf, s...)
}

// writeByte appends a byte to the builder
func (tsb *tinyStringBuilder) writeByte(b byte) {
	tsb.buf = append(tsb.buf, b)
}

// writeRune appends a rune to the builder using manual UTF-8 encoding
func (tsb *tinyStringBuilder) writeRune(r rune) {
	// Manual UTF-8 encoding to avoid unicode/utf8 import
	if r < 0x80 {
		tsb.buf = append(tsb.buf, byte(r))
	} else if r < 0x800 {
		tsb.buf = append(tsb.buf, byte(0xC0|(r>>6)), byte(0x80|(r&0x3F)))
	} else if r < 0x10000 {
		tsb.buf = append(tsb.buf, byte(0xE0|(r>>12)), byte(0x80|((r>>6)&0x3F)), byte(0x80|(r&0x3F)))
	} else {
		tsb.buf = append(tsb.buf, byte(0xF0|(r>>18)), byte(0x80|((r>>12)&0x3F)), byte(0x80|((r>>6)&0x3F)), byte(0x80|(r&0x3F)))
	}
}

// string returns the accumulated string using unsafe conversion to avoid allocation
func (tsb *tinyStringBuilder) string() string {
	return unsafeString(tsb.buf)
}

// reset clears the builder for reuse
func (tsb *tinyStringBuilder) reset() {
	tsb.buf = tsb.buf[:0]
}

// len returns the current length of the accumulated string
func (tsb *tinyStringBuilder) len() int {
	return len(tsb.buf)
}

// grow increases the capacity of the buffer if needed
func (tsb *tinyStringBuilder) grow(n int) {
	if cap(tsb.buf)-len(tsb.buf) < n {
		newCap := len(tsb.buf) + n
		if newCap < cap(tsb.buf)*2 {
			newCap = cap(tsb.buf) * 2
		}
		newBuf := make([]byte, len(tsb.buf), newCap)
		copy(newBuf, tsb.buf)
		tsb.buf = newBuf
	}
}
