# Job: deploy-staging

## Usage Information

- **Used in workflows:** 1 times
- **Depends on:** build-docker-images, test-integration
- **Executor:** docker-executor

## Docker Images

- `cimg/base:stable`

## Run Commands

### Command 1

```bash
aws-cli/setup
```

### Command 2

```bash
kubernetes/install-kubectl
```

### Command 3

```bash
# Update Kubernetes manifests with new image tags
sed -i "s/IMAGE_TAG/${CIRCLE_SHA1}/g" k8s/staging/*.yaml

# Apply Kubernetes manifests
kubectl apply -f k8s/staging/ --context=staging-cluster

# Wait for deployment rollout
kubectl rollout status deployment/webapp-frontend --timeout=600s --context=staging-cluster
kubectl rollout status deployment/webapp-backend --timeout=600s --context=staging-cluster

```

### Command 4

```bash
sleep 60
npm run test:smoke -- --env=staging

```

### Command 5

```bash
notify-slack-on-failure
```

## Command Patterns Detected

- **npm:** 1 occurrences
- **sed:** 1 occurrences

## Suggested go-task Conversion

```yaml
deploy-staging:
  deps: [build-docker-images, test-integration]
  cmds:
    - aws-cli/setup
    - kubernetes/install-kubectl
    - # Update Kubernetes manifests with new image tags
sed -i "s/IMAGE_TAG/${CIRCLE_SHA1}/g" k8s/staging/*.yaml

# Apply Kubernetes manifests
kubectl apply -f k8s/staging/ --context=staging-cluster

# Wait for deployment rollout
kubectl rollout status deployment/webapp-frontend --timeout=600s --context=staging-cluster
kubectl rollout status deployment/webapp-backend --timeout=600s --context=staging-cluster

    - sleep 60
npm run test:smoke -- --env=staging

    - notify-slack-on-failure
```

## Navigation

- [← Back to All Jobs](../summaries/all-jobs.md)
- [← Back to Overview](../README.md)
