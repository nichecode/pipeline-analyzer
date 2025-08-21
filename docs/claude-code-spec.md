# Pipeline Analyzer - Claude Code Implementation Specification

## Overview
Implement a complete Go utility that replaces the provided bash script for analyzing CircleCI configs and generating migration documentation. Build everything in one session.

## Project Structure
```
pipeline-analyzer/
├── cmd/analyze/main.go           # CLI entry point
├── internal/
│   ├── circleci/
│   │   ├── types.go              # CircleCI config structs
│   │   ├── parser.go             # YAML parsing
│   │   └── analyzer.go           # Analysis functions
│   ├── markdown/
│   │   ├── generator.go          # Markdown generation
│   │   └── templates.go          # Template functions
│   └── fs/
│       └── writer.go             # File operations
├── go.mod                        # Dependencies
└── testdata/
    └── sample-config.yml         # Test CircleCI config
```

## Dependencies Required
```go
// go.mod dependencies
gopkg.in/yaml.v3  // YAML parsing
flag              // CLI flags (stdlib)
fmt, os, path/filepath, strings, sort // stdlib
```

## Core Types (internal/circleci/types.go)

```go
// Key structs needed to parse CircleCI YAML
type Config struct {
    Version   string                 `yaml:"version"`
    Jobs      map[string]Job         `yaml:"jobs"`
    Workflows map[string]Workflow    `yaml:"workflows"`
    Executors map[string]Executor    `yaml:"executors"`
}

type Job struct {
    Description string                 `yaml:"description"`
    Docker      []DockerConfig        `yaml:"docker"`
    Executor    string                `yaml:"executor"`
    Steps       []Step                `yaml:"steps"`
    Parameters  map[string]interface{} `yaml:"parameters"`
}

type DockerConfig struct {
    Image string `yaml:"image"`
}

type Step struct {
    Run      interface{} `yaml:"run"`      // Can be string or RunConfig
    Checkout interface{} `yaml:"checkout"` // Handle checkout steps
}

type RunConfig struct {
    Command string `yaml:"command"`
    Name    string `yaml:"name"`
}

type Workflow struct {
    Jobs []interface{} `yaml:"jobs"` // Can be strings or job configs
}

type Executor struct {
    Docker []DockerConfig `yaml:"docker"`
}

// Analysis result types
type Analysis struct {
    Jobs           map[string]Job
    Workflows      map[string]Workflow
    Executors      map[string]Executor
    JobUsage       map[string]int
    Dependencies   map[string][]string
    CommandStats   CommandStats
}

type CommandStats struct {
    DockerComposeCount int
    DockerRunCount     int
    ScriptCount        int
    NpmCount           int
    CommandFrequency   map[string]int
}
```

## Core Functions Required

### Parser (internal/circleci/parser.go)
```go
func ParseConfig(filePath string) (*Config, error)
// Parse YAML file into Config struct
```

### Analyzer (internal/circleci/analyzer.go)
```go
func AnalyzeConfig(config *Config) *Analysis
func ExtractRunCommands(job Job) []string
func ClassifyCommands(commands []string) CommandStats
func AnalyzeJobUsage(workflows map[string]Workflow) map[string]int
func ExtractDependencies(workflows map[string]Workflow) map[string][]string
```

### Markdown Generator (internal/markdown/generator.go)
```go
func GenerateAnalysisReport(analysis *Analysis, outputDir string) error
func GenerateJobMarkdown(name string, job Job, outputDir string) error
func GenerateWorkflowMarkdown(name string, workflow Workflow, outputDir string) error
func GenerateMainReadme(analysis *Analysis, outputDir string) error
func GenerateMigrationChecklist(analysis *Analysis, outputDir string) error
```

## Required Output Files (matching bash script exactly)

### Directory Structure
```
circleci-analysis/
├── README.md                     # Main overview
├── MIGRATION-CHECKLIST.md        # Migration guide
├── jobs/                         # Individual job files
│   ├── {job-name}.md
├── workflows/                    # Individual workflow files
│   ├── {workflow-name}.md
└── summaries/                    # Analysis summaries
    ├── all-jobs.md
    ├── job-usage.md
    ├── commands.md
    ├── docker-and-scripts.md
    ├── executors-and-images.md
    └── workflows.md
```

## Key Functionality to Replicate

### Job Analysis
- Extract all unique job definitions
- Generate individual markdown files for each job with:
  - Description, executor, Docker images
  - Full steps in YAML format
  - Extracted run commands
  - Command analysis (count of docker, script, npm commands)
  - Navigation links

### Workflow Analysis  
- Extract workflow structure
- Handle both simple job strings and complex job objects
- Extract job dependencies (requires field)
- Generate workflow markdown with job lists and dependencies

### Command Analysis
- Extract all run commands from all jobs
- Count patterns: docker-compose, docker run, scripts, npm
- Generate frequency table of command patterns (first word)
- Create comprehensive commands analysis markdown

### Usage Analysis
- Count how many times each job is used across workflows
- Identify jobs that others depend on
- Create dependency analysis tables

### Migration Checklist
- Generate suggested go-task structure
- Include job counts and statistics
- Provide step-by-step migration guidance

## CLI Interface
```bash
# Basic usage
./analyze [path-to-config.yml]

# Default to .circleci/config.yml if no path provided
# Output to ./circleci-analysis/ directory
# Print progress and summary stats
```

## Error Handling
- Check if config file exists
- Handle YAML parsing errors gracefully
- Create output directories as needed
- Print helpful error messages

## Implementation Notes

### Workflow Job Parsing
The trickiest part is parsing workflow jobs that can be either:
```yaml
# Simple string
jobs:
  - job-name

# Complex object  
jobs:
  - job-name:
      requires: [other-job]
      context: some-context
```

Handle both cases when extracting job names and dependencies.

### Command Extraction
Extract run commands from steps, handling both formats:
```yaml
# Simple string
- run: "echo hello"

# Complex object
- run:
    command: "echo hello"
    name: "Say hello"
```

### Pattern Matching
Implement Go equivalents of the bash grep patterns:
- docker-compose references
- docker run references  
- script calls (./scripts/, scripts/, .sh files)
- npm/yarn/pnpm commands

## Success Criteria
The Go tool should produce identical markdown output to the bash script, with:
- Same file structure and naming
- Same content and analysis
- Same navigation links
- Equivalent statistics and tables
- Better error handling and performance

## Test Data
Include a sample CircleCI config in testdata/ to validate the implementation works correctly.

## Implementation Approach
1. Start with types and parser to handle YAML
2. Implement analysis functions 
3. Build markdown generators with templates
4. Wire everything together in main.go
5. Test with real config file