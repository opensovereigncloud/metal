/*
Copyright (c) 2021 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved

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
	machinerr "github.com/onmetal/metal-api/internal/errors"
	"github.com/onmetal/metal-api/internal/switches"
	switchv1alpha1 "github.com/onmetal/switch-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

type SwitchReconciler struct {
	client.Client

	Log    logr.Logger
	Scheme *runtime.Scheme
}

// SetupWithManager sets up the controller with the Manager.
func (r *SwitchReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&switchv1alpha1.Switch{}).
		WithEventFilter(r.constructPredicates()).
		Complete(r)
}

func (r *SwitchReconciler) constructPredicates() predicate.Predicate {
	return predicate.Funcs{
		DeleteFunc: func(deleteEvent event.DeleteEvent) bool { return false },
	}
}

func (r *SwitchReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := r.Log.WithValues("switch", req.NamespacedName)

	sw, err := switches.New(ctx, r.Client, reqLogger, req)
	if err != nil {
		return machinerr.GetResultForError(reqLogger, err)
	}

	if err := sw.UpdateLocation(); err != nil {
		return machinerr.GetResultForError(reqLogger, err)
	}

	return machinerr.GetResultForError(reqLogger, nil)
}
