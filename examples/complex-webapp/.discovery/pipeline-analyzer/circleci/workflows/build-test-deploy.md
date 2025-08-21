# Workflow: build-test-deploy

## Job Execution Order

| Job | Dependencies | Context |
|-----|--------------|----------|
| [lint-frontend](../jobs/lint-frontend.md) | None | None |
| [lint-backend](../jobs/lint-backend.md) | None | None |
| [security-scan](../jobs/security-scan.md) | None | None |
| [test-frontend](../jobs/test-frontend.md) | lint-frontend | None |
| [test-backend](../jobs/test-backend.md) | lint-backend | None |
| [test-integration](../jobs/test-integration.md) | test-frontend, test-backend | None |
| [test-e2e](../jobs/test-e2e.md) | test-integration | None |
| [build-frontend](../jobs/build-frontend.md) | test-frontend | None |
| [build-backend](../jobs/build-backend.md) | test-backend | None |
| [build-docker-images](../jobs/build-docker-images.md) | build-frontend, build-backend, security-scan | None |
| [test-performance](../jobs/test-performance.md) | build-frontend | None |
| [deploy-staging](../jobs/deploy-staging.md) | build-docker-images, test-integration | None |
| [hold-for-approval](../jobs/hold-for-approval.md) | build-docker-images, test-performance | None |
| [deploy-production](../jobs/deploy-production.md) | hold-for-approval | None |

## Dependency Graph

```
Independent jobs (run first):
├── lint-frontend
├── lint-backend
├── security-scan

Dependent jobs:
├── test-e2e
│   └── requires: test-integration
├── build-backend
│   └── requires: test-backend
├── test-performance
│   └── requires: build-frontend
├── deploy-production
│   └── requires: hold-for-approval
├── test-backend
│   └── requires: lint-backend
├── build-frontend
│   └── requires: test-frontend
├── build-docker-images
│   └── requires: build-frontend
│   └── requires: build-backend
│   └── requires: security-scan
├── deploy-staging
│   └── requires: build-docker-images
│   └── requires: test-integration
├── hold-for-approval
│   └── requires: build-docker-images
│   └── requires: test-performance
├── test-frontend
│   └── requires: lint-frontend
├── test-integration
│   └── requires: test-frontend
│   └── requires: test-backend
```

## Navigation

- [← Back to Workflows](../summaries/workflows.md)
- [← Back to Overview](../README.md)
