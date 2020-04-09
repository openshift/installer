# Azure MAG UPI Instructions

This brief guide will demonstrate how to use the UPI installer to install OCP 4.X in Azure Government

## Install Procedures

### Prerequisites

Applications
- Azcli
- jq
- yq
- oc client
- kubectl
- openshift-install

Azure Access
- Contributor access to an existing azure government subscription
- Contributor access to an existing azure comercial subscription

Azure Rquirements
- Public DNS zone available in azure comercial
- Public DNS zone available in azure government

### Create Install Config

Run openshift install to create initial cluster config
```shell
openshift-install create install-config
```
Input valid azure public information, this won't be used during the install

Modify the resulting template
```yaml
apiVersion: v1
baseDomain: <valid azure gov dns zone name>
compute:
- hyperthreading: Enabled
  name: worker
  platform: {}
  replicas: 3
controlPlane:
  hyperthreading: Enabled
  name: master
  platform: {}
  replicas: 3
metadata:
  creationTimestamp: null
  name: upiazuremag
networking:
  clusterNetwork:
  - cidr: 10.128.0.0/14
    hostPrefix: 23
  machineNetwork:
  - cidr: 10.0.0.0/16
  networkType: OpenShiftSDN
  serviceNetwork:
  - 172.30.0.0/16
platform:
  azure:
    baseDomainResourceGroupName: <valid azure government rg with dns zone>
    region: <valid azure government region>
publish: External
pullSecret: '<Your Pull Secret>'
sshKey: |
  <Your SSH key>
```

### Copy required files

Copy the following files from $CODE_LOCATION/upi/azure

- 01_vnet.json
- 02_storage.json
- 03_infra.json
- 04_bootstrap.json
- 05_masters.json
- 06_workers.json
- azureGovQuickstart.sh

A quick way to do this is (assuming default code location)
```shell
export CODE_LOCATION=~/go/src/github.com/openshift/installer
cp $CODE_LOCATION/upi/azure/0*.json ./
cp $CODE_LOCATION/upi/azure/azureGovQuickstart.sh ./
```

### Run Azure Gov Quickstart

Be in the directory with the copied files and run
```shell
chmod +x azureGovQuickstart.sh
export WORKER_NODE_COUNT=<Number of workers you want, default of 3>
./azureGovQuickstart.sh -w $WORKER_NODE_COUNT
```

#### Log into azure portal
The installer will prompt you to sign into the azure gov portal with a web browser pop up.

#### Modify the manifests
The installer will pause and ask you to modify the manifests. the following edits need to be made

Quick Summary

1. manifests/cloud-provider-config with the correct information
    - cloud: AzurePublicCloud --> AzureUSGovernmentCloud
    - tenantId: {Azure Public tenant} --> Azure Gov Tenant
    - subscriptionId: {Azure Public Subscription} --> Azure Gov Subscription
2. openshift/99_cloud-creds-secret with the correct information
    - azure_subscription_id: {base64 encoded Public subscription} --> base64 Azure Gov Subscription
    - azure_client_id: {base64 encoded Client ID} --> base64 Azure Gov client id
    - azure_client_secret: {base64 encoded Client Secret} --> base64 Azure Gov Client Secret
    - azure_tenant_id: {base64 encoded tenant id} --> base64 Azure Gov tenant id

Detailed Files

