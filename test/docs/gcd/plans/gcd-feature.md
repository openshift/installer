# Feature Test Plan: OpenShift on Google Cloud Dedicated (GCD)

## 1. Introduction

### 1.1. Overview

This document outlines the testing strategy, objectives, and procedures for
OpenShift IPI (Installer-Provisioned Infrastructure) support on Google Cloud
Dedicated (GCD), also known as Trusted Partner Cloud sovereign regions.

GCD is a fully isolated sovereign cloud environment with no network path to
public Google Cloud. It uses region-specific API endpoints, domain-scoped
project IDs, and a restricted subset of GCP services. The installer must detect
sovereign cloud environments, apply correct defaults, enforce constraints, and
propagate configuration to the cluster components.

### 1.2. Scope

**In Scope**

- Install-config validation for sovereign cloud environments
- Sovereign cloud detection via project ID prefix (`eu0:`)
- Default machine type selection (C3 series for sovereign, N2 for public GCP)
- Default disk type selection (Hyperdisk Balanced for sovereign, PD-SSD for public GCP)
- OS image requirement enforcement for sovereign cloud
- PSC endpoint override restriction for sovereign cloud
- Service account email formatting for domain-scoped project IDs
- Universe domain propagation to CAPG controller
- Cloud provider config with nil tokenURL for non-default universe domains
- GCD credential handling (self-signed JWT via `WithCredentialsJSON`)
- Private cluster deployment on GCD (internal publish only)
- E2E IPI installation on GCD sovereign region (`u-germany-northeast1`)

**Out of Scope**

- UPI (User-Provisioned Infrastructure) on GCD
- S3NS (Cloud de Confiance) region testing (same architecture, different identifiers)
- Post-install Day 2 operations and operator testing
- Upgrade testing from public GCP to GCD
- Multi-region or cross-region scenarios (GCD is single-region)

### 1.3. Key Features

- **OCPSTRAT-3006** - Google Cloud Dedicated support in OpenShift
- Sovereign cloud auto-detection from `eu0:` project ID prefix
- Constrained defaults for machine types, disk types, and OS images
- Private DNS-only deployment (GCD does not support public DNS zones)
- Universe domain-aware credential handling

### 1.4. References

