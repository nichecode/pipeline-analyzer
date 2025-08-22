# CircleCI Analysis Report

**Generated:** 2025-08-22T11:41:24+01:00
**Config:** /Users/nicholas/Projects/pipeline-analyzer/examples/complex-webapp/.circleci/config.yml

## 📊 Overview

- **Unique jobs:** 13
- **Reusable commands:** 4
- **Workflows:** build-test-deploy, nightly-full-test, weekly-security

## 📊 Workflow Overview

```mermaid
flowchart TD
    SECURITY_SCAN["`**security-scan**
Executor: node-executor
• install-dependencies
• npm audit --audit-level=moderate --js...
• if [ -n "$SNYK_TOKEN" ]; then
  npx s...
• ... (1 more)
`"]
    TEST_PERFORMANCE["`**test-performance**
Executor: node-executor
• install-dependencies
• npm run start:prod &
sleep 30
npx lhc...
• npm run test:load -- --reporter json ...
• ... (1 more)
`"]
    TEST_E2E["`**test-e2e**
Executor: e2e-executor
• install-dependencies
• dockerize -wait tcp://localhost:5432 ...
• cd api
python manage.py migrate --set...
• ... (2 more)
`"]
    TEST_INTEGRATION["`**test-integration**
Executor: e2e-executor
• install-dependencies
• dockerize -wait tcp://localhost:5432 ...
• cd api
python manage.py migrate --set...
• ... (2 more)
`"]
    BUILD_FRONTEND["`**build-frontend**
Executor: node-executor
• install-dependencies
• npm run build:prod
npm run build:anal...
• tar -czf build-artifacts.tar.gz dist/...
• ... (1 more)
`"]
    BUILD_BACKEND["`**build-backend**
Executor: python-executor
• cd api
python setup.py sdist bdist_wh...
• cd api
pip install twine
twine check ...
• notify-slack-on-failure
`"]
    BUILD_DOCKER_IMAGES["`**build-docker-images**
Executor: docker-executor
• # Build frontend image
docker build \...
• # Install trivy
curl -sfL https://raw...
• if [ -n "$DOCKER_HUB_TOKEN" ] && [ "$...
• ... (1 more)
`"]
    DEPLOY_STAGING["`**deploy-staging**
Executor: docker-executor
• aws-cli/setup
• kubernetes/install-kubectl
• # Update Kubernetes manifests with ne...
• ... (2 more)
`"]
    TEST_INTEGRATION --> TEST_E2E
    BUILD_FRONTEND --> BUILD_DOCKER_IMAGES
    BUILD_BACKEND --> BUILD_DOCKER_IMAGES
    SECURITY_SCAN --> BUILD_DOCKER_IMAGES
    BUILD_FRONTEND --> TEST_PERFORMANCE
    BUILD_DOCKER_IMAGES --> DEPLOY_STAGING
    TEST_INTEGRATION --> DEPLOY_STAGING

    classDef workflow fill:#e1f5fe,stroke:#01579b,stroke-width:3px
    classDef setup fill:#f3e5f5,stroke:#4a148c,stroke-width:2px
    classDef test fill:#e8f5e8,stroke:#1b5e20,stroke-width:2px
    classDef build fill:#fff3e0,stroke:#e65100,stroke-width:2px
    classDef deploy fill:#e0f2f1,stroke:#004d40,stroke-width:2px
    classDef utility fill:#f1f8e9,stroke:#33691e,stroke-width:2px
    class SECURITY_SCAN test
    class TEST_PERFORMANCE test
    class TEST_E2E test
    class TEST_INTEGRATION test
    class BUILD_FRONTEND build
    class BUILD_BACKEND build
    class BUILD_DOCKER_IMAGES build
    class DEPLOY_STAGING deploy
```

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

### Reusable Commands
Reusable command definitions with analysis and usage patterns:

- [commands/install-dependencies.md](commands/install-dependencies.md)
- [commands/notify-slack-on-failure.md](commands/notify-slack-on-failure.md)
- [commands/restore-npm-cache.md](commands/restore-npm-cache.md)
- [commands/save-npm-cache.md](commands/save-npm-cache.md)

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
| test-integration | 1 | [View Details](jobs/test-integration.md) |
| deploy-production | 1 | [View Details](jobs/deploy-production.md) |
| test-backend | 1 | [View Details](jobs/test-backend.md) |
| lint-backend | 1 | [View Details](jobs/lint-backend.md) |
| build-frontend | 1 | [View Details](jobs/build-frontend.md) |
| build-backend | 1 | [View Details](jobs/build-backend.md) |
| build-docker-images | 1 | [View Details](jobs/build-docker-images.md) |
