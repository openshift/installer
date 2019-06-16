# OpenStack Platform Support

Support for launching clusters on OpenStack is **experimental**.

This document discusses the requirements, current expected behavior, and how to
try out what exists so far.

## OpenStack Requirements

The installer assumes the following about the OpenStack cloud you run against:

* You must create a `clouds.yaml` file with the auth URL and credentials
    necessary to access the OpenStack cloud you want to use.  Information on
    this file can be found at
    https://docs.openstack.org/os-client-config/latest/user/configuration.html
    and it looks like:
```
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
      project_id: <ID of the project. Get it from 'openstack project list'>
      project_domain_name: Default
      user_domain_name: Default         
      auth_url: 'https://10.10.14.22:5001/v2.0'
```

* Swift must be enabled.  The user must have `swiftoperator` permissions and
  `temp-url` support must be enabled. As an OpenStack admin:
  * `openstack role add --user <user> --project <project> swiftoperator`
  * `openstack object store account set --property Temp-URL-Key=superkey`

**NOTE:** Swift is required as the user-data provided by OpenStack is not big
enough to store the ignition config files, so they are served by swift instead.

* You may need to increase the security group related quotas from their default
  values. For example (as an OpenStack admin) `openstack quota set --secgroups 100 --secgroup-rules 1000 <project>`


* A file `install-config.yaml` is required. It defines your cluster topology. Save it because a failed installation will delete it. It looks like this:

