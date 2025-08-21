# Task: build

**Description:** Build all components

## 📋 Task Properties

- **Used by:** 2 other tasks
- **Internal:** false
- **Watch mode:** false

## 🔗 Dependencies

This task depends on:

- [build-frontend](build-frontend.md)
- [build-backend](build-backend.md)

## ⚡ Commands

This is a **dependency-only task** that orchestrates other tasks without running direct commands.

**Execution Flow:**
1. Runs all dependency tasks in the correct order
2. Completes when all dependencies finish successfully

## Navigation

- [← Back to All Tasks](../summaries/all-tasks.md)
- [← Back to Overview](../README.md)
- [Dependency Graph](dependency-graph.md)
