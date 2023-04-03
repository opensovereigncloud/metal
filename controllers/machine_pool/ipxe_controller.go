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
	"text/template"

	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/go-logr/logr"
	"github.com/onmetal/metal-api/apis/machine/v1alpha2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// IpxeReconciler reconciles a Ignition object.
type IpxeReconciler struct {
	client.Client

	Log    logr.Logger
	Scheme *runtime.Scheme
}

type MachineWrapper struct {
	Machine *v1alpha2.Machine `json:"machine"`
}

const ipxeDefaultTemplateName = "ipxe-default"

//+kubebuilder:rbac:groups=machine.onmetal.de,resources=machines,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=machines/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=machines/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete

func (r *IpxeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("namespace", req.NamespacedName).V(1)

	log.Info("fetching template configmaps")
	ipxeDefaultCM := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ipxeDefaultTemplateName,
			Namespace: req.Namespace,
		},
	}
	if err := r.Client.Get(ctx, client.ObjectKeyFromObject(ipxeDefaultCM), ipxeDefaultCM); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("could not get config map, not found")
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}

		log.Error(err, "could not get config map")
		return ctrl.Result{}, err
	}

	machine := &v1alpha2.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
		},
	}
	if err := r.Client.Get(ctx, client.ObjectKeyFromObject(machine), machine); err != nil {
		log.Error(err, "could not get machine")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if machine.Status.Reservation.Reference == nil { // @TODO is it the case for deletion?
		data, err := parseTemplate(ipxeDefaultCM.Data, machine)
		if err != nil {
			log.Error(err, "couldn't parse template")
			return ctrl.Result{}, err
		}

		log.Info("deleting configmap", "name", "ipxe-"+data["name"])
		//configMap, err := r.createConfigMap(data, &req)
		//if err != nil {
		//	return ctrl.Result{}, err
		//}
		//
		//if err := r.Delete(ctx, configMap); err != nil {
		//	log.Error(err, "couldn't delete config map", "resource", req.Name, "namespace", req.Namespace)
		//}

		return ctrl.Result{}, nil
	}

	data, err := parseTemplate(ipxeDefaultCM.Data, machine)
	if err != nil {
		log.Error(err, "couldn't parse template")
		return ctrl.Result{}, err
	}

	configMap := r.createConfigMap(machine.Name, data, &req)
	log.Info("applying ipxe configuration", "ipxe", client.ObjectKeyFromObject(configMap))
	if err := r.Create(ctx, configMap); err != nil {
		log.Error(err, "couldn't create config map", "resource", req.Name, "namespace", req.Namespace)
		return ctrl.Result{}, err
	}

	log.Info("reconciliation finished")
	return ctrl.Result{}, nil
}

func (r *IpxeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha2.Machine{}).
		Complete(r)
}

func parseTemplate(temp map[string]string, machine *v1alpha2.Machine) (map[string]string, error) {
	var tempStr = ""
	for tempKey, tempVal := range temp {
		tempStr += tempKey + ": |\n  " + tempVal + "\n"
	}

	t, err := template.New("temporaryTemplate").Parse(tempStr)
	if err != nil {
		return nil, err
	}

	wrapper := MachineWrapper{
		Machine: machine,
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

func (r *IpxeReconciler) createConfigMap(name string, temp map[string]string, req *ctrl.Request) *corev1.ConfigMap {
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ipxe-" + name,
			Namespace: req.Namespace,
		},
		Data: temp,
	}

	return configMap
}
