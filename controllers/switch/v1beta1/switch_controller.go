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
	"strings"

	"github.com/go-logr/logr"
	ipamv1alpha1 "github.com/onmetal/ipam/api/v1alpha1"
	inventoryv1alpha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	switchv1beta1 "github.com/onmetal/metal-api/apis/switch/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const CIndexedUUID = "spec.uuid"

const (
	CEventTypeNormal = "Normal"

	CSubnetsNotFoundReason   = "NoCorrespondingSubnets"
	CIPsNotFoundReason       = "NoCorrespondingIPs"
	CAdditionalIPAddedReason = "AdditionalIPAdded"

	CInventoryNotFoundMessage    = "no corresponding inventory was found"
	CSwitchConfigNotFoundMessage = "no corresponding switch config was found"
	CSubnetsNotFoundMessage      = "no corresponding subnets was found"
	CIPsNotFoundMessage          = "no corresponding ips was found"
	CAdditionalIPAddedMessage    = "additional ip was added to interface %s"
)

// SwitchReconciler reconciles Switch object corresponding
// to given Inventory object
type SwitchReconciler struct {
	client.Client

	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=switch.onmetal.de,resources=switches,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=switch.onmetal.de,resources=switches/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=switch.onmetal.de,resources=switches/finalizers,verbs=update
//+kubebuilder:rbac:groups=switch.onmetal.de,resources=switchconfigs,verbs=get;list;watch
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=inventories,verbs=get;list;watch
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=inventories/status,verbs=get;list;watch
//+kubebuilder:rbac:groups=ipam.onmetal.de,resources=subnets,verbs=get;list;watch
//+kubebuilder:rbac:groups=ipam.onmetal.de,resources=subnets/status,verbs=get;list;watch
//+kubebuilder:rbac:groups=ipam.onmetal.de,resources=ips,verbs=get;list;watch
//+kubebuilder:rbac:groups=ipam.onmetal.de,resources=ips/status,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *SwitchReconciler) Reconcile(ctx context.Context, req ctrl.Request) (result ctrl.Result, err error) {
	log := r.Log.WithValues("switch", req.NamespacedName)
	result = ctrl.Result{}

	obj := &switchv1beta1.Switch{}
	if err = r.Get(ctx, req.NamespacedName, obj); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info(
				"requested resource not found",
				"name", req.NamespacedName, "kind", "Switch",
			)
		} else {
			log.Error(
				err,
				"failed to get requested resource",
				"name", req.NamespacedName, "kind", "Switch",
			)
		}
		return result, client.IgnoreNotFound(err)
	}
	if !obj.Spec.Managed {
		return
	}

	relatedInventory := &inventoryv1alpha1.Inventory{}
	if err = r.Get(ctx, types.NamespacedName{Namespace: obj.Namespace, Name: obj.Spec.UUID}, relatedInventory); err != nil {
		if !apierrors.IsNotFound(err) {
			log.Error(err, "failed to get resource", "name", obj.Spec.UUID, "kind", "Inventory")
			return
		}
		log.Info("related inventory object not found")
		return
	}

	if !obj.LabelsOK() {
		obj.UpdateSwitchLabels(relatedInventory)
		obj.UpdateSwitchAnnotations(relatedInventory)
		if err = r.Update(ctx, obj); err != nil {
			log.Error(err, "failed to update resource status", "name", obj.Name, "kind", obj.Kind)
			return
		}
		return
	}

	if obj.Status.SwitchState == nil {
		obj.SetInitialStatus(relatedInventory)
		if err = r.Status().Update(ctx, obj); err != nil {
			log.Error(err, "failed to update resource status", "name", obj.Name, "kind", obj.Kind)
			return
		}
		return
	}
	if !obj.StateEqualsTo(switchv1beta1.CSwitchStateInitial) {
		if !obj.InterfacesMatchInventory(relatedInventory) {
			obj.Status.SwitchState = nil
			if err = r.Status().Update(ctx, obj); err != nil {
				log.Error(err, "failed to update resource status", "name", obj.Name, "kind", obj.Kind)
				return
			}
			return
		}
	}

	var relatedConfig *switchv1beta1.SwitchConfig
	switchConfigs := &switchv1beta1.SwitchConfigList{}
	typeLabel, ok := obj.Labels[switchv1beta1.SwitchTypeLabel]
	if !ok {
		typeLabel = "all"
	}
	labelsReq, _ := switchv1beta1.GetLabelSelector(switchv1beta1.SwitchConfigTypeLabel+typeLabel, selection.Exists, []string{})
	selector := labels.NewSelector().Add(*labelsReq)
	opts := &client.ListOptions{
		LabelSelector: selector,
		Limit:         100,
	}
	if err = r.List(ctx, switchConfigs, opts); err != nil {
		log.Error(err, "failed to list resources", "kind", "SwitchConfigList")
		return
	}
	if len(switchConfigs.Items) > 0 {
		relatedConfig = &switchConfigs.Items[0]
	}
	obj.UpdateInterfacesParameters(relatedConfig)

	relatedSwitches := &switchv1beta1.SwitchList{}
	if err = r.List(ctx, relatedSwitches); err != nil {
		log.Error(err, "failed to list resources", "kind", "SwitchList")
		return
	}
	if !obj.ConnectionsOK(relatedSwitches) {
		obj.SetState(switchv1beta1.CSwitchStateProcessing)
		obj.SetConnections(relatedSwitches)
		if err = r.Status().Update(ctx, obj); err != nil {
			log.Error(err, "failed to update resource status", "name", obj.Name, "kind", obj.Kind)
			return
		}
		return
	}

	if !obj.StateEqualsTo(switchv1beta1.CSwitchStateReady) {
		obj.SetRole()
		obj.SetState(switchv1beta1.CSwitchStateReady)
		if err = r.Status().Update(ctx, obj); err != nil {
			log.Error(err, "failed to update resource status", "name", obj.Name, "kind", obj.Kind)
			return
		}
		return
	}

	if err = obj.ResultingIPAMConfig(relatedConfig); err != nil {
		log.Error(err, "failed to compute resulting IPAM configuration", "name", obj.Name, "kind", obj.Kind)
		return
	}

	if !obj.LoopbackSelectorsExist() {
		return
	}

	loopbackIPs, err := r.getRelatedLoopbackIPs(ctx, obj)
	if err != nil {
		log.Error(err, "failed to get list of related loopback IPs")
		return
	}
	//todo: implement IP objects creation in case there are no pre-created IPs
	// if existingAddressesDontMatchLoopbackSubnets() {
	// 	 r.createAbsentLoopbacks()
	// 	 return
	// }
	loopbacksMatchStored := obj.LoopbackIPsMatchStoredIPs(loopbackIPs)
	if !loopbacksMatchStored {
		if err = r.computeLoopbacks(ctx, obj, loopbackIPs); err != nil {
			log.Error(err, "failed to configure loopback addresses", "name", obj.Name, "kind", obj.Kind)
			return
		}
		if err = r.Status().Update(ctx, obj); err != nil {
			log.Error(err, "failed to update resource status", "name", obj.Name, "kind", obj.Kind)
			return
		}
		return
	}

	if !obj.SubnetSelectorsExist() {
		return
	}

	southSubnets, err := r.getRelatedSubnets(ctx, obj)
	if err != nil {
		log.Error(err, "failed to get list of related subnets")
	}
	//todo: implement subnet objects creation in case there are no pre-created subnets
	// if existingSubnetsDontMatchSwitchRangesSubnets() {
	// 	 r.createAbsentSubnets()
	// 	 return
	// }
	subnetsMatchStored := obj.SubnetsMatchStored(southSubnets)
	if !subnetsMatchStored {
		if err = r.computeSubnets(ctx, obj, southSubnets); err != nil {
			log.Error(err, "failed to configure south subnets", "name", obj.Name, "kind", obj.Kind)
			return
		}
		if err = r.Status().Update(ctx, obj); err != nil {
			log.Error(err, "failed to update resource status", "name", obj.Name, "kind", obj.Kind)
			return
		}
		return
	}

	ipAddressesOK := obj.IPaddressesOK(relatedSwitches)
	if !ipAddressesOK {
		if err = r.computeIPAddresses(ctx, obj, relatedSwitches); err != nil {
			log.Error(err, "failed to configure ip addresses", "name", obj.Name, "kind", obj.Kind)
			return
		}
		if err = r.Status().Update(ctx, obj); err != nil {
			log.Error(err, "failed to update resource status", "name", obj.Name, "kind", obj.Kind)
			return
		}
	}

	return
}

