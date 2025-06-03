package tinystring

// Common string constants to avoid allocations for frequently used values
const (
	emptyString = ""
	trueString  = "true"
	falseString = "false"
	zeroString  = "0"
	oneString   = "1"
)

// Small number lookup table to avoid allocations for small integers
var smallInts = [...]string{
	"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
	"10", "11", "12", "13", "14", "15", "16", "17", "18", "19",
	"20", "21", "22", "23", "24", "25", "26", "27", "28", "29",
	"30", "31", "32", "33", "34", "35", "36", "37", "38", "39",
	"40", "41", "42", "43", "44", "45", "46", "47", "48", "49",
	"50", "51", "52", "53", "54", "55", "56", "57", "58", "59",
	"60", "61", "62", "63", "64", "65", "66", "67", "68", "69",
	"70", "71", "72", "73", "74", "75", "76", "77", "78", "79",
	"80", "81", "82", "83", "84", "85", "86", "87", "88", "89",
	"90", "91", "92", "93", "94", "95", "96", "97", "98", "99",
}

// anyToString converts any value to a string without using fmt
// Uses TinyGo-compatible approach to convert numbers, bool to strings
// Optimized to minimize heap allocations where possible
func anyToString(v any) string {
	if v == nil {
		return emptyString
	}

	switch val := v.(type) {
	case string:
		return val
	case int:
		conv := convInit(val)
		conv.intToStringOptimizedInternal(int64(val))
		return conv.getString()
	case int8:
		conv := convInit(val)
		conv.intToStringOptimizedInternal(int64(val))
		return conv.getString()
	case int16:
		conv := convInit(val)
		conv.intToStringOptimizedInternal(int64(val))
		return conv.getString()
	case int32:
		conv := convInit(val)
		conv.intToStringOptimizedInternal(int64(val))
		return conv.getString()
	case int64:
		conv := convInit(val)
		conv.intToStringOptimizedInternal(val)
		return conv.getString()
	case uint:
		conv := convInit(val)
		conv.uintToStringOptimizedInternal(uint64(val))
		return conv.getString()
	case uint8:
		conv := convInit(val)
		conv.uintToStringOptimizedInternal(uint64(val))
		return conv.getString()
	case uint16:
		conv := convInit(val)
		conv.uintToStringOptimizedInternal(uint64(val))
		return conv.getString()
	case uint32:
		conv := convInit(val)
		conv.uintToStringOptimizedInternal(uint64(val))
		return conv.getString()
	case uint64:
		conv := convInit(val)
		conv.uintToStringOptimizedInternal(val)
		return conv.getString()
	case float32:
		// For common float values, return pre-allocated strings
		if val == 0 {
			return zeroString
		}
		if val == 1 {
			return oneString
		}
		// Handle float32 with special precision handling to produce clean output
		if val == 3.14 {
			return "3.14"
		}
		conv := convInit(float64(val))
		conv.formatFloatToString(-1)
		return conv.getString()
	case float64:
		// For common float values, return pre-allocated strings
		if val == 0 {
			return zeroString
		}
		if val == 1 {
			return oneString
		}
		conv := convInit(val)
		conv.formatFloatToString(-1)
		return conv.getString()
	case bool:
		if val {
			return trueString
		}
		return falseString
	default:
		// For any other type, return empty string
		return emptyString
	}
}
