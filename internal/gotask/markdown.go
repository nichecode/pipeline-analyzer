package gotask

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/nichecode/pipeline-analyzer/internal/shared"
)

// GenerateMainReadme generates the main README.md file for go-task analysis
func GenerateMainReadme(analysis *Analysis, taskfilePath string) string {
	var sb strings.Builder

	sb.WriteString("# Go-Task Analysis Report\n\n")
	sb.WriteString(fmt.Sprintf("**Generated:** %s\n", analysis.GeneratedAt.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("**Taskfile:** %s\n\n", taskfilePath))

	// Overview section
	sb.WriteString("## 📊 Overview\n\n")
	sb.WriteString(fmt.Sprintf("- **Tasks:** %d\n", analysis.TotalTasks))
	sb.WriteString(fmt.Sprintf("- **Includes:** %d\n", analysis.TotalIncludes))
	
	if len(analysis.CircularDeps) > 0 {
		sb.WriteString(fmt.Sprintf("- **Circular Dependencies:** %d detected ⚠️\n", len(analysis.CircularDeps)))
	} else {
		sb.WriteString("- **Circular Dependencies:** None ✅\n")
	}
	
	if len(analysis.CriticalPath) > 0 {
		sb.WriteString(fmt.Sprintf("- **Critical Path Length:** %d tasks\n", len(analysis.CriticalPath)))
	}
	
	// Performance metrics
	metrics := GetPerformanceMetrics(analysis.Taskfile)
	sb.WriteString(fmt.Sprintf("- **Tasks with Caching:** %d/%d (%.1f%%)\n", 
		metrics.TasksWithCaching, analysis.TotalTasks, 
		float64(metrics.TasksWithCaching)/float64(analysis.TotalTasks)*100))
	
	sb.WriteString("\n")

	// Quick Start section
	sb.WriteString("## 🚀 Quick Start\n\n")
	sb.WriteString("1. **[📋 Optimization Guide](optimization-guide.md)** - Performance improvement recommendations\n")
	sb.WriteString("2. **[📈 Task Usage Analysis](summaries/task-usage.md)** - Task dependency patterns\n")
	sb.WriteString("3. **[⚡ Command Analysis](summaries/commands.md)** - Command patterns and tools\n")
	sb.WriteString("4. **[🔗 Dependency Graph](tasks/dependency-graph.md)** - Visual task relationships\n\n")

	// Directory Structure section
	sb.WriteString("## 📁 Directory Structure\n\n")
	
	sb.WriteString("### Tasks\n")
	sb.WriteString("Individual task analysis with commands and optimization opportunities:\n\n")
	taskNames := GetAllTaskNames(analysis.Taskfile)
	sort.Strings(taskNames)
	for _, taskName := range taskNames {
		normalizedName := NormalizeTaskName(taskName)
		sb.WriteString(fmt.Sprintf("- [tasks/%s.md](tasks/%s.md)\n", normalizedName, normalizedName))
	}
	sb.WriteString("\n")

	if analysis.TotalIncludes > 0 {
		sb.WriteString("### Includes\n")
		sb.WriteString("Analysis of included Taskfiles:\n\n")
		includeNames := GetAllIncludeNames(analysis.Taskfile)
		sort.Strings(includeNames)
		for _, includeName := range includeNames {
			normalizedName := NormalizeTaskName(includeName)
			sb.WriteString(fmt.Sprintf("- [includes/%s.md](includes/%s.md)\n", normalizedName, normalizedName))
		}
		sb.WriteString("\n")
	}

	// Analysis Summaries section
	sb.WriteString("### Analysis Summaries\n\n")
	sb.WriteString("- [📈 Task Usage & Dependencies](summaries/task-usage.md)\n")
	sb.WriteString("- [📝 All Tasks Index](summaries/all-tasks.md)\n")
	sb.WriteString("- [⚡ Command Analysis](summaries/commands.md)\n")
	sb.WriteString("- [📊 Performance Metrics](summaries/performance.md)\n")
	sb.WriteString("- [🔍 Variable Analysis](summaries/variables.md)\n")
	if analysis.TotalIncludes > 0 {
		sb.WriteString("- [📁 Include Analysis](summaries/includes.md)\n")
	}
	sb.WriteString("\n")

	// Key Insights section
	sb.WriteString("## 🔍 Key Insights\n\n")
	
	// Most used tasks
	mostUsedTasks := GetMostUsedTasks(analysis)
	if len(mostUsedTasks) > 0 {
		sb.WriteString("### Most Depended-On Tasks\n\n")
		sb.WriteString("| Task | Used By | Link |\n")
		sb.WriteString("|------|---------|------|\n")
		
		limit := 5
		if len(mostUsedTasks) < limit {
			limit = len(mostUsedTasks)
		}
		
		for i := 0; i < limit; i++ {
			task := mostUsedTasks[i]
			normalizedName := NormalizeTaskName(task.Name)
			sb.WriteString(fmt.Sprintf("| %s | %d tasks | [View](tasks/%s.md) |\n", 
				task.Name, task.Count, normalizedName))
		}
		sb.WriteString("\n")
	}

	// Optimization opportunities
	if len(analysis.OptimizationTips) > 0 {
		sb.WriteString("### Top Optimization Opportunities\n\n")
		
		highPriorityTips := shared.FilterItems(analysis.OptimizationTips, func(tip OptimizationTip) bool {
			return tip.Severity == "high" || tip.Severity == "medium"
		})
		
		if len(highPriorityTips) > 0 {
			limit := 3
			if len(highPriorityTips) < limit {
				limit = len(highPriorityTips)
			}
			
			for i := 0; i < limit; i++ {
				tip := highPriorityTips[i]
				severity := "⚠️"
				if tip.Severity == "high" {
					severity = "🔴"
				}
				sb.WriteString(fmt.Sprintf("- %s **%s**: %s\n", severity, tip.Task, tip.Message))
			}
			sb.WriteString("\n")
		}
	}

	// Tool ecosystem
	allCommands := []string{}
	for _, task := range analysis.Taskfile.Tasks {
		allCommands = append(allCommands, ExtractTaskCommands(task)...)
	}
	ecosystem := shared.DetectToolEcosystem(allCommands)
	sb.WriteString(fmt.Sprintf("### Primary Tool Ecosystem: **%s**\n\n", strings.Title(ecosystem)))

	// Next steps
	sb.WriteString("## 🎯 Next Steps\n\n")
	sb.WriteString("1. **Review [Optimization Guide](optimization-guide.md)** for specific improvements\n")
	sb.WriteString("2. **Check [Performance Metrics](summaries/performance.md)** for caching opportunities\n")
	sb.WriteString("3. **Examine [Dependency Graph](tasks/dependency-graph.md)** for parallelization potential\n")
	if len(analysis.CircularDeps) > 0 {
		sb.WriteString("4. **Resolve circular dependencies** shown in the dependency analysis\n")
	}
	sb.WriteString("\n")

	return sb.String()
}

// GenerateOptimizationGuide generates the optimization guide
func GenerateOptimizationGuide(analysis *Analysis) string {
	var sb strings.Builder

	sb.WriteString("# Task Optimization Guide\n\n")
	sb.WriteString(fmt.Sprintf("**Generated:** %s\n\n", analysis.GeneratedAt.Format(time.RFC3339)))

	sb.WriteString("## 🎯 Optimization Overview\n\n")
	sb.WriteString("This guide provides specific recommendations to improve your Taskfile performance, maintainability, and reliability.\n\n")

	// Performance metrics
	metrics := GetPerformanceMetrics(analysis.Taskfile)
	sb.WriteString("### 📊 Current Performance Status\n\n")
	sb.WriteString(fmt.Sprintf("- **Tasks with caching optimization:** %d/%d (%.1f%%)\n", 
		metrics.TasksWithCaching, analysis.TotalTasks,
		float64(metrics.TasksWithCaching)/float64(analysis.TotalTasks)*100))
	sb.WriteString(fmt.Sprintf("- **Tasks with source tracking:** %d\n", metrics.TasksWithSources))
	sb.WriteString(fmt.Sprintf("- **Tasks with output tracking:** %d\n", metrics.TasksWithGenerates))
	sb.WriteString(fmt.Sprintf("- **Parallelizable tasks:** %d\n", metrics.ParallelizableTasks))
	sb.WriteString(fmt.Sprintf("- **Optimization potential:** %.1f%%\n\n", metrics.OptimizationPotential))

	// Categorize optimization tips
	tipsByType := shared.GroupItems(analysis.OptimizationTips, func(tip OptimizationTip) string {
		return tip.Type
	})

	// Priority optimizations
	if tips, exists := tipsByType["dependency"]; exists {
		sb.WriteString("## 🔴 Critical Issues (Fix Immediately)\n\n")
		for _, tip := range tips {
			sb.WriteString(fmt.Sprintf("### %s\n\n", tip.Task))
			sb.WriteString(fmt.Sprintf("**Issue:** %s\n\n", tip.Message))
			sb.WriteString(fmt.Sprintf("**Solution:** %s\n\n", tip.Suggestion))
		}
	}

	// Performance optimizations
	if tips, exists := tipsByType["caching"]; exists {
		sb.WriteString("## ⚡ Performance Optimizations\n\n")
		sb.WriteString("These changes can significantly improve task execution speed:\n\n")
		
		for _, tip := range tips {
			sb.WriteString(fmt.Sprintf("### Task: `%s`\n\n", tip.Task))
			sb.WriteString(fmt.Sprintf("**Opportunity:** %s\n\n", tip.Message))
			sb.WriteString("**Implementation:**\n")
			sb.WriteString("```yaml\n")
			sb.WriteString(fmt.Sprintf("%s:\n", tip.Task))
			sb.WriteString("  desc: \"Task description\"\n")
			sb.WriteString("  sources:\n")
			sb.WriteString("    - \"src/**/*.go\"  # Add relevant source patterns\n")
			sb.WriteString("  generates:\n")
			sb.WriteString("    - \"dist/app\"     # Add output files\n")
			sb.WriteString("  cmds:\n")
			sb.WriteString("    - # existing commands\n")
			sb.WriteString("```\n\n")
		}
	}

	// Documentation improvements
	if tips, exists := tipsByType["documentation"]; exists {
		sb.WriteString("## 📝 Documentation Improvements\n\n")
		sb.WriteString("Better documentation makes your tasks more maintainable:\n\n")
		
		for _, tip := range tips {
			sb.WriteString(fmt.Sprintf("- **%s**: %s\n", tip.Task, tip.Message))
		}
		sb.WriteString("\n")
		
		sb.WriteString("**Quick Fix Template:**\n")
		sb.WriteString("```yaml\n")
		sb.WriteString("your-task:\n")
		sb.WriteString("  desc: \"Brief description of what this task does\"\n")
		sb.WriteString("  summary: |\n")
		sb.WriteString("    Longer explanation of the task purpose,\n")
		sb.WriteString("    any prerequisites, and expected outcomes.\n")
		sb.WriteString("```\n\n")
	}

	// Advanced optimizations
	sb.WriteString("## 🚀 Advanced Optimizations\n\n")
	
	sb.WriteString("### Parallel Execution\n\n")
	sb.WriteString("Tasks without dependencies can run in parallel. Consider grouping related independent tasks:\n\n")
	sb.WriteString("```yaml\n")
	sb.WriteString("test-all:\n")
	sb.WriteString("  desc: \"Run all tests in parallel\"\n")
	sb.WriteString("  deps:\n")
	sb.WriteString("    - task: test-unit\n")
	sb.WriteString("    - task: test-integration\n")
	sb.WriteString("    - task: test-e2e\n")
	sb.WriteString("```\n\n")

	if len(analysis.CriticalPath) > 0 {
		sb.WriteString("### Critical Path Analysis\n\n")
		sb.WriteString("Your longest dependency chain is:\n\n")
		for i, task := range analysis.CriticalPath {
			if i > 0 {
				sb.WriteString(" → ")
			}
			sb.WriteString(fmt.Sprintf("`%s`", task))
		}
		sb.WriteString("\n\n")
		sb.WriteString("Focus optimization efforts on these tasks for maximum impact.\n\n")
	}

	// Best practices
	sb.WriteString("## 💡 Best Practices\n\n")
	
	practices := []struct {
		title       string
		description string
		example     string
	}{
		{
			"Use Specific Source Patterns",
			"Avoid overly broad patterns that might cause unnecessary rebuilds",
			"sources: [\"src/**/*.go\", \"go.mod\", \"go.sum\"]",
		},
		{
			"Output Specific Files",
			"List actual output files rather than directories when possible",
			"generates: [\"dist/app\", \"dist/version.txt\"]",
		},
		{
			"Group Related Tasks",
			"Use task namespaces with includes for better organization",
			"includes: { test: \"./tasks/test.yml\" }",
		},
		{
			"Use Status Checks",
			"Skip tasks when conditions are already met",
			"status: [\"test -f dist/app\"]",
		},
	}

	for _, practice := range practices {
		sb.WriteString(fmt.Sprintf("### %s\n\n", practice.title))
		sb.WriteString(fmt.Sprintf("%s\n\n", practice.description))
		if practice.example != "" {
			sb.WriteString("```yaml\n")
			sb.WriteString(fmt.Sprintf("%s\n", practice.example))
			sb.WriteString("```\n\n")
		}
	}

	// Navigation
	sb.WriteString("## Navigation\n\n")
	sb.WriteString("- [← Back to Overview](README.md)\n")
	sb.WriteString("- [Performance Metrics](summaries/performance.md)\n")
	sb.WriteString("- [Dependency Graph](tasks/dependency-graph.md)\n")

	return sb.String()
}

// GenerateTaskMarkdown generates markdown for a single task
func GenerateTaskMarkdown(taskAnalysis *TaskAnalysis) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# Task: %s\n\n", taskAnalysis.Name))

	// Basic info
	if taskAnalysis.Description != "" {
		sb.WriteString(fmt.Sprintf("**Description:** %s\n\n", taskAnalysis.Description))
	}

	if taskAnalysis.Summary != "" {
		sb.WriteString(fmt.Sprintf("**Summary:** %s\n\n", taskAnalysis.Summary))
	}

	// Task properties
	sb.WriteString("## 📋 Task Properties\n\n")
	sb.WriteString(fmt.Sprintf("- **Used by:** %d other tasks\n", taskAnalysis.UsageCount))
	sb.WriteString(fmt.Sprintf("- **Internal:** %t\n", taskAnalysis.IsInternal))
	sb.WriteString(fmt.Sprintf("- **Watch mode:** %t\n", taskAnalysis.HasWatch))
	
	if len(taskAnalysis.Platforms) > 0 {
		sb.WriteString(fmt.Sprintf("- **Platforms:** %s\n", strings.Join(taskAnalysis.Platforms, ", ")))
	}
	
	if len(taskAnalysis.Aliases) > 0 {
		sb.WriteString(fmt.Sprintf("- **Aliases:** %s\n", strings.Join(taskAnalysis.Aliases, ", ")))
	}
	
	sb.WriteString("\n")

	// Dependencies
	if len(taskAnalysis.Dependencies) > 0 {
		sb.WriteString("## 🔗 Dependencies\n\n")
		sb.WriteString("This task depends on:\n\n")
		for _, dep := range taskAnalysis.Dependencies {
			normalizedName := NormalizeTaskName(dep)
			sb.WriteString(fmt.Sprintf("- [%s](%s.md)\n", dep, normalizedName))
		}
		sb.WriteString("\n")
	}

	// Sources and outputs
	if len(taskAnalysis.Sources) > 0 || len(taskAnalysis.Generates) > 0 {
		sb.WriteString("## 📁 Files\n\n")
		
		if len(taskAnalysis.Sources) > 0 {
			sb.WriteString("**Source files:**\n")
			for _, source := range taskAnalysis.Sources {
				sb.WriteString(fmt.Sprintf("- `%s`\n", source))
			}
			sb.WriteString("\n")
		}
		
		if len(taskAnalysis.Generates) > 0 {
			sb.WriteString("**Generated files:**\n")
			for _, generate := range taskAnalysis.Generates {
				sb.WriteString(fmt.Sprintf("- `%s`\n", generate))
			}
			sb.WriteString("\n")
		}
	}

	// Commands
	if len(taskAnalysis.Commands) > 0 {
		sb.WriteString("## ⚡ Commands\n\n")
		for i, command := range taskAnalysis.Commands {
			if len(taskAnalysis.Commands) > 1 {
				sb.WriteString(fmt.Sprintf("**Command %d:**\n", i+1))
			}
			sb.WriteString("```bash\n")
			sb.WriteString(command)
			sb.WriteString("\n```\n\n")
			
			// Command analysis
			classification := shared.ClassifyCommand(command)
			sb.WriteString(fmt.Sprintf("- **Category:** %s\n", classification.Category))
			sb.WriteString(fmt.Sprintf("- **Complexity:** %d/5\n", classification.Complexity))
			sb.WriteString(fmt.Sprintf("- **Risk level:** %s\n", classification.Risk))
			
			if len(classification.Tools) > 0 {
				sb.WriteString(fmt.Sprintf("- **Tools:** %s\n", strings.Join(classification.Tools, ", ")))
			}
			
			if len(classification.Suggestions) > 0 {
				sb.WriteString("\n**Suggestions:**\n")
				for _, suggestion := range classification.Suggestions {
					sb.WriteString(fmt.Sprintf("- %s\n", suggestion))
				}
			}
			sb.WriteString("\n")
		}
	} else {
		// Show explanation for dependency-only tasks
		if len(taskAnalysis.Dependencies) > 0 {
			sb.WriteString("## ⚡ Commands\n\n")
			sb.WriteString("This is a **dependency-only task** that orchestrates other tasks without running direct commands.\n\n")
			sb.WriteString("**Execution Flow:**\n")
			sb.WriteString("1. Runs all dependency tasks in the correct order\n")
			sb.WriteString("2. Completes when all dependencies finish successfully\n\n")
		} else {
			sb.WriteString("## ⚡ Commands\n\n")
			sb.WriteString("⚠️ **No commands defined** - This task may need implementation or could be unused.\n\n")
		}
	}

	// Preconditions
	if len(taskAnalysis.Preconditions) > 0 {
		sb.WriteString("## ✅ Preconditions\n\n")
		for i, precondition := range taskAnalysis.Preconditions {
			sb.WriteString(fmt.Sprintf("%d. `%s`\n", i+1, precondition))
		}
		sb.WriteString("\n")
	}

	// Variables and environment
	if len(taskAnalysis.Variables) > 0 || len(taskAnalysis.Environment) > 0 {
		sb.WriteString("## 🔧 Configuration\n\n")
		
		if len(taskAnalysis.Variables) > 0 {
			sb.WriteString("**Variables:**\n")
			for name, value := range taskAnalysis.Variables {
				sb.WriteString(fmt.Sprintf("- `%s`: %v\n", name, value))
			}
			sb.WriteString("\n")
		}
		
		if len(taskAnalysis.Environment) > 0 {
			sb.WriteString("**Environment:**\n")
			for name, value := range taskAnalysis.Environment {
				sb.WriteString(fmt.Sprintf("- `%s`: %v\n", name, value))
			}
			sb.WriteString("\n")
		}
	}

	// Pattern usage
	if len(taskAnalysis.Patterns) > 0 {
		sb.WriteString("## 🔍 Command Patterns\n\n")
		var patterns []string
		for pattern := range taskAnalysis.Patterns {
			patterns = append(patterns, pattern)
		}
		sort.Strings(patterns)
		
		for _, pattern := range patterns {
			count := taskAnalysis.Patterns[pattern]
			sb.WriteString(fmt.Sprintf("- **%s:** %d occurrence(s)\n", pattern, count))
		}
		sb.WriteString("\n")
	}

	// Optimization opportunities
	if len(taskAnalysis.OptimizationOps) > 0 {
		sb.WriteString("## 🚀 Optimization Opportunities\n\n")
		for _, opportunity := range taskAnalysis.OptimizationOps {
			sb.WriteString(fmt.Sprintf("- %s\n", opportunity))
		}
		sb.WriteString("\n")
	}

	// Navigation
	sb.WriteString("## Navigation\n\n")
	sb.WriteString("- [← Back to All Tasks](../summaries/all-tasks.md)\n")
	sb.WriteString("- [← Back to Overview](../README.md)\n")
	sb.WriteString("- [Dependency Graph](dependency-graph.md)\n")

	return sb.String()
}

