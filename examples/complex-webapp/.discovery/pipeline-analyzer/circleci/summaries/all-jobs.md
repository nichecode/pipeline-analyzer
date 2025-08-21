# All Jobs Index

Total jobs found: **13**

| Job Name | Description | Usage Count | Dependencies |
|----------|-------------|-------------|---------------|
| [build-backend](../jobs/build-backend.md) | *No description* | 1 | test-backend |
| [build-docker-images](../jobs/build-docker-images.md) | *No description* | 1 | build-frontend, build-backend, security-scan |
| [build-frontend](../jobs/build-frontend.md) | *No description* | 1 | test-frontend |
| [deploy-production](../jobs/deploy-production.md) | *No description* | 1 | hold-for-approval |
| [deploy-staging](../jobs/deploy-staging.md) | *No description* | 1 | build-docker-images, test-integration |
| [lint-backend](../jobs/lint-backend.md) | *No description* | 1 | None |
| [lint-frontend](../jobs/lint-frontend.md) | *No description* | 1 | None |
| [security-scan](../jobs/security-scan.md) | *No description* | 3 | None |
| [test-backend](../jobs/test-backend.md) | *No description* | 1 | lint-backend |
| [test-e2e](../jobs/test-e2e.md) | *No description* | 2 | test-integration |
| [test-frontend](../jobs/test-frontend.md) | *No description* | 1 | lint-frontend |
| [test-integration](../jobs/test-integration.md) | *No description* | 1 | test-frontend, test-backend |
| [test-performance](../jobs/test-performance.md) | *No description* | 2 | build-frontend |

## Navigation

- [← Back to Overview](../README.md)
- [Job Usage Analysis](job-usage.md)
