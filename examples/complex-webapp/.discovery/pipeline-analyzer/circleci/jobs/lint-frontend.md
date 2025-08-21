# Job: lint-frontend

## Usage Information

- **Used in workflows:** 1 times
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
npm run lint -- --format junit --output-file test-results/eslint/results.xml
npm run lint:css -- --formatter junit --output test-results/stylelint/results.xml

```

### Command 3

```bash
npm run format:check
```

### Command 4

```bash
notify-slack-on-failure
```

## Command Patterns Detected

- **npm:** 1 occurrences

## Suggested go-task Conversion

```yaml
lint-frontend:
  cmds:
    - install-dependencies
    - npm run lint -- --format junit --output-file test-results/eslint/results.xml
npm run lint:css -- --formatter junit --output test-results/stylelint/results.xml

    - npm run format:check
    - notify-slack-on-failure
```

## Navigation

- [← Back to All Jobs](../summaries/all-jobs.md)
- [← Back to Overview](../README.md)
