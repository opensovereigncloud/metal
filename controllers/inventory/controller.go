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

package inventory

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

const (
	CUndefinedRole = "Undefined"
	CLeafRole      = "Leaf"
	CSpineRole     = "Spine"
)

const (
	CSwitchType    = "Switch"
	CSonicSwitchOs = "SONiC"
)

// Reconciler reconciles a Switch object
type Reconciler struct {
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
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("inventory", req.NamespacedName)
	r.Log.Info("starting reconciliation")

	inventory := &inventoriesv1alpha1.Inventory{}
	if err := r.Get(ctx, req.NamespacedName, inventory); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}
	if inventory.Spec.Host.Type != CSwitchType {
		return ctrl.Result{}, nil
	}
	switches := &switchv1alpha1.SwitchList{}
	switchRes, exists := switchResourceExists(strings.ToLower(inventory.Spec.System.SerialNumber), switches)
	if !exists {
		r.Log.Info("switch resource does not exists")
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
	r.Log.Info("reconciliation finished")

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
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
	sw.Namespace = inv.Namespace
	sw.Spec.Hostname = inv.Spec.Host.Name
	sw.Spec.Ports = inv.Spec.NICs.Count
	sw.Spec.SwitchPorts = countSwitchPorts(inv.Spec.NICs.NICs)
	sw.Spec.SwitchDistro = &switchv1alpha1.SwitchDistroSpec{
		OS:      CSonicSwitchOs,
		Version: inv.Spec.Distro.CommitId,
		ASIC:    inv.Spec.Distro.AsicType,
	}
	sw.Spec.SwitchChassis = &switchv1alpha1.SwitchChassisSpec{
		Manufacturer: inv.Spec.System.Manufacturer,
		SKU:          inv.Spec.System.ProductSKU,
		Serial:       inv.Spec.System.SerialNumber,
		ChassisID:    getChassisId(inv.Spec.NICs.NICs),
	}
	sw.Spec.Interfaces, sw.Spec.Role = setInterfaces(inv.Spec.NICs.NICs)
	return sw, nil
}

func setInterfaces(nics []inventoriesv1alpha1.NICSpec) ([]*switchv1alpha1.InterfaceSpec, string) {
	role := CUndefinedRole
	interfaces := make([]*switchv1alpha1.InterfaceSpec, 0)
	for _, nic := range nics {
		iface, neighbourExists, machinesConnected := buildInterface(&nic)
		if neighbourExists {
			if !machinesConnected && role == CUndefinedRole {
				role = CSpineRole
			}
			if machinesConnected && (role == CSpineRole || role == CUndefinedRole) {
				role = CLeafRole
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
	}
	if len(nic.LLDPs) != 0 {
		neighbourExists = true
		lldpData := nic.LLDPs[0]
		iface.LLDPChassisID = lldpData.ChassisID
		iface.LLDPSystemName = lldpData.SystemName
		iface.LLDPPortID = lldpData.PortID
		iface.LLDPPortDescription = lldpData.PortDescription
		//todo: check neighbour type using advertised LLDP capabilities.
		//  If station capability advertised - change "machineConnected" to true
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

func getChassisId(nics []inventoriesv1alpha1.NICSpec) string {
	for _, nic := range nics {
		if nic.Name == "eth0" {
			return nic.MACAddress
		}
	}
	return ""
}
