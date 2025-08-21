# Job: test-integration

## Usage Information

- **Used in workflows:** 1 times
- **Depends on:** test-frontend, test-backend
- **Executor:** e2e-executor

## Docker Images

- `cimg/node:18.17.0-browsers`
- `cimg/postgres:14.0`
- `redis:6.2-alpine`

## Run Commands

### Command 1

```bash
install-dependencies
```

### Command 2

```bash
dockerize -wait tcp://localhost:5432 -timeout 2m
dockerize -wait tcp://localhost:6379 -timeout 1m

```

### Command 3

```bash
cd api
python manage.py migrate --settings=config.settings.test
python manage.py loaddata fixtures/test_data.json --settings=config.settings.test

```

### Command 4

```bash
cd api && python manage.py runserver 0.0.0.0:8000 --settings=config.settings.test &
npm run start:test &
sleep 30

```

### Command 5

```bash
npm run test:integration -- \
  --reporter=xunit \
  --reporter-options output=test-results/integration/results.xml

```

### Command 6

```bash
mkdir -p test-results/logs
docker logs $(docker ps -aq) > test-results/logs/containers.log 2>&1 || true

```

### Command 7

```bash
notify-slack-on-failure
```

## Command Patterns Detected

- **docker:** 1 occurrences
- **npm:** 1 occurrences
- **python:** 1 occurrences
- **test:** 1 occurrences

## Suggested go-task Conversion

```yaml
test-integration:
  deps: [test-frontend, test-backend]
  cmds:
    - install-dependencies
    - dockerize -wait tcp://localhost:5432 -timeout 2m
dockerize -wait tcp://localhost:6379 -timeout 1m

    - cd api
python manage.py migrate --settings=config.settings.test
python manage.py loaddata fixtures/test_data.json --settings=config.settings.test

    - cd api && python manage.py runserver 0.0.0.0:8000 --settings=config.settings.test &
npm run start:test &
sleep 30

    - npm run test:integration -- \
  --reporter=xunit \
  --reporter-options output=test-results/integration/results.xml

    - mkdir -p test-results/logs
docker logs $(docker ps -aq) > test-results/logs/containers.log 2>&1 || true

    - notify-slack-on-failure
```

## Navigation

- [← Back to All Jobs](../summaries/all-jobs.md)
- [← Back to Overview](../README.md)
