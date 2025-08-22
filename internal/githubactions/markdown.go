package githubactions

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/nichecode/pipeline-analyzer/internal/shared"
)

// MarkdownGenerator generates markdown documentation for GitHub Actions
type MarkdownGenerator struct{}

// NewMarkdownGenerator creates a new markdown generator
func NewMarkdownGenerator() *MarkdownGenerator {
	return &MarkdownGenerator{}
}

// GenerateMainReadme generates the main README for GitHub Actions analysis
func (g *MarkdownGenerator) GenerateMainReadme(results []*AnalysisResult, configPath string) string {
	totalJobs := 0
	totalSteps := 0
	totalWorkflows := len(results)

	for _, result := range results {
		totalJobs += len(result.Jobs)
		totalSteps += result.TotalSteps
	}

	content := fmt.Sprintf(`# GitHub Actions Analysis Report

**Generated:** %s
**Config Path:** %s

## 📊 Overview

- **Workflow files:** %d
- **Total jobs:** %d  
- **Total steps:** %d

`, time.Now().Format(time.RFC3339), configPath, totalWorkflows, totalJobs, totalSteps)

	// Add workflow diagram
	diagram := generateWorkflowDiagram(results)
	if diagram != "" {
		content += "## 📊 Workflow Overview\n\n"
		content += diagram
		content += "\n"
	}

	content += `## 🚀 Quick Start for go-task Migration

1. **[🔄 go-task Migration Guide](summaries/go-task-migration.md)** - Convert workflows to go-task
2. **[🛠️ Actions Usage](summaries/actions-usage.md)** - Actions used across workflows
3. **[⚡ Commands Analysis](summaries/commands-analysis.md)** - Commands that could become tasks

## 📁 Directory Structure

### Workflows
Individual workflow analysis:

`

	// List workflow files
	for _, result := range results {
		workflowName := result.Config.Name
		if workflowName == "" {
			workflowName = "Untitled Workflow"
		}
		content += fmt.Sprintf("- [workflows/%s.md](workflows/%s.md) - %s\n", 
			sanitizeFilename(workflowName), sanitizeFilename(workflowName), workflowName)
	}

	content += `
### Jobs
Individual job analysis with go-task conversion opportunities:

`

	// List all jobs across workflows
	var allJobs []JobAnalysis
	for _, result := range results {
		allJobs = append(allJobs, result.Jobs...)
	}

	// Sort jobs by name for consistent output
	sort.Slice(allJobs, func(i, j int) bool {
		return allJobs[i].Name < allJobs[j].Name
	})

	for _, job := range allJobs {
		content += fmt.Sprintf("- [jobs/%s.md](jobs/%s.md) - %s\n", 
			sanitizeFilename(job.Name), sanitizeFilename(job.Name), job.Name)
	}

	content += `
### Analysis Summaries

- [🔄 go-task Migration Guide](summaries/go-task-migration.md) - **Start here for build refactoring**
- [🛠️ Actions Usage](summaries/actions-usage.md) - GitHub Actions used
- [🏃 Runners Analysis](summaries/runners-analysis.md) - Runner usage patterns  
- [⚡ Commands Analysis](summaries/commands-analysis.md) - Commands suitable for go-task

## 🎯 Build Refactoring Recommendations

`

	// Generate consolidated recommendations
	allRecommendations := make(map[string]int)
	for _, result := range results {
		for _, rec := range result.Recommendations {
			allRecommendations[rec]++
		}
	}

	if len(allRecommendations) > 0 {
		for rec, count := range allRecommendations {
			if count > 1 {
				content += fmt.Sprintf("- %s *(affects %d workflows)*\n", rec, count)
			} else {
				content += fmt.Sprintf("- %s\n", rec)
			}
		}
	} else {
		content += "- No specific recommendations at this time\n"
	}

	// Add popular actions summary
	content += g.generatePopularActionsSection(results)

	content += `
## 🔍 Next Steps for CI/CD Refactoring

1. **Review [go-task Migration Guide](summaries/go-task-migration.md)** for refactoring strategy
2. **Identify commands** that can be moved to go-task from [Commands Analysis](summaries/commands-analysis.md)
3. **Create Taskfile.yml** with extracted commands for local execution
4. **Update workflows** to call ` + "`task <task-name>`" + ` instead of direct commands
5. **Test locally** by running tasks before pushing to CI

This approach makes your builds **portable, testable locally, and CI-agnostic**.
`

	return content
}

