# Job: test-backend

## Usage Information

- **Used in workflows:** 1 times
- **Depends on:** lint-backend
- **Executor:** python-executor

## Docker Images

- `cimg/python:3.11`
- `cimg/postgres:14.0`
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
dockerize -wait tcp://localhost:5432 -timeout 1m
```

### Command 3

```bash
cd api
python manage.py migrate --settings=config.settings.test

```

### Command 4

```bash
cd api
python -m pytest \
  --junitxml=../test-results/pytest/results.xml \
  --cov=. \
  --cov-report=html:../coverage/python \
  --cov-report=xml:../coverage/python/coverage.xml \
  tests/

```

### Command 5

```bash
notify-slack-on-failure
```

## Command Patterns Detected

- **./:** 1 occurrences
- **pip:** 1 occurrences
- **python:** 1 occurrences
- **test:** 1 occurrences

## Suggested go-task Conversion

```yaml
test-backend:
  deps: [lint-backend]
  cmds:
    - python -m pip install --upgrade pip
pip install -r api/requirements.txt
pip install -r api/requirements-dev.txt

    - dockerize -wait tcp://localhost:5432 -timeout 1m
    - cd api
python manage.py migrate --settings=config.settings.test

    - cd api
python -m pytest \
  --junitxml=../test-results/pytest/results.xml \
  --cov=. \
  --cov-report=html:../coverage/python \
  --cov-report=xml:../coverage/python/coverage.xml \
  tests/

    - notify-slack-on-failure
```

## Navigation

- [← Back to All Jobs](../summaries/all-jobs.md)
- [← Back to Overview](../README.md)
