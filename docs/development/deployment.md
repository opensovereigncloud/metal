# Deployment
Operator can be deployed with `kubectl`, `kustomize` or `Helm`. 

You are welcome to choose an option that is suitable for you.

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

### With Helm

    # install release "onmetal-switch" to "onmetal" namespace in case of using local docker registry
    helm install onmetal-switch ./deploy/ -n onmetal --create-namespace
    
    # in case of using another docker registry, specify <registry_name>, <image_name> and <image_tag>
    helm install --set manager.image.repository=<registry_name>/<image_name> --set manager.image.tag="<image_tag>" onmetal-switch deploy/ -n onmetal
    
    # remove release "onmetal-switch" from "onmetal" namespace
    helm uninstall onmetal-switch -n onmetal