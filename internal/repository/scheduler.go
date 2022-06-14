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

	inventoriesv1alpha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	machinev1alpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	"github.com/onmetal/metal-api/internal/entity"
	metalerr "github.com/onmetal/metal-api/pkg/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

const pageListLimit = 1000

type MachineSchedulerRepo struct {
	ctrlclient.Client
}

func NewMachineSchedulerRepo(c ctrlclient.Client) *MachineSchedulerRepo {
	return &MachineSchedulerRepo{
		Client: c,
	}
}

func (m *MachineSchedulerRepo) Schedule(ctx context.Context, e entity.Reservation) error {
	metalAssignment := &machinev1alpha2.MachineAssignment{}
	if err := m.Client.Get(ctx, types.NamespacedName{
		Name: e.OrderName, Namespace: e.OrderNamespace}, metalAssignment); err != nil {
		return err
	}
	machine := &machinev1alpha2.Machine{}
	if err := m.Client.Get(ctx, types.NamespacedName{
		Name: e.RequestName, Namespace: e.RequestNamespace}, machine); err != nil {
		return err
	}

	machine.Status.Reservation.Status = entity.ReservationStatusPending
	machine.Status.Reservation.Reference = prepareReferenceSpec(e.OrderName, e.OrderNamespace)

	if err := m.Client.Status().Update(ctx, machine); err != nil {
		return err
	}

	metalAssignment.Status.Reference = getObjectReference(machine)
	metalAssignment.Status.State = entity.ReservationStatusReserved

	return m.Client.Status().Update(ctx, metalAssignment)
}

func (m *MachineSchedulerRepo) DeleteSchedule(ctx context.Context, e entity.Reservation) error {
	machine := &machinev1alpha2.Machine{}
	if err := m.Client.Get(ctx, types.NamespacedName{
		Name: e.RequestName, Namespace: e.RequestNamespace}, machine); err != nil {
		return err
	}
	machine.Status.Reservation.Status = entity.ReservationStatusAvailable
	machine.Status.Reservation.Reference = nil

	return m.Client.Status().Update(ctx, machine)
}

func (m *MachineSchedulerRepo) IsScheduled(ctx context.Context, e entity.Order) bool {
	metalAssignment := &machinev1alpha2.MachineAssignment{}
	if err := m.Client.Get(ctx, types.NamespacedName{
		Name: e.Name, Namespace: e.Namespace}, metalAssignment); err != nil {
		return true
	}
	return metalAssignment.Status.Reference != nil
}

func (m *MachineSchedulerRepo) FindVacantDevice(ctx context.Context, e entity.Order) (entity.Reservation, error) {
	metalAssignment := &machinev1alpha2.MachineAssignment{}
	if err := m.Client.Get(ctx,
		types.NamespacedName{Name: e.Name, Namespace: e.Namespace}, metalAssignment); err != nil {
		return entity.Reservation{}, err
	}

	size := metalAssignment.Spec.MachineClass.Name
	metalSelector, err := getLabelSelectorForAvailableMachine(size)
	if err != nil {
		return entity.Reservation{}, err
	}

	continueToken := ""
	metalList := &machinev1alpha2.MachineList{}
	for {
		opts := &client.ListOptions{
			LabelSelector: metalSelector,
			Limit:         pageListLimit,
			Continue:      continueToken,
		}

		if err := m.Client.List(ctx, metalList, opts); err != nil {
			return entity.Reservation{}, err
		}

		for m := range metalList.Items {
			if !isTolerated(metalAssignment.Spec.Tolerations, metalList.Items[m].Spec.Taints) {
				continue
			}
			if metalList.Items[m].Status.Reservation.Reference != nil {
				continue
			}
			if metalList.Items[m].Status.Health != machinev1alpha2.MachineStateHealthy {
				continue
			}
			return entity.Reservation{
				OrderName:        metalAssignment.Name,
				OrderNamespace:   metalAssignment.Namespace,
				RequestName:      metalList.Items[m].Name,
				RequestNamespace: metalList.Items[m].Namespace,
			}, nil
		}

		if metalList.Continue == "" || metalList.RemainingItemCount == nil ||
			*metalList.RemainingItemCount == 0 {
			break
		}
	}

	return entity.Reservation{}, metalerr.NotFound("machines for request")
}

func getLabelSelectorForAvailableMachine(size string) (labels.Selector, error) {
	sizeLabel := inventoriesv1alpha1.CLabelPrefix + size
	labelSizeRequirement, err := labels.NewRequirement(sizeLabel, selection.Exists, []string{})
	if err != nil {
		return nil, err
	}

	labelNotReserved, err := labels.NewRequirement(machinev1alpha2.LeasedLabel, selection.DoesNotExist, []string{})
	if err != nil {
		return nil, err
	}
	return labels.NewSelector().
			Add(*labelSizeRequirement).
			Add(*labelNotReserved),
		nil
}

func isTolerated(tolerations []machinev1alpha2.Toleration, taints []machinev1alpha2.Taint) bool {
	if len(taints) == 0 {
		return true
	}
	if len(tolerations) != len(taints) {
		return false
	}
	tolerated := 0
	for t := range tolerations {
		for taint := range taints {
			if !toleratesTaint(tolerations[t], taints[taint]) {
				continue
			}
			tolerated++
		}
	}
	return tolerated == len(taints)
}

func toleratesTaint(toleration machinev1alpha2.Toleration, taint machinev1alpha2.Taint) bool {
	if toleration.Effect != taint.Effect {
		return false
	}

	if toleration.Key != taint.Key {
		return false
	}

	switch toleration.Operator {
	case "", machinev1alpha2.TolerationOpEqual: // empty operator means Equal
		return toleration.Value == taint.Value
	case machinev1alpha2.TolerationOpExists:
		return true
	default:
		return false
	}
}

func getObjectReference(m *machinev1alpha2.Machine) *machinev1alpha2.ResourceReference {
	return &machinev1alpha2.ResourceReference{
		APIVersion: m.APIVersion,
		Kind:       m.Kind,
		Name:       m.Name,
		Namespace:  m.Namespace,
	}
}
