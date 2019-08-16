# OpenStack Platform Support

Support for launching clusters on OpenStack is **experimental**.

This document discusses the requirements, current expected behavior, and how to
try out what exists so far.

In addition, it covers the installation with the default CNI (OpenShiftSDN),
as well as with the Kuryr SDN

## Kuryr SDN

Kuryr is a CNI plug-in that uses Neutron and Octavia to provide networking
for pods and services. It is primarily designed for OpenShift clusters that
run on OpenStack virtual machines. Kuryr improves the network performance by
plugging OCP pods into OpenStack SDN. In addition, it provides
interconnectivity between OCP pods and OpenStack virtual instances.

Kuryr is recommended for OpenShift Container Platform deployments on
encapsulated OpenStack tenant networks in order to avoid double encapsulation,
such as running an encapsulated OpenShift SDN over an OpenStack network.

Conversely, using Kuryr does not make sense in the following cases:

   * You use provider networks or tenant VLANs.
   * The deployment will use many services on a few hypervisors. Each
   OpenShift service creates an Octavia Amphora virtual machine in OpenStack
   that hosts a required load balancer.
   * UDP services are needed.


## Openstack Credentials

There are two ways to pass your credentials to the installer, with a clouds.yaml file or with environment variables. You can also use a combination of the two, but be aware that clouds.yaml file has precident over the environment variables you set.

The installer will look for a clouds.yaml file in the following locations in order:
1. OS_CLIENT_CONFIG_FILE
2. Current directory
3. unix-specific user config directory (~/.config/openstack/clouds.yaml)
4. unix-specific site config directory (/etc/openstack/clouds.yaml)

In many openstack distributions, you can get a clouds.yaml file through Horizon. If you cant, then you can make a `clouds.yaml` file yourself. Information on
    this file can be found at https://docs.openstack.org/openstacksdk/latest/user/config/configuration.html#config-files
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
      auth_url: 'https://10.10.14.22:5001/v2.0'
