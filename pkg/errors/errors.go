// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

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
	resourceNotExist StatusReason = "resource not exist"
	uuidNotExist     StatusReason = "uuid not exist"
	notFound         StatusReason = "not found"
	notLabeled       StatusReason = "not labeled"
	notAMachine      StatusReason = "not a machine"
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

func IsUUIDNotExist(err error) bool { return ReasonForError(err) == uuidNotExist }

func IsNotAMachine(err error) bool { return ReasonForError(err) == notAMachine }

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

func UUIDNotExist(s string) *MachineError {
	return &MachineError{
		ErrStatus: Reason{
			Message: fmt.Sprintf("object: %s, exist but doesn't contain uuid", s), StatusReason: uuidNotExist,
		},
	}
}

func NotFound(msg string) *MachineError {
	return &MachineError{
		ErrStatus: Reason{Message: fmt.Sprintf("%s, not found", msg), StatusReason: notFound},
	}
}

func NotSizeLabeled() *MachineError {
	return &MachineError{
		ErrStatus: Reason{Message: "size labels are not present yet", StatusReason: notLabeled},
	}
}

func Unknown() *MachineError {
	return &MachineError{
		ErrStatus: Reason{Message: "unknown error happened", StatusReason: unknown},
	}
}
