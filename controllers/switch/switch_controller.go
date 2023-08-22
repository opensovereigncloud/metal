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
	"crypto/md5"
	"fmt"
	"net/netip"
	"reflect"
	"strings"

	"github.com/go-logr/logr"
	ipamv1alpha1 "github.com/onmetal/ipam/api/v1alpha1"
	inventoryv1alpha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	switchv1beta1 "github.com/onmetal/metal-api/apis/switch/v1beta1"
	"github.com/onmetal/metal-api/pkg/constants"
	"github.com/onmetal/metal-api/pkg/errors"
	"github.com/onmetal/metal-api/pkg/stateproc"
	switchespkg "github.com/onmetal/metal-api/pkg/switches"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// SwitchReconciler reconciles Switch object corresponding
// to given Inventory object.
type SwitchReconciler struct {
	client.Client

	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

type SwitchClient struct {
	client.Client
	logr.Logger

	ctx                context.Context
	updateSpec         bool
	initialObjectState *switchv1beta1.Switch
}

func (r *SwitchReconciler) newSwitchClient(ctx context.Context) *SwitchClient {
	return &SwitchClient{
		Client: r.Client,
		Logger: r.Log,
		ctx:    ctx,
	}
}

// +kubebuilder:rbac:groups=switch.onmetal.de,resources=switches,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=switch.onmetal.de,resources=switches/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=switch.onmetal.de,resources=switches/finalizers,verbs=update
// +kubebuilder:rbac:groups=switch.onmetal.de,resources=switchconfigs,verbs=get;list;watch
// +kubebuilder:rbac:groups=machine.onmetal.de,resources=inventories,verbs=get;list;watch
// +kubebuilder:rbac:groups=machine.onmetal.de,resources=inventories/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=ipam.onmetal.de,resources=subnets,verbs=get;list;watch;create;update;patch;delete;deletecollection
// +kubebuilder:rbac:groups=ipam.onmetal.de,resources=subnets/status,verbs=get
// +kubebuilder:rbac:groups=ipam.onmetal.de,resources=ips,verbs=get;list;create;update;patch;delete;deletecollection
// +kubebuilder:rbac:groups=ipam.onmetal.de,resources=ips/status,verbs=get

func (r *SwitchReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	nestedCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	obj := &switchv1beta1.Switch{}
	if err := r.Get(nestedCtx, req.NamespacedName, obj); err != nil {
		switch {
		case apierrors.IsNotFound(err):
			r.Log.Info("requested Switch object not found", "name", req.NamespacedName)
		default:
			r.Log.Info("failed to get requested Switch object", "name", req.NamespacedName, "error", err)
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if !obj.GetDeletionTimestamp().IsZero() {
		if !controllerutil.ContainsFinalizer(obj, constants.SwitchFinalizer) {
			return ctrl.Result{}, nil
		}
		err := r.finalize(ctx, obj)
		if err != nil {
			return ctrl.Result{Requeue: true}, err
		}
		controllerutil.RemoveFinalizer(obj, constants.SwitchFinalizer)
		err = r.Update(ctx, obj)
		return ctrl.Result{}, err
	}
	if !controllerutil.ContainsFinalizer(obj, constants.SwitchFinalizer) {
		controllerutil.AddFinalizer(obj, constants.SwitchFinalizer)
		err := r.Update(ctx, obj)
		return ctrl.Result{}, err
	}

	cl := r.newSwitchClient(ctx)
	cl.initialObjectState = obj.DeepCopy()
	result, err := cl.reconcile(nestedCtx, obj)
	return result, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *SwitchReconciler) SetupWithManager(mgr ctrl.Manager) error {
	log := mgr.GetLogger().WithName("switch-controller-setup")

	// setting up the label selector predicate to filter:
	// - inventoried switches
	labelSelectorPredicate, err := predicate.LabelSelectorPredicate(metav1.LabelSelector{
		MatchExpressions: []metav1.LabelSelectorRequirement{
			{
				Key:      constants.InventoriedLabel,
				Operator: metav1.LabelSelectorOpIn,
				Values:   []string{"true"},
			},
		},
	})
	if err != nil {
		log.Error(err, "failed to setup predicates")
		return err
	}

	// predicate to filter switch object update which was not caused
	// by conditions lastUpdateTimestamp change.
	discoverObjectChangesPredicate := predicate.Funcs{
		UpdateFunc: detectChangesPredicate,
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&switchv1beta1.Switch{}).
		WithOptions(controller.Options{
			// RateLimiter:  workqueue.NewItemExponentialFailureRateLimiter(time.Millisecond*500, time.Minute),
			RecoverPanic: pointer.Bool(true),
		}).
		WithEventFilter(predicate.And(labelSelectorPredicate, discoverObjectChangesPredicate)).
		Watches(&switchv1beta1.Switch{}, &handler.Funcs{
			UpdateFunc: r.handleSwitchUpdateEvent,
			DeleteFunc: r.handleSwitchDeleteEvent,
		}).
		Complete(r)
}

func detectChangesPredicate(e event.UpdateEvent) bool {
	objOld, okOld := e.ObjectOld.(*switchv1beta1.Switch)
	objNew, okNew := e.ObjectNew.(*switchv1beta1.Switch)
	if !okOld || !okNew {
		return false
	}
	return switchespkg.ObjectChanged(objOld, objNew)
}

func (r *SwitchReconciler) handleSwitchUpdateEvent(_ context.Context, e event.UpdateEvent, q workqueue.RateLimitingInterface) {
	objOld, okOld := e.ObjectOld.(*switchv1beta1.Switch)
	objNew, okNew := e.ObjectNew.(*switchv1beta1.Switch)
	if !okOld || !okNew {
		return
	}
	// if switch object has no changes, which affect neighbors, then there is no need to
	// enqueue it's neighbors for reconciliation.
	if !switchespkg.ObjectChanged(objOld, objNew) {
		return
	}
	switchesQueue := make(map[string]struct{})
	for _, nicData := range objOld.Status.Interfaces {
		if !switchespkg.NeighborIsSwitch(nicData) {
			continue
		}
		switchesQueue[nicData.Peer.GetObjectReferenceName()] = struct{}{}
	}
	for _, nicData := range objNew.Status.Interfaces {
		if !switchespkg.NeighborIsSwitch(nicData) {
			continue
		}
		switchesQueue[nicData.Peer.GetObjectReferenceName()] = struct{}{}
	}
	for name := range switchesQueue {
		q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
			Namespace: objNew.Namespace,
			Name:      name,
		}})
	}
}

