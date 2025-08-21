package discovery

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/nichecode/pipeline-analyzer/internal/circleci"
	"github.com/nichecode/pipeline-analyzer/internal/githubactions"
	"github.com/nichecode/pipeline-analyzer/internal/gotask"
)

// AnalysisResult represents the result of analyzing a build tool
type AnalysisResult struct {
	Tool      BuildTool `json:"tool"`
	Success   bool      `json:"success"`
	Error     string    `json:"error,omitempty"`
	OutputDir string    `json:"output_dir"`
}

// Analyzer coordinates analysis of all discovered build tools
type Analyzer struct {
	repository   *Repository
	discoveryDir string
}

// NewAnalyzer creates a new analyzer
func NewAnalyzer(repo *Repository, discoveryDir string) *Analyzer {
	return &Analyzer{
		repository:   repo,
		discoveryDir: discoveryDir,
	}
}

// AnalyzeAll analyzes all discovered build tools
func (a *Analyzer) AnalyzeAll() ([]AnalysisResult, error) {
	var results []AnalysisResult

	fmt.Printf("🔍 Discovered %d build tools in repository\n\n", len(a.repository.BuildTools))

	for _, tool := range a.repository.BuildTools {
		fmt.Printf("--- Analyzing %s ---\n", tool.Name)
		
		result := a.analyzeTool(tool)
		results = append(results, result)

		if result.Success {
			fmt.Printf("✅ %s analysis completed successfully\n", tool.Name)
			fmt.Printf("📁 Output: %s\n\n", result.OutputDir)
		} else {
			fmt.Printf("❌ %s analysis failed: %s\n\n", tool.Name, result.Error)
		}
	}

	return results, nil
}

// analyzeTool analyzes a specific build tool
func (a *Analyzer) analyzeTool(tool BuildTool) AnalysisResult {
	result := AnalysisResult{
		Tool: tool,
	}

	// Get the full config path
	configPath := filepath.Join(a.repository.RootPath, tool.ConfigPath)
	
	// Verify config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		result.Error = fmt.Sprintf("config file not found: %s", configPath)
		return result
	}

	// Determine output directory
	outputDir := filepath.Join(a.discoveryDir, tool.Type)
	result.OutputDir = outputDir

	// Analyze based on tool type
	switch tool.Type {
	case "circleci":
		err := a.analyzeCircleCI(configPath, outputDir)
		if err != nil {
			result.Error = err.Error()
		} else {
			result.Success = true
		}

	case "gotask":
		err := a.analyzeGoTask(configPath, outputDir)
		if err != nil {
			result.Error = err.Error()
		} else {
			result.Success = true
		}

	case "github-actions":
		err := a.analyzeGitHubActions(configPath, outputDir)
		if err != nil {
			result.Error = err.Error()
		} else {
			result.Success = true
		}

	default:
		// For now, we'll create a basic analysis for unsupported tools
		err := a.analyzeGeneric(tool, outputDir)
		if err != nil {
			result.Error = err.Error()
		} else {
			result.Success = true
		}
	}

	return result
}

// analyzeCircleCI runs CircleCI analysis
func (a *Analyzer) analyzeCircleCI(configPath, outputDir string) error {
	// Parse the configuration
	config, err := circleci.ParseConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to parse CircleCI config: %w", err)
	}

	// Validate the configuration
	if err := circleci.IsValidConfig(config); err != nil {
		return fmt.Errorf("invalid CircleCI configuration: %w", err)
	}

	fmt.Printf("✅ Configuration parsed successfully\n")
	fmt.Printf("   - Version: %s\n", config.Version)
	fmt.Printf("   - Jobs: %d\n", len(config.Jobs))
	fmt.Printf("   - Workflows: %d\n", len(config.Workflows))
	fmt.Printf("   - Executors: %d\n", len(config.Executors))

	// Validate output directory
	if err := circleci.ValidateOutputDir(outputDir); err != nil {
		return fmt.Errorf("output directory validation failed: %w", err)
	}

	// Perform analysis
	analysis := circleci.AnalyzeConfig(config)

	// Create writer and generate all files
	writer := circleci.NewWriter(outputDir)
	if err := writer.WriteAllFiles(analysis, configPath); err != nil {
		return fmt.Errorf("failed to write analysis files: %w", err)
	}

	// Clean up empty directories
	writer.CleanupEmptyDirs()

	return nil
}

// analyzeGoTask runs go-task analysis
func (a *Analyzer) analyzeGoTask(configPath, outputDir string) error {
	// Parse the taskfile
	taskfile, err := gotask.ParseTaskfile(configPath)
	if err != nil {
		return fmt.Errorf("failed to parse taskfile: %w", err)
	}

	// Validate the taskfile
	if err := gotask.IsValidTaskfile(taskfile); err != nil {
		return fmt.Errorf("invalid taskfile: %w", err)
	}

	fmt.Printf("✅ Taskfile parsed successfully\n")
	fmt.Printf("   - Version: %s\n", taskfile.Version)
	fmt.Printf("   - Tasks: %d\n", len(taskfile.Tasks))
	fmt.Printf("   - Includes: %d\n", len(taskfile.Includes))
	fmt.Printf("   - Global Variables: %d\n", len(taskfile.Vars))

	// Validate output directory
	if err := gotask.ValidateOutputDir(outputDir); err != nil {
		return fmt.Errorf("output directory validation failed: %w", err)
	}

	// Perform analysis
	analysis := gotask.AnalyzeTaskfile(taskfile)

	// Create writer and generate all files
	writer := gotask.NewWriter(outputDir)
	if err := writer.WriteAllFiles(analysis, configPath); err != nil {
		return fmt.Errorf("failed to write analysis files: %w", err)
	}

	return nil
}

