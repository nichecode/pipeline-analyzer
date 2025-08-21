# Task: test

**Description:** Run all tests

## 📋 Task Properties

- **Used by:** 2 other tasks
- **Internal:** false
- **Watch mode:** false

## 🔗 Dependencies

This task depends on:

- [test-frontend](test-frontend.md)
- [test-backend](test-backend.md)
- [test-integration](test-integration.md)

## ⚡ Commands

This is a **dependency-only task** that orchestrates other tasks without running direct commands.

**Execution Flow:**
1. Runs all dependency tasks in the correct order
2. Completes when all dependencies finish successfully

## 🚀 Optimization Opportunities

- Consider adding sources and generates for caching

## Navigation

- [← Back to All Tasks](../summaries/all-tasks.md)
- [← Back to Overview](../README.md)
- [Dependency Graph](dependency-graph.md)
