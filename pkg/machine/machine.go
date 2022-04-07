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

package machine

import (
	"context"

	"github.com/go-logr/logr"
	machinev1alpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	requestv1alpha1 "github.com/onmetal/metal-api/apis/request/v1alpha1"
	metalerr "github.com/onmetal/metal-api/pkg/errors"
	oobonmetal "github.com/onmetal/oob-controller/api/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type Machiner interface {
	Reserve(*machinev1alpha2.Machine, *requestv1alpha1.Request) error
	CheckIn(*machinev1alpha2.Machine) error
	CheckOut(*machinev1alpha2.Machine) error

	GetMachine(string, string) (*machinev1alpha2.Machine, error)

	FindVacantMachine(*requestv1alpha1.Request) (*machinev1alpha2.Machine, error)

	UpdateSpec(*machinev1alpha2.Machine) error
	UpdateStatus(*machinev1alpha2.Machine) error
}

const (
	LeasedSizeLabel   = "machine.onmetal.de/leased-size"
	LeasedPoolLabel   = "machine.onmetal.de/leasing-pool"
	LeasedLabel       = "machine.onmetal.de/leased"
	MetalRequestLabel = "machine.onmetal.de/metal-request"
)

const pageListLimit = 1000

type Machine struct {
	ctrlclient.Client

	ctx      context.Context
	log      logr.Logger
	recorder record.EventRecorder
}

func New(ctx context.Context, c ctrlclient.Client, l logr.Logger,
	r record.EventRecorder) Machiner {
	return &Machine{
		Client:   c,
		ctx:      ctx,
		log:      l,
		recorder: r,
	}
}

func (m *Machine) Reserve(machine *machinev1alpha2.Machine, request *requestv1alpha1.Request) error {
	machine.Labels[LeasedLabel] = "true"
	machine.Labels[MetalRequestLabel] = request.Name

	return m.UpdateSpec(machine)
}

func (m *Machine) DeleteReservation(machine *machinev1alpha2.Machine, request *requestv1alpha1.Request) error {
	delete(machine.Labels, LeasedLabel)
	delete(machine.Labels, MetalRequestLabel)

	return m.UpdateSpec(machine)
}

func (m *Machine) CheckIn(machineObj *machinev1alpha2.Machine) error {
	return m.changeServerPowerState(machineObj.Name, machineObj.Namespace, machinev1alpha2.MachinePowerStateON)
}

func (m *Machine) CheckOut(machine *machinev1alpha2.Machine) error {
	return m.changeServerPowerState(machine.Name, machine.Namespace, machinev1alpha2.MachinePowerStateOFF)
}

func (m *Machine) GetMachine(name, namespace string) (*machinev1alpha2.Machine, error) {
	obj := &machinev1alpha2.Machine{}
	if err := m.Client.Get(m.ctx, types.NamespacedName{Name: name, Namespace: namespace}, obj); err != nil {
		return nil, err
	}
	return obj, nil
}

func (m *Machine) FindVacantMachine(metalRequest *requestv1alpha1.Request) (*machinev1alpha2.Machine, error) {
	size, pool := metalRequest.Spec.MachineClass.Name, metalRequest.Spec.MachinePool.Name
	metalSelector, err := getLabelSelectorForAvailableMachine(size, pool)
	if err != nil {
		return &machinev1alpha2.Machine{}, err
	}

	continueToken := ""
	metalList := &machinev1alpha2.MachineList{}
	for {
		opts := &client.ListOptions{
			LabelSelector: metalSelector,
			Namespace:     metalRequest.Namespace,
			Limit:         pageListLimit,
			Continue:      continueToken,
		}

		if err := m.Client.List(m.ctx, metalList, opts); err != nil {
			return &machinev1alpha2.Machine{}, err
		}

		for m := range metalList.Items {
			if !isTolerated(metalRequest.Spec.Tolerations, metalList.Items[m].Spec.Taints) {
				continue
			}
			if metalList.Items[m].Status.Health != machinev1alpha2.MachineStateHealthy {
				continue
			}
			return &metalList.Items[m], nil
		}

		if metalList.Continue == "" || metalList.RemainingItemCount == nil ||
			*metalList.RemainingItemCount == 0 {
			break
		}
	}

	return &machinev1alpha2.Machine{}, metalerr.NotFound(string(requestv1alpha1.Machine))
}

