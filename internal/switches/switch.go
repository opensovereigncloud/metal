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

package switches

import (
	"context"

	"github.com/go-logr/logr"
	machinev1alpha1 "github.com/onmetal/metal-api/apis/machine/v1alpha1"
	switchv1alpha1 "github.com/onmetal/metal-api/apis/switches/v1alpha1"
	machinerr "github.com/onmetal/metal-api/internal/errors"
	"github.com/onmetal/metal-api/internal/provider"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type Switch struct {
	ctrlclient.Client
	*switchv1alpha1.Switch

	ctx context.Context
	log logr.Logger
}

func New(ctx context.Context, c ctrlclient.Client, l logr.Logger, req ctrl.Request) (*Switch, error) {
	swobj, err := provider.Get(ctx, c, req.Name, req.Namespace, "switch")
	if err != nil {
		return &Switch{}, err
	}
	sw, ok := swobj.(*switchv1alpha1.Switch)
	if !ok {
		return &Switch{}, machinerr.CastType()
	}

	return &Switch{
		Client: c,
		Switch: sw,
		ctx:    ctx,
		log:    l,
	}, nil
}

func (s *Switch) UpdateLocation() error {
	machinesobj, err := provider.List(s.ctx, s.Client, provider.Machine)
	if err != nil {
		return err
	}
	machineList, ok := machinesobj.(*machinev1alpha1.MachineList)
	if !ok {
		return machinerr.CastType()
	}
	machines := machineList.Items
	for m := range machines {
		for i := range machines[m].Status.Interfaces {
			if _, ok := s.Switch.Status.Interfaces[machines[m].Status.Interfaces[i].LLDPPortDescription]; !ok {
				continue
			}
			if isUpdateRequired(&machines[m]) {
				machines[m].Spec.Location = s.updateMachineLocation()
				return s.Client.Update(s.ctx, &machines[m])
			}
			return nil
		}
	}
	return nil
}

func (s *Switch) updateMachineLocation() machinev1alpha1.Location {
	return machinev1alpha1.Location{
		DataHall: s.Switch.Spec.Location.Room,
		Row:      s.Switch.Spec.Location.Row,
		Rack:     s.Switch.Spec.Location.Rack,
	}
}

func isUpdateRequired(m *machinev1alpha1.Machine) bool {
	return true
}
