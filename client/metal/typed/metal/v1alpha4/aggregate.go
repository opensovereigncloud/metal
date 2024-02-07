// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

// Code generated by client-gen. DO NOT EDIT.

package v1alpha4

import (
	"context"
	json "encoding/json"
	"fmt"
	"time"

	v1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"
	metalv1alpha4 "github.com/ironcore-dev/metal/client/applyconfiguration/metal/v1alpha4"
	scheme "github.com/ironcore-dev/metal/client/metal/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// AggregatesGetter has a method to return a AggregateInterface.
// A group's client should implement this interface.
type AggregatesGetter interface {
	Aggregates(namespace string) AggregateInterface
}

// AggregateInterface has methods to work with Aggregate resources.
type AggregateInterface interface {
	Create(ctx context.Context, aggregate *v1alpha4.Aggregate, opts v1.CreateOptions) (*v1alpha4.Aggregate, error)
	Update(ctx context.Context, aggregate *v1alpha4.Aggregate, opts v1.UpdateOptions) (*v1alpha4.Aggregate, error)
	UpdateStatus(ctx context.Context, aggregate *v1alpha4.Aggregate, opts v1.UpdateOptions) (*v1alpha4.Aggregate, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha4.Aggregate, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha4.AggregateList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha4.Aggregate, err error)
	Apply(ctx context.Context, aggregate *metalv1alpha4.AggregateApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha4.Aggregate, err error)
	ApplyStatus(ctx context.Context, aggregate *metalv1alpha4.AggregateApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha4.Aggregate, err error)
	AggregateExpansion
}

// aggregates implements AggregateInterface
type aggregates struct {
	client rest.Interface
	ns     string
}

// newAggregates returns a Aggregates
func newAggregates(c *MetalV1alpha4Client, namespace string) *aggregates {
	return &aggregates{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the aggregate, and returns the corresponding aggregate object, and an error if there is any.
func (c *aggregates) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha4.Aggregate, err error) {
	result = &v1alpha4.Aggregate{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("aggregates").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Aggregates that match those selectors.
func (c *aggregates) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha4.AggregateList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha4.AggregateList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("aggregates").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested aggregates.
func (c *aggregates) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("aggregates").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a aggregate and creates it.  Returns the server's representation of the aggregate, and an error, if there is any.
func (c *aggregates) Create(ctx context.Context, aggregate *v1alpha4.Aggregate, opts v1.CreateOptions) (result *v1alpha4.Aggregate, err error) {
	result = &v1alpha4.Aggregate{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("aggregates").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(aggregate).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a aggregate and updates it. Returns the server's representation of the aggregate, and an error, if there is any.
func (c *aggregates) Update(ctx context.Context, aggregate *v1alpha4.Aggregate, opts v1.UpdateOptions) (result *v1alpha4.Aggregate, err error) {
	result = &v1alpha4.Aggregate{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("aggregates").
		Name(aggregate.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(aggregate).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *aggregates) UpdateStatus(ctx context.Context, aggregate *v1alpha4.Aggregate, opts v1.UpdateOptions) (result *v1alpha4.Aggregate, err error) {
	result = &v1alpha4.Aggregate{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("aggregates").
		Name(aggregate.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(aggregate).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the aggregate and deletes it. Returns an error if one occurs.
func (c *aggregates) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("aggregates").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *aggregates) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("aggregates").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched aggregate.
func (c *aggregates) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha4.Aggregate, err error) {
	result = &v1alpha4.Aggregate{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("aggregates").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// Apply takes the given apply declarative configuration, applies it and returns the applied aggregate.
func (c *aggregates) Apply(ctx context.Context, aggregate *metalv1alpha4.AggregateApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha4.Aggregate, err error) {
	if aggregate == nil {
		return nil, fmt.Errorf("aggregate provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(aggregate)
	if err != nil {
		return nil, err
	}
	name := aggregate.Name
	if name == nil {
		return nil, fmt.Errorf("aggregate.Name must be provided to Apply")
	}
	result = &v1alpha4.Aggregate{}
	err = c.client.Patch(types.ApplyPatchType).
		Namespace(c.ns).
		Resource("aggregates").
		Name(*name).
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *aggregates) ApplyStatus(ctx context.Context, aggregate *metalv1alpha4.AggregateApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha4.Aggregate, err error) {
	if aggregate == nil {
		return nil, fmt.Errorf("aggregate provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(aggregate)
	if err != nil {
		return nil, err
	}

	name := aggregate.Name
	if name == nil {
		return nil, fmt.Errorf("aggregate.Name must be provided to Apply")
	}

	result = &v1alpha4.Aggregate{}
	err = c.client.Patch(types.ApplyPatchType).
		Namespace(c.ns).
		Resource("aggregates").
		Name(*name).
		SubResource("status").
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
