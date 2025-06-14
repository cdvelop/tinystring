package tinystring

import (
	"unsafe"
)

// Minimal reflectlite integration for TinyString JSON functionality
// This file contains essential reflection capabilities adapted from internal/reflectlite
// All functions and types are prefixed with 'ref' to avoid API pollution

// NOTE: refValue struct has been completely eliminated and fused into conv
// This maintains separation of concerns while eliminating duplication

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

// refValueOf returns a new conv initialized to the concrete value stored in i
// This replaces the old refValue-based function
func refValueOf(i any) *conv {
	c := &conv{separator: "_"}
	if i == nil {
		return c
	}
	c.initFromValue(i)
	return c
}

// refEface is the header for an interface{} value
// NOTE: This definition is now in convert.go - commenting out to avoid duplication
/*
type refEface struct {
	typ  *refType
	data unsafe.Pointer
}
*/

// ifaceIndir reports whether t is stored indirectly in an interface value
func ifaceIndir(t *refType) bool {
	return t.kind&kindDirectIface == 0
}

// Type returns the type of v
func (c *conv) Type() *refType {
	return c.typ
}

// refElem returns the value that the interface c contains or that the pointer c points to
func (c *conv) refElem() *conv {
	k := c.refKind()
	switch k {
	case tpInterface:
		var eface refEface
		if c.typ.kind&kindDirectIface != 0 {
			eface = refEface{typ: nil, data: c.ptr}
		} else {
			eface = *(*refEface)(c.ptr)
		}
		if eface.typ == nil {
			return &conv{}
		}
		result := &conv{separator: "_"}
		result.typ = eface.typ
		result.ptr = eface.data
		result.flag = refFlag(eface.typ.Kind())
		if ifaceIndir(eface.typ) {
			result.flag |= flagIndir
		}
		return result
	case tpPointer:
		// Handle pointer dereferencing
		var ptr unsafe.Pointer
		if c.flag&flagIndir != 0 {
			// This is a pointer field from a struct - need to dereference to get the actual pointer
			ptr = *(*unsafe.Pointer)(c.ptr)
		} else {
			// This is a direct pointer from interface{}
			// c.ptr contains the pointer value itself (the address it points to)
			ptr = c.ptr
		}

		if ptr == nil {
			return &conv{}
		}

		elemType := c.typ.Elem()
		if elemType == nil {
			return &conv{}
		}

		// Create proper flags for the element
		// The element is addressable since we're dereferencing a pointer
		fl := c.flag&flagRO | flagAddr | refFlag(elemType.Kind())

		// For elements accessed through pointers, we don't need flagIndir
		// because ptr already points to the actual data
		result := &conv{separator: "_"}
		result.typ = elemType
		result.ptr = ptr
		result.flag = fl
		return result
	}
	panic("reflect: call of reflect.Value.Elem on " + c.Type().Kind().String() + " value")
}

// refNumField returns the number of fields in the struct c
func (c *conv) refNumField() int {
	c.mustBe(tpStruct)
	tt := (*refStructType)(unsafe.Pointer(c.typ))
	return len(tt.fields)
}

// refField returns the i'th field of the struct c
func (c *conv) refField(i int) *conv {
	if c.refKind() != tpStruct {
		panic("reflect: call of reflect.Value.Field on " + c.refKind().String() + " value")
	}
	tt := (*refStructType)(unsafe.Pointer(c.typ))
	if uint(i) >= uint(len(tt.fields)) {
		panic("reflect: Field index out of range")
	}
	field := &tt.fields[i]
	ptr := add(c.ptr, field.offset, "same as non-reflect &v.field")

	// Inherit read-only flags from parent, but allow assignment if parent allows it
	fl := c.flag&(flagRO) | refFlag(field.typ.Kind()) | flagAddr
	// For struct fields, flagIndir is needed only for pointer fields
	// because ptr points to the field location containing the pointer.
	// For other field types, ptr points directly to the field value.
	if field.typ.Kind() == tpPointer {
		fl |= flagIndir
	}

	result := &conv{separator: "_"}
	result.typ = field.typ
	result.ptr = ptr
	result.flag = fl
	return result
}

// refSetString sets c's underlying value to x
func (c *conv) refSetString(x string) {
	c.mustBeAssignable()
	c.mustBe(tpString)
	ptr := c.ptr
	if c.flag&flagIndir != 0 {
		ptr = *(*unsafe.Pointer)(ptr)
	}
	*(*string)(ptr) = x
}

// refSetInt sets c's underlying value to x
func (c *conv) refSetInt(x int64) {
	c.mustBeAssignable()
	ptr := c.ptr
	if c.flag&flagIndir != 0 {
		ptr = *(*unsafe.Pointer)(ptr)
	}
	switch c.refKind() {
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
		panic("reflect: call of reflect.Value.SetInt on " + c.refKind().String() + " value")
	}
}

