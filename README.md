# Openshift Installer

## Supported Platforms

* AWS
* [Libvirt with KVM](docs/dev/libvirt-howto.md)
* OpenStack (Experimental)

## Quick Start

First, install all [build dependencies](docs/dev/dependencies.md).

Clone this repository to `src/github.com/openshift/installer` in your [GOPATH](https://golang.org/cmd/go/#hdr-GOPATH_environment_variable). Then build the `openshift-install` binary with:

```sh
hack/build.sh
```

This will create `bin/openshift-install`. This binary can then be invoked to create an OpenShift cluster, like so:

```sh
bin/openshift-install create cluster
```

The installer requires the terraform binary either alongside openshift-install or in `$PATH`.
If you don't have [terraform](https://www.terraform.io/), run the following to create `bin/terraform`:

```sh
hack/get-terraform.sh
```

The installer will show a series of prompts for user-specific information and use reasonable defaults for everything else. In non-interactive contexts, prompts can be bypassed by providing appropriately-named environment variables. Refer to the [user documentation](docs/user) for more information.

### Connect to the cluster

Details for connecting to your new cluster are printed by the `openshift-install` binary upon completion, and are also available in the `.openshift_install.log` file.

Example output:

```sh
INFO Waiting 10m0s for the openshift-console route to be created...
INFO Install complete!
INFO Run 'export KUBECONFIG=/path/to/auth/kubeconfig' to manage the cluster with 'oc', the OpenShift CLI.
INFO The cluster is ready when 'oc login -u kubeadmin -p 5char-5char-5char-5char' succeeds (wait a few minutes).
INFO Access the OpenShift web-console here: https://console-openshift-console.apps.${OPENSHIFT_INSTALL_CLUSTER_NAME}.${OPENSHIFT_INSTALL_BASE_DOMAIN}:6443
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
