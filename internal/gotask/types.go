package gotask

import "time"

// Taskfile represents the root structure of a Taskfile.yml
type Taskfile struct {
	Version   string                    `yaml:"version"`
	Output    interface{}               `yaml:"output"`
	Method    string                    `yaml:"method"`
	Includes  map[string]Include        `yaml:"includes"`
	Vars      map[string]interface{}    `yaml:"vars"`
	Env       map[string]interface{}    `yaml:"env"`
	Tasks     map[string]Task           `yaml:"tasks"`
	Silent    bool                      `yaml:"silent"`
	Dotenv    []string                  `yaml:"dotenv"`
	Run       string                    `yaml:"run"`
	Interval  string                    `yaml:"interval"`
	Set       []string                  `yaml:"set"`
	Shopt     []string                  `yaml:"shopt"`
}

// Include represents an included Taskfile configuration
type Include struct {
	Taskfile   string                    `yaml:"taskfile"`
	Dir        string                    `yaml:"dir"`
	Optional   bool                      `yaml:"optional"`
	Flatten    bool                      `yaml:"flatten"`
	Internal   bool                      `yaml:"internal"`
	Aliases    []string                  `yaml:"aliases"`
	Excludes   []string                  `yaml:"excludes"`
	Vars       map[string]interface{}    `yaml:"vars"`
	Checksum   string                    `yaml:"checksum"`
}

// Task represents a single task definition
type Task struct {
	Cmds          []interface{}             `yaml:"cmds"`
	Cmd           string                    `yaml:"cmd"`
	Deps          []interface{}             `yaml:"deps"`
	Desc          string                    `yaml:"desc"`
	Summary       string                    `yaml:"summary"`
	Prompt        string                    `yaml:"prompt"`
	Aliases       []string                  `yaml:"aliases"`
	Sources       []string                  `yaml:"sources"`
	Generates     []string                  `yaml:"generates"`
	Status        []string                  `yaml:"status"`
	Preconditions []interface{}             `yaml:"preconditions"`
	Requires      map[string]string         `yaml:"requires"`
	Watch         bool                      `yaml:"watch"`
	Platforms     []string                  `yaml:"platforms"`
	Silent        bool                      `yaml:"silent"`
	Internal      bool                      `yaml:"internal"`
	Vars          map[string]interface{}    `yaml:"vars"`
	Env           map[string]interface{}    `yaml:"env"`
	Run           string                    `yaml:"run"`
	IgnoreError   bool                      `yaml:"ignore_error"`
}

// Command represents a command that can be either string or object
type Command struct {
	Cmd         string                    `yaml:"cmd"`
	Silent      bool                      `yaml:"silent"`
	IgnoreError bool                      `yaml:"ignore_error"`
	Platforms   []string                  `yaml:"platforms"`
	Set         []string                  `yaml:"set"`
	Shopt       []string                  `yaml:"shopt"`
	Defer       bool                      `yaml:"defer"`
	Vars        map[string]interface{}    `yaml:"vars"`
	Env         map[string]interface{}    `yaml:"env"`
}

// Dependency represents a task dependency
type Dependency struct {
	Task string                    `yaml:"task"`
	Vars map[string]interface{}    `yaml:"vars"`
}

// Precondition represents a precondition check
type Precondition struct {
	Sh  string `yaml:"sh"`
	Msg string `yaml:"msg"`
}

// Analysis represents the analysis results for a Taskfile
type Analysis struct {
	Taskfile           *Taskfile
	TaskUsage          map[string]int
	TaskDependencies   map[string][]string
	CommandPatterns    map[string]PatternCount
	VariableUsage      map[string][]string
	EnvironmentUsage   map[string][]string
	IncludeAnalysis    map[string]*IncludeAnalysis
	TotalTasks         int
	TotalIncludes      int
	CircularDeps       [][]string
	CriticalPath       []string
	OptimizationTips   []OptimizationTip
	GeneratedAt        time.Time
}

// PatternCount tracks pattern usage across tasks
type PatternCount struct {
	Count int
	Tasks []string
}

// IncludeAnalysis represents analysis of an included Taskfile
type IncludeAnalysis struct {
	Path         string
	Namespace    string
	TaskCount    int
	Dependencies []string
}

// TaskAnalysis represents detailed analysis of a single task
type TaskAnalysis struct {
	Name            string
	Description     string
	Summary         string
	Commands        []string
	Dependencies    []string
	Sources         []string
	Generates       []string
	Variables       map[string]interface{}
	Environment     map[string]interface{}
	Patterns        map[string]int
	UsageCount      int
	IsInternal      bool
	Platforms       []string
	Aliases         []string
	Preconditions   []string
	HasWatch        bool
	OptimizationOps []string
}

// DependencyGraph represents the task dependency structure
type DependencyGraph struct {
	Tasks    []string
	Edges    map[string][]string
	Levels   map[string]int
	Cycles   [][]string
}

// OptimizationTip represents a suggestion for improvement
type OptimizationTip struct {
	Type        string
	Task        string
	Message     string
	Severity    string
	Suggestion  string
}

// Variable represents variable usage tracking
type Variable struct {
	Name         string
	Value        interface{}
	Type         string
	UsedInTasks  []string
	IsGlobal     bool
	IsEnvironment bool
}

// CommandClassification represents command categorization
type CommandClassification struct {
	Category    string
	Tools       []string
	Complexity  int
	Risk        string
	Suggestions []string
}

// PerformanceMetrics represents performance-related analysis
type PerformanceMetrics struct {
	TasksWithCaching     int
	TasksWithSources     int
	TasksWithGenerates   int
	ParallelizableTasks  int
	OptimizationPotential float64
}

// SecurityAnalysis represents security-related findings
type SecurityAnalysis struct {
	PotentialSecrets   []string
	UnsafeCommands     []string
	ExternalDependencies []string
	Recommendations    []string
}

// ComparisonResult represents comparison between different task systems
type ComparisonResult struct {
	TaskEquivalencies  map[string]string
	MigrationEffort    string
	FeatureGaps        []string
	Recommendations    []string
	ComplexityScore    float64
}