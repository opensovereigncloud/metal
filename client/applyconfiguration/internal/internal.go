// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

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
- name: com.github.ironcore-dev.metal.api.v1alpha1.ConsoleProtocol
  map:
    fields:
    - name: name
      type:
        scalar: string
      default: ""
    - name: port
      type:
        scalar: numeric
      default: 0
- name: com.github.ironcore-dev.metal.api.v1alpha1.Machine
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
        namedType: com.github.ironcore-dev.metal.api.v1alpha1.MachineSpec
      default: {}
    - name: status
      type:
        namedType: com.github.ironcore-dev.metal.api.v1alpha1.MachineStatus
      default: {}
- name: com.github.ironcore-dev.metal.api.v1alpha1.MachineClaim
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
        namedType: com.github.ironcore-dev.metal.api.v1alpha1.MachineClaimSpec
      default: {}
    - name: status
      type:
        namedType: com.github.ironcore-dev.metal.api.v1alpha1.MachineClaimStatus
      default: {}
- name: com.github.ironcore-dev.metal.api.v1alpha1.MachineClaimNetworkInterface
  map:
    fields:
    - name: name
      type:
        scalar: string
      default: ""
    - name: prefix
      type:
        namedType: com.github.ironcore-dev.metal.api.v1alpha1.Prefix
- name: com.github.ironcore-dev.metal.api.v1alpha1.MachineClaimSpec
  map:
    fields:
    - name: ignitionSecretRef
      type:
        namedType: io.k8s.api.core.v1.LocalObjectReference
    - name: image
      type:
        scalar: string
      default: ""
    - name: machineRef
      type:
        namedType: io.k8s.api.core.v1.LocalObjectReference
    - name: machineSelector
      type:
        namedType: io.k8s.apimachinery.pkg.apis.meta.v1.LabelSelector
    - name: networkInterfaces
      type:
        list:
          elementType:
            namedType: com.github.ironcore-dev.metal.api.v1alpha1.MachineClaimNetworkInterface
          elementRelationship: atomic
    - name: power
      type:
        scalar: string
      default: ""
- name: com.github.ironcore-dev.metal.api.v1alpha1.MachineClaimStatus
  map:
    fields:
    - name: phase
      type:
        scalar: string
- name: com.github.ironcore-dev.metal.api.v1alpha1.MachineNetworkInterface
  map:
    fields:
    - name: IPRef
      type:
        namedType: io.k8s.api.core.v1.LocalObjectReference
    - name: macAddress
      type:
        scalar: string
      default: ""
    - name: name
      type:
        scalar: string
      default: ""
    - name: switchRef
      type:
        namedType: io.k8s.api.core.v1.LocalObjectReference
- name: com.github.ironcore-dev.metal.api.v1alpha1.MachineSpec
  map:
    fields:
    - name: asn
      type:
        scalar: string
    - name: inventoryRef
      type:
        namedType: io.k8s.api.core.v1.LocalObjectReference
    - name: locatorLED
      type:
        scalar: string
    - name: loopbackAddressRef
      type:
        namedType: io.k8s.api.core.v1.LocalObjectReference
    - name: machineClaimRef
      type:
        namedType: io.k8s.api.core.v1.ObjectReference
    - name: oobRef
      type:
        namedType: io.k8s.api.core.v1.LocalObjectReference
      default: {}
    - name: power
      type:
        scalar: string
    - name: uuid
      type:
        scalar: string
      default: ""
- name: com.github.ironcore-dev.metal.api.v1alpha1.MachineStatus
  map:
    fields:
    - name: conditions
      type:
        list:
          elementType:
            namedType: io.k8s.apimachinery.pkg.apis.meta.v1.Condition
          elementRelationship: associative
          keys:
          - type
    - name: locatorLED
      type:
        scalar: string
    - name: manufacturer
      type:
        scalar: string
    - name: networkInterfaces
      type:
        list:
          elementType:
            namedType: com.github.ironcore-dev.metal.api.v1alpha1.MachineNetworkInterface
          elementRelationship: atomic
    - name: power
      type:
        scalar: string
    - name: serialNumber
      type:
        scalar: string
    - name: shutdownDeadline
      type:
        namedType: io.k8s.apimachinery.pkg.apis.meta.v1.Time
    - name: sku
      type:
        scalar: string
    - name: state
      type:
        scalar: string
