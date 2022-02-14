/*
Copyright 2021 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved

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

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var switchLog = logf.Log.WithName("switch-resource")

func (in *Switch) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(in).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-switch-onmetal-de-v1alpha1-switch,mutating=true,failurePolicy=fail,sideEffects=None,groups=switch.onmetal.de,resources=switches,verbs=create,versions=v1alpha1,name=mswitch.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &Switch{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (in *Switch) Default() {
	switchLog.Info("default", "name", in.Name)
	if in.Labels == nil {
		in.Labels = map[string]string{}
	}
	if _, ok := in.Labels[LabelChassisId]; !ok {
		if in.Spec.Chassis.ChassisID != CEmptyString {
			in.Labels[LabelChassisId] = MacToLabel(in.Spec.Chassis.ChassisID)
		}
	}
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-switch-onmetal-de-v1alpha1-switch,mutating=false,failurePolicy=fail,sideEffects=None,groups=switch.onmetal.de,resources=switches,verbs=create;update,versions=v1alpha1,name=vswitch.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &Switch{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (in *Switch) ValidateCreate() error {
	switchLog.Info("validate create", "name", in.Name)
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (in *Switch) ValidateUpdate(old runtime.Object) error {
	switchLog.Info("validate update", "name", in.Name)
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (in *Switch) ValidateDelete() error {
	switchLog.Info("validate delete", "name", in.Name)
	return nil
}