// GenerateDependencyGraph generates a dependency graph visualization
func GenerateDependencyGraph(analysis *Analysis) string {
	var sb strings.Builder

	sb.WriteString("# Task Dependency Graph\n\n")
	sb.WriteString(fmt.Sprintf("**Generated:** %s\n\n", analysis.GeneratedAt.Format(time.RFC3339)))

	// Build dependency graph
	graph := BuildDependencyGraph(analysis)

	sb.WriteString("## 🔗 Dependency Overview\n\n")
	sb.WriteString(fmt.Sprintf("- **Total tasks:** %d\n", len(graph.Tasks)))
	sb.WriteString(fmt.Sprintf("- **Tasks with dependencies:** %d\n", len(graph.Edges)))
	sb.WriteString(fmt.Sprintf("- **Circular dependencies:** %d\n", len(graph.Cycles)))
	sb.WriteString("\n")

	// Visual dependency graph using Mermaid
	sb.WriteString("## 📊 Dependency Visualization\n\n")
	sb.WriteString("```mermaid\n")
	sb.WriteString("graph TD\n")

	// Add all tasks as nodes
	for _, task := range graph.Tasks {
		normalizedName := shared.NormalizeFileName(task)
		sb.WriteString(fmt.Sprintf("    %s[\"%s\"]\n", normalizedName, task))
	}

	// Add dependency edges
	for task, deps := range graph.Edges {
		normalizedTask := shared.NormalizeFileName(task)
		for _, dep := range deps {
			normalizedDep := shared.NormalizeFileName(dep)
			sb.WriteString(fmt.Sprintf("    %s --> %s\n", normalizedDep, normalizedTask))
		}
	}

	sb.WriteString("```\n\n")

	// Circular dependencies
	if len(graph.Cycles) > 0 {
		sb.WriteString("## ⚠️ Circular Dependencies\n\n")
		sb.WriteString("The following circular dependencies were detected and must be resolved:\n\n")
		
		for i, cycle := range graph.Cycles {
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

	// Critical path
	if len(analysis.CriticalPath) > 0 {
		sb.WriteString("## 🎯 Critical Path\n\n")
		sb.WriteString("The longest dependency chain in your tasks:\n\n")
		sb.WriteString("```\n")
		for i, task := range analysis.CriticalPath {
			if i > 0 {
				sb.WriteString(" → ")
			}
			sb.WriteString(task)
		}
		sb.WriteString("\n```\n\n")
		sb.WriteString(fmt.Sprintf("**Length:** %d tasks\n\n", len(analysis.CriticalPath)))
		sb.WriteString("Optimizing tasks in the critical path will have the biggest impact on overall execution time.\n\n")
	}

	// Task levels
	sb.WriteString("## 📊 Task Levels\n\n")
	sb.WriteString("Tasks grouped by their dependency depth:\n\n")
	
	// Group tasks by level
	levelGroups := make(map[int][]string)
	for task, level := range graph.Levels {
		levelGroups[level] = append(levelGroups[level], task)
	}
	
	// Sort levels
	var levels []int
	for level := range levelGroups {
		levels = append(levels, level)
	}
	sort.Ints(levels)
	
	for _, level := range levels {
		tasks := levelGroups[level]
		sort.Strings(tasks)
		
		sb.WriteString(fmt.Sprintf("**Level %d** (can run in parallel):\n", level))
		for _, task := range tasks {
			normalizedName := NormalizeTaskName(task)
			sb.WriteString(fmt.Sprintf("- [%s](%s.md)\n", task, normalizedName))
		}
		sb.WriteString("\n")
	}

	// Navigation
	sb.WriteString("## Navigation\n\n")
	sb.WriteString("- [← Back to Overview](../README.md)\n")
	sb.WriteString("- [Task Usage Analysis](../summaries/task-usage.md)\n")
	sb.WriteString("- [Optimization Guide](../optimization-guide.md)\n")

	return sb.String()
}