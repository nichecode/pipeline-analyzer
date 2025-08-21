# Job: test-frontend

**Workflow:** CI/CD Pipeline  
**Runner:** ubuntu-latest  
**Estimated Duration:** ~4 min  
**Caching Enabled:** false

## 📊 Overview

- **Steps:** 6
- **Run commands:** 3
- **Actions used:** 3
- **Dependencies:** 

## ⚡ Commands (go-task candidates)

These commands could be extracted into go-task:

1. `npm ci`
2. `npm run lint`
3. `npm run test:coverage`

### Suggested go-task conversion:

```yaml
version: '3'
tasks:
  test-frontend:
    desc: "test-frontend task"
    cmds:
      - npm ci
      - npm run lint
      - npm run test:coverage
```

## 🛠️ GitHub Actions Used

- `actions/checkout@v4`
- `actions/setup-node@v4`
- `codecov/codecov-action@v4`

## 🔍 Navigation

- [← Back to GitHub Actions Overview](../README.md)
- [go-task Migration Guide](../summaries/go-task-migration.md)
