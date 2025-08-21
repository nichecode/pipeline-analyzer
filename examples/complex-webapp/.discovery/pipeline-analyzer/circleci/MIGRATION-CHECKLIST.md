# CircleCI → go-task Migration Checklist

**Generated:** 2025-08-21T07:59:05+01:00

## 🎯 Migration Overview

This checklist guides you through converting your CircleCI configuration to a local go-task setup.

### 📊 Current State

- **13 jobs** to migrate
- **3 workflows** to understand
- **13 jobs** use Docker
- **0 jobs** use other executors

## 🔄 Migration Steps

### 1. **Understand job dependencies**
Review [Job Usage Analysis](summaries/job-usage.md) to see which jobs depend on others.

### 2. **Examine Docker images and executors**
Check [Docker & Scripts](summaries/docker-and-scripts.md) for container requirements.

### 3. **Analyze command patterns**
Study [Commands Analysis](summaries/commands.md) to understand build patterns.

### 4. **Start with high-impact jobs**
Begin with the most frequently used jobs from the overview.

### 5. **Convert commands to go-task**
Transform run commands to task format - see individual job files.

### 6. **Test task equivalents** locally

## Key Files to Examine

- [Job Usage Analysis](summaries/job-usage.md) - How jobs are reused and their dependencies
- [All Jobs](summaries/all-jobs.md) - Complete job list with descriptions
- [Docker & Scripts](summaries/docker-and-scripts.md) - Docker/script patterns
- [Workflows](summaries/workflows.md) - Workflow structure and job orchestration

## Suggested go-task Structure

```yaml
version: '3'

tasks:
  build-backend:
    desc: "Migrated from CircleCI job"
    deps: [test-backend]
    cmds:
      - # Convert run commands from jobs/build-backend.md

  build-docker-images:
    desc: "Migrated from CircleCI job"
    deps: [build-frontend, build-backend, security-scan]
    cmds:
      - # Convert run commands from jobs/build-docker-images.md

  build-frontend:
    desc: "Migrated from CircleCI job"
    deps: [test-frontend]
    cmds:
      - # Convert run commands from jobs/build-frontend.md

  deploy-production:
    desc: "Migrated from CircleCI job"
    deps: [hold-for-approval]
    cmds:
      - # Convert run commands from jobs/deploy-production.md

  deploy-staging:
    desc: "Migrated from CircleCI job"
    deps: [build-docker-images, test-integration]
    cmds:
      - # Convert run commands from jobs/deploy-staging.md

  lint-backend:
    desc: "Migrated from CircleCI job"
    cmds:
      - # Convert run commands from jobs/lint-backend.md

  lint-frontend:
    desc: "Migrated from CircleCI job"
    cmds:
      - # Convert run commands from jobs/lint-frontend.md

  security-scan:
    desc: "Migrated from CircleCI job"
    cmds:
      - # Convert run commands from jobs/security-scan.md

  test-backend:
    desc: "Migrated from CircleCI job"
    deps: [lint-backend]
    cmds:
      - # Convert run commands from jobs/test-backend.md

  test-e2e:
    desc: "Migrated from CircleCI job"
    deps: [test-integration]
    cmds:
      - # Convert run commands from jobs/test-e2e.md

  test-frontend:
    desc: "Migrated from CircleCI job"
    deps: [lint-frontend]
    cmds:
      - # Convert run commands from jobs/test-frontend.md

  test-integration:
    desc: "Migrated from CircleCI job"
    deps: [test-frontend, test-backend]
    cmds:
      - # Convert run commands from jobs/test-integration.md

  test-performance:
    desc: "Migrated from CircleCI job"
    deps: [build-frontend]
    cmds:
      - # Convert run commands from jobs/test-performance.md

```

## Navigation

- [← Back to Overview](../README.md)
- [All Jobs](summaries/all-jobs.md)
- [Job Usage Analysis](summaries/job-usage.md)
