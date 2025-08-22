package circleci

import (
	"regexp"
	"sort"
	"strings"
	"time"
)

// AnalyzeConfig performs comprehensive analysis of CircleCI configuration
func AnalyzeConfig(config *Config) *Analysis {
	analysis := &Analysis{
		Config:           config,
		JobUsage:         make(map[string]int),
		JobDependencies:  make(map[string][]string),
		CommandPatterns:  make(map[string]PatternCount),
		ExecutorUsage:    make(map[string][]string),
		CommandUsage:     make(map[string]int),
		ReusableCommands: make(map[string]*CommandAnalysis),
		TotalJobs:        len(config.Jobs),
		TotalWorkflows:   len(config.Workflows),
		TotalCommands:    len(config.Commands),
		GeneratedAt:      time.Now(),
	}

	// Analyze reusable commands
	analyzeReusableCommands(config, analysis)

	// Analyze job usage in workflows
	analyzeJobUsage(config, analysis)
	
	// Analyze job dependencies
	analyzeJobDependencies(config, analysis)
	
	// Analyze command patterns
	analyzeCommandPatterns(config, analysis)
	
	// Analyze executor usage
	analyzeExecutorUsage(config, analysis)

	return analysis
}

// analyzeJobUsage counts how many times each job is used across workflows
func analyzeJobUsage(config *Config, analysis *Analysis) {
	for _, workflow := range config.Workflows {
		jobs := ExtractWorkflowJobs(workflow)
		for _, job := range jobs {
			analysis.JobUsage[job.Name]++
		}
	}
}

// analyzeJobDependencies extracts job dependency relationships
func analyzeJobDependencies(config *Config, analysis *Analysis) {
	for _, workflow := range config.Workflows {
		jobs := ExtractWorkflowJobs(workflow)
		for _, job := range jobs {
			if len(job.Requires) > 0 {
				analysis.JobDependencies[job.Name] = job.Requires
			}
		}
	}
}

// analyzeCommandPatterns analyzes patterns in run commands
func analyzeCommandPatterns(config *Config, analysis *Analysis) {
	patterns := map[string]*regexp.Regexp{
		"docker-compose": regexp.MustCompile(`docker-compose`),
		"docker":         regexp.MustCompile(`\bdocker\s+`),
		"npm":           regexp.MustCompile(`npm\s+`),
		"yarn":          regexp.MustCompile(`yarn\s+`),
		"go":            regexp.MustCompile(`go\s+`),
		"python":        regexp.MustCompile(`python\s+`),
		"pip":           regexp.MustCompile(`pip\s+`),
		"make":          regexp.MustCompile(`make\s+`),
		"curl":          regexp.MustCompile(`curl\s+`),
		"wget":          regexp.MustCompile(`wget\s+`),
		"git":           regexp.MustCompile(`git\s+`),
		"ssh":           regexp.MustCompile(`ssh\s+`),
		"scp":           regexp.MustCompile(`scp\s+`),
		"rsync":         regexp.MustCompile(`rsync\s+`),
		"tar":           regexp.MustCompile(`tar\s+`),
		"zip":           regexp.MustCompile(`zip\s+`),
		"unzip":         regexp.MustCompile(`unzip\s+`),
		"echo":          regexp.MustCompile(`echo\s+`),
		"cat":           regexp.MustCompile(`cat\s+`),
		"grep":          regexp.MustCompile(`grep\s+`),
		"sed":           regexp.MustCompile(`sed\s+`),
		"awk":           regexp.MustCompile(`awk\s+`),
		"test":          regexp.MustCompile(`test\s+`),
		"./":            regexp.MustCompile(`\./[^\s]+`),
	}

	for jobName, job := range config.Jobs {
		commands := ExtractCommands(job.Steps)
		for _, command := range commands {
			for patternName, pattern := range patterns {
				if pattern.MatchString(command) {
					if _, exists := analysis.CommandPatterns[patternName]; !exists {
						analysis.CommandPatterns[patternName] = PatternCount{Jobs: []string{}}
					}
					
					pc := analysis.CommandPatterns[patternName]
					pc.Count++
					
					// Add job to list if not already present
					found := false
					for _, existingJob := range pc.Jobs {
						if existingJob == jobName {
							found = true
							break
						}
					}
					if !found {
						pc.Jobs = append(pc.Jobs, jobName)
					}
					
					analysis.CommandPatterns[patternName] = pc
				}
			}
		}
	}
}

