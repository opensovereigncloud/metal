// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package switches

import (
	"bytes"
	"go/build"
	"os"
	"path/filepath"
	"testing"
	"time"

	ipamv1alpha1 "github.com/ironcore-dev/ipam/api/ipam/v1alpha1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/utils/ptr"

	metalv1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"
	"github.com/ironcore-dev/metal/pkg/constants"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestLabelFromFieldRef(t *testing.T) {
	t.Parallel()
	obj := &metalv1alpha4.NetworkSwitch{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sample-switch",
			Namespace: "metal-api",
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "NetworkSwitch",
			APIVersion: "v1beta1",
		},
		Spec: metalv1alpha4.NetworkSwitchSpec{
			Managed:   ptr.To(true),
			Cordon:    ptr.To(false),
			TopSpine:  ptr.To(true),
			ScanPorts: ptr.To(true),
		},
	}
	fieldSelector := &metalv1alpha4.FieldSelectorSpec{
		LabelKey: ptr.To("metal.ironcore.dev/object-owner"),
		FieldRef: &v1.ObjectFieldSelector{
			APIVersion: "v1beta1",
			FieldPath:  "metadata.name",
		},
	}
	expected := map[string]string{"metal.ironcore.dev/object-owner": "sample-switch"}
	label, err := labelFromFieldRef(obj, fieldSelector)
	assert.Nil(t, err)
	assert.Equal(t, expected, label)
}

func TestLabelFromFieldRefFail(t *testing.T) {
	t.Parallel()
	obj := &metalv1alpha4.NetworkSwitch{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sample-switch",
			Namespace: "metal-api",
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "NetworkSwitch",
			APIVersion: "v1beta1",
		},
		Spec: metalv1alpha4.NetworkSwitchSpec{
			Managed:   ptr.To(true),
			Cordon:    ptr.To(false),
			TopSpine:  ptr.To(true),
			ScanPorts: ptr.To(true),
		},
	}
	fieldSelector := &metalv1alpha4.FieldSelectorSpec{
		LabelKey: ptr.To("metal.ironcore.dev/object-owner"),
		FieldRef: &v1.ObjectFieldSelector{
			APIVersion: "v1",
			FieldPath:  "metadata.name",
		},
	}
	expected := "API version mismatch: expected v1beta1, actual v1"
	_, err := labelFromFieldRef(obj, fieldSelector)
	assert.NotNil(t, err)
	assert.Equal(t, expected, err.Error())
}

func TestCalculateASN(t *testing.T) {
	t.Parallel()
	loopbacksSamples := []*metalv1alpha4.IPAddressSpec{

		{
			Address:       ptr.To("100.64.0.1"),
			AddressFamily: ptr.To(constants.IPv4AF),
		},
		{
			Address:       ptr.To("fd00:afc0:e013:1003:ffff::"),
			AddressFamily: ptr.To(constants.IPv6AF),
		},
	}
	asn, err := CalculateASN(loopbacksSamples)
	expected := uint32(4_204_194_305)
	assert.Nil(t, err)
	assert.Equal(t, expected, asn)

	loopbacksSamples = []*metalv1alpha4.IPAddressSpec{
		{
			Address:       ptr.To("fd00:afc0:e013:1003:ffff::"),
			AddressFamily: ptr.To(constants.IPv6AF),
		},
	}
	asn, err = CalculateASN(loopbacksSamples)
	assert.Equal(t, uint32(0), asn)
	assert.NotNil(t, err)

	loopbacksSamples = []*metalv1alpha4.IPAddressSpec{
		{
			Address:       ptr.To("100.64.999.1"),
			AddressFamily: ptr.To(constants.IPv4AF),
		},
		{
			Address:       ptr.To("fd00:afc0:e013:1003:ffff::"),
			AddressFamily: ptr.To(constants.IPv6AF),
		},
	}
	asn, err = CalculateASN(loopbacksSamples)
	assert.Equal(t, uint32(0), asn)
	assert.NotNil(t, err)
}

