// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

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
	IPEventReason  = "IpUpdated"
	IPEventMessage = "IP object %s changed"
)

type IPTracker struct {
	client.Client

	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=metal.ironcore.dev,resources=networkswitches,verbs=get;list;watch
// +kubebuilder:rbac:groups=ipam.onmetal.de,resources=ips,verbs=get;list;watch
// +kubebuilder:rbac:groups=ipam.onmetal.de,resources=ips/status,verbs=get
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

func (r *IPTracker) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	obj := &ipamv1alpha1.IP{}
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	owner := obj.Labels[constants.IPAMObjectOwnerLabel]
	networkSwitch := &metalv1alpha4.NetworkSwitch{}
	if err := r.Get(ctx, types.NamespacedName{Name: owner, Namespace: req.Namespace}, networkSwitch); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	r.Recorder.Event(networkSwitch, v1.EventTypeNormal, IPEventReason, fmt.Sprintf(IPEventMessage, req.Name))
	return ctrl.Result{}, nil
}

func (r *IPTracker) SetupWithManager(mgr ctrl.Manager) error {
	labelPredicate, err := predicate.LabelSelectorPredicate(metav1.LabelSelector{
		MatchLabels: map[string]string{constants.IPAMObjectPurposeLabel: constants.IPAMLoopbackPurpose},
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
		For(&ipamv1alpha1.IP{}).
		WithEventFilter(predicate.And(labelPredicate, eventPredicate)).
		Complete(r)
}

func (r *IPTracker) setupPredicates() predicate.Predicate {
	return predicate.Funcs{
		DeleteFunc: r.deleteHandler,
		UpdateFunc: r.updateHandler,
	}
}

func (r *IPTracker) updateHandler(e event.UpdateEvent) bool {
	ip, ok := e.ObjectNew.(*ipamv1alpha1.IP)
	if !ok {
		return false
	}
	return ip.Status.State == ipamv1alpha1.CFinishedIPState
}

func (r *IPTracker) deleteHandler(e event.DeleteEvent) bool {
	ip, ok := e.Object.(*ipamv1alpha1.IP)
	if !ok {
		return false
	}
	if e.DeleteStateUnknown {
		return true
	}
	return ip.Status.State == ipamv1alpha1.CFinishedIPState
}