func (r *SwitchReconciler) handleSwitchDeleteEvent(_ context.Context, e event.DeleteEvent, q workqueue.RateLimitingInterface) {
	obj, ok := e.Object.(*switchv1beta1.Switch)
	if !ok {
		return
	}
	switchesQueue := make(map[string]struct{})
	for _, nicData := range obj.Status.Interfaces {
		if !switchespkg.NeighborIsSwitch(nicData) {
			continue
		}
		switchesQueue[nicData.Peer.GetObjectReferenceName()] = struct{}{}
	}
	for name := range switchesQueue {
		q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
			Namespace: obj.Namespace,
			Name:      name,
		}})
	}
}

func (s *SwitchClient) reconcile(ctx context.Context, obj *switchv1beta1.Switch) (ctrl.Result, error) {
	if !obj.GetManaged() {
		return ctrl.Result{}, nil
	}
	proc := stateproc.NewGenericStateProcessor[*switchv1beta1.Switch](s, s.Logger).
		RegisterHandler(s.preprocessingCheck).
		RegisterHandler(s.initialize).
		RegisterHandler(s.updateInterfaces).
		RegisterHandler(s.updateNeighbors).
		RegisterHandler(s.updateLayerAndRole).
		RegisterHandler(s.updateConfigRef).
		RegisterHandler(s.updatePortParameters).
		RegisterHandler(s.updateLoopbacks).
		RegisterHandler(s.updateASN).
		RegisterHandler(s.updateSubnets).
		RegisterHandler(s.updateIPAddresses).
		RegisterHandler(s.setStateReady)

	err := proc.Compute(obj)
	if errors.IsMissingRequirements(err) {
		return ctrl.Result{}, err
	}
	if errors.IsInvalidConfigSelector(err) {
		switchespkg.UpdateSwitchConfigSelector(obj)
		s.updateSpec = true
	}
	return s.patch(ctx, obj)
}

func (s *SwitchClient) patch(ctx context.Context, obj *switchv1beta1.Switch) (ctrl.Result, error) {
	obj.ManagedFields = make([]metav1.ManagedFieldsEntry, 0)
	if s.updateSpec {
		err := s.Patch(ctx, obj, client.Apply, switchespkg.PatchOpts)
		if err != nil {
			s.Info("failed to update switch configuration", "name", obj.NamespacedName(), "error", err)
		}
		return ctrl.Result{}, err
	}
	if err := s.Status().Patch(ctx, obj, client.Apply, switchespkg.PatchOpts); err != nil {
		s.Info("failed to update switch configuration", "name", obj.NamespacedName(), "error", err)
		return ctrl.Result{}, err
	}
	return ctrl.Result{Requeue: obj.StateNotReady()}, nil
}

func (r *SwitchReconciler) finalize(ctx context.Context, obj *switchv1beta1.Switch) error {
	selector := labels.NewSelector()
	purposeReq, _ := labels.NewRequirement(constants.IPAMObjectPurposeLabel, selection.In, []string{constants.IPAMSwitchPortPurpose})
	ownerReq, _ := labels.NewRequirement(constants.IPAMObjectOwnerLabel, selection.In, []string{obj.Name})
	selector = selector.Add(*purposeReq).Add(*ownerReq)
	opts := client.ListOptions{
		LabelSelector: selector,
		Namespace:     obj.Namespace,
	}
	delOpts := &client.DeleteAllOfOptions{
		ListOptions: opts,
	}
	err := r.DeleteAllOf(ctx, &ipamv1alpha1.Subnet{}, delOpts, client.InNamespace(obj.Namespace))
	return err
}

func (s *SwitchClient) preprocessingCheck(obj *switchv1beta1.Switch) stateproc.StateFuncResult {
	var result stateproc.StateFuncResult
	if !obj.StatePending() {
		return result
	}
	if obj.GetInventoryRef() == constants.EmptyString {
		result.Break = true
		result.Err = errors.NewSwitchError(errors.ErrorReasonMissingRequirements, errors.MessageMissingInventoryRef, nil)
	}
	if switchespkg.SwitchConfigSelectorInvalid(obj) {
		result.Break = true
		result.Err = errors.NewSwitchError(errors.ErrorReasonInvalidConfigSelector, errors.MessageInvalidConfigSelector, nil)
	}
	return result
}

func (s *SwitchClient) initialize(obj *switchv1beta1.Switch) stateproc.StateFuncResult {
	if obj.Uninitialized() {
		obj.Status = switchv1beta1.SwitchStatus{
			ConfigRef:         nil,
			ASN:               nil,
			TotalPorts:        nil,
			SwitchPorts:       nil,
			Role:              nil,
			Interfaces:        make(map[string]*switchv1beta1.InterfaceSpec),
			LoopbackAddresses: make([]*switchv1beta1.IPAddressSpec, 0),
			Subnets:           make([]*switchv1beta1.SubnetSpec, 0),
			Message:           nil,
		}
		obj.SetLayer(255)
		if obj.GetTopSpine() {
			obj.SetLayer(0)
		}
		switchespkg.SetState(obj, constants.SwitchStateInitial, constants.EmptyString)
	}
	if obj.GetTopSpine() && obj.GetLayer() != 0 {
		obj.SetLayer(0)
	}
	obj.SetCondition(constants.ConditionInitialized, true)
	return stateproc.StateFuncResult{}
}

