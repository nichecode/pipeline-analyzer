# Task: deploy-staging

**Description:** Deploy to staging environment

## 📋 Task Properties

- **Used by:** 0 other tasks
- **Internal:** false
- **Watch mode:** false

## 🔗 Dependencies

This task depends on:

- [build-docker](build-docker.md)
- [test](test.md)

## ⚡ Commands

```bash
kubectl apply -f k8s/staging/ --context=staging-cluster
```

- **Category:** deployment
- **Complexity:** 1/5
- **Risk level:** low
- **Tools:** kubectl

## 🔍 Command Patterns

- **kubectl:** 1 occurrence(s)

## 🚀 Optimization Opportunities

- Consider adding sources and generates for caching

## Navigation

- [← Back to All Tasks](../summaries/all-tasks.md)
- [← Back to Overview](../README.md)
- [Dependency Graph](dependency-graph.md)
