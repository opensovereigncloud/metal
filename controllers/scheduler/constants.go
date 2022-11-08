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
