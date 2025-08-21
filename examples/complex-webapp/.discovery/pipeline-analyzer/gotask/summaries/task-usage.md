# Task Usage Analysis

## Most Frequently Used Tasks

| Rank | Task Name | Used By | Link |
|------|-----------|---------|------|
| 1 | install-backend | 4 tasks | [install-backend](../tasks/install-backend.md) |
| 2 | build-docker | 4 tasks | [build-docker](../tasks/build-docker.md) |
| 3 | install-frontend | 4 tasks | [install-frontend](../tasks/install-frontend.md) |
| 4 | test | 2 tasks | [test](../tasks/test.md) |
| 5 | build | 2 tasks | [build](../tasks/build.md) |
| 6 | start-services | 2 tasks | [start-services](../tasks/start-services.md) |
| 7 | install | 2 tasks | [install](../tasks/install.md) |
| 8 | test-backend | 1 tasks | [test-backend](../tasks/test-backend.md) |
| 9 | lint-backend | 1 tasks | [lint-backend](../tasks/lint-backend.md) |
| 10 | test-e2e | 1 tasks | [test-e2e](../tasks/test-e2e.md) |
| 11 | test-frontend | 1 tasks | [test-frontend](../tasks/test-frontend.md) |
| 12 | test-integration | 1 tasks | [test-integration](../tasks/test-integration.md) |
| 13 | build-frontend | 1 tasks | [build-frontend](../tasks/build-frontend.md) |
| 14 | build-backend | 1 tasks | [build-backend](../tasks/build-backend.md) |
| 15 | lint-frontend | 1 tasks | [lint-frontend](../tasks/lint-frontend.md) |

## Task Dependencies

Tasks with dependencies:

| Task | Dependencies |
|------|-------------|
| [build](../tasks/build.md) | [build-frontend](../tasks/build-frontend.md), [build-backend](../tasks/build-backend.md) |
| [build-backend](../tasks/build-backend.md) | [install-backend](../tasks/install-backend.md) |
| [build-docker](../tasks/build-docker.md) | [build](../tasks/build.md) |
| [build-frontend](../tasks/build-frontend.md) | [install-frontend](../tasks/install-frontend.md) |
| [deploy-prod](../tasks/deploy-prod.md) | [build-docker](../tasks/build-docker.md), [test](../tasks/test.md), [test-e2e](../tasks/test-e2e.md) |
| [deploy-staging](../tasks/deploy-staging.md) | [build-docker](../tasks/build-docker.md), [test](../tasks/test.md) |
| [dev](../tasks/dev.md) | [install](../tasks/install.md) |
| [install](../tasks/install.md) | [install-frontend](../tasks/install-frontend.md), [install-backend](../tasks/install-backend.md) |
| [lint](../tasks/lint.md) | [lint-frontend](../tasks/lint-frontend.md), [lint-backend](../tasks/lint-backend.md) |
| [lint-backend](../tasks/lint-backend.md) | [install-backend](../tasks/install-backend.md) |
| [lint-frontend](../tasks/lint-frontend.md) | [install-frontend](../tasks/install-frontend.md) |
| [performance-test](../tasks/performance-test.md) | [start-services](../tasks/start-services.md) |
| [security-scan](../tasks/security-scan.md) | [build-docker](../tasks/build-docker.md) |
| [start-services](../tasks/start-services.md) | [build-docker](../tasks/build-docker.md) |
| [test](../tasks/test.md) | [test-frontend](../tasks/test-frontend.md), [test-backend](../tasks/test-backend.md), [test-integration](../tasks/test-integration.md) |
| [test-backend](../tasks/test-backend.md) | [install-backend](../tasks/install-backend.md) |
| [test-e2e](../tasks/test-e2e.md) | [build](../tasks/build.md), [start-services](../tasks/start-services.md) |
| [test-frontend](../tasks/test-frontend.md) | [install-frontend](../tasks/install-frontend.md) |
| [test-integration](../tasks/test-integration.md) | [install](../tasks/install.md) |

## 🚀 Independent Tasks

These tasks have no dependencies and can run in parallel:

- [clean](../tasks/clean.md)
- [install-backend](../tasks/install-backend.md)
- [install-frontend](../tasks/install-frontend.md)
- [stop-services](../tasks/stop-services.md)

## Navigation

- [← Back to Overview](../README.md)
- [All Tasks Index](all-tasks.md)
- [Dependency Graph](../tasks/dependency-graph.md)
