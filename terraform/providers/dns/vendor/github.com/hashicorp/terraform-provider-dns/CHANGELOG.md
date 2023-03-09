## 3.2.4 (February 06, 2023)

BUG FIXES:

* provider: Prevent panic by returning errors from miekg/dns to practitioners ([#233](https://github.com/hashicorp/terraform-provider-dns/pull/233))
* resource/dns_a_record_set: Strip leading zeros from IPv4 addresses ([#233](https://github.com/hashicorp/terraform-provider-dns/pull/233))
* resource/dns_aaaa_record_set: Strip leading zeros from IPv4 addresses ([#233](https://github.com/hashicorp/terraform-provider-dns/pull/233))

## 3.2.3 (March 21, 2022)

BUG FIXES:

* provider: Prevented potential EDNS TCP KeepAlive timeout issues ([#187](https://github.com/hashicorp/terraform-provider-dns/pull/187))
* provider: Prevented potential EDNS Expire issues ([#187](https://github.com/hashicorp/terraform-provider-dns/pull/195))
* provider: Prevented "cannot unmarshal DNS" error for responses without EDNS and greater than 512 bytes

## 3.2.2 (March 21, 2022)

NOTES:

* This release was not fully completed. Use 3.2.3 instead.

BUG FIXES:

* provider: Prevented potential EDNS TCP KeepAlive timeout issues ([#187](https://github.com/hashicorp/terraform-provider-dns/pull/187))
* provider: Prevented potential EDNS Expire issues ([#187](https://github.com/hashicorp/terraform-provider-dns/pull/195))
* provider: Prevented "cannot unmarshal DNS" error for responses without EDNS and greater than 512 bytes

## 3.2.1 (July 08, 2021)

DEPENDENCIES:

* Bump github.com/hashicorp/terraform-plugin-sdk/v2 from v2.6.1 to v2.7.0 ([#153](https://github.com/terraform-providers/terraform-provider-dns/issues/153))
* Bump github.com/miekg/dns from v1.1.36-0.20210109083720-731b191cabd1 to v1.1.43 ([#155](https://github.com/terraform-providers/terraform-provider-dns/issues/155))
* Bump github.com/jcmturner/gokrb5/v8 from v8.4.1 to v8.4.2 ([#155](https://github.com/terraform-providers/terraform-provider-dns/issues/155))

## 3.2.0 (June 14, 2021)

NEW FEATURES:

* Add debug logging for resource DNS messages ([#145](https://github.com/terraform-providers/terraform-provider-dns/issues/145))

## 3.1.0 (February 19, 2021)

Binary releases of this provider now include the darwin-arm64 platform. 

NEW FEATURES:

* Support for GSSAPI/Kerberos signed updates ([#30](https://github.com/terraform-providers/terraform-provider-dns/issues/30))

## 3.0.1 (January 11, 2021)

BUG FIXES:

* Prevent multiple TSIG being added during retries ([#116](https://github.com/terraform-providers/terraform-provider-dns/issues/116))

## 3.0.0 (October 14, 2020)

Binary releases of this provider will now include the linux-arm64 platform.

BREAKING CHANGES:

* Upgrade to version 2 of the Terraform Plugin SDK, which drops support for Terraform 0.11. This provider will continue to work as expected for users of Terraform 0.11, which will not download the new version. ([#110](https://github.com/terraform-providers/terraform-provider-dns/issues/110))

## 2.2.0 (July 24, 2019)

* **New Data Source:** `dns_srv_record_set` [#70](https://github.com/terraform-providers/terraform-provider-dns/issues/70)
* **New Resource:** `dns_srv_record_set` [#70](https://github.com/terraform-providers/terraform-provider-dns/issues/70)
* This release includes a stable version of Terraform SDK v0.12.5.
  The provider should still retain full backwards compatibility with Terraform v0.11.x.
* Fix SOA detection logic to cover DNS servers returning a non-SOA record. [#79](https://github.com/terraform-providers/terraform-provider-dns/issues/79)

## 2.1.1 (May 01, 2019)

* This release includes an upgraded Terraform SDK, for the sake of aligning the versions of the SDK amongst released providers, as we lead up to Core v0.12. This should have no noticeable impact on the provider.

## 2.1.0 (April 17, 2019)

NEW FEATURES:

* **New Data Source:** `dns_mx_record_set` ([#71](https://github.com/terraform-providers/terraform-provider-dns/issues/71))
* **New Resource:** `dns_mx_record_set` ([#71](https://github.com/terraform-providers/terraform-provider-dns/issues/71))
* **New Resource:** `dns_txt_record_set` ([#72](https://github.com/terraform-providers/terraform-provider-dns/issues/72))
* All resources can now be imported ([#37](https://github.com/terraform-providers/terraform-provider-dns/issues/37))
* Allow the creation of apex records ([#69](https://github.com/terraform-providers/terraform-provider-dns/issues/69))
* Retry DNS queries on timeout ([#68](https://github.com/terraform-providers/terraform-provider-dns/issues/68))

IMPROVEMENTS:

* The provider is now compatible with Terraform v0.12, while retaining compatibility with prior versions.

## 2.0.0 (May 25, 2018)

BACKWARDS INCOMPATIBILITIES / NOTES:

* Prior versions of the provider would sign requests when sending updates to a DNS server but would not sign the requests to read those values back on subsequent refreshes. For consistency, now _read_ requests are also signed for managed resources in this provider. This does not apply to the data sources, which continue to just send normal unsigned DNS requests as before.

NEW FEATURES:

* Use signed requests when refreshing managed resources ([#35](https://github.com/terraform-providers/terraform-provider-dns/issues/35))
* data/dns_ptr_record_set: Implement data source for PTR record. ([#32](https://github.com/terraform-providers/terraform-provider-dns/issues/32))

BUGS FIXED:

* Normalize IP addresses before comparing them, so non-canonical forms don't cause errant diffs ([#13](https://github.com/terraform-providers/terraform-provider-dns/issues/13))
* Validates zone names are fully qualified and that record names are not as these mistakes seems to be a common source of misconfiguration ([#36](https://github.com/terraform-providers/terraform-provider-dns/issues/36))
* Properly handle IPv6 IP addresses as the update host. Previously this would create an invalid connection address due to not properly constructing the address format. ([#22](https://github.com/terraform-providers/terraform-provider-dns/issues/22))
* When refreshing DNS record resources, `NXDOMAIN` errors are now properly marked as deletions in state rather than returning an error, thus allowing Terraform to plan to re-create the missing records. ([#33](https://github.com/terraform-providers/terraform-provider-dns/issues/33))
* Now checks the type of record returned to prevent unexpected values causing a panic ([#39](https://github.com/terraform-providers/terraform-provider-dns/issues/39))

## 1.0.0 (September 15, 2017)

* No changes from 0.1.1; just adjusting to [the new version numbering scheme](https://www.hashicorp.com/blog/hashicorp-terraform-provider-versioning/).

## 0.1.1 (August 28, 2017)

NEW FEATURES:

* **`dns_aaaa_record_set` data source** for fetching IPv6 address records ([#9](https://github.com/terraform-providers/terraform-provider-dns/issues/9))
* **`dns_ns_record_set` data source** for fetching nameserver records ([#10](https://github.com/terraform-providers/terraform-provider-dns/issues/10))
* **`dns_ns_record_set` resource** for creating new nameserver records via the DNS update protocol ([#10](https://github.com/terraform-providers/terraform-provider-dns/issues/10))

## 0.1.0 (June 20, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
