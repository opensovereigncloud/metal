# Machine Concept

## Resources

Machine is a custom resource which provides abstract knowledge about underlay server.

Resource has a `status`, that may have `healthy` or `unhealthy` value.
If resource has an `unhealthy` state, that's mean that one or both related objects `inventory` and `oob` are missing.
If resource has been processed successfully, i.e. precessing has been `healthy`, means that corresponding machine might be used for booking.

## Machine

Machine is a top level resource that relays on `inventory` and `oob`.

It contains information about location and interface status of corresponding server and provides data about status of that server.

Machine exist to gave possibility order specific server without direct interaction with them.
