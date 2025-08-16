package tinystring

import "io"

// =============================================================================
// FORMAT TEMPLATE SYSTEM - Printf-style formatting operations
// =============================================================================

// Fmt formats a string using a printf-style format string and arguments.
// Example: Fmt("Hello %s", "world") returns "Hello world"
func Fmt(format string, args ...any) string {
	// Inline unifiedFormat logic - eliminated wrapper function
	out := getConv() // Always obtain from pool
	out.wrFormat(buffOut, format, args...)
	return out.String()
}

// Fprintf formats according to a format specifier and writes to w.
// It returns the number of bytes written and any write error encountered.
// Example: Fprintf(os.Stdout, "Hello %s\n", "world")
func Fprintf(w io.Writer, format string, args ...any) (n int, err error) {
	// Obtain converter from pool
	c := getConv()
	defer c.putConv() // Ensure cleanup

	// Use existing wrFormat to populate buffer
	c.wrFormat(buffOut, format, args...)

	// Check for formatting errors
	if c.hasContent(buffErr) {
		return 0, c
	}

	// Write to io.Writer
	data := c.getBytes(buffOut)
	return w.Write(data)
}

// Sscanf parses formatted text from a string using printf-style format specifiers.
// It returns the number of items successfully parsed and any error encountered.
// Example: Sscanf("!3F U+003F question", "!%x U+%x %s", &pos, &enc.uv, &enc.name)
func Sscanf(src string, format string, args ...any) (n int, err error) {
	// Obtain converter from pool
	c := getConv()
	defer c.putConv() // Ensure cleanup

	// Reuse parsing logic with format pattern matching
	n = c.scanWithFormat(src, format, args...)

	// Check for parsing errors
	if c.hasContent(buffErr) {
		return n, c
	}

	return n, nil
}

// applyWidthAndAlignment applies width formatting and alignment to a string
func (c *Conv) applyWidthAndAlignment(str string, width int, leftAlign bool) string {
	if width <= 0 {
		return str
	}

	strLen := len(str)
	pad := width - strLen

	if leftAlign {
		// Para alineación a la izquierda, agregar padding solo si pad > 0
		if pad > 0 {
			return str + spaces(pad)
		}
		return str
	} else if pad > 0 {
		return spaces(pad) + str
	} else if strLen > width {
		// Truncar si el string es más largo que el ancho
		return str[:width]
	}
	return str
}

// wrFormat applies printf-style formatting to arguments and writes to specified buffer destination.
// Universal method with dest-first parameter order - follows buffer API architecture
func (c *Conv) wrFormat(dest buffDest, format string, args ...any) {
	eSz := 0
	for _, arg := range args {
		switch arg.(type) {
		case int, int8, int16, int32, int64:
			eSz += 16 // Estimate for integers
		case uint, uint8, uint16, uint32, uint64:
			eSz += 16 // Estimate for unsigned integers
		case float64, float32:
			eSz += 24 // Estimate for floats
		default:
			eSz += 16 // Default estimate
		}
	}
	// Reset buffer at start BEFORE capacity estimation to avoid contamination
	c.rstBuffer(dest)

	argIndex := 0

	for i := 0; i < len(format); i++ {
		if format[i] == '%' {
			i++

			// Parse format specifier using shared helper
			formatChar, param, formatSpec, width, leftAlign, newI := c.parseFormatSpecifier(format, i)
			i = newI

			// Handle literal %
			if formatChar == '%' {
				c.wrByte(dest, '%')
				continue
			}

			// Validate format specifier using shared validation
			if !c.isValidWriteFormatChar(formatChar) {
				c.wrErr(D.Format, D.Provided, D.Not, D.Supported, byte(formatChar))
				return
			}
			if argIndex >= len(args) {
				c.wrErr(D.Argument, D.Missing, formatSpec)
				return
			}

			// Format value using shared helper
			arg := args[argIndex]
			str := c.formatValue(arg, formatChar, param, formatSpec)
			if c.hasContent(buffErr) {
				return
			}

			// Apply width and alignment if needed
			str = c.applyWidthAndAlignment(str, width, leftAlign)
			argIndex++
			c.wrBytes(dest, []byte(str))
			continue
		} else {
			c.wrByte(dest, format[i])
		}
	}

	if !c.hasContent(buffErr) {
		// Final output is ready in dest buffer
		c.kind = K.String
	}
}

