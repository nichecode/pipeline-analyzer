package circleci

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

// ParseConfig parses a CircleCI config file
func ParseConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return &config, nil
}

// ExtractWorkflowJobs extracts jobs from workflow with proper type handling
func ExtractWorkflowJobs(workflow Workflow) []WorkflowJob {
	var jobs []WorkflowJob

	for _, jobInterface := range workflow.Jobs {
		switch job := jobInterface.(type) {
		case string:
			// Simple job name
			jobs = append(jobs, WorkflowJob{Name: job})
		case map[string]interface{}:
			// Complex job object
			for jobName, jobConfig := range job {
				wfJob := WorkflowJob{Name: jobName}

				if config, ok := jobConfig.(map[string]interface{}); ok {
					if requires, ok := config["requires"].([]interface{}); ok {
						for _, req := range requires {
							if reqStr, ok := req.(string); ok {
								wfJob.Requires = append(wfJob.Requires, reqStr)
							}
						}
					}

					if context, ok := config["context"].([]interface{}); ok {
						for _, ctx := range context {
							if ctxStr, ok := ctx.(string); ok {
								wfJob.Context = append(wfJob.Context, ctxStr)
							}
						}
					} else if context, ok := config["context"].(string); ok {
						wfJob.Context = append(wfJob.Context, context)
					}

					if filters, ok := config["filters"].(map[string]interface{}); ok {
						wfJob.Filters = filters
					}
				}

				jobs = append(jobs, wfJob)
			}
		}
	}

	return jobs
}

// ExtractCommands extracts run commands from job steps
func ExtractCommands(steps []interface{}) []string {
	var commands []string

	for _, stepInterface := range steps {
		switch step := stepInterface.(type) {
		case string:
			// Simple step name (like "checkout")
			if step != "checkout" && step != "setup_remote_docker" {
				commands = append(commands, step)
			}
		case map[string]interface{}:
			// Complex step object
			if runStep, ok := step["run"]; ok {
				cmd := extractRunCommand(runStep)
				if cmd != "" {
					commands = append(commands, cmd)
				}
			}
		}
	}

	return commands
}

// extractRunCommand extracts command from run step (handles both string and object forms)
func extractRunCommand(runStep interface{}) string {
	switch run := runStep.(type) {
	case string:
		// Simple run: "echo hello"
		return run
	case map[string]interface{}:
		// Complex run with command field
		if command, ok := run["command"].(string); ok {
			return command
		}
	}
	return ""
}

// ExtractDockerImages extracts Docker images from job configuration
func ExtractDockerImages(job Job) []string {
	var images []string

	for _, docker := range job.Docker {
		if docker.Image != "" {
			images = append(images, docker.Image)
		}
	}

	return images
}

// GetExecutorImages gets Docker images from executor configuration
func GetExecutorImages(config *Config, executorName string) []string {
	if executor, ok := config.Executors[executorName]; ok {
		var images []string
		for _, docker := range executor.Docker {
			if docker.Image != "" {
				images = append(images, docker.Image)
			}
		}
		return images
	}
	return nil
}

// NormalizeJobName removes special characters from job names for file naming
func NormalizeJobName(name string) string {
	// Replace invalid filename characters
	name = strings.ReplaceAll(name, "/", "-")
	name = strings.ReplaceAll(name, ":", "-")
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ToLower(name)
	return name
}

// IsValidConfig checks if the config has the basic required structure
func IsValidConfig(config *Config) error {
	if config.Version == "" {
		return fmt.Errorf("missing version field")
	}

	if len(config.Jobs) == 0 {
		return fmt.Errorf("no jobs defined")
	}

	return nil
}

// GetAllJobNames returns all unique job names from the configuration
func GetAllJobNames(config *Config) []string {
	jobNames := make([]string, 0, len(config.Jobs))
	for name := range config.Jobs {
		jobNames = append(jobNames, name)
	}
	return jobNames
}

// GetAllWorkflowNames returns all workflow names from the configuration
func GetAllWorkflowNames(config *Config) []string {
	workflowNames := make([]string, 0, len(config.Workflows))
	for name := range config.Workflows {
		workflowNames = append(workflowNames, name)
	}
	return workflowNames
}

// DeepCopy creates a deep copy of interface{} for safe manipulation
func DeepCopy(src interface{}) interface{} {
	if src == nil {
		return nil
	}

	original := reflect.ValueOf(src)
	copy := reflect.New(original.Type()).Elem()
	copyRecursive(original, copy)
	return copy.Interface()
}

func copyRecursive(original, copy reflect.Value) {
	switch original.Kind() {
	case reflect.Interface:
		originalValue := original.Elem()
		if !originalValue.IsValid() {
			return
		}
		copyValue := reflect.New(originalValue.Type()).Elem()
		copyRecursive(originalValue, copyValue)
		copy.Set(copyValue)

	case reflect.Ptr:
		originalValue := original.Elem()
		if !originalValue.IsValid() {
			return
		}
		copy.Set(reflect.New(originalValue.Type()))
		copyRecursive(originalValue, copy.Elem())

	case reflect.Map:
		copy.Set(reflect.MakeMap(original.Type()))
		for _, key := range original.MapKeys() {
			originalValue := original.MapIndex(key)
			copyValue := reflect.New(originalValue.Type()).Elem()
			copyRecursive(originalValue, copyValue)
			copy.SetMapIndex(key, copyValue)
		}

	case reflect.Slice:
		copy.Set(reflect.MakeSlice(original.Type(), original.Len(), original.Cap()))
		for i := 0; i < original.Len(); i++ {
			copyRecursive(original.Index(i), copy.Index(i))
		}

	default:
		copy.Set(original)
	}
}
