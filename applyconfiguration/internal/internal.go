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

package internal

import (
	"fmt"
	"sync"

	typed "sigs.k8s.io/structured-merge-diff/v4/typed"
)

func Parser() *typed.Parser {
	parserOnce.Do(func() {
		var err error
		parser, err = typed.NewParser(schemaYAML)
		if err != nil {
			panic(fmt.Sprintf("Failed to parse schema: %v", err))
		}
	})
	return parser
}

var parserOnce sync.Once
var parser *typed.Parser
var schemaYAML = typed.YAMLObject(`types:
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.AdditionalIPSpec
  map:
    fields:
    - name: address
      type:
        scalar: string
    - name: parentSubnet
      type:
        namedType: io.k8s.apimachinery.pkg.apis.meta.v1.LabelSelector
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.AddressFamiliesMap
  map:
    fields:
    - name: ipv4
      type:
        scalar: boolean
    - name: ipv6
      type:
        scalar: boolean
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.Addresses
  map:
    fields:
    - name: ipv4
      type:
        list:
          elementType:
            namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.IPAddrSpec
          elementRelationship: atomic
    - name: ipv6
      type:
        list:
          elementType:
            namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.IPAddrSpec
          elementRelationship: atomic
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.Aggregate
  map:
    fields:
    - name: apiVersion
      type:
        scalar: string
    - name: kind
      type:
        scalar: string
    - name: metadata
      type:
        namedType: io.k8s.apimachinery.pkg.apis.meta.v1.ObjectMeta
      default: {}
    - name: spec
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.AggregateSpec
      default: {}
    - name: status
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.AggregateStatus
      default: {}
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.AggregateItem
  map:
    fields:
    - name: aggregate
      type:
        scalar: string
    - name: sourcePath
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.JSONPath
      default: {}
    - name: targetPath
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.JSONPath
      default: {}
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.AggregateSpec
  map:
    fields:
    - name: aggregates
      type:
        list:
          elementType:
            namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.AggregateItem
          elementRelationship: atomic
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.AggregateStatus
  map:
    elementType:
      scalar: untyped
      list:
        elementType:
          namedType: __untyped_atomic_
        elementRelationship: atomic
      map:
        elementType:
          namedType: __untyped_deduced_
        elementRelationship: separable
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.AggregationResults
  map:
    elementType:
      scalar: untyped
      list:
        elementType:
          namedType: __untyped_atomic_
        elementRelationship: atomic
      map:
        elementType:
          namedType: __untyped_deduced_
        elementRelationship: separable
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.Benchmark
  map:
    fields:
    - name: apiVersion
      type:
        scalar: string
    - name: kind
      type:
        scalar: string
    - name: metadata
      type:
        namedType: io.k8s.apimachinery.pkg.apis.meta.v1.ObjectMeta
      default: {}
    - name: spec
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.BenchmarkSpec
      default: {}
    - name: status
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.BenchmarkStatus
      default: {}
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.BenchmarkDeviation
  map:
    fields:
    - name: name
      type:
        scalar: string
      default: ""
    - name: value
      type:
        scalar: string
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.BenchmarkResult
  map:
    fields:
    - name: name
      type:
        scalar: string
      default: ""
    - name: value
      type:
        scalar: numeric
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.BenchmarkSpec
  map:
    fields:
    - name: benchmarks
      type:
        map:
          elementType:
            list:
              elementType:
                namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.BenchmarkResult
              elementRelationship: atomic
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.BenchmarkStatus
  map:
    fields:
    - name: machine_deviation
      type:
        map:
          elementType:
            list:
              elementType:
                namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.BenchmarkDeviation
              elementRelationship: atomic
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.BlockSpec
  map:
    fields:
    - name: model
      type:
        scalar: string
    - name: name
      type:
        scalar: string
    - name: partitionTable
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.PartitionTableSpec
    - name: rotational
      type:
        scalar: boolean
      default: false
    - name: size
      type:
        scalar: numeric
    - name: system
      type:
        scalar: string
    - name: type
      type:
        scalar: string
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.CPUSpec
  map:
    fields:
    - name: addressSizes
      type:
        scalar: string
    - name: bogoMips
      type:
        namedType: io.k8s.apimachinery.pkg.api.resource.Quantity
      default: {}
    - name: bugs
      type:
        list:
          elementType:
            scalar: string
          elementRelationship: atomic
    - name: cacheAlignment
      type:
        scalar: numeric
    - name: cacheSize
      type:
        scalar: string
    - name: clFlushSize
      type:
        scalar: numeric
    - name: cores
      type:
        scalar: numeric
    - name: cpuIdLevel
      type:
        scalar: numeric
    - name: family
      type:
        scalar: string
    - name: flags
      type:
        list:
          elementType:
            scalar: string
          elementRelationship: atomic
    - name: fpu
      type:
        scalar: boolean
      default: false
    - name: fpuException
      type:
        scalar: boolean
      default: false
    - name: logicalIds
      type:
        list:
          elementType:
            scalar: numeric
          elementRelationship: atomic
    - name: mhz
      type:
        namedType: io.k8s.apimachinery.pkg.api.resource.Quantity
      default: {}
    - name: microcode
      type:
        scalar: string
    - name: model
      type:
        scalar: string
    - name: modelName
      type:
        scalar: string
    - name: physicalId
      type:
        scalar: numeric
    - name: powerManagement
      type:
        scalar: string
    - name: siblings
      type:
        scalar: numeric
    - name: stepping
      type:
        scalar: string
    - name: vendorId
      type:
        scalar: string
    - name: vmxFlags
      type:
        list:
          elementType:
            scalar: string
          elementRelationship: atomic
    - name: wp
      type:
        scalar: boolean
      default: false
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.ConditionSpec
  map:
    fields:
    - name: lastTransitionTimestamp
      type:
        scalar: string
    - name: lastUpdateTimestamp
      type:
        scalar: string
    - name: message
      type:
        scalar: string
    - name: name
      type:
        scalar: string
    - name: reason
      type:
        scalar: string
    - name: state
      type:
        scalar: boolean
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.DistroSpec
  map:
    fields:
    - name: asicType
      type:
        scalar: string
    - name: buildBy
      type:
        scalar: string
    - name: buildDate
      type:
        scalar: string
    - name: buildNumber
      type:
        scalar: numeric
    - name: buildVersion
      type:
        scalar: string
    - name: commitID
      type:
        scalar: string
    - name: debianVersion
      type:
        scalar: string
    - name: kernelVersion
      type:
        scalar: string
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.FieldSelectorSpec
  map:
    fields:
    - name: fieldRef
      type:
        namedType: io.k8s.api.core.v1.ObjectFieldSelector
    - name: labelKey
      type:
        scalar: string
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.GeneralIPAMSpec
  map:
    fields:
    - name: addressFamily
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.AddressFamiliesMap
    - name: carrierSubnets
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.IPAMSelectionSpec
    - name: loopbackAddresses
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.IPAMSelectionSpec
    - name: loopbackSubnets
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.IPAMSelectionSpec
    - name: southSubnets
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.IPAMSelectionSpec
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.HostSpec
  map:
    fields:
    - name: name
      type:
        scalar: string
      default: ""
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.IPAMSelectionSpec
  map:
    fields:
    - name: fieldSelector
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.FieldSelectorSpec
    - name: labelSelector
      type:
        namedType: io.k8s.apimachinery.pkg.apis.meta.v1.LabelSelector
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.IPAMSpec
  map:
    fields:
    - name: loopbackAddresses
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.IPAMSelectionSpec
    - name: southSubnets
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.IPAMSelectionSpec
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.IPAddrSpec
  map:
    elementType:
      scalar: untyped
      list:
        elementType:
          namedType: __untyped_atomic_
        elementRelationship: atomic
      map:
        elementType:
          namedType: __untyped_deduced_
        elementRelationship: separable
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.IPAddressSpec
  map:
    fields:
    - name: address
      type:
        scalar: string
    - name: addressFamily
      type:
        scalar: string
    - name: extraAddress
      type:
        scalar: boolean
    - name: name
      type:
        scalar: string
    - name: namespace
      type:
        scalar: string
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.IPMISpec
  map:
    fields:
    - name: ipAddress
      type:
        scalar: string
    - name: macAddress
      type:
        scalar: string
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.Identity
  map:
    fields:
    - name: asset
      type:
        scalar: string
    - name: internal
      type:
        list:
          elementType:
            namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.Internal
          elementRelationship: atomic
    - name: serial_number
      type:
        scalar: string
    - name: sku
      type:
        scalar: string
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.Interface
  map:
    fields:
    - name: addresses
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.Addresses
      default: {}
    - name: lanes
      type:
        scalar: numeric
    - name: moved
      type:
        scalar: boolean
    - name: name
      type:
        scalar: string
    - name: peer
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.Peer
      default: {}
    - name: switch_reference
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.ResourceReference
    - name: unknown
      type:
        scalar: boolean
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.InterfaceOverridesSpec
  map:
    fields:
    - name: fec
      type:
        scalar: string
    - name: ip
      type:
        list:
          elementType:
            namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.AdditionalIPSpec
          elementRelationship: atomic
    - name: ipv4MaskLength
      type:
        scalar: numeric
    - name: ipv6Prefix
      type:
        scalar: numeric
    - name: lanes
      type:
        scalar: numeric
    - name: mtu
      type:
        scalar: numeric
    - name: name
      type:
        scalar: string
    - name: state
      type:
        scalar: string
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.InterfaceSpec
  map:
    fields:
    - name: direction
      type:
        scalar: string
    - name: fec
      type:
        scalar: string
    - name: ip
      type:
        list:
          elementType:
            namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.IPAddressSpec
          elementRelationship: atomic
    - name: ipv4MaskLength
      type:
        scalar: numeric
    - name: ipv6Prefix
      type:
        scalar: numeric
    - name: lanes
      type:
        scalar: numeric
    - name: macAddress
      type:
        scalar: string
    - name: mtu
      type:
        scalar: numeric
    - name: peer
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.PeerSpec
    - name: speed
      type:
        scalar: numeric
    - name: state
      type:
        scalar: string
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.InterfacesSpec
  map:
    fields:
    - name: defaults
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.PortParametersSpec
    - name: overrides
      type:
        list:
          elementType:
            namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.InterfaceOverridesSpec
          elementRelationship: atomic
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.Internal
  map:
    fields:
    - name: name
      type:
        scalar: string
    - name: value
      type:
        scalar: string
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.Inventory
  map:
    fields:
    - name: apiVersion
      type:
        scalar: string
    - name: kind
      type:
        scalar: string
    - name: metadata
      type:
        namedType: io.k8s.apimachinery.pkg.apis.meta.v1.ObjectMeta
      default: {}
    - name: spec
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.InventorySpec
      default: {}
    - name: status
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.InventoryStatus
      default: {}
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.InventorySpec
  map:
    fields:
    - name: blocks
      type:
        list:
          elementType:
            namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.BlockSpec
          elementRelationship: atomic
    - name: cpus
      type:
        list:
          elementType:
            namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.CPUSpec
          elementRelationship: atomic
    - name: distro
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.DistroSpec
    - name: host
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.HostSpec
    - name: ipmis
      type:
        list:
          elementType:
            namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.IPMISpec
          elementRelationship: atomic
    - name: memory
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.MemorySpec
    - name: nics
      type:
        list:
          elementType:
            namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.NICSpec
          elementRelationship: atomic
    - name: numa
      type:
        list:
          elementType:
            namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.NumaSpec
          elementRelationship: atomic
    - name: pciDevices
      type:
        list:
          elementType:
            namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.PCIDeviceSpec
          elementRelationship: atomic
    - name: system
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.SystemSpec
    - name: virt
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.VirtSpec
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.InventoryStatus
  map:
    fields:
    - name: computed
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.AggregationResults
      default: {}
    - name: inventoryStatuses
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.InventoryStatuses
      default: {}
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.InventoryStatuses
  map:
    fields:
    - name: ready
      type:
        scalar: boolean
      default: false
    - name: requestsCount
      type:
        scalar: numeric
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.JSONPath
  map:
    elementType:
      scalar: untyped
      list:
        elementType:
          namedType: __untyped_atomic_
        elementRelationship: atomic
      map:
        elementType:
          namedType: __untyped_deduced_
        elementRelationship: separable
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.LLDPSpec
  map:
    fields:
    - name: capabilities
      type:
        list:
          elementType:
            scalar: string
          elementRelationship: atomic
    - name: chassisId
      type:
        scalar: string
    - name: portDescription
      type:
        scalar: string
    - name: portId
      type:
        scalar: string
    - name: systemDescription
      type:
        scalar: string
    - name: systemName
      type:
        scalar: string
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.LoopbackAddresses
  map:
    fields:
    - name: ipv4
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.IPAddrSpec
      default: {}
    - name: ipv6
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.IPAddrSpec
      default: {}
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.Machine
  map:
    fields:
    - name: apiVersion
      type:
        scalar: string
    - name: kind
      type:
        scalar: string
    - name: metadata
      type:
        namedType: io.k8s.apimachinery.pkg.apis.meta.v1.ObjectMeta
      default: {}
    - name: spec
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.MachineSpec
      default: {}
    - name: status
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.MachineStatus
      default: {}
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.MachineSpec
  map:
    fields:
    - name: description
      type:
        scalar: string
    - name: hostname
      type:
        scalar: string
    - name: identity
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.Identity
      default: {}
    - name: inventory_requested
      type:
        scalar: boolean
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.MachineStatus
  map:
    fields:
    - name: health
      type:
        scalar: string
    - name: network
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.Network
      default: {}
    - name: orphaned
      type:
        scalar: boolean
    - name: reboot
      type:
        scalar: string
    - name: reservation
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.Reservation
      default: {}
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.MemorySpec
  map:
    fields:
    - name: total
      type:
        scalar: numeric
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.NDPSpec
  map:
    fields:
    - name: ipAddress
      type:
        scalar: string
    - name: macAddress
      type:
        scalar: string
    - name: state
      type:
        scalar: string
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.NICSpec
  map:
    fields:
    - name: activeFEC
      type:
        scalar: string
    - name: lanes
      type:
        scalar: numeric
    - name: lldps
      type:
        list:
          elementType:
            namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.LLDPSpec
          elementRelationship: atomic
    - name: macAddress
      type:
        scalar: string
    - name: mtu
      type:
        scalar: numeric
    - name: name
      type:
        scalar: string
    - name: ndps
      type:
        list:
          elementType:
            namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.NDPSpec
          elementRelationship: atomic
    - name: pciAddress
      type:
        scalar: string
    - name: speed
      type:
        scalar: numeric
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.Network
  map:
    fields:
    - name: asn
      type:
        scalar: numeric
    - name: interfaces
      type:
        list:
          elementType:
            namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.Interface
          elementRelationship: atomic
    - name: loopback_addresses
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.LoopbackAddresses
      default: {}
    - name: ports
      type:
        scalar: numeric
    - name: redundancy
      type:
        scalar: string
    - name: unknown_ports
      type:
        scalar: numeric
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.NetworkSwitch
  map:
    fields:
    - name: apiVersion
      type:
        scalar: string
    - name: kind
      type:
        scalar: string
    - name: metadata
      type:
        namedType: io.k8s.apimachinery.pkg.apis.meta.v1.ObjectMeta
      default: {}
    - name: spec
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.NetworkSwitchSpec
      default: {}
    - name: status
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.NetworkSwitchStatus
      default: {}
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.NetworkSwitchSpec
  map:
    fields:
    - name: configSelector
      type:
        namedType: io.k8s.apimachinery.pkg.apis.meta.v1.LabelSelector
    - name: cordon
      type:
        scalar: boolean
    - name: interfaces
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.InterfacesSpec
    - name: inventoryRef
      type:
        namedType: io.k8s.api.core.v1.LocalObjectReference
    - name: ipam
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.IPAMSpec
    - name: managed
      type:
        scalar: boolean
    - name: scanPorts
      type:
        scalar: boolean
    - name: topSpine
      type:
        scalar: boolean
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.NetworkSwitchStatus
  map:
    fields:
    - name: asn
      type:
        scalar: numeric
    - name: conditions
      type:
        list:
          elementType:
            namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.ConditionSpec
          elementRelationship: atomic
    - name: configRef
      type:
        namedType: io.k8s.api.core.v1.LocalObjectReference
      default: {}
    - name: interfaces
      type:
        map:
          elementType:
            namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.InterfaceSpec
    - name: layer
      type:
        scalar: numeric
      default: 0
    - name: loopbackAddresses
      type:
        list:
          elementType:
            namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.IPAddressSpec
          elementRelationship: atomic
    - name: message
      type:
        scalar: string
    - name: role
      type:
        scalar: string
    - name: routingConfigTemplate
      type:
        namedType: io.k8s.api.core.v1.LocalObjectReference
      default: {}
    - name: state
      type:
        scalar: string
    - name: subnets
      type:
        list:
          elementType:
            namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.SubnetSpec
          elementRelationship: atomic
    - name: switchPorts
      type:
        scalar: numeric
    - name: totalPorts
      type:
        scalar: numeric
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.NumaSpec
  map:
    fields:
    - name: cpus
      type:
        list:
          elementType:
            scalar: numeric
          elementRelationship: atomic
    - name: distances
      type:
        list:
          elementType:
            scalar: numeric
          elementRelationship: atomic
    - name: id
      type:
        scalar: numeric
      default: 0
    - name: memory
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.MemorySpec
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.ObjectReference
  map:
    fields:
    - name: name
      type:
        scalar: string
    - name: namespace
      type:
        scalar: string
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.PCIDeviceDescriptionSpec
  map:
    fields:
    - name: id
      type:
        scalar: string
      default: ""
    - name: name
      type:
        scalar: string
      default: ""
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.PCIDeviceSpec
  map:
    fields:
    - name: address
      type:
        scalar: string
    - name: busId
      type:
        scalar: string
    - name: class
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.PCIDeviceDescriptionSpec
    - name: interface
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.PCIDeviceDescriptionSpec
    - name: subclass
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.PCIDeviceDescriptionSpec
    - name: subtype
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.PCIDeviceDescriptionSpec
    - name: subvendor
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.PCIDeviceDescriptionSpec
    - name: type
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.PCIDeviceDescriptionSpec
    - name: vendor
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.PCIDeviceDescriptionSpec
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.PartitionSpec
  map:
    fields:
    - name: id
      type:
        scalar: string
    - name: name
      type:
        scalar: string
    - name: size
      type:
        scalar: numeric
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.PartitionTableSpec
  map:
    fields:
    - name: partitions
      type:
        list:
          elementType:
            namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.PartitionSpec
          elementRelationship: atomic
    - name: type
      type:
        scalar: string
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.Peer
  map:
    fields:
    - name: lldp_chassis_id
      type:
        scalar: string
    - name: lldp_port_description
      type:
        scalar: string
    - name: lldp_port_id
      type:
        scalar: string
    - name: lldp_system_name
      type:
        scalar: string
    - name: resource_reference
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.ResourceReference
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.PeerSpec
  map:
    fields:
    - name: chassisId
      type:
        scalar: string
    - name: name
      type:
        scalar: string
    - name: namespace
      type:
        scalar: string
    - name: portDescription
      type:
        scalar: string
    - name: portId
      type:
        scalar: string
    - name: systemName
      type:
        scalar: string
    - name: type
      type:
        scalar: string
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.PortParametersSpec
  map:
    fields:
    - name: fec
      type:
        scalar: string
    - name: ipv4MaskLength
      type:
        scalar: numeric
    - name: ipv6Prefix
      type:
        scalar: numeric
    - name: lanes
      type:
        scalar: numeric
    - name: mtu
      type:
        scalar: numeric
    - name: state
      type:
        scalar: string
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.RegionSpec
  map:
    fields:
    - name: availabilityZone
      type:
        scalar: string
    - name: name
      type:
        scalar: string
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.Reservation
  map:
    fields:
    - name: class
      type:
        scalar: string
    - name: reference
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.ResourceReference
    - name: status
      type:
        scalar: string
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.ResourceReference
  map:
    fields:
    - name: apiVersion
      type:
        scalar: string
    - name: kind
      type:
        scalar: string
    - name: name
      type:
        scalar: string
    - name: namespace
      type:
        scalar: string
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.SubnetSpec
  map:
    fields:
    - name: addressFamily
      type:
        scalar: string
    - name: cidr
      type:
        scalar: string
    - name: network
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.ObjectReference
    - name: region
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.RegionSpec
    - name: subnet
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.ObjectReference
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.SwitchConfig
  map:
    fields:
    - name: apiVersion
      type:
        scalar: string
    - name: kind
      type:
        scalar: string
    - name: metadata
      type:
        namedType: io.k8s.apimachinery.pkg.apis.meta.v1.ObjectMeta
      default: {}
    - name: spec
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.SwitchConfigSpec
      default: {}
    - name: status
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.SwitchConfigStatus
      default: {}
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.SwitchConfigSpec
  map:
    fields:
    - name: ipam
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.GeneralIPAMSpec
    - name: portsDefaults
      type:
        namedType: com.github.ironcore-dev.metal.apis.metal.v1alpha4.PortParametersSpec
    - name: routingConfigTemplate
      type:
        namedType: io.k8s.api.core.v1.LocalObjectReference
    - name: switches
      type:
        namedType: io.k8s.apimachinery.pkg.apis.meta.v1.LabelSelector
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.SwitchConfigStatus
  map:
    elementType:
      scalar: untyped
      list:
        elementType:
          namedType: __untyped_atomic_
        elementRelationship: atomic
      map:
        elementType:
          namedType: __untyped_deduced_
        elementRelationship: separable
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.SystemSpec
  map:
    fields:
    - name: id
      type:
        scalar: string
    - name: manufacturer
      type:
        scalar: string
    - name: productSku
      type:
        scalar: string
    - name: serialNumber
      type:
        scalar: string
- name: com.github.ironcore-dev.metal.apis.metal.v1alpha4.VirtSpec
  map:
    fields:
    - name: vmType
      type:
        scalar: string
- name: io.k8s.api.core.v1.LocalObjectReference
  map:
    fields:
    - name: name
      type:
        scalar: string
    elementRelationship: atomic
- name: io.k8s.api.core.v1.ObjectFieldSelector
  map:
    fields:
    - name: apiVersion
      type:
        scalar: string
    - name: fieldPath
      type:
        scalar: string
      default: ""
    elementRelationship: atomic
- name: io.k8s.apimachinery.pkg.api.resource.Quantity
  scalar: untyped
- name: io.k8s.apimachinery.pkg.apis.meta.v1.FieldsV1
  map:
    elementType:
      scalar: untyped
      list:
        elementType:
          namedType: __untyped_atomic_
        elementRelationship: atomic
      map:
        elementType:
          namedType: __untyped_deduced_
        elementRelationship: separable
- name: io.k8s.apimachinery.pkg.apis.meta.v1.LabelSelector
  map:
    fields:
    - name: matchExpressions
      type:
        list:
          elementType:
            namedType: io.k8s.apimachinery.pkg.apis.meta.v1.LabelSelectorRequirement
          elementRelationship: atomic
    - name: matchLabels
      type:
        map:
          elementType:
            scalar: string
    elementRelationship: atomic
- name: io.k8s.apimachinery.pkg.apis.meta.v1.LabelSelectorRequirement
  map:
    fields:
    - name: key
      type:
        scalar: string
      default: ""
    - name: operator
      type:
        scalar: string
      default: ""
    - name: values
      type:
        list:
          elementType:
            scalar: string
          elementRelationship: atomic
- name: io.k8s.apimachinery.pkg.apis.meta.v1.ManagedFieldsEntry
  map:
    fields:
    - name: apiVersion
      type:
        scalar: string
    - name: fieldsType
      type:
        scalar: string
    - name: fieldsV1
      type:
        namedType: io.k8s.apimachinery.pkg.apis.meta.v1.FieldsV1
    - name: manager
      type:
        scalar: string
    - name: operation
      type:
        scalar: string
    - name: subresource
      type:
        scalar: string
    - name: time
      type:
        namedType: io.k8s.apimachinery.pkg.apis.meta.v1.Time
- name: io.k8s.apimachinery.pkg.apis.meta.v1.ObjectMeta
  map:
    fields:
    - name: annotations
      type:
        map:
          elementType:
            scalar: string
    - name: creationTimestamp
      type:
        namedType: io.k8s.apimachinery.pkg.apis.meta.v1.Time
      default: {}
    - name: deletionGracePeriodSeconds
      type:
        scalar: numeric
    - name: deletionTimestamp
      type:
        namedType: io.k8s.apimachinery.pkg.apis.meta.v1.Time
    - name: finalizers
      type:
        list:
          elementType:
            scalar: string
          elementRelationship: associative
    - name: generateName
      type:
        scalar: string
    - name: generation
      type:
        scalar: numeric
    - name: labels
      type:
        map:
          elementType:
            scalar: string
    - name: managedFields
      type:
        list:
          elementType:
            namedType: io.k8s.apimachinery.pkg.apis.meta.v1.ManagedFieldsEntry
          elementRelationship: atomic
    - name: name
      type:
        scalar: string
    - name: namespace
      type:
        scalar: string
    - name: ownerReferences
      type:
        list:
          elementType:
            namedType: io.k8s.apimachinery.pkg.apis.meta.v1.OwnerReference
          elementRelationship: associative
          keys:
          - uid
    - name: resourceVersion
      type:
        scalar: string
    - name: selfLink
      type:
        scalar: string
    - name: uid
      type:
        scalar: string
- name: io.k8s.apimachinery.pkg.apis.meta.v1.OwnerReference
  map:
    fields:
    - name: apiVersion
      type:
        scalar: string
      default: ""
    - name: blockOwnerDeletion
      type:
        scalar: boolean
    - name: controller
      type:
        scalar: boolean
    - name: kind
      type:
        scalar: string
      default: ""
    - name: name
      type:
        scalar: string
      default: ""
    - name: uid
      type:
        scalar: string
      default: ""
    elementRelationship: atomic
- name: io.k8s.apimachinery.pkg.apis.meta.v1.Time
  scalar: untyped
- name: __untyped_atomic_
  scalar: untyped
  list:
    elementType:
      namedType: __untyped_atomic_
    elementRelationship: atomic
  map:
    elementType:
      namedType: __untyped_atomic_
    elementRelationship: atomic
- name: __untyped_deduced_
  scalar: untyped
  list:
    elementType:
      namedType: __untyped_atomic_
    elementRelationship: atomic
  map:
    elementType:
      namedType: __untyped_deduced_
    elementRelationship: separable
`)