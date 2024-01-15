/*
Copyright (c) 2023 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved

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

package switches

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"

	metalv1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"
	"github.com/ironcore-dev/metal/pkg/auxiliary"
	"github.com/ironcore-dev/metal/pkg/constants"
)

type StateUpdateFunction func(*metalv1alpha4.NetworkSwitch, *SwitchEnvironment) *StateUpdateResult

type StateUpdateEvent struct {
	eventType string
	reason    string
	message   string
}

func NewStateUpdateEvent(eventType, reason, message string) *StateUpdateEvent {
	return &StateUpdateEvent{eventType, reason, message}
}

func (in *StateUpdateEvent) Unwrap() (string, string, string) {
	return in.eventType, in.reason, in.message
}

type StateUpdateResult struct {
	err            error
	condition      string
	reason         string
	verboseMessage string
}

func NewStateUpdateResult() *StateUpdateResult {
	return &StateUpdateResult{}
}

func (in *StateUpdateResult) SetError(value error) *StateUpdateResult {
	in.err = value
	return in
}

func (in *StateUpdateResult) SetCondition(value string) *StateUpdateResult {
	in.condition = value
	return in
}

func (in *StateUpdateResult) SetReason(value string) *StateUpdateResult {
	in.reason = value
	return in
}

func (in *StateUpdateResult) SetMessage(value string) *StateUpdateResult {
	in.verboseMessage = value
	return in
}

func (in *StateUpdateResult) ToEvent() auxiliary.Event {
	eventType := v1.EventTypeWarning
	if in.err == nil {
		eventType = v1.EventTypeNormal
	}
	return NewStateUpdateEvent(eventType, in.reason, in.verboseMessage)
}

func (in *StateUpdateResult) Unwrap() (string, string, string, error) {
	return in.condition, in.reason, in.verboseMessage, in.err
}

type SwitchStateWriter struct {
	rec   record.EventRecorder
	env   *SwitchEnvironment
	funcs []StateUpdateFunction
}

func NewSwitchStateWriter(
	rec record.EventRecorder,
) *SwitchStateWriter {
	return &SwitchStateWriter{
		rec:   rec,
		env:   nil,
		funcs: nil,
	}
}

func (in *SwitchStateWriter) SetEnvironment(env *SwitchEnvironment) *SwitchStateWriter {
	in.env = env
	return in
}

func (in *SwitchStateWriter) Environment() *SwitchEnvironment {
	return in.env
}

func (in *SwitchStateWriter) RegisterStateFunc(f StateUpdateFunction) *SwitchStateWriter {
	in.funcs = append(in.funcs, f)
	return in
}

func (in *SwitchStateWriter) WriteState(o client.Object) auxiliary.Result {
	result := NewStateUpdateResult()
	obj, _ := o.(*metalv1alpha4.NetworkSwitch)
	if in.env == nil {
		result.SetCondition("").SetReason("UndefinedEnvironment").
			SetMessage("reconciliation interrupted due to undefined environment")
		return result
	}
	for _, f := range in.funcs {
		result = f(obj, in.env)
		event := result.ToEvent()
		in.WriteCondition(obj, result)
		in.WriteEvent(obj, event)
		if result.err != nil {
			return result
		}
	}
	return result
}

func (in *SwitchStateWriter) WriteCondition(o auxiliary.ObjectWithConditions, res auxiliary.Result) {
	condition, reason, message, err := res.Unwrap()
	if condition != constants.ConditionReady {
		o.UpdateCondition(condition, reason, message, err == nil)
	}
}

func (in *SwitchStateWriter) WriteEvent(src runtime.Object, e auxiliary.Event) {
	eventType, reason, message := e.Unwrap()
	in.rec.Event(src, eventType, reason, message)
}
