// /*
// Copyright (c) 2021 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// */

package repository

import (
	"context"
	"net"
	"strings"

	ipamv1alpha1 "github.com/onmetal/ipam/api/v1alpha1"
	inventoriesv1alpha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	machinev1alpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	switchv1beta1 "github.com/onmetal/metal-api/apis/switch/v1beta1"
	"github.com/onmetal/metal-api/internal/entity"
	machinerr "github.com/onmetal/metal-api/pkg/errors"
	oobv1 "github.com/onmetal/oob-operator/api/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	machineSizeName = "machine"
	switchSizeName  = "switch"
)

const (
	onePort = 1 + iota
	twoPorts
)

const subnetSize = 30

const defaultNumberOfInterfaces = 2

type DeviceOnboarding struct {
	client ctrlclient.Client
	device ctrlclient.Object
}

func NewOnboardingRepo(c ctrlclient.Client) *DeviceOnboarding {
	return &DeviceOnboarding{client: c}
}

func (o *DeviceOnboarding) Create(ctx context.Context) error {
	if err := o.client.Create(ctx, o.device); err != nil {
		if apierrors.IsAlreadyExists(err) {
			return nil
		}
		return err
	}
	return nil
}

func (o *DeviceOnboarding) InitializationStatus(ctx context.Context, e entity.Onboarding) entity.Initialization {
	inventory := &inventoriesv1alpha1.Inventory{
		ObjectMeta: metav1.ObjectMeta{
			Name:      e.RequestName,
			Namespace: e.RequestNamespace,
		},
	}

	if err := o.client.Get(ctx, ctrlclient.ObjectKeyFromObject(inventory), inventory); err != nil {
		return entity.Initialization{
			Require: false,
			Error:   err,
		}
	}

	machine := inventory.Labels[inventoriesv1alpha1.GetSizeMatchLabel(machineSizeName)]
	switches := inventory.Labels[inventoriesv1alpha1.GetSizeMatchLabel(switchSizeName)]

	switch {
	case machine != "":
		m := &machinev1alpha2.Machine{
			ObjectMeta: metav1.ObjectMeta{
				Name:      e.RequestName,
				Namespace: e.InitializationObjectNamespace,
			},
		}

		if err := o.client.Get(ctx, ctrlclient.ObjectKeyFromObject(m), m); err != nil {
			return entity.Initialization{Require: true, Error: nil}
		}

		return entity.Initialization{Require: false, Error: nil}
	case switches != "":
		s := &switchv1beta1.Switch{
			ObjectMeta: metav1.ObjectMeta{
				Name:      e.RequestName,
				Namespace: e.InitializationObjectNamespace,
			},
		}

		if err := o.client.Get(ctx, ctrlclient.ObjectKeyFromObject(s), s); err != nil {
			return entity.Initialization{Require: true, Error: nil}
		}

		return entity.Initialization{Require: false, Error: nil}
	default:
		return entity.Initialization{Require: false, Error: nil}
	}
}

func (o *DeviceOnboarding) Prepare(ctx context.Context, e entity.Onboarding) error {
	inventory := &inventoriesv1alpha1.Inventory{
		ObjectMeta: metav1.ObjectMeta{
			Name:      e.RequestName,
			Namespace: e.RequestNamespace,
		},
	}
	if err := o.client.Get(ctx, ctrlclient.ObjectKeyFromObject(inventory), inventory); err != nil {
		return err
	}

	if !o.IsSizeLabeled(inventory.Labels) {
		return machinerr.NotSizeLabeled()
	}

	if inventory.Spec.System.ID == "" {
		return machinerr.UUIDNotExist(inventory.Name)
	}

	if _, ok := inventory.Labels[inventoriesv1alpha1.GetSizeMatchLabel(machineSizeName)]; ok {
		o.device = o.prepareMachine(inventory.Spec.System.ID, e)
	}
	return nil
}

func (o *DeviceOnboarding) IsSizeLabeled(labels map[string]string) bool {
	machine := labels[inventoriesv1alpha1.GetSizeMatchLabel(machineSizeName)]
	switches := labels[inventoriesv1alpha1.GetSizeMatchLabel(switchSizeName)]
	return machine != "" || switches != ""
}

