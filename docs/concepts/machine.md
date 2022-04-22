# Machine Concept

## Resources

Machine controller managed one custom resource: `Machine`.


Resource has a `status`, that may have `healthy` or `unhealthy` value.
If resource has an `unhealthy` state, that's mean that one or both related objects `inventory` and `oob` are missing.
If resource has been processed successfully, i.e. precessing has been `healthy`, means that corresponding machine might be used for booking.

## Machine

Machine is a top level resource that relays on `inventory` and `oob`.

It contains information about location and interface status of corresponding server and provides data about status of that server.

Also Machine CR is used for server manipulation like power manipulation. And it has specific fields for that.
