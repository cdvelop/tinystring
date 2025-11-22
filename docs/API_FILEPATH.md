# Filepath Package Equivalents

Replace common `filepath` package functions with TinyString equivalents:

| Go Standard | TinyString Equivalent |
|-------------|----------------------|
| `filepath.Base()` | `Convert(path).PathBase().String()` |
| `filepath.Join()` | `PathJoin("a", "b", "c").String()` — variadic function, zero heap allocation for ≤8 elements |

## PathBase (fluent API)

Use `Convert(path).PathBase().String()` to get the last element of a path.

Examples:

```go
Convert("/a/b/c.txt").PathBase().String() // -> "c.txt"
Convert("folder/file.txt").PathBase().String()   // -> "file.txt"
Convert("").PathBase().String()           // -> "."
Convert(`c:\\file program\\app.exe`).PathBase().String() // -> "app.exe"
```

## PathJoin (cross-platform path joining)

Standalone function with variadic string arguments.
Returns *Conv for method chaining with transformations like ToLower().
Uses fixed array for zero heap allocation (≤8 elements).
Detects separator ("/" or "\\") automatically and avoids duplicates.

Examples:

```go
PathJoin("a", "b", "c").String()            // -> "a/b/c"
PathJoin("/root", "sub", "file").String()   // -> "/root/sub/file"
PathJoin(`C:\dir`, "file").String()         // -> `C:\dir\file`
PathJoin(`\\server`, "share", "file").String() // -> `\\server\share\file`

// Typical use: normalize path case with ToLower() in the same chain
PathJoin("A", "B", "C").ToLower().String() // -> "a/b/c"
```

## Path Extension

Get the file extension (including the leading dot) from a path. Use the
fluent API form `Convert(path).PathExt().String()` which reads the path
from the Conv buffer and returns only the extension (or empty string).

Examples:

```go
Convert("file.txt").PathExt().String()          // -> ".txt"
Convert("/path/to/archive.tar.gz").PathExt().String() // -> ".gz"
Convert(".bashrc").PathExt().String()           // -> ""  (hidden file, no ext)
Convert("noext").PathExt().String()             // -> ""
Convert(`C:\\dir\\app.exe`).PathExt().String() // -> ".exe"

// Typical use: normalize extension case in the same chain. For example,
// when the extension is uppercase you can lower-case it immediately:
Convert("file.TXT").PathExt().ToLower().String() // -> ".txt"
```