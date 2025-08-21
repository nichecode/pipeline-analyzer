# Task: dev

**Description:** Start development environment

## 📋 Task Properties

- **Used by:** 0 other tasks
- **Internal:** false
- **Watch mode:** false

## 🔗 Dependencies

This task depends on:

- [install](install.md)

## ⚡ Commands

### Command 1

```bash
docker-compose up -d postgres redis
```

- **Category:** containerization
- **Complexity:** 1/5
- **Risk level:** low
- **Tools:** docker-compose

### Command 2

```bash
npm run start &
```

- **Category:** package-management
- **Complexity:** 2/5
- **Risk level:** low
- **Tools:** npm

### Command 3

```bash
cd api && python manage.py runserver &
```

- **Category:** runtime
- **Complexity:** 3/5
- **Risk level:** low
- **Tools:** python

### Command 4

```bash
echo "Frontend http://localhost:3000"
```

- **Category:** utility
- **Complexity:** 1/5
- **Risk level:** low
- **Tools:** shell

### Command 5

```bash
echo "Backend http://localhost:8000"
```

- **Category:** utility
- **Complexity:** 1/5
- **Risk level:** low
- **Tools:** shell

## 🔍 Command Patterns

- **docker-compose:** 1 occurrence(s)
- **npm:** 1 occurrence(s)
- **python:** 1 occurrence(s)

## 🚀 Optimization Opportunities

- Consider adding sources and generates for caching

## Navigation

- [← Back to All Tasks](../summaries/all-tasks.md)
- [← Back to Overview](../README.md)
- [Dependency Graph](dependency-graph.md)
