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

package persistence

import (
	"context"

	machine "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	domain "github.com/onmetal/metal-api/domain/machine"
	usecase "github.com/onmetal/metal-api/usecase/onboarding"
	"github.com/onmetal/metal-api/usecase/onboarding/dto"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	onePort = 1 + iota
	twoPorts
)

const subnetSize = 30
const defaultNumberOfInterfaces = 2

type MachineRepository struct {
	client ctrlclient.Client
}

func NewMachineRepository(client ctrlclient.Client) *MachineRepository {
	return &MachineRepository{client: client}
}

func (r *MachineRepository) Create(inventory dto.Inventory) error {
	machineObj := prepareMachine(inventory)
	if err := r.
		client.
		Create(
			context.Background(),
			machineObj); err != nil {
		if apierrors.IsAlreadyExists(err) {
			return usecase.MachineAlreadyOnboarded(inventory.UUID)
		}
		return err
	}
	return nil
}

func (r *MachineRepository) Update(machine domain.Machine) error {
	machineObj, err := r.extractMachineFromCluster(machine.UUID, machine.Namespace)
	if err != nil {
		return err
	}

	r.updateMachine(machineObj, machine)

	if err := r.
		client.
		Update(
			context.Background(),
			machineObj); err != nil {
		return err
	}
	r.updateMachineStatus(machineObj, machine)
	return r.
		client.
		Status().
		Update(
			context.Background(),
			machineObj)
}
func (r *MachineRepository) Get(request dto.Request) (domain.Machine, error) {
	machineObj, err := r.extractMachineFromCluster(request.Name, request.Namespace)
	if apierrors.IsNotFound(err) {
		return domain.Machine{}, usecase.MachineNotFound(request.Name)
	}
	if err != nil {
		return domain.Machine{}, err
	}
	return domain.NewMachine(
		machineObj.Name,
		machineObj.Namespace,
		machineObj.Spec.Identity.SKU,
		machineObj.Spec.Identity.SerialNumber,
		machineObj.Status.Interfaces,
		machineObj.Labels), nil
}

func (r *MachineRepository) extractMachineFromCluster(name, namespace string) (*machine.Machine, error) {
	obj := &machine.Machine{}
	if err := r.
		client.
		Get(
			context.Background(),
			types.NamespacedName{
				Namespace: namespace,
				Name:      name,
			},
			obj); err != nil {
		return nil, err
	}
	return obj, nil
}

func (r *MachineRepository) updateMachine(machineObj *machine.Machine, machine domain.Machine) {
	machineObj.Labels = copySizeLabels(machineObj.Labels, machine.Size)

	machineObj.Spec.Identity.SKU = machine.SKU
	machineObj.Spec.Identity.SerialNumber = machine.SerialNumber
}

func (r *MachineRepository) updateMachineStatus(machineObj *machine.Machine, domainMachine domain.Machine) {
	if machineObj.Status.Reservation.Status == "" {
		machineObj.Status.Reservation.Status = machine.ReservationStatusAvailable
	}

	machineObj.Status.Interfaces = domainMachine.Interfaces

	machineObj.Status.Health = updateHealthStatus(machineObj.Status.Interfaces)

	machineObj.Status.Network = networkStatus(machineObj.Status.Interfaces)
}

func updateHealthStatus(interfaces []machine.Interface) machine.MachineState {
	if len(interfaces) < defaultNumberOfInterfaces {
		return machine.MachineStateUnhealthy
	}
	return machine.MachineStateHealthy
}

func prepareMachine(inventory dto.Inventory) *machine.Machine {
	return &machine.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      inventory.UUID,
			Namespace: inventory.Namespace,
		},
		Spec: machine.MachineSpec{
			Identity: machine.Identity{
				SKU:          inventory.ProductSKU,
				SerialNumber: inventory.SerialNumber,
			},
		},
		Status: machine.MachineStatus{},
	}
}

func copySizeLabels(machineLabels, sizeLabels map[string]string) map[string]string {
	if machineLabels == nil {
		return sizeLabels
	}
	for key, value := range sizeLabels {
		machineLabels[key] = value
	}
	return machineLabels
}

func networkStatus(machineInterfaces []machine.Interface) machine.Network {
	return machine.Network{
		Ports:        len(machineInterfaces),
		Redundancy:   getNetworkRedundancy(machineInterfaces),
		UnknownPorts: countUnknownPorts(machineInterfaces),
	}
}

func getNetworkRedundancy(machineInterfaces []machine.Interface) string {
	switch {
	case len(machineInterfaces) == onePort:
		return machine.InterfaceRedundancySingle
	case len(machineInterfaces) >= twoPorts:
		return machine.InterfaceRedundancyHighAvailability
	default:
		return machine.InterfaceRedundancyNone
	}
}

func countUnknownPorts(machineInterfaces []machine.Interface) int {
	var count int
	for machinePort := range machineInterfaces {
		if !machineInterfaces[machinePort].Unknown {
			continue
		}
		count++
	}
	return count
}
