/*
Copyright 2021.

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

package v1alpha1

import (
	"strings"

	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	"github.com/onmetal/switch-operator/util"
)

// log is for logging in this package.
var switchAssignmentLog = logf.Log.WithName("switchassignment-resource")

func (swa *SwitchAssignment) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(swa).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-switch-onmetal-de-onmetal-de-v1alpha1-switchassignment,mutating=true,failurePolicy=fail,sideEffects=None,groups=switch.onmetal.de,resources=switchassignments,verbs=create;update,versions=v1alpha1,name=mswitchassignment.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &SwitchAssignment{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (swa *SwitchAssignment) Default() {
	switchAssignmentLog.Info("default", "name", swa.Name)
	swa.Labels = map[string]string{}
	swa.Labels[util.LabelSerial] = swa.Spec.Serial
	swa.Labels[util.LabelChassisId] = strings.ReplaceAll(swa.Spec.ChassisID, ":", "-")
}
