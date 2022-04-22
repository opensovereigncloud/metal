package reserve

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Reserver interface {
	Reserve(*metav1.Object, *metav1.Object) error
	DeleteReservation(machine *machinev1alpha2.Machine) error
	CheckIn(*machinev1alpha2.Machine) error
	CheckOut(*machinev1alpha2.Machine) error
}
