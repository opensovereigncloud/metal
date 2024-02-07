<p>Packages:</p>
<ul>
<li>
<a href="#metal.ironcore.dev%2fv1alpha4">metal.ironcore.dev/v1alpha4</a>
</li>
</ul>
<h2 id="metal.ironcore.dev/v1alpha4">metal.ironcore.dev/v1alpha4</h2>
Resource Types:
<ul><li>
<a href="#metal.ironcore.dev/v1alpha4.Aggregate">Aggregate</a>
</li><li>
<a href="#metal.ironcore.dev/v1alpha4.Benchmark">Benchmark</a>
</li><li>
<a href="#metal.ironcore.dev/v1alpha4.Inventory">Inventory</a>
</li><li>
<a href="#metal.ironcore.dev/v1alpha4.Machine">Machine</a>
</li><li>
<a href="#metal.ironcore.dev/v1alpha4.NetworkSwitch">NetworkSwitch</a>
</li><li>
<a href="#metal.ironcore.dev/v1alpha4.Size">Size</a>
</li><li>
<a href="#metal.ironcore.dev/v1alpha4.SwitchConfig">SwitchConfig</a>
</li></ul>
<h3 id="metal.ironcore.dev/v1alpha4.Aggregate">Aggregate
</h3>
<div>
<p>Aggregate is the Schema for the aggregates API.</p>
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
metal.ironcore.dev/v1alpha4
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code><br/>
string
</td>
<td><code>Aggregate</code></td>
</tr>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">
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
<a href="#metal.ironcore.dev/v1alpha4.AggregateSpec">
AggregateSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>aggregates</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.AggregateItem">
[]AggregateItem
</a>
</em>
</td>
<td>
<p>Aggregates is a list of aggregates required to be computed</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.AggregateStatus">
AggregateStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha4.Benchmark">Benchmark
</h3>
<div>
<p>Benchmark is the Schema for the machines API.</p>
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
metal.ironcore.dev/v1alpha4
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code><br/>
string
</td>
<td><code>Benchmark</code></td>
</tr>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">
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
<a href="#metal.ironcore.dev/v1alpha4.BenchmarkSpec">
BenchmarkSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>benchmarks</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.Benchmarks">
map[string]./apis/metal/v1alpha4.Benchmarks
</a>
</em>
</td>
<td>
<p>Benchmarks is the collection of benchmarks.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.BenchmarkStatus">
BenchmarkStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha4.Inventory">Inventory
</h3>
<div>
<p>Inventory is the Schema for the inventories API.</p>
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
metal.ironcore.dev/v1alpha4
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code><br/>
string
</td>
<td><code>Inventory</code></td>
</tr>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">
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
<a href="#metal.ironcore.dev/v1alpha4.InventorySpec">
InventorySpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>system</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.SystemSpec">
SystemSpec
</a>
</em>
</td>
<td>
<p>System contains DMI system information</p>
</td>
</tr>
<tr>
<td>
<code>ipmis</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.IPMISpec">
[]IPMISpec
</a>
</em>
</td>
<td>
<p>IPMIs contains info about IPMI interfaces on the host</p>
</td>
</tr>
<tr>
<td>
<code>blocks</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.BlockSpec">
[]BlockSpec
</a>
</em>
</td>
<td>
<p>Blocks contains info about block devices on the host</p>
</td>
</tr>
<tr>
<td>
<code>memory</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.MemorySpec">
MemorySpec
</a>
</em>
</td>
<td>
<p>Memory contains info block devices on the host</p>
</td>
</tr>
<tr>
<td>
<code>cpus</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.CPUSpec">
[]CPUSpec
</a>
</em>
</td>
<td>
<p>CPUs contains info about cpus, cores and threads</p>
</td>
</tr>
<tr>
<td>
<code>numa</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.NumaSpec">
[]NumaSpec
</a>
</em>
</td>
<td>
<p>NUMA contains info about cpu/memory topology</p>
</td>
</tr>
<tr>
<td>
<code>pciDevices</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.PCIDeviceSpec">
[]PCIDeviceSpec
</a>
</em>
</td>
<td>
<p>PCIDevices contains info about devices accessible through</p>
</td>
</tr>
<tr>
<td>
<code>nics</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.NICSpec">
[]NICSpec
</a>
</em>
</td>
<td>
<p>NICs contains info about network interfaces and network discovery</p>
</td>
</tr>
<tr>
<td>
<code>virt</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.VirtSpec">
VirtSpec
</a>
</em>
</td>
<td>
<p>Virt is a virtualization detected on host</p>
</td>
</tr>
<tr>
<td>
<code>host</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.HostSpec">
HostSpec
</a>
</em>
</td>
<td>
<p>Host contains info about inventorying object</p>
</td>
</tr>
<tr>
<td>
<code>distro</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.DistroSpec">
DistroSpec
</a>
</em>
</td>
<td>
<p>Distro contains info about OS distro</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.InventoryStatus">
InventoryStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha4.Machine">Machine
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
metal.ironcore.dev/v1alpha4
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
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">
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
<a href="#metal.ironcore.dev/v1alpha4.MachineSpec">
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
<a href="#metal.ironcore.dev/v1alpha4.Identity">
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
<a href="#metal.ironcore.dev/v1alpha4.MachineStatus">
MachineStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha4.NetworkSwitch">NetworkSwitch
</h3>
<div>
<p>NetworkSwitch is the Schema for switches API.</p>
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
metal.ironcore.dev/v1alpha4
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code><br/>
string
</td>
<td><code>NetworkSwitch</code></td>
</tr>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">
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
<a href="#metal.ironcore.dev/v1alpha4.NetworkSwitchSpec">
NetworkSwitchSpec
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
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#localobjectreference-v1-core">
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
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#labelselector-v1-meta">
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
<p>Managed is a flag defining whether NetworkSwitch object would be processed during reconciliation</p>
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
<p>Cordon is a flag defining whether NetworkSwitch object is taken offline</p>
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
<p>TopSpine is a flag defining whether NetworkSwitch is a top-level spine switch</p>
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
<a href="#metal.ironcore.dev/v1alpha4.IPAMSpec">
IPAMSpec
</a>
</em>
</td>
<td>
<p>IPAM refers to selectors for subnets which will be used for NetworkSwitch object</p>
</td>
</tr>
<tr>
<td>
<code>interfaces</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.InterfacesSpec">
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
<a href="#metal.ironcore.dev/v1alpha4.NetworkSwitchStatus">
NetworkSwitchStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha4.Size">Size
</h3>
<div>
<p>Size is the Schema for the sizes API.</p>
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
metal.ironcore.dev/v1alpha4
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code><br/>
string
</td>
<td><code>Size</code></td>
</tr>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">
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
<a href="#metal.ironcore.dev/v1alpha4.SizeSpec">
SizeSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>constraints</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.ConstraintSpec">
[]ConstraintSpec
</a>
</em>
</td>
<td>
<p>Constraints is a list of selectors based on machine properties.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.SizeStatus">
SizeStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha4.SwitchConfig">SwitchConfig
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
metal.ironcore.dev/v1alpha4
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
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">
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
<a href="#metal.ironcore.dev/v1alpha4.SwitchConfigSpec">
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
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#labelselector-v1-meta">
Kubernetes meta/v1.LabelSelector
</a>
</em>
</td>
<td>
<p>Switches contains label selector to pick up NetworkSwitch objects</p>
</td>
</tr>
<tr>
<td>
<code>portsDefaults</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.PortParametersSpec">
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
<a href="#metal.ironcore.dev/v1alpha4.GeneralIPAMSpec">
GeneralIPAMSpec
</a>
</em>
</td>
<td>
<p>IPAM refers to selectors for subnets which will be used for NetworkSwitch object</p>
</td>
</tr>
<tr>
<td>
<code>routingConfigTemplate</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#localobjectreference-v1-core">
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
<a href="#metal.ironcore.dev/v1alpha4.SwitchConfigStatus">
SwitchConfigStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha4.AdditionalIPSpec">AdditionalIPSpec
</h3>
<div>
<p>AdditionalIPSpec defines IP address and selector for subnet where address should be reserved.</p>
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
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#labelselector-v1-meta">
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
<h3 id="metal.ironcore.dev/v1alpha4.AddressFamiliesMap">AddressFamiliesMap
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.GeneralIPAMSpec">GeneralIPAMSpec</a>)
</p>
<div>
<p>AddressFamiliesMap contains flags regarding what IP address families should be used.</p>
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
<h3 id="metal.ironcore.dev/v1alpha4.Addresses">Addresses
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.Interface">Interface</a>)
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
<a href="#metal.ironcore.dev/v1alpha4.IPAddrSpec">
[]IPAddrSpec
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
<a href="#metal.ironcore.dev/v1alpha4.IPAddrSpec">
[]IPAddrSpec
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha4.AggregateItem">AggregateItem
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.AggregateSpec">AggregateSpec</a>)
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
<code>sourcePath</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.JSONPath">
JSONPath
</a>
</em>
</td>
<td>
<p>SourcePath is a path in Inventory spec aggregate will be applied to</p>
</td>
</tr>
<tr>
<td>
<code>targetPath</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.JSONPath">
JSONPath
</a>
</em>
</td>
<td>
<p>TargetPath is a path in Inventory status <code>computed</code> field</p>
</td>
</tr>
<tr>
<td>
<code>aggregate</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.AggregateType">
AggregateType
</a>
</em>
</td>
<td>
<p>Aggregate defines whether collection values should be aggregated
for constraint checks, in case if path defines selector for collection</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha4.AggregateSpec">AggregateSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.Aggregate">Aggregate</a>)
</p>
<div>
<p>AggregateSpec defines the desired state of Aggregate.</p>
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
<code>aggregates</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.AggregateItem">
[]AggregateItem
</a>
</em>
</td>
<td>
<p>Aggregates is a list of aggregates required to be computed</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha4.AggregateStatus">AggregateStatus
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.Aggregate">Aggregate</a>)
</p>
<div>
<p>AggregateStatus defines the observed state of Aggregate.</p>
</div>
<h3 id="metal.ironcore.dev/v1alpha4.AggregateType">AggregateType
(<code>string</code> alias)</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.AggregateItem">AggregateItem</a>, <a href="#metal.ironcore.dev/v1alpha4.ConstraintSpec">ConstraintSpec</a>)
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
<tbody><tr><td><p>&#34;avg&#34;</p></td>
<td></td>
</tr><tr><td><p>&#34;count&#34;</p></td>
<td></td>
</tr><tr><td><p>&#34;max&#34;</p></td>
<td></td>
</tr><tr><td><p>&#34;min&#34;</p></td>
<td></td>
</tr><tr><td><p>&#34;sum&#34;</p></td>
<td></td>
</tr></tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha4.AggregationResults">AggregationResults
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.InventoryStatus">InventoryStatus</a>)
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
<code>-</code><br/>
<em>
map[string]interface{}
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha4.BenchmarkDeviation">BenchmarkDeviation
</h3>
<div>
<p>BenchmarkDeviation is a deviation between old value and the new one.</p>
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
<p>Name is the specific benchmark name. e.g. <code>fio-1k</code>.</p>
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
<p>Value is the exact result of specific benchmark.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha4.BenchmarkDeviations">BenchmarkDeviations
(<code>[]./apis/metal/v1alpha4.BenchmarkDeviation</code> alias)</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.BenchmarkStatus">BenchmarkStatus</a>)
</p>
<div>
</div>
<h3 id="metal.ironcore.dev/v1alpha4.BenchmarkResult">BenchmarkResult
</h3>
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
<p>Name is the specific benchmark name. e.g. <code>fio-1k</code>.</p>
</td>
</tr>
<tr>
<td>
<code>value</code><br/>
<em>
uint64
</em>
</td>
<td>
<p>Value is the exact result of specific benchmark.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha4.BenchmarkSpec">BenchmarkSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.Benchmark">Benchmark</a>)
</p>
<div>
<p>BenchmarkSpec contains machine benchmark results.</p>
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
<code>benchmarks</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.Benchmarks">
map[string]./apis/metal/v1alpha4.Benchmarks
</a>
</em>
</td>
<td>
<p>Benchmarks is the collection of benchmarks.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha4.BenchmarkStatus">BenchmarkStatus
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.Benchmark">Benchmark</a>)
</p>
<div>
<p>BenchmarkStatus contains machine benchmarks deviations.</p>
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
<code>machine_deviation</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.BenchmarkDeviations">
map[string]./apis/metal/v1alpha4.BenchmarkDeviations
</a>
</em>
</td>
<td>
<p>BenchmarkDeviations shows the difference between last and current benchmark results.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha4.Benchmarks">Benchmarks
(<code>[]./apis/metal/v1alpha4.BenchmarkResult</code> alias)</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.BenchmarkSpec">BenchmarkSpec</a>)
</p>
<div>
</div>
<h3 id="metal.ironcore.dev/v1alpha4.BlockSpec">BlockSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.InventorySpec">InventorySpec</a>)
</p>
<div>
<p>BlockSpec contains info about block device.</p>
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
<p>Name is a name of the device registered by Linux Kernel</p>
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
<p>Type refers to data carrier form-factor</p>
</td>
</tr>
<tr>
<td>
<code>rotational</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Rotational shows whether disk is solid state or not</p>
</td>
</tr>
<tr>
<td>
<code>system</code><br/>
<em>
string
</em>
</td>
<td>
<p>Bus is a type of hardware interface used to connect the disk to the system</p>
</td>
</tr>
<tr>
<td>
<code>model</code><br/>
<em>
string
</em>
</td>
<td>
<p>Model is a unique hardware part identifier</p>
</td>
</tr>
<tr>
<td>
<code>size</code><br/>
<em>
uint64
</em>
</td>
<td>
<p>Size is a disk space available in bytes</p>
</td>
</tr>
<tr>
<td>
<code>partitionTable</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.PartitionTableSpec">
PartitionTableSpec
</a>
</em>
</td>
<td>
<p>PartitionTable is a partition table currently written to the disk</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha4.CPUSpec">CPUSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.InventorySpec">InventorySpec</a>)
</p>
<div>
<p>CPUSpec contains info about CPUs on hsot machine.</p>
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
<code>physicalId</code><br/>
<em>
uint64
</em>
</td>
<td>
<p>PhysicalID is an ID of physical CPU</p>
</td>
</tr>
<tr>
<td>
<code>logicalIds</code><br/>
<em>
[]uint64
</em>
</td>
<td>
<p>LogicalIDs is a collection of logical CPU nums related to the physical CPU (required for NUMA)</p>
</td>
</tr>
<tr>
<td>
<code>cores</code><br/>
<em>
uint64
</em>
</td>
<td>
<p>Cores is a number of physical cores</p>
</td>
</tr>
<tr>
<td>
<code>siblings</code><br/>
<em>
uint64
</em>
</td>
<td>
<p>Siblings is a number of logical CPUs/threads</p>
</td>
</tr>
<tr>
<td>
<code>vendorId</code><br/>
<em>
string
</em>
</td>
<td>
<p>VendorID is a manufacturer identifire</p>
</td>
</tr>
<tr>
<td>
<code>family</code><br/>
<em>
string
</em>
</td>
<td>
<p>Family refers to processor type</p>
</td>
</tr>
<tr>
<td>
<code>model</code><br/>
<em>
string
</em>
</td>
<td>
<p>Model is a reference id of the model</p>
</td>
</tr>
<tr>
<td>
<code>modelName</code><br/>
<em>
string
</em>
</td>
<td>
<p>ModelName is a common name of the processor</p>
</td>
</tr>
<tr>
<td>
<code>stepping</code><br/>
<em>
string
</em>
</td>
<td>
<p>Stepping is an iteration of the architecture</p>
</td>
</tr>
<tr>
<td>
<code>microcode</code><br/>
<em>
string
</em>
</td>
<td>
<p>Microcode is a firmware reference</p>
</td>
</tr>
<tr>
<td>
<code>mhz</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/api/resource#Quantity">
k8s.io/apimachinery/pkg/api/resource.Quantity
</a>
</em>
</td>
<td>
<p>MHz is a logical core frequency</p>
</td>
</tr>
<tr>
<td>
<code>cacheSize</code><br/>
<em>
string
</em>
</td>
<td>
<p>CacheSize is an L2 cache size</p>
</td>
</tr>
<tr>
<td>
<code>fpu</code><br/>
<em>
bool
</em>
</td>
<td>
<p>FPU defines if CPU has a Floating Point Unit</p>
</td>
</tr>
<tr>
<td>
<code>fpuException</code><br/>
<em>
bool
</em>
</td>
<td>
<p>FPUException</p>
</td>
</tr>
<tr>
<td>
<code>cpuIdLevel</code><br/>
<em>
uint64
</em>
</td>
<td>
<p>CPUIDLevel</p>
</td>
</tr>
<tr>
<td>
<code>wp</code><br/>
<em>
bool
</em>
</td>
<td>
<p>WP tells if WP bit is present</p>
</td>
</tr>
<tr>
<td>
<code>flags</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Flags defines a list of low-level computing capabilities</p>
</td>
</tr>
<tr>
<td>
<code>vmxFlags</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>VMXFlags defines a list of virtualization capabilities</p>
</td>
</tr>
<tr>
<td>
<code>bugs</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Bugs contains a list of known hardware bugs</p>
</td>
</tr>
<tr>
<td>
<code>bogoMips</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/api/resource#Quantity">
k8s.io/apimachinery/pkg/api/resource.Quantity
</a>
</em>
</td>
<td>
<p>BogoMIPS is a synthetic performance metric</p>
</td>
</tr>
<tr>
<td>
<code>clFlushSize</code><br/>
<em>
uint64
</em>
</td>
<td>
<p>CLFlushSize size for cache line flushing feature</p>
</td>
</tr>
<tr>
<td>
<code>cacheAlignment</code><br/>
<em>
uint64
</em>
</td>
<td>
<p>CacheAlignment is a cache size</p>
</td>
</tr>
<tr>
<td>
<code>addressSizes</code><br/>
<em>
string
</em>
</td>
<td>
<p>AddressSizes is an info about address transition system</p>
</td>
</tr>
<tr>
<td>
<code>powerManagement</code><br/>
<em>
string
</em>
</td>
<td>
<p>PowerManagement</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha4.ConditionSpec">ConditionSpec
</h3>
<div>
<p>ConditionSpec contains current condition of port parameters.</p>
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
<h3 id="metal.ironcore.dev/v1alpha4.ConnectionsMap">ConnectionsMap
(<code>map[uint8]*./apis/metal/v1alpha4.NetworkSwitchList</code> alias)</h3>
<div>
</div>
<h3 id="metal.ironcore.dev/v1alpha4.ConstraintSpec">ConstraintSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.SizeSpec">SizeSpec</a>)
</p>
<div>
<p>ConstraintSpec contains conditions of contraint that should be applied on resource.</p>
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
<code>path</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.JSONPath">
JSONPath
</a>
</em>
</td>
<td>
<p>Path is a path to the struct field constraint will be applied to</p>
</td>
</tr>
<tr>
<td>
<code>agg</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.AggregateType">
AggregateType
</a>
</em>
</td>
<td>
<p>Aggregate defines whether collection values should be aggregated
for constraint checks, in case if path defines selector for collection</p>
</td>
</tr>
<tr>
<td>
<code>eq</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.ConstraintValSpec">
ConstraintValSpec
</a>
</em>
</td>
<td>
<p>Equal contains an exact expected value</p>
</td>
</tr>
<tr>
<td>
<code>neq</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.ConstraintValSpec">
ConstraintValSpec
</a>
</em>
</td>
<td>
<p>NotEqual contains an exact not expected value</p>
</td>
</tr>
<tr>
<td>
<code>lt</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/api/resource#Quantity">
k8s.io/apimachinery/pkg/api/resource.Quantity
</a>
</em>
</td>
<td>
<p>LessThan contains an highest expected value, exclusive</p>
</td>
</tr>
<tr>
<td>
<code>lte</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/api/resource#Quantity">
k8s.io/apimachinery/pkg/api/resource.Quantity
</a>
</em>
</td>
<td>
<p>LessThan contains an highest expected value, inclusive</p>
</td>
</tr>
<tr>
<td>
<code>gt</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/api/resource#Quantity">
k8s.io/apimachinery/pkg/api/resource.Quantity
</a>
</em>
</td>
<td>
<p>LessThan contains an lowest expected value, exclusive</p>
</td>
</tr>
<tr>
<td>
<code>gte</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/api/resource#Quantity">
k8s.io/apimachinery/pkg/api/resource.Quantity
</a>
</em>
</td>
<td>
<p>GreaterThanOrEqual contains an lowest expected value, inclusive</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha4.ConstraintValSpec">ConstraintValSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.ConstraintSpec">ConstraintSpec</a>)
</p>
<div>
<p>ConstraintValSpec is a wrapper around value for constraint.
Since it is not possilble to set oneOf/anyOf through kubebuilder
markers, type is set to number here, and patched with kustomize
See <a href="https://github.com/kubernetes-sigs/kubebuilder/issues/301">https://github.com/kubernetes-sigs/kubebuilder/issues/301</a></p>
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
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>-</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/api/resource#Quantity">
k8s.io/apimachinery/pkg/api/resource.Quantity
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha4.DistroSpec">DistroSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.InventorySpec">InventorySpec</a>)
</p>
<div>
<p>DistroSpec contains info about distro.</p>
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
<code>buildVersion</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>debianVersion</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>kernelVersion</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>asicType</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>commitID</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>buildDate</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>buildNumber</code><br/>
<em>
uint32
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>buildBy</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha4.FieldSelectorSpec">FieldSelectorSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.IPAMSelectionSpec">IPAMSelectionSpec</a>)
</p>
<div>
<p>FieldSelectorSpec contains label key and field path where to get label value for search.</p>
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
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectfieldselector-v1-core">
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
<h3 id="metal.ironcore.dev/v1alpha4.GeneralIPAMSpec">GeneralIPAMSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.SwitchConfigSpec">SwitchConfigSpec</a>)
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
<a href="#metal.ironcore.dev/v1alpha4.AddressFamiliesMap">
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
<a href="#metal.ironcore.dev/v1alpha4.IPAMSelectionSpec">
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
<a href="#metal.ironcore.dev/v1alpha4.IPAMSelectionSpec">
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
<a href="#metal.ironcore.dev/v1alpha4.IPAMSelectionSpec">
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
<a href="#metal.ironcore.dev/v1alpha4.IPAMSelectionSpec">
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
<h3 id="metal.ironcore.dev/v1alpha4.HostSpec">HostSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.InventorySpec">InventorySpec</a>)
</p>
<div>
<p>HostSpec contains type of inventorying object and in case it is a switch - SONiC version.</p>
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
<p>Hostname contains hostname</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha4.IPAMSelectionSpec">IPAMSelectionSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.GeneralIPAMSpec">GeneralIPAMSpec</a>, <a href="#metal.ironcore.dev/v1alpha4.IPAMSpec">IPAMSpec</a>)
</p>
<div>
<p>IPAMSelectionSpec contains label selector and address family.</p>
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
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#labelselector-v1-meta">
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
<a href="#metal.ironcore.dev/v1alpha4.FieldSelectorSpec">
FieldSelectorSpec
</a>
</em>
</td>
<td>
<p>FieldSelector contains label key and field path where to get label value for search.
If FieldSelector is used as part of IPAM configuration in SwitchConfig object it will
reference to the field path in related NetworkSwitch object. If FieldSelector is used as part of IPAM
configuration in NetworkSwitch object, it will reference to the field path in the same object</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha4.IPAMSpec">IPAMSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.NetworkSwitchSpec">NetworkSwitchSpec</a>)
</p>
<div>
<p>IPAMSpec contains selectors for subnets and loopback IPs and
definition of address families which should be claimed.</p>
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
<a href="#metal.ironcore.dev/v1alpha4.IPAMSelectionSpec">
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
<a href="#metal.ironcore.dev/v1alpha4.IPAMSelectionSpec">
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
<h3 id="metal.ironcore.dev/v1alpha4.IPAddrSpec">IPAddrSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.Addresses">Addresses</a>, <a href="#metal.ironcore.dev/v1alpha4.LoopbackAddresses">LoopbackAddresses</a>)
</p>
<div>
<p>IPAddrSpec defines interface&rsquo;s ip address info.</p>
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
<a href="https://pkg.go.dev/net/netip#Prefix">
net/netip.Prefix
</a>
</em>
</td>
<td>
<p>Address refers to the ip address value</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha4.IPAddressSpec">IPAddressSpec
</h3>
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
<code>ObjectReference</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.ObjectReference">
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
<h3 id="metal.ironcore.dev/v1alpha4.IPMISpec">IPMISpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.InventorySpec">InventorySpec</a>)
</p>
<div>
<p>IPMISpec contains info about IPMI module.</p>
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
<code>ipAddress</code><br/>
<em>
string
</em>
</td>
<td>
<p>IPAddress is an IP address assigned to IPMI network interface</p>
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
<p>MACAddress is a MAC address of IPMI&rsquo;s network interface</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha4.Identity">Identity
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.MachineSpec">MachineSpec</a>)
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
<a href="#metal.ironcore.dev/v1alpha4.Internal">
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
<h3 id="metal.ironcore.dev/v1alpha4.Interface">Interface
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.Network">Network</a>)
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
<a href="#metal.ironcore.dev/v1alpha4.ResourceReference">
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
<a href="#metal.ironcore.dev/v1alpha4.Addresses">
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
<a href="#metal.ironcore.dev/v1alpha4.Peer">
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
<h3 id="metal.ironcore.dev/v1alpha4.InterfaceOverridesSpec">InterfaceOverridesSpec
</h3>
<div>
<p>InterfaceOverridesSpec contains overridden parameters for certain switch port.</p>
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
<a href="#metal.ironcore.dev/v1alpha4.PortParametersSpec">
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
<a href="#metal.ironcore.dev/v1alpha4.*./apis/metal/v1alpha4.AdditionalIPSpec">
[]*./apis/metal/v1alpha4.AdditionalIPSpec
</a>
</em>
</td>
<td>
<p>IP contains a list of additional IP addresses for interface</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha4.InterfaceSpec">InterfaceSpec
</h3>
<div>
<p>InterfaceSpec defines the state of switch&rsquo;s interface.</p>
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
<a href="#metal.ironcore.dev/v1alpha4.PortParametersSpec">
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
<a href="#metal.ironcore.dev/v1alpha4.*./apis/metal/v1alpha4.IPAddressSpec">
[]*./apis/metal/v1alpha4.IPAddressSpec
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
<a href="#metal.ironcore.dev/v1alpha4.PeerSpec">
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
<h3 id="metal.ironcore.dev/v1alpha4.InterfacesSpec">InterfacesSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.NetworkSwitchSpec">NetworkSwitchSpec</a>)
</p>
<div>
<p>InterfacesSpec contains definitions for general switch ports&rsquo; configuration.</p>
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
<a href="#metal.ironcore.dev/v1alpha4.PortParametersSpec">
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
<a href="#metal.ironcore.dev/v1alpha4.*./apis/metal/v1alpha4.InterfaceOverridesSpec">
[]*./apis/metal/v1alpha4.InterfaceOverridesSpec
</a>
</em>
</td>
<td>
<p>Overrides contains set of parameters which should be overridden for listed switch ports</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha4.Internal">Internal
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.Identity">Identity</a>)
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
<h3 id="metal.ironcore.dev/v1alpha4.InventorySpec">InventorySpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.Inventory">Inventory</a>, <a href="#metal.ironcore.dev/v1alpha4.ValidationInventory">ValidationInventory</a>)
</p>
<div>
<p>InventorySpec contains result of inventorization process on the host.</p>
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
<code>system</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.SystemSpec">
SystemSpec
</a>
</em>
</td>
<td>
<p>System contains DMI system information</p>
</td>
</tr>
<tr>
<td>
<code>ipmis</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.IPMISpec">
[]IPMISpec
</a>
</em>
</td>
<td>
<p>IPMIs contains info about IPMI interfaces on the host</p>
</td>
</tr>
<tr>
<td>
<code>blocks</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.BlockSpec">
[]BlockSpec
</a>
</em>
</td>
<td>
<p>Blocks contains info about block devices on the host</p>
</td>
</tr>
<tr>
<td>
<code>memory</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.MemorySpec">
MemorySpec
</a>
</em>
</td>
<td>
<p>Memory contains info block devices on the host</p>
</td>
</tr>
<tr>
<td>
<code>cpus</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.CPUSpec">
[]CPUSpec
</a>
</em>
</td>
<td>
<p>CPUs contains info about cpus, cores and threads</p>
</td>
</tr>
<tr>
<td>
<code>numa</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.NumaSpec">
[]NumaSpec
</a>
</em>
</td>
<td>
<p>NUMA contains info about cpu/memory topology</p>
</td>
</tr>
<tr>
<td>
<code>pciDevices</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.PCIDeviceSpec">
[]PCIDeviceSpec
</a>
</em>
</td>
<td>
<p>PCIDevices contains info about devices accessible through</p>
</td>
</tr>
<tr>
<td>
<code>nics</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.NICSpec">
[]NICSpec
</a>
</em>
</td>
<td>
<p>NICs contains info about network interfaces and network discovery</p>
</td>
</tr>
<tr>
<td>
<code>virt</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.VirtSpec">
VirtSpec
</a>
</em>
</td>
<td>
<p>Virt is a virtualization detected on host</p>
</td>
</tr>
<tr>
<td>
<code>host</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.HostSpec">
HostSpec
</a>
</em>
</td>
<td>
<p>Host contains info about inventorying object</p>
</td>
</tr>
<tr>
<td>
<code>distro</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.DistroSpec">
DistroSpec
</a>
</em>
</td>
<td>
<p>Distro contains info about OS distro</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha4.InventoryStatus">InventoryStatus
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.Inventory">Inventory</a>)
</p>
<div>
<p>InventoryStatus defines the observed state of Inventory.</p>
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
<code>computed</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.AggregationResults">
AggregationResults
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>inventoryStatuses</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.InventoryStatuses">
InventoryStatuses
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha4.InventoryStatuses">InventoryStatuses
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.InventoryStatus">InventoryStatus</a>)
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
<code>ready</code><br/>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>requestsCount</code><br/>
<em>
int
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha4.JSONPath">JSONPath
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.AggregateItem">AggregateItem</a>, <a href="#metal.ironcore.dev/v1alpha4.ConstraintSpec">ConstraintSpec</a>)
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
<code>-</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha4.LLDPCapabilities">LLDPCapabilities
(<code>string</code> alias)</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.LLDPSpec">LLDPSpec</a>)
</p>
<div>
<p>LLDPCapabilities</p>
</div>
<h3 id="metal.ironcore.dev/v1alpha4.LLDPSpec">LLDPSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.NICSpec">NICSpec</a>)
</p>
<div>
<p>LLDPSpec is an entry received by network interface by Link Layer Discovery Protocol.</p>
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
<p>ChassisID is a neighbour box identifier</p>
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
<p>SystemName is given name to the neighbour box</p>
</td>
</tr>
<tr>
<td>
<code>systemDescription</code><br/>
<em>
string
</em>
</td>
<td>
<p>SystemDescription is a short description of the neighbour box</p>
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
<p>PortID is a hardware identifier of the link port</p>
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
<p>PortDescription is a short description of the link port</p>
</td>
</tr>
<tr>
<td>
<code>capabilities</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.LLDPCapabilities">
[]LLDPCapabilities
</a>
</em>
</td>
<td>
<p>Capabilities is a list of LLDP capabilities advertised by neighbor</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha4.LoopbackAddresses">LoopbackAddresses
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.Network">Network</a>)
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
<a href="#metal.ironcore.dev/v1alpha4.IPAddrSpec">
IPAddrSpec
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
<a href="#metal.ironcore.dev/v1alpha4.IPAddrSpec">
IPAddrSpec
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha4.MachineSpec">MachineSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.Machine">Machine</a>)
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
<a href="#metal.ironcore.dev/v1alpha4.Identity">
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
<h3 id="metal.ironcore.dev/v1alpha4.MachineState">MachineState
(<code>string</code> alias)</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.MachineStatus">MachineStatus</a>)
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
<h3 id="metal.ironcore.dev/v1alpha4.MachineStatus">MachineStatus
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.Machine">Machine</a>)
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
<a href="#metal.ironcore.dev/v1alpha4.MachineState">
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
<a href="#metal.ironcore.dev/v1alpha4.Network">
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
<a href="#metal.ironcore.dev/v1alpha4.Reservation">
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
<h3 id="metal.ironcore.dev/v1alpha4.MemorySpec">MemorySpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.InventorySpec">InventorySpec</a>, <a href="#metal.ironcore.dev/v1alpha4.NumaSpec">NumaSpec</a>)
</p>
<div>
<p>MemorySpec contains info about RAM on host.</p>
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
<code>total</code><br/>
<em>
uint64
</em>
</td>
<td>
<p>Total is a total amount of RAM on host</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha4.NDPSpec">NDPSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.NICSpec">NICSpec</a>)
</p>
<div>
<p>NDPSpec is an entry received by IPv6 Neighbour Discovery Protocol.</p>
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
<code>ipAddress</code><br/>
<em>
string
</em>
</td>
<td>
<p>IPAddress is an IPv6 address of a neighbour</p>
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
<p>MACAddress is an MAC address of a neighbour</p>
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
<p>State is a state of discovery</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha4.NICSpec">NICSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.InventorySpec">InventorySpec</a>)
</p>
<div>
<p>NICSpec contains info about network interfaces.</p>
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
<p>Name is a name of the device registered by Linux Kernel</p>
</td>
</tr>
<tr>
<td>
<code>pciAddress</code><br/>
<em>
string
</em>
</td>
<td>
<p>PCIAddress is the PCI bus address network interface is connected to</p>
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
<p>MACAddress is the MAC address of network interface</p>
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
<p>MTU is refers to Maximum Transmission Unit</p>
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
<p>Speed is a speed of network interface in Mbits/s</p>
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
<p>Lanes is a number of used lanes (if supported)</p>
</td>
</tr>
<tr>
<td>
<code>activeFEC</code><br/>
<em>
string
</em>
</td>
<td>
<p>ActiveFEC is an active error correction mode</p>
</td>
</tr>
<tr>
<td>
<code>lldps</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.LLDPSpec">
[]LLDPSpec
</a>
</em>
</td>
<td>
<p>LLDP is a collection of LLDP messages received by the network interface</p>
</td>
</tr>
<tr>
<td>
<code>ndps</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.NDPSpec">
[]NDPSpec
</a>
</em>
</td>
<td>
<p>NDP is a collection of NDP messages received by the network interface</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha4.Network">Network
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.MachineStatus">MachineStatus</a>)
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
<a href="#metal.ironcore.dev/v1alpha4.Interface">
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
<a href="#metal.ironcore.dev/v1alpha4.LoopbackAddresses">
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
<h3 id="metal.ironcore.dev/v1alpha4.NetworkSwitchSpec">NetworkSwitchSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.NetworkSwitch">NetworkSwitch</a>)
</p>
<div>
<p>NetworkSwitchSpec contains desired state of resulting NetworkSwitch configuration.</p>
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
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#localobjectreference-v1-core">
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
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#labelselector-v1-meta">
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
<p>Managed is a flag defining whether NetworkSwitch object would be processed during reconciliation</p>
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
<p>Cordon is a flag defining whether NetworkSwitch object is taken offline</p>
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
<p>TopSpine is a flag defining whether NetworkSwitch is a top-level spine switch</p>
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
<a href="#metal.ironcore.dev/v1alpha4.IPAMSpec">
IPAMSpec
</a>
</em>
</td>
<td>
<p>IPAM refers to selectors for subnets which will be used for NetworkSwitch object</p>
</td>
</tr>
<tr>
<td>
<code>interfaces</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.InterfacesSpec">
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
<h3 id="metal.ironcore.dev/v1alpha4.NetworkSwitchStatus">NetworkSwitchStatus
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.NetworkSwitch">NetworkSwitch</a>)
</p>
<div>
<p>NetworkSwitchStatus contains observed state of NetworkSwitch.</p>
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
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#localobjectreference-v1-core">
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
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#localobjectreference-v1-core">
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
<a href="#metal.ironcore.dev/v1alpha4.*./apis/metal/v1alpha4.InterfaceSpec">
map[string]*./apis/metal/v1alpha4.InterfaceSpec
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
<a href="#metal.ironcore.dev/v1alpha4.*./apis/metal/v1alpha4.SubnetSpec">
[]*./apis/metal/v1alpha4.SubnetSpec
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
<a href="#metal.ironcore.dev/v1alpha4.*./apis/metal/v1alpha4.IPAddressSpec">
[]*./apis/metal/v1alpha4.IPAddressSpec
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
<a href="#metal.ironcore.dev/v1alpha4.*./apis/metal/v1alpha4.ConditionSpec">
[]*./apis/metal/v1alpha4.ConditionSpec
</a>
</em>
</td>
<td>
<p>Condition contains state of port parameters</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha4.NumaSpec">NumaSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.InventorySpec">InventorySpec</a>)
</p>
<div>
<p>NumaSpec describes NUMA node.</p>
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
<code>id</code><br/>
<em>
int
</em>
</td>
<td>
<p>ID is NUMA node ID.</p>
</td>
</tr>
<tr>
<td>
<code>cpus</code><br/>
<em>
[]int
</em>
</td>
<td>
<p>CPUs is a list of CPU logical IDs in current numa node.</p>
</td>
</tr>
<tr>
<td>
<code>distances</code><br/>
<em>
[]int
</em>
</td>
<td>
<p>Distances contains distances to other nodes. Element index corresponds to NUMA node ID.</p>
</td>
</tr>
<tr>
<td>
<code>memory</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.MemorySpec">
MemorySpec
</a>
</em>
</td>
<td>
<p>Memory contains info about NUMA node memory setup.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha4.ObjectReference">ObjectReference
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.IPAddressSpec">IPAddressSpec</a>, <a href="#metal.ironcore.dev/v1alpha4.PeerSpec">PeerSpec</a>, <a href="#metal.ironcore.dev/v1alpha4.SubnetSpec">SubnetSpec</a>)
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
<h3 id="metal.ironcore.dev/v1alpha4.PCIDeviceDescriptionSpec">PCIDeviceDescriptionSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.PCIDeviceSpec">PCIDeviceSpec</a>)
</p>
<div>
<p>PCIDeviceDescriptionSpec contains one of the options that is describing the PCI device.</p>
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
<code>id</code><br/>
<em>
string
</em>
</td>
<td>
<p>ID is a hexadecimal identifier of device property , that corresponds to the value from PCIIDs database</p>
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
<p>Name is a string value of property extracted from PCIID DB</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha4.PCIDeviceSpec">PCIDeviceSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.InventorySpec">InventorySpec</a>)
</p>
<div>
<p>PCIDeviceSpec contains description of PCI device.</p>
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
<code>busId</code><br/>
<em>
string
</em>
</td>
<td>
<p>BusID is an ID of PCI bus on the board device is attached to.</p>
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
<p>Address is an ID of device on PCI bus.</p>
</td>
</tr>
<tr>
<td>
<code>vendor</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.PCIDeviceDescriptionSpec">
PCIDeviceDescriptionSpec
</a>
</em>
</td>
<td>
<p>Vendor refers to manufacturer ore device trademark.</p>
</td>
</tr>
<tr>
<td>
<code>subvendor</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.PCIDeviceDescriptionSpec">
PCIDeviceDescriptionSpec
</a>
</em>
</td>
<td>
<p>Subvendor usually refers to the platform or co-manufacturer. E.g. Lenovo board manufactured for Intel platform (by Intel spec).</p>
</td>
</tr>
<tr>
<td>
<code>type</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.PCIDeviceDescriptionSpec">
PCIDeviceDescriptionSpec
</a>
</em>
</td>
<td>
<p>Type shows device&rsquo;s designation.</p>
</td>
</tr>
<tr>
<td>
<code>subtype</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.PCIDeviceDescriptionSpec">
PCIDeviceDescriptionSpec
</a>
</em>
</td>
<td>
<p>Subtype shows device&rsquo;s subsystem.</p>
</td>
</tr>
<tr>
<td>
<code>class</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.PCIDeviceDescriptionSpec">
PCIDeviceDescriptionSpec
</a>
</em>
</td>
<td>
<p>Class refers to generic device designation.</p>
</td>
</tr>
<tr>
<td>
<code>subclass</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.PCIDeviceDescriptionSpec">
PCIDeviceDescriptionSpec
</a>
</em>
</td>
<td>
<p>Subclass narrows the designation scope.</p>
</td>
</tr>
<tr>
<td>
<code>interface</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.PCIDeviceDescriptionSpec">
PCIDeviceDescriptionSpec
</a>
</em>
</td>
<td>
<p>ProgrammingInterface specifies communication protocols.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha4.PartitionSpec">PartitionSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.PartitionTableSpec">PartitionTableSpec</a>)
</p>
<div>
<p>PartitionSpec contains info about partition.</p>
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
<code>id</code><br/>
<em>
string
</em>
</td>
<td>
<p>ID is a GUID of GPT partition or number for MBR partition</p>
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
<p>Name is a human readable name given to the partition</p>
</td>
</tr>
<tr>
<td>
<code>size</code><br/>
<em>
uint64
</em>
</td>
<td>
<p>Size is a size of partition in bytes</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha4.PartitionTableSpec">PartitionTableSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.BlockSpec">BlockSpec</a>)
</p>
<div>
<p>PartitionTableSpec contains info about partition table on block device.</p>
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
<code>type</code><br/>
<em>
string
</em>
</td>
<td>
<p>Type is a format of partition table</p>
</td>
</tr>
<tr>
<td>
<code>partitions</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.PartitionSpec">
[]PartitionSpec
</a>
</em>
</td>
<td>
<p>Partitions are active partition records on disk</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha4.Peer">Peer
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.Interface">Interface</a>)
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
<a href="#metal.ironcore.dev/v1alpha4.ResourceReference">
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
<h3 id="metal.ironcore.dev/v1alpha4.PeerInfoSpec">PeerInfoSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.PeerSpec">PeerSpec</a>)
</p>
<div>
<p>PeerInfoSpec contains LLDP info about peer.</p>
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
<h3 id="metal.ironcore.dev/v1alpha4.PeerSpec">PeerSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.InterfaceSpec">InterfaceSpec</a>)
</p>
<div>
<p>PeerSpec defines peer info.</p>
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
<a href="#metal.ironcore.dev/v1alpha4.ObjectReference">
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
<a href="#metal.ironcore.dev/v1alpha4.PeerInfoSpec">
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
<h3 id="metal.ironcore.dev/v1alpha4.PortParametersSpec">PortParametersSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.InterfaceOverridesSpec">InterfaceOverridesSpec</a>, <a href="#metal.ironcore.dev/v1alpha4.InterfaceSpec">InterfaceSpec</a>, <a href="#metal.ironcore.dev/v1alpha4.InterfacesSpec">InterfacesSpec</a>, <a href="#metal.ironcore.dev/v1alpha4.SwitchConfigSpec">SwitchConfigSpec</a>)
</p>
<div>
<p>PortParametersSpec contains a set of parameters of switch port.</p>
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
<h3 id="metal.ironcore.dev/v1alpha4.RegionSpec">RegionSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.SubnetSpec">SubnetSpec</a>)
</p>
<div>
<p>RegionSpec defines region info.</p>
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
<h3 id="metal.ironcore.dev/v1alpha4.Reservation">Reservation
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.MachineStatus">MachineStatus</a>)
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
<a href="#metal.ironcore.dev/v1alpha4.ResourceReference">
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
<h3 id="metal.ironcore.dev/v1alpha4.ResourceReference">ResourceReference
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.Interface">Interface</a>, <a href="#metal.ironcore.dev/v1alpha4.Peer">Peer</a>, <a href="#metal.ironcore.dev/v1alpha4.Reservation">Reservation</a>)
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
<h3 id="metal.ironcore.dev/v1alpha4.SizeSpec">SizeSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.Size">Size</a>)
</p>
<div>
<p>SizeSpec defines the desired state of Size.</p>
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
<code>constraints</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.ConstraintSpec">
[]ConstraintSpec
</a>
</em>
</td>
<td>
<p>Constraints is a list of selectors based on machine properties.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha4.SizeStatus">SizeStatus
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.Size">Size</a>)
</p>
<div>
<p>SizeStatus defines the observed state of Size.</p>
</div>
<h3 id="metal.ironcore.dev/v1alpha4.SubnetSpec">SubnetSpec
</h3>
<div>
<p>SubnetSpec defines switch&rsquo;s subnet info.</p>
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
<a href="#metal.ironcore.dev/v1alpha4.ObjectReference">
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
<a href="#metal.ironcore.dev/v1alpha4.ObjectReference">
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
<a href="#metal.ironcore.dev/v1alpha4.RegionSpec">
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
<h3 id="metal.ironcore.dev/v1alpha4.SwitchConfigSpec">SwitchConfigSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.SwitchConfig">SwitchConfig</a>)
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
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#labelselector-v1-meta">
Kubernetes meta/v1.LabelSelector
</a>
</em>
</td>
<td>
<p>Switches contains label selector to pick up NetworkSwitch objects</p>
</td>
</tr>
<tr>
<td>
<code>portsDefaults</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.PortParametersSpec">
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
<a href="#metal.ironcore.dev/v1alpha4.GeneralIPAMSpec">
GeneralIPAMSpec
</a>
</em>
</td>
<td>
<p>IPAM refers to selectors for subnets which will be used for NetworkSwitch object</p>
</td>
</tr>
<tr>
<td>
<code>routingConfigTemplate</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#localobjectreference-v1-core">
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
<h3 id="metal.ironcore.dev/v1alpha4.SwitchConfigStatus">SwitchConfigStatus
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.SwitchConfig">SwitchConfig</a>)
</p>
<div>
<p>SwitchConfigStatus contains observed state of SwitchConfig.</p>
</div>
<h3 id="metal.ironcore.dev/v1alpha4.SystemSpec">SystemSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.InventorySpec">InventorySpec</a>)
</p>
<div>
<p>SystemSpec contains DMI system information.</p>
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
<code>id</code><br/>
<em>
string
</em>
</td>
<td>
<p>ID is a UUID of a system board</p>
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
<p>Manufacturer refers to the company that produced the product</p>
</td>
</tr>
<tr>
<td>
<code>productSku</code><br/>
<em>
string
</em>
</td>
<td>
<p>ProductSKU is a product&rsquo;s Stock Keeping Unit</p>
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
<p>SerialNumber contains serial number of a system</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha4.ValidationInventory">ValidationInventory
</h3>
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
<code>spec</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.InventorySpec">
InventorySpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>system</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.SystemSpec">
SystemSpec
</a>
</em>
</td>
<td>
<p>System contains DMI system information</p>
</td>
</tr>
<tr>
<td>
<code>ipmis</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.IPMISpec">
[]IPMISpec
</a>
</em>
</td>
<td>
<p>IPMIs contains info about IPMI interfaces on the host</p>
</td>
</tr>
<tr>
<td>
<code>blocks</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.BlockSpec">
[]BlockSpec
</a>
</em>
</td>
<td>
<p>Blocks contains info about block devices on the host</p>
</td>
</tr>
<tr>
<td>
<code>memory</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.MemorySpec">
MemorySpec
</a>
</em>
</td>
<td>
<p>Memory contains info block devices on the host</p>
</td>
</tr>
<tr>
<td>
<code>cpus</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.CPUSpec">
[]CPUSpec
</a>
</em>
</td>
<td>
<p>CPUs contains info about cpus, cores and threads</p>
</td>
</tr>
<tr>
<td>
<code>numa</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.NumaSpec">
[]NumaSpec
</a>
</em>
</td>
<td>
<p>NUMA contains info about cpu/memory topology</p>
</td>
</tr>
<tr>
<td>
<code>pciDevices</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.PCIDeviceSpec">
[]PCIDeviceSpec
</a>
</em>
</td>
<td>
<p>PCIDevices contains info about devices accessible through</p>
</td>
</tr>
<tr>
<td>
<code>nics</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.NICSpec">
[]NICSpec
</a>
</em>
</td>
<td>
<p>NICs contains info about network interfaces and network discovery</p>
</td>
</tr>
<tr>
<td>
<code>virt</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.VirtSpec">
VirtSpec
</a>
</em>
</td>
<td>
<p>Virt is a virtualization detected on host</p>
</td>
</tr>
<tr>
<td>
<code>host</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.HostSpec">
HostSpec
</a>
</em>
</td>
<td>
<p>Host contains info about inventorying object</p>
</td>
</tr>
<tr>
<td>
<code>distro</code><br/>
<em>
<a href="#metal.ironcore.dev/v1alpha4.DistroSpec">
DistroSpec
</a>
</em>
</td>
<td>
<p>Distro contains info about OS distro</p>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="metal.ironcore.dev/v1alpha4.VirtSpec">VirtSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metal.ironcore.dev/v1alpha4.InventorySpec">InventorySpec</a>)
</p>
<div>
<p>VirtSpec contains info about detected host virtualization.</p>
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
<code>vmType</code><br/>
<em>
string
</em>
</td>
<td>
<p>VMType is a type of virtual machine engine</p>
</td>
</tr>
</tbody>
</table>
<hr/>
<p><em>
Generated with <code>gen-crd-api-reference-docs</code>
on git commit <code>26e566b</code>.
</em></p>
