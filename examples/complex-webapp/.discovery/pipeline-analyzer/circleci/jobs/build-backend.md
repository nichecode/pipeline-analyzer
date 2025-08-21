# Job: build-backend

## Usage Information

- **Used in workflows:** 1 times
- **Depends on:** test-backend
- **Executor:** python-executor

## Docker Images

- `cimg/python:3.11`

## Run Commands

### Command 1

```bash
cd api
python setup.py sdist bdist_wheel

```

### Command 2

```bash
cd api
pip install twine
twine check dist/*

```

### Command 3

```bash
notify-slack-on-failure
```

## Command Patterns Detected

- **pip:** 1 occurrences
- **python:** 1 occurrences

## Suggested go-task Conversion

```yaml
build-backend:
  deps: [test-backend]
  cmds:
    - cd api
python setup.py sdist bdist_wheel

    - cd api
pip install twine
twine check dist/*

    - notify-slack-on-failure
```

## Navigation

- [← Back to All Jobs](../summaries/all-jobs.md)
- [← Back to Overview](../README.md)
