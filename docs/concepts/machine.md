# Machine Concept

The Machine concept is a custom resource that represents a server and provides abstract knowledge about its underlying infrastructure.

The Machine resource relies on two key components: `inventory` and `oob`.

It provides information about the server's location and interface status, allowing users to retrieve relevant data about the server's status without directly interacting with it.

The primary purpose of the Machine resource is to enable users to order a specific server without requiring direct interaction with the physical machine itself. By leveraging the abstracted knowledge provided by the Machine resource, users can conveniently select and book servers based on their specific requirements.

## Resources

A Machine resource has a status attribute, which indicates the overall health of the server. 
The status can have two values: `healthy` or `unhealthy`.

If the status is `unhealthy`, it means that one or both of the related objects, namely inventory and oob (out-of-band management), are missing or network interfaces are not redundant. In this state, the server may not be fully operational or accessible.

If the status is `healthy`, it signifies that the corresponding `machine` has been processed successfully.
The processing involves verifying the presence and functionality of both the `inventory` and `oob` components and also verifying is there two network interfaces connected to the different switches or not.
A `healthy` status implies that the `machine` is in a usable state and can be booked or utilized.