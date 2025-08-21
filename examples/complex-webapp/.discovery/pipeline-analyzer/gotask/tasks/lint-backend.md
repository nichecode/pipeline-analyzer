# Task: lint-backend

**Description:** Lint backend code

## 📋 Task Properties

- **Used by:** 1 other tasks
- **Internal:** false
- **Watch mode:** false

## 🔗 Dependencies

This task depends on:

- [install-backend](install-backend.md)

## ⚡ Commands

### Command 1

```bash
cd api && flake8 .
```

- **Category:** utility
- **Complexity:** 3/5
- **Risk level:** low
- **Tools:** shell

### Command 2

```bash
cd api && black --check .
```

- **Category:** utility
- **Complexity:** 3/5
- **Risk level:** low
- **Tools:** shell

### Command 3

```bash
cd api && mypy .
```

- **Category:** utility
- **Complexity:** 3/5
- **Risk level:** low
- **Tools:** shell

## 🚀 Optimization Opportunities

- Consider adding sources and generates for caching

## Navigation

- [← Back to All Tasks](../summaries/all-tasks.md)
- [← Back to Overview](../README.md)
- [Dependency Graph](dependency-graph.md)