func (s *SwitchClient) updateInterfaces(obj *switchv1beta1.Switch) stateproc.StateFuncResult {
	result := stateproc.StateFuncResult{}
	if obj.GetInventoryRef() == constants.EmptyString {
		obj.SetCondition(constants.ConditionInterfacesOK, false).
			SetReason(errors.ErrorReasonMissingRequirements.String()).
			SetMessage(errors.MessageMissingInventoryRef)
		switchespkg.SetState(obj, constants.SwitchStatePending, errors.StateMessageMissingRequirements)
		result.Break = true
		result.Err = errors.NewSwitchError(errors.ErrorReasonMissingInventoryRef, errors.MessageMissingInventoryRef, nil)
		return result
	}
	inventory := &inventoryv1alpha1.Inventory{}
	err := s.Get(s.ctx, types.NamespacedName{
		Namespace: obj.Namespace,
		Name:      obj.Spec.InventoryRef.Name,
	}, inventory)
	if err == nil {
		switchespkg.ApplyInterfacesFromInventory(obj, inventory)
		obj.SetCondition(constants.ConditionInterfacesOK, true)
		switchespkg.SetState(obj, constants.SwitchStateProcessing, constants.EmptyString)
		return result
	}
	reason := errors.ErrorReasonRequestFailed
	message := errors.MessageRequestFailedWithKind("Inventory")
	if apierrors.IsNotFound(err) {
		reason = errors.ErrorReasonObjectNotExist
		message = errors.MessageObjectNotExistWithKind("Inventory")
	}
	obj.SetCondition(constants.ConditionInterfacesOK, false).
		SetReason(reason.String()).
		SetMessage(message)
	switchespkg.SetState(obj, constants.SwitchStateInvalid, errors.StateMessageRequestRelatedObjectsFailed)
	result.Break = true
	result.Err = errors.NewSwitchError(reason, message, err)
	return result
}

func (s *SwitchClient) updateConfigRef(obj *switchv1beta1.Switch) stateproc.StateFuncResult {
	result := stateproc.StateFuncResult{}
	switchConfig, err := s.getSwitchConfig(obj)
	if err == nil {
		obj.SetConfigRef(switchConfig.Name)
		if !switchConfig.RoutingConfigTemplateIsEmpty() {
			obj.SetRoutingConfigTemplate(switchConfig.GetRoutingConfigTemplate())
		}
		obj.SetCondition(constants.ConditionConfigRefOK, true)
		switchespkg.SetState(obj, constants.SwitchStateProcessing, constants.EmptyString)
		return result
	}
	reason := errors.ErrorReasonRequestFailed
	message := errors.MessageRequestFailedWithKind("SwitchConfig")
	if err.Error() == errors.MessageObjectNotExistWithKind("SwitchConfig") {
		reason = errors.ErrorReasonObjectNotExist
		message = errors.MessageFailedToDiscoverConfig
	}
	if err.Error() == errors.MessageToManyCandidatesFoundWithKind("SwitchConfig") {
		reason = errors.ErrorReasonTooManyCandidates
		message = errors.MessageFailedToDiscoverConfig
	}
	obj.SetCondition(constants.ConditionConfigRefOK, false).
		SetReason(reason.String()).
		SetMessage(message)
	obj.SetConfigRef(constants.EmptyString)
	switchespkg.SetState(obj, constants.SwitchStatePending, message)
	result.Break = true
	result.Err = errors.NewSwitchError(reason, message, err)
	return result
}

func (s *SwitchClient) updatePortParameters(obj *switchv1beta1.Switch) stateproc.StateFuncResult {
	result := stateproc.StateFuncResult{}
	config := &switchv1beta1.SwitchConfig{}
	if err := s.Get(s.ctx, types.NamespacedName{
		Namespace: obj.Namespace,
		Name:      obj.Status.ConfigRef.Name,
	}, config); err != nil {
		obj.SetCondition(constants.ConditionPortParametersOK, false).
			SetReason(errors.ErrorReasonRequestFailed.String()).
			SetMessage(errors.MessageRequestFailedWithKind("SwitchConfigList"))
		switchespkg.SetState(obj, constants.SwitchStateInvalid, errors.StateMessageRequestRelatedObjectsFailed)
		result.Break = true
		result.Err = errors.NewSwitchError(
			errors.ErrorReasonRequestFailed,
			errors.MessageRequestFailedWithKind("SwitchConfig"),
			err,
		)
		return result
	}
	switches, err := s.getSwitches(obj.Namespace)
	if err != nil {
		obj.SetCondition(constants.ConditionPortParametersOK, false).
			SetReason(errors.ErrorReasonRequestFailed.String()).
			SetMessage(errors.MessageRequestFailedWithKind("SwitchList"))
		switchespkg.SetState(obj, constants.SwitchStateInvalid, errors.StateMessageRequestRelatedObjectsFailed)
		result.Break = true
		result.Err = errors.NewSwitchError(
			errors.ErrorReasonRequestFailed,
			errors.MessageRequestFailedWithKind("SwitchList"),
			err,
		)
		return result
	}
	switchespkg.ApplyInterfaceParams(obj, config)
	switchespkg.InheritInterfaceParams(obj, switches)
	switchespkg.AlignInterfacesWithParams(obj)
	obj.SetCondition(constants.ConditionPortParametersOK, true)
	switchespkg.SetState(obj, constants.SwitchStateProcessing, constants.EmptyString)
	return result
}

