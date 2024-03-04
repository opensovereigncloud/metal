// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package v1beta1

import (
	"context"
	"crypto/md5"
	"fmt"
	"net"
	"strings"

	"github.com/go-logr/logr"
	ipamv1alpha1 "github.com/ironcore-dev/ipam/api/ipam/v1alpha1"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/tools/reference"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	metalv1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"
	"github.com/ironcore-dev/metal/pkg/constants"
	switchespkg "github.com/ironcore-dev/metal/pkg/switches"
)

// IPAMReconciler reconciles NetworkSwitch object
// and creates required IPAM objects.
type IPAMReconciler struct {
	client.Client

	Log                     logr.Logger
	Scheme                  *runtime.Scheme
	Recorder                record.EventRecorder
	SwitchPortsIPAMDisabled bool
}

// +kubebuilder:rbac:groups=metal.ironcore.dev,resources=networkswitches,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=metal.ironcore.dev,resources=networkswitches/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=ipam.ironcore.dev,resources=subnets,verbs=get;list;watch;create;update;patch;delete;deletecollection
// +kubebuilder:rbac:groups=ipam.ironcore.dev,resources=subnets/status,verbs=get
// +kubebuilder:rbac:groups=ipam.ironcore.dev,resources=ips,verbs=get;list;watch;create;update;patch;delete;deletecollection
// +kubebuilder:rbac:groups=ipam.ironcore.dev,resources=ips/status,verbs=get

func (r *IPAMReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	obj := &metalv1alpha4.NetworkSwitch{}
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	ref, err := reference.GetReference(r.Scheme, obj)
	if err != nil {
		r.Log.Error(err, "failed to construct reference", "request", req.NamespacedName)
		return ctrl.Result{}, err
	}

	log := r.Log.WithValues("object", *ref)
	log.Info("reconciliation started")
	requestCtx := logr.NewContext(ctx, log)
	return r.reconcileRequired(requestCtx, obj)
}

func (r *IPAMReconciler) reconcileRequired(ctx context.Context, obj *metalv1alpha4.NetworkSwitch) (ctrl.Result, error) {
	if !obj.GetDeletionTimestamp().IsZero() {
		return ctrl.Result{}, nil
	}
	return r.reconcileManaged(ctx, obj)
}

func (r *IPAMReconciler) reconcileManaged(ctx context.Context, obj *metalv1alpha4.NetworkSwitch) (ctrl.Result, error) {
	if !obj.Managed() {
		log := logr.FromContextOrDiscard(ctx)
		log.WithValues("reason", constants.ReasonUnmanagedSwitch).Info("reconciliation finished")
		return ctrl.Result{}, nil
	}
	return r.reconcile(ctx, obj)
}

func (r *IPAMReconciler) reconcile(ctx context.Context, obj *metalv1alpha4.NetworkSwitch) (ctrl.Result, error) {
	log := logr.FromContextOrDiscard(ctx)
	svc := switchespkg.NewSwitchEnvironmentSvc(r.Client, log)
	env := svc.GetEnvironment(ctx, obj)
	if env.Config == nil {
		log.Info("no corresponding SwitchConfig object found")
		return ctrl.Result{}, nil
	}
	return r.reconcileCleanup(ctx, obj, svc)
}

func (r *IPAMReconciler) reconcileCleanup(
	ctx context.Context,
	obj *metalv1alpha4.NetworkSwitch,
	svc *switchespkg.SwitchEnvironmentSvc,
) (ctrl.Result, error) {
	log := logr.FromContextOrDiscard(ctx)
	log.Info("cleaning up failed loopback IP objects")
	if err := cleanupFailedLoopbackIPs(ctx, obj, svc); err != nil {
		return ctrl.Result{}, err
	}
	log.Info("cleaning up failed switch ports Subnet objects")
	if err := cleanupFailedSwitchPortSubnets(ctx, obj, svc); err != nil {
		return ctrl.Result{}, err
	}
	log.Info("cleaning up failed south Subnet objects")
	if err := cleanupFailedSouthSubnets(ctx, obj, svc); err != nil {
		return ctrl.Result{}, err
	}
	return r.reconcileIPAM(ctx, obj, svc)
}

