# Task: build-backend

**Description:** Build backend package

## 📋 Task Properties

- **Used by:** 1 other tasks
- **Internal:** false
- **Watch mode:** false

## 🔗 Dependencies

This task depends on:

- [install-backend](install-backend.md)

## 📁 Files

**Source files:**
- `api/**/*.py`
- `api/setup.py`

**Generated files:**
- `api/dist/**`

## ⚡ Commands

```bash
cd api && python setup.py sdist bdist_wheel
```

- **Category:** runtime
- **Complexity:** 3/5
- **Risk level:** low
- **Tools:** python

## 🔍 Command Patterns

- **python:** 1 occurrence(s)

## Navigation

- [← Back to All Tasks](../summaries/all-tasks.md)
- [← Back to Overview](../README.md)
- [Dependency Graph](dependency-graph.md)
