# TinyString

TinyString is a lightweight Go library that provides conv manipulation with a fluid API, specifically designed for small devices and web applications using TinyGo as the target compiler.

## Key Features

- 🚀 **Fluid and chainable API** - Easy to use and readable operations
- 🔄 **Common conv transformations** - All essential string operations included
- 🧵 **Concurrency safe** - Thread-safe operations
- 📦 **Zero standard library dependencies** - No `fmt`, `strings`, or `strconv` imports
- 🎯 **TinyGo compatible** - Optimized for minimal binary size and embedded systems
- 🌐 **Web-ready** - Perfect for WebAssembly deployments
- 🔄 **Universal type conversion** - Support for converting any data type to string
- ⚡ **Performance optimized** - Manual implementations avoid standard library overhead
- �️ **Easily extensible** - Clean architecture for adding new operations

## Why TinyString?

**Go is an incredibly pleasant language to write**, and discovering that it can run in the browser through WebAssembly was simply **amazing**! However, as a relatively new language in the WebAssembly ecosystem, not all the tooling exists yet for optimal web deployment.

While Go's standard library does an excellent job for traditional applications, **the main barrier to WebAssembly adoption with Go is the enormous binary sizes** when compiling for small devices or web deployment. This is particularly problematic because:

🌐 **Every project needs string manipulation** - Working with strings, numbers, errors, and basic data transformations is universal  
📦 **Standard library bloat** - Traditional libraries (`fmt`, `strings`, `strconv`) add significant overhead  
🚀 **WebAssembly adoption barrier** - Large binaries slow down web app loading and hurt user experience  
💾 **Memory constraints** - Small devices and edge computing require efficient resource usage  

**TinyString aims to solve this fundamental problem** by being:
- 🎯 **Small and lightweight** - Minimal binary footprint for web deployment
- ⚡ **Memory efficient** - Optimized allocation patterns
- 🔧 **Easy to use** - Familiar, chainable API that Go developers love
- 🌍 **WebAssembly-first** - Designed specifically for modern web deployment needs

TinyString provides all essential string operations with **manual implementations** that:

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

TinyString is specifically designed for environments where binary size matters. Unlike other string libraries that import standard packages like `fmt`, `strings`, and `strconv`, TinyString implements all functionality manually to achieve significant binary size reductions and optimal WebAssembly compatibility.


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
conv := tinystring.Convert("MÍ téxtO").RemoveTilde().String()
// Result: "MI textO"

// Examples with other data types
numText := tinystring.Convert(42).String()
// Result: "42"

boolText := tinystring.Convert(true).ToUpper().String()
// Result: "TRUE"

floatText := tinystring.Convert(3.14).String()
// Result: "3.14"

// Chaining operations
conv := tinystring.Convert("Él Múrcielago Rápido")
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

- `Convert(v any)`: Initialize conv processing with any data type (string, *string, int, float, bool, etc.). When using a string pointer (*string) along with the `Apply()` method, the original string will be modified directly, avoiding extra memory allocations.
- `Apply()`: Updates the original string pointer with the current content. This method should be used when you want to modify the original string directly without additional allocations.
- `String()`: Returns the content of the conv as a string without modifying any original pointers.
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
- `Contains(conv, search string)`: Checks if a string contains another, returns boolean (e.g. `Contains("hello world", "world")` -> `true`)
- `CountOccurrences(conv, search string)`: Counts how many times a string appears in another (e.g. `CountOccurrences("hello hello world", "hello")` -> `2`)
- `Repeat(n int)`: Repeats the string n times (e.g. "abc".Repeat(3) -> "abcabcabc")
- `Truncate(maxWidth any, reservedChars ...any)`: Truncates conv so that it does not exceed the specified width, adding ellipsis if necessary. If the conv is shorter or equal, it remains unchanged. The maxWidth parameter accepts any numeric type. The reservedChars parameter is optional and also accepts any numeric type. (e.g. "Hello, World!".Truncate(10) -> "Hello, ..." or "Hello, World!".Truncate(10, 3) -> "Hell...")
- `TruncateName(maxCharsPerWord any, maxWidth any)`: Truncates names and surnames in a user-friendly way for displaying in limited spaces like chart labels. It adds abbreviation dots where appropriate and handles the first word specially when there are more than 2 words. Parameters: maxCharsPerWord (maximum characters per word), maxWidth (maximum total length). (e.g. Convert("Jeronimo Dominguez").TruncateName(3, 15) -> "Jer. Dominguez")
- `RoundDecimals(decimals int)`: Rounds a numeric value to the specified number of decimal places with ceiling rounding by default (e.g. `Convert(3.154).RoundDecimals(2).String()` -> `"3.16"`)
- `Down()`: Modifies rounding behavior to floor rounding (must be used after RoundDecimals, e.g. `Convert(3.154).RoundDecimals(2).Down().String()` -> `"3.15"`)
- `FormatNumber()`: Formats a number with thousand separators and removes trailing zeros after the decimal point (e.g. `Convert(2189009.00).FormatNumber().String()` -> `"2.189.009"`)
- `Format(format string, args ...any)`: Static function for sprintf-style string formatting with support for %s, %d, %f, %b, %x, %o, %v, %% specifiers (e.g. `Format("Hello %s, you have %d messages", "John", 5)` -> `"Hello John, you have 5 messages"`)
- `StringError()`: Returns both the string result and any error that occurred during processing (e.g. `result, err := Convert("123").ToInt(); conv, err2 := Convert(result).StringError()`)
- `Quote()`: Wraps the string content in quotes with proper escaping of special characters (e.g. `Convert("hello").Quote().String()` -> `"\"hello\""`)
- `ToBool()`: Converts conv content to boolean, supporting string boolean values and numeric values where non-zero = true (e.g. `Convert("true").ToBool()` -> `true, nil` or `Convert(42).ToBool()` -> `true, nil`)
- `ToInt(base ...int)`: Converts conv content to integer with optional base, supports float truncation (e.g. `Convert("123").ToInt()` -> `123, nil`)
- `ToUint(base ...int)`: Converts conv content to unsigned integer with optional base, supports float truncation (e.g. `Convert("456").ToUint()` -> `456, nil`)  
- `ToFloat()`: Converts conv content to float64 (e.g. `Convert("3.14").ToFloat()` -> `3.14, nil`)


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
result := tinystring.Convert(true).String()
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
result := tinystring.Convert(42).String()
// Result: "42"

