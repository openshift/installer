# Installing Tectonic on Microsoft Azure

This guide steps through installing a Tectonic enterprise-ready Kubernetes cluster on Microsoft Azure. The resulting basic cluster has 1 controller and 3 worker nodes and is ready for application deployment, monitoring, and scaling.

## Prerequisites

### Azure account

Before beginning, create an [Azure account][azure-home] with a valid credit card.

### Azure CLI tool

The [`az` command line API tool][az] will be used to retrieve and generate credentials granting access to the Installer:

```sh
$ curl -L https://aka.ms/InstallAzureCli | bash
```

### Command line comfort

This guide assumes intermediate comfort with the Unix command line, including such tasks as setting environment variables and configuring `ssh-agent` with a key. The Tectonic installation process on Azure happens in the terminal.

### Domain name

This guide requires a domain name available for configuration, and the ability to follow steps to delegate it or a subdomain to Azure DNS name service.

## Configuring Azure DNS

The Azure DNS service allows you to perform DNS management, traffic management, availability monitoring and domain registration. DNS management is the only feature of Azure DNS required to install Tectonic.

### Create an Azure DNS Zone

When creating an Azure DNS Zone, enter a domain or subdomain that you own and can manage.

The Tectonic installation requires a domain or subdomain name in which it will create two [sub-]subdomains: one for the Tectonic console, and one for the Kubernetes API server. This allows Tectonic to access and use the listed domain. This tutorial employs the domain name `example.com`. The string `example.com` should be replaced with the domain or subdomain name configured in this step wherever it appears later in the tutorial.

1. From the menu at left, selet *DNS Zones*
2. Click *Add* to create a new zone.
3. Enter an existing, registered domain or subdomain name.
4. Click *Create*.

Azure provides 4 DNS nameservers for the new zone. The domain or sub-domain must be [configured to use these nameservers][azure-dns-delegate]. Visit the domain registrar to add the Azure NS records.

1. Go to the domain registrar’s website.
2. Go to the DNS settings page and enter the four Azure nameservers as NS records for the domain or subdomain.
3. Save the updated domain settings.

Note that it may take from a few minutes to several hours for the changes to take effect, depending on the TTL setting of existing NS records.

To verify which nameservers are associated with your domain, use a tool like `dig` or `nslookup`. If no nameservers are returned when you look up your domain, changes may still be pending. Here's an example command:

```bash
$ dig -t ns [example.com]
```

The nameservers are set up correctly when the lookup yields the four hostnames provided by Azure.

## Generating Azure authentication assets

Execute `az login` to obtain an authentication token. See the [Azure CLI docs][azlogin] for more information. Once logged in, note the `id` value of the output from the `az login` command. This is a simple way to retrieve the Subscription ID for the Azure account.

### Add Active Directory Service Principal role assignment

Next, add a new Active Directory (AD) Service Principal (SP) role assignment to grant Installer access to Azure:

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

Export the following environment variables with values obtained from the output of the SP role assignment. As noted above, `ARM_SUBSCRIPTION_ID` is the `id` of the Azure account returned by `az login`.

```
# The id value in az login output
$ export ARM_SUBSCRIPTION_ID=azure-acct-sub-id
# The appID value in az ad output
$ export ARM_CLIENT_ID=generated-app-id
# The password value in az ad output
$ export ARM_CLIENT_SECRET=generated-pass
# The tenant value in az ad output
$ export ARM_TENANT_ID=generated-tenant
```

## Add a key to ssh-agent

The next step in preparing the environment for installation is to add the key to be used for logging in to each cluster node during initialization to the local `ssh-agent`.

### ssh-agent

Ensure `ssh-agent` is running by listing the known keys:

```bash
$ ssh-add -L
```

Add the SSH private key that will be used for the Tectonic installation to `ssh-agent`:

```bash
$ ssh-add ~/.ssh/id_rsa
```

Verify that the SSH key identity is available to the ssh-agent:

```bash
$ ssh-add -L
```

## Tectonic Installer

### Register for Tectonic

Register for a [Tectonic Account][register], free for up to 10 nodes. The Tectonic license and container registry pull secret are required during installation. Download these credentials as the files `license.txt` and `config.json` from your Tectonic Account.

### Download and extract Tectonic Installer

Open a new terminal and run the following commands to download and extract Tectonic Installer:

```bash
$ curl -O https://releases.tectonic.com/tectonic-1.7.1-tectonic.2.tar.gz
$ tar xzvf tectonic-1.7.1-tectonic.2.tar.gz
$ cd tectonic
```

### Initialize and configure Installer and Terraform

#### Set INSTALLER_PATH

Start by setting `INSTALLER_PATH` to the location of the installation host's Tectonic installer platform. The platform should be one of `darwin` or `linux`.

```bash
$ export INSTALLER_PATH=$(pwd)/tectonic-installer/linux/installer # Edit for linux or darwin
$ export PATH=$PATH:$(pwd)/tectonic-installer/linux
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

#### Create a cluster build directory

Choose a cluster name to identify the cluster. Export an environment variable with the chosen cluster name. In this tutorial, `my-cluster` is used.

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

These are the basic values that must be adjusted for each Tectonic deployment on Azure. See the details of each value in the comments in the `terraform.tfvars` file.

* `tectonic_admin_email` - For the initial Console login
* `tectonic_admin_password_hash` - Use [`bcrypt-tool`][bcrypt-tool] to encrypt password
* `tectonic_azure_client_secret` - As in `ARM_CLIENT_SECRET` above
* `tectonic_azure_ssh_key` - Full path to the public key part of the key added to `ssh-agent` above
* `tectonic_azure_location` - e.g., `centralus`
* `tectonic_base_domain` - The DNS domain or subdomain delegated to an Azure DNS zone above
* `tectonic_azure_external_dns_zone_id` - Value of `id` in `az network dns zone list` output
* `tectonic_cluster_name` - Usually matches `$CLUSTER` as set above
* `tectonic_license_path` - Full path to `tectonic-license.txt` file downloaded from Tectonic account
* `tectonic_pull_secret_path` - Full path to `config.json` container pull secret file downloaded from Tectonic account

## Deploy the cluster

Validate the plan before deploying:

```
$ terraform plan -var-file=build/${CLUSTER}/terraform.tfvars platforms/azure
```

Deploy the cluster – aka `apply`:

```
$ terraform apply -var-file=build/${CLUSTER}/terraform.tfvars platforms/azure
```

The apply step will run for some time and prints status on the standard output.

## Access the cluster

When `terraform apply` is complete, the Tectonic console will be available at `https://my-cluster.example.com`, as configured in the cluster build's variables file.

[**NEXT:** Deploying an application on Tectonic][first-app]


[az]: https://docs.microsoft.com/en-us/cli/azure/install-azure-cli
[azlogin]: https://docs.microsoft.com/en-us/azure/xplat-cli-connect
[azure-dns-delegate]: https://docs.microsoft.com/en-us/azure/dns/dns-delegate-domain-azure-dns
[azure-home]: https://azure.microsoft.com/
[azure-vars]: ../../variables/azure.md
[bcrypt-tool]: https://github.com/coreos/bcrypt-tool/releases
[first-app]: first-app.md
[register]: https://account.coreos.com/signup/summary/tectonic-2016-12
[vars]: ../../variables/config.md
