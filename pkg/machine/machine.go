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
	"k8s.io/client-go/tools/record"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type Machiner interface {
	PatchSpec(*machinev1alpha2.Machine) error
	PatchStatus(*machinev1alpha2.Machine) error
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
	if len(machine.ManagedFields) == 0 {
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
