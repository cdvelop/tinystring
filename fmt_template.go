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
	// Reset buffer at start to avoid concatenation issues
	c.rstBuffer(dest) // Use API instead of manual manipulation

	// Pre-calculate buffer size to reduce reallocations
	eSz := len(format)
	for _, arg := range args {
		switch arg.(type) {
		case string:
			eSz += 32 // Estimate for strings
		case int, int64, int32:
			eSz += 16 // Estimate for integers
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
			case 'd':
				formatChar, param, formatSpec = 'd', 10, "%d"
			case 'f':
				formatChar, param, formatSpec = 'f', precision, "%f"
			case 'o':
				formatChar, param, formatSpec = 'o', 8, "%o"
			case 'b':
				formatChar, param, formatSpec = 'b', 2, "%b"
			case 'x':
				formatChar, param, formatSpec = 'x', 16, "%x"
			case 'v':
				formatChar, param, formatSpec = 'v', 0, "%v"
			case 's':
				formatChar, param, formatSpec = 's', 0, "%s"
			case '%':
				c.wrByte(dest, '%')
				continue
			default:
				c.wrErr(D.Format, D.Specifier, D.Not, D.Supported, format[i])
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
				case 'd', 'o', 'b', 'x':
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
					case uint:
						intVal = int64(v)
						ok = true
					case uint8:
						intVal = int64(v)
						ok = true
					case uint16:
						intVal = int64(v)
						ok = true
					case uint32:
						intVal = int64(v)
						ok = true
					case uint64:
						if v <= 9223372036854775807 {
							intVal = int64(v)
							ok = true
						}
					}
					if ok {
						c.rstBuffer(buffWork)
						if param == 10 {
							c.wrIntBase(buffWork, intVal, 10, true)
						} else {
							c.wrIntBase(buffWork, intVal, param, true)
						}
						str = c.getString(buffWork)
					} else {
						c.wrErr(D.Invalid, D.Type, D.Of, D.Argument, formatSpec)
						c.kind = KErr
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
						c.kind = KErr
						return
					}
				case 's':
					if strVal, ok := arg.(string); ok {
						str = strVal
					} else {
						c.wrErr(D.Invalid, D.Type, D.Of, D.Argument, formatSpec)
						c.kind = KErr
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
		c.kind = KString
	} else {
		c.kind = KErr
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
