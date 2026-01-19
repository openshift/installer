# Install: User-Provided Infrastructure
The steps for performing a user-provided infrastructure install are outlined here. 
You are also free to create the required resources through other methods; the snippets
(below) are just an example.

_Note_: The example assumes that the `publish` value in the install-config is set to `External`.
There are comments in the snippets for installs where `publish` is set to `Internal`.

## Prerequisites
* all prerequisites from [README](README.md)
* the following binaries installed and in $PATH:
  * gcloud
  * gsutil
* gcloud authenticated to an account with [additional](iam.md) roles:
  * Service Account Key Admin

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
export MASTER_SUBNET_CIDR='10.0.0.0/17'
export WORKER_SUBNET_CIDR='10.0.128.0/17'

export KUBECONFIG=auth/kubeconfig
export CLUSTER_NAME=$(jq -r .clusterName metadata.json)
export INFRA_ID=$(jq -r .infraID metadata.json)
export PROJECT_NAME=$(jq -r .gcp.projectID metadata.json)
export REGION=$(jq -r .gcp.region metadata.json)
export ZONE_0=$(gcloud compute regions describe ${REGION} --format=json | jq -r .zones[0] | cut -d "/" -f9)
export ZONE_1=$(gcloud compute regions describe ${REGION} --format=json | jq -r .zones[1] | cut -d "/" -f9)
export ZONE_2=$(gcloud compute regions describe ${REGION} --format=json | jq -r .zones[2] | cut -d "/" -f9)
export ZONES="${ZONE_0} ${ZONE_1} ${ZONE_2}"

export MASTER_IGNITION=$(cat master.ign)
export WORKER_IGNITION=$(cat worker.ign)
```

## Create the VPC
Create the VPC, network, and subnets for the cluster.
This step can be skipped if installing into a pre-existing VPC, such as a [Shared VPC (XPN)][sharedvpc].

Create a shell script: `01_vpc.sh`, or create a single shell script to contain all snippets in this file.

```console
# create network                                                                                                                                                                                                                                                              
gcloud compute networks create ${INFRA_ID}-network --subnet-mode=custom

# create subnets                                                                                                                                                                                                                                                              
gcloud compute networks subnets create ${INFRA_ID}-master-subnet --network=${INFRA_ID}-network --range=${MASTER_SUBNET_CIDR} --region=${REGION}
gcloud compute networks subnets create ${INFRA_ID}-worker-subnet --network=${INFRA_ID}-network --range=${WORKER_SUBNET_CIDR} --region=${REGION}

# create router                                                                                                                                                                                                                                                               
gcloud compute routers create ${INFRA_ID}-router --network=${INFRA_ID}-network --region=${REGION}

# create nats                                                                                                                                                                                                                                                                 
gcloud compute routers nats create ${INFRA_ID}-nat-master --router=${INFRA_ID}-router --auto-allocate-nat-external-ips --nat-custom-subnet-ip-ranges=${INFRA_ID}-master-subnet --region=${REGION}
gcloud compute routers nats create ${INFRA_ID}-nat-worker --router=${INFRA_ID}-router --auto-allocate-nat-external-ips --nat-custom-subnet-ip-ranges=${INFRA_ID}-worker-subnet --region=${REGION}

