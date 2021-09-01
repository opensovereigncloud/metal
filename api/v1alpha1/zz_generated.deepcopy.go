//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Aggregate) DeepCopyInto(out *Aggregate) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Aggregate.
func (in *Aggregate) DeepCopy() *Aggregate {
	if in == nil {
		return nil
	}
	out := new(Aggregate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Aggregate) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AggregateItem) DeepCopyInto(out *AggregateItem) {
	*out = *in
	out.SourcePath = in.SourcePath
	out.TargetPath = in.TargetPath
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AggregateItem.
func (in *AggregateItem) DeepCopy() *AggregateItem {
	if in == nil {
		return nil
	}
	out := new(AggregateItem)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AggregateList) DeepCopyInto(out *AggregateList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Aggregate, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AggregateList.
func (in *AggregateList) DeepCopy() *AggregateList {
	if in == nil {
		return nil
	}
	out := new(AggregateList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AggregateList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AggregateSpec) DeepCopyInto(out *AggregateSpec) {
	*out = *in
	if in.Aggregates != nil {
		in, out := &in.Aggregates, &out.Aggregates
		*out = make([]AggregateItem, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AggregateSpec.
func (in *AggregateSpec) DeepCopy() *AggregateSpec {
	if in == nil {
		return nil
	}
	out := new(AggregateSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AggregateStatus) DeepCopyInto(out *AggregateStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AggregateStatus.
func (in *AggregateStatus) DeepCopy() *AggregateStatus {
	if in == nil {
		return nil
	}
	out := new(AggregateStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BlockSpec) DeepCopyInto(out *BlockSpec) {
	*out = *in
	if in.PartitionTable != nil {
		in, out := &in.PartitionTable, &out.PartitionTable
		*out = new(PartitionTableSpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BlockSpec.
func (in *BlockSpec) DeepCopy() *BlockSpec {
	if in == nil {
		return nil
	}
	out := new(BlockSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BlockTotalSpec) DeepCopyInto(out *BlockTotalSpec) {
	*out = *in
	if in.Blocks != nil {
		in, out := &in.Blocks, &out.Blocks
		*out = make([]BlockSpec, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BlockTotalSpec.
func (in *BlockTotalSpec) DeepCopy() *BlockTotalSpec {
	if in == nil {
		return nil
	}
	out := new(BlockTotalSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CPUSpec) DeepCopyInto(out *CPUSpec) {
	*out = *in
	if in.LogicalIDs != nil {
		in, out := &in.LogicalIDs, &out.LogicalIDs
		*out = make([]uint64, len(*in))
		copy(*out, *in)
	}
	out.MHz = in.MHz.DeepCopy()
	if in.Flags != nil {
		in, out := &in.Flags, &out.Flags
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.VMXFlags != nil {
		in, out := &in.VMXFlags, &out.VMXFlags
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Bugs != nil {
		in, out := &in.Bugs, &out.Bugs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	out.BogoMIPS = in.BogoMIPS.DeepCopy()
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CPUSpec.
func (in *CPUSpec) DeepCopy() *CPUSpec {
	if in == nil {
		return nil
	}
	out := new(CPUSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CPUTotalSpec) DeepCopyInto(out *CPUTotalSpec) {
	*out = *in
	if in.CPUs != nil {
		in, out := &in.CPUs, &out.CPUs
		*out = make([]CPUSpec, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CPUTotalSpec.
func (in *CPUTotalSpec) DeepCopy() *CPUTotalSpec {
	if in == nil {
		return nil
	}
	out := new(CPUTotalSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConstraintSpec) DeepCopyInto(out *ConstraintSpec) {
	*out = *in
	if in.Equal != nil {
		in, out := &in.Equal, &out.Equal
		*out = new(ConstraintValSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.NotEqual != nil {
		in, out := &in.NotEqual, &out.NotEqual
		*out = new(ConstraintValSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.LessThan != nil {
		in, out := &in.LessThan, &out.LessThan
		x := (*in).DeepCopy()
		*out = &x
	}
	if in.LessThanOrEqual != nil {
		in, out := &in.LessThanOrEqual, &out.LessThanOrEqual
		x := (*in).DeepCopy()
		*out = &x
	}
	if in.GreaterThan != nil {
		in, out := &in.GreaterThan, &out.GreaterThan
		x := (*in).DeepCopy()
		*out = &x
	}
	if in.GreaterThanOrEqual != nil {
		in, out := &in.GreaterThanOrEqual, &out.GreaterThanOrEqual
		x := (*in).DeepCopy()
		*out = &x
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConstraintSpec.
func (in *ConstraintSpec) DeepCopy() *ConstraintSpec {
	if in == nil {
		return nil
	}
	out := new(ConstraintSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConstraintValSpec) DeepCopyInto(out *ConstraintValSpec) {
	*out = *in
	if in.Literal != nil {
		in, out := &in.Literal, &out.Literal
		*out = new(string)
		**out = **in
	}
	if in.Numeric != nil {
		in, out := &in.Numeric, &out.Numeric
		x := (*in).DeepCopy()
		*out = &x
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConstraintValSpec.
func (in *ConstraintValSpec) DeepCopy() *ConstraintValSpec {
	if in == nil {
		return nil
	}
	out := new(ConstraintValSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DistroSpec) DeepCopyInto(out *DistroSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DistroSpec.
func (in *DistroSpec) DeepCopy() *DistroSpec {
	if in == nil {
		return nil
	}
	out := new(DistroSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HostSpec) DeepCopyInto(out *HostSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HostSpec.
func (in *HostSpec) DeepCopy() *HostSpec {
	if in == nil {
		return nil
	}
	out := new(HostSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IPMISpec) DeepCopyInto(out *IPMISpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IPMISpec.
func (in *IPMISpec) DeepCopy() *IPMISpec {
	if in == nil {
		return nil
	}
	out := new(IPMISpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Inventory) DeepCopyInto(out *Inventory) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Inventory.
func (in *Inventory) DeepCopy() *Inventory {
	if in == nil {
		return nil
	}
	out := new(Inventory)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Inventory) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InventoryList) DeepCopyInto(out *InventoryList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Inventory, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InventoryList.
func (in *InventoryList) DeepCopy() *InventoryList {
	if in == nil {
		return nil
	}
	out := new(InventoryList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *InventoryList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InventorySpec) DeepCopyInto(out *InventorySpec) {
	*out = *in
	if in.System != nil {
		in, out := &in.System, &out.System
		*out = new(SystemSpec)
		**out = **in
	}
	if in.IPMIs != nil {
		in, out := &in.IPMIs, &out.IPMIs
		*out = make([]IPMISpec, len(*in))
		copy(*out, *in)
	}
	if in.Blocks != nil {
		in, out := &in.Blocks, &out.Blocks
		*out = new(BlockTotalSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Memory != nil {
		in, out := &in.Memory, &out.Memory
		*out = new(MemorySpec)
		**out = **in
	}
	if in.CPUs != nil {
		in, out := &in.CPUs, &out.CPUs
		*out = new(CPUTotalSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.NICs != nil {
		in, out := &in.NICs, &out.NICs
		*out = new(NICTotalSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Virt != nil {
		in, out := &in.Virt, &out.Virt
		*out = new(VirtSpec)
		**out = **in
	}
	if in.Host != nil {
		in, out := &in.Host, &out.Host
		*out = new(HostSpec)
		**out = **in
	}
	if in.Distro != nil {
		in, out := &in.Distro, &out.Distro
		*out = new(DistroSpec)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InventorySpec.
func (in *InventorySpec) DeepCopy() *InventorySpec {
	if in == nil {
		return nil
	}
	out := new(InventorySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InventoryStatus) DeepCopyInto(out *InventoryStatus) {
	*out = *in
	in.Computed.DeepCopyInto(&out.Computed)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InventoryStatus.
func (in *InventoryStatus) DeepCopy() *InventoryStatus {
	if in == nil {
		return nil
	}
	out := new(InventoryStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JSONPath) DeepCopyInto(out *JSONPath) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JSONPath.
func (in *JSONPath) DeepCopy() *JSONPath {
	if in == nil {
		return nil
	}
	out := new(JSONPath)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LLDPSpec) DeepCopyInto(out *LLDPSpec) {
	*out = *in
	if in.Capabilities != nil {
		in, out := &in.Capabilities, &out.Capabilities
		*out = make([]LLDPCapabilities, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LLDPSpec.
func (in *LLDPSpec) DeepCopy() *LLDPSpec {
	if in == nil {
		return nil
	}
	out := new(LLDPSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MemorySpec) DeepCopyInto(out *MemorySpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MemorySpec.
func (in *MemorySpec) DeepCopy() *MemorySpec {
	if in == nil {
		return nil
	}
	out := new(MemorySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NDPSpec) DeepCopyInto(out *NDPSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NDPSpec.
func (in *NDPSpec) DeepCopy() *NDPSpec {
	if in == nil {
		return nil
	}
	out := new(NDPSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NICSpec) DeepCopyInto(out *NICSpec) {
	*out = *in
	if in.LLDPs != nil {
		in, out := &in.LLDPs, &out.LLDPs
		*out = make([]LLDPSpec, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.NDPs != nil {
		in, out := &in.NDPs, &out.NDPs
		*out = make([]NDPSpec, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NICSpec.
func (in *NICSpec) DeepCopy() *NICSpec {
	if in == nil {
		return nil
	}
	out := new(NICSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NICTotalSpec) DeepCopyInto(out *NICTotalSpec) {
	*out = *in
	if in.NICs != nil {
		in, out := &in.NICs, &out.NICs
		*out = make([]NICSpec, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NICTotalSpec.
func (in *NICTotalSpec) DeepCopy() *NICTotalSpec {
	if in == nil {
		return nil
	}
	out := new(NICTotalSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PartitionSpec) DeepCopyInto(out *PartitionSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PartitionSpec.
func (in *PartitionSpec) DeepCopy() *PartitionSpec {
	if in == nil {
		return nil
	}
	out := new(PartitionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PartitionTableSpec) DeepCopyInto(out *PartitionTableSpec) {
	*out = *in
	if in.Partitions != nil {
		in, out := &in.Partitions, &out.Partitions
		*out = make([]PartitionSpec, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PartitionTableSpec.
func (in *PartitionTableSpec) DeepCopy() *PartitionTableSpec {
	if in == nil {
		return nil
	}
	out := new(PartitionTableSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Size) DeepCopyInto(out *Size) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Size.
func (in *Size) DeepCopy() *Size {
	if in == nil {
		return nil
	}
	out := new(Size)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Size) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SizeList) DeepCopyInto(out *SizeList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Size, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SizeList.
func (in *SizeList) DeepCopy() *SizeList {
	if in == nil {
		return nil
	}
	out := new(SizeList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SizeList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SizeSpec) DeepCopyInto(out *SizeSpec) {
	*out = *in
	if in.Constraints != nil {
		in, out := &in.Constraints, &out.Constraints
		*out = make([]ConstraintSpec, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SizeSpec.
func (in *SizeSpec) DeepCopy() *SizeSpec {
	if in == nil {
		return nil
	}
	out := new(SizeSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SizeStatus) DeepCopyInto(out *SizeStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SizeStatus.
func (in *SizeStatus) DeepCopy() *SizeStatus {
	if in == nil {
		return nil
	}
	out := new(SizeStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SystemSpec) DeepCopyInto(out *SystemSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SystemSpec.
func (in *SystemSpec) DeepCopy() *SystemSpec {
	if in == nil {
		return nil
	}
	out := new(SystemSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtSpec) DeepCopyInto(out *VirtSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtSpec.
func (in *VirtSpec) DeepCopy() *VirtSpec {
	if in == nil {
		return nil
	}
	out := new(VirtSpec)
	in.DeepCopyInto(out)
	return out
}
