# Install Tectonic on Azure with Terraform

Following this guide will deploy a Tectonic cluster within your Azure account.

Generally, the Azure platform templates adhere to the standards defined by the project [conventions][conventions] and [generic platform requirements][generic]. This document aims to document the implementation details specific to the Azure platform.

<p style="background:#d9edf7; padding: 10px;" class="text-info"><strong>Alpha:</strong> These modules and instructions are currently considered alpha. See the <a href="../../platform-lifecycle.md">platform life cycle</a> for more details.</p>

## Prerequsities

* **Terraform**: Tectonic Installer includes and requires a specific version of Terraform. This is included in the Tectonic Installer tarball. See the [Tectonic Installer release notes][release-notes] for information about which Terraform versions are compatible.
* **DNS**: Set up your DNS zone in a resource group called `tectonic-dns-group` or specify a different resource group using the `tectonic_azure_dns_resource_group` variable below. We use a separate resource group assuming that you have a zone that you already want to use. Follow the [docs to set one up][azure-dns].
* **Tectonic Account**: Register for a [Tectonic Account][register], which is free for up to 10 nodes. You must provide the cluster license and pull secret during installation.
* **Azure CLI**: The Azure Command line interface is required to generate Azure credentials.

## Getting Started

### Download and extract Tectonic Installer

Open a new terminal, and run the following commands to download and extract Tectonic Installer.

```bash
$ curl -O https://releases.tectonic.com/tectonic-1.6.7-tectonic.1.tar.gz # download
$ tar xzvf tectonic-1.6.7-tectonic.1.tar.gz # extract the tarball
$ cd tectonic
```

### Initialize and configure Terraform

Start by setting the `INSTALLER_PATH` to the location of your platform's Tectonic installer. The platform should be `darwin` or `linux`.

```bash
$ export INSTALLER_PATH=$(pwd)/tectonic-installer/darwin/installer # Edit the platform name.
$ export PATH=$PATH:$(pwd)/tectonic-installer/darwin # Put the `terraform` binary in our PATH
```

Make a copy of the Terraform configuration file for our system. Do not share this configuration file as it is specific to your machine.

```bash
$ sed "s|<PATH_TO_INSTALLER>|$INSTALLER_PATH|g" terraformrc.example > .terraformrc
$ export TERRAFORM_CONFIG=$(pwd)/.terraformrc
```

Next, get the modules that Terraform will use to create the cluster resources:

```
$ terraform get platforms/azure
Get: file:///Users/tectonic-installer/modules/azure/vnet
Get: file:///Users/tectonic-installer/modules/azure/etcd
Get: file:///Users/tectonic-installer/modules/azure/master
Get: file:///Users/tectonic-installer/modules/azure/worker
Get: file:///Users/tectonic-installer/modules/azure/dns
Get: file:///Users/tectonic-installer/modules/bootkube
Get: file:///Users/tectonic-installer/modules/tectonic
```

Generate credentials using the Azure CLI. If you're not logged in, execute `az login` first. See the [docs][login] for more info.

```
$ az ad sp create-for-rbac -n "http://tectonic" --role contributor
Retrying role assignment creation: 1/24
Retrying role assignment creation: 2/24
{
 "appId": "generated-app-id",
 "displayName": "azure-cli-2017-01-01",
 "name": "http://tectonic-coreos",
 "password": "generated-pass",
 "tenant": "generated-tenant"
}
```

Export variables that correspond to the data that was just generated. The subscription is your Azure Subscription ID.

```
$ export ARM_SUBSCRIPTION_ID=abc-123-456
$ export ARM_CLIENT_ID=generated-app-id
$ export ARM_CLIENT_SECRET=generated-pass
$ export ARM_TENANT_ID=generated-tenant
```

Now we're ready to specify our cluster configuration.

## Customize the deployment

Customizations to the base installation live in `examples/terraform.tfvars.azure`. Export a variable that will be your cluster identifier:

```
$ export CLUSTER=my-cluster
```

Create a build directory to hold your customizations and copy the example file into it:

```
$ mkdir -p build/${CLUSTER}
$ cp examples/terraform.tfvars.azure build/${CLUSTER}/terraform.tfvars
```

Edit the parameters with your Azure details, domain name, license, etc. [View all of the Azure specific options][azure-vars] and [the common Tectonic variables][vars].

## Deploy the cluster

Test out the plan before deploying everything:

