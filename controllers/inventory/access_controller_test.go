// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrlruntime "sigs.k8s.io/controller-runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"

	controllers "github.com/ironcore-dev/metal/controllers/inventory"
	"github.com/ironcore-dev/metal/controllers/inventory/fake"
)

func accessReconciler(a *assert.Assertions, objects ...ctrlclient.Object) *controllers.AccessReconciler {
	fakeClient, err := fake.NewFakeWithObjects(objects...)
	a.NotEmpty(fakeClient)
	a.Nil(err)
	return &controllers.AccessReconciler{
		Client:             fakeClient,
		BootstrapAPIServer: "",
		Log:                ctrlruntime.Log.WithName("controllers").WithName("TestAccess"),
	}
}

func TestReconciler(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	inventoryName, inventoryNamespace := "b127954c-3475-22b2-b91c-62d8b4f8cd3f", "default"
	inventory := fake.InventoryObject(inventoryName, inventoryNamespace)

	fakeIP := "127.0.0.1"
	var fakeIPPort int32 = 6443
	endpoint := fake.IPObjectEndpoint(fakeIP, fakeIPPort)

	reconciler := accessReconciler(a, inventory, endpoint)
	request := reconcilerRequest(inventoryName, inventoryNamespace)

	result, err := reconciler.Reconcile(context.Background(), request)
	a.Nil(err, "musr reconcile without error")
	a.True(result.IsZero())
}
func TestInventoryServiceAccount(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	inventoryName, inventoryNamespace := "b127954c-3475-22b2-b91c-62d8b4f8cd3f", "default"
	inventory := fake.InventoryObject(inventoryName, inventoryNamespace)

	fakeIP := "127.0.0.1"
	var fakeIPPort int32 = 6443
	endpoint := fake.IPObjectEndpoint(fakeIP, fakeIPPort)

	reconciler := accessReconciler(a, inventory, endpoint)

	request := reconcilerRequest(inventoryName, inventoryNamespace)

	inventory, err := reconciler.ExtractInventory(context.Background(), request)
	a.NotEmpty(inventory)
	a.Nil(err, "must return inventory")

	inventoryServiceAccount := reconciler.InventoryServiceAccount(inventory)

	a.Equal(inventoryServiceAccount.Name, controllers.InventoryServiceAccountPrefix+inventory.Name)
	a.Equal(inventoryServiceAccount.Namespace, inventoryNamespace)
	a.NotEmpty(inventoryServiceAccount)
}

func TestKubeconfigForServer(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	inventoryName, inventoryNamespace := "b127954c-3475-22b2-b91c-62d8b4f8cd3f", "default"
	inventoryTokenSecret := machineInventoryTokenSecret(inventoryName, inventoryNamespace)

	fakeIP := "127.0.0.1"
	var fakeIPPort int32 = 6443
	endpoint := fake.IPObjectEndpoint(fakeIP, fakeIPPort)

	reconciler := accessReconciler(a, endpoint)

	secret, err := reconciler.KubeconfigForServer(
		context.Background(),
		inventoryTokenSecret,
		"fake",
		"default",
	)
	a.Nil(err, "must create kubeconfig for server")
	a.NotEmpty(secret.Clusters[reconciler.ClusterName()], "cluster info must exist")
	a.NotEmpty(secret.AuthInfos["fake"], "auth info must exist")
}
func TestAPIServerEndpointV4(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	fakeIP := "127.0.0.1"
	var fakeIPPort int32 = 6443
	endpoint := fake.IPObjectEndpoint(fakeIP, fakeIPPort)

	serverAPI := fmt.Sprintf("https://%s:%d", fakeIP, fakeIPPort)

	reconciler := accessReconciler(a, endpoint)

	apiServerEndpoint, err := reconciler.APIServerEndpoint(context.Background())
	a.Nil(err, "must return api endpoint")
	a.Equal(serverAPI, apiServerEndpoint)
}

func TestAPIServerEndpointV6(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	fakeIPv6 := "fd94:105f:be82:47fa::"
	var fakeIPv6Port int32 = 6443
	endpointsV6 := fake.RawEndpoints("fd94:105f:be82:47fa::", fakeIPv6Port)

	serverAPI := fmt.Sprintf("https://[%s]:%d", fakeIPv6, fakeIPv6Port)

	kubernetesAPIAddress := controllers.KubernetesAPIAddressFromEndpoint(endpointsV6.Subsets[0])
	a.NotEqual("", kubernetesAPIAddress)

	kubernetesAPIPort := controllers.KubernetesAPIPortFromEndpoint(endpointsV6.Subsets[0].Ports)
	a.NotEqual(0, kubernetesAPIPort)
	apiEndpoint := controllers.APIEndpoint{
		Address: kubernetesAPIAddress,
		Port:    kubernetesAPIPort,
	}
	apiServerEndpoint := controllers.SanitizeAPIServerEndpointString(apiEndpoint)
	a.Equal(serverAPI, apiServerEndpoint)
}

func TestKubernetesAPIAddressFromEndpoint(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	fakeIP := "127.0.0.1"
	var fakeIPPort int32 = 6443

	endpoints := fake.RawEndpoints(fakeIP, fakeIPPort)
	a.NotEmpty(endpoints.Subsets, "endpoint must not be empty")

	result := controllers.KubernetesAPIAddressFromEndpoint(endpoints.Subsets[0])
	a.Equal(fakeIP, result)
}
func TestKubernetesAPIPortFromEndpoint(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	fakeIP := "127.0.0.1"
	var fakeIPPort int32 = 6443

	endpoints := fake.RawEndpoints(fakeIP, fakeIPPort)
	a.NotEmpty(endpoints.Subsets, "endpoint must not be empty")

	result := controllers.KubernetesAPIPortFromEndpoint(endpoints.Subsets[0].Ports)
	a.Equal(fakeIPPort, result)
}

func reconcilerRequest(name, namespace string) ctrlruntime.Request {
	return ctrlruntime.Request{
		NamespacedName: types.NamespacedName{
			Namespace: namespace,
			Name:      name,
		},
	}
}
func machineInventoryTokenSecret(name, namespace string) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: map[string][]byte{
			"ca.crt": {1, 2, 3},
			"token":  {1, 2, 3},
		},
	}
}
