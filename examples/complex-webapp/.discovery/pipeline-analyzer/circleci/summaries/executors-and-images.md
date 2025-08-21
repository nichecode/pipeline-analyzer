# Executors & Images Analysis

## Image/Executor Usage

| Image/Executor | Jobs Using |
|----------------|------------|
| `cimg/postgres:14.0` | 1 jobs |
| `cimg/python:3.11` | 1 jobs |
| `docker-executor (cimg/base:stable)` | 3 jobs |
| `e2e-executor (cimg/node:18.17.0-browsers)` | 2 jobs |
| `e2e-executor (cimg/postgres:14.0)` | 2 jobs |
| `e2e-executor (redis:6.2-alpine)` | 2 jobs |
| `node-executor (cimg/node:18.17.0)` | 5 jobs |
| `python-executor (cimg/python:3.11)` | 3 jobs |

## Detailed Job Assignments

### cimg/postgres:14.0

- [test-backend](../jobs/test-backend.md)

### cimg/python:3.11

- [test-backend](../jobs/test-backend.md)

### docker-executor (cimg/base:stable)

- [build-docker-images](../jobs/build-docker-images.md)
- [deploy-production](../jobs/deploy-production.md)
- [deploy-staging](../jobs/deploy-staging.md)

### e2e-executor (cimg/node:18.17.0-browsers)

- [test-e2e](../jobs/test-e2e.md)
- [test-integration](../jobs/test-integration.md)

### e2e-executor (cimg/postgres:14.0)

- [test-e2e](../jobs/test-e2e.md)
- [test-integration](../jobs/test-integration.md)

### e2e-executor (redis:6.2-alpine)

- [test-e2e](../jobs/test-e2e.md)
- [test-integration](../jobs/test-integration.md)

### node-executor (cimg/node:18.17.0)

- [build-frontend](../jobs/build-frontend.md)
- [lint-frontend](../jobs/lint-frontend.md)
- [security-scan](../jobs/security-scan.md)
- [test-frontend](../jobs/test-frontend.md)
- [test-performance](../jobs/test-performance.md)

### python-executor (cimg/python:3.11)

- [build-backend](../jobs/build-backend.md)
- [lint-backend](../jobs/lint-backend.md)
- [test-backend](../jobs/test-backend.md)

## Navigation

- [← Back to Overview](../README.md)
- [Docker & Scripts](docker-and-scripts.md)