# [optional] These lines will streamline the deprovision process                                                                                                                                                                                                              
cat << EOF >> "${deprovision_commands_file}"                                                                                                                                                                                                                                  
gcloud compute routers nats delete -q ${INFRA_ID}-nat-master --router ${INFRA_ID}-router --region ${REGION}                                                                                                                                                                   
gcloud compute routers nats delete -q ${INFRA_ID}-nat-worker --router ${INFRA_ID}-router --region ${REGION}                                                                                                                                                                   
gcloud compute routers delete -q ${INFRA_ID}-router --region ${region}                                                                                                                                                                                                        
gcloud compute networks subnets delete -q ${INFRA_ID}-master-subnet --region ${REGION}                                                                                                                                                                                        
gcloud compute networks subnets delete -q ${INFRA_ID}-worker-subnet --region ${REGION}                                                                                                                                                                                        
gcloud compute networks delete -q ${INFRA_ID}-network                                                                                                                                                                                                                         
EOF
```

- `deprovision_commands_file`: the filename (including the path) where the deprovision commands will be added
- `INFRA_ID`: the name of the cluster followed by a unique hash
- `MASTER_SUBNET_CIDR`: the CIDR for the master subnet (for example 10.0.0.0/17)
- `WORKER_SUBNET_CIDR`:the CIDR for the worker subnet (for example 10.0.128.0/17)
- `REGION`: the region to deploy the cluster into (for example us-east1)

Execute the commands in the file above to create the network resources.

```sh
./01_vpc.sh
```

## Configure VPC variables
Configure the variables based on the VPC created with `01_vpc.sh`.
If you are using a pre-existing VPC, such as a [Shared VPC (XPN)][sharedvpc], set these to the `.selfLink` of the targeted resources.

```sh
export CLUSTER_NETWORK=$(gcloud compute networks describe ${INFRA_ID}-network --format json | jq -r .selfLink)
export CONTROL_SUBNET=$(gcloud compute networks subnets describe ${INFRA_ID}-master-subnet --region=${REGION} --format json | jq -r .selfLink)
export COMPUTE_SUBNET=$(gcloud compute networks subnets describe ${INFRA_ID}-worker-subnet --region=${REGION} --format json | jq -r .selfLink)
```

## Create DNS zones and load balancers
Create the DNS zone and load balancers for the cluster.
You can exclude the DNS zone or external load balancer by removing their associated section(s) from the `02_infra.sh`.
If you choose to exclude the DNS zone, you will need to create it some other way and ensure it is populated with the necessary records as documented below.

If you are installing into a [Shared VPC (XPN)][sharedvpc],
exclude the DNS section as it must be created in the host project.

Create a shell script: `02_infra.sh`, or add to the original.

```console
# ---------------------------------------------------------
# Create the Private DNS Zone
# ---------------------------------------------------------
gcloud dns managed-zones create ${INFRA_ID}-private-zone --description="Pre-created private DNS zone" --dns-name=${CLUSTER_NAME}.${BASE_DOMAIN}. --visibility="private" --networks=${CLUSTER_NETWORK}

# ---------------------------------------------------------
# Create the resources for the external load balancer
# Note: This can be skipped for internal installs
# ---------------------------------------------------------
gcloud compute addresses create ${INFRA_ID}-cluster-public-ip --region=${REGION}
sleep 3
address_selflink=$(gcloud compute addresses describe ${INFRA_ID}-cluster-public-ip --region=${REGION} --format=json | jq -r .selfLink)

gcloud compute http-health-checks create ${INFRA_ID}-api-http-health-check --port=6080 --request-path="/readyz"
sleep 3
hc_selflink=$(gcloud compute http-health-checks describe ${INFRA_ID}-api-http-health-check --format=json | jq -r .selfLink)

gcloud compute target-pools create ${INFRA_ID}-api-target-pool --http-health-check=${hc_selflink} --region=${REGION}
sleep 3
tp_selflink=$(gcloud compute target-pools describe ${INFRA_ID}-api-target-pool --region=${REGION} --format=json | jq -r .selfLink)

gcloud compute forwarding-rules create ${INFRA_ID}-api-forwarding-rule --region=${REGION} --address=${address_selflink} --target-pool=${tp_selflink} --ports=6443
sleep 3

# [optional] These lines will streamline the deprovision process 
cat << EOF >> "${deprovision_commands_file}"
gcloud compute forwarding-rules delete -q ${INFRA_ID}-api-forwarding-rule --region=${REGION}
gcloud compute target-pools delete -q ${INFRA_ID}-api-target-pool --region=${REGION}
gcloud compute http-health-checks delete -q ${INFRA_ID}-api-http-health-check
gcloud compute addresses delete -q ${INFRA_ID}-cluster-public-ip --region=${REGION}
EOF

