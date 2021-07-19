package v1alpha1

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	"github.com/onmetal/k8s-inventory/api/v1alpha1"
)

const (
	CAggregatesResourceType = "aggregates"
)

type AggregateInterface interface {
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1alpha1.Aggregate, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1alpha1.AggregateList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Create(ctx context.Context, aggregate *v1alpha1.Aggregate, opts metav1.CreateOptions) (*v1alpha1.Aggregate, error)
	Update(ctx context.Context, aggregate *v1alpha1.Aggregate, opts metav1.UpdateOptions) (*v1alpha1.Aggregate, error)
	UpdateStatus(ctx context.Context, aggregate *v1alpha1.Aggregate, opts metav1.UpdateOptions) (*v1alpha1.Aggregate, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (*v1alpha1.Aggregate, error)
}

type aggregateClient struct {
	restClient rest.Interface
	ns         string
}

func (c *aggregateClient) Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1alpha1.Aggregate, error) {
	result := &v1alpha1.Aggregate{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource(CAggregatesResourceType).
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(result)

	return result, err
}

func (c *aggregateClient) List(ctx context.Context, opts metav1.ListOptions) (*v1alpha1.AggregateList, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result := &v1alpha1.AggregateList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource(CAggregatesResourceType).
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)

	return result, err
}

func (c *aggregateClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	watcher, err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource(CAggregatesResourceType).
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)

	return watcher, err
}

func (c *aggregateClient) Create(ctx context.Context, aggregate *v1alpha1.Aggregate, opts metav1.CreateOptions) (*v1alpha1.Aggregate, error) {
	result := &v1alpha1.Aggregate{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource(CAggregatesResourceType).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(aggregate).
		Do(ctx).
		Into(result)

	return result, err
}

func (c *aggregateClient) Update(ctx context.Context, aggregate *v1alpha1.Aggregate, opts metav1.UpdateOptions) (*v1alpha1.Aggregate, error) {
	result := &v1alpha1.Aggregate{}
	err := c.restClient.Put().
		Namespace(c.ns).
		Resource(CAggregatesResourceType).
		Name(aggregate.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(aggregate).
		Do(ctx).
		Into(result)

	return result, err
}

func (c *aggregateClient) UpdateStatus(ctx context.Context, aggregate *v1alpha1.Aggregate, opts metav1.UpdateOptions) (*v1alpha1.Aggregate, error) {
	result := &v1alpha1.Aggregate{}
	err := c.restClient.Put().
		Namespace(c.ns).
		Resource(CAggregatesResourceType).
		Name(aggregate.Name).
		SubResource(CStatusSubresource).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(aggregate).
		Do(ctx).
		Into(result)

	return result, err
}

func (c *aggregateClient) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.restClient.Delete().
		Namespace(c.ns).
		Resource(CAggregatesResourceType).
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

func (c *aggregateClient) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}

	return c.restClient.Delete().
		Namespace(c.ns).
		Resource(CAggregatesResourceType).
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

func (c *aggregateClient) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (*v1alpha1.Aggregate, error) {
	result := &v1alpha1.Aggregate{}
	err := c.restClient.Patch(pt).
		Namespace(c.ns).
		Resource(CAggregatesResourceType).
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)

	return result, err
}
