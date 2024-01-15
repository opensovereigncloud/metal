<p>Packages:</p>
<ul>
<li>
<a href="#metal.ironcore.dev%2fv1beta1">metal.ironcore.dev/v1beta1</a>
</li>
</ul>
<h2 id="metal.ironcore.dev/v1beta1">metal.ironcore.dev/v1beta1</h2>
Resource Types:
<ul><li>
<a href="#metal.ironcore.dev/v1beta1.Switch">Switch</a>
</li><li>
<a href="#metal.ironcore.dev/v1beta1.SwitchConfig">SwitchConfig</a>
</li></ul>
<h3 id="metal.ironcore.dev/v1beta1.Switch">Switch
</h3>
<div>
<p>Switch is the Schema for switches API.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code><br/>
string</td>
<td>
<code>
metal.ironcore.dev/v1beta1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code><br/>
string
</td>
<td><code>Switch</code></td>
</tr>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://v1-21.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#metal.ironcore.dev/v1beta1.SwitchSpec">
SwitchSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>inventoryRef</code><br/>
<em>
<a href="https://v1-21.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#localobjectreference-v1-core">
Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<p>InventoryRef contains reference to corresponding inventory object
Empty InventoryRef means that there is no corresponding Inventory object</p>
</td>
</tr>
<tr>
<td>
<code>configSelector</code><br/>
<em>
<a href="https://v1-21.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#labelselector-v1-meta">
Kubernetes meta/v1.LabelSelector
</a>
</em>
</td>
<td>
<p>ConfigSelector contains selector to filter out corresponding SwitchConfig.
If the selector is not defined, it will be populated by defaulting webhook
with MatchLabels item, containing &lsquo;metal.ironcore.dev/layer&rsquo; key with value
equals to object&rsquo;s .status.layer.</p>
</td>
</tr>
<tr>
<td>
<code>managed</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Managed is a flag defining whether Switch object would be processed during reconciliation</p>
</td>
</tr>
<tr>
<td>
<code>cordon</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Cordon is a flag defining whether Switch object is taken offline</p>
</td>
</tr>
<tr>
<td>
<code>topSpine</code><br/>
<em>
bool
</em>
</td>
<td>
<p>TopSpine is a flag defining whether Switch is a top-level spine switch</p>
</td>
</tr>
<tr>
<td>
<code>scanPorts</code><br/>
<em>
bool
</em>
</td>
<td>
<p>ScanPorts is a flag defining whether to run periodical scanning on switch ports</p>
</td>
</tr>
<tr>
<td>
<code>ipam</code><br/>
<em>
<a href="#metal.ironcore.dev/v1beta1.IPAMSpec">
IPAMSpec
</a>
</em>
</td>
<td>
<p>IPAM refers to selectors for subnets which will be used for Switch object</p>
</td>
</tr>
<tr>
<td>
<code>interfaces</code><br/>
<em>
<a href="#metal.ironcore.dev/v1beta1.InterfacesSpec">
InterfacesSpec
</a>
</em>
</td>
<td>
<p>Interfaces contains general configuration for all switch ports</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#metal.ironcore.dev/v1beta1.SwitchStatus">
SwitchStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1beta1.SwitchConfig">SwitchConfig
</h3>
<div>
<p>SwitchConfig is the Schema for switch config API.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code><br/>
string</td>
<td>
<code>
metal.ironcore.dev/v1beta1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code><br/>
string
</td>
<td><code>SwitchConfig</code></td>
</tr>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://v1-21.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#metal.ironcore.dev/v1beta1.SwitchConfigSpec">
SwitchConfigSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>switches</code><br/>
<em>
<a href="https://v1-21.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#labelselector-v1-meta">
Kubernetes meta/v1.LabelSelector
</a>
</em>
</td>
<td>
<p>Switches contains label selector to pick up Switch objects</p>
</td>
</tr>
<tr>
<td>
<code>portsDefaults</code><br/>
<em>
<a href="#metal.ironcore.dev/v1beta1.PortParametersSpec">
PortParametersSpec
</a>
</em>
</td>
<td>
<p>PortsDefaults contains switch port parameters which will be applied to all ports of the switches
which fit selector conditions</p>
</td>
</tr>
<tr>
<td>
<code>ipam</code><br/>
<em>
<a href="#metal.ironcore.dev/v1beta1.GeneralIPAMSpec">
GeneralIPAMSpec
</a>
</em>
</td>
<td>
<p>IPAM refers to selectors for subnets which will be used for Switch object</p>
</td>
</tr>
<tr>
<td>
<code>routingConfigTemplate</code><br/>
<em>
<a href="https://v1-21.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#localobjectreference-v1-core">
Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<p>RoutingConfigTemplate contains the reference to the ConfigMap object which contains go-template for FRR config</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#metal.ironcore.dev/v1beta1.SwitchConfigStatus">
SwitchConfigStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1beta1.AdditionalIPSpec">AdditionalIPSpec
</h3>
<div>
<p>AdditionalIPSpec defines IP address and selector for subnet where address should be reserved</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>address</code><br/>
<em>
string
</em>
</td>
<td>
<p>Address contains additional IP address that should be assigned to the interface</p>
</td>
</tr>
<tr>
<td>
<code>parentSubnet</code><br/>
<em>
<a href="https://v1-21.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#labelselector-v1-meta">
Kubernetes meta/v1.LabelSelector
</a>
</em>
</td>
<td>
<p>ParentSubnet contains label selector to pick up IPAM objects</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1beta1.AddressFamiliesMap">AddressFamiliesMap
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1beta1.GeneralIPAMSpec">GeneralIPAMSpec</a>)
</p>
<div>
<p>AddressFamiliesMap contains flags regarding what IP address families should be used</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>ipv4</code><br/>
<em>
bool
</em>
</td>
<td>
<p>IPv4 is a flag defining whether IPv4 is used or not</p>
</td>
</tr>
<tr>
<td>
<code>ipv6</code><br/>
<em>
bool
</em>
</td>
<td>
<p>IPv6 is a flag defining whether IPv6 is used or not</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1beta1.ConditionSpec">ConditionSpec
</h3>
<div>
<p>ConditionSpec contains current condition of port parameters</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<p>Name reflects the name of the condition</p>
</td>
</tr>
<tr>
<td>
<code>state</code><br/>
<em>
bool
</em>
</td>
<td>
<p>State reflects the state of the condition</p>
</td>
</tr>
<tr>
<td>
<code>lastUpdateTimestamp</code><br/>
<em>
string
</em>
</td>
<td>
<p>LastUpdateTimestamp reflects the last timestamp when condition was updated</p>
</td>
</tr>
<tr>
<td>
<code>lastTransitionTimestamp</code><br/>
<em>
string
</em>
</td>
<td>
<p>LastTransitionTimestamp reflects the last timestamp when condition changed state from one to another</p>
</td>
</tr>
<tr>
<td>
<code>reason</code><br/>
<em>
string
</em>
</td>
<td>
<p>Reason reflects the reason of condition state</p>
</td>
</tr>
<tr>
<td>
<code>message</code><br/>
<em>
string
</em>
</td>
<td>
<p>Message reflects the verbose message about the reason</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1beta1.ConnectionsMap">ConnectionsMap
(<code>map[uint8]*..SwitchList</code> alias)</h3>
<div>
</div>
<h3 id="metal.ironcore.dev/v1beta1.FieldSelectorSpec">FieldSelectorSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1beta1.IPAMSelectionSpec">IPAMSelectionSpec</a>)
</p>
<div>
<p>FieldSelectorSpec contains label key and field path where to get label value for search</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>labelKey</code><br/>
<em>
string
</em>
</td>
<td>
<p>LabelKey contains label key</p>
</td>
</tr>
<tr>
<td>
<code>fieldRef</code><br/>
<em>
<a href="https://v1-21.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#objectfieldselector-v1-core">
Kubernetes core/v1.ObjectFieldSelector
</a>
</em>
</td>
<td>
<p>FieldRef contains reference to the field of resource where to get label&rsquo;s value</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1beta1.GeneralIPAMSpec">GeneralIPAMSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1beta1.SwitchConfigSpec">SwitchConfigSpec</a>)
</p>
<div>
<p>GeneralIPAMSpec contains definition of selectors, used to filter
required IPAM objects.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>addressFamily</code><br/>
<em>
<a href="#metal.ironcore.dev/v1beta1.AddressFamiliesMap">
AddressFamiliesMap
</a>
</em>
</td>
<td>
<p>AddressFamily contains flags to define which address families are used for switch subnets</p>
</td>
</tr>
<tr>
<td>
<code>carrierSubnets</code><br/>
<em>
<a href="#metal.ironcore.dev/v1beta1.IPAMSelectionSpec">
IPAMSelectionSpec
</a>
</em>
</td>
<td>
<p>CarrierSubnets contains label selector for Subnet object where switch&rsquo;s south subnet
should be reserved</p>
</td>
</tr>
<tr>
<td>
<code>loopbackSubnets</code><br/>
<em>
<a href="#metal.ironcore.dev/v1beta1.IPAMSelectionSpec">
IPAMSelectionSpec
</a>
</em>
</td>
<td>
<p>LoopbackSubnets contains label selector for Subnet object where switch&rsquo;s loopback
IP addresses should be reserved</p>
</td>
</tr>
<tr>
<td>
<code>southSubnets</code><br/>
<em>
<a href="#metal.ironcore.dev/v1beta1.IPAMSelectionSpec">
IPAMSelectionSpec
</a>
</em>
</td>
<td>
<p>SouthSubnets defines selector for subnets object which will be assigned to switch</p>
</td>
</tr>
<tr>
<td>
<code>loopbackAddresses</code><br/>
<em>
<a href="#metal.ironcore.dev/v1beta1.IPAMSelectionSpec">
IPAMSelectionSpec
</a>
</em>
</td>
<td>
<p>LoopbackAddresses defines selector for IP objects which should be referenced as switch&rsquo;s loopback addresses</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1beta1.IPAMSelectionSpec">IPAMSelectionSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1beta1.GeneralIPAMSpec">GeneralIPAMSpec</a>, <a href="#metal.ironcore.dev/v1beta1.IPAMSpec">IPAMSpec</a>)
</p>
<div>
<p>IPAMSelectionSpec contains label selector and address family</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>labelSelector</code><br/>
<em>
<a href="https://v1-21.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#labelselector-v1-meta">
Kubernetes meta/v1.LabelSelector
</a>
</em>
</td>
<td>
<p>LabelSelector contains label selector to pick up IPAM objects</p>
</td>
</tr>
<tr>
<td>
<code>fieldSelector</code><br/>
<em>
<a href="#metal.ironcore.dev/v1beta1.FieldSelectorSpec">
FieldSelectorSpec
</a>
</em>
</td>
<td>
<p>FieldSelector contains label key and field path where to get label value for search.
If FieldSelector is used as part of IPAM configuration in SwitchConfig object it will
reference to the field path in related Switch object. If FieldSelector is used as part of IPAM
configuration in Switch object, it will reference to the field path in the same object</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1beta1.IPAMSpec">IPAMSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1beta1.SwitchSpec">SwitchSpec</a>)
</p>
<div>
<p>IPAMSpec contains selectors for subnets and loopback IPs and
definition of address families which should be claimed</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>southSubnets</code><br/>
<em>
<a href="#metal.ironcore.dev/v1beta1.IPAMSelectionSpec">
IPAMSelectionSpec
</a>
</em>
</td>
<td>
<p>SouthSubnets defines selector for subnet object which will be assigned to switch</p>
</td>
</tr>
<tr>
<td>
<code>loopbackAddresses</code><br/>
<em>
<a href="#metal.ironcore.dev/v1beta1.IPAMSelectionSpec">
IPAMSelectionSpec
</a>
</em>
</td>
<td>
<p>LoopbackAddresses defines selector for IP object which will be assigned to switch&rsquo;s loopback interface</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1beta1.IPAddressSpec">IPAddressSpec
</h3>
<div>
<p>IPAddressSpec defines interface&rsquo;s ip address info</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>ObjectReference</code><br/>
<em>
<a href="#metal.ironcore.dev/v1beta1.ObjectReference">
ObjectReference
</a>
</em>
</td>
<td>
<p>
(Members of <code>ObjectReference</code> are embedded into this type.)
</p>
<p>Contains information to locate the referenced object</p>
</td>
</tr>
<tr>
<td>
<code>address</code><br/>
<em>
string
</em>
</td>
<td>
<p>Address refers to the ip address value</p>
</td>
</tr>
<tr>
<td>
<code>extraAddress</code><br/>
<em>
bool
</em>
</td>
<td>
<p>ExtraAddress is a flag defining whether address was added as additional by user</p>
</td>
</tr>
<tr>
<td>
<code>addressFamily</code><br/>
<em>
string
</em>
</td>
<td>
<p>AddressFamily refers to the AF of IP address</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1beta1.InterfaceOverridesSpec">InterfaceOverridesSpec
</h3>
<div>
<p>InterfaceOverridesSpec contains overridden parameters for certain switch port</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>PortParametersSpec</code><br/>
<em>
<a href="#metal.ironcore.dev/v1beta1.PortParametersSpec">
PortParametersSpec
</a>
</em>
</td>
<td>
<p>
(Members of <code>PortParametersSpec</code> are embedded into this type.)
</p>
<p>Contains port parameters overrides</p>
</td>
</tr>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<p>Name refers to switch port name</p>
</td>
</tr>
<tr>
<td>
<code>ip</code><br/>
<em>
<a href="#metal.ironcore.dev/v1beta1.*..AdditionalIPSpec">
[]*..AdditionalIPSpec
</a>
</em>
</td>
<td>
<p>IP contains a list of additional IP addresses for interface</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1beta1.InterfaceSpec">InterfaceSpec
</h3>
<div>
<p>InterfaceSpec defines the state of switch&rsquo;s interface</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>PortParametersSpec</code><br/>
<em>
<a href="#metal.ironcore.dev/v1beta1.PortParametersSpec">
PortParametersSpec
</a>
</em>
</td>
<td>
<p>
(Members of <code>PortParametersSpec</code> are embedded into this type.)
</p>
<p>Contains port parameters</p>
</td>
</tr>
<tr>
<td>
<code>macAddress</code><br/>
<em>
string
</em>
</td>
<td>
<p>MACAddress refers to the interface&rsquo;s hardware address
validation pattern</p>
</td>
</tr>
<tr>
<td>
<code>speed</code><br/>
<em>
uint32
</em>
</td>
<td>
<p>Speed refers to interface&rsquo;s speed</p>
</td>
</tr>
<tr>
<td>
<code>ip</code><br/>
<em>
<a href="#metal.ironcore.dev/v1beta1.*..IPAddressSpec">
[]*..IPAddressSpec
</a>
</em>
</td>
<td>
<p>IP contains a list of IP addresses that are assigned to interface</p>
</td>
</tr>
<tr>
<td>
<code>direction</code><br/>
<em>
string
</em>
</td>
<td>
<p>Direction refers to the interface&rsquo;s connection &lsquo;direction&rsquo;</p>
</td>
</tr>
<tr>
<td>
<code>peer</code><br/>
<em>
<a href="#metal.ironcore.dev/v1beta1.PeerSpec">
PeerSpec
</a>
</em>
</td>
<td>
<p>Peer refers to the info about device connected to current switch port</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1beta1.InterfacesSpec">InterfacesSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1beta1.SwitchSpec">SwitchSpec</a>)
</p>
<div>
<p>InterfacesSpec contains definitions for general switch ports&rsquo; configuration</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>defaults</code><br/>
<em>
<a href="#metal.ironcore.dev/v1beta1.PortParametersSpec">
PortParametersSpec
</a>
</em>
</td>
<td>
<p>Defaults contains switch port parameters which will be applied to all ports of the switches</p>
</td>
</tr>
<tr>
<td>
<code>overrides</code><br/>
<em>
<a href="#metal.ironcore.dev/v1beta1.*..InterfaceOverridesSpec">
[]*..InterfaceOverridesSpec
</a>
</em>
</td>
<td>
<p>Overrides contains set of parameters which should be overridden for listed switch ports</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1beta1.ObjectReference">ObjectReference
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1beta1.IPAddressSpec">IPAddressSpec</a>, <a href="#metal.ironcore.dev/v1beta1.PeerSpec">PeerSpec</a>, <a href="#metal.ironcore.dev/v1beta1.SubnetSpec">SubnetSpec</a>)
</p>
<div>
<p>ObjectReference contains enough information to let you locate the
referenced object across namespaces.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<p>Name contains name of the referenced object</p>
</td>
</tr>
<tr>
<td>
<code>namespace</code><br/>
<em>
string
</em>
</td>
<td>
<p>Namespace contains namespace of the referenced object</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1beta1.PeerInfoSpec">PeerInfoSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1beta1.PeerSpec">PeerSpec</a>)
</p>
<div>
<p>PeerInfoSpec contains LLDP info about peer</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>chassisId</code><br/>
<em>
string
</em>
</td>
<td>
<p>ChassisID refers to the chassis identificator - either MAC-address or system uuid
validation pattern</p>
</td>
</tr>
<tr>
<td>
<code>systemName</code><br/>
<em>
string
</em>
</td>
<td>
<p>SystemName refers to the advertised peer&rsquo;s name</p>
</td>
</tr>
<tr>
<td>
<code>portId</code><br/>
<em>
string
</em>
</td>
<td>
<p>PortID refers to the advertised peer&rsquo;s port ID</p>
</td>
</tr>
<tr>
<td>
<code>portDescription</code><br/>
<em>
string
</em>
</td>
<td>
<p>PortDescription refers to the advertised peer&rsquo;s port description</p>
</td>
</tr>
<tr>
<td>
<code>type</code><br/>
<em>
string
</em>
</td>
<td>
<p>Type refers to the peer type</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1beta1.PeerSpec">PeerSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1beta1.InterfaceSpec">InterfaceSpec</a>)
</p>
<div>
<p>PeerSpec defines peer info</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>ObjectReference</code><br/>
<em>
<a href="#metal.ironcore.dev/v1beta1.ObjectReference">
ObjectReference
</a>
</em>
</td>
<td>
<p>
(Members of <code>ObjectReference</code> are embedded into this type.)
</p>
<p>Contains information to locate the referenced object</p>
</td>
</tr>
<tr>
<td>
<code>PeerInfoSpec</code><br/>
<em>
<a href="#metal.ironcore.dev/v1beta1.PeerInfoSpec">
PeerInfoSpec
</a>
</em>
</td>
<td>
<p>
(Members of <code>PeerInfoSpec</code> are embedded into this type.)
</p>
<p>Contains LLDP info about peer</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1beta1.PortParametersSpec">PortParametersSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1beta1.InterfaceOverridesSpec">InterfaceOverridesSpec</a>, <a href="#metal.ironcore.dev/v1beta1.InterfaceSpec">InterfaceSpec</a>, <a href="#metal.ironcore.dev/v1beta1.InterfacesSpec">InterfacesSpec</a>, <a href="#metal.ironcore.dev/v1beta1.SwitchConfigSpec">SwitchConfigSpec</a>)
</p>
<div>
<p>PortParametersSpec contains a set of parameters of switch port</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>lanes</code><br/>
<em>
uint32
</em>
</td>
<td>
<p>Lanes refers to a number of lanes used by switch port</p>
</td>
</tr>
<tr>
<td>
<code>mtu</code><br/>
<em>
uint32
</em>
</td>
<td>
<p>MTU refers to maximum transmission unit value which should be applied on switch port</p>
</td>
</tr>
<tr>
<td>
<code>ipv4MaskLength</code><br/>
<em>
uint32
</em>
</td>
<td>
<p>IPv4MaskLength defines prefix of subnet where switch port&rsquo;s IPv4 address should be reserved</p>
</td>
</tr>
<tr>
<td>
<code>ipv6Prefix</code><br/>
<em>
uint32
</em>
</td>
<td>
<p>IPv6Prefix defines prefix of subnet where switch port&rsquo;s IPv6 address should be reserved</p>
</td>
</tr>
<tr>
<td>
<code>fec</code><br/>
<em>
string
</em>
</td>
<td>
<p>FEC refers to forward error correction method which should be applied on switch port</p>
</td>
</tr>
<tr>
<td>
<code>state</code><br/>
<em>
string
</em>
</td>
<td>
<p>State defines default state of switch port</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1beta1.RegionSpec">RegionSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1beta1.SubnetSpec">SubnetSpec</a>)
</p>
<div>
<p>RegionSpec defines region info</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<p>Name refers to the switch&rsquo;s region</p>
</td>
</tr>
<tr>
<td>
<code>availabilityZone</code><br/>
<em>
string
</em>
</td>
<td>
<p>AvailabilityZone refers to the switch&rsquo;s availability zone</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1beta1.SubnetSpec">SubnetSpec
</h3>
<div>
<p>SubnetSpec defines switch&rsquo;s subnet info</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>subnet</code><br/>
<em>
<a href="#metal.ironcore.dev/v1beta1.ObjectReference">
ObjectReference
</a>
</em>
</td>
<td>
<p>Contains information to locate the referenced object</p>
</td>
</tr>
<tr>
<td>
<code>network</code><br/>
<em>
<a href="#metal.ironcore.dev/v1beta1.ObjectReference">
ObjectReference
</a>
</em>
</td>
<td>
<p>Contains information to locate the referenced object</p>
</td>
</tr>
<tr>
<td>
<code>cidr</code><br/>
<em>
string
</em>
</td>
<td>
<p>CIDR refers to subnet CIDR
validation pattern</p>
</td>
</tr>
<tr>
<td>
<code>region</code><br/>
<em>
<a href="#metal.ironcore.dev/v1beta1.RegionSpec">
RegionSpec
</a>
</em>
</td>
<td>
<p>Region refers to switch&rsquo;s region</p>
</td>
</tr>
<tr>
<td>
<code>addressFamily</code><br/>
<em>
string
</em>
</td>
<td>
<p>AddressFamily refers to the AF of subnet</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1beta1.SwitchConfigSpec">SwitchConfigSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1beta1.SwitchConfig">SwitchConfig</a>)
</p>
<div>
<p>SwitchConfigSpec contains desired configuration for selected switches.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>switches</code><br/>
<em>
<a href="https://v1-21.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#labelselector-v1-meta">
Kubernetes meta/v1.LabelSelector
</a>
</em>
</td>
<td>
<p>Switches contains label selector to pick up Switch objects</p>
</td>
</tr>
<tr>
<td>
<code>portsDefaults</code><br/>
<em>
<a href="#metal.ironcore.dev/v1beta1.PortParametersSpec">
PortParametersSpec
</a>
</em>
</td>
<td>
<p>PortsDefaults contains switch port parameters which will be applied to all ports of the switches
which fit selector conditions</p>
</td>
</tr>
<tr>
<td>
<code>ipam</code><br/>
<em>
<a href="#metal.ironcore.dev/v1beta1.GeneralIPAMSpec">
GeneralIPAMSpec
</a>
</em>
</td>
<td>
<p>IPAM refers to selectors for subnets which will be used for Switch object</p>
</td>
</tr>
<tr>
<td>
<code>routingConfigTemplate</code><br/>
<em>
<a href="https://v1-21.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#localobjectreference-v1-core">
Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<p>RoutingConfigTemplate contains the reference to the ConfigMap object which contains go-template for FRR config</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1beta1.SwitchConfigStatus">SwitchConfigStatus
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1beta1.SwitchConfig">SwitchConfig</a>)
</p>
<div>
<p>SwitchConfigStatus contains observed state of SwitchConfig</p>
</div>
<h3 id="metal.ironcore.dev/v1beta1.SwitchSpec">SwitchSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1beta1.Switch">Switch</a>)
</p>
<div>
<p>SwitchSpec contains desired state of resulting Switch configuration</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>inventoryRef</code><br/>
<em>
<a href="https://v1-21.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#localobjectreference-v1-core">
Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<p>InventoryRef contains reference to corresponding inventory object
Empty InventoryRef means that there is no corresponding Inventory object</p>
</td>
</tr>
<tr>
<td>
<code>configSelector</code><br/>
<em>
<a href="https://v1-21.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#labelselector-v1-meta">
Kubernetes meta/v1.LabelSelector
</a>
</em>
</td>
<td>
<p>ConfigSelector contains selector to filter out corresponding SwitchConfig.
If the selector is not defined, it will be populated by defaulting webhook
with MatchLabels item, containing &lsquo;metal.ironcore.dev/layer&rsquo; key with value
equals to object&rsquo;s .status.layer.</p>
</td>
</tr>
<tr>
<td>
<code>managed</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Managed is a flag defining whether Switch object would be processed during reconciliation</p>
</td>
</tr>
<tr>
<td>
<code>cordon</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Cordon is a flag defining whether Switch object is taken offline</p>
</td>
</tr>
<tr>
<td>
<code>topSpine</code><br/>
<em>
bool
</em>
</td>
<td>
<p>TopSpine is a flag defining whether Switch is a top-level spine switch</p>
</td>
</tr>
<tr>
<td>
<code>scanPorts</code><br/>
<em>
bool
</em>
</td>
<td>
<p>ScanPorts is a flag defining whether to run periodical scanning on switch ports</p>
</td>
</tr>
<tr>
<td>
<code>ipam</code><br/>
<em>
<a href="#metal.ironcore.dev/v1beta1.IPAMSpec">
IPAMSpec
</a>
</em>
</td>
<td>
<p>IPAM refers to selectors for subnets which will be used for Switch object</p>
</td>
</tr>
<tr>
<td>
<code>interfaces</code><br/>
<em>
<a href="#metal.ironcore.dev/v1beta1.InterfacesSpec">
InterfacesSpec
</a>
</em>
</td>
<td>
<p>Interfaces contains general configuration for all switch ports</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1beta1.SwitchStatus">SwitchStatus
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1beta1.Switch">Switch</a>)
</p>
<div>
<p>SwitchStatus contains observed state of Switch</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>configRef</code><br/>
<em>
<a href="https://v1-21.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#localobjectreference-v1-core">
Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<p>ConfigRef contains reference to corresponding SwitchConfig object
Empty ConfigRef means that there is no corresponding SwitchConfig object</p>
</td>
</tr>
<tr>
<td>
<code>routingConfigTemplate</code><br/>
<em>
<a href="https://v1-21.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#localobjectreference-v1-core">
Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<p>RoutingConfigTemplate contains the reference to the ConfigMap object which contains go-template for FRR config.
This field reflects the corresponding field of the related SwitchConfig object.</p>
</td>
</tr>
<tr>
<td>
<code>asn</code><br/>
<em>
uint32
</em>
</td>
<td>
<p>ASN contains current autonomous system number defined for switch</p>
</td>
</tr>
<tr>
<td>
<code>totalPorts</code><br/>
<em>
uint32
</em>
</td>
<td>
<p>TotalPorts refers to total number of ports</p>
</td>
</tr>
<tr>
<td>
<code>switchPorts</code><br/>
<em>
uint32
</em>
</td>
<td>
<p>SwitchPorts refers to the number of ports excluding management interfaces, loopback etc.</p>
</td>
</tr>
<tr>
<td>
<code>role</code><br/>
<em>
string
</em>
</td>
<td>
<p>Role refers to switch&rsquo;s role</p>
</td>
</tr>
<tr>
<td>
<code>layer</code><br/>
<em>
uint32
</em>
</td>
<td>
<p>Layer refers to switch&rsquo;s current position in connection hierarchy</p>
</td>
</tr>
<tr>
<td>
<code>interfaces</code><br/>
<em>
<a href="#metal.ironcore.dev/v1beta1.*..InterfaceSpec">
map[string]*..InterfaceSpec
</a>
</em>
</td>
<td>
<p>Interfaces refers to switch&rsquo;s interfaces configuration</p>
</td>
</tr>
<tr>
<td>
<code>subnets</code><br/>
<em>
<a href="#metal.ironcore.dev/v1beta1.*..SubnetSpec">
[]*..SubnetSpec
</a>
</em>
</td>
<td>
<p>Subnets refers to the switch&rsquo;s south subnets</p>
</td>
</tr>
<tr>
<td>
<code>loopbackAddresses</code><br/>
<em>
<a href="#metal.ironcore.dev/v1beta1.*..IPAddressSpec">
[]*..IPAddressSpec
</a>
</em>
</td>
<td>
<p>LoopbackAddresses refers to the switch&rsquo;s loopback addresses</p>
</td>
</tr>
<tr>
<td>
<code>state</code><br/>
<em>
string
</em>
</td>
<td>
<p>State is the current state of corresponding object or process</p>
</td>
</tr>
<tr>
<td>
<code>message</code><br/>
<em>
string
</em>
</td>
<td>
<p>Message contains a brief description of the current state</p>
</td>
</tr>
<tr>
<td>
<code>conditions</code><br/>
<em>
<a href="#metal.ironcore.dev/v1beta1.*..ConditionSpec">
[]*..ConditionSpec
</a>
</em>
</td>
<td>
<p>Condition contains state of port parameters</p>
</td>
</tr>
</tbody>
</table>
<hr/>
<p><em>
Generated with <code>gen-crd-api-reference-docs</code>
on git commit <code>95c3af5</code>.
</em></p>
