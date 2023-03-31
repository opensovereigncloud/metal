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

package controllers

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	kubeadmv1beta2 "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta2"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	machinev1alpha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
)

const (
	inventoryServiceAccountPrefix = "inventory-"
	kubeconfigSecretPrefix        = "kubeconfig-inventory-"
	inventoryNamespace            = "metal-api-system"
	inventoryKind                 = "Inventory"
	kubeadmConfigConfigMapName    = "kubeadm-config"
	kubeSystemNamespace           = "kube-system"
	clusterStatusKey              = "ClusterStatus"
)

// AccessReconciler reconciles an Inventory object for creating a dedicated kubeconfig.
type AccessReconciler struct {
	client.Client

	BootstrapAPIServer string
	Log                logr.Logger
	Scheme             *runtime.Scheme
}

// +kubebuilder:rbac:groups=machine.onmetal.de,resources=inventories,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=machine.onmetal.de,resources=inventories/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=machine.onmetal.de,resources=inventories/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=serviceaccounts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=rolebindings,verbs=get;list;watch;create;update;patch;delete

func (r *AccessReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("access", req.NamespacedName)

	inventory := &machinev1alpha1.Inventory{}

	err := r.Get(ctx, req.NamespacedName, inventory)
	if apierrors.IsNotFound(err) {
		log.Error(err, "requested inventory resource not found", "name", req.NamespacedName)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	if err != nil {
		log.Error(err, "unable to get inventory resource", "name", req.NamespacedName)
		return ctrl.Result{}, err
	}

	machineInventoryServiceAccountName := inventoryServiceAccountPrefix + inventory.Name
	machineInventoryNamespace := inventoryNamespace
	machineInventoryServiceAccount := &corev1.ServiceAccount{
		ObjectMeta: ctrl.ObjectMeta{
			Namespace: machineInventoryNamespace,
			Name:      machineInventoryServiceAccountName,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         machinev1alpha1.GroupVersion.Version,
					Kind:               inventoryKind,
					Name:               inventory.Name,
					UID:                inventory.UID,
					Controller:         pointer.Bool(true),
					BlockOwnerDeletion: pointer.Bool(true),
				},
			},
		},
	}

	currentMachineInventoryServiceAccount := &corev1.ServiceAccount{}
	err = r.Client.Get(ctx, client.ObjectKeyFromObject(machineInventoryServiceAccount), currentMachineInventoryServiceAccount)
	if err != nil && apierrors.IsNotFound(err) {
		err = r.Client.Create(ctx, machineInventoryServiceAccount)
		if err != nil {
			log.Error(err, "unable to create machine inventory service account", "name", machineInventoryServiceAccountName)
			return ctrl.Result{}, err
		}
	} else if err != nil {
		log.Error(err, "unable to get machine inventory service account", "name", machineInventoryServiceAccountName)
		return ctrl.Result{}, err
	}

	saTokenSecret := &corev1.Secret{
		ObjectMeta: ctrl.ObjectMeta{
			Namespace: machineInventoryNamespace,
			Name:      machineInventoryServiceAccountName,
			Annotations: map[string]string{
				corev1.ServiceAccountNameKey: machineInventoryServiceAccountName,
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         machinev1alpha1.GroupVersion.Version,
					Kind:               inventoryKind,
					Name:               inventory.Name,
					UID:                inventory.UID,
					Controller:         pointer.Bool(true),
					BlockOwnerDeletion: pointer.Bool(true),
				},
			},
		},
		Type: corev1.SecretTypeServiceAccountToken,
	}

	err = r.Client.Get(ctx, client.ObjectKeyFromObject(saTokenSecret), saTokenSecret)
	if err != nil && apierrors.IsNotFound(err) {
		err = r.Client.Create(ctx, saTokenSecret)
		if err != nil {
			log.Error(err, "unable to create machine inventory service account secret", "name", machineInventoryServiceAccountName)
			return ctrl.Result{}, err
		}
	} else if err != nil {
		log.Error(err, "unable to get machine inventory service account secret", "name", machineInventoryServiceAccountName)
		return ctrl.Result{}, err
	}

	machineInventoryRole := &rbacv1.Role{
		ObjectMeta: ctrl.ObjectMeta{
			Namespace: machineInventoryNamespace,
			Name:      machineInventoryServiceAccountName,
		},
	}

	err = r.Client.Get(ctx, client.ObjectKeyFromObject(machineInventoryRole), machineInventoryRole)
	if err != nil && apierrors.IsNotFound(err) {
		machineInventoryRole = &rbacv1.Role{
			ObjectMeta: ctrl.ObjectMeta{
				Namespace: machineInventoryNamespace,
				Name:      machineInventoryServiceAccountName,
				OwnerReferences: []metav1.OwnerReference{
					{
						APIVersion:         machinev1alpha1.GroupVersion.Version,
						Kind:               inventoryKind,
						Name:               inventory.Name,
						UID:                inventory.UID,
						Controller:         pointer.Bool(true),
						BlockOwnerDeletion: pointer.Bool(true),
					},
				},
			},
			Rules: []rbacv1.PolicyRule{
				{
					Verbs: []string{
						"get",
						"create",
						"update",
						"patch",
					},
					APIGroups: []string{
						"machine.onmetal.de",
					},
					Resources: []string{
						"inventories",
					},
				},
			},
		}

		err = r.Client.Create(ctx, machineInventoryRole)
		if err != nil {
			logr.FromContextOrDiscard(ctx).Error(err, "unable to create machine inventory role", "name", machineInventoryServiceAccountName)
			return ctrl.Result{}, err
		}
	} else if err != nil {
		logr.FromContextOrDiscard(ctx).Error(err, "unable to get machine inventory role", "name", machineInventoryServiceAccountName)
		return ctrl.Result{}, err
	}

	machineInventoryRoleBinding := &rbacv1.RoleBinding{
		ObjectMeta: ctrl.ObjectMeta{
			Namespace: machineInventoryNamespace,
			Name:      machineInventoryServiceAccountName,
		},
	}
	err = r.Client.Get(ctx, client.ObjectKeyFromObject(machineInventoryRoleBinding), machineInventoryRoleBinding)
	if err != nil && apierrors.IsNotFound(err) {
		machineRoleBinding := &rbacv1.RoleBinding{
			ObjectMeta: ctrl.ObjectMeta{
				Namespace: machineInventoryNamespace,
				Name:      machineInventoryServiceAccountName,
				OwnerReferences: []metav1.OwnerReference{
					{
						APIVersion:         machinev1alpha1.GroupVersion.Version,
						Kind:               inventoryKind,
						Name:               inventory.Name,
						UID:                inventory.UID,
						Controller:         pointer.Bool(true),
						BlockOwnerDeletion: pointer.Bool(true),
					},
				},
			},
			Subjects: []rbacv1.Subject{
				{
					Kind:      "ServiceAccount",
					APIGroup:  corev1.GroupName,
					Name:      machineInventoryServiceAccount.Name,
					Namespace: machineInventoryServiceAccount.Namespace,
				},
			},
			RoleRef: rbacv1.RoleRef{
				APIGroup: rbacv1.GroupName,
				Kind:     "Role",
				Name:     machineInventoryRole.Name,
			},
		}

		err = r.Client.Create(ctx, machineRoleBinding)
		if err != nil {
			logr.FromContextOrDiscard(ctx).Error(err, "unable to create machine inventory role binding", "name", machineInventoryServiceAccountName)
			return ctrl.Result{}, err
		}
	} else if err != nil {
		logr.FromContextOrDiscard(ctx).Error(err, "unable to get rolebinding", "name", machineInventoryServiceAccountName)
		return ctrl.Result{}, err
	}

	err = r.Client.Get(ctx, client.ObjectKeyFromObject(machineInventoryServiceAccount), machineInventoryServiceAccount)
	if err != nil {
		logr.FromContextOrDiscard(ctx).Error(err, "unable to get machine inventory service account", "name", machineInventoryServiceAccountName)
		return ctrl.Result{}, err
	}

	machineInventoryTokenSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: machineInventoryServiceAccount.Namespace,
			Name:      machineInventoryServiceAccountName,
		},
	}
	err = r.Client.Get(ctx, client.ObjectKeyFromObject(machineInventoryTokenSecret), machineInventoryTokenSecret)
	if err != nil {
		logr.FromContextOrDiscard(ctx).Error(err, "unable to get machine inventory token secret", "name", machineInventoryServiceAccountName)
		return ctrl.Result{}, err
	}

	serverString := ""
	if r.BootstrapAPIServer == "" {
		kubeadmConfigConfigMap := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: kubeSystemNamespace,
				Name:      kubeadmConfigConfigMapName,
			},
		}
		err = r.Client.Get(ctx, client.ObjectKeyFromObject(kubeadmConfigConfigMap), kubeadmConfigConfigMap)
		if err != nil {
			logr.FromContextOrDiscard(ctx).Error(err, "unable to get kubeadm configmap", "name", machineInventoryServiceAccountName)
			return ctrl.Result{}, err
		}

		clusterStatusString, ok := kubeadmConfigConfigMap.Data[clusterStatusKey]
		if !ok {
			err := errors.New("ClusterStatus is missing in kubeadm configuration")
			logr.FromContextOrDiscard(ctx).Error(err, "", "name", machineInventoryServiceAccountName)
			return ctrl.Result{}, err
		}

		clusterStatus := &kubeadmv1beta2.ClusterStatus{}
		clusterStatusDecoder := k8syaml.NewYAMLOrJSONDecoder(strings.NewReader(clusterStatusString), 4096)
		err = clusterStatusDecoder.Decode(clusterStatus)
		if err != nil {
			logr.FromContextOrDiscard(ctx).Error(err, "unable to deserialize cluster status", "name", machineInventoryServiceAccountName)
			return ctrl.Result{}, err
		}

		if len(clusterStatus.APIEndpoints) == 0 {
			err = errors.New("cluster status has no API endpoints specified")
			logr.FromContextOrDiscard(ctx).Error(err, "", "name", machineInventoryServiceAccountName)
			return ctrl.Result{}, err
		}

		for _, v := range clusterStatus.APIEndpoints {
			serverString = fmt.Sprintf("https://%s:%d", v.AdvertiseAddress, v.BindPort)
			break
		}
	} else {
		serverString = r.BootstrapAPIServer
	}

	kubeconfig := clientcmdapi.Config{
		Kind:       "Config",
		APIVersion: "v1",
		Clusters: map[string]*clientcmdapi.Cluster{
			"onmetal": {
				CertificateAuthorityData: machineInventoryTokenSecret.Data["ca.crt"],
				Server:                   serverString,
			},
		},
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			machineInventoryServiceAccount.Name: {
				Token: string(machineInventoryTokenSecret.Data["token"]),
			},
		},
		Contexts: map[string]*clientcmdapi.Context{
			machineInventoryServiceAccount.Name + "@" + "onmetal": {
				Cluster:   "onmetal",
				AuthInfo:  machineInventoryServiceAccount.Name,
				Namespace: inventoryNamespace,
			},
		},
		CurrentContext: machineInventoryServiceAccount.Name + "@" + "onmetal",
	}

	kubeconfigBytes, err := clientcmd.Write(kubeconfig)
	if err != nil {
		logr.FromContextOrDiscard(ctx).Error(err, "unable to write kubeconfig", "name", machineInventoryServiceAccountName)
		return ctrl.Result{}, err
	}

	kubeconfigName := kubeconfigSecretPrefix + inventory.Name
	kubeconfigSecret := &corev1.Secret{
		ObjectMeta: ctrl.ObjectMeta{
			Namespace: machineInventoryNamespace,
			Name:      kubeconfigName,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         machinev1alpha1.GroupVersion.Version,
					Kind:               inventoryKind,
					Name:               inventory.Name,
					UID:                inventory.UID,
					Controller:         pointer.Bool(true),
					BlockOwnerDeletion: pointer.Bool(true),
				},
			},
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			"kubeconfig": kubeconfigBytes,
		},
	}
	currentKubeconfigSecret := &corev1.Secret{}
	err = r.Client.Get(ctx, client.ObjectKeyFromObject(kubeconfigSecret), currentKubeconfigSecret)
	if err != nil && apierrors.IsNotFound(err) {
		if err = r.Client.Create(ctx, kubeconfigSecret); err != nil {
			logr.FromContextOrDiscard(ctx).Error(err, "unable to create machine inventory kubeconfig", "name", kubeconfigName)
			return ctrl.Result{}, err
		}
	} else if err != nil {
		logr.FromContextOrDiscard(ctx).Error(err, "unable to get machine inventory kubeconfig", "name", kubeconfigName)
		return ctrl.Result{}, err
	} else {
		if !reflect.DeepEqual(currentKubeconfigSecret.Data["kubeconfig"], kubeconfigSecret.Data["kubeconfig"]) {
			logr.FromContextOrDiscard(ctx).Info("update machine inventory kubeconfig", "name", kubeconfigName)
			currentKubeconfigSecret.Data = kubeconfigSecret.Data
			err = r.Client.Update(ctx, currentKubeconfigSecret)
			if err != nil {
				logr.FromContextOrDiscard(ctx).Error(err, "unable to update machine inventory kubeconfig", "name", kubeconfigName)
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AccessReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&machinev1alpha1.Inventory{}).
		Complete(r)
}