// refSetUint sets c's underlying value to x
func (c *conv) refSetUint(x uint64) {
	c.mustBeAssignable()
	ptr := c.ptr
	if c.flag&flagIndir != 0 {
		ptr = *(*unsafe.Pointer)(ptr)
	}
	switch c.refKind() {
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
		c.err = errorType("reflect: call of reflect.Value.SetUint on " + c.refKind().String() + " value")
	}
}

// refSetFloat sets c's underlying value to x
func (c *conv) refSetFloat(x float64) {
	c.mustBeAssignable()
	ptr := c.ptr
	if c.flag&flagIndir != 0 {
		ptr = *(*unsafe.Pointer)(ptr)
	}
	switch c.refKind() {
	case tpFloat32:
		*(*float32)(ptr) = float32(x)
	case tpFloat64:
		*(*float64)(ptr) = x
	default:
		c.err = errorType("reflect: call of reflect.Value.SetFloat on " + c.refKind().String() + " value")
	}
}

// refSetBool sets c's underlying value to x
func (c *conv) refSetBool(x bool) {
	c.mustBeAssignable()
	c.mustBe(tpBool)
	ptr := c.ptr
	if c.flag&flagIndir != 0 {
		ptr = *(*unsafe.Pointer)(ptr)
	}
	*(*bool)(ptr) = x
}

// refSet assigns x to the value c
// c must be addressable and must not have been obtained by accessing unexported struct fields
func (c *conv) refSet(x *conv) {
	c.mustBeAssignable()
	if c.err != "" {
		return
	}
	x.mustBeExported() // do not let unexported x leak
	if x.err != "" {
		c.err = x.err
		return
	}

	// For pointer types, we need to copy the pointer value itself
	if c.refKind() == tpPointer {
		// c.ptr points to the pointer variable
		// We need to set the pointer variable to the value that x represents
		if x.refKind() == tpPointer {
			// Copy pointer value from x to c
			*(*unsafe.Pointer)(c.ptr) = *(*unsafe.Pointer)(x.ptr)
		} else {
			// x is not a pointer, this shouldn't happen in normal cases
			typedmemmove(c.typ, c.ptr, x.ptr)
		}
	} else {
		// For non-pointer types, copy the value
		typedmemmove(c.typ, c.ptr, x.ptr)
	}
}

