package tinystring

import "sync"

// Phase 7: Conv Pool for memory optimization
// Reuse conv objects to eliminate the 53.67% allocation hotspot from newConv()
var convPool = sync.Pool{
	New: func() any {
		return &conv{
			separator: "_", // default separator
		}
	},
}

// getConv gets a reusable conv from the pool
func getConv() *conv {
	return convPool.Get().(*conv)
}

// putConv returns a conv to the pool after resetting it
func (c *conv) putConv() {
	// Reset all fields to default state
	c.stringVal = ""
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
	c.resetBuffer()

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

// resetBuffer resets the buffer length while keeping capacity
func (c *conv) resetBuffer() {
	c.buf = c.buf[:0]
}

// getReusableBuffer returns a buffer with specified capacity, reusing existing if possible
func (c *conv) getReusableBuffer(capacity int) []byte {
	c.ensureCapacity(capacity)
	c.resetBuffer()
	return c.buf
}

// Phase 9: Optimized buffer management with direct string interning
func (c *conv) setStringFromBuffer() {
	var resultStr string
	if len(c.buf) == 0 {
		resultStr = ""
	} else {
		// Direct string interning for small strings to avoid double allocation
		if len(c.buf) <= 32 {
			resultStr = internStringFromBytes(c.buf) // Direct from bytes, no temp string
		} else {
			// For large strings, direct allocation is still better
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

// getRuneBuffer gets a reusable rune buffer from the pool
func getRuneBuffer(capacity int) []rune {
	bufInterface := runeBufferPool.Get()
	buf := bufInterface.([]rune)

	// Reset the buffer
	buf = buf[:0]
	// Grow if needed
	if capacity > cap(buf) {
		// If requested capacity is much larger, allocate new buffer
		if capacity > cap(buf)*2 {
			runeBufferPool.Put(buf[:0]) // SA6002: sync.Pool expects interface{}
			return make([]rune, 0, capacity)
		}
		// Otherwise, grow the buffer
		runeBufferPool.Put(buf[:0]) // SA6002: sync.Pool expects interface{}
		return make([]rune, 0, capacity)
	}

	return buf
}

// putRuneBuffer returns a rune buffer to the pool
func putRuneBuffer(buf *[]rune) {
	if buf == nil {
		return
	}
	// Only pool buffers that aren't too large to avoid memory leaks
	if cap(*buf) <= defaultBufCap*4 {
		resetBuf := (*buf)[:0]
		runeBufferPool.Put(resetBuf) // SA6002: sync.Pool expects interface{}
	}
}

// newBuf creates an optimally-sized buffer for common string operations
func (t *conv) newBuf(sizeMultiplier int) (string, []byte) {
	str := t.getString()
	if len(str) == 0 {
		return str, nil
	}
	bufSize := len(str) * sizeMultiplier
	if bufSize < 16 {
		bufSize = 16 // Minimum useful buffer size
	}
	return str, make([]byte, 0, bufSize)
}