// parseFormatSpecifier extracts format specifier and parameters from format string
// Returns formatChar, param, formatSpec, width, leftAlign, and new index position
func (c *Conv) parseFormatSpecifier(format string, i int) (formatChar rune, param int, formatSpec string, width int, leftAlign bool, newI int) {
	// Parse flags and width
	for i < len(format) && (format[i] == '-' || (format[i] >= '0' && format[i] <= '9')) {
		if format[i] == '-' {
			leftAlign = true
			i++
		}
		// Parse width
		w := 0
		for i < len(format) && format[i] >= '0' && format[i] <= '9' {
			w = w*10 + int(format[i]-'0')
			i++
		}
		if w > 0 {
			width = w
		}
	}
	// Parse precision for floats
	precision := -1
	if i < len(format) && format[i] == '.' {
		i++
		p := 0
		for i < len(format) && format[i] >= '0' && format[i] <= '9' {
			p = p*10 + int(format[i]-'0')
			i++
		}
		precision = p
	}
	if i >= len(format) {
		return 0, 0, "", 0, false, i
	}

	// Parse format character and return parameters
	switch format[i] {
	case 'c':
		formatChar, param, formatSpec = 'c', 0, "%c"
	case 'U':
		formatChar, param, formatSpec = 'U', 0, "%U"
	case 'd':
		formatChar, param, formatSpec = 'd', 10, "%d"
	case 'u':
		formatChar, param, formatSpec = 'u', 10, "%u"
	case 'f':
		formatChar, param, formatSpec = 'f', precision, "%f"
	case 'e':
		formatChar, param, formatSpec = 'e', precision, "%e"
	case 'E':
		formatChar, param, formatSpec = 'E', precision, "%E"
	case 'g':
		formatChar, param, formatSpec = 'g', precision, "%g"
	case 'G':
		formatChar, param, formatSpec = 'G', precision, "%G"
	case 'o':
		formatChar, param, formatSpec = 'o', 8, "%o"
	case 'O':
		formatChar, param, formatSpec = 'O', 8, "%O"
	case 'b':
		formatChar, param, formatSpec = 'b', 2, "%b"
	case 'B':
		formatChar, param, formatSpec = 'B', 2, "%B"
	case 'x':
		formatChar, param, formatSpec = 'x', 16, "%x"
	case 'X':
		formatChar, param, formatSpec = 'X', 16, "%X"
	case 'p':
		formatChar, param, formatSpec = 'p', 0, "%p"
	case 't':
		formatChar, param, formatSpec = 't', 0, "%t"
	case 'v':
		formatChar, param, formatSpec = 'v', 0, "%v"
	case 'q':
		formatChar, param, formatSpec = 'q', 0, "%q"
	case 's':
		formatChar, param, formatSpec = 's', 0, "%s"
	case '%':
		formatChar, param, formatSpec = '%', 0, "%%"
	default:
		formatChar, param, formatSpec = rune(format[i]), 0, ""
	}

	return formatChar, param, formatSpec, width, leftAlign, i
}

// isValidFormatChar validates format characters for both read and write operations
func (c *Conv) isValidFormatChar(ch rune) bool {
	switch ch {
	case 'c', 'U', 'd', 'u', 'f', 'e', 'E', 'g', 'G', 'o', 'O', 'b', 'B', 'x', 'X', 'p', 't', 'v', 'q', 's', '%':
		return true
	default:
		return false
	}
}

// isValidWriteFormatChar validates format characters for write operations (reuses isValidFormatChar)
func (c *Conv) isValidWriteFormatChar(ch rune) bool {
	return c.isValidFormatChar(ch)
}

// spaces returns a string with n spaces
func spaces(n int) string {
	if n <= 0 {
		return ""
	}
	b := make([]byte, n)
	for i := range b {
		b[i] = ' '
	}
	return string(b)
}

// wrInvalidTypeErr writes an invalid type error for the given format spec
func (c *Conv) wrInvalidTypeErr(formatSpec string) {
	c.wrErr(D.Invalid, D.Type, D.Of, D.Argument, formatSpec)
}

