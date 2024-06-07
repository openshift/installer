#!/usr/bin/env bash

wait_for_ha_api() {
    while :
    do
        is_topology_ha
        retcode=$?
        if [[ $retcode -eq 2 ]]
        then
            echo "topology is not HighlyAvailable, no need to wait for API availability"
            return 0
        fi
        if [[ $retcode -eq 0 ]]
        then
            ## HA topology, we can start the wait loop for API availability
            break
        fi

        ## error happened, so let's retry after 5s
        sleep 5
    done

    echo "Waiting for at least 2 available IP addresses for the default/kubernetes service"
    while ! is_api_available
    do
        sleep 5
    done
}

## 0 - HA control plane 'HighlyAvailable'
## 1 - error condition
## 2 - other topology
is_topology_ha() {
    output=$(oc --kubeconfig="$KUBECONFIG" get infrastructures cluster -o jsonpath='{.status.controlPlaneTopology}' 2>&1 )
    # shellcheck disable=SC2124
    status=$?
    if [[ $status -ne 0 ]]
    then
        echo "The following error happened while retrieving infrastructures/cluster object"
        echo "$output"
        return 1 # unexpected error condition
    fi

    if [[ -z $output ]]
    then
        echo "status.infrastructureTopology of the infrastructures/cluster object is empty"
        return 1 # unexpected error condition
    fi

    if [[ $output == "HighlyAvailable" ]]
    then
        return 0 ## HA control plane
    fi

    return 2 ## non HA control plane
}

##
## for HA cluster, we mark the bootstrap process as complete when there
## are at least two IP addresses available to the endpoints
## of the default/kubernetes service object.
## TODO: move this to kas operator as a subcommand of the render command
is_api_available() {
    output=$(oc --kubeconfig="$KUBECONFIG" get endpoints kubernetes --namespace=default -o jsonpath='{range @.subsets[*]}{range @.addresses[*]}{.ip}{" "}' 2>&1 )
    # shellcheck disable=SC2124
    status=$?
    if [[ $status -ne 0 ]]
    then
	echo "The following error happened while retrieving the default/kubernetes endpoint object"
	echo "$output"
	return 1
    fi
    
    echo "Got the following addresses for the default/kubernetes endpoint object: $output"
    count=$(echo "$output" | wc -w)
    if [[ ! $count -gt 1 ]]
    then
	return 1
    fi
    
    echo "Got at least 2 available addresses for the default/kubernetes service"
    return 0
}
