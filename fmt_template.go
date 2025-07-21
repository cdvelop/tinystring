package tinystring

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

// wrFormat applies printf-style formatting to arguments and writes to specified buffer destination.
// Universal method with dest-first parameter order - follows buffer API architecture
func (c *conv) wrFormat(dest buffDest, format string, args ...any) {
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
			// start := i // Removed unused variable
			i++
			// Parse flags and width
			leftAlign := false
			width := 0
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
				break
			}
			var formatChar rune
			var param int
			var formatSpec string
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
				c.wrByte(dest, '%')
				continue
			default:
				c.wrErr(D.Format, D.Provided, D.Not, D.Supported, format[i])
				return
			}
			if format[i] != '%' {
				if argIndex >= len(args) {
					c.wrErr(D.Argument, D.Missing, formatSpec)
					return
				}
				arg := args[argIndex]
				var str string
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
						str = string(ch)
					} else {
						c.wrErr(D.Invalid, D.Type, D.Of, D.Argument, "%c")
						return
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
						hex := ""
						c.rstBuffer(buffWork)
						c.wrIntBase(buffWork, int64(code), 16, false, true)
						hex = c.getString(buffWork)
						// Pad to at least 4 digits
						for len(hex) < 4 {
							hex = "0" + hex
						}
						str = "U+" + hex
					} else {
						c.wrErr(D.Invalid, D.Type, D.Of, D.Argument, "%U")
						return
					}
				case 'p':
					// Pointer formatting: always print '0x' for any pointer value
					str = "0x"
				case 'g', 'G':
					// Compact float formatting (manual, no stdlib)
					var floatVal float64
					var ok bool
					switch v := arg.(type) {
					case float64:
						floatVal = v
						ok = true
					case float32:
						floatVal = float64(v)
						ok = true
					}
					if ok {
						c.rstBuffer(buffWork)
						compact := formatCompactFloat(floatVal, param, formatChar == 'G')
						c.wrString(buffWork, compact)
						str = c.getString(buffWork)
					} else {
						c.wrErr(D.Invalid, D.Type, D.Of, D.Argument, formatSpec)
						return
					}
				case 'e', 'E':
					// Scientific notation (manual, no stdlib)
					var floatVal float64
					var ok bool
					switch v := arg.(type) {
					case float64:
						floatVal = v
						ok = true
					case float32:
						floatVal = float64(v)
						ok = true
					}
					if ok {
						c.rstBuffer(buffWork)
						sci := formatScientific(floatVal, param, formatChar == 'E')
						c.wrString(buffWork, sci)
						str = c.getString(buffWork)
					} else {
						c.wrErr(D.Invalid, D.Type, D.Of, D.Argument, formatSpec)
						return
					}
				case 'q':
					// Quoted string or rune
					var ok bool
					switch v := arg.(type) {
					case string:
						str = "\"" + v + "\""
						ok = true
					case rune:
						str = "'" + string(v) + "'"
						ok = true
					case byte:
						str = "'" + string(rune(v)) + "'"
						ok = true
					}
					if !ok {
						c.wrErr(D.Invalid, D.Type, D.Of, D.Argument, formatSpec)
						return
					}
				case 't':
					// Boolean formatting
					var ok bool
					var bval bool
					switch v := arg.(type) {
					case bool:
						bval = v
						ok = true
					}
					if ok {
						if bval {
							str = "true"
						} else {
							str = "false"
						}
					} else {
						c.wrErr(D.Invalid, D.Type, D.Of, D.Argument, formatSpec)
						return
					}
				case 'd', 'o', 'b', 'x', 'O', 'B', 'X':
					var intVal int64
					var ok bool
					switch v := arg.(type) {
					case int:
						intVal = int64(v)
						ok = true
					case int8:
						intVal = int64(v)
						ok = true
					case int16:
						intVal = int64(v)
						ok = true
					case int32:
						intVal = int64(v)
						ok = true
					case int64:
						intVal = v
						ok = true
					case uint, uint8, uint16, uint32, uint64:
						intVal = int64(toUint64(v))
						ok = true
					}
					if ok {
						c.rstBuffer(buffWork)
						// Use uppercase for 'X', 'O', 'B'
						upper := formatChar == 'X' || formatChar == 'O' || formatChar == 'B'
						if param == 10 {
							c.wrIntBase(buffWork, intVal, 10, true, upper)
						} else {
							c.wrIntBase(buffWork, intVal, param, true, upper)
						}
						str = c.getString(buffWork)
					} else {
						c.wrErr(D.Invalid, D.Type, D.Of, D.Argument, formatSpec)
						return
					}
				case 'u':
					var uintVal uint64
					var ok bool
					switch v := arg.(type) {
					case uint:
						uintVal = uint64(v)
						ok = true
					case uint8:
						uintVal = uint64(v)
						ok = true
					case uint16:
						uintVal = uint64(v)
						ok = true
					case uint32:
						uintVal = uint64(v)
						ok = true
					case uint64:
						uintVal = v
						ok = true
					case int, int8, int16, int32, int64:
						// Accept signed as unsigned for %u
						uintVal = uint64(toInt64(v))
						ok = true
					}
					if ok {
						c.rstBuffer(buffWork)
						c.wrUintBase(buffWork, uintVal, 10)
						str = c.getString(buffWork)
					} else {
						c.wrErr(D.Invalid, D.Type, D.Of, D.Argument, formatSpec)
						return
					}
				case 'f':
					if floatVal, ok := arg.(float64); ok {
						c.rstBuffer(buffWork)
						if param >= 0 {
							c.wrFloatWithPrecision(buffWork, floatVal, param)
						} else {
							c.wrFloat64(buffWork, floatVal)
						}
						str = c.getString(buffWork)
					} else {
						c.wrErr(D.Invalid, D.Type, D.Of, D.Argument, formatSpec)
						return
					}
				case 's':
					if strVal, ok := arg.(string); ok {
						str = strVal
					} else {
						c.wrErr(D.Invalid, D.Type, D.Of, D.Argument, formatSpec)
						return
					}
				case 'v':
					c.rstBuffer(buffWork)
					if errVal, ok := arg.(error); ok {
						c.wrString(buffWork, errVal.Error())
						str = c.getString(buffWork)
					} else {
						c.anyToBuff(buffWork, arg)
						if c.hasContent(buffErr) {
							return
						}
						str = c.getString(buffWork)
					}
				}
				// Apply width and alignment if needed
				if width > 0 {
					strLen := len(str)
					pad := width - strLen
					if leftAlign {
						// Para alineación a la izquierda, agregar padding solo si pad > 0
						if pad > 0 {
							str = str + spaces(pad)
						}
					} else if pad > 0 {
						str = spaces(pad) + str
					} else if strLen > width {
						// Truncar si el string es más largo que el ancho
						str = str[:width]
					}
				}
				argIndex++
				c.wrBytes(dest, []byte(str))
				continue
			}
		} else {
			c.wrByte(dest, format[i])
		}
	}

	if !c.hasContent(buffErr) {
		// Final output is ready in dest buffer
		c.Kind = K.String
	}
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
