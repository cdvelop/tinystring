package fmt

// Kind represents the specific Kind of type that a Type represents (private)
// Unified with convert.go Kind, using K prefix for fmt naming convention.
//
// IMPORTANT: The order and values of Kind must NOT be changed.
// These values are used in tinyreflect, a minimal version of reflectlite from the Go standard library.
// Keeping the order and values identical ensures compatibility with code and data shared between tinystring and tinyreflect.
type Kind uint8

// Kind exposes the Kind constants as fields for external use, while keeping the underlying type and values private.
var K = struct {
	Invalid       Kind
	Bool          Kind
	Int           Kind
	Int8          Kind
	Int16         Kind
	Int32         Kind
	Int64         Kind
	Uint          Kind
	Uint8         Kind
	Uint16        Kind
	Uint32        Kind
	Uint64        Kind
	Uintptr       Kind
	Float32       Kind
	Float64       Kind
	Complex64     Kind
	Complex128    Kind
	Array         Kind
	Chan          Kind
	Func          Kind
	Interface     Kind
	Map           Kind
	Pointer       Kind
	Slice         Kind
	String        Kind
	Struct        Kind
	UnsafePointer Kind
}{
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
}

// kindNames provides string representations for each Kind value
var kindNames = []string{
	"invalid",        // 0
	"bool",           // 1
	"int",            // 2
	"int8",           // 3
	"int16",          // 4
	"int32",          // 5
	"int64",          // 6
	"uint",           // 7
	"uint8",          // 8
	"uint16",         // 9
	"uint32",         // 10
	"uint64",         // 11
	"uintptr",        // 12
	"float32",        // 13
	"float64",        // 14
	"complex64",      // 15
	"complex128",     // 16
	"array",          // 17
	"chan",           // 18
	"func",           // 19
	"interface",      // 20
	"map",            // 21
	"ptr",            // 22
	"slice",          // 23
	"string",         // 24
	"struct",         // 25
	"unsafe.Pointer", // 26
}

// String returns the name of the Kind as a string
func (k Kind) String() string {
	if int(k) >= 0 && int(k) < len(kindNames) {
		return kindNames[k]
	}
	return "invalid"
}
