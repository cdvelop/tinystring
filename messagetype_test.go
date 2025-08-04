package tinystring

import (
	"testing"
)

func TestStringTypeDetection(t *testing.T) {
	t.Run("Empty string", func(t *testing.T) {
		msg, msgType := Convert("").StringType()
		if msgType != M.Normal {
			t.Errorf("Expected Normal for empty string, got %v", msgType)
		}
		if msg != "" {
			t.Errorf("Expected empty string, got %q", msg)
		}
	})

	t.Run("Error keywords", func(t *testing.T) {
		errorKeywords := []string{
			"This is an error message",
			"Operation failed",
			"exit status 1",
			"variable undeclared",
			"function undefined",
			"fatal exception",
		}
		for _, keyword := range errorKeywords {
			msg, msgType := Convert(keyword).StringType()
			if msgType != M.Error {
				t.Errorf("Expected Error for keyword %q, got %v", keyword, msgType)
			}
			if msg != keyword {
				t.Errorf("Expected message to be unchanged, got %q", msg)
			}
		}
	})

	t.Run("Success keywords", func(t *testing.T) {
		successKeywords := []string{
			"Success! Operation completed",
			"success",
			"Operation completed",
			"Build successful",
			"Task done",
		}
		for _, keyword := range successKeywords {
			msg, msgType := Convert(keyword).StringType()
			if msgType != M.Success {
				t.Errorf("Expected Success for keyword %q, got %v", keyword, msgType)
			}
			if msg != keyword {
				t.Errorf("Expected message to be unchanged, got %q", msg)
			}
		}
	})

	t.Run("Info keywords", func(t *testing.T) {
		infoKeywords := []string{
			"Info: Starting process",
			"... initializing ...",
			"starting up",
			"initializing system",
		}
		for _, keyword := range infoKeywords {
			_, msgType := Convert(keyword).StringType()
			if msgType != M.Info {
				t.Errorf("Expected Info for keyword %q, got %v", keyword, msgType)
			}
		}
	})

	t.Run("Warning keywords", func(t *testing.T) {
		warningKeywords := []string{
			"Warning: disk space low",
			"warn user",
		}
		for _, keyword := range warningKeywords {
			_, msgType := Convert(keyword).StringType()
			if msgType != M.Warning {
				t.Errorf("Expected Warning for keyword %q, got %v", keyword, msgType)
			}
		}
	})

	t.Run("Normal message", func(t *testing.T) {
		_, msgType := Convert("Hello world").StringType()
		if msgType != M.Normal {
			t.Errorf("Expected Normal for generic message, got %v", msgType)
		}
	})
}
