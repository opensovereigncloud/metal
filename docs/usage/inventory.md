# Inventory Usage

## Extending API

### Adding a new version

One should not modify API once it got to be used. 

Instead, in order to introduce breaking changes to the API, a new API version should be created.

First, move the existing controller to the different file, as generator will try to put a new controller into the same location, e.g.

    mv controllers/inventory/inventory_controller.go controllers/inventory/inventory_v1alpha1_controller.go

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
      name: inventories.metal.ironcore.dev
    path: crd_patch.yaml
```

When you are ready to drop the support for the API version, give CRD a `+kubebuilder:skipversion` marker, 
or just remove it completely from the code base.

This includes:
- API objects
- client
- controller

## Consuming API

Package provides a client library written in go for programmatic interactions with API.  

Clients for corresponding API versions are located in `clientset/inventory` and act similar to [client-go](https://github.com/kubernetes/client-go).

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

	apiv1alpha1 "github.com/ironcore-dev/metal/apis/inventory/v1alpha1"
	clientv1alpha1 "github.com/ironcore-dev/metal/clientset/inventory/v1alpha1"
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
    
	apiv1alpha1 "github.com/ironcore-dev/metal/apis/inventory/v1alpha1"
	clientv1alpha1 "github.com/ironcore-dev/metal/clientset/inventory/v1alpha1"
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
