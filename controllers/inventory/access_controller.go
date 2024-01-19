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
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/utils/ptr"
	ctrlruntime "sigs.k8s.io/controller-runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"

	metalv1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"
)

const (
	InventoryServiceAccountPrefix = "inventory-"
	kubeconfigSecretPrefix        = "kubeconfig-inventory-"
	inventoryKind                 = "Inventory"
	kubernetesServiceName         = "kubernetes"
	kubernetesServiceNamespace    = "default"
)

type APIEndpoint struct {
	Address string
	Port    int32
}

// SetupWithManager sets up the controller with the Manager.
func (r *AccessReconciler) SetupWithManager(mgr ctrlruntime.Manager) error {
	return ctrlruntime.NewControllerManagedBy(mgr).
		For(&metalv1alpha4.Inventory{}).
		Complete(r)
}

// AccessReconciler reconciles an Inventory object for creating a dedicated kubeconfig.
type AccessReconciler struct {
	ctrlclient.Client

	BootstrapAPIServer string
	Log                logr.Logger
	Scheme             *k8sruntime.Scheme
}

// +kubebuilder:rbac:groups=metal.ironcore.dev,resources=inventories,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=metal.ironcore.dev,resources=inventories/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=metal.ironcore.dev,resources=inventories/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=serviceaccounts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=endpoints,verbs=get;list;watch
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=rolebindings,verbs=get;list;watch;create;update;patch;delete

func (r *AccessReconciler) Reconcile(ctx context.Context, req ctrlruntime.Request) (ctrlruntime.Result, error) {
	log := r.Log.WithValues("access", req.NamespacedName)

	inventory, err := r.ExtractInventory(ctx, req)
	if err != nil {
		log.Error(err, "unable to get inventory resource", "error", err)
		return ctrlruntime.Result{}, ctrlclient.IgnoreNotFound(err)
	}

	inventoryServiceAccount := r.InventoryServiceAccount(inventory)

	clusterAccessSecret, err := r.CreateAccountForServerToClusterAccess(
		ctx,
		inventoryServiceAccount,
		inventory)
	if err != nil {
		log.Info("unable to create service account for cluster access", "error", err)
		return ctrlruntime.Result{}, err
	}

	newKubeconfigSecret, err := r.ClientKubeconfigSecret(
		ctx,
		inventory,
		clusterAccessSecret,
		inventoryServiceAccount.Name,
		inventory.Namespace)
	if err != nil {
		log.Info("unable get client kubeconfig secret", "error", err)
		return ctrlruntime.Result{}, err
	}

	currentKubeconfigSecret, err := r.CurrentKubeConfig(ctx, newKubeconfigSecret)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			log.Error(err, "failed to get current kubeconfig for inventory",
				"name", newKubeconfigSecret.Name)
			return ctrlruntime.Result{}, err
		}
		if err = r.CreateNewKubeconfigForClusterAccess(ctx, newKubeconfigSecret); err != nil {
			log.Error(err, "unable to create inventory kubeconfig",
				"name", newKubeconfigSecret.Name)
			return ctrlruntime.Result{}, err
		}
		return ctrlruntime.Result{}, nil
	}

	if r.UpdateNeeded(currentKubeconfigSecret, newKubeconfigSecret) {
		if err = r.UpdateExistingKubeconfigForClusterAccess(
			ctx,
			currentKubeconfigSecret,
			newKubeconfigSecret); err != nil {
			log.Info("unable update existing kubeconfig for cluster access", "error", err)
			return ctrlruntime.Result{}, err
		}
	}
	return ctrlruntime.Result{}, nil
}

func (r *AccessReconciler) CurrentKubeConfig(
	ctx context.Context,
	kubeconfigSecret *corev1.Secret) (*corev1.Secret, error) {
	currentKubeconfigSecret := &corev1.Secret{}
	err := r.
		Client.
		Get(
			ctx,
			ctrlclient.ObjectKeyFromObject(kubeconfigSecret),
			currentKubeconfigSecret)
	if err != nil {
		return nil, err
	}
	return currentKubeconfigSecret, nil
}

