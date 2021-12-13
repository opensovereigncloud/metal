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
<h3 id="benchmark.onmetal.de/v1alpha3.BenchmarkDeviation">BenchmarkDeviation
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
string
</em>
</td>
<td>
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
<code>machine_deviation</code><br/>
<em>
<a href="#benchmark.onmetal.de/v1alpha3.[]..BenchmarkDeviation">
map[string][]..BenchmarkDeviation
</a>
</em>
</td>
<td>
<p>MachineDeviation shows the difference between last and current benchmark results.</p>
</td>
</tr>
</tbody>
</table>
<hr/>
<p><em>
Generated with <code>gen-crd-api-reference-docs</code>
on git commit <code>172f4e7</code>.
</em></p>