// SetupWithManager sets up the controller with the Manager.
func (r *SwitchReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// indexing switch's .spec.uuid field
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &switchv1beta1.Switch{}, CIndexedUUID, func(raw client.Object) []string {
		var res []string
		obj := raw.(*switchv1beta1.Switch)
		uuid := obj.Spec.UUID
		if uuid == "" {
			return res
		}
		res = append(res, obj.Spec.UUID)
		return res
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&switchv1beta1.Switch{}).
		Watches(&source.Kind{Type: &switchv1beta1.Switch{}}, &handler.Funcs{
			UpdateFunc: r.handleSwitchUpdate,
		}).
		// watches for inventories to handle UPDATE events of corresponding objects
		// on UPDATE event: if interfaces changed enqueue corresponding switch otherwise do nothing
		Watches(&source.Kind{Type: &inventoryv1alpha1.Inventory{}}, &handler.Funcs{
			UpdateFunc: r.handleInventoryUpdate,
		}).
		// watches for switchconfigs to handle CREATE and UPDATE events of objects that referencing labels
		// corresponding to existing on current switch object
		// on CREATE: enqueue all corresponding switches
		// on UPDATE: enqueue all corresponding switches
		Watches(&source.Kind{Type: &switchv1beta1.SwitchConfig{}}, &handler.Funcs{
			CreateFunc: r.handleSwitchConfigCreate,
			UpdateFunc: r.handleSwitchConfigUpdate,
		}).
		// watches for subnets to handle CREATE and UPDATE events of objects that referenced by switch objects
		// on CREATE: enqueue corresponding switch(es) if exist
		// on UPDATE: if labels was updated enqueue corresponding switch(es) if exist
		Watches(&source.Kind{Type: &ipamv1alpha1.Subnet{}}, &handler.Funcs{
			CreateFunc: r.handleSubnetCreate,
			UpdateFunc: r.handleSubnetUpdate,
		}).
		// watches for ips to handle CREATE and UPDATE events of object that referenced by switch objects
		// on CREATE: enqueue corresponding switch if exists
		// on UPDATE: if labels was updated enqueue corresponding switch if exists
		Watches(&source.Kind{Type: &ipamv1alpha1.IP{}}, &handler.Funcs{
			CreateFunc: r.handleIPCreate,
			UpdateFunc: r.handleIPUpdate,
		}).
		Complete(r)
}

