#!/bin/bash -x

function patchit {
    # allow etcd-operator to start the etcd cluster without minimum of 3 master nodes
    oc --kubeconfig ./auth/kubeconfig patch etcd cluster -p='{"spec": {"unsupportedConfigOverrides": {"useUnsupportedUnsafeNonHANonProductionUnstableEtcd": true}}}' --type=merge || return 1

    # allow cluster-authentication-operator to deploy OAuthServer without minimum of 3 master nodes
    oc --kubeconfig ./auth/kubeconfig patch authentications.operator.openshift.io cluster -p='{"spec": {"managementState": "Managed", "unsupportedConfigOverrides": {"useUnsupportedUnsafeNonHANonProductionUnstableOAuthServer": true}}}' --type=merge || return 1

    # Mark etcd-quorum-guard as unmanaged
    oc --kubeconfig ./auth/kubeconfig patch clusterversion/version --type='merge' -p "$(cat <<- EOF
 spec:
    overrides:
      - group: apps/v1
        kind: Deployment
        name: etcd-quorum-guard
        namespace: openshift-etcd
        unmanaged: true
EOF
)" || return 1

    # scale down etcd-quorum-guard
    oc --kubeconfig ./auth/kubeconfig scale --replicas=0 deployment/etcd-quorum-guard -n openshift-etcd || return 1

    return 0
}

while ! patchit; do
    echo "Waiting to try again..."
    sleep 10
done
touch patch.done

