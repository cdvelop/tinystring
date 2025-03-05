# TinyText

TinyText is a lightweight Go library that provides text manipulation with a fluid API, without external dependencies or standard library dependencies.

## Features

- 🚀 Fluid and chainable API
- 🔄 Common text transformations
- 🧵 Concurrency safe
- 📦 No external dependencies
- 🎯 Easily extensible

## Installation

```bash
go get github.com/cdvelop/tinytext
```

## Usage

```go
import "github.com/cdvelop/tinytext"

// Basic example
text := tinytext.Convert("MÍ téxtO").RemoveTilde().String()
// Result: "MI textO"

// Chaining operations
text := tinytext.Convert("Él Múrcielago Rápido")
    .RemoveTilde()
    .CamelCaseLower()
    .String()
// Result: "elMurcielagoRapido"
```

### Available Operations

- `RemoveTilde()`: Removes accents and diacritics (e.g. "café" -> "cafe")
- `ToLower()`: Converts to lowercase (e.g. "HELLO" -> "hello")
- `ToUpper()`: Converts to uppercase (e.g. "hello" -> "HELLO")
- `CamelCaseLower()`: Converts to camelCase (e.g. "hello world" -> "helloWorld")
- `CamelCaseUpper()`: Convert to UpperCase (e.g. "hello world" -> "HelloWorld")
- `ToSnakeCaseLower()`: Converts to snake_case (e.g. "hello world" -> "hello_world"), With  Other Params: `ToSnakeCaseLower("-")` -> "hello-world" 
- `ToSnakeCaseUpper()`: Convert to SNAKE_CASE (e.g. "hello world" -> "HELLO_WORLD"), With Other Params: `ToSnakeCaseUpper("-")` -> "HELLO-WORLD"


### Examples

```go
// Remove accents
tinytext.Convert("áéíóú").RemoveTilde().String()
// Result: "aeiou"

// Convert to camelCase
tinytext.Convert("hello world").CamelCaseLower().String()
// Result: "helloWorld"

// Combining operations
tinytext.Convert("HÓLA MÚNDO")
    .RemoveTilde()
    .ToLower()
    .String()
// Result: "hola mundo"
```

## Contributing

Contributions are welcome. Please open an issue to discuss proposed changes.

## License

MIT License