# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Pipeline Analyzer is a comprehensive Go tool that auto-discovers and analyzes multiple build tools in repositories (CircleCI, GitHub Actions, go-task, npm, Docker, etc.). It generates detailed migration guides specifically focused on refactoring CI/CD pipelines to use go-task as the unified command runner, enabling local execution and CI platform portability.

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

This is a CLI tool built in Go with auto-discovery and multi-platform analysis:

- **cmd/pipeline-analyzer/main.go** - CLI entry point (auto-discovery only mode)
- **internal/discovery/** - Repository scanning and build tool discovery  
- **internal/circleci/** - CircleCI config parsing and analysis
- **internal/githubactions/** - GitHub Actions workflow parsing and analysis
- **internal/gotask/** - go-task analysis with optimization recommendations
- **internal/shared/** - Common command pattern recognition and utilities

### Key Design Decisions

1. **Auto-Discovery Architecture**: Scans repositories for multiple build tools automatically
2. **YAML Parsing Strategy**: Uses `interface{}` types to handle flexible config syntaxes across different platforms
3. **go-task Migration Focus**: All analysis targets practical go-task refactoring opportunities
4. **Unified Command Runner Pattern**: Promotes CI → go-task → commands architecture
5. **Self-Contained Packages**: Each build tool analyzer is independent and self-contained

## Current Code Structure

```
cmd/pipeline-analyzer/main.go      # Auto-discovery CLI (repository path only)
internal/
├── discovery/
│   ├── scanner.go                 # Multi-tool repository scanning  
│   └── analyzer.go                # Coordinates analysis of all discovered tools
├── circleci/
│   ├── types.go                   # CircleCI config structs
│   ├── parser.go                  # CircleCI YAML parsing
│   ├── analyzer.go                # Job analysis and recommendations
│   ├── generator.go               # Markdown generation
│   ├── templates.go               # CircleCI-specific templates
│   └── writer.go                  # File output management
├── githubactions/
│   ├── types.go                   # GitHub Actions workflow structs
│   ├── parser.go                  # Workflow YAML parsing
│   ├── analyzer.go                # Job/step analysis with go-task opportunities
│   ├── markdown.go                # GitHub Actions markdown generation
│   └── writer.go                  # File output management
├── gotask/
│   ├── types.go                   # Taskfile structs
│   ├── parser.go                  # Taskfile YAML parsing
│   ├── analyzer.go                # Task dependency and optimization analysis
│   ├── markdown.go                # Task documentation generation
│   ├── summaries.go               # Performance and usage summaries
│   └── writer.go                  # File output management
└── shared/
    ├── patterns.go                # Command pattern recognition (npm, docker, etc.)
    └── utils.go                   # Common utilities
```

## Dependencies

- `gopkg.in/yaml.v3` - YAML parsing for all config files (CircleCI, GitHub Actions, Taskfiles)
- Standard library packages: `flag`, `fmt`, `os`, `path/filepath`, `strings`, `sort`, `regexp`, `time`

## Build Tool Support

**Fully Analyzed:**
- **CircleCI** - Complete job, workflow, executor, and command analysis
- **GitHub Actions** - Workflow, job, step, and action analysis with go-task migration guides
- **go-task** - Task dependency analysis, performance optimization, and caching recommendations
- **npm** - Package.json script analysis

**Detected Only:** 
- Docker, Python, PHP, Java, Rust, Maven, Gradle, Terraform, Makefile

## Output Requirements

The tool generates structured analysis in `.discovery/pipeline-analyzer/` containing:

**Per Build Tool:**
- `README.md` - Overview and quick start guide
- `migration-checklist.md` / `optimization-guide.md` - Actionable improvement guides
- `jobs/` or `tasks/` - Individual component analysis
- `summaries/` - Command patterns, usage statistics, optimization opportunities

**Cross-Tool:**
- Main `README.md` - Overview of all discovered tools
- `index.html` - Interactive navigation interface
- Links between equivalent functionality across tools

**Key Focus Areas:**
- go-task migration opportunities in CI tools
- Command pattern analysis for task extraction  
- Local vs CI execution guidance
- Docker image and dependency optimization
- Cross-platform build portability