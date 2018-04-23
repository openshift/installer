# Tectonic Installer

Tectonic is built on pure-upstream Kubernetes but has an opinion on the best way to install and run a Kubernetes cluster. This project helps you install a Kubernetes cluster the "Tectonic Way". It provides good defaults, enables install automation, and is customizable to meet your infrastructure needs.

Goals of the project:

- Install Kubernetes clusters
- Secure by default (uses TLS, RBAC by default, OIDC AuthN, etcd)
- Automatable install process for scripts and CI/CD
- Deploy on any infrastructure: Amazon AWS, Microsoft Azure, OpenStack, Google Cloud, bare metal
- Run on any OS: Container Linux (the default), [RHEL][rhel-installation], Ubuntu, and others
- Customizable and modular: Change DNS providers, security settings, authentication providers
- Highly Available by default: Deploy all Kubernetes components HA, use etcd Operator

*Note*: the project has recently undergone some rearchitecting to support our goal of providing automatic operations, most notably automatic updates, to Kubernetes clusters. The master branch of the project reflects this new design approach and currently provides support only for AWS. In order to deploy Tectonic to other platforms, e.g. Azure, bare metal, OpenStack, etc, please checkout the [track-1](https://github.com/coreos/tectonic-installer/tree/track-1) branch of this project, which maintains support for the previous architecture and more platforms.

## Getting Started

**To use a tested release** on a supported platform, follow the links below.

**To hack or modify** the templates or add a new platform, use the scripts in this repo to boot and tear down clusters.

### Official releases

See the official Tectonic documentation:

- [AWS using a GUI](https://coreos.com/tectonic/docs/latest/install/aws/) [[**stable**][platform-lifecycle]]
- [AWS using Terraform CLI](https://coreos.com/tectonic/docs/latest/install/aws/aws-terraform.html) [[**stable**][platform-lifecycle]]
- [Azure using Terraform](https://coreos.com/tectonic/docs/latest/install/azure/azure-terraform.html) [[**stable**][platform-lifecycle]]
- [Bare metal using a GUI](https://coreos.com/tectonic/docs/latest/install/bare-metal/) [[**stable**][platform-lifecycle]]
- [Bare metal using Terraform CLI](https://coreos.com/tectonic/docs/latest/install/bare-metal/metal-terraform.html) [[**stable**][platform-lifecycle]]

### Hacking

These instructions can be used for AWS:

1. Build the project
    ```shell
    bazel build tarball
    ```

    *Note*: the project can optionally be built without installing Bazel, provided Docker is installed:
    ```shell
    docker run --rm -v $PWD:$PWD -w $PWD quay.io/coreos/tectonic-builder:bazel-v0.3 bazel --output_base=.cache build tarball
    ```

2. Extract the tarball
    ```shell
    tar -zxf bazel-bin/tectonic-dev.tar.gz
    cd tectonic-dev
    ```

3. Add binaries to $PATH
    ```shell
    export PATH=$(pwd)/installer:$PATH
    ```

4. Edit Tectonic configuration file including the $CLUSTER_NAME
    ```shell
    $EDITOR examples/tectonic.aws.yaml
    ```

5. Init Tectonic CLI
    ```shell
    tectonic init --config=examples/tectonic.aws.yaml
    ```

6. Install Tectonic cluster
    ```shell
    tectonic install --dir=$CLUSTER_NAME
    ```

7. Teardown Tectonic cluster
    ```shell
    tectonic destroy --dir=$CLUSTER_NAME
    ```

#### Tests

See [tests/README.md](tests/README.md).


[platform-lifecycle]: https://coreos.com/tectonic/docs/latest/platform-lifecycle.html
[release-notes]: https://coreos.com/tectonic/releases/
[rhel-installation]: https://coreos.com/tectonic/docs/latest/install/rhel/installing-workers.html
