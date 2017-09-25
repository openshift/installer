# Tectonic Installer

Tectonic is built on pure-upstream Kubernetes but has an opinion on the best way to install and run a Kubernetes cluster. This project helps you install a Kubernetes cluster the "Tectonic Way". It provides good defaults, enables install automation, and is customizable to meet your infrastructure needs.

Goals of the project:

- Install [Self-Hosted Kubernetes Clusters](https://github.com/kubernetes/community/blob/master/contributors/design-proposals/self-hosted-kubernetes.md)
- Secure by default (uses TLS, RBAC by default, OIDC AuthN, etcd)
- Automatable install process for scripts and CI/CD
- Deploy on any infrastructure: Amazon AWS, Microsoft Azure, OpenStack, Google Cloud, bare metal
- Run on any OS: Container Linux (the default), [RHEL](Documentation/install/rhel/installing-workers.md#installing-tectonic-workers-on-red-hat-enterprise-linux), Ubuntu, and others
- Customizable and modular: Change DNS providers, security settings, authentication providers
- Highly Available by default: Deploy all Kubernetes components HA, use etcd Operator

Check the [ROADMAP](ROADMAP.md) for details on where the project is headed.

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

- [OpenStack via Terraform](Documentation/install/openstack/openstack-terraform.md) [[**alpha**][platform-lifecycle]]
- [VMware via Terraform](Documentation/install/vmware/vmware-terraform.md) [[**alpha**][platform-lifecycle]]

**Go and Source**

[Install Go](https://golang.org/doc/install) if not already installed.

Then get the Tectonic Installer source code:

```
go get github.com/coreos/tectonic-installer
cd $(go env GOPATH)/src/github.com/coreos/tectonic-installer
```

**Terraform**

The Tectonic Installer releases include a build of [Terraform](https://terraform.io). See the [Tectonic Installer release notes][release-notes] for information about which Terraform versions are compatible.

The [latest Terraform binary](https://www.terraform.io/downloads.html) may not always work as Tectonic Installer, which sometimes relies on bug fixes or features not yet available in the official Terraform release.

**Yarn (optional)**

The [Yarn](https://yarnpkg.com) JavaScript package manager is required for building the frontend code. On OS X, install using Homebrew: `brew install yarn`.

#### Common Usage

**Choose your platform**

First, set the `PLATFORM=` environment variable. This example will use `PLATFORM=azure`.

- `PLATFORM=openstack` [OpenStack via Terraform](Documentation/install/openstack/openstack-terraform.md) [[**alpha**][platform-lifecycle]]
- `PLATFORM=vmware` [VMware via Terraform](Documentation/install/vmware/vmware-terraform.md) [[**alpha**][platform-lifecycle]]

**Initiate the Cluster Configuration**

Use `make` to create a new directory `build/<cluster-name>` to hold all module references, Terraform state files, and custom variable files.

```
PLATFORM=azure CLUSTER=my-cluster make localconfig
```

**Configure Cluster**

Set variables in the `build/<cluster-name>/terraform.tfvars` file as needed. Available variables are found in the `platforms/<PLATFORM>/config.tf` and `platforms/<PLATFORM>/variables.tf` files.

Examples for each platform can be found in [the examples directory](examples/).

**Terraform Lifecycle**

`plan`, `apply`, and `destroy` are provided as `make` targets to ease the build directory and custom binary complexity.

```
PLATFORM=azure CLUSTER=my-cluster make plan
```

```
PLATFORM=azure CLUSTER=my-cluster make apply
```

```
PLATFORM=azure CLUSTER=my-cluster make destroy
```

#### Tests

We have different set of tests:

##### Basic tests

Our basic set of tests includes:
- Code linting
- UI tests
- Backend unit tests

They are run on **every** PR.

##### Smoke tests

In addition to our basic set of tests we have smoke tests. These test the
Tectonic installer on our supported platforms.
- AWS
- Azure
- Bare metal

They can be run on a PR by applying the *run-smoke-tests* GitHub label.

Further details can be found in our [Jenkinsfile](./Jenkinsfile) which serves as
the single source of truth.

To run a smoke test locally you need to set the following environment variables:

```
CLUSTER
AWS_ACCESS_KEY_ID
AWS_SECRET_ACCESS_KEY
TF_VAR_tectonic_aws_ssh_key
TF_VAR_tectonic_aws_region
TF_VAR_tectonic_license_path
TF_VAR_tectonic_pull_secret_path
TF_VAR_base_domain
```

Make sure both the *Tectonic pull secret* as well as the *Tectonic license* is
saved somewhere in the repository folder. Only the repository folder will be
mounted into the Docker container where the tests will be executed in. The test
framework will not be able to read any files outside the repository folder
during test execution.

Once the environment variables are set, run `make tests/smoke
TEST=spec/aws_spec.rb`.


[platform-lifecycle]: Documentation/platform-lifecycle.md
[release-notes]: https://coreos.com/tectonic/releases/
