# Usage

## Resources

Machine controller managed one custom resource: `Machine`.


Resource has a `status`, that may have `healthy` or `unhealthy` value.
If resource has an `unhealthy` state, that's mean that one or both related objects `inventory` and `oob` are missing.
If resource has been processed successfully, i.e. precessing has been `healthy`, means that corresponding machine might be used for booking.

## Machine

Machine is a top level resource that relays on `inventory` and `oob`.

It contains information about location and interface status of corresponding server and provides data about status of that server.

Also Machine CR is used for server manipulation like power manipulation. And it has specific fields for that.

A proper Machine CR should be formed following the rules below. 

```yaml
apiVersion: machine.onmetal.de/v1alpha1
kind: Machine
metadata:
  name: machine-sample
spec:
  location:
    // Datacenter - name of building where machine lies
	// Optional
    // String
	datacenter: ""
	// DataHall - name of room in Datacenter where machine lies
	// Optional
    // String
	dataHall: ""
	// Shelf - defines place for server in Datacenter (an alternative name of Rack)
	// Optional
	// String
	shelf: ""
	// Slot - defines switch location in rack (an alternative name for Row)
	// Optional
	// String
	slot: "" 
	// HU - is a unit of measure defined 44.45 mm
	// Optional
	// String
	hu: ""
	// Row - switch location in rack
	// Optional
	// Int16
	row: 0
	// Rack - is a place for server in DataCenter
	// Optional
	// Int16
	rack: 0
  identity:
  	// SKU - stock keeping unit. The label allows vendors automatically track the movement of inventory
	// Optional
	// String
	sku: ""
	// SerialNumber - unique machine number
	// Optional
	// String
	serial_number: ""
	// Optional
	// String
	asset: ""

	// ScanPorts - trigger manual port scan
	// Bool
  scan_ports: false
  action:
	// PowerState - defines desired machine power state
	// Optional
	// Validation Pattern=`^(?:On|Reset|ResetImmediate|Off|OffImmediate)$`
	// String
	power_state: ""
```