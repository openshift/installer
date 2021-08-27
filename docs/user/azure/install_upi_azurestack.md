# Install: User-Provided Infrastructure on Azure Stack Hub

The steps for performing a user-provided infrastructure install are outlined here. Several
[Azure Resource Manager][azuretemplates] templates are provided to assist in
completing these steps or to help model your own. You are also free to create
the required resources through other methods; the templates are just an
example.

## Prerequisites

* all prerequisites from [README](README.md)
* the following binaries installed and in $PATH:
  * [openshift-install][openshiftinstall]
    * It is recommended that the OpenShift installer CLI version is the same of the cluster being deployed. The version used in this example is 4.9 GA.
  * [az (Azure CLI)][azurecli] installed and aunthenticated
    * `az` should be [configured to connect to the Azure Stack Hub instance][configurecli]
    * Commands flags and structure may vary between `az` versions. The recommended version used in this example is 2.26.1.
  * python3
  * [jq][jqjson]
  * [yq][yqyaml] (N.B. there are multiple versions of `yq`, some with different syntaxes.)

## Create an install config

Create an install config, `install-config.yaml`. Here is a minimal example:

```yaml
apiVersion: v1
baseDomain: <example.com>
compute:
- name: worker
  platform: {}
  replicas: 0
metadata:
  name: padillon
platform:
  azure:
    armEndpoint: <azurestack-arm-endpoint>
    baseDomainResourceGroupName: <resource-group-for-example.com>
    cloudName: AzureStackCloud
    region: <azurestack-region>
pullSecret: <redacted>
sshKey: |
  <pubkey>
```

We'll be providing the compute machines ourselves, so we set compute replicas to 0.

