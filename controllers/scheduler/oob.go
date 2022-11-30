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

package scheduler

import (
	"context"

	"github.com/onmetal/metal-api/apis/machine/v1alpha2"
	oobv1 "github.com/onmetal/oob-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *Reconciler) getOOBMachine(ctx context.Context, machine *v1alpha2.Machine) (*oobv1.OOB, error) {
	oobMachine := &oobv1.OOB{
		ObjectMeta: metav1.ObjectMeta{
			Name:      machine.Status.OOB.Reference.Name,
			Namespace: machine.Status.OOB.Reference.Namespace,
		},
	}

	err := r.Client.Get(ctx, client.ObjectKeyFromObject(oobMachine), oobMachine)
	if err != nil {
		return nil, err
	}

	return oobMachine, nil
}
