package main

import (
	"fmt"
	"os"
	"strings" // Only for section finding in README
	"time"

	"github.com/cdvelop/tinystring"
)

// ReportGenerator handles README and documentation generation
type ReportGenerator struct {
	ReadmePath string
	TempPath   string
}

// NewReportGenerator creates a new report generator
func NewReportGenerator(readmePath string) *ReportGenerator {
	return &ReportGenerator{
		ReadmePath: readmePath,
		TempPath:   readmePath + ".tmp",
	}
}

// UpdateREADMEWithBinaryData updates README with binary size comparison data
func (r *ReportGenerator) UpdateBinaryData(binaries []BinaryInfo) error {
	LogInfo("Updating README with binary size analysis...")

	content, err := r.generateBinarySizeSection(binaries)
	if err != nil {
		return tinystring.Err(err)
	}

	return r.updateREADMESection("Binary Size Comparison", content)
}

// UpdateREADMEWithMemoryData updates README with memory benchmark data
func (r *ReportGenerator) UpdateMemoryData(comparisons []MemoryComparison) error {
	LogInfo("Updating README with memory allocation analysis...")

	content, err := r.generateMemorySection(comparisons)
	if err != nil {
		return fmt.Errorf("failed to generate memory section: %v", err)
	}

	return r.updateREADMESection("Memory Usage Comparison", content)
}

