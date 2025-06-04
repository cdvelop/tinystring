package tinystring

// Optimized string builder using minimal memory allocations
// Based on Go's strings.Builder but with TinyGo compatibility and zero dependencies

// getBuilder returns a new string builder to minimize allocations
// Pool removed for thread safety and TinyGo compatibility
func getBuilder() *tinyStringBuilder {
	return &tinyStringBuilder{
		buf: make([]byte, 0, 64), // Start with reasonable capacity
	}
}

// putBuilder is now a no-op since we removed the pool
// Kept for API compatibility
func putBuilder(builder *tinyStringBuilder) {
	// No-op: pool removed for thread safety and minimal dependencies
}

// Optimized string builder
type tinyStringBuilder struct {
	buf []byte
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
