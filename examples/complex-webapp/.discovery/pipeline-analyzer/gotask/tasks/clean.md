# Task: clean

**Description:** Clean all build artifacts

## 📋 Task Properties

- **Used by:** 0 other tasks
- **Internal:** false
- **Watch mode:** false

## ⚡ Commands

**Command 1:**
```bash
rm -rf dist/
```

- **Category:** utility
- **Complexity:** 1/5
- **Risk level:** high
- **Tools:** shell

**Command 2:**
```bash
rm -rf node_modules/
```

- **Category:** utility
- **Complexity:** 1/5
- **Risk level:** high
- **Tools:** shell

**Command 3:**
```bash
rm -rf api/dist/
```

- **Category:** utility
- **Complexity:** 1/5
- **Risk level:** high
- **Tools:** shell

**Command 4:**
```bash
docker-compose down --volumes --remove-orphans
```

- **Category:** containerization
- **Complexity:** 1/5
- **Risk level:** low
- **Tools:** docker-compose

**Command 5:**
```bash
docker image prune -f
```

- **Category:** containerization
- **Complexity:** 1/5
- **Risk level:** low
- **Tools:** docker

## 🔍 Command Patterns

- **docker:** 1 occurrence(s)
- **docker-compose:** 1 occurrence(s)

## 🚀 Optimization Opportunities

- Consider adding sources and generates for caching

## Navigation

- [← Back to All Tasks](../summaries/all-tasks.md)
- [← Back to Overview](../README.md)
- [Dependency Graph](dependency-graph.md)
