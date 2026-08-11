# Platform Feature Matrix

This matrix shows which installer features are tested for each platform and
how they are tested. Use it to identify coverage gaps and prioritize new test
documentation. The matrix is maintained manually and should be updated when
features, tests, or CI jobs change.

## Platform Key

| Abbr | Full Name | Notes |
|------|-----------|-------|
| AWS | Amazon Web Services | |
| AZ | Microsoft Azure | Includes Azure Stack |
| GCP | Google Cloud Platform | Public cloud |
| GCD | Google Cloud Dedicated | Sovereign cloud (Trusted Partner Cloud) |
| BM | Bare Metal | IPI and Agent-based |
| IBM | IBM Cloud | |
| NTX | Nutanix | |
| OST | Red Hat OpenStack Platform | |
| PVC | IBM PowerVC | Shares OpenStack CAPI provider |
| PVS | IBM Power Virtual Server | |
| VSP | VMware vSphere | |
| OVT | oVirt / RHV | Legacy, destroy-only |

## Legend

| Code | Meaning |
|------|---------|
| AT | Automatically Tested - covered by unit tests or CI e2e jobs |
| MT | Manually Tested - manual test procedures documented |
| AT/MT | Both automatic and manual tests exist |
| PT | Partially Tested - some paths tested, gaps remain |
| NT | Not Tested - feature is implemented but no test coverage found |
| NA | Not Applicable - feature not implemented for this platform |

## Feature Matrix

| Feature | AWS | AZ | GCP | GCD | BM | IBM | NTX | OST | PVC | PVS | VSP | OVT |
|---------|-----|-----|-----|-----|-----|-----|-----|-----|-----|-----|-----|-----|
| **Core Installation** | | | | | | | | | | | | |
| IPI Install | AT[1] | AT[2] | AT[3] | AT[4] | AT[5] | AT[6] | AT[6] | AT[7] | AT[6] | AT[6] | AT[8] | NA |
| UPI Install | AT[9] | AT[9] | AT[9] | NA | NA | NA | NA | AT[9] | NA | NA | AT[9] | NT[10] |
| Agent-based Install | AT[11] | AT[11] | AT[11] | NA | AT[11] | NA | AT[11] | AT[11] | NA | NA | AT[11] | NA |
| Single Node (SNO) | AT | NT | NT | NA | AT | NA | NT | NT | NA | NA | NT | NA |
| Cluster Destroy | NT[12] | NT | AT[13] | NT[14] | NT | NT | NT | NT | NT | AT | AT | NT |
| **Networking** | | | | | | | | | | | | |
| Existing VPC / VNet | AT | AT | AT | NT | NA | AT | AT | AT | NT | AT | AT | NA |
| Shared VPC / XPN | AT | AT | AT[15] | NA | NA | NA | NA | NA | NA | NA | NA | NA |
| Private Cluster | AT | AT | AT | AT[16] | NA | NT | NA | NT | NA | NT | NA | NA |
| Proxy Support | NT[17] | NT[17] | NT[17] | NT[17] | NT[17] | NT[17] | NT[17] | AT | NT[17] | NT[17] | NT[17] | NA |
| Dual Stack (IPv4+IPv6) | AT | AT | NA | NA | AT | NA | AT | AT | NT | NA | AT | NA |
| User-Provisioned DNS | AT | AT | AT | NA | NA | NA | NA | NA | NA | NA | NA | NA |
| **Compute** | | | | | | | | | | | | |
| Default Machine Types | AT | AT | AT | AT/MT[18] | NA | AT | NT | AT | NT | AT | NA | AT |
| Custom Disk Types | AT | AT | AT | AT[19] | NA | NT | NA | NA | NA | NA | NT | NA |
| Confidential Compute | AT | AT | AT | NA | NA | NA | NA | NA | NA | NA | NA | NA |
| Secure Boot / UEFI | NA | AT | NT[20] | NA | AT | NA | AT | NA | NA | NA | NA | NA |
| Dedicated Hosts | AT[21] | NA | NA | NA | NA | AT | NA | NA | NA | NA | NA | NA |
| Custom OS Image | NT | NT | AT[22] | AT/MT[22] | AT | NA | NT | NT | NT | NT | NT | NA |
| **Storage and Security** | | | | | | | | | | | | |
| KMS / Disk Encryption | AT | AT | AT | NA | NA | AT | NA | NA | NA | NA | NA | NA |
| User Tags / Labels | AT | AT | AT | NA | NA | NA | NA | NA | NA | NA | NT[23] | NA |
| Service Endpoints / PSC | AT | NA | AT[24] | NA[25] | NA | AT | NA | NA | NA | AT | NA | NA |
| FIPS Mode | NT[26] | NT[26] | NT[26] | NT[26] | NT[26] | NT[26] | NT[26] | NT[26] | NT[26] | NT[26] | NT[26] | NA |
| **Infrastructure** | | | | | | | | | | | | |
| CAPI Provisioning | AT | AT | AT | AT | NA[27] | AT | AT | AT | AT[28] | AT | AT | NA |
| Load Balancer Config | NT | NA | NA | NA | NT | NA | NT | NT | NT | NA | NT | NT |
| Failure Domains / Zones | AT | AT | AT | NA | NA | NA | AT | NA | NA | NA | AT | NA |
| Feature Gates | AT | NT | NT | NA | NT | NA | NT | NT | NA | NA | NT | NA |
| **Sovereign / Gov Cloud** | | | | | | | | | | | | |
| Sovereign Cloud Support | AT[29] | NT[30] | AT[31] | AT/MT[32] | NA | NA | NA | NA | NA | NA | NA | NA |
| **Test Documentation** | | | | | | | | | | | | |
| Test Plan | No | No | No | Yes[33] | No | No | No | No | No | No | No | No |
| Test Cases | No | No | No | Yes[33] | No | No | No | No | No | No | No | No |

