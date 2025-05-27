package tinystring

// Simple memory optimization without external pools
// Using static buffers and manual buffer management for better performance

// Simple string builder that avoids unnecessary allocations
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

// string returns the accumulated string
func (tsb *tinyStringBuilder) string() string {
	return string(tsb.buf)
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