// formatValue formats a single value according to format character
func (c *Conv) formatValue(arg interface{}, formatChar rune, param int, formatSpec string) string {
	switch formatChar {
	case 'c':
		// Character formatting: accept rune, byte, int
		var ch rune
		var ok bool
		switch v := arg.(type) {
		case rune:
			ch = v
			ok = true
		case byte:
			ch = rune(v)
			ok = true
		case int:
			ch = rune(v)
			ok = true
		}
		if ok {
			return string(ch)
		} else {
			c.wrInvalidTypeErr("%c")
			return ""
		}
	case 'U':
		// Unicode code point formatting: U+XXXX (always uppercase hex, at least 4 digits)
		var r rune
		var ok bool
		switch v := arg.(type) {
		case rune:
			r = v
			ok = true
		case int:
			r = rune(v)
			ok = true
		}
		if ok {
			code := int(r)
			c.rstBuffer(buffWork)
			c.wrIntBase(buffWork, int64(code), 16, false, true)
			// Pad to at least 4 digits by checking buffer length directly
			for c.workLen < 4 {
				// Prepend '0' by shifting existing content
				if c.workLen+1 > len(c.work) {
					c.work = append(c.work, 0) // Expand capacity if needed
				}
				// Shift existing content right
				copy(c.work[1:c.workLen+1], c.work[:c.workLen])
				c.work[0] = '0'
				c.workLen++
			}
			// Build "U+" prefix + hex directly in output
			return "U+" + c.getString(buffWork) // Only allocation when needed
		} else {
			c.wrInvalidTypeErr("%U")
			return ""
		}
	case 'p':
		// Pointer formatting: always print '0x' for any pointer value
		return "0x"
	case 'g', 'G':
		// Compact float formatting (manual, no stdlib)
		if floatVal, ok := c.toFloat64(arg); ok {
			c.rstBuffer(buffWork)
			compact := formatCompactFloat(floatVal, param, formatChar == 'G')
			c.wrString(buffWork, compact)
			return c.getString(buffWork) // Keep for compatibility with formatFloat usage
		} else {
			c.wrInvalidTypeErr(formatSpec)
			return ""
		}
	case 'e', 'E':
		// Scientific notation (manual, no stdlib)
		if floatVal, ok := c.toFloat64(arg); ok {
			c.rstBuffer(buffWork)
			sci := formatScientific(floatVal, param, formatChar == 'E')
			c.wrString(buffWork, sci)
			return c.getString(buffWork)
		} else {
			c.wrInvalidTypeErr(formatSpec)
			return ""
		}
	case 'q':
		// Quoted string or rune
		switch v := arg.(type) {
		case string:
			return "\"" + v + "\""
		case rune:
			return "'" + string(v) + "'"
		case byte:
			return "'" + string(rune(v)) + "'"
		}
		c.wrInvalidTypeErr(formatSpec)
		return ""
	case 't':
		// Boolean formatting
		if bval, ok := arg.(bool); ok {
			if bval {
				return "true"
			} else {
				return "false"
			}
		} else {
			c.wrInvalidTypeErr(formatSpec)
			return ""
		}
	case 'd', 'o', 'b', 'x', 'O', 'B', 'X':
		if intVal, ok := c.toInt64(arg); ok {
			c.rstBuffer(buffWork)
			// Use uppercase for 'X', 'O', 'B'
			upper := formatChar == 'X' || formatChar == 'O' || formatChar == 'B'
			if param == 10 {
				c.wrIntBase(buffWork, intVal, 10, true, upper)
			} else {
				c.wrIntBase(buffWork, intVal, param, true, upper)
			}
			return c.getString(buffWork)
		} else {
			c.wrInvalidTypeErr(formatSpec)
			return ""
		}
	case 'u':
		if uintVal, ok := c.toUint64(arg); ok {
			c.rstBuffer(buffWork)
			c.wrUintBase(buffWork, uintVal, 10)
			return c.getString(buffWork)
		} else {
			c.wrInvalidTypeErr(formatSpec)
			return ""
		}
	case 'f':
		if floatVal, ok := c.toFloat64(arg); ok {
			c.rstBuffer(buffWork)
			if param >= 0 {
				c.wrFloatWithPrecision(buffWork, floatVal, param)
			} else {
				c.wrFloat64(buffWork, floatVal)
			}
			return c.getString(buffWork)
		} else {
			c.wrInvalidTypeErr(formatSpec)
			return ""
		}
	case 's':
		if strVal, ok := arg.(string); ok {
			return strVal
		} else {
			c.wrInvalidTypeErr(formatSpec)
			return ""
		}
	case 'v':
		c.rstBuffer(buffWork)
		if errVal, ok := arg.(error); ok {
			c.wrString(buffWork, errVal.Error())
			return c.getString(buffWork)
		} else {
			c.anyToBuff(buffWork, arg)
			if c.hasContent(buffErr) {
				return ""
			}
			return c.getString(buffWork)
		}
	}
	return ""
}

