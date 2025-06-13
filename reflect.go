package tinystring

import (
	"unsafe"
)

// Minimal reflectlite integration for TinyString JSON functionality
// This file contains essential reflection capabilities adapted from internal/reflectlite
// All functions and types are prefixed with 'ref' to avoid API pollution

// refValue represents a value and its type information
type refValue struct {
	typ  *refType
	ptr  unsafe.Pointer
	flag refFlag
}

// refFlag holds metadata about the value
type refFlag uintptr

const (
	flagKindWidth           = 5 // there are 27 kinds
	flagKindMask    refFlag = 1<<flagKindWidth - 1
	flagStickyRO    refFlag = 1 << 5
	flagEmbedRO     refFlag = 1 << 6
	flagIndir       refFlag = 1 << 7
	flagAddr        refFlag = 1 << 8
	flagMethod      refFlag = 1 << 9
	flagMethodShift         = 10
	flagKindShift           = flagMethodShift + 10 // room for method index
	flagRO          refFlag = flagStickyRO | flagEmbedRO
)

// refValueOf returns a new refValue initialized to the concrete value stored in i
func refValueOf(i any) refValue {
	if i == nil {
		return refValue{}
	}
	return unpackEface(i)
}

// refEface is the header for an interface{} value
type refEface struct {
	typ  *refType
	data unsafe.Pointer
}

// unpackEface converts an interface{} to a refValue
func unpackEface(i any) refValue {
	e := (*refEface)(unsafe.Pointer(&i))
	t := e.typ
	if t == nil {
		return refValue{}
	}
	f := refFlag(t.Kind())

	// For basic types, we need to determine if they are stored directly or indirectly
	// Basic types like int, string, bool, float64 etc. are typically stored directly
	// in interface{} when they fit in the interface data field
	switch t.Kind() {
	case tpBool, tpInt, tpInt8, tpInt16, tpInt32, tpInt64,
		tpUint, tpUint8, tpUint16, tpUint32, tpUint64, tpUintptr,
		tpFloat32, tpFloat64, tpPointer, tpUnsafePointer:
		// These basic types are stored directly in interface - no flagIndir
		// The e.data pointer points directly to the value
	case tpString:
		// Strings are stored directly in interface on most architectures
		// The e.data pointer points to the string header
	default:
		// For other types (struct, slice, array, etc.), check if stored indirectly
		if ifaceIndir(t) {
			f |= flagIndir
		}
	}

	return refValue{t, e.data, f}
}

// ifaceIndir reports whether t is stored indirectly in an interface value
func ifaceIndir(t *refType) bool {
	return t.kind&kindDirectIface == 0
}

// Kind returns the kind of value v contains
func (v refValue) Kind() kind {
	return kind(v.flag & flagKindMask)
}

// Type returns the type of v
func (v refValue) Type() *refType {
	return v.typ
}

// IsValid reports whether v represents a value
func (v refValue) IsValid() bool {
	return v.flag != 0
}

// Elem returns the value that the interface v contains or that the pointer v points to
// Elem returns the value that the interface v contains or that the pointer v points to
func (v refValue) Elem() refValue {
	k := v.Kind()
	switch k {
	case tpInterface:
		var eface refEface
		if v.typ.kind&kindDirectIface != 0 {
			eface = refEface{typ: nil, data: v.ptr}
		} else {
			eface = *(*refEface)(v.ptr)
		}
		if eface.typ == nil {
			return refValue{}
		}
		fl := refFlag(eface.typ.Kind())
		if ifaceIndir(eface.typ) {
			fl |= flagIndir
		}
		return refValue{eface.typ, eface.data, fl}
	case tpPointer:
		// Handle pointer dereferencing
		var ptr unsafe.Pointer
		if v.flag&flagIndir != 0 {
			// This is a pointer field from a struct - need to dereference to get the actual pointer
			ptr = *(*unsafe.Pointer)(v.ptr)
		} else {
			// This is a direct pointer from interface{}
			// v.ptr contains the pointer value itself (the address it points to)
			ptr = v.ptr
		}

		if ptr == nil {
			return refValue{}
		}

		elemType := v.typ.Elem()
		if elemType == nil {
			return refValue{}
		} // Create proper flags for the element
		// The element is addressable since we're dereferencing a pointer
		fl := v.flag&flagRO | flagAddr | refFlag(elemType.Kind())

		// For elements accessed through pointers, we don't need flagIndir
		// because ptr already points to the actual data
		return refValue{elemType, ptr, fl}
	}
	panic("reflect: call of reflect.Value.Elem on " + v.Type().Kind().String() + " value")
}