- name: com.github.ironcore-dev.metal.api.v1alpha1.OOB
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
        namedType: com.github.ironcore-dev.metal.api.v1alpha1.OOBSpec
      default: {}
    - name: status
      type:
        namedType: com.github.ironcore-dev.metal.api.v1alpha1.OOBStatus
      default: {}
- name: com.github.ironcore-dev.metal.api.v1alpha1.OOBSecret
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
        namedType: com.github.ironcore-dev.metal.api.v1alpha1.OOBSecretSpec
      default: {}
    - name: status
      type:
        namedType: com.github.ironcore-dev.metal.api.v1alpha1.OOBSecretStatus
      default: {}
- name: com.github.ironcore-dev.metal.api.v1alpha1.OOBSecretSpec
  map:
    fields:
    - name: expirationTime
      type:
        namedType: io.k8s.apimachinery.pkg.apis.meta.v1.Time
    - name: macAddress
      type:
        scalar: string
      default: ""
    - name: password
      type:
        scalar: string
      default: ""
    - name: username
      type:
        scalar: string
      default: ""
- name: com.github.ironcore-dev.metal.api.v1alpha1.OOBSecretStatus
  map:
    fields:
    - name: conditions
      type:
        list:
          elementType:
            namedType: io.k8s.apimachinery.pkg.apis.meta.v1.Condition
          elementRelationship: associative
          keys:
          - type
- name: com.github.ironcore-dev.metal.api.v1alpha1.OOBSpec
  map:
    fields:
    - name: consoleProtocol
      type:
        namedType: com.github.ironcore-dev.metal.api.v1alpha1.ConsoleProtocol
    - name: endpointRef
      type:
        namedType: io.k8s.api.core.v1.LocalObjectReference
    - name: flags
      type:
        map:
          elementType:
            scalar: string
    - name: macAddress
      type:
        scalar: string
      default: ""
    - name: protocol
      type:
        namedType: com.github.ironcore-dev.metal.api.v1alpha1.Protocol
    - name: secretRef
      type:
        namedType: io.k8s.api.core.v1.LocalObjectReference
- name: com.github.ironcore-dev.metal.api.v1alpha1.OOBStatus
  map:
    fields:
    - name: conditions
      type:
        list:
          elementType:
            namedType: io.k8s.apimachinery.pkg.apis.meta.v1.Condition
          elementRelationship: associative
          keys:
          - type
    - name: firmwareVersion
      type:
        scalar: string
    - name: manufacturer
      type:
        scalar: string
    - name: serialNumber
      type:
        scalar: string
    - name: sku
      type:
        scalar: string
    - name: state
      type:
        scalar: string
    - name: type
      type:
        scalar: string
- name: com.github.ironcore-dev.metal.api.v1alpha1.Prefix
  scalar: untyped
- name: com.github.ironcore-dev.metal.api.v1alpha1.Protocol
  map:
    fields:
    - name: name
      type:
        scalar: string
      default: ""
    - name: port
      type:
        scalar: numeric
      default: 0
- name: io.k8s.api.core.v1.LocalObjectReference
  map:
    fields:
    - name: name
      type:
        scalar: string
    elementRelationship: atomic
- name: io.k8s.api.core.v1.ObjectReference
  map:
    fields:
    - name: apiVersion
      type:
        scalar: string
    - name: fieldPath
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
    - name: resourceVersion
      type:
        scalar: string
    - name: uid
      type:
        scalar: string
    elementRelationship: atomic
- name: io.k8s.apimachinery.pkg.apis.meta.v1.Condition
  map:
    fields:
    - name: lastTransitionTime
      type:
        namedType: io.k8s.apimachinery.pkg.apis.meta.v1.Time
    - name: message
      type:
        scalar: string
      default: ""
    - name: observedGeneration
      type:
        scalar: numeric
    - name: reason
      type:
        scalar: string
      default: ""
    - name: status
      type:
        scalar: string
      default: ""
    - name: type
      type:
        scalar: string
      default: ""
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
