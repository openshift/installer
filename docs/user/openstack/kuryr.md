# Kuryr

Kuryr is a CNI plug-in that uses Neutron and Octavia to provide networking for pods and services. It is primarily designed for OpenShift clusters that run on OpenStack virtual machines. Kuryr improves the network performance by plugging
OpenShift pods into OpenStack SDN. In addition, it provides interconnectivity between OpenShift pods and OpenStack virtual instances.

Kuryr is recommended for OpenShift deployments on encapsulated OpenStack tenant networks in order to avoid double encapsulation, such as running an encapsulated OpenShift SDN over an OpenStack network.

Conversely, using Kuryr does not make sense in the following cases:

* You use provider networks or tenant VLANs.
* The deployment will use many services on a few hypervisors. Each OpenShift service creates an Octavia Amphora virtual machine in OpenStack that hosts a required load balancer.
* UDP services are needed.

[Learn more about Kuryr](https://docs.openstack.org/kuryr-kubernetes/latest/).

## Requirements when enabling Kuryr

When using Kuryr SDN, as the pods, services, namespaces, network policies, etc., are using resources from the OpenStack Quota, the minimum requirements are higher. Kuryr also has some additional requirements on top of what a default install requires. At a bare minimum you need the following quota to run a *default* cluster:

* Floating IPs: 3 (plus the expected number of services of LoadBalancer type)
* Security Groups: 100 (1 needed per network policy)
* Security Group Rules: 500
* Routers: 1
* Subnets: 100 (1 needed per namespace)
* Networks: 100 (1 needed per namespace)
* Ports: 1000
* RAM: 112 GB
* vCPUs: 28
* Volume Storage: 175 GB
* Instances: 7

## Increase Quota

As highlighted in the minimum quota recommendations, when using Kuryr SDN, there is a need for increasing the quotas as pods, services, namespaces, network policies are using OpenStack resources. So, as an administrator,the next quotas should be increased for the selected project:

```sh
openstack quota set --secgroups 100 --secgroup-rules 500 --ports 500 --subnets 100 --networks 100 <project>
```

## Neutron Configuration

Kuryr CNI makes use of the Neutron Trunks extension to plug containers into the OpenStack SDN, so the `trunks` extension must be enabled for Kuryr to properly work.

In addition, if the default ML2/OVS Neutron driver is used, the firewall must be set to `openvswitch` instead of `ovs_hybrid` so that security groups are enforced on trunk subports and Kuryr can properly handle Network Policies.

## Octavia

Kuryr SDN uses OpenStack Octavia LBaaS to implement OpenShift services. Thus the OpenStack environment must have Octavia components installed and configured if Kuryr SDN is used.

**NOTE:** Depending on your OpenStack environment Octavia may not support UDP listeners, which means there is no support for UDP services if Kuryr SDN is used.

## Installing with Kuryr SDN

To deploy with Kuryr SDN instead of the default OpenShift SDN, you simply need to modify the `install-config.yaml` file to include `Kuryr` as the desired `networking.networkType` and proceed with the same steps as with the default OpenShift SDN:

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

**NOTE:** both `trunkSupport` and `octaviaSupport` are automatically discovered by the installer, so there is no need to set them. But if your environment doesn't meet both requirements Kuryr SDN will not properly work, as trunks are needed to connect the pods to the OpenStack network and Octavia to create the OpenShift services.

### Known limitations of installing with Kuryr SDN

There are known limitations when using Kuryr SDN:

* There is an amphora load balancer VM being deployed per OpenShift service with the default Octavia load balancer driver (amphora driver). If the environment is resource constrained it could be a problem to create a large amount of services as each one will create a new VM.
* Depending on the Octavia version, UDP listeners are not supported. This means that OpenShift UDP services are not supported unless OpenStack version is Train or newer.
* Depending on the Octavia version, there is a known limitation of Octavia not supporting listeners on different protocols (e.g., UDP and TCP) on the same port. Thus services exposing the same port for different protocols are not supported unless OpenStack version is Train or newer.
* Due to the above UDP limitations of Octavia, Kuryr is forcing pods to use TCP for DNS resolution (`use-vc` option at `resolv.conf`) if the OpenStack version is not Train or newer. This may be a problem for pods running Go applications compiled with `CGO_DEBUG` flag disabled as that forces to use the `go` resolver that is only using UDP and is not considering the `use-vc` option added by Kuryr to the `resolv.conf`. This is a problem also for musl-based containers as it's resolver does not support `use-vc` option. This would include e.g., images build from `alpine`.
