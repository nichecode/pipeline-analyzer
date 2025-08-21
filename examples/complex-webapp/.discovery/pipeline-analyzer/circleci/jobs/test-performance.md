# Job: test-performance

## Usage Information

- **Used in workflows:** 2 times
- **Depends on:** build-frontend
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
npm run start:prod &
sleep 30
npx lhci autorun --config=.lighthouserc.json

```

### Command 3

```bash
npm run test:load -- --reporter json > load-test-results.json

```

### Command 4

```bash
notify-slack-on-failure
```

## Command Patterns Detected

- **npm:** 1 occurrences

## Suggested go-task Conversion

```yaml
test-performance:
  deps: [build-frontend]
  cmds:
    - install-dependencies
    - npm run start:prod &
sleep 30
npx lhci autorun --config=.lighthouserc.json

    - npm run test:load -- --reporter json > load-test-results.json

    - notify-slack-on-failure
```

## Navigation

- [← Back to All Jobs](../summaries/all-jobs.md)
- [← Back to Overview](../README.md)
