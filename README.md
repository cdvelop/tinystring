# TinyString

TinyString is a lightweight Go library that provides text manipulation with a fluid API, specifically designed for small devices and web applications using TinyGo as the target compiler.

## Key Features

- 🚀 **Fluid and chainable API** - Easy to use and readable operations
- 🔄 **Common text transformations** - All essential string operations included
- 🧵 **Concurrency safe** - Thread-safe operations
- 📦 **Zero standard library dependencies** - No `fmt`, `strings`, or `strconv` imports
- 🎯 **TinyGo compatible** - Optimized for minimal binary size and embedded systems
- 🌐 **Web-ready** - Perfect for WebAssembly deployments
- 🔄 **Universal type conversion** - Support for converting any data type to string
- ⚡ **Performance optimized** - Manual implementations avoid standard library overhead
- �️ **Easily extensible** - Clean architecture for adding new operations

## Why TinyString?

Traditional Go string libraries rely heavily on the standard library (`fmt`, `strings`, `strconv`), which can significantly increase binary size when using TinyGo for embedded systems or WebAssembly targets. TinyString provides all essential string operations with **manual implementations** that:

- ✅ Reduce binary size by avoiding standard library imports
- ✅ Ensure TinyGo compatibility without compilation issues  
- ✅ Provide predictable performance without hidden allocations
- ✅ Enable deployment to resource-constrained environments
- ✅ Support WebAssembly with minimal overhead


## Installation

```bash
go get github.com/cdvelop/tinystring
```

## TinyGo Optimization

TinyString is specifically designed for environments where binary size matters. Unlike other string libraries that import standard packages like `fmt`, `strings`, and `strconv`, TinyString implements all functionality manually to:

### Binary Size Comparison
```bash
# Traditional approach with standard library
go build -o app-standard main.go     # ~2.1MB binary
tinygo build -o app-standard.wasm -target wasm main.go  # 580.6KB WebAssembly

# TinyString approach  
go build -o app-tiny main.go         # ~1.2MB binary  
tinygo build -o app-tiny.wasm -target wasm main.go      # 229.5KB WebAssembly
```

#### WebAssembly Optimization Comparison

| Optimization Level | Standard Library | TinyString | Size Reduction |
|-------------------|------------------|------------|----------------|
| Default TinyGo optimization (-opt=z) | 580.6KB | 229.5KB | 60.5% |
| Maximum size optimization | 141.1KB | 22.2KB | 84.3% |
| Optimized for speed over size | 815.9KB | 321.3KB | 60.6% |
| No optimization, best for debugging | 1.8MB | 664.3KB | 63.8% |


### Target Environments
- **Embedded Systems**: Microcontrollers with limited flash memory
- **WebAssembly**: Minimal WASM modules for web applications  
- **IoT Devices**: Resource-constrained edge computing
- **Mobile Apps**: Reduced app size for faster downloads
- **Serverless Functions**: Cold start optimization

## Technical Implementation

TinyString achieves its goals through **manual implementations** of commonly used standard library functions:

### Replaced Standard Library Functions
| Standard Library | TinyString Implementation | Purpose |
|-----------------|---------------------------|---------|
| `strconv.ParseFloat` | `parseFloatManual` | String to float conversion |
| `strconv.FormatFloat` | `floatToStringManual` | Float to string conversion |
| `strconv.ParseInt` | `parseIntManual` | String to integer conversion |
| `strconv.FormatInt` | `intToStringOptimized` | Integer to string conversion |
| `strings.IndexByte` | `indexByteManual` | Byte search in strings |
| `fmt.Sprintf` | Custom `sprintf` implementation | String formatting |

### Performance Benefits
- **No hidden allocations** from standard library calls
- **Predictable memory usage** with manual control
- **Reduced binary bloat** from unused standard library code
- **TinyGo compatibility** without compilation warnings
- **Custom optimizations** for common use cases

## Usage

```go
import "github.com/cdvelop/tinystring"

// Basic example with string
text := tinystring.Convert("MÍ téxtO").RemoveTilde().String()
// Result: "MI textO"

// Examples with other data types
numText := tinystring.Convert(42).String()
// Result: "42"

boolText := tinystring.Convert(true).ToUpper().String()
// Result: "TRUE"

floatText := tinystring.Convert(3.14).String()
// Result: "3.14"

// Chaining operations
text := tinystring.Convert("Él Múrcielago Rápido")
    .RemoveTilde()
    .CamelCaseLower()
    .String()
// Result: "elMurcielagoRapido"

// Working with string pointers (avoids extra allocations)
// This method reduces memory allocations by modifying the original string directly
originalText := "Él Múrcielago Rápido"
tinystring.Convert(&originalText).RemoveTilde().CamelCaseLower().Apply()
// originalText is now modified directly: "elMurcielagoRapido"
```

