# Task: security-scan

**Description:** Run security scans

## 📋 Task Properties

- **Used by:** 0 other tasks
- **Internal:** false
- **Watch mode:** false

## 🔗 Dependencies

This task depends on:

- [build-docker](build-docker.md)

## ⚡ Commands

**Command 1:**
```bash
npm audit --audit-level=moderate
```

- **Category:** package-management
- **Complexity:** 1/5
- **Risk level:** low
- **Tools:** npm

**Command 2:**
```bash
if [ -n "$SNYK_TOKEN" ]; then npx snyk test; fi
```

- **Category:** utility
- **Complexity:** 3/5
- **Risk level:** low
- **Tools:** shell

**Command 3:**
```bash
trivy image complex-webapp-frontend:latest
```

- **Category:** utility
- **Complexity:** 1/5
- **Risk level:** low
- **Tools:** shell

**Command 4:**
```bash
trivy image complex-webapp-backend:latest
```

- **Category:** utility
- **Complexity:** 1/5
- **Risk level:** low
- **Tools:** shell

## 🔍 Command Patterns

- **npm:** 1 occurrence(s)

## 🚀 Optimization Opportunities

- Consider adding sources and generates for caching

## Navigation

- [← Back to All Tasks](../summaries/all-tasks.md)
- [← Back to Overview](../README.md)
- [Dependency Graph](dependency-graph.md)