```
apiVersion: v1
baseDomain: example.com
controlPlane:
  name: master
  platform:
    openstack:
      type: m1.master
  replicas: 3
compute:
- name: worker
  platform:
    openstack:
      type: m1.master
  replicas: 1
metadata:
  name: test-cluster
networking:
  clusterNetworks:
  - cidr: 10.128.0.0/14
    hostSubnetLength: 9
  machineCIDR: 10.0.0.0/16
  serviceCIDR: 172.30.0.0/16
  type: OpenShiftSDN
platform:
  openstack:
    cloud: openstack
    region: RegionOne
    externalNetwork: openshift-external
    flavorName: m1.master
    computeFlavor: m1.master
    baseDomainResourceGroupName: os4-common
    lbFloatingIP: "<FLOATING EXTERNAL IP ADDRESS>"
pullSecret: '{"auths":{"cloud.openshift.com": ...}'
sshKey: ssh-rsa AAAA...
```
### Pull secret
The pull Secret [can be obtained here](https://cloud.redhat.com/openshift/install).

### ssh

The sshKey in the `install-config.yaml` is required for ssh'ing into the virtual machines.

If you only have a private IP on your VMs you can do that to ssh into them from the OpenStack Controller:

```
sudo ip netns ls

$ sudo ip netns ls
qdhcp-d3904a4d-f25d-4145-bc85-f2be0573ffd7 (id: 12)
qrouter-a34...
...
```

Grab the ID **qdhcp**-`<YOUR OPENSHIFT NETWORK ID>` (in this case: qdhcp-d3904a4d-f25d-4145-bc85-f2be0573ffd7)

```
sudo ip netns exec qdhcp-<YOUR NETWORK ID> ssh -i <PRIVATE SSH KEY FILE> core@<PRIVATE IP OF YOUR VM>
```

* The installer requires a proper RHCOS image in the OpenStack cluster or project:
`openstack image create --container-format=bare --disk-format=qcow2 --file rhcos-${RHCOSVERSION}-openstack.qcow2 rhcos-${RHCOSVERSION}`

**NOTE:** Depending on your OpenStack environment you can upload the RHCOS image
as `raw` or `qcow2`. See [Disk and container formats for images](https://docs.openstack.org/image-guide/image-formats.html) for more information. The installer looks for an image named rhcos. This could be overridden via the `OPENSHIFT_INSTALL_OS_IMAGE_OVERRIDE` environment variable if for instance your cloud operator provides the image under a different name.

* The public network should be created by the OSP admin. Verify the name/ID of the 'External' network:
```
openstack network list --long -c ID -c Name -c "Router Type"
+--------------------------------------+----------------+-------------+
| ID                                   | Name           | Router Type |
+--------------------------------------+----------------+-------------+
| 148a8023-62a7-4672-b018-003462f8d7dc | public_network | External    |
+--------------------------------------+----------------+-------------+
```

**NOTE:** If the `neutron` `trunk` service plugin is enabled, trunk port will be created by default. for more information, please refer to [neutron trunk port](https://wiki.openstack.org/wiki/Neutron/TrunkPort).

### Isolated Development

If you would like to set up an isolated development environment, you may use a
bare metal host running CentOS 7.  The following repository includes some
instructions and scripts to help with creating a single-node OpenStack
development environment for running the installer.  Please refer to the
documentation in that repository for further details.

* https://github.com/shiftstack-dev-tools/ocp-doit

## OpenShift API Access

All the OpenShift nodes are created in an OpenStack tenant network and as such, can't be accessed directly. The installer does not create any floating IP addresses.

However, the installer does need access to the OpenShift's API as it is being deployed.

There are two ways you can handle this.

### Bring Your Own Cluster IP

We recommend you create a floating IP address ahead of time, add the
API record to your own DNS and let the installer use that address.

First, create the floating IP:

    $ openstack floating ip create <external network>

Note the actual IP address. We will use `10.19.115.117` throughout this
document.

Next, add the `api.<cluster name>.<cluster domain>` and `*.apps.<cluster
name>.<cluster domain>` name records pointing to that floating IP to your DNS:

    api.ostest.shiftstack.com IN A 10.19.115.117
    *.apps.ostest.shiftstack.com  IN  A  10.19.115.117

If you don't have a DNS server under your control, you finish the installation
by adding the following to your `/etc/hosts`:

    10.19.115.117 api.ostest.shiftstack.com
    10.19.115.117 console-openshift-console.apps.ostest.shiftstack.com

**NOTE:** *this will make the API accessible only to you. This is fine for your
own testing (and it is enough for the installation to succeed), but it is not
enough for a production deployment.*

Finally, add the floating IP address to `install-config.yaml`.

It should be under `platform.openstack.lbFloatingIP`. For example:

```yaml
apiVersion: v1beta2
baseDomain: shiftstack.com
clusterID:  3f47f546-c010-4c46-895c-c8fce6cf0451
# ...
platform:
  openstack:
    cloud:            standalone
    externalNetwork:  public
    region:           regionOne
    computeFlavor:    m1.medium
    lbFloatingIP:     "10.19.115.117"
```

This will let you do a fully unattended end to end deployment.


### No Floating IP

If you don't want to pre-create a floating IP address, you will still want to create the API DNS record or the installer will fail waiting for the API.

Without the floating IP, you won't know the right IP address of the server ahead of time, so you will have to wait for it to come up and create the DNS records then:

    $ watch openstack server list

Wait for the `<cluster name>-api` server comes up and you can make your changes then.

**WARNING:** The installer will fail if it can't reach the bootstrap OpenShift API in 30 minutes.

Even if the installer times out, the OpenShift cluster should still come up. Once the bootstrapping process is in place, it should all run to completion.

So you should be able to deploy OpenShift without any floating IP addresses and DNS records and create everything yourself after the cluster is up.


## Current Expected Behavior

As mentioned, OpenStack support is still experimental. Currently:

* Deploys an isolated tenant network
* Deploys a instance used as a 'service VM' that hosts a load balancer for the
OpenShift API and as an internal DNS for the instances
* Deploys a bootstrap instance to bootstrap the OpenShift cluster
* Once the masters are deployed, the bootstrap instance is destroyed
* Deploys 3 master nodes
* The OpenShift UI is served at `https://<cluster name>-api.<cluster domain>` (but you need to create that DNS record yourself)

The installer should finish successfully, though it is still undergoing development and things might break from time to time.

### Workarounds

#### Switch over from bootstrap control plane to production control plane

During the installation a support VM will be spawned. It contains a load balancer (HAPROXY) and a DNS Server (CoreDNS).

On this VM runs a service, which periodically checks the production cluster for its control plane being available. In this case the IP address of the bootstrap control plane will be removed from the load balancer's config, the load balancer will be restarted. 

Sometimes this watcher service starts to flutter because during the installation process the API servers on the production control plane might be restarted. This can delay the installation procedure a lot.

Workaround: Disable the watcher script on the support VM (with HAPROXY) immediately after it has been spawned by the installer until the control plane (all API servers on the master VMs) is available:

```
sudo systemctl stop haproxy-watcher.timer
sudo systemctl stop haproxy-watcher
sudo systemctl disable haproxy-watcher.timer
sudo systemctl disable haproxy-watcher
```

SSH into the masters and try to get a 'not authorized' answer from the API servers:

```
curl -k https://localhost:6443
```

If the response is no `curl: (7) Failed connect to localhost:6443; Connection refused` you can enable the watcher on the support VM again:

```
sudo systemctl start haproxy-watcher.timer
sudo systemctl start haproxy-watcher
sudo systemctl enable haproxy-watcher.timer
sudo systemctl enable haproxy-watcher
```

#### Increase timeouts

If you are regulary are running into timeout problems you can increase the timeouts in `cmd/openshift-install/create.go`. In the future it might by better to make them configurable in the install-config.yaml script.


#### External DNS

While deploying the cluster, the installer will hang trying to reach the API as
the node running the installer cannot resolve the service VM (the cluster
should still come up successfully within the isolated network).

You can add the service VM floating IP address at the top of your `/etc/resolv.conf`:

```
$ cat /etc/resolv.conf
# Generated by NetworkManager
search example.com
# OpenShift Service VM DNS:
nameserver 10.19.115.117

# Your previous DNS servers:
nameserver 83.240.0.215
nameserver 83.240.0.136
```

(the service VM floating IP is `10.19.115.117` in this example)

If you don't want to update your DNS config, you can add a couple of entries in your `/etc/hosts` file instead:

```
$ cat /etc/hosts
127.0.0.1   localhost localhost.localdomain localhost4 localhost4.localdomain4
::1         localhost localhost.localdomain localhost6 localhost6.localdomain6
10.19.115.117 <cluster-name>-api.<domain>
10.19.115.117 console-openshift-console.apps.<cluster-name>.<domain>
```

If you do expose the cluster, the installer should complete successfully.

It will print the console URL, username and password and you should be able to go there and log in.

```
INFO Install complete!
INFO To access the cluster as the system:admin user when using 'oc', run 'export KUBECONFIG=/path/to/installer/auth/kubeconfig'
INFO Access the OpenShift web-console here: https://console-openshift-console.apps.ostest.shiftstack.com
INFO Login to the console with user: kubeadmin, password: 5char-5char-5char-5char
```

## Using an External Load Balancer

This documents how to shift from the API VM load balancer, which is
intended for initial cluster deployment and not highly available, to an
external load balancer.

The load balancer must serve ports 6443, 443, and 80 to any users of
the system.  Port 22623 is for serving ignition startup configurations
to the OpenShift nodes and should not be reachable outside of the cluster.

The first step is to add floating IPs to all the master nodes:

* `openstack floating ip create --port master-port-0 <public network>`
* `openstack floating ip create --port master-port-1 <public network>`
* `openstack floating ip create --port master-port-2 <public network>`

Once complete you can see your floating IPs using:

* `openstack server list`

These floating IPs can then be used by the load balancer to access
the cluster.  An example haproxy configuration for port 6443 is below.
The other port configurations are identical.

```
listen <cluster name>-api-6443
    bind 0.0.0.0:6443
    mode tcp
    balance roundrobin
    server <cluster name>-master-2 <floating ip>:6443 check
    server <cluster name>-master-0 <floating ip>:6443 check
    server <cluster name>-master-1 <floating ip>:6443 check
```

The next step is to allow network access from the load balancer network
to the master nodes:

* `openstack security group rule create master --remote-ip <load balancer CIDR> --ingress --protocol tcp --dst-port 6443`
* `openstack security group rule create master --remote-ip <load balancer CIDR> --ingress --protocol tcp --dst-port 443`
* `openstack security group rule create master --remote-ip <load balancer CIDR> --ingress --protocol tcp --dst-port 80`

You could also specify a specific IP address with /32 if you wish.

You can verify the operation of the load balancer now if you wish, using the
curl commands given below.

Now the DNS entry for <cluster name>-api.<base domain> needs to be updated
to point to the new load balancer:

* `<load balancer ip> <cluster-name>-api.<base domain>`

The external load balancer should now be operational along with your own
DNS solution. It's best to test this configuration before removing
the API. The following curl command is an example of how
to check functionality:

`curl https://<loadbalancer-ip>:6443/version --insecure`

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

Another useful thing to check is that the ignition configurations are
only available from within the deployment. The following command should
only succeed from a node in the OpenShift cluster:

* `curl https://<loadbalancer ip>:22623/config/master --insecure`

Now that the DNS and load balancer has been moved, we can take down the existing
api VM:

* `openstack server delete <cluster name>-api`

## Disambiguating the External Network

The installer assumes that the name of the external network is unique.  In case
there is more than one network with the same name as the desired external
network, it’s possible to provide a UUID to specify which network should be
used.

```
$ env TF_VAR_openstack_external_network_id="6a32627e-d98d-40d8-9324-5da7cf1452fc" \
> bin/openshift-install create cluster
```

## Troubleshooting

See the [troubleshooting installer issues in OpenStack](./troubleshooting.md) guide.

## Reporting Issues

Please see the [Issue Tracker][issues_openstack] for current known issues.
Please report a new issue if you do not find an issue related to any trouble
you’re having.

[issues_openstack]: https://github.com/openshift/installer/issues?utf8=%E2%9C%93&q=is%3Aissue+is%3Aopen+openstack