// generateBinarySizeSection creates the binary size comparison section
func (r *ReportGenerator) generateBinarySizeSection(binaries []BinaryInfo) (string, error) {
	var content strings.Builder

	content.WriteString("## Binary Size Comparison\n\n")
	content.WriteString("[Standard Library Example](benchmark/bench-binary-size/standard-lib/main.go) | [TinyString Example](benchmark/bench-binary-size/tinystring-lib/main.go)\n\n")
	content.WriteString("<!-- This table is automatically generated from build-and-measure.sh -->\n")
	content.WriteString("*Last updated: " + time.Now().Format("2006-01-02 15:04:05") + "*\n\n")

	// Group binaries by optimization level
	optimizations := getOptimizationConfigs()

	content.WriteString("| Build Type | Parameters | Standard Library<br/>`go build` | TinyString<br/>`tinygo build` | Size Reduction | Performance |\n")
	content.WriteString("|------------|------------|------------------|------------|----------------|-------------|\n")

	var allImprovements []float64
	var maxImprovement float64
	var totalSavings int64

	for _, opt := range optimizations {
		// Find matching binaries for this optimization level
		standardNative := findBinaryByPattern(binaries, "standard", "native", opt.Suffix)
		tinystringNative := findBinaryByPattern(binaries, "tinystring", "native", opt.Suffix)
		standardWasm := findBinaryByPattern(binaries, "standard", "wasm", opt.Suffix)
		tinystringWasm := findBinaryByPattern(binaries, "tinystring", "wasm", opt.Suffix)

		// Build type icons and names
		buildIcon := getBuildTypeIcon(opt.Name)
		parameters := getBuildParameters(opt.Name, false)    // Native
		wasmParameters := getBuildParameters(opt.Name, true) // WASM

		// Native builds
		if standardNative.Name != "" && tinystringNative.Name != "" {
			improvementPercent := calculateImprovementPercent(standardNative.Size, tinystringNative.Size)
			sizeDiff := standardNative.Size - tinystringNative.Size
			performanceIndicator := getPerformanceIndicator(improvementPercent)

			content.WriteString(fmt.Sprintf("| %s **%s Native** | `%s` | %s | %s | **-%s** | %s **%.1f%%** |\n",
				buildIcon, capitalizeFirst(opt.Name), parameters,
				standardNative.SizeStr, tinystringNative.SizeStr,
				FormatSize(sizeDiff), performanceIndicator, improvementPercent))

			allImprovements = append(allImprovements, improvementPercent)
			if improvementPercent > maxImprovement {
				maxImprovement = improvementPercent
			}
			totalSavings += sizeDiff
		}

		// WebAssembly builds
		if standardWasm.Name != "" && tinystringWasm.Name != "" {
			improvementPercent := calculateImprovementPercent(standardWasm.Size, tinystringWasm.Size)
			sizeDiff := standardWasm.Size - tinystringWasm.Size
			performanceIndicator := getPerformanceIndicator(improvementPercent)

			content.WriteString(fmt.Sprintf("| 🌐 **%s WASM** | `%s` | %s | %s | **-%s** | %s **%.1f%%** |\n",
				capitalizeFirst(opt.Name), wasmParameters,
				standardWasm.SizeStr, tinystringWasm.SizeStr,
				FormatSize(sizeDiff), performanceIndicator, improvementPercent))

			allImprovements = append(allImprovements, improvementPercent)
			if improvementPercent > maxImprovement {
				maxImprovement = improvementPercent
			}
			totalSavings += sizeDiff
		}
	}

	// Calculate averages
	var avgImprovement float64
	var avgWasmImprovement float64
	var avgNativeImprovement float64
	var wasmCount, nativeCount int

	for i, opt := range optimizations {
		standardNative := findBinaryByPattern(binaries, "standard", "native", opt.Suffix)
		tinystringNative := findBinaryByPattern(binaries, "tinystring", "native", opt.Suffix)
		standardWasm := findBinaryByPattern(binaries, "standard", "wasm", opt.Suffix)
		tinystringWasm := findBinaryByPattern(binaries, "tinystring", "wasm", opt.Suffix)

		if standardNative.Name != "" && tinystringNative.Name != "" {
			improvement := calculateImprovementPercent(standardNative.Size, tinystringNative.Size)
			avgNativeImprovement += improvement
			nativeCount++
		}

		if standardWasm.Name != "" && tinystringWasm.Name != "" {
			improvement := calculateImprovementPercent(standardWasm.Size, tinystringWasm.Size)
			avgWasmImprovement += improvement
			wasmCount++
		}
		_ = i
	}

	if len(allImprovements) > 0 {
		for _, imp := range allImprovements {
			avgImprovement += imp
		}
		avgImprovement /= float64(len(allImprovements))
	}

	if nativeCount > 0 {
		avgNativeImprovement /= float64(nativeCount)
	}
	if wasmCount > 0 {
		avgWasmImprovement /= float64(wasmCount)
	}

	// Performance summary
	content.WriteString("\n### 🎯 Performance Summary\n\n")
	content.WriteString(fmt.Sprintf("- 🏆 **Peak Reduction: %.1f%%** (Best optimization)\n", maxImprovement))
	if wasmCount > 0 {
		content.WriteString(fmt.Sprintf("- ✅ **Average WebAssembly Reduction: %.1f%%**\n", avgWasmImprovement))
	}
	if nativeCount > 0 {
		content.WriteString(fmt.Sprintf("- ✅ **Average Native Reduction: %.1f%%**\n", avgNativeImprovement))
	}
	content.WriteString(fmt.Sprintf("- 📦 **Total Size Savings: %s across all builds**\n\n", FormatSize(totalSavings)))

	content.WriteString("#### Performance Legend\n")
	content.WriteString("- ❌ Poor (<5% reduction)\n")
	content.WriteString("- ➖ Fair (5-15% reduction)\n")
	content.WriteString("- ✅ Good (15-70% reduction)\n")
	content.WriteString("- 🏆 Outstanding (>70% reduction)\n\n")

	return content.String(), nil
}

