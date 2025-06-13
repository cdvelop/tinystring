# TinyString - Library Context for LLM Maintenance

## What is TinyString?

**Zero-dependency Go string manipulation library** for TinyGo/WebAssembly with fluid API. **80% smaller binaries** by replacing standard library with manual implementations.

## Core Problem & Solution

**Problem**: Standard library (`fmt`, `strings`, `strconv`, `encoding/json`) creates huge WebAssembly binaries  
**Solution**: Manual implementations in 5,062 lines of Go code, eliminating stdlib imports entirely

## Critical Constraints

- ❌ **NEVER import**: `fmt`, `strings`, `strconv`, `reflect`, `encoding/json`, `errors`
- ✅ **Only allowed**: `unsafe` and essential runtime packages
- 🔄 **Pattern**: All operations via `Convert(any).Method().String()` chain
- 🎯 **Priority**: Binary size over runtime performance

## API Pattern & Standard Library Replacements

```go
// Entry point: Convert any type, chain operations, get result
result := tinystring.Convert(input).Method1().Method2().String()

// Standard library replacements:
fmt.Sprintf(format, args...)     → Format(format, args...).String()
strings.ToLower(s)               → Convert(s).ToLower().String()  
strconv.Itoa(i)                  → Convert(i).String()
strconv.ParseInt(s, base, bits)  → Convert(s).ToInt(base)
strings.Split(s, sep)            → Split(s, sep)
strings.Join(slice, sep)         → Convert(slice).Join(sep).String()
```

## Complete Method Reference

| Category | Methods | Location | Stdlib Equivalent |
|----------|---------|----------|-------------------|
| **Core** | `Convert(any)`, `String()`, `Apply()` | convert.go | Entry point |
| **Text** | `ToLower()`, `ToUpper()`, `RemoveTilde()`, `Capitalize()` | capitalize.go | strings.ToLower/ToUpper |
| **Case** | `CamelCaseLower()`, `CamelCaseUpper()`, `ToSnakeCaseLower()` | capitalize.go | Custom transforms |
| **Strings** | `Split()`, `Join()`, `Replace()`, `TrimPrefix()`, `Repeat()` | split.go, join.go, replace.go, repeat.go | strings.* |
| **Search** | `Contains()`, `CountOccurrences()` | contain.go | strings.Contains/Count |
| **Numbers** | `ToInt()`, `ToUint()`, `ToFloat()`, `ToBool()`, `RoundDecimals()` | numeric.go, bool.go | strconv.Parse* |  
| **Format** | `Format()`, `FormatNumber()`, `Quote()` | format.go, numeric.go, quote.go | fmt.Sprintf |
| **Advanced** | `Truncate()`, `TruncateName()`, `ParseKeyValue()` | truncate.go, parse.go | Custom logic |
| **JSON** | `JsonEncode()`, `JsonDecode()` | json_encode.go, json_decode.go | encoding/json |
| **Error** | `StringError()`, `Errorf()`, `Err()` | error.go | fmt.Errorf |

## File Structure & Implementation Size

| File | Lines | Purpose | Key Methods |
|------|-------|---------|-------------|
| **convert.go** | 486 | Core engine, type system | `Convert()`, type interfaces, `Apply()` |
| **format.go** | 798 | Printf replacement | `Format()` (replaces fmt.Sprintf) |
| **numeric.go** | 799 | Number operations | `ToInt()`, `ToFloat()`, `RoundDecimals()`, `FormatNumber()` |
| **reflect.go** | 789 | Custom reflection | `refType`, `refValue` for JSON |
| **json_decode.go** | 531 | JSON parsing | `JsonDecode()` (replaces json.Unmarshal) |
| **json_encode.go** | 363 | JSON generation | `JsonEncode()` (replaces json.Marshal) |
| **abi.go** | 276 | Type definitions | `kind` enum, struct cache |
| **capitalize.go** | 218 | Text transforms | `Capitalize()`, `CamelCase*()`, `ToSnakeCase*()` |
| **truncate.go** | 174 | Smart truncation | `Truncate()`, `TruncateName()` |
| **replace.go** | 97 | String replacement | `Replace()` (replaces strings.Replace) |
| **mapping.go** | 94 | Character maps | Accent removal tables |
| **split.go** | 83 | String splitting | `Split()` (replaces strings.Split) |
| **error.go** | 78 | Error system | `Errorf()`, `Err()` (replaces fmt.Errorf) |
| **join.go** | 65 | String joining | `Join()` (replaces strings.Join) |
| **parse.go** | 58 | Parsing utilities | `ParseKeyValue()` |
| **bool.go** | 55 | Boolean conversion | `ToBool()` (replaces strconv.ParseBool) |
| **quote.go** | 42 | String quoting | `Quote()` (replaces strconv.Quote) |
| **contain.go** | 33 | Search operations | `Contains()`, `CountOccurrences()` |
| **repeat.go** | 22 | String repetition | `Repeat()` (replaces strings.Repeat) |
| **numeric_convert.go** | 1 | Placeholder | Empty file |

**Total: 5,062 lines** replacing 6 standard library packages

## Usage Examples & Binary Impact

```go
// Memory optimization: modify original string
text := "José María García"
tinystring.Convert(&text).RemoveTilde().CamelCaseLower().Apply()
// text = "joseMaria Garcia"

// Type conversion chain
result := tinystring.Convert(42.7859).RoundDecimals(2).FormatNumber().String()
// result = "42.79"

// JSON operations without encoding/json
bytes, err := tinystring.Convert(&user).JsonEncode()
err = tinystring.Convert(jsonData).JsonDecode(&user)
```

**Binary Size Reduction**: 76.5% - 87.6% smaller WebAssembly builds  
**Trade-off**: Higher runtime memory usage for smaller distribution size

## Test Organization & Debug Files

The TinyString project includes comprehensive test coverage with specialized debug and diagnostic tests:

### Debug Test Files
- **`json_debug_test.go`**: JSON encode/decode diagnostics, pointer-to-struct field tests, Convert() pointer handling
- **`reflect_debug_test.go`**: Reflection/field access diagnostics, field setter operations, corruption testing

### Core Test Files (by functionality)
- **JSON**: `json_encode_test.go`, `json_decode_test.go`, `json_data_test.go`
- **Reflection**: `reflect_test.go` 
- **Numeric**: `numeric_test.go`, `numeric_convert.go`
- **Text**: `capitalize_test.go`, `format_test.go`, `quote_test.go`
- **String Operations**: `contain_test.go`, `convert_test.go`, `join_test.go`, `parse_test.go`, `repeat_test.go`
- **Utilities**: `abi_test.go`, `bool_test.go`, `concurrency_test.go`

**Total Test Coverage**: ~95% with comprehensive edge case testing and diagnostic utilities
