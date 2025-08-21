package circleci

import "time"

// Config represents the main CircleCI configuration
type Config struct {
	Version   string                 `yaml:"version"`
	Jobs      map[string]Job         `yaml:"jobs"`
	Workflows map[string]Workflow    `yaml:"workflows"`
	Executors map[string]Executor    `yaml:"executors"`
	Commands  map[string]Command     `yaml:"commands"`
	Orbs      map[string]interface{} `yaml:"orbs"`
}

// Job represents a CircleCI job
type Job struct {
	Description string                 `yaml:"description"`
	Docker      []DockerConfig         `yaml:"docker"`
	Machine     interface{}            `yaml:"machine"`
	MacOS       interface{}            `yaml:"macos"`
	Executor    string                 `yaml:"executor"`
	Steps       []interface{}          `yaml:"steps"`
	Environment map[string]interface{} `yaml:"environment"`
	WorkingDir  string                 `yaml:"working_directory"`
	Parallelism int                    `yaml:"parallelism"`
}

// Workflow represents a CircleCI workflow
type Workflow struct {
	Jobs     []interface{} `yaml:"jobs"`
	Triggers []interface{} `yaml:"triggers"`
}

// WorkflowJob represents a job within a workflow (can be string or object)
type WorkflowJob struct {
	Name     string
	Requires []string
	Context  []string
	Filters  map[string]interface{}
}

// Executor represents a reusable executor
type Executor struct {
	Docker      []DockerConfig         `yaml:"docker"`
	Machine     interface{}            `yaml:"machine"`
	MacOS       interface{}            `yaml:"macos"`
	Environment map[string]interface{} `yaml:"environment"`
	WorkingDir  string                 `yaml:"working_directory"`
}

// DockerConfig represents Docker configuration
type DockerConfig struct {
	Image       string                 `yaml:"image"`
	Name        string                 `yaml:"name"`
	Entrypoint  []string               `yaml:"entrypoint"`
	Command     []string               `yaml:"command"`
	User        string                 `yaml:"user"`
	Environment map[string]interface{} `yaml:"environment"`
	Auth        map[string]interface{} `yaml:"auth"`
}

// Command represents a reusable command
type Command struct {
	Description string        `yaml:"description"`
	Parameters  interface{}   `yaml:"parameters"`
	Steps       []interface{} `yaml:"steps"`
}

// RunStep represents a run step
type RunStep struct {
	Name               string                 `yaml:"name"`
	Command            string                 `yaml:"command"`
	Shell              string                 `yaml:"shell"`
	Environment        map[string]interface{} `yaml:"environment"`
	Background         bool                   `yaml:"background"`
	WorkingDirectory   string                 `yaml:"working_directory"`
	NoOutputTimeout    string                 `yaml:"no_output_timeout"`
	When               string                 `yaml:"when"`
}

// Analysis represents the analysis results
type Analysis struct {
	Config           *Config
	JobUsage         map[string]int
	JobDependencies  map[string][]string
	CommandPatterns  map[string]PatternCount
	ExecutorUsage    map[string][]string
	TotalJobs        int
	TotalWorkflows   int
	GeneratedAt      time.Time
}

// PatternCount tracks pattern usage
type PatternCount struct {
	Count int
	Jobs  []string
}

// JobAnalysis represents analysis for a single job
type JobAnalysis struct {
	Name         string
	Description  string
	Commands     []string
	DockerImages []string
	Executor     string
	Dependencies []string
	UsageCount   int
	Patterns     map[string]int
}

// WorkflowAnalysis represents analysis for a single workflow
type WorkflowAnalysis struct {
	Name string
	Jobs []WorkflowJob
}