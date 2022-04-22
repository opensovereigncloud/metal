# Inventory Concept

Inventory is a controller that allows to store results of machine inventarization process done with [inventory CLI](https://github.com/onmetal/inventory) in a form of k8s resource.

Inventory implements corresponding resource specification, controller and golang client for it.

### Libraries

Apart from libraries required to developm the operator itself, project has some other includes:

- [messagediff](https://github.com/d4l3k/messagediff) - used to compute the diff for resources on update 