- Strategy: [OCPSTRAT-3006](https://redhat.atlassian.net/browse/OCPSTRAT-3006)
- Installer PR: [openshift/installer#10630](https://github.com/openshift/installer/pull/10630)
- CI job PR: [openshift/release#81238](https://github.com/openshift/release/pull/81238)

## 2. Testing Strategy

### 2.1. Schedule

| Milestone | Date / Sprint | Notes |
|-----------|---------------|-------|
| Unit test coverage complete | TBD | All sovereign cloud code paths covered |
| CI job green | TBD | `e2e-gcd-ovn-private-techpreview` passing |
| Feature complete | TBD | GCD IPI install succeeds end-to-end |

### 2.2. Test Types

- **Unit Tests:** Validate sovereign cloud detection, defaults, validation
  rules, and configuration generation in isolation
- **Integration Tests:** Verify install-config parsing, asset generation, and
  credential handling for sovereign cloud configurations
- **E2E Tests:** Full IPI installation on GCD sovereign region with cluster
  health verification
- **Negative Tests:** Confirm that unsupported configurations (e.g., PSC
  endpoints, public DNS, wrong machine types) are correctly rejected

### 2.3. Test Environments

- **Unit/Integration:** Local Go test environment, no cloud credentials needed
- **E2E CI:** Prow job `e2e-gcd-ovn-private-techpreview` using `gcd` cluster
  profile with:
  - Region: `u-germany-northeast1`
  - Base domain: `ci.gcd.devcluster.openshift.com`
  - Feature set: `TechPreviewNoUpgrade`
  - Publish: `Internal` (private cluster)
  - Machine type: `c3-standard-4`
  - OS image: `rhcos10` (uploaded to GCD project)
  - Bastion: Fedora CoreOS with GVNIC on C3

## 3. Test Areas and Test Cases

### 3.1. Key Test Activities

| Activity | Type | Est. Duration |
|----------|------|---------------|
| Sovereign cloud detection unit tests | Unit | < 1 hour |
| Default values unit tests | Unit | < 1 hour |
| Validation unit tests | Unit | < 1 hour |
| Service account formatting unit tests | Unit | < 30 min |
| Install-config validation integration tests | Integration | < 1 hour |
| E2E IPI install on GCD | E2E | ~2 hours |

### 3.2. Test Cases

Test case documentation is maintained alongside this plan in the
[`cases/`](../cases/) directory, with references to the corresponding Go test
files.

| Test Case | Doc | Test File(s) |
|-----------|-----|--------------|
| GCD Sovereign Cloud Install | [gcd_sovereign_install.md](../cases/gcd_sovereign_install.md) | Multiple (see doc) |

### 3.3. Unit Test Coverage Map

The following table maps each GCD-specific behavior to the source code and
test file that covers it.

| Behavior | Source | Test File |
|----------|--------|-----------|
| Sovereign cloud detection from project ID | `pkg/types/gcp/platform.go` `GetCloudEnvironment()` | `pkg/types/gcp/platform_test.go` |
| Non-default universe domain check | `pkg/types/gcp/platform.go` `IsNonDefaultUniverseDomain()` | `pkg/types/gcp/platform_test.go` |
| Service account email for domain-scoped projects | `pkg/types/gcp/platform.go` `GetDefaultServiceAccount()` | `pkg/types/gcp/platform_test.go` |
| Default machine type (C3 for sovereign) | `pkg/asset/installconfig/gcp/validation.go` `DefaultInstanceTypeForArchAndProjectID()` | `pkg/asset/installconfig/gcp/validation_test.go` |
| Default disk type (Hyperdisk Balanced for sovereign) | `pkg/types/gcp/machinepools.go` `DefaultDiskTypeForInstanceAndProjectID()` | `pkg/types/gcp/machinepools_test.go` |
| OS image required for sovereign cloud | `pkg/types/gcp/validation/machinepool.go` `ValidateOSImageForSovereignCloud()` | `pkg/types/gcp/validation/machinepool_test.go` |
| PSC endpoint forbidden for sovereign cloud | `pkg/asset/installconfig/gcp/validation.go` `validateServiceEndpointOverride()` | `pkg/asset/installconfig/gcp/validation_test.go` |
| Universe domain credential options | `pkg/asset/installconfig/gcp/services.go` `CredentialOptions()` | `pkg/asset/installconfig/gcp/services_test.go` |
| Cloud provider config nil tokenURL | `pkg/asset/manifests/cloudproviderconfig.go` | `pkg/asset/manifests/cloudproviderconfig_test.go` |
| CAPG controller universe domain env var | `pkg/clusterapi/system.go` | `pkg/clusterapi/system_test.go` |

## 4. Test Details

### 4.1. Sovereign Cloud Detection

The installer identifies GCD environments by inspecting the project ID prefix.
Project IDs with the `eu0:` prefix are classified as sovereign cloud. This
detection drives all downstream defaults and validation.

**Key invariant:** Organization-scoped public GCP projects (e.g.,
`myorg:my-project`) must NOT be classified as sovereign. Only the known prefix
`eu0` triggers sovereign behavior.

### 4.2. Default Values for Sovereign Cloud

When a sovereign cloud environment is detected:

| Setting | Public GCP | Sovereign Cloud (GCD) |
|---------|------------|-----------------------|
| Default instance type (x86) | `n2-standard-4` | `c3-standard-4` |
| Default instance type (ARM64) | `t2a-standard-4` | `c3-standard-4` |
| Default disk type (C3) | `hyperdisk-balanced` | `hyperdisk-balanced` |
| OS image | Optional | Required (name + project) |
| PSC endpoint override | Allowed | Forbidden |
| Token URL in cloud-provider config | Standard OAuth2 URL | `"nil"` |
| CAPG `GOOGLE_CLOUD_UNIVERSE_DOMAIN` | Not set | Set from credentials |

### 4.3. E2E Test Execution

The E2E test job `e2e-gcd-ovn-private-techpreview` performs a full IPI
installation on GCD. The workflow:

1. **Credential validation** (`ipi-conf-gcd-creds`) - Validates GCD service
   account key contains `universe_domain` field
2. **VPC provisioning** (`gcp-provision-vpc`) - Creates VPC, subnets, Cloud NAT
   in `u-germany-northeast1`
3. **Bastion provisioning** (`gcp-provision-bastionhost`) - Creates C3 bastion
   with Fedora CoreOS (GVNIC-enabled image)
4. **Proxy setup** - Configures HTTP proxy on bastion for private cluster access
5. **Install-config generation** (`ipi-conf-gcd`) - Generates install-config
   with GCD-specific settings
6. **Cluster installation** - Runs `openshift-install create cluster`
7. **E2E tests** - Runs conformance tests against the deployed cluster
8. **Deprovision** - Destroys cluster, bastion, and VPC

### 4.4. GCD-Specific Requirements in E2E

| Requirement | How it is handled |
|-------------|-------------------|
| Machine types: C3/M3 only | Default `c3-standard-4`, bastion also C3 |
| Disks: Hyperdisk Balanced only | Default from `DefaultDiskTypeForInstanceAndProjectID()` |
| DNS: private zones only | `PUBLISH=Internal` in job config |
| Images: must exist in GCD project | `DEFAULT_MACHINE_OSIMAGE=rhcos10` (pre-uploaded) |
| Bastion needs GVNIC | Fedora CoreOS image uploaded with `--guest-os-features=GVNIC` |
| Project IDs use `eu0:` prefix | Cluster profile provides `eu0:` project |
| Single region | `u-germany-northeast1` |
| Universe domain | Read from `gce.json` credentials, set via `gcloud config` and env var |

## 5. Risks

| Risk | Impact | Mitigation |
|------|--------|------------|
| GCD environment availability | E2E tests cannot run if GCD region is down | Unit tests cover all logic without cloud access |
| RHCOS image not uploaded to GCD | Install fails on missing boot image | Pre-upload RHCOS to GCD project, track in CI config |
| OAuth2 token endpoint unavailable in GCD | Credential refresh fails | Use `WithCredentialsJSON` for self-signed JWT (PR #10630) |
| New `eu0:` prefix not in known list | Future sovereign prefixes not detected | `sovereignCloudProjectPrefixes` slice is extensible |
| C3 machine quota in GCD | Insufficient quota blocks install | Monitor quota, request increases proactively |

## 6. Exit Criteria

The GCD feature will be considered tested when:

- All unit tests for sovereign cloud detection, defaults, and validation pass
- No open regressions in existing GCP test suites
- E2E CI job `e2e-gcd-ovn-private-techpreview` passes (cluster installs and
  conformance tests succeed)
- All negative validation cases (PSC endpoint, missing OS image, wrong machine
  types) are covered by unit tests
- Service account email formatting is verified for domain-scoped project IDs
