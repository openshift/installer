#!/usr/bin/env bash

wait_for_ha_api() {
    if [ "$BOOTSTRAP_INPLACE" = true ]
    then
	return 0
    fi

    echo "Waiting for at least 2 available IP addresses for the default/kubernetes service"
    while ! is_api_available
    do
	sleep 5
    done
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
