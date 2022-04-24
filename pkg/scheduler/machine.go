package scheduler

import (
	"context"

	"github.com/go-logr/logr"
	machinev1alpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	requestv1alpha1 "github.com/onmetal/metal-api/apis/request/v1alpha1"
	machineclient "github.com/onmetal/metal-api/pkg/machine"
	"github.com/onmetal/metal-api/pkg/provider"
	"github.com/onmetal/metal-api/pkg/reserve"
	"k8s.io/client-go/tools/record"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type machine struct {
	ctrlclient.Client
	machineclient.Machiner

	ctx      context.Context
	log      logr.Logger
	recorder record.EventRecorder
}

func New(ctx context.Context, c ctrlclient.Client, l logr.Logger, recorder record.EventRecorder) Scheduler {
	mClient := machineclient.New(ctx, c, l, recorder)
	return &machine{
		Client:   c,
		recorder: recorder,
		ctx:      ctx,
		log:      l,
		Machiner: mClient,
	}
}

func (m *machine) Schedule(metalRequest *requestv1alpha1.Request) error {
	machineForRequest, err := m.Machiner.FindVacantMachine(metalRequest)
	if err != nil {
		return err
	}

	r := reserve.NewMachineReserver(m.ctx, m.Client, m.log, m.recorder, machineForRequest)
	if err := r.Reserve(metalRequest.Name); err != nil {
		return err
	}

	metalRequest.Status.State = machinev1alpha2.RequestStateReserved
	metalRequest.Status.Reference = getObjectReference(machineForRequest)

	return m.Status().Update(m.ctx, metalRequest)
}

func (m *machine) DeleteScheduling(metalRequest *requestv1alpha1.Request) error {
	if metalRequest.Status.Reference == nil {
		m.log.Info("machine reference not found", "request", metalRequest.Name)
		return nil
	}

	machineObj := &machinev1alpha2.Machine{}
	machineName, machineNamespase := metalRequest.Status.Reference.Name, metalRequest.Status.Reference.Namespace
	if err := provider.GetObject(m.ctx, machineName, machineNamespase, m.Client, machineObj); err != nil {
		return err
	}

	r := reserve.NewMachineReserver(m.ctx, m.Client, m.log, m.recorder, machineObj)
	return r.DeleteReservation()
}

func getObjectReference(m *machinev1alpha2.Machine) *requestv1alpha1.ResourceReference {
	return &requestv1alpha1.ResourceReference{
		APIVersion: m.APIVersion,
		Kind:       m.Kind,
		Name:       m.Name,
		Namespace:  m.Namespace,
	}
}