func (s *SwitchClient) updateNeighbors(obj *switchv1beta1.Switch) stateproc.StateFuncResult {
	result := stateproc.StateFuncResult{}
	switches, err := s.getSwitches(obj.Namespace)
	if err != nil {
		obj.SetCondition(constants.ConditionNeighborsOK, false).
			SetReason(errors.ErrorReasonRequestFailed.String()).
			SetMessage(errors.MessageRequestFailedWithKind("SwitchList"))
		switchespkg.SetState(obj, constants.SwitchStateInvalid, errors.StateMessageRequestRelatedObjectsFailed)
		result.Break = true
		result.Err = errors.NewSwitchError(
			errors.ErrorReasonRequestFailed,
			errors.MessageRequestFailedWithKind("SwitchList"),
			err,
		)
		return result
	}
	for _, item := range switches.Items {
		for _, nicData := range obj.Status.Interfaces {
			if nicData.Peer == nil {
				continue
			}
			if nicData.Peer.PeerInfoSpec == nil {
				continue
			}
			if reflect.DeepEqual(nicData.Peer.PeerInfoSpec, &switchv1beta1.PeerInfoSpec{}) {
				continue
			}
			peerChassisID := nicData.Peer.PeerInfoSpec.GetChassisID()
			if strings.ReplaceAll(peerChassisID, ":", "") != item.Annotations[constants.HardwareChassisIDAnnotation] {
				continue
			}
			nicData.Peer.SetObjectReference(item.Name, item.Namespace)
		}
	}
	obj.SetCondition(constants.ConditionNeighborsOK, true)
	switchespkg.SetState(obj, constants.SwitchStateProcessing, constants.EmptyString)
	return result
}

func (s *SwitchClient) updateLayerAndRole(obj *switchv1beta1.Switch) stateproc.StateFuncResult {
	result := stateproc.StateFuncResult{}
	switches, err := s.getSwitches(obj.Namespace)
	if err != nil {
		obj.SetCondition(constants.ConditionLayerAndRoleOK, false).
			SetReason(errors.ErrorReasonRequestFailed.String()).
			SetMessage(errors.MessageRequestFailedWithKind("SwitchList"))
		switchespkg.SetState(obj, constants.SwitchStateInvalid, errors.StateMessageRequestRelatedObjectsFailed)
		result.Break = true
		result.Err = errors.NewSwitchError(
			errors.ErrorReasonRequestFailed,
			errors.MessageRequestFailedWithKind("SwitchList"),
			err,
		)
		return result
	}
	switchespkg.ComputeLayer(obj, switches)
	if obj.GetLayer() == 255 {
		obj.SetCondition(constants.ConditionLayerAndRoleOK, false).
			SetReason(errors.ErrorReasonFailedToComputeLayer.String()).
			SetMessage(errors.MessageFailedToComputeLayer)
		switchespkg.SetState(obj, constants.SwitchStateInvalid, errors.StateMessageRelatedObjectsStateInvalid)
		result.Break = true
		result.Err = errors.NewSwitchError(
			errors.ErrorReasonFailedToComputeLayer,
			errors.StateMessageRelatedObjectsStateInvalid,
			nil,
		)
		return result
	}
	switchespkg.SetRole(obj)
	obj.SetCondition(constants.ConditionLayerAndRoleOK, true)
	switchespkg.SetState(obj, constants.SwitchStateProcessing, constants.EmptyString)
	return result
}

func (s *SwitchClient) updateLoopbacks(obj *switchv1beta1.Switch) stateproc.StateFuncResult {
	result := stateproc.StateFuncResult{}
	loopbacks := &ipamv1alpha1.IPList{}
	af, err := s.getIPAMObjectsList(obj, loopbacks)
	if err != nil {
		obj.SetCondition(constants.ConditionLoopbacksOK, false).
			SetReason(errors.ErrorReasonRequestFailed.String()).
			SetMessage(errors.MessageRequestFailedWithKind("IPList"))
		switchespkg.SetState(obj, constants.SwitchStateInvalid, errors.StateMessageRequestRelatedObjectsFailed)
		result.Break = true
		result.Err = errors.NewSwitchError(
			errors.ErrorReasonRequestFailed,
			errors.MessageRequestFailedWithKind("IPList"),
			err,
		)
		return result
	}
	loopbacksToApply, addressFamiliesMap := processLoopbacks(loopbacks, af)
	afOK := switchespkg.AddressFamiliesMatchConfig(true, af.GetIPv6(), addressFamiliesMap)
	if len(loopbacksToApply) == 0 || !afOK {
		_ = s.createLoopbackIPs(obj)
		obj.SetCondition(constants.ConditionLoopbacksOK, false).
			SetReason(errors.ErrorReasonObjectNotExist.String()).
			SetMessage(errors.MessageMissingLoopbackV4IP)
		switchespkg.SetState(obj, constants.SwitchStateInvalid, errors.StateMessageMissingRequirements)
		result.Break = true
		result.Err = errors.NewSwitchError(errors.ErrorReasonObjectNotExist, errors.MessageMissingLoopbackV4IP, nil)
		return result
	}
	obj.Status.LoopbackAddresses = make([]*switchv1beta1.IPAddressSpec, len(loopbacksToApply))
	copy(obj.Status.LoopbackAddresses, loopbacksToApply)
	obj.SetCondition(constants.ConditionLoopbacksOK, true)
	switchespkg.SetState(obj, constants.SwitchStateProcessing, constants.EmptyString)
	return result
}

