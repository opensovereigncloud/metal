package errors

import (
	"fmt"

	"github.com/onmetal/metal-api/internal/constants"
)

const (
	StateMessageRequestRelatedObjectsFailed string = "failed to request related objects, check conditions for details"
	StateMessageMissingRequirements         string = "some requirements are missing, check conditions for details"
)

const (
	ErrorReasonMissingRequirements  StateProcErrorReason = "MissingRequirements"
	ErrorReasonMissingInventoryRef  StateProcErrorReason = "MissingInventoryRef"
	ErrorReasonMissingTypeLabel     StateProcErrorReason = "MissingTypeLabel"
	ErrorReasonRequestFailed        StateProcErrorReason = "APIRequestFailed"
	ErrorReasonObjectNotExist       StateProcErrorReason = "ObjectNotExist"
	ErrorReasonASNCalculationFailed StateProcErrorReason = "ASNCalculationFailed"
	ErrorReasonIPAssignmentFailed   StateProcErrorReason = "IPAssignmentFailed"
	ErrorReasonDataOutdated         StateProcErrorReason = "DataOutdated"
)

const (
	MessageMissingTypeLabel          string = "missing requirements: label switch.onmetal.de/type"
	MessageMissingInventoryRef       string = "missing requirements: reference to corresponding Inventory at .spec.InventoryRef.name"
	MessageMissingLoopbackV4IP       string = "missing requirements: IP object of V4 address family to be assigned to loopback interface"
	MessageRequestFailed             string = "failed to get requested object"
	MessageObjectNotExist            string = "requested object does not exist"
	MessageFailedToAssignIPAddresses string = "failed to assign IP addresses to switch ports"
	MessageParseIPFailed             string = "failed to parse IP address"
	MessageParseCIDRFailed           string = "failed to parse CIDR"
	MessageInvalidInputType          string = "invalid input type"
	MessageMissingAPIVersion         string = "missing API version"
	MessageAPIVersionMismatch        string = "API version mismatch"
	MessageDuplicatedIPAddress       string = "duplicated IP address"
	MessageDuplicatedSubnet          string = "duplicated subnet"
	MessageFieldSelectorNotDefined   string = "field selector is not defined"
	MessageUnmarshallingFailed       string = "failed to unmarshal bytes to map"
	MessageMarshallingFailed         string = "failed to marshal input to bytes"
	MessageInvalidFieldPath          string = "invalid field path"
)

type SwitchError struct {
	reason  StateProcErrorReason
	message string
	error   error
}

func NewSwitchError(reason StateProcErrorReason, message string, err error) *SwitchError {
	return &SwitchError{reason, message, err}
}

func (s *SwitchError) Reason() StateProcErrorReason {
	return s.reason
}

func (s *SwitchError) Message() string {
	return s.message
}

func (s *SwitchError) Error() string {
	if s.error == nil {
		return constants.EmptyString
	}
	return s.error.Error()
}

func IsMissingRequirements(err StateProcError) bool {
	return err != nil && err.Reason() == ErrorReasonMissingRequirements
}

func MessageRequestFailedWithKind(kind string) string {
	return fmt.Sprintf("%s: %s", MessageRequestFailed, kind)
}

func MessageObjectNotExistWithKind(kind string) string {
	return fmt.Sprintf("%s: %s", MessageObjectNotExist, kind)
}
