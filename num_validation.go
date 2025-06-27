package tinystring

// Phase 11: Extended fast parsing for common integers (0-99999)
// parseSmallInt optimizes parsing of small integers using direct byte access
// Returns the parsed integer and nil if successful, otherwise returns 0 and non-nil error
// Expanded from 999 to 99999 to handle more common integer patterns
func (c *conv) parseSmallInt(s string) (int, error) {
	if len(s) == 0 {
		return 0, c.wrErr(D.String, D.Empty)
	}

	var out int
	var negative bool

	// Check for negative sign
	i := 0
	if s[0] == '-' {
		negative = true
		i = 1
		if len(s) == 1 {
			return 0, c.wrErr(D.Format, D.Invalid)
		}
	} else if s[0] == '+' {
		i = 1
		if len(s) == 1 {
			return 0, c.wrErr(D.Format, D.Invalid)
		}
	}

	// Manual parsing with overflow check for up to 99999
	for ; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return 0, c.wrErr(D.Format, D.Invalid)
		}
		digit := int(s[i] - '0')
		// Check for overflow before multiplication
		if out > (99999-digit)/10 {
			return 0, c.wrErr(D.Number, D.Overflow)
		}
		out = out*10 + digit
	}

	if negative {
		out = -out
	}

	return out, nil
}
