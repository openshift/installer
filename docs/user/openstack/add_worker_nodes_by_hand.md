# Adding Compute Nodes By Hand
The preferred method of adding compute nodes to an OpenShift cluster is
to use the worker MachineSet with [Machine API Operator](https://github.com/openshift/machine-api-operator).

It is also possible to add a compute node through other tools like ansible. For example, 
the compute nodes may be added as described in [Installing OpenShift on OpenStack User-Provisioned Infrastructure][1];
or when scaling up the compute nodes using [Openshift-ansible][2].

When doing so, there are a few important points that should be taken into account. 

## Apply Necessary Security Groups To Ports
When adding a new compute, the necessary security groups need to be applied to 
allow the node access the cluster API endpoints and services.   

If the node is directly attached to the machineNetwork,  the pre-existing worker security groups maybe 
applied to the port attaching the node. The existing security rules should not require modifications.  

If the nodes is attached to a different network than the machineNetwork,
then new security group rules will have to be configured because the security group rules deployed by default on a cluster
match on both destination and source networks.


## Configure Allowed Address Pairs On Node Port
The port used by the new compute node needs to be configured to accept traffic for the ingress VIP.

The value for the ingress VIP can be found by running the command:

```
$ openstack port show <cluster-id>-ingress-port -f value -c fixed_ips
10.0.128.7
```

Configure the compute node's port to accept traffic for the ingress VIP by
setting the allowed address pairs.

```
$ openstack port set 97f73ba7-e104-49e2-ad7a-d4b440acc57e --allowed-address ip-address='10.0.128.7'

$ openstack port show 97f73ba7-e104-49e2-ad7a-d4b440acc57e -c allowed_address_pairs
   +-----------------------+----------------------------------------------------------+
   | Field                 | Value                                                    |
   +-----------------------+----------------------------------------------------------+
   | allowed_address_pairs | ip_address='10.0.128.7', mac_address='fa:16:3e:0d:35:f8' |
   +-----------------------+----------------------------------------------------------+
```



[1]: https://github.com/openshift/installer/blob/master/docs/user/openstack/install_upi.md#installing-openshift-on-openstack-user-provisioned-infrastructure
[2]: https://github.com/openshift/openshift-ansible

## Access To API Endpoints Using FQD Hostnames.
The node needs to have access to the cluster API endpoints using
Fully Qualified Domain hostnames.  The node's DNS configuration should resolve the hostname to an IP address which is accessible
to the node. Network security should be configured to allow access.

DNS resolution is a system administration task and deploying a compute node
by hand does not require a particular DNS configuration. For example, adding the necessary entries to /etc/hosts
file is one possible way. For other possible ways to configure DNS resolution
consult your operating system's documentation.

Before starting the deployment process
correct DNS configuration and access should be verified by using the appropriate tools.  For example, if the cluster is named
my.cluster.com, then the following _curl_ command maybe used to check the necessary access: 

```curl --insecure https://api.url.text:6443```
where ```api.url.text``` is a hostname that resolves to an IP address the node can use to reach the API endpoints.
```curl``` is an example and not the only tool that maybe used to verify access.