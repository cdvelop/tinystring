
# Memory Allocation Tips (from fasthttp)

## Zero-Allocation Conversions entre `[]byte` y `string`

En código crítico para el rendimiento, convertir entre `[]byte` y `string` usando las funciones estándar de Go puede ser ineficiente por las asignaciones de memoria. Para evitarlo, se pueden usar conversiones **unsafe** que no asignan memoria:

```go
// Convierte []byte a string sin asignación
func UnsafeString(b []byte) string {
    // #nosec G103
    return *(*string)(unsafe.Pointer(&b))
}

// Convierte string a []byte sin asignación
func UnsafeBytes(s string) []byte {
    // #nosec G103
    return unsafe.Slice(unsafe.StringData(s), len(s))
}
```

**Advertencia:** Estas conversiones rompen la seguridad de tipos de Go. No modifiques el `[]byte` retornado por `UnsafeBytes(s string)` si la string original sigue en uso, ya que las strings son inmutables en Go y pueden ser compartidas en tiempo de ejecución.

## Trucos con buffers `[]byte`

- Las funciones estándar de Go aceptan buffers nil:

```go
var src []byte
// Ambos buffers no inicializados

dst = append(dst, src...)  // válido si dst y/o src son nil
copy(dst, src)             // válido si dst y/o src son nil
(string(src) == "")        // true si src es nil
(len(src) == 0)            // true si src es nil
src = src[:0]              // funciona con src nil

for i, ch := range src {   // no entra si src es nil
    doSomething(i, ch)
}
```

- No es necesario hacer nil checks para buffers `[]byte`:

```go
srcLen := len(src) // en vez de if src != nil { srcLen = len(src) }
```

- Puedes hacer append de un string a un buffer `[]byte`:

```go
dst = append(dst, "foobar"...)
```

- Un buffer `[]byte` puede extenderse hasta su capacidad:

```go
buf := make([]byte, 100)
a := buf[:10]  // len(a) == 10, cap(a) == 100
b := a[:100]   // válido, cap(a) == 100
```

- Todas las funciones de fasthttp aceptan buffers nil:

```go
statusCode, body, err := fasthttp.Get(nil, "http://google.com/")
uintBuf := fasthttp.AppendUint(nil, 1234)
```

- Conversiones sin asignación entre string y `[]byte`:

```go
func b2s(b []byte) string {
    return *(*string)(unsafe.Pointer(&b))
}

func s2b(s string) (b []byte) {
    bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
    return b
}
```

**Advertencia:** Es un método unsafe, el string y el buffer comparten los mismos bytes. ¡No modifiques el buffer si el string sigue vivo!

## Buenas prácticas para evitar asignaciones

- Reutiliza objetos y buffers `[]byte` tanto como sea posible.
- Usa [sync.Pool](https://pkg.go.dev/sync#Pool) para reutilización eficiente.
- Evita conversiones innecesarias entre `[]byte` y `string`.
- Verifica tu código con el [race detector](https://go.dev/doc/articles/race_detector.html).
- Escribe tests y benchmarks para los caminos críticos.
