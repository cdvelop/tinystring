package tinystring

// KV extracts the value after the first delimiter. If not found, returns an error.
// Usage: Convert("key:value").KV(":") => "value", nil
// If no delimiter is provided, uses ":" by default.
func (c *conv) KV(delimiters ...string) (string, error) {
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
