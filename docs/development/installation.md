# Installation

## Prerequisites
The following tools are required to make changes on that package.

- [make](https://www.gnu.org/software/make/) - to execute build goals
- [golang](https://golang.org/) - to compile the code
- [minikube](https://minikube.sigs.k8s.io/) or access to k8s cluster - to deploy and test operator
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) - to interact with k8s cluster via CLI
- [kustomize](https://kustomize.io/) - to generate deployment configs
- [kubebuilder](https://book.kubebuilder.io) - framework to build operators

 > If you have to build Docker images on your host,
you also need to have [Docker](https://www.docker.com/) or its alternative installed.

## Required operators
Current operator depends on the following operators:
- [onmetal/ipam](https://github.com/onmetal/ipam)
- [onmetal/oob-operator](https://github.com/onmetal/oob-operator)


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

## Inventory Usage

`./config/samples/` directory contains examples of manifests. They can be used to try out the controller.

    # apply inventory config
    kubectl apply -f config/samples/machine_v1alpha1_inventory.yaml
    # apply size config
    kubectl apply -f config/samples/machine_v1alpha1_size.yaml
    # apply aggregate config
    kubectl apply -f config/samples/machine_v1alpha1_aggregate.yaml
    # get resources
    kubectl get inventories
    # get sample resource
    kubectl describe inventory a967954c-3475-11b2-a85c-84d8b4f8cd2d

### Clean

After development is done, clean up local environment.

    # generate deployment config and delete corresponding entities
    kustomize build config/default | kubectl delete -f -
    # remove registry bridge
    docker stop registry-bridge
    # stop minikube
    minikube stop