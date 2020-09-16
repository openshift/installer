# GCP Platform Customization

Beyond the [platform-agnostic `install-config.yaml` properties](../customization.md#platform-customization), the installer supports additional, GCP-specific properties.

* `projectID` (required string): The project where the cluster should be created.
* `region` (required string): The GCP region where the cluster should be created.
* `network` (optional string): The name of an existing GCP VPC where the cluster infrastructure should be provisioned.
* `controlPlaneSubnet` (optional string): The name of an existing GCP subnet which should be used by the cluster control plane.
* `computeSubnet` (optional string): The name of an existing GCP subnet which should be used by the cluster nodes.
* `defaultMachinePlatform` (optional object): Default [GCP-specific machine pool properties](#machine-pools) which apply to [machine pools](../customization.md#machine-pools) that do not define their own GCP-specific properties.
* `licenses` (optional list of strings): A list of license URLs (https) that should be applied to the compute images (as defined in [the API][compute-images]). The use of this property in combination with any mechanism that results in using pre-built images (such as the current OPENSHIFT_INSTALL_OS_IMAGE_OVERRIDE) is forbidden. Also, note that use of these URLs will force the installer to copy the source image before being used. An example of this license is the one that enables [nested virtualization][gcp-nested]. A full list of available licenses can be retrieved using [the license API][license-api].
* `labels` (optional object): A map of key-value pairs to apply to resources that are created. Not all GCP resources support labels.
  Additionally, the following requirements apply:
    * Keys and values cannot be longer than be 63 characters each.
    * Keys and values can only contain lowercase letters, numeric characters, underscores, and hyphens. International characters are allowed.
    * Keys must start with a lowercase letter and cannot be empty.

## Machine pools

* `type` (optional string): The [GCP machine type][machine-type].
* `zones` (optional array of strings): The availability zones used for machines in the pool.
* `osDisk` (optional object):
    * `diskSizeGB` (optional integer): The size of the disk in gigabytes (GB) (Minimum: 16GB, Maximum: 65536GB).
    * `diskType` (optional string): The type of disk (allowed values are: `pd-ssd`, and `pd-standard`. Default: `pd-ssd`).

## Installing to Existing Networks & Subnetworks

The installer can use an existing VPC and subnets when provisioning an OpenShift cluster. If one of `network`, `controlPlaneSubnet`, or `computeSubnet` is specified, all must be specified ([see example below](#pre-existing-networks--subnets)). Furthermore, each of the networks must belong to the project specified by `projectID`, and the subnets must belong to the specified cluster `region`. The installer will use these existing networks when creating infrastructure such as VM instances, load balancers, firewall rules, and DNS zones.

### Cluster Isolation

In a scenario where multiple clusters are installed to the same VPC network, the installer maintains cluster isolation by using firewall rules which specify allowed sources and destinations through network tags. By tagging each Compute Engine VM instance with a unique cluster id and creating corresponding firewall rules, the installer ensures that a cluster's control plane is only accessible by its own member nodes.

By design, possible inter-cluster access is limited to:
* The API, which is globally available with an external publishing strategy or available throughout the network in an internal publishing strategy
* Debugging tools; i.e. ports on VM instances are open to the `machineCidr` for SSH & ICMP

## Examples

Some example `install-config.yaml` are shown below.
For examples of platform-agnostic configuration fragments, see [here](../customization.md#examples).

### Minimal

An example minimal GCP install config is:

```yaml
apiVersion: v1
baseDomain: example.com
metadata:
  name: example-cluster
platform:
  gcp:
    project: example-project
    region: us-east1
    osDisk:
      diskType: pd-ssd
      diskSizeGB: 120
pullSecret: '{"auths": ...}'
sshKey: ssh-ed25519 AAAA...
```

### Custom machine pools

An example GCP install config with custom machine pools:

```yaml
apiVersion: v1
baseDomain: example.com
compute:
- name: worker
  platform: 
    gcp:
      type: n2-standard-2
      zones:
      - us-central1-a
      - us-central1-c
      osDisk:
        diskType: pd-standard
        diskSizeGB: 128
  replicas: 3
controlPlane:
  name: master
  platform:
    gcp:
      type: n2-standard-4
      zones:
      - us-central1-a
      - us-central1-c
      osDisk:
        diskType: pd-ssd
        diskSizeGB: 1024
  replicas: 3
metadata:
  name: example-cluster
platform:
  gcp:
    projectID: openshift-dev-installer
    region: us-central1
pullSecret: '{"auths": ...}'
sshKey: ssh-ed25519 AAAA...
```

### Pre-existing Networks & Subnets

An example GCP install config utilizing an existing network and subnets:

```yaml
apiVersion: v1
baseDomain: example.com
metadata:
  name: example-cluster
platform:
  gcp:
    projectID: example-project
    region: us-east1
    computeSubnet: example-worker-subnet
    controlPlaneSubnet: example-controlplane-subnet
    network: example-network
pullSecret: '{"auths": ...}'
sshKey: ssh-ed25519 AAAA...
```

### Nested virtualization

An example GCP install config enabling [GCP's nested virtualization license][gcp-nested]:

```yaml
apiVersion: v1
baseDomain: example.com
platform:
  gcp:
    projectID: example-project
    region: us-east1
    licenses:
    - https://compute.googleapis.com/compute/v1/projects/vm-options/global/licenses/enable-vmx
```

[machine-type]: https://cloud.google.com/compute/docs/machine-types
[compute-images]: https://cloud.google.com/compute/docs/reference/rest/v1/images
[gcp-nested]: https://cloud.google.com/compute/docs/instances/enable-nested-virtualization-vm-instances
[license-api]: https://cloud.google.com/compute/docs/reference/rest/v1/licenses/list
