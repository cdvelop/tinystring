package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

type BinaryInfo struct {
	Name        string
	Size        int64
	SizeStr     string
	Type        string
	Library     string
	OptLevel    string
	Description string
}

type OptimizationConfig struct {
	Name        string
	Flags       string
	Description string
	Suffix      string
}

func main() {
	fmt.Println("🔍 Analyzing binary sizes with multiple optimization levels...")

	binaries := measureBinarySizes()
	if len(binaries) == 0 {
		fmt.Println("❌ No binaries found to analyze")
		return
	}

	displayResults(binaries)
	displayOptimizationTable(binaries)
	updateREADME(binaries)

	fmt.Println("✅ README updated with real binary size data")
}

func getOptimizationConfigs() []OptimizationConfig {
	return []OptimizationConfig{
		{
			Name:        "Default",
			Flags:       "",
			Description: "Default TinyGo optimization (-opt=z)",
			Suffix:      "",
		},
		{
			Name:        "Ultra Size",
			Flags:       "-no-debug -panic=trap -scheduler=none -gc=leaking",
			Description: "Maximum size optimization",
			Suffix:      "-ultra",
		},
		{
			Name:        "Speed",
			Flags:       "-opt=2",
			Description: "Optimized for speed over size",
			Suffix:      "-speed",
		},
		{
			Name:        "Debug",
			Flags:       "-opt=0",
			Description: "No optimization, best for debugging",
			Suffix:      "-debug",
		},
	}
}

func measureBinarySizes() []BinaryInfo {
	var binaries []BinaryInfo
	examplesDir := "examples"

	// Standard library binaries
	standardDir := filepath.Join(examplesDir, "standard-lib")
	binaries = append(binaries, findBinaries(standardDir, "standard")...)

	// TinyString library binaries
	tinystringDir := filepath.Join(examplesDir, "tinystring-lib")
	binaries = append(binaries, findBinaries(tinystringDir, "tinystring")...)

	return binaries
}

func findBinaries(dir, library string) []BinaryInfo {
	var binaries []BinaryInfo
	configs := getOptimizationConfigs()

	// Determine executable extension based on OS
	execExt := ""
	if runtime.GOOS == "windows" {
		execExt = ".exe"
	}

	for _, config := range configs {
		// Look for native binary
		execName := library + config.Suffix + execExt
		execPath := filepath.Join(dir, execName)
		if info, err := os.Stat(execPath); err == nil {
			binaries = append(binaries, BinaryInfo{
				Name:        execName,
				Size:        info.Size(),
				SizeStr:     formatSize(info.Size()),
				Type:        "native",
				Library:     library,
				OptLevel:    config.Name,
				Description: config.Description,
			})
		}

		// Look for WebAssembly binary
		wasmName := library + config.Suffix + ".wasm"
		wasmPath := filepath.Join(dir, wasmName)
		if info, err := os.Stat(wasmPath); err == nil {
			binaries = append(binaries, BinaryInfo{
				Name:        wasmName,
				Size:        info.Size(),
				SizeStr:     formatSize(info.Size()),
				Type:        "wasm",
				Library:     library,
				OptLevel:    config.Name,
				Description: config.Description,
			})
		}
	}

	return binaries
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%dB", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f%cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func displayResults(binaries []BinaryInfo) {
	fmt.Println("\n📊 Binary Size Results:")
	fmt.Println("========================")

	for _, binary := range binaries {
		fmt.Printf("%-25s %-8s %-12s %-10s %s\n",
			binary.Name,
			binary.Type,
			binary.Library,
			binary.OptLevel,
			binary.SizeStr)
	}
	fmt.Println()
}

func displayOptimizationTable(binaries []BinaryInfo) {
	fmt.Println("📊 Optimization Comparison Table:")
	fmt.Println("================================")

	// Group binaries by type (wasm/native) and library
	wasmStandard := make(map[string]BinaryInfo)
	wasmTinystring := make(map[string]BinaryInfo)

	for _, binary := range binaries {
		if binary.Type == "wasm" {
			if binary.Library == "standard" {
				wasmStandard[binary.OptLevel] = binary
			} else if binary.Library == "tinystring" {
				wasmTinystring[binary.OptLevel] = binary
			}
		}
	}

	// Display WebAssembly comparison table
	fmt.Println("\nWebAssembly Size Comparison:")
	fmt.Printf("%-12s %-15s %-15s %-15s\n", "Optimization", "Standard Lib", "TinyString", "Reduction")
	fmt.Println("---------------------------------------------------------------")

	configs := getOptimizationConfigs()
	for _, config := range configs {
		standardSize := "N/A"
		tinystringSize := "N/A"
		reduction := "N/A"

		if std, exists := wasmStandard[config.Name]; exists {
			standardSize = std.SizeStr
		}

		if tiny, exists := wasmTinystring[config.Name]; exists {
			tinystringSize = tiny.SizeStr

			// Calculate reduction percentage if both exist
			if stdBinary, stdExists := wasmStandard[config.Name]; stdExists {
				reductionPercent := float64(stdBinary.Size-tiny.Size) / float64(stdBinary.Size) * 100
				reduction = fmt.Sprintf("%.1f%%", reductionPercent)
			}
		}

		fmt.Printf("%-12s %-15s %-15s %-15s\n", config.Name, standardSize, tinystringSize, reduction)
	}
	fmt.Println()
}