// analyzeExecutorUsage tracks which jobs use which executors/images
func analyzeExecutorUsage(config *Config, analysis *Analysis) {
	for jobName, job := range config.Jobs {
		// Direct Docker images
		images := ExtractDockerImages(job)
		for _, image := range images {
			analysis.ExecutorUsage[image] = append(analysis.ExecutorUsage[image], jobName)
		}
		
		// Executor references
		if job.Executor != "" {
			executorImages := GetExecutorImages(config, job.Executor)
			for _, image := range executorImages {
				key := job.Executor + " (" + image + ")"
				analysis.ExecutorUsage[key] = append(analysis.ExecutorUsage[key], jobName)
			}
		}
	}
}

// AnalyzeJob performs detailed analysis of a single job
func AnalyzeJob(config *Config, jobName string, analysis *Analysis) *JobAnalysis {
	job, exists := config.Jobs[jobName]
	if !exists {
		return nil
	}

	jobAnalysis := &JobAnalysis{
		Name:         jobName,
		Description:  job.Description,
		Commands:     ExtractCommands(job.Steps),
		DockerImages: ExtractDockerImages(job),
		Executor:     job.Executor,
		UsageCount:   analysis.JobUsage[jobName],
		Patterns:     make(map[string]int),
	}

	// Add executor images if applicable
	if job.Executor != "" {
		executorImages := GetExecutorImages(config, job.Executor)
		jobAnalysis.DockerImages = append(jobAnalysis.DockerImages, executorImages...)
	}

	// Count pattern usage in this job
	for pattern, patternCount := range analysis.CommandPatterns {
		for _, patternJob := range patternCount.Jobs {
			if patternJob == jobName {
				jobAnalysis.Patterns[pattern]++
			}
		}
	}

	// Get dependencies
	if deps, exists := analysis.JobDependencies[jobName]; exists {
		jobAnalysis.Dependencies = deps
	}

	return jobAnalysis
}

// AnalyzeWorkflow performs detailed analysis of a single workflow
func AnalyzeWorkflow(config *Config, workflowName string) *WorkflowAnalysis {
	workflow, exists := config.Workflows[workflowName]
	if !exists {
		return nil
	}

	return &WorkflowAnalysis{
		Name: workflowName,
		Jobs: ExtractWorkflowJobs(workflow),
	}
}

// GetMostUsedJobs returns jobs sorted by usage frequency
func GetMostUsedJobs(analysis *Analysis) []struct {
	Name  string
	Count int
} {
	type jobUsage struct {
		Name  string
		Count int
	}

	var usage []jobUsage
	for name, count := range analysis.JobUsage {
		usage = append(usage, jobUsage{Name: name, Count: count})
	}

	sort.Slice(usage, func(i, j int) bool {
		return usage[i].Count > usage[j].Count
	})

	// Convert to return type
	result := make([]struct {
		Name  string
		Count int
	}, len(usage))
	
	for i, job := range usage {
		result[i] = struct {
			Name  string
			Count int
		}{Name: job.Name, Count: job.Count}
	}

	return result
}