## Footnotes

1. 29 e2e CI jobs in openshift/release
2. 15 e2e CI jobs including Azure Stack
3. 16 e2e CI jobs covering IPI, UPI, shared VPC, KMS, custom DNS, and more
4. 1 e2e CI job: `e2e-gcd-ovn-private-techpreview`
5. 10 e2e CI jobs for bare metal IPI
6. 1 e2e CI job
7. 8 e2e CI jobs including dual-stack, proxy, and external LB
8. 16 e2e CI jobs including multi-vCenter, static IPs, and host groups
9. UPI templates in `upi/<platform>/` with corresponding e2e UPI CI jobs
10. UPI templates exist in `upi/ovirt/` but no CI coverage found
11. 11 agent-based CI jobs covering compact, HA, SNO, dual-stack, and PXE
    configurations across multiple platforms
12. AWS destroy has 9 source files in `pkg/destroy/aws/` with 0 test files
13. GCP destroy has unit tests: `gcp_test.go`, `policybinding_test.go`,
    `quota_test.go`
14. GCD shares GCP destroy code in `pkg/destroy/gcp/`; no GCD-specific
    destroy tests
15. GCP shared VPC via `NetworkProjectID`; dedicated CI jobs
    (`e2e-gcp-ovn-xpn`, `e2e-gcp-xpn-dedicated-dns-project`, etc.)
16. GCD is private-only (no public DNS zones in sovereign cloud); the e2e
    job uses `Publish: Internal`
17. Proxy is configured at the platform-agnostic `installconfig.Proxy`
    level; no platform-specific proxy validation tests found. OpenStack is
    the exception with a dedicated `e2e-openstack-proxy` CI job.
18. GCD defaults to C3 machine types (only series available in sovereign
    regions); covered by unit tests in `machinepools_test.go` and manual
    verification scenario in test docs
19. GCD defaults to `hyperdisk-balanced` (only disk type available in
    sovereign regions); covered by `machinepools_test.go`
20. GCP has a `SecureBoot` field in machine pool but no validation tests
    for it
21. AWS `DedicatedHost` placement behind `AWSDedicatedHosts` feature gate;
    IBM Cloud also supports dedicated hosts with `DedicatedHosts` field
22. GCP/GCD `OSImage` (name + project) is required for sovereign cloud and
    optional for public GCP; validated in `machinepool_test.go`. GCD test
    docs include manual verification of OS image on deployed nodes.
23. vSphere supports `tagIDs` in failure domain topology but has no
    dedicated tag validation tests
24. GCP uses `PSCEndpoint` (Private Service Connect); IBM Cloud and Power VS
    use `ServiceEndpoints` for API endpoint overrides
25. GCD forbids PSC endpoint overrides; this constraint is validated in unit
    tests (reflected in the Sovereign Cloud Support row)
26. FIPS is a platform-agnostic boolean in install-config; no
    platform-specific FIPS validation tests exist in the installer
27. Bare Metal uses custom infrastructure provisioning (not Cluster API)
28. PowerVC shares the OpenStack CAPI provider
29. AWS GovCloud regions are supported; European Sovereign Cloud is behind
    the `AWSEuropeanSovereignCloudInstall` feature gate with test coverage
    in `featuregates_test.go`
30. Azure supports `AzureUSGovernmentCloud`, `AzureChinaCloud`, and
    `AzureGermanCloud` cloud names but has no dedicated sovereign/gov cloud
    test coverage
31. GCP sovereign cloud detection (`GetCloudEnvironment`) tested in
    `pkg/types/gcp/platform_test.go`
32. GCD has full sovereign cloud test coverage: unit tests for detection,
    defaults, validation, and config generation, plus manual test
    procedures documented in `test/docs/gcd/cases/gcd_sovereign_install.md`
33. Test plan at `test/docs/gcd/plans/gcd-feature.md` and test cases at
    `test/docs/gcd/cases/gcd_sovereign_install.md` covering 17 scenarios
    (13 automated, 4 manual)

## How to Update This Matrix

Update this matrix when:

- A new platform is added to the installer
- A feature gains or loses test coverage
- New test documentation is added under `test/docs/`
- CI jobs are added, removed, or renamed in openshift/release

To determine the correct status code for a cell:

1. Check for unit tests: look for `_test.go` files in
   `pkg/types/<platform>/validation/`, `pkg/asset/installconfig/<platform>/`,
   `pkg/destroy/<platform>/`, and `pkg/infrastructure/<platform>/`
2. Check for e2e CI jobs: search the openshift/release repo for job
   definitions referencing the platform
3. Check for manual test docs: look for files under
   `test/docs/<platform>/cases/` with `MANUAL_TESTS_START` markers
4. If the feature has code but no tests of any kind, use NT
5. If the feature is not implemented for the platform, use NA

When adding a footnote, use the next available number and add the
explanation to the Footnotes section. Keep footnotes factual and brief.

## References

- [Test documentation structure and conventions](README.md)
- [GCD feature test plan](gcd/plans/gcd-feature.md)
- [GCD sovereign install test cases](gcd/cases/gcd_sovereign_install.md)
- CI job definitions: openshift/release repo, branch config for
  `openshift-installer-main`