### Available Operations

- `Convert(v any)`: Initialize text processing with any data type (string, *string, int, float, bool, etc.). When using a string pointer (*string) along with the `Apply()` method, the original string will be modified directly, avoiding extra memory allocations.
- `Apply()`: Updates the original string pointer with the current content. This method should be used when you want to modify the original string directly without additional allocations.
- `String()`: Returns the content of the text as a string without modifying any original pointers.
- `RemoveTilde()`: Removes accents and diacritics (e.g. "café" -> "cafe") 
- `ToLower()`: Converts to lowercase (e.g. "HELLO" -> "hello")
- `ToUpper()`: Converts to uppercase (e.g. "hello" -> "HELLO")
- `Capitalize()`: Capitalizes the first letter of each word (e.g. "hello world" -> "Hello World")
- `CamelCaseLower()`: Converts to camelCase (e.g. "hello world" -> "helloWorld")
- `CamelCaseUpper()`: Convert to UpperCase (e.g. "hello world" -> "HelloWorld")
- `ToSnakeCaseLower()`: Converts to snake_case (e.g. "hello world" -> "hello_world"), With  Other Params: `ToSnakeCaseLower("-")` -> "hello-world" 
- `ToSnakeCaseUpper()`: Convert to SNAKE_CASE (e.g. "hello world" -> "HELLO_WORLD"), With Other Params: `ToSnakeCaseUpper("-")` -> "HELLO-WORLD"
- `Split(data, separator string)`: Divides a string by a separator and returns a slice of substrings
- `Join(sep ...string)`: Joins elements of a string slice with a specified separator (default: space). (e.g. `Convert([]string{"Hello", "World"}).Join()` -> `"Hello World"` or `Convert([]string{"Hello", "World"}).Join("-")` -> `"Hello-World"`)
- `ParseKeyValue(input string, delimiter string)`: Extracts the value from a key:value string format (e.g. `ParseKeyValue("name:John")` -> `"John", nil`)
- `Replace(old, new any, n ...int)`: Replaces occurrences of a substring. If n is provided, replaces up to n occurrences. If n < 0 or not provided, replaces all. The old and new parameters can be of any type (string, int, float, bool) and will be converted to string automatically. (e.g. "hello world" -> "hello universe" or "value 123 here" -> "value 456 here")
- `TrimPrefix(prefix string)`: Removes a specified prefix from the beginning of a string (e.g. "prefix-content" -> "content")
- `TrimSuffix(suffix string)`: Removes a specified suffix from the end of a string (e.g. "file.txt" -> "file")
- `Trim()`: Removes spaces from the beginning and end of a string (e.g. "  hello  " -> "hello")
- `Contains(text, search string)`: Checks if a string contains another, returns boolean (e.g. `Contains("hello world", "world")` -> `true`)
- `CountOccurrences(text, search string)`: Counts how many times a string appears in another (e.g. `CountOccurrences("hello hello world", "hello")` -> `2`)
- `Repeat(n int)`: Repeats the string n times (e.g. "abc".Repeat(3) -> "abcabcabc")
- `Truncate(maxWidth any, reservedChars ...any)`: Truncates text so that it does not exceed the specified width, adding ellipsis if necessary. If the text is shorter or equal, it remains unchanged. The maxWidth parameter accepts any numeric type. The reservedChars parameter is optional and also accepts any numeric type. (e.g. "Hello, World!".Truncate(10) -> "Hello, ..." or "Hello, World!".Truncate(10, 3) -> "Hell...")
- `TruncateName(maxCharsPerWord any, maxWidth any)`: Truncates names and surnames in a user-friendly way for displaying in limited spaces like chart labels. It adds abbreviation dots where appropriate and handles the first word specially when there are more than 2 words. Parameters: maxCharsPerWord (maximum characters per word), maxWidth (maximum total length). (e.g. Convert("Jeronimo Dominguez").TruncateName(3, 15) -> "Jer. Dominguez")
- `RoundDecimals(decimals int)`: Rounds a numeric value to the specified number of decimal places with ceiling rounding by default (e.g. `Convert(3.154).RoundDecimals(2).String()` -> `"3.16"`)
- `Down()`: Modifies rounding behavior to floor rounding (must be used after RoundDecimals, e.g. `Convert(3.154).RoundDecimals(2).Down().String()` -> `"3.15"`)
- `FormatNumber()`: Formats a number with thousand separators and removes trailing zeros after the decimal point (e.g. `Convert(2189009.00).FormatNumber().String()` -> `"2.189.009"`)
- `Format(format string, args ...any)`: Static function for sprintf-style string formatting with support for %s, %d, %f, %b, %x, %o, %v, %% specifiers (e.g. `Format("Hello %s, you have %d messages", "John", 5)` -> `"Hello John, you have 5 messages"`)
- `StringError()`: Returns both the string result and any error that occurred during processing (e.g. `result, err := Convert("123").ToInt(); text, err2 := FromInt(result).StringError()`)
- `Quote()`: Wraps the string content in quotes with proper escaping of special characters (e.g. `Convert("hello").Quote().String()` -> `"\"hello\""`)
- `ToBool()`: Converts text content to boolean, supporting string boolean values and numeric values where non-zero = true (e.g. `Convert("true").ToBool()` -> `true, nil` or `Convert(42).ToBool()` -> `true, nil`)
- `FromBool(value bool)`: Creates a new Text instance from a boolean value (e.g. `FromBool(true).String()` -> `"true"`)
- `ToInt(base ...int)`: Converts text content to integer with optional base, supports float truncation (e.g. `Convert("123").ToInt()` -> `123, nil`)
- `ToUint(base ...int)`: Converts text content to unsigned integer with optional base, supports float truncation (e.g. `Convert("456").ToUint()` -> `456, nil`)  
- `ToFloat()`: Converts text content to float64 (e.g. `Convert("3.14").ToFloat()` -> `3.14, nil`)
- `FromInt(value int, base ...int)`: Creates a new Text instance from an integer with optional base (e.g. `FromInt(42).String()` -> `"42"`)
- `FromUint(value uint, base ...int)`: Creates a new Text instance from an unsigned integer with optional base (e.g. `FromUint(123).String()` -> `"123"`)
- `FromFloat(value float64)`: Creates a new Text instance from a float64 value (e.g. `FromFloat(3.14).String()` -> `"3.14"`)