func (o *DeviceOnboarding) GatherData(ctx context.Context, e entity.Onboarding) error {
	inventory := &inventoriesv1alpha1.Inventory{
		ObjectMeta: metav1.ObjectMeta{
			Name:      e.RequestName,
			Namespace: e.RequestNamespace,
		},
	}
	if err := o.client.Get(ctx, ctrlclient.ObjectKeyFromObject(inventory), inventory); err != nil {
		return err
	}

	machine := &machinev1alpha2.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      inventory.Name,
			Namespace: e.InitializationObjectNamespace,
		},
	}
	if err := o.client.Get(ctx, ctrlclient.ObjectKeyFromObject(machine), machine); err != nil {
		return err
	}

	machine.Labels = copySizeLabelsToMachine(machine.Labels, inventory.Labels)
	machine.Spec = gatherMachineSpecData(inventory, machine.Spec)

	if err := o.client.Update(ctx, machine); err != nil {
		return err
	}
	machine.Status = o.gatherMachineStatusData(ctx, inventory, machine.Status)

	return o.client.Status().Update(ctx, machine)
}

func (o *DeviceOnboarding) prepareMachine(uuid string, e entity.Onboarding) *machinev1alpha2.Machine {
	return &machinev1alpha2.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      uuid,
			Namespace: e.InitializationObjectNamespace,
			Labels:    map[string]string{machinev1alpha2.UUIDLabel: uuid},
		},
		Spec: machinev1alpha2.MachineSpec{InventoryRequested: true},
	}
}

func gatherMachineSpecData(inventory *inventoriesv1alpha1.Inventory,
	machineSpec machinev1alpha2.MachineSpec) machinev1alpha2.MachineSpec {
	machineSpec.InventoryRequested = false
	machineSpec.Identity = updateIdentity(inventory)

	return machineSpec
}

func updateIdentity(i *inventoriesv1alpha1.Inventory) machinev1alpha2.Identity {
	return machinev1alpha2.Identity{
		SKU:          i.Spec.System.ProductSKU,
		SerialNumber: i.Spec.System.SerialNumber,
	}
}

func copySizeLabelsToMachine(machineLabels, inventoryLabels map[string]string) map[string]string {
	for key, value := range inventoryLabels {
		if !strings.Contains(key, inventoriesv1alpha1.CLabelPrefix) {
			continue
		}
		machineLabels[key] = value
	}
	return machineLabels
}

type metadata struct {
	typeMeta   metav1.TypeMeta
	objectMeta metav1.ObjectMeta
}

func (o *DeviceOnboarding) gatherMachineStatusData(ctx context.Context,
	inventory *inventoriesv1alpha1.Inventory,
	machineStatus machinev1alpha2.MachineStatus) machinev1alpha2.MachineStatus {
	inventoryRef := metadata{
		typeMeta:   inventory.TypeMeta,
		objectMeta: inventory.ObjectMeta,
	}
	if !machineStatus.Inventory.Exist || machineStatus.Inventory.Reference == nil {
		machineStatus.Inventory = updateInventoryResourceReference(inventoryRef)
	}
	if !machineStatus.OOB.Exist || machineStatus.OOB.Reference == nil {
		oobRef := o.findMachineOOBByLabel(ctx, inventory.Spec.System.ID)
		machineStatus.OOB = updateOOBResourceReference(oobRef)
	}

	if machineStatus.Reservation.Status == "" {
		machineStatus.Reservation.Status = "Available"
	}
	machineStatus.Interfaces = o.updateMachineInterfaces(ctx, inventory, machineStatus.Interfaces)

	machineStatus.Network = o.updateNetworkStatus(machineStatus.Interfaces)

	machineStatus.Health = updateHealthStatus(machineStatus)
	if machineStatus.Health == machinev1alpha2.MachineStateHealthy {
		machineStatus.Orphaned = false
	} else {
		machineStatus.Orphaned = true
	}

	return machineStatus
}

func updateInventoryResourceReference(i metadata) machinev1alpha2.ObjectReference {
	return machinev1alpha2.ObjectReference{
		Exist: true,
		Reference: &machinev1alpha2.ResourceReference{
			Kind: i.typeMeta.Kind, APIVersion: i.typeMeta.APIVersion,
			Name: i.objectMeta.Name, Namespace: i.objectMeta.Namespace},
	}
}

