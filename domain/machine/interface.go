// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package domain

import metalv1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"

func NewMachineInterface(
	interfaceName string,
	systemName string,
	chassisID string,
	portID string,
	portDescription string,
) metalv1alpha4.Interface {
	return metalv1alpha4.Interface{
		Name:    interfaceName,
		Unknown: false,
		Peer: metalv1alpha4.Peer{
			LLDPSystemName:      systemName,
			LLDPChassisID:       chassisID,
			LLDPPortID:          portID,
			LLDPPortDescription: portDescription,
		},
	}
}
