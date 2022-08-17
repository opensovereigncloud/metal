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

	"github.com/onmetal/metal-api/common/types/base"
	"github.com/onmetal/metal-api/pkg/provider"
	"k8s.io/apimachinery/pkg/labels"
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

func NewClient(c ctrlclient.Client) (*ClientImpl, error) {
	if c == nil {
		return &ClientImpl{}, fmt.Errorf("%w", errClientNotDefined)
	}
	return &ClientImpl{client: c}, nil
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
		if err := c.client.
			Status().
			Update(context.Background(), k8sObject); err != nil {
			return err
		}
		return c.client.Update(context.Background(), k8sObject)
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

func (c *ClientImpl) Get(obj any, sa base.Metadata) error {
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

func (c *ClientImpl) List(objList any, listOptions *provider.ListOptions) error {
	switch k8sObjectList := objList.(type) {
	case ctrlclient.ObjectList:
		listOptionFilter := &ctrlclient.ListOptions{
			LabelSelector: ctrlclient.MatchingLabelsSelector{Selector: labels.SelectorFromSet(listOptions.Filter)},
			Continue:      listOptions.Pagination,
		}
		return c.client.List(
			context.Background(),
			k8sObjectList,
			listOptionFilter)
	default:
		return fmt.Errorf("%w, %T", errObjectTypeNotDefined, objList)
	}
}
