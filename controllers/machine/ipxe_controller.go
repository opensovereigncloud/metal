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
	"strings"
	"text/template"

	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/onmetal/onmetal-image/oci/remote"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"

	"github.com/go-logr/logr"
	"github.com/onmetal/metal-api/apis/machine/v1alpha2"
	onmetalimage "github.com/onmetal/onmetal-image"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// IpxeReconciler reconciles a Ignition object.
type IpxeReconciler struct {
	client.Client

	Log         logr.Logger
	Scheme      *runtime.Scheme
	ImageParser ImageParser
}

type MachineWrapper struct {
	Machine *v1alpha2.Machine `json:"machine"`
}

const (
	IpxeDefaultTemplateName = "ipxe-default"
	OnmetalImage            = "ghcr.io/onmetal/onmetal-image/gardenlinux:1099"
)

//+kubebuilder:rbac:groups=machine.onmetal.de,resources=machines,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=machines/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=machines/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete

func (r *IpxeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("namespace", req.NamespacedName).V(0)

	log.Info("reconcile started")

	_, err := r.ImageParser.GetDescription()
	if err != nil {
		log.Error(err, "could not get image description")
		return ctrl.Result{}, err
	}

	log.Info("fetching template configmaps")
	ipxeDefaultCM := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      IpxeDefaultTemplateName,
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

	err = r.Client.Get(ctx, client.ObjectKeyFromObject(configMap), configMap)
	if apierrors.IsNotFound(err) {
		log.Info("config map for machine not found, create new ipxe configuration", "ipxe", client.ObjectKeyFromObject(configMap))

		if err := r.Create(ctx, configMap); err != nil {
			log.Error(err, "couldn't create config map")
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, nil
	}

	if err != nil {
		log.Error(err, "could not get config map")
		return ctrl.Result{}, err
	}

	// @TODO update CM

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

type ImageParser interface {
	GetDescription() (ImageDescription, error)
}

type ImageDescription struct {
	KernelDigest    string
	InitRAMFsDigest string
	RootFSDigest    string
	CommandLine     string
}

type OnmetalImageParser struct {
	Log      logr.Logger
	Registry *remote.Registry
}

func (p *OnmetalImageParser) GetDescription() (ImageDescription, error) {
	var imageDescription ImageDescription

	imageDescription, err := p.describeImage()
	if err != nil {
		p.Log.Error(err, "could not describe image")
		return imageDescription, err
	}

	return imageDescription, nil
}

func (p *OnmetalImageParser) describeImage() (ImageDescription, error) {
	var imageDescription ImageDescription

	onmetalImage, err := p.getOnmetalImage()
	if err != nil {
		p.Log.Error(err, "could not get onmetal image")
		return imageDescription, err
	}

	if onmetalImage.Kernel != nil {
		imageDescription.KernelDigest = p.formatDigest(string(onmetalImage.Kernel.Descriptor().Digest))
	}

	if onmetalImage.InitRAMFs != nil {
		imageDescription.InitRAMFsDigest = p.formatDigest(string(onmetalImage.InitRAMFs.Descriptor().Digest))
	}

	if onmetalImage.RootFS != nil {
		imageDescription.RootFSDigest = p.formatDigest(string(onmetalImage.RootFS.Descriptor().Digest))
	}

	imageDescription.CommandLine = onmetalImage.Config.CommandLine

	return imageDescription, nil
}

func (p *OnmetalImageParser) getOnmetalImage() (*onmetalimage.Image, error) {
	ociImage, err := p.Registry.Resolve(context.Background(), OnmetalImage)
	if err != nil {
		p.Log.Error(err, "registry resolving failed")
		return nil, err
	}

	onmetalImage, err := onmetalimage.ResolveImage(context.Background(), ociImage)
	if err != nil {
		p.Log.Error(err, "image resolving failed")
		return nil, err
	}

	return onmetalImage, nil
}

// remove sha256 prefix.
func (p *OnmetalImageParser) formatDigest(digest string) string {
	separatedStrings := strings.Split(digest, ":")

	if len(separatedStrings) > 1 {
		return separatedStrings[1]
	}

	return ""
}
