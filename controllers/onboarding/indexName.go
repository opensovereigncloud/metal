// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	ipam "github.com/onmetal/ipam/api/v1alpha1"
	oob "github.com/onmetal/oob-operator/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	metalv1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"
)

func machineIndex(rawObj client.Object) []string {
	obj, _ := rawObj.(*metalv1alpha4.Machine)
	return []string{obj.ObjectMeta.Name}
}

func inventoryIndex(rawObj client.Object) []string {
	obj, _ := rawObj.(*metalv1alpha4.Inventory)
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

func subnetIndex(rawObj client.Object) []string {
	obj, _ := rawObj.(*ipam.Subnet)
	return []string{obj.ObjectMeta.Name}
}
