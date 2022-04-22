package machine

import (
	"testing"

	machinev1alpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	requestv1alpha1 "github.com/onmetal/metal-api/apis/request/v1alpha1"
)

func TestIsTolerated(t *testing.T) {
	tolerations := []requestv1alpha1.Toleration{
		{
			Key:    "uuid",
			Value:  "123",
			Effect: machinev1alpha2.TaintEffectNotAvailable,
		},
	}
	taints := []machinev1alpha2.Taint{
		{
			Key:    "uuid",
			Value:  "123",
			Effect: machinev1alpha2.TaintEffectNotAvailable,
		},
	}

	if !isTolerated(tolerations, taints) {
		t.Error("should be tolerated")
	}
}

func TestIsNotTolerated(t *testing.T) {
	tolerations := []requestv1alpha1.Toleration{
		{
			Key:    "uuid",
			Value:  "123",
			Effect: machinev1alpha2.TaintEffectNotAvailable,
		},
	}
	taints := []machinev1alpha2.Taint{
		{
			Key:    "uuid",
			Value:  "123",
			Effect: machinev1alpha2.TaintEffectNotAvailable,
		},
		{
			Key:    "uuid",
			Value:  "123",
			Effect: machinev1alpha2.TaintEffectError,
		},
	}

	if isTolerated(tolerations, taints) {
		t.Error("should not be tolerated")
	}
}

func TestIsNotToleratedEqual(t *testing.T) {
	tolerations := []requestv1alpha1.Toleration{
		{
			Key:    "uuid",
			Value:  "123",
			Effect: machinev1alpha2.TaintEffectNotAvailable,
		},
		{
			Key:    "uuid",
			Value:  "1234",
			Effect: machinev1alpha2.TaintEffectSuspended,
		},
	}
	taints := []machinev1alpha2.Taint{
		{
			Key:    "uuid",
			Value:  "123",
			Effect: machinev1alpha2.TaintEffectNotAvailable,
		},
		{
			Key:    "uuid",
			Value:  "1234",
			Effect: machinev1alpha2.TaintEffectError,
		},
	}

	if isTolerated(tolerations, taints) {
		t.Error("should not be tolerated")
	}
}
