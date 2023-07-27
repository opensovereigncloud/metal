/*
Copyright (c) 2022 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved

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

package stateproc

import (
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/onmetal/metal-api/pkg/constants"
	"github.com/onmetal/metal-api/pkg/errors"
)

type GenericStateProcessor[T client.Object] struct {
	Client    client.Client
	Functions []func(T) StateFuncResult
	Log       logr.Logger
}

func NewGenericStateProcessor[T client.Object](
	cl client.Client,
	log logr.Logger,
) *GenericStateProcessor[T] {
	return &GenericStateProcessor[T]{
		Client:    cl,
		Functions: make([]func(T) StateFuncResult, 0),
		Log:       log,
	}
}

func (gsp *GenericStateProcessor[T]) RegisterHandler(handler func(T) StateFuncResult) *GenericStateProcessor[T] {
	gsp.Functions = append(gsp.Functions, handler)
	return gsp
}

func (gsp *GenericStateProcessor[T]) Compute(obj T) errors.StateProcError {
	for _, f := range gsp.Functions {
		res := f(obj)
		err := res.Err
		if err != nil {
			gsp.log(err, obj.GetName())
		}
		if res.Break {
			return err
		}
	}
	return nil
}

func (gsp *GenericStateProcessor[T]) log(err errors.StateProcError, name string) {
	switch err.Error() {
	case constants.EmptyString:
		gsp.Log.Info(
			"object processing was interrupted",
			"name", name,
			"reason", err.Reason(),
			"verbose", err.Message(),
		)
	default:
		gsp.Log.Error(
			err,
			"object processing was interrupted",
			"name", name,
			"reason", err.Reason(),
			"verbose", err.Message(),
		)
	}
}

type StateFuncResult struct {
	Err   errors.StateProcError
	Break bool
}
