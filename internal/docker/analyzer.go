package docker

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// AnalyzeDocker performs comprehensive analysis of Docker configurations
func AnalyzeDocker(rootPath string) (*DockerAnalysis, error) {
	dockerfiles, composeFiles, err := DiscoverDockerFiles(rootPath)
	if err != nil {
		return nil, fmt.Errorf("failed to discover Docker files: %w", err)
	}

	analysis := &DockerAnalysis{
		Dockerfiles:   []*DockerfileAnalysis{},
		Summary:       &DockerSummary{},
		GeneratedAt:   time.Now(),
	}

	// Analyze Dockerfiles
	for _, dockerfilePath := range dockerfiles {
		dockerfileAnalysis, err := ParseDockerfile(dockerfilePath)
		if err != nil {
			// Log error but continue with other files
			fmt.Printf("Warning: Failed to parse Dockerfile %s: %v\n", dockerfilePath, err)
			continue
		}
		analysis.Dockerfiles = append(analysis.Dockerfiles, dockerfileAnalysis)
	}

	// Analyze docker-compose files
	if len(composeFiles) > 0 {
		// Use the first compose file found (most common case)
		composeAnalysis, err := ParseDockerCompose(composeFiles[0])
		if err != nil {
			fmt.Printf("Warning: Failed to parse docker-compose %s: %v\n", composeFiles[0], err)
		} else {
			analysis.DockerCompose = composeAnalysis
		}
	}

	// Analyze Docker usage across other tools
	analysis.Usage = analyzeDockerUsage(rootPath, analysis)

	// Generate summary
	analysis.Summary = generateDockerSummary(analysis)

	return analysis, nil
}

// generateDockerSummary creates a high-level summary of the Docker analysis
func generateDockerSummary(analysis *DockerAnalysis) *DockerSummary {
	summary := &DockerSummary{
		TotalDockerfiles:       len(analysis.Dockerfiles),
		HasDockerCompose:       analysis.DockerCompose != nil,
		UniqueBaseImages:       []string{},
		MostCommonInstructions: []string{},
		Recommendations:        []string{},
	}

	if analysis.DockerCompose != nil {
		summary.ServiceCount = analysis.DockerCompose.ServiceCount
	}

	// Count multi-stage builds and collect base images
	baseImageMap := make(map[string]int)
	instructionMap := make(map[string]int)
	totalSecurityIssues := 0
	totalOptimizationIssues := 0

	for _, dockerfile := range analysis.Dockerfiles {
		if dockerfile.MultiStage {
			summary.MultiStageBuilds++
		}

		// Collect base images
		for _, baseImage := range dockerfile.BaseImages {
			baseImageMap[baseImage]++
		}

		// Count instructions
		for _, stage := range dockerfile.Stages {
			for _, instruction := range stage.Instructions {
				instructionMap[instruction.Instruction]++
			}
		}

		// Count security issues
		if dockerfile.SecurityScan != nil {
			if dockerfile.SecurityScan.RunAsRoot {
				totalSecurityIssues++
			}
			if dockerfile.SecurityScan.UsesLatestTags {
				totalSecurityIssues++
			}
			if !dockerfile.SecurityScan.HasHealthCheck {
				totalSecurityIssues++
			}

			totalSecurityIssues += len(dockerfile.SecurityScan.SecurityRecommendations)
			totalOptimizationIssues += len(dockerfile.SecurityScan.BestPracticeIssues)
		}
	}

	// Extract unique base images
	for baseImage := range baseImageMap {
		summary.UniqueBaseImages = append(summary.UniqueBaseImages, baseImage)
	}

	// Find most common instructions
	summary.MostCommonInstructions = findTopInstructions(instructionMap, 5)

	summary.SecurityIssues = totalSecurityIssues
	summary.OptimizationIssues = totalOptimizationIssues

	// Generate high-level recommendations
	generateSummaryRecommendations(summary, analysis)

	// Calculate overall score
	summary.OverallScore = calculateOverallScore(summary)

	return summary
}

// findTopInstructions finds the most commonly used Docker instructions
func findTopInstructions(instructionMap map[string]int, limit int) []string {
	type instructionCount struct {
		instruction string
		count       int
	}

	var counts []instructionCount
	for instruction, count := range instructionMap {
		counts = append(counts, instructionCount{instruction, count})
	}

	// Simple bubble sort (fine for small datasets)
	for i := 0; i < len(counts); i++ {
		for j := i + 1; j < len(counts); j++ {
			if counts[i].count < counts[j].count {
				counts[i], counts[j] = counts[j], counts[i]
			}
		}
	}

	var result []string
	maxLimit := limit
	if len(counts) < maxLimit {
		maxLimit = len(counts)
	}

	for i := 0; i < maxLimit; i++ {
		result = append(result, counts[i].instruction)
	}

	return result
}