func (r *IPAMReconciler) reconcileIPAM(
	ctx context.Context,
	obj *metalv1alpha4.NetworkSwitch,
	svc *switchespkg.SwitchEnvironmentSvc,
) (ctrl.Result, error) {
	log := logr.FromContextOrDiscard(ctx)
	env := svc.Env
	switch {
	case env.LoopbackIPs == nil:
		log.Info("processing loopbacks")
		if err := processLoopbacks(ctx, obj, svc); err != nil {
			log.Error(err, "failed to create loopback IP objects")
			return ctrl.Result{}, err
		}
	case env.SouthSubnets == nil:
		log.Info("processing south subnets")
		if err := processSouthSubnets(ctx, obj, svc); err != nil {
			log.Error(err, "failed to create Subnet objects")
			return ctrl.Result{}, err
		}
	case env.SwitchPortSubnets == nil:
		if r.SwitchPortsIPAMDisabled {
			break
		}
		log.Info("processing switch ports Subnet objects")
		if err := processSwitchPortsSubnets(ctx, obj, svc); err != nil {
			log.Error(err, "failed to create switch ports Subnet objects")
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *IPAMReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// label predicate to filter only NetworkSwitch object,
	// which were already onboarded
	labelPredicate, err := predicate.LabelSelectorPredicate(metav1.LabelSelector{
		MatchLabels: map[string]string{constants.InventoriedLabel: "true"},
	})
	if err != nil {
		r.Log.Error(err, "failed to setup predicates")
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&metalv1alpha4.NetworkSwitch{}, builder.WithPredicates(labelPredicate)).
		WithOptions(controller.Options{
			RecoverPanic: ptr.To(true),
		}).
		// watches for ipam.Subnet objects are required to trigger reconciliation
		// in case related ipam.Subnet object defining switch's loopbacks or carrier
		// subnet was updated or being deleted
		Watches(&ipamv1alpha1.Subnet{}, handler.Funcs{
			UpdateFunc: r.handleSubnetUpdateEvent,
			DeleteFunc: r.handleSubnetDeleteEvent,
		}).
		Complete(r)
}

func (r *IPAMReconciler) handleSubnetUpdateEvent(
	ctx context.Context,
	e event.UpdateEvent,
	q workqueue.RateLimitingInterface,
) {
	r.Log.WithValues("handler", "SubnetUpdateEvent")
	subnet, ok := e.ObjectNew.(*ipamv1alpha1.Subnet)
	if !ok {
		return
	}
	if subnet.Status.State != ipamv1alpha1.CFinishedSubnetState {
		return
	}
	switches := r.switchesToEnqueueOnSubnetEvent(ctx, subnet)
	if switches == nil {
		return
	}
	for _, item := range switches.Items {
		q.Add(reconcile.Request{NamespacedName: item.NamespacedName()})
	}
}

func (r *IPAMReconciler) handleSubnetDeleteEvent(ctx context.Context, e event.DeleteEvent, q workqueue.RateLimitingInterface) {
	r.Log.WithValues("handler", "SubnetDeleteEvent")
	subnet, ok := e.Object.(*ipamv1alpha1.Subnet)
	if !ok {
		return
	}
	if subnet.Status.State != ipamv1alpha1.CFinishedSubnetState {
		return
	}
	switches := r.switchesToEnqueueOnSubnetEvent(ctx, subnet)
	if switches == nil {
		return
	}
	for _, item := range switches.Items {
		q.Add(reconcile.Request{NamespacedName: item.NamespacedName()})
	}
}

func (r *IPAMReconciler) switchesToEnqueueOnSubnetEvent(
	ctx context.Context,
	subnet *ipamv1alpha1.Subnet,
) *metalv1alpha4.NetworkSwitchList {
	switches := &metalv1alpha4.NetworkSwitchList{}
	if err := r.List(ctx, switches); err != nil {
		r.Log.Error(err, "failed to list NetworkSwitch objects")
		return nil
	}
	switchConfigs := &metalv1alpha4.SwitchConfigList{}
	if err := r.List(ctx, switchConfigs); err != nil {
		r.Log.Error(err, "failed to list SwitchConfig objects")
		return nil
	}
	for _, item := range switchConfigs.Items {
		switch {
		case switchespkg.IPAMSelectorMatchLabels(nil, item.Spec.IPAM.CarrierSubnets, subnet.Labels):
			return switches
		case switchespkg.IPAMSelectorMatchLabels(nil, item.Spec.IPAM.LoopbackSubnets, subnet.Labels):
			return switches
		}
	}

	return nil
}

func processLoopbacks(
	ctx context.Context,
	obj *metalv1alpha4.NetworkSwitch,
	svc *switchespkg.SwitchEnvironmentSvc,
) error {
	counter := 0
	cfg := svc.Env.Config
	params := cfg.Spec.IPAM.LoopbackAddresses
	if obj.Spec.IPAM != nil && obj.Spec.IPAM.LoopbackAddresses != nil {
		params = obj.Spec.IPAM.LoopbackAddresses
	}
	loopbacks := &ipamv1alpha1.IPList{}
	if err := svc.ListIPAMObjects(ctx, obj, params, loopbacks); err != nil {
		return err
	}
	expectedAF := cfg.Spec.IPAM.AddressFamily
	existingAFFlag := existingLoopbacksAddressFamilies(loopbacks, expectedAF)
	if missedAFFlag, ok := switchespkg.AddressFamiliesMatchConfig(
		true, expectedAF.GetIPv6(), existingAFFlag); !ok {
		svc.Log.Info("discrepancy between required and existing IP objects in part of address families")
		created, err := createIPs(ctx, obj, svc, missedAFFlag)
		if err != nil {
			return err
		}
		counter += created
	}
	svc.Log.WithValues("count", counter).Info("loopback IP objects created")
	return nil
}

func existingLoopbacksAddressFamilies(
	list *ipamv1alpha1.IPList,
	af *metalv1alpha4.AddressFamiliesMap,
) int {
	afEnabledFlag := 0
	for _, item := range list.Items {
		if item.Status.State != ipamv1alpha1.CFinishedIPState {
			continue
		}
		if !af.GetIPv6() && item.Status.Reserved.Net.Is6() {
			continue
		}
		afEnabledFlag = switchespkg.ComputeAFFlag(
			item.Status.Reserved.Net.Is4(), item.Status.Reserved.Net.Is6(), afEnabledFlag)
	}
	return afEnabledFlag
}

func createIPs(
	ctx context.Context,
	obj *metalv1alpha4.NetworkSwitch,
	svc *switchespkg.SwitchEnvironmentSvc,
	afFlag int,
) (int, error) {
	counter := 0
	loopbackSubnets := &ipamv1alpha1.SubnetList{}
	params := svc.Env.Config.Spec.IPAM.LoopbackSubnets
	if err := svc.ListIPAMObjects(ctx, obj, params, loopbackSubnets); err != nil {
		return counter, err
	}
	labelsToApply, err := switchespkg.ResultingLabels(
		obj, obj.Spec.IPAM.GetLoopbacksSelection(), svc.Env.Config.Spec.IPAM.LoopbackAddresses)
	if err != nil {
		return counter, err
	}
	for _, item := range loopbackSubnets.Items {
		if item.Status.State != ipamv1alpha1.CFinishedSubnetState {
			continue
		}
		// check whether loopbacks subnet has free address
		if resource.NewQuantity(1, resource.DecimalSI).Cmp(item.Status.CapacityLeft) > 1 {
			continue
		}
		switch {
		case item.Status.Type == ipamv1alpha1.CIPv4SubnetType && afFlag == 2:
			continue
		case item.Status.Type == ipamv1alpha1.CIPv6SubnetType && afFlag == 1:
			continue
		}
		ip := buildIPObject(obj, item, labelsToApply)
		err = svc.Create(ctx, ip)
		switch {
		case apierrors.IsAlreadyExists(err):
			continue
		case err != nil:
			return counter, err
		default:
			counter += 1
		}
	}

	return counter, nil
}

func buildIPObject(
	obj *metalv1alpha4.NetworkSwitch,
	parent ipamv1alpha1.Subnet,
	labels map[string]string,
) *ipamv1alpha1.IP {
	ipObject := &ipamv1alpha1.IP{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-lo-%s", obj.Name, strings.ToLower(string(parent.Status.Type))),
			Namespace: obj.Namespace,
			Labels:    labels,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: obj.APIVersion,
					Kind:       obj.Kind,
					Name:       obj.Name,
					UID:        obj.UID},
			},
		},
		Spec: ipamv1alpha1.IPSpec{
			Subnet: v1.LocalObjectReference{Name: parent.Name},
			Consumer: &ipamv1alpha1.ResourceReference{
				APIVersion: obj.APIVersion,
				Kind:       obj.Kind,
				Name:       obj.Name,
			},
		},
	}
	return ipObject
}

