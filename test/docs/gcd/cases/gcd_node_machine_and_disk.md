# Test Case: GCD nodes use sovereign machine and disk defaults

| Field | Value |
|-------|-------|
| **Feature** | [OCPSTRAT-3006](https://redhat.atlassian.net/browse/OCPSTRAT-3006) - Google Cloud Dedicated |
| **Component** | Installer / GCP platform |
| **Type** | Manual |
| **Priority** | P1 |
| **Test plan** | [gcd-feature.md](../plans/gcd-feature.md) |

GCD sovereign regions only offer C3, M3, and A3 Edge machines, and only
Hyperdisk Balanced disks. This confirms the installer's sovereign defaults
(`c3-standard-4` on `hyperdisk-balanced`) reached the cloud and that GCD
provisioned what was asked for.

## Setup

- A GCD cluster installed with `openshift-install create cluster`, using a
  domain-scoped project ID and a `u-` region.
- `oc` authenticated to the cluster (through the bastion proxy).
- `gcloud` authenticated to the GCD project.

```bash
export GOOGLE_CLOUD_UNIVERSE_DOMAIN=apis-berlin-build0.goog
export PROJECT_ID="eu0:<gcd-project>"
export INFRA_ID="$(oc get -o jsonpath='{.status.infrastructureName}' infrastructure cluster)"
```

## Test

### Step

List the cluster instances and their machine types.

```bash
gcloud compute instances list --project="${PROJECT_ID}" \
  --filter="name~${INFRA_ID}" \
  --format="table(name, machineType.basename(), zone)"
```

### Expect

Every control plane and compute instance is `c3-standard-4`, spread across
`u-germany-northeast1-a`, `-b`, and `-c`.

### Step

List the disks belonging to the cluster.

```bash
gcloud compute disks list --project="${PROJECT_ID}" \
  --filter="name~${INFRA_ID}" \
  --format="table(name, type.basename(), sizeGb, zone)"
```

### Expect

Every boot disk is of type `hyperdisk-balanced`.

### Step

Check the cluster's own view of its nodes.

```bash
oc get nodes
```

### Expect

All control plane and compute nodes are `Ready`.

## Cleanup

None. Every step is read-only.