Azure Stack is not supported by the interactive wizard, but you can use public Azure credentials to create an install config with [the usual approach](install.md#create-configuration) and then edit according to the example above.

## Credentials

Both Azure and Azure Stack credentials are stored by the installer at `~/.azure/osServicePrincipal.json`. The installer will request the required information if no credentials are found.

```console
$ openshift-install create manifests
? azure subscription id xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
? azure tenant id xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
? azure service principal client id xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
? azure service principal client secret xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
INFO Saving user credentials to "/home/user_id/.azure/osServicePrincipal.json"
```

### Extract data from install config

Some data from the install configuration file will be used on later steps. Export them as environment variables with:

```sh
export CLUSTER_NAME=$(yq -r .metadata.name install-config.yaml)
export AZURE_REGION=$(yq -r .platform.azure.region install-config.yaml)
export SSH_KEY=$(yq -r .sshKey install-config.yaml | xargs)
export BASE_DOMAIN=$(yq -r .baseDomain install-config.yaml)
export BASE_DOMAIN_RESOURCE_GROUP=$(yq -r .platform.azure.baseDomainResourceGroupName install-config.yaml)
```

## Create manifests

Create manifests to enable customizations that are not exposed via the install configuration.

```console
$ openshift-install create manifests
INFO Credentials loaded from file "/home/user_id/.azure/osServicePrincipal.json"
INFO Consuming "Install Config" from target directory
WARNING Making control-plane schedulable by setting MastersSchedulable to true for Scheduler cluster settings
```

### Remove control plane machines and machinesets

Remove the control plane machines and compute machinesets from the manifests.
We'll be providing those ourselves and don't want to involve the [machine-API operator][machine-api-operator].

```sh
rm -f openshift/99_openshift-cluster-api_master-machines-*.yaml
rm -f openshift/99_openshift-cluster-api_worker-machineset-*.yaml
```

### Make control-plane nodes unschedulable

Currently [emptying the compute pools](#empty-the-compute-pool) makes control-plane nodes schedulable.
But due to a [Kubernetes limitation][kubernetes-service-load-balancers-exclude-masters], router pods running on control-plane nodes will not be reachable by the ingress load balancer.
Update the scheduler configuration to keep router pods and other workloads off the control-plane nodes:

```sh
python3 -c '
import yaml;
path = "manifests/cluster-scheduler-02-config.yml";
data = yaml.full_load(open(path));
data["spec"]["mastersSchedulable"] = False;
open(path, "w").write(yaml.dump(data, default_flow_style=False))'
```

### Remove DNS Zones

We don't want [the ingress operator][ingress-operator] to create DNS records (we're going to do it manually) so we need to remove
the `publicZone` section from the DNS configuration in manifests.

```sh
python3 -c '
import yaml;
path = "manifests/cluster-dns-02-config.yml";
data = yaml.full_load(open(path));
del data["spec"]["publicZone"];
open(path, "w").write(yaml.dump(data, default_flow_style=False))'
```

### Resource Group Name and Infra ID

The OpenShift cluster has been assigned an identifier in the form of `<cluster_name>-<random_string>`. This identifier, called "Infra ID", will be used as
the base name of most resources that will be created in this example. Export the Infra ID as an environment variable that will be used later in this example:

```sh
export INFRA_ID=$(yq -r .status.infrastructureName manifests/cluster-infrastructure-02-config.yml)
```

Also, all resources created in this Azure deployment will exist as part of a [resource group][azure-resource-group]. The resource group name is also
based on the Infra ID, in the form of `<cluster_name>-<random_string>-rg`. Export the resource group name to an environment variable that will be used later:

```sh
export RESOURCE_GROUP=$(yq -r .status.platformStatus.azure.resourceGroupName manifests/cluster-infrastructure-02-config.yml)
```

**Optional:** it's possible to choose any other name for the Infra ID and/or the resource group, but in that case some adjustments in manifests are needed.
A Python script is provided to help with these adjustments. Export the `INFRA_ID` and the `RESOURCE_GROUP` environment variables with the desired names, copy the
[`setup-manifests.py`](../../../upi/azure/setup-manifests.py) script locally and invoke it with:

```sh
python3 setup-manifests.py $RESOURCE_GROUP $INFRA_ID
```

### Create Cluster Credentials

Azure Stack Hub can only operate in `Manual` credentials mode. So we must set the Cloud Credential Operator to Manual:

```shell
cat >> manifests/cco-configmap.yaml << EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: cloud-credential-operator-config
  namespace: openshift-cloud-credential-operator
  annotations:
    release.openshift.io/create-only: "true"
data:
  disabled: "true"
EOF
```

Then follow the [official documentation for creating manual credentials][manual-credentials].

Please follow the instructions above, but the result should be similar to this example (note the filenames are arbitrary):

```console
$ cat manifests/*credentials-secret.yaml
apiVersion: v1
kind: Secret
metadata:
    name: azure-cloud-credentials
    namespace: openshift-cloud-controller-manager
stringData:
  azure_subscription_id: <subscription-id>
  azure_client_id: <client-id>
  azure_client_secret: <secret>
  azure_tenant_id: <tenant>
  azure_resource_prefix: <$INFRA_ID>
  azure_resourcegroup: <$RESOURCE_GROUP>
  azure_region: <$REGION>
apiVersion: v1
kind: Secret
metadata:
    name: installer-cloud-credentials
    namespace: openshift-image-registry
stringData:
  azure_subscription_id: <subscription-id>
  azure_client_id: <client-id>
  azure_client_secret: <secret>
  azure_tenant_id: <tenant>
  azure_resource_prefix: <$INFRA_ID>
  azure_resourcegroup: <$RESOURCE_GROUP>
  azure_region: <$REGION>
apiVersion: v1
kind: Secret
metadata:
    name: cloud-credentials
    namespace: openshift-ingress-operator
stringData:
  azure_subscription_id: <subscription-id>
  azure_client_id: <client-id>
  azure_client_secret: <secret>
  azure_tenant_id: <tenant>
  azure_resource_prefix: <$INFRA_ID>
  azure_resourcegroup: <$RESOURCE_GROUP>
  azure_region: <$REGION>
apiVersion: v1
kind: Secret
metadata:
  name: azure-cloud-credentials
  namespace: openshift-machine-api
stringData:
  azure_subscription_id: <subscription-id>
  azure_client_id: <client-id>
  azure_client_secret: <secret>
  azure_tenant_id: <tenant>
  azure_resource_prefix: <$INFRA_ID>
  azure_resourcegroup: <$RESOURCE_GROUP>
  azure_region: <$REGION>
```

## Create ignition configs

Now we can create the bootstrap ignition configs:

```console
$ openshift-install create ignition-configs
INFO Consuming Openshift Manifests from target directory
INFO Consuming Worker Machines from target directory
INFO Consuming Common Manifests from target directory
INFO Consuming Master Machines from target directory
```

After running the command, several files will be available in the directory.

```console
$ tree
.
├── auth
│   └── kubeconfig
├── bootstrap.ign
├── master.ign
├── metadata.json
└── worker.ign
```

## Create The Resource Group

Use the command below to create the resource group in the selected Azure region:

```sh
az group create --name $RESOURCE_GROUP --location $AZURE_REGION
```

## Upload the files to a Storage Account

The deployment steps will read the Red Hat Enterprise Linux CoreOS virtual hard disk (VHD) image and the bootstrap ignition config file
from a blob. Create a storage account that will be used to store them and export its key as an environment variable.

```sh
az storage account create -g $RESOURCE_GROUP --location $AZURE_REGION --name ${CLUSTER_NAME}sa --kind Storage --sku Standard_LRS
export ACCOUNT_KEY=`az storage account keys list -g $RESOURCE_GROUP --account-name ${CLUSTER_NAME}sa --query "[0].value" -o tsv`
```

### Copy the cluster image

In order to create VMs, the RHCOS VHD must be available in the Azure Stack
environment. The VHD should be downloaded locally, decompressed, and uploaded to a
storage blob.

First, download and decompress the VHD, note that the decompressed file is 16GB:

```sh
$ export COMPRESSED_VHD_URL=$(curl -s https://raw.githubusercontent.com/openshift/installer/release-4.9/data/data/rhcos-amd64.json | jq -r '(.baseURI + .images.azurestack.path)')
$ $ curl -O -L $COMPRESSED_VHD_URL 
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100  983M  100  983M    0     0  9416k      0  0:01:46  0:01:46 --:--:-- 10.5M
$ gunzip rhcos-49.84.202107010027-0-azurestack.x86_64.vhd.gz
```

Next, create a container for the VHD:

```sh
az storage container create --name vhd --account-name ${CLUSTER_NAME}sa
```

As mentioned above, the VHD size is massive. If you have a fast upload speed,
you may be satisfied to upload the whole 16GB file:

```sh
az storage blob upload --account-name ${CLUSTER_NAME}sa --account-key $ACCOUNT_KEY -c vhd -n "rhcos.vhd" -f rhcos-49.84.202107010027-0-azurestack.x86_64.vhd
```

If your local connection speed is too slow to upload a 16 GB file, consider using a VM in your Azure Stack Hub instance or another cloud provider.

### Upload the bootstrap ignition

Create a blob storage container and upload the generated `bootstrap.ign` file:

```sh
az storage container create --name files --account-name "${CLUSTER_NAME}sa" --public-access blob --account-key "$ACCOUNT_KEY"
az storage blob upload --account-name "${CLUSTER_NAME}sa" --account-key "$ACCOUNT_KEY" -c "files" -f "bootstrap.ign" -n "bootstrap.ign"
```

## Create the DNS zones

A few DNS records are required for clusters that use user-provisioned infrastructure. Feel free to choose the DNS strategy that fits you best.

This example adds records to an Azure Stack public zone. For external (internet) visibility, your authoritative DNS zone, such as [Azure's own DNS solution][azure-dns],
can delegate to the DNS nameservers for you Azure Stack environment.
Note that the public zone doesn't necessarily need to exist in the same resource group of the
cluster deployment itself and may even already exist in your organization for the desired base domain. If that's the case, you can skip the public DNS
zone creation step, but make sure the install config generated earlier [reflects that scenario](customization.md#cluster-scoped-properties).

Create the new *public* DNS zone in the resource group exported in the `BASE_DOMAIN_RESOURCE_GROUP` environment variable, or just skip this step if you're going
to use one that already exists in your organization:

```sh
az network dns zone create -g "$BASE_DOMAIN_RESOURCE_GROUP" -n "${CLUSTER_NAME}.${BASE_DOMAIN}"
```

## Deployment

The key parts of this UPI deployment are the [Azure Resource Manager][azuretemplates] templates, which are responsible
for deploying most resources. They're provided as a few json files following the "NN_name.json" pattern. In the
next steps we're going to deploy each one of them in order, using [az (Azure CLI)][azurecli] and providing the expected parameters.

## Deploy the Virtual Network

In this example we're going to create a Virtual Network and subnets specifically for the OpenShift cluster. You can skip this step
if the cluster is going to live in a VNet already existing in your organization, or you can edit the `01_vnet.json` file to your
own needs (e.g. change the subnets address prefixes in CIDR format).

Copy the [`01_vnet.json`](../../../upi/azurestack/01_vnet.json) ARM template locally.

Create the deployment using the `az` client:

```sh
az deployment group create -g $RESOURCE_GROUP \
  --template-file "01_vnet.json" \
  --parameters baseName="$INFRA_ID"
```

## Deploy the image

Copy the [`02_storage.json`](../../../upi/azurestack/02_storage.json) ARM template locally.

Create the deployment using the `az` client:

```sh
export VHD_BLOB_URL=`az storage blob url --account-name ${CLUSTER_NAME}sa --account-key $ACCOUNT_KEY -c vhd -n "rhcos.vhd" -o tsv`

az deployment group create -g $RESOURCE_GROUP \
  --template-file "02_storage.json" \
  --parameters vhdBlobURL="$VHD_BLOB_URL" \
  --parameters baseName="$INFRA_ID"
```

## Deploy the load balancers

Copy the [`03_infra.json`](../../../upi/azurestack/03_infra.json) ARM template locally.

Deploy the load balancer and public IP addresses using the `az` client:

```sh
az deployment group create -g $RESOURCE_GROUP \
  --template-file "03_infra.json" \
  --parameters baseName="$INFRA_ID"
```

Create an `api` and `api-int` DNS record in the *public* zone for the API public load balancer. Note that the `BASE_DOMAIN_RESOURCE_GROUP` must point to the resource group where the public DNS zone exists.

```sh
export PUBLIC_IP=$(az network public-ip list -g "$RESOURCE_GROUP" --query "[?name=='${INFRA_ID}-master-pip'] | [0].ipAddress" -o tsv)
az network dns record-set a add-record -g "$RESOURCE_GROUP" -z "${CLUSTER_NAME}.${BASE_DOMAIN}" -n api -a "$PUBLIC_IP" --ttl 60
az network dns record-set a add-record -g "$RESOURCE_GROUP" -z "${CLUSTER_NAME}.${BASE_DOMAIN}" -n api-int -a "$PUBLIC_IP" --ttl 60
```

## Launch the temporary cluster bootstrap

Copy the [`04_bootstrap.json`](../../../upi/azurestack/04_bootstrap.json) ARM template locally.

Create the deployment using the `az` client:

```sh
export BOOTSTRAP_URL=$(az storage blob url --account-name "${INFRA_ID}sa" --account-key "$ACCOUNT_KEY" -c "files" -n "bootstrap.ign" -o tsv)
export BOOTSTRAP_IGNITION=$(jq -rcnM --arg v "3.2.0" --arg url "$BOOTSTRAP_URL" '{ignition:{version:$v,config:{replace:{source:$url}}}}' | base64 | tr -d '\n')

az deployment group create --verbose -g "$RESOURCE_GROUP" \
  --template-file "04_bootstrap.json" \
  --parameters bootstrapIgnition="$BOOTSTRAP_IGNITION" \
  --parameters sshKeyData="$SSH_KEY" \
  --parameters baseName="$INFRA_ID" \
  --parameters diagnosticsStorageAccountName="${INFRA_ID}sa"
```

## Launch the permanent control plane

Copy the [`05_masters.json`](../../../upi/azurestack/05_masters.json) ARM template locally.

Create the deployment using the `az` client:

```sh
export MASTER_IGNITION=$(cat master.ign | base64 | tr -d '\n')

az deployment group create -g "$RESOURCE_GROUP" \
  --template-file "05_masters.json" \
  --parameters masterIgnition="$MASTER_IGNITION" \
  --parameters sshKeyData="$SSH_KEY" \
  --parameters baseName="$INFRA_ID" \
  --parameters masterVMSize="Standard_DS4_v2" \
  --parameters diskSizeGB="1023" \
  --parameters diagnosticsStorageAccountName="${INFRA_ID}sa"
```

## Wait for the bootstrap complete

Wait until cluster bootstrapping has completed:

```console
$ openshift-install wait-for bootstrap-complete --log-level debug
DEBUG OpenShift Installer v4.n
DEBUG Built from commit 6b629f0c847887f22c7a95586e49b0e2434161ca
INFO Waiting up to 30m0s for the Kubernetes API at https://api.cluster.basedomain.com:6443...
DEBUG Still waiting for the Kubernetes API: the server could not find the requested resource
DEBUG Still waiting for the Kubernetes API: the server could not find the requested resource
DEBUG Still waiting for the Kubernetes API: Get https://api.cluster.basedomain.com:6443/version?timeout=32s: dial tcp: connect: connection refused
INFO API v1.14.n up
INFO Waiting up to 30m0s for bootstrapping to complete...
DEBUG Bootstrap status: complete
INFO It is now safe to remove the bootstrap resources
```

Once the bootstrapping process is complete you can deallocate and delete bootstrap resources:

```sh
az network nsg rule delete -g $RESOURCE_GROUP --nsg-name ${INFRA_ID}-nsg --name bootstrap_ssh_in
az vm stop -g $RESOURCE_GROUP --name ${INFRA_ID}-bootstrap
az vm deallocate -g $RESOURCE_GROUP --name ${INFRA_ID}-bootstrap
az vm delete -g $RESOURCE_GROUP --name ${INFRA_ID}-bootstrap --yes
az disk delete -g $RESOURCE_GROUP --name ${INFRA_ID}-bootstrap_OSDisk --no-wait --yes
az network nic delete -g $RESOURCE_GROUP --name ${INFRA_ID}-bootstrap-nic --no-wait
az storage blob delete --account-key $ACCOUNT_KEY --account-name ${CLUSTER_NAME}sa --container-name files --name bootstrap.ign
az network public-ip delete -g $RESOURCE_GROUP --name ${INFRA_ID}-bootstrap-ssh-pip
```

## Access the OpenShift API

You can now use the `oc` or `kubectl` commands to talk to the OpenShift API. The admin credentials are in `auth/kubeconfig`. For example:

```sh
export KUBECONFIG="$PWD/auth/kubeconfig"
oc get nodes
oc get clusteroperator
```

Note that only the API will be up at this point. The OpenShift web console will run on the compute nodes.

## Launch compute nodes

You may create compute nodes by launching individual instances discretely or by automated processes outside the cluster (e.g. Auto Scaling Groups).
You can also take advantage of the built in cluster scaling mechanisms and the machine API in OpenShift.

In this example, we'll manually launch three instances via the provided ARM template. Additional instances can be launched by editing the `06_workers.json` file.

Copy the [`06_workers.json`](../../../upi/azure/06_workers.json) ARM template locally.

Create the deployment using the `az` client:

```sh
export WORKER_IGNITION=`cat worker.ign | base64 | tr -d '\n'`

az deployment group create -g $RESOURCE_GROUP \
  --template-file "06_workers.json" \
  --parameters workerIgnition="$WORKER_IGNITION" \
  --parameters sshKeyData="$SSH_KEY" \
  --parameters baseName="$INFRA_ID" \
  --parameters diagnosticsStorageAccountName="${INFRA_ID}sa"
```

### Approve the worker CSRs

Even after they've booted up, the workers will not show up in `oc get nodes`.

Instead, they will create certificate signing requests (CSRs) which need to be approved. Eventually, you should see `Pending` entries looking like the ones below.
You can use `watch oc get csr -A` to watch until the pending CSR's are available.

```console
$ oc get csr -A
NAME        AGE    REQUESTOR                                                                   CONDITION
csr-8bppf   2m8s   system:serviceaccount:openshift-machine-config-operator:node-bootstrapper   Pending
csr-dj2w4   112s   system:serviceaccount:openshift-machine-config-operator:node-bootstrapper   Pending
csr-ph8s8   11s    system:serviceaccount:openshift-machine-config-operator:node-bootstrapper   Pending
csr-q7f6q   19m    system:node:master01                                                        Approved,Issued
csr-5ztvt   19m    system:node:master02                                                        Approved,Issued
csr-576l2   19m    system:node:master03                                                        Approved,Issued
csr-htmtm   19m    system:serviceaccount:openshift-machine-config-operator:node-bootstrapper   Approved,Issued
csr-wpvxq   19m    system:serviceaccount:openshift-machine-config-operator:node-bootstrapper   Approved,Issued
csr-xpp49   19m    system:serviceaccount:openshift-machine-config-operator:node-bootstrapper   Approved,Issued
```

You should inspect each pending CSR with the `oc describe csr <name>` command and verify that it comes from a node you recognise. If it does, they can be approved:

```console
$ oc adm certificate approve csr-8bppf csr-dj2w4 csr-ph8s8
certificatesigningrequest.certificates.k8s.io/csr-8bppf approved
certificatesigningrequest.certificates.k8s.io/csr-dj2w4 approved
certificatesigningrequest.certificates.k8s.io/csr-ph8s8 approved
```

Approved nodes should now show up in `oc get nodes`, but they will be in the `NotReady` state. They will create a second CSR which must also be reviewed and approved.
Repeat the process of inspecting the pending CSR's and approving them.

Once all CSR's are approved, the node should switch to `Ready` and pods will be scheduled on it.

```console
$ oc get nodes
NAME       STATUS   ROLES    AGE     VERSION
master01   Ready    master   23m     v1.14.6+cebabbf7a
master02   Ready    master   23m     v1.14.6+cebabbf7a
master03   Ready    master   23m     v1.14.6+cebabbf7a
node01     Ready    worker   2m30s   v1.14.6+cebabbf7a
node02     Ready    worker   2m35s   v1.14.6+cebabbf7a
node03     Ready    worker   2m34s   v1.14.6+cebabbf7a
```

### Add the Ingress DNS Records

Create DNS records in the public zone pointing at the ingress load balancer. Use A, CNAME, etc. records, as you see fit.
You can create either a wildcard `*.apps.{baseDomain}.` or [specific records](#specific-route-records) for every route (more on the specific records below).

First, wait for the ingress default router to create a load balancer and populate the `EXTERNAL-IP` column:

```console
$ oc -n openshift-ingress get service router-default
NAME             TYPE           CLUSTER-IP      EXTERNAL-IP     PORT(S)                      AGE
router-default   LoadBalancer   172.30.20.10   35.130.120.110   80:32288/TCP,443:31215/TCP   20
```

Add a `*.apps` record to the *public* DNS zone:

```sh
export PUBLIC_IP_ROUTER=`oc -n openshift-ingress get service router-default --no-headers | awk '{print $4}'`
az network dns record-set a add-record -g $BASE_DOMAIN_RESOURCE_GROUP -z ${CLUSTER_NAME}.${BASE_DOMAIN} -n *.apps -a $PUBLIC_IP_ROUTER --ttl 300
```

Or, in case of adding this cluster to an already existing public zone, use instead:

```sh
export PUBLIC_IP_ROUTER=`oc -n openshift-ingress get service router-default --no-headers | awk '{print $4}'`
az network dns record-set a add-record -g $BASE_DOMAIN_RESOURCE_GROUP -z ${BASE_DOMAIN} -n *.apps.${CLUSTER_NAME} -a $PUBLIC_IP_ROUTER --ttl 300
```

#### Specific route records

If you prefer to add explicit domains instead of using a wildcard, you can create entries for each of the cluster's current routes. Use the command below to check what they are:

```console
$ oc get --all-namespaces -o jsonpath='{range .items[*]}{range .status.ingress[*]}{.host}{"\n"}{end}{end}' routes
oauth-openshift.apps.cluster.basedomain.com
console-openshift-console.apps.cluster.basedomain.com
downloads-openshift-console.apps.cluster.basedomain.com
alertmanager-main-openshift-monitoring.apps.cluster.basedomain.com
grafana-openshift-monitoring.apps.cluster.basedomain.com
prometheus-k8s-openshift-monitoring.apps.cluster.basedomain.com
```

## Wait for the installation complete

Wait until cluster is ready:

```console
$ openshift-install wait-for install-complete --log-level debug
DEBUG Built from commit 6b629f0c847887f22c7a95586e49b0e2434161ca
INFO Waiting up to 30m0s for the cluster at https://api.cluster.basedomain.com:6443 to initialize...
DEBUG Still waiting for the cluster to initialize: Working towards 4.2.12: 99% complete, waiting on authentication, console, monitoring
DEBUG Still waiting for the cluster to initialize: Working towards 4.2.12: 100% complete
DEBUG Cluster is initialized
INFO Waiting up to 10m0s for the openshift-console route to be created...
DEBUG Route found in openshift-console namespace: console
DEBUG Route found in openshift-console namespace: downloads
DEBUG OpenShift console route is created
INFO Install complete!
INFO To access the cluster as the system:admin user when using 'oc', run 'export KUBECONFIG=${PWD}/auth/kubeconfig'
INFO Access the OpenShift web-console here: https://console-openshift-console.apps.cluster.basedomain.com
INFO Login to the console with user: kubeadmin, password: REDACTED
```

[azuretemplates]: https://docs.microsoft.com/en-us/azure/azure-resource-manager/template-deployment-overview
[openshiftinstall]: https://github.com/openshift/installer
[azurecli]: https://docs.microsoft.com/en-us/cli/azure/
[configurecli]: https://docs.microsoft.com/en-us/azure-stack/user/azure-stack-version-profiles-azurecli2?view=azs-2102&tabs=ad-lin#connect-with-azure-cli
[jqjson]: https://stedolan.github.io/jq/
[yqyaml]: https://kislyuk.github.io/yq/
[ingress-operator]: https://github.com/openshift/cluster-ingress-operator
[machine-api-operator]: https://github.com/openshift/machine-api-operator
[azure-identity]: https://docs.microsoft.com/en-us/azure/architecture/framework/security/identity
[azure-resource-group]: https://docs.microsoft.com/en-us/azure/azure-resource-manager/management/overview#resource-groups
[azure-dns]: https://docs.microsoft.com/en-us/azure/dns/dns-overview
[kubernetes-service-load-balancers-exclude-masters]: https://github.com/kubernetes/kubernetes/issues/65618
[manual-credentials]: https://docs.openshift.com/container-platform/4.8/installing/installing_azure/manually-creating-iam-azure.html
[azure-vhd-utils]: https://github.com/microsoft/azure-vhd-utils
