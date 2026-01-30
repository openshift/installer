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
          udsName: "/etc/kubernetes/bootstrap-configs/konnectivity-server.socket"
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
             udsName: "/etc/kubernetes/bootstrap-configs/konnectivity-server.socket"
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

1. **KAS → Konnectivity Server**: The bootstrap KAS connects to the Konnectivity server via **Unix Domain Socket (UDS)** at `/etc/kubernetes/bootstrap-configs/konnectivity-server.socket`. This path is already mounted in the KAS pod (as the `config` volume), so no additional volume mounts are needed. Since both run on the same bootstrap node, UDS is the recommended transport and no authentication is required.

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
  --uds-name=/etc/kubernetes/bootstrap-configs/konnectivity-server.socket  # KAS connects via UDS (already mounted in KAS pod)
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

### Task 1.8: Investigate Konnectivity Agent Authentication Options

**Status**: ✅ Complete

**Objective**: Determine how to authenticate Konnectivity agents to the server, leveraging existing bootstrap certificate infrastructure where possible.

**Context**: The current PoC implementation fails because the Konnectivity agent requires a CA certificate to validate the server connection:

```
Error: failed to run proxy connection with failed to read CA cert : read .: is a directory
```

#### Findings

##### 1. Konnectivity Binary Command-Line Flags (Validated via Execution)

The actual command-line flags were validated by executing the binaries from the OpenShift release payload:

**proxy-agent flags** (certificate-related):
| Flag | Default | Description |
|------|---------|-------------|
| `--ca-cert` | `""` (empty) | "If non-empty the CAs we use to validate clients." |
| `--agent-cert` | `""` (empty) | "If non-empty secure communication with this cert." |
| `--agent-key` | `""` (empty) | "If non-empty secure communication with this key." |

**proxy-server flags** (certificate-related):
| Flag | Default | Description |
|------|---------|-------------|
| `--cluster-cert` | `""` (empty) | "If non-empty secure communication with this cert." (agent-facing) |
| `--cluster-key` | `""` (empty) | "If non-empty secure communication with this key." (agent-facing) |
| `--cluster-ca-cert` | `""` (empty) | "If non-empty the CA we use to validate Agent clients." |
| `--server-cert` | `""` (empty) | "If non-empty secure communication with this cert." (KAS-facing) |
| `--server-key` | `""` (empty) | "If non-empty secure communication with this key." (KAS-facing) |
| `--server-ca-cert` | `""` (empty) | "If non-empty the CA we use to validate KAS clients." |
| `--uds-name` | `""` (empty) | "uds-name should be empty for TCP traffic. For UDS set to its name." |

**Key finding**: All certificate flags default to empty strings with "If non-empty" descriptions, suggesting unauthenticated mode should be supported. However, **validation testing revealed this is broken** (see below).

##### 2. Root Cause of the Errors (Validated via Testing)

Two separate issues were identified through container execution testing:

**Issue A: Upstream Bug - CA Cert Path Handling**

The error `failed to read CA cert : read .: is a directory` is caused by a bug in the upstream `apiserver-network-proxy` code:

```go
// In pkg/util/certificates.go - getCACertPool()
os.ReadFile(filepath.Clean(caFile))
```

When `caFile` is empty (`""`), `filepath.Clean("")` returns `"."` (current directory), and `os.ReadFile(".")` fails with "read .: is a directory".

**Validation:**
```bash
# Empty --ca-cert triggers the bug immediately
$ podman run --rm <image> /usr/bin/proxy-agent --proxy-server-host=127.0.0.1 --proxy-server-port=8091
Error: failed to run proxy connection with failed to read CA cert : read .: is a directory

# With valid CA cert, agent runs correctly (just needs server)
$ podman run --rm -v /tmp/ca.crt:/tmp/ca.crt:ro <image> /usr/bin/proxy-agent \
    --ca-cert=/tmp/ca.crt --proxy-server-host=127.0.0.1 --proxy-server-port=8091
E0130 "cannot connect once" err="...dial tcp 127.0.0.1:8091: connect: connection refused"
```

**Conclusion**: Unauthenticated mode is **broken upstream**. A valid CA certificate is required.

**Issue B: Empty Bootstrap IP**

The current PoC has `--proxy-server-host=` (empty) because `{{.BootstrapNodeIP}}` is not populated at ignition generation time:

```bash
# From pod describe output:
Args:
  --proxy-server-host=           # Empty!
  --proxy-server-port=8091
```

**Root Cause Analysis:**

