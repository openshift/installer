# Konnectivity Bootstrap Integration Plan

## Goal

Enable the bootstrap kube-apiserver (KAS) to access webhooks hosted in the pod network by deploying Konnectivity in the bootstrap environment.

**Scope**: All platforms, proof-of-concept quality. Prioritize a working demonstration over edge case handling.

## Architecture Overview

1. Deploy a Konnectivity server in the bootstrap environment
2. Configure the bootstrap KAS with an EgressSelectorConfiguration to proxy cluster traffic through the local Konnectivity server
3. Deploy a DaemonSet that runs a Konnectivity agent on all nodes, connecting back to the bootstrap Konnectivity server
4. Remove the Konnectivity agent DaemonSet during bootstrap teardown

## Assumptions to Validate

- [x] The bootstrap KAS will be able to access cluster-hosted webhooks via Konnectivity (validated in Task 1.7)
- [x] Non-bootstrap KAS instances will not be impacted and will continue routing normally (validated in Task 1.7)

---

## Phase 1: Investigation

All investigation tasks must be completed before proceeding to implementation.

### Task 1.1: Understand Bootstrap Environment Creation

**Objective**: Identify where and how the bootstrap environment is created in the installer.

**Status**: ✅ Complete

#### Findings

##### Bootstrap Ignition Assets

The bootstrap environment is created through a chain of assets in `pkg/asset/ignition/bootstrap/`:

| File | Responsibility |
|------|----------------|
| [common.go](pkg/asset/ignition/bootstrap/common.go) | Core bootstrap config generation (774 lines). Defines `bootstrapTemplateData`, gathers dependencies, adds files/units from `data/data/bootstrap/` |
| [bootstrap.go](pkg/asset/ignition/bootstrap/bootstrap.go) | Main `Bootstrap` asset wrapper, produces `bootstrap.ign` |
| [bootstrap_in_place.go](pkg/asset/ignition/bootstrap/bootstrap_in_place.go) | Single-node variant, produces `bootstrap-in-place-for-live-iso.ign` |
| [cvoignore.go](pkg/asset/ignition/bootstrap/cvoignore.go) | CVO manifest overrides |

Platform-specific modules exist in subdirectories (e.g., `bootstrap/baremetal/`, `bootstrap/aws/`).

##### Bootstrap KAS Configuration

The bootstrap KAS is **not pre-generated as a static pod manifest**. Instead, it is dynamically rendered during bootstrap execution by the `cluster-kube-apiserver-operator`:

1. The bootstrap machine boots with `bootstrap.ign`
2. `bootkube.service` starts and runs `/usr/local/bin/bootkube.sh`
3. The **kube-apiserver-bootstrap** stage (lines 247-279 in [bootkube.sh.template](data/data/bootstrap/files/usr/local/bin/bootkube.sh.template)) renders the KAS:
   ```bash
   KUBE_APISERVER_OPERATOR_IMAGE=$(image_for cluster-kube-apiserver-operator)
   bootkube_podman_run \
     --name kube-apiserver-render \
     "${KUBE_APISERVER_OPERATOR_IMAGE}" \
     /usr/bin/cluster-kube-apiserver-operator render \
       --manifest-etcd-serving-ca=etcd-ca-bundle.crt \
       --manifest-etcd-server-urls="${ETCD_ENDPOINTS}" \
       ...
   ```
4. The rendered config is copied to `/etc/kubernetes/bootstrap-configs/kube-apiserver-config.yaml`
5. Static pod manifests are placed in `/etc/kubernetes/manifests/` for the bootstrap kubelet

**Note on bootstrap kubelet**: The bootstrap node runs a kubelet in "standalone" mode ([kubelet.sh.template](data/data/bootstrap/files/usr/local/bin/kubelet.sh.template)). It watches `--pod-manifest-path=/etc/kubernetes/manifests` for static pods but has no `--kubeconfig`, so it never connects to the API server and does not register as a Node object.

##### Bootstrap Ignition Generation Flow

```
InstallConfig
    ↓
├─→ Manifests (cluster-config.yaml, etc.)
├─→ TLS Assets (177+ certificate/key pairs)
├─→ KubeConfigs (Admin, Kubelet, Loopback)
├─→ Machines (Master, Worker definitions)
├─→ ReleaseImage
└─→ RHCOS Image
    ↓
Common.generateConfig() → bootstrap.ign
```

