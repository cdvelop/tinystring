package tinystring

// Kind represents the specific Kind of type that a Type represents (private)
// Unified with convert.go Kind, using K prefix for TinyString naming convention
type kind uint8

// Kind exposes the kind constants as fields for external use, while keeping the underlying type and values private.
var Kind = struct {
	Array      kind
	Bool       kind
	Byte       kind
	Chan       kind
	Complex128 kind
	Complex64  kind
	Err        kind
	Float32    kind
	Float64    kind
	Func       kind
	Int        kind
	Int16      kind
	Int32      kind
	Int64      kind
	Int8       kind
	Interface  kind
	Invalid    kind
	Map        kind
	Pointer    kind
	Slice      kind
	String     kind
	SliceStr   kind
	Struct     kind
	Uint       kind
	Uint16     kind
	Uint32     kind
	Uint64     kind
	Uint8      kind
	Uintptr    kind
	UnsafePtr  kind
}{
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29,
}

// kindNames provides string representations for each kind value
var kindNames = []string{
	"array",
	"bool",
	"[]byte",
	"chan",
	"complex128",
	"complex64",
	"err",
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
	"[]string",
	"struct",
	"uint",
	"uint16",
	"uint32",
	"uint64",
	"uint8",
	"uintptr",
	"unsafe.Pointer",
}

// String returns the name of the kind as a string
func (k kind) String() string {
	if int(k) >= 0 && int(k) < len(kindNames) {
		return kindNames[k]
	}
	return "invalid"
}
