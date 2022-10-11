// nolint
package v1alpha1

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	"github.com/onmetal/metal-api/apis/inventory/v1alpha1"
)

const (
	CInventoriesResourceType = "inventories"
)

type InventoryInterface interface {
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1alpha1.Inventory, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1alpha1.InventoryList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Create(ctx context.Context, inventory *v1alpha1.Inventory, opts metav1.CreateOptions) (*v1alpha1.Inventory, error)
	Update(ctx context.Context, inventory *v1alpha1.Inventory, opts metav1.UpdateOptions) (*v1alpha1.Inventory, error)
	UpdateStatus(ctx context.Context, inventory *v1alpha1.Inventory, opts metav1.UpdateOptions) (*v1alpha1.Inventory, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (*v1alpha1.Inventory, error)
}

type inventoryClient struct {
	restClient rest.Interface
	ns         string
}

func (c *inventoryClient) Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1alpha1.Inventory, error) {
	result := &v1alpha1.Inventory{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource(CInventoriesResourceType).
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(result)

	return result, err
}

func (c *inventoryClient) List(ctx context.Context, opts metav1.ListOptions) (*v1alpha1.InventoryList, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result := &v1alpha1.InventoryList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource(CInventoriesResourceType).
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)

	return result, err
}

func (c *inventoryClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	watcher, err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource(CInventoriesResourceType).
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)

	return watcher, err
}

func (c *inventoryClient) Create(ctx context.Context, inventory *v1alpha1.Inventory, opts metav1.CreateOptions) (*v1alpha1.Inventory, error) {
	result := &v1alpha1.Inventory{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource(CInventoriesResourceType).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(inventory).
		Do(ctx).
		Into(result)

	return result, err
}

func (c *inventoryClient) Update(ctx context.Context, inventory *v1alpha1.Inventory, opts metav1.UpdateOptions) (*v1alpha1.Inventory, error) {
	result := &v1alpha1.Inventory{}
	err := c.restClient.Put().
		Namespace(c.ns).
		Resource(CInventoriesResourceType).
		Name(inventory.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(inventory).
		Do(ctx).
		Into(result)

	return result, err
}

func (c *inventoryClient) UpdateStatus(ctx context.Context, inventory *v1alpha1.Inventory, opts metav1.UpdateOptions) (*v1alpha1.Inventory, error) {
	result := &v1alpha1.Inventory{}
	err := c.restClient.Put().
		Namespace(c.ns).
		Resource(CInventoriesResourceType).
		Name(inventory.Name).
		SubResource(CStatusSubresource).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(inventory).
		Do(ctx).
		Into(result)

	return result, err
}

func (c *inventoryClient) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.restClient.Delete().
		Namespace(c.ns).
		Resource(CInventoriesResourceType).
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

func (c *inventoryClient) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}

	return c.restClient.Delete().
		Namespace(c.ns).
		Resource(CInventoriesResourceType).
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

func (c *inventoryClient) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (*v1alpha1.Inventory, error) {
	result := &v1alpha1.Inventory{}
	err := c.restClient.Patch(pt).
		Namespace(c.ns).
		Resource(CInventoriesResourceType).
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)

	return result, err
}
