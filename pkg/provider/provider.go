package provider

import (
	"context"

	"k8s.io/apimachinery/pkg/types"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func GetObject[T ctrlclient.Object](ctx context.Context,
	name, namespace string, c ctrlclient.Client, obj T) error {
	if err := c.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, obj); err != nil {
		return err
	}
	return nil
}
