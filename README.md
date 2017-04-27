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

Note: This repo does not yet have all Tectonic Installer features imported. This will happen over the coming weeks as we are able to make some of the surrounding infrastructure public as well. This notice will be removed once the AWS and Baremetal graphical installer code has been fully integrated.

Checkout the [ROADMAP](ROADMAP.md) for details on where the project is headed.

## Getting Started

To use the installer you can either use an official release (starting March 29, 2017), or hack on the scripts in this repo.

### Official releases

See the official Tectonic documentation:

- [AWS Cloud Formation](https://coreos.com/tectonic/docs/latest/install/aws/) [[**stable**][platform-lifecycle]]
- [Bare-Metal](https://coreos.com/tectonic/docs/latest/install/bare-metal/) [[**stable**][platform-lifecycle]]

### Hacking

#### Common Usage

At a high level, using the installer follows the workflow below. See each platform guide for specifics.

**Install Terraform**

This project is built on Terraform and requires version 0.9.4. Download and install an [official Terraform binary](https://releases.hashicorp.com/terraform/0.8.8/) for your OS or use your favorite package manager.

**Choose your platform**

The example below will use `PLATFORM=azure` but you can set the value to something different. Also, as you configure the cluster refer to the linked documentation to find the configuration parameters.

- `PLATFORM=aws` [AWS via Terraform](Documentation/install/aws/aws-terraform.md) [[**alpha**][platform-lifecycle]]
- `PLATFORM=azure` [Azure via Terraform](Documentation/install/azure/azure-terraform.md) [[**alpha**][platform-lifecycle]]
- `PLATFORM=metal` [Bare-Metal via Terraform](Documentation/install/metal/metal-terraform.md) [[**alpha**][platform-lifecycle]]
- `PLATFORM=openstack` [OpenStack via Terraform](Documentation/install/openstack/openstack-terraform.md) [[**alpha**][platform-lifecycle]]

**Initiate the Cluster Configuration**

This will create a new directory `build/<cluster-name>` which holds all module references, Terraform state files, and custom variable files.

```
PLATFORM=azure CLUSTER=my-cluster make localconfig
```

**Configure Cluster**

Set variables in the `terraform.tfvars` file as needed, or you will be prompted. Available variables can be found in the `config.tf` and `variables.tf` files present in the `platforms/<PLATFORM>` directory.

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

[platform-lifecycle]: Documentation/platform-lifecycle.md
