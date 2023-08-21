package controllers

import (
	ipam "github.com/onmetal/ipam/api/v1alpha1"
	inventories "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	machine "github.com/onmetal/metal-api/apis/machine/v1alpha3"
	oob "github.com/onmetal/oob-operator/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func machineIndex(rawObj client.Object) []string {
	obj, _ := rawObj.(*machine.Machine)
	return []string{obj.ObjectMeta.Name}
}

func inventoryIndex(rawObj client.Object) []string {
	obj, _ := rawObj.(*inventories.Inventory)
	return []string{obj.ObjectMeta.Name}
}

func oobIndex(rawObj client.Object) []string {
	obj, _ := rawObj.(*oob.OOB)
	return []string{obj.ObjectMeta.Name}
}

func ipIndex(rawObj client.Object) []string {
	obj, _ := rawObj.(*ipam.IP)
	return []string{obj.ObjectMeta.Name}
}