func (r *SwitchReconciler) handleSwitchUpdate(e event.UpdateEvent, q workqueue.RateLimitingInterface) {
	prevObj := e.ObjectOld.(*switchv1beta1.Switch)
	currObj := e.ObjectNew.(*switchv1beta1.Switch)
	switchesQueue := make(map[string]struct{})
	for _, nicData := range prevObj.Status.Interfaces {
		if nicData.Peer == nil {
			continue
		}
		if nicData.Peer.ObjectReference == nil {
			continue
		}
		if nicData.Peer.PeerInfoSpec == nil {
			continue
		}
		if nicData.Peer.PeerInfoSpec.Type != switchv1beta1.CPeerTypeSwitch {
			continue
		}
		switchesQueue[nicData.Peer.ObjectReference.Name] = struct{}{}
	}
	for _, nicData := range currObj.Status.Interfaces {
		if nicData.Peer == nil {
			continue
		}
		if nicData.Peer.ObjectReference == nil {
			continue
		}
		if nicData.Peer.PeerInfoSpec == nil {
			continue
		}
		if nicData.Peer.PeerInfoSpec.Type != switchv1beta1.CPeerTypeSwitch {
			continue
		}
		switchesQueue[nicData.Peer.ObjectReference.Name] = struct{}{}
	}
	for name := range switchesQueue {
		q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
			Namespace: currObj.Namespace,
			Name:      name,
		}})
	}
}

