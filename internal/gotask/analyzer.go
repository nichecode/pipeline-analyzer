package gotask

import (
	"regexp"
	"sort"
	"strings"
	"time"
)

// AnalyzeTaskfile performs comprehensive analysis of a Taskfile
func AnalyzeTaskfile(taskfile *Taskfile) *Analysis {
	analysis := &Analysis{
		Taskfile:         taskfile,
		TaskUsage:        make(map[string]int),
		TaskDependencies: make(map[string][]string),
		CommandPatterns:  make(map[string]PatternCount),
		VariableUsage:    make(map[string][]string),
		EnvironmentUsage: make(map[string][]string),
		IncludeAnalysis:  make(map[string]*IncludeAnalysis),
		TotalTasks:       len(taskfile.Tasks),
		TotalIncludes:    len(taskfile.Includes),
		GeneratedAt:      time.Now(),
	}

	// Analyze task dependencies and usage
	analyzeTaskDependencies(taskfile, analysis)
	
	// Analyze command patterns
	analyzeCommandPatterns(taskfile, analysis)
	
	// Analyze variables and environment
	analyzeVariables(taskfile, analysis)
	
	// Analyze includes
	analyzeIncludes(taskfile, analysis)
	
	// Build dependency graph and detect cycles
	analysis.CircularDeps = detectCircularDependencies(analysis.TaskDependencies)
	analysis.CriticalPath = findCriticalPath(analysis.TaskDependencies)
	
	// Generate optimization tips
	analysis.OptimizationTips = generateOptimizationTips(taskfile, analysis)

	return analysis
}

// analyzeTaskDependencies analyzes task dependencies and usage patterns
func analyzeTaskDependencies(taskfile *Taskfile, analysis *Analysis) {
	for taskName, task := range taskfile.Tasks {
		deps := ExtractTaskDependencies(task)
		analysis.TaskDependencies[taskName] = deps
		
		// Count how many tasks depend on each task
		for _, dep := range deps {
			analysis.TaskUsage[dep]++
		}
	}
}

// analyzeCommandPatterns analyzes patterns in task commands
func analyzeCommandPatterns(taskfile *Taskfile, analysis *Analysis) {
	patterns := map[string]*regexp.Regexp{
		"docker":         regexp.MustCompile(`\bdocker\s+`),
		"docker-compose": regexp.MustCompile(`docker-compose`),
		"go":            regexp.MustCompile(`\bgo\s+`),
		"npm":           regexp.MustCompile(`\bnpm\s+`),
		"yarn":          regexp.MustCompile(`\byarn\s+`),
		"make":          regexp.MustCompile(`\bmake\s+`),
		"git":           regexp.MustCompile(`\bgit\s+`),
		"curl":          regexp.MustCompile(`\bcurl\s+`),
		"wget":          regexp.MustCompile(`\bwget\s+`),
		"kubectl":       regexp.MustCompile(`\bkubectl\s+`),
		"helm":          regexp.MustCompile(`\bhelm\s+`),
		"terraform":     regexp.MustCompile(`\bterraform\s+`),
		"ansible":       regexp.MustCompile(`\bansible\s+`),
		"python":        regexp.MustCompile(`\bpython\s+`),
		"pip":           regexp.MustCompile(`\bpip\s+`),
		"mvn":           regexp.MustCompile(`\bmvn\s+`),
		"gradle":        regexp.MustCompile(`\bgradle\s+`),
		"cargo":         regexp.MustCompile(`\bcargo\s+`),
		"rustc":         regexp.MustCompile(`\brustc\s+`),
		"./":            regexp.MustCompile(`\./[^\s]+`),
		"ssh":           regexp.MustCompile(`\bssh\s+`),
		"scp":           regexp.MustCompile(`\bscp\s+`),
		"rsync":         regexp.MustCompile(`\brsync\s+`),
	}

	for taskName, task := range taskfile.Tasks {
		commands := ExtractTaskCommands(task)
		for _, command := range commands {
			for patternName, pattern := range patterns {
				if pattern.MatchString(command) {
					if _, exists := analysis.CommandPatterns[patternName]; !exists {
						analysis.CommandPatterns[patternName] = PatternCount{Tasks: []string{}}
					}
					
					pc := analysis.CommandPatterns[patternName]
					pc.Count++
					
					// Add task to list if not already present
					found := false
					for _, existingTask := range pc.Tasks {
						if existingTask == taskName {
							found = true
							break
						}
					}
					if !found {
						pc.Tasks = append(pc.Tasks, taskName)
					}
					
					analysis.CommandPatterns[patternName] = pc
				}
			}
		}
	}
}

