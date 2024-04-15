# Controller Integration Tests

Primarily this folder contains integration tests, used to exercise the ASO controller to ensure that everything works properly. This includes tests for resources as well as for the controller itself. (This folder also contains resource registration, in `controller_resources.go` and `controller_resources_gen.go`).

## Test authorship

Always write a new test to exercise a new resource or resource version.

We shouldn't assume that, just because an earlier version worked fine, the new version will as well. There's always the possibility that the way a particular resource provider behaves with a new version will require a change in the was we interact with it.

As a minimum, we want to have tests for

* the latest `stable` version of the resource;
* the prior `stable` version of the resource; and
* the latest `preview` version of the resource.

Given that we don't want to have to maintain tests for every version of every resource, and each additional test makes our CI test suite take longer, consider removing tests for older versions of resources when we add tests for newer versions. This is a judgment call, and we recommend discussion with the team first.

## File naming

Test files should use the following naming convention:

``` bash
<group>_<subject>_<scenario>_<version>_test.go
```

Where:

* `<group>` is the Kubernetes resource group under test.  
  Lowercase. Use the actual name of the group, e.g. `compute`, not `Compute`, nor `Microsoft.Compute`.
* `<subject>` is the subject under test, often a kubernetes resource within the given group.  
  Lowercase. Use the actual name of the resource, e.g. `account`, not `Account`.
* `<scenario>` is a (very) short indicator of the test scenario.  
  Typically camelCase for clarity.
  Common scenarios include `crud` (short for Create, Read, Update and Delete).
* `<version>` is the Kubernetes API version of the resource under test.  
  Use the actual version in full, not a contraction, and not the API version, e.g. `v1api20210101`, not `v1alpha1` or `2021-01-01`.

All parts are optional and should be omitted if they don't apply to the tests in that file.

If you have tests for multiple scenarios, create multiple files, one per scenario.

Some examples to illustrate:

| Filename                                                             | Breakdown                                                                                                          | Check that ASO can correctly                                                                                 |
| -------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------ |
| `batch_account_crud_v1api20210101_test.go`                           | Group: batch<br/>Subject: account<br/>Scenario: crud<br/>Version:&nbsp;v1api20210101                                  | create and maintain a Batch Account.                                                                         |
| `servicebus_namespace_basicSkuCrud_v1api20210101preview_test.go`     | Group: servicebus<br/>Subject: namespace<br/>Scenario: basicSkuCrud<br/>Version:&nbsp;v1api20210101preview            | create and maintain a Basic SKU ServiceBus Namespace.                                                        |
| `servicebus_namespace_standardSkuCrud_v1api20211101_test.go`         | Group: servicebus<br/>Subject: namespace<br/>Scenario: standardSkuCrud<br/>Version:&nbsp;v1api20211101                | create and maintain a Standard SKU ServiceBus Namespace.                                                     |
| `documentdb_databaseaccount_mongodbCrud_v1api20210515_test.go`       | Group: documentdb<br/>Subject: databaseaccount<br/>Scenario: mongodbCrud<br/>Version:&nbsp;v1api20210515              | create and maintain a CosmosDB account in MongoDB mode                                                       |
| `documentdb_databaseaccount_secretsFromAzure_v1api20210515_test.go` | Group: documentdb<br/>Subject: databaseaccount<br/>Scenario: cosmosdbSecretsFromAzure<br/>Version:&nbsp;v1api20210515 | Test that ASO can correctly create maintain a CosmosDB Database Account and retrieve access keys from Azure. |

### Additional Notes

The filename for all resource tests should include the version of the primary resource under test, even if it's currently the only version of the resource. We want to avoid the need to rename existing tests when a new version of a resource is imported. Never assume that ASO will only ever have one version of any specific resource.

Why use group-kind-version and not group-version-kind for naming? Not every test is strictly related to a specific resource version. Ordering of parts is based on making it as easy as possible for future maintainers to scan the list of test files and find the one they want.