func (r *SwitchReconciler) handleInventoryUpdate(e event.UpdateEvent, q workqueue.RateLimitingInterface) {
	log := r.Log.WithValues("switch-controller-watch", "inventory update")
	prevObj := e.ObjectOld.(*inventoryv1alpha1.Inventory)
	currObj := e.ObjectNew.(*inventoryv1alpha1.Inventory)
	existingLabels := currObj.GetLabels()
	if len(existingLabels) == 0 {
		return
	}
	sizeLabel := inventoryv1alpha1.GetSizeMatchLabel(switchv1beta1.CSwitchSizeName)
	if _, ok := existingLabels[sizeLabel]; !ok {
		return
	}
	if reflect.DeepEqual(prevObj.Spec.NICs, currObj.Spec.NICs) {
		return
	}
	switches := &switchv1beta1.SwitchList{}
	labelsReq, _ := switchv1beta1.GetLabelSelector(switchv1beta1.InventoryRefLabel, selection.Equals, []string{currObj.Name})
	selector := labels.NewSelector().Add(*labelsReq)
	opts := &client.ListOptions{
		LabelSelector: selector,
		Limit:         100,
	}
	if err := r.List(context.Background(), switches, opts); err != nil {
		log.Error(err, "failed to list resources", "kind", "SwitchList")
		return
	}
	if len(switches.Items) == 0 {
		return
	}
	q.Add(reconcile.Request{NamespacedName: switches.Items[0].GetNamespacedName()})
}

func (r *SwitchReconciler) handleSwitchConfigCreate(e event.CreateEvent, q workqueue.RateLimitingInterface) {
	obj := e.Object.(*switchv1beta1.SwitchConfig)
	r.enqueueDependenciesOnSwitchConfigEvent(obj, q)
}

func (r *SwitchReconciler) handleSwitchConfigUpdate(e event.UpdateEvent, q workqueue.RateLimitingInterface) {
	log := r.Log.WithValues("switch-controller-watch", "switchconfig update")
	prevObj := e.ObjectOld.(*switchv1beta1.SwitchConfig)
	currObj := e.ObjectNew.(*switchv1beta1.SwitchConfig)
	if reflect.DeepEqual(prevObj.Labels, currObj.Labels) {
		// assume that changes in spec, so it's needed to put related switches into reconciliation queue
		r.enqueueDependenciesOnSwitchConfigEvent(currObj, q)
		return
	}
	switches := &switchv1beta1.SwitchList{}
	_, allInPrev := prevObj.Labels[switchv1beta1.SwitchConfigTypeLabel+"all"]
	_, allInCurr := prevObj.Labels[switchv1beta1.SwitchConfigTypeLabel+"all"]
	// if type-all in previous or current labels then need to reconcile all switches
	switch allInPrev || allInCurr {
	case true:
		if err := r.List(context.Background(), switches); err != nil {
			log.Error(err, "failed to list resources", "kind", "SwitchList")
			return
		}
	case false:
		relatedTypes := make([]string, 0)
		swcTypeLabelPrefix := switchv1beta1.SwitchConfigTypeLabel
		for k := range prevObj.Labels {
			_, labelExistsInCurr := currObj.Labels[k]
			if strings.Contains(k, swcTypeLabelPrefix) && !labelExistsInCurr {
				ltype := strings.ReplaceAll(k, swcTypeLabelPrefix, "")
				relatedTypes = append(relatedTypes, ltype)
			}
		}
		for k := range currObj.Labels {
			_, labelExistsInCurr := prevObj.Labels[k]
			if strings.Contains(k, swcTypeLabelPrefix) && !labelExistsInCurr {
				ltype := strings.ReplaceAll(k, swcTypeLabelPrefix, "")
				relatedTypes = append(relatedTypes, ltype)
			}
		}
		swTypeLabelKey := switchv1beta1.SwitchTypeLabel
		labelsReq, _ := switchv1beta1.GetLabelSelector(swTypeLabelKey, selection.In, relatedTypes)
		selector := labels.NewSelector().Add(*labelsReq)
		opts := &client.ListOptions{
			LabelSelector: selector,
			Limit:         100,
		}
		if err := r.List(context.Background(), switches, opts); err != nil {
			log.Error(err, "failed to list resources", "kind", "SwitchList")
			return
		}
	}

	for _, item := range switches.Items {
		q.Add(reconcile.Request{NamespacedName: item.GetNamespacedName()})
	}
}

