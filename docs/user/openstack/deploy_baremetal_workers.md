# Deploying OpenShift bare-metal workers on OpenStack cloud provider


## Table of Contents
- [Deploying OpenShift bare-metal workers on OpenStack cloud provider](#deploying-openshift-bare-metal-workers-on-openstack-cloud-provider)
  - [Table of Contents](#table-of-contents)
  - [Common prerequisites](#common-prerequisites)
  - [Considerations when deploying bare-metal workers](#considerations-when-deploying-bare-metal-workers)
  - [Deploying cluster with BM workers on tenant network deployed by the installer](#deploying-cluster-with-bm-workers-on-tenant-network-deployed-by-the-installer)
  - [Deploying cluster with BM workers on preexisting network](#deploying-cluster-with-bm-workers-on-preexisting-network)
  - [Deploying cluster with BM machines only on preexisting network](#deploying-cluster-with-bm-machines-only-on-preexisting-network)
  - [Known issues](#known-issues)


## Common prerequisites

* [Ironic bare metal service][2] is enabled and accessible through the [OpenStack Compute API][3].
* Bare-metal machines are available as an [OpenStack flavor][1].
* At the time bare-metal machines are booted, the image used to boot them is available from the [Glance OpenStack image service][4].
If given a URL, the installer will fetch that image from the location specified and upload it to Glance. See relevant
[documentation][5].


## Considerations when deploying bare-metal workers

* Long boot times for bare-metal machines may result in deployment timeout

    The first time a bare-metal server is booted the boot time could be significantly longer that
    the first time a VM server is booted. This longer boot time could exceed the timeout settings of
    the installer. If that occurs, the deployment will fail with a timeout error.  The deployment maybe
    restarted and completed by re-running  the installer with the appropriate *wait-for*
    command. For example:

        ./openshift-install wait-for install-complete --log-level debug


## Deploying cluster with BM workers on tenant network deployed by the installer

The initial cluster is deployed using bare-metal workers only. Bare-metal workers and control plane VMs are all
attached to the tenant network provisioned by the installer.

- Requirements:
    - By default, OpenStack networks support attaching both VMs and bare-metal machines to them.
    * [Ironic bare metal service][2] can listen for, and PXE-boot machines in tenant networks

- Create install-config.yaml:

    - Set `compute.[worker].platform.openstack.type` to the bare-metal server flavor.
    - Set `controlPlane.platform.openstack.type` to the VM flavor which will be used by the control plane nodes.

        For example:

                controlPlane:
                   platform:
                     openstack:
                       type: <vmComputeFlavorForMasters>

                 ... other settings

                 compute:
                 - architecture: amd64
                   hyperthreading: Enabled
                   name: worker
                   platform:
                     openstack:
                       type: <baremetalComputeFlavor>
                   replicas: 3

                 ... other settings



- Run the openshift installer:

      ./openshift-install create cluster --log-level debug

- wait for the installer to complete.

    If the installer times out waiting for the bare-metal workers to complete booting,
    restart the installation using [the appropriate *wait-for* command](#considerations-when-deploying-bare-metal-workers).

## Deploying cluster with BM workers on preexisting network

Initial cluster is deployed using bare-metal workers only. Bare-metal workers are attached to
a preexisting network.

- Requirements:
    - An OpenStack subnet has been pre-provisioned which supports attaching both VMs and bare-metal servers to it.

- Create install-config.yaml:

    - Set `compute.[worker].platform.openstack.type` to the bare-metal server flavor.

    - Set `controlPlane.platform.openstack.type` to the VM flavor which will be used by the control plane nodes.

    - Set `platform.openstack.machinesSubnet` to the UUID of the pre-provisioned subnet.

        For example:

                controlPlane:
                   platform:
                     openstack:
                       type: <vmComputeFlavorForMasters>

                 ... other settings

                 compute:
                 - architecture: amd64
                   hyperthreading: Enabled
                   name: worker
                   platform:
                     openstack:
                       type: <baremetalComputeFlavor>
                   replicas: 3

                 ... other settings

                 platform:
                   openstack:
                     machinesSubnet: <uuidOfPreprovisionedSubnet>

- Run the openshift installer:

        ./openshift-install create cluster --log-level debug

- wait for the installer to complete.

    If the installer times out waiting for the bare-metal workers to complete booting,
    restart the installation using [the appropriate *wait-for* command](#considerations-when-deploying-bare-metal-workers).

## Deploying cluster with BM machines only on preexisting network

Initial cluster is deployed using bare-metal machines only for both Control
Plane and Compute nodes. The cluster is deployed to a preexisting network.

- Requirements:
    - An OpenStack subnet has been pre-provisioned which supports attaching bare-metal servers to it.

- Create install-config.yaml:

    - Set `compute.[worker].platform.openstack.type` to a bare-metal server flavor.

    - Set `controlPlane.platform.openstack.type` to a bare-metal server flavor.

    - Set `platform.openstack.machinesSubnet` to the UUID of the pre-provisioned subnet.

        For example:

                controlPlane:
                   platform:
                     openstack:
                       type: <bmComputeFlavorForMasters>

                 ... other settings

                 compute:
                 - architecture: amd64
                   hyperthreading: Enabled
                   name: worker
                   platform:
                     openstack:
                       type: <bmComputeFlavorForWorkers>
                   replicas: 3

                 ... other settings

                 platform:
                   openstack:
                     machinesSubnet: <uuidOfPreprovisionedSubnet>

- Run the openshift installer:

        ./openshift-install create cluster --log-level debug

- wait for the installer to complete.

    If the installer times out waiting for the bare-metal workers to complete booting,
    restart the installation using [the appropriate *wait-for* command](#considerations-when-deploying-bare-metal-workers).

## Known issues

Bare metal nodes are not supported on clusters that use Kuryr.


[1]: <https://docs.openstack.org/nova/latest/user/flavors.html> "In OpenStack, flavors define the compute, memory, and storage capacity of nova computing instances"
[2]: <https://docs.openstack.org/ironic/latest/>
[3]: <https://docs.openstack.org/api-ref/compute/>
[4]: <https://docs.openstack.org/glance/latest/>
[5]: <https://github.com/openshift/installer/blob/master/docs/user/openstack/customization.md#image-overrides>
