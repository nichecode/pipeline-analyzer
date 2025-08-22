package circleci

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/nichecode/pipeline-analyzer/internal/shared"
)

// GenerateMainReadme generates the main README.md file
func GenerateMainReadme(analysis *Analysis, configPath string) string {
	var sb strings.Builder

	sb.WriteString("# CircleCI Analysis Report\n\n")
	sb.WriteString(fmt.Sprintf("**Generated:** %s\n", analysis.GeneratedAt.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("**Config:** %s\n\n", configPath))

	// Overview section
	sb.WriteString("## 📊 Overview\n\n")
	sb.WriteString(fmt.Sprintf("- **Unique jobs:** %d\n", analysis.TotalJobs))
	sb.WriteString(fmt.Sprintf("- **Reusable commands:** %d\n", analysis.TotalCommands))

	if analysis.TotalWorkflows > 0 {
		workflowNames := GetAllWorkflowNames(analysis.Config)
		sb.WriteString(fmt.Sprintf("- **Workflows:** %s\n", strings.Join(workflowNames, ", ")))
	} else {
		sb.WriteString("- **Workflows:** None found\n")
	}
	sb.WriteString("\n")

	// Workflow diagram
	diagram := generateCircleCIDiagram(analysis)
	if diagram != "" {
		sb.WriteString("## 📊 Workflow Overview\n\n")
		sb.WriteString(diagram)
		sb.WriteString("\n")
	}

	// Quick Start section
	sb.WriteString("## 🚀 Quick Start\n\n")
	sb.WriteString("1. **[📋 Migration Checklist](migration-checklist.md)** - Your step-by-step guide\n")
	sb.WriteString("2. **[📈 Job Usage Analysis](summaries/job-usage.md)** - Job reuse patterns and dependencies\n")
	sb.WriteString("3. **[⚡ Commands Analysis](summaries/commands.md)** - All run commands for conversion\n\n")

	// Directory Structure section
	sb.WriteString("## 📁 Directory Structure\n\n")

	sb.WriteString("### Jobs\n")
	sb.WriteString("Individual job analysis with run commands and configuration:\n\n")
	jobNames := GetAllJobNames(analysis.Config)
	sort.Strings(jobNames)
	for _, jobName := range jobNames {
		normalizedName := NormalizeJobName(jobName)
		sb.WriteString(fmt.Sprintf("- [jobs/%s.md](jobs/%s.md)\n", normalizedName, normalizedName))
	}
	sb.WriteString("\n")

	if analysis.TotalCommands > 0 {
		sb.WriteString("### Reusable Commands\n")
		sb.WriteString("Reusable command definitions with analysis and usage patterns:\n\n")
		commandNames := GetAllCommandNames(analysis.Config)
		sort.Strings(commandNames)
		for _, commandName := range commandNames {
			normalizedName := NormalizeJobName(commandName)
			sb.WriteString(fmt.Sprintf("- [commands/%s.md](commands/%s.md)\n", normalizedName, normalizedName))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("### Workflows\n")
	sb.WriteString("Workflow structure and job dependencies:\n\n")
	if analysis.TotalWorkflows > 0 {
		workflowNames := GetAllWorkflowNames(analysis.Config)
		sort.Strings(workflowNames)
		for _, workflowName := range workflowNames {
			normalizedName := NormalizeJobName(workflowName)
			sb.WriteString(fmt.Sprintf("- [workflows/%s.md](workflows/%s.md)\n", normalizedName, normalizedName))
		}
	} else {
		sb.WriteString("- No workflows found\n")
	}
	sb.WriteString("\n")

	// Analysis Summaries section
	sb.WriteString("### Analysis Summaries\n\n")
	sb.WriteString("- [📈 Job Usage & Dependencies](summaries/job-usage.md)\n")
	sb.WriteString("- [📝 All Jobs Index](summaries/all-jobs.md)\n")
	sb.WriteString("- [⚡ Commands Analysis](summaries/commands.md)\n")
	sb.WriteString("- [🐳 Docker & Scripts](summaries/docker-and-scripts.md)\n")
	sb.WriteString("- [⚙️ Executors & Images](summaries/executors-and-images.md)\n")
	sb.WriteString("- [🔄 Workflows Index](summaries/workflows.md)\n\n")

	// Next Steps section
	sb.WriteString("## 🎯 Next Steps\n\n")
	sb.WriteString("1. **Start with [Migration Checklist](migration-checklist.md)**\n")
	sb.WriteString("2. **Review most frequently used jobs** from [job usage analysis](summaries/job-usage.md)\n")
	sb.WriteString("3. **Examine job dependencies** to understand execution order\n")
	sb.WriteString("4. **Begin converting** highest-impact jobs to go-task format\n\n")

	// Most Used Jobs section
	sb.WriteString("## 🔍 Most Used Jobs\n\n")
	mostUsedJobs := GetMostUsedJobs(analysis)
	if len(mostUsedJobs) > 0 {
		sb.WriteString("| Job | Usage Count | Link |\n")
		sb.WriteString("|-----|-------------|------|\n")

		// Show top 10
		limit := 10
		if len(mostUsedJobs) < limit {
			limit = len(mostUsedJobs)
		}

		for i := 0; i < limit; i++ {
			job := mostUsedJobs[i]
			normalizedName := NormalizeJobName(job.Name)
			sb.WriteString(fmt.Sprintf("| %s | %d | [View Details](jobs/%s.md) |\n",
				job.Name, job.Count, normalizedName))
		}
	} else {
		sb.WriteString("No workflow data available\n")
	}

	return sb.String()
}

// GenerateMigrationChecklist generates the migration checklist
func GenerateMigrationChecklist(analysis *Analysis) string {
	var sb strings.Builder

	sb.WriteString("# CircleCI → go-task Migration Checklist\n\n")
	sb.WriteString(fmt.Sprintf("**Generated:** %s\n\n", analysis.GeneratedAt.Format(time.RFC3339)))

	sb.WriteString("## 🎯 Migration Overview\n\n")
	sb.WriteString("This checklist guides you through converting your CircleCI configuration to a local go-task setup.\n\n")

	// Quick Stats
	sb.WriteString("### 📊 Current State\n\n")
	sb.WriteString(fmt.Sprintf("- **%d jobs** to migrate\n", analysis.TotalJobs))
	sb.WriteString(fmt.Sprintf("- **%d workflows** to understand\n", analysis.TotalWorkflows))

	dockerJobs, otherJobs := CountDockerUsage(analysis.Config)
	sb.WriteString(fmt.Sprintf("- **%d jobs** use Docker\n", dockerJobs))
	sb.WriteString(fmt.Sprintf("- **%d jobs** use other executors\n", otherJobs))
	sb.WriteString("\n")

	// Step-by-step process
	sb.WriteString("## 🔄 Migration Steps\n\n")
	sb.WriteString("### 1. **Understand job dependencies**\n")
	sb.WriteString("Review [Job Usage Analysis](summaries/job-usage.md) to see which jobs depend on others.\n\n")

	sb.WriteString("### 2. **Examine Docker images and executors**\n")
	sb.WriteString("Check [Docker & Scripts](summaries/docker-and-scripts.md) for container requirements.\n\n")

	sb.WriteString("### 3. **Analyze command patterns**\n")
	sb.WriteString("Study [Commands Analysis](summaries/commands.md) to understand build patterns.\n\n")

	sb.WriteString("### 4. **Start with high-impact jobs**\n")
	sb.WriteString("Begin with the most frequently used jobs from the overview.\n\n")

	sb.WriteString("### 5. **Convert commands to go-task**\n")
	sb.WriteString("Transform run commands to task format - see individual job files.\n\n")

	sb.WriteString("### 6. **Test task equivalents** locally\n\n")

	// Key Files section
	sb.WriteString("## Key Files to Examine\n\n")
	sb.WriteString("- [Job Usage Analysis](summaries/job-usage.md) - How jobs are reused and their dependencies\n")
	sb.WriteString("- [All Jobs](summaries/all-jobs.md) - Complete job list with descriptions\n")
	sb.WriteString("- [Docker & Scripts](summaries/docker-and-scripts.md) - Docker/script patterns\n")
	sb.WriteString("- [Workflows](summaries/workflows.md) - Workflow structure and job orchestration\n\n")

	// Suggested go-task Structure
	sb.WriteString("## Suggested go-task Structure\n\n")
	sb.WriteString("```yaml\n")
	sb.WriteString("version: '3'\n\n")
	sb.WriteString("tasks:\n")

	jobNames := GetAllJobNames(analysis.Config)
	sort.Strings(jobNames)
	for _, jobName := range jobNames {
		job := analysis.Config.Jobs[jobName]
		sb.WriteString(fmt.Sprintf("  %s:\n", jobName))

		if job.Description != "" {
			sb.WriteString(fmt.Sprintf("    desc: \"%s\"\n", job.Description))
		} else {
			sb.WriteString("    desc: \"Migrated from CircleCI job\"\n")
		}

		// Try to detect dependencies
		if deps, hasDeps := analysis.JobDependencies[jobName]; hasDeps && len(deps) > 0 {
			sb.WriteString(fmt.Sprintf("    deps: [%s]\n", strings.Join(deps, ", ")))
		}

		sb.WriteString("    cmds:\n")
		sb.WriteString(fmt.Sprintf("      - # Convert run commands from jobs/%s.md\n", NormalizeJobName(jobName)))
		sb.WriteString("\n")
	}
	sb.WriteString("```\n\n")

	// Navigation
	sb.WriteString("## Navigation\n\n")
	sb.WriteString("- [← Back to Overview](../README.md)\n")
	sb.WriteString("- [All Jobs](summaries/all-jobs.md)\n")
	sb.WriteString("- [Job Usage Analysis](summaries/job-usage.md)\n")

	return sb.String()
}

// GenerateJobMarkdown generates markdown for a single job
func GenerateJobMarkdown(jobAnalysis *JobAnalysis) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# Job: %s\n\n", jobAnalysis.Name))

	if jobAnalysis.Description != "" {
		sb.WriteString(fmt.Sprintf("**Description:** %s\n\n", jobAnalysis.Description))
	}

	// Usage info
	sb.WriteString("## Usage Information\n\n")
	sb.WriteString(fmt.Sprintf("- **Used in workflows:** %d times\n", jobAnalysis.UsageCount))

	if len(jobAnalysis.Dependencies) > 0 {
		sb.WriteString(fmt.Sprintf("- **Depends on:** %s\n", strings.Join(jobAnalysis.Dependencies, ", ")))
	}

	if jobAnalysis.Executor != "" {
		sb.WriteString(fmt.Sprintf("- **Executor:** %s\n", jobAnalysis.Executor))
	}
	sb.WriteString("\n")

	// Docker Images
	if len(jobAnalysis.DockerImages) > 0 {
		sb.WriteString("## Docker Images\n\n")
		for _, image := range jobAnalysis.DockerImages {
			sb.WriteString(fmt.Sprintf("- `%s`\n", image))
		}
		sb.WriteString("\n")
	}

	// Commands
	if len(jobAnalysis.Commands) > 0 {
		sb.WriteString("## Run Commands\n\n")
		for i, command := range jobAnalysis.Commands {
			sb.WriteString(fmt.Sprintf("### Command %d\n\n", i+1))
			sb.WriteString("```bash\n")
			sb.WriteString(command)
			sb.WriteString("\n```\n\n")
		}
	}

	// Pattern Usage
	if len(jobAnalysis.Patterns) > 0 {
		sb.WriteString("## Command Patterns Detected\n\n")
		var patterns []string
		for pattern := range jobAnalysis.Patterns {
			patterns = append(patterns, pattern)
		}
		sort.Strings(patterns)

		for _, pattern := range patterns {
			count := jobAnalysis.Patterns[pattern]
			sb.WriteString(fmt.Sprintf("- **%s:** %d occurrences\n", pattern, count))
		}
		sb.WriteString("\n")
	}

	// go-task conversion suggestion
	sb.WriteString("## Suggested go-task Conversion\n\n")
	sb.WriteString("```yaml\n")
	sb.WriteString(fmt.Sprintf("%s:\n", jobAnalysis.Name))
	if jobAnalysis.Description != "" {
		sb.WriteString(fmt.Sprintf("  desc: \"%s\"\n", jobAnalysis.Description))
	}
	if len(jobAnalysis.Dependencies) > 0 {
		sb.WriteString(fmt.Sprintf("  deps: [%s]\n", strings.Join(jobAnalysis.Dependencies, ", ")))
	}
	sb.WriteString("  cmds:\n")

	if len(jobAnalysis.Commands) > 0 {
		for _, command := range jobAnalysis.Commands {
			// Simple command formatting for go-task
			sb.WriteString(fmt.Sprintf("    - %s\n", command))
		}
	} else {
		sb.WriteString("    - # Add commands here\n")
	}
	sb.WriteString("```\n\n")

	// Navigation
	sb.WriteString("## Navigation\n\n")
	sb.WriteString("- [← Back to All Jobs](../summaries/all-jobs.md)\n")
	sb.WriteString("- [← Back to Overview](../README.md)\n")

	return sb.String()
}

// GenerateWorkflowMarkdown generates markdown for a single workflow
func GenerateWorkflowMarkdown(workflowAnalysis *WorkflowAnalysis) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# Workflow: %s\n\n", workflowAnalysis.Name))

	sb.WriteString("## Job Execution Order\n\n")

	if len(workflowAnalysis.Jobs) > 0 {
		sb.WriteString("| Job | Dependencies | Context |\n")
		sb.WriteString("|-----|--------------|----------|\n")

		for _, job := range workflowAnalysis.Jobs {
			deps := "None"
			if len(job.Requires) > 0 {
				deps = strings.Join(job.Requires, ", ")
			}

			context := "None"
			if len(job.Context) > 0 {
				context = strings.Join(job.Context, ", ")
			}

			normalizedName := NormalizeJobName(job.Name)
			jobLink := fmt.Sprintf("[%s](../jobs/%s.md)", job.Name, normalizedName)
			sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n", jobLink, deps, context))
		}
	} else {
		sb.WriteString("No jobs found in this workflow.\n")
	}
	sb.WriteString("\n")

	// Dependency Graph
	sb.WriteString("## Dependency Graph\n\n")
	sb.WriteString("```\n")

	// Create a simple text-based dependency graph
	independentJobs := []string{}
	dependentJobs := map[string][]string{}

	for _, job := range workflowAnalysis.Jobs {
		if len(job.Requires) == 0 {
			independentJobs = append(independentJobs, job.Name)
		} else {
			dependentJobs[job.Name] = job.Requires
		}
	}

	if len(independentJobs) > 0 {
		sb.WriteString("Independent jobs (run first):\n")
		for _, job := range independentJobs {
			sb.WriteString(fmt.Sprintf("├── %s\n", job))
		}
		sb.WriteString("\n")
	}

	if len(dependentJobs) > 0 {
		sb.WriteString("Dependent jobs:\n")
		for job, deps := range dependentJobs {
			sb.WriteString(fmt.Sprintf("├── %s\n", job))
			for _, dep := range deps {
				sb.WriteString(fmt.Sprintf("│   └── requires: %s\n", dep))
			}
		}
	}

	sb.WriteString("```\n\n")

	// Navigation
	sb.WriteString("## Navigation\n\n")
	sb.WriteString("- [← Back to Workflows](../summaries/workflows.md)\n")
	sb.WriteString("- [← Back to Overview](../README.md)\n")

	return sb.String()
}

// GenerateCommandMarkdown generates markdown for a reusable command
func GenerateCommandMarkdown(cmdAnalysis *CommandAnalysis) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# Reusable Command: %s\n\n", cmdAnalysis.Name))

	// Description
	if cmdAnalysis.Description != "" {
		sb.WriteString(fmt.Sprintf("**Description:** %s\n\n", cmdAnalysis.Description))
	}

	// Usage
	sb.WriteString("## 📊 Usage Statistics\n\n")
	sb.WriteString(fmt.Sprintf("- **Used by:** %d job(s)\n", cmdAnalysis.UsageCount))
	
	if len(cmdAnalysis.Parameters) > 0 {
		sb.WriteString(fmt.Sprintf("- **Parameters:** %d\n", len(cmdAnalysis.Parameters)))
	}
	sb.WriteString("\n")

	// Parameters
	if len(cmdAnalysis.Parameters) > 0 {
		sb.WriteString("## ⚙️ Parameters\n\n")
		sb.WriteString("| Parameter | Type | Description |\n")
		sb.WriteString("|-----------|------|-------------|\n")
		
		for paramName, paramDef := range cmdAnalysis.Parameters {
			paramType := "string"
			description := "No description"
			
			if paramMap, ok := paramDef.(map[string]interface{}); ok {
				if t, exists := paramMap["type"]; exists {
					paramType = fmt.Sprintf("%v", t)
				}
				if d, exists := paramMap["description"]; exists {
					description = fmt.Sprintf("%v", d)
				}
			}
			
			sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n", paramName, paramType, description))
		}
		sb.WriteString("\n")
	}

	// Commands/Steps
	if len(cmdAnalysis.Commands) > 0 {
		sb.WriteString("## ⚡ Steps\n\n")
		for i, command := range cmdAnalysis.Commands {
			if len(cmdAnalysis.Commands) > 1 {
				sb.WriteString(fmt.Sprintf("**Step %d:**\n", i+1))
			}
			sb.WriteString("```bash\n")
			sb.WriteString(command)
			sb.WriteString("\n```\n\n")
		}
	} else {
		sb.WriteString("## ⚡ Steps\n\n")
		sb.WriteString("This command uses non-run steps (e.g., cache operations, setup steps).\n\n")
	}

	// Command patterns
	if len(cmdAnalysis.Patterns) > 0 {
		sb.WriteString("## 🔍 Tool Usage\n\n")
		for pattern, count := range cmdAnalysis.Patterns {
			if count > 0 {
				sb.WriteString(fmt.Sprintf("- **%s:** %d occurrence(s)\n", pattern, count))
			}
		}
		sb.WriteString("\n")
	}

	// Navigation
	sb.WriteString("## Navigation\n\n")
	sb.WriteString("- [← Back to Commands](../summaries/commands.md)\n")
	sb.WriteString("- [← Back to Overview](../README.md)\n")

	return sb.String()
}

