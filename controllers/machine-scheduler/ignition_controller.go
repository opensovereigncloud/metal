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

package controllers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/go-logr/logr"
	"github.com/onmetal/metal-api/apis/machine/v1alpha2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	errTemplateNameIsMissed = errors.New("template is missing required 'name' field")
)

// IgnitionReconciler reconciles a Ignition object.
type IgnitionReconciler struct {
	client.Client

	Log    logr.Logger
	Scheme *runtime.Scheme
}

type MachineWrapper struct {
	Machine           *v1alpha2.Machine           `json:"machine"`
	MachineAssignment *v1alpha2.MachineAssignment `json:"machineAssignment"`
}

const (
	ignitionFieldOwner = client.FieldOwner("metal-api.onmetal.de/ignition")

	templateName       = "ipxe-template"
	secretTemplateName = "ipxe-secret-template"
)

//+kubebuilder:rbac:groups=machine.machine.onmetal.de,resources=ignitions,verbs=get;list;watch
//+kubebuilder:rbac:groups=machine.machine.onmetal.de,resources=ignitions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
//nolint:funlen,gocognit,nestif,cyclop //TODO: linter disabled but we need to fix the problems.
func (r *IgnitionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := r.Log.WithValues("namespace", req.NamespacedName)

	reqLogger.V(1).Info("reconciling ignition", "ignition", req)

	var err error
	var err2 error

	reqLogger.V(1).Info("fetching template configmaps")
	var templateCM *corev1.ConfigMap
	if err = r.Client.Get(ctx, types.NamespacedName{
		Namespace: req.Namespace,
		Name:      templateName}, templateCM); err != nil {
		reqLogger.Error(err, "template config map is not available in the current namespace",
			"template name", templateName,
			"namespace", req.Namespace)
	}
	var secretTemplateCM *corev1.ConfigMap
	if err2 = r.Client.Get(ctx, types.NamespacedName{
		Namespace: req.Namespace,
		Name:      secretTemplateName}, secretTemplateCM); err != nil {
		reqLogger.Error(err, "template config map is not available in the current namespace",
			"secret template name", secretTemplateName,
			"namespace", req.Namespace)
	}

	if err != nil && err2 != nil {
		return ctrl.Result{}, err
	}

	reqLogger.V(1).Info("fetching machine assignment resource", "machine assignment", req)
	machineAssignment := &v1alpha2.MachineAssignment{}
	if err = r.Get(ctx, req.NamespacedName, machineAssignment); err != nil {
		reqLogger.Error(err, "couldn't get machine assignment in namespace",
			"machine assignment", req.Name, "namespace", req.Namespace)
		return ctrl.Result{}, err
	}

	if machineAssignment.Status.MachineRef.Name == "" {
		reqLogger.V(1).Info("machine is not yet reserved")
		return ctrl.Result{}, nil
	}

	machine := &v1alpha2.Machine{}
	if err = r.Get(ctx, types.NamespacedName{
		Namespace: machineAssignment.Status.MachineRef.Namespace,
		Name:      machineAssignment.Status.MachineRef.Name,
	}, machine); err != nil {
		reqLogger.Error(err, "couldn't get assigned machine in namespace", "machine", req.Name, "namespace", req.Namespace)
		return ctrl.Result{}, err
	}

	if !machineAssignment.DeletionTimestamp.IsZero() {
		if templateCM != nil {
			data, err := parseTemplate(templateCM.Data, machine, machineAssignment)
			if err != nil {
				reqLogger.Error(err, "couldn't parse template")
				return ctrl.Result{}, err
			}

			reqLogger.V(1).Info("deleting configmap", "name", "ipxe-"+data["name"])
			configMap, err := r.createConfigMap(data, &req)
			if err != nil {
				return ctrl.Result{}, err
			}

			if err := r.Delete(ctx, configMap); err != nil {
				reqLogger.Error(err, "couldn't delete config map", "resource", req.Name, "namespace", req.Namespace)
			}
		}

		if secretTemplateCM != nil {
			data, err := parseTemplate(templateCM.Data, machine, machineAssignment)
			if err != nil {
				reqLogger.Error(err, "couldn't parse template")
				return ctrl.Result{}, err
			}

			reqLogger.V(1).Info("deleting secret", "name", "ipxe-"+data["name"])
			secret, err := r.createSecret(data, &req)
			if err != nil {
				return ctrl.Result{}, err
			}
			if err := r.Delete(ctx, secret); err != nil {
				reqLogger.Error(err, "couldn't delete secret", "resource", req.Name, "namespace", req.Namespace)
			}
		}
		return ctrl.Result{}, nil
	}

	reqLogger.V(2).Info("resources", "machine assignment",
		fmt.Sprintf("%+v", machineAssignment), "machine", fmt.Sprintf("%+v", machine))

	if templateCM != nil {
		data, err := parseTemplate(templateCM.Data, machine, machineAssignment)
		if err != nil {
			reqLogger.Error(err, "couldn't parse template")
			return ctrl.Result{}, err
		}

		configMap, err := r.createConfigMap(data, &req)
		if err != nil {
			return ctrl.Result{}, err
		}

		reqLogger.V(1).Info("applying ignition configuration", "ignition", client.ObjectKeyFromObject(configMap))
		if err := r.Patch(ctx, configMap, client.Apply, ignitionFieldOwner, client.ForceOwnership); err != nil {
			reqLogger.Error(err, "couldn't create config map", "resource", req.Name, "namespace", req.Namespace)
			return ctrl.Result{}, err
		}
	}

	if secretTemplateCM != nil {
		data, err := parseTemplate(secretTemplateCM.Data, machine, machineAssignment)
		if err != nil {
			reqLogger.Error(err, "couldn't parse template")
			return ctrl.Result{}, err
		}

		secret, err := r.createSecret(data, &req)
		if err != nil {
			return ctrl.Result{}, err
		}

		reqLogger.V(1).Info("applying ignition secret configuration", "ignition secret", client.ObjectKeyFromObject(secret))
		if err := r.Patch(ctx, secret, client.Apply, ignitionFieldOwner, client.ForceOwnership); err != nil {
			reqLogger.Error(err, "couldn't create secret", "resource", req.Name, "namespace", req.Namespace)
			return ctrl.Result{}, err
		}
	}

	reqLogger.V(1).Info("reconciliation finished")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *IgnitionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha2.MachineAssignment{}).
		Complete(r)
}

