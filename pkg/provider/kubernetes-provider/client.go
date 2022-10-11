// /*
// Copyright (c) 2021 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// */

package kubernetesprovider

import (
	"context"
	"errors"
	"fmt"

	"github.com/onmetal/metal-api/types/common"
	typesv1 "k8s.io/apimachinery/pkg/types"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	errObjectTypeNotDefined = errors.New("set object type not defined")
	errClientNotDefined     = errors.New("kubernetes client not defined")
)

type ClientImpl struct {
	client ctrlclient.Client
}

func NewClient(c ctrlclient.Client) (ctrlclient.Client, error) {
	if c == nil {
		return nil, errClientNotDefined
	}
	return c, nil
}

func (c *ClientImpl) Create(obj any) error {
	switch k8sObject := obj.(type) {
	case ctrlclient.Object:
		return c.client.Create(context.Background(), k8sObject)
	default:
		return fmt.Errorf("%w, %T", errObjectTypeNotDefined, obj)
	}
}

func (c *ClientImpl) Update(obj any) error {
	switch k8sObject := obj.(type) {
	case ctrlclient.Object:
		ctx := context.Background()
		return c.
			client.
			Update(ctx, k8sObject)
	default:
		return fmt.Errorf("%w, %T", errObjectTypeNotDefined, obj)
	}
}

func (c *ClientImpl) StatusUpdate(obj any) error {
	switch k8sObject := obj.(type) {
	case ctrlclient.Object:
		ctx := context.Background()
		return c.client.
			Status().
			Update(ctx, k8sObject)

	default:
		return fmt.Errorf("%w, %T", errObjectTypeNotDefined, obj)
	}
}

func (c *ClientImpl) Delete(obj any) error {
	switch k8sObject := obj.(type) {
	case ctrlclient.Object:
		return c.client.Delete(context.Background(), k8sObject)
	default:
		return fmt.Errorf("%w, %T", errObjectTypeNotDefined, obj)
	}
}

func (c *ClientImpl) Patch(obj any, patch []byte) error {
	switch k8sObject := obj.(type) {
	case ctrlclient.Object:
		return c.client.Patch(
			context.Background(),
			k8sObject,
			ctrlclient.RawPatch(typesv1.MergePatchType, patch))
	default:
		return fmt.Errorf("%w, %T", errObjectTypeNotDefined, obj)
	}
}

func (c *ClientImpl) Get(obj any, sa common.Metadata) error {
	switch k8sObject := obj.(type) {
	case ctrlclient.Object:
		return c.client.Get(
			context.Background(),
			typesv1.NamespacedName{
				Name:      sa.Name(),
				Namespace: sa.Namespace(),
			}, k8sObject)
	default:
		return fmt.Errorf("%w, %T", errObjectTypeNotDefined, obj)
	}
}