// generatePopularActionsSection creates a section showing most used actions
func (g *MarkdownGenerator) generatePopularActionsSection(results []*AnalysisResult) string {
	actionUsage := make(map[string]int)
	
	for _, result := range results {
		for action, count := range result.ActionUsage {
			actionUsage[action] += count
		}
	}

	if len(actionUsage) == 0 {
		return ""
	}

	content := `
## 🔧 Most Used Actions

| Action | Usage Count |
|--------|-------------|
`

	// Sort actions by usage count (descending)
	type actionCount struct {
		action string
		count  int
	}

	var actions []actionCount
	for action, count := range actionUsage {
		actions = append(actions, actionCount{action, count})
	}

	sort.Slice(actions, func(i, j int) bool {
		return actions[i].count > actions[j].count
	})

	// Show top 10 actions
	limit := 10
	if len(actions) < limit {
		limit = len(actions)
	}

	for i := 0; i < limit; i++ {
		content += fmt.Sprintf("| %s | %d |\n", actions[i].action, actions[i].count)
	}

	return content
}

// GenerateWorkflowAnalysis generates analysis for a specific workflow
func (g *MarkdownGenerator) GenerateWorkflowAnalysis(result *AnalysisResult) string {
	workflowName := result.Config.Name
	if workflowName == "" {
		workflowName = "Untitled Workflow"
	}

	content := fmt.Sprintf(`# Workflow: %s

**File:** %s  
**Generated:** %s

## 📊 Overview

- **Jobs:** %d
- **Total steps:** %d
- **Actions used:** %d unique actions
- **Runners used:** %d unique runners

## 🔄 go-task Migration Opportunities

`, workflowName, filepath.Base(result.FilePath), 
		result.GeneratedAt.Format(time.RFC3339), 
		len(result.Jobs), result.TotalSteps, 
		len(result.ActionUsage), len(result.RunnerUsage))

	// Show specific recommendations for this workflow
	if len(result.Recommendations) > 0 {
		for _, rec := range result.Recommendations {
			content += fmt.Sprintf("- %s\n", rec)
		}
	} else {
		content += "- No specific migration opportunities identified\n"
	}

	content += `
## 📋 Jobs Overview

| Job | Runner | Steps | Commands | Actions |
|-----|--------|-------|----------|---------|
`

	for _, job := range result.Jobs {
		content += fmt.Sprintf("| [%s](../jobs/%s.md) | %s | %d | %d | %d |\n",
			job.Name, sanitizeFilename(job.Name), job.Runner, 
			job.StepCount, len(job.RunCommands), len(job.ActionsUsed))
	}

	// Add command patterns analysis
	if len(result.CommandPatterns) > 0 {
		content += `
## ⚡ Command Patterns Detected

These command patterns are good candidates for go-task consolidation:

`
		for pattern, commands := range result.CommandPatterns {
			content += fmt.Sprintf("### %s (%d commands)\n\n", pattern, len(commands))
			for _, cmd := range commands {
				content += fmt.Sprintf("- `%s`\n", cmd)
			}
			content += "\n"
		}
	}

	content += `
## 🔍 Navigation

- [← Back to GitHub Actions Overview](../README.md)
- [All Jobs](../summaries/actions-usage.md)
- [go-task Migration Guide](../summaries/go-task-migration.md)
`

	return content
}

