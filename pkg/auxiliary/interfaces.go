// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

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
