#!/usr/bin/env bash


#[root@localhost manifests]# crictl exec -it $(crictl ps --name etcdctl -q)  etcdctl member list
#29b43fbdc0a560b5, started, master0, https://10.x.x.100:2380, https://10.x.x.100:2379, false
#65058609be25f3a3, started, etcd-bootstrap, https://10.x.x.10:2380, https://10.x.x.10:2379, false
#dd70d9cc2bef90e1, started, master1, https://10.x.x.101:2380, https://10.x.x.101:2379, false
#f3ee2af0de233621, started, master2, https://10.x.x.102:2380, https://10.x.x.102:2379, false

wait_for_etcd_ha() {
    while :
    do
        # check topology is not HighlyAvailable or not
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
            "topology is HighlyAvailable, start to check etcd ha"
        fi
        # get member list and check for bootstrap members
        MEMBER_LIST=$(sudo crictl exec -it $(crictl ps --name etcdctl -q)  etcdctl member list 2>/dev/null || true)
        if [ -z "$MEMBER_LIST" ]; then
            echo "Unable to get member list, retrying..."
            sleep 5
            continue
        fi

         # Check if any member contains "bootstrap" in the name
        BOOTSTRAP_MEMBERS=$(echo "$MEMBER_LIST" | grep -i bootstrap || true)          
        if [ -z "$BOOTSTRAP_MEMBERS" ]; then
            # unexpected error
            echo "No bootstrap member found in cluster" 
            return 1
        
        # counting health members（status is started）
        # this check should consider the 
        # https://github.com/openshift/cluster-etcd-operator/blob/main/pkg/operator/bootstrapteardown/bootstrap_teardown_controller.go#L173
        HEALTHY_COUNT=$(echo "$MEMBER_LIST" | grep -c 'started')
        if [ "$HEALTHY_COUNT" -eq 4 ]; then
            echo "All etcd members are ready"
        else 
            echo "Waiting for all members to be ready"
            sleep 5
            continue

        # remove that members static pod
        rm -f /etc/kubernetes/manifests/etcd-member-pod.yaml
        echo "etcd static pod successfully removed"
        return 0
    done
}
