# Azure MAG UPI Instructions

This brief guide will demonstrate how to use the UPI installer to install OCP 4.X in Azure Government

## Notice of supportability

CloudFit Software, partnering with Microsoft and Red Hat are actively working towards support exceptions that will cover OCP 4.3 implementations running on Azure Government implementations.  This information will be published and updated periodically.
Until there is a support excpetion Red Hat recommends that this only be used in proof of concept (POC) environments, and does not recommend this for production use cases as it will not be supported.


## Install Procedures

### Prerequisites

Applications in your $PATH
- Azcli
- jq
- yq
- oc client
- kubectl
- openshift-install

Azure Access
- Contributor access to an existing azure government subscription
- Ability to create the following resource types
  - Disk
  - DNS Zone
  - Image
  - Load Balancer
  - Managed Identity
  - Network Interface
    - With a public ip binding (this will be resolved in a later release)
  - Network Security Group
  - Private DNS Zone
  - Public IP address
  - Storage Account
  - Virtual Machine
  - Virtual Network

Azure Rquirements
- Public DNS zone available in azure government

### Create Install Config

Here is a sample install config, create a install-config.yaml file in an empty folder with this template to get started.

```yaml
#install-config.yaml
apiVersion: v1
baseDomain: <YOUR_PUBLIC_DNS_ZONE>
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
  name: <YOUR_CLUSTER_NAME>
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
    baseDomainResourceGroupName: <YOUR_DNS_ZONE_RG>
    region: <YOUR_AZURE_GOVERNMENT_REGION>
publish: External
pullSecret: '<YOUR_PULL_SECRET>'
sshKey: |
 <YOUR_SSH_KEY>
```

Create or modify this file to ensure the right azure environment secrets get passed, ~/.azure/osServicePrincipal.json. All of these values are in plain text, they will be translated to base64 encoded secrets during the install.

```json
{
  "subscriptionId":"<YOUR_AZURE_GOV_SUBSCRIPTION_ID",
  "clientId":"<YOUR_AZURE_GOV_CLIENT_ID>",
  "clientSecret":"<YOUR_AZURE_GOV_CLIENT_SECRET>",
  "tenantId":"<YOUR_AZURE_GOV_TENANT_ID>"
}
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

and make the azureGovQuickstart.sh file executable.

A quick way to do this is (assuming default code location)
```shell
export CODE_LOCATION=~/go/src/github.com/openshift/installer
cp $CODE_LOCATION/upi/azure/0*.json ./
cp $CODE_LOCATION/upi/azure/azureGovQuickstart.sh ./
chmod +x azureGovQuickstart.sh
```

### Run Azure Gov Quickstart

Be in the directory with the copied files and run
```shell
export WORKER_NODE_COUNT=<Number of workers you want, default of 3>
./azureGovQuickstart.sh -w $WORKER_NODE_COUNT
```

#### Log into azure portal
The installer will prompt you to sign into the azure gov portal with a web browser pop up.

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

## Supported Azure Environment Overrides

|Azure Environment 	| AZURE_ENVIRONMENT Override Text |
|------------------	|-------------------------------	|
| USGovernmentCloud	| AZUREUSGOVERNMENTCLOUD 	        |
| GermanCloud      	| AZUREGERMANCLOUD       	        |
| ChinaCloud       	| AZURECHINACLOUD        	        |