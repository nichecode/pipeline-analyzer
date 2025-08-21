# Pipeline Analyzer

A comprehensive build tool analyzer that helps teams refactor CI/CD pipelines for local development and cross-platform portability. Automatically discovers CircleCI, GitHub Actions, go-task, npm, and other build tools in repositories, then generates detailed migration guides for consolidating CI commands into go-task for local execution.

## 🎯 Purpose

**Build System Refactoring** - Transform your CI/CD setup to this architecture:
```
CircleCI/GitHub Actions → calls `task build` → Taskfile.yml runs commands
```

**Benefits:**
- ✅ **Locally executable** - Run the same commands locally as in CI
- ✅ **CI-agnostic** - Easy to switch between GitHub Actions ↔ CircleCI  
- ✅ **Portable and testable** - No more "works in CI but not locally"

## 🚀 Quick Start

### Installation

```bash
go install github.com/nichecode/pipeline-analyzer/cmd/pipeline-analyzer@latest
```

### Usage

```bash
# Analyze current repository
pipeline-analyzer .

# Analyze specific repository
pipeline-analyzer /path/to/repo
```

The tool will:
1. **Auto-discover** all build tools in the repository
2. **Create `.discovery/pipeline-analyzer/`** with comprehensive analysis
3. **Generate migration guides** for converting to go-task
4. **Provide HTML and markdown** navigation

## 📊 Supported Build Tools

- **CircleCI** - Complete workflow and job analysis with Docker image tracking
- **GitHub Actions** - Workflow analysis with go-task migration recommendations  
- **go-task** - Task optimization and dependency analysis
- **npm** - Package.json script analysis
- **Docker** - Container and compose file detection
- **Python, PHP, Java, Rust** - Package manager detection

## 🔍 Example Output

After running `pipeline-analyzer .`, you'll get:

```
.discovery/pipeline-analyzer/
├── README.md                 # Main overview
├── index.html                # Interactive navigation  
├── circleci/
│   ├── README.md            # CircleCI analysis
│   ├── migration-checklist.md
│   └── summaries/           # Job usage, Docker images, commands
├── github-actions/
│   ├── README.md            # GitHub Actions analysis  
│   ├── summaries/
│   │   └── go-task-migration.md  # ⭐ Key migration guide
├── gotask/
│   ├── README.md            # go-task analysis
│   └── optimization-guide.md
└── npm/
    └── README.md            # npm scripts analysis
```

## 💡 Key Features

- **Migration-First Analysis** - Focuses on actionable go-task conversions
- **Command Pattern Detection** - Identifies commands suitable for task extraction
- **Cross-Platform Comparison** - Shows equivalent tasks across different tools
- **Local Testing Guidance** - Instructions for running tasks locally
- **HTML Navigation** - Interactive browsing of analysis results
