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

package controllers

import (
	"context"
	"k8s.io/client-go/tools/record"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	inventoriesv1alpha1 "github.com/onmetal/k8s-inventory/api/v1alpha1"
	switchv1alpha1 "github.com/onmetal/switch-operator/api/v1alpha1"
)

const Switch = "Switch"

// Reconciler reconciles a Switch object
type Reconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=switch.onmetal.de.onmetal.de,resources=switches,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=switch.onmetal.de.onmetal.de,resources=switches/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=switch.onmetal.de.onmetal.de,resources=switches/finalizers,verbs=update
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=inventories,verbs=get;list;watch;update;patch
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=inventories/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=events,verbs=create;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Switch object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("switch", req.NamespacedName)

	// your logic here:
	// get list of inventories with host type == Switch
	switches := &switchv1alpha1.SwitchList{}
	inventories := getSwitchesInventories()
	for _, inv := range inventories {
		switchRes, exists := switchResourceExists(inv.Spec.System.SerialNumber, switches)
		if !exists {
			r.Log.Info("switch resource does not exists")
			preparedSwitch := getPreparedSwitch(switchRes, inv)
			if err := r.Client.Create(ctx, preparedSwitch); err != nil {
				r.Log.Error(err, "failed to create switch resource")
				return ctrl.Result{}, err
			}
		}
	}
	//todo: update connections

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&switchv1alpha1.Switch{}).
		WithEventFilter(r.constructPredicates()).
		Complete(r)
}

func (r *Reconciler) constructPredicates() predicate.Predicate {
	return predicate.Funcs{
		UpdateFunc: onUpdate,
	}
}

func onUpdate(e event.UpdateEvent) bool {
	oldObject := e.ObjectOld.(*switchv1alpha1.Switch)
	newObject := e.ObjectNew.(*switchv1alpha1.Switch)
	return !reflect.DeepEqual(oldObject.Spec, newObject.Spec)
}

func getSwitchesInventories() []*inventoriesv1alpha1.Inventory {
	result := make([]*inventoriesv1alpha1.Inventory, 0)
	inventories := inventoriesv1alpha1.InventoryList{}
	for _, inv := range inventories.Items {
		if inv.Spec.Host.Type == Switch {
			result = append(result, &inv)
		}
	}
	return result
}

func switchResourceExists(name string, switches *switchv1alpha1.SwitchList) (*switchv1alpha1.Switch, bool) {
	for _, switchRes := range switches.Items {
		if switchRes.ObjectMeta.Name == name {
			return &switchRes, true
		}
	}
	return &switchv1alpha1.Switch{}, false
}

func getPreparedSwitch(sw *switchv1alpha1.Switch, inv *inventoriesv1alpha1.Inventory) *switchv1alpha1.Switch {
	sw.ObjectMeta.Name = inv.Spec.System.SerialNumber
	sw.Spec.ID = inv.Spec.System.SerialNumber // in future will be refactored to use UUID
	sw.Spec.Ports = inv.Spec.NICs.Count
	return sw
}
