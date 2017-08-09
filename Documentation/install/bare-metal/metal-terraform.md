# Install Tectonic on Bare-Metal with Terraform

Following this guide will deploy a Tectonic cluster on virtual or physical hardware.

## Prerequisites

* Terraform: Tectonic Installer includes and requires a specific version of Terraform. This is included in the Tectonic Installer tarball. See the [Tectonic Installer release notes][release-notes] for information about which Terraform versions are compatible.
* [Matchbox v0.6+](https://github.com/coreos/matchbox/releases) installation with the gRPC API enabled. See [installation](https://coreos.com/matchbox/docs/latest/deployment.html).
* Matchbox TLS client credentials
* PXE network boot environment with DHCP, TFTP, and DNS services. See [network-setup](https://coreos.com/matchbox/docs/latest/network-setup.html).
* DNS records for the Kubernetes controller(s) and Tectonic Ingress worker(s). See [DNS](https://coreos.com/tectonic/docs/latest/install/bare-metal#networking).
* Machines with BIOS options set to boot from disk normally, but PXE prior to installation.
* Machines with known MAC addresses and stable domain names.
* Tectonic Account - Register for a [Tectonic Account][register], which is free for up to 10 nodes. You will need to provide the cluster license and pull secret below.
* `ipmitool` or `virt-install` will be used to actually boot the machines.
* A SSH keypair whose private key is present in your system's ssh-agent.

## Getting Started

### Download and extract Tectonic Installer

Open a new terminal, and run the following commands to download and extract Tectonic Installer.

```bash
$ curl -O https://releases.tectonic.com/tectonic-1.7.1-tectonic.1.tar.gz
$ tar xzvf tectonic-1.7.1-tectonic.1.tar.gz
$ cd tectonic
```

### Initialize and configure Terraform

Start by setting the `INSTALLER_PATH` to the location of your platform's Tectonic installer. The platform should be `linux` or `darwin`.

```bash
$ export INSTALLER_PATH=$(pwd)/tectonic-installer/linux/installer
$ export PATH=$PATH:$(pwd)/tectonic-installer/linux
```

Make a copy of the Terraform configuration file for our system. Do not share this configuration file as it is specific to your machine.

```bash
$ sed "s|<PATH_TO_INSTALLER>|$INSTALLER_PATH|g" terraformrc.example > .terraformrc
$ export TERRAFORM_CONFIG=$(pwd)/.terraformrc
```

Next, get the modules that Terraform will use to create the cluster resources:

```bash
$ terraform get ./platforms/metal
```

Now we're ready to specify our cluster configuration.

## Customize the deployment

Create a build directory to hold your customizations and copy the example file into it:

```
$ export CLUSTER=my-cluster
$ mkdir -p build/${CLUSTER}
$ cp examples/terraform.tfvars.metal build/${CLUSTER}/terraform.tfvars
```

Customizations should be made to `build/${CLUSTER}/terraform.tfvars`. Edit the following variables to correspond to your matchbox installation:

* `tectonic_matchbox_http_url`
* `tectonic_matchbox_rpc_endpoint`
* `tectonic_matchbox_client_cert`
* `tectonic_matchbox_client_key`
* `tectonic_matchboc_ca`

Edit additional variables to specify DNS records, list machines, and set a SSH key and Tectonic Console email and password.

Several variables are currently required, but their values are not used.

* `tectonic_are_domain`
* `tectonic_master_count`
* `tectonic_worker_count`
* `tectonic_etcd_count`

## Deploy the cluster

Test out the plan before deploying everything:

```
$ terraform plan -var-file=build/${CLUSTER}/terraform.tfvars platforms/metal
```

Next, deploy the cluster:

```
$ terraform apply -var-file=build/${CLUSTER}/terraform.tfvars platforms/metal
```

This will write machine profiles and matcher groups to the matchbox service.

## Power On

Power on the machines with `ipmitool` or `virt-install`. Machines will PXE boot, install Container Linux to disk, and reboot.

```
ipmitool -H node1.example.com -U USER -P PASS power off
ipmitool -H node1.example.com -U USER -P PASS chassis bootdev pxe
ipmitool -H node1.example.com -U USER -P PASS power on
```

Terraform will wait for the disk installation and reboot to complete and then be able to copy credentials to the nodes to bootstrap the cluster. You may see `null_resource.kubeconfig.X: Still creating...` during this time.

Run `terraform apply` until all tasks complete. Your Tectonic cluster should be ready. If you encounter any issues, check the known issues and workarounds below.

## Access the cluster

The Tectonic Console should be up and running after the containers have downloaded. You can access it at the DNS name configured in your variables file.

Inside of the `/generated` folder you should find any credentials, including the CA if generated, and a kubeconfig. You can use this to control the cluster with `kubectl`:

```
$ export KUBECONFIG=generated/auth/kubeconfig
$ kubectl cluster-info
```

## Work with the cluster

For more information on working with installed clusters, see [Scaling Tectonic bare metal clusters][scale-metal], and [Uninstalling Tectonic][uninstall].

## Known issues and workarounds

See the [troubleshooting][troubleshooting] document for workarounds for bugs that are being tracked.


[conventions]: ../../conventions.md
[generic]: ../../generic-platform.md
[register]: https://account.coreos.com/signup/summary/tectonic-2016-12
[account]: https://account.coreos.com
[vars]: ../../variables/config.md
[troubleshooting]: ../../troubleshooting/faq.md
[uninstall]: uninstall.md
[scale-metal]: ../../admin/bare-metal-scale.md
[release-notes]: https://coreos.com/tectonic/releases/
