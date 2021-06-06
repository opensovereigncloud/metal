# switch-operator
switch-operator is a k8s controller that includes the following custom resources:  
- `Switch` defines switch object created from collected inventory and handles switch's state
- `SwitchAssignment` defines created by user assignment of top-level spine role

### Required tools
The following tools are required to make changes on that package.

- [make](https://www.gnu.org/software/make/) - to execute build goals
- [golang](https://golang.org/) - to compile the code
- [minikube](https://minikube.sigs.k8s.io/) or access to k8s cluster - to deploy and test operator
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) - to interact with k8s cluster via CLI
- [kustomize](https://kustomize.io/) - to generate deployment configs
- [kubebuilder](https://book.kubebuilder.io) - framework to build operators
- [operator framework](https://operatorframework.io/) - framework to maintain project structure
- [helm](https://helm.sh/) - to work with helm charts

If you have to build Docker images on your host,
you also need to have [Docker](https://www.docker.com/) or its alternative installed.

### Prepare environment
If you have access to the docker registry and k8s installation that you can use for development purposes, you may skip
corresponding steps. Otherwise, create a local instance of k8s and registry.

    # start minikube
    minikube start
    # enable registry
    minikube addons enable registry
    # run proxy to registry
    docker run --rm -d --name registry-bridge --network=host alpine ash -c "apk add socat && socat TCP-LISTEN:5000,reuseaddr,fork TCP:$(minikube ip):5000"

### Webhooks
If certificates manager is not configured in your environment, run the following command to install it:

    make install-cert-manager

### Build and install
In order to build and deploy, execute following command set, where `localhost:5000` in `IMG` arguments is a URL to access registry.
Please note that during the build of the docker image, detected tests will be run. To use existing cluster for tests (for example it
could be minikube cluster), you need to set environment variable:

    export USE_EXISTING_CLUSTER=true

    # install crds. You also need to install k8s-inventory/inventory crds by yourself
    make install
    # build container and push to local registry
    make docker-build docker-push IMG="localhost:5000/switch-operator:latest" GIT_USER=yourusername GIT_PASSWORD=youraccesstoken
    # deploy controller
    make deploy IMG="localhost:5000/switch-operator:latest"

Check `Makefile` for the full list of `make` goals with descriptions.

### Test
To run tests:

    # run tests on existing cluster
    USE_EXISTING_CLUSTER=true make test
    # run tests on dedicated cluster
    USE_EXISTING_CLUSTER=false make test

### Use
`./config/samples/` directory contains examples of manifests. They can be used to try out the controller.

    # apply config
    kubectl apply -f config/samples/switch.onmetal.de_v1alpha1_switch.yaml
    # get resources
    kubectl get switches.machine.onmetal.de -n onmetal
    # get sample resource
    kubectl describe switches.switch.onmetal.de -n onmetal 8223ab8c-ad85-cabf-4c75-7217629ffece

### Clean
After development is done, clean up local environment.

    # generate deployment config and delete corresponding entities
    kustomize build config/default | kubectl delete -f -
    # remove registry bridge
    docker stop registry-bridge
    # stop minikube
    minikube stop

## Deployment
Operator can be deployed with kubectl, kustomize or Helm. Use may choose one that is more suitable.

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

    # install release "onmetal-switch" to "onmetal" namespace
    helm install onmetal-switch ./deploy/ -n onmetal --create-namespace
    
    # remove release "onmetal-switch" from "onmetal" namespace
    helm uninstall onmetal-switch -n onmetal