# k8s-inventory
k8s operator for inventory CRD

### Prepare environment

    # start minikube
    minikube start
    # enable registry
    minikube addons enable registry
    # run proxy to registry
    docker run --rm -d --name registry-bridge --network=host alpine ash -c "apk add socat && socat TCP-LISTEN:5000,reuseaddr,fork TCP:$(minikube ip):5000"

### Build and install

    # generate configs
    make install
    # build container and push to local registry
    make docker-build docker-push IMG="localhost:5000/k8s-inventory:latest"
    # deploy controller
    make deploy IMG="localhost:5000/k8s-inventory:latest"

### Use

    # apply config
    kubectl apply -f config/samples/machine_v1alpha1_inventory.yaml
    # get resources
    kubectl get inventories
    # get sample resource
    kubectl describe inventory inventory-sample

### Clean

    # clean
    kustomize build config/default | kubectl delete -f -
    # remove registry bridge
    docker stop registry-bridge
    # stop minikube
    minikube stop