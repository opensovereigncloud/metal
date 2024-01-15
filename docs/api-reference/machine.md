<p>Packages:</p>
<ul>
<li>
<a href="#metal.ironcore.dev%2fv1alpha3">metal.ironcore.dev/v1alpha3</a>
</li>
</ul>
<h2 id="metal.ironcore.dev/v1alpha3">metal.ironcore.dev/v1alpha3</h2>
Resource Types:
<ul><li>
<a href="#metal.ironcore.dev/v1alpha3.Machine">Machine</a>
</li></ul>
<h3 id="metal.ironcore.dev/v1alpha3.Machine">Machine
</h3>
<div>
<p>Machine - is the data structure for a Machine resource.
It contains an aggregated information from Inventory and OOB resources.</p>
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
metal.ironcore.dev/v1alpha3
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code><br/>
string
</td>
<td><code>Machine</code></td>
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
<a href="#metal.ironcore.dev/v1alpha3.MachineSpec">
MachineSpec
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
<em>(Optional)</em>
<p>Hostname - defines machine domain name</p>
</td>
</tr>
<tr>
<td>
<code>description</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Description - summary info about machine</p>
</td>
</tr>
<tr>
<td>
<code>identity</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha3.Identity">
Identity
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Identity - defines machine hardware info</p>
</td>
</tr>
<tr>
<td>
<code>inventory_requested</code><br/>
<em>
bool
</em>
</td>
<td>
<p>InventoryRequested - defines if inventory requested or not</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha3.MachineStatus">
MachineStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha3.Addresses">Addresses
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha3.Interface">Interface</a>)
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
<code>ipv4</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha3.IPAddressSpec">
[]IPAddressSpec
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>ipv6</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha3.IPAddressSpec">
[]IPAddressSpec
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha3.IPAddressSpec">IPAddressSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha3.Addresses">Addresses</a>, <a href="#metal.ironcore.dev/v1alpha3.LoopbackAddresses">LoopbackAddresses</a>)
</p>
<div>
<p>IPAddressSpec defines interface&rsquo;s ip address info.</p>
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
<code>-</code><br/>
<em>
net/netip.Prefix
</em>
</td>
<td>
<p>Address refers to the ip address value</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha3.Identity">Identity
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha3.MachineSpec">MachineSpec</a>)
</p>
<div>
<p>Identity - defines hardware information about machine.</p>
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
<code>sku</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>SKU - stock keeping unit. The label allows vendors automatically track the movement of inventory</p>
</td>
</tr>
<tr>
<td>
<code>serial_number</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>SerialNumber - unique machine number</p>
</td>
</tr>
<tr>
<td>
<code>asset</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>internal</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha3.Internal">
[]Internal
</a>
</em>
</td>
<td>
<p>Deprecated</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha3.Interface">Interface
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha3.Network">Network</a>)
</p>
<div>
<p>Interface - defines information about machine interfaces.</p>
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
<em>(Optional)</em>
<p>Name - machine interface name</p>
</td>
</tr>
<tr>
<td>
<code>switch_reference</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha3.ResourceReference">
ResourceReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>SwitchReference - defines unique switch identification</p>
</td>
</tr>
<tr>
<td>
<code>addresses</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha3.Addresses">
Addresses
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>IPv4 - defines machine IPv4 address</p>
</td>
</tr>
<tr>
<td>
<code>peer</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha3.Peer">
Peer
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Peer - defines lldp peer info.</p>
</td>
</tr>
<tr>
<td>
<code>lanes</code><br/>
<em>
uint32
</em>
</td>
<td>
<em>(Optional)</em>
<p>Lane - defines number of lines per interface</p>
</td>
</tr>
<tr>
<td>
<code>moved</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Moved  - defines if interface was reconnected to another switch or not</p>
</td>
</tr>
<tr>
<td>
<code>unknown</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Unknown - defines information availability about interface</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha3.Internal">Internal
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha3.Identity">Identity</a>)
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
<code>value</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha3.LoopbackAddresses">LoopbackAddresses
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha3.Network">Network</a>)
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
<code>ipv4</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha3.IPAddressSpec">
IPAddressSpec
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>ipv6</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha3.IPAddressSpec">
IPAddressSpec
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha3.MachineSpec">MachineSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha3.Machine">Machine</a>)
</p>
<div>
<p>MachineSpec - defines the desired spec of Machine.</p>
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
<em>(Optional)</em>
<p>Hostname - defines machine domain name</p>
</td>
</tr>
<tr>
<td>
<code>description</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Description - summary info about machine</p>
</td>
</tr>
<tr>
<td>
<code>identity</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha3.Identity">
Identity
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Identity - defines machine hardware info</p>
</td>
</tr>
<tr>
<td>
<code>inventory_requested</code><br/>
<em>
bool
</em>
</td>
<td>
<p>InventoryRequested - defines if inventory requested or not</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha3.MachineState">MachineState
(<code>string</code> alias)</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha3.MachineStatus">MachineStatus</a>)
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
<tbody><tr><td><p>&#34;Healthy&#34;</p></td>
<td><p>MachineStateHealthy - When State is <code>Healthy</code> Machine` is allowed to be booked.</p>
</td>
</tr><tr><td><p>&#34;Unhealthy&#34;</p></td>
<td><p>MachineStateUnhealthy - When State is <code>Unhealthy</code>` Machine isn&rsquo;t allowed to be booked.</p>
</td>
</tr></tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha3.MachineStatus">MachineStatus
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha3.Machine">Machine</a>)
</p>
<div>
<p>MachineStatus - defines machine aggregated info.</p>
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
<code>reboot</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Reboot - defines machine reboot status</p>
</td>
</tr>
<tr>
<td>
<code>health</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha3.MachineState">
MachineState
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Health - defines machine condition.
&ldquo;healthy&rdquo; if both OOB and Inventory are presented and &ldquo;unhealthy&rdquo; if one of them isn&rsquo;t</p>
</td>
</tr>
<tr>
<td>
<code>network</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha3.Network">
Network
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Network - defines machine network status</p>
</td>
</tr>
<tr>
<td>
<code>reservation</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha3.Reservation">
Reservation
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Reservation - defines machine reservation state and reference object.</p>
</td>
</tr>
<tr>
<td>
<code>orphaned</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Orphaned - defines machine condition whether OOB or Inventory is missing or not</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha3.Network">Network
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha3.MachineStatus">MachineStatus</a>)
</p>
<div>
<p>Network - defines machine network status.</p>
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
<code>asn</code><br/>
<em>
uint32
</em>
</td>
<td>
<p>ASN - defines calculated Autonomous system Number.</p>
</td>
</tr>
<tr>
<td>
<code>redundancy</code><br/>
<em>
string
</em>
</td>
<td>
<p>Redundancy - defines machine redundancy status.
Available values: &ldquo;Single&rdquo;, &ldquo;HighAvailability&rdquo; or &ldquo;None&rdquo;</p>
</td>
</tr>
<tr>
<td>
<code>ports</code><br/>
<em>
int
</em>
</td>
<td>
<p>Ports - defines number of machine ports</p>
</td>
</tr>
<tr>
<td>
<code>unknown_ports</code><br/>
<em>
int
</em>
</td>
<td>
<p>UnknownPorts - defines number of machine interface without info</p>
</td>
</tr>
<tr>
<td>
<code>interfaces</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha3.Interface">
[]Interface
</a>
</em>
</td>
<td>
<p>Interfaces - defines machine interfaces info</p>
</td>
</tr>
<tr>
<td>
<code>loopback_addresses</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha3.LoopbackAddresses">
LoopbackAddresses
</a>
</em>
</td>
<td>
<p>Loopbacks refers to the switch&rsquo;s loopback addresses</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha3.Peer">Peer
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha3.Interface">Interface</a>)
</p>
<div>
<p>Peer - contains machine neighbor information collected from LLDP.</p>
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
<code>lldp_system_name</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>LLDPSystemName - defines switch name obtained from Link Layer Discovery Protocol
layer 2 neighbor discovery protocol</p>
</td>
</tr>
<tr>
<td>
<code>lldp_chassis_id</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>LLDPChassisID - defines switch ID for chassis obtained from Link Layer Discovery Protocol</p>
</td>
</tr>
<tr>
<td>
<code>lldp_port_id</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>LLDPPortID - defines switch port ID obtained from Link Layer Discovery Protocol</p>
</td>
</tr>
<tr>
<td>
<code>lldp_port_description</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>LLDPPortDescription - defines switch definition obtained from Link Layer Discovery Protocol</p>
</td>
</tr>
<tr>
<td>
<code>resource_reference</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha3.ResourceReference">
ResourceReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>ResourceReference refers to the related resource definition</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha3.Reservation">Reservation
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha3.MachineStatus">MachineStatus</a>)
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
<code>status</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Status - defines Machine Order state provided by OOB Machine Resources</p>
</td>
</tr>
<tr>
<td>
<code>class</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Class - defines what class the mahchine was reserved under</p>
</td>
</tr>
<tr>
<td>
<code>reference</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha3.ResourceReference">
ResourceReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Reference - defines underlying referenced object.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha3.ResourceReference">ResourceReference
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha3.Interface">Interface</a>, <a href="#metal.ironcore.dev/v1alpha3.Peer">Peer</a>, <a href="#metal.ironcore.dev/v1alpha3.Reservation">Reservation</a>)
</p>
<div>
<p>ResourceReference defines related resource info.</p>
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
<em>(Optional)</em>
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
<em>(Optional)</em>
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
<em>(Optional)</em>
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
<em>(Optional)</em>
<p>Namespace refers to the resource namespace</p>
</td>
</tr>
</tbody>
</table>
<hr/>
<p><em>
Generated with <code>gen-crd-api-reference-docs</code>
on git commit <code>95c3af5</code>.
</em></p>
