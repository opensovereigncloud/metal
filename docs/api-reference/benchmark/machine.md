<p>Packages:</p>
<ul>
<li>
<a href="#benchmark.onmetal.de%2fv1alpha3">benchmark.onmetal.de/v1alpha3</a>
</li>
</ul>
<h2 id="benchmark.onmetal.de/v1alpha3">benchmark.onmetal.de/v1alpha3</h2>
Resource Types:
<ul></ul>
<h3 id="benchmark.onmetal.de/v1alpha3.Benchmark">Benchmark
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
</td>
</tr>
</tbody>
</table>
<h3 id="benchmark.onmetal.de/v1alpha3.Deviation">Deviation
</h3>
<p>
(<em>Appears on:</em><a href="#benchmark.onmetal.de/v1alpha3.MachineStatus">MachineStatus</a>)
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
<code>disks</code><br/>
<em>
<a href="#benchmark.onmetal.de/v1alpha3.DiskDeviation">
[]DiskDeviation
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>networks</code><br/>
<em>
<a href="#benchmark.onmetal.de/v1alpha3.NetworkDeviation">
[]NetworkDeviation
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="benchmark.onmetal.de/v1alpha3.DiskDeviation">DiskDeviation
</h3>
<p>
(<em>Appears on:</em><a href="#benchmark.onmetal.de/v1alpha3.Deviation">Deviation</a>)
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
<p>Name contains full device name (like &ldquo;/dev/hda&rdquo; etc)</p>
</td>
</tr>
<tr>
<td>
<code>value</code><br/>
<em>
<a href="#benchmark.onmetal.de/v1alpha3.DiskValue">
[]DiskValue
</a>
</em>
</td>
<td>
<p>Results contains disk benchmark results.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="benchmark.onmetal.de/v1alpha3.DiskValue">DiskValue
</h3>
<p>
(<em>Appears on:</em><a href="#benchmark.onmetal.de/v1alpha3.DiskDeviation">DiskDeviation</a>)
</p>
<div>
<p>DiskValue contains block (device) changes.</p>
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
<code>ioPattern</code><br/>
<em>
string
</em>
</td>
<td>
<p>IOPattern defines type of I/O pattern (like &ldquo;read/write/readwrite&rdquo; etc)
more types could be found here: <a href="https://fio.readthedocs.io/en/latest/fio_doc.html#cmdoption-arg-readwrite">https://fio.readthedocs.io/en/latest/fio_doc.html#cmdoption-arg-readwrite</a></p>
</td>
</tr>
<tr>
<td>
<code>smallBlockReadIops</code><br/>
<em>
string
</em>
</td>
<td>
<p>SmallBlockReadIOPS contains benchmark result for read IOPS with small block size (device specified block size)</p>
</td>
</tr>
<tr>
<td>
<code>smallBlockWriteIops</code><br/>
<em>
string
</em>
</td>
<td>
<p>SmallBlockWriteIOPS contains benchmark result for write IOPS with small block size (device specified block size)</p>
</td>
</tr>
<tr>
<td>
<code>bandwidthReadIops</code><br/>
<em>
string
</em>
</td>
<td>
<p>BandwidthReadIOPS contains benchmark result for read IOPS with large block size (much larger then device specified block size)</p>
</td>
</tr>
<tr>
<td>
<code>bandwidthWriteIops</code><br/>
<em>
string
</em>
</td>
<td>
<p>BandwidthWriteIOPS contains benchmark result for write IOPS with large block size (much larger then device specified block size)</p>
</td>
</tr>
</tbody>
</table>
<h3 id="benchmark.onmetal.de/v1alpha3.Machine">Machine
</h3>
<div>
<p>Machine is the Schema for the machines API</p>
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
<a href="#benchmark.onmetal.de/v1alpha3.MachineSpec">
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
<code>benchmarks</code><br/>
<em>
<a href="#benchmark.onmetal.de/v1alpha3.[]..Benchmark">
map[string][]..Benchmark
</a>
</em>
</td>
<td>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#benchmark.onmetal.de/v1alpha3.MachineStatus">
MachineStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="benchmark.onmetal.de/v1alpha3.MachineSpec">MachineSpec
</h3>
<p>
(<em>Appears on:</em><a href="#benchmark.onmetal.de/v1alpha3.Machine">Machine</a>)
</p>
<div>
<p>MachineSpec contains machine benchmark results.</p>
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
<a href="#benchmark.onmetal.de/v1alpha3.[]..Benchmark">
map[string][]..Benchmark
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="benchmark.onmetal.de/v1alpha3.MachineStatus">MachineStatus
</h3>
<p>
(<em>Appears on:</em><a href="#benchmark.onmetal.de/v1alpha3.Machine">Machine</a>)
</p>
<div>
<p>MachineStatus contains machine benchmarks deviations.</p>
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
<code>deviation</code><br/>
<em>
<a href="#benchmark.onmetal.de/v1alpha3.Deviation">
Deviation
</a>
</em>
</td>
<td>
<p>Deviation shows the difference between last and current benchmark results.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="benchmark.onmetal.de/v1alpha3.NetworkDeviation">NetworkDeviation
</h3>
<p>
(<em>Appears on:</em><a href="#benchmark.onmetal.de/v1alpha3.Deviation">Deviation</a>)
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
<p>Name defines a name of network device</p>
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
<p>Results contains disk benchmark results.</p>
</td>
</tr>
</tbody>
</table>
<hr/>
<p><em>
Generated with <code>gen-crd-api-reference-docs</code>
on git commit <code>674c389</code>.
</em></p>
