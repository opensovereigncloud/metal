// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package switches

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	ipamv1alpha1 "github.com/onmetal/ipam/api/v1alpha1"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/yaml"

	metalv1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"
	"github.com/ironcore-dev/metal/pkg/constants"
)

var basePath = filepath.Join("..", "..", "test_samples", "switch")
var statuses = make(map[string]metalv1alpha4.NetworkSwitchStatus)

func loadSwitches() []*metalv1alpha4.NetworkSwitch {
	list := make([]*metalv1alpha4.NetworkSwitch, 0)
	samplesPath := filepath.Join(basePath, "switches")
	samples, _ := GetTestSamples(samplesPath)
	for _, sample := range samples {
		raw, _ := os.ReadFile(sample)
		obj := &metalv1alpha4.NetworkSwitch{}
		sampleYaml := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(raw), len(raw))
		_ = sampleYaml.Decode(obj)
		list = append(list, obj)
	}
	return list
}

func loadInventories() map[string]*metalv1alpha4.Inventory {
	inventoriesMap := make(map[string]*metalv1alpha4.Inventory)
	samplesPath := filepath.Join(basePath, "inventories")
	samples, _ := GetTestSamples(samplesPath)
	for _, sample := range samples {
		raw, _ := os.ReadFile(sample)
		obj := &metalv1alpha4.Inventory{}
		sampleYaml := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(raw), len(raw))
		_ = sampleYaml.Decode(obj)
		inventoriesMap[obj.Name] = obj
	}
	return inventoriesMap
}

func loadConfigs() map[string]*metalv1alpha4.SwitchConfig {
	configsMap := make(map[string]*metalv1alpha4.SwitchConfig)
	samplesPath := filepath.Join(basePath, "switch_configs")
	samples, _ := GetTestSamples(samplesPath)
	for _, sample := range samples {
		raw, _ := os.ReadFile(sample)
		obj := &metalv1alpha4.SwitchConfig{}
		sampleYaml := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(raw), len(raw))
		_ = sampleYaml.Decode(obj)
		configsMap[obj.Name] = obj
	}
	return configsMap
}

func loadLoopbacks() map[string]*ipamv1alpha1.IPList {
	loopbacks := make(map[string]*ipamv1alpha1.IPList)
	samplesPath := filepath.Join(basePath, "switch_ipam_objects", "loopbacks")
	samples, _ := GetTestSamples(samplesPath)
	for _, sample := range samples {
		raw, _ := os.ReadFile(sample)
		obj := &ipamv1alpha1.IP{}
		sampleYaml := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(raw), len(raw))
		_ = sampleYaml.Decode(obj)
		name := obj.Labels[constants.IPAMObjectOwnerLabel]
		if list, ok := loopbacks[name]; ok {
			list.Items = append(list.Items, *obj)
			continue
		}
		loopbacks[name] = &ipamv1alpha1.IPList{Items: []ipamv1alpha1.IP{*obj}}
	}
	return loopbacks
}

func loadSubnets() map[string]*ipamv1alpha1.SubnetList {
	subnets := make(map[string]*ipamv1alpha1.SubnetList)
	samplesPath := filepath.Join(basePath, "switch_ipam_objects", "subnets")
	samples, _ := GetTestSamples(samplesPath)
	for _, sample := range samples {
		raw, _ := os.ReadFile(sample)
		obj := &ipamv1alpha1.Subnet{}
		sampleYaml := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(raw), len(raw))
		_ = sampleYaml.Decode(obj)
		name := obj.Labels[constants.IPAMObjectOwnerLabel]
		if list, ok := subnets[name]; ok {
			list.Items = append(list.Items, *obj)
			continue
		}
		subnets[name] = &ipamv1alpha1.SubnetList{Items: []ipamv1alpha1.Subnet{*obj}}
	}
	return subnets
}

func copySwitchList(src []*metalv1alpha4.NetworkSwitch) []metalv1alpha4.NetworkSwitch {
	dst := make([]metalv1alpha4.NetworkSwitch, 0)
	for _, item := range src {
		dst = append(dst, *item)
	}
	return dst
}

