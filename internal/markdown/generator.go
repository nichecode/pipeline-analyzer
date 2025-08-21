package markdown

import (
	"fmt"
	"sort"
	"strings"

	"github.com/nichecode/pipeline-analyzer/internal/circleci"
)

// GenerateAllJobsIndex generates the all-jobs.md summary
func GenerateAllJobsIndex(analysis *circleci.Analysis) string {
	var sb strings.Builder

	sb.WriteString("# All Jobs Index\n\n")
	sb.WriteString(fmt.Sprintf("Total jobs found: **%d**\n\n", analysis.TotalJobs))

	jobNames := circleci.GetAllJobNames(analysis.Config)
	sort.Strings(jobNames)

	sb.WriteString("| Job Name | Description | Usage Count | Dependencies |\n")
	sb.WriteString("|----------|-------------|-------------|---------------|\n")

	for _, jobName := range jobNames {
		job := analysis.Config.Jobs[jobName]
		description := job.Description
		if description == "" {
			description = "*No description*"
		}
		
		usageCount := analysis.JobUsage[jobName]
		if usageCount == 0 {
			usageCount = 0 // Explicit for unused jobs
		}

		deps := "None"
		if jobDeps, hasDeps := analysis.JobDependencies[jobName]; hasDeps && len(jobDeps) > 0 {
			deps = strings.Join(jobDeps, ", ")
		}

		normalizedName := circleci.NormalizeJobName(jobName)
		jobLink := fmt.Sprintf("[%s](../jobs/%s.md)", jobName, normalizedName)
		
		sb.WriteString(fmt.Sprintf("| %s | %s | %d | %s |\n", 
			jobLink, description, usageCount, deps))
	}

	sb.WriteString("\n## Navigation\n\n")
	sb.WriteString("- [← Back to Overview](../README.md)\n")
	sb.WriteString("- [Job Usage Analysis](job-usage.md)\n")

	return sb.String()
}

// GenerateJobUsageAnalysis generates job-usage.md
func GenerateJobUsageAnalysis(analysis *circleci.Analysis) string {
	var sb strings.Builder

	sb.WriteString("# Job Usage Analysis\n\n")

	// Most used jobs
	sb.WriteString("## Most Frequently Used Jobs\n\n")
	mostUsedJobs := circleci.GetMostUsedJobs(analysis)
	
	if len(mostUsedJobs) > 0 {
		sb.WriteString("| Rank | Job Name | Usage Count | Link |\n")
		sb.WriteString("|------|----------|-------------|------|\n")
		
		for i, job := range mostUsedJobs {
			normalizedName := circleci.NormalizeJobName(job.Name)
			jobLink := fmt.Sprintf("[%s](../jobs/%s.md)", job.Name, normalizedName)
			sb.WriteString(fmt.Sprintf("| %d | %s | %d | %s |\n", 
				i+1, job.Name, job.Count, jobLink))
		}
	} else {
		sb.WriteString("No job usage data found.\n")
	}
	sb.WriteString("\n")

	// Job Dependencies
	sb.WriteString("## Job Dependencies\n\n")
	if len(analysis.JobDependencies) > 0 {
		sb.WriteString("Jobs with dependencies:\n\n")
		sb.WriteString("| Job | Dependencies |\n")
		sb.WriteString("|-----|-------------|\n")
		
		var depJobs []string
		for job := range analysis.JobDependencies {
			depJobs = append(depJobs, job)
		}
		sort.Strings(depJobs)
		
		for _, job := range depJobs {
			deps := analysis.JobDependencies[job]
			normalizedName := circleci.NormalizeJobName(job)
			jobLink := fmt.Sprintf("[%s](../jobs/%s.md)", job, normalizedName)
			sb.WriteString(fmt.Sprintf("| %s | %s |\n", jobLink, strings.Join(deps, ", ")))
		}
	} else {
		sb.WriteString("No job dependencies found.\n")
	}
	sb.WriteString("\n")

	// Unused jobs
	sb.WriteString("## Unused Jobs\n\n")
	unusedJobs := []string{}
	for jobName, count := range analysis.JobUsage {
		if count == 0 {
			unusedJobs = append(unusedJobs, jobName)
		}
	}
	
	// Also check for jobs not in any workflow
	allJobs := circleci.GetAllJobNames(analysis.Config)
	for _, jobName := range allJobs {
		if _, exists := analysis.JobUsage[jobName]; !exists {
			unusedJobs = append(unusedJobs, jobName)
		}
	}
	
	if len(unusedJobs) > 0 {
		sort.Strings(unusedJobs)
		sb.WriteString("Jobs that are not used in any workflow:\n\n")
		for _, job := range unusedJobs {
			normalizedName := circleci.NormalizeJobName(job)
			jobLink := fmt.Sprintf("[%s](../jobs/%s.md)", job, normalizedName)
			sb.WriteString(fmt.Sprintf("- %s\n", jobLink))
		}
	} else {
		sb.WriteString("All jobs are used in workflows.\n")
	}
	sb.WriteString("\n")

	sb.WriteString("## Navigation\n\n")
	sb.WriteString("- [← Back to Overview](../README.md)\n")
	sb.WriteString("- [All Jobs Index](all-jobs.md)\n")

	return sb.String()
}

