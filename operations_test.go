package fmt

import "testing"

func TestLastIndex(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		substr string
		want   int
	}{
		// Casos básicos
		{
			name:   "substring found at end",
			s:      "hello world",
			substr: "world",
			want:   6,
		},
		{
			name:   "substring found at beginning",
			s:      "hello world",
			substr: "hello",
			want:   0,
		},
		{
			name:   "substring found in middle",
			s:      "hello world hello",
			substr: "hello",
			want:   12, // última ocurrencia
		},
		{
			name:   "substring not found",
			s:      "hello world",
			substr: "xyz",
			want:   -1,
		},
		{
			name:   "single character found",
			s:      "hello",
			substr: "l",
			want:   3, // última 'l'
		},
		{
			name:   "single character at end",
			s:      "hello",
			substr: "o",
			want:   4,
		},

		// Casos edge
		{
			name:   "empty substring",
			s:      "hello",
			substr: "",
			want:   5, // len(s)
		},
		{
			name:   "empty string with empty substring",
			s:      "",
			substr: "",
			want:   0,
		},
		{
			name:   "substring longer than string",
			s:      "hi",
			substr: "hello",
			want:   -1,
		},
		{
			name:   "identical strings",
			s:      "hello",
			substr: "hello",
			want:   0,
		},
		{
			name:   "substring not found in empty string",
			s:      "",
			substr: "hello",
			want:   -1,
		},

		// Casos con múltiples ocurrencias
		{
			name:   "multiple occurrences",
			s:      "abcabcabc",
			substr: "abc",
			want:   6, // última ocurrencia
		},
		{
			name:   "overlapping patterns",
			s:      "aaaa",
			substr: "aa",
			want:   2, // última posición donde empieza "aa"
		},
		{
			name:   "repeated single character",
			s:      "aaaa",
			substr: "a",
			want:   3, // última 'a'
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LastIndex(tt.s, tt.substr)
			if got != tt.want {
				t.Errorf("LastIndex(%q, %q) = %d, want %d", tt.s, tt.substr, got, tt.want)
			}
		})
	}
}

// TestLastIndexFileExtensions prueba el caso de uso específico de extracción de extensiones
// como se usa en DocPDF.RegisterImageOptions
func TestLastIndexFileExtensions(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     string
		wantPos  int
		hasError bool
	}{
		{
			name:     "simple jpg file",
			filename: "image.jpg",
			want:     "jpg",
			wantPos:  5,
			hasError: false,
		},
		{
			name:     "simple png file",
			filename: "photo.png",
			want:     "png",
			wantPos:  5,
			hasError: false,
		},
		{
			name:     "file with multiple dots",
			filename: "backup.image.jpg",
			want:     "jpg",
			wantPos:  12,
			hasError: false,
		},
		{
			name:     "complex filename with version",
			filename: "document.v2.backup.pdf",
			want:     "pdf",
			wantPos:  18,
			hasError: false,
		},
		{
			name:     "path with directories",
			filename: "/home/user/docs/file.txt",
			want:     "txt",
			wantPos:  20,
			hasError: false,
		},
		{
			name:     "windows path",
			filename: "C:\\Users\\docs\\image.gif",
			want:     "gif",
			wantPos:  19,
			hasError: false,
		},
		{
			name:     "file without extension",
			filename: "README",
			want:     "",
			wantPos:  -1,
			hasError: true,
		},
		{
			name:     "hidden file with extension",
			filename: ".hidden.txt",
			want:     "txt",
			wantPos:  7,
			hasError: false,
		},
		{
			name:     "multiple extensions common case",
			filename: "archive.tar.gz",
			want:     "gz",
			wantPos:  11,
			hasError: false,
		},
		{
			name:     "url-like filename",
			filename: "https://example.com/image.png",
			want:     "png",
			wantPos:  25,
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simular el comportamiento de DocPDF.RegisterImageOptions
			pos := LastIndex(tt.filename, ".")

			if tt.hasError {
				if pos != -1 {
					t.Errorf("Expected no dot found (pos = -1), got pos = %d", pos)
				}
				return
			}

			if pos != tt.wantPos {
				t.Errorf("LastIndex(%q, \".\") = %d, want %d", tt.filename, pos, tt.wantPos)
				return
			}

			if pos >= 0 {
				extension := tt.filename[pos+1:]
				if extension != tt.want {
					t.Errorf("Extension extracted: %q, want %q", extension, tt.want)
				}
			}
		})
	}
}

// TestLastIndexImageTypeDetection simula exactamente el caso de uso de DocPDF
func TestLastIndexImageTypeDetection(t *testing.T) {
	// Función que simula el comportamiento de RegisterImageOptions
	getImageType := func(fileStr string) (string, error) {
		pos := LastIndex(fileStr, ".")
		if pos < 0 {
			return "", Errf("image file has no extension and no type was specified: %s", fileStr)
		}
		return fileStr[pos+1:], nil
	}

	tests := []struct {
		name        string
		filename    string
		expectedExt string
		expectError bool
	}{
		{
			name:        "JPEG image",
			filename:    "photo.jpg",
			expectedExt: "jpg",
			expectError: false,
		},
		{
			name:        "PNG image",
			filename:    "screenshot.png",
			expectedExt: "png",
			expectError: false,
		},
		{
			name:        "GIF image",
			filename:    "animation.gif",
			expectedExt: "gif",
			expectError: false,
		},
		{
			name:        "Complex filename",
			filename:    "user.profile.backup.jpeg",
			expectedExt: "jpeg",
			expectError: false,
		},
		{
			name:        "Path with extension",
			filename:    "/var/www/images/logo.png",
			expectedExt: "png",
			expectError: false,
		},
		{
			name:        "File without extension should error",
			filename:    "imagefile",
			expectedExt: "",
			expectError: true,
		},
		{
			name:        "Hidden file with extension",
			filename:    ".DS_Store.backup.jpg",
			expectedExt: "jpg",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ext, err := getImageType(tt.filename)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for filename %q, but got none", tt.filename)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for filename %q: %v", tt.filename, err)
				return
			}

			if ext != tt.expectedExt {
				t.Errorf("getImageType(%q) = %q, want %q", tt.filename, ext, tt.expectedExt)
			}
		})
	}
}

// Benchmark para comparar performance
func BenchmarkLastIndex(b *testing.B) {
	s := "this is a very long string with multiple occurrences of the word test and more test words"
	substr := "test"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LastIndex(s, substr)
	}
}

func BenchmarkLastIndexFileExtension(b *testing.B) {
	filename := "/very/long/path/to/some/deeply/nested/directory/with/a/very/long/filename.extension"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pos := LastIndex(filename, ".")
		if pos >= 0 {
			_ = filename[pos+1:]
		}
	}
}