func processLoopbacks(
	list *ipamv1alpha1.IPList,
	af *switchv1beta1.AddressFamiliesMap) ([]*switchv1beta1.IPAddressSpec, map[ipamv1alpha1.SubnetAddressType]*bool) {
	addressFamiliesMap := map[ipamv1alpha1.SubnetAddressType]*bool{
		ipamv1alpha1.CIPv4SubnetType: nil,
		ipamv1alpha1.CIPv6SubnetType: nil,
	}
	loopbacksToApply := make([]*switchv1beta1.IPAddressSpec, 0)
	for _, item := range list.Items {
		var ipAF string
		if item.Status.State != ipamv1alpha1.CFinishedIPState {
			continue
		}
		if !af.GetIPv6() && item.Spec.IP.Net.Is6() {
			continue
		}
		switch {
		case item.Spec.IP.Net.Is4():
			addressFamiliesMap[ipamv1alpha1.CIPv4SubnetType] = pointer.Bool(true)
			ipAF = constants.IPv4AF
		case item.Spec.IP.Net.Is6():
			addressFamiliesMap[ipamv1alpha1.CIPv6SubnetType] = pointer.Bool(true)
			ipAF = constants.IPv6AF
		}
		ip := &switchv1beta1.IPAddressSpec{}
		ip.SetObjectReference(item.Name, item.Namespace)
		ip.SetAddress(item.Spec.IP.String())
		ip.SetAddressFamily(ipAF)
		ip.SetExtraAddress(false)
		loopbacksToApply = append(loopbacksToApply, ip)
	}
	return loopbacksToApply, addressFamiliesMap
}

func (s *SwitchClient) updateASN(obj *switchv1beta1.Switch) stateproc.StateFuncResult {
	result := stateproc.StateFuncResult{}
	asn, err := switchespkg.CalculateASN(obj.Status.LoopbackAddresses)
	if err != nil {
		obj.SetCondition(constants.ConditionAsnOK, false).
			SetReason(errors.ErrorReasonASNCalculationFailed.String()).
			SetMessage(err.Error())
		switchespkg.SetState(obj, constants.SwitchStateInvalid, err.Error())
		result.Break = true
		result.Err = errors.NewSwitchError(
			errors.ErrorReasonASNCalculationFailed,
			constants.EmptyString,
			err,
		)
		return result
	}
	obj.SetASN(asn)
	obj.SetCondition(constants.ConditionAsnOK, true)
	switchespkg.SetState(obj, constants.SwitchStateProcessing, constants.EmptyString)
	return result
}

func (s *SwitchClient) updateSubnets(obj *switchv1beta1.Switch) stateproc.StateFuncResult {
	result := stateproc.StateFuncResult{}
	subnets := &ipamv1alpha1.SubnetList{}
	af, err := s.getIPAMObjectsList(obj, subnets)
	if err != nil {
		obj.SetCondition(constants.ConditionSubnetsOK, false).
			SetReason(errors.ErrorReasonRequestFailed.String()).
			SetMessage(errors.MessageRequestFailedWithKind("SubnetList"))
		switchespkg.SetState(obj, constants.SwitchStateInvalid, errors.StateMessageRequestRelatedObjectsFailed)
		result.Break = true
		result.Err = errors.NewSwitchError(
			errors.ErrorReasonRequestFailed,
			errors.MessageRequestFailedWithKind("SubnetList"),
			err,
		)
		return result
	}
	subnetsToApply, addressFamiliesMap := processSubnets(obj, subnets, af)
	afOK := switchespkg.AddressFamiliesMatchConfig(af.GetIPv4(), af.GetIPv6(), addressFamiliesMap)
	if len(subnetsToApply) == 0 || !afOK {
		_ = s.createSubnets(obj, addressFamiliesMap)
		obj.SetCondition(constants.ConditionSubnetsOK, false).
			SetReason(errors.ErrorReasonObjectNotExist.String()).
			SetMessage(errors.MessageObjectNotExistWithKind("Subnet"))
		switchespkg.SetState(obj, constants.SwitchStateInvalid, errors.StateMessageMissingRequirements)
		result.Break = true
		result.Err = errors.NewSwitchError(
			errors.ErrorReasonObjectNotExist,
			errors.MessageObjectNotExistWithKind("Subnet"),
			nil,
		)
		return result
	}
	obj.Status.Subnets = make([]*switchv1beta1.SubnetSpec, len(subnetsToApply))
	copy(obj.Status.Subnets, subnetsToApply)
	obj.SetCondition(constants.ConditionSubnetsOK, true)
	switchespkg.SetState(obj, constants.SwitchStateProcessing, constants.EmptyString)
	return result
}

func processSubnets(
	obj *switchv1beta1.Switch,
	list *ipamv1alpha1.SubnetList,
	af *switchv1beta1.AddressFamiliesMap) ([]*switchv1beta1.SubnetSpec, map[ipamv1alpha1.SubnetAddressType]*bool) {
	addressFamiliesMap := map[ipamv1alpha1.SubnetAddressType]*bool{
		ipamv1alpha1.CIPv4SubnetType: nil,
		ipamv1alpha1.CIPv6SubnetType: nil,
	}
	subnetsToApply := make([]*switchv1beta1.SubnetSpec, 0)
	for _, item := range list.Items {
		if item.Status.State == ipamv1alpha1.CFailedSubnetState {
			continue
		}
		if item.Status.State == ipamv1alpha1.CProcessingSubnetState {
			addressFamiliesMap[item.Status.Type] = pointer.Bool(true)
			continue
		}
		if (!af.GetIPv4() && item.Status.Reserved.IsIPv4()) || (!af.GetIPv6() && item.Status.Reserved.IsIPv6()) {
			continue
		}
		requiredCapacity := switchespkg.GetTotalAddressesCount(obj.Status.Interfaces, item.Status.Type)
		// Have to change condition: replace CapacityLeft to Capacity.
		// After creation of child subnet objects for switch ports, value of
		// the CapacityLeft field becomes reduced, thus the check is always failed.
		if requiredCapacity.Cmp(item.Status.Capacity) > 0 {
			continue
		}
		addressFamiliesMap[item.Status.Type] = pointer.Bool(true)
		subnet := &switchv1beta1.SubnetSpec{}
		subnet.SetSubnetObjectRef(item.Name, item.Namespace)
		subnet.SetNetworkObjectRef(item.Spec.Network.Name, item.Namespace)
		subnet.SetCIDR(item.Status.Reserved.Net.String())
		subnet.SetAddressFamily(string(item.Status.Type))
		subnetsToApply = append(subnetsToApply, subnet)
	}
	return subnetsToApply, addressFamiliesMap
}

