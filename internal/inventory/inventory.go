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
	machinev1alpha1 "github.com/onmetal/metal-api/apis/machine/v1alpha1"
	switchv1alpha1 "github.com/onmetal/metal-api/apis/switches/v1alpha1"
	machinerr "github.com/onmetal/metal-api/internal/errors"
	"github.com/onmetal/metal-api/internal/provider"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	machineLeasedByLabel      = "machine.onmetal.de/leased-by"
	defaultNumberOfInterfaces = 2
	healthy                   = true
	subnetSize                = 30
	machineSizeName           = "machine"
)

type Inventory struct {
	ctrlclient.Client
	*inventoriesv1alpha1.Inventory

	ctx     context.Context
	log     logr.Logger
	machine *machinev1alpha1.Machine
}

func New(ctx context.Context, c ctrlclient.Client, l logr.Logger, req ctrl.Request) (*Inventory, error) {
	invobj, err := provider.Get(ctx, c, req.Name, req.Namespace, provider.Inventory)
	if err != nil {
		return nil, err
	}
	inventory, ok := invobj.(*inventoriesv1alpha1.Inventory)
	if !ok {
		return &Inventory{}, machinerr.CastType()
	}
	labels := inventory.GetLabels()
	if len(labels) == 0 {
		return nil, machinerr.NotLabeled()
	}
	machineSizeLabel := inventoriesv1alpha1.GetSizeMatchLabel(machineSizeName)
	if _, ok := labels[machineSizeLabel]; !ok {
		return nil, machinerr.NotAMachine()
	}
	return &Inventory{
		Client:    c,
		Inventory: inventory,
		ctx:       ctx,
		log:       l,
	}, nil
}

func (i *Inventory) Update() error {
	machineobj, err := provider.Get(i.ctx, i.Client, i.Spec.System.ID, i.Namespace, provider.Machine)
	if err != nil {
		return err
	}
	machine, ok := machineobj.(*machinev1alpha1.Machine)
	if !ok {
		return machinerr.CastType()
	}
	i.machine = machine
	if updErr := i.updateMachineInterfaces(); updErr != nil {
		return updErr
	}
	if interfaceName, chassisID, required := updateRequired(i.machine); required {
		i.machine.Spec.Location = i.updateLocation(interfaceName, chassisID)
	}
	i.updateIdentity()

	if i.isLeasedLabels() {
		i.machine.Spec.Action.PowerState = machinev1alpha1.MachinPowerStateON
	}

	i.machine.Spec.InventoryRequested = false
	return i.Client.Update(i.ctx, i.machine)
}

func (i *Inventory) updateMachineInterfaces() error {
	i.machine.Status.Interfaces = i.getInterfaces()
	return i.updateInventoryStatus(healthy)
}

func (i *Inventory) updateLocation(interfaceName, chassisID string) machinev1alpha1.Location {
	label := map[string]string{switchv1alpha1.LabelChassisId: strings.ReplaceAll(chassisID, ":", "-")}
	obj, err := provider.GetByLabel(i.ctx, i.Client, label, provider.Switch)
	if err != nil {
		i.log.Info("can't find switch for machine", "chassisID", chassisID, "error", err)
		return i.machine.Spec.Location
	}
	s, ok := obj.(*switchv1alpha1.Switch)
	if !ok {
		return i.machine.Spec.Location
	}
	if s.Spec.Location == nil {
		return i.machine.Spec.Location
	}
	if _, ok := s.Status.Interfaces[interfaceName]; !ok {
		return i.machine.Spec.Location
	}
	return machinev1alpha1.Location{
		DataHall: s.Spec.Location.Room,
		Row:      s.Spec.Location.Row,
		Rack:     s.Spec.Location.Rack,
	}
}

func (i *Inventory) updateIdentity() {
	i.machine.Spec.Identity.SKU = i.Spec.System.ProductSKU
	i.machine.Spec.Identity.SerialNumber = i.Spec.System.SerialNumber
}

func (i *Inventory) isLeasedLabels() bool {
	_, ok := i.Inventory.Labels[machineLeasedByLabel]
	return ok
}

