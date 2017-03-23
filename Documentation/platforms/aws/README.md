# Install Tectonic on AWS with Terraform

Following this guide will deploy a Tectonic cluster within your AWS account.

Generally, the AWS platform templates adhere to the standards defined by the project [conventions][conventions] and [generic platform requirements][generic]. This document aims to document the implementation details specific to the Azure platform.

## Prerequsities

 - **DNS** - Ensure that the DNS zone is already created and available in route53 for the account. For example if the `tectonic_base_domain` is set to `kube.example.com` a route53 zone must exist for this domain and the AWS nameservers must be configured for the domain.
 - **Make** - This guide uses `make` to download a customized version of Terraform, which is pinned to a specific version and includes required plugins.
 - **Tectonic Account** - Register for a [Tectonic Account][register], which is free for up to 10 nodes. You will need to provide the cluster license and pull secret below.

## Getting Started

First, download Terraform with via `make`. This will download the pinned Terraform binary and modules:

```
$ cd tectonic-installer
$ make terraform-download
```

After downloading, you will need to source this new binary in your `$PATH`. This is important, especially if you have another verison of Terraform installed. Run this command to add it to your path:

```
$ export PATH=/path/to/tectonic-installer/bin/terraform/terraform:$PATH
```

You can double check that you're using the binary that was just downloaded:

```
$ which terraform
/Users/coreos/tectonic-installer/bin/terraform/terraform
```

Next, get the modules that Terraform will use to create the cluster resources:

```
$ terraform get platforms/aws
```

Configure your Azure credentials. See the [AWS docs][env] for details.

```
$ export AWS_ACCESS_KEY_ID=
$ export AWS_SECRET_ACCESS_KEY=
```

Set your desired region:

```
$ export AWS_REGION=
```

Now we're ready to specify our cluster configuration.

## Customize the deployment

Use this example to customize your cluster configuration. A few fields require special consideration:

 - **tectonic_base_domain** - domain name that is set up with in a resource group, as described in the prerequisites.
 - **tectonic_pull_secret_path** - path on disk to your downloaded pull secret. You can find this on your [Account dashboard][account].
 - **tectonic_license_path** - path on disk to your downloaded Tectonic license. You can find this on your [Account dashboard][account].
 - **tectonic_admin_password_hash** - generate a hash with the [bcrypt-hash tool][bcrypt] that will be used for your admin user.

Here's an example of the full file:

**build/<cluster>/terraform.tfvars**

```
TODO: insert me
```

## Deploy the cluster

Test out the plan before deploying everything:

```
$ PLATFORM=aws CLUSTER=my-cluster make plan
```

Next, deploy the cluster:

```
$ PLATFORM=aws CLUSTER=my-cluster make apply
```

This should run for a little bit, and when complete, your Tectonic cluster should be ready.

If you encounter any issues, check the known issues and workarounds below.

To delete your cluster, run:

```
$ PLATFORM=aws CLUSTER=my-cluster make destroy
```

### Known issues and workarounds

See the [troubleshooting][troubleshooting] document for work arounds for bugs that are being tracked.

[conventions]: ../conventions.md
[generic]: ../generic-platform.md
[env]: http://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html#cli-environment
[register]: https://account.tectonic.com/signup/summary/tectonic-2016-12
[account]: https://account.tectonic.com
[bcrypt]: https://github.com/coreos/bcrypt-tool/releases/tag/v1.0.0