func (s *SwitchClient) updateIPAddresses(obj *switchv1beta1.Switch) stateproc.StateFuncResult {
	var err error
	result := stateproc.StateFuncResult{}
	for name, data := range obj.Status.Interfaces {
		if !strings.HasPrefix(name, constants.SwitchPortNamePrefix) {
			continue
		}
		switch data.GetDirection() {
		case constants.DirectionNorth:
			if data.Peer == nil {
				continue
			}
			err = s.updateNorthIPs(name, data, obj)
		case constants.DirectionSouth:
			err = s.updateSouthIPs(name, data, obj)
		}
	}
	if err != nil {
		obj.SetCondition(constants.ConditionIPAddressesOK, false).
			SetReason(errors.ErrorReasonIPAssignmentFailed.String()).
			SetMessage(err.Error())
		switchespkg.SetState(obj, constants.SwitchStateInvalid, errors.MessageFailedToAssignIPAddresses)
		result.Break = true
		result.Err = errors.NewSwitchError(
			errors.ErrorReasonIPAssignmentFailed,
			errors.MessageFailedToAssignIPAddresses,
			err,
		)
		return result
	}
	obj.SetCondition(constants.ConditionIPAddressesOK, true)
	switchespkg.SetState(obj, constants.SwitchStateProcessing, constants.EmptyString)
	return result
}

// IP objects' creation commented out due to decision that at the moment it is not required.
// If necessary it might be used.
func (s *SwitchClient) updateNorthIPs(
	// NIC name is only needed in case of IP object creation
	_ string,
	data *switchv1beta1.InterfaceSpec,
	obj *switchv1beta1.Switch,
) error {
	switches, err := s.getSwitches(obj.Namespace)
	if err != nil {
		return err
	}
	ipsToApply := make([]*switchv1beta1.IPAddressSpec, 0)
	for _, item := range switches.Items {
		if item.Name != data.Peer.GetObjectReferenceName() {
			continue
		}
		peerNICData := switchespkg.GetPeerData(item.Status.Interfaces, data.Peer.GetPortDescription(), data.Peer.GetPortID())
		if peerNICData == nil {
			continue
		}
		requestedIPs := switchespkg.RequestIPs(peerNICData)
		ipsToApply = append(ipsToApply, requestedIPs...)
		data.IP = make([]*switchv1beta1.IPAddressSpec, len(ipsToApply))
		copy(data.IP, ipsToApply)
		// if err := s.createIPs(obj, nic, ipsToApply); err != nil {
		// 	return err
		// }
	}
	return nil
}

// IP objects' creation commented out due to decision that at the moment it is not required.
// If necessary it might be used.
func (s *SwitchClient) updateSouthIPs(
	nic string,
	data *switchv1beta1.InterfaceSpec,
	obj *switchv1beta1.Switch,
) error {
	ipsToApply := make([]*switchv1beta1.IPAddressSpec, 0)
	extraIPs, err := switchespkg.GetExtraIPs(obj, nic)
	if err != nil {
		return err
	}
	ipsToApply = append(ipsToApply, extraIPs...)
	computedIPs, subnetsToCreate, err := switchespkg.GetComputedIPs(obj, nic, data)
	if err != nil {
		return err
	}
	ipsToApply = append(ipsToApply, computedIPs...)
	data.IP = make([]*switchv1beta1.IPAddressSpec, len(ipsToApply))
	copy(data.IP, ipsToApply)
	// if err := s.createIPs(obj, nic, computedIPs); err != nil {
	// 	return err
	// }
	_ = s.createSwitchPortsSubnets(obj, nic, subnetsToCreate)
	return nil
}

func (s *SwitchClient) setStateReady(obj *switchv1beta1.Switch) stateproc.StateFuncResult {
	switchespkg.SetState(obj, constants.SwitchStateReady, constants.EmptyString)
	return stateproc.StateFuncResult{}
}

func (s *SwitchClient) getSwitches(namespace string) (*switchv1beta1.SwitchList, error) {
	switches := &switchv1beta1.SwitchList{}
	inventoriedLabelReq, _ := labels.NewRequirement(constants.InventoriedLabel, selection.Exists, []string{})
	selector := labels.NewSelector().Add(*inventoriedLabelReq)
	opts := &client.ListOptions{
		LabelSelector: selector,
		Namespace:     namespace,
		Limit:         100,
	}
	if err := s.List(s.ctx, switches, opts); err != nil {
		return nil, err
	}
	return switches, nil
}

func (s *SwitchClient) getSwitchConfig(obj *switchv1beta1.Switch) (*switchv1beta1.SwitchConfig, error) {
	switchConfigs := &switchv1beta1.SwitchConfigList{}
	selector, err := metav1.LabelSelectorAsSelector(obj.GetConfigSelector())
	if err != nil {
		return nil, err
	}
	opts := &client.ListOptions{
		LabelSelector: selector,
		Namespace:     obj.Namespace,
		Limit:         100,
	}
	if err = s.List(s.ctx, switchConfigs, opts); err != nil {
		return nil, err
	}
	if len(switchConfigs.Items) == 0 {
		return nil, switchespkg.NewProcessingError(errors.MessageObjectNotExistWithKind("SwitchConfig"))
	}
	if len(switchConfigs.Items) > 1 {
		return nil, switchespkg.NewProcessingError(errors.MessageToManyCandidatesFoundWithKind("SwitchConfig"))
	}
	return &switchConfigs.Items[0], nil
}

// ----------------------------------------
// IPAM related functions for objects creation
// ----------------------------------------

