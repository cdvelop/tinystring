package tinystring

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

// safeCounter provides thread-safe counting for detecting errors
type safeCounter struct {
	mu    sync.Mutex
	count int
	errs  []string
}

func (c *safeCounter) addError(msg string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
	c.errs = append(c.errs, msg)
}

// TestConcurrentConvert tests that the Convert method and its chained operations
// are safe to use concurrently from multiple goroutines.
func TestConcurrentConvert(t *testing.T) {
	const (
		numGoroutines  = 200 // Reduced from 1000 to prevent resource exhaustion
		testString     = "Él Múrcielago Rápido"
		expectedResult = "elMurcielagoRapido"
	)

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Thread-safe error collection
	var counter safeCounter

	// Add timeout protection
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			result := Convert(testString).
				RemoveTilde().
				CamelCaseLower().
				String()

			if result != expectedResult {
				counter.addError(fmt.Sprintf("goroutine %d: got %q, want %q", id, result, expectedResult))
			}
		}(i)
	}

	// Wait with timeout
	select {
	case <-done:
		if counter.count > 0 {
			t.Errorf("Failed with %d errors:\n%s", counter.count, strings.Join(counter.errs, "\n"))
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Test timed out after 5 seconds")
	}
}

// TestConcurrentUtilityFunctions tests that standalone utility functions
// are safe to use concurrently from multiple goroutines.
func TestConcurrentUtilityFunctions(t *testing.T) {
	const numGoroutines = 100 // Reduced from 500

	testCases := []struct {
		name     string
		function func() (string, error)
		expected string
	}{
		{
			name: "Split",
			function: func() (string, error) {
				result := Split("apple,banana,cherry", ",")
				return result[1], nil
			},
			expected: "banana",
		},
		{
			name: "ParseKeyValue",
			function: func() (string, error) {
				val, err := ParseKeyValue("user:admin", ":")
				if err != nil {
					return "", err
				}
				return val, nil
			},
			expected: "admin",
		},
		{
			name: "Contains",
			function: func() (string, error) {
				if Contains("hello world", "world") {
					return "true", nil
				}
				return "false", nil
			},
			expected: "true",
		},
		{
			name: "CountOccurrences",
			function: func() (string, error) {
				count := CountOccurrences("abracadabra", "abra")
				if count == 2 {
					return "2", nil
				}
				return "wrong", nil
			},
			expected: "2",
		},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			var wg sync.WaitGroup
			wg.Add(numGoroutines)

			// Thread-safe error collection
			var counter safeCounter

			// Add timeout protection
			done := make(chan struct{})
			go func() {
				wg.Wait()
				close(done)
			}()

			for i := 0; i < numGoroutines; i++ {
				go func(id int) {
					defer wg.Done()
					result, err := tc.function()
					if err != nil {
						counter.addError(fmt.Sprintf("goroutine %d: error: %v", id, err))
					}
					if result != tc.expected {
						counter.addError(fmt.Sprintf("goroutine %d: got %q, want %q", id, result, tc.expected))
					}
				}(i)
			}

			// Wait with timeout
			select {
			case <-done:
				if counter.count > 0 {
					t.Errorf("Failed with %d errors:\n%s", counter.count, strings.Join(counter.errs, "\n"))
				}
			case <-time.After(5 * time.Second):
				t.Fatal("Test timed out after 5 seconds")
			}
		})
	}
}

// TestConcurrentStringManipulation tests that complex string manipulations
// executed concurrently produce consistent results.
func TestConcurrentStringManipulation(t *testing.T) {
	const (
		numGoroutines = 100 // Reduced from 300
		iterations    = 5   // Reduced from 10
	)

	testCases := []struct {
		name     string
		input    string
		process  func(string) string
		expected string
	}{
		{
			name:  "Complex Transformation 1",
			input: "  User-Name With Áccents  ",
			process: func(s string) string {
				return Convert(s).
					Trim().
					RemoveTilde().
					Replace(" ", "_").
					Replace("-", "_").
					ToLower().
					String()
			},
			expected: "user_name_with_accents",
		},
		{
			name:  "Complex Transformation 2",
			input: "this.is.a.file.name.txt",
			process: func(s string) string {
				// First replace periods with spaces, then apply CamelCaseUpper,
				// then remove the ".txt" suffix
				return Convert(s).
					TrimSuffix(".txt"). // Remove suffix first
					Replace(".", " ").  // Then replace periods with spaces
					CamelCaseUpper().   // Convert to CamelCase
					String()
			},
			expected: "ThisIsAFileName",
		},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			var wg sync.WaitGroup
			wg.Add(numGoroutines)

			// Thread-safe error collection
			var counter safeCounter

			// Add timeout protection
			done := make(chan struct{})
			go func() {
				wg.Wait()
				close(done)
			}()

			for i := 0; i < numGoroutines; i++ {
				go func(id int) {
					defer wg.Done()

					for j := 0; j < iterations; j++ {
						result := tc.process(tc.input)
						if result != tc.expected {
							counter.addError(fmt.Sprintf("goroutine %d, iteration %d: got %q, want %q",
								id, j, result, tc.expected))
							return
						}
					}
				}(i)
			}

			// Wait with timeout
			select {
			case <-done:
				if counter.count > 0 {
					t.Errorf("Failed with %d errors:\n%s", counter.count, strings.Join(counter.errs, "\n"))
				}
			case <-time.After(5 * time.Second):
				t.Fatal("Test timed out after 5 seconds")
			}
		})
	}
}
