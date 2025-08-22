# Commands Analysis for go-task Migration

## ⚡ Command Patterns Suitable for go-task

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

- `docker build -t webapp-backend:${{ github.sha }} -f docker/backend/Dockerfile .`
- `docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \`
- `docker build -t webapp-frontend:${{ github.sha }} -f docker/frontend/Dockerfile .`

### npm (6 commands)

- `npm ci`
- `npm run lint`
- `npm run test:coverage`
- `npm run build:prod`
- `npm audit --audit-level=moderate`


## 🔍 Navigation

- [← Back to GitHub Actions Overview](../README.md)
- [go-task Migration Guide](go-task-migration.md)
