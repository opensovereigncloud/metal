// Copyright 2023 OnMetal authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package providers

import (
	"context"
	"net/netip"

	"k8s.io/apimachinery/pkg/labels"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"

	metalv1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"
	"github.com/ironcore-dev/metal/pkg/constants"
	"github.com/ironcore-dev/metal/usecase/onboarding/dto"
)

type SwitchRepository struct {
	client ctrlclient.Client
}

func NewSwitchRepository(
	client ctrlclient.Client,
) *SwitchRepository {
	return &SwitchRepository{
		client: client,
	}
}

func (s *SwitchRepository) ByChassisID(chassisID string) (dto.SwitchInfo, error) {
	label := map[string]string{
		constants.LabelChassisID: chassisID,
	}
	labelBasedOptions := labelBasedOptions(label)
	sw, err := s.extractSwitchFromCluster(labelBasedOptions)
	if err != nil {
		return dto.SwitchInfo{}, err
	}
	if sw.StateNotReady() {
		return dto.SwitchInfo{}, errSwitchIsNotReady
	}
	return toSwitchInfo(sw), nil
}

func (s *SwitchRepository) extractSwitchFromCluster(
	listOptions *ctrlclient.ListOptions,
) (*metalv1alpha4.NetworkSwitch, error) {
	obj := &metalv1alpha4.NetworkSwitchList{}
	if err := s.
		client.
		List(
			context.Background(),
			obj,
			listOptions,
		); err != nil {
		return nil, err
	}
	if len(obj.Items) == 0 {
		return nil, errNotFound
	}
	return &obj.Items[0], nil
}

func toSwitchInfo(sw *metalv1alpha4.NetworkSwitch) dto.SwitchInfo {
	var lanes uint32
	if sw.Spec.Interfaces != nil && sw.Spec.Interfaces.Defaults != nil {
		lanes = sw.Spec.Interfaces.Defaults.GetLanes()
	}
	return dto.SwitchInfo{
		Name:           sw.Name,
		Lanes:          lanes,
		InterfacesInfo: toSwitchInterfaces(sw),
	}
}

func toSwitchInterfaces(sw *metalv1alpha4.NetworkSwitch) map[string]dto.Interface {
	swInterfaces := make(map[string]dto.Interface, len(sw.Status.Interfaces))
	for k, v := range sw.Status.Interfaces {
		if v == nil {
			continue
		}
		swInterfaces[k] = dto.Interface{IP: toSwitchIP(v.IP)}
	}
	return swInterfaces
}

func toSwitchIP(ips []*metalv1alpha4.IPAddressSpec) []netip.Prefix {
	switchIPs := make([]netip.Prefix, 0, len(ips))
	for ip := range ips {
		if ips[ip] == nil || ips[ip].Address == nil {
			continue
		}
		prefix, err := netip.ParsePrefix(*ips[ip].Address)
		if err != nil {
			continue
		}
		switchIPs = append(switchIPs, prefix)
	}
	return switchIPs
}
func labelBasedOptions(label map[string]string) *ctrlclient.ListOptions {
	return &ctrlclient.ListOptions{
		LabelSelector: ctrlclient.MatchingLabelsSelector{Selector: labels.SelectorFromSet(label)},
	}
}
