package scheduler

import (
	machinev1alpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
)

type Scheduler interface {
	Schedule(*machinev1alpha2.MachineAssignment) error
	DeleteScheduling(*machinev1alpha2.MachineAssignment) error
}

func IsAldreadyScheduled(metalRequest *machinev1alpha2.MachineAssignment) bool {
	return metalRequest.Status.Reference != nil
}
