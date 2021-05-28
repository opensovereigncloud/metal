package util

import (
	"math"
	"time"
)

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

func GetNeededMaskLength(addressesCount float64) uint8 {
	pow := 2.0
	for math.Pow(2, pow) < addressesCount {
		pow++
	}
	return 32 - uint8(pow)
}
