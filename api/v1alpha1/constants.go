package v1alpha1

import "time"

const (
	CLabelPrefix        = "switch.onmetal.de/"
	CLabelSerial        = "serial"
	CLabelChassisId     = "chassisId"
	CLabelConnChassisId = "connection-chassisId"
)

const (
	CUndefinedRole = "Undefined"
	CLeafRole      = "Leaf"
	CSpineRole     = "Spine"
)

const (
	CMachineType   = "Machine"
	CSwitchType    = "Switch"
	CSonicSwitchOs = "SONiC"
)

const CStationCapability = "Station"

const CRequeueInterval = time.Duration(5) * time.Second

const (
	CIPv4AddressesPerPort = 16
	CIPv6AddressesPerPort = 8
)

const Namespace = "onmetal"

var LabelSerial = CLabelPrefix + CLabelSerial
var LabelChassisId = CLabelPrefix + CLabelChassisId
var ConnectionLabelChassisId = CLabelPrefix + CLabelConnChassisId
