<p>Packages:</p>
<ul>
<li>
<a href="#machine.onmetal.de%2fv1alpha1">machine.onmetal.de/v1alpha1</a>
</li>
</ul>
<h2 id="machine.onmetal.de/v1alpha1">machine.onmetal.de/v1alpha1</h2>
Resource Types:
<ul></ul>
<h3 id="machine.onmetal.de/v1alpha1.Aggregate">Aggregate
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
<a href="#machine.onmetal.de/v1alpha1.AggregateSpec">
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
<a href="#machine.onmetal.de/v1alpha1.AggregateItem">
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
<a href="#machine.onmetal.de/v1alpha1.AggregateStatus">
AggregateStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="machine.onmetal.de/v1alpha1.AggregateItem">AggregateItem
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha1.AggregateSpec">AggregateSpec</a>)
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
<a href="#machine.onmetal.de/v1alpha1.JSONPath">
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
<a href="#machine.onmetal.de/v1alpha1.JSONPath">
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
<a href="#machine.onmetal.de/v1alpha1.AggregateType">
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
<h3 id="machine.onmetal.de/v1alpha1.AggregateSpec">AggregateSpec
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha1.Aggregate">Aggregate</a>)
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
<a href="#machine.onmetal.de/v1alpha1.AggregateItem">
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
<h3 id="machine.onmetal.de/v1alpha1.AggregateStatus">AggregateStatus
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha1.Aggregate">Aggregate</a>)
</p>
<div>
<p>AggregateStatus defines the observed state of Aggregate.</p>
</div>
<h3 id="machine.onmetal.de/v1alpha1.AggregateType">AggregateType
(<code>string</code> alias)</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha1.AggregateItem">AggregateItem</a>, <a href="#machine.onmetal.de/v1alpha1.ConstraintSpec">ConstraintSpec</a>)
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
<h3 id="machine.onmetal.de/v1alpha1.AggregationResults">AggregationResults
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha1.InventoryStatus">InventoryStatus</a>)
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
<h3 id="machine.onmetal.de/v1alpha1.BlockSpec">BlockSpec
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha1.InventorySpec">InventorySpec</a>)
</p>
<div>
<p>BlockSpec contains info about block device</p>
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
<a href="#machine.onmetal.de/v1alpha1.PartitionTableSpec">
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
<h3 id="machine.onmetal.de/v1alpha1.CPUSpec">CPUSpec
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha1.InventorySpec">InventorySpec</a>)
</p>
<div>
<p>CPUSpec contains info about CPUs on hsot machine</p>
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
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/api/resource#Duration">
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
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/api/resource#Duration">
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
<h3 id="machine.onmetal.de/v1alpha1.ConstraintSpec">ConstraintSpec
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha1.SizeSpec">SizeSpec</a>)
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
<a href="#machine.onmetal.de/v1alpha1.JSONPath">
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
<a href="#machine.onmetal.de/v1alpha1.AggregateType">
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
<a href="#machine.onmetal.de/v1alpha1.ConstraintValSpec">
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
<a href="#machine.onmetal.de/v1alpha1.ConstraintValSpec">
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
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/api/resource#Duration">
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
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/api/resource#Duration">
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
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/api/resource#Duration">
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
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/api/resource#Duration">
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
<h3 id="machine.onmetal.de/v1alpha1.ConstraintValSpec">ConstraintValSpec
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha1.ConstraintSpec">ConstraintSpec</a>)
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
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/api/resource#Duration">
k8s.io/apimachinery/pkg/api/resource.Quantity
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="machine.onmetal.de/v1alpha1.DistroSpec">DistroSpec
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha1.InventorySpec">InventorySpec</a>)
</p>
<div>
<p>DistroSpec contains info about distro</p>
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
<h3 id="machine.onmetal.de/v1alpha1.HostSpec">HostSpec
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha1.InventorySpec">InventorySpec</a>)
</p>
<div>
<p>HostSpec contains type of inventorying object and in case it is a switch - SONiC version</p>
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
<h3 id="machine.onmetal.de/v1alpha1.IPMISpec">IPMISpec
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha1.InventorySpec">InventorySpec</a>)
</p>
<div>
<p>IPMISpec contains info about IPMI module</p>
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
<h3 id="machine.onmetal.de/v1alpha1.Inventory">Inventory
</h3>
<div>
<p>Inventory is the Schema for the inventories API</p>
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
<a href="#machine.onmetal.de/v1alpha1.InventorySpec">
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
<a href="#machine.onmetal.de/v1alpha1.SystemSpec">
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
<a href="#machine.onmetal.de/v1alpha1.IPMISpec">
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
<a href="#machine.onmetal.de/v1alpha1.BlockSpec">
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
<a href="#machine.onmetal.de/v1alpha1.MemorySpec">
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
<a href="#machine.onmetal.de/v1alpha1.CPUSpec">
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
<a href="#machine.onmetal.de/v1alpha1.NumaSpec">
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
<a href="#machine.onmetal.de/v1alpha1.PCIDeviceSpec">
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
<a href="#machine.onmetal.de/v1alpha1.NICSpec">
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
<a href="#machine.onmetal.de/v1alpha1.VirtSpec">
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
<a href="#machine.onmetal.de/v1alpha1.HostSpec">
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
<a href="#machine.onmetal.de/v1alpha1.DistroSpec">
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
<a href="#machine.onmetal.de/v1alpha1.InventoryStatus">
InventoryStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="machine.onmetal.de/v1alpha1.InventorySpec">InventorySpec
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha1.Inventory">Inventory</a>, <a href="#machine.onmetal.de/v1alpha1.ValidationInventory">ValidationInventory</a>)
</p>
<div>
<p>InventorySpec contains result of inventorization process on the host</p>
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
<a href="#machine.onmetal.de/v1alpha1.SystemSpec">
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
<a href="#machine.onmetal.de/v1alpha1.IPMISpec">
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
<a href="#machine.onmetal.de/v1alpha1.BlockSpec">
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
<a href="#machine.onmetal.de/v1alpha1.MemorySpec">
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
<a href="#machine.onmetal.de/v1alpha1.CPUSpec">
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
<a href="#machine.onmetal.de/v1alpha1.NumaSpec">
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
<a href="#machine.onmetal.de/v1alpha1.PCIDeviceSpec">
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
<a href="#machine.onmetal.de/v1alpha1.NICSpec">
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
<a href="#machine.onmetal.de/v1alpha1.VirtSpec">
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
<a href="#machine.onmetal.de/v1alpha1.HostSpec">
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
<a href="#machine.onmetal.de/v1alpha1.DistroSpec">
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
<h3 id="machine.onmetal.de/v1alpha1.InventoryStatus">InventoryStatus
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha1.Inventory">Inventory</a>)
</p>
<div>
<p>InventoryStatus defines the observed state of Inventory</p>
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
<a href="#machine.onmetal.de/v1alpha1.AggregationResults">
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
<a href="#machine.onmetal.de/v1alpha1.InventoryStatuses">
InventoryStatuses
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="machine.onmetal.de/v1alpha1.InventoryStatuses">InventoryStatuses
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha1.InventoryStatus">InventoryStatus</a>)
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
<h3 id="machine.onmetal.de/v1alpha1.JSONPath">JSONPath
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha1.AggregateItem">AggregateItem</a>, <a href="#machine.onmetal.de/v1alpha1.ConstraintSpec">ConstraintSpec</a>)
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
<h3 id="machine.onmetal.de/v1alpha1.LLDPCapabilities">LLDPCapabilities
(<code>string</code> alias)</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha1.LLDPSpec">LLDPSpec</a>)
</p>
<div>
<p>LLDPCapabilities</p>
</div>
<h3 id="machine.onmetal.de/v1alpha1.LLDPSpec">LLDPSpec
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha1.NICSpec">NICSpec</a>)
</p>
<div>
<p>LLDPSpec is an entry received by network interface by Link Layer Discovery Protocol</p>
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
<a href="#machine.onmetal.de/v1alpha1.LLDPCapabilities">
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
<h3 id="machine.onmetal.de/v1alpha1.MemorySpec">MemorySpec
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha1.InventorySpec">InventorySpec</a>, <a href="#machine.onmetal.de/v1alpha1.NumaSpec">NumaSpec</a>)
</p>
<div>
<p>MemorySpec contains info about RAM on host</p>
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
<h3 id="machine.onmetal.de/v1alpha1.NDPSpec">NDPSpec
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha1.NICSpec">NICSpec</a>)
</p>
<div>
<p>NDPSpec is an entry received by IPv6 Neighbour Discovery Protocol</p>
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
<h3 id="machine.onmetal.de/v1alpha1.NICSpec">NICSpec
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha1.InventorySpec">InventorySpec</a>)
</p>
<div>
<p>NICSpec contains info about network interfaces</p>
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
<a href="#machine.onmetal.de/v1alpha1.LLDPSpec">
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
<a href="#machine.onmetal.de/v1alpha1.NDPSpec">
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
<h3 id="machine.onmetal.de/v1alpha1.NumaSpec">NumaSpec
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha1.InventorySpec">InventorySpec</a>)
</p>
<div>
<p>NumaSpec describes NUMA node</p>
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
<a href="#machine.onmetal.de/v1alpha1.MemorySpec">
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
<h3 id="machine.onmetal.de/v1alpha1.PCIDeviceDescriptionSpec">PCIDeviceDescriptionSpec
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha1.PCIDeviceSpec">PCIDeviceSpec</a>)
</p>
<div>
<p>PCIDeviceDescriptionSpec contains one of the options that is describing the PCI device</p>
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
<h3 id="machine.onmetal.de/v1alpha1.PCIDeviceSpec">PCIDeviceSpec
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha1.InventorySpec">InventorySpec</a>)
</p>
<div>
<p>PCIDeviceSpec contains description of PCI device</p>
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
<a href="#machine.onmetal.de/v1alpha1.PCIDeviceDescriptionSpec">
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
<a href="#machine.onmetal.de/v1alpha1.PCIDeviceDescriptionSpec">
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
<a href="#machine.onmetal.de/v1alpha1.PCIDeviceDescriptionSpec">
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
<a href="#machine.onmetal.de/v1alpha1.PCIDeviceDescriptionSpec">
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
<a href="#machine.onmetal.de/v1alpha1.PCIDeviceDescriptionSpec">
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
<a href="#machine.onmetal.de/v1alpha1.PCIDeviceDescriptionSpec">
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
<a href="#machine.onmetal.de/v1alpha1.PCIDeviceDescriptionSpec">
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
<h3 id="machine.onmetal.de/v1alpha1.PartitionSpec">PartitionSpec
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha1.PartitionTableSpec">PartitionTableSpec</a>)
</p>
<div>
<p>PartitionSpec contains info about partition</p>
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
<h3 id="machine.onmetal.de/v1alpha1.PartitionTableSpec">PartitionTableSpec
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha1.BlockSpec">BlockSpec</a>)
</p>
<div>
<p>PartitionTableSpec contains info about partition table on block device</p>
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
<a href="#machine.onmetal.de/v1alpha1.PartitionSpec">
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
<h3 id="machine.onmetal.de/v1alpha1.Size">Size
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
<a href="#machine.onmetal.de/v1alpha1.SizeSpec">
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
<a href="#machine.onmetal.de/v1alpha1.ConstraintSpec">
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
<a href="#machine.onmetal.de/v1alpha1.SizeStatus">
SizeStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="machine.onmetal.de/v1alpha1.SizeSpec">SizeSpec
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha1.Size">Size</a>)
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
<a href="#machine.onmetal.de/v1alpha1.ConstraintSpec">
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
<h3 id="machine.onmetal.de/v1alpha1.SizeStatus">SizeStatus
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha1.Size">Size</a>)
</p>
<div>
<p>SizeStatus defines the observed state of Size.</p>
</div>
<h3 id="machine.onmetal.de/v1alpha1.SystemSpec">SystemSpec
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha1.InventorySpec">InventorySpec</a>)
</p>
<div>
<p>SystemSpec contains DMI system information</p>
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
<h3 id="machine.onmetal.de/v1alpha1.ValidationInventory">ValidationInventory
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
<a href="#machine.onmetal.de/v1alpha1.InventorySpec">
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
<a href="#machine.onmetal.de/v1alpha1.SystemSpec">
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
<a href="#machine.onmetal.de/v1alpha1.IPMISpec">
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
<a href="#machine.onmetal.de/v1alpha1.BlockSpec">
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
<a href="#machine.onmetal.de/v1alpha1.MemorySpec">
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
<a href="#machine.onmetal.de/v1alpha1.CPUSpec">
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
<a href="#machine.onmetal.de/v1alpha1.NumaSpec">
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
<a href="#machine.onmetal.de/v1alpha1.PCIDeviceSpec">
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
<a href="#machine.onmetal.de/v1alpha1.NICSpec">
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
<a href="#machine.onmetal.de/v1alpha1.VirtSpec">
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
<a href="#machine.onmetal.de/v1alpha1.HostSpec">
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
<a href="#machine.onmetal.de/v1alpha1.DistroSpec">
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
<h3 id="machine.onmetal.de/v1alpha1.VirtSpec">VirtSpec
</h3>
<p>
(<em>Appears on:</em><a href="#machine.onmetal.de/v1alpha1.InventorySpec">InventorySpec</a>)
</p>
<div>
<p>VirtSpec contains info about detected host virtualization</p>
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
on git commit <code>1aa370a</code>.
</em></p>
