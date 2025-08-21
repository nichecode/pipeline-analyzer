package circleci

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Writer handles file system operations for analysis output
type Writer struct {
	outputDir string
}

// NewWriter creates a new Writer with the specified output directory
func NewWriter(outputDir string) *Writer {
	return &Writer{outputDir: outputDir}
}

// CreateDirectories creates the required directory structure
func (w *Writer) CreateDirectories() error {
	dirs := []string{
		w.outputDir,
		filepath.Join(w.outputDir, "jobs"),
		filepath.Join(w.outputDir, "workflows"),
		filepath.Join(w.outputDir, "summaries"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// WriteAllFiles writes all analysis files to the output directory
func (w *Writer) WriteAllFiles(analysis *Analysis, configPath string) error {
	// Create directories
	if err := w.CreateDirectories(); err != nil {
		return err
	}

	// Write main files
	if err := w.writeMainFiles(analysis, configPath); err != nil {
		return err
	}

	// Write job files
	if err := w.writeJobFiles(analysis); err != nil {
		return err
	}

	// Write workflow files
	if err := w.writeWorkflowFiles(analysis); err != nil {
		return err
	}

	// Write summary files
	if err := w.writeSummaryFiles(analysis); err != nil {
		return err
	}

	return nil
}

// writeMainFiles writes README.md and migration-checklist.md
func (w *Writer) writeMainFiles(analysis *Analysis, configPath string) error {
	// Write README.md
	readme := GenerateMainReadme(analysis, configPath)
	if err := w.writeFile("README.md", readme); err != nil {
		return fmt.Errorf("failed to write README.md: %w", err)
	}

	// Write migration-checklist.md
	checklist := GenerateMigrationChecklist(analysis)
	if err := w.writeFile("migration-checklist.md", checklist); err != nil {
		return fmt.Errorf("failed to write migration-checklist.md: %w", err)
	}

	return nil
}

// writeJobFiles writes individual job analysis files
func (w *Writer) writeJobFiles(analysis *Analysis) error {
	jobNames := GetAllJobNames(analysis.Config)

	for _, jobName := range jobNames {
		jobAnalysis := AnalyzeJob(analysis.Config, jobName, analysis)
		if jobAnalysis == nil {
			continue
		}

		content := GenerateJobMarkdown(jobAnalysis)
		normalizedName := NormalizeJobName(jobName)
		filename := fmt.Sprintf("jobs/%s.md", normalizedName)

		if err := w.writeFile(filename, content); err != nil {
			return fmt.Errorf("failed to write job file %s: %w", filename, err)
		}
	}

	return nil
}

// writeWorkflowFiles writes individual workflow analysis files
func (w *Writer) writeWorkflowFiles(analysis *Analysis) error {
	workflowNames := GetAllWorkflowNames(analysis.Config)

	for _, workflowName := range workflowNames {
		workflowAnalysis := AnalyzeWorkflow(analysis.Config, workflowName)
		if workflowAnalysis == nil {
			continue
		}

		content := GenerateWorkflowMarkdown(workflowAnalysis)
		normalizedName := NormalizeJobName(workflowName)
		filename := fmt.Sprintf("workflows/%s.md", normalizedName)

		if err := w.writeFile(filename, content); err != nil {
			return fmt.Errorf("failed to write workflow file %s: %w", filename, err)
		}
	}

	return nil
}

// writeSummaryFiles writes all summary analysis files
func (w *Writer) writeSummaryFiles(analysis *Analysis) error {
	summaries := map[string]string{
		"all-jobs.md":             GenerateAllJobsIndex(analysis),
		"job-usage.md":            GenerateJobUsageAnalysis(analysis),
		"commands.md":             GenerateCommandsAnalysis(analysis),
		"docker-and-scripts.md":   GenerateDockerAndScriptsAnalysis(analysis),
		"executors-and-images.md": GenerateExecutorsAndImagesAnalysis(analysis),
		"workflows.md":            GenerateWorkflowsIndex(analysis),
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

// DetermineOutputDir determines the appropriate output directory
func DetermineOutputDir(customDir string) (string, error) {
	if customDir != "" {
		return customDir, nil
	}

	// Check if we're in a git repository
	if isGitRepo() {
		return ".discovery/pipeline-analyzer/circleci", nil
	}

	return "circleci-analysis", nil
}

// isGitRepo checks if the current directory is a git repository
func isGitRepo() bool {
	_, err := os.Stat(".git")
	return err == nil
}

// PrintSummary prints a summary of the analysis results
func PrintSummary(analysis *Analysis, outputDir string) {
	fmt.Printf("\n✅ Analysis complete!\n\n")
	fmt.Printf("📁 **Output directory:** %s\n", outputDir)
	fmt.Printf("📖 **Start here:** %s/README.md\n\n", outputDir)

	fmt.Printf("🔗 **Key entry points:**\n")
	fmt.Printf("   - 📋 Migration guide: %s/migration-checklist.md\n", outputDir)
	fmt.Printf("   - 📈 Job analysis: %s/summaries/job-usage.md\n", outputDir)
	fmt.Printf("   - 📝 All jobs: %s/summaries/all-jobs.md\n\n", outputDir)

	fmt.Printf("📊 **Analysis summary:**\n")
	fmt.Printf("   - Jobs analyzed: %d\n", analysis.TotalJobs)
	fmt.Printf("   - Workflows found: %d\n", analysis.TotalWorkflows)

	dockerJobs, otherJobs := CountDockerUsage(analysis.Config)
	fmt.Printf("   - Docker jobs: %d\n", dockerJobs)
	fmt.Printf("   - Other executors: %d\n\n", otherJobs)

	fmt.Printf("🚀 **Ready to start your go-task migration!**\n")
	fmt.Printf("   Open %s/README.md in your browser or markdown viewer.\n", outputDir)
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

// CleanupEmptyDirs removes empty directories from the output
func (w *Writer) CleanupEmptyDirs() error {
	dirs := []string{
		filepath.Join(w.outputDir, "jobs"),
		filepath.Join(w.outputDir, "workflows"),
		filepath.Join(w.outputDir, "summaries"),
	}

	for _, dir := range dirs {
		if isEmpty, err := isDirEmpty(dir); err == nil && isEmpty {
			os.Remove(dir)
		}
	}

	return nil
}

// isDirEmpty checks if a directory is empty
func isDirEmpty(dirname string) (bool, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdir(1)
	if err == nil {
		return false, nil // Directory has at least one file
	}

	return true, nil // Directory is empty
}

// GetRelativePath converts absolute path to relative if possible
func GetRelativePath(path string) string {
	if wd, err := os.Getwd(); err == nil {
		if rel, err := filepath.Rel(wd, path); err == nil && !strings.HasPrefix(rel, "..") {
			return rel
		}
	}
	return path
}
