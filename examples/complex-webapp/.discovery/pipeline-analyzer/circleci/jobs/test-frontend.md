# Job: test-frontend

## Usage Information

- **Used in workflows:** 1 times
- **Depends on:** lint-frontend
- **Executor:** node-executor

## Docker Images

- `cimg/node:18.17.0`

## Run Commands

### Command 1

```bash
install-dependencies
```

### Command 2

```bash
mkdir -p test-results/jest
npm run test:unit -- \
  --ci \
  --coverage \
  --watchAll=false \
  --testResultsProcessor="jest-junit" \
  --coverageReporters="lcov" "text" "cobertura" \
  --maxWorkers=2

```

### Command 3

```bash
if [ -n "$CODECOV_TOKEN" ]; then
  npx codecov
fi

```

### Command 4

```bash
notify-slack-on-failure
```

## Command Patterns Detected

- **npm:** 1 occurrences

## Suggested go-task Conversion

```yaml
test-frontend:
  deps: [lint-frontend]
  cmds:
    - install-dependencies
    - mkdir -p test-results/jest
npm run test:unit -- \
  --ci \
  --coverage \
  --watchAll=false \
  --testResultsProcessor="jest-junit" \
  --coverageReporters="lcov" "text" "cobertura" \
  --maxWorkers=2

    - if [ -n "$CODECOV_TOKEN" ]; then
  npx codecov
fi

    - notify-slack-on-failure
```

## Navigation

- [← Back to All Jobs](../summaries/all-jobs.md)
- [← Back to Overview](../README.md)
