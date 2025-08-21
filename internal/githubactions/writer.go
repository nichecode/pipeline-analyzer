package githubactions

import (
	"fmt"
	"os"
	"path/filepath"
)

// Writer handles file writing for GitHub Actions analysis
type Writer struct {
	outputDir string
}

// NewWriter creates a new GitHub Actions writer
func NewWriter(outputDir string) *Writer {
	return &Writer{
		outputDir: outputDir,
	}
}

// WriteAllFiles writes all analysis files
func (w *Writer) WriteAllFiles(results []*AnalysisResult, configPath string) error {
	// Create necessary directories
	dirs := []string{
		filepath.Join(w.outputDir, "workflows"),
		filepath.Join(w.outputDir, "jobs"),
		filepath.Join(w.outputDir, "summaries"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Write main README
	if err := w.writeMainReadme(results, configPath); err != nil {
		return fmt.Errorf("failed to write main README: %w", err)
	}

	// Write individual workflow files
	for _, result := range results {
		if err := w.writeWorkflowFiles(result); err != nil {
			return fmt.Errorf("failed to write workflow files: %w", err)
		}
	}

	// Write summary files
	if err := w.writeSummaryFiles(results); err != nil {
		return fmt.Errorf("failed to write summary files: %w", err)
	}

	return nil
}

// writeMainReadme writes the main README file
func (w *Writer) writeMainReadme(results []*AnalysisResult, configPath string) error {
	generator := NewMarkdownGenerator()
	content := generator.GenerateMainReadme(results, configPath)
	
	readmePath := filepath.Join(w.outputDir, "README.md")
	return os.WriteFile(readmePath, []byte(content), 0644)
}

// writeWorkflowFiles writes individual workflow analysis files
func (w *Writer) writeWorkflowFiles(result *AnalysisResult) error {
	generator := NewMarkdownGenerator()
	
	// Write workflow overview
	workflowContent := generator.GenerateWorkflowAnalysis(result)
	workflowPath := filepath.Join(w.outputDir, "workflows", 
		fmt.Sprintf("%s.md", sanitizeFilename(result.Config.Name)))
	
	if err := os.WriteFile(workflowPath, []byte(workflowContent), 0644); err != nil {
		return err
	}

	// Write individual job files
	for _, job := range result.Jobs {
		jobContent := generator.GenerateJobAnalysis(job, result.Config.Name)
		jobPath := filepath.Join(w.outputDir, "jobs", 
			fmt.Sprintf("%s.md", sanitizeFilename(job.Name)))
		
		if err := os.WriteFile(jobPath, []byte(jobContent), 0644); err != nil {
			return err
		}
	}

	return nil
}

// writeSummaryFiles writes summary analysis files
func (w *Writer) writeSummaryFiles(results []*AnalysisResult) error {
	generator := NewMarkdownGenerator()
	
	summaries := map[string]string{
		"actions-usage.md":    generator.GenerateActionsUsage(results),
		"runners-analysis.md": generator.GenerateRunnersAnalysis(results),
		"commands-analysis.md": generator.GenerateCommandsAnalysis(results),
		"go-task-migration.md": generator.GenerateGoTaskMigration(results),
	}

	for filename, content := range summaries {
		path := filepath.Join(w.outputDir, "summaries", filename)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", filename, err)
		}
	}

	return nil
}

// sanitizeFilename sanitizes filename for filesystem
func sanitizeFilename(name string) string {
	if name == "" {
		return "untitled"
	}
	
	// Replace common problematic characters
	replacements := map[rune]string{
		'/': "-",
		'\\': "-", 
		':': "-",
		'*': "-",
		'?': "-",
		'"': "-",
		'<': "-",
		'>': "-",
		'|': "-",
		' ': "-",
	}
	
	result := ""
	for _, char := range name {
		if replacement, exists := replacements[char]; exists {
			result += replacement
		} else {
			result += string(char)
		}
	}
	
	return result
}