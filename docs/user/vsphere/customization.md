# vSphere Platform Customization

Beyond the [platform-agnostic `install-config.yaml` properties](../customization.md#platform-customization), the installer supports additional, vSphere-specific properties.

## Cluster-scoped properties

* `vCenter` (required string): The domain name or IP address of the vCenter.
* `username` (required string): The username to use to connect to the vCenter.
* `password` (required string): The password to use to connect to the vCenter.
* `datacenter` (required string): The name of the datacenter to use in the vCenter.
* `defaultDatastore` (required string): The default datastore to use for provisioning volumes.
* `folder` (optional string): The absolute path of an existing folder where the installer should create VMs. The absolute path is of the form `/example_datacenter/vm/example_folder/example_subfolder`. If a value is specified, the folder must exist. If no value is specified, a folder named with the cluster ID will be created in the `datacenter` VM folder.
* `resourcePool` (optional string): The absolute path of an existing resource pool where the installer should create VMs. The absolute path is of the form `/example_datacenter/host/example_cluster/Resources/example_resource_pool/optionally_sub_resource_pool`. If a value is specified, the resource pool must exist. If no value is specified, resources will be installed in the root of the cluster `/example_datacenter/host/example_cluster/Resources`.

## Machine pools

* `osDisk` (optional object):
    * `diskSizeGB` (optional integer): The size of the disk in gigabytes (GB).
* `cpus` (optional integer): The total number of virtual processor cores to assign a vm.
* `coresPerSocket` (optional integer): The number of cores per socket in a vm. The number of vCPUs on the vm will be cpus/coresPerSocket (default is 1).
* `memoryMB` (optional integer): The size of a VM's memory in megabytes.
* `disk_type` (optional string): DiskType is the name of the disk provisioning type, valid values are thin, thick, and eagerZeroedThick. When not specified, it will be set according to the default storage policy of vsphere.

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

## vCenter credential types

`platform.vSphere.credentialType` accepts `global` or `component-scoped`. If it is omitted, `global` is used for backward compatibility.

With `global`, each `vcenters` entry must contain `user` and `password`; `componentCredentials` must not be set. These credentials are used by the installer and by the generated vSphere credential Secret.

With `component-scoped`, `user` and `password` must be omitted from every `vcenters` entry. Each entry must contain complete credentials for `machineManagement`, `storage`, `cloudControllerManager`, and `vsphereProblemDetector`:

```yaml
platform:
  vsphere:
    credentialType: component-scoped
    vcenters:
    - server: vcenter.example.com
      datacenters:
      - DC1
      componentCredentials:
        machineManagement:
          user: ocp-machine-api@vsphere.local
          password: machine-api-password
        storage:
          user: ocp-csi@vsphere.local
          password: csi-password
        cloudControllerManager:
          user: ocp-ccm@vsphere.local
          password: ccm-password
        vsphereProblemDetector:
          user: ocp-diagnostics@vsphere.local
          password: diagnostics-password
```

The `machineManagement` credential is used for installer-side vCenter operations, including VM provisioning and cleanup. Its privileges must therefore include the installer provisioning privileges described in [Privileges](privileges.md). Each component credential must have the privileges required by its corresponding OpenShift component. Credentials are stored in the install configuration and generated Secrets; protect the install-config directory accordingly.

The installer creates these four Secrets for component-scoped manual installations:

| Component | Secret | Namespace |
|---|---|---|
| Machine API | `vsphere-cloud-credentials` | `openshift-machine-api` |
| vSphere CSI | `vmware-vsphere-cloud-credentials` | `openshift-cluster-csi-drivers` |
| Cloud Controller Manager | `vsphere-cloud-credentials` | `openshift-cloud-controller-manager` |
| Problem Detector | `vsphere-cloud-credentials` | `openshift-cluster-storage-operator` |

Each Secret contains `<vcenter-server>.username` and `<vcenter-server>.password` keys for every configured vCenter, with the credentials selected for that component. These Secrets are intended for `credentialsMode: Manual`; the Cloud Credential Operator does not reconcile them in that mode.
