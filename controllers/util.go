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

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	switchv1alpha1 "github.com/onmetal/switch-operator/api/v1alpha1"
)

type switchProcessor struct {
	startPoint processorStep
}

type processorStep interface {
	setNext(processorStep)
	getNext() processorStep
	execute(*switchv1alpha1.Switch, *SwitchReconciler, context.Context) error
}

type step struct {
	next processorStep
}

type preparationStep step
type interfacesStep step
type creationStep step
type assignmentStep step
type connectionsStep step
type subnetsStep step
type addressesStep step
type deletionStep step
type specUpdateState step
type statusUpdateStep step

func (p *deletionStep) setNext(processorStep) {}

func (p *deletionStep) getNext() processorStep {
	return p.next
}

func (p *deletionStep) execute(obj *switchv1alpha1.Switch, r *SwitchReconciler, ctx context.Context) error {
	if err := r.finalize(ctx, obj); err != nil {
		r.Log.Error(err, "failed to finalize resource", "name", obj.NamespacedName())
		return err
	}
	return nil
}

func (p *specUpdateState) setNext(processorStep) {}

func (p *specUpdateState) getNext() processorStep {
	return p.next
}

func (p *specUpdateState) execute(obj *switchv1alpha1.Switch, r *SwitchReconciler, ctx context.Context) error {
	if err := r.Update(ctx, obj); err != nil {
		r.Log.Error(err, "failed to update resource", "name", obj.NamespacedName())
		return err
	}
	return nil
}

func (p *statusUpdateStep) setNext(processorStep) {}

func (p *statusUpdateStep) getNext() processorStep {
	return p.next
}

func (p *statusUpdateStep) execute(obj *switchv1alpha1.Switch, r *SwitchReconciler, ctx context.Context) error {
	if err := r.Status().Update(ctx, obj); err != nil {
		r.Log.Error(err, "failed to update resource status", "name", obj.NamespacedName())
		return err
	}
	return nil
}

func (p *preparationStep) setNext(next processorStep) {
	p.next = next
}

func (p *preparationStep) getNext() processorStep {
	return p.next
}

func (p *preparationStep) execute(obj *switchv1alpha1.Switch, r *SwitchReconciler, ctx context.Context) error {
	if !controllerutil.ContainsFinalizer(obj, switchv1alpha1.CSwitchFinalizer) {
		controllerutil.AddFinalizer(obj, switchv1alpha1.CSwitchFinalizer)
		p.setNext(&specUpdateState{})
		return nil
	}
	p.setNext(&interfacesStep{})
	return nil
}

func (p *interfacesStep) setNext(next processorStep) {
	p.next = next
}

func (p *interfacesStep) getNext() processorStep {
	return p.next
}

func (p *interfacesStep) execute(obj *switchv1alpha1.Switch, r *SwitchReconciler, ctx context.Context) error {
	if !obj.InterfacesUpdated(r.Background.switches) {
		obj.UpdateInterfaces(r.Background.switches)
		p.setNext(&specUpdateState{})
		return nil
	}
	p.setNext(&creationStep{})
	return nil
}

func (p *creationStep) setNext(next processorStep) {
	p.next = next
}

func (p *creationStep) getNext() processorStep {
	return p.next
}

func (p *creationStep) execute(obj *switchv1alpha1.Switch, r *SwitchReconciler, ctx context.Context) error {
	if obj.Status.State == switchv1alpha1.EmptyString {
		obj.FillStatusOnCreate()
		p.setNext(&statusUpdateStep{})
		return nil
	}
	p.setNext(&assignmentStep{})
	return nil
}

func (p *assignmentStep) setNext(next processorStep) {
	p.next = next
}

func (p *assignmentStep) getNext() processorStep {
	return p.next
}

func (p *assignmentStep) execute(obj *switchv1alpha1.Switch, r *SwitchReconciler, ctx context.Context) error {
	if r.Background.assignment != nil {
		if r.Background.assignment.Status.State != switchv1alpha1.StateFinished {
			r.Background.assignment.FillStatus(switchv1alpha1.StateFinished, &switchv1alpha1.LinkedSwitchSpec{
				Name:      obj.Name,
				Namespace: obj.Namespace,
			})
			if err := r.Status().Update(ctx, r.Background.assignment); err != nil {
				r.Log.Error(err, "failed to update resource", "kind", r.Background.assignment.Kind, "name", r.Background.assignment.NamespacedName())
				return err
			}
			obj.Status.ConnectionLevel = 0
		}
	}
	p.setNext(&connectionsStep{})
	return nil
}

func (p *connectionsStep) setNext(next processorStep) {
	p.next = next
}

func (p *connectionsStep) getNext() processorStep {
	return p.next
}

func (p *connectionsStep) execute(obj *switchv1alpha1.Switch, r *SwitchReconciler, ctx context.Context) error {
	ok := obj.PeersProcessingFinished(r.Background.switches, r.Background.assignment)
	if ok {
		obj.Status.State = switchv1alpha1.StateDefineAddresses
		if r.Background.switches.AllConnectionsOk() {
			p.setNext(&subnetsStep{})
			return nil
		}
		p.setNext(&statusUpdateStep{})
	} else {
		obj.Status.State = switchv1alpha1.StateDefinePeers
		obj.UpdatePeersData(r.Background.switches)
		obj.UpdateConnectionLevel(r.Background.switches)
		p.setNext(&statusUpdateStep{})
		return nil
	}
	return nil
}

func (p *subnetsStep) setNext(next processorStep) {
	p.next = next
}

func (p *subnetsStep) getNext() processorStep {
	return p.next
}

func (p *subnetsStep) execute(obj *switchv1alpha1.Switch, r *SwitchReconciler, ctx context.Context) error {
	if obj.AddressesDefined() {
		if obj.Status.State != switchv1alpha1.StateFinished {
			obj.Status.State = switchv1alpha1.StateFinished
			p.setNext(&statusUpdateStep{})
			return nil
		}
	} else {
		if err := r.defineSubnets(ctx, obj, r.Background.switches, r.Background.assignment); err != nil {
			r.Log.Error(err, "failed to define south subnets")
			return err
		}
		p.setNext(&addressesStep{})
		return nil
	}
	return nil
}

func (p *addressesStep) setNext(next processorStep) {
	p.next = next
}

func (p *addressesStep) getNext() processorStep {
	return p.next
}

func (p *addressesStep) execute(obj *switchv1alpha1.Switch, r *SwitchReconciler, ctx context.Context) error {
	obj.UpdateSouthInterfacesAddresses()
	obj.UpdateNorthInterfacesAddresses(r.Background.switches)
	p.setNext(&specUpdateState{})
	return nil
}

func (c *switchProcessor) launch(obj *switchv1alpha1.Switch, r *SwitchReconciler, ctx context.Context) (ctrl.Result, error) {
	return executeStep(c.startPoint, obj, r, ctx)
}

func executeStep(step processorStep, obj *switchv1alpha1.Switch, r *SwitchReconciler, ctx context.Context) (ctrl.Result, error) {
	if err := step.execute(obj, r, ctx); err != nil {
		return ctrl.Result{}, err
	}
	next := step.getNext()
	if next != nil {
		return executeStep(next, obj, r, ctx)
	}
	return ctrl.Result{}, nil
}
