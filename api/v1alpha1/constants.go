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

package v1alpha1

import "time"

const (
	CLabelPrefix    = "switch.onmetal.de/"
	CLabelSerial    = "serial"
	CLabelChassisId = "chassisId"

	CUndefinedRole = "Undefined"
	CLeafRole      = "Leaf"
	CSpineRole     = "Spine"

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
)

var LabelSerial = CLabelPrefix + CLabelSerial
var LabelChassisId = CLabelPrefix + CLabelChassisId
