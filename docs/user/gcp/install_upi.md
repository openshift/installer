# Install: User-Provided Infrastructure
The steps for performing a user-provided infrastructure install are outlined here. Several
[Deployment Manager][deploymentmanager] templates are provided to assist in
completing these steps or to help model your own. You are also free to create
the required resources through other methods; the templates are just an
example.

## Prerequisites
* all prerequisites from [README](README.md)
* the following binaries installed and in $PATH:
  * gcloud
  * gsutil
* gcloud authenticated to an account with [additional](iam.md) roles:
  * Deployment Manager Editor
  * Service Account Key Admin
* the following API Services enabled:
  * Cloud Deployment Manager V2 API (deploymentmanager.googleapis.com)

## Create Ignition configs
The machines will be started manually. Therefore, it is required to generate
the bootstrap and machine Ignition configs and store them for later steps.
Use a [staged install](../overview.md#multiple-invocations) to enable desired customizations.

### Create an install config
Create an install configuration as for [the usual approach](install.md#create-configuration).

If you are installing into a [Shared VPC (XPN)][sharedvpc],
skip this step and create the `install-config.yaml` manually using the documentation references/examples.
The installer will not be able to access the public DNS zone in the host project for the base domain prompt.

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

```sh
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

If you are installing into a [Shared VPC (XPN)][sharedvpc],
`publish` must be set to `Internal`.
The installer will not be able to access the public DNS zone for the the base domain in the host project, which is required for External clusters.
This can be reversed in a step [below](enable-external-ingress-optional).

```sh
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

### Create manifests
Create manifest to enable customizations which are not exposed via the install configuration.

```console
$ openshift-install create manifests
INFO Consuming "Install Config" from target directory
```

### Remove control plane machines
Remove the control plane machines from the manifests.
We'll be providing those ourselves and don't want to involve [the machine-API operator][machine-api-operator].

```sh
rm -f openshift/99_openshift-cluster-api_master-machines-*.yaml
```

### Remove compute machinesets (optional)
If you do not want the cluster to provision compute machines, remove the compute machinesets from the manifests as well.

```sh
rm -f openshift/99_openshift-cluster-api_worker-machineset-*.yaml
```

### Make control-plane nodes unschedulable
Currently [emptying the compute pools](#empty-compute-pools) makes control-plane nodes schedulable.
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
export MASTER_SUBNET_CIDR='10.0.0.0/19'
export WORKER_SUBNET_CIDR='10.0.32.0/19'

export KUBECONFIG=auth/kubeconfig
export CLUSTER_NAME=$(jq -r .clusterName metadata.json)
export INFRA_ID=$(jq -r .infraID metadata.json)
export PROJECT_NAME=$(jq -r .gcp.projectID metadata.json)
export REGION=$(jq -r .gcp.region metadata.json)
export ZONE_0=$(gcloud compute regions describe ${REGION} --format=json | jq -r .zones[0] | cut -d "/" -f9)
export ZONE_1=$(gcloud compute regions describe ${REGION} --format=json | jq -r .zones[1] | cut -d "/" -f9)
export ZONE_2=$(gcloud compute regions describe ${REGION} --format=json | jq -r .zones[2] | cut -d "/" -f9)

export MASTER_IGNITION=$(cat master.ign)
export WORKER_IGNITION=$(cat worker.ign)
```

## Create the VPC
Create the VPC, network, and subnets for the cluster.
This step can be skipped if installing into a pre-existing VPC, such as a [Shared VPC (XPN)][sharedvpc].

Copy [`01_vpc.py`](../../../upi/gcp/01_vpc.py) locally.

Create a resource definition file: `01_vpc.yaml`

```console
$ cat <<EOF >01_vpc.yaml
imports:
- path: 01_vpc.py
resources:
- name: cluster-vpc
  type: 01_vpc.py
  properties:
    infra_id: '${INFRA_ID}'
    region: '${REGION}'
    master_subnet_cidr: '${MASTER_SUBNET_CIDR}'
    worker_subnet_cidr: '${WORKER_SUBNET_CIDR}'
EOF
```
- `infra_id`: the infrastructure name (INFRA_ID above)
- `region`: the region to deploy the cluster into (for example us-east1)
- `master_subnet_cidr`: the CIDR for the master subnet (for example 10.0.0.0/19)
- `worker_subnet_cidr`: the CIDR for the worker subnet (for example 10.0.32.0/19)

Create the deployment using gcloud.

```sh
gcloud deployment-manager deployments create ${INFRA_ID}-vpc --config 01_vpc.yaml
```

## Configure VPC variables
Configure the variables based on the VPC created with `01_vpc.yaml`.
If you are using a pre-existing VPC, such as a [Shared VPC (XPN)][sharedvpc], set these to the `.selfLink` of the targeted resources.

```sh
export CLUSTER_NETWORK=$(gcloud compute networks describe ${INFRA_ID}-network --format json | jq -r .selfLink)
export CONTROL_SUBNET=$(gcloud compute networks subnets describe ${INFRA_ID}-master-subnet --region=${REGION} --format json | jq -r .selfLink)
export COMPUTE_SUBNET=$(gcloud compute networks subnets describe ${INFRA_ID}-worker-subnet --region=${REGION} --format json | jq -r .selfLink)
```

## Create DNS entries and load balancers
Create the DNS zone and load balancers for the cluster.
You can exclude the DNS zone or external load balancer by removing their associated section(s) from the `02_infra.yaml`.
If you choose to exclude the DNS zone, you will need to create it some other way and ensure it is populated with the necessary records as documented below.

If you are installing into a [Shared VPC (XPN)][sharedvpc],
exclude the DNS section as it must be created in the host project.

Copy [`02_dns.py`](../../../upi/gcp/02_dns.py) locally.
Copy [`02_lb_ext.py`](../../../upi/gcp/02_lb_ext.py) locally.
Copy [`02_lb_int.py`](../../../upi/gcp/02_lb_int.py) locally.

Create a resource definition file: `02_infra.yaml`

```console
$ cat <<EOF >02_infra.yaml
imports:
- path: 02_dns.py
- path: 02_lb_ext.py
- path: 02_lb_int.py
resources:
- name: cluster-dns
  type: 02_dns.py
  properties:
    infra_id: '${INFRA_ID}'
    cluster_domain: '${CLUSTER_NAME}.${BASE_DOMAIN}'
    cluster_network: '${CLUSTER_NETWORK}'
- name: cluster-lb-ext
  type: 02_lb_ext.py
  properties:
    infra_id: '${INFRA_ID}'
    region: '${REGION}'
- name: cluster-lb-int
  type: 02_lb_int.py
  properties:
    cluster_network: '${CLUSTER_NETWORK}'
    control_subnet: '${CONTROL_SUBNET}'
    infra_id: '${INFRA_ID}'
    region: '${REGION}'
    zones:
    - '${ZONE_0}'
    - '${ZONE_1}'
    - '${ZONE_2}'
EOF
```
- `infra_id`: the infrastructure name (INFRA_ID above)
- `region`: the region to deploy the cluster into (for example us-east1)
- `cluster_domain`: the domain for the cluster (for example openshift.example.com)
- `cluster_network`: the URI to the cluster network
- `control_subnet`: the URI to the control subnet
- `zones`: the zones to deploy the control plane instances into (for example us-east1-b, us-east1-c, us-east1-d)

Create the deployment using gcloud.

```sh
gcloud deployment-manager deployments create ${INFRA_ID}-infra --config 02_infra.yaml
```

## Configure infra variables
If you excluded the `cluster-lb-ext` section above, then skip `CLUSTER_PUBLIC_IP`.

```sh
export CLUSTER_IP=$(gcloud compute addresses describe ${INFRA_ID}-cluster-ip --region=${REGION} --format json | jq -r .address)
export CLUSTER_PUBLIC_IP=$(gcloud compute addresses describe ${INFRA_ID}-cluster-public-ip --region=${REGION} --format json | jq -r .address)
```

## Add DNS entries
The templates do not create DNS entries due to limitations of Deployment Manager, so we must create them manually.

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
If you deployed external load balancers with `02_infra.yaml`, you can deploy external DNS entries.

```sh
if [ -f transaction.yaml ]; then rm transaction.yaml; fi
gcloud dns record-sets transaction start --zone ${BASE_DOMAIN_ZONE_NAME}
gcloud dns record-sets transaction add ${CLUSTER_PUBLIC_IP} --name api.${CLUSTER_NAME}.${BASE_DOMAIN}. --ttl 60 --type A --zone ${BASE_DOMAIN_ZONE_NAME}
gcloud dns record-sets transaction execute --zone ${BASE_DOMAIN_ZONE_NAME}
```

## Create firewall rules and IAM roles
Create the firewall rules and IAM roles for the cluster.
You can exclude either of these by removing their associated section(s) from the `02_infra.yaml`.
If you choose to do so, you will need to create the required resources some other way.
Details about these resources can be found in the imported python templates.

If you are installing into a [Shared VPC (XPN)][sharedvpc],
exclude the firewall section as they must be created in the host project.

Copy [`03_firewall.py`](../../../upi/gcp/03_firewall.py) locally.
Copy [`03_iam.py`](../../../upi/gcp/03_iam.py) locally.

Create a resource definition file: `03_security.yaml`

```console
$ cat <<EOF >03_security.yaml
imports:
- path: 03_firewall.py
- path: 03_iam.py
resources:
- name: cluster-firewall
  type: 03_firewall.py
  properties:
    allowed_external_cidr: '0.0.0.0/0'
    infra_id: '${INFRA_ID}'
    cluster_network: '${CLUSTER_NETWORK}'
    network_cidr: '${NETWORK_CIDR}'
- name: cluster-iam
  type: 03_iam.py
  properties:
    infra_id: '${INFRA_ID}'
EOF
```
- `allowed_external_cidr`: limits access to the cluster API and ssh to the bootstrap host. (for example External: 0.0.0.0/0, Internal: ${NETWORK_CIDR})
- `infra_id`: the infrastructure name (INFRA_ID above)
- `region`: the region to deploy the cluster into (for example us-east1)
- `cluster_network`: the URI to the cluster network
- `network_cidr`: the CIDR of the vpc network (for example 10.0.0.0/16)

Create the deployment using gcloud.

```sh
gcloud deployment-manager deployments create ${INFRA_ID}-security --config 03_security.yaml
```

## Configure security variables
Configure the variables based on the `03_security.yaml` deployment.
If you excluded the IAM section, ensure these are set to the `.email` of their associated resources.

```sh
export MASTER_SERVICE_ACCOUNT=$(gcloud iam service-accounts list --filter "email~^${INFRA_ID}-m@${PROJECT_NAME}." --format json | jq -r '.[0].email')
export WORKER_SERVICE_ACCOUNT=$(gcloud iam service-accounts list --filter "email~^${INFRA_ID}-w@${PROJECT_NAME}." --format json | jq -r '.[0].email')
```

## Add required roles to IAM service accounts
The templates do not create the policy bindings due to limitations of Deployment Manager, so we must create them manually.

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
export IMAGE_SOURCE=$(curl https://raw.githubusercontent.com/openshift/installer/master/data/data/rhcos.json | jq -r .gcp.url)
gcloud compute images create "${INFRA_ID}-rhcos-image" --source-uri="${IMAGE_SOURCE}"
export CLUSTER_IMAGE=$(gcloud compute images describe ${INFRA_ID}-rhcos-image --format json | jq -r .selfLink)
```

## Upload the bootstrap.ign to a new bucket
Create a bucket and upload the bootstrap.ign file.

```sh
gsutil mb gs://${INFRA_ID}-bootstrap-ignition
gsutil cp bootstrap.ign gs://${INFRA_ID}-bootstrap-ignition/
```

Create a signed URL for the bootstrap instance to use to access the Ignition
config. Export the URL from the output as a variable.

```sh
export BOOTSTRAP_IGN=$(gsutil signurl -d 1h service-account-key.json gs://${INFRA_ID}-bootstrap-ignition/bootstrap.ign | grep "^gs:" | awk '{print $5}')
```

## Launch temporary bootstrap resources

Copy [`04_bootstrap.py`](../../../upi/gcp/04_bootstrap.py) locally.

Create a resource definition file: `04_bootstrap.yaml`

```console
$ cat <<EOF >04_bootstrap.yaml
imports:
- path: 04_bootstrap.py
resources:
- name: cluster-bootstrap
  type: 04_bootstrap.py
  properties:
    infra_id: '${INFRA_ID}'
    region: '${REGION}'
    zone: '${ZONE_0}'
    cluster_network: '${CLUSTER_NETWORK}'
    control_subnet: '${CONTROL_SUBNET}'
    image: '${CLUSTER_IMAGE}'
    machine_type: 'n1-standard-4'
    root_volume_size: '128'
    bootstrap_ign: '${BOOTSTRAP_IGN}'
EOF
```
- `infra_id`: the infrastructure name (INFRA_ID above)
- `region`: the region to deploy the cluster into (for example us-east1)
- `zone`: the zone to deploy the bootstrap instance into (for example us-east1-b)
- `cluster_network`: the URI to the cluster network
- `control_subnet`: the URI to the control subnet
- `image`: the URI to the RHCOS image
- `machine_type`: the machine type of the instance (for example n1-standard-4)
- `bootstrap_ign`: the URL output when creating a signed URL above.

You can add custom tags to `04_bootstrap.py` as needed

```console
            'tags': {
                'items': [
                    context.properties['infra_id'] + '-master',
                    context.properties['infra_id'] + '-bootstrap',
                    'my-custom-tag-example'
                ]
            },
```

Create the deployment using gcloud.

```sh
gcloud deployment-manager deployments create ${INFRA_ID}-bootstrap --config 04_bootstrap.yaml
```

## Add the bootstrap instance to the load balancers
The templates do not manage load balancer membership due to limitations of Deployment
Manager, so we must add the bootstrap node manually.

### Add bootstrap instance to internal load balancer instance group

```sh
gcloud compute instance-groups unmanaged add-instances ${INFRA_ID}-bootstrap-instance-group --zone=${ZONE_0} --instances=${INFRA_ID}-bootstrap
```

### Add bootstrap instance group to the internal load balancer backend service

```sh
gcloud compute backend-services add-backend ${INFRA_ID}-api-internal-backend-service --region=${REGION} --instance-group=${INFRA_ID}-bootstrap-instance-group --instance-group-zone=${ZONE_0}
```

## Launch permanent control plane

Copy [`05_control_plane.py`](../../../upi/gcp/05_control_plane.py) locally.

Create a resource definition file: `05_control_plane.yaml`

```console
$ cat <<EOF >05_control_plane.yaml
imports:
- path: 05_control_plane.py
resources:
- name: cluster-control-plane
  type: 05_control_plane.py
  properties:
    infra_id: '${INFRA_ID}'
    zones:
    - '${ZONE_0}'
    - '${ZONE_1}'
    - '${ZONE_2}'
    control_subnet: '${CONTROL_SUBNET}'
    image: '${CLUSTER_IMAGE}'
    machine_type: 'n1-standard-4'
    root_volume_size: '128'
    service_account_email: '${MASTER_SERVICE_ACCOUNT}'
    ignition: '${MASTER_IGNITION}'
EOF
```
- `infra_id`: the infrastructure name (INFRA_ID above)
- `region`: the region to deploy the cluster into (for example us-east1)
- `zones`: the zones to deploy the control plane instances into (for example us-east1-b, us-east1-c, us-east1-d)
- `control_subnet`: the URI to the control subnet
- `image`: the URI to the RHCOS image
- `machine_type`: the machine type of the instance (for example n1-standard-4)
- `service_account_email`: the email address for the master service account created above
- `ignition`: the contents of the master.ign file

You can add custom tags to `05_control_plane.py` as needed

```console
            'tags': {
                'items': [
                    context.properties['infra_id'] + '-master',
                    'my-custom-tag-example'
                ]
            },
```

Create the deployment using gcloud.

```sh
gcloud deployment-manager deployments create ${INFRA_ID}-control-plane --config 05_control_plane.yaml
```

## Configure control plane variables

```sh
export MASTER0_IP=$(gcloud compute instances describe ${INFRA_ID}-m-0 --zone ${ZONE_0} --format json | jq -r .networkInterfaces[0].networkIP)
export MASTER1_IP=$(gcloud compute instances describe ${INFRA_ID}-m-1 --zone ${ZONE_1} --format json | jq -r .networkInterfaces[0].networkIP)
export MASTER2_IP=$(gcloud compute instances describe ${INFRA_ID}-m-2 --zone ${ZONE_2} --format json | jq -r .networkInterfaces[0].networkIP)
```

## Add DNS entries for control plane etcd
The templates do not manage DNS entries due to limitations of Deployment Manager, so we must add the etcd entries manually.

If you are installing into a [Shared VPC (XPN)][sharedvpc],
use the `--account` and `--project` parameters to perform these actions in the host project.

```sh
if [ -f transaction.yaml ]; then rm transaction.yaml; fi
gcloud dns record-sets transaction start --zone ${INFRA_ID}-private-zone
gcloud dns record-sets transaction add ${MASTER0_IP} --name etcd-0.${CLUSTER_NAME}.${BASE_DOMAIN}. --ttl 60 --type A --zone ${INFRA_ID}-private-zone
gcloud dns record-sets transaction add ${MASTER1_IP} --name etcd-1.${CLUSTER_NAME}.${BASE_DOMAIN}. --ttl 60 --type A --zone ${INFRA_ID}-private-zone
gcloud dns record-sets transaction add ${MASTER2_IP} --name etcd-2.${CLUSTER_NAME}.${BASE_DOMAIN}. --ttl 60 --type A --zone ${INFRA_ID}-private-zone
gcloud dns record-sets transaction add \
  "0 10 2380 etcd-0.${CLUSTER_NAME}.${BASE_DOMAIN}." \
  "0 10 2380 etcd-1.${CLUSTER_NAME}.${BASE_DOMAIN}." \
  "0 10 2380 etcd-2.${CLUSTER_NAME}.${BASE_DOMAIN}." \
  --name _etcd-server-ssl._tcp.${CLUSTER_NAME}.${BASE_DOMAIN}. --ttl 60 --type SRV --zone ${INFRA_ID}-private-zone
gcloud dns record-sets transaction execute --zone ${INFRA_ID}-private-zone
```

## Add control plane instances to load balancers
The templates do not manage load balancer membership due to limitations of Deployment
Manager, so we must add the control plane nodes manually.

### Add control plane instances to internal load balancer instance groups

```sh
gcloud compute instance-groups unmanaged add-instances ${INFRA_ID}-master-${ZONE_0}-instance-group --zone=${ZONE_0} --instances=${INFRA_ID}-m-0
gcloud compute instance-groups unmanaged add-instances ${INFRA_ID}-master-${ZONE_1}-instance-group --zone=${ZONE_1} --instances=${INFRA_ID}-m-1
gcloud compute instance-groups unmanaged add-instances ${INFRA_ID}-master-${ZONE_2}-instance-group --zone=${ZONE_2} --instances=${INFRA_ID}-m-2
```

### Add control plane instances to external load balancer target pools (optional)
If you deployed external load balancers with `02_infra.yaml`, add the control plane instances to the target pool.

```sh
gcloud compute target-pools add-instances ${INFRA_ID}-api-target-pool --instances-zone="${ZONE_0}" --instances=${INFRA_ID}-m-0
gcloud compute target-pools add-instances ${INFRA_ID}-api-target-pool --instances-zone="${ZONE_1}" --instances=${INFRA_ID}-m-1
gcloud compute target-pools add-instances ${INFRA_ID}-api-target-pool --instances-zone="${ZONE_2}" --instances=${INFRA_ID}-m-2
```

## Launch additional compute nodes
You may create compute nodes by launching individual instances discretely
or by automated processes outside the cluster (e.g. Auto Scaling Groups). You
can also take advantage of the built in cluster scaling mechanisms and the
machine API in OpenShift, as mentioned [above](#create-ignition-configs). In
this example, we'll manually launch one instance via the Deployment Manager
template. Additional instances can be launched by including additional
resources of type 06_worker.py in the file.

Copy [`06_worker.py`](../../../upi/gcp/06_worker.py) locally.

Create a resource definition file: `06_worker.yaml`
```console
$ cat <<EOF >06_worker.yaml
imports:
- path: 06_worker.py
resources:
- name: 'w-0'
  type: 06_worker.py
  properties:
    infra_id: '${INFRA_ID}'
    zone: '${ZONE_0}'
    compute_subnet: '${COMPUTE_SUBNET}'
    image: '${CLUSTER_IMAGE}'
    machine_type: 'n1-standard-4'
    root_volume_size: '128'
    service_account_email: '${WORKER_SERVICE_ACCOUNT}'
    ignition: '${WORKER_IGNITION}'
- name: 'w-1'
  type: 06_worker.py
  properties:
    infra_id: '${INFRA_ID}'
    zone: '${ZONE_1}'
    compute_subnet: '${COMPUTE_SUBNET}'
    image: '${CLUSTER_IMAGE}'
    machine_type: 'n1-standard-4'
    root_volume_size: '128'
    service_account_email: '${WORKER_SERVICE_ACCOUNT}'
    ignition: '${WORKER_IGNITION}'
EOF
```
- `name`: the name of the compute node (for example w-a-0)
- `infra_id`: the infrastructure name (INFRA_ID above)
- `region`: the region to deploy the cluster into (for example us-east1)
- `zone`: the zone to deploy the worker node into (for example us-east1-b)
- `compute_subnet`: the URI to the compute subnet
- `image`: the URI to the RHCOS image
- `machine_type`: The machine type of the instance (for example n1-standard-4)
- `service_account_email`: the email address for the worker service account created above
- `ignition`: the contents of the worker.ign file

You can add custom tags to `06_worker.py` as needed

```console
            'tags': {
                'items': [
                    context.properties['infra_id'] + '-worker',
                    'my-custom-tag-example'
                ]
            },
```

Create the deployment using gcloud.

```sh
gcloud deployment-manager deployments create ${INFRA_ID}-worker --config 06_worker.yaml
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
gcloud compute backend-services remove-backend ${INFRA_ID}-api-internal-backend-service --region=${REGION} --instance-group=${INFRA_ID}-bootstrap-instance-group --instance-group-zone=${ZONE_0}
gsutil rm gs://${INFRA_ID}-bootstrap-ignition/bootstrap.ign
gsutil rb gs://${INFRA_ID}-bootstrap-ignition
gcloud deployment-manager deployments delete ${INFRA_ID}-bootstrap
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

export ROUTER_IP=$(oc -n openshift-ingress get service router-default --no-headers | awk '{print $4}')
```

### Add the internal *.apps DNS record

```sh
if [ -f transaction.yaml ]; then rm transaction.yaml; fi
gcloud dns record-sets transaction start --zone ${INFRA_ID}-private-zone
gcloud dns record-sets transaction add ${ROUTER_IP} --name \*.apps.${CLUSTER_NAME}.${BASE_DOMAIN}. --ttl 300 --type A --zone ${INFRA_ID}-private-zone
gcloud dns record-sets transaction execute --zone ${INFRA_ID}-private-zone
```

### Add the external *.apps DNS record (optional)

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
Add a single firewall rule to allow the gce health checks to access all of the services.
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
NAME      VERSION   AVAILABLE   PROGRESSING   SINCE   STATUS
version             False       True          24m     Working towards 4.2.0-0.okd-2019-08-05-204819: 99% complete

$ oc get clusteroperators
NAME                                       VERSION                         AVAILABLE   PROGRESSING   DEGRADED   SINCE
authentication                             4.2.0-0.okd-2019-08-05-204819   True        False         False      6m18s
cloud-credential                           4.2.0-0.okd-2019-08-05-204819   True        False         False      17m
cluster-autoscaler                         4.2.0-0.okd-2019-08-05-204819   True        False         False      80s
console                                    4.2.0-0.okd-2019-08-05-204819   True        False         False      3m57s
dns                                        4.2.0-0.okd-2019-08-05-204819   True        False         False      22m
image-registry                             4.2.0-0.okd-2019-08-05-204819   True        False         False      5m4s
ingress                                    4.2.0-0.okd-2019-08-05-204819   True        False         False      4m38s
insights                                   4.2.0-0.okd-2019-08-05-204819   True        False         False      21m
kube-apiserver                             4.2.0-0.okd-2019-08-05-204819   True        False         False      12m
kube-controller-manager                    4.2.0-0.okd-2019-08-05-204819   True        False         False      12m
kube-scheduler                             4.2.0-0.okd-2019-08-05-204819   True        False         False      11m
machine-api                                4.2.0-0.okd-2019-08-05-204819   True        False         False      18m
machine-config                             4.2.0-0.okd-2019-08-05-204819   True        False         False      22m
marketplace                                4.2.0-0.okd-2019-08-05-204819   True        False         False      5m38s
monitoring                                 4.2.0-0.okd-2019-08-05-204819   True        False         False      86s
network                                    4.2.0-0.okd-2019-08-05-204819   True        False         False      14m
node-tuning                                4.2.0-0.okd-2019-08-05-204819   True        False         False      6m8s
openshift-apiserver                        4.2.0-0.okd-2019-08-05-204819   True        False         False      6m48s
openshift-controller-manager               4.2.0-0.okd-2019-08-05-204819   True        False         False      12m
openshift-samples                          4.2.0-0.okd-2019-08-05-204819   True        False         False      67s
operator-lifecycle-manager                 4.2.0-0.okd-2019-08-05-204819   True        False         False      15m
operator-lifecycle-manager-catalog         4.2.0-0.okd-2019-08-05-204819   True        False         False      15m
operator-lifecycle-manager-packageserver   4.2.0-0.okd-2019-08-05-204819   True        False         False      6m48s
service-ca                                 4.2.0-0.okd-2019-08-05-204819   True        False         False      17m
service-catalog-apiserver                  4.2.0-0.okd-2019-08-05-204819   True        False         False      6m18s
service-catalog-controller-manager         4.2.0-0.okd-2019-08-05-204819   True        False         False      6m19s
storage                                    4.2.0-0.okd-2019-08-05-204819   True        False         False      6m20s

$ oc get pods --all-namespaces
NAMESPACE                                               NAME                                                                READY     STATUS      RESTARTS   AGE
kube-system                                             etcd-member-ip-10-0-3-111.us-east-2.compute.internal                1/1       Running     0          35m
kube-system                                             etcd-member-ip-10-0-3-239.us-east-2.compute.internal                1/1       Running     0          37m
kube-system                                             etcd-member-ip-10-0-3-24.us-east-2.compute.internal                 1/1       Running     0          35m
openshift-apiserver-operator                            openshift-apiserver-operator-6d6674f4f4-h7t2t                       1/1       Running     1          37m
openshift-apiserver                                     apiserver-fm48r                                                     1/1       Running     0          30m
openshift-apiserver                                     apiserver-fxkvv                                                     1/1       Running     0          29m
openshift-apiserver                                     apiserver-q85nm                                                     1/1       Running     0          29m
...
openshift-service-ca-operator                           openshift-service-ca-operator-66ff6dc6cd-9r257                      1/1       Running     0          37m
openshift-service-ca                                    apiservice-cabundle-injector-695b6bcbc-cl5hm                        1/1       Running     0          35m
openshift-service-ca                                    configmap-cabundle-injector-8498544d7-25qn6                         1/1       Running     0          35m
openshift-service-ca                                    service-serving-cert-signer-6445fc9c6-wqdqn                         1/1       Running     0          35m
openshift-service-catalog-apiserver-operator            openshift-service-catalog-apiserver-operator-549f44668b-b5q2w       1/1       Running     0          32m
openshift-service-catalog-controller-manager-operator   openshift-service-catalog-controller-manager-operator-b78cr2lnm     1/1       Running     0          31m
```

[deploymentmanager]: https://cloud.google.com/deployment-manager/docs
[ingress-operator]: https://github.com/openshift/cluster-ingress-operator
[kubernetes-service-load-balancers-exclude-masters]: https://github.com/kubernetes/kubernetes/issues/65618
[machine-api-operator]: https://github.com/openshift/machine-api-operator
[sharedvpc]: https://cloud.google.com/vpc/docs/shared-vpc
