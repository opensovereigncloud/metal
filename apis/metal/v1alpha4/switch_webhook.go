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

package v1alpha4

import (
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/ironcore-dev/metal/pkg/constants"
)

// log is for logging in this package.
var switchlog = logf.Log.WithName("networkswitch-resource")

// validateInventoryLabels validates inventoried and inventory-ref labels with following conditions:
// - inventoried label is always used with inventory-ref label
// - if no inventoried label then inventory-ref label is not allowed
// - inventoried label can be only "false" or "true"
// - inventory-ref label can be only valid UUID.
func validateInventoryLabels(switchObject *NetworkSwitch) (err error) {
	inventoried, inventoriedOk := switchObject.Labels["metal.ironcore.dev/inventoried"]
	inventoryRef, inventoryRefOk := switchObject.Labels["metal.ironcore.dev/inventory-ref"]

	switch {
	case !inventoriedOk && !inventoryRefOk:
		return
	case !inventoriedOk && inventoryRefOk:
		return errors.New("inventory-ref label is set but inventoried label is not set")
	case inventoriedOk && !inventoryRefOk:
		return errors.New("inventoried label is set but inventory-ref label is not set")
	}

	_, err = uuid.Parse(inventoryRef)
	if err != nil {
		return fmt.Errorf("inventory-ref label must be a valid UUID, but current value is %s", inventoryRef)
	}

	if inventoried == "true" || inventoried == "false" {
		return
	}

	return fmt.Errorf("inventoried label must be set to true or false, but current value is %s", inventoried)
}

// validateOverrides validates switch interface overrides with following conditions:
// MTU, FEC, Lanes can be only set for non-north interfaces.
func validateOverrides(currentState *NetworkSwitch, newState *NetworkSwitch) (err error) {
	var errList field.ErrorList
	if newState.Spec.Interfaces == nil {
		return
	}
	if newState.Spec.Interfaces.Overrides == nil {
		return
	}
	newInterfaces := newState.Spec.Interfaces.Overrides
	currentInterfaces := currentState.Status.Interfaces

	if len(currentInterfaces) == 0 {
		switchlog.Info(
			"interface override webhook",
			"Current switch object does not contain interfaces information in its Status, looks like it was just created, skipping validation",
			currentState.Name)
		return
	}

	for _, newInterface := range newInterfaces {
		currentInterface, ok := currentInterfaces[newInterface.GetName()]
		if !ok {
			err = errors.Errorf(
				"%s interface override for update was not found in the interfaces status, probably one does not exists",
				newInterface.GetName())
			return
		}

		if currentInterface.GetDirection() == constants.DirectionSouth {
			continue
		}

		interfaceChanged := false
		switch {
		case newInterface.FEC != nil && newInterface.GetFEC() != currentInterface.GetFEC():
			interfaceChanged = true
		case newInterface.Lanes != nil && newInterface.GetLanes() != currentInterface.GetLanes():
			interfaceChanged = true
		case newInterface.MTU != nil && newInterface.GetMTU() != currentInterface.GetMTU():
			interfaceChanged = true
		}

		if interfaceChanged {
			errList = append(
				errList,
				field.Invalid(
					field.NewPath("spec.interfaces.overrides"),
					newInterface.GetName(),
					"Changing FEC, MTU, Lanes are not allowed for north interfaces"))
		}
	}

	if len(errList) > 0 {
		return apierrors.NewInvalid(
			schema.GroupKind{
				Group: SchemeGroupVersion.Group,
				Kind:  "NetworkSwitch"},
			newState.Name,
			errList)
	}

	return
}

func (in *NetworkSwitch) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(in).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-metal-ironcore-dev-v1alpha4-networkswitch,mutating=true,failurePolicy=fail,sideEffects=None,groups=metal.ironcore.dev,resources=networkswitches,verbs=create;update,versions=v1alpha4,name=mnetworkswitch.v1alpha4.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &NetworkSwitch{}

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (in *NetworkSwitch) Default() {
	switchlog.Info("default", "name", in.Name)
	in.setDefaultConfigSelector()
}

func (in *NetworkSwitch) setDefaultConfigSelector() {
	selector := in.GetConfigSelector()
	if selector == nil {
		if in.GetLayer() == 255 {
			return
		}
		layerAsString := strconv.Itoa(int(in.GetLayer()))
		in.Spec.ConfigSelector = &metav1.LabelSelector{
			MatchLabels: map[string]string{constants.SwitchConfigLayerLabel: layerAsString},
		}
		return
	}
	_, ok := selector.MatchLabels[constants.SwitchConfigLayerLabel]
	if ok && len(selector.MatchLabels) <= 1 {
		layerAsString := strconv.Itoa(int(in.GetLayer()))
		in.Spec.ConfigSelector = &metav1.LabelSelector{
			MatchLabels: map[string]string{constants.SwitchConfigLayerLabel: layerAsString},
		}
	}
	if ok && len(selector.MatchLabels) > 1 {
		delete(selector.MatchLabels, constants.SwitchConfigLayerLabel)
	}
}

// +kubebuilder:webhook:path=/validate-metal-ironcore-dev-v1alpha4-networkswitch,mutating=false,failurePolicy=fail,sideEffects=None,groups=metal.ironcore.dev,resources=networkswitches,verbs=create;update,versions=v1alpha4,name=vnetworkswitch.v1alpha4.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &NetworkSwitch{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (in *NetworkSwitch) ValidateCreate() (warnings admission.Warnings, err error) {
	switchlog.Info("validate create", "name", in.Name)

	err = validateInventoryLabels(in)
	return
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (in *NetworkSwitch) ValidateUpdate(old runtime.Object) (warnings admission.Warnings, err error) {
	switchlog.Info("validate update", "name", in.Name)

	err = validateInventoryLabels(in)
	if err != nil {
		return
	}

	currentState, ok := old.(*NetworkSwitch)
	if !ok {
		err = errors.New("failed to cast previous object version to NetworkSwitch resource type")
		return
	}

	if currentState.Spec.InventoryRef == nil && in.Spec.InventoryRef == nil {
		return
	}
	if currentState.Spec.InventoryRef == nil && in.Spec.InventoryRef != nil {
		return
	}
	if currentState.Spec.InventoryRef != nil && in.Spec.InventoryRef == nil {
		err = errors.New("cannot change inventory reference, operation denied")
		return
	}
	if currentState.GetInventoryRef() != in.GetInventoryRef() {
		err = errors.New("cannot change inventory reference, operation denied")
		return
	}

	err = validateOverrides(currentState, in)

	return
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (in *NetworkSwitch) ValidateDelete() (warnings admission.Warnings, err error) {
	switchlog.Info("validate delete", "name", in.Name)
	return
}