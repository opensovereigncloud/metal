package v1alpha1

import (
	"fmt"
	"strings"

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

// +kubebuilder:webhook:path=/validate-machine-onmetal-de-v1alpha1-aggregate,mutating=false,failurePolicy=fail,sideEffects=None,groups=machine.onmetal.de,resources=aggregates,verbs=create;update,versions=v1alpha1,name=vaggregate.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &Aggregate{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (in *Aggregate) ValidateCreate() error {
	aggregatelog.Info("validate create", "name", in.Name)
	return in.validate()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
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

func accountPath(s map[string]interface{}, tokenizedPath []string) error {
	lastTokenIdx := len(tokenizedPath) - 1
	var prevVal interface{} = s
	for idx, token := range tokenizedPath {
		theMap, ok := prevVal.(map[string]interface{})
		// if previous value is empty struct, but there are still tokens
		// then parent path is used to set value, and it is not possible
		// to set value in child path
		if !ok {
			return errors.Errorf("can not use path %s to set value, as parent path %s already used to set value", strings.Join(tokenizedPath, "."), strings.Join(tokenizedPath[:idx+1], "."))
		}

		currVal, ok := theMap[token]
		// if value is not set
		if !ok {
			// if it is the last token, then set empty struct
			if idx == lastTokenIdx {
				theMap[token] = struct{}{}
				// otherwise create a map
			} else {
				theMap[token] = make(map[string]interface{})
			}
			// if value is set
		} else {
			_, ok := currVal.(map[string]interface{})
			// if value is not map and it is the last token,
			// then there is a duplicate path
			if !ok && idx == lastTokenIdx {
				return errors.Errorf("duplicate path %s", strings.Join(tokenizedPath, "."))
			}
			// if it is map and it is the last token,
			// then there is a child path used to set value
			if ok && idx == lastTokenIdx {
				return errors.Errorf("can not use path %s to set value, as there is a child path", strings.Join(tokenizedPath, "."))
			}
			// if it is not a map and it is not the last token,
			// then there is parent path used to set value
			if !ok && idx != lastTokenIdx {
				return errors.Errorf("can not use path %s to set value, as parent path %s already used to set value", strings.Join(tokenizedPath, "."), strings.Join(tokenizedPath[:idx+1], "."))
			}
			// if it is a map and it is not the last token,
			// then continue
		}
		prevVal = theMap[token]
	}

	return nil
}
