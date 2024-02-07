// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package providers

import (
	"context"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"

	metalv1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"
	"github.com/ironcore-dev/metal/common/types/events"
	ipdomain "github.com/ironcore-dev/metal/domain/address"
	domain "github.com/ironcore-dev/metal/domain/machine"
	reservdomain "github.com/ironcore-dev/metal/domain/reservation"
)

const (
	onePort = 1 + iota
	twoPorts
)

type MachineRepository struct {
	client               ctrlclient.Client
	domainEventPublisher events.DomainEventPublisher
}

func NewMachineRepository(
	client ctrlclient.Client,
	domainEventPublisher events.DomainEventPublisher,
) *MachineRepository {
	return &MachineRepository{
		client:               client,
		domainEventPublisher: domainEventPublisher,
	}
}
func (r *MachineRepository) Save(machine domain.Machine) error {
	uuidOptions := ctrlclient.MatchingFields{
		"metadata.name": machine.UUID,
	}
	machineObj, err := r.extractMachineFromCluster(uuidOptions)
	if err != nil {
		if errors.Is(err, errNotFound) {
			return r.Create(machine)
		}
		return err
	}
	return r.updateMachineSpecAndStatus(machineObj, machine)
}

func (r *MachineRepository) Create(machine domain.Machine) error {
	machineObj := prepareMachine(machine)
	if err := r.
		client.
		Create(
			context.Background(),
			machineObj); err != nil {
		return err
	}
	r.domainEventPublisher.Publish(machine.PopEvents()...)
	return nil
}

func (r *MachineRepository) Update(machine domain.Machine) error {
	uuidOptions := ctrlclient.MatchingFields{
		"metadata.name": machine.UUID,
	}
	machineObj, err := r.extractMachineFromCluster(uuidOptions)
	if err != nil {
		return err
	}
	return r.updateMachineSpecAndStatus(machineObj, machine)
}

