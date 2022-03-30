<p>Packages:</p>
<ul>
<li>
<a href="#machine.onmetal.de%2fv1alpha2">machine.onmetal.de/v1alpha2</a>
</li>
</ul>
<h2 id="machine.onmetal.de/v1alpha2">machine.onmetal.de/v1alpha2</h2>
Resource Types:
<ul></ul>
<h3 id="machine.onmetal.de/v1alpha2.IPAddressSpec">IPAddressSpec
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha2.Interface">Interface</a>)
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
<p>Address refers to the ip address value</p>
</td>
</tr>
<tr>
<td>
<code>resource_reference</code><br/>
<em>
<a href="#machine.onmetal.de/v1alpha2.ResourceReference">
ResourceReference
</a>
</em>
</td>
<td>
<p>ResourceReference refers to the related resource definition</p>
</td>
</tr>
</tbody>
</table>
<h3 id="machine.onmetal.de/v1alpha2.Identity">Identity
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha2.MachineSpec">MachineSpec</a>)
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
<a href="#machine.onmetal.de/v1alpha2.Internal">
[]Internal
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="machine.onmetal.de/v1alpha2.Interface">Interface
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha2.MachineStatus">MachineStatus</a>)
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
<a href="#machine.onmetal.de/v1alpha2.ResourceReference">
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
<code>ipv4</code><br/>
<em>
<a href="#machine.onmetal.de/v1alpha2.IPAddressSpec">
IPAddressSpec
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
<code>ipv6</code><br/>
<em>
<a href="#machine.onmetal.de/v1alpha2.IPAddressSpec">
IPAddressSpec
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>IPv6 - defines machine IPv6 address</p>
</td>
</tr>
<tr>
<td>
<code>peer</code><br/>
<em>
<a href="#machine.onmetal.de/v1alpha2.Peer">
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
byte
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
<h3 id="machine.onmetal.de/v1alpha2.Internal">Internal
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha2.Identity">Identity</a>)
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
<h3 id="machine.onmetal.de/v1alpha2.Machine">Machine
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
<a href="#machine.onmetal.de/v1alpha2.MachineSpec">
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
<code>taints</code><br/>
<em>
<a href="#machine.onmetal.de/v1alpha2.Taint">
[]Taint
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Taints - defines list of Taint that applied on the Machine</p>
</td>
</tr>
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
<a href="#machine.onmetal.de/v1alpha2.Identity">
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
<a href="#machine.onmetal.de/v1alpha2.MachineStatus">
MachineStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="machine.onmetal.de/v1alpha2.MachineSpec">MachineSpec
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha2.Machine">Machine</a>)
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
<code>taints</code><br/>
<em>
<a href="#machine.onmetal.de/v1alpha2.Taint">
[]Taint
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Taints - defines list of Taint that applied on the Machine</p>
</td>
</tr>
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
<a href="#machine.onmetal.de/v1alpha2.Identity">
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
<h3 id="machine.onmetal.de/v1alpha2.MachineState">MachineState
(<code>string</code> alias)</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha2.MachineStatus">MachineStatus</a>)
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
<td><p>When State is <code>Healthy</code> Machine` is allowed to be booked.</p>
</td>
</tr><tr><td><p>&#34;Unhealthy&#34;</p></td>
<td><p>When State is <code>Unhealthy</code>` Machine isn&rsquo;t allowed to be booked.</p>
</td>
</tr></tbody>
</table>
<h3 id="machine.onmetal.de/v1alpha2.MachineStatus">MachineStatus
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha2.Machine">Machine</a>)
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
<code>interfaces</code><br/>
<em>
<a href="#machine.onmetal.de/v1alpha2.Interface">
[]Interface
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Interfaces - defines machine interfaces info</p>
</td>
</tr>
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
<a href="#machine.onmetal.de/v1alpha2.MachineState">
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
<a href="#machine.onmetal.de/v1alpha2.Network">
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
<code>oob</code><br/>
<em>
<a href="#machine.onmetal.de/v1alpha2.ObjectReference">
ObjectReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>OOB - defines status of OOB</p>
</td>
</tr>
<tr>
<td>
<code>inventory</code><br/>
<em>
<a href="#machine.onmetal.de/v1alpha2.ObjectReference">
ObjectReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Inventory - defines status of Inventory</p>
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
<h3 id="machine.onmetal.de/v1alpha2.Network">Network
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha2.MachineStatus">MachineStatus</a>)
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
<code>redundancy</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
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
<em>(Optional)</em>
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
<em>(Optional)</em>
<p>UnknownPorts - defines number of machine interface without info</p>
</td>
</tr>
</tbody>
</table>
<h3 id="machine.onmetal.de/v1alpha2.ObjectReference">ObjectReference
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha2.MachineStatus">MachineStatus</a>)
</p>
<div>
<p>ObjectReference - defines object reference status and additional information</p>
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
<code>exist</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Exist - defines where referenced object exist or not</p>
</td>
</tr>
<tr>
<td>
<code>reference</code><br/>
<em>
<a href="#machine.onmetal.de/v1alpha2.ResourceReference">
ResourceReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Reference - defines underlaying referenced object e.g. Inventory or OOB kind.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="machine.onmetal.de/v1alpha2.Peer">Peer
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha2.Interface">Interface</a>)
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
<code>lldp_chassi_id</code><br/>
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
<a href="#machine.onmetal.de/v1alpha2.ResourceReference">
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
<h3 id="machine.onmetal.de/v1alpha2.ResourceReference">ResourceReference
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha2.IPAddressSpec">IPAddressSpec</a>, <a href="#machine.onmetal.de/v1alpha2.Interface">Interface</a>, <a href="#machine.onmetal.de/v1alpha2.ObjectReference">ObjectReference</a>, <a href="#machine.onmetal.de/v1alpha2.Peer">Peer</a>)
</p>
<div>
<p>ResourceReference defines related resource info</p>
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
<h3 id="machine.onmetal.de/v1alpha2.Taint">Taint
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha2.MachineSpec">MachineSpec</a>)
</p>
<div>
<p>Taint represents taint that can be applied to the machine.
The machine this Taint is attached to has the &ldquo;effect&rdquo; on
any pod that does not tolerate the Taint.</p>
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
<code>key</code><br/>
<em>
string
</em>
</td>
<td>
<p>Key - applied to the machine.</p>
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
<p>Value - corresponding to the taint key.</p>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#machine.onmetal.de/v1alpha2.TaintStatus">
TaintStatus
</a>
</em>
</td>
<td>
<p>Effect - defines taint effect on the Machine.
Valid effects are NotAvailable and Suspended.</p>
</td>
</tr>
<tr>
<td>
<code>time_added</code><br/>
<em>
<a href="https://v1-21.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#time-v1-meta">
Kubernetes meta/v1.Time
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>TimeAdded represents the time at which the taint was added.
It is only written for NoExecute taints.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="machine.onmetal.de/v1alpha2.TaintStatus">TaintStatus
(<code>string</code> alias)</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha2.Taint">Taint</a>)
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
<tbody><tr><td><p>&#34;Error&#34;</p></td>
<td><p>When Machine status is Error it&rsquo;s impossible to order machine</p>
</td>
</tr><tr><td><p>&#34;NotAvailable&#34;</p></td>
<td><p>When Machine status is NotAvailable it&rsquo;s not possible to order it.</p>
</td>
</tr><tr><td><p>&#34;Suspended&#34;</p></td>
<td><p>When Machine status is Suspended that&rsquo;s meant that there is some issues
and need to run stress test.</p>
</td>
</tr></tbody>
</table>
<hr/>
<p><em>
Generated with <code>gen-crd-api-reference-docs</code>
on git commit <code>effbf2c</code>.
</em></p>
