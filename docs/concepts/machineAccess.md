# Machine Access Concepts

Kubeconfig is a configuration file that contains authentication and authorization information for a Kubernetes cluster. 
The Access Controller provides that file for every observed server.

The Access Controller is an important component Metal-API that provides the necessary access controls and permissions for service accounts, secrets, and kubeconfig. It ensures that every server in the cluster can update its own Inventory.


The Access Controller is a crucial component of Metal-API that plays a vital role in creating, maintaining, and supporting service accounts and secrets for every server in the cluster, regardless of its size.

It serves as a fundamental mechanism for ensuring a secure and efficient way to obtain server Inventory within the Kubernetes cluster.

## Resources

The main purpose of Access Controller is to ensure that every server in the Kubernetes cluster can update its own inventory into the cluster. It achieves this by providing the necessary access controls and permissions for server. It creates service accounts, secrets, bindings and kubeconfig.

The Access Controller concept involves the following key resources:

1. Service Accounts: 
    The Access Controller creates and manages service accounts for each server in the cluster. Service accounts enable secure authentication and authorization for server operations within the cluster.

2. Secrets: 
    The Access Controller is responsible for handling secrets, which contain sensitive information such as authentication tokens, passwords, or certificates. It securely manages and provides access to these secrets for the servers in the cluster.

3. Bindings: 
    To establish the necessary connections between service accounts and their associated permissions, the Access Controller creates and maintains bindings. Bindings define the relationships between service accounts and the resources they are authorized to access.

4. Kubeconfig: 
    The Access Controller generates a kubeconfig file for each observed server. Kubeconfig is a configuration file that contains the necessary authentication and authorization information to interact with a Kubernetes cluster. By providing kubeconfig files, the Access Controller enables servers to securely update their own Inventory within the cluster.

The primary purpose of the Access Controller is to ensure that every server in the Kubernetes cluster has the appropriate access controls and permissions to manage its own Inventory. 
It accomplishes this by creating and managing service accounts, secrets, bindings, and kubeconfig files. 
This enables secure and efficient server interactions within the cluster while maintaining the overall security and integrity of the system.