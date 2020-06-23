# OpenShift Installer

## Supported Platforms

* [AWS](docs/user/aws/README.md)
* [AWS (UPI)](docs/user/aws/install_upi.md)
* [Azure](docs/user/azure/README.md)
* [Bare Metal (UPI)](docs/user/metal/install_upi.md)
* [Bare Metal (IPI) (Experimental)](docs/user/metal/install_ipi.md)
* [GCP](docs/user/gcp/README.md)
* [GCP (UPI)](docs/user/gcp/install_upi.md)
* [Libvirt with KVM](docs/dev/libvirt/README.md) (development only)
* [OpenStack](docs/user/openstack/README.md)
* [OpenStack (UPI) (Experimental)](docs/user/openstack/install_upi.md)
* [oVirt](docs/user/ovirt/install_ipi.md)
* [vSphere](docs/user/vsphere/README.md)
* [vSphere (UPI)](docs/user/vsphere/install_upi.md)

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
INFO To access the cluster as the system:admin user when using 'oc', run 'export KUBECONFIG=/path/to/installer/auth/kubeconfig'
INFO Access the OpenShift web-console here: https://console-openshift-console.apps.${CLUSTER_NAME}.${BASE_DOMAIN}:6443
INFO Login to the console with user: kubeadmin, password: 5char-5char-5char-5char
```

### Cleanup

Destroy the cluster and release associated resources with:

```sh
openshift-install destroy cluster
```

Note that you almost certainly also want to clean up the installer state files too, including `auth/`, `terraform.tfstate`, etc.
The best thing to do is always pass the `--dir` argument to `install` and `destroy`.
And if you want to reinstall from scratch, `rm -rf` the asset directory beforehand.
