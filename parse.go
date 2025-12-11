package fmt

// ExtractValue extracts the value after the first delimiter. If not found, returns an error.
// Usage: Convert("key:value").ExtractValue(":") => "value", nil
// If no delimiter is provided, uses ":" by default.
func (c *Conv) ExtractValue(delimiters ...string) (string, error) {
	src := c.String()
	d := ":"
	if len(delimiters) > 0 && delimiters[0] != "" {
		d = delimiters[0]
	}
	if src == d {
		return "", nil
	}
	_, after, found := c.splitByDelimiterWithBuffer(src, d)
	if !found {
		return "", c.wrErr(D.Format, D.Invalid, D.Delimiter, D.Not, D.Found)
	}
	return after, nil
}

// TagValue searches for the value of a key in a Go struct tag-like string.
// Example: Convert(`json:"name" Label:"Nombre"`).TagValue("Label") => "Nombre", true
func (c *Conv) TagValue(key string) (string, bool) {
	src := c.GetString(BuffOut)

	// Reutilizar splitStr para dividir por espacios
	parts := c.splitStr(src)

	for _, part := range parts {
		// Split by ':' using existing function
		k, v, found := c.splitByDelimiterWithBuffer(part, ":")
		if !found {
			continue
		}

		if k == key {
			// Remove quotes if present
			if len(v) >= 2 && v[0] == '"' && v[len(v)-1] == '"' {
				v = v[1 : len(v)-1]
			}
			return v, true
		}
	}
	return "", false
}
