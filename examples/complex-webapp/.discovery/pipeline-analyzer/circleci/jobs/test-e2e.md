# Job: test-e2e

## Usage Information

- **Used in workflows:** 2 times
- **Depends on:** test-integration
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
python manage.py loaddata fixtures/e2e_data.json --settings=config.settings.test

```

### Command 4

```bash
cd api && python manage.py runserver 0.0.0.0:8000 --settings=config.settings.test &
npm run start:e2e &
sleep 45

```

### Command 5

```bash
npm run test:e2e -- \
  --reporter spec \
  --reporter junit \
  --reporter-options mochaFile=test-results/e2e/results.xml

```

### Command 6

```bash
if [ -d "cypress/screenshots" ]; then
  mkdir -p test-results/screenshots
  cp -r cypress/screenshots/* test-results/screenshots/
fi

```

### Command 7

```bash
notify-slack-on-failure
```

## Command Patterns Detected

- **npm:** 1 occurrences
- **python:** 1 occurrences
- **test:** 1 occurrences

## Suggested go-task Conversion

```yaml
test-e2e:
  deps: [test-integration]
  cmds:
    - install-dependencies
    - dockerize -wait tcp://localhost:5432 -timeout 2m
dockerize -wait tcp://localhost:6379 -timeout 1m

    - cd api
python manage.py migrate --settings=config.settings.test
python manage.py loaddata fixtures/e2e_data.json --settings=config.settings.test

    - cd api && python manage.py runserver 0.0.0.0:8000 --settings=config.settings.test &
npm run start:e2e &
sleep 45

    - npm run test:e2e -- \
  --reporter spec \
  --reporter junit \
  --reporter-options mochaFile=test-results/e2e/results.xml

    - if [ -d "cypress/screenshots" ]; then
  mkdir -p test-results/screenshots
  cp -r cypress/screenshots/* test-results/screenshots/
fi

    - notify-slack-on-failure
```

## Navigation

- [← Back to All Jobs](../summaries/all-jobs.md)
- [← Back to Overview](../README.md)
