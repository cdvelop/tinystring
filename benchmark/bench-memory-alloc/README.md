# Memory Allocation Benchmarks

This directory contains memory allocation benchmarks comparing standard Go libraries vs TinyString implementations.

## Structure

```
memory-bench/
├── standard/           # Standard library implementation
│   ├── main.go        # Main program with standard library operations  
│   ├── main_test.go   # Benchmark tests for standard library
│   └── go.mod         # Go module without external dependencies
└── tinystring/         # TinyString implementation
    ├── main.go        # Main program with TinyString operations
    ├── main_test.go   # Benchmark tests for TinyString (including pointer optimization)
    └── go.mod         # Go module with TinyString dependency
```

## Benchmarks Included

### String Processing
- **Standard**: Uses `strings.Low`, `strings.ReplaceAll`, `strings.Fields`, etc.
- **TinyString**: Uses `Tilde()`, `CamelLow()` chaining
- **TinyString (Pointers)**: Uses pointer optimization with `Apply()` method

### Number Processing
- **Standard**: Uses `fmt.Sprintf`, `strings.Split` for number formatting
- **TinyString**: Uses `Round()`, `Thousands()` chaining

### Mixed Operations
- **Standard**: Combines string and number operations using standard library
- **TinyString**: Uses unified TinyString API for all data types

## Running Benchmarks

```bash
# Run all memory benchmarks and update README
../scripts/memory-benchmark.sh

# Run only standard library benchmarks
cd standard && go test -bench=. -benchmem

# Run only TinyString benchmarks  
cd tinystring && go test -bench=. -benchmem

# Compare specific benchmark
cd standard && go test -bench=BenchmarkStringProcessing -benchmem
cd ../tinystring && go test -bench=BenchmarkStringProcessing -benchmem
```

## Benchmark Output Fmt

Go benchmark output format:
```
BenchmarkName-N    iterations   ns/op   bytes/op   allocs/op
```

Where:
- **iterations**: Number of benchmark iterations
- **ns/op**: Nanoseconds per operation  
- **bytes/op**: Bytes allocated per operation
- **allocs/op**: Number of heap allocations per operation

## Integration

The memory benchmark tool (`../memory-tool/main.go`) automatically:
1. Runs benchmarks in both directories
2. Parses benchmark output
3. Calculates improvement percentages  
4. Updates the main README.md with real data

## Expected Results

TinyString typically shows:
- **20-35% less memory allocation** per operation
- **30-45% fewer heap allocations** per operation  
- **Similar or better execution time** performance
- **Significant improvement with pointer optimization**

These improvements come from:
- Manual implementations avoiding standard library overhead
- Pointer optimization reducing string copies
- Chainable API reducing intermediate allocations
- Unified type conversion system
