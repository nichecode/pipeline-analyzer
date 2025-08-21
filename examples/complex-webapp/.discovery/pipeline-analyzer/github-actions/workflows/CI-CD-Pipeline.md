# Workflow: CI/CD Pipeline

**File:** ci.yml  
**Generated:** 2025-08-21T10:04:30+01:00

## 📊 Overview

- **Jobs:** 3
- **Total steps:** 16
- **Actions used:** 4 unique actions
- **Runners used:** 1 unique runners

## 🔄 go-task Migration Opportunities

- 💡 Consider creating go-task equivalents for repeated command patterns
- 📦 Multiple npm commands detected - consider consolidating into go-task
- 🐳 Docker commands found - go-task could simplify container management

## 📋 Jobs Overview

| Job | Runner | Steps | Commands | Actions |
|-----|--------|-------|----------|---------|
| [test-frontend](../jobs/test-frontend.md) | ubuntu-latest | 6 | 3 | 3 |
| [test-backend](../jobs/test-backend.md) | ubuntu-latest | 4 | 4 | 2 |
| [build-and-deploy](../jobs/build-and-deploy.md) | ubuntu-latest | 6 | 8 | 2 |

## ⚡ Command Patterns Detected

These command patterns are good candidates for go-task consolidation:

### npm (6 commands)

- `npm ci`
- `npm run lint`
- `npm run test:coverage`
- `npm ci`
- `npm run build:prod`
- `npm audit --audit-level=moderate`

### pip (2 commands)

- `pip install -r api/requirements.txt`
- `pip install -r api/requirements-dev.txt`

### shell (3 commands)

- `cd api`
- `aquasec/trivy image webapp-frontend:${{ github.sha }}`
- `echo "Deploying to staging..."`

### python (1 commands)

- `python -m pytest tests/ --cov=. --cov-report=xml`

### docker (3 commands)

- `docker build -t webapp-frontend:${{ github.sha }} -f docker/frontend/Dockerfile .`
- `docker build -t webapp-backend:${{ github.sha }} -f docker/backend/Dockerfile .`
- `docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \`


## 🔍 Navigation

- [← Back to GitHub Actions Overview](../README.md)
- [All Jobs](../summaries/actions-usage.md)
- [go-task Migration Guide](../summaries/go-task-migration.md)
