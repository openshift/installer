# OpenStack Platform Support

This document discusses the requirements, current expected behavior, and how to try out what exists so far.
In addition, it covers the installation with the default CNI (OpenShiftSDN), as well as with the Kuryr SDN.

## Table of Contents

- [OpenStack Platform Support](#openstack-platform-support)
  - [Table of Contents](#table-of-contents)
  - [Reference Documents](#reference-documents)
  - [OpenStack Requirements](#openstack-requirements)
    - [Master Nodes](#master-nodes)
    - [Worker Nodes](#worker-nodes)
    - [Bootstrap Node](#bootstrap-node)
    - [Image Registry Requirements](#image-registry-requirements)
    - [Disk Requirements](#disk-requirements)
    - [Neutron Public Network](#neutron-public-network)
    - [Nova Metadata Service](#nova-metadata-service)
    - [Glance Service](#glance-service)
  - [OpenStack Credentials](#openstack-credentials)
    - [Self Signed OpenStack CA certificates](#self-signed-openstack-ca-certificates)
  - [Standalone Single-Node Development Environment](#standalone-single-node-development-environment)
  - [Running The Installer](#running-the-installer)
    - [Known Issues](#known-issues)
    - [Initial Setup](#initial-setup)
    - [API Access](#api-access)
      - [Using Floating IPs](#using-floating-ips)
        - [Create API and Ingress Floating IP Addresses](#create-api-and-ingress-floating-ip-addresses)
        - [Create API and Ingress DNS Records](#create-api-and-ingress-dns-records)
        - [External API Access](#external-api-access)
        - [External Ingress (apps) Access](#external-ingress-apps-access)
      - [Without Floating IPs](#without-floating-ips)
    - [Running a Deployment](#running-a-deployment)
    - [Current Expected Behavior](#current-expected-behavior)
    - [Checking Cluster Status](#checking-cluster-status)
    - [Destroying The Cluster](#destroying-the-cluster)
  - [Post Install Operations](#post-install-operations)
    - [Adding a MachineSet](#adding-a-machineset)
      - [Using a Server Group](#using-a-server-group)
      - [Setting Nova Availability Zones](#setting-nova-availability-zones)
    - [Using a Custom External Load Balancer](#using-a-custom-external-load-balancer)
    - [Reconfiguring cloud provider](#reconfiguring-cloud-provider)
      - [Modifying cloud provider options](#modifying-cloud-provider-options)
      - [Refreshing a CA Certificate](#refreshing-a-ca-certificate)
  - [Reporting Issues](#reporting-issues)

## Reference Documents

- [Privileges](privileges.md)
- [Known Issues and Workarounds](known-issues.md)
- [Using the OSP 4 installer with Kuryr](kuryr.md)
- [Troubleshooting your cluster](troubleshooting.md)
- [Customizing your install](customization.md)
- [Installing OpenShift on OpenStack User-Provisioned Infrastructure](install_upi.md)
- [Learn about the OpenShift on OpenStack networking infrastructure design](../../design/openstack/networking-infrastructure.md)
- [Deploying OpenShift bare-metal workers](deploy_baremetal_workers.md)
- [Deploying OpenShift single root I/O virtualization (SRIOV) workers](deploy_sriov_workers.md)
- [Provider Networks](provider_networks.md)
- [How to set affinity rules for workers at install-time](affinity.md)

## OpenStack Requirements

In order to run the latest version of the installer in OpenStack, at a bare minimum you need the following quota to run a *default* cluster. While it is possible to run the cluster with fewer resources than this, it is not recommended. Certain cases, such as deploying [without FIPs](#without-floating-ips), or deploying with an [external load balancer](#using-an-external-load-balancer) are documented below, and are not included in the scope of this recommendation. If you are planning on using Kuryr, or want to learn more about it, please read through the [Kuryr documentation](kuryr.md). **NOTE: The installer has been tested and developed on Red Hat OSP 13.**

For a successful installation it is required:

- Floating IPs: 2 (plus one that will be created and destroyed by the Installer during the installation process)
- Security Groups: 3
- Security Group Rules: 60
- Routers: 1
- Subnets: 1
- Server Groups: 1
- RAM: 112 GB
- vCPUs: 28
- Volume Storage: 175 GB
- Instances: 7
- Depending on the type of [image registry backend](#image-registry-requirements) either 1 Swift container or an additional 100 GB volume.
- OpenStack resource tagging

**NOTE:** The installer will check OpenStack quota limits to make sure that the requested resources can be created. Note that it won't check for resource availability in the cloud, but only on the quotas.

You may need to increase the security group related quotas from their default values. For example (as an OpenStack administrator):

```sh
openstack quota set --secgroups 8 --secgroup-rules 100 <project>`
```

Once you configure the quota for your project, please ensure that the user for the installer has the proper [privileges](privileges.md).

### Master Nodes

The default deployment stands up 3 master nodes, which is the minimum amount required for a cluster. For each master node you stand up, you will need 1 instance, and 1 port available in your quota. They should be assigned a flavor with at least 16 GB RAM, 4 vCPUs, and 25 GB Disk (or Root Volume). It is theoretically possible to run with a smaller flavor, but be aware that if it takes too long to stand up services, or certain essential services crash, the installer could time out, leading to a failed install.

The Master Nodes are placed in a single Server Group with "soft anti-affinity" policy; the machines will therefore be creted on separate hosts when possible.

### Worker Nodes

The default deployment stands up 3 worker nodes. Worker nodes host the applications you run on OpenShift. The flavor assigned to the worker nodes should have at least 2 vCPUs, 8 GB RAM and 25 GB Disk (or Root Volume). However, if you are experiencing `Out Of Memory` issues, or your installs are timing out, try increasing the size of your flavor to match the master nodes: 4 vCPUs and 16 GB RAM.

See the [OpenShift documentation](https://docs.openshift.com/container-platform/4.4/architecture/control-plane.html#defining-workers_control-plane) for more information on the worker nodes.

### Bootstrap Node

The bootstrap node is a temporary node that is responsible for standing up the control plane on the masters. Only one bootstrap node will be stood up and it will be deprovisioned once the production control plane is ready. To do so, you need 1 instance, and 1 port. We recommend a flavor with a minimum of 16 GB RAM, 4 vCPUs, and 25 GB Disk (or Root Volume).

### Image Registry Requirements

If Swift is available in the cloud where the installation is being performed, it is used as the default backend for the OpenShift image registry. At the time of installation only an empty container is created without loading any data. Later on, for the system to work properly, you need to have enough free space to store the container images.

In this case the user must have `swiftoperator` permissions. As an OpenStack administrator:

```sh
openstack role add --user <user> --project <project> swiftoperator
```

If Swift is not available, the [PVC](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#persistentvolumeclaims) storage is used as the backend. For this purpose, a persistent volume of 100 GB will be created in Cinder and mounted to the image registry pod during the installation.

**Note:** If you are deploying a cluster in an Availability Zone where Swift isn't available but where Cinder is, it is recommended to deploy the Image Registry with Cinder backend. It will try to schedule the volume into the same AZ as the Nova zone where the PVC is located; otherwise it'll pick the default availability zone. If needed, the Image registry can be moved to another availability zone by a day 2 operation.

If you want to force Cinder to be used as a backend for the Image Registry, you need to remove the `swiftoperator` permissions. As an OpenStack administrator:

```sh
openstack role remove --user <user> --project <project> swiftoperator
```

**Note:** Since Cinder supports only [ReadWriteOnce](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#access-modes) access mode, it's not possible to have more than one replica of the image registry pod.

### Disk Requirements

Etcd, which runs on the control plane nodes, has [disk requirements](https://github.com/etcd-io/etcd/blob/master/Documentation/op-guide/hardware.md#disks) that need to be met to ensure the stability of the cluster.

Generally speaking, it is advised to choose for the control plane nodes a flavour that is backed by SSD in order to reduce latency.

If the ephemeral disk that gets attached to instances of the chosen flavor does not meet etcd requirements, check if the cloud has a more performant volume type and use a [custom `install-config.yaml`](customization.md) to deploy the control plane with root volumes. However, please note that Ceph RBD (and any other network-attached storage) can result in unpredictable network latencies. Prefer PCI passthrough of an NVM device instead.

In order to **measure the performance of your disk**, you can use [fio](https://github.com/axboe/fio):

```shell
sudo podman run \
	--volume "/var/lib/etcd:/mount:z" \
	docker.io/ljishen/fio \
		--directory=/mount \
		--name=iotest \
		--size=22m \
		--bs=2300 \
		--fdatasync=1 \
		--ioengine=sync \
		--rw=write
```

The command must be run as superuser.

In the command output, look for the the 99th percentile of [fdatasync](https://linux.die.net/man/2/fdatasync) durations (`fsync/fdatasync/sync_file_range` -> `sync percentiles`). The number must be less than 10ms (or 10000µs: fio fluidly adjusts the scale between ms/µs/ns depending on the numbers).

**Look for spikes.** Even if the baseline latency looks good, there may be spikes where it comes up, triggering issues that result in API being unavailable.

**Prometheus collects etcd-specific metrics.**

Once the cluster is up, Prometheus provides useful metrics here:

```
https://prometheus-k8s-openshift-monitoring.apps.<cluster name>.<domain name>/graph?g0.range_input=2h&g0.stacked=0&g0.expr=histogram_quantile(0.99%2C%20rate(etcd_disk_wal_fsync_duration_seconds_bucket%5B5m%5D))&g0.tab=0&g1.range_input=2h&g1.expr=histogram_quantile(0.99%2C%20rate(etcd_disk_backend_commit_duration_seconds_bucket%5B5m%5D))&g1.tab=0&g2.range_input=2h&g2.expr=etcd_server_health_failures&g2.tab=0
```

Click "Login with OpenShift", enter `kubeadmin` and the password printed out by the installer.

The units are in seconds and should stay under 10ms (0.01s) at all times. The `etcd_health` graph should remain at 0.

In order to collect relevant information interactively, **run the conformance tests**:

```
git clone https://github.com/openshift/origin/
make WHAT=cmd/openshift-tests
export KUBECONFIG=<path/to/kubeconfig>
_output/local/bin/linux/amd64/openshift-tests run openshift/conformance/parallel
```

The entire test suite takes over an hour to complete. Run it and check the Prometheus logs afterwards.

### Neutron Public Network

The public network should be created by the OpenStack administrator. Verify the name/ID of the 'External' network:

```sh
openstack network list --long -c ID -c Name -c "Router Type"
+--------------------------------------+----------------+-------------+
| ID                                   | Name           | Router Type |
+--------------------------------------+----------------+-------------+
| 148a8023-62a7-4672-b018-003462f8d7dc | public_network | External    |
+--------------------------------------+----------------+-------------+
```

**NOTE:** If the `neutron` `trunk` service plug-in is enabled, trunk port will be created by default. For more information, please refer to [neutron trunk port](https://wiki.openstack.org/wiki/Neutron/TrunkPort).

### Nova Metadata Service

Nova [metadata service](https://docs.openstack.org/nova/latest/user/metadata.html#metadata-service) must be enabled and available at `http://169.254.169.254`. Currently the service is used to deliver Ignition config files to Nova instances and provide information about the machine to `kubelet`.

### Glance Service

The creation of images in Glance should be available to the user. Now Glance is used for two things:

- Right after the installation starts, the installer automatically uploads the actual `RHCOS` binary image to Glance with the name `<clusterID>-rhcos`. The image exists throughout the life of the cluster and is removed along with it.

- The installer stores bootstrap ignition configs in a temporary image called `<clusterID>-ignition`. This is not a canonical use of the service, but this solution allows us to unify the installation process, since Glance is available on all OpenStack clouds, unlike Swift. The image exists for a limited period of time while the bootstrap process is running (normally 10-30 minutes), and then is automatically deleted.

## OpenStack Credentials

You must have a `clouds.yaml` file in order to run the installer. The installer will look for a `clouds.yaml` file in the following locations in order:

1. Value of `OS_CLIENT_CONFIG_FILE` environment variable
2. Current directory
3. unix-specific user config directory (`~/.config/openstack/clouds.yaml`)
4. unix-specific site config directory (`/etc/openstack/clouds.yaml`)

In many OpenStack distributions, you can generate a `clouds.yaml` file through Horizon. Otherwise, you can make a `clouds.yaml` file yourself.
Information on this file can be found [here](https://docs.openstack.org/openstacksdk/latest/user/config/configuration.html#config-files) and it looks like:

```yaml
clouds:
  shiftstack:
    auth:
      auth_url: http://10.10.14.42:5000/v3
      project_name: shiftstack
      username: shiftstack_user
      password: XXX
      user_domain_name: Default
      project_domain_name: Default
  dev-evn:
    region_name: RegionOne
    auth:
      username: 'devuser'
      password: XXX
      project_name: 'devonly'
      auth_url: 'https://10.10.14.22:5001/v2.0'
```

The file can contain information about several clouds. For instance, the example above describes two clouds: `shiftstack` and `dev-evn`.
In order to determine which cloud to use, the user can either specify it in the `install-config.yaml` file under `platform.openstack.cloud` or with `OS_CLOUD` environment variable. If both are omitted, then the cloud name defaults to `openstack`.

### Self Signed OpenStack CA certificates

If your OpenStack cluster uses self signed CA certificates for endpoint authentication, add the `cacert` key to your `clouds.yaml`. Its value should be a valid path to your CA cert, and the file should be readable by the user who runs the installer. The path can be either absolute, or relative to the current working directory while running the installer.

For example:

```yaml
clouds:
  shiftstack:
    auth: ...
    cacert: "/etc/pki/ca-trust/source/anchors/ca.crt.pem"
```

## Standalone Single-Node Development Environment

If you would like to set up an isolated development environment, you may use a bare metal host running CentOS 7. The following repository includes some instructions and scripts to help with creating a single-node OpenStack development environment for running the installer. Please refer to [this documentation](https://github.com/shiftstack-dev-tools/ocp-doit) for further details.

## Running The Installer

### Known Issues

OpenStack support has [known issues](known-issues.md). We will be documenting workarounds until we are able to resolve these bugs in the upcoming releases. To see the latest status of any bug, read through bugzilla or github link provided in that bug's description. If you know of a possible workaround that hasn't been documented yet, please comment in that bug's tracking link so we can address it as soon as possible. Also note that any bug listed in these documents is already a top priority issue for the dev team, and will be resolved as soon as possible. If you find more bugs during your runs, please read the section on [issue reporting](#reporting-issues).

### Initial Setup

Please head to [openshift.com/try](https://www.openshift.com/try) to get the latest versions of the installer, and instructions to run it.

Before running the installer, we recommend you create a directory for each cluster you plan to deploy. See the documents on the [recommended workflow](../overview.md#multiple-invocations) for more information about why you should do it this way.

```sh
mkdir ostest
cp install-config.yaml ostest/install-config.yaml
```

### API Access

All the OpenShift nodes get created in an OpenStack tenant network and as such, can't be accessed directly in most OpenStack deployments. We will briefly explain how to set up access to the OpenShift API with and without floating IP addresses.

#### Using Floating IPs

This method allows you to attach two floating IP (FIP) addresses to endpoints in OpenShift.

A standard deployment uses three floating IP addresses in total:

1. External access to the OpenShift API
2. External access to the workloads (apps) running on the OpenShift cluster
3. Temporary IP address for bootstrap log collection

The first two addresses (API and Ingress) are generally created up-front and have the corresponding DNS records resolve to them.

The third floating IP is created automatically by the installer and will be destroyed along with all the other bootstrap resources. If the bootstrapping process fails, the installer will try to SSH into the bootstrap node and collect the logs.

##### Create API and Ingress Floating IP Addresses

The deployed OpenShift cluster will need two floating IP addresses, one to attach to the API load balancer (lb FIP), and one for the OpenShift applications (apps FIP). Note that the LB FIP is the IP address you will add to your `install-config.yaml` or select in the interactive installer prompt.

You can create them like so:

```sh
openstack floating ip create --description "API <cluster name>.<base domain>" <external network>
# => <lb FIP>
openstack floating ip create --description "Ingress <cluster name>.<base domain>" <external network>
# => <apps FIP>
```

**NOTE:** These IP addresses will **not** show up attached to any particular server (e.g. when running `openstack server list`). Similarly, the API and Ingress ports will always be in the `DOWN` state.

This is because the ports are not attached to the servers directly. Instead, their fixed IP addresses are managed by keepalived. This has no record in Neutron's database and as such, is not visible to OpenStack.

*The network traffic will flow through even though the IPs and ports do not show up in the servers*.

For more details, read the [OpenShift on OpenStack networking infrastructure design document](../../design/openstack/networking-infrastructure.md).

##### Create API and Ingress DNS Records

You will also need to add the following records to your DNS:

```dns
api.<cluster name>.<base domain>.  IN  A  <lb FIP>
*.apps.<cluster name>.<base domain>.  IN  A  <apps FIP>
```

If you're unable to create and publish these DNS records, you can add them to your `/etc/hosts` file.

```dns
<lb FIP> api.<cluster name>.<base domain>
<apps FIP> console-openshift-console.apps.<cluster name>.<base domain>
<apps FIP> integrated-oauth-server-openshift-authentication.apps.<cluster name>.<base domain>
<apps FIP> oauth-openshift.apps.<cluster name>.<base domain>
<apps FIP> prometheus-k8s-openshift-monitoring.apps.<cluster name>.<base domain>
<apps FIP> grafana-openshift-monitoring.apps.<cluster name>.<base domain>
<apps FIP> <app name>.apps.<cluster name>.<base domain>
```

**WARNING:** *this workaround will make the API accessible only to the computer with these `/etc/hosts` entries. This is fine for your own testing (and it is enough for the installation to succeed), but it is not enough for a production deployment. In addition, if you create new OpenShift apps or routes, you will have to add their entries too, because `/etc/hosts` does not support wildcard entries.*

##### External API Access

If you have specified the API floating IP (either via the installer prompt or by adding the `apiFloatingIP` entry in your `install-config.yaml`) the installer will attach the Floating IP address to the `api-port` automatically.

If you have created the API DNS record, you should be able access the OpenShift API.

##### External Ingress (apps) Access

In the same manner, you may have specified an Ingress floating IP by adding the `ingressFloatingIP` entry in your `install-config.yaml`, in which case the installer attaches the Floating IP address to the `ingress-port` automatically.

If `ingressFloatingIP` is empty or absent in `install-config.yaml`, the Ingress port will be created but not attached to any floating IP. You can manually attach the Ingress floating IP to the ingress-port after the cluster is created.

That can be done in the following steps:

```sh
openstack port show <cluster name>-<clusterID>-ingress-port
```

Then attach the FIP to it:

```sh
openstack floating ip set --port <cluster name>-<clusterID>-ingress-port <apps FIP>
```

This assumes the floating IP and corresponding `*.apps` DNS record exists.


#### Without Floating IPs

If you cannot or don't want to pre-create a floating IP address, the installation should still succeed, however the installer will fail waiting for the API.

**WARNING:** The installer will fail if it can't reach the bootstrap OpenShift API in 20 minutes.

Even if the installer times out, the OpenShift cluster should still come up. Once the bootstrapping process is in place, it should all run to completion. So you should be able to deploy OpenShift without any floating IP addresses and DNS records and create everything yourself after the cluster is up.

### Running a Deployment

To run the installer, you have the option of using the interactive wizard, or providing your own `install-config.yaml` file for it. The wizard is the easier way to run the installer, but passing your own `install-config.yaml` enables you to use more fine grained customizations. If you are going to create your own `install-config.yaml`, read through the available [OpenStack customizations](customization.md). For information on running the installer with Kuryr, see the [Kuryr docs](kuryr.md).

```sh
./openshift-install create cluster --dir ostest
```

If you want to create an install config without deploying a cluster, you can use the command:

```sh
./openshift-install create install-config --dir ostest
```

### Current Expected Behavior

Currently:

- Deploys an isolated tenant network
- Deploys a bootstrap instance to bootstrap the OpenShift cluster
- Deploys 3 master nodes
- Once the masters are deployed, the bootstrap instance is destroyed
- Deploys 3 worker nodes

Look for a message like this to verify that your install succeeded:

```txt
INFO Install complete!
INFO To access the cluster as the system:admin user when using 'oc', run 'export KUBECONFIG=/home/stack/ostest/auth/kubeconfig'
INFO Access the OpenShift web-console here: https://console-openshift-console.apps.ostest.shiftstack.com
INFO Login to the console with user: kubeadmin, password: xxx
```

### Checking Cluster Status

If you want to see the status of the apps and services in your cluster during, or after a deployment, first export your administrator's kubeconfig:

```sh
export KUBECONFIG=ostest/auth/kubeconfig
```

After a finished deployment, there should be a node for each master and worker server created. You can check this with the command:

```sh
oc get nodes
```

To see the version of your OpenShift cluster, do:

```sh
oc get clusterversion
```

To see the status of you operators, do:

```sh
oc get clusteroperator
```

Finally, to see all the running pods in your cluster, you can do:

```sh
oc get pods -A
```
### Destroying The Cluster

To destroy the cluster, point it to your cluster with this command:

```sh
./openshift-install --log-level debug destroy cluster --dir ostest
```

Then, you can delete the folder containing the cluster metadata:

```sh
rm -rf ostest/
```
## Post Install Operations

### Adding a MachineSet

Groups of Compute nodes are managed using the [MachineSet][machine-set-code] resource. It is possible to create additional MachineSets post-install, for example to assign workloads to specific machines.

When running on OpenStack, the MachineSet has platform-specific fields under `spec.template.spec.providerSpec.value`. For more information about the values that you can set in the `providerSpec`, see [the API definition](provider-spec-definition).

```yaml
apiVersion: machine.openshift.io/v1beta1
kind: MachineSet
metadata:
  labels:
    machine.openshift.io/cluster-api-cluster: <infrastructure_ID>
    machine.openshift.io/cluster-api-machine-role: <node_role>
    machine.openshift.io/cluster-api-machine-type: <node_role>
  name: <infrastructure_ID>-<node_role>
  namespace: openshift-machine-api
spec:
  replicas: <number_of_replicas>
  selector:
    matchLabels:
      machine.openshift.io/cluster-api-cluster: <infrastructure_ID>
      machine.openshift.io/cluster-api-machineset: <infrastructure_ID>-<node_role>
  template:
    metadata:
      labels:
        machine.openshift.io/cluster-api-cluster: <infrastructure_ID>
        machine.openshift.io/cluster-api-machine-role: <node_role>
        machine.openshift.io/cluster-api-machine-type: <node_role>
        machine.openshift.io/cluster-api-machineset: <infrastructure_ID>-<node_role>
    spec:
      providerSpec:
        value:
          apiVersion: openstackproviderconfig.openshift.io/v1alpha1
          cloudName: openstack
          cloudsSecret:
            name: openstack-cloud-credentials
            namespace: openshift-machine-api
          flavor: <nova_flavor>
          image: <glance_image_name_or_location>
          serverGroupID: <UUID of the pre-created Nova server group (optional)>
          kind: OpenstackProviderSpec
          networks:
          - filter: {}
            subnets:
            - filter:
                name: <subnet_name>
                tags: openshiftClusterID=<infrastructure_ID>
          securityGroups:
          - filter: {}
            name: <infrastructure_ID>-<node_role>
          serverMetadata:
            Name: <infrastructure_ID>-<node_role>
            openshiftClusterID: <infrastructure_ID>
          tags:
          - openshiftClusterID=<infrastructure_ID>
          trunk: true
          userDataSecret:
            name: <node_role>-user-data
          availabilityZone: <optional_openstack_availability_zone>
```

[provider-spec-definition]: https://github.com/openshift/cluster-api-provider-openstack/blob/155384b859c5b2fb5b7f11c9111d3f8e4f3066bd/pkg/apis/openstackproviderconfig/v1alpha1/types.go#L31

#### Defining a MachineSet That Uses Multiple Networks

To define a MachineSet with multiple networks, the `primarySubnet` value in the `providerSpec` must be set to the OpenStack subnet that you want the Kubernetes endpoints of the nodes to be published on. For most use cases, this is the same subnet as the [machinesSubnet](./customization.md#cluster-scoped-properties) in the `install-config.yaml`.

 After you set the subnet, add all of the networks that you want to attach to your machines to the `Networks` list in `providerSpec`. You must also add the network that the primary subnet is part of to this list.

#### Using a Server Group

In order to hint the Nova scheduler to spread the Machines across different
hosts, first create a Server Group with the [desired
policy][server-group-docs]:

```shell
openstack server group create --policy=anti-affinity <server-group-name>
## OR ##
openstack --os-compute-api-version=2.15 server group create --policy=soft-anti-affinity <server-group-name>
```

If the command is successful, the OpenStack CLI will return the ID of the newly
created Server Group. Paste it in the optional `serverGroupID` property of the
MachineSet.

#### Setting Nova Availability Zones

In order to use Availability Zones, create one MachineSet per target
Availability Zone, and set the Availability Zone in the `availabilityZone`
property of the MachineSet.

**NOTE:** Note when deploying with `Kuryr` there is an Octavia API loadbalancer VM that will not fulfill the Availability Zones restrictions due to Octavia lack of support for it. In addition, if Octavia only has the amphora provider instead of also the OVN-Octavia provider, all the OpenShift services will be backed up by Octavia Load Balancer VMs which will not fulfill the Availability Zone restrictions either.

[machine-set-code]: https://github.com/openshift/cluster-api-provider-openstack/blob/master/pkg/apis/openstackproviderconfig/v1alpha1/types.go
[server-group-docs]: https://docs.openstack.org/api-ref/compute/?expanded=create-server-group-detail#create-server-group

### Using a Custom External Load Balancer

You can shift ingress/egress traffic from the default OpenShift on OpenStack load balancer to a load balancer that you provide. To do so, the instance that it runs from must be able to access every machine in your cluster. You might ensure this access by creating the instance on a subnet that is within your cluster's network, and then attaching a router interface to that subnet from the `OpenShift-external-router` [object/instance/whatever]. This can also be accomplished by attaching floating ips to the machines you want to add to your load balancer.

#### External Facing OpenShift Services

Add the following external facing services to your new load balancer:

- The master nodes serve the OpenShift API on port 6443 using TCP.
- The apps hosted on the worker nodes are served on ports 80, and 443. They are both served using TCP.

Note: Make sure the instance that your new load balancer is running on has security group rules that allow TCP traffic over these ports.

#### HAProxy Example Load Balancer Config

The following `HAProxy` config file demonstrates a basic configuration for an external load balancer:

```haproxy
listen <cluster-name>-api-6443
        bind 0.0.0.0:6443
        mode tcp
        balance roundrobin
        server <cluster-name>-master-0 192.168.0.154:6443 check
        server <cluster-name>-master-1 192.168.0.15:6443 check
        server <cluster-name>-master-2 192.168.3.128:6443 check
listen <cluster-name>-apps-443
        bind 0.0.0.0:443
        mode tcp
        balance roundrobin
        server <cluster-name>-worker-0 192.168.3.18:443 check
        server <cluster-name>-worker-1 192.168.2.228:443 check
        server <cluster-name>-worker-2 192.168.1.253:443 check
listen <cluster-name>-apps-80
        bind 0.0.0.0:80
        mode tcp
        balance roundrobin
        server <cluster-name>-worker-0 192.168.3.18:80 check
        server <cluster-name>-worker-1 192.168.2.228:80 check
        server <cluster-name>-worker-2 192.168.1.253:80 check
```

#### DNS Lookups

To ensure that your API and apps are accessible through your load balancer, [create or update your DNS entries](#create-api-and-ingress-dns-records) for those endpoints. To use your new load balancing service for external traffic, make sure the IP address for these DNS entries is the IP address your load balancer is reachable at.

```dns
<load balancer ip> api.<cluster-name>.<base domain>
<load balancer ip> apps.<cluster-name>.base domain>
```

#### Verifying that the API is Reachable

One good way to test whether or not you can reach the API is to run the `oc` command. If you can't do that easily, you can use this curl command:

```sh
curl https://api.<cluster-name>.<base domain>:6443/version --insecure
```

Result:

```json
{
  "major": "1",
  "minor": "19",
  "gitVersion": "v1.19.2+4abb4a7",
  "gitCommit": "4abb4a77838037b8dbb8e4ca34e63c4a129654c8",
  "gitTreeState": "clean",
  "buildDate": "2020-11-12T05:46:36Z",
  "goVersion": "go1.15.2",
  "compiler": "gc",
  "platform": "linux/amd64"
}
```

Note: The versions in the sample output may differ from your own. As long as you get a JSON payload response, the API is accessible.

#### Verifying that Apps Reachable

The simplest way to verify that apps are reachable is to open the OpenShift console in a web browser. If you don't have access to a web browser, query the console with the following curl command:

```sh
curl http://console-openshift-console.apps.<cluster-name>.<base domain> -I -L --insecure
```


Result:

```http
HTTP/1.1 302 Found
content-length: 0
location: https://console-openshift-console.apps.<cluster-name>.<base domain>/
cache-control: no-cacheHTTP/1.1 200 OK
referrer-policy: strict-origin-when-cross-origin
set-cookie: csrf-token=39HoZgztDnzjJkq/JuLJMeoKNXlfiVv2YgZc09c3TBOBU4NI6kDXaJH1LdicNhN1UsQWzon4Dor9GWGfopaTEQ==; Path=/; Secure
x-content-type-options: nosniff
x-dns-prefetch-control: off
x-frame-options: DENY
x-xss-protection: 1; mode=block
date: Tue, 17 Nov 2020 08:42:10 GMT
content-type: text/html; charset=utf-8
set-cookie: 1e2670d92730b515ce3a1bb65da45062=9b714eb87e93cf34853e87a92d6894be; path=/; HttpOnly; Secure; SameSite=None
cache-control: private
```

### Reconfiguring cloud provider

If you need to update the OpenStack cloud provider configuration you can edit the ConfigMap containing it:

```sh
oc edit configmap -n openshift-config cloud-provider-config
```

**NOTE:** It can take a while to reconfigure the cluster depending on the size of it. The reconfiguration is completed once no node is getting `SchedulingDisabled` taint anymore.

There are several things you can change:

#### Modifying cloud provider options

If you need to modify the direct cloud provider options, then edit the `config` key in the ConfigMap. A brief list of possible options is shown in [Cloud Provider configuration](./customization.md#cloud-provider-configuration) section.


#### Refreshing a CA Certificate

If you ran the installer with a [custom CA certificate](#self-signed-openstack-ca-certificates), then this certificate can be changed while the cluster is running. To change your certificate, edit the value of the `ca-cert.pem` key in the `cloud-provider-config` configmap with a valid PEM certificate.

## Reporting Issues

Please see the [Issue Tracker][issues_openstack] for current known issues.
Please report a new issue if you do not find an issue related to any trouble you’re having.

[issues_openstack]: https://github.com/openshift/installer/issues?utf8=%E2%9C%93&q=is%3Aissue+is%3Aopen+openstack
