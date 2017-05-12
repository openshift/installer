# Tectonic Installer
[![Build Status](https://jenkins-tectonic-installer-public.prod.coreos.systems/buildStatus/icon?job=coreos%20-%20tectonic-installer/tectonic-installer/master)](https://jenkins-tectonic-installer-public.prod.coreos.systems/job/coreos%20-%20tectonic-installer/job/tectonic-installer/job/master)

Tectonic is built on pure-upstream Kubernetes but has an opinion on the best way to install and run a Kubernetes cluster. This project helps you install a Kubernetes cluster the "Tectonic Way". It provides good defaults, enables install automation, and is customizable to meet your infrastructure needs.

Goals of the project:

- Installation of [Self-Hosted Kubernetes Cluster](https://github.com/kubernetes/community/blob/master/contributors/design-proposals/self-hosted-kubernetes.md)
- Secure by default (use TLS, RBAC by default, OIDC AuthN, etcd)
- Automatable install process for scripts and CI/CD
- Deploy Tectonic on any infrastructure (Amazon, Azure, OpenStack, GCP, etc)
- Runs Tectonic on any OS (Container Linux, RHEL, CentOS, etc)
- Customizable and modular (change DNS providers, security settings, etc)
- HA by default (deploy all Kubernetes components HA, use etcd Operator)

Checkout the [ROADMAP](ROADMAP.md) for details on where the project is headed.

## Getting Started

**To use a tested release** on an supported platform, follow the links below.

**To hack or modify** the templates or add a new platform, use the scripts in this repo to boot and tear down clusters.

### Official releases

See the official Tectonic documentation:

- [AWS using a GUI](https://coreos.com/tectonic/docs/latest/install/aws/) [[**stable**][platform-lifecycle]]
- [AWS using Terraform CLI](https://coreos.com/tectonic/docs/latest/install/aws/aws-terraform.html) [[**stable**][platform-lifecycle]]
- [Bare metal using a GUI](https://coreos.com/tectonic/docs/latest/install/bare-metal/) [[**stable**][platform-lifecycle]]
- [Bare metal using Terraform CLI](https://coreos.com/tectonic/docs/latest/install/bare-metal/metal-terraform.html) [[**stable**][platform-lifecycle]]

### Hacking

In order to successfully build this project, you must first of all place it according to the Go workspace convention, i.e. at `$GOPATH/src/github.com/coreos/tectonic-installer`. If you don't set `$GOPATH`, it should by default be at `$HOME/go`.

#### Requirements

To build Tectonic Installer, you will need to install the following requirements:

##### Terraform

This project is built on [Terraform](http://terraform.io) and requires version 0.9.4. Download and install an [official Terraform binary](https://releases.hashicorp.com/terraform/0.9.4/) for your OS or use your favorite package manager.

##### Yarn

You need the [Yarn](https://yarnpkg.com) JavaScript package manager. If you're on OS X, you can install it via Homebrew: `brew install yarn`.

#### Common Usage

At a high level, using the installer follows the workflow below. See each platform guide for specifics.

**Choose your platform**

The example below will use `PLATFORM=azure` but you can set the value to something different. Also, as you configure the cluster refer to the linked documentation to find the configuration parameters.

- `PLATFORM=azure` [Azure via Terraform](Documentation/install/azure/azure-terraform.md) [[**alpha**][platform-lifecycle]]
- `PLATFORM=openstack` [OpenStack via Terraform](Documentation/install/openstack/openstack-terraform.md) [[**alpha**][platform-lifecycle]]

**Initiate the Cluster Configuration**

This will create a new directory `build/<cluster-name>` which holds all module references, Terraform state files, and custom variable files.

```
PLATFORM=azure CLUSTER=my-cluster make localconfig
```

**Configure Cluster**

Set variables in the `terraform.tfvars` file as needed, or you will be prompted. Available variables can be found in the `config.tf` and `variables.tf` files present in the `platforms/<PLATFORM>` directory.
Examples for each platform can be found in [the examples directory](examples/).

**Terraform Lifecycle**

Plan, apply, and destroy are provided as Make targets to make working with the build directory and custom binary easier.

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

Tests are run for all approved pull requests via Jenkins. See the [Jenkinsfile](./Jenkinsfile) for details.

Tests can be run locally by:


**AWS**

```
export PLATFORM="aws"
export AWS_REGION="us-east-1"
export TF_VAR_tectonic_cluster_name=my-smoke-test
export TF_VAR_tectonic_license_path=/path/to/license.txt
export TF_VAR_tectonic_pull_secret_path=/path/to/pull-secret.json

make localconfig
ln -sf ../../test/aws.tfvars build/${TF_VAR_tectonic_cluster_name}/terraform.tfvars
make plan
make apply
make destroy
```

[platform-lifecycle]: Documentation/platform-lifecycle.md
