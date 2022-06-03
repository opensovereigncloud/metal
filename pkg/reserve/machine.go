package reserve

import (
	"context"

	"github.com/go-logr/logr"
	machinev1alpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	metalerr "github.com/onmetal/metal-api/pkg/errors"
	oobonmetal "github.com/onmetal/oob-controller/api/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/record"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

const fieldManager = "reserver"

type Machine struct {
	ctrlclient.Client
	*machinev1alpha2.Machine

	ctx      context.Context
	log      logr.Logger
	recorder record.EventRecorder
}

func NewMachineReserver(ctx context.Context, c ctrlclient.Client,
	l logr.Logger, recorder record.EventRecorder,
	kind *machinev1alpha2.Machine) *Machine {
	return &Machine{
		Client:   c,
		ctx:      ctx,
		log:      l,
		Machine:  kind,
		recorder: recorder,
	}
}

func (m *Machine) Reserve(requestName string) error {
	m.Labels[machinev1alpha2.LeasedLabel] = "true"
	m.Labels[machinev1alpha2.MetalRequestLabel] = requestName

	if err := m.Client.Update(m.ctx, m.Machine); err != nil {
		return err
	}

	m.Machine.Status.Reservation.RequestState = machinev1alpha2.RequestStatePending
	m.Machine.Status.Reservation.Reference = prepareRefenceSpec(requestName, m.Namespace)

	return m.Client.Status().Patch(m.ctx, m.Machine, ctrlclient.Merge, &ctrlclient.PatchOptions{
		FieldManager: fieldManager,
	})
}

func (m *Machine) DeleteReservation() error {
	delete(m.Labels, machinev1alpha2.LeasedLabel)
	delete(m.Labels, machinev1alpha2.MetalRequestLabel)

	if err := m.Client.Update(m.ctx, m.Machine); err != nil {
		return err
	}

	m.Machine.Status.Reservation.RequestState = machinev1alpha2.RequestStateAvailable
	m.Machine.Status.Reservation.Reference = &machinev1alpha2.ResourceReference{}

	return m.Client.Status().Update(m.ctx, m.Machine)
}

func (m *Machine) CheckIn() error {
	return nil
	// return m.changeServerPowerState(m.Name, m.Namespace, machinev1alpha2.MachinePowerStateON)
}

func (m *Machine) CheckOut() error {
	return nil
	// return m.changeServerPowerState(m.Name, m.Namespace, machinev1alpha2.MachinePowerStateOFF)
}

func (m *Machine) changeServerPowerState(name, namespace, powerState string) error {
	oobObj, err := m.getOOBMachineByUUIDLabel(name, namespace)
	if err != nil {
		return err
	}

	if powerState == oobObj.Spec.PowerState && powerState == machinev1alpha2.MachinePowerStateON {
		powerState = machinev1alpha2.MachinePowerStateReset
	}

	oobObj.Spec.PowerState = powerState

	m.log.Info("oob state changed", "uuid", "namespace", oobObj.Name, oobObj.Namespace)
	return m.Client.Patch(m.ctx, oobObj, ctrlclient.Merge, &ctrlclient.PatchOptions{
		FieldManager: "machine-controller",
	})
}

func (m *Machine) getOOBMachineByUUIDLabel(name, namespace string) (*oobonmetal.Machine, error) {
	oobs := &oobonmetal.MachineList{}
	listOptions := &ctrlclient.ListOptions{
		Namespace: namespace,
		LabelSelector: ctrlclient.MatchingLabelsSelector{
			Selector: labels.SelectorFromSet(map[string]string{
				machinev1alpha2.UUIDLabel: name,
			})}}
	if err := m.Client.List(m.ctx, oobs, listOptions); err != nil {
		return nil, err
	}

	if len(oobs.Items) == 0 {
		return nil, metalerr.ObjectNotFound(name, "oob")
	}
	return &oobs.Items[0], nil
}

func prepareRefenceSpec(requestName, namespace string) *machinev1alpha2.ResourceReference {
	return &machinev1alpha2.ResourceReference{
		Name: requestName, Namespace: namespace,
	}
}
