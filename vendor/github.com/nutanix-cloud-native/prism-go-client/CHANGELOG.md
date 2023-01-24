# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
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

