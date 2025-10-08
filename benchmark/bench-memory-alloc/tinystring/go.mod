module memory-bench-tinystring

go 1.25.2

require (
	benchmark/shared v0.0.0
	github.com/cdvelop/tinystring v0.0.0
)

// Use local TinyString module
replace github.com/cdvelop/tinystring => ../../..

// Use local shared module
replace benchmark/shared => ../../shared
