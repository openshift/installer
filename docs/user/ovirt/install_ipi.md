![oVirt Logo](./images/oVirt-logo.png#center)

**Table of Contents**

- [Install using oVirt platform provider](#install-using-ovirt-platform-provider)
  * [Overview](#overview)
  * [Prerequisite](#prerequisite)
  * [Minimum resources](#minimum-resources)
  * [Install](#install)
    + [Minimum permission for installation](#minimum-permission-for-installation)
    + [ovirt-config.yaml](#ovirt-configyaml)
    + [ovirt-credentials](#ovirt-credentials)
    + [Bootstrap VM](#bootstrap-vm)
    + [Install using the wizard](#install-using-the-wizard)
    + [Install in stages when customization is needed](#install-in-stages-when-customization-is-needed)

# Install using oVirt platform provider 

## Overview

This provider enables the OpenShift Installer to provision VM resources in an 
oVirt data center, that will be used as worker and masters of the clusters. It 
will also create the bootstrap machine, and the configuration needed to get 
the initial cluster running by supplying DNS a service and load balancing, all 
using static pods. 
This work is related to the Bare-Metal provider because oVirt does not supply 
DNS and LB services but is a platform provider. See also [Bare Metal IPI Networking Infrastructure]
 

## Prerequisite

1. oVirt/RHV version 4.3.9.4 or later. 
2. Allocate 2 IP on the VM network:
    - IP for the internal kubernetes api, that all components will interact with 
    - IP for the Ingress, the load balancer in front of the cluster apps 
    To work with this provider one must supply 2 IPs excluded from the MAC range
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


## Minimum resources

The default master/worker:
- 4 CPUs
- 16 RAM
- 120 GB disk

For 3 masters/3 workers, the target Cluster **must have at least**:
- 96RAM
- 24vCPUs
- 720GiB storage
- Storage that is fast enough for etcd, [using-fio-to-tell-whether-your-storage-is-fast-enough-for-etcd](https://www.ibm.com/cloud/blog/using-fio-to-tell-whether-your-storage-is-fast-enough-for-etcd)

> Worker count can be reduced to 2 in `install-config.yaml` in case needed.

The cluster will create by default 1 bootstrap, 3 master, and 3 workers machines. 
By the time the first worker is up the bootstrap VM should be destroyed, and this 
is included in the minimum resources calculation.


## Install 

### Minimum permission for installation

It's **not recommended** to users use admin@internal during the installation. Instead, create an exclusive user to install and manage OCP on oVirt.

The minimum permissions are:
- DiskOperator
- DiskCreator
- UserTemplateBasedVm
- TemplateOwner
- TemplateCreator
- ClusterAdmin  (on the specific cluster targeted for OCP deployment)

There is an [ansible playbook available](https://github.com/oVirt/ocp-on-ovirt/tree/master/installer-tools/ocpadmin) which helps to setup an internal user and group with the minimum privileges to run the openshift-install on oVirt.

### ovirt-config.yaml

The ovirt-config.yaml is created under ${HOME}/.ovirt directory by the installer.
It contains all information how the installer connects to oVirt and can be re-used
if required to re-trigger a new installation.

Below the description of all config options in ovirt-config.yaml.

| Name           | Value                          | Type     | Example                                                                                                |
| ---------------|:------------------------------:|:--------:|:------------------------------------------------------------------------------------------------------:|
| ovirt_url      | URL for Engine API             | string   | https://engine.fqdn.home/ovirt-engine/api                                                              |
| ovirt_fqdn     | Engine FQDN                    | string   | engine.fqdn.home                                                                                       |
| ovirt_username | User to connect with Engine    | string   | admin@internal                                                                                         |
| ovirt_password | Password for the user provided | string   | superpass                                                                                              |
| ovirt_insecure | TLS verification disabled      | boolean  | false                                                                                                  |
| ovirt_ca_bundle| CA Bundle                      | string   | -----BEGIN CERTIFICATE----- MIIDvTCCAqWgAwIBAgICEAA.... ----- END CERTIFICATE -----                    |
| ovirt_pem_url  | PEM URL                        | string   | https://engine.fqdn.home/ovirt-engine/services/pki-resource?resource=ca-certificate&format=X509-PEM-CA |

### ovirt-credentials
During installation ${HOME}/.ovirt/ovirt-config.yaml is converted to a **secret** named as **ovirt-credentials**
and every openshift component with permission can use it.

$ oc get secrets --all-namespaces | grep ovirt-credentials
```
kube-system                                        ovirt-credentials
openshift-machine-api                              ovirt-credentials
```

$ oc get secret ovirt-credentials -n kube-system -o yaml
$ oc get secret ovirt-credentials -n openshift-machine-api -o yaml
```
apiVersion: v1
data:
  ovirt_ca_bundle: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUR2VENDFTE1Ba0dBMVV...
  ovirt_cafile: ""
  ovirt_insecure: Zm2U=
  ovirt_password: cmGF0
  ovirt_url: aHR0cHM3Z5lL2FwaQ==
  ovirt_username: YWRtaJuYWw=
kind: Secret
metadata:
  creationTimestamp: "2020-07-30T15:03:06Z"
  managedFields:
  - apiVersion: v1
    fieldsType: FieldsV1
    fieldsV1:
      f:data:
        .: {}
        f:ovirt_ca_bundle: {}
        f:ovirt_cafile: {}
        f:ovirt_insecure: {}
        f:ovirt_password: {}
        f:ovirt_url: {}
        f:ovirt_username: {}
      f:type: {}
    manager: cluster-bootstrap
    operation: Update
    time: "2020-07-30T15:03:06Z"
  name: ovirt-credentials
  namespace: kube-system
  resourceVersion: "94"
  selfLink: /api/v1/namespaces/kube-system/secrets/ovirt-credentials
  uid: 642dbc91-12eb-4111-baa7-d79cbc9b79e4
type: Opaque
```

### Bootstrap VM

The bootstrap will perform ignition fully and will advertise the IP in the
pre-login msg. Go to Engine webadmin UI, and open the console of the bootstrap
VM to get it.


### Install using the wizard 

At this stage the installer can create a cluster by gathering all the information 
using a wizard:
```console
$ openshift-install create cluster --dir=install_dir
? SSH Public Key /home/user/.ssh/id_dsa.pub
? Platform ovirt
? Engine FQDN[:PORT] [? for help] ovirt-engine-fqdn
? Enter ovirt-engine username admin@internal
? Enter password ***
? oVirt cluster xxxx
? oVirt storage xxxx
? oVirt network xxxx
? Internal API virtual IP 10.0.0.1
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

When the all prompts are done the installer will create ${HOME}/.ovirt/ovirt-config.yaml
containing all required information about the connection with Engine.
The installation process will create a temporary VM which will trigger bootstrap VM
for later create three masters nodes. The masters nodes will create all services and
checks required. Finally, the cluster will create the three workers node.

In the end the installer finishes and the cluster should be up.

To access the cluster as the system:admin user: 

```console
$ export KUBECONFIG=$PWD/install_dir/auth/kubeconfig
$ oc get nodes
```

[Bare Metal IPI Networking Infrastructure]: https://github.com/openshift/installer/blob/master/docs/design/baremetal/networking-infrastructure.md

#### Installing OpenShift on RHV/oVirt in *insecure* mode

<!-- Do not change this title as it is used in the code to point users to the right place -->

Starting OpenShift 4.7 we are sunsetting the “insecure” option from the OpenShift Installer. Starting with this version, the installer only supports installation methods from the user interface that lead to using verified certificates.

This change also means that setting up the CA certificate for RHV is no longer required before running the installer. The installer will ask you for confirmation about the certificate and store the CA certificate for use during the installation.

Should you, nevertheless, require an installation without certificate verification you can create a file named ovirt-config.yaml in the .ovirt directory in your home directory (~/.ovirt/ovirt-config.yaml) before running the installer with the following content:

```yaml
ovirt_url: https://ovirt.example.com/ovirt-engine/api
ovirt_fqdn: ovirt.example.com
ovirt_pem_url: ""
ovirt_username: admin@internal
ovirt_password: super-secret-password
ovirt_insecure: true
```

Please note that this option is **not recommended** as it will allow a potential attacker to perform a Man-in-the-Middle attack and capture sensitive credentials on the network.
