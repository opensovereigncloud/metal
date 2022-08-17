# Scheduler Deployment

Scheduler can be deployed with kubectl or kustomize. Use may choose one that is more suitable.

### With kubectl using kustomize configs

    # deploy
    kubectl apply -k deploy/scheduler/default
    # remove
    kubectl delete -k deploy/scheduler/default

### With kustomize

    # build and apply
    kustomize build deploy/scheduler/default | kubectl apply -f -
    # build and remove
    kustomize build deploy/scheduler/default | kubectl delete -f -