// GetMostUsedPatterns returns command patterns sorted by frequency
func GetMostUsedPatterns(analysis *Analysis) []struct {
	Pattern string
	Count   int
	Jobs    []string
} {
	type patternUsage struct {
		Pattern string
		Count   int
		Jobs    []string
	}

	var usage []patternUsage
	for pattern, patternCount := range analysis.CommandPatterns {
		usage = append(usage, patternUsage{
			Pattern: pattern,
			Count:   patternCount.Count,
			Jobs:    patternCount.Jobs,
		})
	}

	sort.Slice(usage, func(i, j int) bool {
		return usage[i].Count > usage[j].Count
	})

	// Convert to return type
	result := make([]struct {
		Pattern string
		Count   int
		Jobs    []string
	}, len(usage))
	
	for i, pattern := range usage {
		result[i] = struct {
			Pattern string
			Count   int
			Jobs    []string
		}{Pattern: pattern.Pattern, Count: pattern.Count, Jobs: pattern.Jobs}
	}

	return result
}

// GetCommandFrequency analyzes first word frequency in commands
func GetCommandFrequency(config *Config) map[string]int {
	frequency := make(map[string]int)
	
	for _, job := range config.Jobs {
		commands := ExtractCommands(job.Steps)
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

// CountDockerUsage counts jobs using Docker vs other executors
func CountDockerUsage(config *Config) (int, int) {
	dockerJobs := 0
	otherJobs := 0
	
	for _, job := range config.Jobs {
		if len(job.Docker) > 0 || job.Executor != "" {
			dockerJobs++
		} else {
			otherJobs++
		}
	}
	
	return dockerJobs, otherJobs
}

// GetJobsByPattern returns jobs that match a specific pattern
func GetJobsByPattern(analysis *Analysis, pattern string) []string {
	if patternCount, exists := analysis.CommandPatterns[pattern]; exists {
		return patternCount.Jobs
	}
	return nil
}

// GetAllCommandNames returns all reusable command names from the config
func GetAllCommandNames(config *Config) []string {
	commandNames := make([]string, 0, len(config.Commands))
	for name := range config.Commands {
		commandNames = append(commandNames, name)
	}
	return commandNames
}

// analyzeReusableCommands analyzes reusable command definitions
func analyzeReusableCommands(config *Config, analysis *Analysis) {
	for cmdName, command := range config.Commands {
		// Extract commands from command steps
		commands := ExtractCommands(command.Steps)
		
		// Create command analysis
		cmdAnalysis := &CommandAnalysis{
			Name:        cmdName,
			Description: command.Description,
			Commands:    commands,
			UsageCount:  0, // Will be updated when counting usage
			Parameters:  make(map[string]interface{}),
			Patterns:    make(map[string]int),
		}
		
		// Parse parameters if they exist
		if command.Parameters != nil {
			if params, ok := command.Parameters.(map[string]interface{}); ok {
				cmdAnalysis.Parameters = params
			}
		}
		
		// Count command patterns in this reusable command
		patterns := map[string]*regexp.Regexp{
			"docker":    regexp.MustCompile(`\bdocker\s+`),
			"npm":       regexp.MustCompile(`\bnpm\s+`),
			"yarn":      regexp.MustCompile(`\byarn\s+`),
			"git":       regexp.MustCompile(`\bgit\s+`),
			"curl":      regexp.MustCompile(`\bcurl\s+`),
			"aws":       regexp.MustCompile(`\baws\s+`),
			"kubectl":   regexp.MustCompile(`\bkubectl\s+`),
			"helm":      regexp.MustCompile(`\bhelm\s+`),
			"python":    regexp.MustCompile(`\bpython\s+`),
			"pip":       regexp.MustCompile(`\bpip\s+`),
		}
		
		for patternName, pattern := range patterns {
			for _, command := range commands {
				if pattern.MatchString(command) {
					cmdAnalysis.Patterns[patternName]++
				}
			}
		}
		
		analysis.ReusableCommands[cmdName] = cmdAnalysis
	}
	
	// Count reusable command usage in jobs
	for _, job := range config.Jobs {
		for _, stepInterface := range job.Steps {
			if stepMap, ok := stepInterface.(map[string]interface{}); ok {
				for stepKey := range stepMap {
					if _, exists := config.Commands[stepKey]; exists {
						analysis.CommandUsage[stepKey]++
						analysis.ReusableCommands[stepKey].UsageCount++
					}
				}
			}
		}
	}
}