### Enhanced Type Conversion and Formatting

TinyString now supports comprehensive type conversion with error handling and advanced formatting features:

#### String Formatting with Variable Arguments

```go
// Format strings with sprintf-style formatting
result := tinystring.Format("Hello %s, you have %d messages", "John", 5)
// Result: "Hello John, you have 5 messages"

// Support for various format specifiers
result := tinystring.Format("Number: %d, Float: %.2f, Bool: %v", 42, 3.14159, true)
// Result: "Number: 42, Float: 3.14, Bool: true"

// Format with hex, binary, and octal
result := tinystring.Format("Hex: %x, Binary: %b, Octal: %o", 255, 10, 8)
// Result: "Hex: ff, Binary: 1010, Octal: 10"
```

#### Advanced Numeric Rounding

```go
// Default rounding behavior (ceiling/up rounding)
result := tinystring.Convert(3.154).RoundDecimals(2).String()
// Result: "3.16"

// Explicit down rounding using Down() method
result := tinystring.Convert(3.154).RoundDecimals(2).Down().String() 
// Result: "3.15"

// Chain with other operations
result := tinystring.Convert(123.987).RoundDecimals(1).Down().ToUpper().String()
// Result: "123.9"
```

#### Boolean Conversion

```go
// String to boolean conversion
result, err := tinystring.Convert("true").ToBool()
// Result: true, err: nil

// Numeric to boolean (non-zero = true)
result, err := tinystring.Convert(42).ToBool()
// Result: true, err: nil

result, err := tinystring.Convert(0).ToBool() 
// Result: false, err: nil

// Boolean to string
result := tinystring.FromBool(true).String()
// Result: "true"
```

#### Enhanced Numeric Conversions

```go
// Integer conversions with float handling
result, err := tinystring.Convert("123").ToInt()
// Result: 123, err: nil

result, err := tinystring.Convert("123.45").ToInt() // Truncates float
// Result: 123, err: nil

// Unsigned integer conversions  
result, err := tinystring.Convert("456").ToUint()
// Result: 456, err: nil

result, err := tinystring.Convert("789.99").ToUint() // Truncates float
// Result: 789, err: nil

// Float conversions
result, err := tinystring.Convert("3.14159").ToFloat()
// Result: 3.14159, err: nil

// Creating from numeric types
result := tinystring.FromInt(42).String()
// Result: "42"

result := tinystring.FromUint(123).String() 
// Result: "123"

result := tinystring.FromFloat(3.14).String()
// Result: "3.14"
```

