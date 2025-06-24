package tinystring

// =============================================================================
// FORMAT TEMPLATE SYSTEM - Printf-style formatting operations
// =============================================================================

// Fmt formats a string using a printf-style format string and arguments.
// Example: Fmt("Hello %s", "world") returns "Hello world"
func Fmt(format string, args ...any) *conv {
	// Inline unifiedFormat logic - eliminated wrapper function
	out := getConv() // Always obtain from pool
	out.wrFormat(buffOut, format, args...)
	return out
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
			if i+1 < len(format) {
				i++ // Handle precision for floats (e.g., "%.2f")
				precision := -1
				if format[i] == '.' {
					i++
					start := i
					for i < len(format) && format[i] >= '0' && format[i] <= '9' {
						i++
					}
					if start < i {
						// Parse precision directly without creating new conv
						precisionStr := format[start:i]
						precision = 0
						for _, char := range precisionStr {
							if char >= '0' && char <= '9' {
								precision = precision*10 + int(char-'0')
							}
						}
					}
				} // Handle format specifiers
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
				} // Common format handling logic for all specifiers except '%'
				if format[i] != '%' {
					// Inline handleFormat logic
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

						// Handle all integer types
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
							if v <= 9223372036854775807 { // Max int64
								intVal = int64(v)
								ok = true
							}
						}

						if ok {
							// Clear work buffer before use to prevent contamination
							c.rstBuffer(buffWork)

							if param == 10 {
								c.wrInt(buffWork, intVal)
							} else {
								c.wrInt64Base(buffWork, intVal, param)
							}
							str = c.getString(buffWork)
						} else {
							c.wrErr(D.Invalid, D.Type, D.Of, D.Argument, formatSpec)
							c.kind = KErr
							return
						}
					case 'f':
						if floatVal, ok := arg.(float64); ok {
							// Clear work buffer before use
							c.rstBuffer(buffWork)

							// Apply precision directly to float value if specified
							if param >= 0 {
								c.wrFloatWithPrecision(buffWork, floatVal, param)
							} else {
								// Convert float to string without precision limit
								c.wrFloat(buffWork, floatVal)
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
						// Clear work buffer before use
						c.rstBuffer(buffWork)

						// Special handling for error types in %v format
						if errVal, ok := arg.(error); ok {
							c.wrString(buffWork, errVal.Error())
							str = c.getString(buffWork)
						} else {
							// Use anyToBuff for proper conversion following buffer API
							c.anyToBuff(buffWork, arg)
							if c.hasContent(buffErr) {
								return
							}
							str = c.getString(buffWork)
						}

					}

					argIndex++
					c.wrBytes(dest, []byte(str))
					continue
				}
			} else {
				c.wrByte(dest, format[i])
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
