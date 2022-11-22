# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

### Changed

## [0.9.0] - 2022-09-27

### Added

- Added support for tracing via OpenTelemetry.

## [0.8.1] - 2022-09-26

### Changed

- Fixed bug for http go where response handler was overwritten in context object.

## [0.8.0] - 2022-09-22

### Added

- Added support for constructing a proxy authenticated client.

## [0.7.2] - 2022-09-09

### Changed

- Updated reference to abstractions.

## [0.7.1] - 2022-09-07

### Added

- Added support for additional status codes.

## [0.7.0] - 2022-08-24

### Added

- Adds context param in send async methods

## [0.6.2] - 2022-08-30

### Added

- Default 100 secs timeout for all request with a default context.

## [0.6.1] - 2022-08-29

### Changed

- Fixed a bug where an error would be returned for a 201 response with described response.

## [0.6.0] - 2022-08-17

### Added

- Adds a chaos handler optional middleware for tests

## [0.5.2] - 2022-06-27

### Changed

- Fixed an issue where response error was ignored for Patch calls

## [0.5.1] - 2022-06-07

### Changed

- Updated abstractions and yaml dependencies.

## [0.5.0] - 2022-05-26

### Added

- Adds support for enum or enum collections responses

## [0.4.1] - 2022-05-19

### Changed

- Fixed a bug where CAE support would leak connections when retrying.

## [0.4.0] - 2022-05-18

### Added

- Adds support for continuous access evaluation.

## [0.3.0] - 2022-04-19

### Changed

- Upgraded to abstractions 0.4.0.
- Upgraded to go 18.

## [0.2.0] - 2022-04-08

### Added

- Added support for decoding special characters in query parameters names.

## [0.1.0] - 2022-03-30

### Added

- Initial tagged release of the library.
