# Job Usage Analysis

## Most Frequently Used Jobs

| Rank | Job Name | Usage Count | Link |
|------|----------|-------------|------|
| 1 | security-scan | 3 | [security-scan](../jobs/security-scan.md) |
| 2 | test-performance | 2 | [test-performance](../jobs/test-performance.md) |
| 3 | test-e2e | 2 | [test-e2e](../jobs/test-e2e.md) |
| 4 | build-docker-images | 1 | [build-docker-images](../jobs/build-docker-images.md) |
| 5 | deploy-staging | 1 | [deploy-staging](../jobs/deploy-staging.md) |
| 6 | test-frontend | 1 | [test-frontend](../jobs/test-frontend.md) |
| 7 | test-backend | 1 | [test-backend](../jobs/test-backend.md) |
| 8 | build-backend | 1 | [build-backend](../jobs/build-backend.md) |
| 9 | build-frontend | 1 | [build-frontend](../jobs/build-frontend.md) |
| 10 | test-integration | 1 | [test-integration](../jobs/test-integration.md) |
| 11 | hold-for-approval | 1 | [hold-for-approval](../jobs/hold-for-approval.md) |
| 12 | deploy-production | 1 | [deploy-production](../jobs/deploy-production.md) |
| 13 | lint-frontend | 1 | [lint-frontend](../jobs/lint-frontend.md) |
| 14 | lint-backend | 1 | [lint-backend](../jobs/lint-backend.md) |

## Job Dependencies

Jobs with dependencies:

| Job | Dependencies |
|-----|-------------|
| [build-backend](../jobs/build-backend.md) | test-backend |
| [build-docker-images](../jobs/build-docker-images.md) | build-frontend, build-backend, security-scan |
| [build-frontend](../jobs/build-frontend.md) | test-frontend |
| [deploy-production](../jobs/deploy-production.md) | hold-for-approval |
| [deploy-staging](../jobs/deploy-staging.md) | build-docker-images, test-integration |
| [hold-for-approval](../jobs/hold-for-approval.md) | build-docker-images, test-performance |
| [test-backend](../jobs/test-backend.md) | lint-backend |
| [test-e2e](../jobs/test-e2e.md) | test-integration |
| [test-frontend](../jobs/test-frontend.md) | lint-frontend |
| [test-integration](../jobs/test-integration.md) | test-frontend, test-backend |
| [test-performance](../jobs/test-performance.md) | build-frontend |

## Unused Jobs

All jobs are used in workflows.

## Navigation

- [← Back to Overview](../README.md)
- [All Jobs Index](all-jobs.md)
