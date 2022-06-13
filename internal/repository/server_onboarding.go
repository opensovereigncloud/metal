package repository

import (
	"context"
	"errors"

	inventoriesv1alpha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	machinev1alpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	"github.com/onmetal/metal-api/internal/entity"
	machinerr "github.com/onmetal/metal-api/pkg/errors"
	oobv1 "github.com/onmetal/oob-controller/api/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type ServerOnboardingRepo struct {
	client    ctrlclient.Client
	inventory *inventoriesv1alpha1.Inventory
}

func NewServerOnboardingRepo(c ctrlclient.Client) *ServerOnboardingRepo {
	return &ServerOnboardingRepo{
		client: c,
	}
}

func (o *ServerOnboardingRepo) Create(ctx context.Context) error {
	if err := o.client.Create(ctx, o.inventory); err != nil {
		if apierrors.IsAlreadyExists(err) {
			return nil
		}
		return err
	}
	return nil
}

func (o *ServerOnboardingRepo) IsInitialized(ctx context.Context, e entity.Onboarding) bool {
	oobObj := &oobv1.Machine{}
	if err := o.client.Get(ctx, types.NamespacedName{
		Name: e.RequestName, Namespace: e.RequestNamespace}, oobObj); err != nil {
		return true
	}

	inventory := &inventoriesv1alpha1.Inventory{}
	if err := o.client.Get(ctx, types.NamespacedName{
		Name: oobObj.Status.UUID, Namespace: e.InitializationObjectNamespace}, inventory); err != nil {
		if apierrors.IsNotFound(err) {
			return false
		}
		return false
	}
	return true
}

func (o *ServerOnboardingRepo) Prepare(ctx context.Context, e entity.Onboarding) error {
	oobObj := &oobv1.Machine{}
	if err := o.client.Get(ctx, types.NamespacedName{
		Name: e.RequestName, Namespace: e.RequestNamespace}, oobObj); err != nil {
		return err
	}

	if oobObj.Status.UUID == "" {
		return machinerr.UUIDNotExist(e.RequestName)
	}

	e.InitializationObjectName = oobObj.Status.UUID
	o.inventory = prepareInventory(e)

	return nil
}

func (o *ServerOnboardingRepo) GatherData(ctx context.Context, e entity.Onboarding) error {
	oob := &oobv1.Machine{}
	if err := o.client.Get(ctx, types.NamespacedName{
		Name: e.RequestName, Namespace: e.RequestNamespace}, oob); err != nil {
		return err
	}

	inventory := &inventoriesv1alpha1.Inventory{}
	if err := o.client.Get(ctx, types.NamespacedName{
		Name: oob.Status.UUID, Namespace: e.InitializationObjectNamespace}, inventory); err != nil {
		return err
	}

	if o.IsSizeLabeled(inventory.Labels) {
		inventory.Status.InventoryStatuses.RequestsCount = 0
		return o.client.Update(ctx, inventory)
	}

	if inventory.Status.InventoryStatuses.RequestsCount > 1 {
		return errors.New("machine was booted but inventory not appeared")
	}

	if err := o.enableOOBMachineForInventory(ctx, oob); err != nil {
		return err
	}

	inventory.Status.InventoryStatuses.RequestsCount = 1
	return o.client.Update(ctx, inventory)
}

func (o *ServerOnboardingRepo) IsSizeLabeled(labels map[string]string) bool {
	machine := labels[inventoriesv1alpha1.GetSizeMatchLabel(machineSizeName)]
	switches := labels[inventoriesv1alpha1.GetSizeMatchLabel(switchSizeName)]
	return machine != "" || switches != ""
}

func (o *ServerOnboardingRepo) enableOOBMachineForInventory(ctx context.Context, oobObj *oobv1.Machine) error {
	oobObj.Spec.PowerState = getPowerState(oobObj.Spec.PowerState)
	oobObj.Labels = setUpLabels(oobObj)
	return o.client.Update(ctx, oobObj)
}

func prepareInventory(e entity.Onboarding) *inventoriesv1alpha1.Inventory {
	return &inventoriesv1alpha1.Inventory{
		ObjectMeta: metav1.ObjectMeta{
			Name:      e.InitializationObjectName,
			Namespace: e.InitializationObjectNamespace,
		},
		Spec: inventoriesv1alpha1.InventorySpec{
			System: &inventoriesv1alpha1.SystemSpec{
				ID: e.InitializationObjectName,
			},
			Host: &inventoriesv1alpha1.HostSpec{
				Name: "",
			},
		},
	}
}

func getPowerState(state string) string {
	switch state {
	case "On":
		// In case when machine already running Reset is required.
		// Machine should be started from scratch.
		// return "Reset"
		return state
	default:
		return "On"
	}
}

func setUpLabels(oobObj *oobv1.Machine) map[string]string {
	if oobObj.Labels == nil {
		return map[string]string{machinev1alpha2.UUIDLabel: oobObj.Status.UUID}
	}
	if _, ok := oobObj.Labels[machinev1alpha2.UUIDLabel]; !ok {
		oobObj.Labels[machinev1alpha2.UUIDLabel] = oobObj.Status.UUID
	}
	return oobObj.Labels
}