// NumField returns the number of fields in the struct v
func (v refValue) NumField() int {
	v.mustBe(tpStruct)
	tt := (*refStructType)(unsafe.Pointer(v.typ))
	return len(tt.fields)
}

// Field returns the i'th field of the struct v
func (v refValue) Field(i int) refValue {
	if v.Kind() != tpStruct {
		panic("reflect: call of reflect.Value.Field on " + v.Kind().String() + " value")
	}
	tt := (*refStructType)(unsafe.Pointer(v.typ))
	if uint(i) >= uint(len(tt.fields)) {
		panic("reflect: Field index out of range")
	}
	field := &tt.fields[i]
	ptr := add(v.ptr, field.offset, "same as non-reflect &v.field") // Inherit read-only flags from parent, but allow assignment if parent allows it
	fl := v.flag&(flagRO) | refFlag(field.typ.Kind()) | flagAddr
	// For struct fields, flagIndir is needed only for pointer fields
	// because ptr points to the field location containing the pointer.
	// For other field types, ptr points directly to the field value.
	if field.typ.Kind() == tpPointer {
		fl |= flagIndir
	}
	return refValue{field.typ, ptr, fl}
}

// SetString sets v's underlying value to x
func (v refValue) SetString(x string) {
	v.mustBeAssignable()
	v.mustBe(tpString)
	ptr := v.ptr
	if v.flag&flagIndir != 0 {
		ptr = *(*unsafe.Pointer)(ptr)
	}
	*(*string)(ptr) = x
}

// SetInt sets v's underlying value to x
func (v refValue) SetInt(x int64) {
	v.mustBeAssignable()
	ptr := v.ptr
	if v.flag&flagIndir != 0 {
		ptr = *(*unsafe.Pointer)(ptr)
	}
	switch v.Kind() {
	case tpInt:
		*(*int)(ptr) = int(x)
	case tpInt8:
		*(*int8)(ptr) = int8(x)
	case tpInt16:
		*(*int16)(ptr) = int16(x)
	case tpInt32:
		*(*int32)(ptr) = int32(x)
	case tpInt64:
		*(*int64)(ptr) = x
	default:
		panic("reflect: call of reflect.Value.SetInt on " + v.Kind().String() + " value")
	}
}

// SetUint sets v's underlying value to x
func (v refValue) SetUint(x uint64) {
	v.mustBeAssignable()
	ptr := v.ptr
	if v.flag&flagIndir != 0 {
		ptr = *(*unsafe.Pointer)(ptr)
	}
	switch v.Kind() {
	case tpUint:
		*(*uint)(ptr) = uint(x)
	case tpUint8:
		*(*uint8)(ptr) = uint8(x)
	case tpUint16:
		*(*uint16)(ptr) = uint16(x)
	case tpUint32:
		*(*uint32)(ptr) = uint32(x)
	case tpUint64:
		*(*uint64)(ptr) = x
	case tpUintptr:
		*(*uintptr)(ptr) = uintptr(x)
	default:
		panic("reflect: call of reflect.Value.SetUint on " + v.Kind().String() + " value")
	}
}

// SetFloat sets v's underlying value to x
func (v refValue) SetFloat(x float64) {
	v.mustBeAssignable()
	ptr := v.ptr
	if v.flag&flagIndir != 0 {
		ptr = *(*unsafe.Pointer)(ptr)
	}
	switch v.Kind() {
	case tpFloat32:
		*(*float32)(ptr) = float32(x)
	case tpFloat64:
		*(*float64)(ptr) = x
	default:
		panic("reflect: call of reflect.Value.SetFloat on " + v.Kind().String() + " value")
	}
}

