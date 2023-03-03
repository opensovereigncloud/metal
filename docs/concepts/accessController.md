Access Controller Concepts

Access Controller is a component of Metal-API that creates, maintains, and supports 
service accounts and secrets for every sized Inventory in cluster.
It is a crucial component for maintaining a secure and efficient way to obtain 
server Inventory in cluster.

The main purpose of Access Controller is to ensure that every server in the Kubernetes cluster 
can update its own inventory into the cluster. It achieves this by providing the 
necessary access controls and permissions for server. It creates service accounts, secrets, bindings and kubeconfig.

Kubeconfig is a configuration file that contains authentication and authorization information for a Kubernetes cluster. 
The Access Controller provides that file for every observed server.

The Access Controller is an important component Metal-API that provides 
the necessary access controls and permissions for service accounts, secrets, and kubeconfig. 
It ensures that every server in the cluster can update its own Inventory.