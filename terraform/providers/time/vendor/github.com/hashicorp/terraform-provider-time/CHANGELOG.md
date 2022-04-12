## 0.7.2 (July 01, 2021)

BUG FIXES:

* resource/time_sleep: Prevent `context deadline exceeded` error when timeout duration is configured above 20 minutes ([#45](https://github.com/hashicorp/terraform-provider-time/issues/45))

## 0.7.1 (May 04, 2021)

BUG FIXES:

* provider: Ensure `darwin/arm64` platform is included in releases

## 0.7.0 (February 19, 2021)

Binary releases of this provider now include the darwin-arm64 platform. This version contains no further changes.

## 0.6.0 (October 04, 2020)

BREAKING CHANGES:

* Dropped support for Terraform 0.11 and lower ([#16](https://github.com/hashicorp/terraform-provider-time/issues/16))

ENHANCEMENTS

* Made `time_sleep` context aware, allowing easier early cancellation ([#16](https://github.com/hashicorp/terraform-provider-time/issues/16))

## 0.5.0 (May 13, 2020)

FEATURES

* **New Resource:** `time_sleep` ([#12](https://github.com/hashicorp/terraform-provider-time/issues/12))

# 0.4.0 (April 21, 2020)

BREAKING CHANGES:

* resource/time_offset: `keepers` argument renamed to `triggers`
* resource/time_offset: Remove non-RFC3339 RFC and `unixdate` attribute
* resource/time_rotating: `keepers` argument renamed to `triggers`
* resource/time_rotating: Remove non-RFC3339 RFC and `unixdate` attributes
* resource/time_static: `keepers` argument renamed to `triggers`
* resource/time_static: Remove non-RFC3339 RFC and `unixdate` attributes

# v0.3.0

ENHANCEMENTS:

* resource/time_offset: Add `keepers` argument
* resource/time_rotating: Add `keepers` argument
* resource/time_static: Add `keepers` argument

BUG FIXES:

* resource/time_offset: Ensure `base_rfc3339` is always set in Terraform state during creation, even if unconfigured

# v0.2.0

BREAKING CHANGES:

* resource/time_static: The `expiration_` arguments have been moved to the new `time_rotating` resource as `rotation_` arguments.

FEATURES:

* **New Resource:** `time_offset`
* **New Resource:** `time_rotating`

# v0.1.0

FEATURES:

* **New Resource:** `time_static`