Key operations in `generateConfig()`:
1. Initialize Ignition v3.2 config
2. Add files from `data/data/bootstrap/files/` (processed as Go templates)
3. Add systemd units from `data/data/bootstrap/systemd/`
4. Inject manifests, TLS assets, and kubeconfigs via `addParentFiles()`
5. Apply platform-specific customizations
6. Serialize to JSON

##### Data Files and Templates

Bootstrap templates live in `data/data/bootstrap/`:

```
bootstrap/
├── files/
│   └── usr/local/bin/
│       ├── bootkube.sh.template      # Main bootstrap orchestrator
│       ├── kubelet.sh.template
│       └── ...
├── systemd/common/units/
│   ├── bootkube.service              # Runs bootkube.sh
│   ├── kubelet.service.template
│   └── ...
└── <platform>/                       # Platform-specific overrides
```

##### Key Insight for Konnectivity Integration

Since the bootstrap KAS is rendered at runtime by `cluster-kube-apiserver-operator`, adding EgressSelectorConfiguration requires passing a config override file to the render command.

**Approach** (validated in Task 1.1.1):
1. Create an egress selector config template in `data/data/bootstrap/files/`
2. Add `--config-override-files=<path>` to the existing render invocation in `bootkube.sh.template`
3. The operator's built-in config merging will inject `egressSelectorConfiguration` into the final KubeAPIServerConfig

The Konnectivity server itself can be added as a static pod manifest in `/etc/kubernetes/manifests/`, similar to how the bootstrap KAS, etcd, and other control plane components run.

---

### Task 1.1.1: Investigate KAS Operator Render Command Extensibility

**Objective**: Determine if `cluster-kube-apiserver-operator render` can be extended to output additional KAS arguments (specifically EgressSelectorConfiguration), or if post-processing is required.

**Source**: `/home/mbooth/src/openshift/cluster-kube-apiserver-operator`

**Status**: ✅ Complete

#### Findings

##### Built-in Extensibility Mechanism

**Good news**: The render command already supports config overrides via `--config-override-files` flag.

The flag can be specified multiple times and each file is:
1. Rendered as a Go text/template
2. Merged into the final `KubeAPIServerConfig` using `resourcemerge.MergeProcessConfig()`

Merge order: `defaultConfig` → `bootstrapOverrides` → `--config-override-files` (in order provided)

##### Render Command Output

| Output | Format | Location |
|--------|--------|----------|
| KubeAPIServerConfig | YAML | `--config-output-file` (e.g., `/assets/kube-apiserver-bootstrap/config`) |
| Bootstrap manifests | YAML | `--asset-output-dir/bootstrap-manifests/` |
| Runtime manifests | YAML | `--asset-output-dir/manifests/` |

##### EgressSelectorConfiguration Injection

EgressSelectorConfiguration is a top-level field in `KubeAPIServerConfig`:

```yaml
apiVersion: kubecontrolplane.config.openshift.io/v1
kind: KubeAPIServerConfig
egressSelectorConfiguration:
  egressSelections:
  - name: "cluster"
    connection:
      proxyProtocol: "HTTPConnect"
      transport:
        uds:
          udsName: "/etc/kubernetes/konnectivity-server/konnectivity-server.socket"
```

##### Recommended Approach: Config Override File ✅

**No operator modifications required.** Use the existing `--config-override-files` mechanism:

1. Create `data/data/bootstrap/files/opt/openshift/egress-selector-config.yaml.template`:
   ```yaml
   apiVersion: kubecontrolplane.config.openshift.io/v1
   kind: KubeAPIServerConfig
   egressSelectorConfiguration:
     egressSelections:
     - name: "cluster"
       connection:
         proxyProtocol: "HTTPConnect"
         transport:
           uds:
             udsName: "/etc/kubernetes/konnectivity-server/konnectivity-server.socket"
   ```

2. Modify [bootkube.sh.template](data/data/bootstrap/files/usr/local/bin/bootkube.sh.template) to add the flag:
   ```bash
   /usr/bin/cluster-kube-apiserver-operator render \
     ... existing args ...
     --config-override-files=/assets/egress-selector-config.yaml \
     ...
   ```

