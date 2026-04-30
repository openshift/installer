# Custom Machine Pools at Day 0 via InstallConfig

## Background

The `Compute` field in `InstallConfig` is already a `[]MachinePool` array, and the `Worker` asset already loops over all entries to generate MachineSets — so the plumbing is partially there. The blockers today are:

1. Validation rejects any pool name that isn't `"worker"` or `"edge"` (`pkg/types/validation/installconfig.go`)
2. No `MachineConfigPool` (MCP) manifest is generated for pools the installer doesn't know about — MCO won't manage the nodes without one at day 0
3. Custom pool machines need a pool-specific pointer ignition stub so the MCS serves the right MachineConfigs

---

## Decisions

| Decision | Detail |
|---|---|
| Custom pool naming | Any valid DNS label is allowed; no prescribed naming convention |
| Max custom pools | **5** custom pools per install config |
| Platform support | GCP only for the initial prototype |
| MachineSet role | Custom pool MachineSets use the `worker` role for tags, subnets, and MAPI labels — only the user-data secret differs |
| Worker balancing | Only auto-balance when the user has not explicitly set a worker replica count |

---

## Phase 1 — Allow Custom Pool Names ✅

**Files changed:**
- `pkg/types/machinepools.go`
- `pkg/types/validation/installconfig.go`
- `pkg/types/defaults/machinepools.go`

**What was done:**

1. **Reserved-name set** added to `machinepools.go`:
   ```go
   var ReservedMachinePoolNames = sets.New(
       MachinePoolComputeRoleName,      // "worker"
       MachinePoolEdgeRoleName,         // "edge"
       MachinePoolControlPlaneRoleName, // "master"
       MachinePoolArbiterRoleName,      // "arbiter"
   )

   func IsCustomPool(poolName string) bool {
       return poolName != "" && !ReservedMachinePoolNames.Has(poolName)
   }
   ```

2. **`validateCompute()`** relaxed: name must be a valid DNS label, must not collide with reserved names, GCP-only platform restriction, at most 5 custom pools.

3. **Defaults**: custom pools default to **0 replicas** so users must opt in. `IsCustomPool("")` returns false so unnamed pools (test fixtures) are unaffected.

---

## Phase 2 — MachineConfigPool Manifest Generation ✅

During bootstrap, MCO runs in bootstrap mode and generates `master` and `worker` MCPs itself. The bootstrap controller ingests installer-written manifests from `openshift/` (via `cp openshift/* /etc/mcc/bootstrap/` in `bootkube.sh`). MCO has no knowledge of custom pool names, so the installer writing the MCP into `openshift/` is the **only** path for a custom pool to be recognised during bootstrap.

**Files created/modified:**
- `pkg/asset/machines/machineconfigpool/manifest.go` — pure helper functions, no Asset interface; mirrors the `machineconfig/` package pattern
- `pkg/asset/machines/worker.go` — calls MCP generation per custom pool

**MCP structure generated** (per custom pool named `<pool-name>`):

```yaml
apiVersion: machineconfiguration.openshift.io/v1
kind: MachineConfigPool
metadata:
  name: <pool-name>
spec:
  machineConfigSelector:
    matchExpressions:
      - key: machineconfiguration.openshift.io/role
        operator: In
        values: [worker, <pool-name>]   # inherits worker configs
  nodeSelector:
    matchLabels:
      node-role.kubernetes.io/<pool-name>: ""
```

`worker` is included in `matchExpressions` so the custom pool inherits base worker MachineConfigs. MCO uses `nodeSelector` (not `machineSelector`) to assign nodes; pool placement on day 2 is driven by the node's initial MCS request to `/config/<pool-name>`.

File naming: `99_openshift-machineconfig_<pool-name>-mcp.yaml`

---

## Phase 3 — Custom Pool Ignition and User Data Secret ✅

**The pointer ignition stub must be pool-specific.** `pkg/asset/ignition/machine/node.go:pointerIgnitionConfig` embeds the role as the MCS URL path (`https://api-int.<cluster>.<domain>:22623/config/<role>`). A custom pool machine booting with `/config/worker` would receive worker MachineConfigs, not its own pool's configs.

