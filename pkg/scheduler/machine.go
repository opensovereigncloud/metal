package scheduler

import (
	"context"

	"github.com/go-logr/logr"
	machinev1alpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	requestv1alpha1 "github.com/onmetal/metal-api/apis/request/v1alpha1"
	machineclient "github.com/onmetal/metal-api/pkg/machine"
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

func NewMachine(ctx context.Context, c ctrlclient.Client, l logr.Logger, recorder record.EventRecorder) Scheduler {
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
	machineForRequest, err := m.FindVacantMachine(metalRequest)
	if err != nil {
		return err
	}

	if err := m.Reserve(machineForRequest, metalRequest); err != nil {
		return err
	}

	metalRequest.Status.State = requestv1alpha1.Reserved
	metalRequest.Status.Reference = getObjectReference(machineForRequest)

	return m.Status().Update(m.ctx, metalRequest)
}

func (m *machine) DeleteScheduling(metalRequest *requestv1alpha1.Request) error {
	machineObj, err := m.GetMachine(metalRequest.Status.Reference.Name, metalRequest.Status.Reference.Namespace)
	if err != nil {
		return err
	}

	return m.CheckOut(machineObj)
}

func getObjectReference(m *machinev1alpha2.Machine) *requestv1alpha1.ResourceReference {
	return &requestv1alpha1.ResourceReference{
		APIVersion: m.APIVersion,
		Kind:       m.Kind,
		Name:       m.Name,
		Namespace:  m.Namespace,
	}
}
