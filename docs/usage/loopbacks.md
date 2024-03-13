# Loopback Usage

To create loopbacks, it is necessary to prepare the environment by configuring loopback subnets. Below are examples of loopback subnet configurations for both IPv4 and IPv6:

## IPv4

For IPv4 loopback subnets, you need to create a subnet with a specified CIDR range and label. 

Here's an example in YAML:
```yaml
apiVersion: ipam.metal.ironcore.dev/v1alpha1
kind: Subnet
metadata:
  name: ipv4-cidr-subnet-sample
  namespace: metal
  labels:
    "loopback": "loopback"
spec:
  cidr: "10.0.0.0/16"
  network:
    name: "test-looback"

```

In this example, we've defined an IPv4 loopback subnet with a CIDR range of 10.0.0.0/16 and labeled it as "loopback." You can customize the CIDR range as needed.

### IPv6

For IPv6 loopback subnets, you should create a parent loopback subnet with a specified CIDR range and label. Based on this parent subnet, new subnets with a prefix size of /64 will be created. 

Here's an example in YAML:
```yaml
apiVersion: ipam.metal.ironcore.dev/v1alpha1
kind: Subnet
metadata:
  name: ipv6-cidr-subnet-sample
  namespace: metal
  labels:
    "loopback": "loopback"
spec:
  cidr: "1a10:afc0:e003:0000::/52"
  network:
    name: "test-looback"
```
In this example, we've defined an IPv6 parent loopback subnet with a CIDR range of 1a10:afc0:e003:0000::/52 and labeled it as "loopback." Subnets with a `/64` prefix size will be created based on this parent subnet.


### Configuration

The label value used for loopback subnets, such as "loopback" in the examples above, can be configured using a flag - `loopback_subnet_value_name`. This allows you to customize the label according to your specific requirements.
