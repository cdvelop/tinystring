package main

// TestData contains shared test data used across all benchmarks
type TestData struct {
	Strings []string
	Numbers []int
	Floats  []float64
}

// BenchmarkConfig holds configuration for benchmark runs
type BenchmarkConfig struct {
	Iterations int
	DataSize   int
	Parallel   bool
	MemoryTest bool
}

// TestScenarios defines common test scenarios for benchmarks
type TestScenario struct {
	Name        string
	Description string
	Input       interface{}
	Expected    interface{}
}