##### Key Implementation Details

- The bootstrap KAS reads config via `--openshift-config` flag from a YAML file
- EgressSelectorConfiguration in the merged YAML is picked up automatically
- Konnectivity server URL should use `127.0.0.1` since it runs on the same bootstrap node
- Transport can be plain TCP (simplest) or TLS (requires additional cert provisioning)

##### Key Files

| Component | Location |
|-----------|----------|
| Render command | `cluster-kube-apiserver-operator/pkg/cmd/render/render.go` |
| Config merge logic | `vendor/github.com/openshift/library-go/pkg/operator/render/options/generic.go:210-259` |
| Bootstrap KAS pod template | `cluster-kube-apiserver-operator/bindata/bootkube/bootstrap-manifests/kube-apiserver-pod.yaml` |
| Current render invocation | `installer/data/data/bootstrap/files/usr/local/bin/bootkube.sh.template:247-269`

---

### Task 1.2: Understand Bootstrap KAS Network Limitations

**Objective**: Determine why the bootstrap KAS cannot route to the pod network.

**Status**: ✅ Complete

#### Findings

##### Bootstrap Network Topology

The bootstrap node operates in **host network mode only**:

1. **All containers use host networking** - From [bootkube.sh.template](data/data/bootstrap/files/usr/local/bin/bootkube.sh.template) line 24-27:
   ```bash
   bootkube_podman_run() {
       # we run all commands in the host-network to prevent IP conflicts with
       # end-user infrastructure.
       podman run --quiet --net=host --rm --log-driver=k8s-file "${@}"
   }
   ```

2. **Bootstrap kubelet has no CNI** - The standalone kubelet ([kubelet.sh.template](data/data/bootstrap/files/usr/local/bin/kubelet.sh.template)) runs without:
   - `--cni-bin-dir`
   - `--cni-conf-dir`
   - Any network plugin configuration

3. **Static pods run in host network namespace** - They get host IPs, not pod network IPs

##### Why Pod Network is Unreachable

| Bootstrap Node | Cluster Nodes |
|----------------|---------------|
| Host network only (`--net=host`) | Pod network + host network |
| No CNI plugin | OVN-Kubernetes or OpenShift-SDN |
| Infrastructure network (e.g., 10.0.0.0/24) | Pod CIDR routable (e.g., 10.128.0.0/14) |
| No routes to pod/service CIDRs | Full CNI routing stack |

**Root cause**: The cluster-network-operator is deployed **after** bootstrap completes. During bootstrap:
- No CNI plugin is installed
- No VXLAN/Geneve tunnels exist
- No routes to pod CIDR (10.128.0.0/14) or service CIDR (172.30.0.0/16)

##### Bootstrap Lifecycle vs Network Operator

1. **Bootstrap phase**: Bootstrap KAS runs on host network, applies manifests via cluster-bootstrap
2. **CVO deploys operators**: Including the cluster-network-operator
3. **Pod network becomes operational**: CNI is configured, pod CIDR is routable on cluster nodes
4. **Teardown conditions met**: Required control plane pods running, CEO complete, etc.
5. **Bootstrap node exits**: After all teardown conditions are satisfied

The pod network is operational *before* bootstrap completes. However, the bootstrap node itself never gains access to the pod network - it remains on host networking throughout its lifecycle. This is intentional: the bootstrap node is temporary and isolated from the production network by design.

##### Bootstrap Teardown Conditions

The bootstrap node marks itself complete (`/opt/openshift/.bootkube.done`) after these conditions are met:

1. **cluster-bootstrap succeeds** - Applies manifests and waits for required pods:
   ```
   openshift-kube-apiserver/kube-apiserver
   openshift-kube-scheduler/openshift-kube-scheduler
   openshift-kube-controller-manager/kube-controller-manager
   openshift-cluster-version/cluster-version-operator
   ```

2. **CVO overrides restored** - Patches `clusterversion.config.openshift.io/version`

3. **CEO (Cluster Etcd Operator) completes** - Runs `wait-for-ceo` command

4. **API DNS checks pass** - Both `API_URL` and `API_INT_URL` are reachable

For HA clusters, the installer's `wait-for-bootstrap-complete` stage additionally waits for:
- At least 2 nodes with successful KAS revision rollout (checked via `kubeapiservers.operator.openshift.io/cluster` status)

