# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added

### Changed
## [v0.5.0] - 2024-07-29
### Added
- Added v4 Categories beta APIs to v4 client
- Added v4 VolumeGroups beta APIs to v4 client
- Added a cache for v4 clients in v4 package

### Changed
- Updated v4 API clients from v4 alpha to v4 beta  APIs
- Handle trust bundle in v4 client cache GetOrCreate by setting VerifySSL

## [0.4.0] - 2024-05-03
### Added
- Added support for v4 client creation.
- Added support for getting information about an AZ given a uuid.
- Added support for getting a projection of attributes of entities using the 'groups' API endpoint.
- Added support for creating, deleting, listing, and getting the status of recovery plan jobs.
- Add optional function options for the NewKarbonAPIClient constructor
- Add ClusterRegistration interface in karbon package
- Add ClusterRegistration SetInfo and Cluster Addon SetInfo APIs
- Added support for specifying volume groups by category in a recovery plan create request.
- Added support for specifying primary and recovery clusters in a recovery plan.
- Added WithUserAgent client option for v3 client constructor.
- Added Cache for v3 Clients in v3 package.

### Changed
- Change the MetaService interface methods to take context.Context as a parameter
- Local environment provider now fetches port from `NUTANIX_PORT` environment variable
- Add logic to internal.Client for auto retry once after refreshing auth cookie on a 401 response in case of session auth.

### Removed
- remove the unexported method from the v3 service interface enabling mocking

## [0.3.4] - 2022-11-24
### Changed
- Bugfix: Stop explicit base64 decoding of BinaryData from ConfigMap in Kubernetes env provider

## [0.3.3] - 2022-11-24
### Changed
- Kubernetes env provider can now read the trust bundle from both BinaryData and Data


## [0.3.2] - 2022-11-04
### Changed
- Bugfix: Fix the kubebuilder enum annotations for NutanixTrustBundleKind

## [0.3.1] - 2022-11-03
### Added
- Add `AdditionalTrustBundle` property to the `NutanixPrismEndpoint` struct in environment/credential/types.go
- Add `AdditionalTrustBundle` property to the `ManagementEndpoint` struct in environment/types/types/go
- Add `WithPEMEncodedCertBundle` ClientOption for handling PEM Blocks in the v3 Constructor

### Changed
- Add license header to generated file and add a makefile target for generate
- `NUTANIX_INSECURE` and `NUTANIX_ADDITIONAL_TRUST_BUNDLE` environment variables are used to hydrate environment/local provider
- Store the certpool in the client to allow injecting multiple certificates using the `WithCertificate` option
- Use `hashicorp/go-cleanhttp` as the underlying http client constructor

## [0.3.0] - 2022-09-27
### Added
- Added the "environment" package which provides abstraction for retrieving settings like API endpoints, credentials and their sources to evolve independently from clients
- Add logr based configurable logging for internal.NewClient
- Add `WithCertificate` functional option for v3 client constructor
- Add `WithRoundTripper` functional option for v3 client to add custom interceptors
- Add `WithLogger` functional option for v3 client
- Add support for Nutanix-style credentials to "secretdir" environment provider

### Changed
- The http client has been moved from pkg/nutanix to repo root
- The fc stubs have been moved from pkg/nutanix/fc to fc
- The foundation stubs have been moved from pkg/nutanix/foundation to foundation
- The karbon stubs have been moved from pkg/nutanix/karbon to karbon
- The v3 stubs have been moved from pkg/nutanix/v3 to v3
- The underlying http client is moved from root package to internal
- The root package is renamed from `prism_go_client` to `prismgoclient`
- Modify NewClient constructor to use functional options
- NewRequest and NewAuthRequest methods on internal.Client don't admit context in params
- v3 client constructor now takes functional options as parameters
- v3 client constructor returns error instead of failing silently
- Add context to v3 interface method parameters to explicitly propagate context
- Bugfix in secretdir environment provider for tolerating symlinks

### Removed
- remove internal.NewBaseClient constructor

## [0.2.0] - 2022-06-14
### Added
- Add the fc stubs from nutanix/terraform-provider-nutanix
- Add the foundation stubs from nutanix/terraform-provider-nutanix
- Add the karbon stubs from nutanix/terraform-provider-nutanix
- Added GetCurrentLoggedInUser in pkg/nutanix/v3/v3_service.go

### Changed
- Updated the http client with the latest from github.com/nutanix/terraform-provider-nutanix
- Updated the v3 stubs with the latest from github.com/nutanix/terraform-provider-nutanix
- Updated the utils package with the latest from github.com/nutanix/terraform-provider-nutanix

### Removed
- Remove the compiled binary for the client from the source code
- Remove debug logs from pkg/nutanix/client.go


## [0.1.0] - 2022-06-08
### Added
- Initial copy of the v3 stubs and http client by [@vnephologist](https://github.com/vnephologist).

### Changed
- Change MessageResource.Details type from `map[string]string` and `map[string]interface{}` to `interface{}`