// analyzeGitHubActions runs GitHub Actions workflow analysis
func (a *Analyzer) analyzeGitHubActions(configPath, outputDir string) error {
	// For GitHub Actions, configPath might be a directory path (.github/workflows/)
	// We need to handle both single file and directory scanning
	
	analyzer := githubactions.NewAnalyzer()
	
	// Check if configPath is a directory or file
	stat, err := os.Stat(configPath)
	if err != nil {
		return fmt.Errorf("failed to stat config path: %w", err)
	}

	var workflowFiles []string
	if stat.IsDir() {
		// Scan directory for workflow files
		files, err := filepath.Glob(filepath.Join(configPath, "*.yml"))
		if err != nil {
			return fmt.Errorf("failed to find yml files: %w", err)
		}
		yamlFiles, err := filepath.Glob(filepath.Join(configPath, "*.yaml"))
		if err != nil {
			return fmt.Errorf("failed to find yaml files: %w", err)
		}
		workflowFiles = append(files, yamlFiles...)
	} else {
		// Single file
		workflowFiles = []string{configPath}
	}

	if len(workflowFiles) == 0 {
		return fmt.Errorf("no workflow files found")
	}

	fmt.Printf("✅ Found %d workflow file(s)\n", len(workflowFiles))

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Analyze each workflow file
	var allResults []*githubactions.AnalysisResult
	for _, workflowFile := range workflowFiles {
		result, err := analyzer.AnalyzeWorkflow(workflowFile)
		if err != nil {
			fmt.Printf("⚠️  Failed to analyze %s: %v\n", filepath.Base(workflowFile), err)
			continue
		}
		allResults = append(allResults, result)
		
		fmt.Printf("   - %s: %d jobs, %d steps\n", 
			filepath.Base(workflowFile), len(result.Jobs), result.TotalSteps)
	}

	if len(allResults) == 0 {
		return fmt.Errorf("failed to analyze any workflow files")
	}

	// Generate markdown documentation for all workflows
	writer := githubactions.NewWriter(outputDir)
	if err := writer.WriteAllFiles(allResults, configPath); err != nil {
		return fmt.Errorf("failed to write analysis files: %w", err)
	}

	return nil
}

// analyzeGeneric creates a basic analysis for unsupported tools
func (a *Analyzer) analyzeGeneric(tool BuildTool, outputDir string) error {
	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create a basic README
	readme := fmt.Sprintf(`# %s Analysis

**Tool Type:** %s  
**Config File:** %s  
**Description:** %s  
**Generated:** %s

## Status

This build tool was discovered but detailed analysis is not yet supported.

## Configuration File

The configuration file is located at: %s

## Next Steps

- Manual review of the configuration file
- Consider adding support for this tool type to the analyzer
- Check tool-specific documentation for optimization opportunities

## Navigation

- [← Back to Discovery Overview](../README.md)
`, tool.Name, tool.Type, tool.ConfigPath, tool.Description, time.Now().Format(time.RFC3339), tool.ConfigPath)

	readmePath := filepath.Join(outputDir, "README.md")
	return os.WriteFile(readmePath, []byte(readme), 0644)
}

// GenerateOverview generates an overview of all discovered tools
func (a *Analyzer) GenerateOverview(results []AnalysisResult) error {
	overviewPath := filepath.Join(a.discoveryDir, "README.md")
	
	content := fmt.Sprintf(`# Pipeline Analysis Discovery Report

**Repository:** %s  
**Generated:** %s  
**Git Repository:** %t

## 🔍 Discovered Build Tools

Found **%d** build tools in this repository:

| Tool | Type | Status | Analysis |
|------|------|--------|----------|
`, a.repository.RootPath, time.Now().Format(time.RFC3339), a.repository.GitRepo, len(results))

	for _, result := range results {
		status := "✅ Success"
		analysis := fmt.Sprintf("[View Analysis](%s/README.md)", result.Tool.Type)
		
		if !result.Success {
			status = "❌ Failed"
			analysis = result.Error
		}

		content += fmt.Sprintf("| %s | %s | %s | %s |\n", 
			result.Tool.Name, result.Tool.Type, status, analysis)
	}

	content += `

## 📁 Directory Structure

Each discovered tool has its own analysis directory:

`

	for _, result := range results {
		if result.Success {
			content += fmt.Sprintf("- [%s/](%s/) - %s analysis\n", 
				result.Tool.Type, result.Tool.Type, result.Tool.Name)
		}
	}

	content += `

## 🚀 Getting Started

1. **Review the overview above** to see all discovered tools
2. **Navigate to specific analyses** using the links in the table
3. **Follow optimization recommendations** in each tool's analysis
4. **Cross-reference between tools** for comprehensive improvements

## 📊 Summary

This automated discovery found multiple build and deployment tools in your repository. Each tool has been analyzed for:

- Configuration validation
- Performance optimization opportunities  
- Security considerations
- Best practices recommendations
- Migration and modernization suggestions

Start with the tools marked as "Success" above for detailed insights.
`

	return os.WriteFile(overviewPath, []byte(content), 0644)
}

