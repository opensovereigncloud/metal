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
	ipamv1alpha1 "github.com/onmetal/ipam/api/v1alpha1"
	"github.com/onmetal/metal-api/apis/switch/v1beta1"
	"math"
	"net"
	"strconv"
	"text/template"

	"github.com/Masterminds/sprig"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	syaml "sigs.k8s.io/yaml"

	"github.com/go-logr/logr"
	"github.com/onmetal/metal-api/apis/machine/v1alpha2"
)

const (
	ignitionFieldOwner    = client.FieldOwner("metal-api.onmetal.de/ignition")
	finalizer             = "metal-api.onmetal.de/ignition"
	configMapTemplateName = "ipxe-template"
	secretTemplateName    = "ignition-template"
	ipxePrefix            = "ipxe-"
	// TODO(flpeter) duplicate code from switch configurer
	asnBase = 4200000000
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
	Hostname          string                      `json:"hostname"`
	IPv6WithoutPrefix string                      `json:"ipv6WithoutPrefix"`
	IPv6              string                      `json:"ipv6"`
	RouterID          string                      `json:"routerID"`
	ASN               string                      `json:"asn"`
}

//+kubebuilder:rbac:groups=machine.machine.onmetal.de,resources=ignitions,verbs=get;list;watch
//+kubebuilder:rbac:groups=machine.machine.onmetal.de,resources=ignitions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *IgnitionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := r.Log.WithValues("namespace", req.NamespacedName)

	reqLogger.V(1).Info("reconciling ignition", "ignition", req)
	configMapTemplateExists := true
	configMaptemplate := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: req.Namespace,
			Name:      configMapTemplateName,
		},
	}
	if err := r.Client.Get(ctx, client.ObjectKeyFromObject(configMaptemplate), configMaptemplate); err != nil {
		reqLogger.Error(err, "ConfigMap template is not available in the current namespace", "name", configMapTemplateName, "namespace", req.Namespace)
		configMapTemplateExists = false
	}
	secretTemplateExists := true
	secretTemplate := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: req.Namespace,
			Name:      secretTemplateName,
		},
	}
	if err := r.Client.Get(ctx, client.ObjectKeyFromObject(secretTemplate), secretTemplate); err != nil {
		reqLogger.Error(err, "Secret template is not available in the current namespace", "name", secretTemplateName, "namespace", req.Namespace)
		secretTemplateExists = false
	}

	if !configMapTemplateExists && !secretTemplateExists {
		return ctrl.Result{}, errors.New("no iPXE temples found")
	}

	reqLogger.V(1).Info("fetching machine assignment resource", "machine assignment", req)
	machineAssignment := &v1alpha2.MachineAssignment{}
	if err := r.Get(ctx, req.NamespacedName, machineAssignment); err != nil {
		reqLogger.Error(err, "couldn't get machine assignment in namespace", "machine assignment", req.Name, "namespace", req.Namespace)
		return ctrl.Result{}, err
	}

	if machineAssignment.Status.MachineRef == nil || machineAssignment.Status.MachineRef.Name == "" {
		reqLogger.V(1).Info("machine is not yet reserved")
		return ctrl.Result{}, nil
	}

	machine := &v1alpha2.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: machineAssignment.Status.MachineRef.Namespace,
			Name:      machineAssignment.Status.MachineRef.Name,
		},
	}
	if err := r.Get(ctx, client.ObjectKeyFromObject(machine), machine); err != nil {
		reqLogger.Error(err, "couldn't get assigned machine in namespace", "machine", req.Name, "namespace", req.Namespace)
		return ctrl.Result{}, err
	}

	// examine DeletionTimestamp to determine if object is under deletion
	if machineAssignment.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.
		if !controllerutil.ContainsFinalizer(machineAssignment, finalizer) {
			controllerutil.AddFinalizer(machineAssignment, finalizer)
			if err := r.Client.Update(ctx, machineAssignment); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		// The object is being deleted
		if controllerutil.ContainsFinalizer(machineAssignment, finalizer) {
			// our finalizer is present, so lets handle any external dependency
			if configMaptemplate != nil {
				configMap := &corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      ipxePrefix + machineAssignment.Status.MachineRef.Name,
						Namespace: machineAssignment.Status.MachineRef.Namespace,
					},
				}

				reqLogger.V(1).Info("deleting configmap", "name", configMap.Name)
				err := r.Delete(ctx, configMap)
				if err != nil && !apierrors.IsNotFound(err) {
					reqLogger.Error(err, "couldn't delete config map", "resource", configMap.Name, "namespace", configMap.Namespace)
					return ctrl.Result{}, err
				}
			}

			if secretTemplate != nil {
				secret := &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      ipxePrefix + machineAssignment.Status.MachineRef.Name,
						Namespace: machineAssignment.Status.MachineRef.Namespace,
					},
				}

				reqLogger.V(1).Info("deleting secret", "name", secret.Name)
				err := r.Delete(ctx, secret)
				if err != nil && !apierrors.IsNotFound(err) {
					reqLogger.Error(err, "couldn't delete secret", "resource", secret.Name, "namespace", secret.Namespace)
					return ctrl.Result{}, err
				}
			}

			// remove our finalizer from the list and update it.
			controllerutil.RemoveFinalizer(machineAssignment, finalizer)
			if err := r.Client.Update(ctx, machineAssignment); err != nil {
				return ctrl.Result{}, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}

	reqLogger.V(2).Info("resources", "machine assignment", fmt.Sprintf("%#v", machineAssignment), "machine", fmt.Sprintf("%#v", machine))

	machineSubnet, err := r.getMachineSubnet(ctx, machineAssignment)
	if err != nil {
		return ctrl.Result{}, err
	}

	machineLoopbackIP, err := r.getMachineLoopbackIP(ctx, machineAssignment)
	if err != nil {
		return ctrl.Result{}, err
	}

	hostname := ""
	if value, ok := machineAssignment.Labels["machine.onmetal.de/compute-name"]; ok {
		hostname = value
	}
	if hostname == "" {
		hostname = machine.Name
	}
	wrapper := MachineWrapper{
		Machine:           machine,
		MachineAssignment: machineAssignment,
		Hostname:          hostname,
		IPv6WithoutPrefix: machineSubnet.Status.Reserved.AsIPAddr().Net.String(),
		IPv6:              machineSubnet.Status.Reserved.String(),
		RouterID:          machineLoopbackIP.String(),
		ASN:               calculateAsn(machineLoopbackIP),
	}

	if configMapTemplateExists {
		data, err := parseTemplate(configMaptemplate.Data, wrapper)
		if err != nil {
			reqLogger.Error(err, "couldn't parse template")
			return ctrl.Result{}, err
		}

		name := ipxePrefix + machineAssignment.Status.MachineRef.Name
		namespace := machineAssignment.Status.MachineRef.Namespace
		configMap, err := r.createConfigMap(data, name, namespace)
		if err != nil {
			return ctrl.Result{}, err
		}

		reqLogger.V(1).Info("applying ignition configuration", "ignition", client.ObjectKeyFromObject(configMap))
		if err := r.Patch(ctx, configMap, client.Apply, ignitionFieldOwner, client.ForceOwnership); err != nil {
			reqLogger.Error(err, "couldn't create config map", "resource", req.Name, "namespace", req.Namespace)
			return ctrl.Result{}, err
		}
	}

	if secretTemplateExists {
		templateData := map[string]string{}
		if machineAssignment.Spec.Ignition != nil {
			machineIgnitionSecret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      machineAssignment.Spec.Ignition.Name,
					Namespace: machineAssignment.Namespace,
				},
			}
			if err := r.Get(ctx, client.ObjectKeyFromObject(machineIgnitionSecret), machineIgnitionSecret); err != nil {
				reqLogger.Error(err, "couldn't get ignition for machine in namespace",
					"machine", machineAssignment.Status.MachineRef.Name,
					"namespace", req.Namespace,
					"ignition", machineAssignment.Spec.Ignition.Name)
				return ctrl.Result{}, err
			}

			if len(machineIgnitionSecret.Data) > 0 {
				for k, v := range machineIgnitionSecret.Data {
					templateData[k] = string(v)
				}
			}
		}

		if len(secretTemplate.Data) > 0 {
			for k, v := range secretTemplate.Data {
				templateData[k] = string(v)
			}
		}
		data, err := parseTemplate(templateData, wrapper)
		if err != nil {
			reqLogger.Error(err, "couldn't parse template")
			return ctrl.Result{}, err
		}

		name := ipxePrefix + machineAssignment.Status.MachineRef.Name
		namespace := machineAssignment.Status.MachineRef.Namespace
		secret, err := r.createSecret(data, name, namespace)
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

