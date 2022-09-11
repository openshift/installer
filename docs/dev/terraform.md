# Embedded Terraform in the Installer

The installer runs terraform in order to provision the infrastructure for a cluster. However, there is a design goal
for the installer that the user performing the installation should not be required to have any software pre-installed
in their environment. This design goal precludes the user from, for example, being required to have any containerization
software available. To get around this for terraform, the installer embeds within its binary the needed terraform
binaries. At installation time, the terraform binaries are extracted from the installer into the local file system.

## Design

### Building terraform binaries

The terraform binary and all of the terraform provider binaries are built from the ./terraform directory. For example,
the aws terraform provider is built from ./terraform/providers/aws. The terraform binary is built from
./terraform/terraform. Each of the terraform providers is compressed into its own zip file to reduce the size of the
installer binary. Terraform can use the zipped terraform providers as is without the installer having to unzip the
providers first.

### Embedding terraform binaries

After the terraform providers are built, they are each copied into the ./pkg/terraform/providers/mirror directory. The
./pkg/terraform/providers directory is embedded into the installer. To make things simpler, we automatically assign
every terraform provider to the openshift/local registry and use a version of 1.0.0. This keeps us from having to track
the actual registry name of the provider or the actual version used. This does not have any effect on the behavior of
terraform since the installer tightly controls which providers are made available. As an example, when building for
linux on amd64, the aws terraform provider would be placed at
./pkg/terraform/providers/mirror/openshift/local/aws/terraform-provider-aws_1.0.0_linux_amd64.zip. The terraform binary
is copied into ./pkg/terraform/terraform to be embedded into the installer as well.

### Extracting terraform binaries

When the installer is provisioning the infrastructure for a cluster, the installer extracts the necessary terraform
binaries to the local file system. Every terraform stage must define which terraform providers are needed for the stage.
The installer will use that information to (1) build a version.tf file to tell terraform which providers to use and (2)
extract only the needed providers.

For example, the set of terraform providers required for the gcp stages are the google and ignition terraform providers.
The installer will extract the terraform binary to ${install-dir}/terraform/bin/terraform. The installer will extract
the google and ignition terraform providers to ${install-dir}/terraform/plugins. The terraform exec is configured to use
${install-dir}/terraform/plugins as the explicit directory for the plugins.

## Adding a new terraform provider

To add a new terraform provider, create a directory under ./terraform/providers. For example, if you want to add a
terraform provider named mycloud, then you would create the ./terraform/providers/mycloud directory.

### Public terraform provider

If the terraform provider you want to add is a public terraform provider, create a tools.go file with an unnamed import
to the main package of the terraform provider. Add a go.mod file with a require statement to the terraform provider
module.

For example, let's say that you are adding the mycloud terraform provider. The mycloud terraform provider is in the
github.com/hashicorp/terraform-provider-mycloud repository. The main package of the terraform provider is
github.com/hashicorp/terraform-provider-mycloud. The version of the terraform provider is v1.2.3.

```go
// ./terraform/providers/mycloud/tools.go
package main

import (
	_ "github.com/hashicorp/terraform-provider-mycloud"
)
```

```go
// ./terraform/providers/mycloud/go.mod
module github.com/openshift/installer/terraform/providers/mycloud

go 1.18

require github.com/hashicorp/terraform-provider-mycloud v1.2.3
```

After creating the directory and adding the tools.go and go.mod file, run the following command.
1. `make -C terraform go-mod-tidy-vendor.mycloud`

When there is a tools.go file, the Makefile will build the terraform provider from the package referenced in the unnamed
import in the tools.go file. For the example mycloud terraform provider, the terraform provider will be built from the
./terraform/providers/mycloud/vendor/github.com/hashicorp/terraform/provider-mycloud directory.

### Custom terraform provider

If the terraform provider you want to add is a custom terraform provider, then create a main.go file for the custom
terraform provider.

When there is a main.go file, the Makefile will build the terraform provider from the top-level package of the custom
terraform provider.

## Update the version of terraform or a provider

To update the version of terraform or a provider, modify the version in the go.mod file for the relevant sub-module.
To update the version of terraform, change the version of the required module in ./terraform/terraform/go.mod. To update
the version of the aws terraform provider, change the version of the required module in ./terraform/providers/aws/go.mod.

After updating the require statement in the go.mod file, perform the following actions in the sub-module.
1. Remove the require stanza in the go.mod file for the indirect requires.
2. Match the replaces in the go.mod file with the replaces from the go.mod file in the upstream terraform provider.
3. Run `make -C terraform go-mod-tidy-vendor`.