**Files created/modified:**
- `pkg/asset/ignition/machine/custompool.go` — generates one `<pool-name>.ign` per custom pool; holds results in `FilesByPool map[string]*asset.File`
- `pkg/asset/machines/worker.go` — generates `<pool-name>-user-data` secret per pool; stores them in `CustomPoolUserDataFiles []*asset.File`

**GCP-specific implementation details:**

Custom pool MachineSets use `"worker"` as the GCP role (not the pool name). This ensures:
- Standard worker firewall tags apply (`<clusterID>-worker`)
- Standard worker subnet is used (`<clusterID>-worker-subnet`)
- MAPI labels are `machine-role: worker` / `machine-type: worker`

Only the user-data secret differs, pointing machines at `/config/<pool-name>` on the MCS. MCO handles pool assignment from there.

`pkg/asset/machines/gcp/machines.go` — `getTags` and `getNetworks` default cases now handle arbitrary roles gracefully (fall through to worker behaviour) as a safety net.

`pkg/asset/manifests/mco.go` — `gcpBootImages()` previously only checked `ic.Compute[0]` for a custom `OSImage`, so a custom pool with an overridden boot image would not disable automatic boot image management. Fixed to loop all compute pools (matching the AWS/Azure pattern), so setting `OSImage` on any custom pool correctly disables MCO's automatic MachineSet boot image updates cluster-wide.

**Worker replica balancing:**
- The worker pool's replica count defaults to 3 (`pkg/types/defaults/machinepools.go`)
- Whether the user explicitly set the worker count is determined by parsing `installConfig.File.Data` (raw YAML, pre-defaults) — if `replicas` is absent for the worker pool in the raw YAML, it was defaulted
- If worker count was **not** explicitly set: worker replicas = `max(0, workerDefault - sum(allCustomPoolReplicas))`
- If worker count **was** explicitly set: both counts are respected as-is
- Example: worker unspecified, custom=2 → worker adjusted to 1
- Example: worker=2 explicit, custom=2 → both kept, worker=2 custom=2
- Example: worker unspecified, custom=4 → worker=0 (warn logged), custom=4

**MachineSet topology note:**
All platforms create **one MachineSet per availability zone** (or failure domain), with total replicas spread evenly (`total / numZones`, remainder distributed round-robin). The replica count on the pool is the *total* across all MachineSets for that pool, not per-MachineSet.

---

## Phase 4 — CAPI Path (deferred)

**File:** `pkg/asset/machines/clusterapi.go`

`ClusterAPIComputeInstall` gate path would need the same treatment as Phase 3. Deferred — CAPI compute is still feature-gated and not broadly enabled.

---

## Phase 5 — Tests and Docs (pending)

1. **Validation tests** (`pkg/types/validation/installconfig_test.go`):
   - Valid custom pool name
   - Reserved name collision
   - Invalid DNS label
   - More than 5 custom pools → validation error
   - Custom pool on non-GCP platform → error
   - Existing `"worker"` / `"edge"` cases still pass

2. **MCP generation tests** (`pkg/asset/machines/machineconfigpool/manifest_test.go`):
   - Correct `nodeSelector` (not `machineSelector`)
   - `matchExpressions` includes both `worker` and pool name

3. **User docs** (`docs/user/`): naming rules, replica defaults, MCP inheritance

---

## Example InstallConfig

```yaml
compute:
  - name: worker
    replicas: 2
  - name: gpu-pool
    replicas: 3
    platform:
      gcp:
        zones: [us-east4-a, us-east4-b]
```

This produces:
- MachineSets: `cluster-worker-*` (2 total replicas), `cluster-gpu-pool-a` (2 replicas), `cluster-gpu-pool-b` (1 replica)
- MCP `gpu-pool` inheriting worker MachineConfigs, selecting nodes with `node-role.kubernetes.io/gpu-pool: ""`
- `gpu-pool.ign` pointer ignition pointing to `/config/gpu-pool` on the MCS
- `gpu-pool-user-data` secret wrapping the custom ignition
