# Cluster API Infrastructure Providers

This directory contains the Cluster API infrastructure providers leveraged
by the Installer to provision cluster infrastructure without an external
control plane.

## Adding a new terraform provider

To add a new Cluster API provider, create a directory under ./cluster-api/providers. For example, if you want to add a
provider named mycloud, then you would create the ./cluster-api/providers/mycloud directory.

Create a tools.go file with an unnamed import to the main package of the provider. 
Add a go.mod file with a require statement to the provider module.

For example, let's say that you are adding the mycloud provider. The mycloud provider is in the
github.com/openshift/cluster-api-provider-mycloud repository. The main package is
github.com/openshift/cluster-api-provider-mycloud. The version is v1.2.3.

```go
// ./cluster-api/providers/mycloud/tools.go
package main

import (
	_ "github.com/openshift/cluster-api-provider-mycloud"
)
```

```go
// ./cluster-api/providers/mycloud/go.mod
module github.com/openshift/installer/cluster-api/providers/mycloud

go 1.20

require github.com/openshift/cluster-api-provider-mycloud v1.2.3
```

After creating the directory and adding the tools.go and go.mod file, run the following command.
1. `make -C terraform go-mod-tidy-vendor.mycloud`

When there is a tools.go file, the Makefile will build the provider from the package referenced in the unnamed
import in the tools.go file. For the example mycloud provider, the provider will be built from the
./cluster-api/providers/mycloud/vendor/github.com/openshift/cluster-api-provider-mycloud directory.
