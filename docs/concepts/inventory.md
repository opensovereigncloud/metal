# Inventory Concept

The Inventory controller is an essential component that facilitates the storage of server inventory results within the Metal-API stack. 
It leverages the inventory CLI to perform machine inventory processes and stores the obtained information as a Kubernetes resource.

## Resources

The Inventory concept involves the following key resources:

1. Inventory CLI: 
    The Inventory controller integrates with the inventory CLI tool, which is responsible for performing inventory generation. The inventory CLI gathers comprehensive data about the server, including hardware specifications,  configuration details, and other relevant information.

2. Inventory Resource: 
    The Inventory controller implements a corresponding resource specification within Kubernetes. This resource is specifically designed to store and manage the results of the server inventory process. It provides a standardized format for organizing and accessing the inventory data within the Metal-API stack.

3. Controllers: 
    The Inventory controller itself is responsible for managing the lifecycle of the inventory resources. It handles operations such as creating, updating, and deleting inventory objects based on the machine inventory results obtained from the inventory CLI. The controller ensures that the inventory data remains consistent and synchronized with the underlying machines.

4. Golang Client: 
    To interact with the Kubernetes cluster, a dedicated Golang client is provided. The client allows for programmatic access to the inventory resources, enabling automation and integration with other components of the Metal-API stack.

The Inventory controller facilitates the seamless storage and management of machine inventory data within the Metal-API stack by leveraging a dedicated Kubernetes resource and providing a controller for handling the inventory resources. This enables efficient tracking and utilization of server `inventory` information throughout the system.

## Controllers

1. The Inventory controller: 
    Responsible for the lifecycle management of inventory resources. It performs operations such as creating, updating, and deleting inventory objects based on the server inventory results obtained from the inventory CLI. The Inventory controller ensures that the inventory data remains consistent and synchronized with the underlying machines. By integrating with the inventory CLI, it enables seamless storage and retrieval of server inventory information within the Metal-API stack.

2. The Size controller: 
    Focuses on updating labels on the inventory based on the specifications provided in the size object. It ensures that the inventory resources are appropriately labeled to reflect their specific size-related characteristics. By updating the labels, the Size controller enables efficient categorization and filtering of the inventory resources based on size-related attributes. This facilitates easier querying and management of machines based on their size requirements.

3. The Aggregate controller: 
    Plays a vital role in calculating aggregate information in a concise and usable format for future use. It analyzes and processes data from various sources within the Metal-API stack to derive meaningful aggregated insights. By consolidating and summarizing relevant information, the Aggregate controller provides a high-level overview and valuable statistics about the server inventory. This aggregated information can be utilized for monitoring, reporting, and decision-making purposes.

4. The Access controller: 
    Is responsible for creating a set of objects that enable servers to update their data securely. It ensures that the necessary access controls and permissions are in place to facilitate secure interactions between the servers and the Metal-API stack. By creating the required objects, such as service accounts, secrets, bindings, and kubeconfig files, the Access controller enables authorized servers to update their inventory and associated data while maintaining the overall security and integrity of the system.

These controllers collectively contribute to the efficient and secure management of the Inventory. 
They handle various aspects such as inventory lifecycle, size labeling, data aggregation, and secure server interactions, ensuring smooth operations and reliable data management within the Metal-API ecosystem.