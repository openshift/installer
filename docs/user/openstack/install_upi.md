# Installing OpenShift on OpenStack User-Provisioned Infrastructure

The User-Provisioned Infrastructure (UPI) process installs OpenShift in stages, providing opportunities for modifications or integrating with existing infrastructure.

It contrasts with the fully-automated Installer-Provisioned Infrastructure (IPI) which creates everything in one go.

With UPI, creating the cloud (OpenStack) resources (e.g. Nova servers, Neutron ports, security groups) is the responsibility of the person deploying OpenShift.

The installer is still used to generate the ignition files and monitor the installation process.

This provides a greater flexibility at the cost of a more explicit and interactive process.

Below is a step-by-step guide to a UPI installation that mimics an automated IPI installation; prerequisites and steps described below should be adapted to the constraints of the target infrastructure.

Please be aware of the [Known Issues](known-issues.md#known-issues-specific-to-user-provisioned-installations)
of this method of installation.

## Table of Contents


- [Installing OpenShift on OpenStack User-Provisioned Infrastructure](#installing-openshift-on-openstack-user-provisioned-infrastructure)
  - [Table of Contents](#table-of-contents)
  - [Prerequisites](#prerequisites)
  - [Install Ansible](#install-ansible)
    - [RHEL](#rhel)
    - [Fedora](#fedora)
  - [OpenShift Configuration Directory](#openshift-configuration-directory)
  - [Red Hat Enterprise Linux CoreOS (RHCOS)](#red-hat-enterprise-linux-coreos-rhcos)
  - [API and Ingress Floating IP Addresses](#api-and-ingress-floating-ip-addresses)
  - [Install Config](#install-config)
    - [Configure the machineNetwork.CIDR apiVIP and ingressVIP](#configure-the-machinenetworkcidr-apivip-and-ingressvip)
    - [Empty Compute Pools](#empty-compute-pools)
    - [Modify NetworkType (Required for Kuryr SDN)](#modify-networktype-required-for-kuryr-sdn)
  - [Edit Manifests](#edit-manifests)
    - [Remove Machines and MachineSets](#remove-machines-and-machinesets)
    - [Make control-plane nodes unschedulable](#make-control-plane-nodes-unschedulable)
  - [Ignition Config](#ignition-config)
    - [Infra ID](#infra-id)
    - [Bootstrap Ignition](#bootstrap-ignition)
      - [Edit the Bootstrap Ignition](#edit-the-bootstrap-ignition)
      - [Upload the Boostrap Ignition](#upload-the-boostrap-ignition)
        - [Example 1: Swift](#example-1-swift)
        - [Example 2: Glance image service](#example-2-glance-image-service)
    - [Create the Bootstrap Ignition Shim](#create-the-bootstrap-ignition-shim)
      - [Ignition file served by server using self-signed certificate](#ignition-file-served-by-server-using-self-signed-certificate)
    - [Master Ignition](#master-ignition)
  - [Network Topology](#network-topology)
    - [Security Groups](#security-groups)
    - [Network, Subnet and external router](#network-subnet-and-external-router)
    - [Subnet DNS (optional)](#subnet-dns-optional)
  - [Bootstrap](#bootstrap)
  - [Control Plane](#control-plane)
    - [Control Plane Trunks (Kuryr SDN)](#control-plane-trunks-kuryr-sdn)
    - [Wait for the Control Plane to Complete](#wait-for-the-control-plane-to-complete)
    - [Access the OpenShift API](#access-the-openshift-api)
    - [Delete the Bootstrap Resources](#delete-the-bootstrap-resources)
  - [Compute Nodes](#compute-nodes)
    - [Compute Nodes Trunks (Kuryr SDN)](#compute-nodes-trunks-kuryr-sdn)
    - [Approve the worker CSRs](#approve-the-worker-csrs)
    - [Wait for the OpenShift Installation to Complete](#wait-for-the-openshift-installation-to-complete)
    - [Compute Nodes with SR-IOV NICs](#compute-nodes-with-sr-iov-nics)
  - [Destroy the OpenShift Cluster](#destroy-the-openshift-cluster)

## Prerequisites

The file `inventory.yaml` contains the variables most likely to need customization.
**NOTE**: some of the default pods (e.g. the `openshift-router`) require at least two nodes so that is the effective minimum.

The requirements for UPI are broadly similar to the [ones for OpenStack IPI][ipi-reqs]:

[ipi-reqs]: ./README.md#openstack-requirements

- OpenStack account with `clouds.yaml`
  - input in the `openshift-install` wizard
- Nova flavors
  - inventory: `os_flavor_master` and `os_flavor_worker`
- An external subnet for external connectivity. Required if any of the floating IPs is set in the inventory.
  - inventory: `os_external_network`
- The `openshift-install` binary
- A subnet range for the Nova servers / OpenShift Nodes, that does not conflict with your existing network
  - inventory: `os_subnet_range`
- A cluster name you will want to use
  - input in the `openshift-install` wizard
- A base domain
  - input in the `openshift-install` wizard
- OpenShift Pull Secret
  - input in the `openshift-install` wizard
- A DNS zone you can configure
  - it must be the resolver for the base domain, for the installer and for the end-user machines
  - it will host two records: for API and apps access

For an installation with Kuryr SDN on UPI, you should also check the requirements which are the same
needed for [OpenStack IPI with Kuryr][ipi-reqs-kuryr]. Please also note that **RHEL 7 nodes are not
supported on deployments configured with Kuryr**. This is because Kuryr container images are based on
RHEL 8 and may not work properly when run on RHEL 7.

[ipi-reqs-kuryr]: ./kuryr.md#requirements-when-enabling-kuryr

## Install Ansible

This repository contains [Ansible playbooks][ansible-upi] to deploy OpenShift on OpenStack.

They can be downloaded from Github with this script:

```sh
RELEASE="release-4.6"; xargs -n 1 curl -O <<< "
        https://raw.githubusercontent.com/openshift/installer/${RELEASE}/upi/openstack/bootstrap.yaml
        https://raw.githubusercontent.com/openshift/installer/${RELEASE}/upi/openstack/common.yaml
        https://raw.githubusercontent.com/openshift/installer/${RELEASE}/upi/openstack/compute-nodes.yaml
        https://raw.githubusercontent.com/openshift/installer/${RELEASE}/upi/openstack/control-plane.yaml
        https://raw.githubusercontent.com/openshift/installer/${RELEASE}/upi/openstack/down-bootstrap.yaml
        https://raw.githubusercontent.com/openshift/installer/${RELEASE}/upi/openstack/down-compute-nodes.yaml
        https://raw.githubusercontent.com/openshift/installer/${RELEASE}/upi/openstack/down-control-plane.yaml
        https://raw.githubusercontent.com/openshift/installer/${RELEASE}/upi/openstack/down-load-balancers.yaml
        https://raw.githubusercontent.com/openshift/installer/${RELEASE}/upi/openstack/down-network.yaml
        https://raw.githubusercontent.com/openshift/installer/${RELEASE}/upi/openstack/down-security-groups.yaml
        https://raw.githubusercontent.com/openshift/installer/${RELEASE}/upi/openstack/down-containers.yaml
        https://raw.githubusercontent.com/openshift/installer/${RELEASE}/upi/openstack/inventory.yaml
        https://raw.githubusercontent.com/openshift/installer/${RELEASE}/upi/openstack/network.yaml
        https://raw.githubusercontent.com/openshift/installer/${RELEASE}/upi/openstack/security-groups.yaml"
```

For installing a different version, change the branch (`release-4.6`)
accordingly (e.g. `release-4.7`). Note that `down-containers.yaml` was only
introduced in `release-4.6`.

**Requirements:**

* Python
* Ansible
* Python modules required in the playbooks. Namely:
  * openstackclient
  * openstacksdk
  * netaddr

### RHEL

From a RHEL 8 box, make sure that the repository origins are all set:

```sh
sudo subscription-manager register # if not done already
sudo subscription-manager attach --pool=$YOUR_POOLID # if not done already
sudo subscription-manager repos --disable=* # if not done already

sudo subscription-manager repos \
  --enable=rhel-8-for-x86_64-baseos-rpms \
  --enable=openstack-16-tools-for-rhel-8-x86_64-rpms \
  --enable=ansible-2.9-for-rhel-8-x86_64-rpms \
  --enable=rhel-8-for-x86_64-appstream-rpms
```

Then install the packages:
```sh
sudo dnf install python3-openstackclient ansible python3-openstacksdk python3-netaddr
```

Make sure that `python` points to Python3:
```sh
sudo alternatives --set python /usr/bin/python3
```

### Fedora

This command installs all required dependencies on Fedora:

```sh
sudo dnf install python3-openstackclient ansible python3-openstacksdk python3-netaddr
```

[ansible-upi]: ../../../upi/openstack "Ansible Playbooks for Openstack UPI"

## OpenShift Configuration Directory

All the configuration files, logs and installation state are kept in a single directory:

```sh
$ mkdir -p openstack-upi
$ cd openstack-upi
```

## Red Hat Enterprise Linux CoreOS (RHCOS)

A proper [RHCOS][rhcos] image in the OpenStack cluster or project is required for successful installation.

Get the RHCOS image for your OpenShift version [here][rhcos-image]. You should download images with the highest version that is less than or equal to the OpenShift version that you install. Use the image versions that match your OpenShift version if they are available.

The OpenStack QCOW2 image is delivered in compressed format and therefore has the `.gz` extension. Unfortunately, compressed image support is not supported in OpenStack. So, you have to decompress the data before uploading it into Glance. The following command will unpack the image and create `rhcos-${RHCOSVERSION}-openstack.qcow2` file without `.gz` extension.

```sh
$ gunzip rhcos-${RHCOSVERSION}-openstack.qcow2.gz
```

Next step is to create a Glance image.

**NOTE:** *This document* will use `rhcos` as the Glance image name, but it's not mandatory.

```sh
$ openstack image create --container-format=bare --disk-format=qcow2 --file rhcos-${RHCOSVERSION}-openstack.qcow2 rhcos
```

**NOTE:** Depending on your OpenStack environment you can upload the RHCOS image as `raw` or `qcow2`. See [Disk and container formats for images](https://docs.openstack.org/image-guide/introduction.html#disk-and-container-formats-for-images) for more information.

[qemu_guest_agent]: https://docs.openstack.org/nova/latest/admin/configuration/hypervisor-kvm.html
If the RHCOS image being used supports it,  the [KVM Qemu Guest Agent][qemu_guest_agent] may be used to enable optional
access between OpenStack KVM hypervisors and the cluster nodes.

To enable this feature, you must add the `hw_qemu_guest_agent=yes` property to the image:

```
$ openstack image rhcos update --property hw_qemu_guest_agent=yes
```

Finally validate that the image was successfully created:

```sh
$ openstack image show rhcos
```

[rhcos]: https://www.openshift.com/learn/coreos/
[rhcos-image]: https://mirror.openshift.com/pub/openshift-v4/dependencies/rhcos/

## API and Ingress Floating IP Addresses

If the variables `os_api_fip`, `os_ingress_fip` and `os_bootstrap_fip` are found in `inventory.yaml`, the corresponding floating IPs will be attached to the API load balancer, to the worker nodes load balancer and to the temporary machine used for the install process, respectively. Note that `os_external_network` is a requirement for those.

**NOTE**: throughout this document, we will use `203.0.113.23` as the public IP address for the OpenShift API endpoint and `203.0.113.19` as the public IP for the ingress (`*.apps`) endpoint. `203.0.113.20` will be the public IP used for the bootstrap machine.

```sh
$ openstack floating ip create --description "OpenShift API" <external>
=> 203.0.113.23
$ openstack floating ip create --description "OpenShift Ingress" <external>
=> 203.0.113.19
$ openstack floating ip create --description "bootstrap machine" <external>
=> 203.0.113.20
```

The OpenShift API (for the OpenShift administrators and app developers) will be at `api.<cluster name>.<cluster domain>` and the Ingress (for the apps' end users) at `*.apps.<cluster name>.<cluster domain>`.

Create these two records in your DNS zone:

```plaintext
api.openshift.example.com.    A 203.0.113.23
*.apps.openshift.example.com. A 203.0.113.19
```

They will need to be available to your developers, end users as well as the OpenShift installer process later in this guide.

## Install Config

Run the `create install-config` subcommand and fill in the desired entries:

```sh
$ openshift-install create install-config
? SSH Public Key </home/user/.ssh/id_rsa.pub>
? Platform <openstack>
? Cloud <openstack>
? ExternalNetwork <external>
? APIFloatingIPAddress <203.0.113.23>
? FlavorName <m1.xlarge>
? Base Domain <example.com>
? Cluster Name <openshift>
```

Most of these are self-explanatory. `Cloud` is the cloud name in your `clouds.yaml` i.e. what's set as your `OS_CLOUD` environment variable.

*Cluster Name* and *Base Domain* will together form the fully qualified domain name which the API interface will expect to the called, and the default name with which OpenShift will expose newly created applications.

Given the values above, the OpenShift API will be available at:

```plaintext
https://api.openshift.example.com:6443/
```

Afterwards, you should have `install-config.yaml` in your current directory:

```sh
$ tree
.
└── install-config.yaml
```

### Configure the machineNetwork.CIDR apiVIP and ingressVIP
The `machineNetwork` represents the OpenStack network which will be used to connect all the OpenShift cluster nodes.
The `machineNetwork.CIDR` defines the IP range, in CIDR notation, from which the installer will choose what IP addresses
to assign the nodes.  The `apiVIP` and `ingressVIP` are the IP addresses the installer will assign to the cluster API and
ingress VIPs, respectively.
In the previous steps, the installer added default values for the `machineNetwork.CIDR`, and then it picked the
5th and 7th IP addresses from that range to assign to `apiVIP` and `ingressVIP`.
`machineNetwork.CIDR` needs to match the IP range specified by `os_subnet_range` in the `inventory.yaml` file.

When the installer creates the manifest files from an existing `install-config.yaml` file, it validates that the
`apiVIP` and `ingressVIP` fall within the IP range specified by `machineNetwork.CIDR`. If they do not, it errors out.
If you change the value of `machineNetwork.CIDR` you must make sure the `apiVIP` and `ingressVIP` values still fall within
the new range. There are two options for setting the `apiVIP` and `ingressVIP`. If you know the values you want to use,
you can specify them in the `install-config.yaml` file. If you want the installer to pick the 5th and 7th IP addresses in the
new range, you need to remove the `apiVIP` and `ingressVIP` entries from the `install-config.yaml` file.

To illustrate the process, we will use '192.0.2.0/24' as an example. It defines a usable IP range from
192.0.2.1 to 192.0.2.254. There are some IP addresses that should be avoided because they are usually taken up or
reserved. For example, the first address (.1) is usually assigned to a router. The DHCP and DNS servers will use a few
more addresses, usually .2, .3, .11 and .12. The actual addresses used by these services depend on the configuration of
the OpenStack deployment in use. You should check your OpenStack deployment.


The following script modifies the value of `machineNetwork.CIDR` in the `install-config.yaml` file.

```sh
$ python -c 'import yaml
path = "install-config.yaml"
data = yaml.safe_load(open(path))
data["networking"]["machineNetwork"][0]["cidr"] = "192.0.2.0/24"
open(path, "w").write(yaml.dump(data, default_flow_style=False))'
```

Next, we need to correct the `apiVIP` and `ingressVIP` values.

The following script will clear the values from the `install-config.yaml` file so that the installer will pick
the 5th and 7th IP addresses in the new range, 192.0.2.5 and 192.0.2.7.

```sh
$ python -c 'import yaml
import sys
path = "install-config.yaml"
data = yaml.safe_load(open(path))
if "apiVIP" in data["platform"]["openstack"]:
   del data["platform"]["openstack"]["apiVIP"]
if "ingressVIP" in data["platform"]["openstack"]:
   del data["platform"]["openstack"]["ingressVIP"]
open(path, "w").write(yaml.dump(data, default_flow_style=False))'
```

If you want to specify the values yourself, you can use the following script, which sets them to 192.0.2.8
and 192.0.2.9.

```sh
$ python -c 'import yaml
import sys
path = "install-config.yaml"
data = yaml.safe_load(open(path))
if "apiVIP" in data["platform"]["openstack"]:
   data["platform"]["openstack"]["apiVIP"] = "192.0.2.8"
if "ingressVIP" in data["platform"]["openstack"]:
   data["platform"]["openstack"]["ingressVIP"] = "192.0.2.9"
open(path, "w").write(yaml.dump(data, default_flow_style=False))'
```

**NOTE**: All the scripts in this guide work with Python 3 as well as Python 2. You can also choose to edit the
`install-config.yaml` file by hand.

### Empty Compute Pools

UPI will not rely on the Machine API for node creation. Instead, we will create the compute nodes ("workers") manually.

We will set their count to `0` in `install-config.yaml`. Look under `compute` -> (first entry) -> `replicas`.

This command will do it for you:

```sh
$ python -c '
import yaml
path = "install-config.yaml"
data = yaml.safe_load(open(path))
data["compute"][0]["replicas"] = 0
open(path, "w").write(yaml.dump(data, default_flow_style=False))'
```

### Modify NetworkType (Required for Kuryr SDN)

By default the `networkType` is set to `OpenShiftSDN` on the `install-config.yaml`.

If an installation with Kuryr is desired, you must modify the `networkType` field.

This command will do it for you:

```sh
$ python -c '
import yaml
path = "install-config.yaml"
data = yaml.safe_load(open(path))
data["networking"]["networkType"] = "Kuryr"
open(path, "w").write(yaml.dump(data, default_flow_style=False))'
```

Also set `os_networking_type` to `Kuryr` in `inventory.yaml`.

## Edit Manifests

We are not relying on the Machine API so we can delete the control plane Machines and compute MachineSets from the manifests.

**WARNING**: The `install-config.yaml` file will be automatically deleted in the next section. If you want to keep it around, copy it elsewhere now!

First, let's turn the install config into manifests:

```sh
$ openshift-install create manifests

$ tree
.
├── manifests
│   ├── 04-openshift-machine-config-operator.yaml
│   ├── cloud-provider-config.yaml
│   ├── cluster-config.yaml
│   ├── cluster-dns-02-config.yml
│   ├── cluster-infrastructure-02-config.yml
│   ├── cluster-ingress-02-config.yml
│   ├── cluster-network-01-crd.yml
│   ├── cluster-network-02-config.yml
│   ├── cluster-proxy-01-config.yaml
│   ├── cluster-scheduler-02-config.yml
│   ├── cvo-overrides.yaml
│   ├── kube-cloud-config.yaml
│   ├── kube-system-configmap-root-ca.yaml
│   ├── machine-config-server-tls-secret.yaml
│   └── openshift-config-secret-pull-secret.yaml
└── openshift
    ├── 99_cloud-creds-secret.yaml
    ├── 99_kubeadmin-password-secret.yaml
    ├── 99_openshift-cluster-api_master-machines-0.yaml
    ├── 99_openshift-cluster-api_master-machines-1.yaml
    ├── 99_openshift-cluster-api_master-machines-2.yaml
    ├── 99_openshift-cluster-api_master-user-data-secret.yaml
    ├── 99_openshift-cluster-api_worker-machineset-0.yaml
    ├── 99_openshift-cluster-api_worker-user-data-secret.yaml
    ├── 99_openshift-machineconfig_master.yaml
    ├── 99_openshift-machineconfig_worker.yaml
    ├── 99_rolebinding-cloud-creds-secret-reader.yaml
    └── 99_role-cloud-creds-secret-reader.yaml

2 directories, 38 files
```

### Remove Machines and MachineSets

Remove the control-plane Machines and compute MachineSets, because we'll be providing those ourselves and don't want to involve the
[machine-API operator][mao]:

```sh
$ rm -f openshift/99_openshift-cluster-api_master-machines-*.yaml openshift/99_openshift-cluster-api_worker-machineset-*.yaml
```
Leave the compute MachineSets in if you want to create compute machines via the machine API. However, some references must be updated in the machineset spec (`openshift/99_openshift-cluster-api_worker-machineset-0.yaml`) to match your environment:

* The OS image: `spec.template.spec.providerSpec.value.image`

[mao]: https://github.com/openshift/machine-api-operator

### Make control-plane nodes unschedulable

Currently [emptying the compute pools][empty-compute-pools] makes control-plane nodes schedulable. But due to a [Kubernetes limitation][kubebug], router pods running on control-plane nodes will not be reachable by the ingress load balancer. Update the scheduler configuration to keep router pods and other workloads off the control-plane nodes:

```sh
$ python -c '
import yaml
path = "manifests/cluster-scheduler-02-config.yml"
data = yaml.safe_load(open(path))
data["spec"]["mastersSchedulable"] = False
open(path, "w").write(yaml.dump(data, default_flow_style=False))'
```

[empty-compute-pools]: #empty-compute-pools
[kubebug]: https://github.com/kubernetes/kubernetes/issues/65618

## Ignition Config

Next, we will turn these manifests into [Ignition][ignition] files. These will be used to configure the Nova servers on boot (Ignition performs a similar function as cloud-init).

```sh
$ openshift-install create ignition-configs
$ tree
.
├── auth
│   ├── kubeadmin-password
│   └── kubeconfig
├── bootstrap.ign
├── master.ign
├── metadata.json
└── worker.ign
```

[ignition]: https://coreos.com/ignition/docs/latest/

### Infra ID

The OpenShift cluster has been assigned an identifier in the form of `<cluster name>-<random string>`. You do not need this for anything, but it is a good idea to keep it around.
You can see the various metadata about your future cluster in `metadata.json`.

The Infra ID is under the `infraID` key:

```sh
$ export INFRA_ID=$(jq -r .infraID metadata.json)
$ echo $INFRA_ID
openshift-qlvwv
```

We'll use the `infraID` as the prefix for all the OpenStack resources we'll create. That way, you'll be able to have multiple deployments in the same OpenStack project without name conflicts.

Make sure your shell session has the `$INFRA_ID` environment variable set when you run the commands later in this document.

### Bootstrap Ignition

#### Edit the Bootstrap Ignition

We need to set the bootstrap hostname explicitly, and in the case of OpenStack using self-signed certificate, the CA cert file. The IPI installer does this automatically, but for now UPI does not.

We will update the ignition file (`bootstrap.ign`) to create the following files:

**`/etc/hostname`**:

```plaintext
openshift-qlvwv-bootstrap
```

(using the `infraID`)

**`/opt/openshift/tls/cloud-ca-cert.pem`** (if applicable).

**NOTE**: We recommend you back up the Ignition files before making any changes!

You can edit the Ignition file manually or run this Python script:

```python
import base64
import json
import os

with open('bootstrap.ign', 'r') as f:
    ignition = json.load(f)

storage = ignition.get('storage', {})
files = storage.get('files', [])

infra_id = os.environ.get('INFRA_ID', 'openshift').encode()
hostname_b64 = base64.standard_b64encode(infra_id + b'-bootstrap\n').decode().strip()
files.append(
{
    'path': '/etc/hostname',
    'mode': 420,
    'contents': {
        'source': 'data:text/plain;charset=utf-8;base64,' + hostname_b64,
    },
})

ca_cert_path = os.environ.get('OS_CACERT', '')
if ca_cert_path:
    with open(ca_cert_path, 'r') as f:
        ca_cert = f.read().encode().strip()
        ca_cert_b64 = base64.standard_b64encode(ca_cert).decode().strip()

    files.append(
    {
        'path': '/opt/openshift/tls/cloud-ca-cert.pem',
        'mode': 420,
        'contents': {
            'source': 'data:text/plain;charset=utf-8;base64,' + ca_cert_b64,
        },
    })

storage['files'] = files
ignition['storage'] = storage

with open('bootstrap.ign', 'w') as f:
    json.dump(ignition, f)
```

Feel free to make any other changes.

#### Upload the Boostrap Ignition

The generated boostrap ignition file tends to be quite large (around 300KB -- it contains all the manifests, master and worker ignitions etc.). This is generally too big to be passed to the server directly (the OpenStack Nova user data limit is 64KB).

To boot it up, we will create a smaller Ignition file that will be passed to Nova as user data and that will download the main ignition file upon execution.

The main file needs to be uploaded to an HTTP(S) location the Bootstrap node will be able to access.

Choose the storage that best fits your needs and availability.

**IMPORTANT**: The `bootstrap.ign` contains sensitive information such as your `clouds.yaml` credentials. It should not be accessible by the public! It will only be used once during the Nova boot of the Bootstrap server. We strongly recommend you restrict the access to that server only and delete the file afterwards.

Possible choices include:

* Swift (see Example 1 below);
* Glance (see Example 2 below);
* Amazon S3;
* Internal web server inside your organization;
* A throwaway Nova server in `$INFRA_ID-nodes` hosting a static web server exposing the file.

In this guide, we will assume the file is at the following URL:

https://static.example.com/bootstrap.ign

##### Example 1: Swift

The `swift` client is needed for enabling listing on the container.
It can be installed by the following command:

```sh
$ sudo dnf install python3-swiftclient
```

Create the `<container_name>` (e.g. $INFRA_ID) container and upload the `bootstrap.ign` file:

```sh
$ openstack container create <container_name> --public
$ openstack object create <container_name> bootstrap.ign
```

Get the `storage_url` from the output:

```sh
$ swift stat -v
```

The URL to be put in the `source` property of the Ignition Shim (see below) will have the following format: `<storage_url>/<container_name>/bootstrap.ign`.

##### Example 2: Glance image service

Create the `<image_name>` image and upload the `bootstrap.ign` file:

```sh
$ openstack image create --disk-format=raw --container-format=bare --file bootstrap.ign <image_name>
```

**NOTE**: Make sure the created image has `active` status.

Copy and save `file` value of the output, it should look like `/v2/images/<image_id>/file`.

Get Glance public URL:

```sh
$ openstack catalog show image
```

By default Glance service doesn't allow anonymous access to the data. So, if you use Glance to store the ignition config, then you also need to provide a valid auth token in the `ignition.config.merge.httpHeaders` field.

The token can be obtained with this command:

```sh
openstack token issue -c id -f value
```

Note that this token can be generated as any OpenStack user with Glance read access; this particular token will only be used for downloading the Ignition file.

The command will return the token to be added to the `ignition.config.merge[0].httpHeaders` property in the Bootstrap Ignition Shim (see [below](#create-the-bootstrap-ignition-shim)):

```json
"httpHeaders": [
	{
		"name": "X-Auth-Token",
		"value": "<token>"
	}
]
```

Combine the public URL with the `file` value to get the link to your bootstrap ignition, in the format `<glance_public_url>/v2/images/<image_id>/file`.

Example of the link to be put in the `source` property of the Ignition Shim (see below): `https://public.glance.example.com:9292/v2/images/b7e2b84e-15cf-440a-a113-3197518da024/file`.

### Create the Bootstrap Ignition Shim

As mentioned before due to Nova user data size limit, we will need to create a new Ignition file that will load the bulk of the Bootstrap node configuration. This will be similar to the existing `master.ign` and `worker.ign` files.

Create a file called `$INFRA_ID-bootstrap-ignition.json` (fill in your `infraID`) with the following contents:

```json
{
  "ignition": {
    "config": {
      "merge": [
        {
          "source": "https://static.example.com/bootstrap.ign",
        }
      ]
    },
    "version": "3.1.0"
  }
}
```

Change the `ignition.config.merge.source` field to the URL hosting the `bootstrap.ign` file you've uploaded previously.

#### Ignition file served by server using self-signed certificate

In order for the bootstrap node to retrieve the ignition file when it is served by a server using self-signed certificate, it is necessary to add the CA certificate to the `ignition.security.tls.certificateAuthorities` in the ignition file. Here is how you might do it.

Encode the certificate to base64:
```sh
$ openssl x509 -in cacert.pem | base64 -w0
```

Add the base64-encoded certificate to the ignition shim:
```json
{
  "ignition": {
    "security": {
      "tls": {
        "certificateAuthorities": [
          {
            "source": "data:text/plain;charset=utf-8;base64,<base64_encoded_certificate>",
          }
        ]
      }
    },
    "version": "3.1.0"
  }
}
```

### Master Ignition

Similar to bootstrap, we need to make sure the hostname is set to the expected value (it must match the name of the Nova server exactly).

Since that value will be different for every master node, we will need to create multiple Ignition files: one for every node.

We will deploy three Control plane (master) nodes. Their Ignition configs can be create like so:

```sh
$ for index in $(seq 0 2); do
    MASTER_HOSTNAME="$INFRA_ID-master-$index\n"
    python -c "import base64, json, sys
ignition = json.load(sys.stdin)
storage = ignition.get('storage', {})
files = storage.get('files', [])
files.append({'path': '/etc/hostname', 'mode': 420, 'contents': {'source': 'data:text/plain;charset=utf-8;base64,' + base64.standard_b64encode(b'$MASTER_HOSTNAME').decode().strip()}})
storage['files'] = files
ignition['storage'] = storage
json.dump(ignition, sys.stdout)" <master.ign >"$INFRA_ID-master-$index-ignition.json"
done
```

This should create files `openshift-qlvwv-master-0-ignition.json`, `openshift-qlvwv-master-1-ignition.json` and `openshift-qlvwv-master-2-ignition.json`.

If you look inside, you will see that they contain very little. In fact, most of the master Ignition is served by the Machine Config Server running on the bootstrap node and the masters contain only enough to know where to look for the rest.

You can make your own changes here.

**NOTE**: The worker nodes do not require any changes to their Ignition, but you can make your own by editing `worker.ign`.

## Network Topology

In this section we'll create all the networking pieces necessary to host the OpenShift cluster: security groups, network, subnet, router, ports.

### Security Groups

```sh
$ ansible-playbook -i inventory.yaml security-groups.yaml
```

The playbook creates one Security group for the Control Plane and one for the Compute nodes, then attaches rules for enabling communication between the nodes.
### Network, Subnet and external router

```sh
$ ansible-playbook -i inventory.yaml network.yaml
```

The playbook creates a network and a subnet. The subnet obeys `os_subnet_range`; however the first ten IP addresses are removed from the allocation pool. These addresses will be used for the VRRP addresses managed by keepalived for high availability. For more information, read the [networking infrastructure design document][net-infra].

Outside connectivity will be provided by attaching the floating IP addresses (IPs in the inventory) to the corresponding routers.

[net-infra]: https://github.com/openshift/installer/blob/master/docs/design/openstack/networking-infrastructure.md

### Subnet DNS (optional)

**NOTE**: This step is optional and only necessary if you want to control the default resolvers your Nova servers will use.

During deployment, the OpenShift nodes will need to be able to resolve public name records to download the OpenShift images and so on. They will also need to resolve the OpenStack API endpoint.

The default resolvers are often set up by the OpenStack administrator in Neutron. However, some deployments do not have default DNS servers set, meaning the servers are not able to resolve any records when they boot.

If you are in this situation, you can add resolvers to your Neutron subnet (`openshift-qlvwv-nodes`). These will be put into `/etc/resolv.conf` on your servers post-boot.

For example, if you want to add the following nameservers: `198.51.100.86` and `198.51.100.87`, you can run this command:

```sh
$ openstack subnet set --dns-nameserver <198.51.100.86> --dns-nameserver <198.51.100.87> "$INFRA_ID-nodes"
```

## Bootstrap

```sh
$ ansible-playbook -i inventory.yaml bootstrap.yaml
```

The playbook sets the *allowed address pairs* on each port attached to our OpenShift nodes.

Since the keepalived-managed IP addresses are not attached to any specific server, Neutron would block their traffic by default. By passing them to `--allowed-address` the traffic can flow freely through.

An additional Floating IP is also attached to the bootstrap port. This is not necessary for the deployment (and we will delete the bootstrap resources afterwards). However, if the bootstrapping phase fails for any reason, the installer will try to SSH in and download the bootstrap log. That will only succeed if the node is reachable (which in general means a floating IP).

After the bootstrap server is active, you can check the console log to see that it is getting the ignition correctly:

```sh
$ openstack console log show "$INFRA_ID-bootstrap"
```

You can also SSH into the server (using its floating IP address) and check on the bootstrapping progress:

```sh
$ ssh core@203.0.113.24
[core@openshift-qlvwv-bootstrap ~]$ journalctl -b -f -u bootkube.service
```

## Control Plane

```sh
$ ansible-playbook -i inventory.yaml control-plane.yaml
```

Our control plane will consist of three nodes. The servers will be passed the `master-?-ignition.json` files prepared earlier.

The playbook places the Control Plane in a Server Group with "soft anti-affinity" policy.

The master nodes should load the initial Ignition and then keep waiting until the bootstrap node stands up the Machine Config Server which will provide the rest of the configuration.

### Control Plane Trunks (Kuryr SDN)

If `os_networking_type` is set to `Kuryr` in the Ansible inventory, the playbook creates the Trunks for Kuryr to plug the containers into the OpenStack SDN.

### Wait for the Control Plane to Complete

When that happens, the masters will start running their own pods, run etcd and join the "bootstrap" cluster. Eventually, they will form a fully operational control plane.

You can monitor this via the following command:

```sh
$ openshift-install wait-for bootstrap-complete
```

Eventually, it should output the following:

```plaintext
INFO API v1.14.6+f9b5405 up
INFO Waiting up to 30m0s for bootstrapping to complete...
```

This means the masters have come up successfully and are joining the cluster.

Eventually, the `wait-for` command should end with:

```plaintext
INFO It is now safe to remove the bootstrap resources
```

### Access the OpenShift API

You can use the `oc` or `kubectl` commands to talk to the OpenShift API. The admin credentials are in `auth/kubeconfig`:

```sh
$ export KUBECONFIG="$PWD/auth/kubeconfig"
$ oc get nodes
$ oc get pods -A
```

**NOTE**: Only the API will be up at this point. The OpenShift UI will run on the compute nodes.

### Delete the Bootstrap Resources

```sh
$ ansible-playbook -i inventory.yaml down-bootstrap.yaml
```

The teardown playbook deletes the bootstrap port and server.

Now the bootstrap floating IP can also be destroyed.

If you haven't done so already, you should also disable the bootstrap Ignition URL.

## Compute Nodes

```sh
$ ansible-playbook -i inventory.yaml compute-nodes.yaml
```

This process is similar to the masters, but the workers need to be approved before they're allowed to join the cluster.

The workers need no ignition override.

### Compute Nodes Trunks (Kuryr SDN)

If `os_networking_type` is set to `Kuryr` in the Ansible inventory, the playbook creates the Trunks for Kuryr to plug the containers into the OpenStack SDN.

### Compute Nodes with SR-IOV NICs

Using single root I/O virtualization (SR-IOV) networking as an additional network in OpenShift can be beneficial for applications that require high bandwidth and low latency. To enable this in your cluster, you will need to install the [SR-IOV Network Operator](https://docs.openshift.com/container-platform/4.6/networking/hardware_networks/installing-sriov-operator.html). If you are not sure whether your cluster supports this feature, please refer to the [SR-IOV hardware networks documentation](https://docs.openshift.com/container-platform/4.6/networking/hardware_networks/about-sriov.html). If you are planning an openstack deployment with SR-IOV networks and need addition resources, check the [OpenStack SR-IOV deployment docs](https://access.redhat.com/documentation/en-us/red_hat_openstack_platform/16.1/html-single/network_functions_virtualization_planning_and_configuration_guide/index#assembly_sriov_parameters). Once you meet these requirements, you can start provisioning an SR-IOV network and subnet in OpenStack.

```sh
openstack network create radio --provider-physical-network radio --provider-network-type vlan --provider-segment 120

openstack subnet create radio --network radio --subnet-range <your CIDR range> --dhcp
```

Your compute nodes will need to have two types of ports for this feature to work. One port needs to connect the node to your OpenShift network so that it can join the cluster and communicate with the other nodes. The other type of port is for your SR-IOV traffic. The OpenShift networking port should be created the same way we normally create ports for compute nodes.

```sh
openstack port create os_port_worker_0 --network <infraID>-network --security-group <infraID>-worker --fixed-ip subnet=<infraID>-nodes,ip-address=<a fixed IP> --allowed-address ip-address=<infraID>-ingress-port
```

The SR-IOV port(s) must be created explicitly by the user and passed as a NIC during instance creation, otherwise the `vnic-type` will not be `direct` and it will not work.

```sh
openstack port create radio_port --vnic-type direct --network radio --fixed-ip subnet=radio,ip-address=<a fixed ip> --tag=radio --disable-port-security
```

When you create your instance, make sure that the SR-IOV port and the OCP port you created for it are added as NICs.

```sh
openstack server create --image <infraID>-rhcos --flavor ocp --user-data <ocp project>/build-artifacts/worker.ign --nic port-id=<os_port_worker_0 ID> --nic port-id=<radio_port ID> --config-drive true worker-<worker_id>.<cluster_name>.<cluster_domain>
```

### Approve the worker CSRs

Even after they've booted up, the workers will not show up in `oc get nodes`.

Instead, they will create certificate signing requests (CSRs) which need to be approved. You can watch for the CSRs here:

```sh
$ watch oc get csr -A
```

Eventually, you should see `Pending` entries looking like this

```sh
$ oc get csr -A
NAME        AGE    REQUESTOR                                                                   CONDITION
csr-2scwb   16m    system:serviceaccount:openshift-machine-config-operator:node-bootstrapper   Approved,Issued
csr-5jwqf   16m    system:node:openshift-qlvwv-master-0                                        Approved,Issued
csr-88jp8   116s   system:serviceaccount:openshift-machine-config-operator:node-bootstrapper   Pending
csr-9dt8f   15m    system:node:openshift-qlvwv-master-1                                        Approved,Issued
csr-bqkw5   16m    system:serviceaccount:openshift-machine-config-operator:node-bootstrapper   Approved,Issued
csr-dpprd   6s     system:serviceaccount:openshift-machine-config-operator:node-bootstrapper   Pending
csr-dtcws   24s    system:serviceaccount:openshift-machine-config-operator:node-bootstrapper   Pending
csr-lj7f9   16m    system:node:openshift-qlvwv-master-2                                        Approved,Issued
csr-lrtlk   15m    system:serviceaccount:openshift-machine-config-operator:node-bootstrapper   Approved,Issued
csr-wkm94   16m    system:serviceaccount:openshift-machine-config-operator:node-bootstrapper   Approved,Issued
```

You should inspect each pending CSR and verify that it comes from a node you recognize:

```sh
$ oc describe csr csr-88jp8
Name:               csr-88jp8
Labels:             <none>
Annotations:        <none>
CreationTimestamp:  Wed, 23 Oct 2019 13:22:51 +0200
Requesting User:    system:serviceaccount:openshift-machine-config-operator:node-bootstrapper
Status:             Pending
Subject:
         Common Name:    system:node:openshift-qlvwv-worker-0
         Serial Number:
         Organization:   system:nodes
Events:  <none>
```

If it does (this one is for `openshift-qlvwv-worker-0` which we've created earlier), you can approve it:

```sh
$ oc adm certificate approve csr-88jp8
```

Approved nodes should now show up in `oc get nodes`, but they will be in the `NotReady` state. They will create a second CSR which you should also review:

```sh
$ oc get csr -A
NAME        AGE     REQUESTOR                                                                   CONDITION
csr-2scwb   17m     system:serviceaccount:openshift-machine-config-operator:node-bootstrapper   Approved,Issued
csr-5jwqf   17m     system:node:openshift-qlvwv-master-0                                         Approved,Issued
csr-7mv4d   13s     system:node:openshift-qlvwv-worker-1                                         Pending
csr-88jp8   3m29s   system:serviceaccount:openshift-machine-config-operator:node-bootstrapper   Approved,Issued
csr-9dt8f   17m     system:node:openshift-qlvwv-master-1                                         Approved,Issued
csr-bqkw5   18m     system:serviceaccount:openshift-machine-config-operator:node-bootstrapper   Approved,Issued
csr-bx7p4   28s     system:node:openshift-qlvwv-worker-0                                         Pending
csr-dpprd   99s     system:serviceaccount:openshift-machine-config-operator:node-bootstrapper   Approved,Issued
csr-dtcws   117s    system:serviceaccount:openshift-machine-config-operator:node-bootstrapper   Approved,Issued
csr-lj7f9   17m     system:node:openshift-qlvwv-master-2                                         Approved,Issued
csr-lrtlk   17m     system:serviceaccount:openshift-machine-config-operator:node-bootstrapper   Approved,Issued
csr-wkm94   18m     system:serviceaccount:openshift-machine-config-operator:node-bootstrapper   Approved,Issued
csr-wqpfd   21s     system:node:openshift-qlvwv-worker-2                                         Pending
```

(we see the CSR approved earlier as well as a new `Pending` one for the same node: `openshift-qlvwv-worker-0`)

And approve:

```sh
$ oc adm certificate approve csr-bx7p4
```

Once this CSR is approved, the node should switch to `Ready` and pods will be scheduled on it.

### Wait for the OpenShift Installation to Complete

Run the following command to verify the OpenShift cluster is fully deployed:

```sh
$ openshift-install --log-level debug wait-for install-complete
```

Upon success, it will print the URL to the OpenShift Console (the web UI) as well as admin username and password to log in.

## Destroy the OpenShift Cluster

```sh
$ ansible-playbook -i inventory.yaml  \
	down-bootstrap.yaml      \
	down-control-plane.yaml  \
	down-compute-nodes.yaml  \
	down-load-balancers.yaml \
	down-containers.yaml     \
	down-network.yaml        \
	down-security-groups.yaml
```

The playbook `down-load-balancers.yaml` idempotently deletes the load balancers created by the Kuryr installation, if any.

**NOTE:** The deletion of load balancers with `provisioning_status` `PENDING-*` is skipped. Make sure to retry the
`down-load-balancers.yaml` playbook once the load balancers have transitioned to `ACTIVE`.

Then, remove the `api` and `*.apps` DNS records.

The floating IPs can also be deleted if not useful any more.
