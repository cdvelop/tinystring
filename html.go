package fmt

// EscapeAttr returns a string safe to place inside an HTML attribute value.
//
// It escapes the following characters:
//
//	& -> &amp;
//	" -> &quot;
//	' -> &#39;
//	< -> &lt;
//	> -> &gt;
//
// Example:
//
//	s := Convert(`Tom & Jerry's "House" <tag>`).EscapeAttr()
//	// s == `Tom &amp; Jerry&#39;s &quot;House&quot; &lt;tag&gt;`
//
// Note: this method performs plain string replacements and does not detect
// existing HTML entities. Calling EscapeAttr on a string that already
// contains entities (for example `&amp;`) will produce double-escaped
// output (`&amp;amp;`). This behavior is intentional and matches a simple
// escape-for-attribute semantics.
func (c *Conv) EscapeAttr() string {
	return c.Replace("&", "&amp;").
		Replace("\"", "&quot;").
		Replace("'", "&#39;").
		Replace("<", "&lt;").
		Replace(">", "&gt;").
		String()
}

// EscapeHTML returns a string safe for inclusion into HTML content.
//
// It escapes the following characters:
//
//	& -> &amp;
//	< -> &lt;
//	> -> &gt;
//	" -> &quot;
//	' -> &#39;
//
// Example:
//
//	s := Convert(`<div class="x">Tom & Jerry's</div>`).EscapeHTML()
//	// s == `&lt;div class=&quot;x&quot;&gt;Tom &amp; Jerry&#39;s&lt;/div&gt;`
//
// Like EscapeAttr, this method uses simple replacements and will double-escape
// existing entities.
func (c *Conv) EscapeHTML() string {
	return c.Replace("&", "&amp;").
		Replace("<", "&lt;").
		Replace(">", "&gt;").
		Replace("\"", "&quot;").
		Replace("'", "&#39;").
		String()
}

// Html creates a string for HTML content, similar to Translate but without automatic spacing.
// It supports two modes:
// 1. Format mode: If the first argument is a string containing '%', it behaves like Fmt.
// 2. Concatenation mode: Otherwise, it concatenates arguments (translating LocStr) without spaces.
//
// Usage:
// Html("div", "span").String() -> "divspan"
// Html("<div class='%s'>", "foo").String() -> "<div class='foo'>"
func Html(values ...any) *Conv {
	// Use unified smart processing
	// separator="", allowStringCode=false (to support 2-letter tags), detectFormat=true
	return GetConv().SmartArgs(BuffOut, "", false, true, values...)
}
