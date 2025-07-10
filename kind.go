package tinystring

// Kind represents the specific Kind of type that a Type represents (private)
// Unified with convert.go Kind, using K prefix for TinyString naming convention
type Kind uint8

// Kind exposes the Kind constants as fields for external use, while keeping the underlying type and values private.
var K = struct {
	Array      Kind
	Bool       Kind
	Bytes      Kind
	Chan       Kind
	Complex128 Kind
	Complex64  Kind
	Float32    Kind
	Float64    Kind
	Func       Kind
	Int        Kind
	Int16      Kind
	Int32      Kind
	Int64      Kind
	Int8       Kind
	Interface  Kind
	Invalid    Kind
	Map        Kind
	Pointer    Kind
	Slice      Kind
	String     Kind
	Struct     Kind
	Uint       Kind
	Uint16     Kind
	Uint32     Kind
	Uint64     Kind
	Uint8      Kind
	Uintptr    Kind
	UnsafePtr  Kind
}{
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27,
}

// kindNames provides string representations for each Kind value
var kindNames = []string{
	"array",
	"bool",
	"[]byte",
	"chan",
	"complex128",
	"complex64",
	"float32",
	"float64",
	"func",
	"int",
	"int16",
	"int32",
	"int64",
	"int8",
	"interface",
	"invalid",
	"map",
	"ptr",
	"slice",
	"string",
	"struct",
	"uint",
	"uint16",
	"uint32",
	"uint64",
	"uint8",
	"uintptr",
	"unsafe.Pointer",
}

// String returns the name of the Kind as a string
func (k Kind) String() string {
	if int(k) >= 0 && int(k) < len(kindNames) {
		return kindNames[k]
	}
	return "invalid"
}
