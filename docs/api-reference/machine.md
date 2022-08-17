<p>Packages:</p>
<ul>
<li>
<a href="#machine.onmetal.de%2fv1alpha2">machine.onmetal.de/v1alpha2</a>
</li>
</ul>
<h2 id="machine.onmetal.de/v1alpha2">machine.onmetal.de/v1alpha2</h2>
Resource Types:
<ul></ul>
<h3 id="machine.onmetal.de/v1alpha2.EFIVar">EFIVar
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha2.MachineAssignmentSpec">MachineAssignmentSpec</a>)
</p>
<div>
<p>EFIVar is a variable to pass to EFI while booting up.</p>
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
<code>uuid</code><br/>
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
<h3 id="machine.onmetal.de/v1alpha2.IPAddress">IPAddress
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha2.NetworkInterfaces">NetworkInterfaces</a>)
</p>
<div>
<p>IP is an IP address.</p>
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
inet.af/netaddr.IP
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="machine.onmetal.de/v1alpha2.IPAddressSpec">IPAddressSpec
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha2.Interface">Interface</a>)
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
<h3 id="machine.onmetal.de/v1alpha2.MachineAssignment">MachineAssignment
</h3>
<div>
<p>MachineAssignment is the Schema for the requests API.</p>
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
<a href="#machine.onmetal.de/v1alpha2.MachineAssignmentSpec">
MachineAssignmentSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>tolerations</code><br/>
<em>
<a href="#machine.onmetal.de/v1alpha2.Toleration">
[]Toleration
</a>
</em>
</td>
<td>
<p>Tolerations define tolerations the Machine has. Only MachinePools whose taints
covered by Tolerations will be considered to run the Machine.</p>
</td>
</tr>
<tr>
<td>
<code>machineClass</code><br/>
<em>
<a href="https://v1-21.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#localobjectreference-v1-core">
Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<p>MachineClass is a reference to the machine class/flavor of the machine.</p>
</td>
</tr>
<tr>
<td>
<code>image</code><br/>
<em>
string
</em>
</td>
<td>
<p>Image is the URL providing the operating system image of the machine.</p>
</td>
</tr>
<tr>
<td>
<code>networkInterfaces</code><br/>
<em>
<a href="#machine.onmetal.de/v1alpha2.NetworkInterfaces">
[]NetworkInterfaces
</a>
</em>
</td>
<td>
<p>Interfaces define a list of network interfaces present on the machine</p>
</td>
</tr>
<tr>
<td>
<code>volumes</code><br/>
<em>
<a href="#machine.onmetal.de/v1alpha2.Volume">
[]Volume
</a>
</em>
</td>
<td>
<p>Volumes are volumes attached to this machine.</p>
</td>
</tr>
<tr>
<td>
<code>ignition</code><br/>
<em>
<a href="#machine.onmetal.de/v1alpha2.ObjectSelector">
ObjectSelector
</a>
</em>
</td>
<td>
<p>Ignition is a reference to a config map containing the ignition YAML for the machine to boot up.
If key is empty, DefaultIgnitionKey will be used as fallback.</p>
</td>
</tr>
<tr>
<td>
<code>efiVars</code><br/>
<em>
<a href="#machine.onmetal.de/v1alpha2.EFIVar">
[]EFIVar
</a>
</em>
</td>
<td>
<p>EFIVars are variables to pass to EFI while booting up.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#machine.onmetal.de/v1alpha2.MachineAssignmentStatus">
MachineAssignmentStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="machine.onmetal.de/v1alpha2.MachineAssignmentSpec">MachineAssignmentSpec
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha2.MachineAssignment">MachineAssignment</a>)
</p>
<div>
<p>MachineAssignmentSpec defines the desired state of Request.</p>
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
<code>tolerations</code><br/>
<em>
<a href="#machine.onmetal.de/v1alpha2.Toleration">
[]Toleration
</a>
</em>
</td>
<td>
<p>Tolerations define tolerations the Machine has. Only MachinePools whose taints
covered by Tolerations will be considered to run the Machine.</p>
</td>
</tr>
<tr>
<td>
<code>machineClass</code><br/>
<em>
<a href="https://v1-21.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#localobjectreference-v1-core">
Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<p>MachineClass is a reference to the machine class/flavor of the machine.</p>
</td>
</tr>
<tr>
<td>
<code>image</code><br/>
<em>
string
</em>
</td>
<td>
<p>Image is the URL providing the operating system image of the machine.</p>
</td>
</tr>
<tr>
<td>
<code>networkInterfaces</code><br/>
<em>
<a href="#machine.onmetal.de/v1alpha2.NetworkInterfaces">
[]NetworkInterfaces
</a>
</em>
</td>
<td>
<p>Interfaces define a list of network interfaces present on the machine</p>
</td>
</tr>
<tr>
<td>
<code>volumes</code><br/>
<em>
<a href="#machine.onmetal.de/v1alpha2.Volume">
[]Volume
</a>
</em>
</td>
<td>
<p>Volumes are volumes attached to this machine.</p>
</td>
</tr>
<tr>
<td>
<code>ignition</code><br/>
<em>
<a href="#machine.onmetal.de/v1alpha2.ObjectSelector">
ObjectSelector
</a>
</em>
</td>
<td>
<p>Ignition is a reference to a config map containing the ignition YAML for the machine to boot up.
If key is empty, DefaultIgnitionKey will be used as fallback.</p>
</td>
</tr>
<tr>
<td>
<code>efiVars</code><br/>
<em>
<a href="#machine.onmetal.de/v1alpha2.EFIVar">
[]EFIVar
</a>
</em>
</td>
<td>
<p>EFIVars are variables to pass to EFI while booting up.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="machine.onmetal.de/v1alpha2.MachineAssignmentStatus">MachineAssignmentStatus
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha2.MachineAssignment">MachineAssignment</a>)
</p>
<div>
<p>MachineAssignmentStatus defines the observed state of Request.</p>
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
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>machineRef</code><br/>
<em>
<a href="#machine.onmetal.de/v1alpha2.ResourceReference">
ResourceReference
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>computeMachineRef</code><br/>
<em>
<a href="#machine.onmetal.de/v1alpha2.ResourceReference">
ResourceReference
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
<code>reservation</code><br/>
<em>
<a href="#machine.onmetal.de/v1alpha2.Reservation">
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
<h3 id="machine.onmetal.de/v1alpha2.NetworkInterfaces">NetworkInterfaces
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha2.MachineAssignmentSpec">MachineAssignmentSpec</a>)
</p>
<div>
<p>Interface is the definition of a single interface.</p>
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
<p>Name is the name of the interface</p>
</td>
</tr>
<tr>
<td>
<code>target</code><br/>
<em>
<a href="https://v1-21.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#localobjectreference-v1-core">
Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<p>Target is the referenced resource of this interface.</p>
</td>
</tr>
<tr>
<td>
<code>priority</code><br/>
<em>
int32
</em>
</td>
<td>
<p>Priority is the priority level of this interface</p>
</td>
</tr>
<tr>
<td>
<code>ip</code><br/>
<em>
<a href="#machine.onmetal.de/v1alpha2.IPAddress">
IPAddress
</a>
</em>
</td>
<td>
<p>IP specifies a concrete IP address which should be allocated from a Subnet</p>
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
<p>ObjectReference - defines object reference status and additional information.</p>
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
<h3 id="machine.onmetal.de/v1alpha2.ObjectSelector">ObjectSelector
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha2.MachineAssignmentSpec">MachineAssignmentSpec</a>)
</p>
<div>
<p>ObjectSelector is a reference to a specific &lsquo;key&rsquo; within a ConfigMap resource.
In some instances, <code>key</code> is a required field.</p>
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
<code>LocalObjectReference</code><br/>
<em>
<a href="https://v1-21.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#localobjectreference-v1-core">
Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<p>
(Members of <code>LocalObjectReference</code> are embedded into this type.)
</p>
<p>The name of the ConfigMap resource being referred to.</p>
</td>
</tr>
<tr>
<td>
<code>key</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>The key of the entry in the ConfigMap resource&rsquo;s <code>data</code> field to be used.
Some instances of this field may be defaulted, in others it may be
required.</p>
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
<h3 id="machine.onmetal.de/v1alpha2.Reservation">Reservation
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha2.MachineStatus">MachineStatus</a>)
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
<code>reference</code><br/>
<em>
<a href="#machine.onmetal.de/v1alpha2.ResourceReference">
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
<h3 id="machine.onmetal.de/v1alpha2.ResourceReference">ResourceReference
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha2.IPAddressSpec">IPAddressSpec</a>, <a href="#machine.onmetal.de/v1alpha2.Interface">Interface</a>, <a href="#machine.onmetal.de/v1alpha2.MachineAssignmentStatus">MachineAssignmentStatus</a>, <a href="#machine.onmetal.de/v1alpha2.ObjectReference">ObjectReference</a>, <a href="#machine.onmetal.de/v1alpha2.Peer">Peer</a>, <a href="#machine.onmetal.de/v1alpha2.Reservation">Reservation</a>)
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
<em>(Optional)</em>
<p>Value - corresponding to the taint key.</p>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#machine.onmetal.de/v1alpha2.TaintEffect">
TaintEffect
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
<h3 id="machine.onmetal.de/v1alpha2.TaintEffect">TaintEffect
(<code>string</code> alias)</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha2.Taint">Taint</a>, <a href="#machine.onmetal.de/v1alpha2.Toleration">Toleration</a>)
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
<td><p>When Machine taint effect is Error it&rsquo;s impossible to order machine. And it requires to run stresstest.</p>
</td>
</tr><tr><td><p>&#34;NoSchedule&#34;</p></td>
<td><p>When Machine taint effect is NoSchedule.</p>
</td>
</tr><tr><td><p>&#34;NotAvailable&#34;</p></td>
<td><p>When Machine taint effect is NotAvailable that&rsquo;s mean that Inventory or OOB not exist.</p>
</td>
</tr><tr><td><p>&#34;Suspended&#34;</p></td>
<td><p>When Machine taint effect is Suspended.</p>
</td>
</tr></tbody>
</table>
<h3 id="machine.onmetal.de/v1alpha2.Toleration">Toleration
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha2.MachineAssignmentSpec">MachineAssignmentSpec</a>)
</p>
<div>
<p>The resource this Toleration is attached to tolerates any taint that matches
the triple <key,value,effect> using the matching operator <operator>.</p>
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
<p>Key is the taint key that the toleration applies to. Empty means match all taint keys.
If the key is empty, operator must be Exists; this combination means to match all values and all keys.</p>
</td>
</tr>
<tr>
<td>
<code>operator</code><br/>
<em>
<a href="#machine.onmetal.de/v1alpha2.TolerationOperator">
TolerationOperator
</a>
</em>
</td>
<td>
<p>Operator represents a key&rsquo;s relationship to the value.
Valid operators are Exists and Equal. Defaults to Equal.
Exist is equivalent to wildcard for value, so that a resource can
tolerate all taints of a particular category.</p>
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
<p>Value is the taint value the toleration matches to.
If the operator is Exists, the value should be empty, otherwise just a regular string.</p>
</td>
</tr>
<tr>
<td>
<code>effect</code><br/>
<em>
<a href="#machine.onmetal.de/v1alpha2.TaintEffect">
TaintEffect
</a>
</em>
</td>
<td>
<p>Effect indicates the taint effect to match. Empty means match all taint effects.
When specified, allowed values are NoSchedule.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="machine.onmetal.de/v1alpha2.TolerationOperator">TolerationOperator
(<code>string</code> alias)</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha2.Toleration">Toleration</a>)
</p>
<div>
<p>A toleration operator is the set of operators that can be used in a toleration.</p>
</div>
<table>
<thead>
<tr>
<th>Value</th>
<th>Description</th>
</tr>
</thead>
<tbody><tr><td><p>&#34;Equal&#34;</p></td>
<td></td>
</tr><tr><td><p>&#34;Exists&#34;</p></td>
<td></td>
</tr></tbody>
</table>
<h3 id="machine.onmetal.de/v1alpha2.Volume">Volume
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha2.MachineAssignmentSpec">MachineAssignmentSpec</a>)
</p>
<div>
<p>Volume defines a volume attachment of a machine.</p>
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
<p>Name is the name of the VolumeAttachment</p>
</td>
</tr>
<tr>
<td>
<code>priority</code><br/>
<em>
int32
</em>
</td>
<td>
<p>Priority is the OS priority of the volume.</p>
</td>
</tr>
<tr>
<td>
<code>VolumeSource</code><br/>
<em>
<a href="#machine.onmetal.de/v1alpha2.VolumeSource">
VolumeSource
</a>
</em>
</td>
<td>
<p>
(Members of <code>VolumeSource</code> are embedded into this type.)
</p>
<p>VolumeAttachmentSource is the source where the storage for the VolumeAttachment resides at.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="machine.onmetal.de/v1alpha2.VolumeClaimSource">VolumeClaimSource
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha2.VolumeSource">VolumeSource</a>)
</p>
<div>
<p>VolumeClaimSource references a VolumeClaim as VolumeAttachment source.</p>
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
<code>ref</code><br/>
<em>
<a href="https://v1-21.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#localobjectreference-v1-core">
Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<p>Ref is a reference to the VolumeClaim.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="machine.onmetal.de/v1alpha2.VolumeSource">VolumeSource
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha2.Volume">Volume</a>)
</p>
<div>
<p>VolumeSource specifies the source to use for a VolumeAttachment.</p>
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
<code>volumeClaim</code><br/>
<em>
<a href="#machine.onmetal.de/v1alpha2.VolumeClaimSource">
VolumeClaimSource
</a>
</em>
</td>
<td>
<p>VolumeClaim instructs the VolumeAttachment to use a VolumeClaim as source for the attachment.</p>
</td>
</tr>
</tbody>
</table>
<hr/>
<p><em>
Generated with <code>gen-crd-api-reference-docs</code>
on git commit <code>99a336f</code>.
</em></p>
