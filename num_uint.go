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

// toUint64 converts various integer types to uint64 with validation
func (c *conv) toUint64(arg any) (uint64, bool) {
	switch v := arg.(type) {
	case uint:
		return uint64(v), true
	case uint8:
		return uint64(v), true
	case uint16:
		return uint64(v), true
	case uint32:
		return uint64(v), true
	case uint64:
		return v, true
	case int:
		return uint64(v), true
	case int8:
		return uint64(v), true
	case int16:
		return uint64(v), true
	case int32:
		return uint64(v), true
	case int64:
		return uint64(v), true
	default:
		return 0, false
	}
}

// Uint converts the value to an unsigned integer with optional base specification.
// If no base is provided, base 10 is used. Supports bases 2-36.
// Returns the converted uint and any error that occurred during conversion.
func (c *conv) Uint(base ...int) (uint, error) {
	val := c.parseIntBase(base...)
	if val < 0 || val > 4294967295 {
		return 0, c.wrErr(D.Number, D.Overflow)
	}
	if c.hasContent(buffErr) {
		return 0, c
	}
	return uint(val), nil
}

// Uint32 extrae el valor del buffer de salida y lo convierte a uint32.
func (c *conv) Uint32(base ...int) (uint32, error) {
	val := c.parseIntBase(base...)
	if val < 0 || val > 4294967295 {
		return 0, c.wrErr(D.Number, D.Overflow)
	}
	if c.hasContent(buffErr) {
		return 0, c
	}
	return uint32(val), nil
}

// Uint64 extrae el valor del buffer de salida y lo convierte a uint64.
func (c *conv) Uint64(base ...int) (uint64, error) {
	val := c.parseIntBase(base...)
	if c.hasContent(buffErr) {
		return 0, c
	}
	return uint64(val), nil
}
