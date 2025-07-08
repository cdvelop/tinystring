package tinystring

// kind represents the specific kind of type that a Type represents (private)
// Unified with convert.go kind, using K prefix for TinyString naming convention
type kind uint8

const (
	KArray kind = iota
	KBool
	KChan
	KComplex128
	KComplex64
	KErr // Error type (separate from KInvalid)
	KFloat32
	KFloat64
	KFunc
	KInt
	KInt16
	KInt32
	KInt64
	KInt8
	KInterface
	KInvalid
	KMap
	KPointer
	KSlice
	KString
	KSliceStr // Slice of strings
	KStruct
	KUint
	KUint16
	KUint32
	KUint64
	KUint8
	KUintptr
	KUnsafePtr
)

// kindNames provides string representations for each kind value
var kindNames = []string{
	"array",
	"bool",
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
