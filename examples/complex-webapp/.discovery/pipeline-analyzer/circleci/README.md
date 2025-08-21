# CircleCI Analysis Report

**Generated:** 2025-08-21T10:47:48+01:00
**Config:** /Users/nicholas/Projects/pipeline-analyzer/examples/complex-webapp/.circleci/config.yml

## 📊 Overview

- **Unique jobs:** 13
- **Workflows:** build-test-deploy, nightly-full-test, weekly-security

## 🚀 Quick Start

1. **[📋 Migration Checklist](migration-checklist.md)** - Your step-by-step guide
2. **[📈 Job Usage Analysis](summaries/job-usage.md)** - Job reuse patterns and dependencies
3. **[⚡ Commands Analysis](summaries/commands.md)** - All run commands for conversion

## 📁 Directory Structure

### Jobs
Individual job analysis with run commands and configuration:

- [jobs/build-backend.md](jobs/build-backend.md)
- [jobs/build-docker-images.md](jobs/build-docker-images.md)
- [jobs/build-frontend.md](jobs/build-frontend.md)
- [jobs/deploy-production.md](jobs/deploy-production.md)
- [jobs/deploy-staging.md](jobs/deploy-staging.md)
- [jobs/lint-backend.md](jobs/lint-backend.md)
- [jobs/lint-frontend.md](jobs/lint-frontend.md)
- [jobs/security-scan.md](jobs/security-scan.md)
- [jobs/test-backend.md](jobs/test-backend.md)
- [jobs/test-e2e.md](jobs/test-e2e.md)
- [jobs/test-frontend.md](jobs/test-frontend.md)
- [jobs/test-integration.md](jobs/test-integration.md)
- [jobs/test-performance.md](jobs/test-performance.md)

### Workflows
Workflow structure and job dependencies:

- [workflows/build-test-deploy.md](workflows/build-test-deploy.md)
- [workflows/nightly-full-test.md](workflows/nightly-full-test.md)
- [workflows/weekly-security.md](workflows/weekly-security.md)

### Analysis Summaries

- [📈 Job Usage & Dependencies](summaries/job-usage.md)
- [📝 All Jobs Index](summaries/all-jobs.md)
- [⚡ Commands Analysis](summaries/commands.md)
- [🐳 Docker & Scripts](summaries/docker-and-scripts.md)
- [⚙️ Executors & Images](summaries/executors-and-images.md)
- [🔄 Workflows Index](summaries/workflows.md)

## 🎯 Next Steps

1. **Start with [Migration Checklist](migration-checklist.md)**
2. **Review most frequently used jobs** from [job usage analysis](summaries/job-usage.md)
3. **Examine job dependencies** to understand execution order
4. **Begin converting** highest-impact jobs to go-task format

## 🔍 Most Used Jobs

| Job | Usage Count | Link |
|-----|-------------|------|
| security-scan | 3 | [View Details](jobs/security-scan.md) |
| test-e2e | 2 | [View Details](jobs/test-e2e.md) |
| test-performance | 2 | [View Details](jobs/test-performance.md) |
| lint-backend | 1 | [View Details](jobs/lint-backend.md) |
| test-frontend | 1 | [View Details](jobs/test-frontend.md) |
| test-backend | 1 | [View Details](jobs/test-backend.md) |
| test-integration | 1 | [View Details](jobs/test-integration.md) |
| build-docker-images | 1 | [View Details](jobs/build-docker-images.md) |
| deploy-staging | 1 | [View Details](jobs/deploy-staging.md) |
| deploy-production | 1 | [View Details](jobs/deploy-production.md) |
