# Install using oVirt platform provider 

## Overview

This provider enables the OpenShift Installer to provision VM resources in an 
oVirt data center, that will be used as worker and masters of the clusters. It 
will also create the bootstrap machine, and the configuration needed to get 
the initial cluster running by supplying DNS a service and load balancing, all 
using static pods. 
This work is related to the Bare-Metal provider because oVirt does not supply 
DNS and LB services but is a platform provider. See also [OpenShift-Metal³ kni-installer](https://github.com/openshift-metal3/kni-installer)
 

## Prerequisite

1. oVirt/RHV version 4.3.9.4 or later. 
2. Allocate 3 IP on the VM network: 
    - IP for the internal kubernetes api, that all components will interact with 
    - IP for the internal DNS service, to bootstrap etcd and to resolve names like 
    `api.$CLUSTER_NAME.$CLUSTER_DOMAIN` and node names 
    - IP for the Ingress, the load balancer in front of the cluster apps 
    To work with this provider one must supply 2 IPs that are related to any MAC 
    in the virtualization env, where the cluster will run. Those IPs will be active 
    by keepalived, on, initially the bootstrap machine, and then the masters, after 
    a fail-over, when the bootstrap is killed. 
    Locate those IP's in the target network. If you want the network details, go to 
    oVirt's webadmin and look for the designated cluster details and its networks. 
    One way to check if an IP is in use is to check if it has ARP associated with it 
    perform this check while on one of the hosts that would run the VMs: 
       ```console
       $ arp 10.35.1.19
       10.35.1.1 (10.35.1.1) -- no entry
       ```
3. Name resolution of `api_vip` from your installing machine 
The installer must resolve the `api_vip` during the installation, as it will 
interact with the API to follow the cluster version progress. 
To make it work prepend the `dns_vip` to your `/etc/resolv.conf` 
    ```
    search example.org
    nameserver $dns_vip
    ```


## Minimum resources

The default master/worker:
- 4 CPUs
- 16 RAM
- 120 GB disk

For 3 masters/3 workers, the target oVirt cluster must have at least:
- 96RAM
- 24vCPUs
- 720GiB storage
- 50IOPS for the Master VM disks, [per ETCD requirement doc](https://github.com/etcd-io/etcd/blob/master/Documentation/op-guide/hardware.md)

> Worker count can be reduced to 2 in `install-config.yaml` in case needed.

The cluster will create by default 1 bootstrap, 3 master, and 3 workers machines. 
By the time the first worker is up the bootstrap VM should be destroyed, and this 
is included in the minimum resources calculation.


## Install 


### Install using the wizard 

At this stage the installer can create a cluster by gathering all the information 
using a wizard:
```console
$ openshift-install create cluster --dir=install_dir
? SSH Public Key /home/user/.ssh/id_dsa.pub
? Platform ovirt
? Enter oVirt's api endpoint URL https://ovirt-engine-fqdn/ovirt-engine/api
? Is the installed oVirt certificate trusted? Yes
? Enter ovirt-engine username admin@internal
? Enter password ***
? oVirt cluster xxxx
? oVirt storage xxxx
? oVirt network xxxx
? Internal API virtual IP 10.0.0.1
? Internal DNS virtual IP 10.0.0.2
? Ingress virtual IP 10.0.0.3
? Base Domain example.org
? Cluster Name test
? Pull Secret [? for help]
INFO Consuming Install Config from target directory
INFO Creating infrastructure resources...
INFO Waiting up to 20m0s for the Kubernetes API at https://api.test.example.org:6443...
INFO API v1.17.1 up
INFO Waiting up to 40m0s for bootstrapping to complete...
INFO Destroying the bootstrap resources...
INFO Waiting up to 30m0s for the cluster at https://api.test.example.org:6443 to initialize...
INFO Waiting up to 10m0s for the openshift-console route to be created...
INFO Install complete!
INFO To access the cluster as the system:admin user when using 'oc', run 'export KUBECONFIG=/home/user/install_dir/auth/kubeconfig'
INFO Access the OpenShift web-console here: https://console-openshift-console.apps.test.example.org
INFO Login to the console with user: kubeadmin, password: xxxxxxxxx
```


### Install in stages when customization is needed 

Start the installation by creating an `install-config` interactively, using a work-dir:

```console
$ openshift-install create install-config --dir=install_dir
``` 

The resulting `install_dir/install-config.yaml` can be further customized if needed.
For general customization please see [docs/user/customization.md](../customization.md#platform-customization)
For ovirt-specific see [customization.md](./customization.md) 
Continue the installation using the install-config in the new folder `install_dir`

```console
$ openshift-install create cluster --dir=install_dir
``` 

When the all prompts are done the installer will create 3 VMs under the oVirt
cluster supplied, and another VM as the bootstrap node. 
The bootstrap will perform ignition fully and will advertise the IP in the
pre-login msg. Go to oVirt webadmin UI, and open the console of the bootstrap
VM to get it. 
In the end the installer finishes and the cluster should be up.

To access the cluster as the system:admin user: 

```console
$ export KUBECONFIG=$PWD/install_dir/auth/kubeconfig
$ oc get nodes
```