func (s *SwitchClient) getIPAMObjectsList(
	obj *switchv1beta1.Switch,
	list client.ObjectList) (*switchv1beta1.AddressFamiliesMap, error) {
	config := &switchv1beta1.SwitchConfig{}
	key := types.NamespacedName{Namespace: obj.Namespace, Name: obj.Status.ConfigRef.Name}
	if err := s.Get(s.ctx, key, config); err != nil {
		return nil, err
	}
	var params *switchv1beta1.IPAMSelectionSpec
	switch list.(type) {
	case *ipamv1alpha1.IPList:
		params = config.Spec.IPAM.LoopbackAddresses
		if obj.Spec.IPAM != nil && obj.Spec.IPAM.LoopbackAddresses != nil {
			params = obj.Spec.IPAM.LoopbackAddresses
		}
	case *ipamv1alpha1.SubnetList:
		params = config.Spec.IPAM.SouthSubnets
		if obj.Spec.IPAM != nil && obj.Spec.IPAM.SouthSubnets != nil {
			params = obj.Spec.IPAM.SouthSubnets
		}
	default:
		return nil, switchespkg.NewProcessingError(errors.MessageInvalidInputType)
	}
	if err := s.listIPAMObjects(obj, params, list); err != nil {
		return nil, err
	}
	return config.Spec.IPAM.AddressFamily, nil
}

func (s *SwitchClient) createLoopbackIPs(obj *switchv1beta1.Switch) error {
	config := &switchv1beta1.SwitchConfig{}
	key := types.NamespacedName{Name: obj.Status.ConfigRef.Name, Namespace: obj.Namespace}
	if err := s.Get(s.ctx, key, config); err != nil {
		return err
	}
	loopbacksSubnets := &ipamv1alpha1.SubnetList{}
	if err := s.listIPAMObjects(obj, config.Spec.IPAM.LoopbackSubnets, loopbacksSubnets); err != nil {
		return err
	}
	labelsToApply, err := switchespkg.ResultingLabels(
		obj, obj.Spec.IPAM.GetLoopbacksSelection(), config.Spec.IPAM.LoopbackAddresses)
	if err != nil {
		return err
	}
	for _, item := range loopbacksSubnets.Items {
		if item.Status.State != ipamv1alpha1.CFinishedSubnetState {
			continue
		}
		// check whether loopbacks subnet has free address
		if resource.NewQuantity(1, resource.DecimalSI).Cmp(item.Status.CapacityLeft) > 1 {
			continue
		}
		// check only ipv6 flag since we need ipv4 loopback anyway to compute ASN
		if !config.Spec.IPAM.AddressFamily.GetIPv6() && item.Status.Type == ipamv1alpha1.CIPv6SubnetType {
			continue
		}
		ipObject, err := s.buildIPObject(obj, item, labelsToApply)
		if err != nil {
			return err
		}
		if err := s.Create(s.ctx, ipObject); err != nil {
			if !apierrors.IsAlreadyExists(err) {
				return err
			}
		}
	}
	return nil
}

func (s *SwitchClient) buildIPObject(
	obj *switchv1beta1.Switch,
	subnet ipamv1alpha1.Subnet,
	labelsToApply map[string]string,
) (*ipamv1alpha1.IP, error) {
	bits := constants.IPv4LoopbackBits
	if subnet.Status.Type == ipamv1alpha1.CIPv6SubnetType {
		bits = constants.IPv6LoopbackBits
	}
	proposedCIDR, err := subnet.ProposeForBits(uint8(bits))
	if err != nil {
		return nil, err
	}
	ip := proposedCIDR.Net.Addr()
	ok, err := s.checkIPAvailable(ip, obj.Namespace)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, switchespkg.NewProcessingError(errors.MessageDuplicatedIPAddress)
	}
	ipObject := &ipamv1alpha1.IP{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-lo-%s", obj.Name, strings.ToLower(string(subnet.Status.Type))),
			Namespace: obj.Namespace,
			Labels:    labelsToApply,
		},
		Spec: ipamv1alpha1.IPSpec{
			Subnet: v1.LocalObjectReference{Name: subnet.Name},
			Consumer: &ipamv1alpha1.ResourceReference{
				Kind:       obj.Kind,
				APIVersion: obj.APIVersion,
				Name:       obj.Name,
			},
			IP: &ipamv1alpha1.IPAddr{
				Net: ip,
			},
		},
	}
	return ipObject, nil
}

func (s *SwitchClient) checkIPAvailable(ip netip.Addr, ns string) (bool, error) {
	selector := labels.NewSelector()
	req, _ := labels.NewRequirement(constants.IPAMObjectPurposeLabel, selection.In, []string{constants.IPAMLoopbackPurpose})
	selector = selector.Add(*req)
	opts := &client.ListOptions{
		LabelSelector: selector,
		Namespace:     ns,
		Limit:         100,
	}
	ips := &ipamv1alpha1.IPList{}
	if err := s.List(s.ctx, ips, opts); err != nil {
		return false, err
	}
	for _, item := range ips.Items {
		if item.Spec.IP.Net == ip {
			return false, nil
		}
	}
	return true, nil
}

func (s *SwitchClient) createSubnets(
	obj *switchv1beta1.Switch,
	requiredAFFound map[ipamv1alpha1.SubnetAddressType]*bool,
) error {
	config := &switchv1beta1.SwitchConfig{}
	key := types.NamespacedName{Name: obj.Status.ConfigRef.Name, Namespace: obj.Namespace}
	if err := s.Get(s.ctx, key, config); err != nil {
		return err
	}
	carrierSubnets := &ipamv1alpha1.SubnetList{}
	if err := s.listIPAMObjects(obj, config.Spec.IPAM.CarrierSubnets, carrierSubnets); err != nil {
		return err
	}
	labelsToApply, err := switchespkg.ResultingLabels(obj, obj.Spec.IPAM.GetSubnetsSelection(), config.Spec.IPAM.SouthSubnets)
	if err != nil {
		return err
	}
	for _, item := range carrierSubnets.Items {
		if item.Status.State != ipamv1alpha1.CFinishedSubnetState {
			continue
		}
		if requiredAFFound[item.Status.Type] == nil {
			continue
		}
		if pointer.BoolDeref(requiredAFFound[item.Status.Type], false) {
			continue
		}
		subnet, err := s.buildSubnetObject(obj, item, labelsToApply)
		if err != nil {
			return err
		}
		if err := s.Create(s.ctx, subnet); err != nil {
			if !apierrors.IsAlreadyExists(err) {
				return err
			}
		}
	}
	return nil
}

