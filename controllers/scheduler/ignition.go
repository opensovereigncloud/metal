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
	"net"

	ipamv1alpha1 "github.com/onmetal/ipam/api/v1alpha1"

	"github.com/onmetal/metal-api/apis/machine/v1alpha2"
	"github.com/onmetal/metal-api/pkg/constants"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type IgnitionWrapper struct {
	Machine           *v1alpha2.Machine           `json:"machine"`
	MachineAssignment *v1alpha2.MachineAssignment `json:"machineAssignment"`
	Hostname          string                      `json:"hostname"`
	IPv6WithoutPrefix string                      `json:"ipv6WithoutPrefix"`
	IPv6              string                      `json:"ipv6"`
	RouterID          string                      `json:"routerID"`
	ASN               string                      `json:"asn"`
}

func (r *Reconciler) reconcileIgnition(ctx context.Context, machineAssignment *v1alpha2.MachineAssignment) error {
	reqLogger := r.Log.WithName("Ignition").
		WithValues("machineAssignment", machineAssignment.Name, "namespace", machineAssignment.Namespace)

	if machineAssignment.Status.MetalComputeRef == nil || machineAssignment.Status.MetalComputeRef.Name == "" {
		reqLogger.Info("machine is not yet reserved")
		return nil
	}

	configMapTemplateExists := true
	configMaptemplate := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: machineAssignment.Namespace,
			Name:      IgnitionConfigMapTemplateName,
		},
	}
	if err := r.Client.Get(ctx, client.ObjectKeyFromObject(configMaptemplate), configMaptemplate); err != nil {
		reqLogger.Error(err, "ConfigMap template is not available in the current namespace",
			"name", IgnitionConfigMapTemplateName, "namespace", machineAssignment.Namespace)
		configMapTemplateExists = false
	}
	secretTemplateExists := true
	secretTemplate := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: machineAssignment.Namespace,
			Name:      IgnitionSecretTemplateName,
		},
	}
	if err := r.Client.Get(ctx, client.ObjectKeyFromObject(secretTemplate), secretTemplate); err != nil {
		reqLogger.Error(err, "Secret template is not available in the current namespace", "name",
			IgnitionSecretTemplateName, "namespace", machineAssignment.Namespace)
		secretTemplateExists = false
	}

	if !configMapTemplateExists && !secretTemplateExists {
		return errors.New("no iPXE temples found")
	}

	machine := &v1alpha2.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: machineAssignment.Status.MetalComputeRef.Namespace,
			Name:      machineAssignment.Status.MetalComputeRef.Name,
		},
	}
	if err := r.Get(ctx, client.ObjectKeyFromObject(machine), machine); err != nil {
		reqLogger.Error(err, "couldn't get assigned machine in namespace",
			"machine", machineAssignment.Status.MetalComputeRef.Name,
			"namespace", machineAssignment.Status.MetalComputeRef.Namespace)
		return err
	}

	machineSubnet, err := r.getMachineSubnet(ctx, machineAssignment)
	if err != nil {
		return err
	}

	machineLoopbackIP, err := r.getMachineLoopbackIP(ctx, machineAssignment)
	if err != nil {
		return err
	}

	hostname := ""
	if value, ok := machineAssignment.Labels[ComputeNameLabel]; ok {
		hostname = value
	}
	if hostname == "" {
		hostname = machine.Name
	}
	wrapper := IgnitionWrapper{
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
			return err
		}

		name := IgnitionIpxePrefix + machineAssignment.Status.MetalComputeRef.Name
		namespace := machineAssignment.Status.MetalComputeRef.Namespace
		configMap := r.createConfigMap(data, name, namespace)
		reqLogger.Info("applying ignition configuration", "configMap", configMap.Name)
		if err := r.Patch(ctx, configMap, client.Apply, IgnitionFieldOwner, client.ForceOwnership); err != nil {
			reqLogger.Error(err, "couldn't create config map", "resource", machineAssignment.Name, "namespace", machineAssignment.Namespace)
			return err
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
					"machine", machineAssignment.Status.MetalComputeRef.Name,
					"namespace", machineAssignment.Namespace,
					"ignition", machineAssignment.Spec.Ignition.Name)
				return err
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
			return err
		}

		name := IgnitionIpxePrefix + machineAssignment.Status.MetalComputeRef.Name
		namespace := machineAssignment.Status.MetalComputeRef.Namespace
		secret := r.createSecret(data, name, namespace)
		reqLogger.Info("applying ignition secret configuration", "secret", secret.Name)
		if err := r.Patch(ctx, secret, client.Apply, IgnitionFieldOwner, client.ForceOwnership); err != nil {
			reqLogger.Error(err, "couldn't create secret", "resource", machineAssignment.Name, "namespace", machineAssignment.Namespace)
			return err
		}
	}

	return nil
}

