# Kuryr

Kuryr is a CNI plug-in that uses Neutron and Octavia to provide networking for pods and services. It is primarily designed for OpenShift clusters that run on OpenStack virtual machines. Kuryr improves the network performance by plugging
OpenShift pods into OpenStack SDN. In addition, it provides interconnectivity between OpenShift pods and OpenStack virtual instances.

Kuryr is recommended for OpenShift deployments on encapsulated OpenStack tenant networks in order to avoid double encapsulation, such as running an encapsulated OpenShift SDN over an OpenStack network.

If you use provider networks or tenant VLANs, you do not need to use Kuryr to
avoid double encapsulation. The performance benefit is negligible. Depending on
your configuration, though, using Kuryr to avoid having two overlays might still
be beneficial.

Kuryr is not recommended in deployments where all of the following criteria are
true:

* The OpenStack version is less than Train.
* The deployment uses UDP services, or a large number of TCP services on few
  hypervisors.

or

* The `ovn-octavia` Octavia driver is disabled.
* The deployment uses a large number of TCP services on few hypervisors.

[Learn more about Kuryr](https://docs.openstack.org/kuryr-kubernetes/latest/).

## Requirements when enabling Kuryr

When using Kuryr SDN, as the pods, services, namespaces, network policies, etc., are using resources from the OpenStack Quota, the minimum requirements are higher. Kuryr also has some additional requirements on top of what a default install requires. At a bare minimum you need the following quota to run a *default* cluster:

* Floating IPs: 3 (plus the expected number of services of LoadBalancer type)
* Security Groups: 250 (1 needed per Service and per NetworkPolicy)
* Security Group Rules: 1000
* Routers: 1
* Subnets: 250 (1 needed per namespace)
* Networks: 250 (1 needed per namespace)
* Ports: 1500
* RAM: 112 GB
* vCPUs: 28
* Volume Storage: 175 GB
* Instances: 7
* Load Balancers: 100 (1 needed per OpenShift service)
* Load Balancer Listeners: 500 (1 per port exposed on the services)
* Load Balancer Pools: 500 (1 per port exposed on the services)

## Increase Quota

As highlighted in the minimum quota recommendations, when using Kuryr SDN, there is a need for increasing the quotas as pods, services, namespaces, network policies are using OpenStack resources. So, as an administrator,the next quotas should be increased for the selected project:

```sh
openstack quota set --secgroups 250 --secgroup-rules 1000 --ports 1500 --subnets 250 --networks 250 <project>
```

**NOTE:** Each Amphora Load Balancer creates a VM. Even if that VM is not part of the user quota, the OpenStack cluster must have enough capacity to allocate those VMs too. A standard OpenShift installation ends up with more that 50 Amphora Load Balancers.

## Neutron Configuration

Kuryr CNI makes use of the Neutron Trunks extension to plug containers into the OpenStack SDN, so the `trunks` extension must be enabled for Kuryr to properly work.

In addition, if the default ML2/OVS Neutron driver is used, the firewall must be set to `openvswitch` instead of `ovs_hybrid` so that security groups are enforced on trunk subports and Kuryr can properly handle Network Policies.

## Octavia

Kuryr SDN uses OpenStack Octavia LBaaS to implement OpenShift services. Thus the OpenStack environment must have Octavia components installed and configured if Kuryr SDN is used.

**NOTE:** Depending on your OpenStack environment Octavia may not support UDP listeners, which means there is no support for UDP services if Kuryr SDN is used.

### The Octavia OVN Driver

Octavia supports multiple provider drivers through the Octavia API.
￼
To see all available Octavia provider drivers, on a command line, enter:

```yaml
$ openstack loadbalancer provider list
+---------+-------------------------------------------------+
| name    | description                                     |
+---------+-------------------------------------------------+
| amphora | The Octavia Amphora driver.                     |
| octavia | Deprecated alias of the Octavia Amphora driver. |
| ovn     | Octavia OVN driver.                             |
+---------+-------------------------------------------------+
```
￼
Beginning with OpenStack Train, the Octavia OVN provider driver (`ovn`) is supported.
`ovn` is an integration driver for the load balancing that Octavia and OVN provide.
It supports basic load balancing capabilities, and is based on OpenFlow rules.
￼
The Amphora provider driver is the default driver. If `ovn` is enabled,
however, Kuryr uses it.￼
If Kuryr uses `ovn` instead of Amphora, it offers the following benefits:
￼
* Decreased resource requirements. Kuryr does not require a load balancer VM
for each Service.
* Reduced network latency.
* Increased service creation speed by using OpenFlow rules instead of a VM for
each Service.
* Distributed load balancing actions across all nodes instead of centralized on
Amphora VMs.

## Installing with Kuryr SDN

To deploy with Kuryr SDN instead of the default OpenShift SDN, you simply need to modify the `install-config.yaml` file to include `Kuryr` as the desired `networking.networkType` and proceed with the same steps as with the default OpenShift SDN:

```yaml
apiVersion: v1
...
networking:
  networkType: Kuryr
  ...
```

**NOTE:** If your environment doesn't support trunks or OpenStack Octavia service isn't available, Kuryr SDN will not properly work, as trunks are needed to connect the pods to the OpenStack network and Octavia to create the OpenShift services.

### Known limitations of installing with Kuryr SDN

There are known limitations when using Kuryr SDN:

* There is an amphora load balancer VM being deployed per OpenShift service with the default Octavia load balancer driver (amphora driver). If the environment is resource constrained it could be a problem to create a large amount of services as each one will create a new VM.
* Depending on the Octavia version, UDP listeners are not supported. This means that OpenShift UDP services are not supported unless OpenStack version is Train or newer.
* Depending on the Octavia version, there is a known limitation of Octavia not supporting listeners on different protocols (e.g., UDP and TCP) on the same port. Thus services exposing the same port for different protocols are not supported unless OpenStack version is Train or newer.
* Due to the above UDP limitations of Octavia, Kuryr is forcing pods to use TCP for DNS resolution (`use-vc` option at `resolv.conf`) if the OpenStack version is not Train or newer. This may be a problem for Go applications compiled without CGO support (i.e. `CGO_ENABLED=0`) as the native Go resolver is using UDP only and is not considering the `use-vc` option added by Kuryr to the `resolv.conf` for `GO` version 1.12 or earlier. This is a problem also for musl-based containers as its resolver does not support `use-vc` option. This would include e.g., images build from `alpine`.
* Kuryr uses Octavia loadbalancer for OpenShift services. If Octavia Amphora driver is used, that means for each OpenShift service an Amphora VM is generated. If the Availability Zone option is used note, due to Octavia limitations, the Load Balancer VMs will not fulfill the Availability Zone restrictions and can be scheduled in any available AZ.

