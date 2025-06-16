module memory-bench-standard

go 1.21

require benchmark/shared v0.0.0

// Use local shared module
replace benchmark/shared => ../../shared