func (r *SwitchReconciler) enqueueDependenciesOnSwitchConfigEvent(obj *switchv1beta1.SwitchConfig, q workqueue.RateLimitingInterface) {
	//todo: rewrite using switchConfig.Spec.Switches label selector
	log := r.Log.WithValues("switch-controller-watch", "switchconfig create or update")
	relatedTypes := make([]string, 0)
	swcTypeLabelPrefix := switchv1beta1.SwitchConfigTypeLabel
	configForAll := false
	for k := range obj.Labels {
		if strings.Contains(k, swcTypeLabelPrefix) {
			ltype := strings.ReplaceAll(k, swcTypeLabelPrefix, "")
			if ltype == "all" {
				configForAll = true
				break
			}
			relatedTypes = append(relatedTypes, ltype)
		}
	}

	switches := &switchv1beta1.SwitchList{}
	switch configForAll {
	case true:
		if err := r.List(context.Background(), switches); err != nil {
			log.Error(err, "failed to list resources", "kind", "SwitchList")
			return
		}
	case false:
		swTypeLabelKey := switchv1beta1.SwitchTypeLabel
		labelsReq, _ := switchv1beta1.GetLabelSelector(swTypeLabelKey, selection.In, relatedTypes)
		selector := labels.NewSelector().Add(*labelsReq)
		opts := &client.ListOptions{
			LabelSelector: selector,
			Limit:         100,
		}
		if err := r.List(context.Background(), switches, opts); err != nil {
			log.Error(err, "failed to list resources", "kind", "SwitchList")
			return
		}
	}
	for _, item := range switches.Items {
		q.Add(reconcile.Request{NamespacedName: item.GetNamespacedName()})
	}
}

func (r *SwitchReconciler) handleSubnetCreate(e event.CreateEvent, q workqueue.RateLimitingInterface) {
	log := r.Log.WithValues("switch-controller-watch", "subnet create")
	var opts *client.ListOptions
	var switches *switchv1beta1.SwitchList
	var subnetOwnerLabel string
	var subnetOwnerLabelExists bool
	subnet := e.Object.(*ipamv1alpha1.Subnet)
	switches = &switchv1beta1.SwitchList{}
	subnetPurposeLabel, subnetPurposeLabelExists := subnet.Labels[switchv1beta1.IPAMObjectPurposeLabel]
	if subnetPurposeLabelExists {
		if subnetPurposeLabel == switchv1beta1.CIPAMPurposeSouthSubnet {
			subnetOwnerLabel, subnetOwnerLabelExists = subnet.Labels[switchv1beta1.IPAMObjectOwnerLabel]
		}
	}
	if subnetOwnerLabelExists {
		opts = &client.ListOptions{
			FieldSelector: fields.SelectorFromSet(map[string]string{CIndexedUUID: subnetOwnerLabel}),
			Limit:         100,
		}
		if err := r.List(context.Background(), switches, opts); err != nil {
			log.Error(err, "failed to list resources", "kind", "SwitchList")
			return
		}
		for _, item := range switches.Items {
			q.Add(reconcile.Request{NamespacedName: item.GetNamespacedName()})
		}
	}
}

func (r *SwitchReconciler) handleSubnetUpdate(e event.UpdateEvent, q workqueue.RateLimitingInterface) {
	log := r.Log.WithValues("switch-controller-watch", "subnet update")
	var opts *client.ListOptions
	var switches *switchv1beta1.SwitchList
	var prevOwnerLabel, currOwnerLabel string
	var prevOwnerLabelExists, currOwnerLabelExists bool
	prevObj := e.ObjectOld.(*ipamv1alpha1.Subnet)
	currObj := e.ObjectNew.(*ipamv1alpha1.Subnet)
	if reflect.DeepEqual(prevObj.Labels, currObj.Labels) {
		return
	}
	switches = &switchv1beta1.SwitchList{}
	prevPurposeLabel, prevPurposeLabelExists := prevObj.Labels[switchv1beta1.IPAMObjectPurposeLabel]
	if prevPurposeLabelExists {
		if prevPurposeLabel == switchv1beta1.CIPAMPurposeSouthSubnet {
			prevOwnerLabel, prevOwnerLabelExists = prevObj.Labels[switchv1beta1.IPAMObjectOwnerLabel]
		}
	}
	if prevOwnerLabelExists {
		opts = &client.ListOptions{
			FieldSelector: fields.SelectorFromSet(map[string]string{CIndexedUUID: prevOwnerLabel}),
			Limit:         100,
		}
		if err := r.List(context.Background(), switches, opts); err != nil {
			log.Error(err, "failed to list resources", "kind", "SwitchList")
			return
		}
		for _, item := range switches.Items {
			q.Add(reconcile.Request{NamespacedName: item.GetNamespacedName()})
		}
	}

	currPurposeLabel, currPurposeLabelExists := currObj.Labels[switchv1beta1.IPAMObjectPurposeLabel]
	if currPurposeLabelExists {
		if currPurposeLabel == switchv1beta1.CIPAMPurposeSouthSubnet {
			currOwnerLabel, currOwnerLabelExists = currObj.Labels[switchv1beta1.IPAMObjectOwnerLabel]
		}
	}
	if currOwnerLabelExists {
		opts = &client.ListOptions{
			FieldSelector: fields.SelectorFromSet(map[string]string{CIndexedUUID: currOwnerLabel}),
			Limit:         100,
		}
		if err := r.List(context.Background(), switches, opts); err != nil {
			log.Error(err, "failed to list resources", "kind", "SwitchList")
			return
		}
		for _, item := range switches.Items {
			q.Add(reconcile.Request{NamespacedName: item.GetNamespacedName()})
		}
	}
}