func (s *SwitchClient) buildSubnetObject(
	obj *switchv1beta1.Switch,
	item ipamv1alpha1.Subnet,
	labelsToApply map[string]string,
) (*ipamv1alpha1.Subnet, error) {
	addressesRequired := switchespkg.GetTotalAddressesCount(obj.Status.Interfaces, item.Status.Type)
	proposedCIDR, err := item.ProposeForCapacity(addressesRequired)
	if err != nil {
		return nil, err
	}
	ok, err := s.checkSubnetAvailable(proposedCIDR, obj.Namespace)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, switchespkg.NewProcessingError(errors.MessageDuplicatedSubnet)
	}
	subnet := &ipamv1alpha1.Subnet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-sn-%s", obj.Name, strings.ToLower(string(item.Status.Type))),
			Namespace: obj.Namespace,
			Labels:    labelsToApply,
		},
		Spec: ipamv1alpha1.SubnetSpec{
			CIDR: proposedCIDR,
			Network: v1.LocalObjectReference{
				Name: item.Spec.Network.Name,
			},
			ParentSubnet: v1.LocalObjectReference{
				Name: item.Name,
			},
			Consumer: &ipamv1alpha1.ResourceReference{
				Kind:       obj.Kind,
				APIVersion: obj.APIVersion,
				Name:       obj.Name,
			},
		},
	}
	return subnet, nil
}

func (s *SwitchClient) checkSubnetAvailable(cidr *ipamv1alpha1.CIDR, ns string) (bool, error) {
	selector := labels.NewSelector()
	req, _ := labels.NewRequirement(constants.IPAMObjectPurposeLabel, selection.In, []string{constants.IPAMSouthSubnetPurpose})
	selector = selector.Add(*req)
	opts := &client.ListOptions{
		LabelSelector: selector,
		Namespace:     ns,
		Limit:         100,
	}
	subnets := &ipamv1alpha1.SubnetList{}
	if err := s.List(s.ctx, subnets, opts); err != nil {
		return false, err
	}
	for _, item := range subnets.Items {
		if cidr.String() == item.Spec.CIDR.String() {
			return false, nil
		}
	}
	return true, nil
}

func (s *SwitchClient) listIPAMObjects(
	obj *switchv1beta1.Switch,
	params *switchv1beta1.IPAMSelectionSpec,
	list client.ObjectList,
) error {
	selector, err := switchespkg.GetSelectorFromIPAMSpec(obj, params)
	if err != nil {
		return err
	}
	opts := &client.ListOptions{
		LabelSelector: selector,
		Namespace:     obj.Namespace,
		Limit:         100,
	}
	if err := s.List(s.ctx, list, opts); err != nil {
		return err
	}
	return nil
}

//nolint:unused
func (s *SwitchClient) createIPs(obj *switchv1beta1.Switch, nic string, ips []*switchv1beta1.IPAddressSpec) error {
	for _, item := range ips {
		prefix, _ := netip.ParsePrefix(item.GetAddress())
		addr := prefix.Addr()
		ip := &ipamv1alpha1.IP{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("%s-%s-%s",
					obj.Name, strings.ToLower(item.GetAddressFamily()), strings.ToLower(nic)),
				Namespace: obj.Namespace,
				Labels: map[string]string{
					constants.IPAMObjectPurposeLabel: constants.IPAMSwitchPortPurpose,
					constants.IPAMObjectOwnerLabel:   obj.Name,
					constants.IPAMObjectNICNameLabel: nic,
				},
			},
			Spec: ipamv1alpha1.IPSpec{
				Subnet: v1.LocalObjectReference{
					Name: item.GetObjectReferenceName(),
				},
				Consumer: &ipamv1alpha1.ResourceReference{
					Kind:       obj.Kind,
					APIVersion: obj.APIVersion,
					Name:       obj.Name,
				},
				IP: &ipamv1alpha1.IPAddr{
					Net: addr,
				},
			},
		}
		if err := s.Create(s.ctx, ip); err != nil {
			if !apierrors.IsAlreadyExists(err) {
				return err
			}
		}
	}
	return nil
}

func (s *SwitchClient) createSwitchPortsSubnets(
	obj *switchv1beta1.Switch,
	nic string,
	subnets []*ipamv1alpha1.SubnetSpec,
) error {
	gvk := obj.GroupVersionKind()
	for _, item := range subnets {
		cidr := item.CIDR.String()
		hash := md5.Sum([]byte(cidr))
		item.Consumer = &ipamv1alpha1.ResourceReference{
			APIVersion: gvk.GroupVersion().String(),
			Kind:       gvk.Kind,
			Name:       obj.Name,
		}
		subnet := &ipamv1alpha1.Subnet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-%s-%x", obj.Name, strings.ToLower(nic), hash[:4]),
				Namespace: obj.Namespace,
				Labels: map[string]string{
					constants.IPAMObjectPurposeLabel: constants.IPAMSwitchPortPurpose,
					constants.IPAMObjectOwnerLabel:   obj.Name,
					constants.IPAMObjectNICNameLabel: nic,
				},
			},
			Spec: *item,
		}
		if err := s.Create(s.ctx, subnet); err != nil {
			if !apierrors.IsAlreadyExists(err) {
				return err
			}
		}
	}
	return nil
}
