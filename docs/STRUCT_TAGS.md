# Struct Tag Extraction

fmt allows extracting values from struct field tags, useful for parsing metadata like JSON tags or custom labels.

## Usage

```go
// Struct tag value extraction (TagValue):
value, found := Convert(`json:"name" Label:"Nombre"`).TagValue("Label") // out: "Nombre", true
value, found := Convert(`json:"name" Label:"Nombre"`).TagValue("xml")   // out: "", false
```

The `TagValue()` method parses the tag string and extracts the value for a specific key. It returns the value and a boolean indicating if the key was found.