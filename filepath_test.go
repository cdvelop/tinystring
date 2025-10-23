package tinystring

import "testing"

func TestJoin(t *testing.T) {
	tests := []struct {
		name string
		elem []string
		want string
	}{
		{"empty", []string{}, ""},
		{"single", []string{"a"}, "a"},
		{"two", []string{"a", "b"}, "a/b"},
		{"three", []string{"a", "b", "c"}, "a/b/c"},
		{"with root", []string{"/root", "sub", "file"}, "/root/sub/file"},
		{"empty elements", []string{"a", "", "b"}, "a/b"},
		{"all empty", []string{"", "", ""}, ""},
		{"trailing slash", []string{"a/", "b"}, "a/b"},
		{"leading slash", []string{"a", "/b"}, "a/b"},
		{"multiple slashes", []string{"a/", "/b"}, "a/b"},
		{"absolute path", []string{"/", "a", "b"}, "/a/b"},
		{"complex", []string{"/var", "log", "app.log"}, "/var/log/app.log"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := PathJoin(tc.elem...)
			if got != tc.want {
				t.Errorf("Join(%v) = %q; want %q", tc.elem, got, tc.want)
			}
		})
	}
}

func TestJoinWindows(t *testing.T) {
	tests := []struct {
		name string
		elem []string
		want string
	}{
		{"windows simple", []string{`C:\`, "dir", "file.txt"}, `C:\dir\file.txt`},
		{"windows path", []string{`C:\Program Files`, "App", "app.exe"}, `C:\Program Files\App\app.exe`},
		{"windows mixed", []string{`D:\data\`, `\logs`, "app.log"}, `D:\data\logs\app.log`},
		{"unc path", []string{`\\server`, "share", "file"}, `\\server\share\file`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := PathJoin(tc.elem...)
			if got != tc.want {
				t.Errorf("Join(%v) = %q; want %q", tc.elem, got, tc.want)
			}
		})
	}
}

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
