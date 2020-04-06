# Azure MAG UPI quickstart

1. run ```openshift-install create install-config```
    - input valid azure public information (this won't be used)
2. Modify install-config with the correct values for azure government
    - platform.azure.region
    - platform.azure.baseDomainResourceGroupName
    - basedomain
3. copy all ARM templates from $CODE_LOCATION/upi/azure
4. copy azureGovQuickstart.sh from $CODE_LOCATION/upi/azure
5. run azureGovQuickstart.sh -w {Number of desired worker nodes}
6. when installer pauses edit the following files
    1. manifests/cloud-provider-config with the correct information
      - cloud: AzurePublicCloud --> AzureUSGovernmentCloud
      - tenantId: {Azure Public tenant} --> Azure Gov Tenant
      - subscriptionId: {Azure Public Subscription} --> Azure Gov Subscription
    2. openshift/99_cloud-creds-secret with the correct information
      - azure_subscription_id: {base64 encoded Public subscription} --> base64 Azure Gov Subscription
      - azure_client_id: {base64 encoded Client ID} --> base64 Azure Gov client id
      - azure_client_secret: {base64 encoded Client Secret} --> base64 Azure Gov Client Secret
      - azure_tenant_id: {base64 encoded tenant id} --> base64 Azure Gov tenant id
7. run ```export KUBECONFIG="$PWD/auth/kubeconfig"```
8. watch for incoming CSRs with ```watch oc get csr -A```
    - will look like this "system:serviceaccount:openshift-machine-config-operator:node-bootstrapper"
9. Approve CSRs with ```oc adm certificate approve [csr IDs]```
10. after CSRs have been approved you should have 3 master nodes and desired worker nodes ```oc get nodes```
11. run ```oc edit configs.imageregistry.operator.openshift.io/cluster ```
    - change managedState: Managed --> Removed
12. add A records to the public and private dns zone for the created node LB
    - *.apps.{base domain}
13. run ```openshift-install wait-for install-complete``` for the username and password of the web console
