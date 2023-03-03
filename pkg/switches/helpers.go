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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go/build"
	"io/fs"
	"math"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"

	gocidr "github.com/apparentlymart/go-cidr/cidr"
	ipamv1alpha1 "github.com/onmetal/ipam/api/v1alpha1"
	"golang.org/x/mod/modfile"
	"inet.af/netaddr"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"

	inventoryv1alpha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	switchv1beta1 "github.com/onmetal/metal-api/apis/switch/v1beta1"
	"github.com/onmetal/metal-api/internal/constants"
	"github.com/onmetal/metal-api/pkg/errors"
)

var (
	// fixme: SONiC switches consider nic's ipv6 address to be a subnet if it is the first (index 0)
	//	 address if nic's subnet, which cause an issue with ip address removal. Hence ipv6 address
	//	 index should be configurable (cmd flag, spec field and whatnot)
	BaseIPv4AddressIndex int = 1
	BaseIPv6AddressIndex int = 0
)

var PatchOpts *client.SubResourcePatchOptions = &client.SubResourcePatchOptions{
	PatchOptions: client.PatchOptions{
		Force:        pointer.Bool(true),
		FieldManager: "metal-api-controller-manager",
	},
}

func ApplyInterfacesFromInventory(obj *switchv1beta1.Switch, inventory *inventoryv1alpha1.Inventory) {
	if obj.Status.Interfaces == nil {
		obj.Status.Interfaces = make(map[string]*switchv1beta1.InterfaceSpec)
	}
	for _, item := range inventory.Spec.NICs {
		if !strings.HasPrefix(item.Name, constants.SwitchPortNamePrefix) {
			continue
		}
		iface := &switchv1beta1.InterfaceSpec{}
		iface.SetMACAddress(item.MACAddress)
		iface.SetSpeed(item.Speed)
		iface.SetDirection(constants.DirectionSouth)
		iface.SetIPEmpty()
		iface.SetPortParametersEmpty()
		if len(item.LLDPs) == 0 {
			obj.Status.Interfaces[item.Name] = iface
			continue
		}
		neighborData := item.LLDPs[0]
		peerData := &switchv1beta1.PeerSpec{
			PeerInfoSpec:    &switchv1beta1.PeerInfoSpec{},
			ObjectReference: &switchv1beta1.ObjectReference{},
		}
		peerData.SetChassisID(neighborData.ChassisID)
		peerData.SetSystemName(neighborData.SystemName)
		peerData.SetPortDescription(neighborData.PortDescription)
		peerData.SetPortID(neighborData.PortID)
		peerData.SetType(func() string {
			for _, capability := range neighborData.Capabilities {
				if capability == constants.LLDPCapabilityStation {
					return constants.NeighborTypeMachine
				}
			}
			return constants.NeighborTypeSwitch
		}())
		iface.Peer = peerData
		obj.Status.Interfaces[item.Name] = iface
	}
	obj.SetSwitchPorts(uint32(len(obj.Status.Interfaces)))
	obj.SetTotalPorts(uint32(len(inventory.Spec.NICs)))
}

func ApplyInterfaceParams(obj *switchv1beta1.Switch, config *switchv1beta1.SwitchConfig) {
	// set interfaces params:
	//   overrides - the highest priority
	//   defaults defined in switch spec - have higher priority then global params
	//   defined in switchConfig - applied if overrides & default are not defined
	basicParameters := config.Spec.PortsDefaults.DeepCopy()
	overriddenParameters := make(map[string]*switchv1beta1.PortParametersSpec)
	if obj.Spec.Interfaces == nil {
		for _, params := range obj.Status.Interfaces {
			params.PortParametersSpec = basicParameters.DeepCopy()
		}
		return
	}
	if obj.Spec.Interfaces.Defaults != nil {
		copyPortParams(obj.Spec.Interfaces.Defaults, basicParameters)
	}
	if obj.Spec.Interfaces.Overrides != nil && len(obj.Spec.Interfaces.Overrides) > 0 {
		for _, item := range obj.Spec.Interfaces.Overrides {
			overriddenParameters[item.GetName()] = item.PortParametersSpec
		}
	}
	for name, params := range obj.Status.Interfaces {
		if override, ok := overriddenParameters[name]; ok {
			copyPortParams(override, params.PortParametersSpec)
			continue
		}
		copyPortParams(basicParameters, params.PortParametersSpec)
	}
}

