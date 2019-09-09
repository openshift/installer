# Handling quota on GCP
All cloud providers have some concept of limits imposed to protect their services from misuse as well as protect customers from mistakes that can leads to excessive charges.  In GCP these will vary by region, project and account.  This is great news in terms of flexibility but can present an initial hurdle if the defaults in your project are too low.

## Baseline usage
A vanilla IPI installation will result in 24 CPUs, 3 Static IPs and 768 GB of storage consumed.

## Common problems
Be sure to consider cluster growth and consumption from other clusters if using a shared projectd.  The most likely areas of contention are CPU, Static IPs and Storage (SSD) quota.  Whenever an installation fails the installer CLI will return the relevant error message stating which quota was exceeded in a particular region.

## Increasing limits
To adjust quotas visit the [GCP console][gcp-console-quota] and make necessary changes.  This will likely involve filing a support ticket so it's best to plan ahead as this is often the most time consuming barrier to your first running cluster.  For more detailed information please refer to the [GCP documentation][gcp-docs-quota].

[gcp-console-quota]: https://console.cloud.google.com/iam-admin/quotas
[gcp-docs-quota]: https://cloud.google.com/compute/quotas