func (i *Inventory) getInterfaces() []machinev1alpha1.Interface {
	interfaces := make([]machinev1alpha1.Interface, 0, defaultNumberOfInterfaces)
	nicsSpec := i.Spec.NICs
	for nic := range nicsSpec {
		if len(nicsSpec[nic].LLDPs) == 0 {
			interfaces = i.baseConnectionInfo(&nicsSpec[nic], interfaces)
			continue
		}
		label := map[string]string{
			switchv1alpha1.LabelChassisId: strings.ReplaceAll(nicsSpec[nic].LLDPs[0].ChassisID, ":", "-"),
		}
		obj, err := provider.GetByLabel(i.ctx, i.Client, label, provider.Switch)
		if machinerr.IsNotFound(err) {
			interfaces = i.baseConnectionInfo(&nicsSpec[nic], interfaces)
			continue
		}
		s, ok := obj.(*switchv1alpha1.Switch)
		if !ok {
			continue
		}
		switchInterface, ok := s.Status.Interfaces[nicsSpec[nic].LLDPs[0].PortDescription]
		if !ok {
			interfaces = i.baseConnectionInfo(&nicsSpec[nic], interfaces)
			continue
		}
		interfaces = i.connectionInfoEnrichment(&nicsSpec[nic], interfaces, s.Name, switchInterface)
	}
	return interfaces
}

func (i *Inventory) connectionInfoEnrichment(nicsSpec *inventoriesv1alpha1.NICSpec,
	interfaces []machinev1alpha1.Interface,
	switchUUID string, switchInterface *switchv1alpha1.InterfaceSpec) []machinev1alpha1.Interface {
	return append(interfaces, machinev1alpha1.Interface{
		Name:                nicsSpec.Name,
		Lane:                switchInterface.Lanes,
		IPv4:                i.getAddress(switchInterface.IPv4.Address),
		IPv6:                i.getAddress(switchInterface.IPv6.Address),
		Moved:               i.getMovedInterface(nicsSpec),
		Unknown:             false,
		SwitchUUID:          switchUUID,
		LLDPSystemName:      nicsSpec.LLDPs[0].SystemName,
		LLDPChassisID:       nicsSpec.LLDPs[0].ChassisID,
		LLDPPortID:          nicsSpec.LLDPs[0].PortID,
		LLDPPortDescription: nicsSpec.LLDPs[0].PortDescription,
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
	interfaces []machinev1alpha1.Interface) []machinev1alpha1.Interface {
	if len(nicsSpec.LLDPs) != 1 {
		i.log.Info("incorrect lldp neighbor count",
			"inventory", i.Name,
			"interface", nicsSpec.Name,
			"count", len(nicsSpec.LLDPs))
		return append(interfaces, machinev1alpha1.Interface{
			Name:    nicsSpec.Name,
			Unknown: true,
		})
	}
	return append(interfaces, machinev1alpha1.Interface{
		Name:                nicsSpec.Name,
		Unknown:             false,
		Moved:               i.getMovedInterface(nicsSpec),
		LLDPSystemName:      nicsSpec.LLDPs[0].SystemName,
		LLDPChassisID:       nicsSpec.LLDPs[0].ChassisID,
		LLDPPortID:          nicsSpec.LLDPs[0].PortID,
		LLDPPortDescription: nicsSpec.LLDPs[0].PortDescription,
	})
}

func (i *Inventory) getMovedInterface(newInterfaceState *inventoriesv1alpha1.NICSpec) bool {
	for m := range i.machine.Status.Interfaces {
		if i.machine.Status.Interfaces[m].Name != newInterfaceState.Name {
			continue
		}
		if i.machine.Status.Interfaces[m].LLDPChassisID != newInterfaceState.LLDPs[0].ChassisID {
			return true
		}
	}
	return false
}

func (i *Inventory) updateInventoryStatus(status bool) error {
	i.machine.Status.Inventory = status
	return i.Client.Status().Update(i.ctx, i.machine)
}

func updateRequired(m *machinev1alpha1.Machine) (switchPort, switchID string, required bool) {
	if len(m.Status.Interfaces) == 0 {
		return "", "", false
	}
	for i := range m.Status.Interfaces {
		if m.Status.Interfaces[i].Moved {
			return m.Status.Interfaces[i].LLDPPortDescription, m.Status.Interfaces[i].LLDPChassisID, true
		}
	}
	return m.Status.Interfaces[0].LLDPPortDescription,
		m.Status.Interfaces[0].LLDPChassisID,
		!(m.Spec.Location.DataHall != "" && m.Spec.Location.Row != 0 && m.Spec.Location.Rack != 0)
}
