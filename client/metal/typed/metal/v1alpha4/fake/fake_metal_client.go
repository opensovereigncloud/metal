// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha4 "github.com/ironcore-dev/metal/client/metal/typed/metal/v1alpha4"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeMetalV1alpha4 struct {
	*testing.Fake
}

func (c *FakeMetalV1alpha4) Aggregates(namespace string) v1alpha4.AggregateInterface {
	return &FakeAggregates{c, namespace}
}

func (c *FakeMetalV1alpha4) Benchmarks(namespace string) v1alpha4.BenchmarkInterface {
	return &FakeBenchmarks{c, namespace}
}

func (c *FakeMetalV1alpha4) Inventories(namespace string) v1alpha4.InventoryInterface {
	return &FakeInventories{c, namespace}
}

func (c *FakeMetalV1alpha4) Machines(namespace string) v1alpha4.MachineInterface {
	return &FakeMachines{c, namespace}
}

func (c *FakeMetalV1alpha4) NetworkSwitches(namespace string) v1alpha4.NetworkSwitchInterface {
	return &FakeNetworkSwitches{c, namespace}
}

func (c *FakeMetalV1alpha4) Sizes(namespace string) v1alpha4.SizeInterface {
	return &FakeSizes{c, namespace}
}

func (c *FakeMetalV1alpha4) SwitchConfigs(namespace string) v1alpha4.SwitchConfigInterface {
	return &FakeSwitchConfigs{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeMetalV1alpha4) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
