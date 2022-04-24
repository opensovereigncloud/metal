package reserve

import (
	"context"

	"github.com/go-logr/logr"
	machinev1alpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	metalerr "github.com/onmetal/metal-api/pkg/errors"
	"github.com/onmetal/metal-api/pkg/machine"
	oobonmetal "github.com/onmetal/oob-controller/api/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/record"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type Machine struct {
	ctrlclient.Client

	ctx      context.Context
	log      logr.Logger
	machine  *machinev1alpha2.Machine
	recorder record.EventRecorder
}

func NewMachineReserver(ctx context.Context, c ctrlclient.Client,
	l logr.Logger, recorder record.EventRecorder, kind *machinev1alpha2.Machine) Reserver {

	return &Machine{
		Client:   c,
		ctx:      ctx,
		log:      l,
		machine:  kind,
		recorder: recorder,
	}
}

func (m *Machine) Reserve(requestName string) error {
	m.machine.Labels[machinev1alpha2.LeasedLabel] = "true"
	m.machine.Labels[machinev1alpha2.MetalRequestLabel] = requestName

	mm := machine.New(m.ctx, m.Client, m.log, m.recorder)
	return mm.UpdateSpec(m.machine)
}

func (m *Machine) DeleteReservation() error {
	delete(m.machine.Labels, machinev1alpha2.LeasedLabel)
	delete(m.machine.Labels, machinev1alpha2.MetalRequestLabel)

	mm := machine.New(m.ctx, m.Client, m.log, m.recorder)
	return mm.UpdateSpec(m.machine)
}

func (m *Machine) CheckIn() error {
	return m.changeServerPowerState(m.machine.Name, m.machine.Namespace, machinev1alpha2.MachinePowerStateON)
}

func (m *Machine) CheckOut() error {
	return m.changeServerPowerState(m.machine.Name, m.machine.Namespace, machinev1alpha2.MachinePowerStateOFF)
}

func (m *Machine) changeServerPowerState(name, namespace, powerState string) error {
	oobObj, err := m.getOOBMachineByUUIDLabel(name, namespace)
	if err != nil {
		return err
	}

	if powerState == oobObj.Spec.PowerState {
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
