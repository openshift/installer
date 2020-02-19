# Install using oVirt platform provider

## Overview
This provider enables the Openshift Installer to provision VM resources in an \
oVirt data center, that will be used as worker and masters of the clusters. It \
will also create the bootstrap machine, and the configuration needed to get \
the initial cluster running by supplying DNS a service and load balancing, all \
using static pods.

This work is related to the Bare-Metal provider because oVirt does not supply
DNS and LB services but is a platform provider. See also [OpenShift-Metal³ kni-installer](https://github.com/openshift-metal3/kni-installer)

## Prerequisite
1. allocate 3 IP on the VM network:

   1.1 IP for the internal kubernetes api, that all components will interact with
   
   1.2 IP for the internal DNS service, to bootstrap etcd and to resolve names like
   `api.$CLUSTER_NAME.$CLUSTER_DOMAIN` and node names
   
   1.3 IP for the Ingress, the load balancer in front of the cluster apps
   
   To work with this provider one must supply 2 IPs that are related to any MAC
   in the virtualization env, where the cluster will run. Those IPs will be active
   by keepalived, on, initially the bootstrap machine, and then the masters, after
   a fail-over, when the bootstrap is killed.
   
   Locate those IP's in the target network. If you want the network details, go to
   oVirt's webadmin and look for the designated cluster details and its networks.
   
   One way to check if an IP is in use is to check if it has ARP associated with it 
   - perform this check while on one of the hosts that would run the VMs:
   ```console
   $ arp 10.35.1.19
   10.35.1.1 (10.35.1.1) -- no entry
   ``` 
2. RHCOS openstack image uploaded and ready as a template in oVirt
    The `openshift-installer` has an hardcoded reference to its relevant supported images.
    This information can be extracted easily and be used to download the qcow images
    and upload it into ovirt-engine.
    All of this is combined in that [ansible playbook](https://gist.github.com/rgolangh/adccf6d6b5eaecaebe0b0aeba9d3331b).  
    >Note: that is a temporary solution and the installer will probably upload it
    using the ovirt terraform provider.

    For upstream users the FCOS image can be used from the built-in glance
    provider - see this ansible role https://github.com/oVirt/ovirt-ansible-image-template
    Use it to make a template from the glance image.

3. Create the ~/.ovirt/ovirt-config.yaml file with contents similar to this.
   ```
   $ cat ~/.ovirt/ovirt-config.yaml 
   ovirt_url: https://[engine.name.fqdn]/ovirt-engine/api
   ovirt_insecure: true
   ovirt_username: admin@internal
   ovirt_password: somePassword
   ```   

## Install
Start the installation by creating an `install-config` interactively, using a work-dir:
```console
$ openshift-install create install-config --dir=cluster-prod
```

The resulting `cluster-prod/install-config.yaml` can be further customized if needed.

Continue the installation using the install-config in the new folder `cluster-prod`
```console
$ openshift-install create cluster --dir=cluster-prod
```

When the all prompts are done the installer will create 3 VMs under the oVirt
cluster supplied, and another VM as the bootstrap node.

The bootstrap will perform ignition fully and will advertise the IP in the
pre-login msg. Go to oVirt webadmin UI, and open the console of the bootstrap
VM to get it.

In the end the installer finishes and the cluster should be up.
To access the cluster as the system:admin user:
 
```console
$ export KUBECONFIG=$PWD/cluster-prod/auth/kubeconfig
$ oc get nodes
```

[image]: https://releases-rhcos.svc.ci.openshift.org/storage/releases/4.1/410.8.20190624.0/rhcos-410.8.20190624.0-openstack.qcow2