// generateSummaryRecommendations generates high-level recommendations
func generateSummaryRecommendations(summary *DockerSummary, analysis *DockerAnalysis) {
	if summary.SecurityIssues > 0 {
		summary.Recommendations = append(summary.Recommendations,
			"Address security issues found in Dockerfiles")
	}

	if summary.OptimizationIssues > 0 {
		summary.Recommendations = append(summary.Recommendations,
			"Optimize Dockerfiles for better performance and caching")
	}

	if summary.MultiStageBuilds == 0 && summary.TotalDockerfiles > 0 {
		summary.Recommendations = append(summary.Recommendations,
			"Consider using multi-stage builds to reduce image sizes")
	}

	if !summary.HasDockerCompose && summary.TotalDockerfiles > 1 {
		summary.Recommendations = append(summary.Recommendations,
			"Consider adding docker-compose.yml for easier multi-container management")
	}

	// Check for latest tags in base images
	hasLatestTags := false
	for _, baseImage := range summary.UniqueBaseImages {
		if !containsVersion(baseImage) {
			hasLatestTags = true
			break
		}
	}

	if hasLatestTags {
		summary.Recommendations = append(summary.Recommendations,
			"Pin base images to specific versions for better reproducibility")
	}
}

// containsVersion checks if an image name contains a version tag
func containsVersion(image string) bool {
	return len(image) > 0 && (image[len(image)-1:] != ":" || image != "latest")
}

// calculateOverallScore calculates an overall score for the Docker setup
func calculateOverallScore(summary *DockerSummary) int {
	score := 100

	// Deduct points for security issues
	score -= summary.SecurityIssues * 10

	// Deduct points for optimization issues
	score -= summary.OptimizationIssues * 5

	// Add points for multi-stage builds
	if summary.MultiStageBuilds > 0 {
		score += 10
	}

	// Add points for having docker-compose
	if summary.HasDockerCompose {
		score += 15
	}

	// Ensure score is within bounds
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score
}

// ValidateOutputDir ensures the output directory exists and is writable
func ValidateOutputDir(outputDir string) error {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Test write permissions
	testFile := filepath.Join(outputDir, ".write_test")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		return fmt.Errorf("output directory is not writable: %w", err)
	}

	// Clean up test file
	os.Remove(testFile)

	return nil
}

// analyzeDockerUsage analyzes how Docker files are used across other build tools
func analyzeDockerUsage(rootPath string, analysis *DockerAnalysis) *DockerUsageAnalysis {
	usage := &DockerUsageAnalysis{
		DockerfileReferences:    []*DockerUsageReference{},
		DockerComposeReferences: []*DockerUsageReference{},
		DockerCommandReferences: []*DockerUsageReference{},
		TotalReferences:         0,
	}

	// Search for Docker usage in various tools
	searchCircleCIUsage(rootPath, analysis, usage)
	searchGitHubActionsUsage(rootPath, analysis, usage)
	searchGoTaskUsage(rootPath, analysis, usage)

	// Update total references
	usage.TotalReferences = len(usage.DockerfileReferences) + 
		len(usage.DockerComposeReferences) + 
		len(usage.DockerCommandReferences)

	return usage
}

// searchCircleCIUsage searches for Docker references in CircleCI configuration
func searchCircleCIUsage(rootPath string, analysis *DockerAnalysis, usage *DockerUsageAnalysis) {
	circleCIConfig := filepath.Join(rootPath, ".circleci", "config.yml")
	if content, err := os.ReadFile(circleCIConfig); err == nil {
		text := string(content)
		
		// Search for docker build commands
		dockerBuildRegex := regexp.MustCompile(`docker\s+build\s+[^\n]*`)
		matches := dockerBuildRegex.FindAllString(text, -1)
		for _, match := range matches {
			usage.DockerCommandReferences = append(usage.DockerCommandReferences, &DockerUsageReference{
				Tool:     "circleci",
				File:     ".circleci/config.yml",
				Location: "job",
				Command:  strings.TrimSpace(match),
				Context:  "Docker build command",
			})
		}
		
		// Search for docker-compose commands
		composeRegex := regexp.MustCompile(`docker-compose\s+[^\n]*`)
		matches = composeRegex.FindAllString(text, -1)
		for _, match := range matches {
			usage.DockerComposeReferences = append(usage.DockerComposeReferences, &DockerUsageReference{
				Tool:     "circleci",
				File:     ".circleci/config.yml",
				Location: "job",
				Command:  strings.TrimSpace(match),
				Context:  "Docker Compose command",
			})
		}
		
		// Search for specific Dockerfile references
		for _, dockerfile := range analysis.Dockerfiles {
			dockerfileRelPath, _ := filepath.Rel(rootPath, dockerfile.FilePath)
			if strings.Contains(text, dockerfileRelPath) {
				usage.DockerfileReferences = append(usage.DockerfileReferences, &DockerUsageReference{
					Tool:     "circleci",
					File:     ".circleci/config.yml",
					Location: "job",
					Command:  dockerfileRelPath,
					Context:  "Dockerfile reference",
				})
			}
		}
	}
}

