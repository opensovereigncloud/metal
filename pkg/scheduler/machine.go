package scheduler

import (
	"context"

	"github.com/go-logr/logr"
	machinev1alpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	machineclient "github.com/onmetal/metal-api/pkg/machine"
	"github.com/onmetal/metal-api/pkg/provider"
	"github.com/onmetal/metal-api/pkg/reserve"
	"k8s.io/client-go/tools/record"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type Machine struct {
	ctrlclient.Client
	machineclient.Machiner

	ctx      context.Context
	log      logr.Logger
	recorder record.EventRecorder
}

func NewMachine(ctx context.Context, c ctrlclient.Client, l logr.Logger, recorder record.EventRecorder) *Machine {
	mClient := machineclient.New(ctx, c, l, recorder)
	return &Machine{
		Client:   c,
		recorder: recorder,
		ctx:      ctx,
		log:      l,
		Machiner: mClient,
	}
}

func (m *Machine) Schedule(metalRequest *machinev1alpha2.MachineAssignment) error {
	machineForRequest, err := m.Machiner.FindVacantMachine(metalRequest)
	if err != nil {
		return err
	}

	var reserver reserve.Reserver //nolint:gosimple
	reserver = reserve.NewMachineReserver(m.ctx, m.Client, m.log, m.recorder, machineForRequest)
	if err := reserver.Reserve(metalRequest.Name); err != nil {
		return err
	}

	metalRequest.Status.Reference = getObjectReference(machineForRequest)
	metalRequest.Status.State = machinev1alpha2.RequestStateReserved

	return m.Status().Patch(m.ctx, metalRequest, ctrlclient.Merge, &ctrlclient.PatchOptions{
		FieldManager: "scheduler",
	})
}

func (m *Machine) DeleteScheduling(metalRequest *machinev1alpha2.MachineAssignment) error {
	if metalRequest.Status.Reference == nil {
		m.log.Info("machine reference not found", "request", metalRequest.Name)
		return nil
	}

	machineObj := &machinev1alpha2.Machine{}
	machineName, machineNamespase := metalRequest.Status.Reference.Name, metalRequest.Status.Reference.Namespace
	if err := provider.GetObject(m.ctx, machineName, machineNamespase, m.Client, machineObj); err != nil {
		return err
	}

	var reserver reserve.Reserver //nolint:gosimple
	reserver = reserve.NewMachineReserver(m.ctx, m.Client, m.log, m.recorder, machineObj)
	return reserver.DeleteReservation()
}

func getObjectReference(m *machinev1alpha2.Machine) *machinev1alpha2.ResourceReference {
	return &machinev1alpha2.ResourceReference{
		APIVersion: m.APIVersion,
		Kind:       m.Kind,
		Name:       m.Name,
		Namespace:  m.Namespace,
	}
}
