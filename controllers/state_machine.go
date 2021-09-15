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
	"reflect"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"

	switchv1alpha1 "github.com/onmetal/switch-operator/api/v1alpha1"
)

const CSwitchRequeueInterval = time.Second

type executable interface {
	execute(*switchv1alpha1.Switch, *switchv1alpha1.Switch) error
	setNext(executable)
	getNext() executable
}

type processingStep struct {
	checkFunc  func(*switchv1alpha1.Switch) bool
	setFunc    func(*switchv1alpha1.Switch) error
	updateFunc func(*switchv1alpha1.Switch) error
	next       executable
}

func newStep(
	stepCheckFunc func(*switchv1alpha1.Switch) bool,
	stepSetFunc func(*switchv1alpha1.Switch) error,
	stepUpdateFunc func(*switchv1alpha1.Switch) error,
	nextStep executable) *processingStep {
	return &processingStep{
		checkFunc:  stepCheckFunc,
		setFunc:    stepSetFunc,
		updateFunc: stepUpdateFunc,
		next:       nextStep,
	}
}

func (s *processingStep) setNext(step executable) {
	s.next = step
}

func (s *processingStep) getNext() executable {
	return s.next
}

func (s *processingStep) execute(obj *switchv1alpha1.Switch, initialState *switchv1alpha1.Switch) error {
	if s.checkFunc != nil {
		if !s.checkFunc(obj) {
			return s.setAndUpdate(obj, initialState)
		}
	} else {
		return s.setAndUpdate(obj, initialState)
	}
	return nil
}

func (s *processingStep) setAndUpdate(obj *switchv1alpha1.Switch, initialState *switchv1alpha1.Switch) error {
	if err := s.setFunc(obj); err != nil {
		return err
	}
	if (!reflect.DeepEqual(initialState.Spec, obj.Spec) || !obj.StatusDeepEqual(*initialState)) && s.updateFunc != nil {
		s.setNext(nil)
		if err := s.updateFunc(obj); err != nil {
			return err
		}
		return nil
	}
	if obj.Status.Configuration != nil && initialState.Status.Configuration != nil {
		if (obj.Status.Configuration.Type != initialState.Status.Configuration.Type && obj.Status.Configuration.Type == switchv1alpha1.CConfigManagementTypeFailed) &&
			s.updateFunc != nil {
			s.setNext(nil)
			if err := s.updateFunc(obj); err != nil {
				return err
			}
			return nil
		}
	}
	return nil
}

type stateMachine struct {
	executionPoint executable
}

func newStateMachine(startPoint executable) *stateMachine {
	return &stateMachine{
		executionPoint: startPoint,
	}
}

func (s *stateMachine) launch(obj *switchv1alpha1.Switch) (ctrl.Result, error) {
	initialState := obj.DeepCopy()
	if err := s.executeStep(obj, initialState); err != nil {
		return ctrl.Result{RequeueAfter: CSwitchRequeueInterval}, err
	}
	if reflect.DeepEqual(initialState.Spec, obj.Spec) && obj.StatusDeepEqual(*initialState) {
		return ctrl.Result{RequeueAfter: CSwitchRequeueInterval}, nil
	}
	return ctrl.Result{}, nil
}

func (s *stateMachine) executeStep(currentState *switchv1alpha1.Switch, initialState *switchv1alpha1.Switch) error {
	if err := s.executionPoint.execute(currentState, initialState); err != nil {
		return err
	}
	s.executionPoint = s.executionPoint.getNext()
	if s.executionPoint != nil {
		return s.executeStep(currentState, initialState)
	}
	return nil
}
