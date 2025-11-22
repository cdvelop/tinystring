# Key-Value Parsing

TinyString provides key-value parsing functionality to extract values from strings with separators.

## Usage

```go
// Key-value parsing with the new API:
value, err := Convert("user:admin").ExtractValue()            // out: "admin", nil
value, err := Convert("count=42").ExtractValue("=")          // out: "42", nil
```

The `ExtractValue()` method splits the string on the first occurrence of the separator (default ":") and returns the value part. If a custom separator is provided, it uses that instead.