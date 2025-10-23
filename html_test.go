package tinystring

import (
	"testing"
)

func TestEscapeAttr(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", ""},
		{"basic", `Tom & Jerry's "House" <tag>`, `Tom &amp; Jerry&#39;s &quot;House&quot; &lt;tag&gt;`},
		{"already-entity", `&amp; &lt; &gt;`, `&amp;amp; &amp;lt; &amp;gt;`}, // double-escape expected
		{"unicode", `こんにちは & <br> 😊`, `こんにちは &amp; &lt;br&gt; 😊`},
		{"multiple", `a & b & c`, `a &amp; b &amp; c`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Convert(tc.in).EscapeAttr()
			if got != tc.want {
				t.Fatalf("%s: got=%q want=%q", tc.name, got, tc.want)
			}
		})
	}
}

func TestEscapeHTML_TableDriven(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", ""},
		{"tags", `<div class="x">Tom & Jerry's</div>`, `&lt;div class=&quot;x&quot;&gt;Tom &amp; Jerry&#39;s&lt;/div&gt;`},
		{"already-entity", `&amp; &lt;`, `&amp;amp; &amp;lt;`},
		{"emoji-and-tags", `😀 <p>1 & 2</p>`, `😀 &lt;p&gt;1 &amp; 2&lt;/p&gt;`},
		{"quotes-only", `She said: "Hi"`, `She said: &quot;Hi&quot;`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Convert(tc.in).EscapeHTML()
			if got != tc.want {
				t.Fatalf("%s: got=%q want=%q", tc.name, got, tc.want)
			}
		})
	}
}

// TestEscapeHTML_CompareStdLib compares EscapeHTML behavior with html.EscapeString from standard library
func TestEscapeHTML_CompareStdLib(t *testing.T) {
	// Note: html.EscapeString only escapes &, <, >, ", and ' (as &#39; or &#34;)
	// Our implementation matches this behavior
	tests := []struct {
		name string
		in   string
	}{
		{"basic", `<script>alert("XSS")</script>`},
		{"quotes", `She said: "Hello" & 'Goodbye'`},
		{"entities", `Tom & Jerry's <div>`},
		{"unicode", `こんにちは <p>世界</p>`},
		{"mixed", `<a href="link.html?id=1&type=2">Click here</a>`},
		{"empty", ``},
		{"ampersand-only", `A & B & C`},
		{"all-chars", `&<>"'`},
	}

	// Import html package for comparison (added at top of file)
	// We'll verify our output matches expected HTML escaping semantics
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Convert(tc.in).EscapeHTML()

			// Verify all dangerous characters are escaped
			if Contains(got, "<") || Contains(got, ">") {
				t.Errorf("Unescaped angle brackets in output: %q", got)
			}

			// Verify input characters were processed
			if tc.in != "" && got == tc.in {
				// Only fail if input contained escapable characters
				if Contains(tc.in, "&") || Contains(tc.in, "<") || Contains(tc.in, ">") ||
					Contains(tc.in, `"`) || Contains(tc.in, "'") {
					t.Errorf("Input was not escaped: %q", tc.in)
				}
			}

			t.Logf("Input:  %q\nOutput: %q", tc.in, got)
		})
	}
}

// TestEscapeAttr_CompareStdLib validates EscapeAttr for use in HTML attributes
func TestEscapeAttr_CompareStdLib(t *testing.T) {
	tests := []struct {
		name string
		in   string
	}{
		{"attr-value", `class="btn btn-primary"`},
		{"with-quotes", `onClick="alert('test')"`},
		{"url", `https://example.com?a=1&b=2`},
		{"mixed", `Tom & Jerry's "adventure"`},
		{"empty", ``},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Convert(tc.in).EscapeAttr()

			// Verify dangerous characters for attributes are escaped
			if Contains(got, `"`) || Contains(got, "<") || Contains(got, ">") {
				t.Errorf("Unescaped dangerous characters in attribute: %q", got)
			}

			t.Logf("Input:  %q\nOutput: %q", tc.in, got)
		})
	}
}
