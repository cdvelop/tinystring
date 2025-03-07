# TinyString

TinyString is a lightweight Go library that provides text manipulation with a fluid API, without external dependencies or standard library dependencies.

## Features

- 🚀 Fluid and chainable API
- 🔄 Common text transformations
- 🧵 Concurrency safe
- 📦 No external dependencies
- 🎯 Easily extensible

## Installation

```bash
go get github.com/cdvelop/tinystring
```

## Usage

```go
import "github.com/cdvelop/tinystring"

// Basic example
text := tinystring.Convert("MÍ téxtO").RemoveTilde().String()
// Result: "MI textO"

// Chaining operations
text := tinystring.Convert("Él Múrcielago Rápido")
    .RemoveTilde()
    .CamelCaseLower()
    .String()
// Result: "elMurcielagoRapido"
```

### Available Operations

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
```

## Contributing

Contributions are welcome. Please open an issue to discuss proposed changes.

## License

MIT License