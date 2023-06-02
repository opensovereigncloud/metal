# Architecture

The Metal-API stack is composed of different parts, 
each responsible for specific aspects of the application's functionality and operations.

The purpose of this architecture is to provide a modular and scalable approach to managing the Metal-API stack. 
Each component plays a specific role in the overall system, ensuring efficient onboarding, switch and machine management, inventory maintenance. 

By dividing the application into distinct parts, the architecture allows for flexibility and extensibility, enabling smooth operation and maintenance of the Metal-API stack.

## Parts

The architecture of Metal-API involves the following key parts:

1. Onboarding Controllers: 
    These controllers handle the process of bringing up a new `machine` object from the out-of-band (oob) management system. The onboarding controllers create an empty `inventory` object and initiate the default server boot-up process. Once the server is booted up, the inventory tool generates comprehensive data about the server, which is then updated into the Kubernetes cluster. 

    Subsequently, the machine onboarding controller creates a `machine` object based on the corresponding inventory, provided it has a machine label associated with it.

2. Machine Controllers: 
    These controllers are responsible for the creation of machine pool objects, which are used for scheduling purposes. Additionally, machine controllers handle machine reservation and power control operations. When reserving a machine, a specialized IPXE controller generates a suitable configuration that will be used during server boot-up.

3. Inventory Controllers: 
    The inventory controllers - are responsible for maintaining and provisioning the inventory resources. 
    The size controller - updates labels on the inventory based on the specifications provided in the `size` object.
    The aggregate controller - calculates aggregate information in a concise and usable format for future use. 
    The access controller - creates a set of objects that allow a server to update its data securely.