The `BootstrapNodeIP` template variable is sourced from the `OPENSHIFT_INSTALL_BOOTSTRAP_NODE_IP` environment variable ([common.go:347-351](pkg/asset/ignition/bootstrap/common.go#L347)):

```go
bootstrapNodeIP := os.Getenv("OPENSHIFT_INSTALL_BOOTSTRAP_NODE_IP")
```

This environment variable is **not automatically set** by the installer for IPI (Installer Provisioned Infrastructure) deployments because:
1. Ignition configs are generated **before** the bootstrap VM is created
2. The bootstrap VM's IP is assigned dynamically by the cloud provider or DHCP
3. The IP isn't known until after Cluster API provisions the infrastructure

**Why `api-int` DNS Cannot Be Used:**

The initial thought was to use `api-int.<cluster-domain>` as the server host since it's a well-known DNS name. However, this approach is **NOT viable**:

- `api-int` resolves to a **VIP** (baremetal/vSphere) or **load balancer** (cloud), NOT specifically to the bootstrap node
- On baremetal, the VIP can move to control plane nodes via VRRP before bootstrap teardown
- On cloud platforms, load balancers distribute traffic to any healthy backend, including control plane nodes
- Konnectivity agents might connect to a node that doesn't have the Konnectivity server running

**Two Templating Layers:**

The DaemonSet manifest in `bootkube.sh.template` uses two templating layers:
1. **Go templates** (`{{.BootstrapNodeIP}}`) - Resolved at `openshift-install create ignition-configs` time
2. **Shell variables** (`${KONNECTIVITY_IMAGE}`) - Resolved at runtime when bootkube.sh runs

The problem is that `{{.BootstrapNodeIP}}` is resolved in Layer 1, when the IP isn't known yet.

**Validation:**
```bash
# Empty host dials to ":8091" (no host)
$ podman run --rm -v /tmp/ca.crt:/tmp/ca.crt:ro <image> /usr/bin/proxy-agent \
    --ca-cert=/tmp/ca.crt --proxy-server-host= --proxy-server-port=8091
E0130 "cannot connect once" err="...dial tcp :8091: connect: connection refused"
```

**Solution: Runtime IP Detection**

Detect the bootstrap node IP at **shell runtime** using `ip route get` instead of relying on the Go template variable. This command queries the kernel routing table and returns the source IP that would be used to reach an external destination:

```bash
# IPv4
BOOTSTRAP_NODE_IP=$(ip route get 1.1.1.1 2>/dev/null | awk '{for(i=1;i<=NF;i++) if($i=="src") print $(i+1); exit}')

# IPv6
BOOTSTRAP_NODE_IP=$(ip -6 route get 2001:4860:4860::8888 2>/dev/null | awk '{for(i=1;i<=NF;i++) if($i=="src") print $(i+1); exit}')
```

This approach:
- Works on any platform - detects the actual IP at runtime
- No dependency on DNS resolution or VIPs
- Uses reliable Linux networking commands
- Handles IPv4/IPv6 correctly via Go template conditional for the `ip` command variant

**See Task 2.7** for the implementation details.

##### 3. Bootstrap TLS Assets Inventory

The installer generates ~10 Certificate Authorities in `pkg/asset/tls/`. Key CAs that could potentially be reused:

| CA | File | Validity | Purpose | Reuse Candidate? |
|----|------|----------|---------|------------------|
| **RootCA** | `tls/root-ca.crt` | 10 years | Signs MCS, journal-gatewayd, IRI certs | ✅ Yes - general bootstrap CA |
| AdminKubeConfigSignerCertKey | `tls/admin-kubeconfig-signer.crt` | 10 years | Signs admin kubeconfig | Maybe |
| KubeletBootstrapCertSigner | `tls/kubelet-bootstrap-kubeconfig-signer.crt` | 10 years | Signs kubelet bootstrap cert | Maybe |
| AggregatorSignerCertKey | `tls/aggregator-signer.crt` | 1 day | Signs aggregator client | No (short-lived) |
| KubeAPIServerToKubeletSignerCertKey | `tls/kube-apiserver-to-kubelet-signer.crt` | 1 year | KAS→kubelet auth | No (wrong purpose) |

**Recommended CA**: Reuse the existing **RootCA** (`/opt/openshift/tls/root-ca.crt` and `.key`) which is already available on the bootstrap node. This avoids creating new installer TLS assets while providing full mTLS support. Certificates signed by RootCA will have 1-day validity, appropriate for the temporary bootstrap phase.

##### 4. HyperShift Konnectivity PKI Implementation

HyperShift creates a dedicated PKI for Konnectivity in [pki/konnectivity.go](hypershift/control-plane-operator/controllers/hostedcontrolplane/pki/konnectivity.go):

| Function | CN | Purpose | Usage |
|----------|----|---------| ------|
| `ReconcileKonnectivitySignerSecret` | `konnectivity-signer` | Self-signed CA | Signs all Konnectivity certs |
| `ReconcileKonnectivityServerSecret` | `konnectivity-server-local` | Server cert for KAS endpoint | `--server-cert/key` |
| `ReconcileKonnectivityClusterSecret` | `konnectivity-server` | Server cert for agent endpoint | `--cluster-cert/key` |
| `ReconcileKonnectivityClientSecret` | `konnectivity-client` | Client cert for KAS | Used in EgressSelectorConfiguration |
| `ReconcileKonnectivityAgentSecret` | `konnectivity-agent` | Client cert for agents | `--agent-cert/key` |

**Key insight**: All agents share a **single client certificate** (`konnectivity-agent`). Per-agent certificates are not required.

##### 5. Certificate Distribution Mechanism (HyperShift)

HyperShift distributes certificates via Kubernetes resources:

1. **Konnectivity CA** → `ConfigMap` named `konnectivity-ca-bundle` in `kube-system`
2. **Agent certificate** → `Secret` named `konnectivity-agent` in `kube-system`

The agent DaemonSet mounts these as volumes:
```yaml
volumes:
- name: agent-certs
  secret:
    secretName: konnectivity-agent    # Contains tls.crt, tls.key
- name: konnectivity-ca
  configMap:
    name: konnectivity-ca-bundle      # Contains ca.crt
```

Container mounts:
```yaml
volumeMounts:
- name: agent-certs
  mountPath: /etc/konnectivity/agent
- name: konnectivity-ca
  mountPath: /etc/konnectivity/ca
```

Agent args reference these paths:
```bash
--ca-cert=/etc/konnectivity/ca/ca.crt
--agent-cert=/etc/konnectivity/agent/tls.crt
--agent-key=/etc/konnectivity/agent/tls.key
```

##### 6. Minimum Viable Certificate Configuration

For the bootstrap scenario, the minimum configuration is:

**Option A: Fully Unauthenticated (Current PoC)** ❌ NOT VIABLE
- Server: No `--cluster-*` flags, no `--server-*` flags (uses UDS for KAS)
- Agent: No `--ca-cert`, `--agent-cert`, `--agent-key` flags
- **Status**: ❌ **Broken due to upstream bug** - `filepath.Clean("")` returns `"."` causing immediate crash

**Option B: Server-Only TLS (Agent validates server, no client auth)** ✅ MINIMUM VIABLE
- Server: `--cluster-cert`, `--cluster-key` (no `--cluster-ca-cert`)
- Agent: `--ca-cert` only (no `--agent-cert`, `--agent-key`)
- Agents verify the server but server accepts any agent
- **Status**: ✅ Works - validated by providing CA cert without agent certs

**Option C: Full mTLS (Production-ready)**
- Server: `--cluster-cert`, `--cluster-key`, `--cluster-ca-cert`
- Agent: `--ca-cert`, `--agent-cert`, `--agent-key`
- Mutual authentication (recommended for production)

#### Recommended Implementation Approach

Generate certificates at **runtime on the bootstrap node** using the existing `RootCA`. This approach is simpler than creating new installer TLS assets and solves the "bootstrap IP not known at ignition time" problem.

1. **Certificate generation** (in `bootkube.sh` at runtime):
   - Use `openssl` to generate Konnectivity server cert (1-day validity, signed by RootCA)
   - Use `openssl` to generate shared agent client cert (1-day validity, signed by RootCA)
   - RootCA is already available at `/opt/openshift/tls/root-ca.{crt,key}`

2. **Certificate deployment**:
   - Server certs: Written to bootstrap node filesystem, mounted into Konnectivity server static pod
   - Agent certs: Packaged into a `Secret` manifest, applied by cluster-bootstrap, mounted by DaemonSet

3. **Runtime generation flow**:
   ```
   bootkube.sh (runtime on bootstrap node)
       │
       ├─→ Generate konnectivity-server.{crt,key} (signed by RootCA)
       │   └─→ Mount in Konnectivity server static pod
       │
       ├─→ Generate konnectivity-agent.{crt,key} (signed by RootCA)
       │   └─→ Package into Secret manifest
       │
       └─→ Write manifests/konnectivity-agent-certs.yaml
           └─→ Applied by cluster-bootstrap alongside DaemonSet
   ```

4. **Files to modify**:
   | File | Purpose |
   |------|---------|
   | `data/data/bootstrap/files/usr/local/bin/bootkube.sh.template` | Certificate generation, server pod, agent Secret |

5. **Benefits**:
   - No new Go TLS assets required
   - Certificates generated when bootstrap IP is known
   - Short 1-day validity appropriate for temporary bootstrap agents
   - Full mTLS authentication between server and agents
   - All agents share a single client certificate

#### Unauthenticated Mode: Validated as Broken ❌

Testing confirmed that unauthenticated mode does **not work** due to an upstream bug:

```bash
# Test image extracted from OCP 4.17 release payload
IMAGE=quay.io/openshift-release-dev/ocp-v4.0-art-dev@sha256:01745f6ea0f7763b86ee62e324402aec620bb4cd5f096a10962264cb4d68cd2d

# Without CA cert - crashes immediately
$ podman run --rm $IMAGE /usr/bin/proxy-agent \
    --proxy-server-host=127.0.0.1 --proxy-server-port=8091
Error: failed to run proxy connection with failed to read CA cert : read .: is a directory

# With CA cert - runs correctly (just needs server to connect to)
$ openssl req -x509 -newkey rsa:2048 -keyout /tmp/ca.key -out /tmp/ca.crt -days 365 -nodes -subj "/CN=test-ca"
$ podman run --rm -v /tmp/ca.crt:/tmp/ca.crt:ro $IMAGE /usr/bin/proxy-agent \
    --ca-cert=/tmp/ca.crt --proxy-server-host=127.0.0.1 --proxy-server-port=8091
E0130 "cannot connect once" err="...dial tcp 127.0.0.1:8091: connect: connection refused"
```

**Conclusion**: TLS certificates are **required** for the Konnectivity agent to function. The PoC must be updated to generate and deploy certificates.

---

## Phase 2: Implementation

### Task 2.1: Generate Konnectivity Certificates at Runtime

**Status**: ⏳ Pending

**Objective**: Generate Konnectivity server and agent certificates at runtime on the bootstrap node using the existing RootCA.

**Approach**: Use `openssl` commands in `bootkube.sh` to generate certificates signed by `/opt/openshift/tls/root-ca.{crt,key}`.

**Files to modify**:

| File | Purpose |
|------|---------|
| `data/data/bootstrap/files/usr/local/bin/bootkube.sh.template` | Add certificate generation before static pod creation |

**Implementation** (add to bootkube.sh.template before Konnectivity server static pod creation):

```bash
    # Generate Konnectivity certificates signed by RootCA (1-day validity)
    KONNECTIVITY_CERT_DIR=/opt/openshift/tls/konnectivity
    mkdir -p "${KONNECTIVITY_CERT_DIR}"

    # Server certificate for agent endpoint
    openssl req -new -newkey rsa:2048 -nodes \
        -keyout "${KONNECTIVITY_CERT_DIR}/server.key" \
        -out "${KONNECTIVITY_CERT_DIR}/server.csr" \
        -subj "/CN=konnectivity-server/O=openshift"

    openssl x509 -req -in "${KONNECTIVITY_CERT_DIR}/server.csr" \
        -CA /opt/openshift/tls/root-ca.crt \
        -CAkey /opt/openshift/tls/root-ca.key \
        -CAcreateserial \
        -out "${KONNECTIVITY_CERT_DIR}/server.crt" \
        -days 1 \
        -extfile <(printf "extendedKeyUsage=serverAuth\nsubjectAltName=IP:${BOOTSTRAP_NODE_IP}")

    # Agent client certificate (shared by all agents)
    openssl req -new -newkey rsa:2048 -nodes \
        -keyout "${KONNECTIVITY_CERT_DIR}/agent.key" \
        -out "${KONNECTIVITY_CERT_DIR}/agent.csr" \
        -subj "/CN=konnectivity-agent/O=openshift"

    openssl x509 -req -in "${KONNECTIVITY_CERT_DIR}/agent.csr" \
        -CA /opt/openshift/tls/root-ca.crt \
        -CAkey /opt/openshift/tls/root-ca.key \
        -CAcreateserial \
        -out "${KONNECTIVITY_CERT_DIR}/agent.crt" \
        -days 1 \
        -extfile <(printf "extendedKeyUsage=clientAuth")

    # Copy CA cert to konnectivity directory for server to validate agents
    cp /opt/openshift/tls/root-ca.crt "${KONNECTIVITY_CERT_DIR}/ca.crt"

    # Clean up CSR files
    rm -f "${KONNECTIVITY_CERT_DIR}"/*.csr

    echo "Generated Konnectivity certificates in ${KONNECTIVITY_CERT_DIR}"
```

**Certificate details**:

| File | CN | ExtKeyUsage | Validity | SAN |
|------|-------|-------------|----------|-----|
| `server.crt` | `konnectivity-server` | `serverAuth` | 1 day | `IP:${BOOTSTRAP_NODE_IP}` |
| `agent.crt` | `konnectivity-agent` | `clientAuth` | 1 day | None |
| `ca.crt` | `root-ca` | N/A | 10 years | N/A (copy of RootCA) |

**Ordering dependency**: Certificate generation must occur **after** `BOOTSTRAP_NODE_IP` detection (Task 2.7) and **before** static pod manifest creation (Task 2.2).

---

### Task 2.1.1: Create Agent Certificate Secret

**Status**: ⏳ Pending

**Objective**: Package the agent certificate into a Kubernetes Secret for deployment to cluster nodes.

**Approach**: Generate the Secret manifest in `bootkube.sh` and write it to `manifests/` for cluster-bootstrap to apply.

**Implementation** (add to bootkube.sh.template after certificate generation):

```bash
    # Create Secret manifest for Konnectivity agent certificates
    cat > manifests/konnectivity-agent-certs.yaml <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: konnectivity-agent-certs
  namespace: kube-system
  labels:
    app: konnectivity-agent
    openshift.io/bootstrap-only: "true"
type: Opaque
data:
  tls.crt: $(base64 -w0 "${KONNECTIVITY_CERT_DIR}/agent.crt")
  tls.key: $(base64 -w0 "${KONNECTIVITY_CERT_DIR}/agent.key")
  ca.crt: $(base64 -w0 /opt/openshift/tls/root-ca.crt)
EOF
    echo "Created Konnectivity agent certificate Secret manifest"
```

**Secret contents**:

| Key | Source | Purpose |
|-----|--------|---------|
| `tls.crt` | `agent.crt` | Agent client certificate |
| `tls.key` | `agent.key` | Agent private key |
| `ca.crt` | `root-ca.crt` | CA for validating server |

**DaemonSet volume mount** (update Task 2.4):

```yaml
volumes:
- name: konnectivity-certs
  secret:
    secretName: konnectivity-agent-certs
...
volumeMounts:
- name: konnectivity-certs
  mountPath: /etc/konnectivity
  readOnly: true
```

**Agent args**:
```bash
--ca-cert=/etc/konnectivity/ca.crt
--agent-cert=/etc/konnectivity/tls.crt
--agent-key=/etc/konnectivity/tls.key
```

---

### Task 2.2: Deploy Konnectivity Server on Bootstrap

**Objective**: Add a Konnectivity server to the bootstrap environment.

**Approach**: Deploy as a **static pod** (consistent with other bootstrap components like KAS, etcd).

**Files to modify**:

| File | Purpose |
|------|---------|
| `data/data/bootstrap/files/usr/local/bin/bootkube.sh.template` | Add `KONNECTIVITY_IMAGE` variable and create static pod manifest inline |

**Image injection approach**:

The Konnectivity image is obtained using the `image_for` shell function (defined in [release-image.sh.template](data/data/bootstrap/files/usr/local/bin/release-image.sh.template)), which extracts images from the release payload.

1. Add to bootkube.sh.template (around line 49-68 with other image definitions):
   ```bash
   KONNECTIVITY_IMAGE=$(image_for apiserver-network-proxy)
   ```

2. Create the static pod manifest inline in bootkube.sh.template (after certificate generation, before the KAS render stage), using the shell variable:
   ```bash
   # Create Konnectivity server static pod manifest
   cat > /etc/kubernetes/manifests/konnectivity-server-pod.yaml <<EOF
   apiVersion: v1
   kind: Pod
   metadata:
     name: konnectivity-server
     namespace: kube-system
   spec:
     hostNetwork: true
     priorityClassName: system-node-critical
     containers:
     - name: konnectivity-server
       image: ${KONNECTIVITY_IMAGE}
       command:
       - /usr/bin/proxy-server
       args:
       - --logtostderr=true
       - --uds-name=/etc/kubernetes/bootstrap-configs/konnectivity-server.socket
       - --cluster-cert=/etc/konnectivity/server.crt
       - --cluster-key=/etc/konnectivity/server.key
       - --cluster-ca-cert=/etc/konnectivity/ca.crt
       - --agent-port=8091
       - --health-port=2041
       - --mode=http-connect
       - --proxy-strategies=destHost,defaultRoute
       volumeMounts:
       - name: config-dir
         mountPath: /etc/kubernetes/bootstrap-configs
       - name: konnectivity-certs
         mountPath: /etc/konnectivity
         readOnly: true
     volumes:
     - name: config-dir
       hostPath:
         path: /etc/kubernetes/bootstrap-configs
         type: Directory
     - name: konnectivity-certs
       hostPath:
         path: /opt/openshift/tls/konnectivity
         type: Directory
   EOF
   ```

   **Certificate arguments**:
   - `--cluster-cert`/`--cluster-key`: Server certificate for agent endpoint (mTLS)
   - `--cluster-ca-cert`: CA certificate to validate connecting agents (RootCA)

   The UDS socket is placed in `/etc/kubernetes/bootstrap-configs/` which is already mounted in the KAS pod as the `config` volume, eliminating the need for any KAS pod manifest modifications.

This follows the same pattern used for other bootstrap components where images are extracted from the release payload and injected into manifests at runtime.

**Startup ordering**: Static pods are started by kubelet in parallel. The KAS will retry connecting to the UDS socket until the Konnectivity server is ready. No explicit ordering mechanism is needed.

---

### Task 2.3: Configure Bootstrap KAS with EgressSelectorConfiguration

**Objective**: Configure the bootstrap KAS to proxy cluster-bound traffic through Konnectivity.

**Approach**: Use `--config-override-files` flag (validated in Task 1.1.1).

**Files to create/modify**:

| File | Purpose |
|------|---------|
| `data/data/bootstrap/files/opt/openshift/egress-selector-config.yaml` | EgressSelectorConfiguration override |
| `data/data/bootstrap/files/usr/local/bin/bootkube.sh.template` | Add `--config-override-files` flag to render command |

**EgressSelectorConfiguration**:
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
          udsName: "/etc/kubernetes/bootstrap-configs/konnectivity-server.socket"
```

**bootkube.sh.template modification** (around line 247-269):
```bash
/usr/bin/cluster-kube-apiserver-operator render \
  ... existing args ...
  --config-override-files=/assets/egress-selector-config.yaml \
  ...
```

**KAS pod volume mount**: No additional volume mounts are needed.

The KAS pod template already mounts `/etc/kubernetes/bootstrap-configs` as the `config` volume (see [kube-apiserver-pod.yaml](cluster-kube-apiserver-operator/bindata/bootkube/bootstrap-manifests/kube-apiserver-pod.yaml)). By placing the Konnectivity UDS socket at `/etc/kubernetes/bootstrap-configs/konnectivity-server.socket`, the KAS pod can access it through the existing volume mount.

This avoids any need for pod manifest post-processing.

---

### Task 2.4: Create Konnectivity Agent DaemonSet

**Objective**: Deploy Konnectivity agents on cluster nodes to connect back to the bootstrap server.

**Important**: The DaemonSet runs on **cluster nodes** (masters/workers), not the bootstrap node. The bootstrap node's kubelet runs in standalone mode and doesn't connect to the API server. The DaemonSet must be deployed into the cluster via the bootstrap KAS.

**Deployment mechanism**: Manifests placed in the `manifests/` directory on the bootstrap node are applied to the cluster by `cluster-bootstrap` (see [bootkube.sh.template:582](data/data/bootstrap/files/usr/local/bin/bootkube.sh.template)). This is how all cluster resources (ConfigMaps, Secrets, CRDs, operators, etc.) are initially deployed.

#### Option A: Go Asset (Recommended for Production)

Create a proper Go asset in `pkg/asset/manifests/` that generates the DaemonSet YAML. This follows the existing pattern for cluster manifests.

**Files to create**:

| File | Purpose |
|------|---------|
| `pkg/asset/manifests/konnectivity.go` | New asset that generates the DaemonSet manifest |

**Advantages**:
- Manifest generated at `openshift-install create manifests` time
- Can be inspected/modified before cluster creation
- Single templating mechanism (Go)
- Follows existing patterns in the codebase
- Easier to test

**Implementation sketch**:
```go
// KonnectivityAgent generates the Konnectivity agent DaemonSet manifest
type KonnectivityAgent struct {
    FileList []*asset.File
}

func (k *KonnectivityAgent) Dependencies() []asset.Asset {
    return []asset.Asset{
        &installconfig.InstallConfig{},
        // Need bootstrap IP from somewhere - may require new dependency
    }
}

func (k *KonnectivityAgent) Generate(ctx context.Context, dependencies asset.Parents) error {
    // Generate DaemonSet YAML with bootstrap IP and image reference
}
```

**Challenge**: The Go asset runs at `openshift-install create manifests` time, before the bootstrap node exists. We need access to:
- Bootstrap node IP (available via install config or derived from network config)
- Konnectivity image reference (available from release image)

#### Option B: Inline in bootkube.sh (Simpler for PoC)

Create the manifest at bootstrap runtime and write it to `manifests/` before cluster-bootstrap runs.

**Files to modify**:

| File | Purpose |
|------|---------|
| `data/data/bootstrap/files/usr/local/bin/bootkube.sh.template` | Add DaemonSet creation |

**Key configuration decisions**:
- **Namespace**: `kube-system` (standard for infrastructure components)
- **hostNetwork**: `true` (agent needs to reach bootstrap node IP)
- **Tolerations**: Tolerate all taints (must run on all nodes including control plane)
- **Bootstrap IP**: Use `${BOOTSTRAP_NODE_IP}` shell variable (detected at runtime, see Task 2.7)
- **Image**: Use `${KONNECTIVITY_IMAGE}` shell variable
- **Certificates**: Mount from `konnectivity-agent-certs` Secret (created in Task 2.1.1)

**Implementation in bootkube.sh.template** (after certificate Secret creation, before cluster-bootstrap):
```bash
# Create Konnectivity agent DaemonSet manifest for cluster deployment
cat > manifests/konnectivity-agent-daemonset.yaml <<EOF
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: konnectivity-agent
  namespace: kube-system
  labels:
    app: konnectivity-agent
    openshift.io/bootstrap-only: "true"
spec:
  selector:
    matchLabels:
      app: konnectivity-agent
  template:
    metadata:
      labels:
        app: konnectivity-agent
    spec:
      hostNetwork: true
      dnsPolicy: Default
      priorityClassName: system-node-critical
      tolerations:
      - operator: Exists
      containers:
      - name: konnectivity-agent
        image: ${KONNECTIVITY_IMAGE}
        command:
        - /usr/bin/proxy-agent
        args:
        - --logtostderr=true
        - --ca-cert=/etc/konnectivity/ca.crt
        - --agent-cert=/etc/konnectivity/tls.crt
        - --agent-key=/etc/konnectivity/tls.key
        - --proxy-server-host=${BOOTSTRAP_NODE_IP}
        - --proxy-server-port=8091
        - --health-server-port=2041
        - --agent-identifiers=default-route=true
        - --keepalive-time=30s
        - --probe-interval=5s
        - --sync-interval=5s
        volumeMounts:
        - name: konnectivity-certs
          mountPath: /etc/konnectivity
          readOnly: true
      volumes:
      - name: konnectivity-certs
        secret:
          secretName: konnectivity-agent-certs
EOF
```

**How this works**:
1. `${BOOTSTRAP_NODE_IP}` is detected at runtime using `ip route get` (Task 2.7)
2. `${KONNECTIVITY_IMAGE}` is resolved when bootkube.sh runs on bootstrap node
3. The manifest is written to `manifests/konnectivity-agent-daemonset.yaml`
4. `cluster-bootstrap` applies the Secret and DaemonSet to the cluster
5. The DaemonSet schedules pods on cluster nodes, which mount the certificate Secret

**Disadvantages**:
- Manifest not visible until bootstrap runtime
- Harder to test outside of a real cluster

#### Decision: Option B (Inline)

For the PoC, use **Option B (inline in bootkube.sh)** to reduce implementation complexity. Option A (Go asset) is deferred to post-PoC work.

---

### Task 2.5: Integrate Teardown with Bootstrap Completion

**Status**: ⏳ Partially Complete (design complete, implementation pending)

**Objective**: Remove all Konnectivity bootstrap-only resources when bootstrap completes.

**Approach**: Add cleanup logic to `pkg/destroy/bootstrap/bootstrap.go` (per Task 1.3).

**Resources to delete**:

| Resource | Namespace | Name | Purpose |
|----------|-----------|------|---------|
| DaemonSet | `kube-system` | `konnectivity-agent` | Agent pods on cluster nodes |
| Secret | `kube-system` | `konnectivity-agent-certs` | Agent TLS certificates |

**Files to modify**:

| File | Purpose |
|------|---------|
| `pkg/destroy/bootstrap/bootstrap.go` | Add `deleteKonnectivityResources()` function |

**Implementation**:
```go
func deleteKonnectivityResources(ctx context.Context, dir string) error {
    // Load kubeconfig from asset directory
    kubeconfigPath := filepath.Join(dir, "auth", "kubeconfig")
    config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
    if err != nil {
        return fmt.Errorf("failed to load kubeconfig: %w", err)
    }

    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        return fmt.Errorf("failed to create Kubernetes client: %w", err)
    }

    namespace := "kube-system"

    // Delete DaemonSet
    err = clientset.AppsV1().DaemonSets(namespace).Delete(ctx, "konnectivity-agent", metav1.DeleteOptions{})
    if err != nil && !apierrors.IsNotFound(err) {
        logrus.Warnf("Failed to delete Konnectivity agent DaemonSet: %v", err)
    } else if err == nil {
        logrus.Info("Deleted Konnectivity agent DaemonSet")
    }

    // Delete Secret
    err = clientset.CoreV1().Secrets(namespace).Delete(ctx, "konnectivity-agent-certs", metav1.DeleteOptions{})
    if err != nil && !apierrors.IsNotFound(err) {
        logrus.Warnf("Failed to delete Konnectivity agent certs Secret: %v", err)
    } else if err == nil {
        logrus.Info("Deleted Konnectivity agent certs Secret")
    }

    return nil
}
```

**Timing**: Called after `provider.DestroyBootstrap()` returns successfully.

**Error handling**: Failures are logged as warnings but do not fail the overall teardown. The bootstrap infrastructure is already destroyed at this point, so these resources are orphaned and can be manually cleaned up if needed.

---

### Task 2.6: Ignition Config Validation

**Status**: ⏳ Needs Re-validation (prior validation was for unauthenticated mode; now includes mTLS)

**Objective**: Validate implementation without deploying a cluster.

#### Prerequisites

Build the installer with the Konnectivity changes:

```bash
SKIP_TERRAFORM=y hack/build.sh
```

#### Test Setup

Create a test directory with a minimal install-config.yaml. The "none" platform is used because it doesn't require cloud provider credentials or infrastructure validation:

```bash
# Create test directory
mkdir -p /tmp/konnectivity-test
cd /tmp/konnectivity-test

# Generate an SSH key for the test
ssh-keygen -t ed25519 -f /tmp/konnectivity-test/test-key -N "" -q

# Create install-config.yaml
cat > install-config.yaml << EOF
apiVersion: v1
baseDomain: example.com
metadata:
  name: test-cluster
networking:
  networkType: OVNKubernetes
  clusterNetwork:
  - cidr: 10.128.0.0/14
    hostPrefix: 23
  serviceNetwork:
  - 172.30.0.0/16
  machineNetwork:
  - cidr: 192.168.1.0/24
controlPlane:
  name: master
  replicas: 3
compute:
- name: worker
  replicas: 2
platform:
  none: {}
pullSecret: '{"auths":{"fake":{"auth":"aWQ6cGFzcwo="}}}'
sshKey: '$(cat /tmp/konnectivity-test/test-key.pub)'
EOF
```

#### Generate Ignition Configs

```bash
cd /tmp/konnectivity-test
/path/to/installer/bin/openshift-install create ignition-configs
```

Expected output:
```
level=warning msg=Release Image Architecture not detected. Release Image Architecture is unknown
level=info msg=Consuming Install Config from target directory
level=info msg=Successfully populated MCS CA cert information: root-ca ...
level=info msg=Successfully populated MCS TLS cert information: root-ca ...
level=info msg=Ignition-Configs created in: . and auth
```

#### Validation Checklist

After generating ignition configs, run these validation commands:

##### 1. EgressSelectorConfiguration file exists

```bash
jq -r '.storage.files[] | select(.path == "/opt/openshift/egress-selector-config.yaml") | .contents.source' \
  bootstrap.ign | sed 's/data:text\/plain;charset=utf-8;base64,//' | base64 -d
```

**Expected output**:
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
          udsName: "/etc/kubernetes/bootstrap-configs/konnectivity-server.socket"
```

##### 2. KONNECTIVITY_IMAGE variable defined

```bash
jq -r '.storage.files[] | select(.path == "/usr/local/bin/bootkube.sh") | .contents.source' \
  bootstrap.ign | sed 's/data:text\/plain;charset=utf-8;base64,//' | base64 -d | grep -A2 "KONNECTIVITY_IMAGE"
```

**Expected output**:
```
KONNECTIVITY_IMAGE=$(image_for apiserver-network-proxy)

mkdir --parents ./{bootstrap-manifests,manifests}
--
    image: ${KONNECTIVITY_IMAGE}
    command:
    - /usr/bin/proxy-server
--
        image: ${KONNECTIVITY_IMAGE}
        command:
        - /usr/bin/proxy-agent
```

##### 3. Konnectivity server static pod manifest creation

```bash
jq -r '.storage.files[] | select(.path == "/usr/local/bin/bootkube.sh") | .contents.source' \
  bootstrap.ign | sed 's/data:text\/plain;charset=utf-8;base64,//' | base64 -d | grep -A40 "konnectivity-server-pod.yaml" | head -45
```

**Expected output** (partial):
```yaml
    cat > /etc/kubernetes/manifests/konnectivity-server-pod.yaml <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: konnectivity-server
  namespace: kube-system
  labels:
    app: konnectivity-server
spec:
  hostNetwork: true
  priorityClassName: system-node-critical
  containers:
  - name: konnectivity-server
    image: ${KONNECTIVITY_IMAGE}
    command:
    - /usr/bin/proxy-server
    args:
    - --logtostderr=true
    - --uds-name=/etc/kubernetes/bootstrap-configs/konnectivity-server.socket
    - --agent-port=8091
    - --health-port=2041
    - --mode=http-connect
    - --proxy-strategies=destHost,defaultRoute
    ...
```

##### 4. Konnectivity agent DaemonSet manifest creation

```bash
jq -r '.storage.files[] | select(.path == "/usr/local/bin/bootkube.sh") | .contents.source' \
  bootstrap.ign | sed 's/data:text\/plain;charset=utf-8;base64,//' | base64 -d | grep -A55 "konnectivity-agent-daemonset.yaml" | head -55
```

**Expected output** (partial):
```yaml
    cat > manifests/konnectivity-agent-daemonset.yaml <<EOF
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: konnectivity-agent
  namespace: kube-system
  labels:
    app: konnectivity-agent
    openshift.io/bootstrap-only: "true"
spec:
  selector:
    matchLabels:
      app: konnectivity-agent
  ...
    spec:
      hostNetwork: true
      dnsPolicy: Default
      priorityClassName: system-node-critical
      tolerations:
      - operator: Exists
      containers:
      - name: konnectivity-agent
        image: ${KONNECTIVITY_IMAGE}
        command:
        - /usr/bin/proxy-agent
        args:
        - --logtostderr=true
        - --proxy-server-host=
        - --proxy-server-port=8091
        ...
```

**Note**: The `--proxy-server-host=` is empty because the "none" platform doesn't have an installer-managed bootstrap IP. On IPI platforms (AWS, GCP, Azure, etc.), this would contain the actual bootstrap node IP address via the `{{.BootstrapNodeIP}}` template variable.

##### 5. KAS render command includes config override

```bash
jq -r '.storage.files[] | select(.path == "/usr/local/bin/bootkube.sh") | .contents.source' \
  bootstrap.ign | sed 's/data:text\/plain;charset=utf-8;base64,//' | base64 -d | grep "config-override-files"
```

**Expected output**:
```
		--config-override-files=/assets/egress-selector-config.yaml
```

#### Cleanup

```bash
rm -rf /tmp/konnectivity-test
```

#### Validation Results Summary

| Check | Status | Notes |
|-------|--------|-------|
| EgressSelectorConfiguration file | ✅ Pass | Correctly embedded at `/opt/openshift/egress-selector-config.yaml` |
| KONNECTIVITY_IMAGE variable | ✅ Pass | Defined and used in both server and agent manifests |
| Konnectivity server static pod | ✅ Pass | Complete manifest with UDS socket and correct args |
| Konnectivity agent DaemonSet | ✅ Pass | Complete manifest with hostNetwork, tolerations |
| --config-override-files flag | ✅ Pass | Added to KAS render command |

#### Limitations

- Runtime behavior (shell script execution) not tested
- Image extraction from release payload not tested
- Actual Konnectivity connectivity not validated
- Bootstrap node IP only populated on IPI platforms with infrastructure provisioning

---

### Task 2.7: Fix Bootstrap IP Detection for Konnectivity Agent

**Status**: ⏳ Pending

**Objective**: Replace the Go template variable `{{.BootstrapNodeIP}}` with runtime shell IP detection.

**Problem**: The `{{.BootstrapNodeIP}}` template variable is only populated from the `OPENSHIFT_INSTALL_BOOTSTRAP_NODE_IP` environment variable, which is not set for IPI deployments where the bootstrap VM IP is assigned dynamically. See Task 1.8 "Issue B" for detailed root cause analysis.

**Solution**: Detect the bootstrap node IP at runtime using `ip route get` when the DaemonSet manifest is generated in bootkube.sh.

**Files to modify**:

| File | Change |
|------|--------|
| `data/data/bootstrap/files/usr/local/bin/bootkube.sh.template` | Add IP detection before DaemonSet heredoc |

**Implementation**:

Before the DaemonSet manifest heredoc (around line 631), add:

```bash
    # Detect bootstrap node IP at runtime using the default route source address
{{- if .UseIPv6ForNodeIP }}
    BOOTSTRAP_NODE_IP=$(ip -6 route get 2001:4860:4860::8888 2>/dev/null | awk '{for(i=1;i<=NF;i++) if($i=="src") print $(i+1); exit}')
{{- else }}
    BOOTSTRAP_NODE_IP=$(ip route get 1.1.1.1 2>/dev/null | awk '{for(i=1;i<=NF;i++) if($i=="src") print $(i+1); exit}')
{{- end }}
    echo "Detected bootstrap node IP: ${BOOTSTRAP_NODE_IP}"
```

Change line 665 from:
```yaml
        - --proxy-server-host={{.BootstrapNodeIP}}
```

To:
```yaml
        - --proxy-server-host=${BOOTSTRAP_NODE_IP}
```

**Why `api-int` DNS cannot be used**:
- `api-int.<cluster-domain>` resolves to a VIP (baremetal) or load balancer (cloud)
- The VIP can move to control plane nodes via VRRP before bootstrap teardown
- Load balancers distribute traffic to any healthy backend, including control plane nodes
- Konnectivity agents might connect to a node without the Konnectivity server

**Verification**:

1. Generate ignition configs with the modified template
2. Decode bootkube.sh from bootstrap.ign and verify:
   - IP detection command is present with correct IPv4/IPv6 selection
   - DaemonSet manifest uses `${BOOTSTRAP_NODE_IP}` shell variable
3. Deploy a test cluster and verify:
   - Bootstrap node logs show detected IP
   - Konnectivity agents on cluster nodes connect successfully

---

## Open Questions

### Resolved

| Question | Answer | Source |
|----------|--------|--------|
| Which port should the Konnectivity server listen on? | UDS socket for KAS, TCP port 8091 for agents | Task 1.5 (HyperShift), Task 1.6 |
| Should the Konnectivity server be a static pod or systemd service? | Static pod (consistent with KAS, etcd) | Task 1.1 |
| What is the expected lifecycle overlap between bootstrap and cluster control plane? | Pod network is operational before bootstrap completes; production KAS instances are running | Task 1.2, Task 1.3 |

### Remaining Uncertainties

These should be investigated during implementation if they become blocking:

1. **Bootstrap node IP discovery** ✅ Resolved (Updated)
   - **Issue**: The agent DaemonSet needs to know the bootstrap node's IP address to connect to the Konnectivity server.
   - **Original Solution** ❌: Use `{{.BootstrapNodeIP}}` Go template variable. This doesn't work because the variable is only populated from the `OPENSHIFT_INSTALL_BOOTSTRAP_NODE_IP` environment variable, which isn't set for IPI deployments.
   - **Updated Solution** ✅: Detect the bootstrap IP at **shell runtime** using `ip route get`. See Task 1.8 "Issue B" for root cause analysis and Task 2.7 for implementation details.

2. **KAS pod volume mount for UDS socket** ✅ Resolved
   - **Issue**: The bootstrap KAS pod needs access to the Konnectivity UDS socket.
   - **Solution**: Place the socket at `/etc/kubernetes/bootstrap-configs/konnectivity-server.socket`. This directory is already mounted in the KAS pod as the `config` volume, so no additional volume mounts or manifest modifications are needed.

3. **Static pod startup timing**
   - **Issue**: If KAS starts before Konnectivity server, initial webhook calls may fail until the server is ready.
   - **Mitigation**: KAS should retry connections to the UDS socket. Webhook calls during early bootstrap may fail and retry.
   - **Risk**: Low - the KAS has retry logic, and webhooks are not typically needed in the earliest bootstrap stages.

4. **Agent connectivity before CNI**
   - **Issue**: Agents run in hostNetwork mode and connect to bootstrap via infrastructure network. This should work, but needs validation.
   - **Mitigation**: Agents use hostNetwork (bypasses CNI) and connect via bootstrap node's infrastructure IP.
   - **Risk**: Low - this is the same network path used by kubelets to reach the bootstrap KAS.

---

## Manual Testing

This section describes manual testing procedures for validating the Konnectivity integration on a real cluster.

### Obtaining the Bootstrap IP

The bootstrap node's public IP can be extracted from the CAPI Machine manifest stored in `.clusterapi_output/`:

```bash
# Find the bootstrap machine manifest and extract the external IP
yq '.status.addresses[] | select(.type == "ExternalIP") | .address' \
  .clusterapi_output/Machine-openshift-cluster-api-guests-*-bootstrap.yaml
```

Alternatively, view all addresses:
```bash
cat .clusterapi_output/Machine-openshift-cluster-api-guests-*-bootstrap.yaml | grep -A10 "addresses:"
```

Note: The `openshift-install gather bootstrap` command will automatically extract this IP if the `--bootstrap` flag is not provided.

### SSH Access to Bootstrap Node

The installer always generates a bootstrap SSH key pair (see [bootstrapsshkeypair.go](pkg/asset/tls/bootstrapsshkeypair.go)). The public key is embedded in the bootstrap ignition and added to the `core` user's authorized_keys alongside any user-provided SSH key.

**Option 1**: Use user-provided SSH key (if configured in install-config.yaml):
```bash
ssh -i ~/.ssh/<your-key> core@<bootstrap-ip>
```

**Option 2**: Extract the generated key from the installer state file:
```bash
# Extract the private key from .openshift_install_state.json
jq -r '.["*tls.BootstrapSSHKeyPair"]["Priv"]' .openshift_install_state.json | base64 -d > bootstrap-ssh.key
chmod 600 bootstrap-ssh.key

# SSH to bootstrap node
ssh -i bootstrap-ssh.key core@<bootstrap-ip>
```

Note: The `openshift-install gather bootstrap` command automatically uses the generated key, so manual extraction is only needed for interactive debugging sessions.

### Bootstrap Service Debugging

If bootstrap changes cause issues, use the service recording mechanism documented in [bootstrap_services.md](docs/dev/bootstrap_services.md):

- Progress tracked in `/var/log/openshift/<service>.json`
- Use `record_service_stage_start` / `record_service_stage_success` functions in shell scripts
- Failed-units captured in log bundle via `openshift-install gather bootstrap`

### Log Gathering on Failure

```bash
openshift-install gather bootstrap --bootstrap <IP> --master <IP1> --master <IP2> --master <IP3>
```

The log bundle includes:
- `bootstrap/journals/bootkube.log` - Main bootstrap orchestration
- `bootstrap/containers/` - Container logs and descriptions
- `bootstrap/pods/` - Render command outputs

### End-to-End Validation Steps

1. Deploy a cluster with the modified installer
2. SSH to the bootstrap node and verify:
   - Konnectivity server static pod is running: `crictl ps | grep konnectivity-server`
   - UDS socket exists: `ls -la /etc/kubernetes/bootstrap-configs/konnectivity-server.socket`
   - Server is listening on agent port: `ss -tlnp | grep 8091`
3. Verify Konnectivity agent pods are running on cluster nodes:
   ```bash
   kubectl get pods -n kube-system -l app=konnectivity-agent
   ```
4. Check agent connectivity to server (from bootstrap node):
   ```bash
   # Check server logs for agent connections
   crictl logs $(crictl ps -q --name konnectivity-server)
   ```
5. Deploy a test webhook in the cluster (e.g., a simple validating webhook)
6. Confirm the bootstrap KAS can reach the webhook (create a resource that triggers the webhook)
7. Complete bootstrap teardown and verify Konnectivity resources are removed:
   ```bash
   # Both should return "not found" after teardown
   kubectl get daemonset -n kube-system konnectivity-agent
   kubectl get secret -n kube-system konnectivity-agent-certs
   ```
8. Confirm production KAS on cluster nodes can still reach webhooks normally

---

## Post-PoC Considerations

The following items are out of scope for the initial proof-of-concept but should be addressed for a production-ready implementation:

### 1. UDS Socket Location

**Current approach**: The UDS socket is placed at `/etc/kubernetes/bootstrap-configs/konnectivity-server.socket` to reuse an existing volume mount in the KAS pod.

**Concerns**:
- This directory is intended for configuration files, not runtime sockets
- The socket persists in the config directory after Konnectivity server terminates
- May cause confusion during debugging or log analysis

**Post-PoC options**:
- Create a dedicated `/etc/kubernetes/konnectivity/` directory and add proper volume mounts to the KAS pod via operator changes
- Use a runtime directory like `/run/konnectivity/` with appropriate volume mounts
- Coordinate with cluster-kube-apiserver-operator maintainers on the preferred approach

### 2. Konnectivity Agent DaemonSet as Go Asset

**Current approach**: The DaemonSet manifest is created inline in bootkube.sh.template using mixed Go/shell templating.

**Concerns**:
- Manifest not visible until bootstrap runtime
- Mixed templating (Go + shell) is confusing and error-prone
- Harder to test and validate
- Doesn't follow established patterns for cluster manifests

**Post-PoC implementation**:
- Create `pkg/asset/manifests/konnectivity.go` as a proper Go asset
- Generate DaemonSet YAML at `openshift-install create manifests` time
- Manifest can be inspected/modified before cluster creation
- Follows existing patterns (similar to other manifests in `pkg/asset/manifests/`)
- Requires solving bootstrap IP availability at manifest generation time

### 3. ~~Agent-to-Server Authentication (mTLS)~~ ✅ Now Implemented

**Status**: This is now implemented in the PoC via Task 2.1 and Task 2.1.1.

The implementation uses the existing `RootCA` to sign both server and agent certificates at runtime on the bootstrap node. Certificates have 1-day validity, appropriate for the temporary bootstrap phase.

