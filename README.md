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

These instructions can be used for the official stable platforms listed above, and for the following alpha/beta platforms:

- [OpenStack via Terraform][openstack-tf] [[**alpha**][platform-lifecycle]]
- [VMware via Terraform][vmware-tf] [[**alpha**][platform-lifecycle]]


1. Build the project
    ```shell
    bazel build tarball
    ```

2. Unzip the tarball
    ```shell
    cd bazel-bin
    tar -zxvf tectonic.tar.gz
    cd tectonic
    ```

3. Add binaries to $PATH
    ```shell
    export PATH=$(pwd)/tectonic-installer/linux/:$PATH
    ```

4. Choose **one** of the platforms:
    ```shell
    export PLATFORM=aws
    export PLATFORM=azure
    export PLATFORM=gcp
    export PLATFORM=govcloud
    export PLATFORM=metal
    export PLATFORM=openstack-neutron
    export PLATFORM=vmware
    ```

5. Edit Tectonic configuration file including the $CLUSTER_NAME
    ```shell
    $EDITOR examples/tectonic.$PLATFORM.yaml`
    ```

6. Init Tectonic CLI
    ```shell
    tectonic init --config=example/tectonic.$PLATFORM.yaml
    ```

7. Install Tectonic cluster
    ```shell
    tectonic install --dir=$CLUSTER_NAME
    ```

8. Teardown Tectonic cluster
    ```shell
    tectonic destroy $CLUSTER_NAME
    ```


#### Tests

See [tests/README.md](tests/README.md).


[openstack-tf]: https://github.com/coreos/tectonic-docs/blob/master/Documentation/install/openstack/openstack-terraform.md
[platform-lifecycle]: https://coreos.com/tectonic/docs/latest/platform-lifecycle.html
[release-notes]: https://coreos.com/tectonic/releases/
[rhel-installation]: https://coreos.com/tectonic/docs/latest/install/rhel/installing-workers.html
[vmware-tf]: https://github.com/coreos/tectonic-docs/blob/master/Documentation/install/vmware/vmware-terraform.md