// generateMemorySection creates the memory allocation comparison section
func (r *ReportGenerator) generateMemorySection(comparisons []MemoryComparison) (string, error) {
	var content strings.Builder

	content.WriteString("## Memory Usage Comparison\n\n")
	content.WriteString("[Standard Library Example](benchmark/bench-memory-alloc/standard) | [TinyString Example](benchmark/bench-memory-alloc/tinystring)\n\n")
	content.WriteString("<!-- This table is automatically generated from memory-benchmark.sh -->\n")
	content.WriteString("*Last updated: " + time.Now().Format("2006-01-02 15:04:05") + "*\n\n")
	content.WriteString("Performance benchmarks comparing memory allocation patterns between standard Go library and TinyString:\n\n")

	// Enhanced table with better styling and icons
	content.WriteString("| 🧪 **Benchmark Category** | 📚 **Library** | 💾 **Memory/Op** | 🔢 **Allocs/Op** | ⏱️ **Time/Op** | 📈 **Memory Trend** | 🎯 **Alloc Trend** | 🏆 **Performance** |\n")
	content.WriteString("|----------------------------|----------------|-------------------|-------------------|-----------------|---------------------|---------------------|--------------------|\n")

	var totalMemoryDiff float64
	var totalAllocDiff float64
	var benchmarkCount int

	for _, comparison := range comparisons {
		if comparison.Standard.Name != "" && comparison.TinyString.Name != "" {
			memImprovement := calculateMemoryImprovement(
				comparison.Standard.BytesPerOp, comparison.TinyString.BytesPerOp)
			allocImprovement := calculateMemoryImprovement(
				comparison.Standard.AllocsPerOp, comparison.TinyString.AllocsPerOp)

			// Calculate percentage changes for tracking
			memPercent := calculateMemoryPercent(comparison.Standard.BytesPerOp, comparison.TinyString.BytesPerOp)
			allocPercent := calculateMemoryPercent(comparison.Standard.AllocsPerOp, comparison.TinyString.AllocsPerOp)

			totalMemoryDiff += memPercent
			totalAllocDiff += allocPercent
			benchmarkCount++

			// Get performance indicators
			memoryIndicator := getMemoryPerformanceIndicator(memPercent)
			allocIndicator := getAllocPerformanceIndicator(allocPercent)
			overallIndicator := getOverallPerformanceIndicator(memPercent, allocPercent)

			// Category with emoji
			categoryIcon := getBenchmarkCategoryIcon(comparison.Category)

			// Standard library row with enhanced styling
			content.WriteString(fmt.Sprintf("| %s **%s** | 📊 Standard | `%s` | `%d` | `%s` | - | - | - |\n",
				categoryIcon,
				comparison.Category,
				FormatSize(comparison.Standard.BytesPerOp),
				comparison.Standard.AllocsPerOp,
				formatNanoTime(comparison.Standard.NsPerOp)))

			// TinyString row with improvements and visual indicators
			content.WriteString(fmt.Sprintf("| | 🚀 TinyString | `%s` | `%d` | `%s` | %s **%s** | %s **%s** | %s |\n",
				FormatSize(comparison.TinyString.BytesPerOp),
				comparison.TinyString.AllocsPerOp,
				formatNanoTime(comparison.TinyString.NsPerOp),
				memoryIndicator, memImprovement,
				allocIndicator, allocImprovement,
				overallIndicator))
		}
	}

	// Calculate averages for summary
	var avgMemoryDiff, avgAllocDiff float64
	if benchmarkCount > 0 {
		avgMemoryDiff = totalMemoryDiff / float64(benchmarkCount)
		avgAllocDiff = totalAllocDiff / float64(benchmarkCount)
	}

	// Performance summary section with enhanced styling
	content.WriteString("\n### 🎯 Performance Summary\n\n")

	// Memory efficiency classification
	memoryClass := getMemoryEfficiencyClass(avgMemoryDiff)
	allocClass := getAllocEfficiencyClass(avgAllocDiff)

	content.WriteString(fmt.Sprintf("- 💾 **Memory Efficiency**: %s (%.1f%% average change)\n", memoryClass, avgMemoryDiff))
	content.WriteString(fmt.Sprintf("- 🔢 **Allocation Efficiency**: %s (%.1f%% average change)\n", allocClass, avgAllocDiff))
	content.WriteString(fmt.Sprintf("- 📊 **Benchmarks Analyzed**: %d categories\n", benchmarkCount))
	content.WriteString("- 🎯 **Optimization Focus**: Binary size reduction vs runtime efficiency\n\n")

	// Enhanced trade-offs analysis with better formatting
	content.WriteString("### ⚖️ Trade-offs Analysis\n\n")
	content.WriteString("The benchmarks reveal important trade-offs between **binary size** and **runtime performance**:\n\n")

	content.WriteString("#### 📦 **Binary Size Benefits** ✅\n")
	content.WriteString("- 🏆 **16-84% smaller** compiled binaries\n")
	content.WriteString("- 🌐 **Superior WebAssembly** compression ratios\n")
	content.WriteString("- 🚀 **Faster deployment** and distribution\n")
	content.WriteString("- 💾 **Lower storage** requirements\n\n")

	content.WriteString("#### 🧠 **Runtime Memory Considerations** ⚠️\n")
	content.WriteString("- 📈 **Higher allocation overhead** during execution\n")
	content.WriteString("- 🗑️ **Increased GC pressure** due to allocation patterns\n")
	content.WriteString("- ⚡ **Trade-off optimizes** for distribution size over runtime efficiency\n")
	content.WriteString("- 🔄 **Different optimization strategy** than standard library\n\n")

	content.WriteString("#### 🎯 **Optimization Recommendations**\n")
	content.WriteString("| 🎯 **Use Case** | 💡 **Recommendation** | 🔧 **Best For** |\n")
	content.WriteString("|-----------------|------------------------|------------------|\n")
	content.WriteString("| 🌐 WebAssembly Apps | ✅ **TinyString** | Size-critical web deployment |\n")
	content.WriteString("| 📱 Embedded Systems | ✅ **TinyString** | Resource-constrained devices |\n")
	content.WriteString("| ☁️ Edge Computing | ✅ **TinyString** | Fast startup and deployment |\n")
	content.WriteString("| 🏢 Memory-Intensive Server | ⚠️ **Standard Library** | High-throughput applications |\n")
	content.WriteString("| 🔄 High-Frequency Processing | ⚠️ **Standard Library** | Performance-critical workloads |\n\n")

	content.WriteString("#### 📊 **Performance Legend**\n")
	content.WriteString("- 🏆 **Excellent** (Better performance)\n")
	content.WriteString("- ✅ **Good** (Acceptable trade-off)\n")
	content.WriteString("- ⚠️ **Caution** (Higher resource usage)\n")
	content.WriteString("- ❌ **Poor** (Significant overhead)\n\n")

	return content.String(), nil
}

