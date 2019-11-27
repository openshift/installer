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
    - [Swift](#swift)
    - [Disk Requirements](#disk-requirements)
    - [Red Hat Enterprise Linux CoreOS (RHCOS)](#red-hat-enterprise-linux-coreos-rhcos)
    - [Neutron Public Network](#neutron-public-network)
    - [Nova Metadata Service](#nova-metadata-service)
  - [OpenStack Credentials](#openstack-credentials)
  - [Standalone Single-Node Development Environment](#standalone-single-node-development-environment)
  - [Running The Installer](#running-the-installer)
    - [Known Issues](#known-issues)
    - [Initial Setup](#initial-setup)
    - [API Access](#api-access)
      - [Using Floating IPs](#using-floating-ips)
      - [Without Floating IPs](#without-floating-ips)
    - [Running a Deployment](#running-a-deployment)
    - [Current Expected Behavior](#current-expected-behavior)
    - [Checking Cluster Status](#checking-cluster-status)
    - [Destroying The Cluster](#destroying-the-cluster)
  - [Using an External Load Balancer](#using-an-external-load-balancer)
  - [Reporting Issues](#reporting-issues)

## Reference Documents

- [Known Issues and Workarounds](known-issues.md)
- [Using the OSP 4 installer with Kuryr](kuryr.md)
- [Troubleshooting your cluster](troubleshooting.md)
- [Customizing your install](customization.md)
- [Learn about the OpenShift on OpenStack networking infrastructure design](../../design/openstack/networking-infrastructure.md)

## OpenStack Requirements

In order to run the latest version of the installer in OpenStack, at a bare minimum you need the following quota to run a *default* cluster. While it is possible to run the cluster with fewer resources than this, it is not recommended. Certain cases, such as deploying [without FIPs](#without-floating-ips), or deploying with an [external load balancer](#using-an-external-load-balancer) are documented below, and are not included in the scope of this recommendation. If you are planning on using Kuryr, or want to learn more about it, please read through the [Kuryr documentation](kuryr.md). **NOTE: The installer has been tested and developed on Red Hat OSP 13.**

For a successful installation it is required:

- Floating IPs: 2
- Security Groups: 3
- Security Group Rules: 60
- Routers: 1
- Subnets: 1
- RAM: 112 GB
- vCPUs: 28
- Volume Storage: 175 GB
- Instances: 7
- Swift containers: 2
- Swift objects: 1
- Available space in Swift: at least 10 MB

**NOTE:** Size depends on the size of the bootstrap ignition file.

You may need to increase the security group related quotas from their default values. For example (as an OpenStack administrator):

```sh
openstack quota set --secgroups 8 --secgroup-rules 100 <project>`
```

### Master Nodes

The default deployment stands up 3 master nodes, which is the minimum amount required for a cluster. For each master node you stand up, you will need 1 instance, and 1 port available in your quota. They should be assigned a flavor with at least 16 GB RAM, 4 vCPUs, and 25 GB Disk. It is theoretically possible to run with a smaller flavor, but be aware that if it takes too long to stand up services, or certain essential services crash, the installer could time out, leading to a failed install.

### Worker Nodes

The default deployment stands up 3 worker nodes. In our testing we determined that 2 was the minimum number of workers you could have to get a successful install, but we don't recommend running with that few. Worker nodes host the applications you run on OpenShift, so it is in your best interest to have more of them. See [here](https://docs.openshift.com/enterprise/3.0/architecture/infrastructure_components/kubernetes_infrastructure.html#node) for more information. The flavor assigned to the worker nodes should have at least 2 vCPUs, 8 GB RAM and 25 GB Disk. However, if you are experiencing `Out Of Memory` issues, or your installs are timing out, you should increase the size of your flavor to match the masters: 4 vCPUs and 16 GB RAM.

### Bootstrap Node

The bootstrap node is a temporary node that is responsible for standing up the control plane on the masters. Only one bootstrap node will be stood up and it will be deprovisioned once the production control plane is ready. To do so, you need 1 instance, and 1 port. We recommend a flavor with a minimum of 16 GB RAM, 4 vCPUs, and 25 GB Disk.

### Swift

Swift is required for installation as the user-data provided by OpenStack Metadata service is not big enough to store the ignition config files, so they are served by Swift instead.

Swift is also used as a backend for the OpenShift image registry, but at the time of installation only an empty container is created without loading any data. Later on, for the system to work properly, you need to have enough free space to store the container images.

The user must have `swiftoperator` permissions and `temp-url` support must be enabled. As an OpenStack administrator:

```sh
openstack role add --user <user> --project <project> swiftoperator
openstack object store account set --property Temp-URL-Key=superkey
```

### Disk Requirements

Etcd runs on the control plane nodes, and has disk requirements that need to be met to ensure the stability of the cluster. If the ephemeral disk that gets attached to instances of the chosen flavor does not meet [etcd requirements](https://github.com/etcd-io/etcd/blob/master/Documentation/op-guide/hardware.md#disks), check if the cloud has a more performant volume type and use a [custom `install-config.yaml`](customization.md) to deploy the control plane with root volumes.

### Red Hat Enterprise Linux CoreOS (RHCOS)

Get the latest RHCOS image [here](https://mirror.openshift.com/pub/openshift-v4/dependencies/rhcos/pre-release/latest/). The installer requires a proper RHCOS image in the OpenStack cluster or project:

```sh
openstack image create --container-format=bare --disk-format=qcow2 --file rhcos-${RHCOSVERSION}-openstack.qcow2 rhcos
```

**NOTE:** Depending on your OpenStack environment you can upload the RHCOS image as `raw` or `qcow2`. See [Disk and container formats for images](https://docs.openstack.org/image-guide/image-formats.html) for more information. At the time of writing, the installer looks for an image named `rhcos`.

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

## Standalone Single-Node Development Environment

If you would like to set up an isolated development environment, you may use a bare metal host running CentOS 7. The following repository includes some instructions and scripts to help with creating a single-node OpenStack development environment for running the installer. Please refer to [this documentation](https://github.com/shiftstack-dev-tools/ocp-doit) for further details.

## Running The Installer

### Known Issues

OpenStack support has [known issues](known-issues.md). We will be documenting workarounds until we are able to resolve these bugs in the upcoming releases. To see the latest status of any bug, read through bugzilla or github link provided in that bug's description. If you know of a possible workaround that hasn't been documented yet, please comment in that bug's tracking link so we can address it as soon as possible. Also note that any bug listed in these documents is already a top priority issue for the dev team, and will be resolved as soon as possible. If you find more bugs during your runs, please read the section on [issue reporting](#reporting-issues).

### Initial Setup

Please head to [try.openshift.com](https://try.openshift.com) to get the latest versions of the installer, and instructions to run it.

Before running the installer, we recommend you create a directory for each cluster you plan to deploy. See the documents on the [recommended workflow](../overview.md#multiple-invocations) for more information about why you should do it this way.

```sh
mkdir ostest
cp install-config.yaml ostest/install-config.yaml
```

### API Access

All the OpenShift nodes get created in an OpenStack tenant network and as such, can't be accessed directly in most OpenStack deployments. We will briefly explain how to set up access to the OpenShift API with and without floating IP addresses.

#### Using Floating IPs

This method allows you to attach two floating IP (FIP) addresses to endpoints in OpenShift.

You will need to create two floating IP addresses, one to attach to the API load balancer (lb FIP), and one for the OpenShift applications (apps FIP). Note that the LB FIP is the same floating IP as the one you added to your `install-config.yaml`. The following command is an example of how to create floating IPs:

```sh
openstack floating ip create <external network>
```

You will also need to add the following records to your DNS:

```dns
api.<cluster name>.<base domain>  IN  A  <lb FIP>
```

If you don't have a DNS server under your control, you should add the records to your `/etc/hosts` file.

**NOTE:** *this will make the API accessible only to you. This is fine for your own testing (and it is enough for the installation to succeed), but it is not enough for a production deployment.*

In order to reach the applications running on your worker nodes, you should attach a floating IP to the `ingress-port` at the end of your install. That can be done in the following steps:

```sh
openstack port show <cluster name>-<clusterID>-ingress-port
```

Then attach the FIP to it:

```sh
openstack floating ip set --port <ingress port id> <apps FIP>
```

Add a wildcard A record for `*apps.` in your DNS:

```dns
*.apps.<cluster name>.<base domain>  IN  A  <apps FIP>
```

Alternatively, if you don't control the DNS, you can add the hostnames in `/etc/hosts`:

```dns
<apps FIP> console-openshift-console.apps.<cluster name>.<base domain>
<apps FIP> integrated-oauth-server-openshift-authentication.apps.<cluster name>.<base domain>
<apps FIP> oauth-openshift.apps.<cluster name>.<base domain>
<apps FIP> prometheus-k8s-openshift-monitoring.apps.<cluster name>.<base domain>
<apps FIP> grafana-openshift-monitoring.apps.<cluster name>.<base domain>
<apps FIP> <app name>.apps.<cluster name>.<base domain>
```

#### Without Floating IPs

If you cannot or don't want to pre-create a floating IP address, the installation should still succeed, however the installer will fail waiting for the API.

**WARNING:** The installer will fail if it can't reach the bootstrap OpenShift API in 30 minutes.

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

## Using an External Load Balancer

This documents how to shift from the internal load balancer, which is intended for internal networking needs, to an external load balancer.

The load balancer must serve ports 6443, 443, and 80 to any users of the system.  Port 22623 is for serving ignition start-up configurations to the OpenShift nodes and should not be reachable outside of the cluster.

The first step is to add floating IPs to all the master nodes:

```sh
openstack floating ip create --port master-port-0 <public network>
openstack floating ip create --port master-port-1 <public network>
openstack floating ip create --port master-port-2 <public network>
```

Once complete you can see your floating IPs using:

```sh
openstack server list
```

These floating IPs can then be used by the load balancer to access the cluster.  An example of HAProxy configuration for port 6443 is below.

```txt
listen <cluster name>-api-6443
    bind 0.0.0.0:6443
    mode tcp
    balance roundrobin
    server <cluster name>-master-2 <floating ip>:6443 check
    server <cluster name>-master-0 <floating ip>:6443 check
    server <cluster name>-master-1 <floating ip>:6443 check
```

The other port configurations are identical.

The next step is to allow network access from the load balancer network to the master nodes:

```sh
openstack security group rule create master --remote-ip <load balancer CIDR> --ingress --protocol tcp --dst-port 6443
openstack security group rule create master --remote-ip <load balancer CIDR> --ingress --protocol tcp --dst-port 443
openstack security group rule create master --remote-ip <load balancer CIDR> --ingress --protocol tcp --dst-port 80
```

You could also specify a specific IP address with /32 if you wish.

You can verify the operation of the load balancer now if you wish, using the curl commands given below.

Now the DNS entry for `api.<cluster name>.<base domain>` needs to be updated to point to the new load balancer:

```dns
<load balancer ip> api.<cluster-name>.<base domain>
```

The external load balancer should now be operational along with your own DNS solution. The following curl command is an example of how to check functionality:

```sh
curl https://<loadbalancer-ip>:6443/version --insecure
```

Result:

```json
{
  "major": "1",
  "minor": "11+",
  "gitVersion": "v1.11.0+ad103ed",
  "gitCommit": "ad103ed",
  "gitTreeState": "clean",
  "buildDate": "2019-01-09T06:44:10Z",
  "goVersion": "go1.10.3",
  "compiler": "gc",
  "platform": "linux/amd64"
}
```

Another useful thing to check is that the ignition configurations are only available from within the deployment. The following command should only succeed from a node in the OpenShift cluster:

```sh
curl https://<loadbalancer ip>:22623/config/master --insecure
```

## Reporting Issues

Please see the [Issue Tracker][issues_openstack] for current known issues.
Please report a new issue if you do not find an issue related to any trouble you’re having.

[issues_openstack]: https://github.com/openshift/installer/issues?utf8=%E2%9C%93&q=is%3Aissue+is%3Aopen+openstack
