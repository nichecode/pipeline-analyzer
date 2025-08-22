# go-task Migration Guide for GitHub Actions

## 🎯 Migration Strategy

This guide helps you refactor your GitHub Actions workflows to use **go-task as the command runner**, making your builds:

- **Locally executable** - Run the same commands locally as in CI
- **CI-agnostic** - Easy to switch between GitHub Actions, CircleCI, etc.
- **Testable** - Debug build issues without pushing to CI

## 🔄 Refactoring Pattern

**Current Pattern:**
```yaml
- name: Build and test
  run: |
    npm install
    npm run build
    npm run test
```

**Target Pattern:**
```yaml
- name: Build and test
  run: task build-and-test
```

With corresponding Taskfile.yml:
```yaml
version: '3'
tasks:
  build-and-test:
    desc: "Build and test the application"
    cmds:
      - npm install
      - npm run build  
      - npm run test
```

## 📋 Migration Checklist

### 1. Analyze Current Commands
Found **14 unique commands** across your workflows:

- `npm run build:prod`
- `docker build -t webapp-frontend:${{ github.sha }} -f docker/frontend/Dockerfile .`
- `docker build -t webapp-backend:${{ github.sha }} -f docker/backend/Dockerfile .`
- `npm ci`
- `pip install -r api/requirements-dev.txt`
- `python -m pytest tests/ --cov=. --cov-report=xml`
- `npm audit --audit-level=moderate`
- `docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \`
- `aquasec/trivy image webapp-frontend:${{ github.sha }}`
- `echo "Deploying to staging..."`
- ... and 4 more

### 2. Create Taskfile.yml

Create a Taskfile.yml in your repository root:

```yaml
version: '3'

vars:
  # Define common variables here

tasks:
  # Extract your commands into logical tasks
  install:
    desc: "Install dependencies"
    cmds:
      - npm install

  build:
    desc: "Build the application"  
    deps: [install]
    cmds:
      - npm run build

  test:
    desc: "Run tests"
    deps: [build]  
    cmds:
      - npm run test

  deploy:
    desc: "Deploy application"
    deps: [test]
    cmds:
      - # Add deployment commands here
```

### 3. Update GitHub Actions Workflows

Replace complex run commands with task calls:

**Before:**
```yaml
- name: Install dependencies
  run: npm install
  
- name: Build  
  run: npm run build
  
- name: Test
  run: npm run test
```

**After:**
```yaml
- name: Install go-task
  run: go install github.com/go-task/task/v3/cmd/task@latest
  
- name: Run build pipeline
  run: task test  # This runs install -> build -> test
```

### 4. Test Locally

Before pushing changes:

```bash
# Install go-task locally
go install github.com/go-task/task/v3/cmd/task@latest

# Test your tasks
task install
task build  
task test
```

## 🚀 Benefits After Migration

- ✅ **Local Development** - Run CI commands locally
- ✅ **Faster Debugging** - No need to push to test builds  
- ✅ **CI Portability** - Easy to switch CI providers
- ✅ **Consistent Environments** - Same commands everywhere
- ✅ **Better Developer Experience** - Clear, documented tasks

## 💡 Pro Tips

1. **Group Related Commands** - Combine related steps into single tasks
2. **Use Dependencies** - Let go-task handle execution order with deps
3. **Add Descriptions** - Use desc for clear task documentation  
4. **Define Variables** - Use vars section for reusable values
5. **Local vs CI** - Use platforms to handle OS differences

## 🔍 Next Steps

1. **Start Small** - Pick one workflow to convert first
2. **Create Tasks** - Extract commands into logical groups  
3. **Update Workflow** - Replace run commands with task calls
4. **Test Locally** - Verify tasks work before committing
5. **Iterate** - Apply pattern to remaining workflows

## 🔍 Navigation

- [← Back to GitHub Actions Overview](../README.md)
- [Commands Analysis](commands-analysis.md)
- [Actions Usage](actions-usage.md)
