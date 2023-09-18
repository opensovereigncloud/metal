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

package auxiliary

import (
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ObjectWithConditions interface {
	UpdateCondition(name, reason, message string, state bool)
}

type Event interface {
	Unwrap() (eventType, reason, message string)
}

type Result interface {
	Unwrap() (condition, reason, message string, err error)
	ToEvent() Event
}

type StateWriter interface {
	WriteState(client.Object) Result
}

type ConditionWriter interface {
	WriteCondition(ObjectWithConditions, Result)
}

type EventWriter interface {
	WriteEvent(runtime.Object, Event)
}
