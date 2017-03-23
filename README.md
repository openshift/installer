# Tectonic Installer

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
- [Baremetal](https://coreos.com/tectonic/docs/latest/install/bare-metal/) [[**stable**][platform-lifecycle]]

### Hacking

#### Choose your platform:

See the platform specific documentation. Then see the *common usage* section below.

- [AWS via Terraform](Documentation/platforms/aws/README.md) [[**alpha**][platform-lifecycle]]
- [Azure via Terraform](Documentation/platforms/azure/README.md) [[**alpha**][platform-lifecycle]]
- [OpenStack via Terraform](Documentation/platforms/openstack/README.md) [[**alpha**][platform-lifecycle]]
- [VMware](Documentation/platforms/vmware/README.md) [[**alpha**][platform-lifecycle]]


#### Common Usage

At a high level, using the installer follows the workflow below. See each platform guide for specifics.

1. Download Terraform

This repo uses a custom build of Terraform which is pinned to a specific version and included required plugins. The easiest way to get the binary is:

```
$ make terraform-download
```

Follow the directions and update your PATH.

2. Initiate Working Directory

This will create a new directory `build/<cluster-name>` which holds all module references, Terraform state files, and custom variable files.

```
PLATFORM=aws CLUSTER=my-cluster make localconfig
```

3. Customize

Set variables in the `terraform.tfvars` file as needed, or you will be prompted. Available variables can be found in the `config.tf` and `variables.tf` files present in the `platforms/<PLATFORM>` directory.

4. Terraform Lifecycle

Plan, apply, and destroy are provided as Make targets to make working with the build directory and custom binary easier.

```
PLATFORM=aws CLUSTER=my-cluster make plan
```

```
PLATFORM=aws CLUSTER=my-cluster make apply
```

```
PLATFORM=aws CLUSTER=my-cluster make destroy
```

[platform-lifecycle]: Documentation/platform-lifecycle.md