//nolint:paralleltest
func TestInitialize(t *testing.T) {
	switches := loadSwitches()
	for _, testObject := range switches {
		result := Initialize(testObject, nil)
		assert.Nil(t, result.err)
		assert.Equal(t, result.condition, constants.ConditionInitialized)
		assert.Equal(t, result.reason, constants.ReasonConditionInitialized)
		assert.Equal(t, result.verboseMessage, constants.MessageConditionInitialized)
		statuses[testObject.Name] = testObject.Status
	}
}

//nolint:paralleltest
func TestUpdateInterfaces(t *testing.T) {
	switches := loadSwitches()
	inventories := loadInventories()
	for _, testObject := range switches {
		testObject.Status = statuses[testObject.Name]
		env := &SwitchEnvironment{}
		result := UpdateInterfaces(testObject, env)
		assert.Equal(t, result.condition, constants.ConditionInterfacesOK)
		assert.Equal(t, result.reason, constants.ErrorReasonMissingRequirements)
		assert.Equal(t, result.verboseMessage, constants.MessageMissingInventory)
		env.Inventory = inventories[testObject.Name]
		result = UpdateInterfaces(testObject, env)
		assert.Equal(t, result.condition, constants.ConditionInterfacesOK)
		assert.Equal(t, result.reason, constants.ReasonConditionInterfacesOK)
		assert.Equal(t, result.verboseMessage, constants.MessageConditionInterfacesOK)
		statuses[testObject.Name] = testObject.Status
	}
}

//nolint:paralleltest
func TestUpdateNeighbors(t *testing.T) {
	switches := loadSwitches()
	for _, testObject := range switches {
		testObject.Status = statuses[testObject.Name]
		env := &SwitchEnvironment{}
		result := UpdateNeighbors(testObject, env)
		assert.Equal(t, result.condition, constants.ConditionNeighborsOK)
		assert.Equal(t, result.reason, constants.ErrorReasonRequestFailed)
		assert.Equal(t, result.verboseMessage, fmt.Sprintf("%s: %s", constants.MessageRequestFailed, "NetworkSwitchList"))
		env.Switches = &metalv1alpha4.NetworkSwitchList{Items: copySwitchList(switches)}
		result = UpdateNeighbors(testObject, env)
		assert.Equal(t, result.condition, constants.ConditionNeighborsOK)
		assert.Equal(t, result.reason, constants.ReasonConditionNeighborsOK)
		assert.Equal(t, result.verboseMessage, constants.MessageConditionNeighborsOK)
		statuses[testObject.Name] = testObject.Status
	}
}

//nolint:paralleltest
func TestUpdateLayerAndRole(t *testing.T) {
	switches := loadSwitches()
	for _, item := range switches {
		item.Status = statuses[item.Name]
	}
	for _, testObject := range switches {
		if !testObject.TopSpine() {
			continue
		}
		env := &SwitchEnvironment{}
		result := UpdateLayerAndRole(testObject, env)
		assert.Equal(t, result.condition, constants.ConditionLayerAndRoleOK)
		assert.Equal(t, result.reason, constants.ErrorReasonRequestFailed)
		assert.Equal(t, result.verboseMessage, fmt.Sprintf("%s: %s", constants.MessageRequestFailed, "NetworkSwitchList"))
		env.Switches = &metalv1alpha4.NetworkSwitchList{Items: copySwitchList(switches)}
		result = UpdateLayerAndRole(testObject, env)
		assert.Equal(t, result.condition, constants.ConditionLayerAndRoleOK)
		assert.Equal(t, result.reason, constants.ReasonConditionLayerAndRoleOK)
		assert.Equal(t, result.verboseMessage, constants.MessageConditionLayerAndRoleOK)
		statuses[testObject.Name] = testObject.Status
	}
	for _, item := range switches {
		if !item.TopSpine() {
			continue
		}
		item.Status = statuses[item.Name]
	}
	for _, testObject := range switches {
		if testObject.TopSpine() {
			continue
		}
		env := &SwitchEnvironment{
			Switches: &metalv1alpha4.NetworkSwitchList{Items: copySwitchList(switches)},
		}
		result := UpdateLayerAndRole(testObject, env)
		assert.Equal(t, result.condition, constants.ConditionLayerAndRoleOK)
		assert.Equal(t, result.reason, constants.ReasonConditionLayerAndRoleOK)
		assert.Equal(t, result.verboseMessage, constants.MessageConditionLayerAndRoleOK)
		statuses[testObject.Name] = testObject.Status
	}
	for _, testObject := range switches {
		if testObject.TopSpine() {
			assert.Equal(t, testObject.GetLayer(), uint32(0))
			continue
		}
		assert.Equal(t, testObject.GetLayer(), uint32(1))
	}
}