// scanWithFormat parses formatted text from a string, reusing wrFormat logic
// Returns the number of items successfully parsed
func (c *Conv) scanWithFormat(src string, format string, args ...any) int {
	srcPos := 0
	fmtPos := 0
	parsed := 0

	for fmtPos < len(format) && srcPos <= len(src) {
		if format[fmtPos] == '%' {
			fmtPos++
			if fmtPos >= len(format) {
				break
			}

			// Parse format specifier using same logic as wrFormat
			formatChar := rune(format[fmtPos])

			// Handle percent literal (%%)
			if formatChar == '%' {
				// This is a literal % character - match it in source
				if srcPos >= len(src) || src[srcPos] != '%' {
					c.wrErr(D.Format, D.Invalid, "literal mismatch")
					return parsed
				}
				srcPos++
				fmtPos++
				continue
			}

			// Validate format specifier (reuse wrFormat validation)
			if !c.isValidFormatChar(formatChar) {
				c.wrErr(D.Format, D.Not, D.Supported, format[fmtPos])
				return parsed
			}

			if parsed >= len(args) {
				c.wrErr(D.Argument, D.Missing)
				return parsed
			}

			// Extract and parse value from source
			valueStr, newPos := c.extractValue(src, srcPos, formatChar)
			if valueStr == "" {
				return parsed
			}

			// Convert and assign using existing conversion logic
			if c.assignParsedValue(valueStr, formatChar, args[parsed]) {
				parsed++
			} else {
				// For type validation errors, preserve the error
				// For parsing failures (empty valueStr from non-parseable input), clear error
				if valueStr != "" {
					// Non-empty valueStr suggests a type validation error, preserve it
					return parsed
				} else {
					// Empty valueStr suggests parsing failure, clear error for partial parsing
					c.rstBuffer(buffErr)
					return parsed
				}
			}

			srcPos = newPos
			fmtPos++
		} else {
			// Literal character - must match (reuse wrFormat literal logic)
			if srcPos >= len(src) || src[srcPos] != format[fmtPos] {
				c.wrErr(D.Format, D.Invalid, "literal mismatch")
				return parsed
			}
			srcPos++
			fmtPos++
		}
	}

	return parsed
}

// parseNumber extracts a number from string starting at pos
func (c *Conv) parseNumber(src string, pos int, allowSign bool) int {
	if allowSign && pos < len(src) && (src[pos] == '-' || src[pos] == '+') {
		pos++
	}
	for pos < len(src) && src[pos] >= '0' && src[pos] <= '9' {
		pos++
	}
	return pos
}

// parseHexNumber extracts a hexadecimal number from string starting at pos
func (c *Conv) parseHexNumber(src string, pos int) int {
	for pos < len(src) && ((src[pos] >= '0' && src[pos] <= '9') ||
		(src[pos] >= 'a' && src[pos] <= 'f') ||
		(src[pos] >= 'A' && src[pos] <= 'F')) {
		pos++
	}
	return pos
}

// extractValue extracts a value from source string based on format character
func (c *Conv) extractValue(src string, pos int, formatChar rune) (string, int) {
	start := pos

	switch formatChar {
	case 'd':
		// Extract decimal number (reuse number parsing logic)
		pos = c.parseNumber(src, pos, true)

	case 'x', 'X':
		// Extract hexadecimal number
		pos = c.parseHexNumber(src, pos)

	case 'f', 'g', 'e':
		// Extract floating point number (reuse float parsing logic)
		pos = c.parseNumber(src, pos, true)
		if pos < len(src) && src[pos] == '.' {
			pos++
			pos = c.parseNumber(src, pos, false)
		}

	case 's':
		// Extract string until whitespace
		for pos < len(src) && src[pos] != ' ' && src[pos] != '\t' &&
			src[pos] != '\n' && src[pos] != '\r' {
			pos++
		}

	case 'c':
		// Extract single character
		if pos < len(src) {
			pos++
		}

	case '%':
		// Literal %
		if pos < len(src) && src[pos] == '%' {
			pos++
			return "%", pos
		}
		c.wrErr(D.Format, D.Invalid, "expected %")
		return "", pos
	}

	if start == pos {
		// No characters extracted - this is not an error for partial parsing
		return "", pos
	}

	return src[start:pos], pos
}

