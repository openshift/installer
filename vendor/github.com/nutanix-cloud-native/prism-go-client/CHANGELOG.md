# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
- Added the "environment" package which provides abstraction for retrieving settings like API endpoints, credentials and their sources to evolve independently from clients
- Add logr based configurable logging for internal.NewClient

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

