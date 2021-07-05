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

package v1alpha1

import (
	"math/big"
	"net"
	"time"

	subnetv1alpha1 "github.com/onmetal/ipam/api/v1alpha1"
)

const (
	CLabelPrefix    = "switch.onmetal.de/"
	CLabelChassisId = "chassisId"

	CLeafRole  = "Leaf"
	CSpineRole = "Spine"

	CMachineType = "Machine"
	CSwitchType  = "Switch"

	CSonicSwitchOs = "SONiC"

	CStationCapability = "Station"

	CAssignmentRequeueInterval = time.Duration(5) * time.Second
	CSwitchRequeueInterval     = time.Duration(15) * time.Second

	CIPv4AddressesPerPort = 16
	CIPv6AddressesPerPort = 8

	CIPv4ZeroNet = "0.0.0.0/0"
	CIPv6ZeroNet = "::/0"

	CNamespace = "onmetal"

	CSwitchFinalizer = "switches.switch.onmetal.de/finalizer"

	CIPv4InterfaceSubnetMask = 30
	CIPv6InterfaceSubnetMask = 127
)

var LabelChassisId = CLabelPrefix + CLabelChassisId

func GetInterfaceSubnetMaskLength(addrType subnetv1alpha1.SubnetAddressType) int {
	if addrType == subnetv1alpha1.CIPv4SubnetType {
		return CIPv4InterfaceSubnetMask
	} else {
		return CIPv6InterfaceSubnetMask
	}
}

//GetMinimalVacantCIDR calculates the minimal suitable network
//from the networks list provided as argument according to the
//needed addresses count. It returns the pointer to the CIDR object.
func GetMinimalVacantCIDR(vacant []subnetv1alpha1.CIDR, addressType subnetv1alpha1.SubnetAddressType, addressesCount int64) *subnetv1alpha1.CIDR {
	zeroNetString := ""
	if addressType == subnetv1alpha1.CIPv4SubnetType {
		zeroNetString = CIPv4ZeroNet
	} else {
		zeroNetString = CIPv6ZeroNet
	}
	_, zeroNet, _ := net.ParseCIDR(zeroNetString)
	minSuitableNet := subnetv1alpha1.CIDRFromNet(zeroNet)
	for _, cidr := range vacant {
		if cidr.AddressCapacity().Cmp(minSuitableNet.AddressCapacity()) < 0 &&
			cidr.AddressCapacity().Cmp(new(big.Int).SetInt64(addressesCount)) >= 0 {
			minSuitableNet = &cidr
		}
	}
	return minSuitableNet
}
