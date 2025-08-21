package githubactions

import "time"

// Workflow represents a GitHub Actions workflow file
type Workflow struct {
	Name string                 `yaml:"name,omitempty"`
	On   interface{}            `yaml:"on,omitempty"` // Can be string, array, or object
	Env  map[string]string      `yaml:"env,omitempty"`
	Jobs map[string]Job         `yaml:"jobs"`
}

// Job represents a job within a workflow
type Job struct {
	Name         string              `yaml:"name,omitempty"`
	RunsOn       interface{}         `yaml:"runs-on,omitempty"` // Can be string or array
	Needs        interface{}         `yaml:"needs,omitempty"`   // Can be string or array
	If           string              `yaml:"if,omitempty"`
	Env          map[string]string   `yaml:"env,omitempty"`
	Strategy     Strategy            `yaml:"strategy,omitempty"`
	Container    interface{}         `yaml:"container,omitempty"` // Can be string or object
	Services     map[string]Service  `yaml:"services,omitempty"`
	Steps        []Step              `yaml:"steps"`
	TimeoutMinutes int               `yaml:"timeout-minutes,omitempty"`
	Permissions  interface{}         `yaml:"permissions,omitempty"`
}

// Strategy defines job strategy (matrix, fail-fast, etc.)
type Strategy struct {
	Matrix     interface{} `yaml:"matrix,omitempty"`
	FailFast   bool        `yaml:"fail-fast,omitempty"`
	MaxParallel int        `yaml:"max-parallel,omitempty"`
}

// Service represents a service container
type Service struct {
	Image       string            `yaml:"image"`
	Env         map[string]string `yaml:"env,omitempty"`
	Ports       []string          `yaml:"ports,omitempty"`
	Volumes     []string          `yaml:"volumes,omitempty"`
	Options     string            `yaml:"options,omitempty"`
	Credentials map[string]string `yaml:"credentials,omitempty"`
}

// Step represents a step within a job
type Step struct {
	Name             string            `yaml:"name,omitempty"`
	Id               string            `yaml:"id,omitempty"`
	If               string            `yaml:"if,omitempty"`
	Uses             string            `yaml:"uses,omitempty"`
	Run              string            `yaml:"run,omitempty"`
	Shell            string            `yaml:"shell,omitempty"`
	With             map[string]interface{} `yaml:"with,omitempty"`
	Env              map[string]string `yaml:"env,omitempty"`
	WorkingDirectory string            `yaml:"working-directory,omitempty"`
	ContinueOnError  bool              `yaml:"continue-on-error,omitempty"`
	TimeoutMinutes   int               `yaml:"timeout-minutes,omitempty"`
}

// AnalysisResult represents the analysis of a GitHub Actions workflow
type AnalysisResult struct {
	Config         *Workflow
	FilePath       string
	GeneratedAt    time.Time
	Jobs           []JobAnalysis
	Workflows      int
	TotalSteps     int
	ActionUsage    map[string]int
	RunnerUsage    map[string]int
	ServiceUsage   map[string]int
	CommandPatterns map[string][]string
	Recommendations []string
	Issues          []string
}

// JobAnalysis contains analysis for a specific job
type JobAnalysis struct {
	Name            string
	Runner          string
	StepCount       int
	RunCommands     []string
	ActionsUsed     []string
	Dependencies    []string
	EstimatedTime   string
	CachingEnabled  bool
	SecurityIssues  []string
	Recommendations []string
}