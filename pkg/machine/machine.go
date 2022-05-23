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
	metalerr "github.com/onmetal/metal-api/pkg/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type Machiner interface {
	FindVacantMachine(*machinev1alpha2.MachineAssignment) (*machinev1alpha2.Machine, error)
	PatchSpec(*machinev1alpha2.Machine) error
	PatchStatus(*machinev1alpha2.Machine) error
}

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

func (m *Machine) FindVacantMachine(metalRequest *machinev1alpha2.MachineAssignment) (*machinev1alpha2.Machine, error) {
	size := metalRequest.Spec.MachineClass.Name
	metalSelector, err := getLabelSelectorForAvailableMachine(size)
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

	return &machinev1alpha2.Machine{}, metalerr.NotFound("machines for request")
}

func (m *Machine) PatchSpec(machine *machinev1alpha2.Machine) error {
	m.updateEventLog("", machine)

	return m.Client.Patch(m.ctx, machine, ctrlclient.Merge, &ctrlclient.PatchOptions{
		FieldManager: "machine",
	})
}

func (m *Machine) PatchStatus(machine *machinev1alpha2.Machine) error {
	m.updateEventLog("", machine)

	m.updateHealthStatus(machine)

	return m.Client.Status().Patch(m.ctx, machine, ctrlclient.Merge, &ctrlclient.PatchOptions{
		FieldManager: "machine",
	})
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

func getLabelSelectorForAvailableMachine(size string) (labels.Selector, error) {
	labelSizeRequirement, err := labels.NewRequirement(machinev1alpha2.LeasedSizeLabel, selection.Equals, []string{size})
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