func TestRequestIPs(t *testing.T) {
	t.Parallel()
	nicSample := &metalv1alpha4.InterfaceSpec{
		IP: []*metalv1alpha4.IPAddressSpec{
			{
				Address:       ptr.To("100.64.0.1/30"),
				AddressFamily: ptr.To(constants.IPv4AF),
			},
			{
				Address:       ptr.To("fd00:afc0:e013:1003:ffff::0/127"),
				AddressFamily: ptr.To(constants.IPv6AF),
				ObjectReference: &metalv1alpha4.ObjectReference{
					Name:      ptr.To("sample"),
					Namespace: ptr.To("default"),
				},
			},
		},
	}
	expectedIPs := []*metalv1alpha4.IPAddressSpec{
		{
			ObjectReference: nil,
			Address:         ptr.To("100.64.0.2/30"),
			AddressFamily:   ptr.To(constants.IPv4AF),
			ExtraAddress:    ptr.To(false),
		},
		{
			ObjectReference: &metalv1alpha4.ObjectReference{
				Name:      ptr.To("sample"),
				Namespace: ptr.To("default"),
			},
			Address:       ptr.To("fd00:afc0:e013:1003:ffff::1/127"),
			AddressFamily: ptr.To(constants.IPv6AF),
			ExtraAddress:  ptr.To(false),
		},
	}
	requestedIPs := RequestIPs(nicSample)
	assert.ElementsMatch(t, expectedIPs, requestedIPs)

	nicSample = &metalv1alpha4.InterfaceSpec{
		IP: []*metalv1alpha4.IPAddressSpec{
			{
				Address:       ptr.To("100.64.0.1/30"),
				AddressFamily: ptr.To(constants.IPv4AF),
			},
			{
				Address:       ptr.To("fd00:afc0:e013:1003:ffff::1/112"),
				AddressFamily: ptr.To(constants.IPv6AF),
				ObjectReference: &metalv1alpha4.ObjectReference{
					Name:      ptr.To("sample"),
					Namespace: ptr.To("default"),
				},
			},
		},
	}
	expectedIPs = []*metalv1alpha4.IPAddressSpec{
		{
			ObjectReference: nil,
			Address:         ptr.To("100.64.0.2/30"),
			AddressFamily:   ptr.To(constants.IPv4AF),
			ExtraAddress:    ptr.To(false),
		},
		{
			ObjectReference: &metalv1alpha4.ObjectReference{
				Name:      ptr.To("sample"),
				Namespace: ptr.To("default"),
			},
			Address:       ptr.To("fd00:afc0:e013:1003:ffff::2/112"),
			AddressFamily: ptr.To(constants.IPv6AF),
			ExtraAddress:  ptr.To(false),
		},
	}
	requestedIPs = RequestIPs(nicSample)
	assert.ElementsMatch(t, expectedIPs, requestedIPs)
}

func TestGetCrdPath(t *testing.T) {
	t.Parallel()
	expected := filepath.Join(build.Default.GOPATH, "pkg/mod/github.com/ironcore-dev/ipam@v0.1.0/config/crd/bases")
	computed, err := GetCrdPath(ipamv1alpha1.Subnet{}, filepath.Join("..", "..", "go.mod"))
	assert.Nil(t, err)
	assert.Equal(t, expected, computed)
}

func TestGetWebhookPath(t *testing.T) {
	t.Parallel()
	expected := filepath.Join(build.Default.GOPATH, "pkg/mod/github.com/ironcore-dev/ipam@v0.1.0/config/webhook")
	computed, err := GetWebhookPath(ipamv1alpha1.Subnet{}, filepath.Join("..", "..", "go.mod"))
	assert.Nil(t, err)
	assert.Equal(t, expected, computed)
}

func TestConditionsUpdated(t *testing.T) {
	t.Parallel()
	tsNow := time.Now()
	tsPast := tsNow.Add(-time.Hour)
	conditionsNow := []*metalv1alpha4.ConditionSpec{
		{
			Name:                    ptr.To(constants.ConditionInitialized),
			State:                   ptr.To(true),
			LastUpdateTimestamp:     ptr.To(tsNow.String()),
			LastTransitionTimestamp: ptr.To(tsNow.String()),
		},
		{
			Name:                    ptr.To(constants.ConditionInterfacesOK),
			State:                   ptr.To(true),
			LastUpdateTimestamp:     ptr.To(tsNow.String()),
			LastTransitionTimestamp: ptr.To(tsNow.String()),
		},
		{
			Name:                    ptr.To(constants.ConditionConfigRefOK),
			State:                   ptr.To(false),
			LastUpdateTimestamp:     ptr.To(tsNow.String()),
			LastTransitionTimestamp: ptr.To(tsNow.String()),
		},
	}
	conditionsPast := []*metalv1alpha4.ConditionSpec{
		{
			Name:                    ptr.To(constants.ConditionInitialized),
			State:                   ptr.To(true),
			LastUpdateTimestamp:     ptr.To(tsNow.String()),
			LastTransitionTimestamp: ptr.To(tsPast.String()),
		},
		{
			Name:                    ptr.To(constants.ConditionInterfacesOK),
			State:                   ptr.To(true),
			LastUpdateTimestamp:     ptr.To(tsNow.String()),
			LastTransitionTimestamp: ptr.To(tsNow.String()),
		},
		{
			Name:                    ptr.To(constants.ConditionConfigRefOK),
			State:                   ptr.To(false),
			LastUpdateTimestamp:     ptr.To(tsNow.String()),
			LastTransitionTimestamp: ptr.To(tsNow.String()),
		},
	}
	actual := conditionsUpdated(conditionsPast, conditionsNow)
	assert.True(t, actual)
}