func (r *AccessReconciler) CreateAccountForServerToClusterAccess(
	ctx context.Context,
	machineInventoryServiceAccount *corev1.ServiceAccount,
	inventory *metalv1alpha4.Inventory) (*corev1.Secret, error) {
	exist, err := r.InventoryServiceAccountExist(ctx, machineInventoryServiceAccount)
	if !exist {
		if err := r.CreateInventoryServiceAccount(ctx, machineInventoryServiceAccount); err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}

	serviceAccountSecretToken := r.InventoryServiceAccountSecretToken(
		inventory.Namespace,
		machineInventoryServiceAccount.Name,
		inventory)

	exist, err = r.InventoryServiceAccountSecretExist(ctx, serviceAccountSecretToken)
	if !exist {
		if err := r.CreateInventoryServiceAccountSecret(ctx, serviceAccountSecretToken); err != nil {
			return nil, errors.Wrapf(err, "unable to create machine inventory service account secret")
		}
	}
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get machine inventory service account secret")
	}

	if err = r.BindInventoryAccountWithPermissions(
		ctx,
		machineInventoryServiceAccount,
		inventory.Namespace,
		inventory); err != nil {
		return nil, err
	}
	return serviceAccountSecretToken, nil
}

func (r *AccessReconciler) BindInventoryAccountWithPermissions(
	ctx context.Context,
	machineInventoryServiceAccount *corev1.ServiceAccount,
	machineInventoryNamespace string,
	inventory *metalv1alpha4.Inventory) error {
	machineInventoryRole := r.InventoryBaseRole(
		machineInventoryNamespace,
		machineInventoryServiceAccount.Name)
	exist, err := r.InventoryServiceAccountRoleExist(ctx, machineInventoryRole)
	if !exist {
		machineInventoryRole = r.RoleForInventoryWithPermissions(
			machineInventoryNamespace,
			machineInventoryServiceAccount.Name,
			inventory)
		if err = r.CreateInventoryRoleWithPermissions(
			ctx,
			machineInventoryRole); err != nil {
			return errors.Wrapf(err, "unable to create machine inventory service account role")
		}
	}
	if err != nil {
		return errors.Wrapf(err, "unable to get machine inventory service account role")
	}

	if err = r.BindPermissionsToTheServiceAccount(
		ctx,
		machineInventoryServiceAccount,
		machineInventoryNamespace,
		inventory,
		machineInventoryRole); err != nil {
		return err
	}
	return nil
}

func (r *AccessReconciler) BindPermissionsToTheServiceAccount(
	ctx context.Context,
	machineInventoryServiceAccount *corev1.ServiceAccount,
	machineInventoryNamespace string,
	inventory *metalv1alpha4.Inventory,
	machineInventoryRole *rbacv1.Role) error {
	machineInventoryRoleBinding := r.BaseInventoryRoleBinding(
		machineInventoryNamespace,
		machineInventoryServiceAccount.Name)
	exist, err := r.InventoryServiceAccountRoleBindingExist(ctx, machineInventoryRoleBinding)
	if !exist {
		inventoryRoleBinding := r.InventoryBindingRoleForServiceAccount(
			machineInventoryRole,
			machineInventoryNamespace,
			machineInventoryServiceAccount,
			inventory)
		if err := r.CreateInventoryServiceAccountBindingRole(
			ctx,
			inventoryRoleBinding,
		); err != nil {
			return errors.Wrapf(err, "unable to create machine inventory service account role binding")
		}
	}
	if err != nil {
		return errors.Wrapf(err, "unable to get machine inventory service account role binding")
	}
	return nil
}

func (r *AccessReconciler) InventorySecretForServiceAccount(
	ctx context.Context, machineInventoryServiceAccount *corev1.ServiceAccount) (*corev1.Secret, error) {
	machineInventorySecret := r.InventorySecret(machineInventoryServiceAccount)
	if err := r.GetKubernetesObject(ctx, machineInventorySecret); err != nil {
		return nil, err
	}
	return machineInventorySecret, nil
}

func (r *AccessReconciler) APIServerEndpoint(ctx context.Context) (string, error) {
	if r.BootstrapAPIServer != "" {
		return r.BootstrapAPIServer, nil
	}
	kubernetesAPIEndpoint, err := r.RetrieveServerAddressAndPortFromCluster(ctx)
	if err != nil {
		return "", err
	}
	return SanitizeAPIServerEndpointString(kubernetesAPIEndpoint), nil
}

func SanitizeAPIServerEndpointString(apiEndpoint APIEndpoint) string {
	if strings.Contains(apiEndpoint.Address, ":") {
		return fmt.Sprintf("https://[%s]:%d",
			apiEndpoint.Address, apiEndpoint.Port)
	}
	return fmt.Sprintf("https://%s:%d",
		apiEndpoint.Address, apiEndpoint.Port)
}