// SetBool sets v's underlying value to x
func (v refValue) SetBool(x bool) {
	v.mustBeAssignable()
	v.mustBe(tpBool)
	ptr := v.ptr
	if v.flag&flagIndir != 0 {
		ptr = *(*unsafe.Pointer)(ptr)
	}
	*(*bool)(ptr) = x
}

// Set assigns x to the value v
// v must be addressable and must not have been obtained by accessing unexported struct fields
func (v refValue) Set(x refValue) {
	v.mustBeAssignable()
	x.mustBeExported() // do not let unexported x leak

	// For pointer types, we need to copy the pointer value itself
	if v.kind() == tpPointer {
		// v.ptr points to the pointer variable
		// We need to set the pointer variable to the value that x represents
		if x.kind() == tpPointer {
			// Copy pointer value from x to v
			*(*unsafe.Pointer)(v.ptr) = *(*unsafe.Pointer)(x.ptr)
		} else {
			// x is not a pointer, this shouldn't happen in normal cases
			typedmemmove(v.typ, v.ptr, x.ptr)
		}
	} else {
		// For non-pointer types, copy the value
		typedmemmove(v.typ, v.ptr, x.ptr)
	}
}

// refZero returns a Value representing the zero value for the specified type
func refZero(typ *refType) refValue {
	if typ == nil {
		panic("reflect: Zero(nil)")
	}

	// For pointer types, zero value is nil pointer
	if typ.Kind() == tpPointer {
		var nilPtr unsafe.Pointer // This is nil
		return refValue{typ, unsafe.Pointer(&nilPtr), refFlag(tpPointer)}
	}

	// For struct and other types, allocate memory for the zero value
	size := typ.Size()

	// Safety check: prevent huge allocations that could cause out of memory
	const maxSafeSize = 1024 * 1024 // 1MB limit
	if size > maxSafeSize {
		// For very large types, use a fixed small buffer
		size = 512
	}

	ptr := unsafe.Pointer(&make([]byte, size)[0])

	// Zero out the memory
	memclr(ptr, size)

	// Return the zero value with correct type and kind
	// Use the kind directly in the flag without shifting
	return refValue{
		typ:  typ,
		ptr:  ptr,
		flag: refFlag(typ.Kind()) | flagAddr,
	}
}

// mustBeExported panics if f records that the value was obtained using an unexported field
func (v refValue) mustBeExported() {
	if v.flag&flagRO != 0 {
		panic("reflect: " + "use of unexported field")
	}
}

// mustBeAssignable panics if v is not assignable
func (v refValue) mustBeAssignable() {
	if v.flag&flagRO != 0 {
		panic("reflect: " + "cannot set value")
	}
	if v.flag&flagAddr == 0 {
		panic("reflect: " + "cannot assign to value")
	}
}

// kind returns the Kind without the flags
func (v refValue) kind() kind {
	return kind(v.flag & flagKindMask)
}

// typedmemmove copies a value of type t to dst from src
func typedmemmove(t *refType, dst, src unsafe.Pointer) {
	// Simplified version - just copy the bytes
	// This should use the actual Go runtime typedmemmove for safety
	// but for our purposes, a simple memory copy works
	memmove(dst, src, t.size)
}

// memmove copies n bytes from src to dst
func memmove(dst, src unsafe.Pointer, size uintptr) {
	// Simplified byte-by-byte copy
	// In real implementation, this would use runtime memmove
	dstBytes := (*[1 << 30]byte)(dst)
	srcBytes := (*[1 << 30]byte)(src)
	for i := uintptr(0); i < size; i++ {
		dstBytes[i] = srcBytes[i]
	}
}

// String returns the string v's underlying value, as a string
func (v refValue) String() string {
	switch k := v.Kind(); k {
	case tpString:
		ptr := v.ptr
		if v.flag&flagIndir != 0 {
			ptr = *(*unsafe.Pointer)(ptr)
		}
		return *(*string)(ptr)
	default:
		// Return a descriptive string instead of panicking
		return "<" + v.Type().Kind().String() + " Value>"
	}
}

