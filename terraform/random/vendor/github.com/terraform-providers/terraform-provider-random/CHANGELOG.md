## 3.1.0 (February 19, 2021)

Binary releases of this provider now include the darwin-arm64 platform. This version contains no further changes.

## 3.0.1 (January 12, 2021)

BUG FIXES:

* `resource_integer`: Integers in state that do not cleanly fit into float64s no longer lose their precision ([#132](https://github.com/terraform-providers/terraform-provider-random/issues/132))

## 3.0.0 (October 09, 2020)

Binary releases of this provider will now include the linux-arm64 platform.

BREAKING CHANGES:

* Upgrade to version 2 of the Terraform Plugin SDK, which drops support for Terraform 0.11. This provider will continue to work as expected for users of Terraform 0.11, which will not download the new version. ([#118](https://github.com/terraform-providers/terraform-provider-random/issues/118))
* Remove deprecated `b64` attribute ([#118](https://github.com/terraform-providers/terraform-provider-random/issues/118))

## 2.3.1 (October 26, 2020)

NOTES: This version is identical to v2.3.0, but has been compiled using Go v1.14.5 to fix https://github.com/hashicorp/terraform-provider-random/issues/120.

## 2.3.0 (July 07, 2020)

NOTES:

* The provider now uses the binary driver for acceptance tests ([#99](https://github.com/terraform-providers/terraform-provider-random/issues/99))

NEW FEATURES:

* Added import handling for `random_string` and `random_password` ([#104](https://github.com/terraform-providers/terraform-provider-random/issues/104))

## 2.2.1 (September 25, 2019)

NOTES:

* The provider has switched to the standalone TF SDK, there should be no noticeable impact on compatibility. ([#76](https://github.com/terraform-providers/terraform-provider-random/issues/76))

## 2.2.0 (August 08, 2019)

NEW FEATURES:

* `random_password` is similar to `random_string` but is marked sensitive for logs and output [[#52](https://github.com/terraform-providers/terraform-provider-random/issues/52)] 

## 2.1.2 (April 30, 2019)

* This release includes another Terraform SDK upgrade intended to align with that being used for other providers as we prepare for the Core v0.12.0 release. It should have no significant changes in behavior for this provider.

## 2.1.1 (April 12, 2019)

* This release includes only a Terraform SDK upgrade intended to align with that being used for other providers as we prepare for the Core v0.12.0 release. It should have no significant changes in behavior for this provider.

## 2.1.0 (March 20, 2019)

IMPROVEMENTS:

* The provider is now compatible with Terraform v0.12, while retaining compatibility with prior versions.

## 2.0.0 (August 15, 2018)

BACKWARDS INCOMPATIBILITIES / NOTES:
* `random_string`: set the ID for random_string resources to "none". Any terraform configuration referring to `random_string.foo.id` will need to be updated to reference `random_string.foo.result` ([#17](https://github.com/terraform-providers/terraform-provider-random/issues/17))

NEW FEATURES:

* `random_uuid` generates random uuid string that is intended to be used as unique identifiers for other resources ([#38](https://github.com/terraform-providers/terraform-provider-random/issues/38))

BUG FIXES: 
* Use UnixNano() instead of Unix() for the current time seed in NewRand() ([#27](https://github.com/terraform-providers/terraform-provider-random/issues/27))
* `random_shuffle`: if `random_shuffle` is given an empty list, it will return an empty list

IMPROVEMENTS:

* Replace ReadPet function in `resource_pet` with schema.Noop ([#34](https://github.com/terraform-providers/terraform-provider-random/issues/34))

## 1.3.1 (May 22, 2018)

BUG FIXES:

* Add migration and new schema version for `resource_string` ([#29](https://github.com/terraform-providers/terraform-provider-random/issues/29))

## 1.3.0 (May 21, 2018)

BUG FIXES:

* `random_integer` now supports update ([#25](https://github.com/terraform-providers/terraform-provider-random/issues/25))

IMPROVEMENTS:

* Add optional minimum character constraints to `random_string` ([#22](https://github.com/terraform-providers/terraform-provider-random/issues/22))

## 1.2.0 (April 03, 2018)

NEW FEATURES:

* `random_integer` and `random_id` are now importable. ([#20](https://github.com/terraform-providers/terraform-provider-random/issues/20))

## 1.1.0 (December 01, 2017)

NEW FEATURES:

* `random_integer` resource generates a single integer within a given range. ([#12](https://github.com/terraform-providers/terraform-provider-random/issues/12))

## 1.0.0 (September 15, 2017)

NEW FEATURES:

* `random_string` resource generates random strings of a given length consisting of letters, digits and symbols. ([#5](https://github.com/terraform-providers/terraform-provider-random/issues/5))

## 0.1.0 (June 21, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
