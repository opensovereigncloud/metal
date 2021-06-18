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

### Required operators
Current operator depends on the following operators:
- [onmetal/k8s-size](https://github.com/onmetal/k8s-size)
- [onmetal/k8s-inventory](https://github.com/onmetal/k8s-inventory)
- [onmetal/k8s-network-global](https://github.com/onmetal/k8s-network-global)
- [onmetal/k8s-subnet](https://github.com/onmetal/k8s-subnet)

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
Webhooks need the certificate manager to be deployed. If there is no certificate manager in your environment, run the following 
command to deploy it:

    make install-cert-manager

### Build and install
In order to build and deploy, execute following command set, where:
- <registry_name>: name of docker registry (i.e. `localhost:5000` for local registry)
- <image_name>: name for docker image
- <image_tag>: tag for docker image
- <username>: username to access git
- <access_token>: token to access git


    # install crds. You also need to install k8s-inventory/inventory crds by yourself
    make install
    # build container and push to local registry
    make docker-build docker-push IMG="<registry_name>/<image_name>:<image_tag>" GIT_USER=<username> GIT_PASSWORD=<access_token>
    # deploy controller
    make deploy IMG="<registry_name>/<image_name>:<image_tag>"

Check `Makefile` for the full list of `make` goals with descriptions.

### Test
To run tests:

    make test

### Use
`./config/samples/` directory contains examples of manifests. They can be used to try out the controller.

    # apply config
    kubectl apply -k config/samples/
    # get resources
    kubectl get switches.machine.onmetal.de -n onmetal
    # get sample resource
    kubectl describe switches.switch.onmetal.de -n onmetal 92b9de0f-19f2-3f3b-95d0-fb668b1d3d3b

### Clean
Clean up local environment, after you've done with development:

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

    # install release "onmetal-switch" to "onmetal" namespace in case of using local docker registry
    helm install onmetal-switch ./deploy/ -n onmetal --create-namespace
    
    # in case of using another docker registry, specify <registry_name>, <image_name> and <image_tag>
    helm install --set manager.image.repository=<registry_name>/<image_name> --set manager.image.tag="<image_tag>" onmetal-switch deploy/ -n onmetal
    
    # remove release "onmetal-switch" from "onmetal" namespace
    helm uninstall onmetal-switch -n onmetal