func processSouthSubnets(
	ctx context.Context,
	obj *metalv1alpha4.NetworkSwitch,
	svc *switchespkg.SwitchEnvironmentSvc,
) error {
	counter := 0
	cfg := svc.Env.Config
	params := cfg.Spec.IPAM.SouthSubnets
	if obj.Spec.IPAM != nil && obj.Spec.IPAM.SouthSubnets != nil {
		params = obj.Spec.IPAM.SouthSubnets
	}
	subnets := &ipamv1alpha1.SubnetList{}
	if err := svc.ListIPAMObjects(ctx, obj, params, subnets); err != nil {
		return err
	}
	expectedAF := cfg.Spec.IPAM.AddressFamily
	existingAFFlag := existingSouthSubnetsAddressFamilies(obj, subnets, expectedAF)
	if missedAFFlag, ok := switchespkg.AddressFamiliesMatchConfig(
		expectedAF.GetIPv4(), expectedAF.GetIPv6(), existingAFFlag); !ok {
		svc.Log.Info("discrepancy between required and existing Subnet objects in part of address families")
		created, err := createSubnets(ctx, obj, svc, missedAFFlag)
		if err != nil {
			return err
		}
		counter += created
	}
	svc.Log.WithValues("count", counter).Info("south Subnet objects created")
	return nil
}

