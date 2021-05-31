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
	"strings"

	"github.com/go-logr/logr"
	inventoriesv1alpha1 "github.com/onmetal/k8s-inventory/api/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	switchv1alpha1 "github.com/onmetal/switch-operator/api/v1alpha1"
)

var Lanes = map[uint32]uint8{
	1000:   1,
	10000:  1,
	25000:  1,
	40000:  4,
	50000:  2,
	100000: 4,
}

// InventoryReconciler reconciles a Switch object
type InventoryReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=switch.onmetal.de,resources=switches,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=switch.onmetal.de,resources=switches/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=switch.onmetal.de,resources=switches/finalizers,verbs=update
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=inventories,verbs=get;list;watch
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=inventories/status,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *InventoryReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("inventory", req.NamespacedName)

	inventory := &inventoriesv1alpha1.Inventory{}
	if err := r.Get(ctx, req.NamespacedName, inventory); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}
	if inventory.Spec.Host.Type != switchv1alpha1.CSwitchType {
		return ctrl.Result{}, nil
	}
	switches := &switchv1alpha1.SwitchList{}
	if err := r.List(ctx, switches); err != nil {
		log.Error(err, "unable to get switches list")
	}
	switchRes, exists := switchResourceExists(strings.ToLower(inventory.Name), switches)
	if !exists {
		preparedSwitch, err := getPreparedSwitch(switchRes, inventory)
		if err != nil {
			r.Log.Error(err, "failed to prepare switch resource for creation")
			return ctrl.Result{}, err
		}
		if err := r.Client.Create(ctx, preparedSwitch); err != nil {
			r.Log.Error(err, "failed to create switch resource")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *InventoryReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&inventoriesv1alpha1.Inventory{}).
		Complete(r)
}

func switchResourceExists(name string, switches *switchv1alpha1.SwitchList) (*switchv1alpha1.Switch, bool) {
	for _, switchRes := range switches.Items {
		if switchRes.Name == name {
			return &switchRes, true
		}
	}
	return &switchv1alpha1.Switch{}, false
}

func getPreparedSwitch(sw *switchv1alpha1.Switch, inv *inventoriesv1alpha1.Inventory) (*switchv1alpha1.Switch, error) {
	sw.Name = inv.Name
	sw.Namespace = switchv1alpha1.Namespace
	if inv.Labels != nil {
		sw.Labels = inv.Labels
	} else {
		sw.Labels = map[string]string{}
	}
	sw.Spec.Hostname = inv.Spec.Host.Name
	sw.Spec.Ports = inv.Spec.NICs.Count
	sw.Spec.SwitchPorts = countSwitchPorts(inv.Spec.NICs.NICs)
	sw.Spec.SwitchDistro = &switchv1alpha1.SwitchDistroSpec{
		OS:      switchv1alpha1.CSonicSwitchOs,
		Version: inv.Spec.Distro.CommitId,
		ASIC:    inv.Spec.Distro.AsicType,
	}
	sw.Spec.SwitchChassis = &switchv1alpha1.SwitchChassisSpec{
		Manufacturer: inv.Spec.System.Manufacturer,
		SKU:          inv.Spec.System.ProductSKU,
		Serial:       inv.Spec.System.SerialNumber,
	}
	label := getChassisId(inv.Spec.NICs.NICs)
	if label != nil {
		sw.Spec.SwitchChassis.ChassisID = label.(string)
		sw.Labels[switchv1alpha1.LabelChassisId] = strings.ReplaceAll(label.(string), ":", "-")
	}
	sw.Spec.Interfaces, sw.Spec.Role = setInterfaces(inv.Spec.NICs.NICs)
	sw.Spec.ConnectionLevel = 255
	return sw, nil
}

func setInterfaces(nics []inventoriesv1alpha1.NICSpec) ([]*switchv1alpha1.InterfaceSpec, string) {
	role := switchv1alpha1.CUndefinedRole
	interfaces := make([]*switchv1alpha1.InterfaceSpec, 0)
	for _, nic := range nics {
		iface, neighbourExists, machinesConnected := buildInterface(&nic)
		if neighbourExists {
			if machinesConnected {
				role = switchv1alpha1.CLeafRole
			}
		}
		interfaces = append(interfaces, iface)
	}
	return interfaces, role
}

func buildInterface(nic *inventoriesv1alpha1.NICSpec) (*switchv1alpha1.InterfaceSpec, bool, bool) {
	neighbourExists := false
	machineConnected := false
	iface := &switchv1alpha1.InterfaceSpec{
		Name:       nic.Name,
		MACAddress: nic.MACAddress,
		Lanes:      Lanes[nic.Speed],
	}
	if len(nic.LLDPs) > 1 {
		neighbourExists = true
		lldpData := nic.LLDPs[1]
		iface.LLDPChassisID = lldpData.ChassisID
		iface.LLDPSystemName = lldpData.SystemName
		iface.LLDPPortID = lldpData.PortID
		iface.LLDPPortDescription = lldpData.PortDescription
		iface.Neighbour = switchv1alpha1.CSwitchType
		for i := range lldpData.Capabilities {
			if lldpData.Capabilities[i] == switchv1alpha1.CStationCapability {
				machineConnected = true
				iface.Neighbour = switchv1alpha1.CMachineType
				break
			}
		}
	}
	return iface, neighbourExists, machineConnected
}

func countSwitchPorts(nics []inventoriesv1alpha1.NICSpec) uint64 {
	count := uint64(0)
	for _, item := range nics {
		if item.PCIAddress == "" {
			count++
		}
	}
	return count
}

func getChassisId(nics []inventoriesv1alpha1.NICSpec) interface{} {
	for _, nic := range nics {
		if nic.Name == "eth0" {
			return nic.MACAddress
		}
	}
	return nil
}