#### String Quoting

```go
// Add quotes around strings with proper escaping
result := tinystring.Convert("hello").Quote().String()
// Result: "\"hello\""

// Handle special characters
result := tinystring.Convert("say \"hello\"").Quote().String()
// Result: "\"say \\\"hello\\\"\""

// Quote with newlines and tabs
result := tinystring.Convert("line1\nline2\ttab").Quote().String()
// Result: "\"line1\\nline2\\ttab\""
```

#### Error Handling with StringError()

```go
// Get both result and error information
result, err := tinystring.Convert("invalid").ToInt()
if err != nil {
    fmt.Printf("Conversion failed: %v", err)
}

// Use StringError() method for operations that might fail
text := tinystring.Convert("123.45").RoundDecimals(2)
result, err := text.StringError()
// Result: "123.45", err: nil (or error if conversion failed)
```

#### Chaining New Operations

```go
// Complex chaining with new functionality
result := tinystring.Format("User %s has %d points", "Alice", 95)
formatted := tinystring.Convert(result).Quote().ToUpper().String()
// Result: "\"USER ALICE HAS 95 POINTS\""

// Numeric processing chain
result := tinystring.FromFloat(123.987)
    .RoundDecimals(2)
    .Down()
    .FormatNumber()
    .String()
// Result: "123.98"

// Type conversion chain
result, err := tinystring.Convert("42")
    .ToInt()
if err == nil {
    formatted := tinystring.FromInt(result * 2).Quote().String()
    // Result: "\"84\""
}
```

### Examples

