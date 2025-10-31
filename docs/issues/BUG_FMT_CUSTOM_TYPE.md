# BUG: Fmt() Doesn't Handle Custom Types with String() Method

**Date**: 2025-10-31  
**Status**: Confirmed  
**Severity**: High  
**Impact**: Breaks fpdf PDF generation (missing `%PDF` header)

## Problem

`tinystring.Fmt()` returns empty string when formatting custom types with `%s`, even if they implement `String()` method. This breaks compatibility with standard `fmt.Sprintf()` behavior.

## Root Cause

In `fmt_template.go:426-433`, the `%s` case only accepts native `string` type:

```go
case 's':
    if strVal, ok := arg.(string); ok {
        return strVal
    } else {
        c.wrInvalidTypeErr(formatSpec)
        return ""  // ❌ Returns empty string for custom types
    }
```

## Test Case (Confirmed)

```go
type customType string
func (c customType) String() string { return string(c) }

Fmt("Version: %s", customType("1.3"))  // Returns: "" (should be "Version: 1.3")
```

**Real-world impact**: `fpdf` calls `outf("%%PDF-%s", pdfVersion)` where `pdfVersion` is custom type → PDF has no header → "Failed to load PDF document"

## Proposed Solution

Modify `formatValue()` case `'s'` to check for types with `String()` method:

```go
case 's':
    // Handle native string
    if strVal, ok := arg.(string); ok {
        return strVal
    }
    // Handle types with String() method (fmt.Stringer interface)
    if stringer, ok := arg.(interface{ String() string }); ok {
        return stringer.String()
    }
    // Fallback: try AnyToBuff (already handles String() internally)
    c.ResetBuffer(BuffWork)
    c.AnyToBuff(BuffWork, arg)
    if c.hasContent(BuffErr) {
        c.wrInvalidTypeErr(formatSpec)
        return ""
    }
    return c.GetString(BuffWork)
```

## Alternative (Simpler)

Reuse existing `AnyToBuff()` which already handles `String()` method:

```go
case 's':
    if strVal, ok := arg.(string); ok {
        return strVal
    }
    // Let AnyToBuff handle it (supports String() method)
    c.ResetBuffer(BuffWork)
    c.AnyToBuff(BuffWork, arg)
    if c.hasContent(BuffErr) {
        return ""
    }
    return c.GetString(BuffWork)
```

## Testing

Run: `cd tinystring && go test -v -run TestFmtWithCustomTypeString`

Expected after fix:
- ✅ `customType("1.3")` with `%s` → "1.3"
- ✅ `Fmt("%%PDF-%s", pdfVersion)` → "%PDF-1.4"
- ✅ PDF generation works in browser

## Impact Analysis

**Before fix:**
- `Fmt("%s", customType)` → `""` (empty string)
- PDF generation broken (no header)
- Incompatible with `fmt.Sprintf()` behavior

**After fix:**
- `Fmt("%s", customType)` → calls `String()` method
- PDF generation works
- Compatible with `fmt.Sprintf()` behavior

## Files to Modify

1. `tinystring/fmt_template.go` - line 426-433 (formatValue case 's')
2. Verify all tests pass: `go test ./...`

## Notes

- `%v` format already works via `AnyToBuff()` 
- `%d`, `%f` etc. have similar issues with custom numeric types
- This fix should apply to all format specifiers that reject custom types
