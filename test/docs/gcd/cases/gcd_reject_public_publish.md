# Test Case: GCD install-config rejects public publish

| Field | Value |
|-------|-------|
| **Feature** | [OCPSTRAT-3006](https://redhat.atlassian.net/browse/OCPSTRAT-3006) - Google Cloud Dedicated |
| **Component** | Installer / GCP platform |
| **Type** | Manual |
| **Priority** | P1 |
| **Test plan** | [gcd-feature.md](../plans/gcd-feature.md) |

GCD has no public DNS, so `publish: External` cannot work there. The installer
must reject it during validation, before any cloud resource is created.
`create manifests` runs that validation without provisioning anything.

## Setup

- An `openshift-install` binary built with GCD support.
- A GCD service account key (`gce.json`) containing a `universe_domain` field,
  exported through `GOOGLE_APPLICATION_CREDENTIALS` or placed in
  `~/.gcp/osServiceAccount.json`.
- No running cluster is needed.

```bash
export INSTALL_DIR="$(mktemp -d)"
```

## Test

### Step

Write an install-config that targets GCD but asks for external publishing.

```bash
cat > "${INSTALL_DIR}/install-config.yaml" <<'EOF'
apiVersion: v1
baseDomain: example.gcd.devcluster.openshift.com
metadata:
  name: test-public-reject
platform:
  gcp:
    projectID: "eu0:<gcd-project>"
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
EOF
```

### Expect

The file is written. Nothing has contacted GCD yet.

### Step

Run validation via manifest generation.

```bash
openshift-install create manifests --dir="${INSTALL_DIR}"
```

### Expect

The command exits non-zero with a validation error naming `publish` or the
missing public DNS zone, and no instances, disks, or DNS zones appear in the
GCD project.

## Cleanup

```bash
rm -rf "${INSTALL_DIR}"
```