func updateOOBResourceReference(oob metadata) machinev1alpha2.ObjectReference {
	return machinev1alpha2.ObjectReference{
		Exist: true,
		Reference: &machinev1alpha2.ResourceReference{
			Kind: oob.typeMeta.Kind, APIVersion: oob.typeMeta.APIVersion,
			Name: oob.objectMeta.Name, Namespace: oob.objectMeta.Namespace},
	}
}

func (o *DeviceOnboarding) findMachineOOBByLabel(ctx context.Context,
	uuid string) metadata {
	obj := &oobv1.OOBList{}
	filter := &ctrlclient.ListOptions{
		LabelSelector: ctrlclient.MatchingLabelsSelector{
			Selector: labels.SelectorFromSet(map[string]string{machinev1alpha2.UUIDLabel: uuid})}}
	if err := o.client.List(ctx, obj, filter); err != nil {
		return metadata{}
	}
	if len(obj.Items) == 0 {
		return metadata{}
	}
	return metadata{
		typeMeta:   obj.Items[0].TypeMeta,
		objectMeta: obj.Items[0].ObjectMeta,
	}
}

func (o *DeviceOnboarding) updateMachineInterfaces(ctx context.Context,
	inventory *inventoriesv1alpha1.Inventory,
	machineInterfaces []machinev1alpha2.Interface,
) []machinev1alpha2.Interface {
	interfaces := make([]machinev1alpha2.Interface, 0, defaultNumberOfInterfaces)
	nicsSpec := inventory.Spec.NICs
	for nic := range nicsSpec {
		if len(nicsSpec[nic].LLDPs) == 0 {
			interfaces = baseConnectionInfo(&nicsSpec[nic], interfaces, machineInterfaces)
			continue
		}

		label := map[string]string{
			switchv1beta1.LabelChassisID: strings.ReplaceAll(nicsSpec[nic].LLDPs[0].ChassisID, ":", "-"),
		}
		s, err := o.getSwitchByLabel(ctx, label)
		if apierrors.IsNotFound(err) || machinerr.IsNotFound(err) {
			interfaces = baseConnectionInfo(&nicsSpec[nic], interfaces, machineInterfaces)
			continue
		}

		switchInterface, ok := s.Status.Interfaces[nicsSpec[nic].LLDPs[0].PortDescription]
		if !ok {
			interfaces = baseConnectionInfo(&nicsSpec[nic], interfaces, machineInterfaces)
			continue
		}
		interfaces = connectionInfoEnrichment(s.ObjectMeta, &nicsSpec[nic], interfaces, s.Name, switchInterface, machineInterfaces)
	}
	return interfaces
}

func (o *DeviceOnboarding) updateNetworkStatus(
	machineInterfaces []machinev1alpha2.Interface) machinev1alpha2.Network {
	return machinev1alpha2.Network{
		Ports:        len(machineInterfaces),
		Redundancy:   getNetworkRedundancy(machineInterfaces),
		UnknownPorts: countUnknownPorts(machineInterfaces),
	}
}

func (o *DeviceOnboarding) getSwitchByLabel(ctx context.Context,
	label map[string]string) (*switchv1beta1.Switch, error) {
	obj := &switchv1beta1.SwitchList{}
	filter := &ctrlclient.ListOptions{
		LabelSelector: ctrlclient.MatchingLabelsSelector{Selector: labels.SelectorFromSet(label)},
	}
	if err := o.client.List(ctx, obj, filter); err != nil {
		return nil, err
	}
	if len(obj.Items) == 0 {
		return nil, machinerr.NotFound("switch")
	}
	return &obj.Items[0], nil
}

