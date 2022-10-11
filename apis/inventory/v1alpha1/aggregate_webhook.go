// /*
// Copyright (c) 2021 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// */

package v1alpha1

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
)

// log is for logging in this package.
var aggregatelog = logf.Log.WithName("aggregate-resource")

func (in *Aggregate) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(in).
		Complete()
}

//nolint:lll
// +kubebuilder:webhook:path=/validate-machine-onmetal-de-v1alpha1-aggregate,mutating=false,failurePolicy=fail,sideEffects=None,groups=machine.onmetal.de,resources=aggregates,verbs=create;update,versions=v1alpha1,name=vaggregate.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &Aggregate{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (in *Aggregate) ValidateCreate() error {
	aggregatelog.Info("validate create", "name", in.Name)
	return in.validate()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (in *Aggregate) ValidateUpdate(old runtime.Object) error {
	aggregatelog.Info("validate update", "name", in.Name)
	return in.validate()
}

func (in *Aggregate) ValidateDelete() error {
	aggregatelog.Info("validate delete", "name", in.Name)
	return nil
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
