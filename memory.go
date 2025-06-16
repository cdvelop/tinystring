package tinystring

import "sync"

// Phase 7: Conv Pool for memory optimization
// Reuse conv objects to eliminate the 53.67% allocation hotspot from newConv()
var convPool = sync.Pool{
	New: func() interface{} {
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
	if len(c.buf) == 0 {
		c.stringVal = ""
	} else {
		// Direct string interning for small strings to avoid double allocation
		if len(c.buf) <= 32 {
			c.stringVal = internStringFromBytes(c.buf) // Direct from bytes, no temp string
		} else {
			// For large strings, direct allocation is still better
			c.stringVal = string(c.buf)
		}
	}
	c.vTpe = typeStr
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
	stringCache   = make([]cachedString, 0, 64) // Cache for interned strings - slice based
	stringCacheMu sync.RWMutex                  // Protect the cache
	maxCacheSize  = 64                          // Maximum cache entries
)

// internStringFromBytes attempts to return a cached string from bytes if available,
// otherwise creates, caches and returns the string. This avoids temporary string allocation.
func internStringFromBytes(b []byte) string {
	// Only intern small strings (most common case)
	if len(b) > 32 || len(b) == 0 {
		return string(b)
	}

	// Fast read-only check first - compare bytes directly to avoid string allocation
	stringCacheMu.RLock()
	for i := range stringCache {
		if len(stringCache[i].str) == len(b) && bytesEqual([]byte(stringCache[i].str), b) {
			stringCacheMu.RUnlock()
			return stringCache[i].ref
		}
	}
	stringCacheMu.RUnlock()

	// Not found, create string and try to add to cache (with write lock)
	s := string(b) // Single string allocation here

	stringCacheMu.Lock()
	defer stringCacheMu.Unlock()

	// Double-check in case another goroutine added it
	for i := range stringCache {
		if stringCache[i].str == s {
			return stringCache[i].ref
		}
	}

	// Add to cache if not full
	if len(stringCache) < maxCacheSize {
		stringCache = append(stringCache, cachedString{
			str: s,
			ref: s,
		})
	}

	return s
}

// bytesEqual compares two byte slices for equality without allocating
func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// internString attempts to return a cached string if available, otherwise caches and returns the string
func internString(s string) string {
	// Only intern small strings (most common case)
	if len(s) > 32 || len(s) == 0 {
		return s
	}

	// Fast read-only check first
	stringCacheMu.RLock()
	for i := range stringCache {
		if stringCache[i].str == s {
			stringCacheMu.RUnlock()
			return stringCache[i].ref
		}
	}
	stringCacheMu.RUnlock()

	// Not found, try to add to cache (with write lock)
	stringCacheMu.Lock()
	defer stringCacheMu.Unlock()

	// Double-check in case another goroutine added it
	for i := range stringCache {
		if stringCache[i].str == s {
			return stringCache[i].ref
		}
	}

	// Add to cache if not full
	if len(stringCache) < maxCacheSize {
		stringCache = append(stringCache, cachedString{
			str: s,
			ref: s,
		})
	}

	return s
}

// newBuf creates an optimally-sized buffer for common string operations
func (t *conv) newBuf(sizeMultiplier int) (string, []byte) {
	str := t.getString()
	if isEmpty(str) {
		return str, nil
	}
	bufSize := len(str) * sizeMultiplier
	if bufSize < 16 {
		bufSize = 16 // Minimum useful buffer size
	}
	return str, make([]byte, 0, bufSize)
}
