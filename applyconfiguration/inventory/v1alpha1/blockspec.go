/*
Copyright (c) 2023 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved

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

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

// BlockSpecApplyConfiguration represents an declarative configuration of the BlockSpec type for use
// with apply.
type BlockSpecApplyConfiguration struct {
	Name           *string                               `json:"name,omitempty"`
	Type           *string                               `json:"type,omitempty"`
	Rotational     *bool                                 `json:"rotational,omitempty"`
	Bus            *string                               `json:"system,omitempty"`
	Model          *string                               `json:"model,omitempty"`
	Size           *uint64                               `json:"size,omitempty"`
	PartitionTable *PartitionTableSpecApplyConfiguration `json:"partitionTable,omitempty"`
}

// BlockSpecApplyConfiguration constructs an declarative configuration of the BlockSpec type for use with
// apply.
func BlockSpec() *BlockSpecApplyConfiguration {
	return &BlockSpecApplyConfiguration{}
}

// WithName sets the Name field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Name field is set to the value of the last call.
func (b *BlockSpecApplyConfiguration) WithName(value string) *BlockSpecApplyConfiguration {
	b.Name = &value
	return b
}

// WithType sets the Type field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Type field is set to the value of the last call.
func (b *BlockSpecApplyConfiguration) WithType(value string) *BlockSpecApplyConfiguration {
	b.Type = &value
	return b
}

// WithRotational sets the Rotational field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Rotational field is set to the value of the last call.
func (b *BlockSpecApplyConfiguration) WithRotational(value bool) *BlockSpecApplyConfiguration {
	b.Rotational = &value
	return b
}

// WithBus sets the Bus field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Bus field is set to the value of the last call.
func (b *BlockSpecApplyConfiguration) WithBus(value string) *BlockSpecApplyConfiguration {
	b.Bus = &value
	return b
}

// WithModel sets the Model field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Model field is set to the value of the last call.
func (b *BlockSpecApplyConfiguration) WithModel(value string) *BlockSpecApplyConfiguration {
	b.Model = &value
	return b
}

// WithSize sets the Size field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Size field is set to the value of the last call.
func (b *BlockSpecApplyConfiguration) WithSize(value uint64) *BlockSpecApplyConfiguration {
	b.Size = &value
	return b
}

// WithPartitionTable sets the PartitionTable field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the PartitionTable field is set to the value of the last call.
func (b *BlockSpecApplyConfiguration) WithPartitionTable(value *PartitionTableSpecApplyConfiguration) *BlockSpecApplyConfiguration {
	b.PartitionTable = value
	return b
}
