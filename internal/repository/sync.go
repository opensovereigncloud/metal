package repository

import (
	"context"

	machinev1alpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	"github.com/onmetal/metal-api/internal/entity"
	"k8s.io/apimachinery/pkg/types"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type AssignmentSynchronizationRepo struct {
	client ctrlclient.Client
}

func NewAssignmentSynchronizationRepo(c ctrlclient.Client) *AssignmentSynchronizationRepo {
	return &AssignmentSynchronizationRepo{
		client: c,
	}
}

func (s *AssignmentSynchronizationRepo) Do(ctx context.Context, e entity.Synchronization) error {
	machine := &machinev1alpha2.Machine{}
	if err := s.client.Get(ctx,
		types.NamespacedName{Name: e.SourceName, Namespace: e.SourceNamespace}, machine); err != nil {
		return err
	}
	metalAssignment := &machinev1alpha2.MachineAssignment{}
	if err := s.client.Get(ctx,
		types.NamespacedName{Name: e.TargetName, Namespace: e.TargetNamespace}, metalAssignment); err != nil {
		return err
	}

	e.SourceStatus = machine.Status.Reservation.Status
	e.TargetStatus = metalAssignment.Status.State

	if !e.IsSyncNeeded() {
		return nil
	}
	metalAssignment.Status.State = e.SourceStatus

	return s.client.Status().Update(ctx, metalAssignment)
}