func existingSouthSubnetsAddressFamilies(
	obj *metalv1alpha4.NetworkSwitch,
	list *ipamv1alpha1.SubnetList,
	af *metalv1alpha4.AddressFamiliesMap,
) int {
	afEnabledFlag := 0
	for _, item := range list.Items {
		if item.Status.State != ipamv1alpha1.CFinishedSubnetState {
			continue
		}
		if (!af.GetIPv4() && item.Status.Reserved.IsIPv4()) || (!af.GetIPv6() && item.Status.Reserved.IsIPv6()) {
			continue
		}
		requiredCapacity := switchespkg.GetTotalAddressesCount(obj.Status.Interfaces, item.Status.Type)
		if requiredCapacity.Cmp(item.Status.Capacity) >= 0 {
			continue
		}
		afEnabledFlag = switchespkg.ComputeAFFlag(
			item.Status.Reserved.IsIPv4(), item.Status.Reserved.IsIPv6(), afEnabledFlag)
	}
	return afEnabledFlag
}

func createSubnets(
	ctx context.Context,
	obj *metalv1alpha4.NetworkSwitch,
	svc *switchespkg.SwitchEnvironmentSvc,
	afFlag int,
) (int, error) {
	counter := 0
	cfg := svc.Env.Config
	params := cfg.Spec.IPAM.CarrierSubnets
	carrierSubnets := &ipamv1alpha1.SubnetList{}
	if err := svc.ListIPAMObjects(ctx, obj, params, carrierSubnets); err != nil {
		return counter, err
	}
	labelsToApply, err := switchespkg.ResultingLabels(
		obj, obj.Spec.IPAM.GetSubnetsSelection(), svc.Env.Config.Spec.IPAM.SouthSubnets)
	if err != nil {
		return counter, err
	}
	for _, item := range carrierSubnets.Items {
		if item.Status.State != ipamv1alpha1.CFinishedSubnetState {
			continue
		}
		switch {
		case item.Status.Type == ipamv1alpha1.CIPv4SubnetType && afFlag == 2:
			continue
		case item.Status.Type == ipamv1alpha1.CIPv6SubnetType && afFlag == 1:
			continue
		}
		subnet := buildSubnetObject(obj, item, labelsToApply)
		err = svc.Create(ctx, subnet)
		switch {
		case apierrors.IsAlreadyExists(err):
			continue
		case err != nil:
			return counter, err
		default:
			counter += 1
		}
	}

	return counter, nil
}

func buildSubnetObject(
	obj *metalv1alpha4.NetworkSwitch,
	parent ipamv1alpha1.Subnet,
	labels map[string]string,
) *ipamv1alpha1.Subnet {
	addressesRequired := switchespkg.GetTotalAddressesCount(obj.Status.Interfaces, parent.Status.Type)
	subnet := &ipamv1alpha1.Subnet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-sn-%s", obj.Name, strings.ToLower(string(parent.Status.Type))),
			Namespace: obj.Namespace,
			Labels:    labels,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: obj.APIVersion,
					Kind:       obj.Kind,
					Name:       obj.Name,
					UID:        obj.UID},
			},
		},
		Spec: ipamv1alpha1.SubnetSpec{
			Capacity:     addressesRequired,
			ParentSubnet: v1.LocalObjectReference{Name: parent.Name},
			Network:      v1.LocalObjectReference{Name: parent.Spec.Network.Name},
			Consumer: &ipamv1alpha1.ResourceReference{
				APIVersion: obj.APIVersion,
				Kind:       obj.Kind,
				Name:       obj.Name,
			},
		},
	}
	return subnet
}

