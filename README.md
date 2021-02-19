# k8s-inventory

k8s-inventory is a k8s operator that allows to store results of machine inventarization process done with [inventory CLI](https://github.com/onmetal/inventory) in a form of k8s resource.

k8s-inventory implements corresponding resource specification, controller and golang client for it.

## Getting started 

### Required tools

Following tools are required to make changes on that package.

- [make](https://www.gnu.org/software/make/) - to execute build goals
- [golang](https://golang.org/) - to compile the code
- [minikube](https://minikube.sigs.k8s.io/) or access to k8s cluster - to deply and test the result
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) - to interact with k8s cluster via CLI
- [kustomize](https://kustomize.io/) - to generate deployment configs
- [kubebuilder](https://book.kubebuilder.io) - framework to build operators
- [operator framework](https://operatorframework.io/) - framework to maintain project structure
- [helm](https://helm.sh/) - to work with helm charts

If you have to build Docker images on your host, you also need to have [Docker](https://www.docker.com/) or its alternative installed.

### Prepare environment

If you have access to the docker registry and k8s installation that you can use for development purposes, you may skip 
corresponding steps.

Otherwise, create a local instance of docker registry and k8s.

    # start minikube
    minikube start
    # enable registry
    minikube addons enable registry
    # run proxy to registry
    docker run --rm -d --name registry-bridge --network=host alpine ash -c "apk add socat && socat TCP-LISTEN:5000,reuseaddr,fork TCP:$(minikube ip):5000"

### Build and install

In order to build and deploy, execute following command set.

    # generate code and configs
    make install
    # build container and push to local registry
    make docker-build docker-push IMG="localhost:5000/k8s-inventory:latest"
    # deploy controller
    make deploy IMG="localhost:5000/k8s-inventory:latest"

Check `Makefile` for the full list of `make` goals with descriptions. 

### Use

`./config/samples/` directory contains examples of manifests. They can be used to try out the controller.

    # apply config
    kubectl apply -f config/samples/machine_v1alpha1_inventory.yaml
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

## Extending API

### Adding a new version

One should not modify API once it got to be used. 

Instead, in order to introduce breaking changes to the API, a new API version should be created.

First, move the existing controller to the different file, as generator will try to put a new controller into the same location, e.g.

    mv controllers/inventory_controller.go controllers/inventory_v1alpha1_controller.go

After that, add a new API version

    operator-sdk create api --group machine --version v1alpha2 --kind Inventory --resource --controller

Do modifications in a new CR, add a new controller to `main.go`.

Following actions should be applied to other parts of project:  
- regenerate code and configs with `make install`
- add a client to client set for the new API version  
- alter Helm chart with new CRD spec

### Deprecating old APIs

Since there is no version deprecation marker available now, old APIs may be deprecated with `kustomize` patches

Describe deprecation status and deprecation warning in patch file, e.g. `crd_patch.yaml`

```
- op: add
  path: "/spec/versions/0/deprecated"
  value: true
- op: add
  path: "/spec/versions/0/deprecationWarning"
  value: "This API version is deprecated. Check documentation for migration instructions."
```

Add patch instructions to `kustomization.yaml`

```
patchesJson6902:
  - target:
      version: v1
      group: apiextensions.k8s.io
      kind: CustomResourceDefinition
      name: inventories.machine.onmetal.de
    path: crd_patch.yaml
```

When you are ready to drop the support for the API version, give CRD a `+kubebuilder:skipversion` marker, 
or just remove it completely from the code base.

This includes:
- API objects
- client
- controller

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

### With Helm chart

    # install release "dev" to "onmetal" namespace
    helm install dev ./chart/ -n onmetal --create-namespace
    # remove release "dev" from "onmetal" namespace
    helm uninstall dev -n onmetal

## Consuming API

Package provides a client library written in go for programmatic interactions with API.  

Clients for corresponding API versions are located in `clientset/` and act similar to [client-go](https://github.com/kubernetes/client-go).

Below are two examples, for inbound (from the pod deployed on a cluster) and outbound (from the program running on 3rd party resources) interactions. 

### Inbound example

```go
import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	apiv1alpha1 "github.com/onmetal/k8s-inventory/api/v1alpha1"
	clientv1alpha1 "github.com/onmetal/k8s-inventory/clientset/v1alpha1"
)

func inbound() error {
	// get config from environment
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	// register CRD types in local client scheme
	if err := apiv1alpha1.AddToScheme(scheme.Scheme); err != nil {
		return errors.Wrap(err, "unable to add registered types to client scheme")
	}

	// create a client from obtained  configuration
	clientset, err := clientv1alpha1.NewForConfig(config)
	if err != nil {
		return errors.Wrap(err, "unable to build clientset from config")
	}

	// get a client for particular namespace
	client := clientset.Inventories("default")

	// request a list of resources
	list, err := client.List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return errors.Wrap(err, "unable to get list of resources")
	}

	// print names of resources
	for _, r := range list.Items {
		fmt.Println(r.Name)
	}

	return nil
}
```

### Outbound example

```go
import (
    "context"
    "fmt"
    "path/filepath"
    
    "github.com/pkg/errors"
    "github.com/spf13/pflag"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/client-go/util/homedir"
    
    "k8s.io/client-go/kubernetes/scheme"
    
    apiv1alpha1 "github.com/onmetal/k8s-inventory/api/v1alpha1"
    clientv1alpha1 "github.com/onmetal/k8s-inventory/clientset/v1alpha1"
)

func outbound() error {
	// make default path to kubeconfig empty
	var kubeconfigDefaultPath string

	// if there is a home directory in environment,
	// alter the defalut path to ~/.kube/config
	if home := homedir.HomeDir(); home != "" {
		kubeconfigDefaultPath = filepath.Join(home, ".kube", "config")
	}

	// configure the kubeconfig CLI flag with dafault value
	kubeconfig := pflag.StringP("kubeconfig", "k", kubeconfigDefaultPath, "path to kubeconfig")
	// configure k8s namespace flag with "default" default value
	namespace := pflag.StringP("namespace", "n", "default", "k8s namespace")
	// parse flags
	pflag.Parse()

	// read in kubeconfig file and build configuration
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return errors.Wrapf(err, "unable to read kubeconfig from path %s", kubeconfig)
	}

	// register CRD types in local client scheme
	if err := apiv1alpha1.AddToScheme(scheme.Scheme); err != nil {
		return errors.Wrap(err, "unable to add registered types to client scheme")
	}

	// create a client from obtained  configuration
	clientset, err := clientv1alpha1.NewForConfig(config)
	if err != nil {
		return errors.Wrap(err, "unable to build clientset from config")
	}

	// get a client for particular namespace
	client := clientset.Inventories(*namespace)
	
	// request a list of resources
	list, err := client.List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return errors.Wrap(err, "unable to get list of resources")
	}
	
	// print names of resources
	for _, r := range list.Items {
		fmt.Println(r.Name)
	}

	return nil
}
```

### Libraries

Apart from libraries required to developm the operator itself, project has some other includes:

- [messagediff](https://github.com/d4l3k/messagediff) - used to compute the diff for resources on update 
