// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package v1alpha4

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	"github.com/ironcore-dev/metal/pkg/constants"
)

func TestSetDefaultConfigSelector(t *testing.T) {
	t.Parallel()
	initialState := &NetworkSwitch{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sample-switch",
			Namespace: "metal-api",
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "NetworkSwitch",
			APIVersion: "v1beta1",
		},
		Spec: NetworkSwitchSpec{
			Managed:   ptr.To(true),
			Cordon:    ptr.To(false),
			TopSpine:  ptr.To(true),
			ScanPorts: ptr.To(true),
		},
		Status: NetworkSwitchStatus{
			Layer: 255,
		},
	}
	testingState := initialState.DeepCopy()
	// no changes should be done - layer is not defined yet
	testingState.setDefaultConfigSelector()
	assert.Equal(t, initialState, testingState)
	initialState = testingState.DeepCopy()

	// layer value applied - selector should be set by defaulting function
	testingState.SetLayer(1)
	testingState.setDefaultConfigSelector()
	assert.NotEqual(t, initialState, testingState)
	appliedSelector := &metav1.LabelSelector{
		MatchLabels: map[string]string{constants.SwitchConfigLayerLabel: "1"},
	}
	assert.Equal(t, appliedSelector, testingState.Spec.ConfigSelector)
	initialState = testingState.DeepCopy()

	// layer value updated - selector should be updated by defaulting function
	testingState.SetLayer(2)
	testingState.setDefaultConfigSelector()
	assert.NotEqual(t, initialState, testingState)
	updatedSelector := &metav1.LabelSelector{
		MatchLabels: map[string]string{constants.SwitchConfigLayerLabel: "2"},
	}
	assert.Equal(t, updatedSelector, testingState.Spec.ConfigSelector)
	initialState = testingState.DeepCopy()

	// selector updated - label related to layer deleted by defaulting function
	testingState.Spec.ConfigSelector.MatchLabels[constants.SwitchTypeLabel] = "spine"
	testingState.setDefaultConfigSelector()
	assert.NotEqual(t, initialState, testingState)
	_, ok := testingState.Spec.ConfigSelector.MatchLabels[constants.SwitchConfigLayerLabel]
	assert.False(t, ok)
	initialState = testingState.DeepCopy()

	// selector flushed - selector updated with layer label by defaulting function
	delete(testingState.Spec.ConfigSelector.MatchLabels, constants.SwitchTypeLabel)
	testingState.setDefaultConfigSelector()
	assert.NotEqual(t, initialState, testingState)
	_, ok = testingState.Spec.ConfigSelector.MatchLabels[constants.SwitchConfigLayerLabel]
	assert.True(t, ok)
	finalSelector := &metav1.LabelSelector{
		MatchLabels: map[string]string{constants.SwitchConfigLayerLabel: "2"},
	}
	assert.Equal(t, finalSelector, testingState.Spec.ConfigSelector)
}
