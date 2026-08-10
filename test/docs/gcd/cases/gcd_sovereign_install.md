# Test Case: GCD Sovereign Cloud Installation

## Metadata

| Field | Value |
|-------|-------|
| **Feature** | [OCPSTRAT-3006](https://redhat.atlassian.net/browse/OCPSTRAT-3006) - Google Cloud Dedicated |
| **Component** | Installer / GCP platform |
| **Test type** | Unit, Integration, E2E |
| **Assignee** | TBD |

## Description

These test cases cover the installer's behavior when targeting a Google Cloud
Dedicated (GCD) sovereign cloud environment. GCD is detected by two conditions:
a domain-scoped project ID (containing `:`) and a sovereign region (prefixed
with `u-`). Both must be present. It differs from public GCP in machine types,
disk types, DNS, OS image requirements, and credential handling.

The cases are organized into groups: sovereign cloud detection, default values,
validation rules, configuration generation, and end-to-end installation.

## Scenarios

| # | Scenario | Status |
|---|----------|--------|
| 1 | [Sovereign cloud detection from project ID and region](#1-sovereign-cloud-detection-from-project-id-and-region) | Automated (unit test) |
| 2 | [Domain-scoped project or sovereign region alone is not sovereign](#2-domain-scoped-project-or-sovereign-region-alone-is-not-sovereign) | Automated (unit test) |
| 3 | [Non-default universe domain detection](#3-non-default-universe-domain-detection) | Automated (unit test) |
| 4 | [Default instance type for sovereign cloud](#4-default-instance-type-for-sovereign-cloud) | Automated (unit test) |
| 5 | [Default disk type for sovereign cloud](#5-default-disk-type-for-sovereign-cloud) | Automated (unit test) |
| 6 | [OS image required for sovereign cloud](#6-os-image-required-for-sovereign-cloud) | Automated (unit test) |
| 7 | [OS image with missing name or project](#7-os-image-with-missing-name-or-project) | Automated (unit test) |
| 8 | [PSC endpoint forbidden for sovereign cloud](#8-psc-endpoint-forbidden-for-sovereign-cloud) | Automated (unit test) |
| 9 | [Service account email for domain-scoped project](#9-service-account-email-for-domain-scoped-project) | Automated (unit test) |
| 10 | [Cloud provider config uses nil tokenURL](#10-cloud-provider-config-uses-nil-tokenurl) | Automated (unit test) |
| 11 | [CAPG controller receives universe domain](#11-capg-controller-receives-universe-domain) | Automated (unit test) |
| 12 | [Credential options for sovereign cloud](#12-credential-options-for-sovereign-cloud) | Automated (unit test) |
| 13 | [E2E: IPI install on GCD sovereign region](#13-e2e-ipi-install-on-gcd-sovereign-region) | Automated (Prow CI) |
| 14 | [Verify GCD cluster nodes use correct machine type and disk](#14-verify-gcd-cluster-nodes-use-correct-machine-type-and-disk) | Manual |
| 15 | [Verify GCD cluster uses private DNS only](#15-verify-gcd-cluster-uses-private-dns-only) | Manual |
| 16 | [Verify GCD cluster destroy removes all resources](#16-verify-gcd-cluster-destroy-removes-all-resources) | Manual |
| 17 | [Verify install-config validation rejects public publish on GCD](#17-verify-install-config-validation-rejects-public-publish-on-gcd) | Manual |

---

## 1. Sovereign cloud detection from project ID and region

- **Test file:** `pkg/types/gcp/platform_test.go`
- **Function under test:** `GetCloudEnvironment(projectID, region)`
- **Automation status:** Automated (unit test)

### Description

Verifies that a sovereign cloud environment is correctly identified when both
conditions are met: the project ID is domain-scoped (contains `:`) and the
region has the sovereign prefix `u-`. Neither condition alone is sufficient.

### Execution

```go
env := gcp.GetCloudEnvironment("eu0:my-project", "u-germany-northeast1")
// Expected: "sovereign"

env = gcp.GetCloudEnvironment("s3ns:my-project", "u-france-east1")
// Expected: "sovereign"
```

### Expected Result

- `GetCloudEnvironment("eu0:my-project", "u-germany-northeast1")` returns `"sovereign"`
- `GetCloudEnvironment("s3ns:my-project", "u-france-east1")` returns `"sovereign"`

### Pass/Fail Criteria

| Project ID | Region | Expected Output |
|-----------|--------|-----------------|
| `"eu0:my-project"` | `"u-germany-northeast1"` | `"sovereign"` |
| `"s3ns:my-project"` | `"u-france-east1"` | `"sovereign"` |
| `"my-project"` | `"us-central1"` | `""` (empty - public GCP) |

---

## 2. Domain-scoped project or sovereign region alone is not sovereign

- **Test file:** `pkg/types/gcp/platform_test.go`
- **Function under test:** `GetCloudEnvironment(projectID, region)`
- **Automation status:** Automated (unit test)

### Description

A domain-scoped project ID without a sovereign region, or a sovereign region
without a domain-scoped project ID, must not be classified as sovereign. Both
conditions are required together. This prevents false positives from
organization-scoped public GCP projects or mistyped regions.

### Execution

```go
// Domain-scoped project but public GCP region - not sovereign
env := gcp.GetCloudEnvironment("other-prefix:my-project", "us-central1")
// Expected: "" (not sovereign)

// Sovereign region but standard project ID - not sovereign
env = gcp.GetCloudEnvironment("my-project", "u-germany-northeast1")
// Expected: "" (not sovereign)
```

### Expected Result

- `GetCloudEnvironment("other-prefix:my-project", "us-central1")` returns `""`
- `GetCloudEnvironment("my-project", "u-germany-northeast1")` returns `""`
- Only the combination of both triggers sovereign detection

### Pass/Fail Criteria

| Project ID | Region | Expected Output |
|-----------|--------|-----------------|
| `"other-prefix:my-project"` | `"us-central1"` | `""` (not sovereign - no `u-` region) |
| `"my-project"` | `"u-germany-northeast1"` | `""` (not sovereign - no `:` in project) |
| `"eu0:my-project"` | `"u-germany-northeast1"` | `"sovereign"` (both conditions met) |

---

## 3. Non-default universe domain detection

- **Test file:** `pkg/types/gcp/platform_test.go`
- **Function under test:** `IsNonDefaultUniverseDomain()`
- **Automation status:** Automated (unit test)

### Description

Verifies that universe domains other than `googleapis.com` are correctly
identified as non-default, which triggers sovereign cloud credential handling.

### Execution

```go
result := gcp.IsNonDefaultUniverseDomain("apis-berlin-build0.goog")
// Expected: true
```

### Pass/Fail Criteria

| Input | Expected Output |
|-------|-----------------|
| `"apis-berlin-build0.goog"` | `true` |
| `"s3nsapis.fr"` | `true` |
| `"googleapis.com"` | `false` |
| `""` | `false` |

---

## 4. Default instance type for sovereign cloud

- **Test file:** `pkg/asset/installconfig/gcp/validation_test.go`
- **Function under test:** `DefaultInstanceTypeForArchAndProjectID(arch, projectID, region)`
- **Automation status:** Automated (unit test)

### Description

GCD sovereign regions only have C3, M3, and A3 Edge machine series available.
The installer must default to `c3-standard-4` for sovereign cloud regardless
of CPU architecture, instead of the public GCP defaults (`n2-standard-4` for
x86, `t2a-standard-4` for ARM64).

### Execution

```go
// Sovereign cloud - always c3-standard-4
instanceType := DefaultInstanceTypeForArchAndProjectID(types.ArchitectureAMD64, "eu0:my-project", "u-germany-northeast1")
// Expected: "c3-standard-4"

instanceType = DefaultInstanceTypeForArchAndProjectID(types.ArchitectureARM64, "eu0:my-project", "u-germany-northeast1")
// Expected: "c3-standard-4"

// Public GCP - architecture-specific
instanceType = DefaultInstanceTypeForArchAndProjectID(types.ArchitectureAMD64, "my-project", "us-central1")
// Expected: "n2-standard-4"
```

### Pass/Fail Criteria

| Architecture | Project ID | Region | Expected Instance Type |
|-------------|------------|--------|----------------------|
| x86_64 | `eu0:my-project` | `u-germany-northeast1` | `c3-standard-4` |
| ARM64 | `eu0:my-project` | `u-germany-northeast1` | `c3-standard-4` |
| x86_64 | `my-project` | `us-central1` | `n2-standard-4` |
| ARM64 | `my-project` | `us-central1` | `t2a-standard-4` |

---

## 5. Default disk type for sovereign cloud

- **Test file:** `pkg/types/gcp/machinepools_test.go`
- **Function under test:** `DefaultDiskTypeForInstanceAndProjectID(instanceType, projectID, region)`
- **Automation status:** Automated (unit test)

### Description

GCD sovereign regions only support Hyperdisk Balanced. The installer must
prefer `hyperdisk-balanced` for sovereign cloud projects, while preferring
`pd-ssd` for public GCP (when the instance type supports both).

### Execution

```go
// Sovereign cloud
diskType := gcp.DefaultDiskTypeForInstanceAndProjectID("c3-standard-4", "eu0:my-project", "u-germany-northeast1")
// Expected: "hyperdisk-balanced"

// Public GCP
diskType = gcp.DefaultDiskTypeForInstanceAndProjectID("c3-standard-4", "my-project", "us-central1")
// Expected: "hyperdisk-balanced" (C3 only supports hyperdisk-balanced)

diskType = gcp.DefaultDiskTypeForInstanceAndProjectID("n2-standard-4", "my-project", "us-central1")
// Expected: "pd-ssd"
```

### Pass/Fail Criteria

| Instance Type | Project ID | Region | Expected Disk Type |
|--------------|------------|--------|-------------------|
| `c3-standard-4` | `eu0:my-project` | `u-germany-northeast1` | `hyperdisk-balanced` |
| `n2-standard-4` | `eu0:my-project` | `u-germany-northeast1` | `hyperdisk-balanced` (sovereign prefers it) |
| `n2-standard-4` | `my-project` | `us-central1` | `pd-ssd` |

---

## 6. OS image required for sovereign cloud

- **Test file:** `pkg/types/gcp/validation/machinepool_test.go`
- **Function under test:** `ValidateOSImageForSovereignCloud()`
- **Automation status:** Automated (unit test)

### Description

GCD does not have access to public GCP image projects. The installer must
require an explicit OS image (both name and project) when deploying to a
sovereign cloud environment. Without it, the cluster nodes cannot boot.

### Execution

```go
// Sovereign cloud without OS image - must fail
platform := &gcp.Platform{ProjectID: "eu0:my-project", Region: "u-germany-northeast1"}
pool := &gcp.MachinePool{} // no OSImage
errs := ValidateOSImageForSovereignCloud(platform, pool, field.NewPath("test"))
// Expected: error "must specify an OS image for sovereign cloud environments"
```

### Pass/Fail Criteria

| Scenario | Expected Result |
|----------|----------------|
| Sovereign cloud, OS image with name + project | No errors |
| Sovereign cloud, nil pool | Error: `osImage` required |
| Sovereign cloud, nil OS image on pool | Error: `osImage` required |
| Sovereign cloud, OS image missing name | Error: `osImage.name` required |
| Sovereign cloud, OS image missing project | Error: `osImage.project` required |
| Public GCP, no OS image | No errors (optional) |

---

## 7. OS image with missing name or project

- **Test file:** `pkg/types/gcp/validation/machinepool_test.go`
- **Function under test:** `ValidateOSImageForSovereignCloud()`
- **Automation status:** Automated (unit test)

### Description

When a sovereign cloud install-config provides an OS image but omits either the
image name or image project, the validation must report a specific error for the
missing field.

### Execution

```go
// Missing name
platform := &gcp.Platform{ProjectID: "eu0:my-project", Region: "u-germany-northeast1"}
pool := &gcp.MachinePool{OSImage: &gcp.OSImage{Project: "eu0-system:rhcos"}}
errs := ValidateOSImageForSovereignCloud(platform, pool, field.NewPath("test"))
// Expected: error on "test.osImage.name"

// Missing project
pool = &gcp.MachinePool{OSImage: &gcp.OSImage{Name: "rhcos-421"}}
errs = ValidateOSImageForSovereignCloud(platform, pool, field.NewPath("test"))
// Expected: error on "test.osImage.project"
```

### Pass/Fail Criteria

| OS Image State | Expected Error Field |
|---------------|---------------------|
| Name present, project missing | `osImage.project: Required` |
| Project present, name missing | `osImage.name: Required` |
| Both present | No errors |
| Both missing (nil OSImage) | `osImage: Required` |

---

## 8. PSC endpoint forbidden for sovereign cloud

- **Test file:** `pkg/asset/installconfig/gcp/validation_test.go`
- **Function under test:** `validateServiceEndpointOverride()`
- **Automation status:** Automated (unit test)

### Description

Private Service Connect (PSC) endpoint overrides are not supported in sovereign
cloud environments. GCD uses its own API endpoint domain
(`apis-berlin-build0.goog`) and custom endpoints would conflict with the
sovereign infrastructure.

### Execution

Provide an install-config with a sovereign project ID and a PSC endpoint:

```yaml
platform:
  gcp:
    projectID: "eu0:my-project"
    region: u-germany-northeast1
    endpoint:
      name: my-psc-endpoint
```

### Expected Result

- Validation returns: `platform.gcp.endpoint.name: Forbidden: endpoint overrides are not supported in sovereign clouds`

### Pass/Fail Criteria

| Scenario | Expected Result |
|----------|----------------|
| Sovereign project + PSC endpoint | `Forbidden` error |
| Sovereign project, no endpoint | No errors |
| Public GCP project + valid PSC endpoint | No errors (allowed) |

---

## 9. Service account email for domain-scoped project

- **Test file:** `pkg/types/gcp/platform_test.go`
- **Function under test:** `GetDefaultServiceAccount()`
- **Automation status:** Automated (unit test)

### Description

Domain-scoped project IDs use a reversed format in service account email
addresses. For a project `eu0:my-project`, the SA email domain becomes
`my-project.eu0.iam.gserviceaccount.com` (the prefix `eu0` and project name
`my-project` are reversed and joined with a dot).

### Execution

```go
platform := &gcp.Platform{ProjectID: "eu0:my-project"}
sa := gcp.GetDefaultServiceAccount(platform, "test-cluster", "master")
// Expected: "test-cluster-m@my-project.eu0.iam.gserviceaccount.com"
```

### Pass/Fail Criteria

| Project ID | Cluster ID | Role | Expected SA Email |
|-----------|-----------|------|-------------------|
| `eu0:my-project` | `test-cluster` | `master` | `test-cluster-m@my-project.eu0.iam.gserviceaccount.com` |
| `eu0:my-project` | `test-cluster` | `worker` | `test-cluster-w@my-project.eu0.iam.gserviceaccount.com` |
| `my-project` | `test-cluster` | `master` | `test-cluster-m@my-project.iam.gserviceaccount.com` |

---

## 10. Cloud provider config uses nil tokenURL

- **Test file:** `pkg/asset/manifests/cloudproviderconfig_test.go`
- **Function under test:** Cloud provider config generation
- **Automation status:** Automated (unit test)

### Description

When the universe domain is not the default `googleapis.com`, the cloud provider
configuration must set `tokenURL` to the literal string `"nil"`. This disables
the standard OAuth2 token endpoint, which is not available in GCD sovereign
regions. Instead, the workload uses self-signed JWTs.

### Expected Result

- Cloud provider config generated for a sovereign cloud install-config contains
  `tokenURL = "nil"`
- Cloud provider config for public GCP does not set `tokenURL` to `"nil"`

---

## 11. CAPG controller receives universe domain

- **Test file:** `pkg/clusterapi/system_test.go`
- **Function under test:** CAPG controller environment setup
- **Automation status:** Automated (unit test)

### Description

When deploying on a non-default universe domain, the Cluster API Provider for
GCP (CAPG) controller must receive the `GOOGLE_CLOUD_UNIVERSE_DOMAIN`
environment variable so it can reach the correct sovereign cloud API endpoints.

### Expected Result

- When universe domain is `apis-berlin-build0.goog`, the CAPG controller
  process has `GOOGLE_CLOUD_UNIVERSE_DOMAIN=apis-berlin-build0.goog` set
- When universe domain is `googleapis.com` or empty, the environment variable
  is not set

---

## 12. Credential options for sovereign cloud

- **Test file:** `pkg/asset/installconfig/gcp/services_test.go`
- **Function under test:** `CredentialOptions()`
- **Automation status:** Automated (unit test)

### Description

GCD service account keys contain a `universe_domain` field (e.g.,
`apis-berlin-build0.goog`). The installer must pass this to the GCP client
libraries and use `WithCredentialsJSON` for self-signed JWT authentication,
since the standard OAuth2 token endpoint (`oauth2.googleapis.com`) is not
reachable from GCD.

### Expected Result

- `CredentialOptions()` extracts the universe domain from the SA key JSON
- Client options include `WithCredentialsJSON` (not `WithTokenSource`)
- Client options include `WithUniverseDomain` set to the credential's domain

---

## 13. E2E: IPI install on GCD sovereign region

- **CI job:** `e2e-gcd-ovn-private-techpreview`
- **Workflow:** `openshift-e2e-gcd`
- **Automation status:** Automated (Prow CI)

### Description

Full end-to-end IPI installation of an OpenShift cluster on GCD. This is the
highest-level validation that all sovereign cloud code paths work together.

### Prerequisites

- GCD service account key in vault (`cluster-secrets-gcd`)
- DNS zone `ci.gcd.devcluster.openshift.com` created in GCD project
- RHCOS image uploaded to GCD project with GVNIC support
- Sufficient C3 machine quota in `u-germany-northeast1`

### Test Steps and Expected Results

| Step | Expected Result |
|:-----|:----------------|
| Validate GCD credentials | `gce.json` contains `universe_domain` field |
| Provision VPC in `u-germany-northeast1` | VPC, subnets, Cloud NAT, firewall rules created |
| Provision C3 bastion with GVNIC image | Bastion accessible, proxy configured |
| Generate install-config | Config uses domain-scoped project, `u-germany-northeast1`, `c3-standard-4`, `hyperdisk-balanced`, RHCOS OS image, `Internal` publish |
| Run `openshift-install create cluster` | Cluster installs successfully on GCD |
| Run conformance tests | Core OpenShift functionality works on GCD |
| Deprovision cluster | All GCD resources destroyed cleanly |
| Deprovision bastion and VPC | Infrastructure cleaned up |

### GCD-Specific Configuration

```yaml
platform:
  gcp:
    projectID: "eu0:<gcd-project>"
    region: u-germany-northeast1
controlPlane:
  platform:
    gcp:
      type: c3-standard-4
      osDisk:
        diskType: hyperdisk-balanced
      osImage:
        name: rhcos10
        project: "eu0:<gcd-project>"
compute:
  - platform:
      gcp:
        type: c3-standard-4
        osDisk:
          diskType: hyperdisk-balanced
        osImage:
          name: rhcos10
          project: "eu0:<gcd-project>"
publish: Internal
featureSet: TechPreviewNoUpgrade
```

### Pass/Fail Criteria

| Check | Pass Condition |
|-------|---------------|
| Cluster install completes | `openshift-install` exits 0 |
| ClusterVersion available | `Available=True`, `Progressing=False` |
| All ClusterOperators healthy | No degraded operators |
| Nodes ready | All control plane and compute nodes `Ready` |
| Conformance tests pass | No unexpected test failures |
| Cluster destroy completes | All GCD resources removed |

---

## 14. Verify GCD cluster nodes use correct machine type and disk

- **Automation status:** Manual
- **Requires:** Running GCD cluster deployed via IPI

### Description

After a successful GCD IPI installation, verify that the control plane and
compute nodes are running on the expected machine type (`c3-standard-4`) and
using Hyperdisk Balanced disks. This confirms that the installer's sovereign
cloud defaults were applied correctly and that GCD provisioned the resources
as requested.

### Prerequisites

- A GCD cluster deployed via `openshift-install create cluster` with a
  domain-scoped project ID and sovereign region (`u-` prefix)
- `oc` CLI authenticated to the cluster
- `gcloud` CLI authenticated to the GCD project with universe domain configured

### Execution

1. Verify node machine types via the GCP console or `gcloud`:

```bash
export GOOGLE_CLOUD_UNIVERSE_DOMAIN=apis-berlin-build0.goog
export INFRA_ID=$(oc get -o jsonpath='{.status.infrastructureName}' infrastructure cluster)

gcloud compute instances list \
  --filter="name~${INFRA_ID}" \
  --format="table(name, machineType.basename(), zone)" \
  --project="eu0:<project-id>"
```

2. Verify disk types for the instances:

```bash
gcloud compute disks list \
  --filter="name~${INFRA_ID}" \
  --format="table(name, type.basename(), sizeGb, zone)" \
  --project="eu0:<project-id>"
```

3. Verify from the OpenShift side:

```bash
oc get nodes -o wide
oc get machines -n openshift-machine-api -o wide
```

### Pass/Fail Criteria

| Check | Pass Condition |
|-------|---------------|
| Control plane machine type | All master instances are `c3-standard-4` |
| Compute machine type | All worker instances are `c3-standard-4` |
| Boot disk type | All boot disks are `hyperdisk-balanced` |
| All nodes Ready | `oc get nodes` shows all nodes in `Ready` state |
| Zone distribution | Instances spread across `u-germany-northeast1-a/b/c` |

---

## 15. Verify GCD cluster uses private DNS only

- **Automation status:** Manual
- **Requires:** Running GCD cluster deployed via IPI

### Description

GCD does not support public DNS zones. Verify that the cluster was deployed
with internal-only DNS, that no public DNS records were created, and that the
API and ingress endpoints are only reachable through private networking.

### Prerequisites

- A GCD cluster deployed with `publish: Internal`
- `oc` CLI authenticated to the cluster (through bastion/proxy)
- `gcloud` CLI authenticated to the GCD project

### Execution

1. Verify the cluster's publish strategy:

```bash
oc get dns cluster -o jsonpath='{.spec.publicZone}'
# Expected: empty (no public zone)

oc get dns cluster -o jsonpath='{.spec.privateZone}'
# Expected: populated with the private zone ID
```

2. Verify DNS zones in GCD:

```bash
gcloud dns managed-zones list \
  --project="eu0:<project-id>" \
  --format="table(name, dnsName, visibility)"
# Expected: only "private" visibility zones
```

3. Verify API server is not publicly accessible:

```bash
# From outside the VPC (not through bastion):
curl -k https://api.<cluster-name>.<base-domain>:6443/healthz
# Expected: connection timeout or DNS resolution failure
```

4. Verify API server is accessible through private network:

```bash
# Through the bastion/proxy:
oc get clusterversion
# Expected: success, shows cluster version
```

### Pass/Fail Criteria

| Check | Pass Condition |
|-------|---------------|
| No public DNS zone | `spec.publicZone` is empty on the DNS cluster resource |
| Private zone exists | `spec.privateZone` is populated |
| GCD DNS zones are private | `gcloud dns managed-zones list` shows only `private` visibility |
| API not publicly reachable | External curl to API endpoint fails |
| API reachable via private network | `oc` commands work through bastion/proxy |

---

## 16. Verify GCD cluster destroy removes all resources

- **Automation status:** Manual
- **Requires:** A GCD cluster that has been destroyed via `openshift-install destroy cluster`

### Description

After running `openshift-install destroy cluster`, verify that all GCD
resources created during installation have been removed. GCD has limited quota
and leftover resources can block future installs and incur costs.

### Prerequisites

- A GCD cluster that was successfully installed and then destroyed
- `gcloud` CLI authenticated to the GCD project
- The cluster's `INFRA_ID` (from the metadata or install log)

### Execution

1. Check for leftover compute instances:

```bash
gcloud compute instances list \
  --filter="name~${INFRA_ID}" \
  --project="eu0:<project-id>"
# Expected: no results
```

2. Check for leftover disks:

```bash
gcloud compute disks list \
  --filter="name~${INFRA_ID}" \
  --project="eu0:<project-id>"
# Expected: no results
```

3. Check for leftover firewall rules:

```bash
gcloud compute firewall-rules list \
  --filter="name~${INFRA_ID}" \
  --project="eu0:<project-id>"
# Expected: no results
```

4. Check for leftover DNS records:

```bash
gcloud dns record-sets list \
  --zone=<cluster-zone> \
  --filter="name~${INFRA_ID}" \
  --project="eu0:<project-id>"
# Expected: no results (or zone itself deleted)
```

5. Check for leftover service accounts:

```bash
gcloud iam service-accounts list \
  --filter="email~${INFRA_ID}" \
  --project="eu0:<project-id>"
# Expected: no results
```

6. Check for leftover health checks and forwarding rules:

```bash
gcloud compute health-checks list \
  --filter="name~${INFRA_ID}" \
  --project="eu0:<project-id>"

gcloud compute forwarding-rules list \
  --filter="name~${INFRA_ID}" \
  --project="eu0:<project-id>"
# Expected: no results for both
```

### Pass/Fail Criteria

| Resource Type | Pass Condition |
|--------------|---------------|
| Compute instances | Zero instances matching `INFRA_ID` |
| Disks | Zero disks matching `INFRA_ID` |
| Firewall rules | Zero rules matching `INFRA_ID` |
| DNS records | Zero records matching `INFRA_ID` |
| Service accounts | Zero SAs matching `INFRA_ID` |
| Health checks | Zero health checks matching `INFRA_ID` |
| Forwarding rules | Zero forwarding rules matching `INFRA_ID` |

---

## 17. Verify install-config validation rejects public publish on GCD

- **Automation status:** Manual
- **Requires:** GCD service account key, `openshift-install` binary

### Description

GCD only supports private DNS zones. Attempting to create a cluster with
`publish: External` should be rejected during install-config validation. This
test verifies the installer catches the misconfiguration before attempting
any cloud API calls.

### Prerequisites

- `openshift-install` binary built from a branch with GCD support
- GCD service account key (`gce.json`) with `universe_domain` field
- No running cluster needed

### Execution

1. Create an install-config with `publish: External` and a GCD project:

```yaml
apiVersion: v1
baseDomain: example.gcd.devcluster.openshift.com
metadata:
  name: test-public-reject
platform:
  gcp:
    projectID: "eu0:<project-id>"
    region: u-germany-northeast1
controlPlane:
  platform:
    gcp:
      type: c3-standard-4
      osImage:
        name: rhcos10
        project: "eu0:<image-project>"
compute:
  - platform:
      gcp:
        type: c3-standard-4
        osImage:
          name: rhcos10
          project: "eu0:<image-project>"
publish: External
pullSecret: '<pull-secret>'
```

2. Run `openshift-install create cluster --dir=<dir>` and observe the output.

### Expected Result

- The installer rejects the configuration with an error indicating that
  public/external publish is not supported for sovereign cloud environments
- No cloud resources are created

### Pass/Fail Criteria

| Check | Pass Condition |
|-------|---------------|
| Validation error | Installer exits with a validation error mentioning publish or DNS |
| No resources created | No instances, disks, or DNS zones created in GCD |

---

## Manual Test Checklist

<!-- MANUAL_TESTS_START -->

This section collects all manual test scenarios for quick reference. Use the
checkboxes to track progress during manual verification runs.

**How to find this section:** search for `MANUAL_TESTS_START` in this file, or
run `grep -l MANUAL_TESTS_START test/docs/gcd/cases/*.md` to find all case
files with manual tests.

### Prerequisites for all manual tests

- [ ] GCD service account key available (`gce.json` with `universe_domain`)
- [ ] `gcloud` CLI installed and configured with `universe_domain=apis-berlin-build0.goog`
- [ ] `oc` CLI installed
- [ ] `openshift-install` binary built with GCD support
- [ ] RHCOS image uploaded to GCD project
- [ ] Sufficient C3 quota in `u-germany-northeast1`

### Checklist

- [ ] **14. Verify GCD cluster nodes use correct machine type and disk**
  - [ ] All control plane nodes are `c3-standard-4`
  - [ ] All compute nodes are `c3-standard-4`
  - [ ] All boot disks are `hyperdisk-balanced`
  - [ ] Instances distributed across `u-germany-northeast1-a/b/c`

- [ ] **15. Verify GCD cluster uses private DNS only**
  - [ ] No public DNS zone on cluster DNS resource
  - [ ] GCD DNS zones are all private visibility
  - [ ] API not reachable from outside VPC
  - [ ] API reachable through bastion/proxy

- [ ] **16. Verify GCD cluster destroy removes all resources**
  - [ ] Zero leftover compute instances
  - [ ] Zero leftover disks
  - [ ] Zero leftover firewall rules
  - [ ] Zero leftover DNS records
  - [ ] Zero leftover service accounts
  - [ ] Zero leftover health checks and forwarding rules

- [ ] **17. Verify install-config validation rejects public publish on GCD**
  - [ ] Installer rejects `publish: External` with validation error
  - [ ] No cloud resources created

<!-- MANUAL_TESTS_END -->

---

## Related Links

- Strategy: [OCPSTRAT-3006](https://redhat.atlassian.net/browse/OCPSTRAT-3006)
- Installer PR: [openshift/installer#10630](https://github.com/openshift/installer/pull/10630)
- CI job PR: [openshift/release#81238](https://github.com/openshift/release/pull/81238)
- Feature test plan: [gcd-feature.md](../plans/gcd-feature.md)
- GCD skill reference: `.claude/skills/gcd/`