func (r *MachineRepository) updateMachineSpecAndStatus(
	machineObj *metalv1alpha4.Machine,
	machine domain.Machine,
) error {
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

func (r *MachineRepository) ByUUID(uuid string) (domain.Machine, error) {
	uuidOptions := ctrlclient.MatchingFields{
		"metadata.name": uuid,
	}
	machineObj, err := r.extractMachineFromCluster(uuidOptions)
	if err != nil {
		return domain.Machine{}, err
	}
	machineID := domain.NewMachineID(machineObj.Labels["id"])
	return domain.Machine{
		ID:           machineID,
		UUID:         machineObj.Name,
		Namespace:    machineObj.Namespace,
		ASN:          machineObj.Status.Network.ASN,
		SKU:          machineObj.Spec.Identity.SKU,
		SerialNumber: machineObj.Spec.Identity.SerialNumber,
		Interfaces:   machineObj.Status.Network.Interfaces,
		Loopbacks:    domainLoopbacks(machineObj.Status.Network.Loopbacks),
		Size:         machineObj.Labels,
	}, nil
}

func (r *MachineRepository) ByID(id domain.MachineID) (domain.Machine, error) {
	uuidOptions := &ctrlclient.ListOptions{
		LabelSelector: labels.SelectorFromSet(map[string]string{"id": id.String()})}
	machineObj, err := r.extractMachineFromCluster(uuidOptions)
	if err != nil {
		return domain.Machine{}, err
	}
	machineID := domain.NewMachineID(machineObj.Labels["id"])
	return domain.Machine{
		ID:           machineID,
		UUID:         machineObj.Name,
		Namespace:    machineObj.Namespace,
		ASN:          machineObj.Status.Network.ASN,
		SKU:          machineObj.Spec.Identity.SKU,
		SerialNumber: machineObj.Spec.Identity.SerialNumber,
		Interfaces:   machineObj.Status.Network.Interfaces,
		Loopbacks:    domainLoopbacks(machineObj.Status.Network.Loopbacks),
		Size:         machineObj.Labels,
	}, nil
}

func domainLoopbacks(loopbacks metalv1alpha4.LoopbackAddresses) domain.Loopbacks {
	return domain.Loopbacks{
		IPv4: ipdomain.Address{
			Prefix: loopbacks.IPv4.Prefix,
		},
		IPv6: ipdomain.Address{
			Prefix: loopbacks.IPv6.Prefix,
		},
	}
}

func (r *MachineRepository) extractMachineFromCluster(options ctrlclient.ListOption) (*metalv1alpha4.Machine, error) {
	obj := &metalv1alpha4.MachineList{}
	if err := r.
		client.
		List(
			context.Background(),
			obj,
			options,
		); err != nil {
		return nil, err
	}
	if len(obj.Items) == 0 {
		return nil, errNotFound
	}
	return &obj.Items[0], nil
}

func (r *MachineRepository) updateMachine(machineObj *metalv1alpha4.Machine, machine domain.Machine) {
	machineObj.Labels = CopySizeLabels(machineObj.Labels, machine.Size)

	machineObj.Spec.Identity.SKU = machine.SKU
	machineObj.Spec.Identity.SerialNumber = machine.SerialNumber
}

func (r *MachineRepository) updateMachineStatus(
	machineObj *metalv1alpha4.Machine,
	domainMachine domain.Machine,
) {
	if machineObj.Status.Reservation.Status == "" {
		machineObj.Status.Reservation.Status = reservdomain.ReservationStatusAvailable
	}
	machineObj.Status.Health = updateHealthStatus(domainMachine.Interfaces)
	machineObj.Status.Network = NetworkStatus(domainMachine)
}

func updateHealthStatus(interfaces []metalv1alpha4.Interface) metalv1alpha4.MachineState {
	if len(interfaces) < 2 {
		return metalv1alpha4.MachineStateUnhealthy
	}
	return metalv1alpha4.MachineStateHealthy
}

func prepareMachine(
	m domain.Machine,
) *metalv1alpha4.Machine {
	m.Size["id"] = m.ID.String()
	return &metalv1alpha4.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.UUID,
			Namespace: m.Namespace,
			Labels:    m.Size,
		},
		Spec: metalv1alpha4.MachineSpec{
			Identity: metalv1alpha4.Identity{
				SKU:          m.SKU,
				SerialNumber: m.SerialNumber,
			},
		},
		Status: metalv1alpha4.MachineStatus{},
	}
}

func CopySizeLabels(
	machineLabels, sizeLabels map[string]string,
) map[string]string {
	if machineLabels == nil {
		return sizeLabels
	}
	for key, value := range sizeLabels {
		machineLabels[key] = value
	}
	return machineLabels
}

func NetworkStatus(
	domainMachine domain.Machine,
) metalv1alpha4.Network {
	return metalv1alpha4.Network{
		ASN:          domainMachine.ASN,
		Ports:        len(domainMachine.Interfaces),
		Redundancy:   NetworkRedundancy(domainMachine.Interfaces),
		UnknownPorts: UnknownPortCount(domainMachine.Interfaces),
		Interfaces:   domainMachine.Interfaces,
		Loopbacks: metalv1alpha4.LoopbackAddresses{
			IPv4: metalv1alpha4.IPAddrSpec{
				Prefix: domainMachine.Loopbacks.IPv4.Prefix,
			},
			IPv6: metalv1alpha4.IPAddrSpec{
				Prefix: domainMachine.Loopbacks.IPv6.Prefix,
			},
		},
	}
}

func NetworkRedundancy(
	machineInterfaces []metalv1alpha4.Interface,
) string {
	switch {
	case len(machineInterfaces) == onePort:
		return domain.InterfaceRedundancySingle
	case len(machineInterfaces) >= twoPorts:
		return domain.InterfaceRedundancyHighAvailability
	default:
		return domain.InterfaceRedundancyNone
	}
}

func UnknownPortCount(machineInterfaces []metalv1alpha4.Interface) int {
	var count int
	for machinePort := range machineInterfaces {
		if !machineInterfaces[machinePort].Unknown {
			continue
		}
		count++
	}
	return count
}
