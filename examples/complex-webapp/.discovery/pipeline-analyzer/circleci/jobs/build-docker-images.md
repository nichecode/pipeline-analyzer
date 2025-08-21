# Job: build-docker-images

## Usage Information

- **Used in workflows:** 1 times
- **Depends on:** build-frontend, build-backend, security-scan
- **Executor:** docker-executor

## Docker Images

- `cimg/base:stable`

## Run Commands

### Command 1

```bash
# Build frontend image
docker build \
  --tag webapp-frontend:${CIRCLE_SHA1} \
  --tag webapp-frontend:latest \
  --build-arg BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ') \
  --build-arg VCS_REF=${CIRCLE_SHA1} \
  --file docker/frontend/Dockerfile .

# Build backend image
docker build \
  --tag webapp-backend:${CIRCLE_SHA1} \
  --tag webapp-backend:latest \
  --build-arg BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ') \
  --build-arg VCS_REF=${CIRCLE_SHA1} \
  --file docker/backend/Dockerfile .

```

### Command 2

```bash
# Install trivy
curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh -s -- -b /tmp

# Scan images
/tmp/trivy image --format json --output frontend-scan.json webapp-frontend:${CIRCLE_SHA1}
/tmp/trivy image --format json --output backend-scan.json webapp-backend:${CIRCLE_SHA1}

# Check for HIGH and CRITICAL vulnerabilities
/tmp/trivy image --severity HIGH,CRITICAL --exit-code 1 webapp-frontend:${CIRCLE_SHA1}
/tmp/trivy image --severity HIGH,CRITICAL --exit-code 1 webapp-backend:${CIRCLE_SHA1}

```

### Command 3

```bash
if [ -n "$DOCKER_HUB_TOKEN" ] && [ "$CIRCLE_BRANCH" == "main" ]; then
  echo "$DOCKER_HUB_TOKEN" | docker login -u "$DOCKER_HUB_USER" --password-stdin
  
  docker push webapp-frontend:${CIRCLE_SHA1}
  docker push webapp-frontend:latest
  docker push webapp-backend:${CIRCLE_SHA1}
  docker push webapp-backend:latest
else
  echo "Skipping Docker push - not main branch or credentials not available"
fi

```

### Command 4

```bash
notify-slack-on-failure
```

## Command Patterns Detected

- **curl:** 1 occurrences
- **docker:** 1 occurrences
- **echo:** 1 occurrences
- **test:** 1 occurrences

## Suggested go-task Conversion

```yaml
build-docker-images:
  deps: [build-frontend, build-backend, security-scan]
  cmds:
    - # Build frontend image
docker build \
  --tag webapp-frontend:${CIRCLE_SHA1} \
  --tag webapp-frontend:latest \
  --build-arg BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ') \
  --build-arg VCS_REF=${CIRCLE_SHA1} \
  --file docker/frontend/Dockerfile .

# Build backend image
docker build \
  --tag webapp-backend:${CIRCLE_SHA1} \
  --tag webapp-backend:latest \
  --build-arg BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ') \
  --build-arg VCS_REF=${CIRCLE_SHA1} \
  --file docker/backend/Dockerfile .

    - # Install trivy
curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh -s -- -b /tmp

# Scan images
/tmp/trivy image --format json --output frontend-scan.json webapp-frontend:${CIRCLE_SHA1}
/tmp/trivy image --format json --output backend-scan.json webapp-backend:${CIRCLE_SHA1}

# Check for HIGH and CRITICAL vulnerabilities
/tmp/trivy image --severity HIGH,CRITICAL --exit-code 1 webapp-frontend:${CIRCLE_SHA1}
/tmp/trivy image --severity HIGH,CRITICAL --exit-code 1 webapp-backend:${CIRCLE_SHA1}

    - if [ -n "$DOCKER_HUB_TOKEN" ] && [ "$CIRCLE_BRANCH" == "main" ]; then
  echo "$DOCKER_HUB_TOKEN" | docker login -u "$DOCKER_HUB_USER" --password-stdin
  
  docker push webapp-frontend:${CIRCLE_SHA1}
  docker push webapp-frontend:latest
  docker push webapp-backend:${CIRCLE_SHA1}
  docker push webapp-backend:latest
else
  echo "Skipping Docker push - not main branch or credentials not available"
fi

    - notify-slack-on-failure
```

## Navigation

- [← Back to All Jobs](../summaries/all-jobs.md)
- [← Back to Overview](../README.md)
