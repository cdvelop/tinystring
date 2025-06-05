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
	content.WriteString("[Standard Library Example](bench-binary-size/standard-lib/main.go) | [TinyString Example](bench-binary-size/tinystring-lib/main.go)\n\n")
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
	content.WriteString("*Last updated: " + time.Now().Format("2006-01-02 15:04:05") + "*\n\n")
	content.WriteString("Performance benchmarks comparing memory allocation patterns:\n\n")

	content.WriteString("| Benchmark | Library | Bytes/Op | Allocs/Op | Time/Op | Memory Improvement | Alloc Improvement |\n")
	content.WriteString("|-----------|---------|----------|-----------|---------|-------------------|------------------|\n")

	for _, comparison := range comparisons {
		if comparison.Standard.Name != "" && comparison.TinyString.Name != "" {
			memImprovement := calculateMemoryImprovement(
				comparison.Standard.BytesPerOp, comparison.TinyString.BytesPerOp)
			allocImprovement := calculateMemoryImprovement(
				comparison.Standard.AllocsPerOp, comparison.TinyString.AllocsPerOp)

			// Standard library row
			content.WriteString(fmt.Sprintf("| **%s** | Standard | %s | %d | %s | - | - |\n",
				comparison.Category,
				FormatSize(comparison.Standard.BytesPerOp),
				comparison.Standard.AllocsPerOp,
				formatNanoTime(comparison.Standard.NsPerOp)))

			// TinyString row with improvements
			content.WriteString(fmt.Sprintf("| | TinyString | %s | %d | %s | **%s** | **%s** |\n",
				FormatSize(comparison.TinyString.BytesPerOp),
				comparison.TinyString.AllocsPerOp,
				formatNanoTime(comparison.TinyString.NsPerOp),
				memImprovement,
				allocImprovement))
		}
	}
	content.WriteString("\n### Trade-offs Analysis\n\n")
	content.WriteString("The benchmarks reveal important trade-offs between binary size and runtime performance:\n\n")
	content.WriteString("**Binary Size Benefits:**\n")
	content.WriteString("- Significantly smaller compiled binaries (16-84% reduction)\n")
	content.WriteString("- Better compression for WebAssembly targets\n")
	content.WriteString("- Reduced distribution and deployment overhead\n\n")
	content.WriteString("**Runtime Memory Considerations:**\n")
	content.WriteString("- Higher memory allocation overhead during execution\n")
	content.WriteString("- Increased GC pressure due to more allocations\n")
	content.WriteString("- Trade-off optimizes for storage/distribution size over runtime efficiency\n\n")
	content.WriteString("**Recommendation:**\n")
	content.WriteString("- Use TinyString for size-constrained environments (embedded, edge computing)\n")
	content.WriteString("- Consider standard library for memory-intensive runtime workloads\n")
	content.WriteString("- Evaluate based on specific deployment constraints\n\n")

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
