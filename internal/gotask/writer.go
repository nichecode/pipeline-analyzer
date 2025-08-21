package gotask

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/nichecode/pipeline-analyzer/internal/shared"
)

// Writer handles file system operations for go-task analysis output
type Writer struct {
	outputDir string
}

// NewWriter creates a new Writer with the specified output directory
func NewWriter(outputDir string) *Writer {
	return &Writer{outputDir: outputDir}
}

// CreateDirectories creates the required directory structure for go-task analysis
func (w *Writer) CreateDirectories() error {
	dirs := []string{
		w.outputDir,
		filepath.Join(w.outputDir, "tasks"),
		filepath.Join(w.outputDir, "summaries"),
	}

	// Add includes directory if needed
	dirs = append(dirs, filepath.Join(w.outputDir, "includes"))

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// WriteAllFiles writes all go-task analysis files to the output directory
func (w *Writer) WriteAllFiles(analysis *Analysis, taskfilePath string) error {
	// Create directories
	if err := w.CreateDirectories(); err != nil {
		return err
	}

	// Write main files
	if err := w.writeMainFiles(analysis, taskfilePath); err != nil {
		return err
	}

	// Write task files
	if err := w.writeTaskFiles(analysis); err != nil {
		return err
	}

	// Write include files
	if err := w.writeIncludeFiles(analysis); err != nil {
		return err
	}

	// Write summary files
	if err := w.writeSummaryFiles(analysis); err != nil {
		return err
	}

	return nil
}

// writeMainFiles writes README.md and optimization-guide.md
func (w *Writer) writeMainFiles(analysis *Analysis, taskfilePath string) error {
	// Write README.md
	readme := GenerateMainReadme(analysis, taskfilePath)
	if err := w.writeFile("README.md", readme); err != nil {
		return fmt.Errorf("failed to write README.md: %w", err)
	}

	// Write optimization-guide.md
	guide := GenerateOptimizationGuide(analysis)
	if err := w.writeFile("optimization-guide.md", guide); err != nil {
		return fmt.Errorf("failed to write optimization-guide.md: %w", err)
	}

	return nil
}

// writeTaskFiles writes individual task analysis files
func (w *Writer) writeTaskFiles(analysis *Analysis) error {
	taskNames := GetAllTaskNames(analysis.Taskfile)

	// Write dependency graph first
	dependencyGraph := GenerateDependencyGraph(analysis)
	if err := w.writeFile("tasks/dependency-graph.md", dependencyGraph); err != nil {
		return fmt.Errorf("failed to write dependency graph: %w", err)
	}

	// Write individual task files
	for _, taskName := range taskNames {
		taskAnalysis := AnalyzeTask(analysis.Taskfile, taskName, analysis)
		if taskAnalysis == nil {
			continue
		}

		content := GenerateTaskMarkdown(taskAnalysis)
		normalizedName := NormalizeTaskName(taskName)
		filename := fmt.Sprintf("tasks/%s.md", normalizedName)

		if err := w.writeFile(filename, content); err != nil {
			return fmt.Errorf("failed to write task file %s: %w", filename, err)
		}
	}

	return nil
}

// writeIncludeFiles writes individual include analysis files
func (w *Writer) writeIncludeFiles(analysis *Analysis) error {
	if len(analysis.IncludeAnalysis) == 0 {
		return nil
	}

	for includeName, includeAnalysis := range analysis.IncludeAnalysis {
		content := generateIncludeMarkdown(includeName, includeAnalysis)
		normalizedName := NormalizeTaskName(includeName)
		filename := fmt.Sprintf("includes/%s.md", normalizedName)

		if err := w.writeFile(filename, content); err != nil {
			return fmt.Errorf("failed to write include file %s: %w", filename, err)
		}
	}

	return nil
}

// writeSummaryFiles writes all summary analysis files
func (w *Writer) writeSummaryFiles(analysis *Analysis) error {
	summaries := map[string]string{
		"all-tasks.md":   GenerateAllTasksIndex(analysis),
		"task-usage.md":  GenerateTaskUsageAnalysis(analysis),
		"commands.md":    GenerateCommandsAnalysis(analysis),
		"performance.md": GeneratePerformanceAnalysis(analysis),
		"variables.md":   GenerateVariableAnalysis(analysis),
	}

	// Add includes summary if there are includes
	if len(analysis.IncludeAnalysis) > 0 {
		summaries["includes.md"] = generateIncludesSummary(analysis)
	}

	for filename, content := range summaries {
		summaryPath := filepath.Join("summaries", filename)
		if err := w.writeFile(summaryPath, content); err != nil {
			return fmt.Errorf("failed to write summary file %s: %w", summaryPath, err)
		}
	}

	return nil
}

// writeFile writes content to a file in the output directory
func (w *Writer) writeFile(filename, content string) error {
	fullPath := filepath.Join(w.outputDir, filename)

	// Create parent directory if it doesn't exist
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", fullPath, err)
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		return fmt.Errorf("failed to write content to %s: %w", fullPath, err)
	}

	return nil
}