// updateREADMESection updates a specific section in the README
func (r *ReportGenerator) updateREADMESection(sectionTitle, newContent string) error {
	// Read current README
	existingContent, err := os.ReadFile(r.ReadmePath)
	if err != nil {
		LogError(fmt.Sprintf("Failed to read README: %v", err))
		return err
	}

	content := string(existingContent)

	// Find section boundaries
	sectionStart := "## " + sectionTitle
	startIndex := strings.Index(content, sectionStart)

	if startIndex == -1 {
		// Section doesn't exist, append to end
		content += "\n" + newContent
	} else {
		// Find next section or end of file
		nextSectionIndex := strings.Index(content[startIndex+len(sectionStart):], "\n## ")
		var endIndex int

		if nextSectionIndex == -1 {
			endIndex = len(content)
		} else {
			endIndex = startIndex + len(sectionStart) + nextSectionIndex
		}

		// Replace the section
		content = content[:startIndex] + newContent + content[endIndex:]
	}

	// Write updated content
	err = os.WriteFile(r.TempPath, []byte(content), 0644)
	if err != nil {
		LogError(fmt.Sprintf("Failed to write temporary README: %v", err))
		return err
	}

	// Replace original with temporary
	err = os.Rename(r.TempPath, r.ReadmePath)
	if err != nil {
		LogError(fmt.Sprintf("Failed to replace README: %v", err))
		return err
	}

	LogSuccess(fmt.Sprintf("Updated README section: %s", sectionTitle))
	return nil
}

// capitalizeFirst capitalizes the first letter of a string
func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[0] >= 'a' && s[0] <= 'z' {
		return string(s[0]-32) + s[1:]
	}
	return s
}

// Helper functions for binary size reporting

// getBuildTypeIcon returns the appropriate icon for build type
func getBuildTypeIcon(optName string) string {
	switch optName {
	case "Default":
		return "🖥️"
	case "Speed":
		return "⚡"
	case "Ultra":
		return "🏁"
	case "Debug":
		return "🔧"
	default:
		return "📦"
	}
}

// getBuildParameters returns the build parameters for different optimization levels
func getBuildParameters(optName string, isWasm bool) string {
	switch optName {
	case "Default":
		if isWasm {
			return "(default -opt=z)"
		}
		return `-ldflags="-s -w"`
	case "Speed":
		if isWasm {
			return "-opt=2 -target wasm"
		}
		return `-ldflags="-s -w"`
	case "Ultra":
		if isWasm {
			return "-no-debug -panic=trap -scheduler=none -gc=leaking -target wasm"
		}
		return `-ldflags="-s -w"`
	case "Debug":
		if isWasm {
			return "-opt=0 -target wasm"
		}
		return `-ldflags="-s -w"`
	default:
		return ""
	}
}

