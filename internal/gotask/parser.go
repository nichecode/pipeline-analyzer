package gotask

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/nichecode/pipeline-analyzer/internal/shared"
	"gopkg.in/yaml.v3"
)

// ParseTaskfile parses a Taskfile.yml file with error recovery
func ParseTaskfile(taskfilePath string) (*Taskfile, error) {
	logger := shared.GetLogger()
	logger.Debug("GoTask", "Starting taskfile parsing", map[string]interface{}{
		"file": taskfilePath,
	})

	data, err := os.ReadFile(taskfilePath)
	if err != nil {
		logger.Error("GoTask", "Failed to read taskfile", map[string]interface{}{
			"file":  taskfilePath,
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to read taskfile: %w", err)
	}

	var taskfile Taskfile
	if err := yaml.Unmarshal(data, &taskfile); err != nil {
		logger.Warn("GoTask", "Standard YAML parsing failed, attempting recovery", map[string]interface{}{
			"file":  taskfilePath,
			"error": err.Error(),
		})
		
		// Try to parse with more lenient approach
		if recoveredTaskfile, recoverErr := parseTaskfileWithRecovery(data, taskfilePath); recoverErr == nil {
			logger.RecoveryAttempt("GoTask", taskfilePath, "flexible_parsing", true)
			return recoveredTaskfile, nil
		} else {
			logger.RecoveryAttempt("GoTask", taskfilePath, "flexible_parsing", false)
		}
		
		// Log detailed parsing error
		logger.ParseError("GoTask", taskfilePath, err, string(data))
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	logger.Info("GoTask", "Taskfile parsed successfully", map[string]interface{}{
		"file":      taskfilePath,
		"version":   taskfile.Version,
		"tasks":     len(taskfile.Tasks),
		"includes":  len(taskfile.Includes),
		"vars":      len(taskfile.Vars),
	})

	return &taskfile, nil
}

// parseTaskfileWithRecovery attempts to parse with basic structure recovery
func parseTaskfileWithRecovery(data []byte, taskfilePath string) (*Taskfile, error) {
	logger := shared.GetLogger()
	logger.Debug("GoTask", "Attempting recovery parsing", map[string]interface{}{
		"file": taskfilePath,
	})
	// Parse as raw interface{} first to check structure
	var raw map[string]interface{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	taskfile := &Taskfile{}
	
	// Safely extract version
	if version, ok := raw["version"].(string); ok {
		taskfile.Version = version
	}
	
	// Safely extract includes with flexible typing
	if includesRaw, ok := raw["includes"]; ok {
		taskfile.Includes = make(map[string]interface{})
		if includesMap, ok := includesRaw.(map[string]interface{}); ok {
			taskfile.Includes = includesMap
		}
	}
	
	// Safely extract vars
	if varsRaw, ok := raw["vars"]; ok {
		if varsMap, ok := varsRaw.(map[string]interface{}); ok {
			taskfile.Vars = varsMap
		}
	}
	
	// Safely extract tasks
	if tasksRaw, ok := raw["tasks"]; ok {
		taskfile.Tasks = make(map[string]Task)
		if tasksMap, ok := tasksRaw.(map[string]interface{}); ok {
			for taskName, taskRaw := range tasksMap {
				if taskMap, ok := taskRaw.(map[string]interface{}); ok {
					task := parseTaskFromMap(taskMap)
					taskfile.Tasks[taskName] = task
				}
			}
		}
	}
	
	return taskfile, nil
}

// parseTaskFromMap safely parses a task from a map
func parseTaskFromMap(taskMap map[string]interface{}) Task {
	task := Task{}
	
	if desc, ok := taskMap["desc"].(string); ok {
		task.Desc = desc
	}
	
	if cmdsRaw, ok := taskMap["cmds"]; ok {
		task.Cmds = []interface{}{}
		switch cmds := cmdsRaw.(type) {
		case []interface{}:
			task.Cmds = cmds
		case string:
			task.Cmds = []interface{}{cmds}
		}
	}
	
	if depsRaw, ok := taskMap["deps"]; ok {
		task.Deps = []interface{}{}
		switch deps := depsRaw.(type) {
		case []interface{}:
			task.Deps = deps
		case string:
			task.Deps = []interface{}{deps}
		}
	}
	
	return task
}

// ExtractTaskCommands extracts all commands from a task
func ExtractTaskCommands(task Task) []string {
	var commands []string

	// Handle single command
	if task.Cmd != "" {
		commands = append(commands, task.Cmd)
	}

	// Handle multiple commands
	for _, cmdInterface := range task.Cmds {
		cmd := extractCommand(cmdInterface)
		if cmd != "" {
			commands = append(commands, cmd)
		}
	}

	return commands
}

// extractCommand extracts command from interface (string or Command object)
func extractCommand(cmdInterface interface{}) string {
	switch cmd := cmdInterface.(type) {
	case string:
		return cmd
	case map[string]interface{}:
		if command, ok := cmd["cmd"].(string); ok {
			return command
		}
	}
	return ""
}

// ExtractTaskDependencies extracts dependencies from a task
func ExtractTaskDependencies(task Task) []string {
	var dependencies []string

	for _, depInterface := range task.Deps {
		dep := extractDependency(depInterface)
		if dep != "" {
			dependencies = append(dependencies, dep)
		}
	}

	return dependencies
}

// extractDependency extracts dependency from interface (string or Dependency object)
func extractDependency(depInterface interface{}) string {
	switch dep := depInterface.(type) {
	case string:
		return dep
	case map[string]interface{}:
		if task, ok := dep["task"].(string); ok {
			return task
		}
	}
	return ""
}

// ExtractPreconditions extracts preconditions from a task
func ExtractPreconditions(task Task) []string {
	var preconditions []string

	for _, precondInterface := range task.Preconditions {
		precond := extractPrecondition(precondInterface)
		if precond != "" {
			preconditions = append(preconditions, precond)
		}
	}

	return preconditions
}

// extractPrecondition extracts precondition from interface
func extractPrecondition(precondInterface interface{}) string {
	switch precond := precondInterface.(type) {
	case string:
		return precond
	case map[string]interface{}:
		if sh, ok := precond["sh"].(string); ok {
			return sh
		}
	}
	return ""
}

// FindTaskfile finds a Taskfile in the given directory
func FindTaskfile(dir string) (string, error) {
	taskfileNames := []string{
		"Taskfile.yml",
		"Taskfile.yaml",
		"taskfile.yml",
		"taskfile.yaml",
	}

	for _, name := range taskfileNames {
		path := filepath.Join(dir, name)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("no taskfile found in %s", dir)
}

// ParseIncludedTaskfile parses an included taskfile
func ParseIncludedTaskfile(include Include, basePath string) (*Taskfile, error) {
	var taskfilePath string

	if filepath.IsAbs(include.Taskfile) {
		taskfilePath = include.Taskfile
	} else {
		baseDir := filepath.Dir(basePath)
		if include.Dir != "" {
			baseDir = filepath.Join(baseDir, include.Dir)
		}
		taskfilePath = filepath.Join(baseDir, include.Taskfile)
	}

	return ParseTaskfile(taskfilePath)
}

// ExtractAllVariables extracts all variables from taskfile and tasks
func ExtractAllVariables(taskfile *Taskfile) map[string]Variable {
	variables := make(map[string]Variable)

	// Global variables
	for name, value := range taskfile.Vars {
		variables[name] = Variable{
			Name:      name,
			Value:     value,
			Type:      getVariableType(value),
			IsGlobal:  true,
		}
	}

	// Global environment variables
	for name, value := range taskfile.Env {
		variables[name] = Variable{
			Name:          name,
			Value:         value,
			Type:          getVariableType(value),
			IsGlobal:      true,
			IsEnvironment: true,
		}
	}

	// Task-specific variables
	for taskName, task := range taskfile.Tasks {
		for name, value := range task.Vars {
			key := fmt.Sprintf("%s.%s", taskName, name)
			variables[key] = Variable{
				Name:        name,
				Value:       value,
				Type:        getVariableType(value),
				UsedInTasks: []string{taskName},
			}
		}

		for name, value := range task.Env {
			key := fmt.Sprintf("%s.%s", taskName, name)
			variables[key] = Variable{
				Name:          name,
				Value:         value,
				Type:          getVariableType(value),
				UsedInTasks:   []string{taskName},
				IsEnvironment: true,
			}
		}
	}

	return variables
}

// getVariableType determines the type of a variable
func getVariableType(value interface{}) string {
	if value == nil {
		return "nil"
	}

	switch v := value.(type) {
	case string:
		// Check if it's a shell command
		if strings.Contains(v, "sh:") {
			return "shell"
		}
		return "string"
	case int, int8, int16, int32, int64:
		return "int"
	case float32, float64:
		return "float"
	case bool:
		return "bool"
	case []interface{}:
		return "array"
	case map[string]interface{}:
		return "object"
	default:
		return reflect.TypeOf(value).String()
	}
}

// IsValidTaskfile validates the basic structure of a taskfile
func IsValidTaskfile(taskfile *Taskfile) error {
	if taskfile.Version == "" {
		return fmt.Errorf("missing version field")
	}

	if len(taskfile.Tasks) == 0 {
		return fmt.Errorf("no tasks defined")
	}

	// Validate version format
	validVersions := []string{"1", "2", "3", "2.1", "2.2", "2.6", "3"}
	validVersion := false
	for _, v := range validVersions {
		if taskfile.Version == v {
			validVersion = true
			break
		}
	}

	if !validVersion {
		return fmt.Errorf("unsupported version: %s", taskfile.Version)
	}

	return nil
}

// GetAllTaskNames returns all task names from the taskfile
func GetAllTaskNames(taskfile *Taskfile) []string {
	taskNames := make([]string, 0, len(taskfile.Tasks))
	for name := range taskfile.Tasks {
		taskNames = append(taskNames, name)
	}
	return taskNames
}

// GetAllIncludeNames returns all include names from the taskfile
func GetAllIncludeNames(taskfile *Taskfile) []string {
	includeNames := make([]string, 0, len(taskfile.Includes))
	for name := range taskfile.Includes {
		includeNames = append(includeNames, name)
	}
	return includeNames
}

// NormalizeTaskName removes special characters from task names for file naming
func NormalizeTaskName(name string) string {
	name = strings.ReplaceAll(name, "/", "-")
	name = strings.ReplaceAll(name, ":", "-")
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, "_", "-")
	name = strings.ToLower(name)
	return name
}

// ExtractTaskPlatforms extracts platform information from a task
func ExtractTaskPlatforms(task Task) []string {
	platforms := make([]string, len(task.Platforms))
	copy(platforms, task.Platforms)
	return platforms
}

// HasSources checks if a task has source files defined
func HasSources(task Task) bool {
	return len(task.Sources) > 0
}

// HasGenerates checks if a task has generates files defined
func HasGenerates(task Task) bool {
	return len(task.Generates) > 0
}

// IsOptimizedForCaching checks if a task is optimized for caching
func IsOptimizedForCaching(task Task) bool {
	return HasSources(task) && HasGenerates(task)
}

// GetTaskComplexity calculates a complexity score for a task
func GetTaskComplexity(task Task) int {
	complexity := 0

	// Base complexity for commands
	complexity += len(task.Cmds)
	if task.Cmd != "" {
		complexity++
	}

	// Additional complexity factors
	complexity += len(task.Deps)
	complexity += len(task.Preconditions)
	complexity += len(task.Status)

	if len(task.Vars) > 0 {
		complexity++
	}

	if len(task.Env) > 0 {
		complexity++
	}

	return complexity
}

// DetectTaskType attempts to classify the task based on its properties
func DetectTaskType(task Task, taskName string) string {
	name := strings.ToLower(taskName)
	commands := ExtractTaskCommands(task)
	
	// Join all commands for pattern matching
	allCommands := strings.ToLower(strings.Join(commands, " "))

	// Classification based on name patterns
	if strings.Contains(name, "build") {
		return "build"
	}
	if strings.Contains(name, "test") {
		return "test"
	}
	if strings.Contains(name, "deploy") {
		return "deploy"
	}
	if strings.Contains(name, "clean") {
		return "cleanup"
	}
	if strings.Contains(name, "lint") || strings.Contains(name, "format") {
		return "quality"
	}
	if strings.Contains(name, "install") || strings.Contains(name, "setup") {
		return "setup"
	}

	// Classification based on command patterns
	if strings.Contains(allCommands, "go build") || strings.Contains(allCommands, "go run") {
		return "build"
	}
	if strings.Contains(allCommands, "go test") || strings.Contains(allCommands, "npm test") {
		return "test"
	}
	if strings.Contains(allCommands, "docker") {
		return "containerization"
	}
	if strings.Contains(allCommands, "kubectl") || strings.Contains(allCommands, "helm") {
		return "deployment"
	}

	return "utility"
}