// GenerateJobAnalysis generates analysis for a specific job
func (g *MarkdownGenerator) GenerateJobAnalysis(job JobAnalysis, workflowName string) string {
	content := fmt.Sprintf(`# Job: %s

**Workflow:** %s  
**Runner:** %s  
**Estimated Duration:** %s  
**Caching Enabled:** %t

## 📊 Overview

- **Steps:** %d
- **Run commands:** %d
- **Actions used:** %d
- **Dependencies:** %s

`, job.Name, workflowName, job.Runner, job.EstimatedTime, job.CachingEnabled,
		job.StepCount, len(job.RunCommands), len(job.ActionsUsed), 
		strings.Join(job.Dependencies, ", "))

	// Show commands that could become go-task tasks
	if len(job.RunCommands) > 0 {
		content += `## ⚡ Commands (go-task candidates)

These commands could be extracted into go-task:

`
		for i, cmd := range job.RunCommands {
			content += fmt.Sprintf("%d. `%s`\n", i+1, cmd)
		}

		content += `
### Suggested go-task conversion:

` + "```yaml" + `
version: '3'
tasks:
  ` + sanitizeTaskName(job.Name) + `:
    desc: "` + job.Name + ` task"
    cmds:
`
		for _, cmd := range job.RunCommands {
			content += fmt.Sprintf("      - %s\n", cmd)
		}

		content += "```\n\n"
	}

	// Show actions used
	if len(job.ActionsUsed) > 0 {
		content += `## 🛠️ GitHub Actions Used

`
		for _, action := range job.ActionsUsed {
			content += fmt.Sprintf("- `%s`\n", action)
		}
		content += "\n"
	}

	// Show recommendations
	if len(job.Recommendations) > 0 {
		content += `## 💡 Optimization Recommendations

`
		for _, rec := range job.Recommendations {
			content += fmt.Sprintf("- %s\n", rec)
		}
		content += "\n"
	}

	// Show security issues
	if len(job.SecurityIssues) > 0 {
		content += `## 🚨 Security Issues

`
		for _, issue := range job.SecurityIssues {
			content += fmt.Sprintf("- %s\n", issue)
		}
		content += "\n"
	}

	content += `## 🔍 Navigation

- [← Back to GitHub Actions Overview](../README.md)
- [go-task Migration Guide](../summaries/go-task-migration.md)
`

	return content
}

// Generate summary files
func (g *MarkdownGenerator) GenerateActionsUsage(results []*AnalysisResult) string {
	actionUsage := make(map[string]int)
	actionJobs := make(map[string][]string)
	
	for _, result := range results {
		for _, job := range result.Jobs {
			for _, action := range job.ActionsUsed {
				actionUsage[action]++
				actionJobs[action] = append(actionJobs[action], job.Name)
			}
		}
	}

	content := `# GitHub Actions Usage Analysis

## 📊 Actions by Popularity

| Action | Usage Count | Jobs Using |
|--------|-------------|------------|
`

	// Sort by usage count
	type actionCount struct {
		action string
		count  int
	}

	var actions []actionCount
	for action, count := range actionUsage {
		actions = append(actions, actionCount{action, count})
	}

	sort.Slice(actions, func(i, j int) bool {
		return actions[i].count > actions[j].count
	})

	for _, ac := range actions {
		jobs := strings.Join(actionJobs[ac.action], ", ")
		if len(jobs) > 100 {
			jobs = jobs[:97] + "..."
		}
		content += fmt.Sprintf("| `%s` | %d | %s |\n", ac.action, ac.count, jobs)
	}

	content += `
## 🔍 Navigation

- [← Back to GitHub Actions Overview](../README.md)
- [go-task Migration Guide](go-task-migration.md)
`

	return content
}

