# TinyString Binary Size Benchmark

This directory contains automated tools to measure and compare binary sizes between standard Go libraries and TinyString implementations.

## Overview

The benchmark system creates two equivalent programs:
- **Standard Library**: Uses `fmt`, `strings`, `strconv` packages
- **TinyString**: Uses only the TinyString library

Both programs are compiled to:
- Native binaries (using `go build`)  
- WebAssembly modules (using `tinygo build`)

## Directory Structure

```
benchmark/
├── examples/
│   ├── standard-lib/          # Example using standard library
│   │   ├── main.go           
│   │   └── go.mod
│   └── tinystring-lib/        # Example using TinyString
│       ├── main.go
│       └── go.mod
├── scripts/
│   ├── build-and-measure.sh   # Main benchmark script
│   ├── clean.sh              # Clean generated files
│   └── update-readme.sh      # Update README only
├── benchmark.go              # Size analysis program
└── README.md                # This file
```

## Requirements

- **Go 1.21+**: For building native binaries
- **TinyGo** (optional): For WebAssembly compilation
  - Install from: https://tinygo.org/getting-started/install/
  - If not installed, only native binaries will be measured

## Usage

### Run Complete Benchmark

```bash
# Make scripts executable (Linux/macOS/WSL)
chmod +x scripts/*.sh

# Run full benchmark and update README
./scripts/build-and-measure.sh
```

### Individual Operations

```bash
# Clean previous builds
./scripts/clean.sh

# Only update README (requires existing binaries)
./scripts/update-readme.sh

# Manual analysis
go run benchmark.go
```

## What It Does

1. **Builds Examples**: Compiles both standard and TinyString versions
2. **Measures Sizes**: Gets exact file sizes of all generated binaries
3. **Updates README**: Automatically replaces the "Binary Size Comparison" section with real data
4. **Shows Results**: Displays size comparison in terminal

## Example Output

```
🚀 Starting binary size benchmark...
✅ TinyGo found: tinygo version 0.30.0
🧹 Cleaning previous files...
📦 Building standard library example...
✅ Standard: Go binary and WebAssembly created
📦 Building TinyString example...
✅ TinyString: Go binary and WebAssembly created
📊 Analyzing sizes and updating README...

📊 Binary Size Results:
========================
standard.exe         native   standard     2.1MB
standard.wasm        wasm     standard     456.2KB
tinystring.exe       native   tinystring   1.2MB
tinystring.wasm      wasm     tinystring   187.4KB

✅ README updated with real binary size data
🎉 Benchmark completed successfully!
```

## How It Works

### Example Programs

Both examples perform identical operations:
- Text case transformations (upper/lower)
- Number to string conversions (int, float)
- String formatting and manipulation
- Text searching and replacement

### Size Measurement

The `benchmark.go` program:
- Scans for generated binaries in `examples/` directories
- Measures file sizes using `os.Stat()`
- Formats sizes in human-readable format (KB, MB)
- Updates the main README.md with real data

### README Integration

The script automatically finds and replaces this section in the main README:

```markdown
### Binary Size Comparison
```bash
# Traditional approach with standard library
go build -o app-standard main.go     # [REAL SIZE] binary
tinygo build -o app-standard.wasm -target wasm main.go  # [REAL SIZE] WebAssembly

# TinyString approach  
go build -o app-tiny main.go         # [REAL SIZE] binary  
tinygo build -o app-tiny.wasm -target wasm main.go      # [REAL SIZE] WebAssembly
```
```

## Troubleshooting

### TinyGo Not Found
If TinyGo is not installed, the benchmark will only measure Go native binaries:
```
❌ TinyGo is not installed. Building only standard Go binaries.
   Install TinyGo from: https://tinygo.org/getting-started/install/
```

### Build Failures
- Ensure you're in the `benchmark/` directory when running scripts
- Check that Go modules are properly initialized
- Verify TinyString library is available in the parent directory

### Permission Issues (Linux/macOS)
```bash
chmod +x scripts/*.sh
```

## Contributing

To add new benchmark scenarios:
1. Create additional example programs in `examples/`
2. Update `benchmark.go` to include new directories
3. Modify the README template in `generateBinarySizeSection()`