// GenerateCommandsAnalysis generates commands.md
func GenerateCommandsAnalysis(analysis *circleci.Analysis) string {
	var sb strings.Builder

	sb.WriteString("# Commands Analysis\n\n")

	// Pattern frequency
	sb.WriteString("## Command Pattern Frequency\n\n")
	patterns := circleci.GetMostUsedPatterns(analysis)
	
	if len(patterns) > 0 {
		sb.WriteString("| Pattern | Occurrences | Jobs Using |\n")
		sb.WriteString("|---------|-------------|------------|\n")
		
		for _, pattern := range patterns {
			jobCount := len(pattern.Jobs)
			sb.WriteString(fmt.Sprintf("| %s | %d | %d jobs |\n", 
				pattern.Pattern, pattern.Count, jobCount))
		}
	} else {
		sb.WriteString("No command patterns detected.\n")
	}
	sb.WriteString("\n")

	// Command frequency by first word
	sb.WriteString("## Most Common Commands\n\n")
	commandFreq := circleci.GetCommandFrequency(analysis.Config)
	
	type cmdFreq struct {
		Command string
		Count   int
	}
	
	var frequencies []cmdFreq
	for cmd, count := range commandFreq {
		frequencies = append(frequencies, cmdFreq{Command: cmd, Count: count})
	}
	
	sort.Slice(frequencies, func(i, j int) bool {
		return frequencies[i].Count > frequencies[j].Count
	})
	
	if len(frequencies) > 0 {
		sb.WriteString("| Command | Usage Count |\n")
		sb.WriteString("|---------|-------------|\n")
		
		// Show top 20
		limit := 20
		if len(frequencies) < limit {
			limit = len(frequencies)
		}
		
		for i := 0; i < limit; i++ {
			freq := frequencies[i]
			sb.WriteString(fmt.Sprintf("| `%s` | %d |\n", freq.Command, freq.Count))
		}
	} else {
		sb.WriteString("No commands found.\n")
	}
	sb.WriteString("\n")

	// Jobs by pattern
	if len(patterns) > 0 {
		sb.WriteString("## Jobs by Pattern\n\n")
		for _, pattern := range patterns {
			if len(pattern.Jobs) > 0 {
				sb.WriteString(fmt.Sprintf("### %s (%d jobs)\n\n", pattern.Pattern, len(pattern.Jobs)))
				
				sort.Strings(pattern.Jobs)
				for _, job := range pattern.Jobs {
					normalizedName := circleci.NormalizeJobName(job)
					jobLink := fmt.Sprintf("[%s](../jobs/%s.md)", job, normalizedName)
					sb.WriteString(fmt.Sprintf("- %s\n", jobLink))
				}
				sb.WriteString("\n")
			}
		}
	}

	sb.WriteString("## Navigation\n\n")
	sb.WriteString("- [← Back to Overview](../README.md)\n")
	sb.WriteString("- [Docker & Scripts](docker-and-scripts.md)\n")

	return sb.String()
}

