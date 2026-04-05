# Install: User-Provisioned Infrastructure
The steps for performing a user-provisioned infrastructure install are outlined here. Several
[Infrastructure Manager][inframanager] templates are provided to assist in
completing these steps or to help model your own. You are also free to create
the required resources through other methods; the templates are just an
example.

## Prerequisites
* all prerequisites from [README](README.md)
* deployment manager to Infrastructure Manager [migration guide](https://docs.cloud.google.com/deployment-manager/docs/dm-convert/convert)
* the following binaries installed and in $PATH:
  * gcloud
  * python
* Python Packages
  * PyYAML
* gcloud authenticated to an account with [additional](iam.md) roles:
  * Cloud Infrastructure Manager Admin
  * Service Account Key Admin
* the following API Services enabled:
  * Cloud Infrastructure Manager API (config.googleapis.com)

## Create Ignition configs
The machines will be started manually. Therefore, it is required to generate
the bootstrap and machine Ignition configs and store them for later steps.
Use a [staged install](../overview.md#multiple-invocations) to enable desired customizations.

### Create an install config
Create an install configuration as for [the usual approach](install.md#create-configuration).

If you are installing into a [Shared VPC (XPN)][sharedvpc],
skip this step and create the `install-config.yaml` manually using the documentation references/examples.

```console
$ openshift-install create install-config
? SSH Public Key /home/user_id/.ssh/id_rsa.pub
? Platform gcp
? Project ID example-project
? Region us-east1
? Base Domain example.com
? Cluster Name openshift
? Pull Secret [? for help]
```

### Empty the compute pool (optional)
If you do not want the cluster to provision compute machines, edit the resulting `install-config.yaml` to set `replicas` to 0 for the `compute` pool.

```console
python -c '
import yaml;
path = "install-config.yaml";
data = yaml.full_load(open(path));
data["compute"][0]["replicas"] = 0;
open(path, "w").write(yaml.dump(data, default_flow_style=False))'
```

```console
compute:
- architecture: amd64
  hyperthreading: Enabled
  name: worker
  platform: {}
  replicas: 0
```

### Enable private cluster setting (optional)
If you want to provision a private cluster, edit the resulting `install-config.yaml` to set `publish` to `Internal`.


```console
python -c '
import yaml;
path = "install-config.yaml";
data = yaml.full_load(open(path));
data["publish"] = "Internal";
open(path, "w").write(yaml.dump(data, default_flow_style=False))'
```

```console
publish: Internal
```

### Export the publish policy

```console
export PUBLISH=$(yq -r .publish install-config.yaml)
```

### Create manifests
Create manifest to enable customizations which are not exposed via the install configuration.

```console
$ openshift-install create manifests
INFO Consuming "Install Config" from target directory
```

### Remove control plane machines
Remove the control plane machines and machinesets from the manifests.
We'll be providing those ourselves and don't want to involve [the machine-API operator][machine-api-operator].

```sh
rm -f openshift/99_openshift-cluster-api_master-machines-*.yaml
rm -f openshift/99_openshift-machine-api_master-control-plane-machine-set.yaml
```

### Remove compute machinesets (optional)
If you do not want the cluster to provision compute machines, remove the compute machinesets from the manifests as well.

```sh
rm -f openshift/99_openshift-cluster-api_worker-machineset-*.yaml
```

### Make control-plane nodes unschedulable
Currently [emptying the compute pools](#empty-the-compute-pool-optional) makes control-plane nodes schedulable.
But due to a [Kubernetes limitation][kubernetes-service-load-balancers-exclude-masters], router pods running on control-plane nodes will not be reachable by the ingress load balancer.
Update the scheduler configuration to keep router pods and other workloads off the control-plane nodes:

```sh
python -c '
import yaml;
path = "manifests/cluster-scheduler-02-config.yml";
data = yaml.full_load(open(path));
data["spec"]["mastersSchedulable"] = False;
open(path, "w").write(yaml.dump(data, default_flow_style=False))'
```

```console
spec:
  mastersSchedulable: false
```

### Remove DNS Zones (optional)
If you don't want [the ingress operator][ingress-operator] to create DNS records on your behalf, remove the `privateZone` and `publicZone` sections from the DNS configuration.
If you do so, you'll need to [add ingress DNS records manually](#add-the-ingress-dns-records) later on.

If you are installing into a [Shared VPC (XPN)][sharedvpc],
remove the `privateZone` section from the DNS configuration.
The `publicZone` will not exist because of `publish: Internal` in `install-config.yaml`.
Remove the `publicZone` line from the command to avoid an error.

```sh
python -c '
import yaml;
path = "manifests/cluster-dns-02-config.yml";
data = yaml.full_load(open(path));
del data["spec"]["publicZone"];
del data["spec"]["privateZone"];
open(path, "w").write(yaml.dump(data, default_flow_style=False))'
```

```console
spec:
  baseDomain: example.com
```

### Update the cloud-provider manifest ([Shared VPC (XPN)][sharedvpc] only)
If you are installing into a [Shared VPC (XPN)][sharedvpc],
update the cloud provider configuration so it understands the network and subnetworks are in a different project (host project).
Otherwise skip this step.

```sh
export HOST_PROJECT="example-shared-vpc"
export HOST_PROJECT_NETWORK_NAME="example-network"
export HOST_PROJECT_COMPUTE_SUBNET_NAME="example-worker-subnet"

sed -i "s/    subnetwork-name.*/    network-project-id = ${HOST_PROJECT}\\n    network-name    = ${HOST_PROJECT_NETWORK_NAME}\\n    subnetwork-name = ${HOST_PROJECT_COMPUTE_SUBNET_NAME}/" manifests/cloud-provider-config.yaml
```

```console
  config: |+
    [global]
    project-id      = example-project
    regional        = true
    multizone       = true
    node-tags       = opensh-ptzzx-master
    node-tags       = opensh-ptzzx-worker
    node-instance-prefix = opensh-ptzzx
    external-instance-groups-prefix = opensh-ptzzx
    network-project-id = example-shared-vpc
    network-name    = example-network
    subnetwork-name = example-worker-subnet
```

### Enable external ingress (optional)
If you are installing into a [Shared VPC (XPN)][sharedvpc],
and you set `publish: Internal` in the `install-config.yaml` but really wanted `publish: External`
then edit the `cluster-ingress-default-ingresscontroller.yaml` manifest to enable external ingress.

```sh
python -c '
import yaml;
path = "manifests/cluster-ingress-default-ingresscontroller.yaml";
data = yaml.full_load(open(path));
data["spec"]["endpointPublishingStrategy"]["loadBalancer"]["scope"] = "External";
open(path, "w").write(yaml.dump(data, default_flow_style=False))'
```

```console
 spec:
  endpointPublishingStrategy:
    loadBalancer:
      scope: External
```

### Create Ignition configs
Now we can create the bootstrap Ignition configs.

```console
$ openshift-install create ignition-configs
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

### Extract infrastructure name from Ignition metadata
By default, Ignition generates a unique cluster identifier comprised of the
cluster name specified during the invocation of the installer and a short
string known internally as the infrastructure name. These values are seeded
in the initial manifests within the Ignition configuration. To use the output
of the default, generated `ignition-configs` extracting the internal
infrastructure name is necessary.

An example of a way to get this is below:

```console
$ jq -r .infraID metadata.json
openshift-vw9j6
```

## Export variables to be used in examples below.

```sh
export BASE_DOMAIN='example.com'
export BASE_DOMAIN_ZONE_NAME='example'
export NETWORK_CIDR='10.0.0.0/16'
export MASTER_SUBNET_CIDR='10.0.0.0/17'
export WORKER_SUBNET_CIDR='10.0.128.0/17'

export KUBECONFIG=auth/kubeconfig
export CLUSTER_NAME=$(jq -r .clusterName metadata.json)
export INFRA_ID=$(jq -r .infraID metadata.json)
export PROJECT_NAME=$(jq -r .gcp.projectID metadata.json)
export REGION=$(jq -r .gcp.region metadata.json)
export ZONE_0=$(gcloud compute regions describe ${REGION} --format=json | jq -r '.zones[0]' | cut -d "/" -f9)
export ZONE_1=$(gcloud compute regions describe ${REGION} --format=json | jq -r '.zones[1]' | cut -d "/" -f9)
export ZONE_2=$(gcloud compute regions describe ${REGION} --format=json | jq -r '.zones[2]' | cut -d "/" -f9)

export MASTER_IGNITION=$(cat master.ign)
export WORKER_IGNITION=$(cat worker.ign)

# Fill in your service account email used for the installation. This may be the service account
# found in the osServiceAccounts.json file or it could be another service account found in the IAM
# section of the google cloud web console.
SERVICE_ACCOUNT_EMAIL=""
export INSTALL_SERVICE_ACCOUNT="projects/${PROJECT_NAME}/serviceAccounts/${SERVICE_ACCOUNT_EMAIL}"

export CLUSTER_DOMAIN="${CLUSTER_NAME}.${BASE_DOMAIN}"
```

## Create the VPC
Create the VPC, network, and subnets for the cluster.
This step can be skipped if installing into a pre-existing VPC, such as a [Shared VPC (XPN)][sharedvpc].

Copy [`01_vpc`](../../../upi/gcp/01_vpc) locally. The directory contains the terraform source file for this stage.

Create the deployment using gcloud.

```console
gcloud infra-manager deployments apply upi-network \
--location=${REGION} \
--input-values=infra_id=${INFRA_ID},project=${PROJECT_NAME},region=${REGION},master_subnet_cidr=${MASTER_SUBNET_CIDR},worker_subnet_cidr=${WORKER_SUBNET_CIDR} \
--project=${PROJECT_NAME} \
--local-source=./01_vpc \
--service-account=${INSTALL_SERVICE_ACCOUNT}
```
- `infra_id`: the infrastructure name (INFRA_ID above)
- `region`/`location`: the region to deploy the cluster into (for example us-east1)
- `project`: Name of the GCP service project
- `master_subnet_cidr`: the CIDR for the master subnet (for example 10.0.0.0/17)
- `worker_subnet_cidr`: the CIDR for the worker subnet (for example 10.0.128.0/17)
- `service-account`: Service account used for gcloud infrastructure creation

**Note**: You may add the `--async` option to immediately apply and allow processing to continue in the background.

## Configure VPC variables
Configure the variables based on the VPC created above through terraform.
If you are using a pre-existing VPC, such as a [Shared VPC (XPN)][sharedvpc], set these to the `.selfLink` of the targeted resources.

```sh
export CLUSTER_NETWORK=$(gcloud compute networks describe ${INFRA_ID}-network --format json | jq -r .selfLink)
export CONTROL_SUBNET=$(gcloud compute networks subnets describe ${INFRA_ID}-master-subnet --region=${REGION} --format json | jq -r .selfLink)
export COMPUTE_SUBNET=$(gcloud compute networks subnets describe ${INFRA_ID}-worker-subnet --region=${REGION} --format json | jq -r .selfLink)
```

## Create DNS entries
Create the DNS zone and load balancers for the cluster.
You can exclude the DNS zone or external load balancer by removing section from `02_dns` and the file `02_lb_ext`.
If you choose to exclude the DNS zone, you will need to create it some other way and ensure it is populated with the necessary records as documented below.

If you are installing into a [Shared VPC (XPN)][sharedvpc],
exclude the DNS section as it must be created in the host project.

Copy [`02_dns`](../../../upi/gcp/02_dns) locally. The directory contains the terraform source file for creating the dns zone during this stage.

Create the deployment using gcloud.

```console
gcloud infra-manager deployments apply upi-dns \
--location=${REGION} \
--input-values=infra_id=${INFRA_ID},project=${PROJECT_NAME},region=${REGION},cluster_domain=${CLUSTER_DOMAIN},cluster_network=${CLUSTER_NETWORK} \
--project=${PROJECT_NAME} \
--local-source=./02_dns \
--service-account=${INSTALL_SERVICE_ACCOUNT}
```

- `cluster_domain`: the domain for the cluster (for example openshift.example.com)
- `cluster_network`: the URI to the cluster network

## Create the external load balancer
**If you are installing a private cluster the following section should be skipped.**
Copy [`02_lb_ext`](../../../upi/gcp/02_lb_ext) locally. The directory contains the terraform source file for creating the external load balancer.

Create the deployment using gcloud.

```console
gcloud infra-manager deployments apply upi-lb-ext \
--location=${REGION} \
--input-values=infra_id=${INFRA_ID},project=${PROJECT_NAME},region=${REGION} \
--project=${PROJECT_NAME} \
--local-source=./02_lb_ext\
--service-account=${INSTALL_SERVICE_ACCOUNT}
```

## Create the internal load balancer

Copy [`02_lb_int`](../../../upi/gcp/02_lb_int) locally. The directory contains the terraform source file for creating the internal load balancer.

Create the deployment using gcloud.

```console
gcloud infra-manager deployments apply upi-lb-int \
--location=${REGION} \
--input-values=infra_id=${INFRA_ID},project=${PROJECT_NAME},region=${REGION},cluster_network=${CLUSTER_NETWORK},control_subnet=${CONTROL_SUBNET},zone_0=${ZONE_0},zone_1=${ZONE_1},zone_2=${ZONE_2} \
--project=${PROJECT_NAME} \
--local-source=./02_lb_int \
--service-account=${INSTALL_SERVICE_ACCOUNT}
```

- `cluster_network`: the URI to the cluster network
- `control_subnet`: the URI to the control subnet
- `zone_0`: the first zone to deploy the control plane instances into (for example us-east1-b)
- `zone_1`: the second zone to deploy the control plane instances into (for example us-east1-c)
- `zone_2`: the third zone to deploy the control plane instances into (for example us-east1-d)

**Note**: Terraform accepts lists, but the infra-manager will only accept scalar values. If you would like to alter the number
of zones, please edit the terraform file.

## Configure infra variables
If you did not create a public load balancer (excluded the `02-lb-ext` section) above, then skip creating/exporting `CLUSTER_PUBLIC_IP`.

```sh
export CLUSTER_IP=$(gcloud compute addresses describe ${INFRA_ID}-cluster-ip --region=${REGION} --format json | jq -r .address)
export CLUSTER_PUBLIC_IP=$(gcloud compute addresses describe ${INFRA_ID}-cluster-public-ip --region=${REGION} --format json | jq -r .address)
```

## Add DNS entries
If you are installing into a [Shared VPC (XPN)][sharedvpc],
use the `--account` and `--project` parameters to perform these actions in the host project.

### Add internal DNS entries

```sh
if [ -f transaction.yaml ]; then rm transaction.yaml; fi
gcloud dns record-sets transaction start --zone ${INFRA_ID}-private-zone
gcloud dns record-sets transaction add ${CLUSTER_IP} --name api.${CLUSTER_NAME}.${BASE_DOMAIN}. --ttl 60 --type A --zone ${INFRA_ID}-private-zone
gcloud dns record-sets transaction add ${CLUSTER_IP} --name api-int.${CLUSTER_NAME}.${BASE_DOMAIN}. --ttl 60 --type A --zone ${INFRA_ID}-private-zone
gcloud dns record-sets transaction execute --zone ${INFRA_ID}-private-zone
```

### Add external DNS entries (optional)
If you deployed external load balancers with `02_lb_ext.tf`, you can deploy external DNS entries.

```sh
if [ -f transaction.yaml ]; then rm transaction.yaml; fi
gcloud dns record-sets transaction start --zone ${BASE_DOMAIN_ZONE_NAME}
gcloud dns record-sets transaction add ${CLUSTER_PUBLIC_IP} --name api.${CLUSTER_NAME}.${BASE_DOMAIN}. --ttl 60 --type A --zone ${BASE_DOMAIN_ZONE_NAME}
gcloud dns record-sets transaction execute --zone ${BASE_DOMAIN_ZONE_NAME}
```

## Create firewall rules and IAM roles
Create the firewall rules and IAM roles for the cluster.

If you are installing into a [Shared VPC (XPN)][sharedvpc], skip the deployment of the firewall rules and IAM service accounts.

Copy [`03_security`](../../../upi/gcp/03_security) locally. The directory contains the terraform source file for this stage.

Create the deployment using gcloud.

```console
gcloud infra-manager deployments apply upi-security \
--location=${REGION} \
--project=${PROJECT_NAME} \
--local-source=./03_security \
--input-values=infra_id=${INFRA_ID},project=${PROJECT_NAME},region=${REGION},cluster_network=${CLUSTER_NETWORK},network_cidr=${NETWORK_CIDR} \
--service-account=${INSTALL_SERVICE_ACCOUNT}
```

- `cluster_network`: the URI to the cluster network
- `network_cidr`: the CIDR of the vpc network (for example 10.0.0.0/16)
- `allowed_external_cidr`: [optional] limits access to the cluster API and ssh to the bootstrap host. (for example External: 0.0.0.0/0, Internal: ${NETWORK_CIDR})

## Configure security variables
Configure the variables based on the `03_security.tf` deployment.
If you excluded the IAM section, ensure these are set to the `.email` of their associated resources.

```sh
export MASTER_SERVICE_ACCOUNT=$(gcloud iam service-accounts list --filter "email~^${INFRA_ID}-m@${PROJECT_NAME}." --format json | jq -r '.[0].email')
export WORKER_SERVICE_ACCOUNT=$(gcloud iam service-accounts list --filter "email~^${INFRA_ID}-w@${PROJECT_NAME}." --format json | jq -r '.[0].email')
```

## Add required roles to IAM service accounts
The templates do not create the policy bindings, but we do it manually below.

If you are installing into a [Shared VPC (XPN)][sharedvpc],
ensure these service accounts have `roles/compute.networkUser` access to each of the host project subnets used by the cluster so the instances can use the networks.
Also ensure the master service account has `roles/compute.networkViewer` access to the host project itself so the gcp-cloud-provider can look for firewall settings as part of ingress controller operations.

```sh
gcloud projects add-iam-policy-binding ${PROJECT_NAME} --member "serviceAccount:${MASTER_SERVICE_ACCOUNT}" --role "roles/compute.instanceAdmin"
gcloud projects add-iam-policy-binding ${PROJECT_NAME} --member "serviceAccount:${MASTER_SERVICE_ACCOUNT}" --role "roles/compute.networkAdmin"
gcloud projects add-iam-policy-binding ${PROJECT_NAME} --member "serviceAccount:${MASTER_SERVICE_ACCOUNT}" --role "roles/compute.securityAdmin"
gcloud projects add-iam-policy-binding ${PROJECT_NAME} --member "serviceAccount:${MASTER_SERVICE_ACCOUNT}" --role "roles/iam.serviceAccountUser"
gcloud projects add-iam-policy-binding ${PROJECT_NAME} --member "serviceAccount:${MASTER_SERVICE_ACCOUNT}" --role "roles/storage.admin"

gcloud projects add-iam-policy-binding ${PROJECT_NAME} --member "serviceAccount:${WORKER_SERVICE_ACCOUNT}" --role "roles/compute.viewer"
gcloud projects add-iam-policy-binding ${PROJECT_NAME} --member "serviceAccount:${WORKER_SERVICE_ACCOUNT}" --role "roles/storage.admin"
```

## Generate a service-account-key for signing the bootstrap.ign url

```sh
gcloud iam service-accounts keys create service-account-key.json --iam-account=${MASTER_SERVICE_ACCOUNT}
```

## Create the cluster image.
Locate the RHCOS image source and create a cluster image.

```sh
if openshift-install coreos print-stream-json 2>/tmp/err.txt >coreos.json; then
  # Note: for OSX users, set environment variable OCP_ARCH.
 jq '.architectures.'"$(echo "$OCP_ARCH" | sed 's/amd64/x86_64/;s/arm64/aarch64/')"'.images.gcp' < coreos.json > gcp.json
  source_image="$(jq -r .name < gcp.json)"
  source_project="$(jq -r .project < gcp.json)"
  rm -f coreos.json gcp.json
  gcloud compute images create "${INFRA_ID}-rhcos-image" --source-image="${source_image}" --source-image-project="${source_project}"
  export CLUSTER_IMAGE=(`gcloud compute images describe ${INFRA_ID}-rhcos-image --format json | jq -r .selfLink`)
else
  export IMAGE_SOURCE=$(curl https://raw.githubusercontent.com/openshift/installer/refs/heads/main/data/data/coreos/coreos-rhel-10.json | jq -r '.architectures.x86_64.images.gcp')
  export IMAGE_NAME=$(echo "${IMAGE_SOURCE}" | jq -r '.name')
  export IMAGE_PROJECT=$(echo "${IMAGE_SOURCE}" | jq -r '.project')
  export CLUSTER_IMAGE=$(gcloud compute images describe ${IMAGE_NAME} --project ${IMAGE_PROJECT} --format json | jq -r .selfLink)
fi
```

## Upload the bootstrap.ign to a new bucket
Create a bucket and upload the bootstrap.ign file.

```sh
gcloud storage buckets create "gs://${INFRA_ID}-bootstrap-ignition"
gcloud storage cp bootstrap.ign "gs://${INFRA_ID}-bootstrap-ignition/"
gcloud storage ls "gs://${INFRA_ID}-bootstrap-ignition/bootstrap.ign"
```

Create a signed URL for the bootstrap instance to use to access the Ignition
config. Export the URL from the output as a variable.

```sh
BOOTSTRAP_IGN="$(gcloud storage sign-url --duration=1h --private-key-file=service-account-key.json "gs://${INFRA_ID}-bootstrap-ignition/bootstrap.ign" | grep "^signed_url:" | awk '{print $2}')"
```

## Launch temporary bootstrap resources

Copy [`04_bootstrap`](../../../upi/gcp/04_bootstrap) locally. The directory contains the terraform source file for this stage.

Create the deployment using gcloud.

### Determine publish policy
```sh
public_cluster=true
if [[ "${PUBLISH}" == "Internal" ]]; then
public_cluster=false
fi
```

```console
gcloud infra-manager deployments apply upi-bootstrap \
--location=${REGION} \
--project=${PROJECT_NAME} \
--local-source=./04_bootstrap \
--input-values=infra_id=${INFRA_ID},project=${PROJECT_NAME},region=${REGION},zone=${ZONE_0},cluster_network=${CLUSTER_NETWORK},subnet=${CONTROL_SUBNET},image=${CLUSTER_IMAGE},bootstrap_ign="${BOOTSTRAP_IGN}",is_public_cluster=${public_cluster} \
--service-account=${INSTALL_SERVICE_ACCOUNT}
```
- `zone`: the zone to deploy the bootstrap instance into (for example us-east1-b)
- `cluster_network`: the URI to the cluster network
- `subnet`: the URI to the control subnet
- `image`: the URI to the RHCOS image
- `machine_type`: [optional] the machine type of the instance (for example n1-standard-4)
- `root_volume_size`: [optional] the size (in GB) for the instance (for example 128)
- `bootstrap_ign`: the URL output when creating a signed URL above.
- `is_public_cluster`: boolean value for creating public clusters.

You can add custom tags to `04_bootstrap.tf` as needed

```console
  tags = [
    "${var.infra_id}-master",
    "${var.infra_id}-bootstrap",
    "custom-tag-example"
  ]
```

## Add the bootstrap instance to the load balancers
Add the bootstrap node manually.

### Add bootstrap instance to internal load balancer instance group

```sh
gcloud compute instance-groups unmanaged add-instances ${INFRA_ID}-bootstrap-ig --zone=${ZONE_0} --instances=${INFRA_ID}-bootstrap
```

### Add bootstrap instance group to the internal load balancer backend service

```sh
gcloud compute backend-services add-backend ${INFRA_ID}-api-internal --region=${REGION} --instance-group=${INFRA_ID}-bootstrap-ig --instance-group-zone=${ZONE_0}
```

## Launch permanent control plane

### Copy the ignition file
The ignition data cannot be passed through the infra-manager `input-values`. The data is read from the `master.ign` file.
The file must be present in the same location as the `.tf` file.

```console
cp master.ign 05_control_plane/master.ign
```

Copy [`05_control_plane`](../../../upi/gcp/05_control_plane) locally. The directory contains the terraform source file for this stage.

Create the deployment using gcloud.

```console
gcloud infra-manager deployments apply upi-control-plane \
--location=${REGION} \
--project=${PROJECT_NAME} \
--local-source=./05_control_plane \
--input-values=infra_id=${INFRA_ID},project=${PROJECT_NAME},region=${REGION},zone_0=${ZONE_0},zone_1=${ZONE_1},zone_2=${ZONE_2},subnet=${CONTROL_SUBNET},image=${CLUSTER_IMAGE},service_account_email=${MASTER_SERVICE_ACCOUNT} \
--service-account=${INSTALL_SERVICE_ACCOUNT}
```
- `zone_0`: the first zone to deploy the control plane instances into (for example us-east1-b)
- `zone_1`: the second zone to deploy the control plane instances into (for example us-east1-c)
- `zone_2`: the third zone to deploy the control plane instances into (for example us-east1-d)
- `subnet`: the URI to the control subnet
- `image`: the URI to the RHCOS image
- `machine_type`: [optional] the machine type of the instance (for example n1-standard-4)
- `disk_size`: [optional] the size (in GB) for the instance (for example 128)
- `disk_type`: [optional] the type of storage disk (for example pd-ssd)
- `service_account_email`: the email address for the master service account created above

**Note**: Terraform accepts lists, but the infra-manager will only accept scalar values. If you would like to alter the number
of zones, please edit the terraform file.

You can add custom tags to `05_control_plane.tf` as needed

```console
  tags = [
    "${var.infra_id}-master",
    "custom-tag-example"
  ]
```

### Remove the temporary ignition file

```console
rm 05_control_plane/master.ign
```

## Add control plane instances to load balancers
Add the control plane nodes manually.

### Add control plane instances to internal load balancer instance groups

```sh
gcloud compute instance-groups unmanaged add-instances ${INFRA_ID}-master-${ZONE_0}-ig --zone=${ZONE_0} --instances=${INFRA_ID}-master-0
gcloud compute instance-groups unmanaged add-instances ${INFRA_ID}-master-${ZONE_1}-ig --zone=${ZONE_1} --instances=${INFRA_ID}-master-1
gcloud compute instance-groups unmanaged add-instances ${INFRA_ID}-master-${ZONE_2}-ig --zone=${ZONE_2} --instances=${INFRA_ID}-master-2
```

### Add control plane instances to external load balancer target pools (optional)
If you deployed external load balancers with `02_lb_ext.tf`, add the control plane instances to the target pool.

```sh
gcloud compute target-pools add-instances ${INFRA_ID}-api-target-pool --instances-zone="${ZONE_0}" --instances=${INFRA_ID}-master-0
gcloud compute target-pools add-instances ${INFRA_ID}-api-target-pool --instances-zone="${ZONE_1}" --instances=${INFRA_ID}-master-1
gcloud compute target-pools add-instances ${INFRA_ID}-api-target-pool --instances-zone="${ZONE_2}" --instances=${INFRA_ID}-master-2
```

## Launch additional compute nodes
You may create compute nodes by launching individual instances discretely
or by automated processes outside the cluster (e.g. Auto Scaling Groups). You
can also take advantage of the built-in cluster scaling mechanisms and the
machine API in OpenShift, as mentioned [above](#create-ignition-configs). In
this example, we'll manually launch two worker instances via the Infrastructure Manager
template. The number of compute instances can be adjusted by increasing/decreasing the
number of `google_compute_instance` resources in `06_worker.tf`.

### Copy the ignition file
The ignition data cannot be passed through the infra-manager `input-values`. The data is read from the `worker.ign` file.
The file must be present in the same location as the `.tf` file.

```console
cp worker.ign 06_worker/worker.ign
```

Copy [`06_worker`](../../../upi/gcp/06_worker) locally. The directory contains the terraform source file for this stage.

Create the deployment using gcloud.

```console
gcloud infra-manager deployments apply upi-worker \
--location=${REGION} \
--project=${PROJECT_NAME} \
--local-source=./06_worker \
--input-values=infra_id=${INFRA_ID},project=${PROJECT_NAME},region=${REGION},zone_0=${ZONE_0},zone_1=${ZONE_1},subnet=${COMPUTE_SUBNET},image=${CLUSTER_IMAGE},service_account_email=${WORKER_SERVICE_ACCOUNT} \
--service-account=${INSTALL_SERVICE_ACCOUNT}
```
- `zone_0`: the first zone to deploy the compute instances into (for example us-east1-b)
- `zone_1`: the second zone to deploy the compute instances into (for example us-east1-c)
- `subnet`: the URI to the compute subnet
- `image`: the URI to the RHCOS image
- `machine_type`: [optional] the machine type of the instance (for example n1-standard-4)
- `disk_size`: [optional] the size (in GB) for the instance (for example 128)
- `disk_type`: [optional] the type of storage disk (for example pd-ssd)
- `service_account_email`: the email address for the worker service account created above

**Note**: Terraform accepts lists, but the infra-manager will only accept scalar values. If you would like to alter the number
of zones, please edit the terraform file.

You can add custom tags to `06_worker.tf` as needed

```console
  tags = [
    "${var.infra_id}-worker",
    "custom-tag-example"
  ]
```

### Remove the temporary ignition file

```console
rm 06_worker/worker.ign
```

## Monitor for `bootstrap-complete`

```console
$ openshift-install wait-for bootstrap-complete
INFO Waiting up to 30m0s for the Kubernetes API at https://api.test.example.com:6443...
INFO API v1.12.4+c53f462 up
INFO Waiting up to 30m0s for the bootstrap-complete event...
```

## Destroy bootstrap resources
At this point, you should delete the bootstrap resources.

If you are installing into a [Shared VPC (XPN)][sharedvpc],
it is safe to remove any bootstrap-specific firewall rules at this time.

```sh
gcloud compute backend-services remove-backend ${INFRA_ID}-api-internal --region=${REGION} --instance-group=${INFRA_ID}-bootstrap-ig --instance-group-zone=${ZONE_0}
gcloud storage rm "gs://${INFRA_ID}-bootstrap-ignition/bootstrap.ign"
gcloud storage buckets delete "gs://${INFRA_ID}-bootstrap-ignition"
gcloud infra-manager deployments delete --project=${PROJECT_NAME} --location=${REGION} upi-bootstrap
```

## Approving the CSR requests for nodes
The CSR requests for client and server certificates for nodes joining the cluster will need to be approved by the administrator.
Nodes that have not been provisioned by the cluster need their associated `system:serviceaccount` certificate approved to join the cluster.
You can view them with:

```console
$ oc get csr
NAME        AGE     REQUESTOR                                                                   CONDITION
csr-8b2br   15m     system:serviceaccount:openshift-machine-config-operator:node-bootstrapper   Approved,Issued
csr-8vnps   15m     system:serviceaccount:openshift-machine-config-operator:node-bootstrapper   Approved,Issued
csr-b96j4   25s     system:node:ip-10-0-52-215.us-east-2.compute.internal                       Approved,Issued
csr-bfd72   5m26s   system:node:ip-10-0-50-126.us-east-2.compute.internal                       Pending
csr-c57lv   5m26s   system:node:ip-10-0-95-157.us-east-2.compute.internal                       Pending
...
```

Administrators should carefully examine each CSR request and approve only the ones that belong to the nodes created by them.
CSRs can be approved by name, for example:

```sh
oc adm certificate approve csr-bfd72
```

## Add the Ingress DNS Records
If you removed the DNS Zone configuration [earlier](#remove-dns-zones), you'll need to manually create some DNS records pointing at the ingress load balancer.
You can create either a wildcard `*.apps.{baseDomain}.` or specific records (more on the specific records below).
You can use A, CNAME, etc. records, as you see fit.

If you are installing into a [Shared VPC (XPN)][sharedvpc],
use the `--account` and `--project` parameters to perform these actions in the host project.

### Wait for the ingress-router to create a load balancer and populate the `EXTERNAL-IP`

```console
$ oc -n openshift-ingress get service router-default
NAME             TYPE           CLUSTER-IP      EXTERNAL-IP      PORT(S)                      AGE
router-default   LoadBalancer   172.30.18.154   35.233.157.184   80:32288/TCP,443:31215/TCP   98
```

### Add the internal *.apps DNS record

```sh
export ROUTER_IP=$(oc -n openshift-ingress get service router-default --no-headers | awk '{print $4}')

if [ -f transaction.yaml ]; then rm transaction.yaml; fi
gcloud dns record-sets transaction start --zone ${INFRA_ID}-private-zone
gcloud dns record-sets transaction add ${ROUTER_IP} --name \*.apps.${CLUSTER_NAME}.${BASE_DOMAIN}. --ttl 300 --type A --zone ${INFRA_ID}-private-zone
gcloud dns record-sets transaction execute --zone ${INFRA_ID}-private-zone
```

### Add the external *.apps DNS record (optional)
If you deployed external load balancers with `02_lb_ext.tf`, you can deploy external DNS entries.

```sh
if [ -f transaction.yaml ]; then rm transaction.yaml; fi
gcloud dns record-sets transaction start --zone ${BASE_DOMAIN_ZONE_NAME}
gcloud dns record-sets transaction add ${ROUTER_IP} --name \*.apps.${CLUSTER_NAME}.${BASE_DOMAIN}. --ttl 300 --type A --zone ${BASE_DOMAIN_ZONE_NAME}
gcloud dns record-sets transaction execute --zone ${BASE_DOMAIN_ZONE_NAME}
```

If you prefer to add explicit domains instead of using a wildcard, you can create entries for each of the cluster's current routes:

```console
$ oc get --all-namespaces -o jsonpath='{range .items[*]}{range .status.ingress[*]}{.host}{"\n"}{end}{end}' routes
oauth-openshift.apps.your.cluster.domain.example.com
console-openshift-console.apps.your.cluster.domain.example.com
downloads-openshift-console.apps.your.cluster.domain.example.com
alertmanager-main-openshift-monitoring.apps.your.cluster.domain.example.com
grafana-openshift-monitoring.apps.your.cluster.domain.example.com
prometheus-k8s-openshift-monitoring.apps.your.cluster.domain.example.com
```

## Add the Ingress firewall rules (optional)
If you are installing into a [Shared VPC (XPN)][sharedvpc],
you'll need to manually create some firewall rules for the ingress services.
These rules would normally be created by the ingress controller via the gcp cloud provider.
When the cloud provider detects Shared VPC (XPN), it will instead emit cluster events informing which firewall rules need to be created.
Either create each rule as requested by the events (option A), or create cluster-wide firewall rules for all services (option B).

Use the `--account` and `--project` parameters to perform these actions in the host project.

### Add firewall rules based on cluster events (option A)
When the cluster is first provisioned, and as services are later created and modified, the gcp cloud provider may generate events informing of firewall rules required to be manually created in order to allow access to these services.

```console
Firewall change required by security admin: `gcloud compute firewall-rules create k8s-fw-a26e631036a3f46cba28f8df67266d55 --network example-network --description "{\"kubernetes.io/service-name\":\"openshift-ingress/router-default\", \"kubernetes.io/service-ip\":\"35.237.236.234\"}\" --allow tcp:443,tcp:80 --source-ranges 0.0.0.0/0 --target-tags exampl-fqzq7-master,exampl-fqzq7-worker --project example-project`
```

Create the firewall rules as instructed.

### Add a cluster-wide health check firewall rule. (option B)
Add a single firewall rule to allow the gce health checks to access all the services.
This enables the ingress load balancers to determine the health status of their instances.

```sh
gcloud compute firewall-rules create --allow='tcp:30000-32767,udp:30000-32767' --network="${CLUSTER_NETWORK}" --source-ranges='130.211.0.0/22,35.191.0.0/16,209.85.152.0/22,209.85.204.0/22' --target-tags="${INFRA_ID}-master,${INFRA_ID}-worker" ${INFRA_ID}-ingress-hc
```

### Add a cluster-wide service firewall rule. (option B)
Add a single firewall rule to allow access to all cluster services.
If you want your cluster to be private, you can use `--source-ranges=${NETWORK_CIDR}`.
This rule may need to be updated accordingly when adding services on ports other than `tcp:80,tcp:443`.

```sh
gcloud compute firewall-rules create --allow='tcp:80,tcp:443' --network="${CLUSTER_NETWORK}" --source-ranges="0.0.0.0/0" --target-tags="${INFRA_ID}-master,${INFRA_ID}-worker" ${INFRA_ID}-ingress
```

## Monitor for cluster completion

```console
$ openshift-install wait-for install-complete
INFO Waiting up to 30m0s for the cluster to initialize...
```

Also, you can observe the running state of your cluster pods:

```console
$ oc get clusterversion
NAME      VERSION                                    AVAILABLE   PROGRESSING   SINCE   STATUS
version   4.21.0-0.nightly-multi-2026-03-25-233540   True        False         101m    Cluster version is 4.21.0-0.nightly-multi-2026-03-25-233540

$ oc get clusteroperators
NAME                                       VERSION                                    AVAILABLE   PROGRESSING   DEGRADED   SINCE   MESSAGE
authentication                             4.21.0-0.nightly-multi-2026-03-25-233540   True        False         False      4h31m   
baremetal                                  4.21.0-0.nightly-multi-2026-03-25-233540   True        False         False      4h55m   
cloud-controller-manager                   4.21.0-0.nightly-multi-2026-03-25-233540   True        False         False      4h58m   
cloud-credential                           4.21.0-0.nightly-multi-2026-03-25-233540   True        False         False      4h59m   
cluster-autoscaler                         4.21.0-0.nightly-multi-2026-03-25-233540   True        False         False      4h55m   
config-operator                            4.21.0-0.nightly-multi-2026-03-25-233540   True        False         False      4h56m   
console                                    4.21.0-0.nightly-multi-2026-03-25-233540   True        False         False      4h37m   
control-plane-machine-set                  4.21.0-0.nightly-multi-2026-03-25-233540   True        False         False      4h55m   
csi-snapshot-controller                    4.21.0-0.nightly-multi-2026-03-25-233540   True        False         False      4h55m   
dns                                        4.21.0-0.nightly-multi-2026-03-25-233540   True        False         False      4h54m   
etcd                                       4.21.0-0.nightly-multi-2026-03-25-233540   True        False         False      4h53m   
image-registry                             4.21.0-0.nightly-multi-2026-03-25-233540   True        False         False      4h43m   
ingress                                    4.21.0-0.nightly-multi-2026-03-25-233540   True        False         False      30m     
insights                                   4.21.0-0.nightly-multi-2026-03-25-233540   True        False         False      4h55m   
kube-apiserver                             4.21.0-0.nightly-multi-2026-03-25-233540   True        False         False      4h49m   
kube-controller-manager                    4.21.0-0.nightly-multi-2026-03-25-233540   True        False         False      4h50m   
kube-scheduler                             4.21.0-0.nightly-multi-2026-03-25-233540   True        False         False      4h53m   
kube-storage-version-migrator              4.21.0-0.nightly-multi-2026-03-25-233540   True        False         False      108m    
machine-api                                4.21.0-0.nightly-multi-2026-03-25-233540   True        False         False      4h51m   
machine-approver                           4.21.0-0.nightly-multi-2026-03-25-233540   True        False         False      4h55m   
machine-config                             4.21.0-0.nightly-multi-2026-03-25-233540   True        False         False      4h56m   
marketplace                                4.21.0-0.nightly-multi-2026-03-25-233540   True        False         False      4h55m   
monitoring                                 4.21.0-0.nightly-multi-2026-03-25-233540   True        False         False      4h36m   
network                                    4.21.0-0.nightly-multi-2026-03-25-233540   True        False         False      4h57m   
node-tuning                                4.21.0-0.nightly-multi-2026-03-25-233540   True        False         False      139m    
olm                                        4.21.0-0.nightly-multi-2026-03-25-233540   True        False         False      107m    
openshift-apiserver                        4.21.0-0.nightly-multi-2026-03-25-233540   True        False         False      4h44m   
openshift-controller-manager               4.21.0-0.nightly-multi-2026-03-25-233540   True        False         False      4h44m   
openshift-samples                          4.21.0-0.nightly-multi-2026-03-25-233540   True        False         False      4h44m   
operator-lifecycle-manager                 4.21.0-0.nightly-multi-2026-03-25-233540   True        False         False      4h55m   
operator-lifecycle-manager-catalog         4.21.0-0.nightly-multi-2026-03-25-233540   True        False         False      4h55m   
operator-lifecycle-manager-packageserver   4.21.0-0.nightly-multi-2026-03-25-233540   True        False         False      4h44m   
service-ca                                 4.21.0-0.nightly-multi-2026-03-25-233540   True        False         False      4h56m   
storage                                    4.21.0-0.nightly-multi-2026-03-25-233540   True        False         False      4h54m 

$ oc get pods --all-namespaces
NAMESPACE                                          NAME                                                                 READY   STATUS             RESTARTS         AGE
cert-manager-operator                              cert-manager-operator-controller-manager-664c6bd4-jjzpb              1/1     Running            0                30m
cert-manager-operator                              konflux-fbc-cert-manager-r2v6b                                       1/1     Running            0                30m
cert-manager                                       cert-manager-cainjector-5db4d8c596-gwg6g                             1/1     Running            0                30m
cert-manager                                       cert-manager-ccbf7b774-tthzs                                         1/1     Running            0                30m
cert-manager                                       cert-manager-webhook-54965d9794-2xsn7                                1/1     Running            0                30m
external-secrets-operator                          external-secrets-operator-controller-manager-6fb7f5f866-9x6c8        1/1     Running            0                30m
external-secrets-operator                          konflux-fbc-eso-5qmgr                                                1/1     Running            0                30m
...
openshift-ovn-kubernetes                           ovnkube-node-zwzsj                                                   8/8     Running            9 (115m ago)     123m
openshift-route-controller-manager                 route-controller-manager-6dfdd544f6-5qx2t                            1/1     Running            0                108m
openshift-route-controller-manager                 route-controller-manager-6dfdd544f6-j7smc                            1/1     Running            0                120m
openshift-route-controller-manager                 route-controller-manager-6dfdd544f6-pkhd7                            1/1     Running            0                114m
openshift-service-ca-operator                      service-ca-operator-777dd59695-spq5f                                 1/1     Running            0                114m
openshift-service-ca                               service-ca-56c867b864-chftf                                          1/1     Running            0                108m
```

[inframanager]: https://docs.cloud.google.com/infrastructure-manager/docs
[ingress-operator]: https://github.com/openshift/cluster-ingress-operator
[kubernetes-service-load-balancers-exclude-masters]: https://github.com/kubernetes/kubernetes/issues/65618
[machine-api-operator]: https://github.com/openshift/machine-api-operator
[sharedvpc]: https://cloud.google.com/vpc/docs/shared-vpc
