package fmt

import (
	"testing"
	"unsafe"
)

func TestGetStringZeroCopy(t *testing.T) {
	tests := []struct {
		name     string
		dest     BuffDest
		input    string
		expected string
	}{
		{"BuffOut basic", BuffOut, "hello world", "hello world"},
		{"BuffWork basic", BuffWork, "test string", "test string"},
		{"BuffErr basic", BuffErr, "error message", "error message"},
		{"BuffOut empty", BuffOut, "", ""},
		{"BuffWork empty", BuffWork, "", ""},
		{"BuffErr empty", BuffErr, "", ""},
		{"BuffOut unicode", BuffOut, "héllo wörld", "héllo wörld"},
		{"BuffWork unicode", BuffWork, "tëst strïng", "tëst strïng"},
		{"BuffErr unicode", BuffErr, "ërror mëssagë", "ërror mëssagë"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := GetConv()
			defer c.PutConv()

			// Write input to the specified buffer
			c.WrString(tt.dest, tt.input)

			// Get zero-copy string
			result := c.GetStringZeroCopy(tt.dest)

			// Check correctness
			if result != tt.expected {
				t.Errorf("GetStringZeroCopy(%v) = %q, want %q", tt.dest, result, tt.expected)
			}

			// Verify zero-allocation: string data should point to buffer
			if len(tt.input) > 0 {
				var bufferPtr unsafe.Pointer
				switch tt.dest {
				case BuffOut:
					bufferPtr = unsafe.Pointer(&c.out[0])
				case BuffWork:
					bufferPtr = unsafe.Pointer(&c.work[0])
				case BuffErr:
					bufferPtr = unsafe.Pointer(&c.err[0])
				}
				stringPtr := unsafe.Pointer(unsafe.StringData(result))
				if stringPtr != bufferPtr {
					t.Errorf("GetStringZeroCopy did not return zero-copy string for %v; pointers differ", tt.dest)
				}
			}

			// Ensure no corruption: string should remain valid
			// Since it's zero-copy, we don't modify buffer after getting string
			if result != tt.expected {
				t.Errorf("String corrupted after GetStringZeroCopy for %v", tt.dest)
			}
		})
	}
}

func BenchmarkGetStringZeroCopy(b *testing.B) {
	c := GetConv()
	defer c.PutConv()

	// Prepare buffer with test data
	testString := "benchmark test string for zero copy"
	c.WrString(BuffOut, testString)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		result := c.GetStringZeroCopy(BuffOut)
		_ = result // Prevent optimization
	}
}