```go
// Remove accents
tinystring.Convert("áéíóú").RemoveTilde().String()
// Result: "aeiou"

// Convert to camelCase
tinystring.Convert("hello world").CamelCaseLower().String()
// Result: "helloWorld"

// Combining operations
tinystring.Convert("HÓLA MÚNDO")
    .RemoveTilde()
    .ToLower()
    .String()
// Result: "hola mundo"

// Converting different data types
tinystring.Convert(123).String()
// Result: "123"

tinystring.Convert(45.67).String()
// Result: "45.67"

tinystring.Convert(true).String()
// Result: "true"

// Convert and transform other data types
tinystring.Convert(456).CamelCaseUpper().String()
// Result: "456"

tinystring.Convert(false).ToUpper().String()
// Result: "FALSE"

// Format number with decimal places
tinystring.Convert(3.12221).RoundDecimals(2).String()
// Result: "3.12"

// Format number with thousand separators
tinystring.Convert(2189009.00).FormatNumber().String()
// Result: "2.189.009"
// Result: "FALSE"

// Split a string by separator
result := tinystring.Split("apple,banana,cherry", ",")
// Result: []string{"apple", "banana", "cherry"}

// Split a string by whitespace (default)
result := tinystring.Split("hello world  test")
// Result: []string{"hello", "world", "test"}

// Split with mixed whitespace characters
result := tinystring.Split("hello\tworld\nnew")
// Result: []string{"hello", "world", "new"}

// Parse key-value string
value, err := tinystring.ParseKeyValue("user:admin")
// Result: value = "admin", err = nil

// Parse with custom delimiter
value, err := tinystring.ParseKeyValue("count=42", "=")
// Result: value = "42", err = nil

// Multiple values with same delimiter
value, err := tinystring.ParseKeyValue("path:usr:local:bin")
// Result: value = "usr:local:bin", err = nil

// Handle error when delimiter is not found
value, err := tinystring.ParseKeyValue("invalidstring")
// Result: value = "", err = error("delimiter ':' not found in string invalidstring")

// Join string slices with default space separator
result := tinystring.Convert([]string{"Hello", "World"}).Join().String()
// Result: "Hello World" 

// Join with custom separator
result := tinystring.Convert([]string{"apple", "banana", "orange"}).Join("-").String()
// Result: "apple-banana-orange"

// Join and chain with other transformations
result := tinystring.Convert([]string{"hello", "world"}).Join().ToUpper().String()
// Result: "HELLO WORLD"

// Replace text
tinystring.Convert("hello world").Replace("world", "universe").String()
// Result: "hello universe"

// Trim prefix and suffix
tinystring.Convert("prefix-content.txt").TrimPrefix("prefix-").TrimSuffix(".txt").String()
// Result: "content"

// Trim spaces and remove file extension
tinystring.Convert("  file.txt  ").Trim().TrimSuffix(".txt").String()
// Result: "file"

// Chain multiple operations
text := tinystring.Convert(" User Name ")
    .Trim()
    .Replace(" ", "_")
    .ToLower()
    .String()
// Result: "user_name"

// Search examples
// Check if a string contains another
result := tinystring.Contains("hello world", "world")
// Result: true

// Count occurrences
count := tinystring.CountOccurrences("abracadabra", "abra")
// Result: 2

// Capitalize each word
tinystring.Convert("hello world").Capitalize().String()
// Result: "Hello World"

// Capitalize with accent removal
tinystring.Convert("hólá múndo")
    .RemoveTilde()
    .Capitalize()
    .String()
// Result: "Hola Mundo"

// Repeat a string multiple times
tinystring.Convert("hello ").Repeat(3).String()
// Result: "hello hello hello "

// Repeat with other transformations
tinystring.Convert("test")
    .ToUpper()
    .Repeat(2)
    .String()
// Result: "TESTTEST"

// Zero or negative repetitions returns an empty string
tinystring.Convert("test").Repeat(0).String()
// Result: ""

// Truncate a long string to specific width
tinystring.Convert("Hello, World!").Truncate(10).String()
// Result: "Hello, ..."

// Truncate with reserved characters (explicitly provided)
tinystring.Convert("Hello, World!").Truncate(10, 3).String()
// Result: "Hell..."

// Text shorter than max width remains unchanged
tinystring.Convert("Hello").Truncate(10).String()
// Result: "Hello"

// Truncate names and surnames for display in charts or limited spaces
tinystring.Convert("Jeronimo Dominguez").TruncateName(3, 15).String()
// Result: "Jer. Dominguez"

// Truncate multiple names and surnames with total length limit
tinystring.Convert("Ana Maria Rodriguez").TruncateName(2, 10).String()
// Result: "An. Mar..."

// Handle first word specially when more than 2 words
tinystring.Convert("Juan Carlos Rodriguez").TruncateName(3, 20).String()
// Result: "Jua. Car. Rodriguez"

// Truncate and transform
tinystring.Convert("hello world")
    .ToUpper()
    .Truncate(8)
    .String()
// Result: "HELLO..."

// Truncate with different numeric types
tinystring.Convert("Hello, World!").Truncate(uint8(10), float64(3)).String()
// Result: "Hell..."

// Chaining truncate and repeat
tinystring.Convert("hello")
    .Truncate(6) // Truncate(6) doesn't change "hello"
    .Repeat(2)
    .String()
// Result: "hellohello"
```


### Working with String Pointers

TinyString supports working directly with string pointers to avoid additional memory allocations. This can be especially useful in performance-critical applications or when processing large volumes of text.

```go
// Create a string variable
text := "Él Múrcielago Rápido"

// Modify it directly using string pointer and Apply()
// No need to reassign the result
Convert(&text).RemoveTilde().ToLower().Apply()

// The original variable is modified
fmt.Println(text) // Output: "el murcielago rapido"

// This approach can reduce memory pressure in high-performance scenarios
// by avoiding temporary string allocations
```

## String() vs Apply()

The library offers two ways to finish a chain of operations:

```go
// 1. Using String() - Returns the result without modifying the original
originalText := "Él Múrcielago Rápido"
result := Convert(&originalText).RemoveTilde().ToLower().String()
fmt.Println(result)        // Output: "el murcielago rapido"
fmt.Println(originalText)  // Output: "Él Múrcielago Rápido" (unchanged)

// 2. Using Apply() - Modifies the original string directly
originalText = "Él Múrcielago Rápido"
Convert(&originalText).RemoveTilde().ToLower().Apply()
fmt.Println(originalText)  // Output: "el murcielago rapido" (modified)
```

Performance benchmarks show that using string pointers can reduce memory allocations when processing large volumes of text, which can be beneficial in high-throughput applications or systems with memory constraints.

```go
// Sample benchmark results:
BenchmarkMassiveProcessingWithoutPointer-16    114458 ops  10689 ns/op  4576 B/op  214 allocs/op
BenchmarkMassiveProcessingWithPointer-16       105290 ops  11434 ns/op  4496 B/op  209 allocs/op
```


## Contributing

Contributions are welcome. Please open an issue to discuss proposed changes.

## License

MIT License