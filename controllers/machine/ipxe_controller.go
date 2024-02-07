// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"bytes"
	"context"
	"strings"
	"text/template"

	metalv1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"

	"k8s.io/utils/ptr"

	computev1alpha1 "github.com/ironcore-dev/ironcore/api/compute/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/go-logr/logr"
	ironcoreimage "github.com/ironcore-dev/ironcore-image"
	"github.com/ironcore-dev/ironcore-image/oci/remote"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// IpxeReconciler reconciles an Ignition object.
type IpxeReconciler struct {
	client.Client

	Log         logr.Logger
	Scheme      *runtime.Scheme
	ImageParser ImageParser
	Templater   Templater
}

const (
	machineKind = "Machine"

	IpxeTemplate string = "    #!ipxe\n\n    kernel https://ghcr.io/layer/{{.KernelDigest}}\n    " +
		"initrd={{.InitRAMFsDigest}}\n    " +
		"gl.url=https://ghcr.io/layer/{{.RootFSDigest}} ignition.config.url=http://2a10:afc0:e013:d000::5b4f/ignition\n    " +
		"{{.CommandLine}}\n\n    " +
		"initrd https://ghcr.io/layer/{{.InitRAMFsDigest}}\n    " +
		"boot"
)

// +kubebuilder:rbac:groups=metal.ironcore.dev,resources=machines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=metal.ironcore.dev,resources=machines/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=metal.ironcore.dev,resources=machines/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete

func (r *IpxeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("namespace", req.NamespacedName).V(0)

	log.Info("reconcile started")

	machine := &metalv1alpha4.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
		},
	}
	if err := r.Client.Get(ctx, client.ObjectKeyFromObject(machine), machine); err != nil {
		log.Error(err, "could not get machine")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if machine.Status.Reservation.Reference == nil {
		log.Info("deleting configmap", "name", "ipxe-"+machine.Name)

		configMap := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "ipxe-" + machine.Name,
				Namespace: req.Namespace,
			},
		}

		if err := r.Delete(ctx, configMap); err != nil && !apierrors.IsNotFound(err) {
			log.Error(err, "couldn't delete config map", "resource", req.Name, "namespace", req.Namespace)
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, nil
	}

	computeMachine := &computev1alpha1.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      machine.Status.Reservation.Reference.Name,
			Namespace: machine.Status.Reservation.Reference.Namespace,
		},
	}
	if err := r.Client.Get(ctx, client.ObjectKeyFromObject(computeMachine), computeMachine); err != nil {
		log.Error(err, "could not get compute machine")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if computeMachine.Spec.Image == "" {
		log.Info("unable to handle ipxe CM, ironcore image url is empty")
		return ctrl.Result{}, nil
	}

	imageDescription, err := r.ImageParser.GetDescription(computeMachine.Spec.Image)
	if err != nil {
		log.Error(err, "could not get image description")
		return ctrl.Result{}, err
	}

	data, err := r.parseTemplate(imageDescription)
	if err != nil {
		log.Error(err, "couldn't parse template")
		return ctrl.Result{}, err
	}

	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ipxe-" + machine.Name,
			Namespace: req.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         metalv1alpha4.SchemeGroupVersion.Version,
					Kind:               machineKind,
					Name:               machine.Name,
					UID:                machine.UID,
					Controller:         ptr.To(true),
					BlockOwnerDeletion: ptr.To(true),
				},
			},
		},
		Data: data,
	}

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

	if err := r.Client.Update(ctx, configMap); err != nil {
		log.Error(err, "could not update config map")
		return ctrl.Result{}, err
	}

	log.Info("reconciliation finished")
	return ctrl.Result{}, nil
}

func (r *IpxeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&metalv1alpha4.Machine{}).
		Complete(r)
}

func (r *IpxeReconciler) parseTemplate(imageDescription ImageDescription) (map[string]string, error) {
	t, err := r.Templater.GetTemplate(IpxeTemplate)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	err = t.Execute(&b, imageDescription)
	if err != nil {
		return nil, err
	}

	tempMap := map[string]string{
		"name": b.String(),
	}

	return tempMap, nil
}

type ImageParser interface {
	GetDescription(url string) (ImageDescription, error)
}

type ImageDescription struct {
	KernelDigest    string
	InitRAMFsDigest string
	RootFSDigest    string
	CommandLine     string
}

type IroncoreImageParser struct {
	Log      logr.Logger
	Registry *remote.Registry
}

func (p *IroncoreImageParser) GetDescription(url string) (ImageDescription, error) {
	var imageDescription ImageDescription

	ironcoreImage, err := p.getIroncoreImage(url)
	if err != nil {
		p.Log.Error(err, "could not get ironcore image")
		return imageDescription, err
	}

	p.describeImage(&imageDescription, ironcoreImage)
	return imageDescription, nil
}

func (p *IroncoreImageParser) describeImage(imageDescription *ImageDescription, ironcoreImage *ironcoreimage.Image) {
	if ironcoreImage.Kernel != nil {
		imageDescription.KernelDigest = p.formatDigest(string(ironcoreImage.Kernel.Descriptor().Digest))
	}

	if ironcoreImage.InitRAMFs != nil {
		imageDescription.InitRAMFsDigest = p.formatDigest(string(ironcoreImage.InitRAMFs.Descriptor().Digest))
	}

	if ironcoreImage.RootFS != nil {
		imageDescription.RootFSDigest = p.formatDigest(string(ironcoreImage.RootFS.Descriptor().Digest))
	}

	imageDescription.CommandLine = ironcoreImage.Config.CommandLine
}

func (p *IroncoreImageParser) getIroncoreImage(url string) (*ironcoreimage.Image, error) {
	ociImage, err := p.Registry.Resolve(context.Background(), url)
	if err != nil {
		p.Log.Error(err, "registry resolving failed")
		return nil, err
	}

	ironcoreImage, err := ironcoreimage.ResolveImage(context.Background(), ociImage)
	if err != nil {
		p.Log.Error(err, "image resolving failed")
		return nil, err
	}

	return ironcoreImage, nil
}

// remove sha256 prefix.
func (p *IroncoreImageParser) formatDigest(digest string) string {
	separatedStrings := strings.Split(digest, ":")

	if len(separatedStrings) > 1 {
		return separatedStrings[1]
	}

	return ""
}

type Templater interface {
	GetTemplate(templateData string) (*template.Template, error)
}

type IpxeTemplater struct {
	template *template.Template
}

func (t *IpxeTemplater) GetTemplate(templateData string) (*template.Template, error) {
	if t.template != nil {
		return t.template, nil
	}

	template, err := template.New("template").Parse(templateData)
	if err != nil {
		return nil, err
	}

	t.template = template
	return t.template, nil
}
