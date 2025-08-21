# Task Optimization Guide

**Generated:** 2025-08-21T10:47:48+01:00

## 🎯 Optimization Overview

This guide provides specific recommendations to improve your Taskfile performance, maintainability, and reliability.

### 📊 Current Performance Status

- **Tasks with caching optimization:** 3/23 (13.0%)
- **Tasks with source tracking:** 4
- **Tasks with output tracking:** 3
- **Parallelizable tasks:** 4
- **Optimization potential:** 87.0%

## ⚡ Performance Optimizations

These changes can significantly improve task execution speed:

### Task: `lint-backend`

**Opportunity:** Task could benefit from caching optimization

**Implementation:**
```yaml
lint-backend:
  desc: "Task description"
  sources:
    - "src/**/*.go"  # Add relevant source patterns
  generates:
    - "dist/app"     # Add output files
  cmds:
    - # existing commands
```

### Task: `build-docker`

**Opportunity:** Task could benefit from caching optimization

**Implementation:**
```yaml
build-docker:
  desc: "Task description"
  sources:
    - "src/**/*.go"  # Add relevant source patterns
  generates:
    - "dist/app"     # Add output files
  cmds:
    - # existing commands
```

### Task: `deploy-prod`

**Opportunity:** Task could benefit from caching optimization

**Implementation:**
```yaml
deploy-prod:
  desc: "Task description"
  sources:
    - "src/**/*.go"  # Add relevant source patterns
  generates:
    - "dist/app"     # Add output files
  cmds:
    - # existing commands
```

### Task: `security-scan`

**Opportunity:** Task could benefit from caching optimization

**Implementation:**
```yaml
security-scan:
  desc: "Task description"
  sources:
    - "src/**/*.go"  # Add relevant source patterns
  generates:
    - "dist/app"     # Add output files
  cmds:
    - # existing commands
```

### Task: `test-integration`

**Opportunity:** Task could benefit from caching optimization

**Implementation:**
```yaml
test-integration:
  desc: "Task description"
  sources:
    - "src/**/*.go"  # Add relevant source patterns
  generates:
    - "dist/app"     # Add output files
  cmds:
    - # existing commands
```

### Task: `lint-frontend`

**Opportunity:** Task could benefit from caching optimization

**Implementation:**
```yaml
lint-frontend:
  desc: "Task description"
  sources:
    - "src/**/*.go"  # Add relevant source patterns
  generates:
    - "dist/app"     # Add output files
  cmds:
    - # existing commands
```

### Task: `deploy-staging`

**Opportunity:** Task could benefit from caching optimization

**Implementation:**
```yaml
deploy-staging:
  desc: "Task description"
  sources:
    - "src/**/*.go"  # Add relevant source patterns
  generates:
    - "dist/app"     # Add output files
  cmds:
    - # existing commands
```

### Task: `performance-test`

**Opportunity:** Task could benefit from caching optimization

**Implementation:**
```yaml
performance-test:
  desc: "Task description"
  sources:
    - "src/**/*.go"  # Add relevant source patterns
  generates:
    - "dist/app"     # Add output files
  cmds:
    - # existing commands
```

### Task: `clean`

**Opportunity:** Task could benefit from caching optimization

**Implementation:**
```yaml
clean:
  desc: "Task description"
  sources:
    - "src/**/*.go"  # Add relevant source patterns
  generates:
    - "dist/app"     # Add output files
  cmds:
    - # existing commands
```

### Task: `dev`

**Opportunity:** Task could benefit from caching optimization

**Implementation:**
```yaml
dev:
  desc: "Task description"
  sources:
    - "src/**/*.go"  # Add relevant source patterns
  generates:
    - "dist/app"     # Add output files
  cmds:
    - # existing commands
```

### Task: `test`

**Opportunity:** Task could benefit from caching optimization

**Implementation:**
```yaml
test:
  desc: "Task description"
  sources:
    - "src/**/*.go"  # Add relevant source patterns
  generates:
    - "dist/app"     # Add output files
  cmds:
    - # existing commands
```

### Task: `test-e2e`

**Opportunity:** Task could benefit from caching optimization

**Implementation:**
```yaml
test-e2e:
  desc: "Task description"
  sources:
    - "src/**/*.go"  # Add relevant source patterns
  generates:
    - "dist/app"     # Add output files
  cmds:
    - # existing commands
```

## 🚀 Advanced Optimizations

### Parallel Execution

Tasks without dependencies can run in parallel. Consider grouping related independent tasks:

```yaml
test-all:
  desc: "Run all tests in parallel"
  deps:
    - task: test-unit
    - task: test-integration
    - task: test-e2e
```

### Critical Path Analysis

Your longest dependency chain is:

`deploy-prod` → `test-e2e` → `start-services` → `build-docker` → `build` → `build-backend` → `install-backend`

Focus optimization efforts on these tasks for maximum impact.

## 💡 Best Practices

### Use Specific Source Patterns

Avoid overly broad patterns that might cause unnecessary rebuilds

```yaml
sources: ["src/**/*.go", "go.mod", "go.sum"]
```

### Output Specific Files

List actual output files rather than directories when possible

```yaml
generates: ["dist/app", "dist/version.txt"]
```

### Group Related Tasks

Use task namespaces with includes for better organization

```yaml
includes: { test: "./tasks/test.yml" }
```

### Use Status Checks

Skip tasks when conditions are already met

```yaml
status: ["test -f dist/app"]
```

## Navigation

- [← Back to Overview](README.md)
- [Performance Metrics](summaries/performance.md)
- [Dependency Graph](tasks/dependency-graph.md)
