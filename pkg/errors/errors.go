/*
Copyright (c) 2021 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved

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

package errors

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const waitForMachineMinute = 5 * time.Minute //nolint:revive //reason: false-positive

type Reason struct {
	Message      string
	StatusReason StatusReason
}

type StatusReason string

const (
	alreadyExist     StatusReason = "already exist"
	resourceNotExist StatusReason = "resource not exist"
	uuidNotExist     StatusReason = "uuid not exist"
	notFound         StatusReason = "not found"
	notLabeled       StatusReason = "not labeled"
	notAMachine      StatusReason = "not a machine"
	castType         StatusReason = "casting failed"
	reservation      StatusReason = "machine reservation impossible"
	unknown          StatusReason = "unknown"
)

type MachineError struct {
	ErrStatus Reason
}

type MachineStatus interface {
	Status() Reason
}

func (e *MachineError) Error() string { return e.ErrStatus.Message }

func (e *MachineError) Status() Reason { return e.ErrStatus }

func IsNotExist(err error) bool { return ReasonForError(err) == resourceNotExist }

func IsNotFound(err error) bool { return ReasonForError(err) == notFound }

func IsAlreadyExists(err error) bool { return ReasonForError(err) == alreadyExist }

func IsUUIDNotExist(err error) bool { return ReasonForError(err) == uuidNotExist }

func IsNotAMachine(err error) bool { return ReasonForError(err) == notAMachine }

func IsNotBookable(err error) bool { return ReasonForError(err) == reservation }

func IsNotSizeLabeled(err error) bool { return ReasonForError(err) == notLabeled }

func ReasonForError(err error) StatusReason {
	if reason := MachineStatus(nil); errors.As(err, &reason) {
		return reason.Status().StatusReason
	}
	return unknown
}

func GetResultForError(reqLogger logr.Logger, err error) (ctrl.Result, error) {
	switch {
	case IsUUIDNotExist(err):
		reqLogger.Info("uuid not provided")
		return ctrl.Result{}, nil
	case IsNotExist(err):
		reqLogger.Info("resource not found", "error", err)
		return ctrl.Result{RequeueAfter: waitForMachineMinute}, nil
	case apierrors.IsConflict(err):
		return ctrl.Result{Requeue: true}, nil
	case IsNotAMachine(err):
		return ctrl.Result{}, nil
	case IsNotFound(err):
		reqLogger.Info("can't operate over object", "error", err)
		return ctrl.Result{}, nil
	case err == nil:
		reqLogger.Info("reconciliation finished")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	default:
		reqLogger.Info("can't operate over object", "error", err)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
}

func AlreadyExist(s string) *MachineError {
	return &MachineError{
		ErrStatus: Reason{
			Message: fmt.Sprintf("resource with uuid: %s, already exist ", s), StatusReason: alreadyExist,
		},
	}
}

func UUIDNotExist(s string) *MachineError {
	return &MachineError{
		ErrStatus: Reason{
			Message: fmt.Sprintf("object: %s, exist but doesn't contain uuid", s), StatusReason: uuidNotExist,
		},
	}
}

func NotExist(s string) *MachineError {
	return &MachineError{
		ErrStatus: Reason{Message: fmt.Sprintf("resource with uuid: %s, not exist", s), StatusReason: resourceNotExist},
	}
}

func NotFound(msg string) *MachineError {
	return &MachineError{
		ErrStatus: Reason{Message: fmt.Sprintf("%s, not found", msg), StatusReason: notFound},
	}
}

func ObjectNotFound(obj, kind string) *MachineError {
	return &MachineError{
		ErrStatus: Reason{Message: fmt.Sprintf("object name: %s, kind: %s not found", obj, kind), StatusReason: notFound},
	}
}

func NotSizeLabeled() *MachineError {
	return &MachineError{
		ErrStatus: Reason{Message: "size labels are not present yet", StatusReason: notLabeled},
	}
}

func NotAMachine() *MachineError {
	return &MachineError{
		ErrStatus: Reason{Message: "not a machine", StatusReason: notAMachine},
	}
}

func CastType() *MachineError {
	return &MachineError{
		ErrStatus: Reason{Message: "object casting failed", StatusReason: castType},
	}
}

func NotBookable() *MachineError {
	return &MachineError{
		ErrStatus: Reason{Message: "impossible to reserve", StatusReason: reservation},
	}
}

func Unknown() *MachineError {
	return &MachineError{
		ErrStatus: Reason{Message: "unknown error happened", StatusReason: unknown},
	}
}
