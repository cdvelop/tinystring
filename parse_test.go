package tinystring_test

import (
	"strings"
	"testing"

	. "github.com/cdvelop/tinystring"
)

func TestParseKeyValue(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		delimiter   string
		wantValue   string
		wantErrText string
	}{
		{
			name:        "Basic key-value with default delimiter",
			input:       "name:John",
			delimiter:   "",
			wantValue:   "John",
			wantErrText: "",
		},
		{
			name:        "No delimiter in string",
			input:       "invalid-string",
			delimiter:   "",
			wantValue:   "",
			wantErrText: "delimiter ':' not found in string invalid-string",
		},
		{
			name:        "Custom delimiter",
			input:       "age=30",
			delimiter:   "=",
			wantValue:   "30",
			wantErrText: "",
		},
		{
			name:        "Value contains delimiter",
			input:       "address:123 Main St:Apt 4",
			delimiter:   "",
			wantValue:   "123 Main St:Apt 4",
			wantErrText: "",
		},
		{
			name:        "Empty input",
			input:       "",
			delimiter:   "",
			wantValue:   "",
			wantErrText: "delimiter ':' not found in string ",
		},
		{
			name:        "Only delimiter",
			input:       ":",
			delimiter:   "",
			wantValue:   "",
			wantErrText: "",
		},
		{
			name:        "Multi-character delimiter",
			input:       "key=>value",
			delimiter:   "=>",
			wantValue:   "value",
			wantErrText: "",
		},
		{
			name:        "Missing custom delimiter",
			input:       "key:value",
			delimiter:   "=",
			wantValue:   "",
			wantErrText: "delimiter '=' not found in string key:value",
		},
		{
			name:        "Empty delimiter uses default",
			input:       "name:John",
			delimiter:   "",
			wantValue:   "John",
			wantErrText: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var delimiters []string
			if tc.delimiter != "" {
				delimiters = append(delimiters, tc.delimiter)
			}

			gotValue, gotErr := ParseKeyValue(tc.input, delimiters...)

			if gotValue != tc.wantValue {
				t.Errorf("ParseKeyValue() value = %q, want %q", gotValue, tc.wantValue)
			}

			// Check error
			if tc.wantErrText == "" {
				if gotErr != nil {
					t.Errorf("ParseKeyValue() error = %v, want nil", gotErr)
				}
			} else {
				if gotErr == nil {
					t.Errorf("ParseKeyValue() error = nil, want error containing %q", tc.wantErrText)
				} else if !strings.Contains(gotErr.Error(), tc.wantErrText) {
					t.Errorf("ParseKeyValue() error = %v, want error containing %q", gotErr, tc.wantErrText)
				}
			}
		})
	}
}
