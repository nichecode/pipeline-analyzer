# Job: test-backend

**Workflow:** CI/CD Pipeline  
**Runner:** ubuntu-latest  
**Estimated Duration:** ~2 min  
**Caching Enabled:** false

## 📊 Overview

- **Steps:** 4
- **Run commands:** 4
- **Actions used:** 2
- **Dependencies:** 

## ⚡ Commands (go-task candidates)

These commands could be extracted into go-task:

1. `pip install -r api/requirements.txt`
2. `pip install -r api/requirements-dev.txt`
3. `cd api`
4. `python -m pytest tests/ --cov=. --cov-report=xml`

### Suggested go-task conversion:

```yaml
version: '3'
tasks:
  test-backend:
    desc: "test-backend task"
    cmds:
      - pip install -r api/requirements.txt
      - pip install -r api/requirements-dev.txt
      - cd api
      - python -m pytest tests/ --cov=. --cov-report=xml
```

## 🛠️ GitHub Actions Used

- `actions/checkout@v4`
- `actions/setup-python@v4`

## 💡 Optimization Recommendations

- ⚡ Consider adding caching to improve build times

## 🔍 Navigation

- [← Back to GitHub Actions Overview](../README.md)
- [go-task Migration Guide](../summaries/go-task-migration.md)