```

If you choose to use environment variables in place of a clouds.yaml, or along side it, consult the following doccumentation:
https://www.terraform.io/docs/providers/openstack/#configuration-reference


## OpenStack Requirements

### Recommended Minimums

In order to run the latest version of the installer in OpenStack, at a bare minimum you need the following quota to run a *default* cluster. While it is possible to run the cluster with fewer resources than this, it is not recommended. Certian edge cases, such as deploying [without FIPs](#without-floating-ips), or deploying with an [external loadbalancer](#using-an-external-load-balancer) are documented below, and are not included in the scope of this recomendation.

   * OpenStack Quota
     * Floating IPs: 3
     * Security Groups: 3
     * Security Group Rules: 60
     * Routers: 1
     * Subnets: 1
     * RAM: 112 Gb
     * VCPU: 28
     * Volume Storage: 175 Gb
     * Instances: 7

### Recommended Minimums with Kuryr SDN

When using Kuryr SDN, as the pods, services, namespaces, network policies, etc.,
are using resources from the OpenStack Quota, the minimum requirements are
higher:

   * OpenStack Quota
     * Floating IPs: 3 (plus the expected number of services of LoadBalancer
       type)
     * Security Groups: 100 (1 needed per network policy)
     * Security Group Rules: 500
     * Routers: 1
     * Subnets: 100 (1 needed per namespace)
     * Networks: 100 (1 needed per namespace)
     * Ports: 500
     * RAM: 112 Gb
     * VCPU: 28
     * Volume Storage: 175 Gb
     * Instances: 7

#### Master Nodes

The default deployment stands up 3 master nodes, which is the minimum amount required for a cluster. For each master node you stand up, you will need 1 instance, and 1 port available in your quota. They should be assigned a flavor with at least 16 Gb RAM, 4 VCPu, and 25 Gb Disk. It is theoretically possible to run with a smaller flavor, but be aware that if it takes too long to stand up services, or certian essential services crash, the installer could time out, leading to a failed install.

#### Worker Nodes

The default deployment stands up 3 worker nodes. In our testing we determined that 2 was the minimum number of workers you could have to get a succesful install, but we don't recommend running with that few. Worker nodes host the apps you run on OpenShift, so it is in your best interest to have more of them. See [here](https://docs.openshift.com/enterprise/3.0/architecture/infrastructure_components/kubernetes_infrastructure.html#node) for more information. The flavor assigned to the worker nodes should have at least 2 VCPUs, 8 Gb RAM and 25 Gb Disk. However, if you are experiencing `Out Of Memory` issues, or your installs are timing out, you should increase the size of your flavor to match the masters: 4 VCPUs and 16 Gb RAM.

#### Bootstrap Node

The bootstrap node is a temporary node that is responsable for standing up the control plane on the masters. Only one bootstrap node will be stood up. To do so, you need 1 instance, and 1 port. We recommend a flavor with a minimum of 16 Gb RAM, 4 VCPUs, and 25 Gb Disk.


### Swift

Swift must be enabled.  The user must have `swiftoperator` permissions and
  `temp-url` support must be enabled. As an OpenStack admin:
  * `openstack role add --user <user> --project <project> swiftoperator`
  * `openstack object store account set --property Temp-URL-Key=superkey`

**NOTE:** Swift is required as the user-data provided by OpenStack is not big
enough to store the ignition config files, so they are served by swift instead.

* You may need to increase the security group related quotas from their default
  values. For example (as an OpenStack admin) `openstack quota set --secgroups 8 --secgroup-rules 100 <project>`

### RHCOS Image

If you do not have a Red Hat Core OS image already, or are looking for the latest,
 [click here](https://mirror.openshift.com/pub/openshift-v4/dependencies/rhcos/pre-release/latest/).

The installer requires a proper RHCOS image in the OpenStack cluster or project:
`openstack image create --container-format=bare --disk-format=qcow2 --file rhcos-${RHCOSVERSION}-openstack.qcow2 rhcos`

**NOTE:** Depending on your OpenStack environment you can upload the RHCOS image
as `raw` or `qcow2`. See [Disk and container formats for images](https://docs.openstack.org/image-guide/image-formats.html) for more information. The installer looks for an image named rhcos. This could be overridden via the `OPENSHIFT_INSTALL_OS_IMAGE_OVERRIDE` environment variable if for instance your cloud operator provides the image under a different name.

### Neutron Public Network

The public network should be created by the OSP admin. Verify the name/ID of the 'External' network:
```
openstack network list --long -c ID -c Name -c "Router Type"
+--------------------------------------+----------------+-------------+
| ID                                   | Name           | Router Type |
+--------------------------------------+----------------+-------------+
| 148a8023-62a7-4672-b018-003462f8d7dc | public_network | External    |
+--------------------------------------+----------------+-------------+
```

**NOTE:** If the `neutron` `trunk` service plugin is enabled, trunk port will be created by default. for more information, please refer to [neutron trunk port](https://wiki.openstack.org/wiki/Neutron/TrunkPort).

### Extra requirements when enabling Kuryr SDN


#### Increase Quota

As highlighted in the minimum quota recommendations, when using Kuryr SDN, there
is a need for increasing the quotas as pods, services, namespaces, network
policies are using OpenStack resources. So, as an admin, the next
quotas should be increased for the selected project:

```
openstack quota set --secgroups 100 --secgroup-rules 500 --ports 500 --subnets 100 --networks 100 <project>
```

#### Neutron Configuration

Kuryr CNI makes use of the Neutron Trunks extension to plug containers into the
OpenStack SDN, so the `trunks` extension must be enabled for Kuryr to properly
work.

In addition, if the default ML2/OVS Neutron driver is used, the firewall must be
set to `openvswitch` instead of `ovs_hybrid` so that security groups are
enforced on trunk subports and Kuryr can properly handle Network Policies.

#### Octavia

Kuryr SDN uses Octavia OpenStack LBaaS to implement OpenShift services.
Thus the OpenStack enviroment must have Octavia components installled and
configured if Kuryr SDN is used.

**NOTE:** Depending on your OpenStack environment Octavia may not support UDP
listeners, which means there is no support for UDP services if kuryr SDN is
used.

### Isolated Development

If you would like to set up an isolated development environment, you may use a
bare metal host running CentOS 7.  The following repository includes some
instructions and scripts to help with creating a single-node OpenStack
development environment for running the installer.  Please refer to the
documentation in that repository for further details.

* https://github.com/shiftstack-dev-tools/ocp-doit

## API Access

All the OpenShift nodes are created in an OpenStack tenant network and as such, can't be accessed directly in most openstack deployments. The installer does not create any floating IP addresses, but does need access to the OpenShift's API as it is being deployed. We will briefly explain how to set up access to the openshift api with and without floating IP addresses.

### Using Floating IPs

This method allows you to attach two floating IP addresses to endpoints in OpenShift.

First, create a floating IP address for the API:

    $ openstack floating ip create <external network>

Next, add the `api.<cluster name>.<cluster domain>` and `*.apps.<cluster
name>.<cluster domain>` name records pointing to that floating IP to your DNS:

    api.example.shiftstack.com IN A <API FIP>

If you don't have a DNS server under your control, you finish the installation
by adding the following to your `/etc/hosts`:

    <API FIP> api.example.shiftstack.com

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
    computeFlavor:    m1.large
    lbFloatingIP:     "<API FIP>"
```

