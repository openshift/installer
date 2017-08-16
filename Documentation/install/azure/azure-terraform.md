# Install Tectonic on Azure with Terraform

This guide deploys a Tectonic cluster on an Azure account.

The Azure platform templates generally adhere to the standards defined by the project [conventions][conventions] and [generic platform requirements][generic]. This document aims to clarify the implementation details specific to the Azure platform.

## Prerequisites

### Terraform

Tectonic Installer includes the required version of Terraform. See the [Tectonic Installer release notes][release-notes] for information about which Terraform versions are compatible.

### DNS

Two methods of providing DNS for the Tectonic installation are supported:

#### Azure-provided DNS

This is Azure's default DNS implementation. For more information, see the [Azure DNS overview][azure-dns].

To use Azure-provided DNS, `tectonic_base_domain` must be set to `""`(empty string).

#### DNS delegation and custom zones via Azure DNS

To configure a custom domain and the associated records in an Azure DNS zone (e.g., `${cluster_name}.foo.bar`):

* The custom domain must be specified using `tectonic_base_domain`
* The domain must be publicly discoverable. The Tectonic installer uses the created record to access the cluster and complete configuration. See the Microsoft Azure documentation for instructions on how to [delegate a domain to Azure DNS][domain-delegation].
* An Azure DNS zone matching the chosen `tectonic_base_domain` must be created prior to running the installer. The full resource ID of the DNS zone must then be referenced in `tectonic_azure_external_dns_zone_id`

### Tectonic Account

Register for a [Tectonic Account][register], free for up to 10 nodes. The cluster license and pull secret are required during installation.

### Azure CLI

The [Azure CLI][azure-cli] is required to generate Azure credentials.

### ssh-agent

Ensure `ssh-agent` is running:
```
$ eval $(ssh-agent)
```

Add the SSH key that will be used for the Tectonic installation to `ssh-agent`:
```
$ ssh-add <path-to-ssh-private-key>
```

Verify that the SSH key identity is available to the ssh-agent:
```
$ ssh-add -L
```

Reference the absolute path of the **_public_** component of the SSH key in `tectonic_azure_ssh_key`.

Without this, terraform is not able to SSH copy the assets and start bootkube.
Also ensure the SSH known_hosts file doesn't have old records for the API DNS name, because key fingerprints will not match.

## Getting Started

### Download and extract Tectonic Installer

Open a new terminal and run the following commands to download and extract Tectonic Installer:

```bash
$ curl -O https://releases.tectonic.com/tectonic-1.7.1-tectonic.1.tar.gz # download
$ tar xzvf tectonic-1.7.1-tectonic.1.tar.gz # extract the tarball
$ cd tectonic
```

### Initialize and configure Terraform

#### Set INSTALLER_PATH

Start by setting `INSTALLER_PATH` to the location of the installation host's Tectonic installer platform. The platform should be one of `darwin` or `linux`.

```bash
$ export INSTALLER_PATH=$(pwd)/tectonic-installer/darwin/installer # Edit the platform name.
$ export PATH=$PATH:$(pwd)/tectonic-installer/darwin # Put the `terraform` binary on PATH
```

#### Copy and configure .terraformrc

Make a copy of the Terraform configuration file. Do not share this configuration file as it is specific to the install host.

```bash
$ sed "s|<PATH_TO_INSTALLER>|$INSTALLER_PATH|g" terraformrc.example > .terraformrc
$ export TERRAFORM_CONFIG=$(pwd)/.terraformrc
```

#### Get Terraform's Azure modules

Next, get the modules for the Azure platform that Terraform will use to create cluster resources:

```
$ terraform get platforms/azure
Get: file:///Users/tectonic-installer/modules/azure/vnet
Get: file:///Users/tectonic-installer/modules/azure/etcd
...
```

### Generate credentials with Azure CLI

Execute `az login` to obtain an authentication token. See the [Azure CLI docs][login] for more information. Once logged in, note the `id` field of the output from the `az login` command. This is a simple way to retrieve the Subscription ID for the Azure account.

Next, add a new role assignment for the Installer to use:

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

Export the following environment variables with values obtained from the output of the role assignment. As noted above, `ARM_SUBSCRIPTION_ID` is the `id` of the Azure account returned by `az login`.

```
# id field in az login output
$ export ARM_SUBSCRIPTION_ID=abc-123-456
# appID field in az ad output
$ export ARM_CLIENT_ID=generated-app-id
# password field in az ad output
$ export ARM_CLIENT_SECRET=generated-pass
# tenant field in az ad output
$ export ARM_TENANT_ID=generated-tenant
```

