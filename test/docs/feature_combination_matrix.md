# Feature Combination Matrix

The [platform feature matrix](platform_feature_matrix.md) tracks one feature at
a time. Most escaped bugs, though, live where two features meet: dual stack is
tested, custom endpoints are tested, and dual stack *with* custom endpoints is
tested by nobody. This matrix tracks the combinations worth testing together.

It is deliberately a short list. Every pair of features is a combination, and
enumerating them is useless; what belongs here is the combinations that are
**meaningful** - ones a real customer configures, or ones where the two
features touch the same installer code.

## When a combination belongs here

Add a row if any of these holds:

- A customer configuration or reference architecture calls for both together
- The two features write to the same asset, manifest, or cloud resource
- One feature changes the defaults or validation of the other
- A past bug came from their interaction

Do not add a row for features that are simply orthogonal. A combination that
nothing connects is noise, and noise makes the list unusable.

## Legend

| Code | Meaning |
|------|---------|
| AT | Automatically Tested - a CI job or unit test exercises the combination |
| MT | Manually Tested - a documented manual case covers it |
| PT | Partially Tested - one platform covers it, others do not |
| NT | Not Tested - the combination is supported but nothing exercises it |
| NA | Not Applicable - the combination cannot be configured |
| TBD | Not yet assessed - **needs an owner** |

## Networking

| Combination | Why it matters | Platforms | Coverage | Evidence |
|-------------|---------------|-----------|----------|----------|
| Dual stack + Agent-based install | Both rewrite the same machine network manifests | BM, VSP, NTX, OST | AT | Agent CI jobs cover dual-stack [1] |
| Dual stack + Azure Stack Hub | Stack Hub has its own networking stack | AZ | TBD | |
| Dual stack + custom service endpoints | Endpoint overrides must resolve on both families | AWS, IBM, PVS | TBD | |
| Private cluster + proxy | The proxy is the only egress path; both rewrite cluster networking | All | PT | Only `e2e-openstack-proxy` exists [2] |
| Private cluster + shared VPC (XPN) | Private zones live in the host project, not the service project | GCP | AT | `e2e-gcp-xpn-dedicated-dns-project` [3] |
| Existing VPC + user-provisioned DNS | Both suppress installer-managed resources | AWS, AZ, GCP | TBD | |
| Private cluster + user-provisioned DNS | No public zone and no installer zone at all | AWS, AZ, GCP | TBD | |

## Security and Compliance

| Combination | Why it matters | Platforms | Coverage | Evidence |
|-------------|---------------|-----------|----------|----------|
| FIPS + private cluster | The common regulated-customer configuration | All | TBD | FIPS has no platform-specific tests at all [4] |
| FIPS + disconnected install | Regulated *and* air-gapped; mirror registry plus FIPS boot | All | TBD | |
| KMS encryption + custom OS image | A customer-supplied image encrypted with a customer key | AWS, AZ, GCP | TBD | |
| Confidential compute + custom disk type | Confidential VMs restrict which disk types are legal | AWS, AZ, GCP | TBD | |
| Secure Boot + custom OS image | A user image must carry the right signatures | AZ, BM, NTX | TBD | |

## Topology

| Combination | Why it matters | Platforms | Coverage | Evidence |
|-------------|---------------|-----------|----------|----------|
| SNO + dual stack | No second node to fall back on if one family misconfigures | BM | AT | Agent CI jobs cover SNO dual-stack [1] |
| SNO + disconnected | The common edge deployment | BM | TBD | |
| Failure domains + existing VPC | Zones must line up with pre-existing subnets | AWS, AZ, GCP, VSP, NTX | TBD | |
| Multi-vCenter + failure domains | Zones spanning vCenters | VSP | AT | Covered by the multi-vCenter CI jobs [5] |

## Sovereign and Restricted Clouds

| Combination | Why it matters | Platforms | Coverage | Evidence |
|-------------|---------------|-----------|----------|----------|
| Sovereign cloud + private cluster | Sovereign clouds have no public DNS, so this is not optional | GCD | AT/MT | [GCD test plan](gcd/plans/gcd-feature.md), [gcd_private_dns_only.md](gcd/cases/gcd_private_dns_only.md) |
| Sovereign cloud + custom OS image | No access to public image projects, so a user image is mandatory | GCD | AT/MT | [gcd_node_machine_and_disk.md](gcd/cases/gcd_node_machine_and_disk.md) |
| Sovereign cloud + constrained machine types | Only C3/M3/A3 exist, which changes the disk defaults too | GCD | AT/MT | [gcd_node_machine_and_disk.md](gcd/cases/gcd_node_machine_and_disk.md) |
| Sovereign cloud + service endpoints | Forbidden combination - validation must reject it | GCD | AT | Unit tested in `validateServiceEndpointOverride()` |
| GovCloud + FIPS | The usual US federal configuration | AWS, AZ | TBD | |

## Footnotes

1. 11 agent-based CI jobs cover compact, HA, SNO, dual-stack, and PXE
   configurations across multiple platforms
2. Proxy is configured at the platform-agnostic `installconfig.Proxy` level;
   OpenStack has the only dedicated proxy CI job
3. GCP shared VPC via `NetworkProjectID`, with `e2e-gcp-ovn-xpn` and
   `e2e-gcp-xpn-dedicated-dns-project`
4. FIPS is a platform-agnostic boolean in install-config with no
   platform-specific validation tests
5. vSphere has 16 e2e CI jobs including multi-vCenter, static IPs, and host
   groups

## How to Update This Matrix

- A feature test plan that ships a feature which meaningfully combines with an
  existing one must add the row here, in the same PR (see
  [README](README.md#test-matrices)).
- `NT` rows are the useful output of this document. Review them when planning
  a release: an `NT` on a combination customers actually run is a gap worth
  funding.
- Cite evidence when you move a cell off `TBD`. A cell with no evidence column
  is indistinguishable from a guess.
