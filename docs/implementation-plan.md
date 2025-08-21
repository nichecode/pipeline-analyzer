# Pipeline Analyzer Implementation Plan

## Overview
Build a Go utility to replace the bash script for analyzing CircleCI configs and generating go-task migration documentation. Focus on practical utility over enterprise-grade robustness.

## Project Structure
```
pipeline-analyzer/
├── cmd/
│   └── analyze/           # Single main command
├── internal/
│   ├── circleci/          # CircleCI parsing & analysis
│   ├── gotask/            # go-task generation  
│   ├── markdown/          # Markdown output
│   └── fs/                # File operations
└── testdata/             # Test fixtures
```

## Implementation Phases

### Phase 1: Foundation & CircleCI Parser
**Goal:** Replace fragile `yq` parsing with proper Go YAML handling

**Prompts:**
1. **Setup basic project structure and dependencies**
   - Initialize the directory structure
   - Add necessary dependencies (yaml.v3, filepath, etc.)
   - Create basic main.go in cmd/analyze

2. **Define CircleCI types**
   - Create internal/circleci/types.go with structs for:
     - Config (top-level)
     - Job (with steps, docker, executor)
     - Workflow (with jobs array/map)
     - Step (run commands, checkout, etc.)
   - Handle the tricky workflow jobs structure (string vs object)

3. **Implement CircleCI parser**
   - Create internal/circleci/parser.go
   - ParseConfig(filePath) function
   - Handle YAML unmarshaling errors gracefully
   - Test with real CircleCI config

4. **Basic file operations**
   - Create internal/fs/ with simple read/write functions
   - ReadFile(path) and WriteFile(path, content)
   - CreateDirectory(path) for output structure

### Phase 2: Analysis Functions
**Goal:** Extract meaningful insights from parsed CircleCI config

**Prompts:**
5. **Job analysis functions**
   - Create internal/circleci/analyzer.go
   - ExtractJobs(config) - get all unique job definitions
   - AnalyzeJobCommands(job) - extract run commands and patterns
   - ClassifyCommands(commands) - docker, scripts, npm, etc.

6. **Workflow analysis**
   - ExtractWorkflows(config) 
   - AnalyzeJobDependencies(workflow) - requires analysis
   - CountJobUsage(workflows) - frequency analysis

7. **Cross-job analysis**
   - Create internal/analysis/ package
   - FindCommonPatterns(jobs) - shared command patterns
   - IdentifySharedDependencies(workflows) - jobs that others depend on
   - GenerateUsageStats(jobs, workflows) - usage frequency

### Phase 3: Markdown Generation
**Goal:** Generate the interconnected documentation files

**Prompts:**
8. **Markdown templates and generators**
   - Create internal/markdown/generator.go
   - GenerateJobMarkdown(job) - individual job files
   - GenerateWorkflowMarkdown(workflow) - workflow files
   - Template functions for consistent navigation links

9. **Summary generators**
   - GenerateJobUsageMarkdown(analysis)
   - GenerateCommandsAnalysisMarkdown(analysis)
   - GenerateAllJobsIndex(jobs)
   - GenerateMigrationChecklist(analysis)

10. **Main README and navigation**
    - GenerateMainReadme(analysis) - overview with stats
    - Ensure all cross-references work correctly
    - Generate directory structure with proper links

### Phase 4: go-task Generation
**Goal:** Auto-generate go-task files from CircleCI analysis

**Prompts:**
11. **go-task types and structure**
    - Create internal/gotask/types.go
    - Define Task, Taskfile structures matching go-task schema
    - Handle the functional grouping (build.yml, test.yml, etc.)

12. **Task generation logic**
    - Create internal/gotask/generator.go
    - ConvertJobToTask(job) - transform CircleCI job to go-task
    - GroupTasksByFunction(tasks) - organize into functional files
    - GenerateTaskfileYAML(tasks) - output proper YAML

13. **Template system for go-task**
    - Create internal/gotask/template.go
    - Templates for main Taskfile.yml (orchestrator)
    - Templates for functional groups (build.yml, test.yml, etc.)
    - Handle dependency mapping from CircleCI requires

### Phase 5: CLI and Integration
**Goal:** Bring it all together with a clean CLI interface

**Prompts:**
14. **Command structure and flags**
    - Enhance cmd/analyze/main.go
    - Add flags: --config, --output-dir, --generate-tasks
    - Basic error handling and user feedback

15. **End-to-end integration**
    - Wire all components together
    - Test with real CircleCI config
    - Ensure output matches bash script functionality
    - Add basic validation

16. **Polish and documentation**
    - Add usage examples
    - Handle edge cases discovered during testing
    - Basic performance improvements if needed

## Success Criteria
- Parse CircleCI configs without fragile shell commands
- Generate same markdown documentation as bash script
- Auto-generate go-task files organized by function
- Reduce migration risk through better analysis
- Run reliably on real CircleCI configs

## Non-Goals
- Enterprise-grade error handling
- Support for every CircleCI feature
- Complex CLI with subcommands
- Community distribution
- Extensive test coverage

## Key Implementation Notes
- Use state-based testing for pure functions
- Keep functions small and focused
- Separate IO from business logic
- Build incrementally and test with real configs
- Prioritize working over perfect