func (g *MarkdownGenerator) GenerateRunnersAnalysis(results []*AnalysisResult) string {
	runnerUsage := make(map[string]int)
	
	for _, result := range results {
		for runner, count := range result.RunnerUsage {
			runnerUsage[runner] += count
		}
	}

	content := `# GitHub Actions Runners Analysis

## 🏃 Runner Usage

| Runner | Jobs Count |
|--------|------------|
`

	for runner, count := range runnerUsage {
		content += fmt.Sprintf("| `%s` | %d |\n", runner, count)
	}

	content += `
## 🔍 Navigation

- [← Back to GitHub Actions Overview](../README.md)
`

	return content
}

func (g *MarkdownGenerator) GenerateCommandsAnalysis(results []*AnalysisResult) string {
	allPatterns := make(map[string][]string)
	
	for _, result := range results {
		for pattern, commands := range result.CommandPatterns {
			allPatterns[pattern] = append(allPatterns[pattern], commands...)
		}
	}

	content := `# Commands Analysis for go-task Migration

## ⚡ Command Patterns Suitable for go-task

`

	for pattern, commands := range allPatterns {
		content += fmt.Sprintf("### %s (%d commands)\n\n", pattern, len(commands))
		
		// Deduplicate commands
		uniqueCommands := make(map[string]bool)
		for _, cmd := range commands {
			uniqueCommands[cmd] = true
		}
		
		for cmd := range uniqueCommands {
			content += fmt.Sprintf("- `%s`\n", cmd)
		}
		content += "\n"
	}

	content += `
## 🔍 Navigation

- [← Back to GitHub Actions Overview](../README.md)
- [go-task Migration Guide](go-task-migration.md)
`

	return content
}

func (g *MarkdownGenerator) GenerateGoTaskMigration(results []*AnalysisResult) string {
	content := `# go-task Migration Guide for GitHub Actions

## 🎯 Migration Strategy

This guide helps you refactor your GitHub Actions workflows to use **go-task as the command runner**, making your builds:

- **Locally executable** - Run the same commands locally as in CI
- **CI-agnostic** - Easy to switch between GitHub Actions, CircleCI, etc.
- **Testable** - Debug build issues without pushing to CI

## 🔄 Refactoring Pattern

**Current Pattern:**
` + "```yaml" + `
- name: Build and test
  run: |
    npm install
    npm run build
    npm run test
` + "```" + `

**Target Pattern:**
` + "```yaml" + `
- name: Build and test
  run: task build-and-test
` + "```" + `

With corresponding Taskfile.yml:
` + "```yaml" + `
version: '3'
tasks:
  build-and-test:
    desc: "Build and test the application"
    cmds:
      - npm install
      - npm run build  
      - npm run test
` + "```" + `

## 📋 Migration Checklist

### 1. Analyze Current Commands
`

	// Collect all unique commands across workflows
	allCommands := make(map[string]bool)
	for _, result := range results {
		for _, job := range result.Jobs {
			for _, cmd := range job.RunCommands {
				allCommands[cmd] = true
			}
		}
	}

	if len(allCommands) > 0 {
		content += fmt.Sprintf("Found **%d unique commands** across your workflows:\n\n", len(allCommands))
		
		count := 0
		for cmd := range allCommands {
			if count >= 10 { // Limit to first 10 for readability
				content += fmt.Sprintf("- ... and %d more\n", len(allCommands)-10)
				break
			}
			content += fmt.Sprintf("- `%s`\n", cmd)
			count++
		}
	}

	content += `
### 2. Create Taskfile.yml

Create a Taskfile.yml in your repository root:

` + "```yaml" + `
version: '3'

vars:
  # Define common variables here

tasks:
  # Extract your commands into logical tasks
  install:
    desc: "Install dependencies"
    cmds:
      - npm install

  build:
    desc: "Build the application"  
    deps: [install]
    cmds:
      - npm run build

  test:
    desc: "Run tests"
    deps: [build]  
    cmds:
      - npm run test

  deploy:
    desc: "Deploy application"
    deps: [test]
    cmds:
      - # Add deployment commands here
` + "```" + `

### 3. Update GitHub Actions Workflows

Replace complex run commands with task calls:

**Before:**
` + "```yaml" + `
- name: Install dependencies
  run: npm install
  
- name: Build  
  run: npm run build
  
- name: Test
  run: npm run test
` + "```" + `

**After:**
` + "```yaml" + `
- name: Install go-task
  run: go install github.com/go-task/task/v3/cmd/task@latest
  
- name: Run build pipeline
  run: task test  # This runs install -> build -> test
` + "```" + `

### 4. Test Locally

Before pushing changes:

` + "```bash" + `
# Install go-task locally
go install github.com/go-task/task/v3/cmd/task@latest

# Test your tasks
task install
task build  
task test
` + "```" + `

## 🚀 Benefits After Migration

- ✅ **Local Development** - Run CI commands locally
- ✅ **Faster Debugging** - No need to push to test builds  
- ✅ **CI Portability** - Easy to switch CI providers
- ✅ **Consistent Environments** - Same commands everywhere
- ✅ **Better Developer Experience** - Clear, documented tasks

## 💡 Pro Tips

1. **Group Related Commands** - Combine related steps into single tasks
2. **Use Dependencies** - Let go-task handle execution order with deps
3. **Add Descriptions** - Use desc for clear task documentation  
4. **Define Variables** - Use vars section for reusable values
5. **Local vs CI** - Use platforms to handle OS differences

## 🔍 Next Steps

1. **Start Small** - Pick one workflow to convert first
2. **Create Tasks** - Extract commands into logical groups  
3. **Update Workflow** - Replace run commands with task calls
4. **Test Locally** - Verify tasks work before committing
5. **Iterate** - Apply pattern to remaining workflows

## 🔍 Navigation

- [← Back to GitHub Actions Overview](../README.md)
- [Commands Analysis](commands-analysis.md)
- [Actions Usage](actions-usage.md)
`

	return content
}

