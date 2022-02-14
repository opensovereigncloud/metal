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
	"reflect"

	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var switchAssignmentLog = logf.Log.WithName("switchassignment-resource")

func (in *SwitchAssignment) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(in).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-switch-onmetal-de-v1alpha1-switchassignment,mutating=true,failurePolicy=fail,sideEffects=None,groups=switch.onmetal.de,resources=switchassignments,verbs=create,versions=v1alpha1,name=mswitchassignment.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &SwitchAssignment{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (in *SwitchAssignment) Default() {
	switchAssignmentLog.Info("default", "name", in.Name)
	if in.Labels == nil {
		in.Labels = map[string]string{}
	}
	if _, ok := in.Labels[LabelChassisId]; !ok {
		in.Labels[LabelChassisId] = MacToLabel(in.Spec.ChassisID)
	}
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-switch-onmetal-de-v1alpha1-switchassignment,mutating=false,failurePolicy=fail,sideEffects=None,groups=switch.onmetal.de,resources=switchassignments,verbs=create;update,versions=v1alpha1,name=vswitchassignment.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &SwitchAssignment{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (in *SwitchAssignment) ValidateCreate() error {
	//switchAssignmentLog.Info("validate create", "name", swa.Name)
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (in *SwitchAssignment) ValidateUpdate(old runtime.Object) error {
	switchAssignmentLog.Info("validate update", "name", in.Name)

	oldSwitchAssignment, ok := old.(*SwitchAssignment)
	if !ok {
		return errors.New("failed to cast previous object version to SwitchAssignment resource type")
	}

	var allErrors field.ErrorList
	if oldSwitchAssignment.Spec.ChassisID != in.Spec.ChassisID {
		allErrors = append(allErrors, field.Invalid(field.NewPath("spec.chassisId"), in.Spec.ChassisID, "Chassis ID change disallowed"))
	}
	if !reflect.DeepEqual(oldSwitchAssignment.Spec.Region, in.Spec.Region) {
		allErrors = append(allErrors, field.Invalid(field.NewPath("spec.region"), in.Spec.Region, "Region change disallowed"))
	}

	if len(allErrors) > 0 {
		return apierrors.NewInvalid(schema.GroupKind{
			Group: GroupVersion.Group,
			Kind:  "SwitchAssignment",
		}, in.Name, allErrors)
	}

	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (in *SwitchAssignment) ValidateDelete() error {
	//switchAssignmentLog.Info("validate delete", "name", swa.Name)

	return nil
}