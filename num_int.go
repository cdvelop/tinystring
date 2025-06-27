package tinystring

func (c *conv) parseIntString(s string, base int, signed bool) int64 {
	// Handle decimal point for float-like input (e.g., "3.14")
	for i := 0; i < len(s); i++ {
		if s[i] == '.' {
			// Try to parse as float, then truncate
			f := c.parseFloat(s)
			if c.hasContent(buffErr) {
				return 0
			}
			return int64(f)
		}
	}
	if base < 2 || base > 36 {
		c.wrErr(D.Base, D.Invalid)
		return 0
	}
	var neg bool
	i := 0
	if len(s) > 0 && s[0] == '-' {
		if !signed {
			c.wrErr(D.Number, D.Negative, D.Not, D.Allowed)
			return 0
		}
		neg = true
		i = 1
		if len(s) == 1 {
			c.wrErr(D.Format, D.Invalid)
			return 0
		}
	} else if len(s) > 0 && s[0] == '+' {
		i = 1
		if len(s) == 1 {
			c.wrErr(D.Format, D.Invalid)
			return 0
		}
	}
	var n int64
	for ; i < len(s); i++ {
		ch := s[i]
		var v byte
		switch {
		case '0' <= ch && ch <= '9':
			v = ch - '0'
		case 'a' <= ch && ch <= 'z':
			v = ch - 'a' + 10
		case 'A' <= ch && ch <= 'Z':
			v = ch - 'A' + 10
		default:
			c.wrErr(D.Format, D.Invalid)
			return 0
		}
		if int(v) >= base {
			c.wrErr(D.Format, D.Invalid)
			return 0
		}
		n = n*int64(base) + int64(v)
	}
	if neg {
		n = -n
	}
	return n
}

// ToInt converts the value to an integer with optional base specification.
// If no base is provided, base 10 is used. Supports bases 2-36.
// Returns the converted integer and any error that occurred during conversion.
func (c *conv) ToInt(base ...int) (int, error) {
	val := c.parseIntBase(base...)
	if val < -2147483648 || val > 2147483647 {
		return 0, c.wrErr(D.Number, D.Overflow)
	}
	if c.hasContent(buffErr) {
		return 0, c
	}
	return int(val), nil
}

// ToUint converts the value to an unsigned integer with optional base specification.
// If no base is provided, base 10 is used. Supports bases 2-36.
// Returns the converted uint and any error that occurred during conversion.
func (c *conv) ToUint(base ...int) (uint, error) {
	val := c.parseIntBase(base...)
	if val < 0 || val > 4294967295 {
		return 0, c.wrErr(D.Number, D.Overflow)
	}
	if c.hasContent(buffErr) {
		return 0, c
	}
	return uint(val), nil
}

// getInt32 extrae el valor del buffer de salida y lo convierte a int32.
// ToInt32 extrae el valor del buffer de salida y lo convierte a int32.
func (c *conv) ToInt32(base ...int) (int32, error) {
	val := c.parseIntBase(base...)
	if val < -2147483648 || val > 2147483647 {
		return 0, c.wrErr(D.Number, D.Overflow)
	}
	if c.hasContent(buffErr) {
		return 0, c
	}
	return int32(val), nil
}

// getInt64 extrae el valor del buffer de salida y lo convierte a int64.
// ToInt64 extrae el valor del buffer de salida y lo convierte a int64.
func (c *conv) ToInt64(base ...int) (int64, error) {
	val := c.parseIntBase(base...)
	if c.hasContent(buffErr) {
		return 0, c
	}
	return val, nil
}

// ToUint32 extrae el valor del buffer de salida y lo convierte a uint32.
func (c *conv) ToUint32(base ...int) (uint32, error) {
	val := c.parseIntBase(base...)
	if val < 0 || val > 4294967295 {
		return 0, c.wrErr(D.Number, D.Overflow)
	}
	if c.hasContent(buffErr) {
		return 0, c
	}
	return uint32(val), nil
}

// ToUint64 extrae el valor del buffer de salida y lo convierte a uint64.
func (c *conv) ToUint64(base ...int) (uint64, error) {
	val := c.parseIntBase(base...)
	if c.hasContent(buffErr) {
		return 0, c
	}
	return uint64(val), nil
}

// wrIntBase writes an integer in the given base to the buffer, with optional uppercase digits
func (c *conv) wrIntBase(dest buffDest, val int64, base int, signed bool, upper ...bool) {
	if base < 2 || base > 36 {
		c.wrErr(D.Base, D.Invalid)
		return
	}
	if val == 0 {
		c.wrString(dest, "0")
		return
	}
	negative := signed && val < 0
	uval := val
	if negative {
		uval = -val
	}
	useUpper := false
	if len(upper) > 0 && upper[0] {
		useUpper = true
	}
	var digits string
	if useUpper {
		digits = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	} else {
		digits = "0123456789abcdef"
	}
	var out [64]byte
	idx := len(out)
	for uval > 0 {
		idx--
		out[idx] = digits[uval%int64(base)]
		uval /= int64(base)
	}
	if negative {
		idx--
		out[idx] = '-'
	}
	c.wrBytes(dest, out[idx:])
}

// parseIntBase reutiliza la lógica de conversión de string a int64, soportando signo y base, y reporta error usando la API interna.
// parseIntBase auto-detects signed/unsigned mode using c.kind and parses the string accordingly.
// It does not take a signed parameter; instead, it checks c.kind (KInt = signed, KUint = unsigned).
func (c *conv) parseIntBase(base ...int) int64 {

	s := c.getString(buffOut)
	baseVal := 10
	if len(base) > 0 {
		baseVal = base[0]
	}
	isSigned := c.kind == KInt
	// Solo permitir negativos en base 10
	if len(s) > 0 && s[0] == '-' {
		if baseVal == 10 {
			isSigned = true
		} else {
			isSigned = false
		}
	}
	return c.parseIntString(s, baseVal, isSigned)
}
