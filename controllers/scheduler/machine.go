package scheduler

import (
	"context"
	"fmt"
	"github.com/onmetal/metal-api/apis/machine/v1alpha2"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *Reconciler) findAvailableMachine(ctx context.Context, machineAssignment *v1alpha2.MachineAssignment) (*v1alpha2.Machine, error) {
	instanceList := &v1alpha2.MachineList{}
	instanceSelector := getLabelSelectorForAvailableMachine(machineAssignment.Spec.MachineSize)
	listOpts := &client.ListOptions{
		LabelSelector: client.MatchingLabelsSelector{Selector: labels.SelectorFromSet(instanceSelector)},
	}

	err := r.Client.List(ctx, instanceList, listOpts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list available machines")
	}

	for _, machine := range instanceList.Items {
		if machine.Status.Reservation.Status == ReservationStatusAvailable {
			return &machine, nil
		}
	}

	return nil, err
}

func (r *Reconciler) getMachine(ctx context.Context, reference *v1alpha2.ResourceReference) (*v1alpha2.Machine, error) {
	machine := &v1alpha2.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      reference.Name,
			Namespace: reference.Namespace,
		},
	}
	err := r.Client.Get(ctx, client.ObjectKeyFromObject(machine), machine)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to get machine '%s' from namespace '%s'", reference.Name, reference.Namespace))
	}

	return machine, nil
}
