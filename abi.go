package tinystring

import "unsafe"

// ABI types consolidated from internal/abi for TinyString JSON functionality
// This eliminates duplication with convert.go types and provides single source of truth

// kind represents the specific kind of type that a Type represents (private)
// Unified with convert.go vTpe, using tp prefix for TinyString naming convention
type kind uint8

const (
	tpInvalid kind = iota
	tpBool
	tpInt
	tpInt8
	tpInt16
	tpInt32
	tpInt64
	tpUint
	tpUint8
	tpUint16
	tpUint32
	tpUint64
	tpUintptr
	tpFloat32
	tpFloat64
	tpComplex64
	tpComplex128
	tpArray
	tpChan
	tpFunc
	tpInterface
	tpMap
	tpPointer
	tpSlice
	tpString
	tpStruct
	tpUnsafePointer

	// TinyString specific types - separate values to avoid conflicts
	tpStrSlice // String slice type (separate from tpSlice)
	tpStrPtr   // String pointer type (separate from tpPointer)
	tpErr      // Error type (separate from tpInvalid)
)

const (
	kindDirectIface = 1 << 5
	kindGCProg      = 1 << 6 // Type.gc points to GC program
	kindMask        = (1 << 5) - 1
)

// String returns the name of k
func (k kind) String() string {
	if int(k) < len(kindNames) {
		return kindNames[k]
	}
	return kindNames[0]
}

var kindNames = []string{
	tpInvalid:       "invalid",
	tpBool:          "bool",
	tpInt:           "int",
	tpInt8:          "int8",
	tpInt16:         "int16",
	tpInt32:         "int32",
	tpInt64:         "int64",
	tpUint:          "uint",
	tpUint8:         "uint8",
	tpUint16:        "uint16",
	tpUint32:        "uint32",
	tpUint64:        "uint64",
	tpUintptr:       "uintptr",
	tpFloat32:       "float32",
	tpFloat64:       "float64",
	tpComplex64:     "complex64",
	tpComplex128:    "complex128",
	tpArray:         "array",
	tpChan:          "chan",
	tpFunc:          "func",
	tpInterface:     "interface",
	tpMap:           "map",
	tpPointer:       "ptr",
	tpSlice:         "slice",
	tpString:        "string",
	tpStruct:        "struct",
	tpUnsafePointer: "unsafe.Pointer",
}

// tFlag is used by a Type to signal what extra type information is available
type tFlag uint8

const (
	tflagUncommon       tFlag = 1 << 0
	tflagExtraStar      tFlag = 1 << 1
	tflagNamed          tFlag = 1 << 2
	tflagRegularMemory  tFlag = 1 << 3
	tflagUnrolledBitmap tFlag = 1 << 4
)

// nameOff is the offset to a name from moduledata.types
type nameOff int32

// typeOff is the offset to a type from moduledata.types
type typeOff int32

// textOff is an offset from the top of a text section
type textOff int32

// refType is the runtime representation of a Go type (adapted from internal/abi)
// Used for JSON struct inspection and field access
type refType struct {
	size       uintptr
	ptrBytes   uintptr // number of (prefix) bytes in the type that can contain pointers
	hash       uint32  // hash of type; avoids computation in hash tables
	tflag      tFlag   // extra type information flags
	align      uint8   // alignment of variable with this type
	fieldAlign uint8   // alignment of struct field with this type
	kind       uint8   // enumeration for C
	// function for comparing objects of this type
	equal     func(unsafe.Pointer, unsafe.Pointer) bool
	gcData    *byte
	str       nameOff // string form
	ptrToThis typeOff // type for pointer to this type, may be zero
}

// refPtrType represents a pointer type
type refPtrType struct {
	refType
	elem *refType // pointer element (pointed at) type
}

// refStructField represents a field in a struct type
type refStructField struct {
	name   refName  // name is always non-empty
	typ    *refType // type of field
	offset uintptr  // byte offset of field
}

// refStructType represents a struct type
type refStructType struct {
	refType
	pkgPath refName
	fields  []refStructField
}

// refSliceType represents a slice type
type refSliceType struct {
	refType
	elem *refType // slice element type
}

// refName is an encoded type name with optional extra data
type refName struct {
	bytes *byte
}

// Kind returns the kind of type
func (t *refType) Kind() kind {
	return kind(t.kind & kindMask)
}

// Size returns the size of data with type t
func (t *refType) Size() uintptr {
	return t.size
}

// Elem returns the element type for pointer, array, channel, map, or slice types
func (t *refType) Elem() *refType {
	switch t.Kind() {
	case tpPointer:
		pt := (*refPtrType)(unsafe.Pointer(t))
		return pt.elem
	case tpArray:
		at := (*refArrayType)(unsafe.Pointer(t))
		return at.elem
	case tpSlice:
		st := (*refSliceType)(unsafe.Pointer(t))
		return st.elem
	// Add other cases as needed
	default:
		return nil
	}
}

// refArrayType represents an array type
type refArrayType struct {
	refType
	elem *refType // array element type
	len  uintptr
}

// NumField returns the number of fields in a struct type
func (t *refStructType) NumField() int {
	return len(t.fields)
}

// Field returns the i'th field of the struct type
func (t *refStructType) Field(i int) refStructField {
	if i < 0 || i >= len(t.fields) {
		panic("reflect: Field index out of range")
	}
	return t.fields[i]
}

// Name returns the name string for refName
func (n refName) Name() string {
	if n.bytes == nil {
		return ""
	}
	i, l := n.readVarint(1)
	return unsafe.String(n.dataChecked(1+i, "non-empty string"), l)
}

// IsExported returns whether the name is exported
func (n refName) IsExported() bool {
	return (*n.bytes)&(1<<0) != 0
}

// readVarint parses a varint as encoded by encoding/binary
func (n refName) readVarint(off int) (int, int) {
	v := 0
	for i := 0; ; i++ {
		x := *n.dataChecked(off+i, "read varint")
		v += int(x&0x7f) << (7 * i)
		if x&0x80 == 0 {
			return i + 1, v
		}
	}
}

// dataChecked does pointer arithmetic on n's bytes
func (n refName) dataChecked(off int, whySafe string) *byte {
	return (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(n.bytes)) + uintptr(off)))
}

// Global cache for struct type information
// Using slice instead of map for TinyGo compatibility
var objects []object

// object represents cached information about a struct type
type object struct {
	snakeName string   // snake_case name of the type
	refType   *refType // reference to the type information
	fields    []field  // cached field information
}

// field represents cached information about a struct field
type field struct {
	name      string   // original field name (e.g., "BirthDate")
	snakeName string   // snake_case name of the field (e.g., "birth_date")
	refType   *refType // type of the field
	offset    uintptr  // byte offset in the struct
	index     int      // field index
}

// findObject searches for cached object by snake_case name
func findObject(snakeName string) *object {
	for i := range objects {
		if objects[i].snakeName == snakeName {
			return &objects[i]
		}
	}
	return nil
}

// addObject adds a new object to the cache
func addObject(obj object) {
	objects = append(objects, obj)
}

// clearObjectCache clears the global object cache - useful for testing
func clearObjectCache() {
	objects = objects[:0] // Clear slice while preserving capacity
}