With the environment set, it's time to specify the deployment details for the cluster.

## Customize the deployment

Possible customizations to the base installation are listed in `examples/terraform.tfvars.azure`. Choose a cluster name to identify the cluster. Export an environment variable with the chosen cluster name. In this example, `my-cluster` is used.

```
$ export CLUSTER=my-cluster
```

Create a build directory for the new cluster and copy the example file into it:

```
$ mkdir -p build/${CLUSTER}
$ cp examples/terraform.tfvars.azure build/${CLUSTER}/terraform.tfvars
```

Edit the parameters in `build/$CLUSTER/terraform.tfvars` with the deployment's Azure details, domain name, license, and pull secret. [View all of the Azure specific options][azure-vars] and [the common Tectonic variables][vars].

### Key values for basic Azure deployment

These are the basic values that must be adjusted for each Tectonic deployment on Azure. See the details of each value in the `terraform.tfvars` file.

* `tectonic_admin_email` - For the initial Console login
* `tectonic_admin_password_hash` - Bcrypted value
* `tectonic_azure_client_secret` - As in `ARM_CLIENT_SECRET` above
* `tectonic_azure_ssh_key` - Full path the the public key part of the key added to `ssh-agent` above
* `tectonic_base_domain` - The DNS domain or subdomain delegated to an Azure DNS zone above
* `tectonic_azure_external_dns_zone_id` - Get with `az network dns zone list`
* `tectonic_cluster_name` - Usually matches `$CLUSTER` as set above
* `tectonic_license_path` - Full path to `tectonic-license.txt` file downloaded from Tectonic account
* `tectonic_pull_secret_path` - Full path to `config.json` container pull secret file downloaded from Tectonic account

## Deploy the cluster

Check the plan before deploying:

```
$ terraform plan -var-file=build/${CLUSTER}/terraform.tfvars platforms/azure
```

Next, deploy the cluster:

```
$ terraform apply -var-file=build/${CLUSTER}/terraform.tfvars platforms/azure
```

This should run for a short time.

## Access the cluster

When `terraform apply` is complete, the Tectonic console will be available at `https://my-cluster.example.com`, as configured in the cluster build's variables file.

### CLI cluster operations with kubectl

Cluster credentials, including any generated CA, are written beneath the `generated/` directory. These credentials allow connections to the cluster with `kubectl`:

```
$ export KUBECONFIG=generated/auth/kubeconfig
$ kubectl cluster-info
```

## Scale an existing cluster on Azure

To scale worker nodes, adjust `tectonic_worker_count` in the cluster build's `terraform.tfvars` file.

Use the `terraform plan` subcommand to check configuration syntax:

```
$ terraform plan \
  -var-file=build/${CLUSTER}/terraform.tfvars \
  -target module.workers \
  platforms/azure
```

Use the `apply` subcommand to deploy the new configuration:

```
$ terraform apply \
  -var-file=build/${CLUSTER}/terraform.tfvars \
  -target module.workers \
  platforms/azure
```

The new nodes should automatically show up in the Tectonic Console shortly after they boot.

## Delete the cluster

Deleting a cluster will remove only the infrastructure elements created by Terraform. For example, an existing DNS resource group is not removed.

To delete the Azure cluster specified in `build/$CLUSTER/terraform.tfvars`, run the following `terraform destroy` command:

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


[azure-cli]: https://docs.microsoft.com/en-us/cli/azure/install-azure-cli
[azure-dns]: https://docs.microsoft.com/en-us/azure/dns/dns-overview
[azure-vars]: ../../variables/azure.md
[bcrypt]: https://github.com/coreos/bcrypt-tool/releases/tag/v1.0.0
[conventions]: ../../conventions.md
[copy-docs]: https://www.terraform.io/docs/commands/apply.html
[domain-delegation]: https://docs.microsoft.com/en-us/azure/dns/dns-delegate-domain-azure-dns
[generic]: ../../generic-platform.md
[install-go]: https://golang.org/doc/install
[login]: https://docs.microsoft.com/en-us/cli/azure/get-started-with-azure-cli
[plan-docs]: https://www.terraform.io/docs/commands/plan.html
[register]: https://account.coreos.com/signup/summary/tectonic-2016-12
[release-notes]: https://coreos.com/tectonic/releases/
[troubleshooting]: ../../troubleshooting/installer-terraform.md
[vars]: ../../variables/config.md