At the time of writing, you will have to create a second floating ip and attach it to the ingress-port if you want to be able to reach *.apps externally.
This can be done after the install completes in three steps:

Get the ID of the ingress port:

```sh
openstack port list | grep "ingress-port"
```

Create and associate a floating IP to the ingress port:

```sh
openstack floating ip create --port <ingress port id> <external network>
```

Add A record in your dns for *apps. in your DNS:

```
*.apps.example.shiftstack.com  IN  A  <ingress FIP>
```
OR add A record in `/etc/hosts`:

```
    <ingress FIP> console-openshift-console.apps.example.shiftstack.com
```

### Without Floating IPs

If you don't want to pre-create a floating IP address, you will still want to create the API DNS record or the installer will fail waiting for the API.
Without the floating IP, you won't know the right IP address of the server ahead of time, so you will have to wait for it to come up and create the DNS records then:

    $ watch openstack server list

Wait for the `<cluster name>-api` server comes up and you can make your changes then.
**WARNING:** The installer will fail if it can't reach the bootstrap OpenShift API in 30 minutes.
Even if the installer times out, the OpenShift cluster should still come up. Once the bootstrapping process is in place, it should all run to completion.
So you should be able to deploy OpenShift without any floating IP addresses and DNS records and create everything yourself after the cluster is up.

## Installing with Kuryr SDN

To deploy with Kuryr SDN instead of the default OpenShift SDN, you simply need
to modify the `install-config.yaml` file to include `Kuryr` as the desired
`networking.networkType` and proceed with the same steps as with the default
OpenShift SDN:

```yaml
apiVersion: v1
...
networking:
  networkType: Kuryr
  ...
platform:
  openstack:
    ...
    trunkSupport: true
    octaviaSupport: true
    ...
```


**NOTE:** both trunkSupport and octaviaSupport are automatically discovered
by the installer, so there is no need to set them. But if your env don't meet
both requirements Kuryr SDN will not properly work, as trunks are needed to
connect the pods to the OpenStack network and Octavia to create the OpenShift
services.

### Known limitations of installing with Kuryr SDN

There are known limitations when using Kuryr SDN:

   * There is an amphora load balancer VM being deployed per OpenShift svc
   with the default Octavia load balancer driver (amphora driver). If the
   environment is resource constrained it could be a problem to create a large
   amount of services.
   * Depending on the Octavia OpenStack version, UDP listeners are not
   supported. This means that OpenShift UDP services are not supported.
   * There is a known limitation of Octavia not supporting listeners on UDP
   and TCP on the same port. Thus services expose the same port for UDP and TCP
   are not supported -- only the TCP listener will be created.
   * Due to the above UDP limitations of Octavia, Kuryr is forcing pods to
   use TCP for DNS resolution (`use-vc` option at `resolv.conf`). This may be a
   problem for pods running go applications compiled with `CGO_DEBUG` flag
   disabled as that forces to use the `go` resolver that is only using UDP and
   is not considering the `use-vc` option added by kuryr to the `resolv.conf`.

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
