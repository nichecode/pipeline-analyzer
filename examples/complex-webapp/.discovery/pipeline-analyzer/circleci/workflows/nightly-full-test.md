# Workflow: nightly-full-test

## Job Execution Order

| Job | Dependencies | Context |
|-----|--------------|----------|
| [test-e2e](../jobs/test-e2e.md) | None | None |
| [test-performance](../jobs/test-performance.md) | None | None |
| [security-scan](../jobs/security-scan.md) | None | None |

## Dependency Graph

```
Independent jobs (run first):
├── test-e2e
├── test-performance
├── security-scan

```

## Navigation

- [← Back to Workflows](../summaries/workflows.md)
- [← Back to Overview](../README.md)
