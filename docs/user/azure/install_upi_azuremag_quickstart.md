# Azure MAG UPI quickstart

1. run ```openshift-install create install-config```
    - input valid azure public information (this won't be used)
2. Modify install-config with the correct values for azure government
    - platform.azure.region
    - platform.azure.baseDomainResourceGroupName
    - basedomain
3. copy all ARM templates from $CODE_LOCATION/upi/azure
4. copy azureGovQuickstart.sh from $CODE_LOCATION/upi/azure
5. run azureGovQuickstart.sh
    - create NS records in parent DNS zone pointing to child zone this will resolve the error of ```lookup [api.clustername] on [public ip] no such host```
6. run ```export KUBECONFIG="$PWD/auth/kubeconfig"```
7. watch for incoming CSRs with ```watch oc get csr -A```
    - will look like this "system:serviceaccount:openshift-machine-config-operator:node-bootstrapper"
8. Approve CSRs with ```oc adm certificate approve [csr IDs]```
9. after CSRs have been approved you should have 3 master nodes and 3 worker nodes ```oc get nodes```