```
$ terraform plan -var-file=build/${CLUSTER}/terraform.tfvars platforms/azure
```

Next, deploy the cluster:

```
$ terraform apply -var-file=build/${CLUSTER}/terraform.tfvars platforms/azure
```

This should run for a little bit, and when complete, your Tectonic cluster should be ready.

If you encounter any issues, check the known issues and workarounds below.

## Access the cluster

The Tectonic Console should be up and running after the containers have downloaded. You can access it at the DNS name configured in your variables file.

Inside of the `/generated` folder you should find any credentials, including the CA if generated, and a kubeconfig. You can use this to control the cluster with `kubectl`:

```
$ KUBECONFIG=generated/auth/kubeconfig
$ kubectl cluster-info
```

## Scale the cluster


To scale worker nodes, adjust `tectonic_worker_count` in `terraform.tfvars`.

Use the `plan` command to check your syntax: 

```
$ terraform plan \
  -var-file=build/${CLUSTER}/terraform.tfvars \
  -target module.workers \
  platforms/azure
```

Once you are ready to make the changes live, use `apply`:

```
$ terraform apply \
  -var-file=build/${CLUSTER}/terraform.tfvars \
  -target module.workers \
  platforms/azure
```

The new nodes should automatically show up in the Tectonic Console shortly after they boot.

## Delete the cluster

Deleting your cluster will remove only the infrastructure elements created by Terraform. If you selected an existing resource group for DNS, this is not touched. To delete, run:

```
$ terraform destroy -var-file=build/${CLUSTER}/terraform.tfvars platforms/azure
```

## Under the hood

### Top-level templates

* The top-level templates that invoke the underlying component modules reside `./platforms/azure`
* Point terraform to this location to start applying: `terraform apply ./platforms/azure`

### Etcd nodes

* Discovery is currently not implemented so DO NOT scale the etcd cluster to more than 1 node, for now.
* Etcd cluster nodes are managed by the terraform module `modules/azure/etcd`
* Node VMs are created as stand-alone instances (as opposed to VM scale sets). This is mostly historical and could change.
* A load-balancer fronts the etcd nodes to provide a simple discovery mechanism, via a VIP + DNS record.
* Currently, the LB is configured with a public IP address. This is not optimal and it should be converted to an internal LB.

### Master nodes

* Master node VMs are managed by the templates in `modules/azure/master`
* An Azure VM Scaling Set resource is used to spin-up multiple identical VM configured as master nodes.
* Master VMs all share the same identical Ignition config
* Master nodes are fronted by one load-balancer for the API one for the Ingress controller.
* The API LB is configured with SourceIP session stickiness, to ensure that TCP (including SSH) sessions from the same client land reliably on the same master node. This allows for provisioning the assets and starting bootkube reliably via SSH.
* a `null_resource` terraform provisioner in the tectonic.tf top-level template will copy the assets and run bootkube automatically on one of the masters.
* make sure the SSH key specified in the tfvars file is also added to the SSH agent on the machine running terraform. Without this, terraform is not able to SSH copy the assets and start bootkube. Also make sure that the SSH known_hosts file doesn't have old records of the API DNS name (fingerprints will not match).

### Worker nodes

* Worker node VMs are managed by the templates in `modules/azure/worker`
* An Azure VM Scaling Set resource is used to spin-up multiple identical VM configured as worker nodes.
* Worker VMs all share the same identical Ignition config
* Worker nodes are not fronted by any LB and don't have public IP addresses. They can be accessed through SSH from any of the master nodes.

## Known issues and workarounds

See the [installer troubleshooting][troubleshooting] document for known problem points and workarounds.


[conventions]: ../../conventions.md
[generic]: ../../generic-platform.md
[register]: https://account.coreos.com/signup/summary/tectonic-2016-12
[account]: https://account.coreos.com
[bcrypt]: https://github.com/coreos/bcrypt-tool/releases/tag/v1.0.0
[plan-docs]: https://www.terraform.io/docs/commands/plan.html
[copy-docs]: https://www.terraform.io/docs/commands/apply.html
[troubleshooting]: ../../troubleshooting/installer-terraform.md
[login]: https://docs.microsoft.com/en-us/cli/azure/get-started-with-azure-cli
[azure-dns]: https://docs.microsoft.com/en-us/azure/dns/dns-getstarted-portal
[vars]: ../../variables/config.md
[azure-vars]: ../../variables/azure.md
[release-notes]: https://coreos.com/tectonic/releases/