//nolint:paralleltest
func TestUpdateConfigRef(t *testing.T) {
	switches := loadSwitches()
	configs := loadConfigs()
	for _, testObject := range switches {
		testObject.Status = statuses[testObject.Name]
		confName := "leafs-config"
		if testObject.TopSpine() {
			confName = "spines-config"
		}
		env := &SwitchEnvironment{}
		result := UpdateConfigRef(testObject, env)
		assert.Equal(t, result.condition, constants.ConditionConfigRefOK)
		assert.Equal(t, result.reason, constants.ErrorReasonMissingRequirements)
		assert.Equal(t, result.verboseMessage, constants.MessageFailedToDiscoverConfig)
		env.Config = configs[confName]
		result = UpdateConfigRef(testObject, env)
		assert.Equal(t, result.condition, constants.ConditionConfigRefOK)
		assert.Equal(t, result.reason, constants.ReasonConditionConfigRefOK)
		assert.Equal(t, result.verboseMessage, constants.MessageConditionConfigRefOK)
		assert.NotEmpty(t, testObject.Status.ConfigRef.Name)
		statuses[testObject.Name] = testObject.Status
	}
}

//nolint:paralleltest
func TestUpdatePortParameters(t *testing.T) {
	switches := loadSwitches()
	configs := loadConfigs()
	for _, item := range switches {
		item.Status = statuses[item.Name]
	}
	for _, testObject := range switches {
		confName := "leafs-config"
		if testObject.TopSpine() {
			confName = "spines-config"
		}
		env := &SwitchEnvironment{}
		result := UpdatePortParameters(testObject, env)
		assert.Equal(t, result.condition, constants.ConditionPortParametersOK)
		assert.Equal(t, result.reason, constants.ErrorReasonMissingRequirements)
		assert.Equal(t, result.verboseMessage, constants.MessageFailedToDiscoverConfig)
		env.Config = configs[confName]
		result = UpdatePortParameters(testObject, env)
		assert.Equal(t, result.condition, constants.ConditionPortParametersOK)
		assert.Equal(t, result.reason, constants.ErrorReasonRequestFailed)
		assert.Equal(t, result.verboseMessage, fmt.Sprintf("%s: %s", constants.MessageRequestFailed, "NetworkSwitchList"))
		env.Switches = &metalv1alpha4.NetworkSwitchList{Items: copySwitchList(switches)}
		result = UpdatePortParameters(testObject, env)
		assert.Equal(t, result.condition, constants.ConditionPortParametersOK)
		assert.Equal(t, result.reason, constants.ReasonConditionPortParametersOK)
		assert.Equal(t, result.verboseMessage, constants.MessageConditionPortParametersOK)
		statuses[testObject.Name] = testObject.Status
	}
}

//nolint:paralleltest
func TestUpdateLoopbacks(t *testing.T) {
	switches := loadSwitches()
	loopbacks := loadLoopbacks()
	for _, testObject := range switches {
		testObject.Status = statuses[testObject.Name]
		env := &SwitchEnvironment{}
		result := UpdateLoopbacks(testObject, env)
		assert.Equal(t, result.condition, constants.ConditionLoopbacksOK)
		assert.Equal(t, result.reason, constants.ErrorReasonMissingRequirements)
		assert.Equal(t, result.verboseMessage, constants.MessageMissingLoopbacks)
		env.LoopbackIPs = loopbacks[testObject.Name]
		result = UpdateLoopbacks(testObject, env)
		assert.Equal(t, result.condition, constants.ConditionLoopbacksOK)
		assert.Equal(t, result.reason, constants.ReasonConditionLoopbacksOK)
		assert.Equal(t, result.verboseMessage, constants.MessageConditionLoopbacksOK)
		assert.Equal(t, len(testObject.Status.LoopbackAddresses), 2)
		statuses[testObject.Name] = testObject.Status
	}
}

