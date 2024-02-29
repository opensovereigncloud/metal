// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package dto

type SubnetInfo struct {
	Name             string
	Namespace        string
	Prefix           int
	ParentSubnetName string
}

func NewSubnetInfo(
	name string,
	namespace string,
	prefix int,
	parentSubnetName string,
) SubnetInfo {
	return SubnetInfo{
		Name:             name,
		Namespace:        namespace,
		Prefix:           prefix,
		ParentSubnetName: parentSubnetName,
	}
}
