<p>Packages:</p>
<ul>
<li>
<a href="#switch.onmetal.de%2fv1alpha1">switch.onmetal.de/v1alpha1</a>
</li>
</ul>
<h2 id="switch.onmetal.de/v1alpha1">switch.onmetal.de/v1alpha1</h2>
Resource Types:
<ul></ul>
<h3 id="switch.onmetal.de/v1alpha1.ChassisSpec">ChassisSpec
</h3>
<p>
(<em>Appears on:</em><a href="#switch.onmetal.de/v1alpha1.SwitchSpec">SwitchSpec</a>)
</p>
<div>
<p>ChassisSpec defines switch&rsquo;s chassis info</p>
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
<code>manufacturer</code><br/>
<em>
string
</em>
</td>
<td>
<p>Manufactirer refers to the switch&rsquo;s manufacturer</p>
</td>
</tr>
<tr>
<td>
<code>serialNumber</code><br/>
<em>
string
</em>
</td>
<td>
<p>SerialNumber refers to the switch&rsquo;s serial number</p>
</td>
</tr>
<tr>
<td>
<code>sku</code><br/>
<em>
string
</em>
</td>
<td>
<p>SKU refers to the switch&rsquo;s stock keeping unit</p>
</td>
</tr>
</tbody>
</table>
<h3 id="switch.onmetal.de/v1alpha1.ConfManagerState">ConfManagerState
(<code>string</code> alias)</h3>
<p>
(<em>Appears on:</em><a href="#switch.onmetal.de/v1alpha1.ConfigurationSpec">ConfigurationSpec</a>)
</p>
<div>
</div>
<table>
<thead>
<tr>
<th>Value</th>
<th>Description</th>
</tr>
</thead>
<tbody><tr><td><p>&#34;active&#34;</p></td>
<td></td>
</tr><tr><td><p>&#34;failed&#34;</p></td>
<td></td>
</tr></tbody>
</table>
<h3 id="switch.onmetal.de/v1alpha1.ConfManagerType">ConfManagerType
(<code>string</code> alias)</h3>
<p>
(<em>Appears on:</em><a href="#switch.onmetal.de/v1alpha1.ConfigurationSpec">ConfigurationSpec</a>)
</p>
<div>
</div>
<table>
<thead>
<tr>
<th>Value</th>
<th>Description</th>
</tr>
</thead>
<tbody><tr><td><p>&#34;local&#34;</p></td>
<td></td>
</tr><tr><td><p>&#34;remote&#34;</p></td>
<td></td>
</tr></tbody>
</table>
<h3 id="switch.onmetal.de/v1alpha1.ConfigurationSpec">ConfigurationSpec
</h3>
<p>
(<em>Appears on:</em><a href="#switch.onmetal.de/v1alpha1.SwitchStatus">SwitchStatus</a>)
</p>
<div>
<p>ConfigurationSpec defines switch&rsquo;s computed configuration</p>
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
<code>managed</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Managed refers to whether switch configuration is managed or not</p>
</td>
</tr>
<tr>
<td>
<code>state</code><br/>
<em>
<a href="#switch.onmetal.de/v1alpha1.SwitchConfState">
SwitchConfState
</a>
</em>
</td>
<td>
<p>State refers to current switch&rsquo;s configuration processing state</p>
</td>
</tr>
<tr>
<td>
<code>managerType</code><br/>
<em>
<a href="#switch.onmetal.de/v1alpha1.ConfManagerType">
ConfManagerType
</a>
</em>
</td>
<td>
<p>Type refers to configuration manager type</p>
</td>
</tr>
<tr>
<td>
<code>managerState</code><br/>
<em>
<a href="#switch.onmetal.de/v1alpha1.ConfManagerState">
ConfManagerState
</a>
</em>
</td>
<td>
<p>State refers to configuration manager state</p>
</td>
</tr>
<tr>
<td>
<code>lastCheck</code><br/>
<em>
string
</em>
</td>
<td>
<p>LastCheck refers to the last timestamp when configuration was applied</p>
</td>
</tr>
</tbody>
</table>
<h3 id="switch.onmetal.de/v1alpha1.ConnectionsMap">ConnectionsMap
(<code>map[uint8]*..SwitchList</code> alias)</h3>
<div>
</div>
<h3 id="switch.onmetal.de/v1alpha1.FECType">FECType
(<code>string</code> alias)</h3>
<p>
(<em>Appears on:</em><a href="#switch.onmetal.de/v1alpha1.InterfaceSpec">InterfaceSpec</a>)
</p>
<div>
</div>
<table>
<thead>
<tr>
<th>Value</th>
<th>Description</th>
</tr>
</thead>
<tbody><tr><td><p>&#34;fc&#34;</p></td>
<td></td>
</tr><tr><td><p>&#34;none&#34;</p></td>
<td></td>
</tr><tr><td><p>&#34;rs&#34;</p></td>
<td></td>
</tr></tbody>
</table>
<h3 id="switch.onmetal.de/v1alpha1.IPAddressSpec">IPAddressSpec
</h3>
<p>
(<em>Appears on:</em><a href="#switch.onmetal.de/v1alpha1.InterfaceSpec">InterfaceSpec</a>, <a href="#switch.onmetal.de/v1alpha1.SwitchStatus">SwitchStatus</a>)
</p>
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
<code>address</code><br/>
<em>
string
</em>
</td>
<td>
<p>Address refers to the ip address value
validation pattern</p>
</td>
</tr>
<tr>
<td>
<code>resourceReference</code><br/>
<em>
<a href="#switch.onmetal.de/v1alpha1.ResourceReferenceSpec">
ResourceReferenceSpec
</a>
</em>
</td>
<td>
<p>ResourceReference refers to the related resource definition</p>
</td>
</tr>
</tbody>
</table>
<h3 id="switch.onmetal.de/v1alpha1.InterfaceSpec">InterfaceSpec
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
<code>fec</code><br/>
<em>
<a href="#switch.onmetal.de/v1alpha1.FECType">
FECType
</a>
</em>
</td>
<td>
<p>FEC refers to the current interface&rsquo;s forward error correction type</p>
</td>
</tr>
<tr>
<td>
<code>mtu</code><br/>
<em>
uint16
</em>
</td>
<td>
<p>MTU refers to the current value of interface&rsquo;s MTU</p>
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
<code>lanes</code><br/>
<em>
byte
</em>
</td>
<td>
<p>Lanes refers to the number of lanes used by interface</p>
</td>
</tr>
<tr>
<td>
<code>state</code><br/>
<em>
<a href="#switch.onmetal.de/v1alpha1.NICState">
NICState
</a>
</em>
</td>
<td>
<p>State refers to the current interface&rsquo;s operational state</p>
</td>
</tr>
<tr>
<td>
<code>ipV4</code><br/>
<em>
<a href="#switch.onmetal.de/v1alpha1.IPAddressSpec">
IPAddressSpec
</a>
</em>
</td>
<td>
<p>IPv4 refers to the interface&rsquo;s IPv4 address</p>
</td>
</tr>
<tr>
<td>
<code>ipV6</code><br/>
<em>
<a href="#switch.onmetal.de/v1alpha1.IPAddressSpec">
IPAddressSpec
</a>
</em>
</td>
<td>
<p>IPv6 refers to the interface&rsquo;s IPv6 address</p>
</td>
</tr>
<tr>
<td>
<code>direction</code><br/>
<em>
<a href="#switch.onmetal.de/v1alpha1.NICDirection">
NICDirection
</a>
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
<a href="#switch.onmetal.de/v1alpha1.PeerSpec">
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
<h3 id="switch.onmetal.de/v1alpha1.LinkedSwitchSpec">LinkedSwitchSpec
</h3>
<p>
(<em>Appears on:</em><a href="#switch.onmetal.de/v1alpha1.SwitchAssignmentStatus">SwitchAssignmentStatus</a>)
</p>
<div>
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
</td>
</tr>
</tbody>
</table>
<h3 id="switch.onmetal.de/v1alpha1.LocationSpec">LocationSpec
</h3>
<p>
(<em>Appears on:</em><a href="#switch.onmetal.de/v1alpha1.SwitchSpec">SwitchSpec</a>)
</p>
<div>
<p>LocationSpec defines switch&rsquo;s location</p>
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
<code>room</code><br/>
<em>
string
</em>
</td>
<td>
<p>Room refers to room name</p>
</td>
</tr>
<tr>
<td>
<code>row</code><br/>
<em>
int16
</em>
</td>
<td>
<p>Row refers to row number</p>
</td>
</tr>
<tr>
<td>
<code>rack</code><br/>
<em>
int16
</em>
</td>
<td>
<p>Rack refers to rack number</p>
</td>
</tr>
<tr>
<td>
<code>hu</code><br/>
<em>
int16
</em>
</td>
<td>
<p>HU refers to height in units</p>
</td>
</tr>
</tbody>
</table>
<h3 id="switch.onmetal.de/v1alpha1.NICDirection">NICDirection
(<code>string</code> alias)</h3>
<p>
(<em>Appears on:</em><a href="#switch.onmetal.de/v1alpha1.InterfaceSpec">InterfaceSpec</a>)
</p>
<div>
</div>
<table>
<thead>
<tr>
<th>Value</th>
<th>Description</th>
</tr>
</thead>
<tbody><tr><td><p>&#34;north&#34;</p></td>
<td></td>
</tr><tr><td><p>&#34;south&#34;</p></td>
<td></td>
</tr></tbody>
</table>
<h3 id="switch.onmetal.de/v1alpha1.NICState">NICState
(<code>string</code> alias)</h3>
<p>
(<em>Appears on:</em><a href="#switch.onmetal.de/v1alpha1.InterfaceSpec">InterfaceSpec</a>)
</p>
<div>
</div>
<table>
<thead>
<tr>
<th>Value</th>
<th>Description</th>
</tr>
</thead>
<tbody><tr><td><p>&#34;down&#34;</p></td>
<td></td>
</tr><tr><td><p>&#34;up&#34;</p></td>
<td></td>
</tr></tbody>
</table>
<h3 id="switch.onmetal.de/v1alpha1.PeerSpec">PeerSpec
</h3>
<p>
(<em>Appears on:</em><a href="#switch.onmetal.de/v1alpha1.InterfaceSpec">InterfaceSpec</a>)
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
<tr>
<td>
<code>resourceReference</code><br/>
<em>
<a href="#switch.onmetal.de/v1alpha1.ResourceReferenceSpec">
ResourceReferenceSpec
</a>
</em>
</td>
<td>
<p>ResourceReference refers to the related resource definition</p>
</td>
</tr>
</tbody>
</table>
<h3 id="switch.onmetal.de/v1alpha1.RegionSpec">RegionSpec
</h3>
<p>
(<em>Appears on:</em><a href="#switch.onmetal.de/v1alpha1.SubnetSpec">SubnetSpec</a>, <a href="#switch.onmetal.de/v1alpha1.SwitchAssignmentSpec">SwitchAssignmentSpec</a>)
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
<h3 id="switch.onmetal.de/v1alpha1.ResourceReferenceSpec">ResourceReferenceSpec
</h3>
<p>
(<em>Appears on:</em><a href="#switch.onmetal.de/v1alpha1.IPAddressSpec">IPAddressSpec</a>, <a href="#switch.onmetal.de/v1alpha1.PeerSpec">PeerSpec</a>, <a href="#switch.onmetal.de/v1alpha1.SubnetSpec">SubnetSpec</a>)
</p>
<div>
<p>ResourceReferenceSpec defines related resource info</p>
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
<em>
string
</em>
</td>
<td>
<p>APIVersion refers to the resource API version</p>
</td>
</tr>
<tr>
<td>
<code>kind</code><br/>
<em>
string
</em>
</td>
<td>
<p>Kind refers to the resource kind</p>
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
<p>Name refers to the resource name</p>
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
<p>Namespace refers to the resource namespace</p>
</td>
</tr>
</tbody>
</table>
<h3 id="switch.onmetal.de/v1alpha1.SoftwarePlatformSpec">SoftwarePlatformSpec
</h3>
<p>
(<em>Appears on:</em><a href="#switch.onmetal.de/v1alpha1.SwitchSpec">SwitchSpec</a>)
</p>
<div>
<p>SoftwarePlatformSpec defines switch&rsquo;s software base</p>
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
<code>onie</code><br/>
<em>
bool
</em>
</td>
<td>
<p>ONIE refers to whether open network installation environment is used</p>
</td>
</tr>
<tr>
<td>
<code>operatingSystem</code><br/>
<em>
string
</em>
</td>
<td>
<p>OperatingSystem refers to switch&rsquo;s operating system</p>
</td>
</tr>
<tr>
<td>
<code>version</code><br/>
<em>
string
</em>
</td>
<td>
<p>Version refers to the operating system version</p>
</td>
</tr>
<tr>
<td>
<code>asic</code><br/>
<em>
string
</em>
</td>
<td>
<p>ASIC refers to the switch&rsquo;s ASIC manufacturer</p>
</td>
</tr>
</tbody>
</table>
<h3 id="switch.onmetal.de/v1alpha1.State">State
(<code>string</code> alias)</h3>
<p>
(<em>Appears on:</em><a href="#switch.onmetal.de/v1alpha1.SwitchAssignmentStatus">SwitchAssignmentStatus</a>)
</p>
<div>
</div>
<table>
<thead>
<tr>
<th>Value</th>
<th>Description</th>
</tr>
</thead>
<tbody><tr><td><p>&#34;finished&#34;</p></td>
<td></td>
</tr><tr><td><p>&#34;pending&#34;</p></td>
<td></td>
</tr><tr><td><p>&#34;deleting&#34;</p></td>
<td></td>
</tr></tbody>
</table>
<h3 id="switch.onmetal.de/v1alpha1.SubnetSpec">SubnetSpec
</h3>
<p>
(<em>Appears on:</em><a href="#switch.onmetal.de/v1alpha1.SwitchStatus">SwitchStatus</a>)
</p>
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
<a href="#switch.onmetal.de/v1alpha1.RegionSpec">
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
<code>resourceReference</code><br/>
<em>
<a href="#switch.onmetal.de/v1alpha1.ResourceReferenceSpec">
ResourceReferenceSpec
</a>
</em>
</td>
<td>
<p>ResourceReference refers to the related resource definition</p>
</td>
</tr>
</tbody>
</table>
<h3 id="switch.onmetal.de/v1alpha1.Switch">Switch
</h3>
<div>
<p>Switch is the Schema for the switches API</p>
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
<a href="#switch.onmetal.de/v1alpha1.SwitchSpec">
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
<code>hostname</code><br/>
<em>
string
</em>
</td>
<td>
<p>Hostname refers to switch hostname</p>
</td>
</tr>
<tr>
<td>
<code>chassis</code><br/>
<em>
<a href="#switch.onmetal.de/v1alpha1.ChassisSpec">
ChassisSpec
</a>
</em>
</td>
<td>
<p>Chassis refers to baremetal box info</p>
</td>
</tr>
<tr>
<td>
<code>softwarePlatform</code><br/>
<em>
<a href="#switch.onmetal.de/v1alpha1.SoftwarePlatformSpec">
SoftwarePlatformSpec
</a>
</em>
</td>
<td>
<p>SoftwarePlatform refers to software info</p>
</td>
</tr>
<tr>
<td>
<code>location</code><br/>
<em>
<a href="#switch.onmetal.de/v1alpha1.LocationSpec">
LocationSpec
</a>
</em>
</td>
<td>
<p>Location refers to the switch&rsquo;s location</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#switch.onmetal.de/v1alpha1.SwitchStatus">
SwitchStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="switch.onmetal.de/v1alpha1.SwitchAssignment">SwitchAssignment
</h3>
<div>
<p>SwitchAssignment is the Schema for the switch assignments API</p>
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
<a href="#switch.onmetal.de/v1alpha1.SwitchAssignmentSpec">
SwitchAssignmentSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>chassisId</code><br/>
<em>
string
</em>
</td>
<td>
<p>ChassisID refers to switch chassis id</p>
</td>
</tr>
<tr>
<td>
<code>region</code><br/>
<em>
<a href="#switch.onmetal.de/v1alpha1.RegionSpec">
RegionSpec
</a>
</em>
</td>
<td>
<p>Region refers to the switch&rsquo;s region</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#switch.onmetal.de/v1alpha1.SwitchAssignmentStatus">
SwitchAssignmentStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="switch.onmetal.de/v1alpha1.SwitchAssignmentSpec">SwitchAssignmentSpec
</h3>
<p>
(<em>Appears on:</em><a href="#switch.onmetal.de/v1alpha1.SwitchAssignment">SwitchAssignment</a>)
</p>
<div>
<p>SwitchAssignmentSpec defines the desired state of SwitchAssignment</p>
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
<p>ChassisID refers to switch chassis id</p>
</td>
</tr>
<tr>
<td>
<code>region</code><br/>
<em>
<a href="#switch.onmetal.de/v1alpha1.RegionSpec">
RegionSpec
</a>
</em>
</td>
<td>
<p>Region refers to the switch&rsquo;s region</p>
</td>
</tr>
</tbody>
</table>
<h3 id="switch.onmetal.de/v1alpha1.SwitchAssignmentStatus">SwitchAssignmentStatus
</h3>
<p>
(<em>Appears on:</em><a href="#switch.onmetal.de/v1alpha1.SwitchAssignment">SwitchAssignment</a>)
</p>
<div>
<p>SwitchAssignmentStatus defines the observed state of SwitchAssignment</p>
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
<code>state</code><br/>
<em>
<a href="#switch.onmetal.de/v1alpha1.State">
State
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>switch</code><br/>
<em>
<a href="#switch.onmetal.de/v1alpha1.LinkedSwitchSpec">
LinkedSwitchSpec
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="switch.onmetal.de/v1alpha1.SwitchConfState">SwitchConfState
(<code>string</code> alias)</h3>
<p>
(<em>Appears on:</em><a href="#switch.onmetal.de/v1alpha1.ConfigurationSpec">ConfigurationSpec</a>)
</p>
<div>
</div>
<table>
<thead>
<tr>
<th>Value</th>
<th>Description</th>
</tr>
</thead>
<tbody><tr><td><p>&#34;applied&#34;</p></td>
<td></td>
</tr><tr><td><p>&#34;in progress&#34;</p></td>
<td></td>
</tr><tr><td><p>&#34;initial&#34;</p></td>
<td></td>
</tr><tr><td><p>&#34;pending&#34;</p></td>
<td></td>
</tr></tbody>
</table>
<h3 id="switch.onmetal.de/v1alpha1.SwitchRole">SwitchRole
(<code>string</code> alias)</h3>
<p>
(<em>Appears on:</em><a href="#switch.onmetal.de/v1alpha1.SwitchStatus">SwitchStatus</a>)
</p>
<div>
</div>
<table>
<thead>
<tr>
<th>Value</th>
<th>Description</th>
</tr>
</thead>
<tbody><tr><td><p>&#34;leaf&#34;</p></td>
<td></td>
</tr><tr><td><p>&#34;spine&#34;</p></td>
<td></td>
</tr></tbody>
</table>
<h3 id="switch.onmetal.de/v1alpha1.SwitchSpec">SwitchSpec
</h3>
<p>
(<em>Appears on:</em><a href="#switch.onmetal.de/v1alpha1.Switch">Switch</a>)
</p>
<div>
<p>SwitchSpec defines the desired state of Switch</p>
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
<code>hostname</code><br/>
<em>
string
</em>
</td>
<td>
<p>Hostname refers to switch hostname</p>
</td>
</tr>
<tr>
<td>
<code>chassis</code><br/>
<em>
<a href="#switch.onmetal.de/v1alpha1.ChassisSpec">
ChassisSpec
</a>
</em>
</td>
<td>
<p>Chassis refers to baremetal box info</p>
</td>
</tr>
<tr>
<td>
<code>softwarePlatform</code><br/>
<em>
<a href="#switch.onmetal.de/v1alpha1.SoftwarePlatformSpec">
SoftwarePlatformSpec
</a>
</em>
</td>
<td>
<p>SoftwarePlatform refers to software info</p>
</td>
</tr>
<tr>
<td>
<code>location</code><br/>
<em>
<a href="#switch.onmetal.de/v1alpha1.LocationSpec">
LocationSpec
</a>
</em>
</td>
<td>
<p>Location refers to the switch&rsquo;s location</p>
</td>
</tr>
</tbody>
</table>
<h3 id="switch.onmetal.de/v1alpha1.SwitchState">SwitchState
(<code>string</code> alias)</h3>
<p>
(<em>Appears on:</em><a href="#switch.onmetal.de/v1alpha1.SwitchStatus">SwitchStatus</a>)
</p>
<div>
</div>
<table>
<thead>
<tr>
<th>Value</th>
<th>Description</th>
</tr>
</thead>
<tbody><tr><td><p>&#34;in progress&#34;</p></td>
<td></td>
</tr><tr><td><p>&#34;initial&#34;</p></td>
<td></td>
</tr><tr><td><p>&#34;ready&#34;</p></td>
<td></td>
</tr></tbody>
</table>
<h3 id="switch.onmetal.de/v1alpha1.SwitchStatus">SwitchStatus
</h3>
<p>
(<em>Appears on:</em><a href="#switch.onmetal.de/v1alpha1.Switch">Switch</a>)
</p>
<div>
<p>SwitchStatus defines the observed state of Switch</p>
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
<code>totalPorts</code><br/>
<em>
uint16
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
uint16
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
<a href="#switch.onmetal.de/v1alpha1.SwitchRole">
SwitchRole
</a>
</em>
</td>
<td>
<p>Role refers to switch&rsquo;s role</p>
</td>
</tr>
<tr>
<td>
<code>connectionLevel</code><br/>
<em>
byte
</em>
</td>
<td>
<p>ConnectionLevel refers to switch&rsquo;s current position in connection hierarchy</p>
</td>
</tr>
<tr>
<td>
<code>interfaces</code><br/>
<em>
<a href="#switch.onmetal.de/v1alpha1.*..InterfaceSpec">
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
<code>subnetV4</code><br/>
<em>
<a href="#switch.onmetal.de/v1alpha1.SubnetSpec">
SubnetSpec
</a>
</em>
</td>
<td>
<p>SubnetV4 refers to the switch&rsquo;s south IPv4 subnet</p>
</td>
</tr>
<tr>
<td>
<code>subnetV6</code><br/>
<em>
<a href="#switch.onmetal.de/v1alpha1.SubnetSpec">
SubnetSpec
</a>
</em>
</td>
<td>
<p>SubnetV6 refers to the switch&rsquo;s south IPv6 subnet</p>
</td>
</tr>
<tr>
<td>
<code>loopbackV4</code><br/>
<em>
<a href="#switch.onmetal.de/v1alpha1.IPAddressSpec">
IPAddressSpec
</a>
</em>
</td>
<td>
<p>LoopbackV4 refers to the switch&rsquo;s loopback IPv4 address</p>
</td>
</tr>
<tr>
<td>
<code>loopbackV6</code><br/>
<em>
<a href="#switch.onmetal.de/v1alpha1.IPAddressSpec">
IPAddressSpec
</a>
</em>
</td>
<td>
<p>LoopbackV6 refers to the switch&rsquo;s loopback IPv6 address</p>
</td>
</tr>
<tr>
<td>
<code>configuration</code><br/>
<em>
<a href="#switch.onmetal.de/v1alpha1.ConfigurationSpec">
ConfigurationSpec
</a>
</em>
</td>
<td>
<p>Configuration refers to how switch&rsquo;s configuration manager is defined</p>
</td>
</tr>
<tr>
<td>
<code>state</code><br/>
<em>
<a href="#switch.onmetal.de/v1alpha1.SwitchState">
SwitchState
</a>
</em>
</td>
<td>
<p>State refers to current switch&rsquo;s processing state</p>
</td>
</tr>
</tbody>
</table>
<hr/>
<p><em>
Generated with <code>gen-crd-api-reference-docs</code>
on git commit <code>7ccba2d</code>.
</em></p>
