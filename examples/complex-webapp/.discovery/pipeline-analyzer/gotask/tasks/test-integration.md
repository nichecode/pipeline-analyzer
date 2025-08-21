# Task: test-integration

**Description:** Run integration tests

## 📋 Task Properties

- **Used by:** 1 other tasks
- **Internal:** false
- **Watch mode:** false

## 🔗 Dependencies

This task depends on:

- [install](install.md)

## ⚡ Commands

### Command 1

```bash
docker-compose up -d postgres redis
```

- **Category:** containerization
- **Complexity:** 1/5
- **Risk level:** low
- **Tools:** docker-compose

### Command 2

```bash
sleep 10
```

- **Category:** utility
- **Complexity:** 1/5
- **Risk level:** low
- **Tools:** shell

### Command 3

```bash
npm run test:integration
```

- **Category:** package-management
- **Complexity:** 1/5
- **Risk level:** low
- **Tools:** npm

### Command 4

```bash
docker-compose down
```

- **Category:** containerization
- **Complexity:** 1/5
- **Risk level:** low
- **Tools:** docker-compose

## 🔍 Command Patterns

- **docker-compose:** 1 occurrence(s)
- **npm:** 1 occurrence(s)

## 🚀 Optimization Opportunities

- Consider adding sources and generates for caching

## Navigation

- [← Back to All Tasks](../summaries/all-tasks.md)
- [← Back to Overview](../README.md)
- [Dependency Graph](dependency-graph.md)
