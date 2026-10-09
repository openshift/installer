# Test Case: GCD cluster uses private DNS only

| Field | Value |
|-------|-------|
| **Feature** | [OCPSTRAT-3006](https://redhat.atlassian.net/browse/OCPSTRAT-3006) - Google Cloud Dedicated |
| **Component** | Installer / GCP platform |
| **Type** | Manual |
| **Priority** | P1 |
| **Test plan** | [gcd-feature.md](../plans/gcd-feature.md) |

GCD does not support public DNS zones, so a GCD cluster must be installed with
`publish: Internal`. This confirms no public zone or public record was created
and that the API is only reachable from inside the VPC.

## Setup

- A GCD cluster installed with `publish: Internal`.
- `oc` authenticated to the cluster (through the bastion proxy).
- `gcloud` authenticated to the GCD project.
- A shell *outside* the cluster VPC, for the external reachability check.

```bash
export GOOGLE_CLOUD_UNIVERSE_DOMAIN=apis-berlin-build0.goog
export PROJECT_ID="eu0:<gcd-project>"
export CLUSTER_NAME="<cluster-name>"
export BASE_DOMAIN="<base-domain>"
```

## Test

### Step

Inspect the cluster DNS config.

```bash
oc get dns cluster -o jsonpath='{.spec.publicZone}{"\n"}{.spec.privateZone}{"\n"}'
```

### Expect

`publicZone` is empty and `privateZone` is populated with the private zone ID.

### Step

List the DNS zones in the GCD project.

```bash
gcloud dns managed-zones list --project="${PROJECT_ID}" \
  --format="table(name, dnsName, visibility)"
```

### Expect

Only zones with `private` visibility are listed.

### Step

From outside the VPC, try to reach the API endpoint.

```bash
curl -k --max-time 15 "https://api.${CLUSTER_NAME}.${BASE_DOMAIN}:6443/healthz"
```

### Expect

The request fails with a DNS resolution failure or a connection timeout.

### Step

From inside the VPC (through the bastion proxy), talk to the cluster.

```bash
oc get clusterversion
```

### Expect

The command succeeds and reports the cluster version.

## Cleanup

None. Every step is read-only.