// generateCircleCIDiagram creates a simple Mermaid diagram for CircleCI workflows
func generateCircleCIDiagram(analysis *Analysis) string {
	if analysis.TotalJobs == 0 {
		return ""
	}
	
	diagram := &shared.MermaidDiagram{
		Title: "CircleCI Pipeline",
	}
	
	// Get main jobs (limit to most important ones)
	jobNames := GetAllJobNames(analysis.Config)
	mostUsed := GetMostUsedJobs(analysis)
	
	// Prioritize frequently used jobs, fall back to all jobs
	var selectedJobs []string
	if len(mostUsed) > 0 {
		for _, jobInfo := range mostUsed {
			selectedJobs = append(selectedJobs, jobInfo.Name)
			if len(selectedJobs) >= 8 { // Limit for readability
				break
			}
		}
	} else {
		// Use first 8 jobs if no usage data
		for i, jobName := range jobNames {
			if i >= 8 {
				break
			}
			selectedJobs = append(selectedJobs, jobName)
		}
	}
	
	// Create nodes for selected jobs
	for _, jobName := range selectedJobs {
		if job, exists := analysis.Config.Jobs[jobName]; exists {
			commands := ExtractCommands(job.Steps)
			if len(commands) > 5 {
				commands = commands[:5]
			}
			
			var executor string
			if job.Docker != nil && len(job.Docker) > 0 {
				executor = fmt.Sprintf("Docker: %s", job.Docker[0].Image)
			} else if job.Machine != nil {
				executor = "Machine executor"
			} else if job.Executor != "" {
				executor = fmt.Sprintf("Executor: %s", job.Executor)
			}
			
			node := shared.MermaidNode{
				ID:          shared.CleanNodeID(jobName),
				Label:       jobName,
				Description: executor,
				Commands:    commands,
				NodeType:    shared.ClassifyNodeType(jobName, commands),
			}
			diagram.Nodes = append(diagram.Nodes, node)
		}
	}
	
	// Add workflow dependencies if available
	if analysis.TotalWorkflows > 0 {
		for _, workflow := range analysis.Config.Workflows {
			for _, jobInterface := range workflow.Jobs {
				// Handle different job formats in workflows
				var jobName string
				var requires []string
				
				switch job := jobInterface.(type) {
				case string:
					jobName = job
				case map[string]interface{}:
					// Extract job name (first key)
					for name := range job {
						jobName = name
						if reqsInterface, hasReqs := job[name].(map[string]interface{}); hasReqs {
							if reqsArray, hasRequires := reqsInterface["requires"]; hasRequires {
								if reqsSlice, ok := reqsArray.([]interface{}); ok {
									for _, req := range reqsSlice {
										if reqStr, ok := req.(string); ok {
											requires = append(requires, reqStr)
										}
									}
								}
							}
						}
						break
					}
				}
				
				if jobName != "" {
					// Find target job node
					var toNode *shared.MermaidNode
					for i := range diagram.Nodes {
						if diagram.Nodes[i].Label == jobName {
							toNode = &diagram.Nodes[i]
							break
						}
					}
					
					// Add dependency edges
					if toNode != nil && len(requires) > 0 {
						for _, reqJob := range requires {
							for i := range diagram.Nodes {
								if diagram.Nodes[i].Label == reqJob {
									edge := shared.MermaidEdge{
										From: diagram.Nodes[i].ID,
										To:   toNode.ID,
									}
									diagram.Edges = append(diagram.Edges, edge)
									break
								}
							}
						}
					}
				}
			}
		}
	}
	
	// If no workflow dependencies, create a simple sequential flow
	if len(diagram.Edges) == 0 && len(diagram.Nodes) > 1 {
		for i := 0; i < len(diagram.Nodes)-1; i++ {
			edge := shared.MermaidEdge{
				From: diagram.Nodes[i].ID,
				To:   diagram.Nodes[i+1].ID,
			}
			diagram.Edges = append(diagram.Edges, edge)
		}
	}
	
	return diagram.Generate()
}