# ---------------------------------------------------------
# Create the resources for the internal load balancer
# ---------------------------------------------------------
gcloud compute addresses create ${INFRA_ID}-cluster-ip --region=${REGION} --subnet=${CONTROL_SUBNET}
sleep 3
int_address_selflink=$(gcloud compute addresses describe ${INFRA_ID}-cluster-ip --region=${REGION} --format=json | jq -r .selfLink)

gcloud compute health-checks create https ${INFRA_ID}-api-internal-health-check --port=6443 --request-path="/readyz"
sleep 3
int_hc_selflink=$(gcloud compute health-checks describe ${INFRA_ID}-api-internal-health-check --format=json | jq -r .selfLink)

gcloud compute backend-services create ${INFRA_ID}-api-internal --region=${REGION} --protocol=TCP --load-balancing-scheme=INTERNAL --health-checks=${int_hc_selflink} --timeout=120
sleep 3
int_bs_selflink=$(gcloud compute backend-services describe ${INFRA_ID}-api-internal --region=${REGION} --format=json | jq -r .selfLink)

ZONES_ARRAY=($ZONES)
for zone in "${ZONES_ARRAY[@]}"; do
    gcloud compute instance-groups unmanaged create ${INFRA_ID}-master-${zone}-ig --zone=${zone}
    sleep 3
    gcloud compute instance-groups unmanaged set-named-ports ${INFRA_ID}-master-${zone}-ig --zone=${zone} --named-ports=ignition:22623,https:6443
    sleep 3
done

gcloud compute forwarding-rules create ${INFRA_ID}-api-internal-forwarding-rule --region=${REGION} --load-balancing-scheme=INTERNAL --ports=6443,22623 --backend-service=${int_bs_selflink} --address=${int_address_selflink} --subnet=${CONTROL_SUBNET}

# [optional] These lines will streamline the deprovision process 
cat << EOF >> "${deprovision_commands_file}"
gcloud compute forwarding-rules delete -q ${infra_id}-api-internal-forwarding-rule --region=${region}
gcloud compute backend-services delete -q ${infra_id}-api-internal --region=${region}
gcloud compute health-checks delete -q https ${infra_id}-api-internal-health-check
gcloud compute addresses delete -q ${infra_id}-cluster-ip --region=${region}
EOF
# [optional] These lines will streamline the deprovision process 
for zone in "${ZONES_ARRAY[@]}"; do
  cat << EOF >> "${deprovision_commands_file}"
gcloud compute instance-groups unmanaged delete -q ${infra_id}-master-${zone}-ig --zone=${zone}
EOF
done
```

- `INFRA_ID`: the infrastructure name (INFRA_ID above)
- `REGION`: the region to deploy the cluster into (for example us-east1)
- `CLUSTER_NAME`: Name of the cluster (from the cluster metadata and install-config)
- `BASE_DOMAIN`: Base domain from the install-config
- `CLUSTER_NETWORK`: the URI to the cluster network
- `CONTROL_SUBNET`: the URI to the control subnet
- `ZONES`: the zones to deploy (consists of all zones declared above ex: $ZONE_0 ...)

## Configure infra variables

```sh
export CLUSTER_IP=$(gcloud compute addresses describe ${INFRA_ID}-cluster-ip --region=${REGION} --format json | jq -r .address)
export CLUSTER_PUBLIC_IP=$(gcloud compute addresses describe ${INFRA_ID}-cluster-public-ip --region=${REGION} --format json | jq -r .address)

export API_INTERNAL_BACKEND_SVC=$(gcloud compute backend-services list --filter="name~${INFRA_ID}-api-internal" --format='value(name)')
```

## Add DNS entries
Create the DNS entries.

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
If you deployed external load balancers with `02_infra.sh`, you can deploy external DNS entries.

```sh
if [ -f transaction.yaml ]; then rm transaction.yaml; fi
gcloud dns record-sets transaction start --zone ${BASE_DOMAIN_ZONE_NAME}
gcloud dns record-sets transaction add ${CLUSTER_PUBLIC_IP} --name api.${CLUSTER_NAME}.${BASE_DOMAIN}. --ttl 60 --type A --zone ${BASE_DOMAIN_ZONE_NAME}
gcloud dns record-sets transaction execute --zone ${BASE_DOMAIN_ZONE_NAME}
```

## Create firewall rules and IAM roles
Create the firewall rules and IAM roles for the cluster.
You can exclude either of these by removing their associated section(s) from the `02_infra.sh`.
If you choose to do so, you will need to create the required resources some other way.

If you are installing into a [Shared VPC (XPN)][sharedvpc],
exclude the firewall section as they must be created in the host project.

Create a shell script: `03_security.sh`, or add to the original.

```console
# ---------------------------------------------------------
# Create the service accounts
# ---------------------------------------------------------
gcloud iam service-accounts create ${INFRA_ID}-m --display-name=${INFRA_ID}-master-node
sleep 3

