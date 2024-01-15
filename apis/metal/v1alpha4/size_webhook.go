/*
Copyright 2021.

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

package v1alpha4

import (
	"strings"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

const (
	CStatusFieldName   = "status"
	CComputedFieldName = "computed"
)

// log is for logging in this package.
var sizelog = logf.Log.WithName("size-resource")

func (in *Size) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(in).
		Complete()
}

// +kubebuilder:webhook:path=/validate-metal-ironcore-dev-v1alpha4-size,mutating=false,failurePolicy=fail,sideEffects=None,groups=metal.ironcore.dev,resources=sizes,verbs=create;update,versions=v1alpha4,name=vsize.v1alpha4.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &Size{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (in *Size) ValidateCreate() (admission.Warnings, error) {
	var warnings admission.Warnings

	sizelog.Info("validate create", "name", in.Name)
	return warnings, in.validate()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (in *Size) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	var warnings admission.Warnings

	sizelog.Info("validate update", "name", in.Name)
	return warnings, in.validate()
}

var CDummyInventorySpec = getDummyInventoryForValidation()

func (in *Size) validate() error {
	ops := make(map[string]int)
	errs := make([]string, 0)

	for _, c := range in.Spec.Constraints {
		pathString := c.Path.String()
		op, ok := ops[pathString]
		if !ok {
			op = 0
		}
		op++
		ops[pathString] = op

		// It is not possible to check on aggregate path existence,
		// as it may be applied dynamically
		// So, instead, we will check that path has a proper prefixs,
		// has at least 4 segments, and matches the schema,
		// i.e. status.computed.{aggregate-name}.{value-key}

		// Moreover, jsonpath doesn't allow to access its parser.
		// Means, it is required to process path in different ways for each case.

		tokens := c.Path.Tokenize()
		if !(len(tokens) >= 4 && tokens[0] == CStatusFieldName && tokens[1] == CComputedFieldName) {
			jp, err := c.Path.ToK8sJSONPath()
			jp.AllowMissingKeys(false)
			if err != nil {
				errs = append(errs, errors.Wrap(err, "unable to parse JSONPath").Error())
			} else if _, err := jp.FindResults(CDummyInventorySpec); err != nil {
				errs = append(errs, errors.Wrap(err, "unable to find results with path").Error())
			}
		}

		if op == 2 {
			err := errors.Errorf("multiple constraints found for field %s", c.Path)
			errs = append(errs, err.Error())
		}

		if c.empty() {
			err := errors.Errorf("constraint for %s does not contains conditions", c.Path)
			errs = append(errs, err.Error())
		}

		if c.hasAggregateAndLiterals() {
			err := errors.New("aggregates can be validated only against numeric values")
			errs = append(errs, err.Error())
		}

		if c.eqAndNeq() {
			err := errors.Errorf("constraint for %s contains both eq and neq conditions", c.Path)
			errs = append(errs, err.Error())
		}

		if c.inclusiveAndExclusive() {
			err := errors.Errorf("constraint for %s contains both gt and gte or lt and lte conditions", c.Path)
			errs = append(errs, err.Error())
		}

		if c.borderAndEq() {
			err := errors.Errorf("constraint for %s contains both gt/gte/lt/lte and eq/neq conditions", c.Path)
			errs = append(errs, err.Error())
		}

		if c.wrongInterval() {
			err := errors.Errorf("constraint for %s lower border is greater than upper border", c.Path)
			errs = append(errs, err.Error())
		}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}

	return nil
}

func (in *ConstraintSpec) hasAggregateAndLiterals() bool {
	return in.Aggregate != "" &&
		(in.Equal != nil && in.Equal.Literal != nil ||
			in.NotEqual != nil && in.NotEqual.Literal != nil)
}

func (in *ConstraintSpec) empty() bool {
	return in.Equal == nil &&
		in.NotEqual == nil &&
		in.GreaterThan == nil &&
		in.GreaterThanOrEqual == nil &&
		in.LessThan == nil &&
		in.LessThanOrEqual == nil
}

func (in *ConstraintSpec) eqAndNeq() bool {
	return in.Equal != nil &&
		in.NotEqual != nil
}

func (in *ConstraintSpec) inclusiveAndExclusive() bool {
	return in.GreaterThan != nil && in.GreaterThanOrEqual != nil ||
		in.LessThan != nil && in.LessThanOrEqual != nil
}

func (in *ConstraintSpec) borderAndEq() bool {
	return (in.GreaterThan != nil || in.GreaterThanOrEqual != nil || in.LessThan != nil || in.LessThanOrEqual != nil) &&
		(in.Equal != nil || in.NotEqual != nil)
}

func (in *ConstraintSpec) wrongInterval() bool {
	var upper *resource.Quantity
	var lower *resource.Quantity

	if in.LessThanOrEqual != nil {
		upper = in.LessThanOrEqual
	}
	if in.LessThan != nil {
		upper = in.LessThan
	}
	if in.GreaterThanOrEqual != nil {
		lower = in.GreaterThanOrEqual
	}
	if in.GreaterThan != nil {
		lower = in.GreaterThan
	}

	if upper == nil || lower == nil {
		return false
	}

	if lower.Cmp(*upper) < 0 {
		return false
	}

	return true
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (in *Size) ValidateDelete() (admission.Warnings, error) {
	var warnings admission.Warnings

	sizelog.Info("validate delete", "name", in.Name)
	return warnings, nil
}