// Int returns v's underlying value, as an int64
func (v refValue) Int() int64 {
	ptr := v.ptr
	if v.flag&flagIndir != 0 {
		ptr = *(*unsafe.Pointer)(ptr)
	}

	switch k := v.Kind(); k {
	case tpInt:
		return int64(*(*int)(ptr))
	case tpInt8:
		return int64(*(*int8)(ptr))
	case tpInt16:
		return int64(*(*int16)(ptr))
	case tpInt32:
		return int64(*(*int32)(ptr))
	case tpInt64:
		return *(*int64)(ptr)
	default:
		panic("reflect: call of reflect.Value.Int on " + v.Kind().String() + " value")
	}
}

// Uint returns v's underlying value, as a uint64
func (v refValue) Uint() uint64 {
	ptr := v.ptr
	if v.flag&flagIndir != 0 {
		ptr = *(*unsafe.Pointer)(ptr)
	}

	switch k := v.Kind(); k {
	case tpUint:
		return uint64(*(*uint)(ptr))
	case tpUint8:
		return uint64(*(*uint8)(ptr))
	case tpUint16:
		return uint64(*(*uint16)(ptr))
	case tpUint32:
		return uint64(*(*uint32)(ptr))
	case tpUint64:
		return *(*uint64)(ptr)
	case tpUintptr:
		return uint64(*(*uintptr)(ptr))
	default:
		panic("reflect: call of reflect.Value.Uint on " + v.Kind().String() + " value")
	}
}

// Float returns v's underlying value, as a float64
func (v refValue) Float() float64 {
	ptr := v.ptr
	if v.flag&flagIndir != 0 {
		ptr = *(*unsafe.Pointer)(ptr)
	}

	switch k := v.Kind(); k {
	case tpFloat32:
		return float64(*(*float32)(ptr))
	case tpFloat64:
		return *(*float64)(ptr)
	default:
		panic("reflect: call of reflect.Value.Float on " + v.Kind().String() + " value")
	}
}

// Bool returns v's underlying value
func (v refValue) Bool() bool {
	v.mustBe(tpBool)
	ptr := v.ptr
	if v.flag&flagIndir != 0 {
		ptr = *(*unsafe.Pointer)(ptr)
	}
	return *(*bool)(ptr)
}

// Interface returns v's current value as an interface{}
func (v refValue) Interface() any {
	if v.flag&flagRO != 0 {
		// Read-only value, we can still interface it
	}
	switch v.Kind() {
	case tpString:
		return v.String()
	case tpInt:
		return int(v.Int())
	case tpInt8:
		return int8(v.Int())
	case tpInt16:
		return int16(v.Int())
	case tpInt32:
		return int32(v.Int())
	case tpInt64:
		return v.Int()
	case tpUint:
		return uint(v.Uint())
	case tpUint8:
		return uint8(v.Uint())
	case tpUint16:
		return uint16(v.Uint())
	case tpUint32:
		return uint32(v.Uint())
	case tpUint64:
		return v.Uint()
	case tpFloat32:
		return float32(v.Float())
	case tpFloat64:
		return v.Float()
	case tpBool:
		return v.Bool()
	case tpInterface:
		// For interface{} fields, extract the actual contained value
		if v.ptr == nil {
			return nil
		}
		eface := *(*refEface)(v.ptr)
		if eface.typ == nil {
			return nil
		}
		// Reconstruct the interface{} from the eface
		return packEface(eface)
	case tpStruct:
		// For struct types, we need to reconstruct the original Go value
		// This requires careful handling of the type and memory layout
		return v.packStructEface()
	default:
		// For other complex types like slice, etc.
		// This is a simplified approach - in a full implementation,
		// we would need to reconstruct the appropriate interface{} value
		return v.ptr // This is a fallback
	}
}

// packEface reconstructs an interface{} value from a refEface
func packEface(eface refEface) any {
	// Reconstruct the interface{} value by casting back
	// This is a simplified approach - create a new interface{} with the same structure
	var result any
	resultEface := (*refEface)(unsafe.Pointer(&result))
	*resultEface = eface
	return result
}