func (r *SwitchReconciler) handleIPCreate(e event.CreateEvent, q workqueue.RateLimitingInterface) {
	log := r.Log.WithValues("switch-controller-watch", "ip create")
	var opts *client.ListOptions
	var switches *switchv1beta1.SwitchList
	var ipOwnerLabel string
	var ipOwnerLabelExists bool
	ip := e.Object.(*ipamv1alpha1.IP)
	switches = &switchv1beta1.SwitchList{}
	ipPurposeLabel, ipPurposeLabelExists := ip.Labels[switchv1beta1.IPAMObjectPurposeLabel]
	if ipPurposeLabelExists {
		if ipPurposeLabel == switchv1beta1.CIPAMPurposeLoopback ||
			ipPurposeLabel == switchv1beta1.CIPAMPurposeInterfaceIP {

			ipOwnerLabel, ipOwnerLabelExists = ip.Labels[switchv1beta1.IPAMObjectOwnerLabel]
		}
	}
	if ipOwnerLabelExists {
		opts = &client.ListOptions{
			FieldSelector: fields.SelectorFromSet(map[string]string{CIndexedUUID: ipOwnerLabel}),
			Limit:         100,
		}
		if err := r.List(context.Background(), switches, opts); err != nil {
			log.Error(err, "failed to list resources", "kind", "SwitchList")
			return
		}
		for _, item := range switches.Items {
			q.Add(reconcile.Request{NamespacedName: item.GetNamespacedName()})
		}
	}
}

func (r *SwitchReconciler) handleIPUpdate(e event.UpdateEvent, q workqueue.RateLimitingInterface) {
	log := r.Log.WithValues("switch-controller-watch", "ip update")
	var opts *client.ListOptions
	var switches *switchv1beta1.SwitchList
	var prevOwnerLabel, currOwnerLabel string
	var prevOwnerLabelExists, currOwnerLabelExists bool
	prevObj := e.ObjectOld.(*ipamv1alpha1.IP)
	currObj := e.ObjectNew.(*ipamv1alpha1.IP)
	if reflect.DeepEqual(prevObj.Labels, currObj.Labels) {
		return
	}
	switches = &switchv1beta1.SwitchList{}
	prevPurposeLabel, prevPurposeLabelExists := prevObj.Labels[switchv1beta1.IPAMObjectPurposeLabel]
	if prevPurposeLabelExists {
		if prevPurposeLabel == switchv1beta1.CIPAMPurposeLoopback ||
			prevPurposeLabel == switchv1beta1.CIPAMPurposeInterfaceIP {

			prevOwnerLabel, prevOwnerLabelExists = prevObj.Labels[switchv1beta1.IPAMObjectOwnerLabel]
		}
	}
	if prevOwnerLabelExists {
		opts = &client.ListOptions{
			FieldSelector: fields.SelectorFromSet(map[string]string{CIndexedUUID: prevOwnerLabel}),
			Limit:         100,
		}
		if err := r.List(context.Background(), switches, opts); err != nil {
			log.Error(err, "failed to list resources", "kind", "SwitchList")
			return
		}
		for _, item := range switches.Items {
			q.Add(reconcile.Request{NamespacedName: item.GetNamespacedName()})
		}
	}

	currPurposeLabel, currPurposeLabelExists := currObj.Labels[switchv1beta1.IPAMObjectPurposeLabel]
	if currPurposeLabelExists {
		if currPurposeLabel == switchv1beta1.CIPAMPurposeLoopback ||
			currPurposeLabel == switchv1beta1.CIPAMPurposeInterfaceIP {

			currOwnerLabel, currOwnerLabelExists = currObj.Labels[switchv1beta1.IPAMObjectOwnerLabel]
		}
	}
	if currOwnerLabelExists {
		opts = &client.ListOptions{
			FieldSelector: fields.SelectorFromSet(map[string]string{CIndexedUUID: currOwnerLabel}),
			Limit:         100,
		}
		if err := r.List(context.Background(), switches, opts); err != nil {
			log.Error(err, "failed to list resources", "kind", "SwitchList")
			return
		}
		for _, item := range switches.Items {
			q.Add(reconcile.Request{NamespacedName: item.GetNamespacedName()})
		}
	}
}

