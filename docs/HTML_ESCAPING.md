# HTML Escaping

TinyString provides two convenience helpers to escape text for HTML:

- `Convert(...).EscapeAttr()` — escape a value for safe inclusion inside an HTML attribute value.
- `Convert(...).EscapeHTML()` — escape a value for safe inclusion inside HTML content.

Both functions perform simple string replacements and will escape the characters: `&`, `<`, `>`, `"`, and `'`.
Note that existing HTML entities will be escaped again (for example `&amp;` -> `&amp;amp;`). This library follows a simple replace-based escaping strategy — if you need entity-aware unescaping/escaping, consider using a full HTML parser.

Examples:

```go
Convert(`Tom & Jerry's "House" <tag>`).EscapeAttr()
// -> `Tom &amp; Jerry&#39;s &quot;House&quot; &lt;tag&gt;`

Convert(`<div>1 & 2</div>`).EscapeHTML()
// -> `&lt;div&gt;1 &amp; 2&lt;/div&gt;`
```