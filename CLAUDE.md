# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Pipeline Analyzer is a Go tool that replaces a fragile bash script for analyzing CircleCI configurations. It parses YAML config files and generates comprehensive markdown documentation to help teams migrate from CircleCI to local development workflows (particularly go-task).

## Development Commands

Uses go-task for build automation. See `Taskfile.yml` for all available tasks.

### Common Tasks
```bash
# Show available tasks
task

# Install dependencies
task deps

# Build the analyzer tool
task build

# Run the analyzer
task run

# Run with custom config
task run -- path/to/config.yml

# Run tests
task test

# Format and lint code
task fmt
task lint

# Run all checks (fmt, lint, test)
task check

# Clean build artifacts
task clean
```

### Development Setup
```bash
# Initial project setup
task setup

# Update dependencies
task update
```

## Architecture

This is a CLI tool built in Go with a clean separation of concerns:

- **cmd/pipeline-analyzer/main.go** - CLI entry point with flag parsing
- **internal/circleci/** - CircleCI config parsing and analysis logic
- **internal/markdown/** - Markdown generation for documentation
- **internal/fs/** - File system operations and directory management

### Key Design Decisions

1. **YAML Parsing Strategy**: Uses `interface{}` types to handle CircleCI's flexible job syntax where jobs can be simple strings or complex objects with dependencies
2. **Pattern Detection**: Replaces fragile bash/grep logic with Go string matching and regex patterns
3. **Output Structure**: Generates analysis in `.discovery/pipeline-analyzer/circleci/` for repo-based analysis

## Expected Code Structure

The project follows this planned structure based on the specification:

```
cmd/pipeline-analyzer/main.go # CLI with --config, --output-dir flags
internal/
├── circleci/
│   ├── types.go              # Config, Job, Workflow structs
│   ├── parser.go             # YAML unmarshaling
│   └── analyzer.go           # Extract commands, count patterns
├── markdown/
│   ├── generator.go          # Job/workflow markdown generation
│   └── templates.go          # Reusable template functions
└── fs/writer.go              # Directory creation, file writing
```

## Dependencies

- `gopkg.in/yaml.v3` - YAML parsing for CircleCI configs
- Standard library packages: `flag`, `fmt`, `os`, `path/filepath`, `strings`, `sort`

## Reference Implementation

The bash script being replaced is located at `docs/circleci-analyzer.sh`. This script demonstrates the exact analysis patterns and output format that the Go implementation should replicate.

## Output Requirements

The tool generates structured analysis in:
- `.discovery/pipeline-analyzer/circleci/` (when in git repo)
- `circleci-analysis/` (fallback)

Key output files include:
- Individual job and workflow markdown files
- Usage statistics and dependency analysis
- Command pattern analysis (docker, scripts, npm counts)
- Migration checklist with go-task suggestions
- Cross-referenced navigation between files