result := tinystring.Convert(123).String() 
// Result: "123"

result := tinystring.Convert(3.14).String()
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
conv := tinystring.Convert("123.45").RoundDecimals(2)
result, err := conv.StringError()
// Result: "123.45", err: nil (or error if conversion failed)
```

#### Chaining New Operations

```go
// Complex chaining with new functionality
result := tinystring.Format("User %s has %d points", "Alice", 95)
formatted := tinystring.Convert(result).Quote().ToUpper().String()
// Result: "\"USER ALICE HAS 95 POINTS\""

// Numeric processing chain
result := tinystring.Convert(123.987)
    .RoundDecimals(2)
    .Down()
    .FormatNumber()
    .String()
// Result: "123.98"

// Type conversion chain
result, err := tinystring.Convert("42")
    .ToInt()
if err == nil {
    formatted := tinystring.Convert(result * 2).Quote().String()
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

// Replace conv
tinystring.Convert("hello world").Replace("world", "universe").String()
// Result: "hello universe"

// Trim prefix and suffix
tinystring.Convert("prefix-content.txt").TrimPrefix("prefix-").TrimSuffix(".txt").String()
// Result: "content"

// Trim spaces and remove file extension
tinystring.Convert("  file.txt  ").Trim().TrimSuffix(".txt").String()
// Result: "file"

// Chain multiple operations
conv := tinystring.Convert(" User Name ")
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

// conv shorter than max width remains unchanged
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

TinyString supports working directly with string pointers to avoid additional memory allocations. This can be especially useful in performance-critical applications or when processing large volumes of conv.

```go
// Create a string variable
conv := "Él Múrcielago Rápido"

// Modify it directly using string pointer and Apply()
// No need to reassign the result
Convert(&conv).RemoveTilde().ToLower().Apply()

// The original variable is modified
fmt.Println(conv) // Output: "el murcielago rapido"

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


## Binary Size Comparison

[Standard Library Example](benchmark/bench-binary-size/standard-lib/main.go) | [TinyString Example](benchmark/bench-binary-size/tinystring-lib/main.go)

<!-- This table is automatically generated from build-and-measure.sh -->
*Last updated: 2025-06-05 21:34:43*

| Build Type | Parameters | Standard Library<br/>`go build` | TinyString<br/>`tinygo build` | Size Reduction | Performance |
|------------|------------|------------------|------------|----------------|-------------|
| 🖥️ **Default Native** | `-ldflags="-s -w"` | 1.6 MB | 1.5 MB | **-67.0 KB** | ❌ **4.2%** |
| 🌐 **Default WASM** | `(default -opt=z)` | 879.1 KB | 671.9 KB | **-207.2 KB** | ✅ **23.6%** |
| 🌐 **Ultra WASM** | `-no-debug -panic=trap -scheduler=none -gc=leaking -target wasm` | 200.6 KB | 94.8 KB | **-105.8 KB** | ✅ **52.8%** |
| 🌐 **Speed WASM** | `-opt=2 -target wasm` | 1.3 MB | 972.2 KB | **-318.4 KB** | ✅ **24.7%** |
| 🌐 **Debug WASM** | `-opt=0 -target wasm` | 3.0 MB | 2.3 MB | **-759.8 KB** | ✅ **24.7%** |

### 🎯 Performance Summary

- 🏆 **Peak Reduction: 52.8%** (Best optimization)
- ✅ **Average WebAssembly Reduction: 31.4%**
- ✅ **Average Native Reduction: 4.2%**
- 📦 **Total Size Savings: 1.4 MB across all builds**

#### Performance Legend
- ❌ Poor (<5% reduction)
- ➖ Fair (5-15% reduction)
- ✅ Good (15-70% reduction)
- 🏆 Outstanding (>70% reduction)


## Memory Usage Comparison

[Standard Library Example](benchmark/bench-memory-alloc/standard) | [TinyString Example](benchmark/bench-memory-alloc/tinystring)

<!-- This table is automatically generated from memory-benchmark.sh -->
*Last updated: 2025-06-05 21:35:11*

Performance benchmarks comparing memory allocation patterns between standard Go library and TinyString:

| 🧪 **Benchmark Category** | 📚 **Library** | 💾 **Memory/Op** | 🔢 **Allocs/Op** | ⏱️ **Time/Op** | 📈 **Memory Trend** | 🎯 **Alloc Trend** | 🏆 **Performance** |
|----------------------------|----------------|-------------------|-------------------|-----------------|---------------------|---------------------|--------------------|
| 📝 **String Processing** | 📊 Standard | `1.2 KB` | `48` | `3.1μs` | - | - | - |
| | 🚀 TinyString | `2.3 KB` | `46` | `9.3μs` | ❌ **96.7% more** | ➖ **4.2% less** | ❌ **Poor** |
| 🔢 **Number Processing** | 📊 Standard | `1.2 KB` | `132` | `4.2μs` | - | - | - |
| | 🚀 TinyString | `2.5 KB` | `120` | `3.8μs` | ❌ **110.7% more** | ✅ **9.1% less** | ❌ **Poor** |
| 🔄 **Mixed Operations** | 📊 Standard | `546 B` | `44` | `2.2μs` | - | - | - |
| | 🚀 TinyString | `1.2 KB` | `46` | `3.5μs` | ❌ **119.8% more** | ➖ **4.5% more** | ❌ **Poor** |
| 📝 **String Processing (Pointer Optimization)** | 📊 Standard | `1.2 KB` | `48` | `3.1μs` | - | - | - |
| | 🚀 TinyString | `2.2 KB` | `38` | `9.0μs` | ❌ **86.0% more** | 🏆 **20.8% less** | ⚠️ **Caution** |

### 🎯 Performance Summary

- 💾 **Memory Efficiency**: ❌ **Poor** (Significant overhead) (103.3% average change)
- 🔢 **Allocation Efficiency**: ✅ **Good** (Allocation efficient) (-7.4% average change)
- 📊 **Benchmarks Analyzed**: 4 categories
- 🎯 **Optimization Focus**: Binary size reduction vs runtime efficiency

### ⚖️ Trade-offs Analysis

The benchmarks reveal important trade-offs between **binary size** and **runtime performance**:

#### 📦 **Binary Size Benefits** ✅
- 🏆 **16-84% smaller** compiled binaries
- 🌐 **Superior WebAssembly** compression ratios
- 🚀 **Faster deployment** and distribution
- 💾 **Lower storage** requirements

#### 🧠 **Runtime Memory Considerations** ⚠️
- 📈 **Higher allocation overhead** during execution
- 🗑️ **Increased GC pressure** due to allocation patterns
- ⚡ **Trade-off optimizes** for distribution size over runtime efficiency
- 🔄 **Different optimization strategy** than standard library

#### 🎯 **Optimization Recommendations**
| 🎯 **Use Case** | 💡 **Recommendation** | 🔧 **Best For** |
|-----------------|------------------------|------------------|
| 🌐 WebAssembly Apps | ✅ **TinyString** | Size-critical web deployment |
| 📱 Embedded Systems | ✅ **TinyString** | Resource-constrained devices |
| ☁️ Edge Computing | ✅ **TinyString** | Fast startup and deployment |
| 🏢 Memory-Intensive Server | ⚠️ **Standard Library** | High-throughput applications |
| 🔄 High-Frequency Processing | ⚠️ **Standard Library** | Performance-critical workloads |

#### 📊 **Performance Legend**
- 🏆 **Excellent** (Better performance)
- ✅ **Good** (Acceptable trade-off)
- ⚠️ **Caution** (Higher resource usage)
- ❌ **Poor** (Significant overhead)


## Contributing

This project is currently being **self-financed** and developed independently. The development, testing, maintenance, and improvements are funded entirely out of my personal resources and time.

If you find this project useful and would like to support its continued development, you can make a donation [here with PayPal](https://paypal.me/cdvelop?country.x=CL&locale.x=es_XC). Your support helps cover:

- 💻 Development time and effort
- 🧪 Testing and quality assurance
- 📚 Documentation improvements
- 🔧 Bug fixes and feature enhancements
- 🌐 Community support and maintenance

Any contribution, however small, is greatly appreciated and directly impacts the project's future development. 🙌

## License

MIT License
