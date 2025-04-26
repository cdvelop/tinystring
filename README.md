# TinyString

TinyString is a lightweight Go library that provides text manipulation with a fluid API, without external dependencies or standard library dependencies.

## Features

- 🚀 Fluid and chainable API
- 🔄 Common text transformations
- 🧵 Concurrency safe
- 📦 No external dependencies
- 🎯 Easily extensible
- 🔄 Support for converting any data type to string

## Installation

```bash
go get github.com/cdvelop/tinystring
```

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
```

### Available Operations

- `Convert(v any)`: Initialize text processing with any data type (string, int, float, bool, etc.)
- `RemoveTilde()`: Removes accents and diacritics (e.g. "café" -> "cafe") 
- `ToLower()`: Converts to lowercase (e.g. "HELLO" -> "hello")
- `ToUpper()`: Converts to uppercase (e.g. "hello" -> "HELLO")
- `Capitalize()`: Capitalizes the first letter of each word (e.g. "hello world" -> "Hello World")
- `CamelCaseLower()`: Converts to camelCase (e.g. "hello world" -> "helloWorld")
- `CamelCaseUpper()`: Convert to UpperCase (e.g. "hello world" -> "HelloWorld")
- `ToSnakeCaseLower()`: Converts to snake_case (e.g. "hello world" -> "hello_world"), With  Other Params: `ToSnakeCaseLower("-")` -> "hello-world" 
- `ToSnakeCaseUpper()`: Convert to SNAKE_CASE (e.g. "hello world" -> "HELLO_WORLD"), With Other Params: `ToSnakeCaseUpper("-")` -> "HELLO-WORLD"
- `Split(data, separator string)`: Divides a string by a separator and returns a slice of substrings
- `ParseKeyValue(input string, delimiter string)`: Extracts the value from a key:value string format (e.g. `ParseKeyValue("name:John")` -> `"John", nil`)
- `Replace(old, new string)`: Replaces all occurrences of a substring (e.g. "hello world" -> "hello universe")
- `TrimPrefix(prefix string)`: Removes a specified prefix from the beginning of a string (e.g. "prefix-content" -> "content")
- `TrimSuffix(suffix string)`: Removes a specified suffix from the end of a string (e.g. "file.txt" -> "file")
- `Trim()`: Removes spaces from the beginning and end of a string (e.g. "  hello  " -> "hello")
- `Contains(text, search string)`: Checks if a string contains another, returns boolean (e.g. `Contains("hello world", "world")` -> `true`)
- `CountOccurrences(text, search string)`: Counts how many times a string appears in another (e.g. `CountOccurrences("hello hello world", "hello")` -> `2`)
- `Repeat(n int)`: Repeats the string n times (e.g. "abc".Repeat(3) -> "abcabcabc")
- `Truncate(maxWidth, reservedChars int)`: Truncates text to a specific width, adding ellipsis if necessary and padding with spaces if text is shorter (e.g. "Hello, World!".Truncate(10, 0) -> "Hello, ...")

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

// Split a string by separator
result := tinystring.Split("apple,banana,cherry", ",")
// Result: []string{"apple", "banana", "cherry"}

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
tinystring.Convert("Hello, World!").Truncate(10, 0).String()
// Result: "Hello, ..."

// Truncate with reserved characters
tinystring.Convert("Hello, World!").Truncate(10, 3).String()
// Result: "Hell..."

// Pad a short string with spaces
tinystring.Convert("Hello").Truncate(10, 0).String()
// Result: "Hello     "

// Truncate and transform
tinystring.Convert("hello world")
    .ToUpper()
    .Truncate(8, 0)
    .String()
// Result: "HELLO..."

// Chaining truncate and repeat
tinystring.Convert("hello")
    .Truncate(6, 0)
    .Repeat(2)
    .String()
// Result: "hello hello "
```

## Contributing

Contributions are welcome. Please open an issue to discuss proposed changes.

## License

MIT License