/*
Copyright (c) 2021 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved

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

package inventory

import (
	"context"
	"net"
	"strings"

	"github.com/go-logr/logr"
	inventoriesv1alpha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	machinev1alpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	switchv1alpha1 "github.com/onmetal/metal-api/apis/switches/v1alpha1"
	machinerr "github.com/onmetal/metal-api/pkg/errors"
	"github.com/onmetal/metal-api/pkg/machine"
	apierr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	subnetSize = 30
)

const defaultNumberOfInterfaces = 2

const (
	machineSizeName = "machine"
)

const (
	onePort = 1 + iota
	twoPorts
)

type Inventory struct {
	ctrlclient.Client
	*inventoriesv1alpha1.Inventory

	ctx      context.Context
	log      logr.Logger
	Machiner machine.Machiner
}

func New(ctx context.Context, c ctrlclient.Client, l logr.Logger, r record.EventRecorder, req ctrl.Request) (*Inventory, error) {
	i, err := getInventory(ctx, c, req.Name, req.Namespace)
	if err != nil {
		return nil, err
	}

	if _, ok := i.Labels[inventoriesv1alpha1.GetSizeMatchLabel(machineSizeName)]; !ok {
		return nil, machinerr.NotAMachine()
	}

	mm := machine.New(ctx, c, l, r)
	return &Inventory{
		Client:    c,
		Inventory: i,
		Machiner:  mm,
		ctx:       ctx,
		log:       l,
	}, nil
}

func getInventory(ctx context.Context, c ctrlclient.Client,
	name, namespace string) (*inventoriesv1alpha1.Inventory, error) {
	obj := &inventoriesv1alpha1.Inventory{}
	if err := c.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, obj); err != nil {
		return nil, err
	}
	return obj, nil
}

func (i *Inventory) UpdateMachine(machineObj *machinev1alpha2.Machine) error {
	machineObj.Spec.InventoryRequested = false
	for key := range i.Labels {
		machineObj.Labels[key] = i.Labels[key]
	}
	if err := i.Machiner.UpdateSpec(machineObj); err != nil {
		return err
	}

	i.updateResourceReference(machineObj)

	i.updateMachineInterfaces(machineObj)

	i.updateIdentity(machineObj)

	i.copySizeLabelsToMachine(machineObj)

	return i.Machiner.UpdateStatus(machineObj)
}

func (i *Inventory) updateMachineInterfaces(m *machinev1alpha2.Machine) {
	m.Status.Interfaces = i.getInterfaces(m)
	m.Status.Network.Ports = len(m.Status.Interfaces)
	m.Status.Network.Redundancy = i.getNetworkRedundancy(m)
	m.Status.Network.UnknownPorts = i.getUnknownPortsCount(m)
}

func (i *Inventory) getNetworkRedundancy(m *machinev1alpha2.Machine) string {
	switch {
	case len(m.Status.Interfaces) == onePort:
		return machinev1alpha2.InterfaceRedundancySingle
	case len(m.Status.Interfaces) >= twoPorts:
		if m.Status.Interfaces[0].Peer.LLDPChassisID != m.Status.Interfaces[1].Peer.LLDPChassisID {
			return machinev1alpha2.InterfaceRedundancyHighAvailability
		}
		return machinev1alpha2.InterfaceRedundancySingle
	default:
		return machinev1alpha2.InterfaceRedundancyNone
	}
}

func (i *Inventory) getUnknownPortsCount(m *machinev1alpha2.Machine) int {
	var count int
	for machinePort := range m.Status.Interfaces {
		if !(m.Status.Interfaces[machinePort].Unknown) {
			continue
		}
		count++
	}
	return count
}

func (i *Inventory) updateIdentity(m *machinev1alpha2.Machine) {
	m.Spec.Identity.SKU = i.Spec.System.ProductSKU
	m.Spec.Identity.SerialNumber = i.Spec.System.SerialNumber
}

func (i *Inventory) copySizeLabelsToMachine(m *machinev1alpha2.Machine) {
	for key, value := range i.Labels {
		if !strings.Contains(key, inventoriesv1alpha1.CLabelPrefix) {
			continue
		}
		m.Labels[key] = value
	}
}

func (i *Inventory) updateResourceReference(m *machinev1alpha2.Machine) {
	if !m.Status.Inventory.Exist || m.Status.Inventory.Reference == nil {
		m.Status.Inventory = i.prepareRefenceSpec()
	}
}

func (i *Inventory) prepareRefenceSpec() machinev1alpha2.ObjectReference {
	return machinev1alpha2.ObjectReference{
		Exist: true,
		Reference: &machinev1alpha2.ResourceReference{
			Kind: i.Kind, APIVersion: i.APIVersion,
			Name: i.Name, Namespace: i.Namespace},
	}
}

func (i *Inventory) getInterfaces(m *machinev1alpha2.Machine) []machinev1alpha2.Interface {
	interfaces := make([]machinev1alpha2.Interface, 0, defaultNumberOfInterfaces)
	nicsSpec := i.Spec.NICs
	for nic := range nicsSpec {
		if len(nicsSpec[nic].LLDPs) == 0 {
			interfaces = i.baseConnectionInfo(&nicsSpec[nic], interfaces, m)
			continue
		}
		label := map[string]string{
			switchv1alpha1.LabelChassisId: strings.ReplaceAll(nicsSpec[nic].LLDPs[0].ChassisID, ":", "-"),
		}
		s, err := i.getSwitchByLabel(i.ctx, i.Client, label)
		if apierr.IsNotFound(err) || machinerr.IsNotFound(err) {
			interfaces = i.baseConnectionInfo(&nicsSpec[nic], interfaces, m)
			continue
		}

		switchInterface, ok := s.Status.Interfaces[nicsSpec[nic].LLDPs[0].PortDescription]
		if !ok {
			interfaces = i.baseConnectionInfo(&nicsSpec[nic], interfaces, m)
			continue
		}
		interfaces = i.connectionInfoEnrichment(s.ObjectMeta, &nicsSpec[nic], interfaces, s.Name, switchInterface, m)
	}
	return interfaces
}

func (i *Inventory) getSwitchByLabel(ctx context.Context, c ctrlclient.Client,
	label map[string]string) (*switchv1alpha1.Switch, error) {
	obj := &switchv1alpha1.SwitchList{}
	filter := &ctrlclient.ListOptions{
		LabelSelector: ctrlclient.MatchingLabelsSelector{Selector: labels.SelectorFromSet(label)},
	}
	if err := c.List(ctx, obj, filter); err != nil {
		return nil, err
	}
	if len(obj.Items) == 0 {
		return nil, machinerr.NotFound("switch")
	}
	return &obj.Items[0], nil
}

func (i *Inventory) connectionInfoEnrichment(sw metav1.ObjectMeta, nicsSpec *inventoriesv1alpha1.NICSpec,
	interfaces []machinev1alpha2.Interface,
	switchUUID string, switchInterface *switchv1alpha1.InterfaceSpec,
	m *machinev1alpha2.Machine) []machinev1alpha2.Interface {
	return append(interfaces, machinev1alpha2.Interface{
		Name:            nicsSpec.Name,
		Lanes:           switchInterface.Lanes,
		IPv4:            &machinev1alpha2.IPAddressSpec{Address: i.getAddress(switchInterface.IPv4.Address)},
		IPv6:            &machinev1alpha2.IPAddressSpec{Address: i.getAddress(switchInterface.IPv6.Address)},
		Moved:           i.getMovedInterface(nicsSpec, m),
		Unknown:         false,
		SwitchReference: &machinev1alpha2.ResourceReference{Kind: "Switch", Namespace: sw.Namespace, Name: sw.Name},
		Peer: &machinev1alpha2.Peer{
			LLDPSystemName:      switchUUID,
			LLDPChassisID:       nicsSpec.LLDPs[0].ChassisID,
			LLDPPortID:          nicsSpec.LLDPs[0].PortID,
			LLDPPortDescription: nicsSpec.LLDPs[0].PortDescription,
		},
	})
}

func (i *Inventory) getAddress(switchIP string) string {
	for s := 0; s < len(switchIP); s++ {
		switch switchIP[s] {
		case '.':
			ip, ipNet, err := net.ParseCIDR(switchIP)
			if err != nil {
				i.log.Info("can't parse ip address", "error", err)
				return ""
			}
			if size, _ := ipNet.Mask.Size(); size < subnetSize {
				i.log.Info("subnet mask less than minimal subnet size", "minimal size", subnetSize,
					"current size", size)
				return ""
			}
			ip = ip.To4()
			ip[3]++
			machineAddr := net.IPNet{
				IP:   ip,
				Mask: ipNet.Mask,
			}
			return machineAddr.String()
		case ':':
			ip, ipNet, err := net.ParseCIDR(switchIP)
			if err != nil {
				i.log.Info("can't parse ip address", "error", err)
				return ""
			}
			ip = ip.To16()
			ip[15]++
			machineAddr := net.IPNet{
				IP:   ip,
				Mask: ipNet.Mask,
			}
			return machineAddr.String()
		}
	}
	return ""
}

func (i *Inventory) baseConnectionInfo(nicsSpec *inventoriesv1alpha1.NICSpec,
	interfaces []machinev1alpha2.Interface, m *machinev1alpha2.Machine) []machinev1alpha2.Interface {
	if len(nicsSpec.LLDPs) != 1 {
		i.log.Info("incorrect lldp neighbor count",
			"inventory", i.Name,
			"interface", nicsSpec.Name,
			"count", len(nicsSpec.LLDPs))
		return append(interfaces, machinev1alpha2.Interface{
			Name:    nicsSpec.Name,
			Unknown: true,
		})
	}
	return append(interfaces, machinev1alpha2.Interface{
		Name:    nicsSpec.Name,
		Unknown: false,
		Moved:   i.getMovedInterface(nicsSpec, m),
		Peer: &machinev1alpha2.Peer{
			LLDPSystemName:      nicsSpec.LLDPs[0].SystemName,
			LLDPChassisID:       nicsSpec.LLDPs[0].ChassisID,
			LLDPPortID:          nicsSpec.LLDPs[0].PortID,
			LLDPPortDescription: nicsSpec.LLDPs[0].PortDescription,
		},
	})
}

func (i *Inventory) getMovedInterface(newInterfaceState *inventoriesv1alpha1.NICSpec,
	m *machinev1alpha2.Machine) bool {
	for mi := range m.Status.Interfaces {
		if m.Status.Interfaces[mi].Name != newInterfaceState.Name {
			continue
		}
		if m.Status.Interfaces[mi].Peer.LLDPChassisID != newInterfaceState.LLDPs[0].ChassisID {
			return true
		}
	}
	return false
}

// func (i *Inventory) updateInventoryStatus(status bool) error {
// 	// i.machine.Status.Inventory.Exist = status
// 	return i.Client.Status().Update(i.ctx, i.machine)
// }

// func updateRequired(m *machinev1alpha2.Machine) (switchPort, switchID string, required bool) {
// 	if len(m.Status.Interfaces) == 0 {
// 		return "", "", false
// 	}
// 	for i := range m.Status.Interfaces {
// 		if m.Status.Interfaces[i].Moved {
// 			return m.Status.Interfaces[i].Peer.LLDPPortDescription, m.Status.Interfaces[i].Peer.LLDPChassisID, true
// 		}
// 	}
// 	return m.Status.Interfaces[0].Peer.LLDPPortDescription,
// 		m.Status.Interfaces[0].Peer.LLDPChassisID
// }
