# Switch HOW-TOs

Switches configuration depends on several custom resources. Important to keep in mind that it depends not only on the objects of particular types (kinds), but also on their proper reconciliation. Hence, it is not enough to only add corresponding custom resource definitions to the cluster and create required objects: it is also required to deploy operators, which reconcile these objects.  

The complete list of custom resources switches are directly depend on is the following:

- Group: `metal.ironcore.dev`
  Version: `v1alpha1`
  Kind: `Inventory`
- Group: `metal.ironcore.dev`
  Version: `v1beta1`
  Kind: `SwitchConfig`
- Group: `ipam.ironcore.dev`
  Version: `v1alpha1`
  Kind: `Subnet`
- Group: `ipam.ironcore.dev`
  Version: `v1alpha1`
  Kind: `IP`

Switches are indirectly depended on the following resources:

- Group: `metal.ironcore.dev`
  Version: `v1alpha1`
  Kind: `Size`
- Group: `ipam.ironcore.dev`
  Version: `v1alpha1`
  Kind: `Network`

Two mentioned above kinds don't  affect switches, but at the same time without `Size` objects inventories will not be automatically labeled as ones related to switches. 
The reference to the `Network` object is mandatory and should be placed in the `Subnet` object's manifest.

In total, aside from mentioned custom resource definitions, the following operators have to be deployed:

- `metal-api`, includes controllers, which process `Size`, `Inventory`, `SwitchConfig`, `Switch`
- `ipam`, includes controllers, which process `Network`, `Subnet`, `IP`

## Dependencies description

- `Size`: these objects are used to properly label `Inventory` objects. It is required that `Inventory` objects, containing switches data, are labeled with `metal.ironcore.dev/size-switch: ""` label. 
Here is an example of `Size` object:

    ```yaml
    apiVersion: metal.ironcore.dev/v1alpha4
    kind: Size
    metadata:
      name: switch
      namespace: metal-api
    spec:
      constraints:
      - eq: broadcom
        path: spec.distro.asicType
    ```

- `Inventory`: these objects are initially the source of truth about switches' interfaces and discovered neighbours. `Switch` objects, which either has no corresponding inventory or does not contain reference to the inventory object will not be configured.
- `SwitchConfig`: these objects contain configuration, common for user-defined switches' types. Example of `SwitchConfig` object:

    ```yaml
    apiVersion: metal.ironcore.dev/v1beta1
    kind: SwitchConfig
    metadata:
      name: spines-config
      namespace: metal
      labels:
        metal.ironcore.dev/layer: "0"
    spec:
      switches:
        matchLabels:
          metal.ironcore.dev/type: "spine"
      portsDefaults:
        lanes: 4
        mtu: 9100
        ipv4MaskLength: 30
        ipv6Prefix: 112
        state: "up"
        fec: "rs"
      ipam:
        addressFamily:
          ipv4: true
          ipv6: true
        carrierSubnets:
          labelSelector:
            matchLabels:
              ipam.ironcore.dev/object-purpose: "switch-carrier"
        loopbackSubnets:
          labelSelector:
            matchLabels:
              ipam.ironcore.dev/object-purpose: "switch-loopbacks"
        southSubnets:
          labelSelector:
            matchLabels:
              ipam.ironcore.dev/object-purpose: "south-subnet"
          fieldSelector:
            labelKey: "ipam.ironcore.dev/object-owner"
            fieldRef:
              fieldPath: "metadata.name"
        loopbackAddresses:
          labelSelector:
            matchLabels:
              ipam.ironcore.dev/object-purpose: "loopback"
          fieldSelector:
            labelKey: "ipam.ironcore.dev/object-owner"
            fieldRef:
              fieldPath: "metadata.name"
    ```

