#!/usr/bin/env bash

# shellcheck disable=SC1091  
. /usr/local/bin/wait-for-ha-api.sh

wait_for_etcd_ha() {
    # check topology is not HighlyAvailable or not
    is_topology_ha
    retcode=$?
    if [[ $retcode -eq 2 ]]; then
        echo "topology is not HighlyAvailable, no need to wait for API availability"
        return 0
    fi
    if [[ $retcode -eq 0 ]]; then
        ## HA topology, we can start the wait loop for API availability
        "topology is HighlyAvailable, start to check etcd ha"
    fi
    # etcdctl container id
    local etcdctl_container_id
    while :
    do
        
        # Get etcdctl container
        etcdctl_container_id=$(crictl ps --name etcdctl -q)
        if [[ -z "$etcdctl_container_id" ]]; then
            echo "Failed to get etcdctl container"
            continue
        fi

        # Get member list and check for bootstrap members
        MEMBER_LIST=$(sudo crictl exec "$etcdctl_container_id"  etcdctl member list 2>/dev/null || true)
        if [[ -z "$MEMBER_LIST" ]]; then
            echo "Unable to get member list, retrying..."
            sleep 5
            continue
        fi

         # Check if any member contains "bootstrap" in the name
        BOOTSTRAP_MEMBER_ID=$(echo "$MEMBER_LIST" | grep -i bootstrap | awk '{print $1}'| tr -d ',' || true)          
        if [[ -z "$BOOTSTRAP_MEMBER_ID" ]]; then
            # unexpected error
            echo "No bootstrap member found in cluster" 
            return 1
        fi
        # counting health members（should be member, not just learner）
        # this check should consider the 
        # https://github.com/openshift/cluster-etcd-operator/blob/main/pkg/operator/bootstrapteardown/bootstrap_teardown_controller.go#L173
        HEALTHY_COUNT=$(echo "$MEMBER_LIST" | grep -c 'false' || true)
        if [[ "$HEALTHY_COUNT" -eq 4 ]]; then
            echo "All etcd members are ready"
        else 
            echo "Waiting for all members to be ready"
            sleep 5
            continue
        fi

        # Remove bootstrap etcd member from etcd cluster
        sudo crictl exec "$etcdctl_container_id" etcdctl member remove "$BOOTSTRAP_MEMBER_ID"
        if [[ $? -eq 0 ]]; then
            echo "Successful to remove bootstrap member"
        else
            # very rare, but we should double check if it's really deleted
            echo "Failed to remove bootstrap member"
            BOOTSTRAP_MEMBER=$(sudo crictl exec "$etcdctl_container_id" etcdctl member list | grep -i bootstrap || true)
            if [[ -z "$BOOTSTRAP_MEMBER" ]]; then
                echo "No bootstrap member found in cluster" 
            else
                echo "Found the bootstrap member"
                continue
            fi 
        fi

        # Remove that members static pod
        rm -f /etc/kubernetes/manifests/etcd-member-pod.yaml
        echo "etcd static pod successfully removed"
        return 0
    done
}
