package scheduler

import (
	"context"
	"github.com/onmetal/metal-api/apis/machine/v1alpha2"
	oobv1 "github.com/onmetal/oob-controller/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *Reconciler) getOOBMachine(ctx context.Context, machine *v1alpha2.Machine) (*oobv1.Machine, error) {
	oobMachine := &oobv1.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      machine.Status.OOB.Reference.Name,
			Namespace: machine.Status.OOB.Reference.Namespace,
		},
	}

	err := r.Client.Get(ctx, client.ObjectKeyFromObject(oobMachine), oobMachine)
	if err != nil {
		return nil, err
	}

	return oobMachine, nil
}
