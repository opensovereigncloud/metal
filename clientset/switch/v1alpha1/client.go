/*
Copyright 2021 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	switchv1alpha1 "github.com/onmetal/metal-api/apis/switches/v1alpha1"
)

type SwitchV1Alpha1Interface interface {
	Projects(namespace string) SwitchInterface
}

type SwitchV1Alpha1Client struct {
	restClient rest.Interface
}

func NewForConfig(c *rest.Config) (*SwitchV1Alpha1Client, error) {
	config := *c
	config.ContentConfig.GroupVersion = &switchv1alpha1.GroupVersion
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
