/*
Copyright (c) 2022 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved

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

package v1beta1

import (
	"context"
	"reflect"
	"time"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	switchv1beta1 "github.com/onmetal/metal-api/apis/switch/v1beta1"
	"github.com/onmetal/metal-api/pkg/constants"
	"github.com/onmetal/metal-api/pkg/errors"
	switchespkg "github.com/onmetal/metal-api/pkg/switches"
)

// SwConfigReconciler reconciles SwitchConfig object corresponding.
type SwConfigReconciler struct {
	client.Client

	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=switch.onmetal.de,resources=switchconfigs,verbs=get;list;watch
// +kubebuilder:rbac:groups=switch.onmetal.de,resources=switches,verbs=get;list
// +kubebuilder:rbac:groups=switch.onmetal.de,resources=switches/status,verbs=get;update;patch

func (r *SwConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	nestedCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	obj := &switchv1beta1.SwitchConfig{}
	if err := r.Get(nestedCtx, req.NamespacedName, obj); err != nil {
		switch {
		case apierrors.IsNotFound(err):
			r.Log.Info("requested SwitchConfig object not found", "name", req.NamespacedName)
		default:
			r.Log.Info("failed to get requested SwitchConfig object", "name", req.NamespacedName, "error", err)
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	result, err := r.updateSwitches(nestedCtx, obj)
	return result, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *SwConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	changesTrackPredicate := predicate.Funcs{
		UpdateFunc: checkSwitchConfigUpdate,
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&switchv1beta1.SwitchConfig{}).
		WithEventFilter(predicate.And(changesTrackPredicate)).
		Complete(r)
}

func checkSwitchConfigUpdate(e event.UpdateEvent) bool {
	objOld, okOld := e.ObjectOld.(*switchv1beta1.SwitchConfig)
	objNew, okNew := e.ObjectNew.(*switchv1beta1.SwitchConfig)
	if !okOld || !okNew {
		return false
	}
	ipamConfigChanged := !reflect.DeepEqual(objOld.Spec.IPAM, objNew.Spec.IPAM)
	portConfigChanged := !reflect.DeepEqual(objOld.Spec.PortsDefaults, objNew.Spec.PortsDefaults)
	labelsChanged := !reflect.DeepEqual(objOld.GetLabels(), objNew.GetLabels())
	return ipamConfigChanged || portConfigChanged || labelsChanged
}

func (r *SwConfigReconciler) updateSwitches(ctx context.Context, obj *switchv1beta1.SwitchConfig) (ctrl.Result, error) {
	result := ctrl.Result{}
	relatedSwitches := &switchv1beta1.SwitchList{}
	if err := r.List(ctx, relatedSwitches); err != nil {
		result.Requeue = true
		return result, err
	}
	labelsSet := labels.Set(obj.Labels)
	for _, item := range relatedSwitches.Items {
		selector, _ := metav1.LabelSelectorAsSelector(item.Spec.ConfigSelector)
		if !selector.Matches(labelsSet) {
			continue
		}
		item.SetCondition(constants.ConditionPortParametersOK, false).
			SetReason(errors.ErrorReasonDataOutdated.String())
		item.ManagedFields = make([]metav1.ManagedFieldsEntry, 0)
		if err := r.Status().Patch(ctx, &item, client.Apply, switchespkg.PatchOpts); err != nil {
			result.RequeueAfter = time.Second * 5
			return result, err
		}
	}
	return ctrl.Result{}, nil
}
