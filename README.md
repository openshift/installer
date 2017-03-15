# Tectonic Installer Kit

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

### Step 1: Sign-up for the Tectonic Free Tier

Sign-up for the [Tectonic Free Tier](https://coreos.com/tectonic).

*Note:* We will make this project flexible enough in the coming weeks to just install Kubernetes without the additional Tectonic Components. Please help make this happen or follow this issue.

### Step 2: Download the Tectonic installer.

```
wget https://releases.tectonic.com/tectonic-X.Y.Z-tectonic.N.tar.gz
tar xzvf tectonic-X.Y.Z-tectonic.N.tar.gz
```

### Step 2: Choose a Platform

- [AWS Cloud Formation](https://coreos.com/tectonic/docs/latest/install/aws/) [[**stable**][platform-lifecycle]]
- [Baremetal](https://coreos.com/tectonic/docs/latest/install/bare-metal/) [[**stable**][platform-lifecycle]]
- [AWS via Terraform](aws/README.md) [[**alpha**][platform-lifecycle]]
- [Azure via Terraform](azure/README.md) [[**alpha**][platform-lifecycle]]
- [OpenStack via Terraform](openstack/README.md) [[**alpha**][platform-lifecycle]]
- [VMware](vmware/README.md) [[**alpha**][platform-lifecycle]]

# Old README stuff

## Getting Started

At this time the Platform SDK relies on the Tectonic Installer to generate all of the Kubernetes assests, certificates, etc. If you don't have a Tectonic installer already [sign-up for one for the free tier](https://coreos.com/tectonic) first, then:

1. Use the Tectonic installer to configure an AWS cluster.
2. Go through the process to create an AWS cluster, do not apply the configuration, but download the assets manually. This is an advanced option on the last screen
3. Unzip the assets in this directory:

```
$ unzip ~/Downloads/<name>-assets.zip
```

## Azure

1. Setup your DNS zone in a resource group called `tectonic-dns-group` or specify a different resource group. We use a separate resource group assuming that you have a zone that you already want to use.
1. Create a folder with the cluster's name under `./build` (e.g. `./build/<cluster-name>`)
1. Copy the `assets-<cluster-name>.zip` to `./boot/<cluster-name>`

```
make PLATFORM=azure CLUSTER=eugene
```

*Common Prerequsities*

1. Configure AWS credentials via environment variables.
[See docs](http://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html#cli-environment)
1. Configure a region by setting `AWS_REGION` environment variable
1. Run through the official Tectonic intaller steps without clicking `Submit` on the last step. 
Instead click on `Manual boot` below to download the assets zip file.
1. Create a folder with the cluster's name under `./build` (e.g. `./build/<cluster-name>`)
1. Copy the `assets-<cluster-name>.zip` to `./boot/<cluster-name>`

### Using Autoscaling groups

1. Ensure all *prerequsities* are met.
1. From the root of the repo, run `make PLATFORM=aws-asg CLUSTER=<cluster-name>`

To clean up run `make destroy PLATFORM=aws-asg CLUSTER=<cluster-name>`


## OpenStack

Prerequsities:

1. The latest Container Linux Alpha (1339.0.0 or later) [uploaded into Glance](https://coreos.com/os/docs/latest/booting-on-openstack.html) and get the image ID
1. Since openstack nova doesn't provide any DNS registration service, AWS Route53 is being used.
Ensure you have a configured `aws` CLI installation.
1. Ensure you have OpenStack credentials set up, i.e. the environment variables `OS_TENANT_NAME`, `OS_USERNAME`, `OS_PASSWORD`, `OS_AUTH_URL`, `OS_REGION_NAME` are set.
1. Create a folder with the cluster's name under `./build` (e.g. `./build/<cluster-name>`)
1. Copy the `assets-<cluster-name>.zip` to `./boot/<cluster-name>`

### Nova network

1. Ensure all *prerequsities* are met.
1. From the root of the repo, run `make PLATFORM=openstack-novanet CLUSTER=<cluster-name>`

To clean up run `make destroy PLATFORM=openstack-novanet CLUSTER=<cluster-name>`

The tectonic cluster will be reachable under `https://<name>.<base_domain>:32000`.

### Neutron network

1. Ensure all *prerequsities* are met.
1. From the root of the repo, run `make PLATFORM=openstack-neutron CLUSTER=<cluster-name>`

To clean up run `make destroy PLATFORM=openstack-neutron CLUSTER=<cluster-name>`

The tectonic cluster will be reachable under `https://<name>.<base_domain>:32000`.

## AWS

*Common Prerequsities*

1. Configure AWS credentials via environment variables. 
[See docs](http://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html#cli-environment)
1. Configure a region by setting `AWS_REGION` environment variable
1. Run through the official Tectonic intaller steps without clicking `Submit` on the last step. 
Instead click on `Manual boot` below to download the assets zip file.
1. Create a folder with the cluster's name under `./build` (e.g. `./build/<cluster-name>`)
1. Copy the `assets-<cluster-name>.zip` to `./boot/<cluster-name>`

### Using Autoscaling groups

1. Ensure all *prerequsities* are met.
1. From the root of the repo, run `make PLATFORM=aws-asg CLUSTER=<cluster-name>`

To clean up run `make destroy PLATFORM=aws-asg CLUSTER=<cluster-name>`

### Without Autoscaling groups

1. Ensure all *prerequsities* are met.
1. From the root of the repo, run `make PLATFORM=aws-noasg CLUSTER=<cluster-name>`

To clean up run `make destroy PLATFORM=aws-noasg CLUSTER=<cluster-name>`

[platform-lifecycle]: Documentation/platform-lifecycle.md
