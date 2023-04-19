package v1beta1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	"github.com/onmetal/metal-api/pkg/constants"
)

func TestSetDefaultConfigSelector(t *testing.T) {
	t.Parallel()
	initialState := &Switch{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sample-switch",
			Namespace: "metal-api",
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "Switch",
			APIVersion: "v1beta1",
		},
		Spec: SwitchSpec{
			Managed:   pointer.Bool(true),
			Cordon:    pointer.Bool(false),
			TopSpine:  pointer.Bool(true),
			ScanPorts: pointer.Bool(true),
		},
		Status: SwitchStatus{
			Layer: pointer.Uint32(255),
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