> :heavy_exclamation_mark: Details about how to set up mapping between `Switch` and `Switchconfig` objects can be found in [concepts](../concepts/switch.md#mapping-between-switch-and-switchconfig-objects)

- `Network`: it is a mandatory object which `Subnet` objects have to reference to;
- `Subnet`: these objects are required to be assigned to switches as their subnets. Switch ports' IP addresses are reserved in these subnets;
- `IP`: these objects require to be assigned to switches as they're loopback addresses. It is also required for switch to have IPv4 loopback address, since it is used to compute switch's ASN for BGP routing;

To learn more about IPAM objects, please refer to the [documentation](https://github.com/ironcore-dev/ipam/blob/main/docs/usage.md)

> :heavy_exclamation_mark: To be consumed by the switch-controller, `Subnet` and `IP` objects have to have labels which match to the selectors defined in corresponding `.spec.ipam` section of `Switch` and `SwitchConfig` objects.

## How to create Switch object

Basically, `onboarding-controller` will create `Switch` object, which is originally based on `Inventory` object, so there is no need to create `Switch` object manually. However, it is possible to pre-create `Switch` objects. In this case `.metadata.name` of the `Switch` object has to be the same as the name of future `Inventory` object.  
`Inventory` object's name is always UUIDv4 and for regular servers can be found in `/sys/class/dmi/id/produc_uuid`. In case of SONiC switches, all of them have the same UUID equals to **03000200-0400-0500-0006-000700080009**, hence it's required to compute UUID. The computation is quite straightforward and consists of two steps:

1. use `uuidgen` command-line tool to generate UUID for namespace:

    ```shell
    > uuidgen --md5 --namespace 00000000-0000-0000-0000-000000000000 --name "ironcore.dev"
    ```

2. next generate UUID using switch's serial number as "name" parameter and UUID generated on previous step as namespace:

    ```shell
    > uuidgen --md5 --namespace 32e12471-57b7-3cbf-b139-033116e1d361 --name <serial number>
    ```

SONiC switch serial number can be found in the file `/sys/class/dmi/id/product_serial` or by using `dmidecode -t system` command.  
After `Switch` object is created it might have no type label - **"metal.ironcore.dev/type"** - which is required to proceed with configuration processing. In this case object's state will be set to `Pending` and there would be the message about the issue which blocks reconciliation.

## How to define connection hierarchy

Hierarchy of switches and machines interconnections is computed automatically, based on LLDP data of each switch port. However, it is required to explicitly define which switches are the top spines. To do this, set the `.spec.topSpine` flag equals to `true` for all switches which are considered to be top spines.

## How to override common switch ports' parameters

Switch port parameters are defined in several places:

- parameters which are common for every switch of certain type are defined in `SwitchConfig` object's `.spec.portsDefaults` field;
- parameters which should be defaults for certain switch are defined in `Switch` object's `spec.interfaces.defaults` field;
- parameters' overrides for certain switch port could be defined in `Switch` object's `spec.interfaces.overrides` field;

During switch reconciliation all parameters from above are merged into resulting parameters which are applied to each switch port. 

The hierarchy of importance:

- Overrides
- Switch-specific defaults
- Global default defined in `SwitchConfig` object

As an example: `SwitchConfig` contains MTU value equals to 1500, MTU value is not defined in `Switch` defaults, but there is an override of MTU value equals to 9100 for switch port "Ethernet120". In this case for all switch ports will be applied MTU equals to 1500, except for switch port "Ethernet120" which will get MTU equals to 9100.

Switch ports can be logically separated into two groups:

- "south" interfaces, which reflect downstream connections;
- "north" interfaces, which reflects upstream connections;

Opposite to "south" interfaces, all "north" interfaces will inherit port parameters from upstream ports, and it is forbidden to set overrides for "north" interfaces.

## How to pin to IPAM objects

IPAM objects have to have proper labels, so the controller can get these objects and use them to configure switches. Labels could be either defined by selectors in `.spec.ipam` fields of `SwitchConfig` and `Switch` objects or default values could be used:

- **"ipam.ironcore.dev/object-purpose"** with values **"south-subnet"** (for `Subnet` object) or **"loopback"** (for `IP` object);
- **"ipam.ironcore.dev/object-owner"** with the switch object name as value;

Pay attention that labelSelector and fieldSelector in `.spec.ipam` configuration have to match labels of `Subnet` and `IP` objects which supposed to be used.

# Troubleshooting

Since `Switch` object reconciliation process depends on various resources and their state, reconciliation process may fail, literally on any of its steps. To make it easier to investigate such failures and not to rely only on operator's logs, `Switch` object status contains the list of items named "Conditions".

Each condition contains its name, state, last update and transition timestamp and, in case of failure, reason and message which are correspondingly short and verbose description of the reason of failure. Here are examples:

   In case of field `.spec.inventoryRef.name` is not filled:
   ```
   conditions:
     ...
     - lastTransitionTimestamp: 2023-03-02 17:06:22.708218085 +0000 UTC m=+525916.982520824
       lastUpdateTimestamp: 2023-03-02 17:06:22.708218085 +0000 UTC m=+525916.982520824
       name: InterfacesOK
       state: false
       reason: MissingRequirements
       message: "missing requirements: reference to corresponding Inventory at .spec.InventoryRef.name"
     ...
   ```
   In case of GET request to kube-apiserver for inventory object is failed for any reason except object does not exist:
   ```
   conditions:
     ...
     - lastTransitionTimestamp: 2023-03-02 17:06:22.708218085 +0000 UTC m=+525916.982520824
       lastUpdateTimestamp: 2023-03-02 17:06:22.708218085 +0000 UTC m=+525916.982520824
       name: InterfacesOK
       state: false
       reason: APIRequestFailed
       message: "failed to get requested object: Inventory"
     ...
   ```
   In case of GET request to kube-apiserver for inventory object is failed due to object does not exist:
   ```
   conditions:
     ...
     - lastTransitionTimestamp: 2023-03-02 17:06:22.708218085 +0000 UTC m=+525916.982520824
       lastUpdateTimestamp: 2023-03-02 17:06:22.708218085 +0000 UTC m=+525916.982520824
       name: InterfacesOK
       state: false
       reason: ObjectNotExist
       message: "requested object does not exist: Inventory"
     ...
   ```
   In case of `Switch` object does not have required label:
   ```
   conditions:
     ...
     - lastTransitionTimestamp: 2023-03-02 17:06:22.708294887 +0000 UTC m=+525916.982597609
       lastUpdateTimestamp: 2023-03-02 17:06:22.708294887 +0000 UTC m=+525916.982597609
       name: ConfigRefOK
       state: false
       reason: MissingRequirements
       message: "missing requirements: label metal.ironcore.dev/type"
     ...
   ```
   In case of GET request to kube-apiserver for switchconfig object is failed for any reason except object does not exist:
   ```
   conditions:
     ...
     - lastTransitionTimestamp: 2023-03-02 17:06:22.708218085 +0000 UTC m=+525916.982520824
       lastUpdateTimestamp: 2023-03-02 17:06:22.708218085 +0000 UTC m=+525916.982520824
       name: ConfigRefOK
       state: false
       reason: APIRequestFailed
       message: "failed to get requested object: SwitchConfig"
     ...
   ```
   In case of GET request to kube-apiserver for switchconfig object is failed due to object does not exist:
   ```
   conditions:
     ...
     - lastTransitionTimestamp: 2023-03-02 17:06:22.708218085 +0000 UTC m=+525916.982520824
       lastUpdateTimestamp: 2023-03-02 17:06:22.708218085 +0000 UTC m=+525916.982520824
       name: ConfigRefOK
       state: false
       reason: ObjectNotExist
       message: "requested object does not exist: SwitchConfig"
     ...
   ```
   In case of GET request to kube-apiserver for switchconfig object is failed for any reason:
   ```
   conditions:
     ...
     - lastTransitionTimestamp: 2023-03-02 17:06:22.708376164 +0000 UTC m=+525916.982678888
       lastUpdateTimestamp: 2023-03-02 17:06:22.708376164 +0000 UTC m=+525916.982678888
       name: PortParametersOK
       state: false
       reason: APIRequestFailed
       message: "failed to get requested object: SwitchConfig"
     ...
   ```
   In case of LIST request to kube-apiserver for switches list is failed for any reason:
   ```
   conditions:
     ...
     - lastTransitionTimestamp: 2023-03-02 17:06:22.708811087 +0000 UTC m=+525916.983113811
       lastUpdateTimestamp: 2023-03-02 17:06:22.708811087 +0000 UTC m=+525916.983113811
       name: NeighborsOK
       state: false
       reason: APIRequestFailed
       message: "failed to get requested object: SwitchList"
     ...
   ```
   In case of LIST request to kube-apiserver for switches list is failed for any reason:
   ```
   conditions:
     ...
     - lastTransitionTimestamp: 2023-03-02 17:06:22.709279735 +0000 UTC m=+525916.983582457
       lastUpdateTimestamp: 2023-03-02 17:06:22.709279735 +0000 UTC m=+525916.983582457
       name: LayerAndRoleOK
       state: false
       reason: APIRequestFailed
       message: "failed to get requested object: SwitchList"
     ...
   ```
   In case of LIST request to kube-apiserver for ips list is failed for any reason:
   ```
   conditions:
     ...
     - lastTransitionTimestamp: 2023-03-02 17:06:22.71269581 +0000 UTC m=+525916.986998617
       lastUpdateTimestamp: 2023-03-02 17:06:22.71269581 +0000 UTC m=+525916.986998617
       name: LoopbacksOK
       state: false
       reason: APIRequestFailed
       message: "failed to get requested object: IPList"
     ...
   ```
   In case of `IP` object, matching to loopback selectors, with address family IPv4 does not exist:
   ```
   conditions:
     ...
     - lastTransitionTimestamp: 2023-03-02 17:06:22.71269581 +0000 UTC m=+525916.986998617
       lastUpdateTimestamp: 2023-03-02 17:06:22.71269581 +0000 UTC m=+525916.986998617
       name: LoopbacksOK
       state: false
       reason: ObjectNotExist
       message: "missing requirements: IP object of V4 address family to be assigned to loopback interface"
     ...
   ```
   In case of parsing of IP address was failed:
   ```
   conditions:
     ...
     - lastTransitionTimestamp: 2023-03-02 17:06:22.712713075 +0000 UTC m=+525916.987015817
       lastUpdateTimestamp: 2023-03-02 17:06:22.712713075 +0000 UTC m=+525916.987015817
       name: AsnOK
       state: false
       reason: ASNCalculationFailed
       message: "failed to parse IP address: <address value>"
     ...
   ```
   In case of LIST request to kube-apiserver for subnets list is failed for any reason:
   ```
   conditions:
     ...
     - lastTransitionTimestamp: 2023-03-02 17:06:22.71269581 +0000 UTC m=+525916.986998617
       lastUpdateTimestamp: 2023-03-02 17:06:22.71269581 +0000 UTC m=+525916.986998617
       name: SubnetssOK
       state: false
       reason: APIRequestFailed
       message: "failed to get requested object: SubnetList"
     ...
   ```
   In case of `Subnet` object, matching to subnet selectors does not exist:
   ```
   conditions:
     ...
     - lastTransitionTimestamp: 2023-03-02 17:06:22.71269581 +0000 UTC m=+525916.986998617
       lastUpdateTimestamp: 2023-03-02 17:06:22.71269581 +0000 UTC m=+525916.986998617
       name: SubnetsOK
       state: false
       reason: ObjectNotExist
       message: "requested object does not exist: Subnet"
     ...
   ```
   In case of parsing of IP address or CIDR was failed:
   ```
   conditions:
     ...
     - lastTransitionTimestamp: 2023-03-02 17:06:22.717401875 +0000 UTC m=+525916.991704668
       lastUpdateTimestamp: 2023-03-02 17:06:22.717401875 +0000 UTC m=+525916.991704668
       name: IPAddressesOK
       state: false
       reason: IPAssignmentFailed
       message: "failed to parse IP address: <address value>" | "failed to parse CIDR: <CIDR value>"
     ...
   ```

