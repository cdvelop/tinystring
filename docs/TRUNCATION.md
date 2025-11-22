# Smart Truncation

```go
// Basic truncation with ellipsis
Convert("Hello, World!").Truncate(10).String()
// out: "Hello, ..."

// Name truncation for UI display
Convert("Jeronimo Dominguez").TruncateName(3, 15).String()
// out: "Jer. Dominguez"

// Advanced name handling
Convert("Juan Carlos Rodriguez").TruncateName(3, 20).String()
// out: "Jua. Car. Rodriguez"
```