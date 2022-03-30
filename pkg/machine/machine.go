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
	machinerr "github.com/onmetal/metal-api/pkg/errors"
	oobonmetal "github.com/onmetal/oob-controller/api/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type Machiner interface {
	GetMachine(string, string) (*machinev1alpha2.Machine, error)
	Reservation(*machinev1alpha2.Machine) error
	UpdateSpec(*machinev1alpha2.Machine) error
	UpdateStatus(*machinev1alpha2.Machine) error
}

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

func (m *Machine) GetMachine(name, namespace string) (*machinev1alpha2.Machine, error) {
	obj := &machinev1alpha2.Machine{}
	if err := m.Client.Get(m.ctx, types.NamespacedName{Name: name, Namespace: namespace}, obj); err != nil {
		return nil, err
	}
	return obj, nil
}

func (m *Machine) Reservation(machine *machinev1alpha2.Machine) error {
	return m.enableServer(machine.Name, machine.Namespace)
}

func (m *Machine) enableServer(name, namespace string) error {
	oobObj, err := m.getOOBMachineByUUIDLabel(name, namespace)
	if err != nil {
		return err
	}

	oobObj.Spec.PowerState = getPowerState(oobObj.Spec.PowerState)

	m.log.Info("oob state changed", "uuid", "namespace", oobObj.Name, oobObj.Namespace)
	return m.Client.Patch(m.ctx, oobObj, ctrlclient.Merge, &ctrlclient.PatchOptions{
		FieldManager: "machine-controller",
	})
}

func (m *Machine) UpdateSpec(machine *machinev1alpha2.Machine) error {
	return m.Client.Update(m.ctx, machine)
}

func (m *Machine) UpdateStatus(machine *machinev1alpha2.Machine) error {
	m.updateEventLog("", machine)

	m.updateHealthStatus(machine)

	return m.Client.Status().Update(m.ctx, machine)
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
		return nil, machinerr.NotFound(name)
	}
	return &oobs.Items[0], nil
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

func getPowerState(state string) string {
	switch state {
	case "On":
		// In case when machine already running Reset is required.
		// Because it will bring machine from scratch.
		return "Reset"
	default:
		return "On"
	}
}
