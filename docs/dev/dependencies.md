# Managing Dependencies

## Build Dependencies

The following dependencies must be installed on your system before you can build the installer.

### Fedora, CentOS, RHEL

```sh
sudo dnf install gcc-c++ zip
```

If you need support for [libvirt destroy](libvirt/README.md#cleanup), you should also install `libvirt-devel`.

### Go

You can follow the official [Go Installation Doc](https://go.dev/doc/install) to install a compatible go version.

## Go Dependencies

We follow a hard flattening approach; i.e. direct and inherited dependencies are installed in the base `vendor/`.

Dependencies are managed with [Go Modules](https://github.com/golang/go/wiki/Modules) but committed directly to the repository.

- Add or update a dependency with `go get <dependency>@<version>`.
- If you want to use a fork of a project or ensure that a dependency is not updated even when another dependency requires a newer version of it, manually add a [replace directive in the go.mod file](https://github.com/golang/go/wiki/Modules#when-should-i-use-the-replace-directive). 
- Run `go mod tidy` to tidy `go.mod` and update `go.sum`, then commit the changes.
- Run `go mod vendor` to re-vendor and then commit updated vendored code separately.

If your vendor bump touched `github.com/openshift/api`, also run `go generate ./pkg/types/installconfig.go` to update [`data/data/install.openshift.io_installconfigs.yaml`](/data/data/install.openshift.io_installconfigs.yaml).

This [guide](https://github.com/golang/go/wiki/Modules#how-to-use-modules) is a great source to learn more about using `go mod`.

For the sake of your fellow reviewers, commit vendored code separately from any other changes.

### CAPI Provider Controller Dependencies

The installer uses Cluster API controllers to provision the cluster insfrastructure (more details [here](./cluster-api.md)). Each provider are vendored to `cluster-api/providers/<provider-name>` (e.g. `cluster-api/providers/aws` for AWS platform).

Whenever there are changes in the CAPI provider projects that need to be available to the installer, we need to bump the provider version in the installer project. The following describes the process to do so. We will use AWS provider (CAPA) as an example, but other providers should follow the same steps.

#### Bump provider `go.mod`

We need to set the version of CAPI provider to the revisions with the changes in the provider `go.mod`.

```bash
# At the project root (i.e. <path>/installer)
cd cluster-api/providers/aws/
# For example:
# - With commit ref: go get sigs.k8s.io/cluster-api-provider-aws/v2@17a09f591
# - With a released version, go get sigs.k8s.io/cluster-api-provider-aws/v2@v2.7.1
go get sigs.k8s.io/cluster-api-provider-aws/v2@<revision>
go mod tidy && go mod vendor
```

**Important**: If CAPI provider project has a `replace` block, we need to make sure provider `go.mod` matches that. For example, the [CAPA `go.mod`](https://github.com/kubernetes-sigs/cluster-api-provider-aws/blob/5b4f7c1acb52392ec4121b4145468a4bbecf26fd/go.mod#L7-L13) and [provider `go.mod`](https://github.com/openshift/installer/blob/73b390363572305eed214627abca83b5b00d70d3/cluster-api/providers/aws/go.mod#L165-L171) must match if any.

#### Bump installer top-level `go.mod`

We need to also set the same version of CAPI provider in the installer top-level `go.mod`.
```bash
# At the project root (i.e. <path>/installer)
# For example:
# - With commit ref: go get sigs.k8s.io/cluster-api-provider-aws/v2@17a09f591
# - With a released version, go get sigs.k8s.io/cluster-api-provider-aws/v2@v2.7.1
go get sigs.k8s.io/cluster-api-provider-aws/v2@<revision>
go mod tidy && go mod vendor
```

#### Update infrastructure CRDs

When bumping version of the CAPI provider, some infrastructure CRDs might be updated and we need to ensure the installer is aware of that as it keeps a copy of the CRD in directory `data/data/cluster-api/`.

An easy way to update the infrastructure manifests is to run `./hack/verify-capi-manifests.sh [provider-dir]`:

```bash
# Update infrastructure manifests for all providers
./hack/verify-capi-manifests.sh
# Or update infrastructure manifest for a specific provider (e.g. aws)
./hack/verify-capi-manifests.sh cluster-api/providers/aws
```

Another way is to manually do it yourself. First, clone the upstream CAPI provider project.

```bash
# We need to have a copy of the CAPI provider project if not yet
git clone git@github.com:kubernetes-sigs/cluster-api-provider-aws.git
cd cluster-api-provider-aws
# Checkout the revision with changes
git fetch origin && git checkout <revision>
```

Next, generate the CRD manifest. The manifest will be available in `out/infrastructure-components.yaml`. We then copy it over to the installer project's corresponding infrastructure manifest.

```bash
make release-manifests
cp out/infrastructure-components.yaml <path-to>/installer/data/data/cluster-api/aws-infrastructure-components.yaml
```

**Important**: If using a released version of CAPI provider, the CRD is available for downloading from the [release page](https://github.com/kubernetes-sigs/cluster-api-provider-aws/releases).