gcloud iam service-accounts create ${INFRA_ID}-w --display-name=${INFRA_ID}-worker-node
sleep 3

# ---------------------------------------------------------
# Create the firewall rules
# ---------------------------------------------------------
allowed_external_cidr="0.0.0.0/0"

gcloud compute firewall-rules create ${INFRA_ID}-bootstrap-in-ssh --network=${CLUSTER_NETWORK} --allow=tcp:22 --source-ranges=${allowed_external_cidr} --target-tags=${INFRA_ID}-bootstrap
gcloud compute firewall-rules create ${INFRA_ID}-api --network=${CLUSTER_NETWORK} --allow=tcp:6443 --source-ranges=${allowed_external_cidr} --target-tags=${INFRA_ID}-master
gcloud compute firewall-rules create ${INFRA_ID}-health-checks --network=${CLUSTER_NETWORK} --allow=tcp:6080,tcp:6443,tcp:22624 --source-ranges=35.191.0.0/16,130.211.0.0/22,209.85.152.0/22,209.85.204.0/22 --target-tags=${INFRA_ID}-master
gcloud compute firewall-rules create ${INFRA_ID}-etcd --network=${CLUSTER_NETWORK} --allow=tcp:2379-2380 --source-tags=${INFRA_ID}-master --target-tags=${INFRA_ID}-master
gcloud compute firewall-rules create ${INFRA_ID}-control-plane --network=${CLUSTER_NETWORK} --allow=tcp:10257,tcp:10259,tcp:22623 --source-tags=${INFRA_ID}-master,${INFRA_ID}-worker --target-tags=${INFRA_ID}-master
gcloud compute firewall-rules create ${INFRA_ID}-internal-network --network=${CLUSTER_NETWORK} --allow=icmp,tcp:22 --source-ranges=${NETWORK_CIDR} --target-tags=${INFRA_ID}-master,${INFRA_ID}-worker
gcloud compute firewall-rules create ${INFRA_ID}-internal-cluster --network=${CLUSTER_NETWORK} --allow=udp:4789,udp:6081,udp:500,udp:4500,esp,tcp:9000-9999,udp:9000-9999,tcp:10250,tcp:30000-32767,udp:30000-32767 --source-tags=${INFRA_ID}-master,${INFRA_ID}-worker --target-tags=${INFRA_ID}-master,${INFRA_ID}-worker