**Key insight**: Bootstrap teardown happens when the *control plane pods* are running on cluster nodes, but the *network operator* may still be initializing. The pod network becomes fully operational after bootstrap exits.

##### Webhook Reachability Problem

When a webhook is deployed in the pod network:
- Bootstrap KAS tries to call webhook at pod IP (e.g., `10.128.2.5:8443`)
- Pod lives in cluster network CIDR
- Bootstrap node has no route to that CIDR
- Result: **connection refused / no route to host**

##### Why CNI Cannot Run on Bootstrap

1. **Dependency cycle**: CNI operator needs API server, but API server needs network for webhooks
2. **Bootstrap is temporary**: Not part of production cluster topology
3. **Design isolation**: Host network prevents conflicts with user infrastructure

---

### Task 1.3: Investigate Bootstrap Teardown Mechanism

**Objective**: Find the existing mechanism for cleaning up bootstrap resources and determine how to integrate Konnectivity agent removal.

**Status**: ✅ Complete

#### Findings

##### Bootstrap Teardown Flow

```
Bootstrap Node (bootkube.sh):
├─ cluster-bootstrap completes
├─ CVO overrides restored
├─ CEO wait completes
├─ touch /opt/openshift/.bootkube.done    ← Bootstrap script complete
└─ bootkube.service exits

Installer Process (cmd/openshift-install/create.go):
├─ waitForBootstrapComplete()             ← Polls ConfigMap status
├─ gatherBootstrapLogs() (optional)
├─ destroybootstrap.Destroy()             ← Deletes bootstrap infrastructure
│   └─ Platform-specific VM/network deletion
└─ WaitForInstallComplete()               ← Waits for cluster ready
```

##### Cluster Resource Cleanup: Current State

**Key finding**: The installer has **NO built-in mechanism** to delete cluster resources (DaemonSets, Deployments, etc.) after bootstrap completes. It only deletes infrastructure (VMs, networks, etc.).

Evidence:
- No kubectl/oc calls in `DestroyBootstrap()`
- No manifest deletion in post-bootstrap flow
- The only cleanup example is removing the MCO bootstrap static pod from the bootstrap node's local filesystem (line 621 in bootkube.sh.template)

##### Integration Point for Konnectivity Cleanup

**Location**: [pkg/destroy/bootstrap/bootstrap.go](pkg/destroy/bootstrap/bootstrap.go)

**Timing**: Delete the DaemonSet **after** infrastructure teardown completes.

**Rationale**:
- Konnectivity tunnel remains available until the bootstrap node is destroyed
- Agents become orphaned when bootstrap disappears (expected)
- We then clean up the orphaned DaemonSet resources
- Cluster API remains accessible via production KAS (already running on cluster nodes)

```go
// In pkg/destroy/bootstrap/bootstrap.go
func Destroy(ctx context.Context, dir string) error {
    // ... existing platform setup ...

    // EXISTING: Platform-specific infrastructure cleanup
    if err := provider.DestroyBootstrap(ctx, dir); err != nil {
        return fmt.Errorf("error destroying bootstrap resources %w", err)
    }

    // NEW: Clean up bootstrap-only cluster resources after infrastructure is gone
    if err := deleteBootstrapClusterResources(ctx, dir); err != nil {
        logrus.Warningf("Failed to clean up bootstrap cluster resources: %v", err)
        // Don't fail - infrastructure is already destroyed
    }

    return nil
}
```

##### Key Files

| Purpose | File |
|---------|------|
| Bootstrap completion detection | [cmd/openshift-install/create.go:164](cmd/openshift-install/create.go) |
| Bootstrap infrastructure destroy | [pkg/destroy/bootstrap/bootstrap.go](pkg/destroy/bootstrap/bootstrap.go) |
| CAPI machine deletion | [pkg/infrastructure/clusterapi/clusterapi.go](pkg/infrastructure/clusterapi/clusterapi.go) |
| Bootstrap script | [bootkube.sh.template:574-655](data/data/bootstrap/files/usr/local/bin/bootkube.sh.template) |

##### Recommendation

