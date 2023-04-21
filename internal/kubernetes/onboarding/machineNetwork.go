// Copyright 2023 OnMetal authors
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

package persistence

import (
	"context"
	"net"
	"strings"

	ipamv1alpha1 "github.com/onmetal/ipam/api/v1alpha1"
	inventoriesv1alpha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	machine "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	switchv1beta1 "github.com/onmetal/metal-api/apis/switch/v1beta1"
	"github.com/onmetal/metal-api/internal/usecase/onboarding/dto"
	"github.com/onmetal/metal-api/pkg/constants"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type MachineInterfaces struct {
	client ctrlclient.Client
}

func NewMachineInterfaces(client ctrlclient.Client) *MachineInterfaces {
	return &MachineInterfaces{client: client}
}

func (o *MachineInterfaces) InterfacesFromInventory(inventory dto.Inventory) []machine.Interface {
	return o.updateMachineInterfaces(inventory.NICs)
}

func (o *MachineInterfaces) updateMachineInterfaces(
	nicsSpec []inventoriesv1alpha1.NICSpec) []machine.Interface {
	interfaces := make([]machine.Interface, 0, defaultNumberOfInterfaces)
	for nic := range nicsSpec {
		if len(nicsSpec[nic].LLDPs) == 0 {
			interfaces = baseConnectionInfo(interfaces, &nicsSpec[nic])
			continue
		}

		label := map[string]string{
			constants.LabelChassisID: strings.ReplaceAll(nicsSpec[nic].LLDPs[0].ChassisID, ":", "-"),
		}

		s, err := o.switchByLabel(label)
		if apierrors.IsNotFound(err) || err != nil {
			interfaces = baseConnectionInfo(interfaces, &nicsSpec[nic])
			continue
		}

		if s.GetState() != constants.SwitchStateReady {
			interfaces = baseConnectionInfo(interfaces, &nicsSpec[nic])
			continue
		}

		switchInterface, ok := s.Status.Interfaces[nicsSpec[nic].LLDPs[0].PortDescription]
		if !ok {
			interfaces = baseConnectionInfo(interfaces, &nicsSpec[nic])
			continue
		}
		interfaces = connectionInfoEnrichment(s.ObjectMeta, &nicsSpec[nic], interfaces, s.Name, switchInterface)
	}
	return interfaces
}
func (o *MachineInterfaces) switchByLabel(label map[string]string) (*switchv1beta1.Switch, error) {
	obj := &switchv1beta1.SwitchList{}
	filter := &ctrlclient.ListOptions{
		LabelSelector: ctrlclient.MatchingLabelsSelector{Selector: labels.SelectorFromSet(label)},
	}
	if err := o.client.List(context.Background(), obj, filter); err != nil {
		return nil, err
	}
	if len(obj.Items) == 0 {
		return nil, errNotFound
	}
	return &obj.Items[0], nil
}

func connectionInfoEnrichment(
	sw metav1.ObjectMeta,
	nicsSpec *inventoriesv1alpha1.NICSpec,
	interfaces []machine.Interface,
	switchUUID string,
	switchInterface *switchv1beta1.InterfaceSpec) []machine.Interface {
	return append(interfaces, machine.Interface{
		Name:            nicsSpec.Name,
		Lanes:           uint8(switchInterface.GetLanes()),
		IPv4:            &machine.IPAddressSpec{Address: address(switchInterface, ipamv1alpha1.CIPv4SubnetType)},
		IPv6:            &machine.IPAddressSpec{Address: address(switchInterface, ipamv1alpha1.CIPv6SubnetType)},
		Unknown:         false,
		SwitchReference: &machine.ResourceReference{Kind: "Switch", Namespace: sw.Namespace, Name: sw.Name},
		Peer: &machine.Peer{
			LLDPSystemName:      switchUUID,
			LLDPChassisID:       nicsSpec.LLDPs[0].ChassisID,
			LLDPPortID:          nicsSpec.LLDPs[0].PortID,
			LLDPPortDescription: nicsSpec.LLDPs[0].PortDescription,
		},
	})
}

func address(switchNIC *switchv1beta1.InterfaceSpec, af ipamv1alpha1.SubnetAddressType) string {
	var ipAddress string
	for _, addr := range switchNIC.IP {
		ip, ipNet, err := net.ParseCIDR(addr.GetAddress())
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

func baseConnectionInfo(
	interfaces []machine.Interface,
	nicsSpec *inventoriesv1alpha1.NICSpec) []machine.Interface {
	if len(nicsSpec.LLDPs) != 1 {
		return append(interfaces, machine.Interface{
			Name:    nicsSpec.Name,
			Unknown: true,
		})
	}
	return append(interfaces, machine.Interface{
		Name:    nicsSpec.Name,
		Unknown: false,
		Peer: &machine.Peer{
			LLDPSystemName:      nicsSpec.LLDPs[0].SystemName,
			LLDPChassisID:       nicsSpec.LLDPs[0].ChassisID,
			LLDPPortID:          nicsSpec.LLDPs[0].PortID,
			LLDPPortDescription: nicsSpec.LLDPs[0].PortDescription,
		},
	})
}
