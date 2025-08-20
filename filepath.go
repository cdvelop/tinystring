package tinystring

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