// generateIncludeMarkdown generates markdown for a single include
func generateIncludeMarkdown(includeName string, includeAnalysis *IncludeAnalysis) string {
	content := fmt.Sprintf("# Include: %s\n\n", includeName)
	content += fmt.Sprintf("**Path:** %s\n", includeAnalysis.Path)
	content += fmt.Sprintf("**Namespace:** %s\n", includeAnalysis.Namespace)
	content += fmt.Sprintf("**Task Count:** %d\n\n", includeAnalysis.TaskCount)

	if len(includeAnalysis.Dependencies) > 0 {
		content += "## Dependencies\n\n"
		for _, dep := range includeAnalysis.Dependencies {
			content += fmt.Sprintf("- %s\n", dep)
		}
		content += "\n"
	}

	// Show individual tasks from the included file
	if len(includeAnalysis.Tasks) > 0 {
		content += "## 📋 Included Tasks\n\n"
		
		// Create a sorted list of task names
		var taskNames []string
		for taskName := range includeAnalysis.Tasks {
			taskNames = append(taskNames, taskName)
		}
		sort.Strings(taskNames)
		
		content += "| Task | Description | Commands | Type |\n"
		content += "|------|-------------|----------|------|\n"
		
		for _, taskName := range taskNames {
			taskAnalysis := includeAnalysis.Tasks[taskName]
			
			description := taskAnalysis.Description
			if description == "" {
				description = "*No description*"
			}
			
			commandCount := len(taskAnalysis.Commands)
			commandsText := fmt.Sprintf("%d commands", commandCount)
			if commandCount == 0 {
				if len(taskAnalysis.Dependencies) > 0 {
					commandsText = "Dependency-only"
				} else {
					commandsText = "No commands"
				}
			}
			
			content += fmt.Sprintf("| **%s** | %s | %s | %s |\n", 
				taskName, description, commandsText, taskAnalysis.Type)
		}
		content += "\n"
		
		// Show detailed commands for each task
		content += "## ⚡ Task Commands\n\n"
		
		for _, taskName := range taskNames {
			taskAnalysis := includeAnalysis.Tasks[taskName]
			
			content += fmt.Sprintf("### %s\n\n", taskName)
			
			if taskAnalysis.Description != "" {
				content += fmt.Sprintf("**Description:** %s\n\n", taskAnalysis.Description)
			}
			
			// Show commands or explain why there are none
			if len(taskAnalysis.Commands) > 0 {
				content += "**Commands:**\n\n"
				for i, command := range taskAnalysis.Commands {
					content += fmt.Sprintf("%d. ```bash\n%s\n```\n\n", i+1, command)
				}
			} else if len(taskAnalysis.Dependencies) > 0 {
				content += "**Type:** Dependency-only task\n\n"
				content += "**Dependencies:**\n"
				for _, dep := range taskAnalysis.Dependencies {
					content += fmt.Sprintf("- %s\n", dep)
				}
				content += "\n"
			} else {
				content += "**Status:** ⚠️ No commands or dependencies defined\n\n"
			}
			
			// Show additional task properties if available
			if len(taskAnalysis.Sources) > 0 {
				content += "**Sources:** " + fmt.Sprintf("%v", taskAnalysis.Sources) + "\n\n"
			}
			if len(taskAnalysis.Generates) > 0 {
				content += "**Generates:** " + fmt.Sprintf("%v", taskAnalysis.Generates) + "\n\n"
			}
		}
	} else if includeAnalysis.TaskCount > 0 {
		content += "## ⚠️ Tasks Not Analyzed\n\n"
		content += fmt.Sprintf("This include contains %d tasks, but they could not be analyzed. ", includeAnalysis.TaskCount)
		content += "This may happen if the included taskfile is not accessible or has parsing errors.\n\n"
	}

	content += "## Navigation\n\n"
	content += "- [← Back to Overview](../README.md)\n"
	content += "- [Include Summary](../summaries/includes.md)\n"

	return content
}