# [optional] These lines will streamline the deprovision process 
cat << EOF >> "${deprovision_commands_file}"
gcloud compute firewall-rules delete -q ${INFRA_ID}-bootstrap-in-ssh ${INFRA_ID}-api ${INFRA_ID}-health-checks ${INFRA_ID}-etcd ${INFRA_ID}-control-plane ${INFRA_ID}-internal-network ${INFRA_ID}-internal-cluster
EOF
```
- `allowed_external_cidr`: limits access to the cluster API and ssh to the bootstrap host. (for example External: 0.0.0.0/0, Internal: ${NETWORK_CIDR})
- `INFRA_ID`: the infrastructure name (INFRA_ID above)
- `REGION`: the region to deploy the cluster into (for example us-east1)
- `CLUSTER_NETWORK`: the URI to the cluster network

## Configure security variables
If you excluded the IAM section, ensure these are set to the `.email` of their associated resources.

```sh
export MASTER_SERVICE_ACCOUNT=$(gcloud iam service-accounts list --filter "email~^${INFRA_ID}-m@${PROJECT_NAME}." --format json | jq -r '.[0].email')
export WORKER_SERVICE_ACCOUNT=$(gcloud iam service-accounts list --filter "email~^${INFRA_ID}-w@${PROJECT_NAME}." --format json | jq -r '.[0].email')
```

## Add required roles to IAM service accounts
Create the policy bindings.

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

Adding policy bindings can take time to apply. It may be helpful to include a backoff function similar to the one below:

```sh
function backoff() {
    local attempt=0
    local failed=0
    while true; do
        eval "$*" && failed=0 || failed=1
        if [[ $failed -eq 0 ]]; then
            break
        fi
        attempt=$(( attempt + 1 ))
        if [[ $attempt -gt 5 ]]; then
            break
        fi
        echo "command failed, retrying in $(( 2 ** attempt )) seconds"
        sleep $(( 2 ** attempt ))
    done
    return $failed
}
```

Pass the gcloud function calls above to the backoff function:

```sh
backoff gcloud projects add-iam-policy-binding "${PROJECT_NAME}" --member "serviceAccount:${MASTER_SERVICE_ACCOUNT}" --role "roles/compute.instanceAdmin" 1> /dev/null
```

## Generate a service-account-key for signing the bootstrap.ign url

```sh
gcloud iam service-accounts keys create service-account-key.json --iam-account=${MASTER_SERVICE_ACCOUNT}
```

## Create the cluster image.
Locate the RHCOS image source and create a cluster image.

```sh
export IMAGE_SOURCE=$(curl https://raw.githubusercontent.com/openshift/installer/master/data/data/coreos/rhcos.json | jq -r '.architectures.x86_64.images.gcp')
export IMAGE_NAME=$(echo "${IMAGE_SOURCE}" | jq -r '.name')
export IMAGE_PROJECT=$(echo "${IMAGE_SOURCE}" | jq -r '.project')
export CLUSTER_IMAGE=$(gcloud compute images describe ${IMAGE_NAME} --project ${IMAGE_PROJECT} --format json | jq -r .selfLink)
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

Create a shell script: `04_bootstrap.sh`, or add to the original.

```console
# ---------------------------------------------------------
# Launch temporary bootstrap resources
# ---------------------------------------------------------
# [optional] for external installs only
gcloud compute addresses create ${INFRA_ID}-bootstrap-public-ip --region=${REGION}
sleep 3
public_ip=$(gcloud compute addresses describe ${INFRA_ID}-bootstrap-public-ip --region=${REGION} --format=json | jq -r .address)

# The following section is for external installs
# If this is an internal install remove the --address parameter and add --no-address
gcloud compute instances create ${INFRA_ID}-bootstrap --boot-disk-size=128GB --image=${CLUSTER_IMAGE} --metadata=^#^user-data='{\"ignition\":{\"config\":{\"replace\":{\"source\":\"${ignition}\"}},\"version\":\"3.2.0\"}}' --machine-type=${MACHINE_TYPE} --zone=${ZONE_0} --tags=${INFRA_ID}-master,${INFRA_ID}-bootstrap --subnet=${CONTROL_SUBNET} --address=${public_ip}
sleep 3

gcloud compute instance-groups unmanaged create ${INFRA_ID}-bootstrap-ig --zone=${ZONE_0}
sleep 3
gcloud compute instance-groups unmanaged set-named-ports ${INFRA_ID}-bootstrap-ig --zone=${ZONE_0} --named-ports=ignition:22623,https:6443
sleep 3

cat << EOF >> "${deprovision_commands_file}"
gcloud compute instance-groups unmanaged delete -q ${INFRA_ID}-bootstrap-ig --zone=${ZONE_0}
gcloud compute instances delete -q ${INFRA_ID}-bootstrap --zone=${ZONE_0}
EOF

# Comment out during internal installs
cat << EOF >> "${deprovision_commands_file}"
gcloud compute addresses delete -q ${INFRA_ID}-bootstrap-public-ip --region=${REGION}
EOF

BOOTSTRAP_INSTANCE_GROUP=$(gcloud compute instance-groups list --filter="name~^${INFRA_ID}-bootstrap-" --format "value(name)")

# Add bootstrap instance to internal load balancer instance group
gcloud compute instance-groups unmanaged add-instances ${BOOTSTRAP_INSTANCE_GROUP} --zone=${ZONE_0} --instances=${INFRA_ID}-bootstrap

# Add bootstrap instance group to the internal load balancer backend service
gcloud compute backend-services add-backend ${API_INTERNAL_BACKEND_SVC} --region=${REGION} --instance-group=${BOOTSTRAP_INSTANCE_GROUP} --instance-group-zone=${ZONE_0}

# The following 2 lines are only intended for external installs.
gcloud compute target-pools add-instances ${INFRA_ID}-api-target-pool --instances-zone=${ZONE_0} --instances=${INFRA_ID}-bootstrap
BOOTSTRAP_IP="$(gcloud compute addresses describe --region "${REGION}" "${INFRA_ID}-bootstrap-public-ip" --format json | jq -r .address)"
# The following line is used for internal installs only
# BOOTSTRAP_IP="$(gcloud compute instances describe --zone "${ZONE_0}" "${INFRA_ID}-bootstrap" --format=json | jq -r .networkInterfaces[].networkIP)"
```
- `INFRA_ID`: the infrastructure name (INFRA_ID above)
- `REGION`: the region to deploy the cluster into (for example us-east1)
- `ZONE_0`: the zone to deploy the bootstrap instance into
- `CLUSTER_NETWORK`: the URI to the cluster network
- `CONTROL_SUBNET`: the URI to the control subnet
- `CLUSTER_IMAGE`: the URI to the RHCOS image
- `MACHINE_TYPE`: the machine type of the instance (for example n1-standard-4)
- `BOOTSTRAP_IGN`: the URL output when creating a signed URL above.

## Launch permanent control plane

Create a shell script: `05_control_plane.sh`, or add to the original.

```console
machine_role="master"
index=0

ZONES_ARRAY=($ZONES)
for zone in "${ZONES_ARRAY[@]}"; do
    gcloud compute instances create ${INFRA_ID}-${machine_role}-${index} --boot-disk-size=${DISK_SIZE}GB --boot-disk-type=pd-ssd --image=${CLUSTER_IMAGE} --metadata=^#^user-data="${MASTER_IGNITION}" --machine-type=${MACHINE_TYPE} --zone=${zone} --no-address --service-account=${MASTER_SERVICE_ACCOUNT} --scopes=https://www.googleapis.com/auth/cloud-platform --tags=${INFRA_ID}-${machine_role} --subnet=${CONTROL_SUBNET}
    sleep 3
    
    cat <<EOF >> "${deprovision_commands_file}"
gcloud compute instances delete -q ${INFRA_ID}-${machine_role}-${index} --zone=${zone}
EOF
  index=$(( $index + 1))
done

# configure control plane variables
MASTER0_IP="$(gcloud compute instances describe "${INFRA_ID}-master-0" --zone "${ZONE_0}" --format json | jq -r .networkInterfaces[0].networkIP)"
MASTER1_IP="$(gcloud compute instances describe "${INFRA_ID}-master-1" --zone "${ZONE_1}" --format json | jq -r .networkInterfaces[0].networkIP)"
MASTER2_IP="$(gcloud compute instances describe "${INFRA_ID}-master-2" --zone "${ZONE_2}" --format json | jq -r .networkInterfaces[0].networkIP)"

GATHER_BOOTSTRAP_ARGS+=('--master' "${MASTER0_IP}" '--master' "${MASTER1_IP}" '--master' "${MASTER2_IP}")

# Add DNS entries for control plane etcd
export PRIVATE_ZONE_NAME="${INFRA_ID}-private-zone"
gcloud dns record-sets transaction start --zone "${PRIVATE_ZONE_NAME}"
gcloud dns record-sets transaction add "${MASTER0_IP}" --name "etcd-0.${CLUSTER_NAME}.${BASE_DOMAIN}." --ttl 60 --type A --zone "${PRIVATE_ZONE_NAME}"
gcloud dns record-sets transaction add "${MASTER1_IP}" --name "etcd-1.${CLUSTER_NAME}.${BASE_DOMAIN}." --ttl 60 --type A --zone "${PRIVATE_ZONE_NAME}"
gcloud dns record-sets transaction add "${MASTER2_IP}" --name "etcd-2.${CLUSTER_NAME}.${BASE_DOMAIN}." --ttl 60 --type A --zone "${PRIVATE_ZONE_NAME}"
gcloud dns record-sets transaction add \
  "0 10 2380 etcd-0.${CLUSTER_NAME}.${BASE_DOMAIN}." \
  "0 10 2380 etcd-1.${CLUSTER_NAME}.${BASE_DOMAIN}." \
  "0 10 2380 etcd-2.${CLUSTER_NAME}.${BASE_DOMAIN}." \
  --name "_etcd-server-ssl._tcp.${CLUSTER_NAME}.${BASE_DOMAIN}." --ttl 60 --type SRV --zone "${PRIVATE_ZONE_NAME}"
gcloud dns record-sets transaction execute --zone "${PRIVATE_ZONE_NAME}"

MASTER_IG_0="$(gcloud compute instance-groups list --filter "name~^${INFRA_ID}-master-${ZONE_0}-" --format "value(name)")"
MASTER_IG_1="$(gcloud compute instance-groups list --filter "name~^${INFRA_ID}-master-${ZONE_1}-" --format "value(name)")"
MASTER_IG_2="$(gcloud compute instance-groups list --filter "name~^${INFRA_ID}-master-${ZONE_2}-" --format "value(name)")"

# Add control plane instances to instance groups
gcloud compute instance-groups unmanaged add-instances ${MASTER_IG_0} --zone=${ZONE_0} --instances=${INFRA_ID}-master-0
gcloud compute instance-groups unmanaged add-instances ${MASTER_IG_1} --zone=${ZONE_1} --instances=${INFRA_ID}-master-1
gcloud compute instance-groups unmanaged add-instances ${MASTER_IG_2} --zone=${ZONE_2} --instances=${INFRA_ID}-master-2

# Add control plane instances to internal load balancer backend-service
gcloud compute backend-services add-backend ${API_INTERNAL_BACKEND_SVC} --region=${REGION} --instance-group=${MASTER_IG_0} --instance-group-zone=${ZONE_0}
gcloud compute backend-services add-backend ${API_INTERNAL_BACKEND_SVC} --region=${REGION} --instance-group=${MASTER_IG_1} --instance-group-zone=${ZONE_1}
gcloud compute backend-services add-backend ${API_INTERNAL_BACKEND_SVC} --region=${REGION} --instance-group=${MASTER_IG_2} --instance-group-zone=${ZONE_2}
```
- `INFRA_ID`: the infrastructure name (INFRA_ID above)
- `REGION`: the region to deploy the cluster into (for example us-east1)
- `ZONES`: the zones to deploy the control plane instances into (for example us-east1-b, us-east1-c, us-east1-d)
- `CONTROL_SUBNET`: the URI to the control subnet
- `CLUSTER_IMAGE`: the URI to the RHCOS image
- `MACHINE_TYPE`: the machine type of the instance (for example n1-standard-4)
- `DISK_SIZE`: The machine boot disk size (in GB).
- `MASTER_SERVICE_ACCOUNT`: the email address for the master service account created above
- `ignition`: the contents of the master.ign file


### Add control plane instances to external load balancer target pools (optional)
If you deployed external load balancers with `02_infra.sh`, add the control plane instances to the target pool.

```sh
gcloud compute target-pools add-instances ${INFRA_ID}-api-target-pool --instances-zone="${ZONE_0}" --instances=${INFRA_ID}-master-0
gcloud compute target-pools add-instances ${INFRA_ID}-api-target-pool --instances-zone="${ZONE_1}" --instances=${INFRA_ID}-master-1
gcloud compute target-pools add-instances ${INFRA_ID}-api-target-pool --instances-zone="${ZONE_2}" --instances=${INFRA_ID}-master-2
```

## Launch additional compute nodes
You may create compute nodes by launching individual instances discretely
or by automated processes outside the cluster (e.g. Auto Scaling Groups). You
can also take advantage of the built in cluster scaling mechanisms and the
machine API in OpenShift, as mentioned [above](#create-ignition-configs). In
this example, we'll manually launch one instance.

Create a resource definition file: `06_worker.sh`
```console
# Available zones and instance zones might be different in region for arm64 machines

if [[ ${BASH_VERSION%%.*} -lt 4 ]]; then
  zones=`(gcloud compute regions describe "${REGION}" --format=json | jq -r '.zones[]' | cut -d "/" -f9)`
  AVAILABILITY_ZONES=( $zones )

  worker_instance_zones=`gcloud compute machine-types list --filter="zone:(${REGION}) AND name=(${MACHINE_TYPE})" --format=json | jq -r '.[].zone'`
  WORKER_INSTANCE_ZONES=( $worker_instance_zones )

  worker_zones=`echo "${AVAILABILITY_ZONES[@]}" "${WORKER_INSTANCE_ZONES[@]}" | sed 's/ /\n/g' | sort -R | uniq -d`
  WORKER_ZONES=( $worker_zones )
else
  mapfile -t AVAILABILITY_ZONES < <(gcloud compute regions describe "${REGION}" --format=json | jq -r '.zones[]' | cut -d "/" -f9)
  mapfile -t WORKER_INSTANCE_ZONES < <(gcloud compute machine-types list --filter="zone:(${REGION}) AND name=(${MACHINE_TYPE})" --format=json | jq -r '.[].zone')
  mapfile -t WORKER_ZONES < <(echo "${AVAILABILITY_ZONES[@]}" "${WORKER_INSTANCE_ZONES[@]}" | sed 's/ /\n/g' | sort -R | uniq -d)
fi

machine_role="worker"
index=0

for zone in "${WORKER_ZONES[@]}"; do
    gcloud compute instances create ${INFRA_ID}-${machine_role}-${index} --boot-disk-size=${DISK_SIZE}GB --boot-disk-type=pd-ssd --image=${CLUSTER_IMAGE} --metadata=^#^user-data="${WORKER_IGNITION}" --machine-type=${MACHINE_TYPE} --zone=${zone} --no-address --service-account=${WORKER_SERVICE_ACCOUNT} --scopes=https://www.googleapis.com/auth/cloud-platform --tags=${INFRA_ID}-${machine_role} --subnet=${COMPUTE_SUBNET}
    sleep 3

    cat << EOF >> "${deprovision_commands_file}"
gcloud compute instances delete -q ${INFRA_ID}-${machine_role}-${index} --zone=${zone}
EOF

    index=$(( $index + 1))
done
```
- `INFRA_ID`: the infrastructure name (INFRA_ID above)
- `COMPUTE_SUBNET`: the URI to the compute subnet
- `CLUSTER_IMAGE`: the URI to the RHCOS image
- `MACHINE_TYPE`: The machine type of the instance (for example n1-standard-4)
- `DISK_SIZE`: The machine boot disk size (in GB).
- `WORKER_SERVICE_ACCOUNT`: the email address for the worker service account created above
- `ignition`: the contents of the worker.ign file
- 
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
gcloud compute backend-services remove-backend ${INFRA_ID}-api-internal --region=${REGION} --instance-group=${INFRA_ID}-bootstrap-instance-group --instance-group-zone=${ZONE_0}
gsutil rm gs://${INFRA_ID}-bootstrap-ignition/bootstrap.ign
gsutil rb gs://${INFRA_ID}-bootstrap-ignition
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

[ingress-operator]: https://github.com/openshift/cluster-ingress-operator
[kubernetes-service-load-balancers-exclude-masters]: https://github.com/kubernetes/kubernetes/issues/65618
[machine-api-operator]: https://github.com/openshift/machine-api-operator
[sharedvpc]: https://cloud.google.com/vpc/docs/shared-vpc