func (m *Machine) UpdateSpec(machine *machinev1alpha2.Machine) error {
	return m.Client.Update(m.ctx, machine)
}

func (m *Machine) UpdateStatus(machine *machinev1alpha2.Machine) error {
	m.updateEventLog("", machine)

	m.updateHealthStatus(machine)

	return m.Client.Status().Update(m.ctx, machine)
}

func (m *Machine) updateEventLog(eventType string, machine *machinev1alpha2.Machine) {
	if m.recorder == nil {
		m.log.Info("event recorder not provided")
		return
	}
	if eventType == "" {
		eventType = "Updated"
	}
	lastEvent := len(machine.ManagedFields) - 1
	m.recorder.Eventf(machine, "Normal", eventType,
		"machine: %s/%s, was updated by %s, at %s",
		machine.Namespace, machine.Name,
		machine.ManagedFields[lastEvent].Manager, machine.ManagedFields[lastEvent].Time)
}

func (m *Machine) updateHealthStatus(machine *machinev1alpha2.Machine) {
	if !machine.Status.OOB.Exist || !machine.Status.Inventory.Exist ||
		len(machine.Status.Interfaces) == 0 {
		machine.Status.Health = machinev1alpha2.MachineStateUnhealthy
		machine.Status.Orphaned = true
	} else {
		machine.Status.Health = machinev1alpha2.MachineStateHealthy
		machine.Status.Orphaned = false
	}
}

func (m *Machine) changeServerPowerState(name, namespace, powerState string) error {
	oobObj, err := m.getOOBMachineByUUIDLabel(name, namespace)
	if err != nil {
		return err
	}

	if powerState == oobObj.Spec.PowerState {
		powerState = machinev1alpha2.MachinePowerStateReset
	}

	oobObj.Spec.PowerState = powerState

	m.log.Info("oob state changed", "uuid", "namespace", oobObj.Name, oobObj.Namespace)
	return m.Client.Patch(m.ctx, oobObj, ctrlclient.Merge, &ctrlclient.PatchOptions{
		FieldManager: "machine-controller",
	})
}

func (m *Machine) getOOBMachineByUUIDLabel(name, namespace string) (*oobonmetal.Machine, error) {
	oobs := &oobonmetal.MachineList{}
	listOptions := &ctrlclient.ListOptions{
		Namespace: namespace,
		LabelSelector: ctrlclient.MatchingLabelsSelector{
			Selector: labels.SelectorFromSet(map[string]string{
				machinev1alpha2.UUIDLabel: name,
			})}}
	if err := m.Client.List(m.ctx, oobs, listOptions); err != nil {
		return nil, err
	}

	if len(oobs.Items) == 0 {
		return nil, metalerr.NotFound(name)
	}
	return &oobs.Items[0], nil
}

func getLabelSelectorForAvailableMachine(size, pool string) (labels.Selector, error) {
	labelSizeRequirement, err := labels.NewRequirement(LeasedSizeLabel, selection.Equals, []string{size})
	if err != nil {
		return nil, err
	}
	labelPoolRequirement, err := labels.NewRequirement(LeasedPoolLabel, selection.Equals, []string{pool})
	if err != nil {
		return nil, err
	}
	labelNotReserved, err := labels.NewRequirement(LeasedLabel, selection.DoesNotExist, []string{})
	if err != nil {
		return nil, err
	}
	return labels.NewSelector().
			Add(*labelPoolRequirement).
			Add(*labelSizeRequirement).
			Add(*labelNotReserved),
		nil
}

func isTolerated(tolerations []requestv1alpha1.Toleration, taints []machinev1alpha2.Taint) bool {
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

func toleratesTaint(toleration requestv1alpha1.Toleration, taint machinev1alpha2.Taint) bool {
	if toleration.Effect != taint.Effect {
		return false
	}

	if toleration.Key != taint.Key {
		return false
	}

	switch toleration.Operator {
	case "", requestv1alpha1.TolerationOpEqual: // empty operator means Equal
		return toleration.Value == taint.Value
	case requestv1alpha1.TolerationOpExists:
		return true
	default:
		return false
	}
}