// GenerateDockerAndScriptsAnalysis generates docker-and-scripts.md
func GenerateDockerAndScriptsAnalysis(analysis *circleci.Analysis) string {
	var sb strings.Builder

	sb.WriteString("# Docker & Scripts Analysis\n\n")

	// Docker usage stats
	dockerJobs, otherJobs := circleci.CountDockerUsage(analysis.Config)
	sb.WriteString("## Docker Usage Overview\n\n")
	sb.WriteString(fmt.Sprintf("- **Jobs using Docker:** %d\n", dockerJobs))
	sb.WriteString(fmt.Sprintf("- **Jobs using other executors:** %d\n", otherJobs))
	sb.WriteString("\n")

	// Docker patterns
	dockerPattern := circleci.GetJobsByPattern(analysis, "docker")
	dockerComposePattern := circleci.GetJobsByPattern(analysis, "docker-compose")
	
	if len(dockerPattern) > 0 {
		sb.WriteString("## Jobs Using Docker Commands\n\n")
		sort.Strings(dockerPattern)
		for _, job := range dockerPattern {
			normalizedName := circleci.NormalizeJobName(job)
			jobLink := fmt.Sprintf("[%s](../jobs/%s.md)", job, normalizedName)
			sb.WriteString(fmt.Sprintf("- %s\n", jobLink))
		}
		sb.WriteString("\n")
	}

	if len(dockerComposePattern) > 0 {
		sb.WriteString("## Jobs Using Docker Compose\n\n")
		sort.Strings(dockerComposePattern)
		for _, job := range dockerComposePattern {
			normalizedName := circleci.NormalizeJobName(job)
			jobLink := fmt.Sprintf("[%s](../jobs/%s.md)", job, normalizedName)
			sb.WriteString(fmt.Sprintf("- %s\n", jobLink))
		}
		sb.WriteString("\n")
	}

	// Script patterns
	scriptPattern := circleci.GetJobsByPattern(analysis, "./")
	if len(scriptPattern) > 0 {
		sb.WriteString("## Jobs Running Local Scripts\n\n")
		sort.Strings(scriptPattern)
		for _, job := range scriptPattern {
			normalizedName := circleci.NormalizeJobName(job)
			jobLink := fmt.Sprintf("[%s](../jobs/%s.md)", job, normalizedName)
			sb.WriteString(fmt.Sprintf("- %s\n", jobLink))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("## Navigation\n\n")
	sb.WriteString("- [← Back to Overview](../README.md)\n")
	sb.WriteString("- [Commands Analysis](commands.md)\n")
	sb.WriteString("- [Executors & Images](executors-and-images.md)\n")

	return sb.String()
}

// GenerateExecutorsAndImagesAnalysis generates executors-and-images.md
func GenerateExecutorsAndImagesAnalysis(analysis *circleci.Analysis) string {
	var sb strings.Builder

	sb.WriteString("# Executors & Images Analysis\n\n")

	if len(analysis.ExecutorUsage) > 0 {
		sb.WriteString("## Image/Executor Usage\n\n")
		sb.WriteString("| Image/Executor | Jobs Using |\n")
		sb.WriteString("|----------------|------------|\n")
		
		var executors []string
		for executor := range analysis.ExecutorUsage {
			executors = append(executors, executor)
		}
		sort.Strings(executors)
		
		for _, executor := range executors {
			jobs := analysis.ExecutorUsage[executor]
			jobCount := len(jobs)
			sb.WriteString(fmt.Sprintf("| `%s` | %d jobs |\n", executor, jobCount))
		}
		sb.WriteString("\n")

		// Detailed breakdown
		sb.WriteString("## Detailed Job Assignments\n\n")
		for _, executor := range executors {
			jobs := analysis.ExecutorUsage[executor]
			if len(jobs) > 0 {
				sb.WriteString(fmt.Sprintf("### %s\n\n", executor))
				sort.Strings(jobs)
				for _, job := range jobs {
					normalizedName := circleci.NormalizeJobName(job)
					jobLink := fmt.Sprintf("[%s](../jobs/%s.md)", job, normalizedName)
					sb.WriteString(fmt.Sprintf("- %s\n", jobLink))
				}
				sb.WriteString("\n")
			}
		}
	} else {
		sb.WriteString("No executor or image usage data found.\n\n")
	}

	sb.WriteString("## Navigation\n\n")
	sb.WriteString("- [← Back to Overview](../README.md)\n")
	sb.WriteString("- [Docker & Scripts](docker-and-scripts.md)\n")

	return sb.String()
}

// GenerateWorkflowsIndex generates workflows.md
func GenerateWorkflowsIndex(analysis *circleci.Analysis) string {
	var sb strings.Builder

	sb.WriteString("# Workflows Index\n\n")
	sb.WriteString(fmt.Sprintf("Total workflows found: **%d**\n\n", analysis.TotalWorkflows))

	if analysis.TotalWorkflows > 0 {
		workflowNames := circleci.GetAllWorkflowNames(analysis.Config)
		sort.Strings(workflowNames)

		sb.WriteString("| Workflow Name | Jobs Count | Link |\n")
		sb.WriteString("|---------------|------------|------|\n")

		for _, workflowName := range workflowNames {
			workflow := analysis.Config.Workflows[workflowName]
			jobs := circleci.ExtractWorkflowJobs(workflow)
			jobCount := len(jobs)
			
			normalizedName := circleci.NormalizeJobName(workflowName)
			workflowLink := fmt.Sprintf("[%s](../workflows/%s.md)", workflowName, normalizedName)
			
			sb.WriteString(fmt.Sprintf("| %s | %d | %s |\n", 
				workflowName, jobCount, workflowLink))
		}
	} else {
		sb.WriteString("No workflows found in the configuration.\n")
	}
	sb.WriteString("\n")

	sb.WriteString("## Navigation\n\n")
	sb.WriteString("- [← Back to Overview](../README.md)\n")
	sb.WriteString("- [Job Usage Analysis](job-usage.md)\n")

	return sb.String()
}