func parseTemplate(temp map[string]string, wrapper MachineWrapper) (map[string]string, error) {
	tempStr, err := syaml.Marshal(temp)
	if err != nil {
		return nil, err
	}

	t, err := template.New("temporaryTemplate").Funcs(sprig.HermeticTxtFuncMap()).Parse(string(tempStr))
	if err != nil {
		return nil, err
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

func (r *IgnitionReconciler) createConfigMap(temp map[string]string, name string, namespace string) (*corev1.ConfigMap, error) {
	configMap := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: temp,
	}
	return configMap, nil
}

func (r *IgnitionReconciler) createSecret(temp map[string]string, name string, namespace string) (*corev1.Secret, error) {
	secret := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		StringData: temp,
	}
	return secret, nil
}

func (r *IgnitionReconciler) getMachineSubnet(ctx context.Context, machineAssignment *v1alpha2.MachineAssignment) (*ipamv1alpha1.Subnet, error) {
	reqLogger := r.Log.WithValues(
		"namespace", machineAssignment.Status.MachineRef.Namespace,
		"name", machineAssignment.Status.MachineRef.Name)

	ownerFilter := client.MatchingLabels{
		v1beta1.IPAMObjectOwnerLabel: machineAssignment.Status.MachineRef.Name,
	}
	machineSubnetList := &ipamv1alpha1.SubnetList{}
	if err := r.List(ctx, machineSubnetList, ownerFilter); err != nil {
		reqLogger.Error(err, "couldn't get subnet owned by machine",
			"owner", machineAssignment.Status.MachineRef.Name)
		return nil, err
	}

	if len(machineSubnetList.Items) == 0 {
		err := errors.New("no subnet found")
		reqLogger.Error(err, "couldn't get subnet owned by machine",
			"owner", machineAssignment.Status.MachineRef.Name)
		return nil, err
	}

	//TODO(flpeter) what if there are more subnets?
	machineSubnet := machineSubnetList.Items[0]
	if machineSubnet.Status.State != ipamv1alpha1.CFinishedSubnetState {
		err := errors.New("subnet state is not Finished")
		reqLogger.Error(err, "couldn't get subnet owned by machine",
			"owner", machineAssignment.Status.MachineRef.Name)
		return nil, err
	}

	if machineSubnet.Status.Reserved == nil {
		err := errors.New("subnet is not reserved")
		reqLogger.Error(err, "couldn't get subnet owned by machine",
			"owner", machineAssignment.Status.MachineRef.Name)
		return nil, err
	}

	return &machineSubnet, nil
}