func (r *SwitchReconciler) getRelatedLoopbackIPs(ctx context.Context, obj *switchv1beta1.Switch) (list *ipamv1alpha1.IPList, err error) {
	list = &ipamv1alpha1.IPList{}
	selector := labels.NewSelector()
	for key, value := range obj.Spec.IPAM.LoopbackAddresses.LabelSelector.MatchLabels {
		req, err := switchv1beta1.GetLabelSelector(
			key,
			selection.Equals,
			[]string{value},
		)
		if err != nil {
			return list, err
		}
		selector = selector.Add(*req)
	}
	opts := &client.ListOptions{
		LabelSelector: selector,
		Namespace:     obj.Namespace,
		Limit:         100,
	}
	if err = r.List(ctx, list, opts); err != nil {
		return
	}
	return
}

func (r *SwitchReconciler) computeLoopbacks(ctx context.Context, obj *switchv1beta1.Switch, list *ipamv1alpha1.IPList) (err error) {
	loopbacks := make([]*switchv1beta1.IPAddressSpec, 0)
	for _, item := range list.Items {
		loopbacks = append(loopbacks, &switchv1beta1.IPAddressSpec{
			ObjectReference: &switchv1beta1.ObjectReference{
				Name:      item.Name,
				Namespace: item.Namespace,
			},
			Address:      item.Status.Reserved.String(),
			ExtraAddress: false,
		})
	}
	obj.Status.LoopbackAddresses = loopbacks
	return
}

func (r *SwitchReconciler) getRelatedSubnets(ctx context.Context, obj *switchv1beta1.Switch) (list *ipamv1alpha1.SubnetList, err error) {
	list = &ipamv1alpha1.SubnetList{}
	selector := labels.NewSelector()
	for key, value := range obj.Spec.IPAM.SouthSubnets.LabelSelector.MatchLabels {
		req, err := switchv1beta1.GetLabelSelector(
			key,
			selection.Equals,
			[]string{value},
		)
		if err != nil {
			return list, err
		}
		selector = selector.Add(*req)
	}
	opts := &client.ListOptions{
		LabelSelector: selector,
		Namespace:     obj.Namespace,
		Limit:         100,
	}
	if err = r.List(ctx, list, opts); err != nil {
		return
	}
	return
}

func (r *SwitchReconciler) computeSubnets(ctx context.Context, obj *switchv1beta1.Switch, list *ipamv1alpha1.SubnetList) (err error) {
	subnets := make([]*switchv1beta1.SubnetSpec, 0)
	for _, item := range list.Items {
		subnets = append(subnets, &switchv1beta1.SubnetSpec{
			ObjectReference: &switchv1beta1.ObjectReference{
				Name:      item.Name,
				Namespace: item.Namespace,
			},
			CIDR: item.Status.Reserved.String(),
		})
	}
	obj.Status.Subnets = subnets
	return
}

func (r *SwitchReconciler) computeIPAddresses(ctx context.Context, obj *switchv1beta1.Switch, list *switchv1beta1.SwitchList) error {
	//todo: creation of related nics' subnets and ips
	extraIPs := obj.GetExtraNICsIPs()
	southIPs, err := obj.GetSouthNICsIP()
	if err != nil {
		return err
	}
	northIPs := obj.GetNorthNICsIP(list)
	for nic, nicData := range obj.Status.Interfaces {
		resultingIPs := make([]*switchv1beta1.IPAddressSpec, 0)
		if ips, ok := extraIPs[nic]; ok {
			resultingIPs = append(resultingIPs, ips...)
		}
		if ips, ok := southIPs[nic]; ok {
			resultingIPs = append(resultingIPs, ips...)
		}
		if ips, ok := northIPs[nic]; ok {
			resultingIPs = append(resultingIPs, ips...)
		}
		nicData.IP = resultingIPs
	}
	return nil
}
