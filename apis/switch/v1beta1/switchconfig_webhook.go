/*
 * Copyright (c) 2022 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package v1beta1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	"github.com/onmetal/metal-api/internal/constants"
)

// log is for logging in this package.
var switchconfiglog = logf.Log.WithName("switchconfig-resource")

func (in *SwitchConfig) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(in).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-switch-onmetal-de-v1beta1-switchconfig,mutating=true,failurePolicy=fail,sideEffects=None,groups=switch.onmetal.de,resources=switchconfigs,verbs=create;update,versions=v1beta1,name=mswitchconfig.v1beta1.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &SwitchConfig{}

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (in *SwitchConfig) Default() {
	switchconfiglog.Info("default", "name", in.Name)
	in.setDefaultIPAMSelectors()
	in.setDefaultPortParams()
}

// +kubebuilder:webhook:path=/validate-switch-onmetal-de-v1beta1-switchconfig,mutating=false,failurePolicy=fail,sideEffects=None,groups=switch.onmetal.de,resources=switchconfigs,verbs=create;update,versions=v1beta1,name=vswitchconfig.v1beta1.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &SwitchConfig{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (in *SwitchConfig) ValidateCreate() error {
	// todo: validate if label(s) with switch type(s) exist, if type != all in.Spec.Switches is not nil and types in labels match switches selector
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (in *SwitchConfig) ValidateUpdate(_ runtime.Object) error {
	// todo: validate if label(s) with switch type(s) exist, if type != all in.Spec.Switches is not nil and types in labels match switches selector
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (in *SwitchConfig) ValidateDelete() error {
	return nil
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
			LabelKey: pointer.String(constants.IPAMObjectOwnerLabel),
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
			LabelKey: pointer.String(constants.IPAMObjectOwnerLabel),
			FieldRef: &v1.ObjectFieldSelector{
				APIVersion: constants.APIVersion,
				FieldPath:  constants.DefaultIPAMFieldRef,
			},
		}
	}

	// defaulting address families
	if in.Spec.IPAM.AddressFamily == nil {
		in.Spec.IPAM.AddressFamily = &AddressFamiliesMap{
			IPv4: pointer.Bool(true),
			IPv6: pointer.Bool(false),
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