func processSwitchPortsSubnets(
	ctx context.Context,
	obj *metalv1alpha4.NetworkSwitch,
	svc *switchespkg.SwitchEnvironmentSvc,
) error {
	c := obj.GetCondition(constants.ConditionIPAddressesOK)
	if !c.GetState() {
		return nil
	}
	counter := 0
	for nic, data := range obj.Status.Interfaces {
		if data.GetDirection() == constants.DirectionNorth {
			continue
		}
		created, err := createSwitchPortSubnets(ctx, obj, svc, nic, data)
		if err != nil {
			return err
		}
		counter += created
	}
	svc.Log.WithValues("count", counter).Info("switch port Subnet objects created")
	return nil
}

func createSwitchPortSubnets(
	ctx context.Context,
	obj *metalv1alpha4.NetworkSwitch,
	svc *switchespkg.SwitchEnvironmentSvc,
	nic string,
	data *metalv1alpha4.InterfaceSpec,
) (int, error) {
	counter := 0
	for _, ip := range data.IP {
		if ip.GetExtraAddress() {
			continue
		}
		parent, err := subnetContainsIP(obj, ip)
		if err != nil {
			return counter, err
		}
		if parent == nil {
			continue
		}
		_, cidr, err := net.ParseCIDR(ip.GetAddress())
		if err != nil {
			return counter, err
		}
		ipamCIDR, err := ipamv1alpha1.CIDRFromString(cidr.String())
		if err != nil {
			return counter, err
		}
		subnet := buildSwitchPortSubnet(
			obj, ipamCIDR, nic, parent.GetSubnetObjectRefName(), parent.GetNetworkObjectRefName())
		err = svc.Create(ctx, subnet)
		switch {
		case apierrors.IsAlreadyExists(err):
			continue
		case err != nil:
			return counter, err
		default:
			counter += 1
		}
	}
	return counter, nil
}

func subnetContainsIP(obj *metalv1alpha4.NetworkSwitch, address *metalv1alpha4.IPAddressSpec) (*metalv1alpha4.SubnetSpec, error) {
	for _, subnet := range obj.Status.Subnets {
		if address.GetAddressFamily() != subnet.GetAddressFamily() {
			continue
		}
		ip, _, err := net.ParseCIDR(address.GetAddress())
		if err != nil {
			return nil, err
		}
		_, cidr, err := net.ParseCIDR(subnet.GetCIDR())
		if err != nil {
			return nil, err
		}
		if cidr.Contains(ip) {
			return subnet, nil
		}
	}
	return nil, nil
}

func buildSwitchPortSubnet(
	obj *metalv1alpha4.NetworkSwitch,
	cidr *ipamv1alpha1.CIDR,
	nic, parent, network string,
) *ipamv1alpha1.Subnet {
	hash := md5.Sum([]byte(cidr.String()))
	subnet := &ipamv1alpha1.Subnet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s-%x", obj.Name, strings.ToLower(nic), hash[:4]),
			Namespace: obj.Namespace,
			Labels: map[string]string{
				constants.IPAMObjectOwnerLabel:       obj.Name,
				constants.IPAMObjectPurposeLabel:     constants.IPAMSwitchPortPurpose,
				constants.IPAMObjectGeneratedByLabel: constants.SwitchManager,
				constants.IPAMObjectNICNameLabel:     nic,
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: obj.APIVersion,
					Kind:       obj.Kind,
					Name:       obj.Name,
					UID:        obj.UID,
				},
			},
		},
		Spec: ipamv1alpha1.SubnetSpec{
			CIDR:         cidr,
			ParentSubnet: v1.LocalObjectReference{Name: parent},
			Network:      v1.LocalObjectReference{Name: network},
			Consumer: &ipamv1alpha1.ResourceReference{
				APIVersion: obj.APIVersion,
				Kind:       obj.Kind,
				Name:       obj.Name,
			},
		},
	}
	return subnet
}

