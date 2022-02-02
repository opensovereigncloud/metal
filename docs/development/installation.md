# Installation
## Prerequisites
The following tools are required to make changes on that package.

- [make](https://www.gnu.org/software/make/) - to execute build goals
- [golang](https://golang.org/) - to compile the code
- [minikube](https://minikube.sigs.k8s.io/) or access to k8s cluster - to deploy and test operator
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) - to interact with k8s cluster via CLI
- [kustomize](https://kustomize.io/) - to generate deployment configs
- [kubebuilder](https://book.kubebuilder.io) - framework to build operators
- [operator framework](https://operatorframework.io/) - framework to maintain project structure
- [helm](https://helm.sh/) - to work with helm charts

 > If you have to build Docker images on your host,
you also need to have [Docker](https://www.docker.com/) or its alternative installed.

## Required operators
Current operator depends on the following operators:
- [onmetal/k8s-inventory](https://github.com/onmetal/k8s-inventory)
- [onmetal/ipam](https://github.com/onmetal/ipam)

## Prepare environment
If you have access to the docker registry and k8s installation that you can use for development purposes, you may skip corresponding steps. 

Otherwise, create a local instance of k8s and registry.

    # start minikube
    minikube start
    # enable registry
    minikube addons enable registry
    # run proxy to registry
    docker run --rm -d --name registry-bridge --network=host alpine ash -c "apk add socat && socat TCP-LISTEN:5000,reuseaddr,fork TCP:$(minikube ip):5000"

## Webhooks
Webhooks need the certificate manager to be deployed. If there is no certificate manager in your environment, run the following 
command to deploy it:

    make install-cert-manager

## Build and install
In order to build and deploy, execute following command set:

> Docker registry is required to build and push an image.

    # used variables:
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

