package githubactions

import (
	"fmt"
	"strings"
	"time"

	"github.com/nichecode/pipeline-analyzer/internal/shared"
)

// Analyzer provides GitHub Actions workflow analysis
type Analyzer struct {
	parser *Parser
}

// NewAnalyzer creates a new GitHub Actions analyzer
func NewAnalyzer() *Analyzer {
	return &Analyzer{
		parser: NewParser(),
	}
}

// AnalyzeWorkflow analyzes a GitHub Actions workflow
func (a *Analyzer) AnalyzeWorkflow(filePath string) (*AnalysisResult, error) {
	logger := shared.GetLogger()
	logger.Debug("GitHubActions", "Starting workflow analysis", map[string]interface{}{
		"file": filePath,
	})

	workflow, err := a.parser.ParseFile(filePath)
	if err != nil {
		logger.Error("GitHubActions", "Failed to parse workflow", map[string]interface{}{
			"file":  filePath,
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to parse workflow: %v", err)
	}

	result := a.analyzeWorkflowData(workflow, filePath)
	
	logger.Info("GitHubActions", "Workflow analyzed successfully", map[string]interface{}{
		"file":       filePath,
		"jobs":       len(result.Jobs),
		"steps":      result.TotalSteps,
		"actions":    len(result.ActionUsage),
		"runners":    len(result.RunnerUsage),
	})

	return result, nil
}

// analyzeWorkflowData performs the actual analysis
func (a *Analyzer) analyzeWorkflowData(workflow *Workflow, filePath string) *AnalysisResult {
	result := &AnalysisResult{
		Config:          workflow,
		FilePath:        filePath,
		GeneratedAt:     time.Now(),
		ActionUsage:     make(map[string]int),
		RunnerUsage:     make(map[string]int),
		ServiceUsage:    make(map[string]int),
		CommandPatterns: make(map[string][]string),
		Workflows:       1,
	}

	// Analyze each job
	for jobName, job := range workflow.Jobs {
		jobAnalysis := a.analyzeJob(jobName, job)
		result.Jobs = append(result.Jobs, jobAnalysis)
		result.TotalSteps += jobAnalysis.StepCount

		// Track runner usage
		runner := a.parser.GetRunnerType(job)
		result.RunnerUsage[runner]++

		// Track action and service usage
		for _, action := range jobAnalysis.ActionsUsed {
			result.ActionUsage[action]++
		}

		// Track service usage
		for serviceName := range job.Services {
			result.ServiceUsage[serviceName]++
		}

		// Analyze command patterns
		a.analyzeCommandPatterns(jobAnalysis.RunCommands, result.CommandPatterns)
	}

	// Generate recommendations
	result.Recommendations = a.generateRecommendations(result)
	result.Issues = a.identifyIssues(result)

	return result
}

// analyzeJob analyzes a specific job
func (a *Analyzer) analyzeJob(jobName string, job Job) JobAnalysis {
	analysis := JobAnalysis{
		Name:         jobName,
		Runner:       a.parser.GetRunnerType(job),
		StepCount:    len(job.Steps),
		Dependencies: a.parser.GetJobDependencies(job),
	}

	var allCommands []string
	var actionsUsed []string

	// Analyze each step
	for _, step := range job.Steps {
		if step.Uses != "" {
			actionsUsed = append(actionsUsed, step.Uses)
		}

		commands := a.parser.ExtractRunCommands(step)
		allCommands = append(allCommands, commands...)

		// Check for caching
		if strings.Contains(step.Uses, "cache") {
			analysis.CachingEnabled = true
		}
	}

	analysis.RunCommands = allCommands
	analysis.ActionsUsed = actionsUsed
	analysis.EstimatedTime = a.estimateJobTime(analysis)
	analysis.Recommendations = a.generateJobRecommendations(analysis)
	analysis.SecurityIssues = a.identifySecurityIssues(analysis)

	return analysis
}

// analyzeCommandPatterns analyzes command patterns for go-task opportunities
func (a *Analyzer) analyzeCommandPatterns(commands []string, patterns map[string][]string) {
	for _, cmd := range commands {
		// Use shared pattern matching
		classification := shared.ClassifyCommand(cmd)
		for _, tool := range classification.Tools {
			patterns[tool] = append(patterns[tool], cmd)
		}
	}
}

// generateRecommendations generates recommendations for go-task migration
func (a *Analyzer) generateRecommendations(result *AnalysisResult) []string {
	var recommendations []string

	// Check for go-task migration opportunities
	if len(result.CommandPatterns) > 0 {
		recommendations = append(recommendations, 
			"💡 Consider creating go-task equivalents for repeated command patterns")
	}

	// Check for npm script opportunities
	if npmCommands, exists := result.CommandPatterns["npm"]; exists && len(npmCommands) > 3 {
		recommendations = append(recommendations, 
			"📦 Multiple npm commands detected - consider consolidating into go-task")
	}

	// Check for Docker opportunities  
	if dockerCommands, exists := result.CommandPatterns["docker"]; exists && len(dockerCommands) > 2 {
		recommendations = append(recommendations, 
			"🐳 Docker commands found - go-task could simplify container management")
	}

	// Check for test command patterns
	hasTests := false
	for pattern := range result.CommandPatterns {
		if strings.Contains(strings.ToLower(pattern), "test") {
			hasTests = true
			break
		}
	}
	if hasTests {
		recommendations = append(recommendations, 
			"🧪 Test commands detected - go-task could provide consistent local/CI testing")
	}

	return recommendations
}

// generateJobRecommendations generates job-specific recommendations
func (a *Analyzer) generateJobRecommendations(job JobAnalysis) []string {
	var recommendations []string

	// Check for caching opportunities
	if !job.CachingEnabled && len(job.RunCommands) > 3 {
		recommendations = append(recommendations, "⚡ Consider adding caching to improve build times")
	}

	// Check for go-task opportunities in commands
	taskOpportunities := 0
	for _, cmd := range job.RunCommands {
		if strings.Contains(cmd, "npm run") || 
		   strings.Contains(cmd, "make") || 
		   strings.Contains(cmd, "scripts/") {
			taskOpportunities++
		}
	}

	if taskOpportunities > 2 {
		recommendations = append(recommendations, 
			"🔄 Multiple script executions - good candidate for go-task consolidation")
	}

	return recommendations
}

// identifySecurityIssues identifies potential security issues
func (a *Analyzer) identifySecurityIssues(job JobAnalysis) []string {
	var issues []string

	// Check for pinned action versions
	for _, action := range job.ActionsUsed {
		if !strings.Contains(action, "@") || strings.Contains(action, "@latest") {
			issues = append(issues, fmt.Sprintf("⚠️ Action '%s' not pinned to specific version", action))
		}
	}

	// Check for sensitive commands
	for _, cmd := range job.RunCommands {
		if strings.Contains(strings.ToLower(cmd), "curl") && strings.Contains(cmd, "|") && strings.Contains(cmd, "sh") {
			issues = append(issues, "🚨 Potential security risk: piping curl to shell")
		}
	}

	return issues
}

// identifyIssues identifies workflow-level issues
func (a *Analyzer) identifyIssues(result *AnalysisResult) []string {
	var issues []string

	// Check for missing workflow name
	if result.Config.Name == "" {
		issues = append(issues, "⚠️ Workflow missing name field")
	}

	// Check for jobs without needs (potential parallelization)
	independentJobs := 0
	for _, job := range result.Jobs {
		if len(job.Dependencies) == 0 {
			independentJobs++
		}
	}

	if independentJobs > 3 {
		issues = append(issues, "💡 Many independent jobs - consider if some should have dependencies")
	}

	return issues
}

// estimateJobTime provides rough time estimate for job
func (a *Analyzer) estimateJobTime(job JobAnalysis) string {
	// Simple heuristic based on step count and command types
	baseTime := job.StepCount * 30 // 30 seconds per step baseline
	
	for _, cmd := range job.RunCommands {
		if strings.Contains(cmd, "npm install") || strings.Contains(cmd, "npm ci") {
			baseTime += 60 // Add minute for installs
		}
		if strings.Contains(cmd, "docker build") {
			baseTime += 180 // Add 3 minutes for Docker builds
		}
		if strings.Contains(cmd, "test") {
			baseTime += 30 // Add 30 seconds for tests
		}
	}

	if baseTime < 60 {
		return "< 1 min"
	} else if baseTime < 300 {
		return fmt.Sprintf("~%d min", baseTime/60)
	} else {
		return "> 5 min"
	}
}