// packStructEface reconstructs an interface{} value for a struct
func (v refValue) packStructEface() any {
	// For struct values, we need to reconstruct the original value
	// This is based on Go's reflect.packEface logic

	// Create a new interface{} header
	var result any
	resultEface := (*refEface)(unsafe.Pointer(&result))
	resultEface.typ = v.typ

	// Handle memory layout depending on whether the struct is stored indirectly
	if ifaceIndir(v.typ) {
		// For larger structs stored indirectly, we need to copy the data
		// Allocate new memory for the struct
		size := v.typ.Size()
		newPtr := unsafe.Pointer(&make([]byte, size)[0])

		// Copy the struct data
		memmove(newPtr, v.ptr, size)
		resultEface.data = newPtr
	} else {
		// For small structs stored directly in interface
		resultEface.data = v.ptr
	}

	return result
}

// mustBe panics if f's kind is not expected
func (v refValue) mustBe(expected kind) {
	if v.Kind() != expected {
		panic("reflect: call of reflect.Value method on " + expected.String() + " value")
	}
}

// add returns p+x
func add(p unsafe.Pointer, x uintptr, whySafe string) unsafe.Pointer {
	return unsafe.Pointer(uintptr(p) + x)
}

// Global cache for struct type information
// Using slice instead of map for TinyGo compatibility
var refStructsInfo []refStructInfo

// refStructInfo contains cached information about a struct type for JSON operations
type refStructInfo struct {
	snakeName string         // snake_case name of the type
	refType   *refType       // reference to the type information
	fields    []refFieldInfo // cached field information
}

// refFieldInfo contains information about a struct field for JSON operations
type refFieldInfo struct {
	name      string   // original field name (e.g., "BirthDate")
	snakeName string   // snake_case name of the field (e.g., "birth_date")
	refType   *refType // type of the field
	offset    uintptr  // byte offset in the struct
	index     int      // field index
}

// getStructInfo fills struct information if not cached, assigns to provided pointer
func getStructInfo(typ *refType, out *refStructInfo) {
	if typ.Kind() != tpStruct {
		return
	}

	// Get unique type name for caching
	typeName := getTypeName(typ)

	// Search in cache first
	for i := range refStructsInfo {
		if refStructsInfo[i].snakeName == typeName {
			*out = refStructsInfo[i]
			return
		}
	}

	// Not in cache, create new struct info
	structType := (*refStructType)(unsafe.Pointer(typ))
	fields := make([]refFieldInfo, len(structType.fields))

	for i, f := range structType.fields {
		fieldName := f.name.Name()
		snakeName := Convert(fieldName).ToSnakeCaseLower().String()

		fields[i] = refFieldInfo{
			name:      fieldName,
			snakeName: snakeName,
			refType:   f.typ,
			offset:    f.offset,
			index:     i,
		}
	}

	// Create new struct info
	newInfo := refStructInfo{
		snakeName: typeName,
		refType:   typ,
		fields:    fields,
	}

	// Add to cache
	refStructsInfo = append(refStructsInfo, newInfo)

	// Assign to output
	*out = newInfo
}

// clearRefStructsCache clears the global struct cache - useful for testing
func clearRefStructsCache() {
	refStructsInfo = refStructsInfo[:0] // Clear slice while preserving capacity
}
func getTypeName(typ *refType) string {
	if typ == nil {
		return "nil"
	}

	// Use type pointer and size to create unique identifier
	// Convert uintptr to string manually since Convert() doesn't handle uintptr
	ptr := uintptr(unsafe.Pointer(typ))
	ptrStr := ""
	if ptr != 0 {
		// Convert uintptr to base-10 string manually
		temp := ptr
		if temp == 0 {
			ptrStr = "0"
		} else {
			digits := ""
			for temp > 0 {
				digit := temp % 10
				digits = string(rune('0'+digit)) + digits
				temp /= 10
			}
			ptrStr = digits
		}
	}

	sizeStr := Convert(int64(typ.size)).String()
	kindStr := typ.Kind().String()

	return kindStr + "_" + sizeStr + "_" + ptrStr
}