// analyzeVariables analyzes variable and environment usage
func analyzeVariables(taskfile *Taskfile, analysis *Analysis) {
	allVars := ExtractAllVariables(taskfile)
	
	for _, variable := range allVars {
		if variable.IsEnvironment {
			analysis.EnvironmentUsage[variable.Name] = variable.UsedInTasks
		} else {
			analysis.VariableUsage[variable.Name] = variable.UsedInTasks
		}
	}
}

// analyzeIncludes analyzes included taskfiles
func analyzeIncludes(taskfile *Taskfile, analysis *Analysis) {
	for name, include := range taskfile.Includes {
		includeAnalysis := &IncludeAnalysis{
			Path:      include.Taskfile,
			Namespace: name,
		}
		
		// Try to parse included taskfile for deeper analysis
		// Note: This would require the actual file to exist
		analysis.IncludeAnalysis[name] = includeAnalysis
	}
}

// detectCircularDependencies detects circular dependencies in tasks
func detectCircularDependencies(dependencies map[string][]string) [][]string {
	var cycles [][]string
	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	
	var dfs func(string, []string) []string
	dfs = func(task string, path []string) []string {
		if recStack[task] {
			// Found a cycle, extract it
			cycleStart := -1
			for i, t := range path {
				if t == task {
					cycleStart = i
					break
				}
			}
			if cycleStart != -1 {
				return append(path[cycleStart:], task)
			}
		}
		
		if visited[task] {
			return nil
		}
		
		visited[task] = true
		recStack[task] = true
		newPath := append(path, task)
		
		for _, dep := range dependencies[task] {
			if cycle := dfs(dep, newPath); cycle != nil {
				return cycle
			}
		}
		
		recStack[task] = false
		return nil
	}
	
	for task := range dependencies {
		if !visited[task] {
			if cycle := dfs(task, []string{}); cycle != nil {
				cycles = append(cycles, cycle)
			}
		}
	}
	
	return cycles
}

// findCriticalPath finds the longest dependency chain
func findCriticalPath(dependencies map[string][]string) []string {
	depths := make(map[string]int)
	paths := make(map[string][]string)
	
	var calculateDepth func(string) int
	calculateDepth = func(task string) int {
		if depth, exists := depths[task]; exists {
			return depth
		}
		
		maxDepth := 0
		longestPath := []string{task}
		
		for _, dep := range dependencies[task] {
			depDepth := calculateDepth(dep)
			if depDepth >= maxDepth {
				maxDepth = depDepth
				longestPath = append([]string{task}, paths[dep]...)
			}
		}
		
		depths[task] = maxDepth + 1
		paths[task] = longestPath
		return depths[task]
	}
	
	maxLength := 0
	var criticalPath []string
	
	for task := range dependencies {
		depth := calculateDepth(task)
		if depth > maxLength {
			maxLength = depth
			criticalPath = paths[task]
		}
	}
	
	return criticalPath
}

