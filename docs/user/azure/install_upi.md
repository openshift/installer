# Install: User Provided Infrastructure (UPI)

The steps for performing a user-provided infrastructure install are outlined here. Several
[Azure Resource Manager][azuretemplates] templates are provided to assist in
completing these steps or to help model your own. You are also free to create
the required resources through other methods; the templates are just an
example.

## Prerequisites

* all prerequisites from [README](README.md)
* the following binaries installed and in $PATH:
  * [openshift-install][openshiftinstall]
    * It is recommended that the OpenShift installer CLI version is the same of the cluster being deployed. The version used in this example is 4.3.0 GA.
  * [az (Azure CLI)][azurecli] installed and aunthenticated
    * Commands flags and structure may vary between `az` versions. The recommended version used in this example is 2.0.80.
  * python3
  * [jq][jqjson]
  * [yq][yqyaml]

## Create an install config

Create an install configuration as for [the usual approach](install.md#create-configuration).

```console
$ openshift-install create install-config
? SSH Public Key /home/user_id/.ssh/id_rsa.pub
? Platform azure
? azure subscription id xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
? azure tenant id xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
? azure service principal client id xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
? azure service principal client secret xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
INFO Saving user credentials to "/home/user_id/.azure/osServicePrincipal.json"
? Region centralus
? Base Domain example.com
? Cluster Name test
? Pull Secret [? for help]
```

Note that we're going to have a new Virtual Network and subnetworks created specifically for this deployment, but it is also possible to use a networking
infrastructure already existing in your organization. Please refer to the [customization instructions](customization.md) for more details about setting
up an install config for that scenario.

### Extract data from install config

Some data from the install configuration file will be used on later steps. Export them as environment variables with:

```sh
export CLUSTER_NAME=`yq -r .metadata.name install-config.yaml`
export AZURE_REGION=`yq -r .platform.azure.region install-config.yaml`
export SSH_KEY=`yq -r .sshKey install-config.yaml | xargs`
export BASE_DOMAIN=`yq -r .baseDomain install-config.yaml`
export BASE_DOMAIN_RESOURCE_GROUP=`yq -r .platform.azure.baseDomainResourceGroupName install-config.yaml`
```

### Empty the compute pool

We'll be providing the compute machines ourselves, so edit the resulting `install-config.yaml` to set `replicas` to 0 for the `compute` pool:

```sh
python3 -c '
import yaml;
path = "install-config.yaml";
data = yaml.full_load(open(path));
data["compute"][0]["replicas"] = 0;
open(path, "w").write(yaml.dump(data, default_flow_style=False))'
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
the `privateZone` and `publicZone` sections from the DNS configuration in manifests.

```sh
python3 -c '
import yaml;
path = "manifests/cluster-dns-02-config.yml";
data = yaml.full_load(open(path));
del data["spec"]["publicZone"];
del data["spec"]["privateZone"];
open(path, "w").write(yaml.dump(data, default_flow_style=False))'
```

### Resource Group Name and Infra ID

The OpenShift cluster has been assigned an identifier in the form of `<cluster_name>-<random_string>`. This identifier, called "Infra ID", will be used as
the base name of most resources that will be created in this example. Export the Infra ID as an environment variable that will be used later in this example:

```sh
export INFRA_ID=`yq -r '.status.infrastructureName' manifests/cluster-infrastructure-02-config.yml`
```

Also, all resources created in this Azure deployment will exist as part of a [resource group][azure-resource-group]. The resource group name is also
based on the Infra ID, in the form of `<cluster_name>-<random_string>-rg`. Export the resource group name to an environment variable that will be user later:

```sh
export RESOURCE_GROUP=`yq -r '.status.platformStatus.azure.resourceGroupName' manifests/cluster-infrastructure-02-config.yml`
```

**Optional:** it's possible to choose any other name for the Infra ID and/or the resource group, but in that case some adjustments in manifests are needed.
A Python script is provided to help with these adjustments. Export the `INFRA_ID` and the `RESOURCE_GROUP` environment variables with the desired names, copy the
[`setup-manifests.py`](../../../upi/azure/setup-manifests.py) script locally and invoke it with:

```sh
python3 setup-manifests.py $RESOURCE_GROUP $INFRA_ID
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

## Create The Resource Group and identity

Use the command below to create the resource group in the selected Azure region:

```sh
az group create --name $RESOURCE_GROUP --location $AZURE_REGION
```

Also, create an identity which will be used to grant the required access to cluster operators:

```sh
az identity create -g $RESOURCE_GROUP -n ${INFRA_ID}-identity
```

## Upload the files to a Storage Account

The deployment steps will read the Red Hat Enterprise Linux CoreOS virtual hard disk (VHD) image and the bootstrap ignition config file
from a blob. Create a storage account that will be used to store them and export its key as an environment variable.

```sh
az storage account create -g $RESOURCE_GROUP --location $AZURE_REGION --name ${CLUSTER_NAME}sa --kind Storage --sku Standard_LRS
export ACCOUNT_KEY=`az storage account keys list -g $RESOURCE_GROUP --account-name ${CLUSTER_NAME}sa --query "[0].value" -o tsv`
```

### Copy the cluster image

Given the size of the RHCOS VHD, it's not possible to run the deployments with this file stored locally on your machine.
We must copy and store it in a storage container instead. To do so, first create a blob storage container and then copy the VHD.

Choose the RHCOS version you'd like to use and export the URL of its VHD to an environment variable. For example, to use the latest release available for the 4.3 version, use:

```sh
export VHD_URL=`curl -s https://raw.githubusercontent.com/openshift/installer/release-4.3/data/data/rhcos.json | jq -r .azure.url`
```

If you'd just like to use the latest _development_ version available (master branch), use:

```sh
export VHD_URL=`curl -s https://raw.githubusercontent.com/openshift/installer/master/data/data/rhcos.json | jq -r .azure.url`
```

Copy the chosen VHD to a blob:

```sh
az storage container create --name vhd --account-name ${CLUSTER_NAME}sa
az storage blob copy start --account-name ${CLUSTER_NAME}sa --account-key $ACCOUNT_KEY --destination-blob "rhcos.vhd" --destination-container vhd --source-uri "$VHD_URL"
```

To track the progress, you can use:

```sh
status="unknown"
while [ "$status" != "success" ]
do
  status=`az storage blob show --container-name vhd --name "rhcos.vhd" --account-name ${CLUSTER_NAME}sa --account-key $ACCOUNT_KEY -o tsv --query properties.copy.status`
  echo $status
done
```

### Upload the bootstrap ignition

Create a blob storage container and upload the generated `bootstrap.ign` file:

```sh
az storage container create --name files --account-name ${CLUSTER_NAME}sa --public-access blob
az storage blob upload --account-name ${CLUSTER_NAME}sa --account-key $ACCOUNT_KEY -c "files" -f "bootstrap.ign" -n "bootstrap.ign"
```

## Create the DNS zones

A few DNS records are required for clusters that use user-provisioned infrastructure. Feel free to choose the DNS strategy that fits you best.

In this example we're going to use [Azure's own DNS solution][azure-dns], so we're going to create a new public DNS zone for external (internet) visibility, and
a private DNS zone for internal cluster resolution. Note that the public zone don't necessarily need to exist in the same resource group of the
cluster deployment itself and may even already exist in your organization for the desired base domain. If that's the case, you can skip the public DNS
zone creation step, but make sure the install config generated earlier [reflects that scenario](customization.md#cluster-scoped-properties).

Create the new *public* DNS zone in the resource group exported in the `BASE_DOMAIN_RESOURCE_GROUP` environment variable, or just skip this step if you're going
to use one that already exists in your organization:

```sh
az network dns zone create -g $BASE_DOMAIN_RESOURCE_GROUP -n ${CLUSTER_NAME}.${BASE_DOMAIN}
```

Create the *private* zone in the same resource group of the rest of this deployment:

```sh
az network private-dns zone create -g $RESOURCE_GROUP -n ${CLUSTER_NAME}.${BASE_DOMAIN}
```

## Grant access to the identity

Grant the *Contributor* role to the Azure identity so that the Ingress Operator can create a public IP and its load balancer. You can do that with:

```sh
export PRINCIPAL_ID=`az identity show -g $RESOURCE_GROUP -n ${INFRA_ID}-identity --query principalId --out tsv`
export RESOURCE_GROUP_ID=`az group show -g $RESOURCE_GROUP --query id --out tsv`
az role assignment create --assignee "$PRINCIPAL_ID" --role 'Contributor' --scope "$RESOURCE_GROUP_ID"
```

## Deployment

The key part of this UPI deployment are the [Azure Resource Manager][azuretemplates] templates, which are responsible
for deploying most resources. They're provided as a few json files named following the "NN_name.json" pattern. In the
next steps we're going to deploy each one of them in order, using [az (Azure CLI)][azurecli] and providing the expected parameters.

## Deploy the Virtual Network

In this example we're going to create a Virtual Network and subnets specifically for the OpenShift cluster. You can skip this step
if the cluster is going to live in a VNet already existing in your organization, or you can edit the `01_vnet.json` file to your
own needs (e.g. change the subnets address prefixes in CIDR format).

Copy the [`01_vnet.json`](../../../upi/azure/01_vnet.json) ARM template locally.

Create the deployment using the `az` client:

```sh
az group deployment create -g $RESOURCE_GROUP \
  --template-file "01_vnet.json" \
  --parameters baseName="$INFRA_ID"
```

Link the VNet just created to the private DNS zone:

```sh
az network private-dns link vnet create -g $RESOURCE_GROUP -z ${CLUSTER_NAME}.${BASE_DOMAIN} -n ${INFRA_ID}-network-link -v "${INFRA_ID}-vnet" -e false
```

## Deploy the image

Copy the [`02_storage.json`](../../../upi/azure/02_storage.json) ARM template locally.

Create the deployment using the `az` client:

```sh
export VHD_BLOB_URL=`az storage blob url --account-name ${CLUSTER_NAME}sa --account-key $ACCOUNT_KEY -c vhd -n "rhcos.vhd" -o tsv`

az group deployment create -g $RESOURCE_GROUP \
  --template-file "02_storage.json" \
  --parameters vhdBlobURL="$VHD_BLOB_URL" \
  --parameters baseName="$INFRA_ID"
```

## Deploy the load balancers

Copy the [`03_infra.json`](../../../upi/azure/03_infra.json) ARM template locally.

Deploy the load balancers and public IP addresses using the `az` client:

```sh
az group deployment create -g $RESOURCE_GROUP \
  --template-file "03_infra.json" \
  --parameters privateDNSZoneName="${CLUSTER_NAME}.${BASE_DOMAIN}" \
  --parameters baseName="$INFRA_ID"
```

Create an `api` DNS record in the *public* zone for the API public load balancer. Note that the `BASE_DOMAIN_RESOURCE_GROUP` must point to the resource group where the public DNS zone exists.

```sh
export PUBLIC_IP=`az network public-ip list -g $RESOURCE_GROUP --query "[?name=='${INFRA_ID}-master-pip'] | [0].ipAddress" -o tsv`
az network dns record-set a add-record -g $BASE_DOMAIN_RESOURCE_GROUP -z ${CLUSTER_NAME}.${BASE_DOMAIN} -n api -a $PUBLIC_IP --ttl 60
```

Or, in case of adding this cluster to an already existing public zone, use instead:

```sh
export PUBLIC_IP=`az network public-ip list -g $RESOURCE_GROUP --query "[?name=='${INFRA_ID}-master-pip'] | [0].ipAddress" -o tsv`
az network dns record-set a add-record -g $BASE_DOMAIN_RESOURCE_GROUP -z ${BASE_DOMAIN} -n api.${CLUSTER_NAME} -a $PUBLIC_IP --ttl 60
```

## Launch the temporary cluster bootstrap

Copy the [`04_bootstrap.json`](../../../upi/azure/04_bootstrap.json) ARM template locally.

Create the deployment using the `az` client:

```sh
export BOOTSTRAP_URL=`az storage blob url --account-name ${CLUSTER_NAME}sa --account-key $ACCOUNT_KEY -c "files" -n "bootstrap.ign" -o tsv`
export BOOTSTRAP_IGNITION=`jq -rcnM --arg v "2.2.0" --arg url $BOOTSTRAP_URL '{ignition:{version:$v,config:{replace:{source:$url}}}}' | base64 -w0`

az group deployment create -g $RESOURCE_GROUP \
  --template-file "04_bootstrap.json" \
  --parameters bootstrapIgnition="$BOOTSTRAP_IGNITION" \
  --parameters sshKeyData="$SSH_KEY" \
  --parameters baseName="$INFRA_ID"
```

## Launch the permanent control plane

Copy the [`05_masters.json`](../../../upi/azure/05_masters.json) ARM template locally.

Create the deployment using the `az` client:

```sh
export MASTER_IGNITION=`cat master.ign | base64`

az group deployment create -g $RESOURCE_GROUP \
  --template-file "05_masters.json" \
  --parameters masterIgnition="$MASTER_IGNITION" \
  --parameters sshKeyData="$SSH_KEY" \
  --parameters privateDNSZoneName="${CLUSTER_NAME}.${BASE_DOMAIN}" \
  --parameters baseName="$INFRA_ID"
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
export WORKER_IGNITION=`cat worker.ign | base64`

az group deployment create -g $RESOURCE_GROUP \
  --template-file "06_workers.json" \
  --parameters workerIgnition="$WORKER_IGNITION" \
  --parameters sshKeyData="$SSH_KEY" \
  --parameters baseName="$INFRA_ID"
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

Create DNS records in the *public* and *private* zones pointing at the ingress load balancer. Use A, CNAME, etc. records, as you see fit.
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

Finally, add a `*.apps` record to the *private* DNS zone:

```sh
export PUBLIC_IP_ROUTER=`oc -n openshift-ingress get service router-default --no-headers | awk '{print $4}'`
az network private-dns record-set a create -g $RESOURCE_GROUP -z ${CLUSTER_NAME}.${BASE_DOMAIN} -n *.apps --ttl 300
az network private-dns record-set a add-record -g $RESOURCE_GROUP -z ${CLUSTER_NAME}.${BASE_DOMAIN} -n *.apps -a $PUBLIC_IP_ROUTER
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
[jqjson]: https://stedolan.github.io/jq/
[yqyaml]: https://yq.readthedocs.io/en/latest/
[ingress-operator]: https://github.com/openshift/cluster-ingress-operator
[machine-api-operator]: https://github.com/openshift/machine-api-operator
[azure-identity]: https://docs.microsoft.com/en-us/azure/architecture/framework/security/identity
[azure-resource-group]: https://docs.microsoft.com/en-us/azure/azure-resource-manager/management/overview#resource-groups
[azure-dns]: https://docs.microsoft.com/en-us/azure/dns/dns-overview
[kubernetes-service-load-balancers-exclude-masters]: https://github.com/kubernetes/kubernetes/issues/65618
