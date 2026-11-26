# Sub-Platform Matrix

The [Platform Feature Matrix](platform_feature_matrix.md) answers "is this
feature tested on Azure?". It cannot answer "is this feature tested on Azure
Stack Hub?", because a single column collapses every variant of a platform into
one cell. This matrix adds that third dimension: feature x platform x
sub-platform.

A **sub-platform** is a variant of a platform that the installer has to treat
differently - a separate partition, cloud environment, API endpoint domain, or
provisioning path. Variants that only change configuration values (a different
region, a different machine type) are not sub-platforms.

Only platforms that actually have variants get a table. One table per parent
platform keeps each one readable and lets teams own their own section.

## Legend

Same codes as the [platform feature matrix](platform_feature_matrix.md), plus
one:

| Code | Meaning |
|------|---------|
| AT | Automatically Tested - covered by unit tests or CI e2e jobs |
| MT | Manually Tested - manual test procedures documented |
| AT/MT | Both automatic and manual tests exist |
| PT | Partially Tested - some paths tested, gaps remain |
| NT | Not Tested - implemented for this sub-platform, no coverage found |
| NA | Not Applicable - not implemented for this sub-platform |
| TBD | Coverage not yet assessed - **needs an owner to fill in** |

`TBD` is not a failure state. A new sub-platform column starts as all `TBD`,
and the platform team replaces cells as they assess them. A cell should never
be guessed: if you have not checked, leave it `TBD`.

## AWS

Partitions are defined in `pkg/types/aws/regions.go`.

| Feature | Commercial | GovCloud | China | ISO (C2S) | ISOB (SC2S) | EU Sovereign |
|---------|-----------|----------|-------|-----------|-------------|--------------|
| Partition identifier | `aws` | `aws-us-gov` | `aws-cn` | `aws-iso` | `aws-iso-b` | `aws` [1] |
| IPI Install | AT | TBD | TBD | TBD | TBD | TBD |
| UPI Install | AT | TBD | TBD | TBD | TBD | TBD |
| Cluster Destroy | NT | TBD | TBD | TBD | TBD | TBD |
| Private Cluster | AT | TBD | TBD | TBD | TBD | TBD |
| Service Endpoints | AT | TBD | TBD | TBD | TBD | TBD |
| Custom OS Image | NT | TBD | TBD | TBD | TBD | TBD |

## Azure

Cloud environments are the `CloudEnvironment` enum in
`pkg/types/azure/platform.go`.

| Feature | Public | US Government | China | German | Stack Hub |
|---------|--------|---------------|-------|--------|-----------|
| `cloudName` value | `AzurePublicCloud` | `AzureUSGovernmentCloud` | `AzureChinaCloud` | `AzureGermanCloud` | `AzureStackCloud` |
| IPI Install | AT | TBD | TBD | TBD | AT [2] |
| UPI Install | AT | TBD | TBD | TBD | TBD |
| Cluster Destroy | NT | TBD | TBD | TBD | TBD |
| Private Cluster | AT | TBD | TBD | TBD | TBD |
| Dual Stack | AT | TBD | TBD | TBD | TBD |
| Custom OS Image | NT | TBD | TBD | TBD | TBD |

## GCP

Sovereign environments are detected by `GetCloudEnvironment()` in
`pkg/types/gcp/platform.go`.

| Feature | Public GCP | GCD (Trusted Partner Cloud) | S3NS (Cloud de Confiance) |
|---------|-----------|-----------------------------|---------------------------|
| Detection | default | domain-scoped project + `u-` region | domain-scoped project + `u-` region |
| IPI Install | AT | AT [3] | TBD [4] |
| UPI Install | AT | NA | NA |
| Cluster Destroy | AT | NT [5] | TBD |
| Private Cluster | AT | AT/MT [3] | TBD |
| Shared VPC / XPN | AT | NA | NA |
| Service Endpoints / PSC | AT | NA [6] | NA [6] |
| Custom OS Image | AT | AT/MT [3] | TBD |
| Default Machine Types | AT | AT/MT [3] | TBD |

## vSphere

| Feature | Single vCenter | Multi-vCenter |
|---------|---------------|---------------|
| IPI Install | AT | AT [7] |
| UPI Install | AT | TBD |
| Failure Domains / Zones | AT | AT [7] |
| Static IPs | AT [7] | TBD |

## Bare Metal

These are provisioning paths rather than cloud variants, but they diverge
enough in the installer to need separate cells.

| Feature | IPI | Agent-based | Image-based (IBI) |
|---------|-----|-------------|-------------------|
| Source | `pkg/infrastructure/baremetal/` | `pkg/asset/agent/` | `pkg/asset/imagebased/` |
| Install | AT | AT | TBD |
| Single Node (SNO) | AT | AT | TBD |
| Dual Stack | AT | AT | TBD |
| Cluster Destroy | NT | NA | TBD |

## OpenStack

| Feature | RHOSP | PowerVC |
|---------|-------|---------|
| IPI Install | AT | AT [8] |
| UPI Install | AT | NA |
| Cluster Destroy | NT | TBD |
| Proxy Support | AT | TBD |
| Dual Stack | AT | NT |

## Platforms without sub-platforms

IBM Cloud, Power VS, and Nutanix each have a single variant today, so their
columns in the [platform feature matrix](platform_feature_matrix.md) are
complete on their own. Add a table here if that changes.

## Footnotes

1. AWS European Sovereign Cloud runs in the standard partition and is gated by
   the `AWSEuropeanSovereignCloudInstall` feature gate
   (`pkg/types/aws/validation/featuregates.go`); it is not a separate partition
   identifier
2. Azure Stack Hub has dedicated validation in
   `pkg/types/azure/validation/` (see `machinepool_test.go`, `disk_test.go`)
   and UPI templates in `upi/azurestack/`
3. See the [GCD feature test plan](gcd/plans/gcd-feature.md) and the manual
   cases in [`gcd/cases/`](gcd/cases/)
4. S3NS shares GCD's architecture and detection logic with different
   identifiers; explicitly out of scope in the GCD test plan
5. GCD shares the GCP destroy code in `pkg/destroy/gcp/` with no GCD-specific
   destroy tests; a manual case covers it
   ([gcd_destroy_cleanup.md](gcd/cases/gcd_destroy_cleanup.md))
6. Sovereign clouds forbid PSC endpoint overrides; the constraint itself is
   unit tested
7. vSphere has 16 e2e CI jobs including multi-vCenter, static IPs, and host
   groups; per-variant attribution of those jobs is not yet broken out
8. PowerVC shares the OpenStack CAPI provider

## How to Update This Matrix

- Adding a sub-platform means adding a column, initialized to `TBD`, not
  guessed values.
- A feature test plan that targets one sub-platform must update that column in
  the same PR (see [README](README.md#test-matrices)).
- When you change a cell away from `TBD`, add a footnote saying what you
  checked - a test file, a CI job name, or a case document.
