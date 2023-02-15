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
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/go-logr/logr"
	ipamv1alpha1 "github.com/onmetal/ipam/api/v1alpha1"
	"inet.af/netaddr"
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
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	inventoryv1alpha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	switchv1beta1 "github.com/onmetal/metal-api/apis/switch/v1beta1"
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

	ctx context.Context
}

func (r *SwitchReconciler) newSwitchClient(ctx context.Context) *SwitchClient {
	return &SwitchClient{
		Client: r.Client,
		ctx:    ctx,
	}
}

// +kubebuilder:rbac:groups=switch.onmetal.de,resources=switches,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=switch.onmetal.de,resources=switches/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=switch.onmetal.de,resources=switches/finalizers,verbs=update
// +kubebuilder:rbac:groups=switch.onmetal.de,resources=switchconfigs,verbs=get;list;watch
// +kubebuilder:rbac:groups=machine.onmetal.de,resources=inventories,verbs=get;list;watch
// +kubebuilder:rbac:groups=machine.onmetal.de,resources=inventories/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=ipam.onmetal.de,resources=subnets,verbs=get;list;watch;create;update;patch
// +kubebuilder:rbac:groups=ipam.onmetal.de,resources=subnets/status,verbs=get
// +kubebuilder:rbac:groups=ipam.onmetal.de,resources=ips,verbs=get;list;create;update;patch
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
	result, err := r.reconcile(nestedCtx, obj)
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
				Key:      InventoriedLabel,
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
			RateLimiter:  workqueue.NewItemExponentialFailureRateLimiter(time.Millisecond*500, time.Minute*5),
			RecoverPanic: pointer.Bool(true),
		}).
		WithEventFilter(predicate.And(labelSelectorPredicate, discoverObjectChangesPredicate)).
		Watches(&source.Kind{Type: &switchv1beta1.Switch{}}, &handler.Funcs{
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
	return reconciliationRequired(objOld, objNew)
}

func (r *SwitchReconciler) handleSwitchUpdateEvent(e event.UpdateEvent, q workqueue.RateLimitingInterface) {
	objOld, okOld := e.ObjectOld.(*switchv1beta1.Switch)
	objNew, okNew := e.ObjectNew.(*switchv1beta1.Switch)
	if !okOld || !okNew {
		return
	}
	// if switch object has no changes, which affect neighbors, then there is no need to
	// enqueue it's neighbors for reconciliation.
	if !reconciliationRequired(objOld, objNew) {
		return
	}
	switchesQueue := make(map[string]struct{})
	for _, nicData := range objOld.Status.Interfaces {
		if !neighborIsSwitch(nicData) {
			continue
		}
		switchesQueue[nicData.Peer.GetObjectReferenceName()] = struct{}{}
	}
	for _, nicData := range objNew.Status.Interfaces {
		if !neighborIsSwitch(nicData) {
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

func (r *SwitchReconciler) handleSwitchDeleteEvent(e event.DeleteEvent, q workqueue.RateLimitingInterface) {
	obj, ok := e.Object.(*switchv1beta1.Switch)
	if !ok {
		return
	}
	switchesQueue := make(map[string]struct{})
	for _, nicData := range obj.Status.Interfaces {
		if !neighborIsSwitch(nicData) {
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

func (r *SwitchReconciler) reconcile(ctx context.Context, obj *switchv1beta1.Switch) (ctrl.Result, error) {
	if !obj.GetDeletionTimestamp().IsZero() {
		return ctrl.Result{}, nil
	}
	if !obj.GetManaged() {
		return ctrl.Result{}, nil
	}
	cl := r.newSwitchClient(ctx)
	proc := NewGenericStateProcessor[*switchv1beta1.Switch](cl, r.Log)
	proc.setFunctions([]func(*switchv1beta1.Switch) StateFuncResult{
		cl.preprocessingCheck,
		cl.initialize,
		cl.updateInterfaces,
		cl.updateConfigRef,
		cl.updatePortParameters,
		cl.updateNeighbors,
		cl.updateLayerAndRole,
		cl.updateLoopbacks,
		cl.updateASN,
		cl.updateSubnets,
		cl.updateIPAddresses,
		cl.setStateReady,
	})

	err := proc.compute(obj)
	if err != nil && err.Error() == ErrorMissingRequirements {
		return ctrl.Result{}, err
	}
	return r.patch(ctx, obj)
}

func (r *SwitchReconciler) patch(ctx context.Context, obj *switchv1beta1.Switch) (ctrl.Result, error) {
	obj.ManagedFields = make([]metav1.ManagedFieldsEntry, 0)
	if err := r.Status().Patch(ctx, obj, client.Apply, patchOpts); err != nil {
		r.Log.Info("failed to update switch configuration", "name", obj.NamespacedName(), "error", err)
		return ctrl.Result{}, err
	}
	switch obj.GetState() {
	case SwitchStateReady:
		return ctrl.Result{}, nil
	default:
		return ctrl.Result{Requeue: true}, nil
	}
}

func (s *SwitchClient) preprocessingCheck(obj *switchv1beta1.Switch) StateFuncResult {
	var result StateFuncResult
	if obj.GetState() != SwitchStatePending {
		return result
	}
	if _, ok := obj.Labels[SwitchTypeLabel]; !ok {
		result.Break = true
		result.Err = reconciliationError(ErrorMissingRequirements)
	}
	return result
}

func (s *SwitchClient) initialize(obj *switchv1beta1.Switch) StateFuncResult {
	if obj.GetState() == EmptyString {
		obj.Status = switchv1beta1.SwitchStatus{
			ConfigRef:         nil,
			ASN:               nil,
			TotalPorts:        nil,
			SwitchPorts:       nil,
			Role:              nil,
			Interfaces:        map[string]*switchv1beta1.InterfaceSpec{},
			LoopbackAddresses: make([]*switchv1beta1.IPAddressSpec, 0),
			Subnets:           make([]*switchv1beta1.SubnetSpec, 0),
			Message:           nil,
		}
		obj.SetLayer(255)
		if obj.GetTopSpine() {
			obj.SetLayer(0)
		}
		setState(obj, SwitchStateInitial, EmptyString)
	}
	obj.SetCondition(ConditionInitialized, true)
	return StateFuncResult{}
}

func (s *SwitchClient) updateInterfaces(obj *switchv1beta1.Switch) StateFuncResult {
	result := StateFuncResult{}
	inventory := &inventoryv1alpha1.Inventory{}
	if err := s.Get(s.ctx, types.NamespacedName{
		Namespace: obj.Namespace,
		Name:      obj.Spec.InventoryRef.Name,
	}, inventory); err != nil {
		reason := ReasonRequestFailed
		message := fmt.Sprintf("%s: Inventory", ErrorFailedToGetRequiredObject)
		if apierrors.IsNotFound(err) {
			reason = ReasonObjectNotExists
			message = fmt.Sprintf("%s: Inventory", ErrorRequiredObjectNotExist)
		}
		obj.SetCondition(ConditionInterfacesOK, false).
			SetReason(reason).
			SetMessage(message)
		setState(obj, SwitchStateInvalid, ErrorFailedToRequestRelatedObjects)
		result.Break = true
		result.Err = reconciliationError(ErrorFailedToGetRequiredObject, "type", "Inventory", "error", err)
		return result
	}
	applyInterfacesFromInventory(obj, inventory)
	obj.SetCondition(ConditionInterfacesOK, true)
	setState(obj, SwitchStateProcessing, EmptyString)
	return result
}

func (s *SwitchClient) updateConfigRef(obj *switchv1beta1.Switch) StateFuncResult {
	result := StateFuncResult{}
	switchConfigs := &switchv1beta1.SwitchConfigList{}
	switchType, ok := obj.Labels[SwitchTypeLabel]
	if !ok {
		obj.SetCondition(ConditionConfigRefOK, false).
			SetReason(ReasonMissingPrerequisites).
			SetMessage(fmt.Sprintf("%s: %s", ErrorTypeLabelMissed, SwitchTypeLabel))
		setState(obj, SwitchStatePending, ErrorMissingRequirements)
		result.Break = true
		result.Err = reconciliationError(ErrorTypeLabelMissed, "label", SwitchTypeLabel)
		return result
	}
	requirements, _ := labels.NewRequirement(SwitchConfigTypeLabelPrefix+switchType, selection.Exists, []string{})
	selector := labels.NewSelector().Add(*requirements)
	opts := &client.ListOptions{
		LabelSelector: selector,
		Namespace:     obj.Namespace,
		Limit:         100,
	}
	if err := s.List(s.ctx, switchConfigs, opts); err != nil {
		obj.SetCondition(ConditionConfigRefOK, false).
			SetReason(ReasonRequestFailed).
			SetMessage(fmt.Sprintf("%s: SwitchConfig", ErrorFailedToListObjects))
		setState(obj, SwitchStateInvalid, ErrorFailedToRequestRelatedObjects)
		result.Break = true
		result.Err = reconciliationError(ErrorFailedToListObjects, "type", "SwitchConfig", "error", err)
		return result
	}
	if len(switchConfigs.Items) == 0 {
		obj.SetCondition(ConditionConfigRefOK, false).
			SetReason(ReasonObjectNotExists).
			SetMessage(fmt.Sprintf("%s: SwitchConfig", ErrorRequiredObjectNotExist))
		setState(obj, SwitchStateInvalid, ErrorMissingRequirements)
		result.Break = true
		result.Err = reconciliationError(ErrorRequiredObjectNotExist, "type", "SwitchConfig")
		return result
	}
	obj.SetConfigRef(switchConfigs.Items[0].Name)
	obj.SetCondition(ConditionConfigRefOK, true)
	setState(obj, SwitchStateProcessing, EmptyString)
	return result
}

func (s *SwitchClient) updatePortParameters(obj *switchv1beta1.Switch) StateFuncResult {
	result := StateFuncResult{}
	config := &switchv1beta1.SwitchConfig{}
	if err := s.Get(s.ctx, types.NamespacedName{
		Namespace: obj.Namespace,
		Name:      obj.Status.ConfigRef.Name,
	}, config); err != nil {
		obj.SetCondition(ConditionPortParametersOK, false).
			SetReason(ReasonRequestFailed).
			SetMessage(fmt.Sprintf("%s: SwitchConfig", ErrorFailedToGetRequiredObject))
		setState(obj, SwitchStateInvalid, ErrorFailedToRequestRelatedObjects)
		result.Break = true
		result.Err = reconciliationError(ErrorFailedToGetRequiredObject, "type", "SwitchConfig", "error", err)
		return result
	}
	applyInterfaceParams(obj, config)
	obj.SetCondition(ConditionPortParametersOK, true)
	setState(obj, SwitchStateProcessing, EmptyString)
	return result
}

func (s *SwitchClient) updateNeighbors(obj *switchv1beta1.Switch) StateFuncResult {
	result := StateFuncResult{}
	switches, err := getSwitches(s.ctx, s.Client, obj.Namespace)
	if err != nil {
		obj.SetCondition(ConditionNeighborsOK, false).
			SetReason(ReasonRequestFailed).
			SetMessage(fmt.Sprintf("%s: Switch", ErrorFailedToListObjects))
		setState(obj, SwitchStateInvalid, ErrorFailedToRequestRelatedObjects)
		result.Break = true
		result.Err = reconciliationError(ErrorFailedToListObjects, "type", "Switch", "error", err)
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
			if strings.ReplaceAll(peerChassisID, ":", "") != item.Annotations[HardwareChassisIDAnnotation] {
				continue
			}
			nicData.Peer.SetObjectReference(item.Name, item.Namespace)
		}
	}
	obj.SetCondition(ConditionNeighborsOK, true)
	setState(obj, SwitchStateProcessing, EmptyString)
	return result
}

func (s *SwitchClient) updateLayerAndRole(obj *switchv1beta1.Switch) StateFuncResult {
	result := StateFuncResult{}
	switches, err := getSwitches(s.ctx, s.Client, obj.Namespace)
	if err != nil {
		obj.SetCondition(ConditionLayerAndRoleOK, false).
			SetReason(ReasonRequestFailed).
			SetMessage(fmt.Sprintf("%s: Switch", ErrorFailedToListObjects))
		setState(obj, SwitchStateInvalid, ErrorFailedToRequestRelatedObjects)
		result.Break = true
		result.Err = reconciliationError(ErrorFailedToListObjects, "type", "Switch", "error", err)
		return result
	}
	computeLayer(obj, switches)
	setRole(obj)
	obj.SetCondition(ConditionLayerAndRoleOK, true)
	setState(obj, SwitchStateProcessing, EmptyString)
	return result
}

func (s *SwitchClient) updateLoopbacks(obj *switchv1beta1.Switch) StateFuncResult {
	result := StateFuncResult{}
	loopbacks := &ipamv1alpha1.IPList{}
	af, err := getIPAMObjectsList(s.ctx, s.Client, obj, loopbacks)
	if err != nil {
		obj.SetCondition(ConditionLoopbacksOK, false).
			SetReason(ReasonRequestFailed).
			SetMessage(fmt.Sprintf("%s: IP", ErrorFailedToListObjects))
		setState(obj, SwitchStateInvalid, ErrorFailedToRequestRelatedObjects)
		result.Break = true
		result.Err = reconciliationError(ErrorFailedToListObjects, "type", "IP", "error", err)
		return result
	}
	loopbacksToApply, addressFamiliesMap := processLoopbacks(loopbacks, af)
	afOK := addressFamiliesMatchConfig(true, af.GetIPv6(), addressFamiliesMap)
	if len(loopbacksToApply) == 0 || !afOK {
		_ = createLoopbackIPs(s.ctx, s.Client, obj)
		obj.SetCondition(ConditionLoopbacksOK, false).
			SetReason(ReasonObjectNotExists).
			SetMessage(fmt.Sprintf("%s: IP v4", ErrorRequiredObjectNotExist))
		setState(obj, SwitchStateInvalid, ErrorMissingRequirements)
		result.Break = true
		result.Err = reconciliationError(ErrorRequiredObjectNotExist, "type", "IP")
		return result
	}
	obj.Status.LoopbackAddresses = make([]*switchv1beta1.IPAddressSpec, len(loopbacksToApply))
	copy(obj.Status.LoopbackAddresses, loopbacksToApply)
	obj.SetCondition(ConditionLoopbacksOK, true)
	setState(obj, SwitchStateProcessing, EmptyString)
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
			ipAF = IPv4AF
		case item.Spec.IP.Net.Is6():
			addressFamiliesMap[ipamv1alpha1.CIPv6SubnetType] = pointer.Bool(true)
			ipAF = IPv6AF
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

func (s *SwitchClient) updateASN(obj *switchv1beta1.Switch) StateFuncResult {
	result := StateFuncResult{}
	asn, err := calculateASN(obj.Status.LoopbackAddresses)
	if err != nil {
		obj.SetCondition(ConditionAsnOK, false).
			SetReason(ReasonASNCalculationFailed).
			SetMessage(err.Error())
		setState(obj, SwitchStateInvalid, err.Error())
		result.Break = true
		result.Err = err
		return result
	}
	obj.SetASN(asn)
	obj.SetCondition(ConditionAsnOK, true)
	setState(obj, SwitchStateProcessing, EmptyString)
	return result
}

func (s *SwitchClient) updateSubnets(obj *switchv1beta1.Switch) StateFuncResult {
	result := StateFuncResult{}
	subnets := &ipamv1alpha1.SubnetList{}
	af, err := getIPAMObjectsList(s.ctx, s.Client, obj, subnets)
	if err != nil {
		obj.SetCondition(ConditionSubnetsOK, false).
			SetReason(ReasonRequestFailed).
			SetMessage(fmt.Sprintf("%s: Subnet", ErrorFailedToListObjects))
		setState(obj, SwitchStateInvalid, ErrorFailedToRequestRelatedObjects)
		result.Break = true
		result.Err = reconciliationError(ErrorFailedToListObjects, "type", "Subnet", "error", err)
		return result
	}
	subnetsToApply, addressFamiliesMap := processSubnets(obj, subnets, af)
	afOK := addressFamiliesMatchConfig(af.GetIPv4(), af.GetIPv6(), addressFamiliesMap)
	if len(subnetsToApply) == 0 || !afOK {
		_ = createSubnets(s.ctx, s.Client, obj, addressFamiliesMap)
		obj.SetCondition(ConditionSubnetsOK, false).
			SetReason(ReasonObjectNotExists).
			SetMessage(fmt.Sprintf("%s: Subnet", ErrorRequiredObjectNotExist))
		setState(obj, SwitchStateInvalid, ErrorMissingRequirements)
		result.Break = true
		result.Err = reconciliationError(ErrorRequiredObjectNotExist, "type", "Subnet")
		return result
	}
	obj.Status.Subnets = make([]*switchv1beta1.SubnetSpec, len(subnetsToApply))
	copy(obj.Status.Subnets, subnetsToApply)
	obj.SetCondition(ConditionSubnetsOK, true)
	setState(obj, SwitchStateProcessing, EmptyString)
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
		requiredCapacity := getTotalAddressesCount(obj.Status.Interfaces, item.Status.Type)
		if requiredCapacity.Cmp(item.Status.CapacityLeft) > 0 {
			continue
		}
		addressFamiliesMap[item.Status.Type] = pointer.Bool(true)
		subnet := &switchv1beta1.SubnetSpec{}
		subnet.SetObjectReference(item.Name, item.Namespace)
		subnet.SetCIDR(item.Status.Reserved.Net.String())
		subnet.SetAddressFamily(string(item.Status.Type))
		subnetsToApply = append(subnetsToApply, subnet)
	}
	return subnetsToApply, addressFamiliesMap
}

func (s *SwitchClient) updateIPAddresses(obj *switchv1beta1.Switch) StateFuncResult {
	var err error
	result := StateFuncResult{}
	for name, data := range obj.Status.Interfaces {
		if !strings.HasPrefix(name, SwitchPortNamePrefix) {
			continue
		}
		switch data.GetDirection() {
		case DirectionNorth:
			if data.Peer == nil {
				continue
			}
			err = updateNorthIPs(s.ctx, s.Client, name, data, obj)
		case DirectionSouth:
			err = updateSouthIPs(s.ctx, name, data, obj)
		}
	}
	if err != nil {
		obj.SetCondition(ConditionIPAddressesOK, false).
			SetReason(ReasonIPAssignmentFailed).
			SetMessage(ErrorFailedIPAddressAssignment)
		setState(obj, SwitchStateInvalid, ErrorFailedIPAddressAssignment)
		result.Break = true
		result.Err = reconciliationError(ErrorFailedIPAddressAssignment, "error", err)
		return result
	}
	obj.SetCondition(ConditionIPAddressesOK, true)
	setState(obj, SwitchStateProcessing, EmptyString)
	return result
}

// IP objects' creation commented out due to decision that at the moment it is not required.
// If necessary it might be used.
func updateNorthIPs(
	ctx context.Context,
	cl client.Client,
	_ string,
	data *switchv1beta1.InterfaceSpec,
	obj *switchv1beta1.Switch,
) error {
	switches, err := getSwitches(ctx, cl, obj.Namespace)
	if err != nil {
		return err
	}
	ipsToApply := make([]*switchv1beta1.IPAddressSpec, 0)
	for _, item := range switches.Items {
		if item.Name != data.Peer.GetObjectReferenceName() {
			continue
		}
		peerNICData := getPeerData(item.Status.Interfaces, data.Peer.GetPortDescription(), data.Peer.GetPortID())
		if peerNICData == nil {
			continue
		}
		requestedIPs := requestIPs(peerNICData)
		ipsToApply = append(ipsToApply, requestedIPs...)
		data.IP = make([]*switchv1beta1.IPAddressSpec, len(ipsToApply))
		copy(data.IP, ipsToApply)
		// if err := r.createIPs(ctx, obj, nic, ipsToApply); err != nil {
		// 	return err
		// }
	}
	return nil
}

// IP objects' creation commented out due to decision that at the moment it is not required.
// If necessary it might be used.
func updateSouthIPs(
	_ context.Context,
	nic string,
	data *switchv1beta1.InterfaceSpec,
	obj *switchv1beta1.Switch,
) error {
	ipsToApply := make([]*switchv1beta1.IPAddressSpec, 0)
	extraIPs, err := getExtraIPs(obj, nic)
	if err != nil {
		return err
	}
	ipsToApply = append(ipsToApply, extraIPs...)
	computedIPs, err := getComputedIPs(obj, nic, data)
	if err != nil {
		return err
	}
	ipsToApply = append(ipsToApply, computedIPs...)
	data.IP = make([]*switchv1beta1.IPAddressSpec, len(ipsToApply))
	copy(data.IP, ipsToApply)
	// if err := r.createIPs(ctx, obj, nic, computedIPs); err != nil {
	// 	return err
	// }
	return nil
}

func (s *SwitchClient) setStateReady(obj *switchv1beta1.Switch) StateFuncResult {
	setState(obj, SwitchStateReady, EmptyString)
	return StateFuncResult{}
}

func getSwitches(ctx context.Context, cl client.Client, ns string) (*switchv1beta1.SwitchList, error) {
	switches := &switchv1beta1.SwitchList{}
	inventoriedLabelReq, _ := labels.NewRequirement(InventoriedLabel, selection.Exists, []string{})
	typedLabelReq, _ := labels.NewRequirement(SwitchTypeLabel, selection.Exists, []string{})
	selector := labels.NewSelector().Add(*inventoriedLabelReq).Add(*typedLabelReq)
	opts := &client.ListOptions{
		LabelSelector: selector,
		Namespace:     ns,
		Limit:         100,
	}
	if err := cl.List(ctx, switches, opts); err != nil {
		return nil, err
	}
	return switches, nil
}

// ----------------------------------------
// IPAM related functions for objects creation
// ----------------------------------------

func getIPAMObjectsList(
	ctx context.Context,
	cl client.Client,
	obj *switchv1beta1.Switch,
	list client.ObjectList) (*switchv1beta1.AddressFamiliesMap, error) {
	config := &switchv1beta1.SwitchConfig{}
	key := types.NamespacedName{Namespace: obj.Namespace, Name: obj.Status.ConfigRef.Name}
	if err := cl.Get(ctx, key, config); err != nil {
		return nil, fmt.Errorf("failed to get related switch config: %w", err)
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
		return nil, reconciliationError(ErrorInvalidInputType)
	}
	if err := listIPAMObjects(ctx, cl, obj, params, list); err != nil {
		return nil, reconciliationError(ErrorFailedToListObjects, "error", err)
	}
	return config.Spec.IPAM.AddressFamily, nil
}

func createLoopbackIPs(ctx context.Context, cl client.Client, obj *switchv1beta1.Switch) error {
	config := &switchv1beta1.SwitchConfig{}
	key := types.NamespacedName{Name: obj.Status.ConfigRef.Name, Namespace: obj.Namespace}
	if err := cl.Get(ctx, key, config); err != nil {
		return err
	}
	loopbacksSubnets := &ipamv1alpha1.SubnetList{}
	if err := listIPAMObjects(ctx, cl, obj, config.Spec.IPAM.LoopbackSubnets, loopbacksSubnets); err != nil {
		return err
	}
	labelsToApply, err := resultingLabels(
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
		ipObject, err := buildIPObject(ctx, cl, obj, item, labelsToApply)
		if err != nil {
			return err
		}
		if err := cl.Create(ctx, ipObject); err != nil {
			if !apierrors.IsAlreadyExists(err) {
				return err
			}
		}
	}
	return nil
}

func buildIPObject(
	ctx context.Context,
	cl client.Client,
	obj *switchv1beta1.Switch,
	subnet ipamv1alpha1.Subnet,
	labelsToApply map[string]string,
) (*ipamv1alpha1.IP, error) {
	bits := IPv4LoopbackBits
	if subnet.Status.Type == ipamv1alpha1.CIPv6SubnetType {
		bits = IPv6LoopbackBits
	}
	proposedCIDR, err := subnet.ProposeForBits(uint8(bits))
	if err != nil {
		return nil, err
	}
	ip := proposedCIDR.Net.IP()
	ok, err := checkIPAvailable(ctx, cl, ip, obj.Namespace)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, reconciliationError(ErrorDuplicateIPAddress)
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

func checkIPAvailable(ctx context.Context, cl client.Client, ip netaddr.IP, ns string) (bool, error) {
	selector := labels.NewSelector()
	req, _ := labels.NewRequirement(IPAMObjectPurposeLabel, selection.In, []string{IPAMLoopbackPurpose})
	selector = selector.Add(*req)
	opts := &client.ListOptions{
		LabelSelector: selector,
		Namespace:     ns,
		Limit:         100,
	}
	ips := &ipamv1alpha1.IPList{}
	if err := cl.List(ctx, ips, opts); err != nil {
		return false, err
	}
	for _, item := range ips.Items {
		if item.Spec.IP.Net == ip {
			return false, nil
		}
	}
	return true, nil
}

func createSubnets(
	ctx context.Context,
	cl client.Client,
	obj *switchv1beta1.Switch,
	requiredAFFound map[ipamv1alpha1.SubnetAddressType]*bool,
) error {
	config := &switchv1beta1.SwitchConfig{}
	key := types.NamespacedName{Name: obj.Status.ConfigRef.Name, Namespace: obj.Namespace}
	if err := cl.Get(ctx, key, config); err != nil {
		return err
	}
	carrierSubnets := &ipamv1alpha1.SubnetList{}
	if err := listIPAMObjects(ctx, cl, obj, config.Spec.IPAM.CarrierSubnets, carrierSubnets); err != nil {
		return err
	}
	labelsToApply, err := resultingLabels(obj, obj.Spec.IPAM.GetSubnetsSelection(), config.Spec.IPAM.SouthSubnets)
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
		subnet, err := buildSubnetObject(ctx, cl, obj, item, labelsToApply)
		if err != nil {
			return err
		}
		if err := cl.Create(ctx, subnet); err != nil {
			if !apierrors.IsAlreadyExists(err) {
				return err
			}
		}
	}
	return nil
}

func buildSubnetObject(
	ctx context.Context,
	cl client.Client,
	obj *switchv1beta1.Switch,
	item ipamv1alpha1.Subnet,
	labelsToApply map[string]string,
) (*ipamv1alpha1.Subnet, error) {
	addressesRequired := getTotalAddressesCount(obj.Status.Interfaces, item.Status.Type)
	proposedCIDR, err := item.ProposeForCapacity(addressesRequired)
	if err != nil {
		return nil, err
	}
	ok, err := checkSubnetAvailable(ctx, cl, proposedCIDR, obj.Namespace)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, reconciliationError(ErrorDuplicateSubnet)
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

func checkSubnetAvailable(ctx context.Context, cl client.Client, cidr *ipamv1alpha1.CIDR, ns string) (bool, error) {
	selector := labels.NewSelector()
	req, _ := labels.NewRequirement(IPAMObjectPurposeLabel, selection.In, []string{IPAMSouthSubnetPurpose})
	selector = selector.Add(*req)
	opts := &client.ListOptions{
		LabelSelector: selector,
		Namespace:     ns,
		Limit:         100,
	}
	subnets := &ipamv1alpha1.SubnetList{}
	if err := cl.List(ctx, subnets, opts); err != nil {
		return false, err
	}
	for _, item := range subnets.Items {
		if cidr.String() == item.Spec.CIDR.String() {
			return false, nil
		}
	}
	return true, nil
}

func listIPAMObjects(
	ctx context.Context,
	cl client.Client,
	obj *switchv1beta1.Switch,
	params *switchv1beta1.IPAMSelectionSpec,
	list client.ObjectList,
) error {
	selector, err := getSelectorFromIPAMSpec(obj, params)
	if err != nil {
		return err
	}
	opts := &client.ListOptions{
		LabelSelector: selector,
		Namespace:     obj.Namespace,
		Limit:         100,
	}
	if err := cl.List(ctx, list, opts); err != nil {
		return err
	}
	return nil
}

//nolint:unused
func (r *SwitchReconciler) createIPs(
	ctx context.Context, obj *switchv1beta1.Switch, nic string, ips []*switchv1beta1.IPAddressSpec) error {
	for _, item := range ips {
		prefix, _ := netaddr.ParseIPPrefix(item.GetAddress())
		addr := prefix.IP()
		ip := &ipamv1alpha1.IP{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("%s-%s-%s",
					obj.Name, strings.ToLower(item.GetAddressFamily()), strings.ToLower(nic)),
				Namespace: obj.Namespace,
				Labels: map[string]string{
					IPAMObjectPurposeLabel: IPAMSwitchPortPurpose,
					IPAMObjectOwnerLabel:   obj.Name,
					IPAMObjectNICNameLabel: nic,
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
		if err := r.Create(ctx, ip); err != nil {
			if !apierrors.IsAlreadyExists(err) {
				return err
			}
		}
	}
	return nil
}
