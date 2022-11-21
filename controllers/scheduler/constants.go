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

package scheduler

import "sigs.k8s.io/controller-runtime/pkg/client"

const (
	SchedulerFinalizer      = "metal-api.onmetal.de/scheduler"
	SchedulerReconcileLabel = "metal-api.onmetal.de/scheduler-reconcile"

	ComputeNameLabel      = "machine.onmetal.de/compute-name"
	ComputeNamespaceLabel = "machine.onmetal.de/compute-namespace"

	ReservationStatusAvailable = "Available"
	ReservationStatusReserved  = "Reserved"
	ReservationStatusPending   = "Pending"
	ReservationStatusError     = "Error"
	ReservationStatusRunning   = "Running"

	IgnitionFieldOwner            = client.FieldOwner("metal-api.onmetal.de/ignition")
	IgnitionConfigMapTemplateName = "ipxe-template"
	IgnitionSecretTemplateName    = "ignition-template"
	IgnitionIpxePrefix            = "ipxe-"
	// TODO(flpeter) duplicate code from switch configurer
	IgnitionASNBase = 4200000000
)
