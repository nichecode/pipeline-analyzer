# Go-Task Analysis Report

**Generated:** 2025-08-21T07:59:05+01:00
**Taskfile:** /Users/nicholas/Projects/pipeline-analyzer/examples/complex-webapp/Taskfile.yml

## 📊 Overview

- **Tasks:** 23
- **Includes:** 0
- **Circular Dependencies:** None ✅
- **Critical Path Length:** 7 tasks
- **Tasks with Caching:** 3/23 (13.0%)

## 🚀 Quick Start

1. **[📋 Optimization Guide](OPTIMIZATION-GUIDE.md)** - Performance improvement recommendations
2. **[📈 Task Usage Analysis](summaries/task-usage.md)** - Task dependency patterns
3. **[⚡ Command Analysis](summaries/commands.md)** - Command patterns and tools
4. **[🔗 Dependency Graph](tasks/dependency-graph.md)** - Visual task relationships

## 📁 Directory Structure

### Tasks
Individual task analysis with commands and optimization opportunities:

- [tasks/build.md](tasks/build.md)
- [tasks/build-backend.md](tasks/build-backend.md)
- [tasks/build-docker.md](tasks/build-docker.md)
- [tasks/build-frontend.md](tasks/build-frontend.md)
- [tasks/clean.md](tasks/clean.md)
- [tasks/deploy-prod.md](tasks/deploy-prod.md)
- [tasks/deploy-staging.md](tasks/deploy-staging.md)
- [tasks/dev.md](tasks/dev.md)
- [tasks/install.md](tasks/install.md)
- [tasks/install-backend.md](tasks/install-backend.md)
- [tasks/install-frontend.md](tasks/install-frontend.md)
- [tasks/lint.md](tasks/lint.md)
- [tasks/lint-backend.md](tasks/lint-backend.md)
- [tasks/lint-frontend.md](tasks/lint-frontend.md)
- [tasks/performance-test.md](tasks/performance-test.md)
- [tasks/security-scan.md](tasks/security-scan.md)
- [tasks/start-services.md](tasks/start-services.md)
- [tasks/stop-services.md](tasks/stop-services.md)
- [tasks/test.md](tasks/test.md)
- [tasks/test-backend.md](tasks/test-backend.md)
- [tasks/test-e2e.md](tasks/test-e2e.md)
- [tasks/test-frontend.md](tasks/test-frontend.md)
- [tasks/test-integration.md](tasks/test-integration.md)

### Analysis Summaries

- [📈 Task Usage & Dependencies](summaries/task-usage.md)
- [📝 All Tasks Index](summaries/all-tasks.md)
- [⚡ Command Analysis](summaries/commands.md)
- [📊 Performance Metrics](summaries/performance.md)
- [🔍 Variable Analysis](summaries/variables.md)

## 🔍 Key Insights

### Most Depended-On Tasks

| Task | Used By | Link |
|------|---------|------|
| build-docker | 4 tasks | [View](tasks/build-docker.md) |
| install-frontend | 4 tasks | [View](tasks/install-frontend.md) |
| install-backend | 4 tasks | [View](tasks/install-backend.md) |
| install | 2 tasks | [View](tasks/install.md) |
| test | 2 tasks | [View](tasks/test.md) |

### Top Optimization Opportunities

- ⚠️ **test-integration**: Task could benefit from caching optimization
- ⚠️ **test-e2e**: Task could benefit from caching optimization
- ⚠️ **build-docker**: Task could benefit from caching optimization

### Primary Tool Ecosystem: **Docker**

## 🎯 Next Steps

1. **Review [Optimization Guide](OPTIMIZATION-GUIDE.md)** for specific improvements
2. **Check [Performance Metrics](summaries/performance.md)** for caching opportunities
3. **Examine [Dependency Graph](tasks/dependency-graph.md)** for parallelization potential

