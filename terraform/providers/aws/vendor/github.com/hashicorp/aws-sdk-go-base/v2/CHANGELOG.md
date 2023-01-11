<!-- markdownlint-disable single-title -->
# v2.0.0 (Unreleased)

# v2.0.0-beta.17 (2022-08-31)

ENHANCEMENTS

* Adds support for `max_attempts` in shared config files. ([#278](https://github.com/hashicorp/aws-sdk-go-base/pull/278))
* Prevents silent failures when `RoleARN` missing from `AssumeRole` or `AssumeRoleWithWebIdentity`. ([#277](https://github.com/hashicorp/aws-sdk-go-base/pull/277))
* Adds support for `SourceIdentity` with `AssumeRole`. ([#311](https://github.com/hashicorp/aws-sdk-go-base/pull/311))
* Adds support for AWS region `me-central-1`. ([#328](https://github.com/hashicorp/aws-sdk-go-base/pull/328))
* Adds support for passing HTTP User-Agent products in `useragent.Context`. ([#318](https://github.com/hashicorp/aws-sdk-go-base/pull/318))

# v2.0.0-beta.16 (2022-04-27)

BREAKING CHANGES

* Removes boolean `SkipEC2MetadataApiCheck` and adds `EC2MetadataServiceEnableState` of type `imds.ClientEnableState`. ([#240](https://github.com/hashicorp/aws-sdk-go-base/pull/240))

ENHANCEMENTS

* Adds support for assuming IAM role with web identity. ([#178](https://github.com/hashicorp/aws-sdk-go-base/pull/178))

# v2.0.0-beta.15 (2022-04-12)

ENHANCEMENTS

* Adds parameter `SuppressDebugLog` to suppress logging. ([#232](https://github.com/hashicorp/aws-sdk-go-base/pull/232))

# v2.0.0-beta.14 (2022-04-07)

ENHANCEMENTS

* Adds support for custom CA bundles in shared config files for AWS SDK for Go v1. ([#226](https://github.com/hashicorp/aws-sdk-go-base/pull/226))

# v2.0.0-beta.13 (2022-03-09)

NOTES

* Filters CR characters out of AWS SDK for Go v1 logs. ([#174](https://github.com/hashicorp/aws-sdk-go-base/pull/174))

# v2.0.0-beta.12 (2022-03-02)

NOTES

* Filters CR characters out of AWS SDK for Go v2 logs. ([#157](https://github.com/hashicorp/aws-sdk-go-base/pull/157))

# v2.0.0-beta.11 (2022-02-28)

BUG FIXES

* No longer overrides shared config and credentials files when using defaults. ([#151](https://github.com/hashicorp/aws-sdk-go-base/pull/151))

# v2.0.0-beta.10 (2022-02-25)

ENHANCEMENTS

* Adds logging for explicitly set authentication parameters. ([#146](https://github.com/hashicorp/aws-sdk-go-base/pull/146))
* Adds warning log when `Profile` and static credentials environment variables are set. ([#146](https://github.com/hashicorp/aws-sdk-go-base/pull/146))

# v2.0.0-beta.9 (2022-02-23)

BUG FIXES

* Now returns an error if an invalid profile is specified. ([#128](https://github.com/hashicorp/aws-sdk-go-base/pull/128))

ENHANCEMENTS

* Retrieves region from IMDS when credentials sourced from IMDS. ([#131](https://github.com/hashicorp/aws-sdk-go-base/pull/131))

# v2.0.0-beta.8 (2022-02-18)

BUG FIXES

* Restores expansion of `~/` in file paths. ([#118](https://github.com/hashicorp/aws-sdk-go-base/pull/118))
* Fixes error when setting custom CA bundle. ([#122](https://github.com/hashicorp/aws-sdk-go-base/pull/122))

ENHANCEMENTS

* Adds expansion of environment variables in file paths. ([#118](https://github.com/hashicorp/aws-sdk-go-base/pull/118))
* Updates list of valid regions. ([#111](https://github.com/hashicorp/aws-sdk-go-base/pull/111))
* Adds parameter `CustomCABundle`. ([#122](https://github.com/hashicorp/aws-sdk-go-base/pull/122))

# v2.0.0-beta.7 (2022-02-14)

BUG FIXES

* Updates HTTP client to correctly handle IMDS authentication from inside a container. ([#116](https://github.com/hashicorp/aws-sdk-go-base/pull/116))

# v2.0.0-beta.6 (2022-02-09)

BREAKING CHANGES

* Removes config parameter `DebugLogging` and always enables logging.
  Client applications are expected to filter logs by setting log levels. ([#97](https://github.com/hashicorp/aws-sdk-go-base/pull/97))

ENHANCEMENTS

* Adds support for setting maximum retries using environment variable `AWS_MAX_ATTEMPTS`. ([#105](https://github.com/hashicorp/aws-sdk-go-base/pull/105))

# v2.0.0-beta.5 (2022-01-31)

BUG FIXES

* Was not correctly setting additional user-agent string parameters on AWS SDK v1 `Session`. ([#95](https://github.com/hashicorp/aws-sdk-go-base/pull/95))

# v2.0.0-beta.4 (2022-01-31)

ENHANCEMENTS

* Adds support for IPv6 IMDS endpoints with parameter `EC2MetadataServiceEndpointMode` and environment variable `AWS_EC2_METADATA_SERVICE_ENDPOINT_MODE`. ([#92](https://github.com/hashicorp/aws-sdk-go-base/pull/92))
* Adds parameter `EC2MetadataServiceEndpoint` and environment variable `AWS_EC2_METADATA_SERVICE_ENDPOINT`.
  Deprecates environment variable `AWS_METADATA_URL`. ([#92](https://github.com/hashicorp/aws-sdk-go-base/pull/92))
* Adds parameter `StsRegion`. ([#91](https://github.com/hashicorp/aws-sdk-go-base/pull/91))
* Adds parameters `UseDualStackEndpoint` and `UseFIPSEndpoint`. ([#88](https://github.com/hashicorp/aws-sdk-go-base/pull/88))

BREAKING CHANGES

* Renames parameter `SkipMetadataApiCheck` to `SkipEC2MetadataApiCheck`. ([#92](https://github.com/hashicorp/aws-sdk-go-base/pull/92))
* Renames assume role parameter `DurationSeconds` to `Duration`. ([#84](https://github.com/hashicorp/aws-sdk-go-base/pull/84))

# v2.0.0-beta.3 (2021-11-03)

ENHANCEMENTS

* Adds parameter `UserAgent` to append to user-agent string. ([#86](https://github.com/hashicorp/aws-sdk-go-base/pull/86))

# v2.0.0-beta.2 (2021-09-27)

ENHANCEMENTS

* Adds parameter `HTTPProxy`. ([#81](https://github.com/hashicorp/aws-sdk-go-base/pull/81))
* Adds parameter `APNInfo` to add APN data to user-agent string. ([#82](https://github.com/hashicorp/aws-sdk-go-base/pull/82))

BREAKING CHANGES

* Moves assume role parameters to `AssumeRole` struct. ([#78](https://github.com/hashicorp/aws-sdk-go-base/pull/78))
