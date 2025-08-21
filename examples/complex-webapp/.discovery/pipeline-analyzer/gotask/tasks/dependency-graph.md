# Task Dependency Graph

**Generated:** 2025-08-21T07:59:05+01:00

## 🔗 Dependency Overview

- **Total tasks:** 23
- **Tasks with dependencies:** 23
- **Circular dependencies:** 0

## 📊 Dependency Visualization

```mermaid
graph TD
    security-scan["security-scan"]
    dev["dev"]
    install["install"]
    test-frontend["test-frontend"]
    lint["lint"]
    build["build"]
    start-services["start-services"]
    deploy-prod["deploy-prod"]
    performance-test["performance-test"]
    install-frontend["install-frontend"]
    test-backend["test-backend"]
    build-frontend["build-frontend"]
    clean["clean"]
    install-backend["install-backend"]
    test-integration["test-integration"]
    test-e2e["test-e2e"]
    build-backend["build-backend"]
    build-docker["build-docker"]
    stop-services["stop-services"]
    test["test"]
    lint-frontend["lint-frontend"]
    lint-backend["lint-backend"]
    deploy-staging["deploy-staging"]
    install-frontend --> lint-frontend
    build-docker --> deploy-staging
    test --> deploy-staging
    install --> dev
    install-frontend --> test-frontend
    lint-frontend --> lint
    lint-backend --> lint
    install-backend --> test-backend
    build --> build-docker
    build-frontend --> build
    build-backend --> build
    install-frontend --> build-frontend
    install --> test-integration
    install-backend --> lint-backend
    build-docker --> security-scan
    start-services --> performance-test
    build --> test-e2e
    start-services --> test-e2e
    install-backend --> build-backend
    test-frontend --> test
    test-backend --> test
    test-integration --> test
    install-frontend --> install
    install-backend --> install
    build-docker --> start-services
    build-docker --> deploy-prod
    test --> deploy-prod
    test-e2e --> deploy-prod
```

## 🎯 Critical Path

The longest dependency chain in your tasks:

```
deploy-prod → test-e2e → start-services → build-docker → build → build-backend → install-backend
```

**Length:** 7 tasks

Optimizing tasks in the critical path will have the biggest impact on overall execution time.

## 📊 Task Levels

Tasks grouped by their dependency depth:

**Level 1** (can run in parallel):
- [clean](clean.md)
- [install-backend](install-backend.md)
- [install-frontend](install-frontend.md)
- [stop-services](stop-services.md)

**Level 2** (can run in parallel):
- [build-backend](build-backend.md)
- [build-frontend](build-frontend.md)
- [install](install.md)
- [lint-backend](lint-backend.md)
- [lint-frontend](lint-frontend.md)
- [test-backend](test-backend.md)
- [test-frontend](test-frontend.md)

**Level 3** (can run in parallel):
- [build](build.md)
- [dev](dev.md)
- [lint](lint.md)
- [test-integration](test-integration.md)

**Level 4** (can run in parallel):
- [build-docker](build-docker.md)
- [test](test.md)

**Level 5** (can run in parallel):
- [deploy-staging](deploy-staging.md)
- [security-scan](security-scan.md)
- [start-services](start-services.md)

**Level 6** (can run in parallel):
- [performance-test](performance-test.md)
- [test-e2e](test-e2e.md)

**Level 7** (can run in parallel):
- [deploy-prod](deploy-prod.md)

## Navigation

- [← Back to Overview](../README.md)
- [Task Usage Analysis](../summaries/task-usage.md)
- [Optimization Guide](../OPTIMIZATION-GUIDE.md)
