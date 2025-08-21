# Job: security-scan

## Usage Information

- **Used in workflows:** 3 times
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
npm audit --audit-level=moderate --json > security-results/npm-audit.json || true
npm audit --audit-level=moderate

```

### Command 3

```bash
if [ -n "$SNYK_TOKEN" ]; then
  npx snyk test --json > security-results/snyk.json || true
  npx snyk test
else
  echo "Skipping Snyk scan - SNYK_TOKEN not set"
fi

```

### Command 4

```bash
notify-slack-on-failure
```

## Command Patterns Detected

- **echo:** 1 occurrences
- **npm:** 1 occurrences
- **test:** 1 occurrences

## Suggested go-task Conversion

```yaml
security-scan:
  cmds:
    - install-dependencies
    - npm audit --audit-level=moderate --json > security-results/npm-audit.json || true
npm audit --audit-level=moderate

    - if [ -n "$SNYK_TOKEN" ]; then
  npx snyk test --json > security-results/snyk.json || true
  npx snyk test
else
  echo "Skipping Snyk scan - SNYK_TOKEN not set"
fi

    - notify-slack-on-failure
```

## Navigation

- [← Back to All Jobs](../summaries/all-jobs.md)
- [← Back to Overview](../README.md)