// refZero returns a conv representing the zero value for the specified type
func refZero(typ *refType) *conv {
	if typ == nil {
		return &conv{err: errorType("reflect: Zero(nil)")}
	}

	c := &conv{separator: "_"}

	// For pointer types, zero value is nil pointer
	if typ.Kind() == tpPointer {
		var nilPtr unsafe.Pointer // This is nil
		c.typ = typ
		c.ptr = unsafe.Pointer(&nilPtr)
		c.flag = refFlag(tpPointer)
		return c
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
	c.typ = typ
	c.ptr = ptr
	c.flag = refFlag(typ.Kind()) | flagAddr

	return c
}

// mustBeExported sets error if c was obtained using an unexported field
func (c *conv) mustBeExported() {
	if c.err != "" {
		return
	}
	if c.flag&flagRO != 0 {
		c.err = errorType("reflect: use of unexported field")
	}
}

// mustBeAssignable sets error if c is not assignable
func (c *conv) mustBeAssignable() {
	if c.err != "" {
		return
	}
	if c.flag&flagRO != 0 {
		c.err = errorType("reflect: cannot set value")
		return
	}
	if c.flag&flagAddr == 0 {
		c.err = errorType("reflect: cannot assign to value")
		return
	}
}

// mustBe sets error if c's kind is not expected
func (c *conv) mustBe(expected kind) {
	if c.err != "" {
		return
	}
	if c.refKind() != expected {
		c.err = errorType("reflect: call of reflect.Value method on " + expected.String() + " value")
	}
}

// refKind returns the Kind without the flags
func (c *conv) refKind() kind {
	return kind(c.flag & flagKindMask)
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

// refIsValid reports whether c represents a value
func (c *conv) refIsValid() bool {
	return c.flag != 0
}

// refInt returns c's underlying value, as an int64
func (c *conv) refInt() int64 {
	if c.err != "" {
		return 0
	}

	ptr := c.ptr
	if c.flag&flagIndir != 0 {
		ptr = *(*unsafe.Pointer)(ptr)
	}

	switch k := c.refKind(); k {
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
		c.err = errorType("reflect: call of reflect.Value.Int on " + c.refKind().String() + " value")
		return 0
	}
}

// refUint returns c's underlying value, as a uint64
func (c *conv) refUint() uint64 {
	if c.err != "" {
		return 0
	}

	ptr := c.ptr
	if c.flag&flagIndir != 0 {
		ptr = *(*unsafe.Pointer)(ptr)
	}

	switch k := c.refKind(); k {
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
		c.err = errorType("reflect: call of reflect.Value.Uint on " + c.refKind().String() + " value")
		return 0
	}
}

// refFloat returns c's underlying value, as a float64
func (c *conv) refFloat() float64 {
	if c.err != "" {
		return 0
	}

	ptr := c.ptr
	if c.flag&flagIndir != 0 {
		ptr = *(*unsafe.Pointer)(ptr)
	}

	switch k := c.refKind(); k {
	case tpFloat32:
		return float64(*(*float32)(ptr))
	case tpFloat64:
		return *(*float64)(ptr)
	default:
		c.err = errorType("reflect: call of reflect.Value.Float on " + c.refKind().String() + " value")
		return 0
	}
}

// refBool returns c's underlying value
func (c *conv) refBool() bool {
	if c.err != "" {
		return false
	}

	c.mustBe(tpBool)
	if c.err != "" {
		return false
	}

	ptr := c.ptr
	if c.flag&flagIndir != 0 {
		ptr = *(*unsafe.Pointer)(ptr)
	}
	return *(*bool)(ptr)
}

// refString returns c's underlying value, as a string
func (c *conv) refString() string {
	if c.err != "" {
		return ""
	}

	if !c.refIsValid() {
		return ""
	}

	// Don't enforce mustBe() - allow reading strings from struct fields
	if c.refKind() != tpString {
		return ""
	}

	ptr := c.ptr
	if c.flag&flagIndir != 0 {
		ptr = *(*unsafe.Pointer)(ptr)
	}
	return *(*string)(ptr)
}

// Interface returns c's current value as an interface{}
func (c *conv) Interface() any {
	if c.err != "" {
		return nil
	}

	if !c.refIsValid() {
		return nil
	}

	switch c.refKind() {
	case tpString:
		return c.refString()
	case tpInt:
		return int(c.refInt())
	case tpInt8:
		return int8(c.refInt())
	case tpInt16:
		return int16(c.refInt())
	case tpInt32:
		return int32(c.refInt())
	case tpInt64:
		return c.refInt()
	case tpUint:
		return uint(c.refUint())
	case tpUint8:
		return uint8(c.refUint())
	case tpUint16:
		return uint16(c.refUint())
	case tpUint32:
		return uint32(c.refUint())
	case tpUint64:
		return c.refUint()
	case tpUintptr:
		return uintptr(c.refUint())
	case tpFloat32:
		return float32(c.refFloat())
	case tpFloat64:
		return c.refFloat()
	case tpBool:
		return c.refBool()
	case tpInterface:
		// For interface{} types, extract the contained value directly
		var eface refEface
		if c.typ.kind&kindDirectIface != 0 {
			eface = refEface{typ: nil, data: c.ptr}
		} else {
			eface = *(*refEface)(c.ptr)
		}
		if eface.typ == nil {
			return nil
		}

		// Create a new interface{} with the contained value
		return *(*any)(unsafe.Pointer(&eface))
	case tpStruct: // For struct types, create an interface{} with the struct value
		// The struct data is stored at c.ptr
		var eface refEface
		eface.typ = c.typ
		eface.data = c.ptr
		return *(*any)(unsafe.Pointer(&eface))
	default:
		// For complex types, return nil for now
		return nil
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

// memclr clears memory at ptr with size bytes
func memclr(ptr unsafe.Pointer, size uintptr) {
	// Simple implementation - zero out the memory
	slice := (*[1 << 30]byte)(ptr)[:size:size]
	for i := range slice {
		slice[i] = 0
	}
}

// refLen returns the length of c
// It panics if c's Kind is not Slice
func (c *conv) refLen() int {
	if c.err != "" {
		return 0
	}
	k := c.refKind()
	switch k {
	case tpSlice:
		// For slices, the length is stored in the slice header
		return (*sliceHeader)(c.ptr).Len
	default:
		c.err = errorType("reflect: call of reflect.Value.Len on " + k.String() + " value")
		return 0
	}
}

// refIndex returns c's i'th element
// It panics if c's Kind is not Slice or if i is out of range
func (c *conv) refIndex(i int) *conv {
	if c.err != "" {
		return &conv{err: c.err}
	}
	k := c.refKind()
	switch k {
	case tpSlice:
		s := (*sliceHeader)(c.ptr)
		if i < 0 || i >= s.Len {
			c.err = errorType("reflect: slice index out of range")
			return &conv{err: c.err}
		}

		// Get element type
		elemType := c.typ.Elem()
		if elemType == nil {
			return &conv{err: errorType("reflect: slice element type is nil")}
		}

		elemSize := elemType.Size()

		// Calculate pointer to element
		elemPtr := unsafe.Pointer(uintptr(s.Data) + uintptr(i)*elemSize)
		// Create new conv for the element
		result := &conv{separator: "_"}
		result.typ = elemType
		result.ptr = elemPtr
		result.flag = refFlag(elemType.Kind())

		// If element is stored indirectly, set the flag
		// Note: strings should never be indirect in slices
		if elemType.Kind() != tpString && elemType.kind&kindDirectIface == 0 {
			result.flag |= flagIndir
		}

		return result
	default:
		c.err = errorType("reflect: call of reflect.Value.Index on " + k.String() + " value")
		return &conv{err: c.err}
	}
}

// sliceHeader is the runtime representation of a slice
type sliceHeader struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}
