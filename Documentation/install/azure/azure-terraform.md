# Install Tectonic on Azure with Terraform

Following this guide will deploy a Tectonic cluster within your Azure account.

Generally, the Azure platform templates adhere to the standards defined by the project [conventions][conventions] and [generic platform requirements][generic]. This document aims to document the implementation details specific to the Azure platform.

## Prerequsities

### Go

Ensure [Go][install-go] is installed.

### Terraform

Tectonic Installer includes and requires a specific version of Terraform. This is included in the Tectonic Installer tarball. See the [Tectonic Installer release notes][release-notes] for information about which Terraform versions are compatible.

### DNS

A few means of providing DNS for your Tectonic installation are supported:

#### Azure-provided DNS

This is Azure's default DNS implementation. For more information, see the [Azure DNS overview][azure-dns].

To use Azure-provided DNS, `tectonic_base_domain` must be set to `""`(empty string).

#### DNS delegation and custom zones via Azure DNS

To configure a custom domain and the associated records in an Azure DNS zone (e.g., `${cluster_name}.foo.bar`):

* The custom domain must be specified using `tectonic_base_domain`
* The domain must be publically discoverable. The Tectonic installer uses the created record to access the cluster and complete configuration. See the Microsoft Azure documentation for instructions on how to [delegate a domain to Azure DNS][domain-delegation].
* An Azure DNS zone which matches `tectonic_base_domain` must be created prior to running the installer. The full resource ID of the DNS zone must then be referenced in `tectonic_azure_external_dns_zone_id`

### Tectonic Account

Register for a [Tectonic Account][register], which is free for up to 10 nodes. You must provide the cluster license and pull secret during installation.

### Azure CLI

The [Azure CLI][azure-cli] is required to generate Azure credentials.

### ssh-agent

Ensure `ssh-agent` is running:
```
$ eval $(ssh-agent)
```

Add the SSH key that will be used for the Tectonic installation to `ssh-agent`:
```
$ ssh-add <path-to-ssh-key>
```

Verify that the SSH key identity is available to the ssh-agent:
```
$ ssh-add -L
```

Reference the absolute path of the **_public_** component of the SSH key in `tectonic_azure_ssh_key`.

Without this, terraform is not able to SSH copy the assets and start bootkube.
Also make sure that the SSH known_hosts file doesn't have old records of the API DNS name (fingerprints will not match).

## Getting Started

### Download and extract Tectonic Installer

Open a new terminal, and run the following commands to download and extract Tectonic Installer.

```bash
$ curl -O https://releases.tectonic.com/tectonic-1.7.1-tectonic.1.tar.gz # download
$ tar xzvf tectonic-1.7.1-tectonic.1.tar.gz # extract the tarball
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

* Etcd cluster nodes are managed by the terraform module `modules/azure/etcd`
* Node VMs are created as an Availability Set (stand-alone instances, deployed across multiple fault domains)
* A load-balancer fronts the etcd nodes to provide a simple discovery mechanism, via a VIP + DNS record.
* Currently, the LB is configured with a public IP address. Future work is planned to convert this to an internal LB.

### Master nodes

* Master node VMs are managed by the templates in `modules/azure/master-as`
* Node VMs are created as an Availability Set (stand-alone instances, deployed across multiple fault domains)
* Master nodes are fronted by one load balancer for the API and one for the Ingress controller.
* The API LB is configured with SourceIP session stickiness, to ensure that TCP (including SSH) sessions from the same client land reliably on the same master node. This allows for provisioning the assets and starting bootkube reliably via SSH.

### Worker nodes

* Worker node VMs are managed by the templates in `modules/azure/worker-as`
* Node VMs are created as an Availability Set (stand-alone instances, deployed across multiple fault domains)
* Worker nodes are not fronted by an LB and don't have public IP addresses. They can be accessed through SSH from any of the master nodes.

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
[azure-dns]: https://docs.microsoft.com/en-us/azure/dns/dns-overview
[vars]: ../../variables/config.md
[azure-vars]: ../../variables/azure.md
[release-notes]: https://coreos.com/tectonic/releases/
[install-go]: https://golang.org/doc/install
[domain-delegation]: https://docs.microsoft.com/en-us/azure/dns/dns-delegate-domain-azure-dns
[azure-cli]: https://docs.microsoft.com/en-us/cli/azure/install-azure-cli
