# Test Case: GCD cluster destroy removes all resources

| Field | Value |
|-------|-------|
| **Feature** | [OCPSTRAT-3006](https://redhat.atlassian.net/browse/OCPSTRAT-3006) - Google Cloud Dedicated |
| **Component** | Installer / GCP platform |
| **Type** | Manual |
| **Priority** | P2 |
| **Test plan** | [gcd-feature.md](../plans/gcd-feature.md) |

GCD projects have tight quota, so leftover resources block the next install and
keep costing money. This confirms `openshift-install destroy cluster` removes
everything it created in the sovereign region.

## Setup

- A GCD cluster that installed successfully, with its asset directory still
  present (`metadata.json` is needed by `destroy`).
- `gcloud` authenticated to the GCD project.

```bash
export GOOGLE_CLOUD_UNIVERSE_DOMAIN=apis-berlin-build0.goog
export PROJECT_ID="eu0:<gcd-project>"
export INSTALL_DIR="<asset-directory>"
export INFRA_ID="$(jq -r .infraID "${INSTALL_DIR}/metadata.json")"
```

## Test

### Step

Destroy the cluster.

```bash
openshift-install destroy cluster --dir="${INSTALL_DIR}" --log-level=info
```

### Expect

The command exits 0 and reports that the cluster was destroyed.

### Step

Look for any leftover resource whose name carries the infra ID.

```bash
for res in "compute instances" "compute disks" "compute firewall-rules" \
           "compute forwarding-rules" "compute health-checks" \
           "compute target-pools" "compute addresses" "compute routers" \
           "compute networks subnets" "compute networks" \
           "iam service-accounts"; do
  echo "== ${res}"
  # shellcheck disable=SC2086
  gcloud ${res} list --project="${PROJECT_ID}" \
    --filter="name~${INFRA_ID}" --format="value(name)"
done
```

### Expect

No names are printed under any heading.

### Step

Check the cluster's private DNS zone, if it still exists.

```bash
ZONE="${INFRA_ID}-private-zone"
if gcloud dns managed-zones describe "${ZONE}" --project="${PROJECT_ID}" &>/dev/null; then
  gcloud dns record-sets list --zone="${ZONE}" --project="${PROJECT_ID}" \
    --filter="name~${INFRA_ID}" --format="value(name)"
else
  echo "zone ${ZONE} already deleted"
fi
```

### Expect

Either the zone is gone, or it holds no records matching the infra ID.

## Cleanup

Delete any resource the previous steps found, then file a bug against the
installer's GCP destroyer (`pkg/destroy/gcp/`) describing which resource type
leaked.
