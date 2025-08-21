# Task: build-docker

**Description:** Build Docker images

## 📋 Task Properties

- **Used by:** 4 other tasks
- **Internal:** false
- **Watch mode:** false

## 🔗 Dependencies

This task depends on:

- [build](build.md)

## ⚡ Commands

### Command 1

```bash
docker build -t complex-webapp-frontend:latest -f docker/frontend/Dockerfile .
```

- **Category:** containerization
- **Complexity:** 1/5
- **Risk level:** low
- **Tools:** docker

**Suggestions:**
- Avoid using 'latest' tag, specify explicit version

### Command 2

```bash
docker build -t complex-webapp-backend:latest -f docker/backend/Dockerfile .
```

- **Category:** containerization
- **Complexity:** 1/5
- **Risk level:** low
- **Tools:** docker

**Suggestions:**
- Avoid using 'latest' tag, specify explicit version

## 🔍 Command Patterns

- **docker:** 1 occurrence(s)

## 🚀 Optimization Opportunities

- Consider adding sources and generates for caching

## Navigation

- [← Back to All Tasks](../summaries/all-tasks.md)
- [← Back to Overview](../README.md)
- [Dependency Graph](dependency-graph.md)
