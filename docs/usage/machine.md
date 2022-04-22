# Machine Usage

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