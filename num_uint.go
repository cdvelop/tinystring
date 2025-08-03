package tinystring

// wrUintBase writes an unsigned integer in the given base to the buffer
func (c *conv) wrUintBase(dest buffDest, value uint64, base int) {
	if base < 2 || base > 36 {
		c.wrString(dest, "0")
		return
	}
	if value == 0 {
		c.wrByte(dest, '0')
		return
	}
	var buf [65]byte
	pos := len(buf)
	for value > 0 {
		pos--
		digit := value % uint64(base)
		if digit < 10 {
			buf[pos] = byte('0' + digit)
		} else {
			buf[pos] = byte('a' + digit - 10)
		}
		value /= uint64(base)
	}
	c.wrBytes(dest, buf[pos:])
}

// helpers for type conversion
func toUint64(v any) uint64 {
	switch x := v.(type) {
	case uint:
		return uint64(x)
	case uint8:
		return uint64(x)
	case uint16:
		return uint64(x)
	case uint32:
		return uint64(x)
	case uint64:
		return x
	}
	return 0
}
func toInt64(v any) int64 {
	switch x := v.(type) {
	case int:
		return int64(x)
	case int8:
		return int64(x)
	case int16:
		return int64(x)
	case int32:
		return int64(x)
	case int64:
		return x
	}
	return 0
}