//nolint:paralleltest
func TestUpdateASN(t *testing.T) {
	switches := loadSwitches()
	for _, testObject := range switches {
		testObject.Status = statuses[testObject.Name]
		result := UpdateASN(testObject, nil)
		assert.Equal(t, result.condition, constants.ConditionAsnOK)
		assert.Equal(t, result.reason, constants.ReasonConditionAsnOK)
		assert.Equal(t, result.verboseMessage, constants.MessageConditionAsnOK)
		assert.NotZero(t, testObject.Status.ASN)
		statuses[testObject.Name] = testObject.Status
	}
}

//nolint:paralleltest
func TestUpdateSubnets(t *testing.T) {
	switches := loadSwitches()
	subnets := loadSubnets()
	for _, testObject := range switches {
		testObject.Status = statuses[testObject.Name]
		env := &SwitchEnvironment{}
		result := UpdateSubnets(testObject, env)
		assert.Equal(t, result.condition, constants.ConditionSubnetsOK)
		assert.Equal(t, result.reason, constants.ErrorReasonMissingRequirements)
		assert.Equal(t, result.verboseMessage, constants.MessageMissingSouthSubnets)
		env.SouthSubnets = subnets[testObject.Name]
		result = UpdateSubnets(testObject, env)
		assert.Equal(t, result.condition, constants.ConditionSubnetsOK)
		assert.Equal(t, result.reason, constants.ReasonConditionSubnetsOK)
		assert.Equal(t, result.verboseMessage, constants.MessageConditionSubnetsOK)
		assert.Equal(t, len(testObject.Status.Subnets), 2)
		statuses[testObject.Name] = testObject.Status
	}
}

//nolint:paralleltest
func TestUpdateSwitchPortIPs(t *testing.T) {
	switches := loadSwitches()
	for _, item := range switches {
		item.Status = statuses[item.Name]
	}
	for _, testObject := range switches {
		if !testObject.TopSpine() {
			continue
		}
		env := &SwitchEnvironment{}
		result := UpdateSwitchPortIPs(testObject, env)
		assert.Equal(t, result.condition, constants.ConditionIPAddressesOK)
		assert.Equal(t, result.reason, constants.ReasonConditionIPAddressesOK)
		assert.Equal(t, result.verboseMessage, constants.MessageConditionIPAddressesOK)
		for _, nicData := range testObject.Status.Interfaces {
			assert.Equal(t, 2, len(nicData.IP))
		}
		statuses[testObject.Name] = testObject.Status
	}
	for _, testObject := range switches {
		if testObject.TopSpine() {
			continue
		}
		env := &SwitchEnvironment{}
		result := UpdateSwitchPortIPs(testObject, env)
		assert.Equal(t, result.condition, constants.ConditionIPAddressesOK)
		assert.Equal(t, result.reason, constants.ErrorReasonIPAssignmentFailed)
		assert.Equal(t, result.verboseMessage, constants.ErrorUpdateSwitchPortIPsFailed)
		env = &SwitchEnvironment{Switches: &metalv1alpha4.NetworkSwitchList{Items: copySwitchList(switches)}}
		result = UpdateSwitchPortIPs(testObject, env)
		assert.Equal(t, result.condition, constants.ConditionIPAddressesOK)
		assert.Equal(t, result.reason, constants.ReasonConditionIPAddressesOK)
		assert.Equal(t, result.verboseMessage, constants.MessageConditionIPAddressesOK)
		for _, nicData := range testObject.Status.Interfaces {
			assert.Equal(t, 2, len(nicData.IP))
		}
		statuses[testObject.Name] = testObject.Status
	}
}