// generateIncludesSummary generates the includes summary
func generateIncludesSummary(analysis *Analysis) string {
	content := "# Include Analysis\n\n"
	content += fmt.Sprintf("Total includes found: **%d**\n\n", analysis.TotalIncludes)

	if len(analysis.IncludeAnalysis) > 0 {
		content += "| Include | Path | Tasks | Namespace |\n"
		content += "|---------|------|-------|--------|\n"

		var includeNames []string
		for name := range analysis.IncludeAnalysis {
			includeNames = append(includeNames, name)
		}
		sort.Strings(includeNames)

		for _, name := range includeNames {
			include := analysis.IncludeAnalysis[name]
			normalizedName := NormalizeTaskName(name)
			includeLink := fmt.Sprintf("[%s](../includes/%s.md)", name, normalizedName)

			content += fmt.Sprintf("| %s | %s | %d | %s |\n",
				includeLink, include.Path, include.TaskCount, include.Namespace)
		}
		content += "\n"
	}

	content += "## Navigation\n\n"
	content += "- [← Back to Overview](../README.md)\n"
	content += "- [All Tasks](all-tasks.md)\n"

	return content
}

// PrintSummary prints a summary of the go-task analysis results
func PrintSummary(analysis *Analysis, outputDir string) {
	fmt.Printf("\n✅ Go-Task analysis complete!\n\n")
	fmt.Printf("📁 **Output directory:** %s\n", outputDir)
	fmt.Printf("📖 **Start here:** %s/README.md\n\n", outputDir)

	fmt.Printf("🔗 **Key entry points:**\n")
	fmt.Printf("   - 📋 Optimization guide: %s/optimization-guide.md\n", outputDir)
	fmt.Printf("   - 📈 Task analysis: %s/summaries/task-usage.md\n", outputDir)
	fmt.Printf("   - ⚡ Performance metrics: %s/summaries/performance.md\n", outputDir)
	fmt.Printf("   - 🔗 Dependency graph: %s/tasks/dependency-graph.md\n\n", outputDir)

	fmt.Printf("📊 **Analysis summary:**\n")
	fmt.Printf("   - Tasks analyzed: %d\n", analysis.TotalTasks)
	fmt.Printf("   - Includes found: %d\n", analysis.TotalIncludes)

	metrics := GetPerformanceMetrics(analysis.Taskfile)
	fmt.Printf("   - Tasks with caching: %d (%.1f%%)\n", 
		metrics.TasksWithCaching, float64(metrics.TasksWithCaching)/float64(analysis.TotalTasks)*100)

	if len(analysis.CircularDeps) > 0 {
		fmt.Printf("   - Circular dependencies: %d ⚠️\n", len(analysis.CircularDeps))
	} else {
		fmt.Printf("   - Circular dependencies: 0 ✅\n")
	}

	fmt.Printf("   - Optimization tips: %d\n\n", len(analysis.OptimizationTips))

	// Show ecosystem
	allCommands := []string{}
	for _, task := range analysis.Taskfile.Tasks {
		allCommands = append(allCommands, ExtractTaskCommands(task)...)
	}
	ecosystem := shared.DetectToolEcosystem(allCommands)
	fmt.Printf("🛠️ **Primary tool ecosystem:** %s\n\n", ecosystem)

	fmt.Printf("🚀 **Ready to optimize your Taskfile!**\n")
	fmt.Printf("   Open %s/README.md in your browser or markdown viewer.\n", outputDir)
}

// DetermineOutputDir determines the appropriate output directory for go-task analysis
func DetermineOutputDir(customDir string, mode string) (string, error) {
	if customDir != "" {
		if mode == "gotask" {
			return filepath.Join(customDir, "gotask"), nil
		}
		return customDir, nil
	}

	// Check if we're in a git repository
	if isGitRepo() {
		return ".discovery/pipeline-analyzer/gotask", nil
	}

	return "gotask-analysis", nil
}

// isGitRepo checks if the current directory is a git repository
func isGitRepo() bool {
	_, err := os.Stat(".git")
	return err == nil
}

// ValidateOutputDir validates that the output directory can be created/written to
func ValidateOutputDir(outputDir string) error {
	// Try to create the directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("cannot create output directory %s: %w", outputDir, err)
	}

	// Try to write a test file
	testFile := filepath.Join(outputDir, ".test")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		return fmt.Errorf("cannot write to output directory %s: %w", outputDir, err)
	}

	// Clean up test file
	os.Remove(testFile)

	return nil
}