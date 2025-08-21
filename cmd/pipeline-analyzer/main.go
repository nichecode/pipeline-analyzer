package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nichecode/pipeline-analyzer/internal/discovery"
	"github.com/nichecode/pipeline-analyzer/internal/shared"
)

var (
	version   = "dev"
	buildTime = "unknown"
)

func main() {
	var (
		repoPath    = flag.String("path", ".", "Path to repository root (default: current directory)")
		showVersion = flag.Bool("version", false, "Show version information")
		help        = flag.Bool("help", false, "Show help")
		debug       = flag.Bool("debug", false, "Enable debug logging")
		logDir      = flag.String("log-dir", "", "Directory to write log files (default: no file logging)")
	)
	flag.Parse()

	// Initialize logging based on flags
	logLevel := shared.LogLevelInfo
	if *debug {
		logLevel = shared.LogLevelDebug
	}
	
	if err := shared.InitLogger(logLevel, *logDir); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		// Continue without file logging
	}
	defer shared.GetLogger().Close()

	if *showVersion {
		fmt.Printf("pipeline-analyzer %s (built %s)\n", version, buildTime)
		return
	}

	if *help {
		printUsage()
		return
	}

	// Handle command line args for repo path
	args := flag.Args()
	if len(args) > 0 {
		*repoPath = args[0]
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(*repoPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Invalid repository path: %v\n", err)
		os.Exit(1)
	}

	// Verify directory exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "❌ Directory does not exist: %s\n", absPath)
		os.Exit(1)
	}

	fmt.Printf("🔍 Scanning repository: %s\n", absPath)

	// Run auto-discovery and analysis
	runAutoDiscovery(absPath)
}

// runAutoDiscovery performs automatic discovery and analysis of all build tools
func runAutoDiscovery(repoPath string) {
	// Create scanner for the repository
	scanner := discovery.NewScanner(repoPath)
	
	// Scan repository and create discovery structure
	repo, discoveryDir, err := scanner.ScanAndCreateStructure()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to scan repository: %v\n", err)
		os.Exit(1)
	}

	// Check if any build tools were found
	if len(repo.BuildTools) == 0 {
		fmt.Printf("⚠️  No supported build tools found in repository\n")
		fmt.Printf("📁 Repository: %s\n", repoPath)
		fmt.Printf("🔍 Searched for: CircleCI, Go Task, GitHub Actions, npm, Composer, Cargo, Maven, Gradle, Makefile, Docker, Python, Terraform\n")
		os.Exit(0)
	}

	fmt.Printf("📁 Discovery directory: %s\n\n", discoveryDir)

	// Create analyzer and run analysis on all discovered tools
	analyzer := discovery.NewAnalyzer(repo, discoveryDir)
	
	results, err := analyzer.AnalyzeAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Analysis failed: %v\n", err)
		os.Exit(1)
	}

	// Generate overview
	if err := analyzer.GenerateOverview(results); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to generate overview: %v\n", err)
		os.Exit(1)
	}

	// Generate HTML index for better navigation
	if err := analyzer.GenerateHTMLIndex(results); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to generate HTML index: %v\n", err)
		os.Exit(1)
	}

	// Print final summary
	printFinalSummary(repo, discoveryDir, results)
}

// printFinalSummary prints the final summary of the analysis
func printFinalSummary(repo *discovery.Repository, discoveryDir string, results []discovery.AnalysisResult) {
	successCount := 0
	for _, result := range results {
		if result.Success {
			successCount++
		}
	}

	fmt.Printf("🎉 Analysis Complete!\n\n")
	fmt.Printf("📊 Summary:\n")
	fmt.Printf("   - Repository: %s\n", repo.RootPath)
	fmt.Printf("   - Build tools found: %d\n", len(repo.BuildTools))
	fmt.Printf("   - Successfully analyzed: %d\n", successCount)
	fmt.Printf("   - Failed: %d\n\n", len(results)-successCount)

	fmt.Printf("📁 Results location: %s\n", discoveryDir)
	fmt.Printf("🚀 Start here: %s/README.md\n\n", discoveryDir)

	if successCount > 0 {
		fmt.Printf("🔗 Quick links:\n")
		for _, result := range results {
			if result.Success {
				fmt.Printf("   - %s: %s/README.md\n", result.Tool.Name, result.OutputDir)
			}
		}
	}
}