```yaml
# manifests/cloud-provider-config.yaml
apiVersion: v1
data:
  config: "{\n\t\"cloud\": \"AzureUSGovernmentCloud\",\n\t\"tenantId\": \"{US_GOV_TENANTID}\",\n\t\"aadClientId\":
    \"\",\n\t\"aadClientSecret\": \"\",\n\t\"aadClientCertPath\": \"\",\n\t\"aadClientCertPassword\":
    \"\",\n\t\"useManagedIdentityExtension\": true,\n\t\"userAssignedIdentityID\":
    \"\",\n\t\"subscriptionId\": \"{US_GOV_SUBSCRIPTION_ID}\",\n\t\"resourceGroup\":
    \"upiazuremag-4mt2r-rg\",\n\t\"location\": \"eastus\",\n\t\"vnetName\": \"upiazuremag-4mt2r-vnet\",\n\t\"vnetResourceGroup\":
    \"upiazuremag-4mt2r-rg\",\n\t\"subnetName\": \"upiazuremag-4mt2r-worker-subnet\",\n\t\"securityGroupName\":
    \"upiazuremag-4mt2r-node-nsg\",\n\t\"routeTableName\": \"upiazuremag-4mt2r-node-routetable\",\n\t\"primaryAvailabilitySetName\":
    \"\",\n\t\"vmType\": \"\",\n\t\"primaryScaleSetName\": \"\",\n\t\"cloudProviderBackoff\":
    true,\n\t\"cloudProviderBackoffRetries\": 0,\n\t\"cloudProviderBackoffExponent\":
    0,\n\t\"cloudProviderBackoffDuration\": 6,\n\t\"cloudProviderBackoffJitter\":
    0,\n\t\"cloudProviderRateLimit\": true,\n\t\"cloudProviderRateLimitQPS\": 6,\n\t\"cloudProviderRateLimitBucket\":
    10,\n\t\"cloudProviderRateLimitQPSWrite\": 6,\n\t\"cloudProviderRateLimitBucketWrite\":
    10,\n\t\"useInstanceMetadata\": true,\n\t\"loadBalancerSku\": \"standard\",\n\t\"excludeMasterFromStandardLB\":
    null,\n\t\"disableOutboundSNAT\": null,\n\t\"maximumLoadBalancerRuleCount\": 0\n}\n"
kind: ConfigMap
metadata:
  creationTimestamp: null
  name: cloud-provider-config
  namespace: openshift-config
```
```yaml
# openshift/99_cloud-creds-secret.yaml
kind: Secret
apiVersion: v1
metadata:
  namespace: kube-system
  name: azure-credentials
data:
  azure_subscription_id: {US_GOV_SUBSCRIPTION_ID | base64 -w0}
  azure_client_id: {US_GOV_CLIENT_ID | base64 -w0}
  azure_client_secret: {US_GOV_CLIENT_SECRET | base64 -w0}
  azure_tenant_id: {US_GOV_TENANT_ID | base64 -w0}
  azure_resource_prefix: <No Change>
  azure_resourcegroup: <No Change>
  azure_region: <No Change>
```
#### Finish Cluster Configuration

Wait 10 minutes for the nodes to requst cluster membership and then run the following.

After the script completes configure kubeconfig, join the work nodes, edit the registry operator, and launch the web console.

```bash
# Configure kubeconfig to authenticate against new openshift
export KUBECONFIG="$PWD/auth/kubeconfig"

# Get CSR ids of pending requests
oc get csr -A

# Approve CSR ids
oc get csr -ojson | jq -r '.items[] | select(.status == {} ) | .metadata.name' | xargs oc adm certificate approve

# Nodes should populate and become ready in a couple of minutes
watch oc get nodes

# Edit image registry operator config
# The storage operator doesn't yet work with Azure Mag so the internal registry has to be disabled
oc edit configs.imageregistry.operator.openshift.io/cluster
# change managedState: Managed --> Removed
```

Add DNS Records in public and private dns zone of the ip address assigned to the new load balancer. 

*.apps --> ip address of new LB

```bash
# Complete the cluster install and get temporary admin password for web console
openshift-install wait-for install-complete
```

#### Add Azure Disk Storage Class 

Add the following storage class to your deployment, for more information see official [docs](https://kubernetes.io/docs/concepts/storage/storage-classes/#azure-disk-storage-class)

```yaml
# azure-disk.yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  annotations:
    description: azure disk
    storageclass.kubernetes.io/is-default-class: "true"
  name: azuredisk
parameters:
  kind: managed
  location: <your_region>
  skuName: <your_sku>
provisioner: kubernetes.io/azure-disk
reclaimPolicy: Delete
volumeBindingMode: Immediate
```
