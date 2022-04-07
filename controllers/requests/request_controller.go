/*
Copyright 2022.

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

	"github.com/go-logr/logr"
	requestv1alpha1 "github.com/onmetal/metal-api/apis/request/v1alpha1"
	metalerr "github.com/onmetal/metal-api/pkg/errors"
	"github.com/onmetal/metal-api/pkg/scheduler"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// RequestReconciler reconciles a Request object
type RequestReconciler struct {
	ctrlclient.Client

	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// SetupWithManager sets up the controller with the Manager.
func (r *RequestReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&requestv1alpha1.Request{}).
		WithEventFilter(r.constructPredicates()).
		Complete(r)
}

func (r *RequestReconciler) constructPredicates() predicate.Predicate {
	return predicate.Funcs{
		UpdateFunc: r.onUpdate,
		DeleteFunc: r.onDelete,
	}
}

//+kubebuilder:rbac:groups=machine.onmetal.de,resources=requests,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=requests/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=requests/finalizers,verbs=update

func (r *RequestReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := r.Log.WithValues("metal-request", req.NamespacedName)

	request, err := r.GetRequest(ctx, req)
	if err != nil {
		return ctrl.Result{}, err
	}

	s := r.newScheduler(ctx, request.Spec.Kind, reqLogger)

	if err := s.Schedule(request); err != nil {
		if metalerr.IsNotFound(err) {
			reqLogger.Info("no objects for reservation found", "error", err)
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *RequestReconciler) GetRequest(ctx context.Context, req ctrl.Request) (*requestv1alpha1.Request, error) {
	metalRequest := &requestv1alpha1.Request{}
	if err := r.Client.Get(ctx,
		types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, metalRequest); err != nil {
		return metalRequest, err
	}
	return metalRequest, nil
}

func (r *RequestReconciler) newScheduler(ctx context.Context,
	kind requestv1alpha1.RequestKind, reqLogger logr.Logger) scheduler.Scheduler {
	switch kind {
	case requestv1alpha1.Machine:
		return scheduler.NewMachine(ctx, r.Client, reqLogger, r.Recorder)
	default:
		return scheduler.NewMachine(ctx, r.Client, reqLogger, r.Recorder)
	}
}

func (r *RequestReconciler) onUpdate(e event.UpdateEvent) bool {
	newObj, ok := e.ObjectNew.(*requestv1alpha1.Request)
	if !ok {
		r.Log.Info("request delete event cast failed")
		return false
	}
	if newObj.Status.State == requestv1alpha1.Reserved {
		return false
	}
	return true
}

func (r *RequestReconciler) onDelete(e event.DeleteEvent) bool {
	obj, ok := e.Object.(*requestv1alpha1.Request)
	if !ok {
		r.Log.Info("request delete event cast failed")
		return false
	}
	ctx := context.Background()

	s := r.newScheduler(ctx, obj.Spec.Kind, r.Log)

	if err := s.DeleteScheduling(obj); err != nil {
		r.Log.Info("scheduling deletion unsuccessful", "error", err)
		return false
	}

	return false
}
