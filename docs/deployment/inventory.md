# Inventory Deployment

Operator can be deployed with kubectl or kustomize. Use may choose one that is more suitable.

### With kubectl using kustomize configs

    # deploy
    kubectl apply -k config/default/
    # remove
    kubectl delete -k config/default/

### With kustomize

    # build and apply
    kustomize build config/default | kubectl apply -f -
    # build and remove
    kustomize build config/default | kubectl delete -f -

