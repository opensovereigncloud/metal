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
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	switchv1alpha1 "github.com/onmetal/switch-operator/api/v1alpha1"
)

const (
	CSwitchesResourceType = "switches"
	CStatusSubresource    = "status"
)

type SwitchInterface interface {
	Get(context.Context, string, metav1.GetOptions) (*switchv1alpha1.Switch, error)
	List(context.Context, metav1.ListOptions) (*switchv1alpha1.SwitchList, error)
	Watch(context.Context, metav1.ListOptions) (watch.Interface, error)
	Create(context.Context, *switchv1alpha1.Switch, metav1.CreateOptions) (*switchv1alpha1.Switch, error)
	Update(context.Context, *switchv1alpha1.Switch, metav1.UpdateOptions) (*switchv1alpha1.Switch, error)
	UpdateStatus(context.Context, *switchv1alpha1.Switch, metav1.UpdateOptions) (*switchv1alpha1.Switch, error)
	Delete(context.Context, string, metav1.DeleteOptions) error
	DeleteCollection(context.Context, metav1.DeleteOptions, metav1.ListOptions) error
	Patch(context.Context, string, types.PatchType, []byte, metav1.PatchOptions, ...string) (*switchv1alpha1.Switch, error)
}

type switchClient struct {
	restClient rest.Interface
	ns         string
}

func (c *switchClient) Get(ctx context.Context, name string, opts metav1.GetOptions) (*switchv1alpha1.Switch, error) {
	result := &switchv1alpha1.Switch{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource(CSwitchesResourceType).
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(result)

	return result, err
}

func (c *switchClient) List(ctx context.Context, opts metav1.ListOptions) (*switchv1alpha1.SwitchList, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result := &switchv1alpha1.SwitchList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource(CSwitchesResourceType).
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)

	return result, err
}

func (c *switchClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	watcher, err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource(CSwitchesResourceType).
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)

	return watcher, err
}

func (c *switchClient) Create(ctx context.Context, obj *switchv1alpha1.Switch, opts metav1.CreateOptions) (*switchv1alpha1.Switch, error) {
	result := &switchv1alpha1.Switch{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource(CSwitchesResourceType).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(obj).
		Do(ctx).
		Into(result)

	return result, err
}

func (c *switchClient) Update(ctx context.Context, obj *switchv1alpha1.Switch, opts metav1.UpdateOptions) (*switchv1alpha1.Switch, error) {
	result := &switchv1alpha1.Switch{}
	err := c.restClient.Put().
		Namespace(c.ns).
		Resource(CSwitchesResourceType).
		Name(obj.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(obj).
		Do(ctx).
		Into(result)

	return result, err
}

func (c *switchClient) UpdateStatus(ctx context.Context, obj *switchv1alpha1.Switch, opts metav1.UpdateOptions) (*switchv1alpha1.Switch, error) {
	result := &switchv1alpha1.Switch{}
	err := c.restClient.Put().
		Namespace(c.ns).
		Resource(CSwitchesResourceType).
		Name(obj.Name).
		SubResource(CStatusSubresource).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(obj).
		Do(ctx).
		Into(result)

	return result, err
}

func (c *switchClient) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.restClient.Delete().
		Namespace(c.ns).
		Resource(CSwitchesResourceType).
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

func (c *switchClient) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}

	return c.restClient.Delete().
		Namespace(c.ns).
		Resource(CSwitchesResourceType).
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

func (c *switchClient) Patch(
	ctx context.Context,
	name string,
	pt types.PatchType,
	data []byte,
	opts metav1.PatchOptions,
	subresources ...string) (*switchv1alpha1.Switch, error) {
	result := &switchv1alpha1.Switch{}
	err := c.restClient.Patch(pt).
		Namespace(c.ns).
		Resource(CSwitchesResourceType).
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)

	return result, err
}