func connectionInfoEnrichment(sw metav1.ObjectMeta, nicsSpec *inventoriesv1alpha1.NICSpec,
	interfaces []machinev1alpha2.Interface,
	switchUUID string, switchInterface *switchv1beta1.InterfaceSpec,
	machineInterfaces []machinev1alpha2.Interface) []machinev1alpha2.Interface {
	return append(interfaces, machinev1alpha2.Interface{
		Name:            nicsSpec.Name,
		Lanes:           switchInterface.Lanes,
		IPv4:            &machinev1alpha2.IPAddressSpec{Address: getAddress(switchInterface, ipamv1alpha1.CIPv4SubnetType)},
		IPv6:            &machinev1alpha2.IPAddressSpec{Address: getAddress(switchInterface, ipamv1alpha1.CIPv6SubnetType)},
		Moved:           getMovedInterface(nicsSpec, machineInterfaces),
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

func getAddress(switchNIC *switchv1beta1.InterfaceSpec, af ipamv1alpha1.SubnetAddressType) string {
	var ipAddress string
	for _, addr := range switchNIC.IP {
		ip, ipNet, err := net.ParseCIDR(addr.Address)
		if err != nil {
			return ipAddress
		}
		if size, _ := ipNet.Mask.Size(); size < subnetSize {
			return ipAddress
		}
		if ipByteRepr := ip.To4(); ipByteRepr != nil && af == ipamv1alpha1.CIPv4SubnetType {
			ipByteRepr[3]++
			machineAddr := net.IPNet{
				IP:   ip,
				Mask: ipNet.Mask,
			}
			return machineAddr.String()
		}
		if ipByteRepr := ip.To16(); ipByteRepr != nil && af == ipamv1alpha1.CIPv6SubnetType {
			ipByteRepr[15]++
			machineAddr := net.IPNet{
				IP:   ip,
				Mask: ipNet.Mask,
			}
			return machineAddr.String()
		}
	}
	return ipAddress
}

func baseConnectionInfo(nicsSpec *inventoriesv1alpha1.NICSpec,
	interfaces []machinev1alpha2.Interface, machineInterfaces []machinev1alpha2.Interface) []machinev1alpha2.Interface {
	if len(nicsSpec.LLDPs) != 1 {
		return append(interfaces, machinev1alpha2.Interface{
			Name:    nicsSpec.Name,
			Unknown: true,
		})
	}
	return append(interfaces, machinev1alpha2.Interface{
		Name:    nicsSpec.Name,
		Unknown: false,
		Moved:   getMovedInterface(nicsSpec, machineInterfaces),
		Peer: &machinev1alpha2.Peer{
			LLDPSystemName:      nicsSpec.LLDPs[0].SystemName,
			LLDPChassisID:       nicsSpec.LLDPs[0].ChassisID,
			LLDPPortID:          nicsSpec.LLDPs[0].PortID,
			LLDPPortDescription: nicsSpec.LLDPs[0].PortDescription,
		},
	})
}

func getMovedInterface(newInterfaceState *inventoriesv1alpha1.NICSpec,
	machineInterfaces []machinev1alpha2.Interface) bool {
	for mi := range machineInterfaces {
		if machineInterfaces[mi].Name != newInterfaceState.Name {
			continue
		}
		if machineInterfaces[mi].Peer == nil || len(newInterfaceState.LLDPs) == 0 {
			continue
		}
		if machineInterfaces[mi].Peer.LLDPChassisID != newInterfaceState.LLDPs[0].ChassisID {
			return true
		}
	}
	return false
}

func getNetworkRedundancy(machineInterfaces []machinev1alpha2.Interface) string {
	switch {
	case len(machineInterfaces) == onePort:
		return machinev1alpha2.InterfaceRedundancySingle
	case len(machineInterfaces) >= twoPorts:
		return machinev1alpha2.InterfaceRedundancyHighAvailability
	default:
		return machinev1alpha2.InterfaceRedundancyNone
	}
}

func countUnknownPorts(machineInterfaces []machinev1alpha2.Interface) int {
	var count int
	for machinePort := range machineInterfaces {
		if !(machineInterfaces[machinePort].Unknown) {
			continue
		}
		count++
	}
	return count
}

func updateHealthStatus(machineStatus machinev1alpha2.MachineStatus) machinev1alpha2.MachineState {
	if !machineStatus.OOB.Exist || !machineStatus.Inventory.Exist ||
		len(machineStatus.Interfaces) < defaultNumberOfInterfaces {
		return machinev1alpha2.MachineStateUnhealthy
	} else {
		return machinev1alpha2.MachineStateHealthy
	}
}
