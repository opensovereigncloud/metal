package v1alpha1

import (
	"github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

const (
	CStatusSubresource = "status"
)

type V1Alpha1Interface interface { //nolint:revive
	Sizes(namespace string) SizeInterface
	Inventories(namespace string) InventoryInterface
	Aggregates(namespace string) AggregateInterface
}

type v1Alpha1Client struct {
	restClient rest.Interface
}

func NewForConfig(c *rest.Config) (V1Alpha1Interface, error) {
	config := *c
	config.ContentConfig.GroupVersion = &v1alpha1.GroupVersion
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &v1Alpha1Client{restClient: client}, nil
}

func (c *v1Alpha1Client) Sizes(namespace string) SizeInterface {
	return &sizeClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}

func (c *v1Alpha1Client) Inventories(namespace string) InventoryInterface {
	return &inventoryClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}

func (c *v1Alpha1Client) Aggregates(namespace string) AggregateInterface {
	return &aggregateClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}