// calculateImprovementPercent calculates the percentage improvement
func calculateImprovementPercent(standardSize, tinystringSize int64) float64 {
	if standardSize <= 0 {
		return 0
	}
	return float64(standardSize-tinystringSize) / float64(standardSize) * 100
}

// getPerformanceIndicator returns the appropriate performance indicator
func getPerformanceIndicator(improvementPercent float64) string {
	switch {
	case improvementPercent < 5:
		return "❌"
	case improvementPercent < 15:
		return "➖"
	case improvementPercent < 70:
		return "✅"
	default:
		return "🏆"
	}
}

// Helper functions for enhanced memory reporting

// calculateMemoryPercent calculates the percentage change in memory usage
func calculateMemoryPercent(standardValue, tinystringValue int64) float64 {
	if standardValue <= 0 {
		return 0
	}
	return float64(tinystringValue-standardValue) / float64(standardValue) * 100
}

// getBenchmarkCategoryIcon returns appropriate icon for benchmark category
func getBenchmarkCategoryIcon(category string) string {
	switch {
	case strings.Contains(category, "String"):
		return "📝"
	case strings.Contains(category, "Number"):
		return "🔢"
	case strings.Contains(category, "Mixed"):
		return "🔄"
	case strings.Contains(category, "Pointer"):
		return "👉"
	default:
		return "🧪"
	}
}

// getMemoryPerformanceIndicator returns indicator for memory performance
func getMemoryPerformanceIndicator(percentChange float64) string {
	switch {
	case percentChange < -20: // 20% improvement (less memory)
		return "🏆"
	case percentChange < -5: // 5% improvement
		return "✅"
	case percentChange < 5: // Similar usage
		return "➖"
	case percentChange < 50: // Up to 50% more
		return "⚠️"
	default: // Over 50% more
		return "❌"
	}
}

// getAllocPerformanceIndicator returns indicator for allocation performance
func getAllocPerformanceIndicator(percentChange float64) string {
	switch {
	case percentChange < -15: // 15% fewer allocations
		return "🏆"
	case percentChange < -5: // 5% fewer allocations
		return "✅"
	case percentChange < 5: // Similar allocations
		return "➖"
	case percentChange < 25: // Up to 25% more
		return "⚠️"
	default: // Over 25% more
		return "❌"
	}
}

// getOverallPerformanceIndicator combines memory and allocation indicators
func getOverallPerformanceIndicator(memPercent, allocPercent float64) string {
	// Average the two percentages for overall assessment
	avgChange := (memPercent + allocPercent) / 2

	switch {
	case avgChange < -15: // Overall improvement
		return "🏆 **Excellent**"
	case avgChange < -5: // Slight improvement
		return "✅ **Good**"
	case avgChange < 15: // Acceptable trade-off
		return "➖ **Fair**"
	case avgChange < 40: // Higher resource usage
		return "⚠️ **Caution**"
	default: // Significant overhead
		return "❌ **Poor**"
	}
}

// getMemoryEfficiencyClass classifies memory efficiency
func getMemoryEfficiencyClass(avgPercent float64) string {
	switch {
	case avgPercent < -10:
		return "🏆 **Excellent** (Lower memory usage)"
	case avgPercent < 0:
		return "✅ **Good** (Memory efficient)"
	case avgPercent < 20:
		return "➖ **Fair** (Acceptable overhead)"
	case avgPercent < 50:
		return "⚠️ **Caution** (Higher memory usage)"
	default:
		return "❌ **Poor** (Significant overhead)"
	}
}

// getAllocEfficiencyClass classifies allocation efficiency
func getAllocEfficiencyClass(avgPercent float64) string {
	switch {
	case avgPercent < -10:
		return "🏆 **Excellent** (Fewer allocations)"
	case avgPercent < 0:
		return "✅ **Good** (Allocation efficient)"
	case avgPercent < 15:
		return "➖ **Fair** (Acceptable allocation pattern)"
	case avgPercent < 35:
		return "⚠️ **Caution** (More allocations)"
	default:
		return "❌ **Poor** (Excessive allocations)"
	}
}
