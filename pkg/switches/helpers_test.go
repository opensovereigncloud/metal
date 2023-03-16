/*
Copyright (c) 2022 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved

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
	"go/build"
	"path/filepath"
	"testing"
	"time"

	ipamv1alpha1 "github.com/onmetal/ipam/api/v1alpha1"

	switchv1beta1 "github.com/onmetal/metal-api/apis/switch/v1beta1"
	"github.com/onmetal/metal-api/pkg/constants"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

func TestLabelFromFieldRef(t *testing.T) {
	t.Parallel()
	obj := &switchv1beta1.Switch{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sample-switch",
			Namespace: "metal-api",
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "Switch",
			APIVersion: "v1beta1",
		},
		Spec: switchv1beta1.SwitchSpec{
			Managed:   pointer.Bool(true),
			Cordon:    pointer.Bool(false),
			TopSpine:  pointer.Bool(true),
			ScanPorts: pointer.Bool(true),
		},
	}
	fieldSelector := &switchv1beta1.FieldSelectorSpec{
		LabelKey: pointer.String("switch.onmetal.de/object-owner"),
		FieldRef: &v1.ObjectFieldSelector{
			APIVersion: "v1beta1",
			FieldPath:  "metadata.name",
		},
	}
	expected := map[string]string{"switch.onmetal.de/object-owner": "sample-switch"}
	label, err := labelFromFieldRef(obj, fieldSelector)
	assert.Nil(t, err)
	assert.Equal(t, expected, label)
}

func TestLabelFromFieldRefFail(t *testing.T) {
	t.Parallel()
	obj := &switchv1beta1.Switch{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sample-switch",
			Namespace: "metal-api",
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "Switch",
			APIVersion: "v1beta1",
		},
		Spec: switchv1beta1.SwitchSpec{
			Managed:   pointer.Bool(true),
			Cordon:    pointer.Bool(false),
			TopSpine:  pointer.Bool(true),
			ScanPorts: pointer.Bool(true),
		},
	}
	fieldSelector := &switchv1beta1.FieldSelectorSpec{
		LabelKey: pointer.String("switch.onmetal.de/object-owner"),
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
	loopbacksSamples := []*switchv1beta1.IPAddressSpec{

		{
			Address:       pointer.String("100.64.0.1"),
			AddressFamily: pointer.String(constants.IPv4AF),
		},
		{
			Address:       pointer.String("fd00:afc0:e013:1003:ffff::"),
			AddressFamily: pointer.String(constants.IPv6AF),
		},
	}
	asn, err := CalculateASN(loopbacksSamples)
	expected := uint32(4_204_194_305)
	assert.Nil(t, err)
	assert.Equal(t, expected, asn)

	loopbacksSamples = []*switchv1beta1.IPAddressSpec{
		{
			Address:       pointer.String("fd00:afc0:e013:1003:ffff::"),
			AddressFamily: pointer.String(constants.IPv6AF),
		},
	}
	asn, err = CalculateASN(loopbacksSamples)
	assert.Equal(t, uint32(0), asn)
	assert.NotNil(t, err)

	loopbacksSamples = []*switchv1beta1.IPAddressSpec{
		{
			Address:       pointer.String("100.64.999.1"),
			AddressFamily: pointer.String(constants.IPv4AF),
		},
		{
			Address:       pointer.String("fd00:afc0:e013:1003:ffff::"),
			AddressFamily: pointer.String(constants.IPv6AF),
		},
	}
	asn, err = CalculateASN(loopbacksSamples)
	assert.Equal(t, uint32(0), asn)
	assert.NotNil(t, err)
}

func TestRequestIPs(t *testing.T) {
	t.Parallel()
	nicSample := &switchv1beta1.InterfaceSpec{
		IP: []*switchv1beta1.IPAddressSpec{
			{
				Address:       pointer.String("100.64.0.1/30"),
				AddressFamily: pointer.String(constants.IPv4AF),
			},
			{
				Address:       pointer.String("fd00:afc0:e013:1003:ffff::0/127"),
				AddressFamily: pointer.String(constants.IPv6AF),
				ObjectReference: &switchv1beta1.ObjectReference{
					Name:      pointer.String("sample"),
					Namespace: pointer.String("default"),
				},
			},
		},
	}
	expectedIPs := []*switchv1beta1.IPAddressSpec{
		{
			ObjectReference: nil,
			Address:         pointer.String("100.64.0.2/30"),
			AddressFamily:   pointer.String(constants.IPv4AF),
			ExtraAddress:    pointer.Bool(false),
		},
		{
			ObjectReference: &switchv1beta1.ObjectReference{
				Name:      pointer.String("sample"),
				Namespace: pointer.String("default"),
			},
			Address:       pointer.String("fd00:afc0:e013:1003:ffff::1/127"),
			AddressFamily: pointer.String(constants.IPv6AF),
			ExtraAddress:  pointer.Bool(false),
		},
	}
	requestedIPs := RequestIPs(nicSample)
	assert.ElementsMatch(t, expectedIPs, requestedIPs)
}

func TestGetCrdPath(t *testing.T) {
	t.Parallel()
	expected := filepath.Join(build.Default.GOPATH, "pkg/mod/github.com/onmetal/ipam@v0.0.20/config/crd/bases")
	computed, err := GetCrdPath(ipamv1alpha1.Subnet{})
	assert.Nil(t, err)
	assert.Equal(t, expected, computed)
}

func TestGetWebhookPath(t *testing.T) {
	t.Parallel()
	expected := filepath.Join(build.Default.GOPATH, "pkg/mod/github.com/onmetal/ipam@v0.0.20/config/webhook")
	computed, err := GetWebhookPath(ipamv1alpha1.Subnet{})
	assert.Nil(t, err)
	assert.Equal(t, expected, computed)
}

func TestConditionsUpdated(t *testing.T) {
	t.Parallel()
	tsNow := time.Now()
	tsPast := tsNow.Add(-time.Hour)
	conditionsNow := []*switchv1beta1.ConditionSpec{
		{
			Name:                    pointer.String(constants.ConditionInitialized),
			State:                   pointer.Bool(true),
			LastUpdateTimestamp:     pointer.String(tsNow.String()),
			LastTransitionTimestamp: pointer.String(tsNow.String()),
		},
		{
			Name:                    pointer.String(constants.ConditionInterfacesOK),
			State:                   pointer.Bool(true),
			LastUpdateTimestamp:     pointer.String(tsNow.String()),
			LastTransitionTimestamp: pointer.String(tsNow.String()),
		},
		{
			Name:                    pointer.String(constants.ConditionConfigRefOK),
			State:                   pointer.Bool(false),
			LastUpdateTimestamp:     pointer.String(tsNow.String()),
			LastTransitionTimestamp: pointer.String(tsNow.String()),
		},
	}
	conditionsPast := []*switchv1beta1.ConditionSpec{
		{
			Name:                    pointer.String(constants.ConditionInitialized),
			State:                   pointer.Bool(true),
			LastUpdateTimestamp:     pointer.String(tsNow.String()),
			LastTransitionTimestamp: pointer.String(tsPast.String()),
		},
		{
			Name:                    pointer.String(constants.ConditionInterfacesOK),
			State:                   pointer.Bool(true),
			LastUpdateTimestamp:     pointer.String(tsNow.String()),
			LastTransitionTimestamp: pointer.String(tsNow.String()),
		},
		{
			Name:                    pointer.String(constants.ConditionConfigRefOK),
			State:                   pointer.Bool(false),
			LastUpdateTimestamp:     pointer.String(tsNow.String()),
			LastTransitionTimestamp: pointer.String(tsNow.String()),
		},
	}
	actual := conditionsUpdated(conditionsPast, conditionsNow)
	assert.True(t, actual)
}