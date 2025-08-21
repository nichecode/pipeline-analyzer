# Job: deploy-production

## Usage Information

- **Used in workflows:** 1 times
- **Depends on:** hold-for-approval
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
# Update Kubernetes manifests
sed -i "s/IMAGE_TAG/${CIRCLE_SHA1}/g" k8s/production/*.yaml

# Rolling deployment with canary
kubectl apply -f k8s/production/ --context=production-cluster

# Canary deployment
kubectl patch deployment webapp-frontend -p '{"spec":{"replicas":1}}' --context=production-cluster
kubectl patch deployment webapp-backend -p '{"spec":{"replicas":1}}' --context=production-cluster

# Wait and validate canary
kubectl rollout status deployment/webapp-frontend --timeout=300s --context=production-cluster
kubectl rollout status deployment/webapp-backend --timeout=300s --context=production-cluster

# Run production smoke tests
sleep 60
npm run test:smoke -- --env=production

# Scale up to full deployment
kubectl patch deployment webapp-frontend -p '{"spec":{"replicas":3}}' --context=production-cluster
kubectl patch deployment webapp-backend -p '{"spec":{"replicas":2}}' --context=production-cluster

# Final rollout check
kubectl rollout status deployment/webapp-frontend --timeout=600s --context=production-cluster
kubectl rollout status deployment/webapp-backend --timeout=600s --context=production-cluster

```

### Command 4

```bash
notify-slack-on-failure
```

## Command Patterns Detected

- **npm:** 1 occurrences
- **sed:** 1 occurrences

## Suggested go-task Conversion

```yaml
deploy-production:
  deps: [hold-for-approval]
  cmds:
    - aws-cli/setup
    - kubernetes/install-kubectl
    - # Update Kubernetes manifests
sed -i "s/IMAGE_TAG/${CIRCLE_SHA1}/g" k8s/production/*.yaml

# Rolling deployment with canary
kubectl apply -f k8s/production/ --context=production-cluster

# Canary deployment
kubectl patch deployment webapp-frontend -p '{"spec":{"replicas":1}}' --context=production-cluster
kubectl patch deployment webapp-backend -p '{"spec":{"replicas":1}}' --context=production-cluster

# Wait and validate canary
kubectl rollout status deployment/webapp-frontend --timeout=300s --context=production-cluster
kubectl rollout status deployment/webapp-backend --timeout=300s --context=production-cluster

# Run production smoke tests
sleep 60
npm run test:smoke -- --env=production

# Scale up to full deployment
kubectl patch deployment webapp-frontend -p '{"spec":{"replicas":3}}' --context=production-cluster
kubectl patch deployment webapp-backend -p '{"spec":{"replicas":2}}' --context=production-cluster

# Final rollout check
kubectl rollout status deployment/webapp-frontend --timeout=600s --context=production-cluster
kubectl rollout status deployment/webapp-backend --timeout=600s --context=production-cluster

    - notify-slack-on-failure
```

## Navigation

- [← Back to All Jobs](../summaries/all-jobs.md)
- [← Back to Overview](../README.md)
