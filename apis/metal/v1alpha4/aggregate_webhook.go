// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package v1alpha4

import (
	"fmt"

	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var aggregatelog = logf.Log.WithName("aggregate-resource")

func (in *Aggregate) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(in).
		Complete()
}

// +kubebuilder:webhook:path=/validate-metal-ironcore-dev-v1alpha4-aggregate,mutating=false,failurePolicy=fail,sideEffects=None,groups=metal.ironcore.dev,resources=aggregates,verbs=create;update,versions=v1alpha4,name=vaggregate.v1alpha4.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &Aggregate{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (in *Aggregate) ValidateCreate() (admission.Warnings, error) {
	var warnings admission.Warnings

	aggregatelog.Info("validate create", "name", in.Name)
	return warnings, in.validate()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (in *Aggregate) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	var warnings admission.Warnings

	aggregatelog.Info("validate update", "name", in.Name)
	return warnings, in.validate()
}

func (in *Aggregate) ValidateDelete() (admission.Warnings, error) {
	var warnings admission.Warnings

	aggregatelog.Info("validate delete", "name", in.Name)
	return warnings, nil
}

func (in *Aggregate) validate() error {
	var allErrs field.ErrorList

	pathAccounter := make(map[string]interface{})
	for i, agg := range in.Spec.Aggregates {
		srcJp, err := agg.SourcePath.ToK8sJSONPath()
		if err != nil {
			msg := errors.Wrap(err, "unable to convert to k8s JSONPath").Error()
			path := field.NewPath(fmt.Sprintf("spec.aggregates[%d].sourcePath", i))
			allErrs = append(allErrs, field.Invalid(path, agg.SourcePath, msg))
		}

		srcJp.AllowMissingKeys(false)
		if _, err := srcJp.FindResults(CDummyInventorySpec); err != nil {
			msg := errors.Wrap(err, "unable to find results with path").Error()
			path := field.NewPath(fmt.Sprintf("spec.aggregates[%d].sourcePath", i))
			allErrs = append(allErrs, field.Invalid(path, agg.SourcePath, msg))
		}

		if err := accountPath(pathAccounter, agg.TargetPath.Tokenize()); err != nil {
			msg := errors.Wrap(err, "unable to insert target path").Error()
			path := field.NewPath(fmt.Sprintf("spec.aggregates[%d].targetPath", i))
			allErrs = append(allErrs, field.Invalid(path, agg.SourcePath, msg))
		}
	}

	if len(allErrs) > 0 {
		gvk := in.GroupVersionKind()
		gk := schema.GroupKind{
			Group: gvk.Group,
			Kind:  gvk.Kind,
		}
		return apierrors.NewInvalid(gk, in.Name, allErrs)
	}

	return nil
}
