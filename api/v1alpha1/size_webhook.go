/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/util/jsonpath"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var sizelog = logf.Log.WithName("size-resource")

func (r *Size) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:path=/validate-machine-onmetal-de-v1alpha1-size,mutating=false,failurePolicy=fail,sideEffects=None,groups=machine.onmetal.de,resources=sizes,verbs=create;update,versions=v1alpha1,name=vsize.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &Size{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Size) ValidateCreate() error {
	sizelog.Info("validate create", "name", r.Name)
	return r.validate()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Size) ValidateUpdate(old runtime.Object) error {
	sizelog.Info("validate update", "name", r.Name)
	return r.validate()
}

var CDummyInventorySpec = getDummyInventorySpec()

func (r *Size) validate() error {
	ops := make(map[string]int, 0)
	errs := make([]string, 0)

	for _, c := range r.Spec.Constraints {
		op, ok := ops[c.Path]
		if !ok {
			op = 0
		}
		op++
		ops[c.Path] = op

		jp := jsonpath.New(c.Path)
		jp.AllowMissingKeys(false)
		if err := jp.Parse(normalizeJSONPath(c.Path)); err != nil {
			errs = append(errs, errors.Wrap(err, "unable to parse JSONPath").Error())
		}

		if _, err := jp.FindResults(CDummyInventorySpec); err != nil {
			errs = append(errs, errors.Wrap(err, "unable to find results with path").Error())
		}

		if op == 2 {
			err := errors.Errorf("multiple constraints found for field %s", c.Path)
			errs = append(errs, err.Error())
		}

		if c.empty() {
			err := errors.Errorf("constraint for %s does not contains conditions", c.Path)
			errs = append(errs, err.Error())
		}

		if c.hasAggregateAndLiterals() {
			err := errors.New("aggregates can be validated only against numeric values")
			errs = append(errs, err.Error())
		}

		if c.eqAndNeq() {
			err := errors.Errorf("constraint for %s contains both eq and neq conditions", c.Path)
			errs = append(errs, err.Error())
		}

		if c.inclusiveAndExclusive() {
			err := errors.Errorf("constraint for %s contains both gt and gte or lt and lte conditions", c.Path)
			errs = append(errs, err.Error())
		}

		if c.borderAndEq() {
			err := errors.Errorf("constraint for %s contains both gt/gte/lt/lte and eq/neq conditions", c.Path)
			errs = append(errs, err.Error())
		}

		if c.wrongInterval() {
			err := errors.Errorf("constraint for %s lower border is greater than upper border", c.Path)
			errs = append(errs, err.Error())
		}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}

	return nil
}

func (r *ConstraintSpec) hasAggregateAndLiterals() bool {
	return r.Aggregate != "" &&
		(r.Equal != nil && r.Equal.Literal != nil ||
			r.NotEqual != nil && r.NotEqual.Literal != nil)
}

func (r *ConstraintSpec) empty() bool {
	return r.Equal == nil &&
		r.NotEqual == nil &&
		r.GreaterThan == nil &&
		r.GreaterThanOrEqual == nil &&
		r.LessThan == nil &&
		r.LessThanOrEqual == nil
}

func (r *ConstraintSpec) eqAndNeq() bool {
	return r.Equal != nil &&
		r.NotEqual != nil
}

func (r *ConstraintSpec) inclusiveAndExclusive() bool {
	return r.GreaterThan != nil && r.GreaterThanOrEqual != nil ||
		r.LessThan != nil && r.LessThanOrEqual != nil
}

func (r *ConstraintSpec) borderAndEq() bool {
	return (r.GreaterThan != nil || r.GreaterThanOrEqual != nil || r.LessThan != nil || r.LessThanOrEqual != nil) &&
		(r.Equal != nil || r.NotEqual != nil)
}

func (r *ConstraintSpec) wrongInterval() bool {
	var upper *resource.Quantity
	var lower *resource.Quantity

	if r.LessThanOrEqual != nil {
		upper = r.LessThanOrEqual
	}
	if r.LessThan != nil {
		upper = r.LessThan
	}
	if r.GreaterThanOrEqual != nil {
		lower = r.GreaterThanOrEqual
	}
	if r.GreaterThan != nil {
		lower = r.GreaterThan
	}

	if upper == nil || lower == nil {
		return false
	}

	if lower.Cmp(*upper) < 0 {
		return false
	}

	return true
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Size) ValidateDelete() error {
	sizelog.Info("validate delete", "name", r.Name)
	return nil
}

// getDummyInventorySpec fills structure with dummy data and used to validate whether path points to existing field
func getDummyInventorySpec() *InventorySpec {
	return &InventorySpec{
		System: &SystemSpec{
			ID:           "",
			Manufacturer: "",
			ProductSKU:   "",
			SerialNumber: "",
		},
		IPMIs: []IPMISpec{
			{
				IPAddress:  "",
				MACAddress: "",
			},
		},
		Blocks: &BlockTotalSpec{
			Count:    0,
			Capacity: 0,
			Blocks: []BlockSpec{
				{
					Name:       "",
					Type:       "",
					Rotational: false,
					Bus:        "",
					Model:      "",
					Size:       0,
					PartitionTable: &PartitionTableSpec{
						Type: "",
						Partitions: []PartitionSpec{
							{
								ID:   "",
								Name: "",
								Size: 0,
							},
						},
					},
				},
			},
		},
		Memory: &MemorySpec{
			Total: 0,
		},
		CPUs: &CPUTotalSpec{
			Sockets: 0,
			Cores:   0,
			Threads: 0,
			CPUs: []CPUSpec{
				{
					PhysicalID: 0,
					LogicalIDs: []uint64{
						0,
					},
					Cores:        0,
					Siblings:     0,
					VendorID:     "",
					Family:       "",
					Model:        "",
					ModelName:    "",
					Stepping:     "",
					Microcode:    "",
					MHz:          *resource.NewScaledQuantity(0, 0),
					CacheSize:    "",
					FPU:          false,
					FPUException: false,
					CPUIDLevel:   0,
					WP:           false,
					Flags: []string{
						"",
					},
					VMXFlags: []string{
						"",
					},
					Bugs: []string{
						"",
					},
					BogoMIPS:        *resource.NewScaledQuantity(0, 0),
					CLFlushSize:     0,
					CacheAlignment:  0,
					AddressSizes:    "",
					PowerManagement: "",
				},
			},
		},
		NICs: &NICTotalSpec{
			Count: 0,
			NICs: []NICSpec{
				{
					Name:       "",
					PCIAddress: "",
					MACAddress: "",
					MTU:        0,
					Speed:      0,
					LLDPs: []LLDPSpec{
						{
							ChassisID:         "",
							SystemName:        "",
							SystemDescription: "",
							PortID:            "",
							PortDescription:   "",
						},
					},
					NDPs: []NDPSpec{
						{
							IPAddress:  "",
							MACAddress: "",
							State:      "",
						},
					},
				},
			},
		},
		Virt: &VirtSpec{
			VMType: "",
		},
		Host: &HostSpec{
			Type: "",
			Name: "",
		},
		Distro: &DistroSpec{
			BuildVersion:  "",
			DebianVersion: "",
			KernelVersion: "",
			AsicType:      "",
			CommitId:      "",
			BuildDate:     "",
			BuildNumber:   0,
			BuildBy:       "",
		},
	}
}

func normalizeJSONPath(jp string) string {
	if strings.HasPrefix(jp, "{.") {
		return jp
	}
	if strings.HasPrefix(jp, ".") {
		return fmt.Sprintf("{%s}", jp)
	}
	return fmt.Sprintf("{.%s}", jp)
}
