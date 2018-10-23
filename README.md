# Openshift Installer

## Supported Platforms

* AWS
* [Libvirt with KVM](docs/dev/libvirt-howto.md)
* OpenStack (Experimental)

## Quick Start

First, install all [build dependencies](docs/dev/dependencies.md).

After cloning this repository, the installer binary will need to be built by running the following:

```sh
hack/build.sh
```

This will create `bin/openshift-install`. This binary can then be invoked to create an OpenShift cluster, like so:

```sh
bin/openshift-install cluster
```

The installer requires the terraform binary either alongside openshift-install or in `$PATH`.
If you don't have [terraform](https://www.terraform.io/), run the following to create `bin/terraform`:

```sh
hack/get-terraform.sh
```

The installer will show a series of prompts for user-specific information (e.g. admin password) and use reasonable defaults for everything else. In non-interactive contexts, prompts can be bypassed by providing appropriately-named environment variables. Refer to the [user documentation](docs/user) for more information.

### Connect to the cluster

#### Console

Shortly after the `cluster` command completes, the OpenShift console will come up at `https://${OPENSHIFT_INSTALL_CLUSTER_NAME}-api.${OPENSHIFT_INSTALL_BASE_DOMAIN}:6443/console/`.
You may need to ignore a certificate warning if you did not configure a certificate authority known to your browser.
Log in using the admin credentials you configured when creating the cluster.

#### Kubeconfig

You can also use the admin kubeconfig which `openshift-install cluster` placed under `--dir` (which defaults to `.`) in `auth/kubeconfig`.
If you launched the cluster with `openshift-install --dir "${DIR}" cluster`, you can use:

```sh
export KUBECONFIG="${DIR}/auth/kubeconfig"
```

### Cleanup

Destroy the cluster and release associated resources with:

```sh
openshift-install destroy-cluster
```
