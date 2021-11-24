# Usage

## Use
`./config/samples/` directory contains examples of manifests. They can be used to try out the controller.

    # apply config
    kubectl apply -k config/samples/
    # get resources
    kubectl get switches.machine.onmetal.de -n onmetal
    # get sample resource
    kubectl describe switches.switch.onmetal.de -n onmetal 92b9de0f-19f2-3f3b-95d0-fb668b1d3d3b

## Clean
Clean up local environment, after you've done with development:

    # generate deployment config and delete corresponding entities
    kustomize build config/default | kubectl delete -f -
    # remove registry bridge
    docker stop registry-bridge
    # stop minikube
    minikube stop

## Test 
To run tests:

    make test