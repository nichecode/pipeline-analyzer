# Job: lint-backend

## Usage Information

- **Used in workflows:** 1 times
- **Executor:** python-executor

## Docker Images

- `cimg/python:3.11`

## Run Commands

### Command 1

```bash
python -m pip install --upgrade pip
pip install -r api/requirements.txt
pip install -r api/requirements-dev.txt

```

### Command 2

```bash
cd api
flake8 . --format=junit-xml --output-file=../test-results/flake8/results.xml
black --check .
mypy . --junit-xml ../test-results/mypy/results.xml

```

## Command Patterns Detected

- **./:** 1 occurrences
- **pip:** 1 occurrences
- **python:** 1 occurrences

## Suggested go-task Conversion

```yaml
lint-backend:
  cmds:
    - python -m pip install --upgrade pip
pip install -r api/requirements.txt
pip install -r api/requirements-dev.txt

    - cd api
flake8 . --format=junit-xml --output-file=../test-results/flake8/results.xml
black --check .
mypy . --junit-xml ../test-results/mypy/results.xml

```

## Navigation

- [← Back to All Jobs](../summaries/all-jobs.md)
- [← Back to Overview](../README.md)