func ComputeLayer(obj *switchv1beta1.Switch, list *switchv1beta1.SwitchList) {
	connectionsMap, keys := buildConnectionMap(list)
	if _, ok := connectionsMap[0]; !ok {
		return
	}

	switch obj.GetTopSpine() {
	case true:
		obj.SetLayer(0)
		for _, nicData := range obj.Status.Interfaces {
			nicData.SetDirection(constants.DirectionSouth)
		}
		return
	case false:
		if obj.GetLayer() != 0 {
			break
		}
		obj.SetLayer(255)
		return
	}

	for _, connectionLevel := range keys {
		if connectionLevel == 255 {
			continue
		}
		if connectionLevel >= obj.GetLayer() {
			continue
		}
		switches := connectionsMap[connectionLevel]
		northPeers := getPeers(obj, switches)
		if len(northPeers.Items) == 0 {
			continue
		}
		obj.SetLayer(connectionLevel + 1)
		setNICsDirections(obj, list)
	}
}

func buildConnectionMap(obj *switchv1beta1.SwitchList) (map[uint32]*switchv1beta1.SwitchList, []uint32) {
	connectionsMap := make(map[uint32]*switchv1beta1.SwitchList)
	keys := make([]uint32, 0)
	for _, item := range obj.Items {
		if reflect.DeepEqual(item.Status, switchv1beta1.SwitchStatus{}) {
			continue
		}
		if item.Status.Layer == nil {
			continue
		}
		layer := item.GetLayer()
		list, ok := connectionsMap[layer]
		if !ok {
			list = &switchv1beta1.SwitchList{}
			list.Items = append(list.Items, item)
			connectionsMap[layer] = list
			keys = append(keys, layer)
			continue
		}
		list.Items = append(list.Items, item)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	return connectionsMap, keys
}

func getPeers(obj *switchv1beta1.Switch, switches *switchv1beta1.SwitchList) *switchv1beta1.SwitchList {
	result := &switchv1beta1.SwitchList{Items: make([]switchv1beta1.Switch, 0)}
	for _, item := range switches.Items {
		for _, nicData := range obj.Status.Interfaces {
			if nicData.Peer == nil {
				continue
			}
			if nicData.Peer.PeerInfoSpec == nil {
				continue
			}
			if reflect.DeepEqual(nicData.Peer.PeerInfoSpec, &switchv1beta1.PeerInfoSpec{}) {
				continue
			}
			peerChassisID := nicData.Peer.PeerInfoSpec.GetChassisID()
			if strings.ReplaceAll(peerChassisID, ":", "") == item.Annotations[constants.HardwareChassisIDAnnotation] {
				result.Items = append(result.Items, item)
			}
		}
	}
	return result
}

func setNICsDirections(obj *switchv1beta1.Switch, switches *switchv1beta1.SwitchList) {
	for _, item := range switches.Items {
		for _, nicData := range obj.Status.Interfaces {
			if nicData.Peer == nil {
				nicData.SetDirection(constants.DirectionSouth)
				continue
			}
			if nicData.Peer.ObjectReference == nil {
				nicData.SetDirection(constants.DirectionSouth)
				continue
			}
			if reflect.DeepEqual(nicData.Peer.PeerInfoSpec, &switchv1beta1.PeerInfoSpec{}) {
				nicData.SetDirection(constants.DirectionSouth)
				continue
			}
			peerChassisID := nicData.Peer.PeerInfoSpec.GetChassisID()
			peerFound := strings.ReplaceAll(peerChassisID, ":", "") == item.Annotations[constants.HardwareChassisIDAnnotation]
			peerLayer := item.GetLayer()
			objLayer := obj.GetLayer()
			peerIsNorth := objLayer > peerLayer
			peerIsSouth := objLayer < peerLayer
			if peerFound && peerIsNorth {
				nicData.SetDirection(constants.DirectionNorth)
				peerNICData := GetPeerData(item.Status.Interfaces, nicData.Peer.GetPortDescription(), nicData.Peer.GetPortID())
				copyPortParams(peerNICData.PortParametersSpec, nicData.PortParametersSpec)
			}
			if peerFound && peerIsSouth {
				nicData.SetDirection(constants.DirectionSouth)
			}
		}
	}
}

func SetRole(in *switchv1beta1.Switch) {
	in.SetRole(constants.SwitchRoleSpine)
	for _, data := range in.Status.Interfaces {
		if data.Peer == nil {
			continue
		}
		if data.Peer.GetType() == constants.NeighborTypeMachine {
			in.SetRole(constants.SwitchRoleLeaf)
			break
		}
	}
}

func copyPortParams(src, dst *switchv1beta1.PortParametersSpec) {
	if src == nil {
		return
	}
	if src.FEC != nil {
		dst.SetFEC(src.GetFEC())
	}
	if src.IPv4MaskLength != nil {
		dst.SetIPv4MaskLength(src.GetIPv4MaskLength())
	}
	if src.IPv6Prefix != nil {
		dst.SetIPv6Prefix(src.GetIPv6Prefix())
	}
	if src.Lanes != nil {
		dst.SetLanes(src.GetLanes())
	}
	if src.MTU != nil {
		dst.SetMTU(src.GetMTU())
	}
	if src.State != nil {
		dst.SetState(src.GetState())
	}
}

func ResultingLabels(
	obj *switchv1beta1.Switch,
	objectSelectors, globalSelectors *switchv1beta1.IPAMSelectionSpec,
) (map[string]string, error) {
	var err error
	result := make(map[string]string)
	selectors := objectSelectors
	if selectors == nil {
		selectors = globalSelectors
	}
	for k, v := range selectors.LabelSelector.MatchLabels {
		result[k] = v
	}
	fieldSelectors, err := labelFromFieldRef(obj, selectors.FieldSelector)
	if err != nil {
		return result, err
	}
	for k, v := range fieldSelectors {
		result[k] = v
	}
	return result, err
}

func GetSelectorFromIPAMSpec(
	obj *switchv1beta1.Switch, spec *switchv1beta1.IPAMSelectionSpec) (labels.Selector, error) {
	selector := labels.NewSelector()
	if spec != nil {
		for key, value := range spec.LabelSelector.MatchLabels {
			req, _ := labels.NewRequirement(key, selection.In, []string{value})
			selector = selector.Add(*req)
		}
		if len(spec.LabelSelector.MatchExpressions) > 0 {
			for _, item := range spec.LabelSelector.MatchExpressions {
				req, _ := labels.NewRequirement(item.Key, selection.Operator(item.Operator), item.Values)
				selector = selector.Add(*req)
			}
		}
		if spec.FieldSelector == nil {
			return selector, nil
		}
		ipamLabelFromFieldRef, err := labelFromFieldRef(*obj, spec.FieldSelector)
		if err != nil {
			return nil, err
		}
		for key, value := range ipamLabelFromFieldRef {
			req, _ := labels.NewRequirement(key, selection.In, []string{value})
			selector = selector.Add(*req)
		}
	}
	return selector, nil
}

func labelFromFieldRef(obj interface{}, src *switchv1beta1.FieldSelectorSpec) (map[string]string, error) {
	if src == nil {
		return nil, NewProcessingError(errors.MessageFieldSelectorNotDefined)
	}
	mapRepr, err := interfaceToMap(obj)
	if err != nil {
		return nil, err
	}
	apiVersion, ok := mapRepr["apiVersion"]
	if !ok {
		return nil, NewProcessingError(errors.MessageMissingAPIVersion)
	}
	if src.FieldRef.APIVersion != "" && apiVersion != src.FieldRef.APIVersion {
		return nil, NewProcessingError(
			"%s: expected %s, actual %s", errors.MessageAPIVersionMismatch, apiVersion, src.FieldRef.APIVersion)
	}
	nested := strings.Split(src.FieldRef.FieldPath, ".")
	label, err := processObjectMap(mapRepr, nested, src)
	if err != nil {
		return nil, err
	}
	return label, nil
}

func processObjectMap(
	repr map[string]interface{}, path []string, src *switchv1beta1.FieldSelectorSpec) (map[string]string, error) {
	var err error
	label := make(map[string]string)
	currentSearchObj := repr
	for i, f := range path {
		v, ok := currentSearchObj[f]
		if !ok {
			return nil, NewProcessingError("%s: %s", errors.MessageInvalidFieldPath, strings.Join(path, "."))
		}
		if i == len(path)-1 {
			switch v.(type) {
			case string:
				label[src.GetLabelKey()] = fmt.Sprintf("%v", v)
				return label, nil
			default:
				return nil, NewProcessingError(errors.MessageInvalidInputType)
			}
		}
		currentSearchObj, err = interfaceToMap(v)
		if err != nil {
			return nil, err
		}
	}
	return label, err
}

func interfaceToMap(i interface{}) (map[string]interface{}, error) {
	var raw []byte
	m := make(map[string]interface{})
	raw, err := json.Marshal(i)
	if err != nil {
		return nil, NewProcessingError(errors.MessageMarshallingFailed)
	}
	err = json.Unmarshal(raw, &m)
	if err != nil {
		return nil, NewProcessingError(errors.MessageUnmarshallingFailed)
	}
	return m, nil
}

func SetState(obj *switchv1beta1.Switch, state, message string) {
	obj.SetState(state)
	obj.SetMessage(message)
}

func CalculateASN(loopbacks []*switchv1beta1.IPAddressSpec) (uint32, error) {
	var result uint32 = 0
	for _, item := range loopbacks {
		if item.GetAddressFamily() != constants.IPv4AF {
			continue
		}
		asn := constants.ASNBase
		addr := net.ParseIP(item.GetAddress())
		if addr == nil {
			return 0, NewProcessingError("%s: %s", errors.MessageParseIPFailed, item.GetAddress())
		}
		asn += uint32(addr[13]) * uint32(math.Pow(2, 16))
		asn += uint32(addr[14]) * uint32(math.Pow(2, 8))
		asn += uint32(addr[15])
		result = asn
		break
	}
	if result == 0 {
		return 0, NewProcessingError(errors.MessageMissingLoopbackV4IP)
	}
	return result, nil
}

func GetExtraIPs(obj *switchv1beta1.Switch, name string) ([]*switchv1beta1.IPAddressSpec, error) {
	ips := make([]*switchv1beta1.IPAddressSpec, 0)
	if obj.Spec.Interfaces != nil && obj.Spec.Interfaces.Overrides != nil {
		for _, item := range obj.Spec.Interfaces.Overrides {
			if item.GetName() != name {
				continue
			}
			if len(item.IP) == 0 {
				continue
			}
			for _, data := range item.IP {
				ip := &switchv1beta1.IPAddressSpec{}
				ip.SetAddress(data.GetAddress())
				ip.SetExtraAddress(true)
				af, err := getAddressFamily(data.GetAddress())
				if err != nil {
					return ips, err
				}
				ip.SetAddressFamily(af)
				ips = append(ips, ip)
			}
		}
	}
	return ips, nil
}

func getAddressFamily(address string) (string, error) {
	addr, err := netaddr.ParseIP(address)
	if err != nil {
		return constants.EmptyString, NewProcessingError("%s: %s", errors.MessageParseIPFailed, address)
	}
	if addr.Is4() {
		return constants.IPv4AF, nil
	}
	return constants.IPv6AF, nil
}

func GetComputedIPs(
	obj *switchv1beta1.Switch, name string, data *switchv1beta1.InterfaceSpec) ([]*switchv1beta1.IPAddressSpec, error) {
	ips := make([]*switchv1beta1.IPAddressSpec, 0)
	for _, subnet := range obj.Status.Subnets {
		cidr, err := ipamv1alpha1.CIDRFromString(subnet.GetCIDR())
		if err != nil {
			return nil, err
		}
		if cidr == nil {
			return nil, NewProcessingError("%s: %s", errors.MessageParseCIDRFailed, subnet.GetCIDR())
		}
		mask := data.GetIPv4MaskLength()
		addrIndex := BaseIPv4AddressIndex
		af := constants.IPv4AF
		if cidr.IsIPv6() {
			mask = data.GetIPv6Prefix()
			addrIndex = BaseIPv6AddressIndex
			af = constants.IPv6AF
		}
		nicSubnet := getInterfaceSubnet(name, constants.SwitchPortNamePrefix, cidr.Net.IPNet(), mask)
		nicAddr, err := gocidr.Host(nicSubnet, addrIndex)
		if err != nil {
			return nil, err
		}
		ip := &switchv1beta1.IPAddressSpec{}
		ip.SetObjectReference(subnet.GetObjectReferenceName(), subnet.GetObjectReferenceNamespace())
		ip.SetAddress(fmt.Sprintf("%s/%d", nicAddr.String(), mask))
		ip.SetExtraAddress(false)
		ip.SetAddressFamily(af)
		ips = append(ips, ip)
	}
	return ips, nil
}

func getInterfaceSubnet(name string, namePrefix string, network *net.IPNet, mask uint32) *net.IPNet {
	index, _ := strconv.Atoi(strings.ReplaceAll(name, namePrefix, ""))
	prefix, _ := network.Mask.Size()
	ifaceNet, _ := gocidr.Subnet(network, int(mask)-prefix, index)
	return ifaceNet
}

func GetPeerData(
	interfaces map[string]*switchv1beta1.InterfaceSpec, portDesc, portID string) *switchv1beta1.InterfaceSpec {
	var nicData *switchv1beta1.InterfaceSpec
	if v, ok := interfaces[portDesc]; ok {
		nicData = v
	} else {
		nicData = interfaces[portID]
	}
	return nicData
}

func RequestIPs(peerNICData *switchv1beta1.InterfaceSpec) []*switchv1beta1.IPAddressSpec {
	requestedAddresses := make([]*switchv1beta1.IPAddressSpec, 0)
	for _, addr := range peerNICData.IP {
		_, cidr, _ := net.ParseCIDR(addr.GetAddress())
		addressIndex := BaseIPv4AddressIndex + 1
		if addr.GetAddressFamily() == constants.IPv6AF {
			addressIndex = BaseIPv6AddressIndex + 1
		}
		ip, _ := gocidr.Host(cidr, addressIndex)
		res := net.IPNet{IP: ip, Mask: cidr.Mask}
		address := &switchv1beta1.IPAddressSpec{}
		address.SetAddress(res.String())
		address.SetExtraAddress(false)
		address.SetAddressFamily(addr.GetAddressFamily())
		if addr.ObjectReference != nil {
			address.SetObjectReference(addr.GetObjectReferenceName(), addr.GetObjectReferenceNamespace())
		}
		requestedAddresses = append(requestedAddresses, address)
	}
	return requestedAddresses
}

func GetTotalAddressesCount(
	ports map[string]*switchv1beta1.InterfaceSpec, af ipamv1alpha1.SubnetAddressType) *resource.Quantity {
	var counter int64 = 0
	for _, item := range ports {
		var bits, ones uint32
		switch af {
		case ipamv1alpha1.CIPv4SubnetType:
			bits = constants.IPv4MaskLength
			ones = item.GetIPv4MaskLength()
		case ipamv1alpha1.CIPv6SubnetType:
			bits = constants.IPv6PrefixLength
			ones = item.GetIPv6Prefix()
		}
		addressesPerLane := math.Pow(2, float64(bits-ones))
		counter += int64(addressesPerLane) * int64(item.GetLanes())
	}
	return resource.NewQuantity(counter, resource.DecimalSI)
}

func AddressFamiliesMatchConfig(ipv4, ipv6 bool, foundAF map[ipamv1alpha1.SubnetAddressType]*bool) bool {
	result := true
	if ipv4 && foundAF[ipamv1alpha1.CIPv4SubnetType] == nil {
		foundAF[ipamv1alpha1.CIPv4SubnetType] = pointer.Bool(false)
		result = false
	}
	if ipv6 && foundAF[ipamv1alpha1.CIPv6SubnetType] == nil {
		foundAF[ipamv1alpha1.CIPv6SubnetType] = pointer.Bool(false)
		result = false
	}
	return result
}

func NeighborIsSwitch(nicData *switchv1beta1.InterfaceSpec) bool {
	if nicData.Peer == nil {
		return false
	}
	if nicData.Peer.ObjectReference == nil {
		return false
	}
	if nicData.Peer.PeerInfoSpec == nil {
		return false
	}
	if nicData.Peer.GetType() != constants.NeighborTypeSwitch {
		return false
	}
	return true
}

func ReconciliationRequired(objOld, objNew *switchv1beta1.Switch) bool {
	labelsChanged := !reflect.DeepEqual(objOld.GetLabels(), objNew.GetLabels())
	annotationsChanged := !reflect.DeepEqual(objOld.GetAnnotations(), objNew.GetAnnotations())
	metadataChanged := labelsChanged || annotationsChanged
	specChanged := !reflect.DeepEqual(objOld.Spec, objNew.Spec)
	conditionsChanged := conditionsUpdated(objOld.Status.Conditions, objNew.Status.Conditions)
	objOld.Status.Conditions = nil
	objNew.Status.Conditions = nil
	statusChanged := !reflect.DeepEqual(objOld.Status, objNew.Status)
	return metadataChanged || specChanged || statusChanged || conditionsChanged
}

func conditionsUpdated(oldData, newData []*switchv1beta1.ConditionSpec) bool {
	if len(oldData) != len(newData) {
		return true
	}
	for _, item := range oldData {
		item.LastUpdateTimestamp = nil
	}
	for _, item := range newData {
		item.LastUpdateTimestamp = nil
	}
	return !reflect.DeepEqual(oldData, newData)
}

// functions used in tests.
func GetCrdPath(crdPackageScheme interface{}) (string, error) {
	globalPackagePath := reflect.TypeOf(crdPackageScheme).PkgPath()
	goModData, err := os.ReadFile(filepath.Join("..", "..", "go.mod"))
	if err != nil {
		return "", err
	}
	goModFile, err := modfile.ParseLax("", goModData, nil)
	if err != nil {
		return "", err
	}
	globalModulePath := ""
	for _, req := range goModFile.Require {
		if strings.HasPrefix(globalPackagePath, req.Mod.Path) {
			globalModulePath = req.Mod.String()
			break
		}
	}
	return filepath.Join(build.Default.GOPATH, "pkg", "mod", globalModulePath, "config", "crd", "bases"), nil
}

func GetWebhookPath(crdPackageScheme interface{}) (string, error) {
	globalPackagePath := reflect.TypeOf(crdPackageScheme).PkgPath()
	goModData, err := os.ReadFile(filepath.Join("..", "..", "go.mod"))
	if err != nil {
		return "", err
	}
	goModFile, err := modfile.ParseLax("", goModData, nil)
	if err != nil {
		return "", err
	}
	globalModulePath := ""
	for _, req := range goModFile.Require {
		if strings.HasPrefix(globalPackagePath, req.Mod.Path) {
			globalModulePath = req.Mod.String()
			break
		}
	}
	return filepath.Join(build.Default.GOPATH, "pkg", "mod", globalModulePath, "config", "webhook"), nil
}

func GetTestSamples(path string) ([]string, error) {
	samples := make([]string, 0)
	if err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			samples = append(samples, path)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return samples, nil
}

func CreateSampleObject(ctx context.Context, c client.Client, obj client.Object, raw []byte) error {
	sampleYaml := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(raw), len(raw))
	err := sampleYaml.Decode(obj)
	if err != nil {
		return err
	}
	if err := c.Create(ctx, obj); err != nil {
		return err
	}
	return nil
}

type ProcessingError struct {
	message string
}

func NewProcessingError(message string, args ...interface{}) *ProcessingError {
	return &ProcessingError{message: fmt.Sprintf(message, args...)}
}

func (e *ProcessingError) Error() string {
	return e.message
}
