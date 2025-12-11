package fmt

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
			got := PathJoin(tc.elem...).String()
			if got != tc.want {
				t.Errorf("PathJoin(%v) = %q; want %q", tc.elem, got, tc.want)
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
			got := PathJoin(tc.elem...).String()
			if got != tc.want {
				t.Errorf("PathJoin(%v) = %q; want %q", tc.elem, got, tc.want)
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
			got := Convert(tc.path).PathBase().String()
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
			got := Convert(tc.path).PathBase().String()
			if got != tc.want {
				t.Fatalf("PathBase(%q) = %q; want %q", tc.path, got, tc.want)
			}
		})
	}
}

func TestPathExt(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{"empty", "", ""},
		{"no extension", "file", ""},
		{"simple extension", "file.txt", ".txt"},
		{"double extension", "archive.tar.gz", ".gz"},
		{"with path", "/path/to/file.txt", ".txt"},
		{"with path no ext", "/path/to/file", ""},
		{"hidden file", ".bashrc", ""},
		{"hidden with ext", ".config.yaml", ".yaml"},
		{"trailing slash", "/path/to/file.txt/", ".txt"},
		{"root", "/", ""},
		{"dot only", ".", ""},
		{"double dot", "..", ""},
		{"multiple dots", "my.file.name.doc", ".doc"},
		{"path with dots", "/my.folder/file.txt", ".txt"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Convert(tc.path).PathExt().String()
			if got != tc.want {
				t.Errorf("PathExt(%q) = %q; want %q", tc.path, got, tc.want)
			}
		})
	}
}

func TestPathExtWindows(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{"windows simple", `C:\file.exe`, ".exe"},
		{"windows path", `C:\Program Files\App\app.exe`, ".exe"},
		{"windows no ext", `C:\dir\file`, ""},
		{"windows double ext", `D:\backup\data.tar.gz`, ".gz"},
		{"unc path", `\\server\share\file.txt`, ".txt"},
		{"windows root", `C:\`, ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Convert(tc.path).PathExt().String()
			if got != tc.want {
				t.Errorf("PathExt(%q) = %q; want %q", tc.path, got, tc.want)
			}
		})
	}
}

func TestPathExtNormalizeCase(t *testing.T) {
	// Typical uses: extension in uppercase normalized to lowercase in the same chain
	got := Convert("file.TXT").PathExt().ToLower().String()
	if got != ".txt" {
		t.Errorf("Normalize ext: got %q; want %q", got, ".txt")
	}

	got = Convert(`C:\DIR\APP.EXE`).PathExt().ToLower().String()
	if got != ".exe" {
		t.Errorf("Normalize ext windows: got %q; want %q", got, ".exe")
	}

	// Already lowercase stays the same
	got = Convert("archive.tar.Gz").PathExt().ToLower().String()
	if got != ".gz" {
		t.Errorf("Normalize mixed case: got %q; want %q", got, ".gz")
	}
}

func TestPathJoinNormalizeCase(t *testing.T) {
	// Typical use: path with mixed case normalized to lowercase
	got := PathJoin("A", "B", "C").ToLower().String()
	if got != "a/b/c" {
		t.Errorf("Normalize path: got %q; want %q", got, "a/b/c")
	}

	got = PathJoin(`C:\Windows`, "System32", "DRIVERS").ToLower().String()
	if got != `c:\windows\system32\drivers` {
		t.Errorf("Normalize windows path: got %q; want %q", got, `c:\windows\system32\drivers`)
	}

	// Mixed case elements
	got = PathJoin("/VAR", "Log", "APP.log").ToLower().String()
	if got != "/var/log/app.log" {
		t.Errorf("Normalize unix path: got %q; want %q", got, "/var/log/app.log")
	}
}
