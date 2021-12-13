<p>Packages:</p>
<ul>
<li>
<a href="#machine.onmetal.de%2fv1alpha1">machine.onmetal.de/v1alpha1</a>
</li>
</ul>
<h2 id="machine.onmetal.de/v1alpha1">machine.onmetal.de/v1alpha1</h2>
Resource Types:
<ul></ul>
<h3 id="machine.onmetal.de/v1alpha1.Action">Action
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha1.MachineSpec">MachineSpec</a>)
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
<code>power_state</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>PowerState - defines desired machine power state</p>
</td>
</tr>
</tbody>
</table>
<h3 id="machine.onmetal.de/v1alpha1.Identity">Identity
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha1.MachineSpec">MachineSpec</a>)
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
<a href="#machine.onmetal.de/v1alpha1.Internal">
[]Internal
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="machine.onmetal.de/v1alpha1.Interface">Interface
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha1.MachineStatus">MachineStatus</a>)
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
<code>switch_uuid</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>SwitchUUID - defines unique switch identification</p>
</td>
</tr>
<tr>
<td>
<code>ipv4</code><br/>
<em>
string
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
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>IPv6 - defines machine IPv6 address</p>
</td>
</tr>
<tr>
<td>
<code>lldp_system_name</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>LLDPSystemName - defines switch name obtained from Link Layer Discovery Protocol -
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
<code>lane</code><br/>
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
<h3 id="machine.onmetal.de/v1alpha1.Internal">Internal
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha1.Identity">Identity</a>)
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
<h3 id="machine.onmetal.de/v1alpha1.Location">Location
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha1.MachineSpec">MachineSpec</a>)
</p>
<div>
<p>Location - defines information about place where machines are stored.</p>
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
<code>datacenter</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Datacenter - name of building where machine lies</p>
</td>
</tr>
<tr>
<td>
<code>data_hall</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>DataHall - name of room in Datacenter where machine lies</p>
</td>
</tr>
<tr>
<td>
<code>shelf</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Shelf - defines place for server in Datacenter (an alternative name of Rack)</p>
</td>
</tr>
<tr>
<td>
<code>slot</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Slot - defines switch location in rack (an alternative name for Row)</p>
</td>
</tr>
<tr>
<td>
<code>hu</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>HU - is a unit of measure defined 44.45 mm</p>
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
<em>(Optional)</em>
<p>Row - switch location in rack</p>
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
<em>(Optional)</em>
<p>Rack - is a place for server in DataCenter</p>
</td>
</tr>
</tbody>
</table>
<h3 id="machine.onmetal.de/v1alpha1.Machine">Machine
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
<a href="#machine.onmetal.de/v1alpha1.MachineSpec">
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
<code>location</code><br/>
<em>
<a href="#machine.onmetal.de/v1alpha1.Location">
Location
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Location - defines machine location in datacenter</p>
</td>
</tr>
<tr>
<td>
<code>action</code><br/>
<em>
<a href="#machine.onmetal.de/v1alpha1.Action">
Action
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Action - defines desired operation on machine</p>
</td>
</tr>
<tr>
<td>
<code>identity</code><br/>
<em>
<a href="#machine.onmetal.de/v1alpha1.Identity">
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
<code>scan_ports</code><br/>
<em>
bool
</em>
</td>
<td>
<p>ScanPorts - trigger manual port scan</p>
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
<a href="#machine.onmetal.de/v1alpha1.MachineStatus">
MachineStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="machine.onmetal.de/v1alpha1.MachineSpec">MachineSpec
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha1.Machine">Machine</a>)
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
<code>location</code><br/>
<em>
<a href="#machine.onmetal.de/v1alpha1.Location">
Location
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Location - defines machine location in datacenter</p>
</td>
</tr>
<tr>
<td>
<code>action</code><br/>
<em>
<a href="#machine.onmetal.de/v1alpha1.Action">
Action
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Action - defines desired operation on machine</p>
</td>
</tr>
<tr>
<td>
<code>identity</code><br/>
<em>
<a href="#machine.onmetal.de/v1alpha1.Identity">
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
<code>scan_ports</code><br/>
<em>
bool
</em>
</td>
<td>
<p>ScanPorts - trigger manual port scan</p>
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
<h3 id="machine.onmetal.de/v1alpha1.MachineStatus">MachineStatus
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha1.Machine">Machine</a>)
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
<a href="#machine.onmetal.de/v1alpha1.Interface">
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
string
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
<a href="#machine.onmetal.de/v1alpha1.Network">
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
<tr>
<td>
<code>inventory</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Inventory - defines status, if Inventory is presented or not</p>
</td>
</tr>
<tr>
<td>
<code>oob</code><br/>
<em>
bool
</em>
</td>
<td>
<p>OOB define status, OOB is presented or not</p>
</td>
</tr>
</tbody>
</table>
<h3 id="machine.onmetal.de/v1alpha1.Network">Network
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha1.MachineStatus">MachineStatus</a>)
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
Available values: &ldquo;Single&rdquo;, &ldquo;High Availability&rdquo; or &ldquo;None&rdquo;</p>
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
<hr/>
<p><em>
Generated with <code>gen-crd-api-reference-docs</code>
on git commit <code>172f4e7</code>.
</em></p>
