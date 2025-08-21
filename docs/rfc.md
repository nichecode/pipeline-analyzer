# Pipeline Analyzer - RFC

## Problem Statement
Replace a fragile bash script that analyzes CircleCI configs with a robust Go utility. The bash script uses `yq` shell commands for YAML parsing and generates interconnected markdown documentation for migration planning.

**Reference Implementation:** See `docs/circleci-analyzer.sh` for the exact bash script to replicate in Go.

## Current Pain Points
- Fragile shell `yq` parsing that breaks on complex YAML structures
- Complex workflow job extraction logic (jobs can be strings or objects)
- Manual string manipulation for pattern detection
- No type safety or proper error handling

## Success Criteria
Generate equivalent analysis output in a clean `.discovery/pipeline-analyzer/circleci/` directory structure with:
- Individual job and workflow markdown files
- Usage statistics and dependency analysis  
- Command pattern analysis (docker, scripts, npm counts)
- Migration checklist with go-task suggestions
- Cross-referenced navigation between files

## Repository Context
- Tool should detect if running in a git repository
- If in repo: output to `.discovery/pipeline-analyzer/circleci/`
- If not in repo: fall back to `circleci-analysis/` directory
- Create overview files at discovery and tool levels

## Architecture Decisions & Rationale

### YAML Parsing Strategy
**Challenge:** CircleCI workflow jobs can be either simple strings or complex objects:
```yaml
jobs:
  - simple-job-name
  - complex-job:
      requires: [other-job]
      context: some-context
```

**Approach:** Use `interface{}` in Go structs, then type switch during analysis
**Why:** Handles both cases cleanly without complex custom unmarshaling

### Command Extraction Pattern
**Challenge:** Run commands have multiple formats:
```yaml
- run: "echo hello"           # Simple string
- run:                        # Complex object
    command: "echo hello"
    name: "Say hello"
```

**Approach:** Handle both in extraction logic with type assertion
**Why:** Preserves all command information while staying flexible

### Pattern Detection Strategy
**Challenge:** Bash script uses grep for docker-compose, script calls, etc.
**Approach:** Use Go string matching with regex patterns
**Why:** More reliable than shell commands, easier to test and extend

## Expected Code Structure
```
cmd/analyze/main.go           # CLI with --config, --output-dir flags
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

## Key Dependencies
- `gopkg.in/yaml.v3` for YAML parsing
- Standard library: `flag`, `fmt`, `os`, `path/filepath`, `strings`, `sort`

## Output Files Required
Match the bash script output exactly:

```
.discovery/pipeline-analyzer/circleci/
├── README.md                 # Overview with stats
├── MIGRATION-CHECKLIST.md    # Migration guidance
├── jobs/{job-name}.md        # Individual job details
├── workflows/{workflow}.md   # Workflow structure
└── summaries/
    ├── all-jobs.md
    ├── job-usage.md          # Usage frequency tables
    ├── commands.md           # Pattern analysis
    ├── docker-and-scripts.md
    ├── executors-and-images.md
    └── workflows.md
```

## Analysis Features Needed
1. **Job Usage Analysis:** Count how many workflows use each job
2. **Dependency Analysis:** Extract `requires` relationships
3. **Command Pattern Analysis:** Count docker, scripts, npm usage per job
4. **Command Frequency:** Most common command patterns (first word)
5. **Executor Analysis:** Which jobs use which Docker images/executors

## CLI Behavior
```bash
./analyze [path-to-config.yml]   # Defaults to .circleci/config.yml
```

Print progress, create directories as needed, handle errors gracefully.

## Implementation Notes
- Keep functions small and focused
- Separate IO from business logic for testability  
- Use descriptive variable names and clear error messages
- Include sample CircleCI config in testdata/ for validation

## Validation
Test with real CircleCI configs to ensure:
- All job types parse correctly
- Workflow dependencies extract properly
- Generated markdown has working cross-references
- Output matches bash script analysis quality