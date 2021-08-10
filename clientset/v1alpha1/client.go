package v1alpha1

import (
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	"github.com/onmetal/switch-operator/api/v1alpha1"
)

type SwitchV1Alpha1Interface interface {
	Projects(namespace string) SwitchInterface
}

type SwitchV1Alpha1Client struct {
	restClient rest.Interface
}

func NewForConfig(c *rest.Config) (*SwitchV1Alpha1Client, error) {
	config := *c
	config.ContentConfig.GroupVersion = &v1alpha1.GroupVersion
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &SwitchV1Alpha1Client{restClient: client}, nil
}

func (c *SwitchV1Alpha1Client) Switches(namespace string) SwitchInterface {
	return &switchClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}