// assignParsedValue converts and assigns a parsed value using existing conversion logic
func (c *Conv) assignParsedValue(valueStr string, formatChar rune, arg any) bool {
	switch formatChar {
	case 'd':
		// Use buffer-based integer conversion instead of creating new Conv
		c.rstBuffer(buffWork)
		c.wrString(buffWork, valueStr)
		c.swapBuff(buffOut, buffErr)  // Save current buffOut
		c.swapBuff(buffWork, buffOut) // Move valueStr to buffOut

		switch ptr := arg.(type) {
		case *int:
			if val, err := c.Int(); err == nil {
				*ptr = val
				c.swapBuff(buffOut, buffWork) // Clear buffOut
				c.swapBuff(buffErr, buffOut)  // Restore original buffOut
				return true
			}
		case *int64:
			if val, err := c.Int64(); err == nil {
				*ptr = val
				c.swapBuff(buffOut, buffWork) // Clear buffOut
				c.swapBuff(buffErr, buffOut)  // Restore original buffOut
				return true
			}
		case *int32:
			if val, err := c.Int32(); err == nil {
				*ptr = val
				c.swapBuff(buffOut, buffWork) // Clear buffOut
				c.swapBuff(buffErr, buffOut)  // Restore original buffOut
				return true
			}
		}
		c.swapBuff(buffOut, buffWork) // Clear buffOut
		c.swapBuff(buffErr, buffOut)  // Restore original buffOut

	case 'x', 'X':
		// Reuse hexadecimal conversion logic from wrFormat
		val := c.parseHexString(valueStr)
		switch ptr := arg.(type) {
		case *int:
			*ptr = int(val)
			return true
		case *int64:
			*ptr = val
			return true
		case *int32:
			*ptr = int32(val)
			return true
		case *uint:
			*ptr = uint(val)
			return true
		case *uint32:
			*ptr = uint32(val)
			return true
		case *uint64:
			*ptr = uint64(val)
			return true
		}

	case 'f', 'g', 'e':
		// Use buffer-based float conversion instead of creating new Conv
		c.rstBuffer(buffWork)
		c.wrString(buffWork, valueStr)
		c.swapBuff(buffOut, buffErr)  // Save current buffOut
		c.swapBuff(buffWork, buffOut) // Move valueStr to buffOut

		switch ptr := arg.(type) {
		case *float64:
			if val, err := c.Float64(); err == nil {
				*ptr = val
				c.swapBuff(buffOut, buffWork) // Clear buffOut
				c.swapBuff(buffErr, buffOut)  // Restore original buffOut
				return true
			}
		case *float32:
			if val, err := c.Float32(); err == nil {
				*ptr = val
				c.swapBuff(buffOut, buffWork) // Clear buffOut
				c.swapBuff(buffErr, buffOut)  // Restore original buffOut
				return true
			}
		}
		c.swapBuff(buffOut, buffWork) // Clear buffOut
		c.swapBuff(buffErr, buffOut)  // Restore original buffOut

	case 's':
		// Direct string assignment
		if ptr, ok := arg.(*string); ok {
			*ptr = valueStr
			return true
		}

	case 'c':
		// Character assignment
		if len(valueStr) > 0 {
			switch ptr := arg.(type) {
			case *rune:
				*ptr = rune(valueStr[0])
				return true
			case *byte:
				*ptr = valueStr[0]
				return true
			}
		}
	}

	c.wrErr(D.Invalid, D.Type, D.Of, D.Argument)
	return false
}

// parseHexString converts hex string to int64 (extracted and optimized from parseScanf)
func (c *Conv) parseHexString(hexStr string) int64 {
	val := int64(0)
	for _, ch := range hexStr {
		val *= 16
		if ch >= '0' && ch <= '9' {
			val += int64(ch - '0')
		} else if ch >= 'a' && ch <= 'f' {
			val += int64(ch - 'a' + 10)
		} else if ch >= 'A' && ch <= 'F' {
			val += int64(ch - 'A' + 10)
		}
	}
	return val
}
