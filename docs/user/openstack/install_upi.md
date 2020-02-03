# Installing OpenShift on OpenStack User-Provisioned Infrastructure

The User-Provisioned Infrastructure (UPI) process installs OpenShift in stages, providing opportunities for modifications or integrating with existing infrastructure.

It contrasts with the fully-automated Installer-Provisioned Infrastructure (IPI) which creates everything in one go.

With UPI, creating the cloud (OpenStack) resources (e.g. Nova servers, Neutron ports, security groups) is the responsibility of the person deploying OpenShift.

The installer is still used to generate the ignition files and monitor the installation process.

This provides a greater flexibility at the cost of a more explicit and interactive process.


## Prerequisites

The requirements for UPI are broadly similar to the [ones for OpenStack IPI][ipi-reqs]:

[ipi-reqs]: ./README.md#openstack-requirements

- OpenStack account with `clouds.yaml`
  - *This document* will use a cloud called `openstack`
- Nova flavors
  - *This document* will use `m1.xlarge` for masters and `m1.large` for workers
- An external subnet you want to use for floating IP addresses
  - *This document* will use `external`
- The `openshift-install` binary in your `$PATH`
- The [Red Hat CoreOS][rhcos] image in Glance
  - *This document* will use `rhcos` as the image name
- A subnet range for the Nova servers / OpenShift Nodes
  - This range must not conflict with your existing network
  - *This document* will use `192.0.2.0/24`
  - **WARNING**: `192.0.2.0/24` is an IP block range reserved for documentation (in [RFC 5737](https://tools.ietf.org/html/rfc5737))
  - Traffic to this range should be blocked by firewalls etc.
  - You must pick your own range -- one that is routable
- A cluster name you will want to use
  - *This document* will use `openshift`
- A base domain
  - *This document* will use `example.com`
- OpenShift Pull Secret
- DNS infrastructure you can configure (you'll need to add two records there)
- A DNS zone you can configure
  - It must be resolvable by the installer and end-user machines
  - You will need to create two DNS records there (for API and apps access)
  - See the [Create Public DNS Records][dns-details] section for more details

[rhcos]: https://www.openshift.com/learn/coreos/
[dns-details]: #create-public-dns-records

The OpenShift API URL will be generated from the cluster name and base domain. E.g.: `https://api.openshift.example.com:6443/`


You can validate most of the above by running the following commands:

```sh
$ export OS_CLOUD=openstack
$ openstack image show rhcos
$ openstack network show external
$ openstack flavor show m1.xlarge
$ openstack flavor show m1.large
$ which openshift-install
```

They should all succeed.

For an installation with Kuryr SDN on UPI, you should also check the requirements which are the same
needed for [OpenStack IPI with Kuryr][ipi-reqs-kuryr]. Please also note that **RHEL 7 nodes are not 
supported on deployments configured with Kuryr**. This is because Kuryr container images are based on
RHEL 8 and may not work properly when run on RHEL 7.

[ipi-reqs-kuryr]: ./kuryr.md#requirements-when-enabling-kuryr


## OpenShift Configuration Directory

All the configuration files, logs and installation state are kept in a single directory:

```sh
$ mkdir -p openstack-upi
$ cd openstack-upi
```

## API and Ingress Floating IP Addresses

```sh
$ openstack floating ip create --description "OpenShift API" <external>
=> 203.0.113.23
$ openstack floating ip create --description "OpenShift Ingress" <external>
=> 203.0.113.19
```

**NOTE**: throughout this document, we will use `203.0.113.23` as the public IP address for the OpenShift API endpoint and `203.0.113.19` as the public IP for the ingress (`*.apps`) endpoint.

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

```
https://api.openshift.example.com:6443/
```

Afterwards, you should have `install-config.yaml` in your current directory:

```sh
$ tree
.
└── install-config.yaml
```

### Fix the Node Subnet

The installer added a default IP range for the OpenShift nodes. It must match the range for the Neutron subnet you'll create later on.

We're going to use a custom subnet to illustrate how that can be done.


Our range will be `192.0.2.0/24` so we need to add that value to
`install-config.yaml`. Look under `networking` -> `machineCIDR`.

This command will do it for you:

```sh
$ python3 -c 'import yaml;
path = "install-config.yaml";
data = yaml.safe_load(open(path));
data["networking"]["machineCIDR"] = "192.0.2.0/24";
open(path, "w").write(yaml.dump(data, default_flow_style=False))'
```

**NOTE**: All the scripts in this guide work with Python 3 as well as Python 2.

### Empty Compute Pools

UPI will not rely on the Machine API for node creation. Instead, we will create the compute nodes ("workers") manually.

We will set their count to `0` in `install-config.yaml`. Look under `compute` -> (first entry) -> `replicas`.

This command will do it for you:

```sh
$ python3 -c '
import yaml;
path = "install-config.yaml";
data = yaml.safe_load(open(path));
data["compute"][0]["replicas"] = 0;
open(path, "w").write(yaml.dump(data, default_flow_style=False))'
```

### Modify NetworkType (Required for Kuryr SDN)

By default the `networkType` is set to `OpenShiftSDN` on the `install-config.yaml`.

If an installation with Kuryr is desired, you must modify the `networkType` field.

This command will do it for you:

```sh
$ python3 -c '
import yaml;
path = "install-config.yaml";
data = yaml.safe_load(open(path));
data["networking"]["networkType"] = "Kuryr";
open(path, "w").write(yaml.dump(data, default_flow_style=False))'
```

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
│   ├── etcd-ca-bundle-configmap.yaml
│   ├── etcd-client-secret.yaml
│   ├── etcd-host-service-endpoints.yaml
│   ├── etcd-host-service.yaml
│   ├── etcd-metric-client-secret.yaml
│   ├── etcd-metric-serving-ca-configmap.yaml
│   ├── etcd-metric-signer-secret.yaml
│   ├── etcd-namespace.yaml
│   ├── etcd-service.yaml
│   ├── etcd-serving-ca-configmap.yaml
│   ├── etcd-signer-secret.yaml
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

[mao]: https://github.com/openshift/machine-api-operator

```sh
$ rm -f openshift/99_openshift-cluster-api_master-machines-*.yaml openshift/99_openshift-cluster-api_worker-machineset-*.yaml
```
You are free to leave the compute MachineSets in if you want to create compute machines via the machine API, but if you do you may need to update the various references (`subnet`, etc.) to match your environment.

### Make control-plane nodes unschedulable

Currently [emptying the compute pools][empty-compute-pools] makes control-plane nodes schedulable. But due to a [Kubernetes limitation][kubebug], router pods running on control-plane nodes will not be reachable by the ingress load balancer. Update the scheduler configuration to keep router pods and other workloads off the control-plane nodes:

[empty-compute-pools]: #empty-compute-pools
[kubebug]: https://github.com/kubernetes/kubernetes/issues/65618

```sh
$ python3 -c '
import yaml;
path = "manifests/cluster-scheduler-02-config.yml"
data = yaml.safe_load(open(path));
data["spec"]["mastersSchedulable"] = False;
open(path, "w").write(yaml.dump(data, default_flow_style=False))'
```


## Ignition Config

Next, we will turn these manifests into [Ignition][ignition] files. These will be used to configure the Nova servers on boot (Ignition performs a similar function as cloud-init).

[ignition]: https://coreos.com/ignition/docs/latest/

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


### Update Bootstrap Ignition

There is some bootstrap configuration we weed to add explicitly to the ignition file. The IPI installer does this automatically via Terraform, but the changes do not show up in the files generated by `create ignition-configs`. This section should go away in the future.

The contents of the files we need to create have to be base64-encoded.

We will create three files:

**`/etc/hostname`**:

```
openshift-qlvwv-bootstrap
```

(using the `infraID`)

**`/etc/NetworkManager/conf.d/dhcp-client.conf`**:

```
[main]
dhcp=dhclient
```

**`/etc/dhcp/dhclient.conf`**:

```
send dhcp-client-identifier = hardware;
prepend domain-name-servers 127.0.0.1;
```

**NOTE**: We recommend you back up the Ignition files before making any changes!

You can edit the Ignition file manually or run this Python script:

```python
import base64
import json
import os

with open('bootstrap.ign', 'r') as f:
    ignition = json.load(f)

files = ignition['storage'].get('files', [])

infra_id = os.environ.get('INFRA_ID', 'openshift').encode()
hostname_b64 = base64.standard_b64encode(infra_id + b'-bootstrap\n').decode().strip()
files.append(
{
    'path': '/etc/hostname',
    'mode': 420,
    'contents': {
        'source': 'data:text/plain;charset=utf-8;base64,' + hostname_b64,
        'verification': {}
    },
    'filesystem': 'root',
})

dhcp_client_conf_b64 = base64.standard_b64encode(b'[main]\ndhcp=dhclient\n').decode().strip()
files.append(
{
    'path': '/etc/NetworkManager/conf.d/dhcp-client.conf',
    'mode': 420,
    'contents': {
        'source': 'data:text/plain;charset=utf-8;base64,' + dhcp_client_conf_b64,
        'verification': {}
        },
    'filesystem': 'root',
})

dhclient_cont_b64 = base64.standard_b64encode(b'send dhcp-client-identifier = hardware;\nprepend domain-name-servers 127.0.0.1;\n').decode().strip()
files.append(
{
    'path': '/etc/dhcp/dhclient.conf',
    'mode': 420,
    'contents': {
        'source': 'data:text/plain;charset=utf-8;base64,' + dhclient_cont_b64,
        'verification': {}
        },
    'filesystem': 'root'
})

ignition['storage']['files'] = files;

with open('bootstrap.ign', 'w') as f:
    json.dump(ignition, f)
```

Feel free to make any other changes.

### Update Master Ignition

Similar to bootstrap, we need to make sure the hostname is set to the expected value (it must match the name of the Nova server exactly).

Since that value will be different for every master node, we will need to create multiple Ignition files: one for every node.

We will deploy three Control plane (master) nodes. Their Ignition configs can be create like so:

```sh
$ for index in $(seq 0 2); do
    MASTER_HOSTNAME="$INFRA_ID-master-$index\n"
    python3 -c "import base64, json, sys;
ignition = json.load(sys.stdin);
files = ignition['storage'].get('files', []);
files.append({'path': '/etc/hostname', 'mode': 420, 'contents': {'source': 'data:text/plain;charset=utf-8;base64,' + base64.standard_b64encode(b'$MASTER_HOSTNAME').decode().strip(), 'verification': {}}, 'filesystem': 'root'});
ignition['storage']['files'] = files;
json.dump(ignition, sys.stdout)" <master.ign >"$INFRA_ID-master-$index-ignition.json"
done
```

This should create files `openshift-qlvwv-master-0-ignition.json`, `openshift-qlvwv-master-1-ignition.json` and `openshift-qlvwv-master-2-ignition.json`.

If you look inside, you will see that they contain very little. In fact, most of the master Ignition is served by the Machine Config Server running on the bootstrap node and the masters contain only enough to know where to look for the rest.

You can make your own changes here.

**NOTE**: The worker nodes do not require any changes to their Ignition, but you can make your own by editing `worker.ign`.


## Network Topology

In this section we'll create all the networking pieces necessary to host the OpenShift cluster: network, subnet, router etc.


### Network and Subnet

We will use the `192.0.2.0/24` subnet range, but remove the first ten IP addresses from the allocation pool.

These addresses will be used for the VRRP addresses managed by keepalived for high availability. For more information, read the [networking infrastructure design document][net-infra].

[net-infra]: https://github.com/openshift/installer/blob/master/docs/design/openstack/networking-infrastructure.md

```sh
$ openstack network create "$INFRA_ID-network" --tag openshiftClusterID="$INFRA_ID"
$ openstack subnet create --subnet-range <192.0.2.0/24> --allocation-pool start=192.0.2.10,end=192.0.2.254 --network "$INFRA_ID-network" --tag openshiftClusterID="$INFRA_ID" "$INFRA_ID-nodes"
```

### Subnet DNS (optional)

During deployment, the OpenShift nodes will need to be able to resolve public name records to download the OpenShift images and so on. They will also need to resolve the OpenStack API endpoint.

The default resolvers are often set up by the OpenStack administrator in Neutron. However, some deployments do not have default DNS servers set, meaning the servers are not able to resolve any records when they boot.

If you are in this situation, you can add resolvers to your Neutron subnet (`openshift-qlvwv-nodes`). These will be put into `/etc/resolv.conf` on your servers post-boot.

For example, if you want to add the following nameservers: `198.51.100.86` and `198.51.100.87`, you can run this command:

```sh
$ openstack subnet set --dns-nameserver <198.51.100.86> --dns-nameserver <198.51.100.87> "$INFRA_ID-nodes"
```

**NOTE**: This step is optional and only necessary if you want to control the default resolvers your Nova servers will use.

### Security Groups

We will need two security groups: one for the master nodes (the control plane) and a separate one for the worker (compute) nodes.


#### Master and Worker Security Groups

```sh
$ openstack security group create "$INFRA_ID-master" --tag openshiftClusterID="$INFRA_ID"
$ openstack security group create "$INFRA_ID-worker" --tag openshiftClusterID="$INFRA_ID"
```
**NOTE**: Tagging security groups is supported from 3.16.0 version of OpenStackClient

#### Master Security Group Rules

```sh
openstack security group rule create --description "ICMP" --protocol icmp "$INFRA_ID-master"
openstack security group rule create --description "machine config server" --protocol tcp --dst-port 22623 --remote-ip 192.0.2.0/24 "$INFRA_ID-master"
openstack security group rule create --description "SSH" --protocol tcp --dst-port 22 "$INFRA_ID-master"
openstack security group rule create --description "DNS (TCP)" --protocol tcp --dst-port 53 --remote-ip 192.0.2.0/24 "$INFRA_ID-master"
openstack security group rule create --description "DNS (UDP)" --protocol udp --dst-port 53 --remote-ip 192.0.2.0/24 "$INFRA_ID-master"
openstack security group rule create --description "mDNS" --protocol udp --dst-port 5353 --remote-ip 192.0.2.0/24 "$INFRA_ID-master"
openstack security group rule create --description "OpenShift API" --protocol tcp --dst-port 6443 "$INFRA_ID-master"
openstack security group rule create --description "VXLAN" --protocol udp --dst-port 4789 --remote-group "$INFRA_ID-master" "$INFRA_ID-master"
openstack security group rule create --description "VXLAN from worker" --protocol udp --dst-port 4789 --remote-group "$INFRA_ID-worker" "$INFRA_ID-master"
openstack security group rule create --description "Geneve" --protocol udp --dst-port 6081 --remote-group "$INFRA_ID-master" "$INFRA_ID-master"
openstack security group rule create --description "Geneve from worker" --protocol udp --dst-port 6081 --remote-group "$INFRA_ID-worker" "$INFRA_ID-master"
openstack security group rule create --description "ovndb" --protocol tcp --dst-port 6641:6642 --remote-group "$INFRA_ID-master" "$INFRA_ID-master"
openstack security group rule create --description "ovndb from worker" --protocol tcp --dst-port 6641:6642 --remote-group "$INFRA_ID-worker" "$INFRA_ID-master"
openstack security group rule create --description "master ingress internal (TCP)" --protocol tcp --dst-port 9000:9999 --remote-group "$INFRA_ID-master" "$INFRA_ID-master"
openstack security group rule create --description "master ingress internal from worker (TCP)" --protocol tcp --dst-port 9000:9999 --remote-group "$INFRA_ID-worker" "$INFRA_ID-master"
openstack security group rule create --description "master ingress internal (UDP)" --protocol udp --dst-port 9000:9999 --remote-group "$INFRA_ID-master" "$INFRA_ID-master"
openstack security group rule create --description "master ingress internal from worker (UDP)" --protocol udp --dst-port 9000:9999 --remote-group "$INFRA_ID-worker" "$INFRA_ID-master"
openstack security group rule create --description "kube scheduler" --protocol tcp --dst-port 10259 --remote-group "$INFRA_ID-master" "$INFRA_ID-master"
openstack security group rule create --description "kube scheduler from worker" --protocol tcp --dst-port 10259 --remote-group "$INFRA_ID-worker" "$INFRA_ID-master"
openstack security group rule create --description "kube controller manager" --protocol tcp --dst-port 10257 --remote-group "$INFRA_ID-master" "$INFRA_ID-master"
openstack security group rule create --description "kube controller manager from worker" --protocol tcp --dst-port 10257 --remote-group "$INFRA_ID-worker" "$INFRA_ID-master"
openstack security group rule create --description "master ingress kubelet secure" --protocol tcp --dst-port 10250 --remote-group "$INFRA_ID-master" "$INFRA_ID-master"
openstack security group rule create --description "master ingress kubelet secure from worker" --protocol tcp --dst-port 10250 --remote-group "$INFRA_ID-worker" "$INFRA_ID-master"
openstack security group rule create --description "etcd" --protocol tcp --dst-port 2379:2380 --remote-group "$INFRA_ID-master" "$INFRA_ID-master"
openstack security group rule create --description "master ingress services (TCP)" --protocol tcp --dst-port 30000:32767 --remote-group "$INFRA_ID-master" "$INFRA_ID-master"
openstack security group rule create --description "master ingress services (TCP) from worker" --protocol tcp --dst-port 30000:32767 --remote-group "$INFRA_ID-worker" "$INFRA_ID-master"
openstack security group rule create --description "master ingress services (UDP)" --protocol udp --dst-port 30000:32767 --remote-group "$INFRA_ID-master" "$INFRA_ID-master"
openstack security group rule create --description "master ingress services (UDP) from worker" --protocol udp --dst-port 30000:32767 --remote-group "$INFRA_ID-worker" "$INFRA_ID-master"
openstack security group rule create --description "VRRP" --protocol vrrp --remote-ip 192.0.2.0/24 "$INFRA_ID-master"
```


#### Worker Security Group Rules

```sh
openstack security group rule create --description "ICMP" --protocol icmp "$INFRA_ID-worker"
openstack security group rule create --description "SSH" --protocol tcp --dst-port 22 "$INFRA_ID-worker"
openstack security group rule create --description "mDNS" --protocol udp --dst-port 5353 --remote-ip 192.0.2.0/24 "$INFRA_ID-worker"
openstack security group rule create --description "Ingress HTTP" --protocol tcp --dst-port 80 "$INFRA_ID-worker"
openstack security group rule create --description "Ingress HTTPS" --protocol tcp --dst-port 443 "$INFRA_ID-worker"
openstack security group rule create --description "router" --protocol tcp --dst-port 1936 --remote-ip 192.0.2.0/24 "$INFRA_ID-worker"
openstack security group rule create --description "VXLAN" --protocol udp --dst-port 4789 --remote-group "$INFRA_ID-worker" "$INFRA_ID-worker"
openstack security group rule create --description "VXLAN from master" --protocol udp --dst-port 4789 --remote-group "$INFRA_ID-master" "$INFRA_ID-worker"
openstack security group rule create --description "Geneve" --protocol udp --dst-port 6081 --remote-group "$INFRA_ID-worker" "$INFRA_ID-worker"
openstack security group rule create --description "Geneve from master" --protocol udp --dst-port 6081 --remote-group "$INFRA_ID-master" "$INFRA_ID-worker"
openstack security group rule create --description "worker ingress internal (TCP)" --protocol tcp --dst-port 9000:9999 --remote-group "$INFRA_ID-worker" "$INFRA_ID-worker"
openstack security group rule create --description "worker ingress internal from master (TCP)" --protocol tcp --dst-port 9000:9999 --remote-group "$INFRA_ID-master" "$INFRA_ID-worker"
openstack security group rule create --description "worker ingress internal (UDP)" --protocol udp --dst-port 9000:9999 --remote-group "$INFRA_ID-worker" "$INFRA_ID-worker"
openstack security group rule create --description "worker ingress internal from master (UDP)" --protocol udp --dst-port 9000:9999 --remote-group "$INFRA_ID-master" "$INFRA_ID-worker"
openstack security group rule create --description "worker ingress kubelet secure" --protocol tcp --dst-port 10250 --remote-group "$INFRA_ID-worker" "$INFRA_ID-worker"
openstack security group rule create --description "worker ingress kubelet secure from master" --protocol tcp --dst-port 10250 --remote-group "$INFRA_ID-master" "$INFRA_ID-worker"
openstack security group rule create --description "worker ingress services (TCP)" --protocol tcp --dst-port 30000:32767 --remote-group "$INFRA_ID-worker" "$INFRA_ID-worker"
openstack security group rule create --description "worker ingress services (TCP) from master" --protocol tcp --dst-port 30000:32767 --remote-group "$INFRA_ID-master" "$INFRA_ID-worker"
openstack security group rule create --description "worker ingress services (UDP)" --protocol udp --dst-port 30000:32767 --remote-group "$INFRA_ID-worker" "$INFRA_ID-worker"
openstack security group rule create --description "worker ingress services (UDP) from master" --protocol udp --dst-port 30000:32767 --remote-group "$INFRA_ID-master" "$INFRA_ID-worker"
openstack security group rule create --description "VRRP" --protocol vrrp --remote-ip 192.0.2.0/24 "$INFRA_ID-worker"
```


### Router

The outside connectivity will be provided by floating IP addresses. Your subnet needs to have a router set up for this to work.

```sh
$ openstack router create "$INFRA_ID-external-router" --tag openshiftClusterID="$INFRA_ID"
$ openstack router set --external-gateway <external> "$INFRA_ID-external-router"
$ openstack router add subnet "$INFRA_ID-external-router" "$INFRA_ID-nodes"
```

**NOTE**: Pass in the same external network name (`external` in the `router create` command above) that you used in the install config and that you used in `openstack floating ip create`.


### Public API and Ingress Access

To provide access to the OpenShift cluster, we will need to create two ports, attach our API and Ingress Floating IPs and publish the appropriate DNS records.

**IMPORTANT:** these ports will need to have specific fixed (i.e. private) IP addresses. They have to be the fifth and seventh addresses from the subnet range. For example:

#### API and Ingress Ports

```sh
$ openstack port create --network "$INFRA_ID-network" --security-group "$INFRA_ID-master" --fixed-ip "subnet=$INFRA_ID-nodes,ip-address=192.0.2.5" --tag openshiftClusterID="$INFRA_ID" "$INFRA_ID-api-port"
$ openstack port create --network "$INFRA_ID-network" --security-group "$INFRA_ID-worker" --fixed-ip "subnet=$INFRA_ID-nodes,ip-address=192.0.2.7" --tag openshiftClusterID="$INFRA_ID" "$INFRA_ID-ingress-port"
```

The fixed IP addresses might differ if you're using a different subnet range.

These ports will always stay `DOWN`. They will never be attached to any Nova server directly. Rather, their fixed IP addresses will be managed by keepalived and moved around if the backing service goes down, providing high availability.

There is another fixed IP address (`192.0.2.6`) which will be used for internal DNS records necessary to install and run the OpenShift cluster. These do not need to be exposed to the outside (they only contain internal IPs) and therefore do not need their own Neutron port either. This address will be managed by keepalived as well.


#### Assign Floating IPs to the Ports

We have created two floating IP addresses at the beginning. We need to add them to the ports:

```sh
$ openstack floating ip set --port "$INFRA_ID-api-port" <203.0.113.23>
$ openstack floating ip set --port "$INFRA_ID-ingress-port" <203.0.113.19>
```

#### Create Public DNS Records

The OpenShift API (for the OpenShift administrators and app developers) will be at `api.<cluster name>.<cluster domain>` and the Ingress (for the apps' end users) at `*.apps.<cluster name>.<cluster domain>`.

Create these two records:

```
api.openshift.example.com.   A 203.0.113.23
*.apps.openshift.example.com. A 203.0.113.19
```


They will need to be available to your developers, end users as well as the OpenShift installer process later in this guide.



## Bootstrap

### Bootstrap Ignition

The generated boostrap ignition file (`bootstrap.ign`) tends to be quite large (around 300KB -- it contains all the manifests, master and worker ignitions etc.). This is generally too big to be passed to the server directly (the OpenStack Nova user data limit is 64KB).

To boot it up, we will need to create a smaller Ignition file that will be passed to Nova as user data and that will download the main ignition file upon execution.

The main file needs to be uploaded to an HTTP(S) location the Bootstrap node will be able to access.

You are free to choose any storage you want. For example:

* Swift Object Storage
  Create the `<container_name>` container and upload the `bootstrap.ign` file:
  ```sh
  $ swift upload <container_name> bootstrap.ign
  ```
  Make the container accessible:
  ```sh
  $ swift post <container_name> --read-acl ".r:*,.rlistings"
  ```
  Get the `StorageURL` from the output:
  ```sh
  $ swift stat -v
  ```
* Amazon S3
* Internal HTTP server inside your organisation
* A short-lived Nova server in `$INFRA_ID-nodes` hosting the file for bootstrapping

In this guide, we will assume the file is at the following URL:

https://static.example.com/bootstrap.ign

**NOTE**: In case the Swift object storage option was chosen the URL will have the following format: `<StorageURL>/<container_name>/bootstrap.ign`

**IMPORTANT**: The `bootstrap.ign` contains sensitive information such as your `clouds.yaml` credentials and TLS certificates. It should **not** be accessible to the public! It will only be used once during the Nova boot of the Bootstrap server. We strongly recommend
you restrict the access to that server only and delete the file afterwards.


### Bootstrap Ignition Shim

As mentioned before due to Nova user data size limit, we will need to create a new Ignition file that will load the bulk of the Bootstrap node configuration. This will be similar to the existing `master.ign` and `worker.ign` files.

Create a file called `$INFRA_ID-bootstrap-ignition.json` (fill in your `infraID`) with the following contents:

```json
{
  "ignition": {
    "config": {
      "append": [
        {
          "source": "https://static.example.com/bootstrap.ign",
          "verification": {}
        }
      ]
    },
    "security": {},
    "timeouts": {},
    "version": "2.2.0"
  },
  "networkd": {},
  "passwd": {},
  "storage": {},
  "systemd": {}
}
```

Change the `ignition.config.append.source` field to the URL hosting the `bootstrap.ign` file you've uploaded previously.


### Bootstrap Port

Generally, it's not necessary to create a port explicitly -- `openstack server create` can do it in one step. However, we need to set the *allowed address pairs* on each port attached to our OpenShift nodes and that cannot be done via the `server create` subcommand.

So we will be creating the ports separately from the servers:

```sh
$ openstack port create --network "$INFRA_ID-network" --security-group "$INFRA_ID-master" --allowed-address ip-address=192.0.2.5 --allowed-address ip-address=192.0.2.6 --allowed-address ip-address=192.0.2.7 --tag openshiftClusterID="$INFRA_ID" "$INFRA_ID-bootstrap-port"
```

Since the keepalived-managed IP addresses are not attached to any specific server, Neutron would block their traffic by default. By passing them to `--allowed-address` the traffic can flow freely through.

We will also add attach an additional Floating IP address to the bootstrap port:

```sh
$ openstack floating ip create --description "Bootstrap IP" external
=> 203.0.113.24
$ openstack floating ip set --port "$INFRA_ID-bootstrap-port" <203.0.113.24>
```

This is not necessary for the deployment (and we will delete the bootstrap resources afterwards). However, if the bootstrapping phase fails for any reason, the installer will try to SSH in and download the bootstrap log. That will only succeed if the node is reachable (which in general means a floating IP).


### Bootstrap Server

Now we can create the bootstrap server which will help us configure up the control plane:

```sh
$ openstack server create --image <rhcos> --flavor <m1.xlarge> --user-data "$INFRA_ID-bootstrap-ignition.json" --port "$INFRA_ID-bootstrap-port" --wait "$INFRA_ID-bootstrap" --property openshiftClusterID="$INFRA_ID"
```

After the server is active, you can check the console log to see that it is getting the ignition correctly:

```
$ openstack console log show "$INFRA_ID-bootstrap"
```

You can also SSH into the server (using its floating IP address) and check on the bootstrapping progress:

```sh
$ ssh core@203.0.113.24
[core@openshift-qlvwv-bootstrap ~]$ journalctl -b -f -u bootkube.service
```


## Control Plane

Similar to the bootstrap, the control plane nodes will need to have their allowed address pairs set which means we have to create the ports separately.

Our control plane will consist of three nodes.

### Control Plane Ports

```sh
for index in $(seq 0 2); do
    openstack port create --network "$INFRA_ID-network" --security-group "$INFRA_ID-master" --allowed-address ip-address=192.0.2.5 --allowed-address ip-address=192.0.2.6 --allowed-address ip-address=192.0.2.7 --tag openshiftClusterID="$INFRA_ID" "$INFRA_ID-master-port-$index"
done
```

### Control Plane Trunks (Required for Kuryr SDN)

We will create the Trunks for Kuryr to plug the containers into the OpenStack SDN.

```sh
for index in $(seq 0 2); do
    openstack network trunk create --parent-port "$INFRA_ID-master-port-$index" "$INFRA_ID-master-trunk-$index"
done
```


### Control Plane Servers

We will create the servers, passing in the `master-?-ignition.json` files prepared earlier:

```sh
for index in $(seq 0 2); do
    openstack server create --image rhcos --flavor m1.xlarge --user-data "$INFRA_ID-master-$index-ignition.json" --port "$INFRA_ID-master-port-$index" --property openshiftClusterID="$INFRA_ID" "$INFRA_ID-master-$index"
done
```

**NOTE**: You might want to do things such as boot from volume, attach specific volumes, etc. instead of the command above. Feel free to adapt it to your needs.


The master nodes should load the initial Ignition and then keep waiting until the bootstrap node stands up the Machine Config Server which will provide the rest of the configuration.

### Wait for the Control Plane to Complete

When that happens, the masters will start running their own pods, run etcd and join the "bootstrap" cluster. Eventually, they will form a fully operational control plane.

You can monitor this via the following command:

```sh
$ openshift-install wait-for bootstrap-complete
```

Eventually, it should output the following:

```
INFO API v1.14.6+f9b5405 up
INFO Waiting up to 30m0s for bootstrapping to complete...
```

This means the masters have come up successfully and are joining the cluster.

Eventually, the `wait-for` command should end with:

```
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

You can now safely delete the bootstrap port, server and floating IP address:

```sh
$ openstack server delete "$INFRA_ID-bootstrap"
$ openstack port delete "$INFRA_ID-bootstrap-port"
$ openstack floating ip delete 203.0.113.24
```

If you hadn't done so already, you should also disable the bootstrap Ignition URL:

https://example.com/bootstrap.ign

## Compute Nodes

This process is similar to the masters, but the workers need to be approved before they're allowed to join the cluster.

We will create three worker nodes here.

**NOTE**: some of the default pods (e.g. the `openshift-router`) require at least two nodes so that is the effective minimum.

The workers need no ignition override -- we can pass the unmodified `worker.ign` as their user data:

### Compute Nodes ports

```sh
for index in $(seq 0 2); do
    openstack port create --network "$INFRA_ID-network" --security-group "$INFRA_ID-worker" --allowed-address ip-address=192.0.2.5 --allowed-address ip-address=192.0.2.6 --allowed-address ip-address=192.0.2.7 --tag openshiftClusterID="$INFRA_ID" "$INFRA_ID-worker-port-$index"
done
```

### Compute Nodes Trunks (Required for Kuryr SDN)

We will create the Trunks for Kuryr to plug the containers into the OpenStack SDN.

```sh
for index in $(seq 0 2); do
    openstack network trunk create --parent-port "$INFRA_ID-worker-port-$index" "$INFRA_ID-worker-trunk-$index"
done
```

### Compute Nodes server

```sh
for index in $(seq 0 2); do
    openstack server create --image rhcos --flavor m1.large --user-data "worker.ign" --port "$INFRA_ID-worker-port-$index" --property openshiftClusterID="$INFRA_ID" "$INFRA_ID-worker-$index"
done
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
csr-5jwqf   16m    system:node:openshift-qlvwv-master-0                                         Approved,Issued
csr-88jp8   116s   system:serviceaccount:openshift-machine-config-operator:node-bootstrapper   Pending
csr-9dt8f   15m    system:node:openshift-qlvwv-master-1                                         Approved,Issued
csr-bqkw5   16m    system:serviceaccount:openshift-machine-config-operator:node-bootstrapper   Approved,Issued
csr-dpprd   6s     system:serviceaccount:openshift-machine-config-operator:node-bootstrapper   Pending
csr-dtcws   24s    system:serviceaccount:openshift-machine-config-operator:node-bootstrapper   Pending
csr-lj7f9   16m    system:node:openshift-qlvwv-master-2                                         Approved,Issued
csr-lrtlk   15m    system:serviceaccount:openshift-machine-config-operator:node-bootstrapper   Approved,Issued
csr-wkm94   16m    system:serviceaccount:openshift-machine-config-operator:node-bootstrapper   Approved,Issued
```

You should inspect each pending CSR and verify that it comes from a node you recognise:

```
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

First, remove the `api` and `*.apps` DNS records.

Then run the following commands to delete all the OpenStack resources we've created:

```sh
$ openstack floating ip delete 203.0.113.23 203.0.113.19
$ openstack server list --column Name --format value | grep "$INFRA_ID" | xargs --no-run-if-empty openstack server delete
```

If Kuryr SDN is used, delete the Trunks and Load Balancers created:
```sh
$ openstack network trunk list --column Name --format value | grep "$INFRA_ID" | xargs --no-run-if-empty openstack network trunk delete

$ for id in $(openstack loadbalancer list -f value -c id); do
    lb=$(openstack loadbalancer show $id |grep openshiftClusterID="$INFRA_ID");
    if [ -n "$lb" ]; then
        openstack loadbalancer delete --cascade $id
    fi
done
```

Proceed with the deletion of the remaining OpenStack resources:
```sh
$ openstack port list --device-owner 'network:router_interface' -c ID -f value | xargs --no-run-if-empty -I {} openstack router remove port "$INFRA_ID"-external-router {}
$ openstack port list --column id  --format value --tags openshiftClusterID="$INFRA_ID" | xargs --no-run-if-empty openstack port delete

$ openstack router unset --external-gateway "$INFRA_ID-external-router"
$ openstack router remove subnet "$INFRA_ID-extrenal-router" "$INFRA_ID-nodes"
$ openstack router delete "$INFRA_ID-external-router"

$ openstack subnet list --c ID --format value --tags openshiftClusterID="$INFRA_ID" | xargs --no-run-if-empty openstack subnet delete
$ openstack subnet pool list --c ID --format value --tags openshiftClusterID="$INFRA_ID" | xargs --no-run-if-empty openstack subnet pool delete
$ openstack network list --c ID --format value --tags openshiftClusterID="$INFRA_ID" | xargs --no-run-if-empty openstack network delete
$ openstack security group list --c ID --format value --tags openshiftClusterID="$INFRA_ID" | xargs --no-run-if-empty openstack security group delete
```
