package tinystring

// looksLikeStruct is a helper function to identify struct types
// This replaces the previous isStruct method with better logic
func (c *conv) looksLikeStruct(value any) bool {
	if value == nil {
		return false
	}

	// Explicitly reject known non-struct types
	switch value.(type) {
	// Basic types - explicitly not structs
	case string, bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return false
	// Slice types - explicitly not structs
	case []string, []bool, []int, []int8, []int16, []int32, []int64, []uint, []uint16, []uint32, []uint64, []float32, []float64, []byte:
		return false
	// Map types - explicitly not structs
	case map[string]string, map[string]any, map[string]int, map[int]string:
		return false
	// Pointer types - explicitly not structs
	case *string, *int, *bool, *float64, **string:
		return false
	// Error type - explicitly not struct
	case error:
		return false
	default:
		// For anything else that's not explicitly known, we conservatively assume it could be a struct
		// This allows structs to be processed, while unknown types will still get an error during reflection
		return true
	}
}