func TestUpdateSwitchConfigSelector(t *testing.T) {
	t.Parallel()
	initialState := &metalv1alpha4.NetworkSwitch{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sample-switch",
			Namespace: "metal-api",
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "NetworkSwitch",
			APIVersion: "v1beta1",
		},
		Spec: metalv1alpha4.NetworkSwitchSpec{
			Managed:   ptr.To(true),
			Cordon:    ptr.To(false),
			TopSpine:  ptr.To(true),
			ScanPorts: ptr.To(true),
		},
		Status: metalv1alpha4.NetworkSwitchStatus{
			Layer: 255,
		},
	}
	testingState := initialState.DeepCopy()
	// no changes should be done - layer is not defined yet
	UpdateSwitchConfigSelector(testingState)
	assert.Equal(t, initialState, testingState)
	initialState = testingState.DeepCopy()

	// layer value applied - selector should be set by defaulting function
	testingState.SetLayer(1)
	UpdateSwitchConfigSelector(testingState)
	assert.NotEqual(t, initialState, testingState)
	appliedSelector := &metav1.LabelSelector{
		MatchLabels: map[string]string{constants.SwitchConfigLayerLabel: "1"},
	}
	assert.Equal(t, appliedSelector, testingState.Spec.ConfigSelector)
	initialState = testingState.DeepCopy()

	// layer value updated - selector should be updated by defaulting function
	testingState.SetLayer(2)
	UpdateSwitchConfigSelector(testingState)
	assert.NotEqual(t, initialState, testingState)
	updatedSelector := &metav1.LabelSelector{
		MatchLabels: map[string]string{constants.SwitchConfigLayerLabel: "2"},
	}
	assert.Equal(t, updatedSelector, testingState.Spec.ConfigSelector)
	initialState = testingState.DeepCopy()

	// selector updated - label related to layer deleted by defaulting function
	testingState.Spec.ConfigSelector.MatchLabels[constants.SwitchTypeLabel] = "spine"
	UpdateSwitchConfigSelector(testingState)
	assert.NotEqual(t, initialState, testingState)
	_, ok := testingState.Spec.ConfigSelector.MatchLabels[constants.SwitchConfigLayerLabel]
	assert.False(t, ok)
	initialState = testingState.DeepCopy()

	// selector flushed - selector updated with layer label by defaulting function
	delete(testingState.Spec.ConfigSelector.MatchLabels, constants.SwitchTypeLabel)
	UpdateSwitchConfigSelector(testingState)
	assert.NotEqual(t, initialState, testingState)
	_, ok = testingState.Spec.ConfigSelector.MatchLabels[constants.SwitchConfigLayerLabel]
	assert.True(t, ok)
	finalSelector := &metav1.LabelSelector{
		MatchLabels: map[string]string{constants.SwitchConfigLayerLabel: "2"},
	}
	assert.Equal(t, finalSelector, testingState.Spec.ConfigSelector)
}

func TestGetTotalAddressesCount(t *testing.T) {
	t.Parallel()
	var (
		q           *resource.Quantity
		addresses   int64
		samplesPath = filepath.Join("..", "..", "test_samples", "switch", "helpers_samples")
	)
	var (
		actualV4addresses int64 = 512
		actualV6addresses int64 = 8388608
	)
	samples, err := GetTestSamples(samplesPath)
	assert.Nil(t, err)
	sampleObjects := make([]*metalv1alpha4.NetworkSwitch, 0)
	for _, f := range samples {
		raw, err := os.ReadFile(f)
		assert.Nil(t, err)
		sampleYaml := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(raw), len(raw))
		obj := &metalv1alpha4.NetworkSwitch{}
		err = sampleYaml.Decode(obj)
		assert.Nil(t, err)
		sampleObjects = append(sampleObjects, obj)
	}

	for _, item := range sampleObjects {
		AlignInterfacesWithParams(item)
		q = GetTotalAddressesCount(item.Status.Interfaces, ipamv1alpha1.CIPv4SubnetType)
		addresses, _ = q.AsInt64()
		assert.Equal(t, addresses, actualV4addresses)
		q = GetTotalAddressesCount(item.Status.Interfaces, ipamv1alpha1.CIPv6SubnetType)
		addresses, _ = q.AsInt64()
		assert.Equal(t, addresses, actualV6addresses)
	}
}

func TestParseInterfaceNameFromSubnet(t *testing.T) {
	t.Parallel()
	sample := "b9a234a5-416b-3d49-a4f8-65b6f30c8ee5-ethernet120-aaaaaaaa"
	expected := "Ethernet120"
	assert.Equal(t, expected, ParseInterfaceNameFromSubnet(sample))
}