// nolint:dupl
func cleanupFailedLoopbackIPs(
	ctx context.Context,
	obj *metalv1alpha4.NetworkSwitch,
	svc *switchespkg.SwitchEnvironmentSvc,
) error {
	filterLabels, err := switchespkg.ResultingLabels(
		obj, obj.Spec.IPAM.GetLoopbacksSelection(), svc.Env.Config.Spec.IPAM.LoopbackAddresses)
	if err != nil {
		return err
	}
	selector := labels.NewSelector()
	for k, v := range filterLabels {
		req, _ := labels.NewRequirement(k, selection.In, []string{v})
		selector = selector.Add(*req)
	}
	opts := &client.ListOptions{
		LabelSelector: selector,
		Namespace:     obj.Namespace,
		Limit:         100,
	}
	ips := &ipamv1alpha1.IPList{}
	if err := svc.List(ctx, ips, opts); err != nil {
		return err
	}
	failedIPFound := 0
	for _, item := range ips.Items {
		if item.Status.State == ipamv1alpha1.CFailedIPState {
			failedIPFound += 1
			svc.Log.Info("cleaning up failed loopback IP", "ip", item.GetName, "state", item.Status.State)
			_ = svc.Delete(ctx, &item)
		}
	}
	if failedIPFound > 0 {
		svc.Log.WithValues("count", failedIPFound).
			Info("loopback IP objects in 'Failed' state discovered")
	}
	return nil
}

// nolint:dupl
func cleanupFailedSwitchPortSubnets(
	ctx context.Context,
	obj *metalv1alpha4.NetworkSwitch,
	svc *switchespkg.SwitchEnvironmentSvc,
) error {
	subnets := &ipamv1alpha1.SubnetList{}
	selector := labels.NewSelector()
	purposeReq, _ := labels.NewRequirement(constants.IPAMObjectPurposeLabel, selection.In, []string{constants.IPAMSwitchPortPurpose})
	ownerReq, _ := labels.NewRequirement(constants.IPAMObjectOwnerLabel, selection.In, []string{obj.Name})
	generatedReq, _ := labels.NewRequirement(constants.IPAMObjectGeneratedByLabel, selection.In, []string{constants.SwitchManager})
	selector = selector.Add(*purposeReq).Add(*ownerReq).Add(*generatedReq)
	opts := &client.ListOptions{
		LabelSelector: selector,
		Namespace:     obj.Namespace,
	}
	if err := svc.List(ctx, subnets, opts); err != nil {
		return err
	}
	failedSubnetsCount := 0
	for _, item := range subnets.Items {
		if item.Status.State == ipamv1alpha1.CFailedSubnetState {
			failedSubnetsCount += 1
			svc.Log.Info("cleaning up failed switch port Subnet", "subnet", item.GetName, "state", item.Status.State)
			_ = svc.Delete(ctx, &item)
		}
	}
	if failedSubnetsCount > 0 {
		svc.Log.WithValues("count", failedSubnetsCount).
			Info("switch port Subnet objects in 'Failed' state cleaned up")
	}
	return nil
}

// nolint:dupl
func cleanupFailedSouthSubnets(
	ctx context.Context,
	obj *metalv1alpha4.NetworkSwitch,
	svc *switchespkg.SwitchEnvironmentSvc,
) error {
	filterLabels, err := switchespkg.ResultingLabels(
		obj, obj.Spec.IPAM.GetSubnetsSelection(), svc.Env.Config.Spec.IPAM.SouthSubnets)
	if err != nil {
		return err
	}
	selector := labels.NewSelector()
	for k, v := range filterLabels {
		req, _ := labels.NewRequirement(k, selection.In, []string{v})
		selector = selector.Add(*req)
	}
	opts := &client.ListOptions{
		LabelSelector: selector,
		Namespace:     obj.Namespace,
		Limit:         100,
	}
	subnets := &ipamv1alpha1.SubnetList{}
	if err := svc.List(ctx, subnets, opts); err != nil {
		return err
	}
	failedSubnetsCount := 0
	for _, item := range subnets.Items {
		if item.Status.State == ipamv1alpha1.CFailedSubnetState {
			failedSubnetsCount += 1
			svc.Log.Info("cleaning up failed south Subnet", "subnet", item.GetName, "state", item.Status.State)
			_ = svc.Delete(ctx, &item)
		}
	}
	if failedSubnetsCount > 0 {
		svc.Log.WithValues("count", failedSubnetsCount).
			Info("south Subnet objects in 'Failed' state discovered")
	}
	return nil
}
