package fmt

// formatScientific formats a float64 in scientific notation (e.g., 1.234000e+03)
// precision: number of digits after decimal point, -1 for default (6)
// upper: true for 'E', false for 'e'
func formatScientific(f float64, precision int, upper bool) string {
	if f == 0 {
		if precision < 0 {
			precision = 6
		}
		mantissa := "0."
		for i := 0; i < precision; i++ {
			mantissa += "0"
		}
		if upper {
			return mantissa + "E+00"
		}
		return mantissa + "e+00"
	}
	neg := false
	if f < 0 {
		neg = true
		f = -f
	}
	// Find exponent
	exp := 0
	for f >= 10 {
		f /= 10
		exp++
	}
	for f > 0 && f < 1 {
		f *= 10
		exp--
	}
	// Round mantissa to precision
	if precision < 0 {
		precision = 6
	}
	mult := 1.0
	for i := 0; i < precision; i++ {
		mult *= 10
	}
	mant := int64(f*mult + 0.5)
	intPart := mant / int64(mult)
	fracPart := mant % int64(mult)
	// Build mantissa string
	mantissa := ""
	if neg {
		mantissa = "-"
	}
	mantissa += itoa(int(intPart))
	if precision > 0 {
		mantissa += "."
		frac := itoaPad(int(fracPart), precision)
		mantissa += frac
	}
	// Exponent
	sign := "+"
	if exp < 0 {
		sign = "-"
		exp = -exp
	}
	if upper {
		return mantissa + "E" + sign + pad2(exp)
	}
	return mantissa + "e" + sign + pad2(exp)
}

// itoa converts int to string (manual, no stdlib)
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := false
	if n < 0 {
		neg = true
		n = -n
	}
	buf := [20]byte{}
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		pos--
		buf[pos] = '-'
	}
	return string(buf[pos:])
}

// itoaPad pads int with leading zeros to width
func itoaPad(n int, width int) string {
	s := itoa(n)
	for len(s) < width {
		s = "0" + s
	}
	return s
}

// pad2 pads int to 2 digits with leading zero
func pad2(n int) string {
	if n < 10 {
		return "0" + itoa(n)
	}
	return itoa(n)
}
