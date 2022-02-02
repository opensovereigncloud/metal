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

package controllers

import (
	"context"
	"reflect"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"

	switchv1alpha1 "github.com/onmetal/metal-api/apis/switches/v1alpha1"
)

const CSwitchRequeueInterval = time.Second

type executable interface {
	execute(context.Context, *switchv1alpha1.Switch, *switchv1alpha1.Switch) error
	setNext(executable)
	getNext() executable
}

type processingStep struct {
	stateCheckFunc     func(*switchv1alpha1.Switch) bool
	stateComputeFunc   func(context.Context, *switchv1alpha1.Switch) error
	resourceUpdateFunc func(context.Context, *switchv1alpha1.Switch) error
	next               executable
}

func newStep(
	stepCheckFunc func(*switchv1alpha1.Switch) bool,
	stepComputeFunc func(context.Context, *switchv1alpha1.Switch) error,
	stepResUpdateFunc func(context.Context, *switchv1alpha1.Switch) error,
	nextStep executable) *processingStep {
	return &processingStep{
		stateCheckFunc:     stepCheckFunc,
		stateComputeFunc:   stepComputeFunc,
		resourceUpdateFunc: stepResUpdateFunc,
		next:               nextStep,
	}
}

func (s *processingStep) setNext(step executable) {
	s.next = step
}

func (s *processingStep) getNext() executable {
	return s.next
}

func (s *processingStep) execute(ctx context.Context, newState *switchv1alpha1.Switch, oldState *switchv1alpha1.Switch) (err error) {
	if s.stateCheckFunc(newState) {
		return
	}
	if err = s.stateComputeFunc(ctx, newState); err != nil {
		return
	}
	if s.resourceUpdateFunc != nil {
		s.setNext(nil)
		if !reflect.DeepEqual(newState.Status, oldState.Status) {
			err = s.resourceUpdateFunc(ctx, newState)
		}
	}
	return
}

type stateMachine struct {
	executionPoint executable
}

func newStateMachine(startPoint executable) *stateMachine {
	return &stateMachine{
		executionPoint: startPoint,
	}
}

func (s *stateMachine) launch(ctx context.Context, curState *switchv1alpha1.Switch) (ctrl.Result, error) {
	prevState := curState.DeepCopy()
	if err := s.executeStep(ctx, curState, prevState); err != nil {
		return ctrl.Result{RequeueAfter: CSwitchRequeueInterval}, err
	}
	stateChanged := !reflect.DeepEqual(curState.Status, prevState.Status)
	switch stateChanged {
	case true:
		return ctrl.Result{}, nil
	case false:
		if curState.Status.State != switchv1alpha1.CSwitchStateReady {
			return ctrl.Result{RequeueAfter: CSwitchRequeueInterval}, nil
		}
		if curState.Status.Configuration.Managed {
			return ctrl.Result{RequeueAfter: CSwitchRequeueInterval * 5}, nil
		}
	}
	return ctrl.Result{}, nil
}

func (s *stateMachine) executeStep(ctx context.Context, curState *switchv1alpha1.Switch, prevState *switchv1alpha1.Switch) (err error) {
	if err = s.executionPoint.execute(ctx, curState, prevState); err != nil {
		return
	}
	s.executionPoint = s.executionPoint.getNext()
	if s.executionPoint == nil {
		return
	}
	return s.executeStep(ctx, curState, prevState)
}
