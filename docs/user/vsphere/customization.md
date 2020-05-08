# vSphere Platform Customization

Beyond the [platform-agnostic `install-config.yaml` properties](../customization.md#platform-customization), the installer supports additional, vSphere-specific properties.

## Cluster-scoped properties

* `vCenter` (required string): The domain name or IP address of the vCenter.
* `username` (required string): The username to use to connect to the vCenter.
* `password` (required string): The password to use to connect to the vCenter.
* `datacenter` (required string): The name of the datacenter to use in the vCenter.
* `defaultDatastore` (required string): The default datastore to use for provisioning volumes.
* `folder` (optional string): The absolute path of an existing folder where the installer should create VMs. The absolute path is of the form `/example_datacenter/vm/example_folder/example_subfolder`. If a value is specified, the folder must exist. If no value is specified, a folder named with the cluster ID will be created in the `datacenter` VM folder.

## Machine pools

* `osDisk` (optional object):
    * `diskSizeGB` (optional integer): The size of the disk in gigabytes (GB).
* `cpus` (optional integer): The total number of virtual processor cores to assign a vm.
* `coresPerSocket` (optional integer): The number of cores per socket in a vm. The number of vCPUs on the vm will be cpus/coresPerSocket (default is 1).
* `memoryMB` (optional integer): The size of a VM's memory in megabytes.

## Examples

Some example `install-config.yaml` are shown below.
For examples of platform-agnostic configuration fragments, see [here](../customization.md#examples).

### Minimal

An example minimal vSphere install config is:

```yaml
apiVersion: v1
baseDomain: example.com
metadata:
  name: test-cluster
platform:
  vSphere:
    vCenter: your.vcenter.example.com
    username: username
    password: password
    datacenter: datacenter
    defaultDatastore: datastore
pullSecret: '{"auths": ...}'
sshKey: ssh-ed25519 AAAA...
```

### Custom Machine Pools

An example vSphere install config with custom machine pools:
```yaml
apiVersion: v1
baseDomain: example.com
controlPlane:
  name: master
  platform:
    vsphere:
      cpus: 8
      coresPerSocket: 2
      memoryMB: 24576
      osDisk:
        diskSizeGB: 512
  replicas: 3
compute:
- name: worker
  platform:
    vsphere:
      cpus: 8
      coresPerSocket: 2
      memoryMB: 24576
      osDisk:
        diskSizeGB: 512
  replicas: 5
metadata:
  name: test-cluster
platform:
  vSphere:
    vCenter: your.vcenter.example.com
    username: username
    password: password
    datacenter: datacenter
    defaultDatastore: datastore
pullSecret: '{"auths": ...}'
sshKey: ssh-ed25519 AAAA...
```
