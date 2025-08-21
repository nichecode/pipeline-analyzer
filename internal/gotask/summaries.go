package gotask

import (
	"fmt"
	"sort"
	"strings"

	"github.com/nichecode/pipeline-analyzer/internal/shared"
)

// GenerateAllTasksIndex generates the all-tasks.md summary
func GenerateAllTasksIndex(analysis *Analysis) string {
	var sb strings.Builder

	sb.WriteString("# All Tasks Index\n\n")
	sb.WriteString(fmt.Sprintf("Total tasks found: **%d**\n\n", analysis.TotalTasks))

	taskNames := GetAllTaskNames(analysis.Taskfile)
	sort.Strings(taskNames)

	sb.WriteString("| Task Name | Description | Used By | Type | Optimized |\n")
	sb.WriteString("|-----------|-------------|---------|------|----------|\n")

	for _, taskName := range taskNames {
		task := analysis.Taskfile.Tasks[taskName]
		description := task.Desc
		if description == "" {
			description = "*No description*"
		}
		
		usageCount := analysis.TaskUsage[taskName]
		taskType := DetectTaskType(task, taskName)
		optimized := "❌"
		if IsOptimizedForCaching(task) {
			optimized = "✅"
		}

		normalizedName := NormalizeTaskName(taskName)
		taskLink := fmt.Sprintf("[%s](../tasks/%s.md)", taskName, normalizedName)
		
		sb.WriteString(fmt.Sprintf("| %s | %s | %d | %s | %s |\n", 
			taskLink, shared.TruncateString(description, 50), usageCount, taskType, optimized))
	}

	sb.WriteString("\n## Task Categories\n\n")
	
	// Group tasks by type
	taskTypes := make(map[string][]string)
	for taskName, task := range analysis.Taskfile.Tasks {
		taskType := DetectTaskType(task, taskName)
		taskTypes[taskType] = append(taskTypes[taskType], taskName)
	}

	for taskType, tasks := range taskTypes {
		sort.Strings(tasks)
		sb.WriteString(fmt.Sprintf("### %s (%d tasks)\n\n", strings.Title(taskType), len(tasks)))
		for _, task := range tasks {
			normalizedName := NormalizeTaskName(task)
			sb.WriteString(fmt.Sprintf("- [%s](../tasks/%s.md)\n", task, normalizedName))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("## Navigation\n\n")
	sb.WriteString("- [← Back to Overview](../README.md)\n")
	sb.WriteString("- [Task Usage Analysis](task-usage.md)\n")
	sb.WriteString("- [Dependency Graph](../tasks/dependency-graph.md)\n")

	return sb.String()
}

// GenerateTaskUsageAnalysis generates task-usage.md
func GenerateTaskUsageAnalysis(analysis *Analysis) string {
	var sb strings.Builder

	sb.WriteString("# Task Usage Analysis\n\n")

	// Most used tasks
	sb.WriteString("## Most Frequently Used Tasks\n\n")
	mostUsedTasks := GetMostUsedTasks(analysis)
	
	if len(mostUsedTasks) > 0 {
		sb.WriteString("| Rank | Task Name | Used By | Link |\n")
		sb.WriteString("|------|-----------|---------|------|\n")
		
		for i, task := range mostUsedTasks {
			normalizedName := NormalizeTaskName(task.Name)
			taskLink := fmt.Sprintf("[%s](../tasks/%s.md)", task.Name, normalizedName)
			sb.WriteString(fmt.Sprintf("| %d | %s | %d tasks | %s |\n", 
				i+1, task.Name, task.Count, taskLink))
		}
	} else {
		sb.WriteString("No task dependencies found.\n")
	}
	sb.WriteString("\n")

	// Task Dependencies
	sb.WriteString("## Task Dependencies\n\n")
	if len(analysis.TaskDependencies) > 0 {
		sb.WriteString("Tasks with dependencies:\n\n")
		sb.WriteString("| Task | Dependencies |\n")
		sb.WriteString("|------|-------------|\n")
		
		var depTasks []string
		for task := range analysis.TaskDependencies {
			depTasks = append(depTasks, task)
		}
		sort.Strings(depTasks)
		
		for _, task := range depTasks {
			deps := analysis.TaskDependencies[task]
			if len(deps) > 0 {
				normalizedName := NormalizeTaskName(task)
				taskLink := fmt.Sprintf("[%s](../tasks/%s.md)", task, normalizedName)
				
				var depLinks []string
				for _, dep := range deps {
					depNormalized := NormalizeTaskName(dep)
					depLinks = append(depLinks, fmt.Sprintf("[%s](../tasks/%s.md)", dep, depNormalized))
				}
				
				sb.WriteString(fmt.Sprintf("| %s | %s |\n", taskLink, strings.Join(depLinks, ", ")))
			}
		}
	} else {
		sb.WriteString("No task dependencies found.\n")
	}
	sb.WriteString("\n")

	// Circular Dependencies
	if len(analysis.CircularDeps) > 0 {
		sb.WriteString("## ⚠️ Circular Dependencies\n\n")
		sb.WriteString("**CRITICAL:** The following circular dependencies must be resolved:\n\n")
		
		for i, cycle := range analysis.CircularDeps {
			sb.WriteString(fmt.Sprintf("### Cycle %d\n\n", i+1))
			sb.WriteString("```\n")
			for j, task := range cycle {
				if j > 0 {
					sb.WriteString(" → ")
				}
				sb.WriteString(task)
			}
			sb.WriteString("\n```\n\n")
		}
	}

	// Independent tasks (potential parallel execution)
	independentTasks := []string{}
	for taskName := range analysis.Taskfile.Tasks {
		if len(analysis.TaskDependencies[taskName]) == 0 {
			independentTasks = append(independentTasks, taskName)
		}
	}
	
	if len(independentTasks) > 0 {
		sb.WriteString("## 🚀 Independent Tasks\n\n")
		sb.WriteString("These tasks have no dependencies and can run in parallel:\n\n")
		sort.Strings(independentTasks)
		for _, task := range independentTasks {
			normalizedName := NormalizeTaskName(task)
			sb.WriteString(fmt.Sprintf("- [%s](../tasks/%s.md)\n", task, normalizedName))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("## Navigation\n\n")
	sb.WriteString("- [← Back to Overview](../README.md)\n")
	sb.WriteString("- [All Tasks Index](all-tasks.md)\n")
	sb.WriteString("- [Dependency Graph](../tasks/dependency-graph.md)\n")

	return sb.String()
}

// GenerateCommandsAnalysis generates commands.md
func GenerateCommandsAnalysis(analysis *Analysis) string {
	var sb strings.Builder

	sb.WriteString("# Command Analysis\n\n")

	// Pattern frequency
	sb.WriteString("## Command Pattern Frequency\n\n")
	patterns := GetMostUsedPatterns(analysis)
	
	if len(patterns) > 0 {
		sb.WriteString("| Pattern | Occurrences | Tasks Using | Category |\n")
		sb.WriteString("|---------|-------------|-------------|-----------|\n")
		
		commonPatterns := shared.GetCommonPatterns()
		
		for _, pattern := range patterns {
			taskCount := len(pattern.Tasks)
			category := "utility"
			if commonPattern, exists := commonPatterns[pattern.Pattern]; exists {
				category = commonPattern.Category
			}
			sb.WriteString(fmt.Sprintf("| %s | %d | %d | %s |\n", 
				pattern.Pattern, pattern.Count, taskCount, category))
		}
	} else {
		sb.WriteString("No command patterns detected.\n")
	}
	sb.WriteString("\n")

	// Command frequency by first word
	sb.WriteString("## Most Common Commands\n\n")
	commandFreq := GetCommandFrequency(analysis.Taskfile)
	
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
		sb.WriteString("| Command | Usage Count | Description |\n")
		sb.WriteString("|---------|-------------|-------------|\n")
		
		// Show top 20
		limit := 20
		if len(frequencies) < limit {
			limit = len(frequencies)
		}
		
		commonPatterns := shared.GetCommonPatterns()
		
		for i := 0; i < limit; i++ {
			freq := frequencies[i]
			description := "Shell command"
			
			// Try to find description from common patterns
			for _, pattern := range commonPatterns {
				if pattern.Regex.MatchString(freq.Command) {
					description = pattern.Description
					break
				}
			}
			
			sb.WriteString(fmt.Sprintf("| `%s` | %d | %s |\n", freq.Command, freq.Count, description))
		}
	} else {
		sb.WriteString("No commands found.\n")
	}
	sb.WriteString("\n")

	// Tool ecosystem analysis
	allCommands := []string{}
	for _, task := range analysis.Taskfile.Tasks {
		allCommands = append(allCommands, ExtractTaskCommands(task)...)
	}
	
	ecosystem := shared.DetectToolEcosystem(allCommands)
	sb.WriteString("## 🛠️ Tool Ecosystem\n\n")
	sb.WriteString(fmt.Sprintf("**Primary ecosystem detected:** %s\n\n", strings.Title(ecosystem)))

	// Tasks by pattern
	if len(patterns) > 0 {
		sb.WriteString("## Tasks by Pattern\n\n")
		
		// Group by category
		patternsByCategory := make(map[string][]struct {
			Pattern string
			Count   int
			Tasks   []string
		})
		
		commonPatterns := shared.GetCommonPatterns()
		
		for _, pattern := range patterns {
			category := "utility"
			if commonPattern, exists := commonPatterns[pattern.Pattern]; exists {
				category = commonPattern.Category
			}
			patternsByCategory[category] = append(patternsByCategory[category], pattern)
		}
		
		for category, categoryPatterns := range patternsByCategory {
			sb.WriteString(fmt.Sprintf("### %s\n\n", strings.Title(category)))
			
			for _, pattern := range categoryPatterns {
				if len(pattern.Tasks) > 0 {
					sb.WriteString(fmt.Sprintf("**%s** (%d tasks):\n", pattern.Pattern, len(pattern.Tasks)))
					sort.Strings(pattern.Tasks)
					for _, task := range pattern.Tasks {
						normalizedName := NormalizeTaskName(task)
						taskLink := fmt.Sprintf("[%s](../tasks/%s.md)", task, normalizedName)
						sb.WriteString(fmt.Sprintf("- %s\n", taskLink))
					}
					sb.WriteString("\n")
				}
			}
		}
	}

	sb.WriteString("## Navigation\n\n")
	sb.WriteString("- [← Back to Overview](../README.md)\n")
	sb.WriteString("- [Performance Metrics](performance.md)\n")
	sb.WriteString("- [Variable Analysis](variables.md)\n")

	return sb.String()
}

// GeneratePerformanceAnalysis generates performance.md
func GeneratePerformanceAnalysis(analysis *Analysis) string {
	var sb strings.Builder

	sb.WriteString("# Performance Analysis\n\n")

	metrics := GetPerformanceMetrics(analysis.Taskfile)

	sb.WriteString("## 📊 Performance Metrics\n\n")
	sb.WriteString(fmt.Sprintf("- **Total tasks:** %d\n", analysis.TotalTasks))
	sb.WriteString(fmt.Sprintf("- **Tasks with source tracking:** %d (%.1f%%)\n", 
		metrics.TasksWithSources, float64(metrics.TasksWithSources)/float64(analysis.TotalTasks)*100))
	sb.WriteString(fmt.Sprintf("- **Tasks with output tracking:** %d (%.1f%%)\n", 
		metrics.TasksWithGenerates, float64(metrics.TasksWithGenerates)/float64(analysis.TotalTasks)*100))
	sb.WriteString(fmt.Sprintf("- **Tasks with full caching:** %d (%.1f%%)\n", 
		metrics.TasksWithCaching, float64(metrics.TasksWithCaching)/float64(analysis.TotalTasks)*100))
	sb.WriteString(fmt.Sprintf("- **Independent tasks:** %d\n", metrics.ParallelizableTasks))
	sb.WriteString(fmt.Sprintf("- **Optimization potential:** %.1f%%\n\n", metrics.OptimizationPotential))

	// Caching optimization opportunities
	sb.WriteString("## ⚡ Caching Opportunities\n\n")
	
	unoptimizedTasks := []string{}
	for taskName, task := range analysis.Taskfile.Tasks {
		if !IsOptimizedForCaching(task) && GetTaskComplexity(task) > 1 {
			unoptimizedTasks = append(unoptimizedTasks, taskName)
		}
	}

	if len(unoptimizedTasks) > 0 {
		sb.WriteString("The following tasks could benefit from caching optimization:\n\n")
		sb.WriteString("| Task | Complexity | Current Status | Recommendation |\n")
		sb.WriteString("|------|------------|----------------|----------------|\n")
		
		sort.Strings(unoptimizedTasks)
		for _, taskName := range unoptimizedTasks {
			task := analysis.Taskfile.Tasks[taskName]
			complexity := GetTaskComplexity(task)
			
			status := "No tracking"
			if HasSources(task) {
				status = "Sources only"
			}
			if HasGenerates(task) {
				status = "Outputs only"
			}
			
			recommendation := "Add sources + generates"
			if HasSources(task) {
				recommendation = "Add generates"
			} else if HasGenerates(task) {
				recommendation = "Add sources"
			}
			
			normalizedName := NormalizeTaskName(taskName)
			taskLink := fmt.Sprintf("[%s](../tasks/%s.md)", taskName, normalizedName)
			
			sb.WriteString(fmt.Sprintf("| %s | %d | %s | %s |\n", 
				taskLink, complexity, status, recommendation))
		}
		sb.WriteString("\n")
	} else {
		sb.WriteString("✅ All complex tasks are already optimized for caching!\n\n")
	}

	// Performance by task type
	sb.WriteString("## 📊 Performance by Task Type\n\n")
	
	taskTypes := make(map[string]struct {
		total     int
		optimized int
	})
	
	for taskName, task := range analysis.Taskfile.Tasks {
		taskType := DetectTaskType(task, taskName)
		stats := taskTypes[taskType]
		stats.total++
		if IsOptimizedForCaching(task) {
			stats.optimized++
		}
		taskTypes[taskType] = stats
	}
	
	sb.WriteString("| Task Type | Total | Optimized | Percentage |\n")
	sb.WriteString("|-----------|-------|-----------|------------|\n")
	
	for taskType, stats := range taskTypes {
		percentage := float64(stats.optimized) / float64(stats.total) * 100
		sb.WriteString(fmt.Sprintf("| %s | %d | %d | %.1f%% |\n", 
			strings.Title(taskType), stats.total, stats.optimized, percentage))
	}
	sb.WriteString("\n")

	// Parallelization analysis
	sb.WriteString("## 🚀 Parallelization Analysis\n\n")
	
	graph := BuildDependencyGraph(analysis)
	levelGroups := make(map[int][]string)
	for task, level := range graph.Levels {
		levelGroups[level] = append(levelGroups[level], task)
	}
	
	var levels []int
	for level := range levelGroups {
		levels = append(levels, level)
	}
	sort.Ints(levels)
	
	sb.WriteString("Tasks that can run in parallel (by dependency level):\n\n")
	sb.WriteString("| Level | Tasks | Can Parallelize |\n")
	sb.WriteString("|-------|-------|-----------------|\n")
	
	for _, level := range levels {
		tasks := levelGroups[level]
		canParallelize := len(tasks) > 1
		parallelText := "❌"
		if canParallelize {
			parallelText = "✅"
		}
		sb.WriteString(fmt.Sprintf("| %d | %d | %s |\n", level, len(tasks), parallelText))
	}
	sb.WriteString("\n")

	// Performance recommendations
	sb.WriteString("## 💡 Performance Recommendations\n\n")
	
	recommendations := []string{}
	
	if metrics.OptimizationPotential > 50 {
		recommendations = append(recommendations, 
			"🔴 **High Priority**: Over 50% of tasks lack caching optimization")
	} else if metrics.OptimizationPotential > 25 {
		recommendations = append(recommendations, 
			"🟡 **Medium Priority**: Consider optimizing remaining tasks for caching")
	}
	
	if metrics.ParallelizableTasks > analysis.TotalTasks/2 {
		recommendations = append(recommendations, 
			"🚀 **Parallelization**: Many tasks can run independently - consider grouping related tasks")
	}
	
	if len(analysis.CriticalPath) > 5 {
		recommendations = append(recommendations, 
			"⏰ **Critical Path**: Long dependency chain detected - focus optimization on critical path tasks")
	}
	
	if len(recommendations) > 0 {
		for _, rec := range recommendations {
			sb.WriteString(fmt.Sprintf("- %s\n", rec))
		}
		sb.WriteString("\n")
	} else {
		sb.WriteString("✅ Your Taskfile is well-optimized for performance!\n\n")
	}

	sb.WriteString("## Navigation\n\n")
	sb.WriteString("- [← Back to Overview](../README.md)\n")
	sb.WriteString("- [Optimization Guide](../OPTIMIZATION-GUIDE.md)\n")
	sb.WriteString("- [Dependency Graph](../tasks/dependency-graph.md)\n")

	return sb.String()
}

// GenerateVariableAnalysis generates variables.md
func GenerateVariableAnalysis(analysis *Analysis) string {
	var sb strings.Builder

	sb.WriteString("# Variable Analysis\n\n")

	allVars := ExtractAllVariables(analysis.Taskfile)

	sb.WriteString(fmt.Sprintf("**Total variables found:** %d\n\n", len(allVars)))

	// Global variables
	globalVars := []Variable{}
	for _, variable := range allVars {
		if variable.IsGlobal {
			globalVars = append(globalVars, variable)
		}
	}

	if len(globalVars) > 0 {
		sb.WriteString("## 🌍 Global Variables\n\n")
		sb.WriteString("| Name | Type | Value | Environment |\n")
		sb.WriteString("|------|------|-------|-------------|\n")
		
		var globalVarNames []string
		for name, variable := range allVars {
			if variable.IsGlobal {
				globalVarNames = append(globalVarNames, name)
			}
		}
		sort.Strings(globalVarNames)
		
		for _, name := range globalVarNames {
			variable := allVars[name]
			envType := "Variable"
			if variable.IsEnvironment {
				envType = "Environment"
			}
			
			valueStr := fmt.Sprintf("%v", variable.Value)
			if len(valueStr) > 50 {
				valueStr = shared.TruncateString(valueStr, 50)
			}
			
			sb.WriteString(fmt.Sprintf("| `%s` | %s | %s | %s |\n", 
				variable.Name, variable.Type, valueStr, envType))
		}
		sb.WriteString("\n")
	}

	// Task-specific variables
	taskVars := []Variable{}
	for _, variable := range allVars {
		if !variable.IsGlobal {
			taskVars = append(taskVars, variable)
		}
	}

	if len(taskVars) > 0 {
		sb.WriteString("## 📋 Task-Specific Variables\n\n")
		
		// Group by task
		varsByTask := make(map[string][]Variable)
		for _, variable := range taskVars {
			if len(variable.UsedInTasks) > 0 {
				taskName := variable.UsedInTasks[0]
				varsByTask[taskName] = append(varsByTask[taskName], variable)
			}
		}
		
		var taskNames []string
		for taskName := range varsByTask {
			taskNames = append(taskNames, taskName)
		}
		sort.Strings(taskNames)
		
		for _, taskName := range taskNames {
			variables := varsByTask[taskName]
			if len(variables) > 0 {
				normalizedName := NormalizeTaskName(taskName)
				sb.WriteString(fmt.Sprintf("### [%s](../tasks/%s.md)\n\n", taskName, normalizedName))
				
				for _, variable := range variables {
					envType := "Variable"
					if variable.IsEnvironment {
						envType = "Environment"
					}
					
					valueStr := fmt.Sprintf("%v", variable.Value)
					if len(valueStr) > 100 {
						valueStr = shared.TruncateString(valueStr, 100)
					}
					
					sb.WriteString(fmt.Sprintf("- **%s** (%s, %s): `%s`\n", 
						variable.Name, variable.Type, envType, valueStr))
				}
				sb.WriteString("\n")
			}
		}
	}

	// Variable usage patterns
	sb.WriteString("## 📊 Variable Usage Patterns\n\n")
	
	// Type distribution
	typeCount := make(map[string]int)
	for _, variable := range allVars {
		typeCount[variable.Type]++
	}
	
	sb.WriteString("### Variable Types\n\n")
	sb.WriteString("| Type | Count | Percentage |\n")
	sb.WriteString("|------|-------|------------|\n")
	
	var types []string
	for varType := range typeCount {
		types = append(types, varType)
	}
	sort.Strings(types)
	
	totalVars := len(allVars)
	for _, varType := range types {
		count := typeCount[varType]
		percentage := float64(count) / float64(totalVars) * 100
		sb.WriteString(fmt.Sprintf("| %s | %d | %.1f%% |\n", varType, count, percentage))
	}
	sb.WriteString("\n")

	// Environment vs Variables
	envCount := 0
	varCount := 0
	for _, variable := range allVars {
		if variable.IsEnvironment {
			envCount++
		} else {
			varCount++
		}
	}
	
	sb.WriteString("### Environment vs Variables\n\n")
	sb.WriteString("| Type | Count | Percentage |\n")
	sb.WriteString("|------|-------|------------|\n")
	sb.WriteString(fmt.Sprintf("| Variables | %d | %.1f%% |\n", varCount, float64(varCount)/float64(totalVars)*100))
	sb.WriteString(fmt.Sprintf("| Environment | %d | %.1f%% |\n", envCount, float64(envCount)/float64(totalVars)*100))
	sb.WriteString("\n")

	sb.WriteString("## Navigation\n\n")
	sb.WriteString("- [← Back to Overview](../README.md)\n")
	sb.WriteString("- [Command Analysis](commands.md)\n")
	sb.WriteString("- [Performance Analysis](performance.md)\n")

	return sb.String()
}