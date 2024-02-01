/*
Copyright (c) 2024 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved

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
	"fmt"

	"github.com/go-logr/logr"
	ipamv1alpha1 "github.com/onmetal/ipam/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	metalv1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"
	"github.com/ironcore-dev/metal/pkg/constants"
)

const (
	SubnetEventReason  = "IpUpdated"
	SubnetEventMessage = "IP object %s changed"
)

type SubnetTracker struct {
	client.Client

	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=metal.ironcore.dev,resources=networkswitches,verbs=get;list;watch
// +kubebuilder:rbac:groups=ipam.onmetal.de,resources=subnets,verbs=get;list;watch
// +kubebuilder:rbac:groups=ipam.onmetal.de,resources=subnets/status,verbs=get
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

func (r *SubnetTracker) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	obj := &ipamv1alpha1.Subnet{}
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	owner := obj.Labels[constants.IPAMObjectOwnerLabel]
	networkSwitch := &metalv1alpha4.NetworkSwitch{}
	if err := r.Get(ctx, types.NamespacedName{Name: owner, Namespace: req.Namespace}, networkSwitch); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	r.Recorder.Event(networkSwitch, v1.EventTypeNormal, SubnetEventReason, fmt.Sprintf(SubnetEventMessage, req.Name))
	return ctrl.Result{}, nil
}

func (r *SubnetTracker) SetupWithManager(mgr ctrl.Manager) error {
	labelPredicate, err := predicate.LabelSelectorPredicate(metav1.LabelSelector{
		MatchLabels: map[string]string{constants.IPAMObjectPurposeLabel: constants.IPAMSouthSubnetPurpose},
		MatchExpressions: []metav1.LabelSelectorRequirement{{
			Key:      constants.IPAMObjectOwnerLabel,
			Operator: metav1.LabelSelectorOpExists,
		}},
	})
	if err != nil {
		r.Log.Error(err, "failed to setup predicates")
	}
	eventPredicate := r.setupPredicates()

	return ctrl.NewControllerManagedBy(mgr).
		For(&ipamv1alpha1.Subnet{}).
		WithEventFilter(predicate.And(labelPredicate, eventPredicate)).
		Complete(r)
}

func (r *SubnetTracker) setupPredicates() predicate.Predicate {
	return predicate.Funcs{
		DeleteFunc: r.deleteHandler,
		UpdateFunc: r.updateHandler,
	}
}

func (r *SubnetTracker) updateHandler(e event.UpdateEvent) bool {
	subnet, ok := e.ObjectNew.(*ipamv1alpha1.Subnet)
	if !ok {
		return false
	}
	return subnet.Status.State == ipamv1alpha1.CFinishedSubnetState
}

func (r *SubnetTracker) deleteHandler(e event.DeleteEvent) bool {
	subnet, ok := e.Object.(*ipamv1alpha1.Subnet)
	if !ok {
		return false
	}
	if e.DeleteStateUnknown {
		return true
	}
	return subnet.Status.State == ipamv1alpha1.CFinishedSubnetState
}