func (r *Reconciler) ignitionCleanup(ctx context.Context, machineAssignment *v1alpha2.MachineAssignment) error {
	reqLogger := r.Log.WithName("Ignition")

	if machineAssignment.Status.MetalComputeRef != nil && machineAssignment.Status.MetalComputeRef.Name != "" {
		configMap := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      IgnitionIpxePrefix + machineAssignment.Status.MetalComputeRef.Name,
				Namespace: machineAssignment.Status.MetalComputeRef.Namespace,
			},
		}

		reqLogger.Info("deleting configmap", "name", configMap.Name)
		err := r.Delete(ctx, configMap)
		if err != nil && !apierrors.IsNotFound(err) {
			reqLogger.Error(err, "couldn't delete config map", "resource", configMap.Name, "namespace", configMap.Namespace)
			return err
		}

		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      IgnitionIpxePrefix + machineAssignment.Status.MetalComputeRef.Name,
				Namespace: machineAssignment.Status.MetalComputeRef.Namespace,
			},
		}

		reqLogger.Info("deleting secret", "name", secret.Name)
		err = r.Delete(ctx, secret)
		if err != nil && !apierrors.IsNotFound(err) {
			reqLogger.Error(err, "couldn't delete secret", "resource", secret.Name, "namespace", secret.Namespace)
			return err
		}
	}

	return nil
}

func (r *Reconciler) createConfigMap(temp map[string]string, name string, namespace string) *corev1.ConfigMap {
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
	return configMap
}

func (r *Reconciler) createSecret(temp map[string]string, name string, namespace string) *corev1.Secret {
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
	return secret
}

func (r *Reconciler) getMachineSubnet(ctx context.Context, machineAssignment *v1alpha2.MachineAssignment) (*ipamv1alpha1.Subnet, error) {
	reqLogger := r.Log.WithValues(
		"namespace", machineAssignment.Status.MetalComputeRef.Namespace,
		"name", machineAssignment.Status.MetalComputeRef.Name)

	ownerFilter := client.MatchingLabels{
		constants.IPAMObjectOwnerLabel: machineAssignment.Status.MetalComputeRef.Name,
	}
	machineSubnetList := &ipamv1alpha1.SubnetList{}
	if err := r.List(ctx, machineSubnetList, ownerFilter); err != nil {
		reqLogger.Error(err, "couldn't get subnet owned by machine",
			"owner", machineAssignment.Status.MetalComputeRef.Name)
		return nil, err
	}

	if len(machineSubnetList.Items) == 0 {
		err := errors.New("no subnet found")
		reqLogger.Error(err, "couldn't get subnet owned by machine",
			"owner", machineAssignment.Status.MetalComputeRef.Name)
		return nil, err
	}

	// TODO(flpeter) what if there are more subnets?
	machineSubnet := machineSubnetList.Items[0]
	if machineSubnet.Status.State != ipamv1alpha1.CFinishedSubnetState {
		err := errors.New("subnet state is not Finished")
		reqLogger.Error(err, "couldn't get subnet owned by machine",
			"owner", machineAssignment.Status.MetalComputeRef.Name)
		return nil, err
	}

	if machineSubnet.Status.Reserved == nil {
		err := errors.New("subnet is not reserved")
		reqLogger.Error(err, "couldn't get subnet owned by machine",
			"owner", machineAssignment.Status.MetalComputeRef.Name)
		return nil, err
	}

	return &machineSubnet, nil
}

func (r *Reconciler) getMachineLoopbackIP(ctx context.Context, machineAssignment *v1alpha2.MachineAssignment) (net.IP, error) {
	reqLogger := r.Log.WithValues(
		"namespace", machineAssignment.Status.MetalComputeRef.Namespace,
		"name", machineAssignment.Status.MetalComputeRef.Name)

	filter := client.MatchingLabels{
		constants.IPAMObjectOwnerLabel:   machineAssignment.Status.MetalComputeRef.Name,
		constants.IPAMObjectPurposeLabel: constants.IPAMLoopbackPurpose,
	}
	machineIPList := &ipamv1alpha1.IPList{}
	if err := r.List(ctx, machineIPList, filter); err != nil {
		reqLogger.Error(err, "couldn't get loopback ip owned by machine",
			"owner", machineAssignment.Status.MetalComputeRef.Name)
		return nil, err
	}

	if len(machineIPList.Items) == 0 {
		err := errors.New("no ip found")
		reqLogger.Error(err, "couldn't get loopback ip owned by machine",
			"owner", machineAssignment.Status.MetalComputeRef.Name)
		return nil, err
	}

	// TODO(flpeter) what if there are more subnets?
	machineIP := machineIPList.Items[0]
	if machineIP.Status.State != ipamv1alpha1.CFinishedIPState {
		err := errors.New("ip state is not Finished")
		reqLogger.Error(err, "couldn't get loopback ip owned by machine",
			"owner", machineAssignment.Status.MetalComputeRef.Name)
		return nil, err
	}

	if machineIP.Status.Reserved == nil {
		err := errors.New("ip is not reserved")
		reqLogger.Error(err, "couldn't get loopback ip owned by machine",
			"owner", machineAssignment.Status.MetalComputeRef.Name)
		return nil, err
	}

	ip, _, err := net.ParseCIDR(machineIP.Status.Reserved.AsCidr().String())
	if err != nil {
		reqLogger.Error(err, "couldn't get loopback ip owned by machine",
			"owner", machineAssignment.Status.MetalComputeRef.Name)
		return nil, err
	}
	return ip, nil
}
