package tinystring

// Join joins path elements using the appropriate separator.
// Detects Windows paths (backslash) or Unix paths (forward slash).
// Empty elements are ignored. Returns empty string if no valid elements.
//
// Examples:
//
//	Join("a", "b", "c")           // -> "a/b/c"
//	Join("/root", "sub", "file")   // -> "/root/sub/file"
//	Join("C:\\dir", "file")        // -> "C:\dir\file"
//	Join("a", "", "b")            // -> "a/b"
func PathJoin(elem ...string) string {
	c := GetConv()
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
			c.Write(sep)
		}
		c.Write(e)
	}
	return c.String()
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
func PathBase(path string) string {
	if path == "" {
		return "."
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
			return "\\"
		}
	}

	// trim trailing separators
	for len(path) > 1 && path[len(path)-1] == sep {
		path = path[:len(path)-1]
	}

	// if path reduced to a single root separator, return it
	if len(path) == 1 && (path[0] == '/' || path[0] == '\\') {
		return path
	}

	// search from end for last separator
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == sep {
			return path[i+1:]
		}
	}
	return path
}