func (r *AccessReconciler) ExtractKubernetesAPIEndpoint(ctx context.Context) (*corev1.Endpoints, error) {
	kubernetesEndpoints := &corev1.Endpoints{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: kubernetesServiceNamespace,
			Name:      kubernetesServiceName,
		},
	}
	if err := r.GetKubernetesObject(ctx, kubernetesEndpoints); err != nil {
		return nil, err
	}
	return kubernetesEndpoints, nil
}

func (r *AccessReconciler) BaseInventoryRoleBinding(
	machineInventoryNamespace string,
	machineInventoryServiceAccountName string) *rbacv1.RoleBinding {
	return &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: machineInventoryNamespace,
			Name:      machineInventoryServiceAccountName,
		},
	}
}

func (r *AccessReconciler) InventoryBaseRole(
	machineInventoryNamespace string,
	machineInventoryServiceAccountName string) *rbacv1.Role {
	return &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: machineInventoryNamespace,
			Name:      machineInventoryServiceAccountName,
		},
	}
}

func (r *AccessReconciler) ExtractInventory(ctx context.Context,
	req ctrlruntime.Request) (*metalv1alpha4.Inventory, error) {
	inventory := &metalv1alpha4.Inventory{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: req.Namespace,
			Name:      req.Name,
		},
	}
	if err := r.GetKubernetesObject(ctx, inventory); err != nil {
		return nil, err
	}
	return inventory, nil
}

func (r *AccessReconciler) InventorySecret(inventoryServiceAccount *corev1.ServiceAccount) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: inventoryServiceAccount.Namespace,
			Name:      inventoryServiceAccount.Name,
		},
	}
}

func (r *AccessReconciler) UpdateNeeded(
	currentKubeconfigSecret, newkubeconfigSecret *corev1.Secret) bool {
	return !reflect.DeepEqual(currentKubeconfigSecret.Data["kubeconfig"], newkubeconfigSecret.Data["kubeconfig"])
}

func (r *AccessReconciler) UpdateExistingKubeconfigForClusterAccess(
	ctx context.Context,
	currentKubeconfigSecret *corev1.Secret,
	kubeconfigSecret *corev1.Secret) error {
	currentKubeconfigSecret.Data = kubeconfigSecret.Data
	if err := r.Client.Update(ctx, currentKubeconfigSecret); err != nil {
		return errors.Wrapf(err, "unable to update machine inventory kubeconfig")
	}
	return nil
}

func (r *AccessReconciler) ClientKubeconfigSecret(
	ctx context.Context,
	inventory *metalv1alpha4.Inventory,
	clusterAccessSecret *corev1.Secret,
	inventoryServiceAccountName string,
	targetNamespace string) (*corev1.Secret, error) {
	kubeconfig, err := r.KubeconfigForServer(
		ctx,
		clusterAccessSecret,
		inventoryServiceAccountName,
		inventory.Namespace)
	if err != nil {
		return nil, err
	}

	kubeconfigBytes, err := clientcmd.Write(kubeconfig)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to marshal kubeconfig")
	}

	kubeconfigName := kubeconfigSecretPrefix + inventory.Name
	return &corev1.Secret{
		ObjectMeta: ctrlruntime.ObjectMeta{
			Namespace: targetNamespace,
			Name:      kubeconfigName,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         metalv1alpha4.SchemeGroupVersion.Version,
					Kind:               inventoryKind,
					Name:               inventory.Name,
					UID:                inventory.UID,
					Controller:         ptr.To(true),
					BlockOwnerDeletion: ptr.To(true),
				},
			},
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			"kubeconfig": kubeconfigBytes,
		},
	}, nil
}

func (r *AccessReconciler) KubeconfigForServer(
	ctx context.Context,
	inventoryTokenSecret *corev1.Secret,
	inventoryServiceAccountName string,
	targetNamespace string) (clientcmdapi.Config, error) {
	apiServerEndpoint, err := r.APIServerEndpoint(ctx)
	if err != nil {
		return clientcmdapi.Config{}, err
	}
	return r.kubernetesClientConfig(
		inventoryTokenSecret,
		inventoryServiceAccountName,
		targetNamespace,
		apiServerEndpoint), nil
}

