# JSON Comparison Benchmark

Este directorio contiene benchmarks para comparar el rendimiento de operaciones JSON entre la biblioteca estándar de Go (`encoding/json`) y TinyString.

## Estructura

- `data.go`: Contiene las estructuras de datos y funciones de generación de datos de prueba
- `main_test.go`: Benchmarks principales para operaciones de Marshal y Unmarshal
- `errors_test.go`: Benchmarks específicos para casos de error

## Ejecución

Para ejecutar los benchmarks:

```bash
go test -bench=. -benchmem
```

Para un benchmark específico:

```bash
go test -bench=BenchmarkJsonMarshalSingle -benchmem
```

## Casos de Prueba

1. **Marshal (Encoding)**
   - Objeto individual
   - Lotes de 100, 1000, y 10000 objetos
   - Casos de error

2. **Unmarshal (Decoding)**
   - Objeto individual
   - Lotes de 100, 1000, y 10000 objetos
   - Casos de error

3. **Casos de Error**
   - JSON mal formado
   - Tipos incorrectos
   - Valores nulos inesperados
   - JSON truncado
   - Estructuras incompletas
