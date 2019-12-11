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
* python dotmap library: install it with `pip install dotmap`

## Create an install config

Create an install configuration as for [the usual approach](install.md#create-configuration):

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

### Extract data from install config

Some data from the install configuration file will be used on later steps. Export them as environment variables with:

```sh
export CLUSTER_NAME=`yq -r .metadata.name install-config.yaml`
export AZURE_REGION=`yq -r .platform.azure.region install-config.yaml`
export SSH_KEY=`yq -r .sshKey install-config.yaml | xargs`
export BASE_DOMAIN=`yq -r .baseDomain install-config.yaml`
export BASE_DOMAIN_RESOURCE_GROUP=`yq -r .platform.azure.baseDomainResourceGroupName install-config.yaml`
export RESOURCE_GROUP=$CLUSTER_NAME
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

### Customize manifests

The manifests need to reflect the resources to be created by the [Azure Resource Manager][azuretemplates] template, e.g. the
resource group name. Also, we don't want [the ingress operator][ingress-operator] to create DNS records (we're going to
do it manually) so we need to remove the `privateZone` and `publicZone` sections from the DNS configuration in manifests.

A Python script named `setup-manifests.py` is provided to help with these and a few other changes required in manifests. This script
works "as is" for a regular UPI deployment, but can be edited to reflect your own needs.

For example, if you'd like to use a Virtual Network and subnets already existing in your organization (instead of a having new one
created specifically for this deployment), you could either edit the `manifests/cluster-infrastructure-02-config.yml` file manually,
or edit the script to provide their names and the resource group where they exist.

```python3
yamlx['status']['platformStatus']['azure']['networkResourceGroupName'] = "RESOURCE-GROUP-OF-MY-VNET"
yamlx['status']['platformStatus']['azure']['virtualNetwork'] = "MY-VNET-NAME"
yamlx['status']['platformStatus']['azure']['controlPlaneSubnet'] = resource_group + "MY-CONTROL-PLANE-SUBNET-NAME"
yamlx['status']['platformStatus']['azure']['computeSubnet'] = resource_group + "MY-COMPUTE-NODES-SUBNET-NAME"
```

Finally, run the script to have the changes applies to your manifests:

```sh
python3 setup-manifests.py $RESOURCE_GROUP
```

### Remove control plane machines and machinesets

Remove the control plane machines and compute machinesets from the manifests.
We'll be providing those ourselves and don't want to involve the [machine-API operator][machine-api-operator].

```sh
rm -f openshift/99_openshift-cluster-api_master-machines-*.yaml
rm -f openshift/99_openshift-cluster-api_worker-machineset-*.yaml
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
├── 01_vnet.json
├── 02_storage.json
├── 03_infra.json
├── 04_bootstrap.json
├── 05_masters.json
├── 06_workers.json
├── auth
│   ├── kubeadmin-password
│   └── kubeconfig
├── bootstrap.ign
├── master.ign
├── metadata.json
├── setup-bootstrap-ignition.py
├── setup-manifests.py
└── worker.ign
```

## Infra ID

The OpenShift cluster has been assigned an identifier in the form of `<cluster name>-<random string>`. You do not need this for anything in this example, but it is a good idea to keep it around.
At this point you can see the various metadata about your future cluster in `metadata.json`.

The Infra ID is under the `infraID` key:

```console
$ export INFRA_ID=$(jq -r .infraID metadata.json)
$ echo $INFRA_ID
openshift-vw4j5
```

## Create The Resource Group and identity

All resources created in this Azure deployment will exist as part of a [resource group][azure-resource-group]. Use the commands
below to create it in the selected Azure region. In this example we're going to use the cluster name as the unique
resource group name, but feel free to choose any other name and export it in the `RESOURCE_GROUP` environment variable,
which will be used in the subsequent steps.

```sh
az group create --name $RESOURCE_GROUP --location $AZURE_REGION
```

Also, create an identity which will be used to grant the required access to cluster operators.

```sh
az identity create -g $RESOURCE_GROUP -n ${RESOURCE_GROUP}-identity
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

Save the blob URL of the copied image for later use:

```sh
export VHD_BLOB_URL=`az storage blob url --account-name ${CLUSTER_NAME}sa --account-key $ACCOUNT_KEY -c vhd -n "rhcos.vhd" -o tsv`
```

### Upload the bootstrap ignition

Create a blob storage container and upload the generated `bootstrap.ign` file:

```sh
az storage container create --name files --account-name ${CLUSTER_NAME}sa --public-access blob
az storage blob upload --account-name ${CLUSTER_NAME}sa --account-key $ACCOUNT_KEY -c "files" -f "bootstrap.ign" -n "bootstrap.ign"
```

Save the blob URL of the copied file for later use:

```sh
export BOOTSTRAP_URL=`az storage blob url --account-name ${CLUSTER_NAME}sa --account-key $ACCOUNT_KEY -c "files" -n "bootstrap.ign" -o tsv`
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
export PRINCIPAL_ID=`az identity show -g $RESOURCE_GROUP -n ${RESOURCE_GROUP}-identity --query principalId --out tsv`
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

```sh
az group deployment create -g $RESOURCE_GROUP \
  --template-file "01_vnet.json"
```

Link the VNet just created to the private DNS zone:

```sh
az network private-dns link vnet create -g $RESOURCE_GROUP -z ${CLUSTER_NAME}.${BASE_DOMAIN} -n ${CLUSTER_NAME}-network-link -v "${RESOURCE_GROUP}-vnet" -e false
```

## Deploy the image

```sh
az group deployment create -g $RESOURCE_GROUP \
  --template-file "02_storage.json" \
  --parameters vhdBlobURL="${VHD_BLOB_URL}"
```

## Deploy the load balancers

Deploy the load balancers and public IP addresses:

```sh
az group deployment create -g $RESOURCE_GROUP \
  --template-file "03_infra.json" \
  --parameters privateDNSZoneName="${CLUSTER_NAME}.${BASE_DOMAIN}"
```

Create an `api` DNS record in the *public* zone for the API public load balancer. Note that the `BASE_DOMAIN_RESOURCE_GROUP` must point to the resource group where the public DNS zone exists.

```sh
export PUBLIC_IP=`az network public-ip list -g $RESOURCE_GROUP --query "[?name=='${CLUSTER_NAME}-master-pip'] | [0].ipAddress" -o tsv`
az network dns record-set a add-record -g $BASE_DOMAIN_RESOURCE_GROUP -z ${CLUSTER_NAME}.${BASE_DOMAIN} -n api -a $PUBLIC_IP --ttl 60
```

Or, in case of adding this cluster to an already existing public zone, use instead:

```sh
export PUBLIC_IP=`az network public-ip list -g $RESOURCE_GROUP --query "[?name=='${CLUSTER_NAME}-master-pip'] | [0].ipAddress" -o tsv`
az network dns record-set a add-record -g $BASE_DOMAIN_RESOURCE_GROUP -z ${BASE_DOMAIN} -n api.${CLUSTER_NAME} -a $PUBLIC_IP --ttl 60
```

## Launch the temporary cluster bootstrap

```sh
export BOOTSTRAP_IGNITION=`python3 setup-bootstrap-ignition.py $BOOTSTRAP_URL`

az group deployment create -g $RESOURCE_GROUP \
  --template-file "04_bootstrap.json" \
  --parameters bootstrapIgnition="$BOOTSTRAP_IGNITION" \
  --parameters sshKeyData="$SSH_KEY" \
  --parameters privateDNSZoneName="${CLUSTER_NAME}.${BASE_DOMAIN}"
```

## Launch the permanent control plane

```sh
export MASTER_IGNITION=`cat master.ign | base64`

az group deployment create -g $RESOURCE_GROUP \
  --template-file "05_masters.json" \
  --parameters masterIgnition="$MASTER_IGNITION" \
  --parameters sshKeyData="$SSH_KEY" \
  --parameters privateDNSZoneName="${CLUSTER_NAME}.${BASE_DOMAIN}"
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
az vm stop -g $RESOURCE_GROUP --name ${RESOURCE_GROUP}-bootstrap
az vm deallocate -g $RESOURCE_GROUP --name ${RESOURCE_GROUP}-bootstrap
az vm delete -g $RESOURCE_GROUP --name ${RESOURCE_GROUP}-bootstrap --yes
az disk delete -g $RESOURCE_GROUP --name ${RESOURCE_GROUP}-bootstrap_OSDisk --no-wait --yes
az network nic delete -g $RESOURCE_GROUP --name ${RESOURCE_GROUP}-bootstrap-nic --no-wait
az storage blob delete --account-key $ACCOUNT_KEY --account-name ${CLUSTER_NAME}sa --container-name files --name bootstrap.ign
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

```sh
export WORKER_IGNITION=`cat worker.ign | base64`

az group deployment create -g $RESOURCE_GROUP \
  --template-file "06_workers.json" \
  --parameters workerIgnition="$WORKER_IGNITION" \
  --parameters sshKeyData="$SSH_KEY"
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