// generateOptimizationTips generates optimization suggestions
func generateOptimizationTips(taskfile *Taskfile, analysis *Analysis) []OptimizationTip {
	var tips []OptimizationTip
	
	// Check for tasks without caching optimization
	for taskName, task := range taskfile.Tasks {
		if !IsOptimizedForCaching(task) && GetTaskComplexity(task) > 2 {
			tips = append(tips, OptimizationTip{
				Type:       "caching",
				Task:       taskName,
				Message:    "Task could benefit from caching optimization",
				Severity:   "medium",
				Suggestion: "Add 'sources' and 'generates' fields to enable task result caching",
			})
		}
	}
	
	// Check for circular dependencies
	for _, cycle := range analysis.CircularDeps {
		tips = append(tips, OptimizationTip{
			Type:       "dependency",
			Task:       strings.Join(cycle, " -> "),
			Message:    "Circular dependency detected",
			Severity:   "high",
			Suggestion: "Refactor tasks to break the circular dependency",
		})
	}
	
	// Check for unused tasks
	for taskName := range taskfile.Tasks {
		if analysis.TaskUsage[taskName] == 0 {
			// Check if it's a root task (not depended on by others)
			tips = append(tips, OptimizationTip{
				Type:       "usage",
				Task:       taskName,
				Message:    "Task is not used as a dependency",
				Severity:   "low",
				Suggestion: "Consider if this task is needed or should be marked as internal",
			})
		}
	}
	
	// Check for missing descriptions
	for taskName, task := range taskfile.Tasks {
		if task.Desc == "" && !task.Internal {
			tips = append(tips, OptimizationTip{
				Type:       "documentation",
				Task:       taskName,
				Message:    "Task missing description",
				Severity:   "low",
				Suggestion: "Add a description to improve task documentation",
			})
		}
	}
	
	return tips
}

// AnalyzeTask performs detailed analysis of a single task
func AnalyzeTask(taskfile *Taskfile, taskName string, analysis *Analysis) *TaskAnalysis {
	task, exists := taskfile.Tasks[taskName]
	if !exists {
		return nil
	}

	taskAnalysis := &TaskAnalysis{
		Name:            taskName,
		Description:     task.Desc,
		Summary:         task.Summary,
		Commands:        ExtractTaskCommands(task),
		Dependencies:    ExtractTaskDependencies(task),
		Sources:         task.Sources,
		Generates:       task.Generates,
		Variables:       task.Vars,
		Environment:     task.Env,
		Patterns:        make(map[string]int),
		UsageCount:      analysis.TaskUsage[taskName],
		IsInternal:      task.Internal,
		Platforms:       task.Platforms,
		Aliases:         task.Aliases,
		Preconditions:   ExtractPreconditions(task),
		HasWatch:        task.Watch,
	}

	// Count pattern usage in this task
	for pattern, patternCount := range analysis.CommandPatterns {
		for _, patternTask := range patternCount.Tasks {
			if patternTask == taskName {
				taskAnalysis.Patterns[pattern]++
			}
		}
	}

	// Generate optimization opportunities for this task
	if !IsOptimizedForCaching(task) && GetTaskComplexity(task) > 2 {
		taskAnalysis.OptimizationOps = append(taskAnalysis.OptimizationOps, 
			"Consider adding sources and generates for caching")
	}

	if task.Desc == "" {
		taskAnalysis.OptimizationOps = append(taskAnalysis.OptimizationOps, 
			"Add description for better documentation")
	}

	return taskAnalysis
}

// GetMostUsedTasks returns tasks sorted by usage frequency
func GetMostUsedTasks(analysis *Analysis) []struct {
	Name  string
	Count int
} {
	type taskUsage struct {
		Name  string
		Count int
	}

	var usage []taskUsage
	for name, count := range analysis.TaskUsage {
		usage = append(usage, taskUsage{Name: name, Count: count})
	}

	sort.Slice(usage, func(i, j int) bool {
		return usage[i].Count > usage[j].Count
	})

	// Convert to return type
	result := make([]struct {
		Name  string
		Count int
	}, len(usage))
	
	for i, task := range usage {
		result[i] = struct {
			Name  string
			Count int
		}{Name: task.Name, Count: task.Count}
	}

	return result
}

