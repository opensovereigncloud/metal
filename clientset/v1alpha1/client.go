package v1alpha1

import (
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	"github.com/onmetal/k8s-inventory/api/v1alpha1"
)

type InventoryV1Alpha1Interface interface {
	Projects(namespace string) InventoryInterface
}

type InventoryV1Alpha1Client struct {
	restClient rest.Interface
}

func NewForConfig(c *rest.Config) (*InventoryV1Alpha1Client, error) {
	config := *c
	config.ContentConfig.GroupVersion = &v1alpha1.GroupVersion
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &InventoryV1Alpha1Client{restClient: client}, nil
}

func (c *InventoryV1Alpha1Client) Inventories(namespace string) InventoryInterface {
	return &inventoryClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}
