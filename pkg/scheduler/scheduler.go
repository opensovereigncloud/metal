package scheduler

import (
	requestv1alpha1 "github.com/onmetal/metal-api/apis/request/v1alpha1"
)

type Scheduler interface {
	Schedule(*requestv1alpha1.Request) error
	DeleteScheduling(*requestv1alpha1.Request) error
}
