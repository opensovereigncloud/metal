/*
Copyright 2021 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved

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
	"net"

	"github.com/go-logr/logr"
	subnetv1alpha1 "github.com/onmetal/ipam/api/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	switchv1alpha1 "github.com/onmetal/switch-operator/api/v1alpha1"
)

type Background struct {
	switches   *switchv1alpha1.SwitchList
	assignment *switchv1alpha1.SwitchAssignment
}

// SwitchReconciler reconciles a Switch object
type SwitchReconciler struct {
	client.Client
	Log        logr.Logger
	Scheme     *runtime.Scheme
	Recorder   record.EventRecorder
	Background *Background
}

//+kubebuilder:rbac:groups=switch.onmetal.de,resources=switches,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=switch.onmetal.de,resources=switches/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=switch.onmetal.de,resources=switches/finalizers,verbs=update
//+kubebuilder:rbac:groups=switch.onmetal.de,resources=events,verbs=create;patch
//+kubebuilder:rbac:groups=ipam.onmetal.de,resources=subnets,verbs=get;list;watch;update
//+kubebuilder:rbac:groups=ipam.onmetal.de,resources=subnets/status,verbs=get;update;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *SwitchReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("name", req.NamespacedName)
	res := &switchv1alpha1.Switch{}
	if err := r.Get(ctx, req.NamespacedName, res); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("requested switch resource not found")
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
		log.Error(err, "failed to get switch resource")
		return ctrl.Result{}, err
	}

	if err := r.prepareBackground(ctx, res); err != nil {
		return ctrl.Result{}, err
	}

	processor := switchProcessor{}
	if res.DeletionTimestamp != nil {
		processor.startPoint = &deletionStep{}
	} else {
		processor.startPoint = &preparationStep{}
	}
	return processor.launch(res, r, ctx)
}

// SetupWithManager sets up the controller with the Manager.
func (r *SwitchReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&switchv1alpha1.Switch{}).
		Watches(&source.Kind{Type: &switchv1alpha1.SwitchAssignment{}}, handler.Funcs{
			UpdateFunc: r.handleAssignmentUpdate(),
		}).
		Watches(&source.Kind{Type: &switchv1alpha1.Switch{}}, handler.Funcs{
			UpdateFunc: r.handleSwitchUpdate(mgr.GetScheme(), &switchv1alpha1.SwitchList{}),
		}).
		Watches(&source.Kind{Type: &subnetv1alpha1.Subnet{}}, handler.Funcs{
			UpdateFunc: r.handleSubnetUpdate(mgr.GetScheme(), &switchv1alpha1.SwitchList{}),
		}).
		Complete(r)
}

func (r *SwitchReconciler) handleAssignmentUpdate() func(event.UpdateEvent, workqueue.RateLimitingInterface) {
	return func(e event.UpdateEvent, q workqueue.RateLimitingInterface) {
		oldRes := e.ObjectOld.(*switchv1alpha1.SwitchAssignment)
		updRes := e.ObjectNew.(*switchv1alpha1.SwitchAssignment)
		if err := r.enqueueOnAssignmentUpdate(context.Background(), oldRes, updRes, q); err != nil {
			r.Log.Error(err, "failed to trigger reconciliation on switchAssignment resource update")
		}
	}
}

func (r *SwitchReconciler) handleSwitchUpdate(scheme *runtime.Scheme, ro runtime.Object) func(event.UpdateEvent, workqueue.RateLimitingInterface) {
	return func(e event.UpdateEvent, q workqueue.RateLimitingInterface) {
		if err := r.enqueueOnSwitchUpdate(context.Background(), scheme, q, ro); err != nil {
			r.Log.Error(err, "failed to trigger reconciliation on switch resource update")
		}
	}
}

func (r *SwitchReconciler) handleSubnetUpdate(scheme *runtime.Scheme, ro runtime.Object) func(event.UpdateEvent, workqueue.RateLimitingInterface) {
	return func(e event.UpdateEvent, q workqueue.RateLimitingInterface) {
		updRes := e.ObjectNew.(*subnetv1alpha1.Subnet)
		if err := r.enqueueOnSubnetUpdate(context.Background(), updRes, scheme, q, ro); err != nil {
			r.Log.Error(err, "failed to trigger reconciliation on subnet resource update")
		}
	}
}

func (r *SwitchReconciler) enqueueOnAssignmentUpdate(
	ctx context.Context,
	oldRes *switchv1alpha1.SwitchAssignment,
	updRes *switchv1alpha1.SwitchAssignment,
	q workqueue.RateLimitingInterface) error {

	if !controllerutil.ContainsFinalizer(oldRes, switchv1alpha1.CSwitchAssignmentFinalizer) {
		return nil
	}
	if updRes.Status.State == switchv1alpha1.StateFinished {
		return nil
	}
	list := &switchv1alpha1.SwitchList{}
	opts, err := updRes.GetListFilter()
	if err != nil {
		return err
	}
	if err := r.List(ctx, list, opts); err != nil {
		return err
	}
	if len(list.Items) == 0 {
		return nil
	}
	obj := &list.Items[0]
	q.Add(reconcile.Request{NamespacedName: obj.NamespacedName()})
	return nil
}

func (r *SwitchReconciler) enqueueOnSwitchUpdate(
	ctx context.Context,
	scheme *runtime.Scheme,
	q workqueue.RateLimitingInterface,
	ro runtime.Object) error {
	context.Background()

	list := &unstructured.UnstructuredList{}
	gvk, err := apiutil.GVKForObject(ro, scheme)
	if err != nil {
		r.Log.Error(err, "unable to get gvk")
		return err
	}
	list.SetGroupVersionKind(gvk)
	if err := r.List(ctx, list); err != nil {
		r.Log.Error(err, "unable to get list of items")
		return err
	}
	for _, item := range list.Items {
		q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
			Namespace: item.GetNamespace(),
			Name:      item.GetName(),
		}})
	}
	return nil
}

func (r *SwitchReconciler) enqueueOnSubnetUpdate(
	ctx context.Context,
	src *subnetv1alpha1.Subnet,
	scheme *runtime.Scheme,
	q workqueue.RateLimitingInterface,
	ro runtime.Object) error {

	if src.Status.State != subnetv1alpha1.CFinishedSubnetState {
		return nil
	}
	list := &unstructured.UnstructuredList{}
	gvk, err := apiutil.GVKForObject(ro, scheme)
	if err != nil {
		r.Log.Error(err, "unable to get gvk")
		return err
	}
	list.SetGroupVersionKind(gvk)
	if err := r.List(ctx, list); err != nil {
		r.Log.Error(err, "unable to get list of items")
		return err
	}
	for _, item := range list.Items {
		q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
			Namespace: item.GetNamespace(),
			Name:      item.GetName(),
		}})
	}
	return nil
}

func (r *SwitchReconciler) finalize(ctx context.Context, res *switchv1alpha1.Switch) error {
	if controllerutil.ContainsFinalizer(res, switchv1alpha1.CSwitchFinalizer) {
		res.FlushStatusOnDelete()
		if err := r.Status().Update(ctx, res); err != nil {
			r.Log.Error(err, "failed to finalize resource", "name", res.NamespacedName())
			return err
		}

		swa, err := r.findAssignment(ctx, res)
		if err != nil {
			r.Log.Error(err, "failed to lookup for related switch assignment resource",
				"gvk", res.GroupVersionKind(), "name", res.NamespacedName())
		}
		if swa != nil {
			swa.FillStatus(switchv1alpha1.StatePending, &switchv1alpha1.LinkedSwitchSpec{})
			if err := r.Status().Update(ctx, swa); err != nil {
				r.Log.Error(err, "failed to set status on resource creation",
					"gvk", swa.GroupVersionKind(), "name", swa.NamespacedName())
			}
		}

		if res.Spec.SouthSubnetV4 != nil {
			subnet := &subnetv1alpha1.Subnet{}
			if err := r.Get(ctx, types.NamespacedName{
				Namespace: res.Spec.SouthSubnetV4.ParentSubnet.Namespace,
				Name:      res.Spec.SouthSubnetV4.ParentSubnet.Name,
			}, subnet); err != nil {
				r.Log.Error(err, "failed to get subnet resource")
			} else {
				_, network, _ := net.ParseCIDR(res.Spec.SouthSubnetV4.CIDR)
				_ = subnet.Release(&subnetv1alpha1.CIDR{Net: network})
				if err := r.Status().Update(ctx, subnet); err != nil {
					r.Log.Error(err, "failed to update subnet status on reservation release")
				}
			}
		}
		if res.Spec.SouthSubnetV6 != nil {
			subnet := &subnetv1alpha1.Subnet{}
			if err := r.Get(ctx, types.NamespacedName{
				Namespace: res.Spec.SouthSubnetV6.ParentSubnet.Namespace,
				Name:      res.Spec.SouthSubnetV6.ParentSubnet.Name,
			}, subnet); err != nil {
				r.Log.Error(err, "failed to get subnet resource")
			} else {
				_, network, _ := net.ParseCIDR(res.Spec.SouthSubnetV6.CIDR)
				_ = subnet.Release(&subnetv1alpha1.CIDR{Net: network})
				if err := r.Status().Update(ctx, subnet); err != nil {
					r.Log.Error(err, "failed to update subnet status on reservation release")
				}
			}
		}

		controllerutil.RemoveFinalizer(res, switchv1alpha1.CSwitchFinalizer)
		if err := r.Update(ctx, res); err != nil {
			r.Log.Error(err, "failed to update resource on finalizer removal",
				"gvk", res.GroupVersionKind(), "name", res.NamespacedName())
			return err
		}
	}
	return nil
}

func (r *SwitchReconciler) findAssignment(ctx context.Context, sw *switchv1alpha1.Switch) (*switchv1alpha1.SwitchAssignment, error) {
	opts, err := sw.GetListFilter()
	if err != nil {
		r.Log.Error(err, "failed to construct list options object")
		return nil, err
	}
	list := &switchv1alpha1.SwitchAssignmentList{}
	if err := r.List(ctx, list, opts); err != nil {
		r.Log.Error(err, "failed to list resources", "gvk", list.GroupVersionKind().String())
		return nil, err
	}
	if len(list.Items) == 0 {
		return nil, nil
	}
	return &list.Items[0], nil
}

func (r *SwitchReconciler) defineSubnets(
	ctx context.Context,
	sw *switchv1alpha1.Switch,
	list *switchv1alpha1.SwitchList,
	assignment *switchv1alpha1.SwitchAssignment) error {

	regions := make([]string, 0)
	zones := make([]string, 0)
	if assignment == nil {
		topLevelSwitch := list.GetTopLevelSwitch()
		if topLevelSwitch == nil {
			return nil
		}
		swa, err := r.findAssignment(ctx, topLevelSwitch)
		if err != nil {
			return err
		}
		if swa == nil {
			return nil
		}
		assignment = swa
	}
	regions = append(regions, assignment.Spec.Region)
	zones = append(zones, assignment.Spec.AvailabilityZone)
	subnets := &subnetv1alpha1.SubnetList{}
	if err := r.Client.List(ctx, subnets); err != nil {
		return err
	}
	if sw.Spec.SouthSubnetV4 == nil {
		cidr, sn, err := sw.GetSuitableSubnet(subnets, subnetv1alpha1.CIPv4SubnetType, regions, zones)
		if err != nil {
			return err
		}
		if cidr != nil && sn != nil {
			sw.Spec.SouthSubnetV4 = &switchv1alpha1.SwitchSubnetSpec{
				ParentSubnet: &switchv1alpha1.ParentSubnetSpec{Namespace: sn.Namespace, Name: sn.Name},
				CIDR:         cidr.String(),
			}
			if err := r.Status().Update(ctx, sn); err != nil {
				return err
			}
		}
	}
	if sw.Spec.SouthSubnetV6 == nil {
		cidr, sn, err := sw.GetSuitableSubnet(subnets, subnetv1alpha1.CIPv6SubnetType, regions, zones)
		if err != nil {
			return err
		}
		if cidr != nil && sn != nil {
			sw.Spec.SouthSubnetV6 = &switchv1alpha1.SwitchSubnetSpec{
				ParentSubnet: &switchv1alpha1.ParentSubnetSpec{Namespace: sn.Namespace, Name: sn.Name},
				CIDR:         cidr.String(),
			}
			if err := r.Status().Update(ctx, sn); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *SwitchReconciler) prepareBackground(ctx context.Context, sw *switchv1alpha1.Switch) error {
	if r.Background == nil {
		r.Background = &Background{
			switches:   nil,
			assignment: nil,
		}
	}

	list := &switchv1alpha1.SwitchList{}
	if err := r.List(ctx, list); err != nil {
		r.Log.Error(err, "failed to list resources", "gvk", list.GroupVersionKind().String())
		return err
	}
	r.Background.switches = list

	swa, err := r.findAssignment(ctx, sw)
	if err != nil {
		r.Log.Error(err, "failed to get related assignment resource",
			"gvk", sw.GroupVersionKind().String(),
			"name", sw.NamespacedName())
		return err
	}
	r.Background.assignment = swa
	return nil
}
