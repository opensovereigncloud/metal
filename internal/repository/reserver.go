package repository

import (
	"context"

	machinev1alpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	"github.com/onmetal/metal-api/internal/entity"
	metalerr "github.com/onmetal/metal-api/pkg/errors"
	oobv1 "github.com/onmetal/oob-controller/api/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type MachineReserverRepo struct {
	client ctrlclient.Client
}

func NewMachineReserverRepo(c ctrlclient.Client) *MachineReserverRepo {
	return &MachineReserverRepo{
		client: c,
	}
}

func (m *MachineReserverRepo) CreateReservation(ctx context.Context, e entity.Reservation) error {
	machine := &machinev1alpha2.Machine{}
	if err := m.client.Get(ctx, types.NamespacedName{Name: e.RequestName, Namespace: e.RequestNamespace}, machine); err != nil {
		return err
	}

	machine.Status.Reservation.Status = entity.ReservationStatusPending
	machine.Status.Reservation.Reference = prepareReferenceSpec(e.OrderName, e.OrderNamespace)

	return m.client.Status().Update(ctx, machine)
}

func (m *MachineReserverRepo) DeleteReservation(ctx context.Context, e entity.Order) error {
	machine := &machinev1alpha2.Machine{}
	if err := m.client.Get(ctx, types.NamespacedName{
		Name: e.Name, Namespace: e.Namespace}, machine); err != nil {
		return err
	}
	machine.Status.Reservation.Status = entity.ReservationStatusAvailable
	machine.Status.Reservation.Reference = nil

	return m.client.Status().Update(ctx, machine)
}

func (m *MachineReserverRepo) GetReservation(ctx context.Context, e entity.Order) (entity.Reservation, error) {
	machine := &machinev1alpha2.Machine{}
	if err := m.client.Get(ctx, types.NamespacedName{
		Name: e.Name, Namespace: e.Namespace}, machine); err != nil {
		return entity.Reservation{}, err
	}

	if machine.Status.Reservation.Status == entity.ReservationStatusAvailable {
		return entity.Reservation{
			RequestName:      e.Name,
			RequestNamespace: e.Namespace,
			Status:           entity.ReservationStatusAvailable,
		}, nil
	}
	if machine.Status.Reservation.Reference == nil {
		return entity.Reservation{}, metalerr.NotFound("reservation")
	}
	var status entity.ReservationStatus
	if machine.Status.Reservation.Status == "" {
		status = entity.ReservationStatusPending
	}
	return entity.Reservation{
		OrderName:        machine.Status.Reservation.Reference.Name,
		OrderNamespace:   machine.Status.Reservation.Reference.Namespace,
		RequestName:      machine.Name,
		RequestNamespace: machine.Namespace,
		Status:           status,
	}, nil
}

type oobMachine struct {
	client                      ctrlclient.Client
	name, namespace, powerState string
}

func (m *MachineReserverRepo) CheckIn(ctx context.Context, e entity.Reservation) error {
	machine := &machinev1alpha2.Machine{}
	if err := m.client.Get(ctx, types.NamespacedName{
		Name: e.RequestName, Namespace: e.RequestNamespace}, machine); err != nil {
		return err
	}
	if !machine.Status.OOB.Exist {
		return metalerr.NotFound("oob for machine")
	}
	s := oobMachine{
		client:     m.client,
		name:       machine.Status.OOB.Reference.Name,
		namespace:  machine.Status.OOB.Reference.Namespace,
		powerState: machinev1alpha2.MachinePowerStateON,
	}
	return s.changeServerPowerState(ctx)
}

func (m *MachineReserverRepo) CheckOut(ctx context.Context, e entity.Reservation) error {
	machine := &machinev1alpha2.Machine{}
	if err := m.client.Get(ctx, types.NamespacedName{
		Name: e.RequestName, Namespace: e.RequestNamespace}, machine); err != nil {
		return err
	}
	if !machine.Status.OOB.Exist {
		return metalerr.NotFound("oob for machine")
	}
	s := oobMachine{
		client:     m.client,
		name:       machine.Status.OOB.Reference.Name,
		namespace:  machine.Status.OOB.Reference.Namespace,
		powerState: machinev1alpha2.MachinePowerStateOFF,
	}
	return s.changeServerPowerState(ctx)
}

func (m *oobMachine) changeServerPowerState(ctx context.Context) error {
	oobObj, err := m.getOOBMachineByUUIDLabel(ctx)
	if err != nil {
		return err
	}
	if oobObj.Spec.PowerState == machinev1alpha2.MachinePowerStateOFF &&
		m.powerState == machinev1alpha2.MachinePowerStateOFF {
		return nil
	}
	if m.powerState == oobObj.Spec.PowerState && m.powerState == machinev1alpha2.MachinePowerStateON {
		m.powerState = machinev1alpha2.MachinePowerStateReset
	}

	oobObj.Spec.PowerState = m.powerState

	return m.client.Update(ctx, oobObj)
}

func (m *oobMachine) getOOBMachineByUUIDLabel(ctx context.Context) (*oobv1.Machine, error) {
	oob := &oobv1.Machine{}
	if err := m.client.Get(ctx, types.NamespacedName{
		Name:      m.name,
		Namespace: m.namespace,
	}, oob); err != nil {
		return nil, err
	}
	return oob, nil
}

func prepareReferenceSpec(requestName, namespace string) *machinev1alpha2.ResourceReference {
	return &machinev1alpha2.ResourceReference{
		Name: requestName, Namespace: namespace,
	}
}
