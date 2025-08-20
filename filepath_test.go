package tinystring

import "testing"

func TestFilePathBase(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"", "."},
		{"/", "/"},
		{"///", "/"},
		{"a", "a"},
		{"a/b", "b"},
		{"a/b/", "b"},
		{"/a/b/c", "c"},
		{"/a/b/c/", "c"},
		{"///a", "a"},
		{"a//b", "b"},
		{".", "."},
		{"..", ".."},
		{"no/slash_at_end", "slash_at_end"},
		// file extension cases
		{"file.txt", "file.txt"},
		{"dir/file.txt", "file.txt"},
		{"/path/to/archive.tar.gz", "archive.tar.gz"},
		{"dir/.hidden", ".hidden"},
		{"a/b.c/d.e", "d.e"},
		{"/trailing/.bashrc/", ".bashrc"},
	}

	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			got := PathBase(tc.path)
			if got != tc.want {
				t.Fatalf("FilePathBase(%q) = %q; want %q", tc.path, got, tc.want)
			}
		})
	}
}

func TestPathBaseWindows(t *testing.T) {
	tests := []struct{ path, want string }{
		{`C:\`, `\`},
		{`C:\file.exe`, `file.exe`},
		{`C:\Program Files\App\app.exe`, `app.exe`},
		{`C:\dir\\sub\file.txt`, `file.txt`},
		{`\\server\share\file`, `file`},
	}

	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			got := PathBase(tc.path)
			if got != tc.want {
				t.Fatalf("PathBase(%q) = %q; want %q", tc.path, got, tc.want)
			}
		})
	}
}