func printUsage() {
	fmt.Printf("pipeline-analyzer %s - Automated Build Tool Discovery & Analysis\n\n", version)

	fmt.Printf("DESCRIPTION:\n")
	fmt.Printf("  Automatically discovers and analyzes all build tools in a repository.\n")
	fmt.Printf("  Creates comprehensive documentation with optimization recommendations\n")
	fmt.Printf("  in a standardized .discovery folder structure.\n\n")

	fmt.Printf("USAGE:\n")
	fmt.Printf("  pipeline-analyzer [repository-path]\n\n")

	fmt.Printf("EXAMPLES:\n")
	fmt.Printf("  pipeline-analyzer                    # Analyze current directory\n")
	fmt.Printf("  pipeline-analyzer /path/to/repo     # Analyze specific repository\n")
	fmt.Printf("  pipeline-analyzer ../my-project     # Analyze relative path\n")
	fmt.Printf("  pipeline-analyzer --debug --log-dir=logs /repo  # Enable debug logging\n\n")
	
	fmt.Printf("OPTIONS:\n")
	fmt.Printf("  --debug                             Enable debug logging output\n")
	fmt.Printf("  --log-dir=<directory>               Write detailed logs to specified directory\n")
	fmt.Printf("  --version                           Show version information\n")
	fmt.Printf("  --help                              Show this help message\n\n")

	fmt.Printf("SUPPORTED BUILD TOOLS:\n")
	fmt.Printf("  - CircleCI (.circleci/config.yml)\n")
	fmt.Printf("  - Go Task (Taskfile.yml)\n")
	fmt.Printf("  - GitHub Actions (.github/workflows/)\n")
	fmt.Printf("  - npm (package.json)\n")
	fmt.Printf("  - Composer (composer.json)\n")
	fmt.Printf("  - Cargo (Cargo.toml)\n")
	fmt.Printf("  - Maven (pom.xml)\n")
	fmt.Printf("  - Gradle (build.gradle)\n")
	fmt.Printf("  - Makefile (Makefile)\n")
	fmt.Printf("  - Docker (Dockerfile, docker-compose.yml)\n")
	fmt.Printf("  - Python (requirements.txt, pyproject.toml)\n")
	fmt.Printf("  - Terraform (*.tf)\n\n")

	fmt.Printf("OPTIONS:\n")
	flag.PrintDefaults()
	fmt.Printf("\n")

	fmt.Printf("OUTPUT:\n")
	fmt.Printf("  All analysis results are placed in:\n")
	fmt.Printf("  📁 .discovery/pipeline-analyzer/\n")
	fmt.Printf("  ├── README.md                    # Discovery overview\n")
	fmt.Printf("  ├── circleci/                    # CircleCI analysis (if found)\n")
	fmt.Printf("  ├── gotask/                      # Go Task analysis (if found)\n")
	fmt.Printf("  ├── npm/                         # npm analysis (if found)\n")
	fmt.Printf("  └── [tool-name]/                 # Other discovered tools\n\n")

	fmt.Printf("  Each tool directory contains:\n")
	fmt.Printf("  - README.md                      # Tool-specific overview\n")
	fmt.Printf("  - OPTIMIZATION-GUIDE.md         # Performance recommendations\n")
	fmt.Printf("  - tasks/ or jobs/                # Individual analysis files\n")
	fmt.Printf("  - summaries/                     # Analysis summaries\n\n")

	fmt.Printf("GETTING STARTED:\n")
	fmt.Printf("  1. Run: pipeline-analyzer\n")
	fmt.Printf("  2. Open: .discovery/pipeline-analyzer/README.md\n")
	fmt.Printf("  3. Follow the generated recommendations\n\n")
}

func init() {
	flag.Usage = printUsage
}
