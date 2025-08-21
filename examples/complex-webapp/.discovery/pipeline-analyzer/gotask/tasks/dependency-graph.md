# Task Dependency Graph

**Generated:** 2025-08-21T09:52:29+01:00

## 🔗 Dependency Overview

- **Total tasks:** 23
- **Tasks with dependencies:** 23
- **Circular dependencies:** 0

## 📊 Dependency Visualization

```mermaid
graph TD
    test-integration["test-integration"]
    lint-frontend["lint-frontend"]
    stop-services["stop-services"]
    install-backend["install-backend"]
    build-frontend["build-frontend"]
    deploy-staging["deploy-staging"]
    performance-test["performance-test"]
    dev["dev"]
    install-frontend["install-frontend"]
    test-e2e["test-e2e"]
    lint["lint"]
    build-backend["build-backend"]
    start-services["start-services"]
    deploy-prod["deploy-prod"]
    install["install"]
    test-frontend["test-frontend"]
    lint-backend["lint-backend"]
    build["build"]
    build-docker["build-docker"]
    security-scan["security-scan"]
    clean["clean"]
    test["test"]
    test-backend["test-backend"]
    build-docker --> deploy-staging
    test --> deploy-staging
    install-backend --> build-backend
    build-docker --> deploy-prod
    test --> deploy-prod
    test-e2e --> deploy-prod
    build-frontend --> build
    build-backend --> build
    build --> test-e2e
    start-services --> test-e2e
    build-docker --> start-services
    install-frontend --> test-frontend
    build-docker --> security-scan
    install --> test-integration
    install-frontend --> lint-frontend
    start-services --> performance-test
    install --> dev
    install-frontend --> install
    install-backend --> install
    test-frontend --> test
    test-backend --> test
    test-integration --> test
    install-backend --> test-backend
    install-frontend --> build-frontend
    lint-frontend --> lint
    lint-backend --> lint
    install-backend --> lint-backend
    build --> build-docker
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
- [Optimization Guide](../optimization-guide.md)
