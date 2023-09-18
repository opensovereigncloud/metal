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

package switches

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	ipamv1alpha1 "github.com/onmetal/ipam/api/v1alpha1"
	inventoryv1alpha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	switchv1beta1 "github.com/onmetal/metal-api/apis/switch/v1beta1"
	"github.com/onmetal/metal-api/pkg/constants"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/yaml"
)

var basePath = filepath.Join("..", "..", "test_samples", "switch")
var statuses = make(map[string]switchv1beta1.SwitchStatus)

func loadSwitches() []*switchv1beta1.Switch {
	list := make([]*switchv1beta1.Switch, 0)
	samplesPath := filepath.Join(basePath, "switches")
	samples, _ := GetTestSamples(samplesPath)
	for _, sample := range samples {
		raw, _ := os.ReadFile(sample)
		obj := &switchv1beta1.Switch{}
		sampleYaml := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(raw), len(raw))
		_ = sampleYaml.Decode(obj)
		list = append(list, obj)
	}
	return list
}

func loadInventories() map[string]*inventoryv1alpha1.Inventory {
	inventoriesMap := make(map[string]*inventoryv1alpha1.Inventory)
	samplesPath := filepath.Join(basePath, "inventories")
	samples, _ := GetTestSamples(samplesPath)
	for _, sample := range samples {
		raw, _ := os.ReadFile(sample)
		obj := &inventoryv1alpha1.Inventory{}
		sampleYaml := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(raw), len(raw))
		_ = sampleYaml.Decode(obj)
		inventoriesMap[obj.Name] = obj
	}
	return inventoriesMap
}

func loadConfigs() map[string]*switchv1beta1.SwitchConfig {
	configsMap := make(map[string]*switchv1beta1.SwitchConfig)
	samplesPath := filepath.Join(basePath, "switch_configs")
	samples, _ := GetTestSamples(samplesPath)
	for _, sample := range samples {
		raw, _ := os.ReadFile(sample)
		obj := &switchv1beta1.SwitchConfig{}
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

func copySwitchList(src []*switchv1beta1.Switch) []switchv1beta1.Switch {
	dst := make([]switchv1beta1.Switch, 0)
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
		assert.Equal(t, result.verboseMessage, fmt.Sprintf("%s: %s", constants.MessageRequestFailed, "SwitchList"))
		env.Switches = &switchv1beta1.SwitchList{Items: copySwitchList(switches)}
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
		assert.Equal(t, result.verboseMessage, fmt.Sprintf("%s: %s", constants.MessageRequestFailed, "SwitchList"))
		env.Switches = &switchv1beta1.SwitchList{Items: copySwitchList(switches)}
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
			Switches: &switchv1beta1.SwitchList{Items: copySwitchList(switches)},
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
		assert.Equal(t, result.verboseMessage, fmt.Sprintf("%s: %s", constants.MessageRequestFailed, "SwitchList"))
		env.Switches = &switchv1beta1.SwitchList{Items: copySwitchList(switches)}
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
		env = &SwitchEnvironment{Switches: &switchv1beta1.SwitchList{Items: copySwitchList(switches)}}
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