func (r *AccessReconciler) kubernetesClientConfig(
	inventoryTokenSecret *corev1.Secret,
	inventoryServiceAccountName string,
	targetNamespace string,
	apiServerEndpoint string) clientcmdapi.Config {
	return clientcmdapi.Config{
		Kind:       "Config",
		APIVersion: "v1",
		Clusters: map[string]*clientcmdapi.Cluster{
			r.ClusterName(): {
				CertificateAuthorityData: inventoryTokenSecret.Data["ca.crt"],
				Server:                   apiServerEndpoint,
			},
		},
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			inventoryServiceAccountName: {
				Token: string(inventoryTokenSecret.Data["token"]),
			},
		},
		Contexts: map[string]*clientcmdapi.Context{
			inventoryServiceAccountName + "@" + r.ClusterName(): {
				Cluster:   r.ClusterName(),
				AuthInfo:  inventoryServiceAccountName,
				Namespace: targetNamespace,
			},
		},
		CurrentContext: inventoryServiceAccountName + "@" + r.ClusterName(),
	}
}

func (r *AccessReconciler) ClusterName() string {
	return "onmetal"
}

func (r *AccessReconciler) InventoryRoleBinding(
	machineInventoryNamespace string,
	inventory *metalv1alpha4.Inventory,
	machineInventoryServiceAccount *corev1.ServiceAccount,
	machineInventoryRole *rbacv1.Role) *rbacv1.RoleBinding {
	return &rbacv1.RoleBinding{
		ObjectMeta: ctrlruntime.ObjectMeta{
			Namespace: machineInventoryNamespace,
			Name:      machineInventoryServiceAccount.Name,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         metalv1alpha4.SchemeGroupVersion.Version,
					Kind:               inventoryKind,
					Name:               inventory.Name,
					UID:                inventory.UID,
					Controller:         ptr.To(true),
					BlockOwnerDeletion: ptr.To(true),
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
}

func (r *AccessReconciler) RoleForInventoryWithPermissions(
	machineInventoryNamespace string,
	machineInventoryServiceAccountName string,
	inventory *metalv1alpha4.Inventory) *rbacv1.Role {
	return &rbacv1.Role{
		ObjectMeta: ctrlruntime.ObjectMeta{
			Namespace: machineInventoryNamespace,
			Name:      machineInventoryServiceAccountName,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         metalv1alpha4.SchemeGroupVersion.Version,
					Kind:               inventoryKind,
					Name:               inventory.Name,
					UID:                inventory.UID,
					Controller:         ptr.To(true),
					BlockOwnerDeletion: ptr.To(true),
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
					"metal.ironcore.dev",
				},
				Resources: []string{
					"inventories",
				},
			},
		},
	}
}

func (r *AccessReconciler) InventoryServiceAccountSecretToken(
	machineInventoryNamespace string,
	machineInventoryServiceAccountName string,
	inventory *metalv1alpha4.Inventory) *corev1.Secret {
	saTokenSecret := &corev1.Secret{
		ObjectMeta: ctrlruntime.ObjectMeta{
			Namespace: machineInventoryNamespace,
			Name:      machineInventoryServiceAccountName,
			Annotations: map[string]string{
				corev1.ServiceAccountNameKey: machineInventoryServiceAccountName,
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         metalv1alpha4.SchemeGroupVersion.Version,
					Kind:               inventoryKind,
					Name:               inventory.Name,
					UID:                inventory.UID,
					Controller:         ptr.To(true),
					BlockOwnerDeletion: ptr.To(true),
				},
			},
		},
		Type: corev1.SecretTypeServiceAccountToken,
	}
	return saTokenSecret
}

func (r *AccessReconciler) CreateInventoryServiceAccount(
	ctx context.Context,
	inventoryServiceAccount *corev1.ServiceAccount) error {
	return r.Client.Create(ctx, inventoryServiceAccount)
}

func (r *AccessReconciler) InventoryServiceAccount(
	inventory *metalv1alpha4.Inventory) *corev1.ServiceAccount {
	machineInventoryServiceAccountName := InventoryServiceAccountPrefix + inventory.Name
	return &corev1.ServiceAccount{
		ObjectMeta: ctrlruntime.ObjectMeta{
			Namespace: inventory.Namespace,
			Name:      machineInventoryServiceAccountName,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         metalv1alpha4.SchemeGroupVersion.Version,
					Kind:               inventoryKind,
					Name:               inventory.Name,
					UID:                inventory.UID,
					Controller:         ptr.To(true),
					BlockOwnerDeletion: ptr.To(true),
				},
			},
		},
	}
}

func (r *AccessReconciler) InventoryServiceAccountExist(
	ctx context.Context,
	machineInventoryServiceAccount *corev1.ServiceAccount) (bool, error) {
	err := r.GetKubernetesObject(ctx, machineInventoryServiceAccount)
	if apierrors.IsNotFound(err) {
		return false, nil
	}
	return true, err
}

