// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package v1alpha4

import (
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/ironcore-dev/metal/pkg/constants"
)

// log is for logging in this package.
var switchconfiglog = logf.Log.WithName("switchconfig-resource")

func (in *SwitchConfig) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(in).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-metal-ironcore-dev-v1alpha4-switchconfig,mutating=true,failurePolicy=fail,sideEffects=None,groups=metal.ironcore.dev,resources=switchconfigs,verbs=create;update,versions=v1alpha4,name=mswitchconfig.v1alpha4.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &SwitchConfig{}

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (in *SwitchConfig) Default() {
	switchconfiglog.Info("default", "name", in.Name)
	in.setDefaultIPAMSelectors()
	in.setDefaultPortParams()
}

// +kubebuilder:webhook:path=/validate-metal-ironcore-dev-v1alpha4-switchconfig,mutating=false,failurePolicy=fail,sideEffects=None,groups=metal.ironcore.dev,resources=switchconfigs,verbs=create;update,versions=v1alpha4,name=vswitchconfig.v1alpha4.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &SwitchConfig{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (in *SwitchConfig) ValidateCreate() (admission.Warnings, error) {
	if in.Spec.IPAM.CarrierSubnets.FieldSelector != nil {
		return nil, errors.New("field selector is not applicable for carrier subnets")
	}
	if in.Spec.IPAM.LoopbackSubnets.FieldSelector != nil {
		return nil, errors.New("field selector is not applicable for loopback subnets")
	}
	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (in *SwitchConfig) ValidateUpdate(_ runtime.Object) (warnings admission.Warnings, err error) {
	if in.Spec.IPAM.CarrierSubnets.FieldSelector != nil {
		return nil, errors.New("field selector is not applicable for carrier subnets")
	}
	if in.Spec.IPAM.LoopbackSubnets.FieldSelector != nil {
		return nil, errors.New("field selector is not applicable for loopback subnets")
	}
	return
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (in *SwitchConfig) ValidateDelete() (warnings admission.Warnings, err error) {
	return
}

func (in *SwitchConfig) setDefaultIPAMSelectors() {
	// defaulting south subnet selectors
	if in.Spec.IPAM.SouthSubnets == nil {
		in.Spec.IPAM.SouthSubnets = &IPAMSelectionSpec{
			LabelSelector: nil,
			FieldSelector: nil,
		}
	}
	if in.Spec.IPAM.SouthSubnets.LabelSelector == nil {
		in.Spec.IPAM.SouthSubnets.LabelSelector = &metav1.LabelSelector{
			MatchLabels: map[string]string{
				constants.IPAMObjectPurposeLabel: constants.IPAMSouthSubnetPurpose,
			},
		}
	}
	if in.Spec.IPAM.SouthSubnets.FieldSelector == nil {
		in.Spec.IPAM.SouthSubnets.FieldSelector = &FieldSelectorSpec{
			LabelKey: ptr.To(constants.IPAMObjectOwnerLabel),
			FieldRef: &v1.ObjectFieldSelector{
				APIVersion: constants.APIVersion,
				FieldPath:  constants.DefaultIPAMFieldRef,
			},
		}
	}

	// defaulting loopbacks selectors
	if in.Spec.IPAM.LoopbackAddresses == nil {
		in.Spec.IPAM.LoopbackAddresses = &IPAMSelectionSpec{
			LabelSelector: nil,
			FieldSelector: nil,
		}
	}
	if in.Spec.IPAM.LoopbackAddresses.LabelSelector == nil {
		in.Spec.IPAM.LoopbackAddresses.LabelSelector = &metav1.LabelSelector{
			MatchLabels: map[string]string{
				constants.IPAMObjectPurposeLabel: constants.IPAMLoopbackPurpose,
			},
		}
	}
	if in.Spec.IPAM.LoopbackAddresses.FieldSelector == nil {
		in.Spec.IPAM.LoopbackAddresses.FieldSelector = &FieldSelectorSpec{
			LabelKey: ptr.To(constants.IPAMObjectOwnerLabel),
			FieldRef: &v1.ObjectFieldSelector{
				APIVersion: constants.APIVersion,
				FieldPath:  constants.DefaultIPAMFieldRef,
			},
		}
	}

	// defaulting address families
	if in.Spec.IPAM.AddressFamily == nil {
		in.Spec.IPAM.AddressFamily = &AddressFamiliesMap{
			IPv4: ptr.To(true),
			IPv6: ptr.To(false),
		}
	}
}

func (in *SwitchConfig) setDefaultPortParams() {
	if in.Spec.PortsDefaults.Lanes == nil {
		in.Spec.PortsDefaults.SetLanes(4)
	}
	if in.Spec.PortsDefaults.MTU == nil {
		in.Spec.PortsDefaults.SetMTU(9100)
	}
	if in.Spec.PortsDefaults.IPv4MaskLength == nil {
		in.Spec.PortsDefaults.SetIPv4MaskLength(30)
	}
	if in.Spec.PortsDefaults.IPv6Prefix == nil {
		in.Spec.PortsDefaults.SetIPv6Prefix(127)
	}
	if in.Spec.PortsDefaults.FEC == nil {
		in.Spec.PortsDefaults.SetFEC(constants.FECNone)
	}
	if in.Spec.PortsDefaults.State == nil {
		in.Spec.PortsDefaults.SetState(constants.NICUp)
	}
}