// refNew returns a new pointer to a zero value of the given type
// Similar to reflect.New()
func refNew(typ *refType) refValue {
	if typ == nil {
		return refValue{}
	}

	// Allocate memory for the type and zero it
	ptr := refAllocate(typ.size)

	// Create pointer type that points to the original type
	ptrType := &refPtrType{
		refType: refType{
			size:       unsafe.Sizeof(ptr),
			kind:       uint8(tpPointer),
			align:      uint8(unsafe.Alignof(ptr)),
			fieldAlign: uint8(unsafe.Alignof(ptr)),
		},
		elem: typ,
	}

	return refValue{
		typ:  &ptrType.refType,
		ptr:  unsafe.Pointer(&ptr), // ptr to the allocated memory
		flag: flagIndir | flagAddr | refFlag(tpPointer),
	}
}

// refAppend appends values to a slice
// Similar to reflect.Append()
func refAppend(slice refValue, values ...refValue) refValue {
	if slice.Kind() != tpSlice {
		panic("refAppend: not a slice")
	}

	if len(values) == 0 {
		return slice
	}

	// For now, implement simple case with one value
	if len(values) != 1 {
		panic("refAppend: multiple values not implemented")
	}

	value := values[0]

	// Get slice header
	sliceHeader := (*refSliceHeader)(slice.ptr)
	elemType := slice.Type().Elem()
	elemSize := elemType.Size()

	// Create new slice with increased capacity
	newLen := sliceHeader.Len + 1
	newCap := sliceHeader.Cap
	if newCap < newLen {
		newCap = newLen * 2
	}

	// Allocate new data array
	newData := refAllocateArray(elemType, newCap)

	// Copy existing elements
	if sliceHeader.Len > 0 {
		refCopySlice(newData, sliceHeader.Data, sliceHeader.Len, elemSize)
	}

	// Copy new element to the end
	destPtr := unsafe.Add(newData, uintptr(sliceHeader.Len)*elemSize)
	refCopyValue(destPtr, value.ptr, elemType)

	// Create new slice value
	newSliceHeader := &refSliceHeader{
		Data: newData,
		Len:  newLen,
		Cap:  newCap,
	}

	return refValue{
		typ:  slice.typ,
		ptr:  unsafe.Pointer(newSliceHeader),
		flag: slice.flag,
	}
}

// Helper functions for memory operations

// refAllocate allocates memory for a type
func refAllocate(size uintptr) unsafe.Pointer {
	if size == 0 {
		return unsafe.Pointer(&zerobase)
	}
	// Use make([]byte, size) and return pointer to first element
	data := make([]byte, size)
	return unsafe.Pointer(&data[0])
}

// refAllocateArray allocates memory for an array of elements
func refAllocateArray(elemType *refType, count int) unsafe.Pointer {
	if count == 0 {
		return unsafe.Pointer(&zerobase)
	}
	size := elemType.Size() * uintptr(count)
	return refAllocate(size)
}

// refCopySlice copies slice elements
func refCopySlice(dest, src unsafe.Pointer, count int, elemSize uintptr) {
	if count <= 0 {
		return
	}
	// Simple byte copy
	srcBytes := (*[1 << 30]byte)(src)
	destBytes := (*[1 << 30]byte)(dest)
	totalSize := uintptr(count) * elemSize
	copy(destBytes[:totalSize], srcBytes[:totalSize])
}

// refCopyValue copies a single value
func refCopyValue(dest, src unsafe.Pointer, typ *refType) {
	if typ.Size() == 0 {
		return
	}
	srcBytes := (*[1 << 20]byte)(src)
	destBytes := (*[1 << 20]byte)(dest)
	copy(destBytes[:typ.Size()], srcBytes[:typ.Size()])
}

// refSliceHeader represents the runtime representation of a slice
type refSliceHeader struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}

// zerobase is used for zero-sized allocations
var zerobase uintptr

// memclr clears memory at ptr with size bytes
func memclr(ptr unsafe.Pointer, size uintptr) {
	// Simple implementation - zero out the memory
	slice := (*[1 << 30]byte)(ptr)[:size:size]
	for i := range slice {
		slice[i] = 0
	}
}
