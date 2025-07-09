
# TinyString JSON API Integration Instructions

## Goal
Implement minimal, dependency-free JSON encoding and decoding for structs in TinyString, using the new reflection system. The API must be compatible with TinyGo and avoid all dependencies (including encoding/json).

## Core Philosophy
This implementation must adhere to TinyString's core principles as outlined in WHY.md:
- **🏆 Smallest possible binary size** - Minimize WebAssembly footprint
- **📦 Zero dependencies** - No imports beyond standard interfaces
- **🔧 Maximum code reuse** - Leverage existing code in TinyString
- **✅ Full TinyGo compatibility** - No compilation issues with TinyGo

## API Design
- Encoding: `Convert(&struct{}).JsonEncode(w io.Writer) error`
- Decoding: `Convert(r io.Reader).JsonDecode(&struct{}) error` (only pointer for decoding)

## Writer/Reader Types
- Use standard Go interfaces: `io.Writer` and `io.Reader` for maximum compatibility.

## JSON Standard
- Support basic struct tags for field naming (e.g., `json:"field_name"`).
- Only exported fields are encoded/decoded.
- No support for omitempty or advanced tag options.
- Minimal implementation: just enough for basic JSON compatibility and TinyGo support.

## Error Handling
- All errors must use the multilingual error system (`Err(D.Type, ...)`) as in the rest of the library.

## Dependencies
- 100% custom implementation. **No use of encoding/json or any external package.**

## File Structure
- Implement encoding in `jsonencode.go` and decoding in `jsondecode.go`.
- Place tests in `jsonencode_test.go` and `jsondecode_test.go`.

## Implementation Strategy
1. **Code Reuse**: Adapt essential reflection code from `tinyreflect/abi.go` and `tinyreflect/reflect.go` directly into TinyString's `reflect.go`.
2. **Struct Tag Handling**: Implement the minimum needed from `refStructTag` type for JSON tag parsing.
3. **Field Access**: Implement minimal struct field access methods similar to `GetFieldName()`, `GetFieldTag()`, etc.

## Supported Types (as per README.md)
JSON functionality must support the following types:
- **Basic types**: `string`, `bool`
- **Numeric types**: All int/uint variants, float32, float64
- **All basic slices**: `[]string`, `[]bool`, `[]byte`, etc.
- **Structs**: Only with supported field types
- **Maps with string keys**: `map[string]string`, `map[string]int`, etc.
- **Pointers**: Only to supported types above

## Test Coverage Requirements
- Create comprehensive tests that cover 100% of the JSON API
- Test all supported types listed in README.md
- Include test cases for:
  - Simple structs with primitive fields
  - Nested structs
  - Slices of structs
  - Maps with string keys
  - Error handling cases (invalid JSON, type mismatches)
  - TinyGo compatibility

## README Update
- Add a usage example for the new JSON API to the README **after implementation is complete**.

## Implementation Requirements
- Support reading struct field names and JSON tags
- Support encoding/decoding primitive types: string, bool, numbers
- Support basic slices and maps with string keys
- Proper error handling with the dictionary system
- TinyGo compatibility