// generateWorkflowDiagram creates a simple Mermaid diagram for GitHub Actions workflows
func generateWorkflowDiagram(results []*AnalysisResult) string {
	if len(results) == 0 {
		return ""
	}
	
	diagram := &shared.MermaidDiagram{
		Title: "GitHub Actions Workflows",
	}
	
	// Process each workflow (limit to first workflow for simplicity)
	result := results[0]
	
	// Create nodes for jobs
	for _, job := range result.Jobs {
		commands := job.RunCommands
		if len(commands) > 5 {
			commands = commands[:5]
		}
		
		node := shared.MermaidNode{
			ID:          shared.CleanNodeID(job.Name),
			Label:       job.Name,
			Description: fmt.Sprintf("Runner: %s", job.Runner),
			Commands:    commands,
			NodeType:    shared.ClassifyNodeType(job.Name, commands),
		}
		diagram.Nodes = append(diagram.Nodes, node)
	}
	
	// Add simple sequential flow for jobs (GitHub Actions job dependencies are complex)
	for i := 0; i < len(diagram.Nodes)-1; i++ {
		edge := shared.MermaidEdge{
			From: diagram.Nodes[i].ID,
			To:   diagram.Nodes[i+1].ID,
		}
		diagram.Edges = append(diagram.Edges, edge)
	}
	
	return diagram.Generate()
}

// sanitizeTaskName creates a valid go-task name
func sanitizeTaskName(name string) string {
	// Convert to lowercase and replace spaces/special chars with hyphens
	result := strings.ToLower(name)
	result = strings.ReplaceAll(result, " ", "-")
	result = strings.ReplaceAll(result, "_", "-")
	
	// Remove other problematic characters
	validChars := "abcdefghijklmnopqrstuvwxyz0123456789-"
	cleaned := ""
	for _, char := range result {
		if strings.ContainsRune(validChars, char) {
			cleaned += string(char)
		}
	}
	
	return cleaned
}