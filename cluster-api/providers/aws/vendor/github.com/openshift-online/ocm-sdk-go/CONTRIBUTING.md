# Contributing to the OCM SDK

## Releasing a new OCM API Model version

First, all changes to the [ocm-api-model](https://github.com/openshift-online/ocm-api-model) have been defined and reviewed. Then, the client types for the model need to be generated via `make update` target in the `ocm-api-model` project.

Once the client types and api description are merged in [ocm-api-model](https://github.com/openshift-online/ocm-api-model), the version of the ocm-api-model
must be incremented for consumption in ocm-sdk-go generation. The version is defined by the latest git tag.

Once all changes to the OCM API Model have been committed to the main branch and a new tag is pused a new release will be created automatically.Then, you will need to update these changes in the SDK. (See "Updating the OCM SDK" section)

### Validating model updates

If you would like to test the SDK against a *local version* use the following instructions:

Ensure ocm-sdk-go is cloned locally alongside your cloned ocm-api-model directory where changes are made.

Use the following commands to test you're locally generated client types:
```
go mod edit -replace=github.com/openshift-online/ocm-api-model/clientapi=/path/to/your/local/ocm-api-model/clientapi

go mod edit -replace=github.com/openshift-online/ocm-api-model/model=/path/to/your/local/ocm-api-model/model

make update
```

## Updating the OCM SDK

The OCM SDK can be generated simply by running the following after all changes have been made:

```shell
./hack/update-model.sh
make update
```

The `./hack/update-model.sh` script will ensure the `ocm-api-model` modules are all up to date with the latest version across the OCM-SDK project.
To verify that they are all in-sync one can use the `./hack/verify-model-version.sh` script.

One can add an optional commit SHA or version to the `./update-model.sh <vX.Y.Z>` script to update the go modules to a specific version.

Whenever an update is made, ensure that the corresponding example in [examples](examples) is also updated where
necessary. It is *highly recommended* that new endpoints have a new example created.

## Releasing a new OCM SDK Version

Releasing a new version requires submitting an MR for review/merge with an update to the `Version` constant in
[version.go](version.go). Additionally, update the [CHANGES.md](CHANGES.md) file to include the new version and
describe all changes included.

Below is an example CHANGES.md update:

```
== 0.1.39 Oct 7 2019

- Update to model 0.0.9:
  - Add `type` attribute to the `ResourceQuota` type.
  - Add `config_managed` attribute to the `RoleBinding` type.
```

Submit an MR for review/merge with the CHANGES.md and version.go update.

Finally, create and submit a new tag with the new version following the below example:

```shell
git checkout main
git pull
git tag -a -m 'Release 0.1.39' v0.1.39
git push origin v0.1.39
```

Note that a repository administrator may need to push the tag to the repository due to access restrictions.
