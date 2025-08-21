# Commands Analysis for go-task Migration

## ⚡ Command Patterns Suitable for go-task

### pip (2 commands)

- `pip install -r api/requirements-dev.txt`
- `pip install -r api/requirements.txt`

### shell (3 commands)

- `aquasec/trivy image webapp-frontend:${{ github.sha }}`
- `echo "Deploying to staging..."`
- `cd api`

### python (1 commands)

- `python -m pytest tests/ --cov=. --cov-report=xml`

### npm (6 commands)

- `npm run lint`
- `npm run test:coverage`
- `npm ci`
- `npm run build:prod`
- `npm audit --audit-level=moderate`

### docker (3 commands)

- `docker build -t webapp-frontend:${{ github.sha }} -f docker/frontend/Dockerfile .`
- `docker build -t webapp-backend:${{ github.sha }} -f docker/backend/Dockerfile .`
- `docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \`


## 🔍 Navigation

- [← Back to GitHub Actions Overview](../README.md)
- [go-task Migration Guide](go-task-migration.md)
