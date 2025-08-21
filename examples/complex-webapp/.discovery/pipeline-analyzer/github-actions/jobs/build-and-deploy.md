# Job: build-and-deploy

**Workflow:** CI/CD Pipeline  
**Runner:** ubuntu-latest  
**Estimated Duration:** > 5 min  
**Caching Enabled:** false

## 📊 Overview

- **Steps:** 6
- **Run commands:** 8
- **Actions used:** 2
- **Dependencies:** test-frontend, test-backend

## ⚡ Commands (go-task candidates)

These commands could be extracted into go-task:

1. `npm ci`
2. `npm run build:prod`
3. `docker build -t webapp-frontend:${{ github.sha }} -f docker/frontend/Dockerfile .`
4. `docker build -t webapp-backend:${{ github.sha }} -f docker/backend/Dockerfile .`
5. `npm audit --audit-level=moderate`
6. `docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \`
7. `aquasec/trivy image webapp-frontend:${{ github.sha }}`
8. `echo "Deploying to staging..."`

### Suggested go-task conversion:

```yaml
version: '3'
tasks:
  build-and-deploy:
    desc: "build-and-deploy task"
    cmds:
      - npm ci
      - npm run build:prod
      - docker build -t webapp-frontend:${{ github.sha }} -f docker/frontend/Dockerfile .
      - docker build -t webapp-backend:${{ github.sha }} -f docker/backend/Dockerfile .
      - npm audit --audit-level=moderate
      - docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \
      - aquasec/trivy image webapp-frontend:${{ github.sha }}
      - echo "Deploying to staging..."
```

## 🛠️ GitHub Actions Used

- `actions/checkout@v4`
- `actions/setup-node@v4`

## 💡 Optimization Recommendations

- ⚡ Consider adding caching to improve build times

## 🔍 Navigation

- [← Back to GitHub Actions Overview](../README.md)
- [go-task Migration Guide](../summaries/go-task-migration.md)