func parseTemplate(temp map[string]string,
	machine *v1alpha2.Machine,
	machineAssignment *v1alpha2.MachineAssignment) (map[string]string, error) {
	var tempStr = ""
	for tempKey, tempVal := range temp {
		tempStr += tempKey + ": |\n  " + tempVal + "\n"
	}

	t, err := template.New("temporaryTemplate").Parse(tempStr)
	if err != nil {
		return nil, err
	}

	wrapper := MachineWrapper{
		Machine:           machine,
		MachineAssignment: machineAssignment,
	}

	var b bytes.Buffer
	err = t.Execute(&b, wrapper)
	if err != nil {
		return nil, err
	}

	var tempMap = make(map[string]string)
	if err = yaml.Unmarshal(b.Bytes(), &tempMap); err != nil {
		return nil, err
	}

	return tempMap, nil
}

func (r *IgnitionReconciler) createConfigMap(temp map[string]string, req *ctrl.Request) (*corev1.ConfigMap, error) {
	if _, ok := temp["name"]; !ok {
		return nil, errTemplateNameIsMissed
	}
	temp["name"] = strings.TrimSuffix(temp["name"], "\n")
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ipxe-" + temp["name"],
			Namespace: req.Namespace,
		},
		Data: temp,
	}
	return configMap, nil
}

func (r *IgnitionReconciler) createSecret(temp map[string]string, req *ctrl.Request) (*corev1.Secret, error) {
	if _, ok := temp["name"]; !ok {
		return nil, errTemplateNameIsMissed
	}
	temp["name"] = strings.TrimSuffix(temp["name"], "\n")
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ipxe-" + temp["name"],
			Namespace: req.Namespace,
		},
		StringData: temp,
	}
	return secret, nil
}
