# Description
This is a local deployment of argo workflows and argo events with a predeployed event bus.  There is an example workflow in workflows

## How to Use

    go build -o kubby-argo main.go
    ./kubby-argo

The binary should output what ip to go to and what to set KUBECONFIG to access the cluster