// searchGitHubActionsUsage searches for Docker references in GitHub Actions workflows
func searchGitHubActionsUsage(rootPath string, analysis *DockerAnalysis, usage *DockerUsageAnalysis) {
	workflowsDir := filepath.Join(rootPath, ".github", "workflows")
	if files, err := filepath.Glob(filepath.Join(workflowsDir, "*.yml")); err == nil {
		for _, file := range files {
			if content, err := os.ReadFile(file); err == nil {
				text := string(content)
				relFile, _ := filepath.Rel(rootPath, file)
				
				// Search for docker build commands
				dockerBuildRegex := regexp.MustCompile(`docker\s+build\s+[^\n]*`)
				matches := dockerBuildRegex.FindAllString(text, -1)
				for _, match := range matches {
					usage.DockerCommandReferences = append(usage.DockerCommandReferences, &DockerUsageReference{
						Tool:     "github-actions",
						File:     relFile,
						Location: "step",
						Command:  strings.TrimSpace(match),
						Context:  "Docker build command",
					})
				}
				
				// Search for docker-compose commands
				composeRegex := regexp.MustCompile(`docker-compose\s+[^\n]*`)
				matches = composeRegex.FindAllString(text, -1)
				for _, match := range matches {
					usage.DockerComposeReferences = append(usage.DockerComposeReferences, &DockerUsageReference{
						Tool:     "github-actions",
						File:     relFile,
						Location: "step",
						Command:  strings.TrimSpace(match),
						Context:  "Docker Compose command",
					})
				}
				
				// Search for specific Dockerfile references
				for _, dockerfile := range analysis.Dockerfiles {
					dockerfileRelPath, _ := filepath.Rel(rootPath, dockerfile.FilePath)
					if strings.Contains(text, dockerfileRelPath) {
						usage.DockerfileReferences = append(usage.DockerfileReferences, &DockerUsageReference{
							Tool:     "github-actions",
							File:     relFile,
							Location: "step",
							Command:  dockerfileRelPath,
							Context:  "Dockerfile reference",
						})
					}
				}
			}
		}
	}
	
	// Also check for .yaml files
	if files, err := filepath.Glob(filepath.Join(workflowsDir, "*.yaml")); err == nil {
		for _, file := range files {
			if content, err := os.ReadFile(file); err == nil {
				text := string(content)
				relFile, _ := filepath.Rel(rootPath, file)
				
				// Search for docker build commands
				dockerBuildRegex := regexp.MustCompile(`docker\s+build\s+[^\n]*`)
				matches := dockerBuildRegex.FindAllString(text, -1)
				for _, match := range matches {
					usage.DockerCommandReferences = append(usage.DockerCommandReferences, &DockerUsageReference{
						Tool:     "github-actions",
						File:     relFile,
						Location: "step",
						Command:  strings.TrimSpace(match),
						Context:  "Docker build command",
					})
				}
			}
		}
	}
}

// searchGoTaskUsage searches for Docker references in Taskfile configurations  
func searchGoTaskUsage(rootPath string, analysis *DockerAnalysis, usage *DockerUsageAnalysis) {
	taskfiles := []string{
		filepath.Join(rootPath, "Taskfile.yml"),
		filepath.Join(rootPath, "Taskfile.yaml"),
		filepath.Join(rootPath, "taskfile.yml"), 
		filepath.Join(rootPath, "taskfile.yaml"),
	}
	
	for _, taskfile := range taskfiles {
		if content, err := os.ReadFile(taskfile); err == nil {
			text := string(content)
			relFile, _ := filepath.Rel(rootPath, taskfile)
			
			// Search for docker build commands
			dockerBuildRegex := regexp.MustCompile(`docker\s+build\s+[^\n]*`)
			matches := dockerBuildRegex.FindAllString(text, -1)
			for _, match := range matches {
				usage.DockerCommandReferences = append(usage.DockerCommandReferences, &DockerUsageReference{
					Tool:     "gotask",
					File:     relFile,
					Location: "task",
					Command:  strings.TrimSpace(match),
					Context:  "Docker build command",
				})
			}
			
			// Search for docker-compose commands
			composeRegex := regexp.MustCompile(`docker-compose\s+[^\n]*`)
			matches = composeRegex.FindAllString(text, -1)
			for _, match := range matches {
				usage.DockerComposeReferences = append(usage.DockerComposeReferences, &DockerUsageReference{
					Tool:     "gotask",
					File:     relFile,
					Location: "task",
					Command:  strings.TrimSpace(match),
					Context:  "Docker Compose command",
				})
			}
			
			// Search for specific Dockerfile references
			for _, dockerfile := range analysis.Dockerfiles {
				dockerfileRelPath, _ := filepath.Rel(rootPath, dockerfile.FilePath)
				if strings.Contains(text, dockerfileRelPath) {
					usage.DockerfileReferences = append(usage.DockerfileReferences, &DockerUsageReference{
						Tool:     "gotask",
						File:     relFile,
						Location: "task",
						Command:  dockerfileRelPath,
						Context:  "Dockerfile reference",
					})
				}
			}
		}
	}
}