func (r *AccessReconciler) InventoryServiceAccountSecretExist(
	ctx context.Context,
	secret *corev1.Secret) (bool, error) {
	err := r.GetKubernetesObject(ctx, secret)
	if err == nil {
		return true, nil
	}
	if apierrors.IsNotFound(err) {
		return false, nil
	}
	return true, err
}

func (r *AccessReconciler) CreateInventoryServiceAccountSecret(ctx context.Context, secret *corev1.Secret) error {
	return r.CreateObjectInCluster(ctx, secret)
}

func (r *AccessReconciler) InventoryServiceAccountRoleExist(ctx context.Context, role *rbacv1.Role) (bool, error) {
	err := r.GetKubernetesObject(ctx, role)
	if err == nil {
		return true, nil
	}
	if apierrors.IsNotFound(err) {
		return false, nil
	}
	return true, err
}

func (r *AccessReconciler) InventoryServiceAccountRoleBindingExist(
	ctx context.Context, roleBinding *rbacv1.RoleBinding) (bool, error) {
	err := r.GetKubernetesObject(ctx, roleBinding)
	if err == nil {
		return true, nil
	}
	if apierrors.IsNotFound(err) {
		return false, nil
	}
	return true, err
}

func (r *AccessReconciler) CreateInventoryRoleWithPermissions(
	ctx context.Context,
	roleForInventory *rbacv1.Role) error {
	return r.CreateObjectInCluster(ctx, roleForInventory)
}

func (r *AccessReconciler) InventoryBindingRoleForServiceAccount(machineInventoryRole *rbacv1.Role,
	machineInventoryNamespace string,
	machineInventoryServiceAccount *corev1.ServiceAccount,
	inventory *metalv1alpha4.Inventory) *rbacv1.RoleBinding {
	return r.InventoryRoleBinding(
		machineInventoryNamespace,
		inventory,
		machineInventoryServiceAccount,
		machineInventoryRole)
}
func (r *AccessReconciler) CreateInventoryServiceAccountBindingRole(
	ctx context.Context,
	machineRoleBinding *rbacv1.RoleBinding) error {
	return r.CreateObjectInCluster(ctx, machineRoleBinding)
}

func (r *AccessReconciler) CreateObjectInCluster(ctx context.Context, object ctrlclient.Object) error {
	return r.
		Client.
		Create(
			ctx,
			object)
}

func (r *AccessReconciler) GetKubernetesObject(
	ctx context.Context,
	object ctrlclient.Object) error {
	err := r.
		Client.
		Get(
			ctx,
			ctrlclient.ObjectKeyFromObject(object),
			object)
	if err != nil {
		return err
	}
	return nil
}

func (r *AccessReconciler) CreateNewKubeconfigForClusterAccess(
	ctx context.Context,
	newKubeconfigSecret *corev1.Secret) error {
	return r.
		Client.
		Create(
			ctx,
			newKubeconfigSecret)
}

func (r *AccessReconciler) RetrieveServerAddressAndPortFromCluster(ctx context.Context) (APIEndpoint, error) {
	kubernetesEndpoints, err := r.ExtractKubernetesAPIEndpoint(ctx)
	if err != nil {
		return APIEndpoint{}, err
	}
	if len(kubernetesEndpoints.Subsets) == 0 {
		return APIEndpoint{}, errKubernetesEndpointIsEmpty
	}
	kubernetesAPIAddress := KubernetesAPIAddressFromEndpoint(kubernetesEndpoints.Subsets[0])
	if kubernetesAPIAddress == "" {
		return APIEndpoint{}, errKubernetesEndpointAddressIsEmpty
	}
	kubernetesAPIPort := KubernetesAPIPortFromEndpoint(kubernetesEndpoints.Subsets[0].Ports)
	if kubernetesAPIPort == 0 {
		return APIEndpoint{}, errKubernetesEndpointAddressPortIsEmpty
	}
	return APIEndpoint{
		Address: kubernetesAPIAddress,
		Port:    kubernetesAPIPort,
	}, nil
}

func KubernetesAPIAddressFromEndpoint(subset corev1.EndpointSubset) string {
	for address := range subset.Addresses {
		if subset.Addresses[address].IP != "" {
			return subset.Addresses[address].IP
		}
	}
	return ""
}

func KubernetesAPIPortFromEndpoint(ports []corev1.EndpointPort) int32 {
	for p := range ports {
		if ports[p].Port == 0 {
			continue
		}
		return ports[p].Port
	}
	return 0
}
