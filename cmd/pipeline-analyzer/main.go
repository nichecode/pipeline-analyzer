package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/nichecode/pipeline-analyzer/internal/circleci"
	"github.com/nichecode/pipeline-analyzer/internal/fs"
)

var (
	version   = "dev"
	buildTime = "unknown"
)

func main() {
	var (
		configPath = flag.String("config", "", "Path to CircleCI config file (default: .circleci/config.yml)")
		outputDir  = flag.String("output-dir", "", "Output directory for analysis (default: auto-detect)")
		showVersion = flag.Bool("version", false, "Show version information")
		help       = flag.Bool("help", false, "Show help")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("pipeline-analyzer %s (built %s)\n", version, buildTime)
		return
	}

	if *help {
		printUsage()
		return
	}

	// Determine config file path
	configFile := *configPath
	if configFile == "" {
		// Check command line args first
		args := flag.Args()
		if len(args) > 0 {
			configFile = args[0]
		} else {
			// Default to .circleci/config.yml
			configFile = ".circleci/config.yml"
		}
	}

	// Check if config file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "❌ Config file not found: %s\n", configFile)
		fmt.Fprintf(os.Stderr, "\nTry:\n")
		fmt.Fprintf(os.Stderr, "  pipeline-analyzer path/to/config.yml\n")
		fmt.Fprintf(os.Stderr, "  pipeline-analyzer --config path/to/config.yml\n")
		os.Exit(1)
	}

	fmt.Printf("🔍 Analyzing CircleCI configuration: %s\n", configFile)

	// Parse the configuration
	config, err := circleci.ParseConfig(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to parse config: %v\n", err)
		os.Exit(1)
	}

	// Validate the configuration
	if err := circleci.IsValidConfig(config); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Invalid configuration: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Configuration parsed successfully\n")
	fmt.Printf("   - Version: %s\n", config.Version)
	fmt.Printf("   - Jobs: %d\n", len(config.Jobs))
	fmt.Printf("   - Workflows: %d\n", len(config.Workflows))
	fmt.Printf("   - Executors: %d\n", len(config.Executors))

	// Determine output directory
	outputPath, err := fs.DetermineOutputDir(*outputDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to determine output directory: %v\n", err)
		os.Exit(1)
	}

	// Validate output directory
	if err := fs.ValidateOutputDir(outputPath); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Output directory validation failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n📁 Output directory: %s\n", outputPath)

	// Perform analysis
	fmt.Printf("\n🔬 Performing analysis...\n")
	analysis := circleci.AnalyzeConfig(config)

	// Create writer and generate all files
	writer := fs.NewWriter(outputPath)
	
	fmt.Printf("📝 Generating documentation...\n")
	if err := writer.WriteAllFiles(analysis, configFile); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to write analysis files: %v\n", err)
		os.Exit(1)
	}

	// Clean up empty directories
	writer.CleanupEmptyDirs()

	// Print summary
	fs.PrintSummary(analysis, outputPath)
}

func printUsage() {
	fmt.Printf("pipeline-analyzer %s - CircleCI Configuration Analyzer\n\n", version)
	
	fmt.Printf("DESCRIPTION:\n")
	fmt.Printf("  Analyzes CircleCI configurations and generates comprehensive documentation\n")
	fmt.Printf("  to help migrate from CircleCI to local development workflows using go-task.\n\n")
	
	fmt.Printf("USAGE:\n")
	fmt.Printf("  pipeline-analyzer [options] [config-file]\n\n")
	
	fmt.Printf("EXAMPLES:\n")
	fmt.Printf("  pipeline-analyzer                           # Use .circleci/config.yml\n")
	fmt.Printf("  pipeline-analyzer config.yml               # Use specific config file\n")
	fmt.Printf("  pipeline-analyzer --config config.yml      # Alternative syntax\n")
	fmt.Printf("  pipeline-analyzer --output-dir ./analysis  # Custom output directory\n\n")
	
	fmt.Printf("OPTIONS:\n")
	flag.PrintDefaults()
	fmt.Printf("\n")
	
	fmt.Printf("OUTPUT:\n")
	fmt.Printf("  The tool generates analysis in:\n")
	fmt.Printf("  - .discovery/pipeline-analyzer/circleci/  (when in git repo)\n")
	fmt.Printf("  - circleci-analysis/                      (fallback)\n\n")
	
	fmt.Printf("  Key files generated:\n")
	fmt.Printf("  - README.md                    # Overview and navigation\n")
	fmt.Printf("  - MIGRATION-CHECKLIST.md      # Step-by-step migration guide\n")
	fmt.Printf("  - jobs/*.md                    # Individual job analysis\n")
	fmt.Printf("  - workflows/*.md               # Workflow structure\n")
	fmt.Printf("  - summaries/*.md               # Analysis summaries\n\n")
}

func init() {
	flag.Usage = printUsage
}