// GetMostUsedPatterns returns command patterns sorted by frequency
func GetMostUsedPatterns(analysis *Analysis) []struct {
	Pattern string
	Count   int
	Tasks   []string
} {
	type patternUsage struct {
		Pattern string
		Count   int
		Tasks   []string
	}

	var usage []patternUsage
	for pattern, patternCount := range analysis.CommandPatterns {
		usage = append(usage, patternUsage{
			Pattern: pattern,
			Count:   patternCount.Count,
			Tasks:   patternCount.Tasks,
		})
	}

	sort.Slice(usage, func(i, j int) bool {
		return usage[i].Count > usage[j].Count
	})

	// Convert to return type
	result := make([]struct {
		Pattern string
		Count   int
		Tasks   []string
	}, len(usage))
	
	for i, pattern := range usage {
		result[i] = struct {
			Pattern string
			Count   int
			Tasks   []string
		}{Pattern: pattern.Pattern, Count: pattern.Count, Tasks: pattern.Tasks}
	}

	return result
}

// GetCommandFrequency analyzes first word frequency in commands
func GetCommandFrequency(taskfile *Taskfile) map[string]int {
	frequency := make(map[string]int)
	
	for _, task := range taskfile.Tasks {
		commands := ExtractTaskCommands(task)
		for _, command := range commands {
			// Extract first word (command)
			words := strings.Fields(command)
			if len(words) > 0 {
				firstWord := strings.ToLower(words[0])
				// Remove common shell operators
				firstWord = strings.TrimPrefix(firstWord, "sudo")
				firstWord = strings.TrimSpace(firstWord)
				if firstWord != "" && firstWord != "&&" && firstWord != "||" {
					frequency[firstWord]++
				}
			}
		}
	}
	
	return frequency
}

// GetPerformanceMetrics calculates performance-related metrics
func GetPerformanceMetrics(taskfile *Taskfile) PerformanceMetrics {
	metrics := PerformanceMetrics{}
	
	parallelizable := 0
	
	for _, task := range taskfile.Tasks {
		if HasSources(task) {
			metrics.TasksWithSources++
		}
		if HasGenerates(task) {
			metrics.TasksWithGenerates++
		}
		if IsOptimizedForCaching(task) {
			metrics.TasksWithCaching++
		}
		
		// Tasks with no dependencies can be parallelized
		if len(task.Deps) == 0 {
			parallelizable++
		}
	}
	
	metrics.ParallelizableTasks = parallelizable
	
	// Calculate optimization potential (percentage of tasks that could be optimized)
	totalTasks := float64(len(taskfile.Tasks))
	if totalTasks > 0 {
		optimizedTasks := float64(metrics.TasksWithCaching)
		metrics.OptimizationPotential = (totalTasks - optimizedTasks) / totalTasks * 100
	}
	
	return metrics
}

// GetTasksByPattern returns tasks that use a specific pattern
func GetTasksByPattern(analysis *Analysis, pattern string) []string {
	if patternCount, exists := analysis.CommandPatterns[pattern]; exists {
		return patternCount.Tasks
	}
	return nil
}

// BuildDependencyGraph creates a dependency graph structure
func BuildDependencyGraph(analysis *Analysis) *DependencyGraph {
	graph := &DependencyGraph{
		Tasks:  GetAllTaskNames(analysis.Taskfile),
		Edges:  analysis.TaskDependencies,
		Levels: make(map[string]int),
		Cycles: analysis.CircularDeps,
	}
	
	// Calculate dependency levels
	for task := range analysis.TaskDependencies {
		graph.Levels[task] = calculateTaskLevel(task, analysis.TaskDependencies, make(map[string]bool))
	}
	
	return graph
}

// calculateTaskLevel calculates the dependency level of a task
func calculateTaskLevel(task string, dependencies map[string][]string, visiting map[string]bool) int {
	if visiting[task] {
		return 0 // Circular dependency
	}
	
	visiting[task] = true
	defer func() { visiting[task] = false }()
	
	maxLevel := 0
	for _, dep := range dependencies[task] {
		level := calculateTaskLevel(dep, dependencies, visiting)
		if level > maxLevel {
			maxLevel = level
		}
	}
	
	return maxLevel + 1
}