Add Konnectivity cleanup to `pkg/destroy/bootstrap/bootstrap.go`:
1. Create function to delete Konnectivity DaemonSet (and any related resources)
2. Call it **after** `provider.DestroyBootstrap()` returns successfully
3. Handle errors gracefully (warn but don't fail - infrastructure is already gone)
4. Delete order: DaemonSet → ConfigMaps → Namespace (if bootstrap-only)

---

### Task 1.4: Investigate Konnectivity in OpenShift Payload

**Objective**: Determine if Konnectivity components are already available in the OpenShift payload.

**Status**: ✅ Complete (answered by Task 1.5)

#### Findings

**Image name**: `apiserver-network-proxy`

**Availability**: Present in OpenShift release payload

**Binaries included**:
- `/usr/bin/proxy-server` - Konnectivity server
- `/usr/bin/proxy-agent` - Konnectivity agent

**How to obtain**: Use the existing `image_for` function in bootkube.sh:
```bash
KONNECTIVITY_IMAGE=$(image_for apiserver-network-proxy)
```

---

### Task 1.5: Investigate HyperShift Konnectivity Deployment

**Objective**: Examine how HyperShift deploys Konnectivity to inform our implementation.

**Source**: `/home/mbooth/src/openshift/hypershift`

**Status**: ✅ Complete

#### Findings

##### Image Source

**Image name**: `apiserver-network-proxy`
- Available in OpenShift release payload
- Contains both `/usr/bin/proxy-server` and `/usr/bin/proxy-agent`

##### Konnectivity Server Configuration

HyperShift runs the server as a **sidecar container** in the kube-apiserver pod:

```bash
/usr/bin/proxy-server
  --logtostderr=true
  --cluster-cert /etc/konnectivity/cluster/tls.crt
  --cluster-key /etc/konnectivity/cluster/tls.key
  --server-cert /etc/konnectivity/server/tls.crt
  --server-key /etc/konnectivity/server/tls.key
  --server-ca-cert /etc/konnectivity/ca/ca.crt
  --server-port 8090          # Client endpoint (KAS connects here)
  --agent-port 8091           # Agent endpoint (agents connect here)
  --health-port 2041
  --admin-port 8093
  --mode http-connect
  --proxy-strategies destHost,defaultRoute
  --keepalive-time 30s
  --frontend-keepalive-time 30s
```

**Key ports**:
- **8090**: Server endpoint (for KAS/proxy clients via HTTPConnect)
- **8091**: Agent endpoint (where agents connect)
- **2041**: Health check

##### Konnectivity Agent Configuration

HyperShift runs agents as a **DaemonSet** on worker nodes:

```bash
/usr/bin/proxy-agent
  --logtostderr=true
  --ca-cert /etc/konnectivity/ca/ca.crt
  --agent-cert /etc/konnectivity/agent/tls.crt
  --agent-key /etc/konnectivity/agent/tls.key
  --proxy-server-host <server-address>
  --proxy-server-port 8091
  --health-server-port 2041
  --agent-identifiers default-route=true
  --keepalive-time 30s
  --probe-interval 5s
  --sync-interval 5s
  --sync-interval-cap 30s
```

**DaemonSet settings**:
- `hostNetwork: true` (except IBMCloud)
- `dnsPolicy: Default` (use host resolver)
- Tolerates all taints
- Rolling update with 10% maxUnavailable

##### Authentication: mTLS

All Konnectivity communication uses **mutual TLS**. HyperShift generates these certificates:

| Secret | Purpose |
|--------|---------|
| `konnectivity-signer` | Self-signed CA (signs all other certs) |
| `konnectivity-server` | Server cert for client endpoint (port 8090) |
| `konnectivity-cluster` | Server cert for agent endpoint (port 8091) |
| `konnectivity-client` | Client cert for KAS/proxies → server |
| `konnectivity-agent` | Agent cert for agents → server |

##### EgressSelectorConfiguration

HyperShift configures KAS with this egress selector:

```yaml
apiVersion: apiserver.k8s.io/v1beta1
kind: EgressSelectorConfiguration
egressSelections:
- name: controlplane
  connection:
    proxyProtocol: Direct
- name: etcd
  connection:
    proxyProtocol: Direct
- name: cluster
  connection:
    proxyProtocol: HTTPConnect
    transport:
      tcp:
        url: https://konnectivity-server-local:8090
        tlsConfig:
          caBundle: /etc/kubernetes/certs/konnectivity-ca/ca.crt
          clientCert: /etc/kubernetes/certs/konnectivity-client/tls.crt
          clientKey: /etc/kubernetes/certs/konnectivity-client/tls.key
```

**Key insight**: The `cluster` egress type uses `HTTPConnect` protocol with TLS, not plain TCP.

##### Simplifications for Bootstrap POC

For our proof-of-concept, we can simplify:

1. **Server location**: Static pod on bootstrap node (not sidecar)
2. **Agent connection**: Agents connect to bootstrap node IP directly (no Route/Service)
3. **Authentication**: Can start with plain TCP, add mTLS later
4. **Single server**: No need for HA server count

##### Key HyperShift Files

| Component | File |
|-----------|------|
| PKI setup | `control-plane-operator/controllers/hostedcontrolplane/pki/konnectivity.go` |
| Server config | `control-plane-operator/controllers/hostedcontrolplane/v2/kas/deployment.go` |
| Agent DaemonSet | `control-plane-operator/hostedclusterconfigoperator/controllers/resources/konnectivity/reconcile.go` |
| Egress selector | `control-plane-operator/controllers/hostedcontrolplane/v2/assets/kube-apiserver/egress-selector-config.yaml` |

---

### Task 1.6: Research Konnectivity Authentication

**Objective**: Determine how the Konnectivity agent authenticates to the Konnectivity server.

**Status**: ✅ Complete

#### Decision: No Authentication for PoC

For the proof-of-concept, **we will skip authentication entirely**:

1. **KAS → Konnectivity Server**: The bootstrap KAS connects to the Konnectivity server via **Unix Domain Socket (UDS)**. Since both run on the same bootstrap node, UDS is the recommended transport and no authentication is required.

2. **Agents → Konnectivity Server**: Agents connect to the Konnectivity server over **unauthenticated TCP**. The server will accept connections without mTLS client certificates.

##### Security Considerations

This approach is **insecure** but acceptable for a PoC because:
- The bootstrap environment is temporary and isolated
- The Konnectivity server only exists during cluster bootstrap
- Production-ready implementations should use mTLS (as HyperShift does)

##### Server Configuration (Unauthenticated)

```bash
/usr/bin/proxy-server
  --logtostderr=true
  --uds-name=/etc/kubernetes/konnectivity-server/konnectivity-server.socket  # KAS connects via UDS
  --agent-port 8091           # Agents connect here (TCP, no auth)
  --health-port 2041
  --mode http-connect
  --proxy-strategies destHost,defaultRoute
  --keepalive-time 30s
  --frontend-keepalive-time 30s
```

##### Agent Configuration (Unauthenticated)

```bash
/usr/bin/proxy-agent
  --logtostderr=true
  --proxy-server-host <bootstrap-node-ip>
  --proxy-server-port 8091
  --health-server-port 2041
  --agent-identifiers default-route=true
  --keepalive-time 30s
  --probe-interval 5s
  --sync-interval 5s
```

##### Future Work

For production readiness, Task 2.1 (certificate generation) would need to be completed to add mTLS authentication. This is out of scope for the initial PoC.

---

### Task 1.7: Validate Assumptions via Documentation

**Objective**: Verify the architectural assumptions through documentation research.

**Status**: ✅ Complete

#### Assumption 1: EgressSelectorConfiguration Routes Webhook Traffic ✅ Confirmed

The `cluster` egress selector covers **all traffic destined for the cluster**, including:
- Webhooks (ValidatingWebhookConfiguration, MutatingWebhookConfiguration)
- Aggregated API servers (APIService resources)
- Pod operations (exec, attach, logs, port-forward)
- Node and service proxy requests

From the [KEP-1281 Network Proxy](https://github.com/kubernetes/enhancements/blob/master/keps/sig-api-machinery/1281-network-proxy/README.md):
> "Pod requests (and pod sub-resource requests) are meant for the cluster and will be routed based on the 'cluster' NetworkContext."

#### Assumption 2: Non-Bootstrap KAS Routes Normally ✅ Confirmed

The Konnectivity feature is **opt-in** and disabled by default. From KEP-1281:
> "The feature is turned off in the KAS by default. Enabled by adding ConnectivityServiceConfiguration."

When no `--egress-selector-config-file` flag is set, the KAS uses direct connections via standard network routing. This means:
- Production KAS instances on cluster nodes (without EgressSelectorConfiguration) will continue to route traffic directly
- Only the bootstrap KAS (with explicit configuration) will use the Konnectivity tunnel

#### Caveats and Limitations

1. **Flat network requirement**: The proxy requires non-overlapping IP ranges between control plane and cluster networks. This is satisfied in OpenShift since bootstrap and cluster networks are distinct.

2. **DNS resolution**: DNS lookups are performed by the Konnectivity server/agent, not the KAS. This should be transparent but may affect debugging.

3. **Protocol consistency**: The `proxyProtocol` (GRPC or HTTPConnect) must match between the KAS EgressSelectorConfiguration and Konnectivity server arguments.

4. **Transport options**:
   - **UDS (Unix Domain Socket)**: Recommended when KAS and Konnectivity server are co-located. **We will use UDS** for KAS → Konnectivity server communication since both run on the bootstrap node.
   - **TCP**: Used for agent → server communication (agents connect from cluster nodes to bootstrap node)

#### Sources

- [Set up Konnectivity service | Kubernetes](https://kubernetes.io/docs/tasks/extend-kubernetes/setup-konnectivity/)
- [KEP-1281: Network Proxy](https://github.com/kubernetes/enhancements/blob/master/keps/sig-api-machinery/1281-network-proxy/README.md)
- [apiserver-network-proxy (Konnectivity)](https://github.com/kubernetes-sigs/apiserver-network-proxy)
- [Konnectivity in HyperShift](https://hypershift.pages.dev/reference/konnectivity/)

---

## Phase 2: Implementation

### Task 2.1: Generate Konnectivity Certificates

**Status**: ⏭️ Skipped (out of scope for PoC)

**Objective**: Create the necessary certificates for Konnectivity server and agent authentication.

**Note**: Per Task 1.6 decision, the PoC uses unauthenticated TCP connections. Certificate generation is deferred to future production-ready work.

---

### Task 2.2: Deploy Konnectivity Server on Bootstrap

**Objective**: Add a Konnectivity server to the bootstrap environment.

**Deliverables**:
- Create Konnectivity server configuration
- Add server deployment to bootstrap Ignition config (likely as a static pod or systemd unit)
- Configure server to listen for agent connections
- Ensure server starts before or alongside the bootstrap KAS

---

### Task 2.3: Configure Bootstrap KAS with EgressSelectorConfiguration

**Objective**: Configure the bootstrap KAS to proxy cluster-bound traffic through Konnectivity.

**Deliverables**:
- Create EgressSelectorConfiguration for the bootstrap KAS
- Configure the KAS to use Konnectivity for "Cluster" egress traffic
- Add configuration to bootstrap KAS static pod manifest

---

### Task 2.4: Create Konnectivity Agent DaemonSet

**Objective**: Deploy Konnectivity agents on cluster nodes to connect back to the bootstrap server.

**Deliverables**:
- Create DaemonSet manifest for Konnectivity agent
- Configure agent to connect to the bootstrap node's Konnectivity server
- Include necessary certificates/credentials for authentication
- Add manifest to bootstrap resources that get applied to the cluster

---

### Task 2.5: Integrate Teardown with Bootstrap Completion

**Objective**: Remove the Konnectivity agent DaemonSet when bootstrap completes.

**Deliverables**:
- Add Konnectivity DaemonSet deletion to the bootstrap teardown flow
- Ensure clean removal without impacting cluster operation
- Test that non-bootstrap KAS continues operating normally after removal

---

### Task 2.6: Manual Validation

**Objective**: Verify the proof-of-concept works end-to-end.

**Validation steps**:
1. Deploy a cluster with the modified installer
2. Verify Konnectivity server is running on bootstrap node
3. Verify Konnectivity agents are running on cluster nodes
4. Deploy a test webhook in the cluster
5. Confirm the bootstrap KAS can reach the webhook
6. Complete bootstrap teardown and verify agent DaemonSet is removed
7. Confirm non-bootstrap KAS can still reach webhooks normally

---

## Open Questions

- Which port should the Konnectivity server listen on?
- Should the Konnectivity server be a static pod or systemd service?
- What is the expected lifecycle overlap between bootstrap and cluster control plane?