func (r *IgnitionReconciler) getMachineLoopbackIP(ctx context.Context, machineAssignment *v1alpha2.MachineAssignment) (net.IP, error) {
	reqLogger := r.Log.WithValues(
		"namespace", machineAssignment.Status.MachineRef.Namespace,
		"name", machineAssignment.Status.MachineRef.Name)

	filter := client.MatchingLabels{
		v1beta1.IPAMObjectOwnerLabel:   machineAssignment.Status.MachineRef.Name,
		v1beta1.IPAMObjectPurposeLabel: v1beta1.CIPAMPurposeLoopback,
	}
	machineIPList := &ipamv1alpha1.IPList{}
	if err := r.List(ctx, machineIPList, filter); err != nil {
		reqLogger.Error(err, "couldn't get loopback ip owned by machine",
			"owner", machineAssignment.Status.MachineRef.Name)
		return nil, err
	}

	if len(machineIPList.Items) == 0 {
		err := errors.New("no ip found")
		reqLogger.Error(err, "couldn't get loopback ip owned by machine",
			"owner", machineAssignment.Status.MachineRef.Name)
		return nil, err
	}

	//TODO(flpeter) what if there are more subnets?
	machineIP := machineIPList.Items[0]
	if machineIP.Status.State != ipamv1alpha1.CFinishedIPState {
		err := errors.New("ip state is not Finished")
		reqLogger.Error(err, "couldn't get loopback ip owned by machine",
			"owner", machineAssignment.Status.MachineRef.Name)
		return nil, err
	}

	if machineIP.Status.Reserved == nil {
		err := errors.New("ip is not reserved")
		reqLogger.Error(err, "couldn't get loopback ip owned by machine",
			"owner", machineAssignment.Status.MachineRef.Name)
		return nil, err
	}

	ip, _, err := net.ParseCIDR(machineIP.Status.Reserved.AsCidr().String())
	if err != nil {
		reqLogger.Error(err, "couldn't get loopback ip owned by machine",
			"owner", machineAssignment.Status.MachineRef.Name)
		return nil, err
	}
	return ip, nil
}

// TODO(flpeter) duplicate code from switch configurer
func calculateAsn(addr net.IP) string {
	var asn int
	asn += int(addr[13]) * int(math.Pow(2, 16))
	asn += (int(addr[14])) * int(math.Pow(2, 8))
	asn += int(addr[15])
	return strconv.Itoa(asnBase + asn)
}
