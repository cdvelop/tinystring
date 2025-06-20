package tinystring

import "sync"

// Phase 7: Conv Pool for memory optimization
// Reuse conv objects to eliminate the 53.67% allocation hotspot from newConv()
var convPool = sync.Pool{
	New: func() any {
		return &conv{}
	},
}

// getConv gets a reusable conv from the pool
func getConv() *conv {
	return convPool.Get().(*conv)
}

// putConv returns a conv to the pool after resetting it
func (c *conv) putConv() { // Reset all fields to default state
	c.stringVal = ""
	c.intVal = 0
	c.uintVal = 0
	c.floatVal = 0
	c.boolVal = false
	c.stringSliceVal = nil
	c.stringPtrVal = nil
	c.vTpe = typeStr
	c.tmpStr = ""
	c.err = ""
	c.buf = c.buf[:0] // Inline resetBuffer

	convPool.Put(c)
}

// Phase 6.2: Buffer reuse methods for memory optimization
// ensureCapacity ensures the buffer has at least the specified capacity
func (c *conv) ensureCapacity(capacity int) {
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

// getReusableBuffer returns a buffer with specified capacity, reusing existing if possible
func (c *conv) getReusableBuffer(capacity int) []byte {
	c.ensureCapacity(capacity)
	c.buf = c.buf[:0] // Inline resetBuffer
	return c.buf
}

// Phase 13: Optimized buffer management with selective string interning
func (c *conv) setStringFromBuffer() {
	var resultStr string
	if len(c.buf) == 0 {
		resultStr = ""
	} else {
		// Phase 13: Reduce string interning overhead - increase threshold and be more selective
		// Only intern very small strings that are likely to be repeated (error messages, common values)
		if len(c.buf) <= 16 && isLikelyReusable(c.buf) {
			resultStr = internStringFromBytes(c.buf) // Direct from bytes, no temp string
		} else {
			// For most strings, direct allocation is faster than cache overhead
			resultStr = string(c.buf)
		}
	}

	c.stringVal = resultStr

	// If working with string pointer, update the original string
	if c.vTpe == typeStrPtr && c.stringPtrVal != nil {
		*c.stringPtrVal = resultStr
		// Keep the vTpe as stringPtr to maintain the pointer relationship
	} else {
		c.vTpe = typeStr
	}

	c.buf = c.buf[:0] // Reset buffer length, keep capacity
}

// Phase 8.5: String interning cache for common small strings
// Simple string cache for frequently used small strings to reduce allocations
// Using slice-based cache for TinyGo compatibility and better performance
type cachedString struct {
	str string
	ref string // Reference to the same string to avoid duplicates
}

var (
	stringCache    [64]cachedString // Fixed-size array for interned strings - thread safe
	stringCacheLen int              // Current number of entries in cache
	stringCacheMu  sync.RWMutex     // Protect the cache
	maxCacheSize   = 64             // Maximum cache entries
)

// internStringFromBytes attempts to return a cached string from bytes if available,
// otherwise creates, caches and returns the string. This avoids temporary string allocation.
func internStringFromBytes(b []byte) string {
	// Only intern small strings (most common case)
	if len(b) > 32 || len(b) == 0 {
		return string(b)
	}

	// Create string once to avoid multiple allocations
	s := string(b)

	// Fast read-only check first - most cache hits happen here
	stringCacheMu.RLock()
	for i := 0; i < stringCacheLen; i++ {
		if stringCache[i].str == s {
			stringCacheMu.RUnlock()
			return stringCache[i].ref
		}
	}
	stringCacheMu.RUnlock()

	// Not found, try to add to cache with write lock
	stringCacheMu.Lock()
	defer stringCacheMu.Unlock()

	// Double-check pattern: another goroutine might have added it while we waited
	for i := 0; i < stringCacheLen; i++ {
		if stringCache[i].str == s {
			return stringCache[i].ref
		}
	}

	// Add to cache if not full
	if stringCacheLen < maxCacheSize {
		stringCache[stringCacheLen] = cachedString{
			str: s,
			ref: s,
		}
		stringCacheLen++
	}

	return s
}

// Phase 11: Rune Buffer Pool for memory optimization
// Reuse rune buffers to eliminate makeRuneBuf() allocation hotspot
var runeBufferPool = sync.Pool{
	New: func() any {
		// Start with a reasonable default capacity
		return make([]rune, 0, defaultBufCap)
	},
}

// Phase 13: Helper to determine if a byte slice is likely to be reused
// Only intern strings that are probably error messages or common values
func isLikelyReusable(buf []byte) bool {
	// Don't intern if too short or too long
	if len(buf) < 3 || len(buf) > 16 {
		return false
	}

	// Intern common error-like patterns
	s := string(buf)

	// Check if it looks like an error message (contains common error words)
	if contains(s, "invalid") || contains(s, "error") || contains(s, "overflow") ||
		contains(s, "base") || contains(s, "range") || contains(s, "character") {
		return true
	}

	// Check if it's a common value (small numbers, true/false, etc.)
	if len(buf) <= 5 {
		// Numbers 0-99999, true, false, etc.
		if isAllDigits(buf) || s == "true" || s == "false" {
			return true
		}
	}

	return false
}

// Helper to check if all bytes are digits
func isAllDigits(buf []byte) bool {
	for _, b := range buf {
		if b < '0' || b > '9' {
			return false
		}
	}
	return len(buf) > 0
}

// Helper to check if string contains substring (simple implementation)
func contains(s, substr string) bool {
	if len(substr) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
