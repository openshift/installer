Terraform Cloud/Enterprise Go Client
==============================

[![Tests](https://github.com/hashicorp/go-tfe/actions/workflows/ci.yml/badge.svg)](https://github.com/hashicorp/go-tfe/actions/workflows/ci.yml)
[![GitHub license](https://img.shields.io/github/license/hashicorp/go-tfe.svg)](https://github.com/hashicorp/go-tfe/blob/main/LICENSE)
[![GoDoc](https://godoc.org/github.com/hashicorp/go-tfe?status.svg)](https://godoc.org/github.com/hashicorp/go-tfe)
[![Go Report Card](https://goreportcard.com/badge/github.com/hashicorp/go-tfe)](https://goreportcard.com/report/github.com/hashicorp/go-tfe)
[![GitHub issues](https://img.shields.io/github/issues/hashicorp/go-tfe.svg)](https://github.com/hashicorp/go-tfe/issues)

The official Go API client for [Terraform Cloud/Enterprise](https://www.hashicorp.com/products/terraform).

This client supports the [Terraform Cloud V2 API](https://www.terraform.io/docs/cloud/api/index.html).
As Terraform Enterprise is a self-hosted distribution of Terraform Cloud, this
client supports both Cloud and Enterprise use cases. In all package
documentation and API, the platform will always be stated as 'Terraform
Enterprise' - but a feature will be explicitly noted as only supported in one or
the other, if applicable (rare).

## Version Information

Almost always, minor version changes will indicate backwards-compatible features and enhancements. Occasionally, function signature changes that reflect a bug fix may appear as a minor version change. Patch version changes will be used for bug fixes, performance improvements, and otherwise unimpactful changes.

## Installation

Installation can be done with a normal `go get`:

```
go get -u github.com/hashicorp/go-tfe
```

## Usage

```go
import tfe "github.com/hashicorp/go-tfe"
```

Construct a new TFE client, then use the various endpoints on the client to
access different parts of the Terraform Enterprise API. For example, to list
all organizations:

```go
config := &tfe.Config{
	Token: "insert-your-token-here",
}

client, err := tfe.NewClient(config)
if err != nil {
	log.Fatal(err)
}

orgs, err := client.Organizations.List(context.Background(), nil)
if err != nil {
	log.Fatal(err)
}
```

## Documentation

For complete usage of the API client, see the [full package docs](https://pkg.go.dev/github.com/hashicorp/go-tfe).

## API Coverage

This API client covers most of the existing Terraform Cloud API calls and is updated regularly to add new or missing endpoints.

- [x] Account
- [x] Agents
- [x] Agent Pools
- [x] Agent Tokens
- [x] Applies
- [x] Audit Trails
- [x] Changelog
- [x] Comments
- [x] Configuration Versions
- [x] Cost Estimation
- [ ] Feature Sets
- [ ] Invoices
- [x] IP Ranges
- [x] Notification Configurations
- [x] OAuth Clients
- [x] OAuth Tokens
- [x] Organizations
- [x] Organization Memberships
- [x] Organization Tags
- [x] Organization Tokens
- [x] Plan Exports
- [x] Plans
- [x] Policies
- [x] Policy Checks
- [x] Policy Sets
- [x] Policy Set Parameters
- [ ] Private Registry
	- [x] Modules
	- [x] Providers
	- [x] Provider Provider Versions and Platforms
	- [ ] GPG Keys
- [x] Runs
- [x] Run Tasks
- [ ] Run Tasks Integration
- [x] Run Triggers
- [x] SSH Keys
- [x] Stability Policy
- [x] State Versions
- [x] State Version Outputs
- [ ] Subscriptions
- [x] Team Access
- [x] Team Membership
- [x] Team Tokens
- [x] Teams
- [x] User Tokens
- [x] Users
- [x] Variable Sets
- [x] Variables
- [ ] VCS Events
- [x] Workspaces
- [x] Workspace-Specific Variables
- [x] Workspace Resources
- [x] Admin
  - [x] Module Sharing
  - [x] Organizations
  - [x] Runs
  - [x] Settings
  - [x] Terraform Versions
  - [x] Users
  - [x] Workspaces


## Examples

See the [examples directory](https://github.com/hashicorp/go-tfe/tree/main/examples).

## Running tests

See [TESTS.md](docs/TESTS.md).

## Issues and Contributing

See [CONTRIBUTING.md](docs/CONTRIBUTING.md)

## Releases

See [RELEASES.md](docs/RELEASES.md)
