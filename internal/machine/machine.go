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
	"strings"

	"github.com/go-logr/logr"
	machinev1alpha1 "github.com/onmetal/metal-api/apis/machine/v1alpha1"
	machinerr "github.com/onmetal/metal-api/internal/errors"
	"github.com/onmetal/metal-api/internal/provider"
	oobonmetal "github.com/onmetal/oob-controller/api/v1"
	switchv1alpha1 "github.com/onmetal/switch-operator/api/v1alpha1"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type Machine struct {
	ctrlclient.Client
	*machinev1alpha1.Machine

	labels   map[string]string
	ctx      context.Context
	log      logr.Logger
	recorder record.EventRecorder
	request  ctrl.Request
}

const (
	interfaceRedundancySingle           = "single"
	interfaceRedundancyHighAvailability = "high availability"
	interfaceRedundancyNone             = "none"
)

const (
	onePort = 1 + iota
	twoPorts
)

const UUIDLabel = "machine.onmetal.de/uuid"

func New(ctx context.Context, c ctrlclient.Client, l logr.Logger,
	r record.EventRecorder, req ctrl.Request) (*Machine, error) {
	obj, err := provider.Get(ctx, c, req.Name, req.Namespace, provider.Machine)
	if err != nil {
		return nil, err
	}
	machine, ok := obj.(*machinev1alpha1.Machine)
	if !ok {
		return &Machine{}, machinerr.CastType()
	}
	return &Machine{
		Client:   c,
		Machine:  machine,
		ctx:      ctx,
		log:      l,
		recorder: r,
		request:  req,
		labels:   map[string]string{UUIDLabel: machine.Name},
	}, nil
}

func (m *Machine) Update() error {
	if m.Spec.Action.PowerState != "" {
		if err := m.updateOOBPowerState(); err != nil {
			return err
		}
		m.Machine.Status.Reboot = "pending"
		m.addEventOnAction(m.Spec.Action.PowerState)
		m.Spec.Action.PowerState = ""
	}

	if err := m.updateStatus(); err != nil {
		return err
	}

	m.updateEventLog()

	m.compare()

	return m.Client.Update(m.ctx, m.Machine)
}

func (m *Machine) updateStatus() error {
	m.Machine.Status.Redundancy = m.getNetworkRedundancy()
	m.Machine.Status.Ports = len(m.Machine.Status.Interfaces)
	m.Machine.Status.UnknownPorts = m.getUnknownPortsCount()
	if !(m.Machine.Status.OOB && m.Machine.Status.Inventory && len(m.Machine.Status.Interfaces) != 0) {
		m.Machine.Status.Health = "unhealthy"
		m.Machine.Status.Orphaned = true
		return m.Client.Status().Update(m.ctx, m.Machine)
	}
	m.Machine.Status.Health = "healthy"
	m.Machine.Status.Orphaned = false
	return m.Client.Status().Update(m.ctx, m.Machine)
}

func (m *Machine) getNetworkRedundancy() string {
	switch {
	case len(m.Machine.Status.Interfaces) == onePort:
		return interfaceRedundancySingle
	case len(m.Machine.Status.Interfaces) >= twoPorts:
		if m.Machine.Status.Interfaces[0].LLDPChassisID != m.Machine.Status.Interfaces[1].LLDPChassisID {
			return interfaceRedundancyHighAvailability
		}
		return interfaceRedundancySingle
	default:
		return interfaceRedundancyNone
	}
}

func (m *Machine) getUnknownPortsCount() int {
	var count int
	for machinePort := range m.Machine.Status.Interfaces {
		if !(m.Machine.Status.Interfaces[machinePort].Unknown) {
			continue
		}
		count++
	}
	return count
}

func (m *Machine) updateOOBPowerState() error {
	obj, err := provider.GetByLabel(m.ctx, m.Client, m.labels, provider.OOB)
	if err != nil {
		return err
	}
	oob, ok := obj.(*oobonmetal.Machine)
	if !ok {
		return machinerr.CastType()
	}
	if oob.Spec.PowerState == machinev1alpha1.MachinPowerStateON {
		return nil
	}
	oob.Spec.PowerState = m.Spec.Action.PowerState
	m.log.Info("oob state changed", "uuid", m.Name)
	return m.Client.Update(m.ctx, oob)
}

func (m *Machine) compare() { m.CompareSwitchInformation() }

func (m *Machine) CompareSwitchInformation() {
	if len(m.Machine.Status.Interfaces) < 1 {
		return
	}
	label := map[string]string{switchv1alpha1.LabelChassisId: strings.ReplaceAll(m.Machine.Status.Interfaces[0].LLDPChassisID, ":", "-")}
	obj, err := provider.GetByLabel(m.ctx, m.Client, label, provider.Switch)
	if err != nil {
		if machinerr.IsNotExist(err) {
			return
		}
		m.log.Info("switch object get failed", "error", err)
		return
	}
	s, ok := obj.(*switchv1alpha1.Switch)
	if !ok {
		m.log.Info("switch object casting failed")
		return
	}
	if s.Spec.Location == nil {
		return
	}
	m.compareMachineAndSwitchLocation(s.Spec.Location, &m.Spec.Location)
}

func (m *Machine) compareMachineAndSwitchLocation(switchLocation *switchv1alpha1.LocationSpec,
	location *machinev1alpha1.Location) {
	if switchLocation.Room != location.DataHall {
		m.log.Info("location has different meaning",
			"switch datahall", switchLocation.Room,
			"machine datahall", location.DataHall)
	}
	if switchLocation.Rack != location.Rack {
		m.log.Info("location has different meaning",
			"switch rack", switchLocation.Rack,
			"machine rack", location.Rack)
	}
	if switchLocation.Row != location.Row {
		m.log.Info("location has different meaning",
			"switch row", switchLocation.Row,
			"machine row", location.Row)
	}
}

func (m *Machine) addEventOnAction(eventType string) {
	lastEvent := len(m.ManagedFields) - 1
	m.recorder.Eventf(m.Machine, "Normal", eventType,
		"machine state was changed by %s at %s",
		m.ManagedFields[lastEvent].Manager, m.ManagedFields[lastEvent].Time)
}

func (m *Machine) updateEventLog() {
	lastEvent := len(m.ManagedFields) - 1
	m.recorder.Eventf(m.Machine, "Normal", "Updated",
		"machine: %s/%s, was updated by %s, at %s",
		m.Namespace, m.Name,
		m.ManagedFields[lastEvent].Manager, m.ManagedFields[lastEvent].Time)
}
