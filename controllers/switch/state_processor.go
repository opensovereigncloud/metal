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

package v1beta1

import (
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
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

func (gsp *GenericStateProcessor[T]) setFunctions(list []func(T) StateFuncResult) {
	gsp.Functions = append(gsp.Functions, list...)
}

func (gsp *GenericStateProcessor[T]) compute(obj T) error {
	var err error
	for _, f := range gsp.Functions {
		res := f(obj)
		err = res.Err
		if err != nil {
			gsp.log(err, obj.GetName())
		}
		if res.Break {
			break
		}
	}
	return err
}

func (gsp *GenericStateProcessor[T]) log(err error, name string) {
	gsp.Log.V(1).Info(
		"object processing was interrupted due to an error",
		"name", name,
		"error", err.Error(),
	)
}

type StateFuncResult struct {
	Err   error
	Break bool
}