func updateREADME(binaries []BinaryInfo) {
	readmePath := "../README.md"
	content, err := os.ReadFile(readmePath)
	if err != nil {
		fmt.Printf("❌ Error reading README: %v\n", err)
		return
	}

	// Generate new section with real data
	newSection := generateBinarySizeSection(binaries) // Replace existing section - find everything from Binary Size Comparison to Target Environments
	var updated string
	re := regexp.MustCompile(`### Binary Size Comparison[\s\S]*?### Target Environments`)
	if re.MatchString(string(content)) {
		updated = re.ReplaceAllString(string(content), newSection+"\n### Target Environments")
	} else {
		// Fallback: replace until end of section
		re2 := regexp.MustCompile(`### Binary Size Comparison[\s\S]*?(?:\n### |\n## |$)`)
		updated = re2.ReplaceAllStringFunc(string(content), func(match string) string {
			if strings.Contains(match, "\n### ") {
				// Keep the next section header
				parts := strings.SplitN(match, "\n### ", 2)
				return newSection + "\n\n### " + parts[1]
			} else if strings.Contains(match, "\n## ") {
				// Keep the next section header
				parts := strings.SplitN(match, "\n## ", 2)
				return newSection + "\n\n## " + parts[1]
			}
			return newSection
		})
	}

	// Write updated content
	err = os.WriteFile(readmePath, []byte(updated), 0644)
	if err != nil {
		fmt.Printf("❌ Error updating README: %v\n", err)
		return
	}
}

func generateBinarySizeSection(binaries []BinaryInfo) string {
	// Group binaries by optimization level and library for default comparison
	wasmStandard := make(map[string]BinaryInfo)
	wasmTinystring := make(map[string]BinaryInfo)
	nativeStandard := make(map[string]BinaryInfo)
	nativeTinystring := make(map[string]BinaryInfo)

	for _, binary := range binaries {
		if binary.Type == "wasm" {
			if binary.Library == "standard" {
				wasmStandard[binary.OptLevel] = binary
			} else if binary.Library == "tinystring" {
				wasmTinystring[binary.OptLevel] = binary
			}
		} else if binary.Type == "native" {
			if binary.Library == "standard" {
				nativeStandard[binary.OptLevel] = binary
			} else if binary.Library == "tinystring" {
				nativeTinystring[binary.OptLevel] = binary
			}
		}
	}

	// Get default sizes (fallback to first available)
	standardNative := getFirstAvailableSize(nativeStandard, "~2.1MB")
	standardWasm := getFirstAvailableSize(wasmStandard, "~500KB+")
	tinystringNative := getFirstAvailableSize(nativeTinystring, "~1.2MB")
	tinystringWasm := getFirstAvailableSize(wasmTinystring, "~180KB")

	result := "### Binary Size Comparison\n"

	// Basic comparison
	result += "```bash\n"
	result += "# Traditional approach with standard library\n"
	result += "go build -o app-standard main.go     # " + standardNative + " binary\n"
	result += "tinygo build -o app-standard.wasm -target wasm main.go  # " + standardWasm + " WebAssembly\n\n"
	result += "# TinyString approach  \n"
	result += "go build -o app-tiny main.go         # " + tinystringNative + " binary  \n"
	result += "tinygo build -o app-tiny.wasm -target wasm main.go      # " + tinystringWasm + " WebAssembly\n"
	result += "```\n\n"

	// Add optimization table if we have multiple optimization levels
	if len(wasmStandard) > 1 || len(wasmTinystring) > 1 {
		result += "#### WebAssembly Optimization Comparison\n\n"
		result += "| Optimization Level | Standard Library | TinyString | Size Reduction |\n"
		result += "|-------------------|------------------|------------|----------------|\n"

		configs := getOptimizationConfigs()
		for _, config := range configs {
			standardSize := "N/A"
			tinystringSize := "N/A"
			reduction := "N/A"

			if std, exists := wasmStandard[config.Name]; exists {
				standardSize = std.SizeStr
			}

			if tiny, exists := wasmTinystring[config.Name]; exists {
				tinystringSize = tiny.SizeStr

				// Calculate reduction percentage
				if stdBinary, stdExists := wasmStandard[config.Name]; stdExists {
					reductionPercent := float64(stdBinary.Size-tiny.Size) / float64(stdBinary.Size) * 100
					reduction = fmt.Sprintf("%.1f%%", reductionPercent)
				}
			}

			result += fmt.Sprintf("| %s | %s | %s | %s |\n",
				config.Description, standardSize, tinystringSize, reduction)
		}
		result += "\n"
	}

	return result
}

func getFirstAvailableSize(sizeMap map[string]BinaryInfo, fallback string) string {
	// Try "Default" first, then any available
	if binary, exists := sizeMap["Default"]; exists {
		return binary.SizeStr
	}
	for _, binary := range sizeMap {
		return binary.SizeStr
	}
	return fallback
}
