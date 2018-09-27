# Openshift Installer

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
