package fmt

// formatCompactFloat mimics Go's %g/%G: uses %f for normal range, %e/%E for very small/large, trims trailing zeros.
func formatCompactFloat(f float64, precision int, upper bool) string {
	if precision < 0 {
		precision = 6
	}
	absf := f
	if absf < 0 {
		absf = -absf
	}
	// Use scientific for very small or large numbers
	if absf != 0 && (absf < 1e-4 || absf >= 1e6) {
		return formatScientific(f, precision, upper)
	}
	// Use %f, then trim trailing zeros and dot
	mult := 1.0
	for i := 0; i < precision; i++ {
		mult *= 10
	}
	val := int64(f*mult + 0.5)
	intPart := val / int64(mult)
	fracPart := val % int64(mult)
	res := itoa(int(intPart))
	if precision > 0 {
		frac := itoaPad(int(fracPart), precision)
		// TrimSpace trailing zeros
		end := len(frac)
		for end > 0 && frac[end-1] == '0' {
			end--
		}
		if end > 0 {
			res += "." + frac[:end]
		}
	}
	return res
}
