# OpenShift Installer 

## Supported Platforms

* [AWS](docs/user/aws/) ([Official Docs](https://docs.openshift.com/container-platform/latest/installing/installing_aws/preparing-to-install-on-aws.html))
* [Azure](docs/user/azure/) ([Official Docs](https://docs.openshift.com/container-platform/latest/installing/installing_azure/preparing-to-install-on-azure.html))
* [Bare Metal](docs/user/metal/) ([Official Docs](https://docs.openshift.com/container-platform/latest/installing/installing_bare_metal/preparing-to-install-on-bare-metal.html))
* [GCP](docs/user/gcp/) ([Official Docs](https://docs.openshift.com/container-platform/latest/installing/installing_gcp/preparing-to-install-on-gcp.html))
* IBM Cloud ([Official Docs](https://docs.openshift.com/container-platform/latest/installing/installing_ibm_cloud/preparing-to-install-on-ibm-cloud.html))
* Nutanix ([Official Docs](https://docs.openshift.com/container-platform/latest/installing/installing_nutanix/preparing-to-install-on-nutanix.html))
* [OpenStack](docs/user/openstack/) ([Official Docs](https://docs.openshift.com/container-platform/latest/installing/installing_openstack/preparing-to-install-on-openstack.html))
* [Power](docs/user/power/) ([Official Docs](https://docs.openshift.com/container-platform/latest/installing/installing_ibm_power/preparing-to-install-on-ibm-power.html))
* Power VS ([Official Docs](https://docs.openshift.com/container-platform/latest/installing/installing_ibm_powervs/preparing-to-install-on-ibm-power-vs.html))
* [vSphere](docs/user/vsphere/) ([Official Docs](https://docs.openshift.com/container-platform/latest/installing/installing_vsphere/preparing-to-install-on-vsphere.html))
* [z/VM](docs/user/zvm/) ([Official Docs](https://docs.openshift.com/container-platform/latest/installing/installing_ibm_z/preparing-to-install-on-ibm-z.html))

## Quick Start

First, install all [build dependencies](docs/dev/dependencies.md).

Clone this repository. Then build the `openshift-install` binary with:

```sh
hack/build.sh
```

This will create `bin/openshift-install`. This binary can then be invoked to create an OpenShift cluster, like so:

```sh
bin/openshift-install create cluster
```

The installer will show a series of prompts for user-specific information and use reasonable defaults for everything else.
In non-interactive contexts, prompts can be bypassed by [providing an `install-config.yaml`](docs/user/overview.md#multiple-invocations).

If you have trouble, refer to [the troubleshooting guide](docs/user/troubleshooting.md).

### Connect to the cluster

Details for connecting to your new cluster are printed by the `openshift-install` binary upon completion, and are also available in the `.openshift_install.log` file.

Example output:

```sh
INFO Waiting 10m0s for the openshift-console route to be created...
INFO Install complete!
INFO To access the cluster as the system:admin user when using 'oc', run
    export KUBECONFIG=/path/to/installer/auth/kubeconfig
INFO Access the OpenShift web-console here: https://console-openshift-console.apps.${CLUSTER_NAME}.${BASE_DOMAIN}:6443
INFO Login to the console with user: kubeadmin, password: 5char-5char-5char-5char
```

### Cleanup

Destroy the cluster and release associated resources with:

```sh
openshift-install destroy cluster
```

Note that you almost certainly also want to clean up the installer state files too, including `auth/`, `terraform.tfstate`, etc.
The best thing to do is always pass the `--dir` argument to `create` and `destroy`.
And if you want to reinstall from scratch, `rm -rf` the asset directory beforehand.