// GenerateHTMLIndex generates an HTML index for better navigation
func (a *Analyzer) GenerateHTMLIndex(results []AnalysisResult) error {
	indexPath := filepath.Join(a.discoveryDir, "index.html")
	
	content := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Pipeline Analysis Discovery Report</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, sans-serif; margin: 40px; line-height: 1.6; }
        .container { max-width: 800px; margin: 0 auto; }
        .tool-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); gap: 20px; margin: 20px 0; }
        .tool-card { border: 1px solid #ddd; border-radius: 8px; padding: 20px; }
        .tool-card h3 { margin-top: 0; color: #0066cc; }
        .status { padding: 4px 8px; border-radius: 4px; font-size: 0.9em; }
        .success { background: #d4edda; color: #155724; }
        .failed { background: #f8d7da; color: #721c24; }
        a { color: #0066cc; text-decoration: none; }
        a:hover { text-decoration: underline; }
        .section { margin: 30px 0; }
        .stats { background: #f8f9fa; padding: 20px; border-radius: 8px; margin: 20px 0; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🔍 Pipeline Analysis Discovery Report</h1>
        
        <div class="section">
            <p><strong>Repository:</strong> ` + a.repository.RootPath + `</p>
            <p><strong>Generated:</strong> ` + time.Now().Format("2006-01-02 15:04:05") + `</p>
            <p><strong>Git Repository:</strong> ` + fmt.Sprintf("%t", a.repository.GitRepo) + `</p>
        </div>

        <div class="stats">
            <h3>📊 Analysis Summary</h3>
            <p>Found <strong>` + fmt.Sprintf("%d", len(results)) + `</strong> build tools</p>
            <p>Successfully analyzed: <strong>` + fmt.Sprintf("%d", countSuccessful(results)) + `</strong></p>
        </div>

        <div class="section">
            <h2>🔍 Discovered Build Tools</h2>
            
            <div class="tool-grid">`

	for _, result := range results {
		status := "✅ Success"
		statusClass := "success"
		link := fmt.Sprintf(`<a href="%s/README.md">View Analysis →</a>`, result.Tool.Type)
		
		if !result.Success {
			status = "❌ Failed"
			statusClass = "failed"
			link = fmt.Sprintf("<p>Error: %s</p>", result.Error)
		}

		content += fmt.Sprintf(`
                <div class="tool-card">
                    <h3>%s</h3>
                    <p><strong>Type:</strong> %s</p>
                    <p><strong>Status:</strong> <span class="status %s">%s</span></p>
                    %s
                </div>`, 
			result.Tool.Name, result.Tool.Type, statusClass, status, link)
	}

	content += `
            </div>
        </div>

        <div class="section">
            <h2>📁 Directory Structure</h2>
            <ul>`

	for _, result := range results {
		if result.Success {
			content += fmt.Sprintf(`
                <li><a href="%s/">%s/</a> - %s analysis</li>`,
				result.Tool.Type, result.Tool.Type, result.Tool.Name)
		}
	}

	content += `
            </ul>
        </div>

        <div class="section">
            <h2>🚀 Getting Started</h2>
            <ol>
                <li><strong>Review the overview above</strong> to see all discovered tools</li>
                <li><strong>Navigate to specific analyses</strong> using the links in the cards</li>
                <li><strong>Follow optimization recommendations</strong> in each tool's analysis</li>
                <li><strong>Cross-reference between tools</strong> for comprehensive improvements</li>
            </ol>
        </div>

        <div class="section">
            <h2>📊 Summary</h2>
            <p>This automated discovery found multiple build and deployment tools in your repository. Each tool has been analyzed for:</p>
            <ul>
                <li>Configuration validation</li>
                <li>Performance optimization opportunities</li>
                <li>Security considerations</li>
                <li>Best practices recommendations</li>
                <li>Migration and modernization suggestions</li>
            </ul>
            <p>Start with the tools marked as "Success" above for detailed insights.</p>
        </div>

        <div class="section">
            <p><em>Generated by <strong>Pipeline Analyzer</strong> - Auto-Discovery Build Tool Analysis</em></p>
        </div>
    </div>
</body>
</html>`

	return os.WriteFile(indexPath, []byte(content), 0644)
}

// countSuccessful counts the number of successful analysis results
func countSuccessful(results []AnalysisResult) int {
	count := 0
	for _, result := range results {
		if result.Success {
			count++
		}
	}
	return count
}