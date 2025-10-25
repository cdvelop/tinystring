package tinystring

// PathJoin joins path elements using the appropriate separator.
// Accepts variadic string arguments and returns a Conv instance for method chaining.
// Detects Windows paths (backslash) or Unix paths (forward slash).
// Empty elements are ignored.
//
// Usage patterns:
//   - PathJoin("a", "b", "c").String()           // -> "a/b/c"
//   - PathJoin("a", "B", "c").ToLower().String() // -> "a/b/c"
//
// Examples:
//
//	PathJoin("a", "b", "c").String()           // -> "a/b/c"
//	PathJoin("/root", "sub", "file").String()   // -> "/root/sub/file"
//	PathJoin(`C:\dir`, "file").String()        // -> "C:\dir\file"
//	PathJoin("a", "", "b").String()            // -> "a/b"
//	PathJoin("a", "B", "c").ToLower().String() // -> "a/b/c"
func PathJoin(elem ...string) *Conv {
	c := GetConv()

	if len(elem) == 0 {
		return c
	}

	sep := "/"
	// detect separator from first element with a separator
	for _, e := range elem {
		if Index(e, "\\") != -1 {
			sep = "\\"
			break
		}
	}

	for i, e := range elem {
		if e == "" {
			continue
		}

		curr := c.GetString(BuffOut)

		// trim leading separators only if not the first element
		if i > 0 && len(curr) > 0 {
			for len(e) > 0 && (e[0] == '/' || e[0] == '\\') {
				e = e[1:]
			}
		}

		// add separator if needed
		if len(curr) > 0 && !HasSuffix(curr, sep) && e != "" {
			c.WrString(BuffOut, sep)
		}
		c.WrString(BuffOut, e)
	}

	return c
}

// pathClean normalizes a path by detecting the separator and handling special cases.
// Returns the cleaned path and the detected separator.
// This is a helper function used by PathBase and PathExt to avoid code duplication.
func pathClean(path string) (string, byte) {
	if path == "" {
		return ".", '/'
	}

	// prefer backslash if present
	sep := byte('/')
	if Index(path, "\\") != -1 {
		sep = '\\'
	}

	// windows drive root like "C:\" or with only extra separators -> return "\\"
	if sep == '\\' && len(path) >= 2 && path[1] == ':' {
		onlySep := true
		for i := 2; i < len(path); i++ {
			if path[i] != '\\' && path[i] != '/' {
				onlySep = false
				break
			}
		}
		if onlySep {
			return "\\", sep
		}
	}

	// trim trailing separators
	for len(path) > 1 && path[len(path)-1] == sep {
		path = path[:len(path)-1]
	}

	// if path reduced to a single root separator, return it
	if len(path) == 1 && (path[0] == '/' || path[0] == '\\') {
		return path, sep
	}

	return path, sep
}

// extractBase returns the base filename from a cleaned path.
// Returns empty string for special cases like ".", "/", "\\".
func extractBase(cleaned string, sep byte) string {
	// Special cases
	if cleaned == "." || cleaned == "\\" || cleaned == "/" {
		return ""
	}

	// search from end for last separator
	for i := len(cleaned) - 1; i >= 0; i-- {
		if cleaned[i] == sep {
			return cleaned[i+1:]
		}
	}
	// no separator found - whole cleaned path is the base
	return cleaned
}

// PathBase returns the last element of path, similar to
// filepath.Base from the Go standard library. It treats
// trailing slashes specially ("/a/b/" -> "b") and preserves
// a single root slash ("/" -> "/"). An empty path returns ".".
//
// The implementation uses tinystring helpers (HasSuffix and Index)
// to avoid importing the standard library and keep the function
// minimal and TinyGo-friendly.
//
// Examples:
//
//		PathBase("/a/b/c.txt") // -> "c.txt"
//		PathBase("folder/file.txt")   // -> "file.txt"
//		PathBase("")           // -> "."
//	 PathBase("c:\file program\app.exe") // -> "app.exe"
//
// PathBase writes the last element of the path into the Conv output buffer.
// Use it as: Convert(path).PathBase().String() and it behaves similarly to
// filepath.Base. Examples:
//
// Convert("/a/b/c.txt").PathBase().String() // -> "c.txt"
// Convert("folder/file.txt").PathBase().String()   // -> "file.txt"
// Convert("").PathBase().String()           // -> "."
// Convert(`c:\file program\app.exe`).PathBase().String() // -> "app.exe"
func (c *Conv) PathBase() *Conv {
	// read source path from buffer
	src := c.GetString(BuffOut)

	cleaned, sep := pathClean(src)

	// clear output buffer - PathBase will write the resulting base
	c.ResetBuffer(BuffOut)

	base := extractBase(cleaned, sep)
	if base == "" {
		// Special case: write the cleaned value (., /, or \)
		c.WrString(BuffOut, cleaned)
	} else {
		c.WrString(BuffOut, base)
	}

	return c
}

// PathExt extracts the file extension from a path and writes it to the Conv buffer.
// Returns the Conv instance for method chaining.
// An empty extension returns an empty string.
//
// Examples:
//
//	Convert("/a/b/c.txt").PathExt().String() // -> ".txt"
//	Convert("file.tar.gz").PathExt().String() // -> ".gz"
//	Convert("noext").PathExt().String()       // -> ""
//	Convert("C:\\dir\\app.EXE").PathExt().ToLower().String() // -> ".exe"
func (c *Conv) PathExt() *Conv {
	// Read current path from output buffer
	src := c.GetString(BuffOut)

	cleaned, sep := pathClean(src)

	// clear output buffer - PathExt returns only the extension
	c.ResetBuffer(BuffOut)

	// get the base filename using helper
	base := extractBase(cleaned, sep)
	if base == "" {
		// Special cases like ".", "/", "\\" have no extension
		return c
	}

	// special cases: "." and ".." have no extension
	if base == "." || base == ".." {
		return c
	}

	// search for last dot in base filename
	for i := len(base) - 1; i >= 0; i-- {
		if base[i] == '.' {
			// don't count leading dot (hidden files like .bashrc)
			if i == 0 {
				return c
			}
			c.WrString(BuffOut, base[i:])
			return c
		}
	}

	return c
}
