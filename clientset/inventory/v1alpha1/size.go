/*
Copyright (c) 2022 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved

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

	metalv1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"
)

const (
	CSizesResourceType = "sizes"
)

type SizeInterface interface {
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*metalv1alpha4.Size, error)
	List(ctx context.Context, opts metav1.ListOptions) (*metalv1alpha4.SizeList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Create(ctx context.Context, size *metalv1alpha4.Size, opts metav1.CreateOptions) (*metalv1alpha4.Size, error)
	Update(ctx context.Context, size *metalv1alpha4.Size, opts metav1.UpdateOptions) (*metalv1alpha4.Size, error)
	UpdateStatus(ctx context.Context, size *metalv1alpha4.Size, opts metav1.UpdateOptions) (*metalv1alpha4.Size, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (*metalv1alpha4.Size, error)
}

type sizeClient struct {
	restClient rest.Interface
	ns         string
}

func (c *sizeClient) Get(ctx context.Context, name string, opts metav1.GetOptions) (*metalv1alpha4.Size, error) {
	result := &metalv1alpha4.Size{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource(CSizesResourceType).
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(result)

	return result, err
}

func (c *sizeClient) List(ctx context.Context, opts metav1.ListOptions) (*metalv1alpha4.SizeList, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result := &metalv1alpha4.SizeList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource(CSizesResourceType).
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)

	return result, err
}

func (c *sizeClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	watcher, err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource(CSizesResourceType).
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)

	return watcher, err
}

func (c *sizeClient) Create(ctx context.Context, size *metalv1alpha4.Size, opts metav1.CreateOptions) (*metalv1alpha4.Size, error) {
	result := &metalv1alpha4.Size{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource(CSizesResourceType).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(size).
		Do(ctx).
		Into(result)

	return result, err
}

func (c *sizeClient) Update(ctx context.Context, size *metalv1alpha4.Size, opts metav1.UpdateOptions) (*metalv1alpha4.Size, error) {
	result := &metalv1alpha4.Size{}
	err := c.restClient.Put().
		Namespace(c.ns).
		Resource(CSizesResourceType).
		Name(size.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(size).
		Do(ctx).
		Into(result)

	return result, err
}

func (c *sizeClient) UpdateStatus(ctx context.Context, size *metalv1alpha4.Size, opts metav1.UpdateOptions) (*metalv1alpha4.Size, error) {
	result := &metalv1alpha4.Size{}
	err := c.restClient.Put().
		Namespace(c.ns).
		Resource(CSizesResourceType).
		Name(size.Name).
		SubResource(CStatusSubresource).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(size).
		Do(ctx).
		Into(result)

	return result, err
}

func (c *sizeClient) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.restClient.Delete().
		Namespace(c.ns).
		Resource(CSizesResourceType).
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

func (c *sizeClient) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}

	return c.restClient.Delete().
		Namespace(c.ns).
		Resource(CSizesResourceType).
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

func (c *sizeClient) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (*metalv1alpha4.Size, error) {
	result := &metalv1alpha4.Size{}
	err := c.restClient.Patch(pt).
		Namespace(c.ns).
		Resource(CSizesResourceType).
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)

	return result, err
}
