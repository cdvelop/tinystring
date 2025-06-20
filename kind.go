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
