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

## Phase 2: Implementation

### Task 2.1: Generate Konnectivity Certificates

**Status**: ⏭️ Skipped (out of scope for PoC)

**Objective**: Create the necessary certificates for Konnectivity server and agent authentication.

**Note**: Per Task 1.6 decision, the PoC uses unauthenticated TCP connections. Certificate generation is deferred to future production-ready work.

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

2. Create the static pod manifest inline in bootkube.sh.template (before the KAS render stage), using the shell variable:
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
       volumeMounts:
       - name: config-dir
         mountPath: /etc/kubernetes/bootstrap-configs
     volumes:
     - name: config-dir
       hostPath:
         path: /etc/kubernetes/bootstrap-configs
         type: Directory
   EOF
   ```

   The socket is placed in `/etc/kubernetes/bootstrap-configs/` which is already mounted in the KAS pod as the `config` volume, eliminating the need for any KAS pod manifest modifications.

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
- **Bootstrap IP**: Use `{{.BootstrapNodeIP}}` Go template variable
- **Image**: Use `${KONNECTIVITY_IMAGE}` shell variable

**Implementation in bootkube.sh.template** (before cluster-bootstrap, after manifest collection):
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
      tolerations:
      - operator: Exists
      containers:
      - name: konnectivity-agent
        image: ${KONNECTIVITY_IMAGE}
        command:
        - /usr/bin/proxy-agent
        args:
        - --logtostderr=true
        - --proxy-server-host={{.BootstrapNodeIP}}
        - --proxy-server-port=8091
        - --health-server-port=2041
        - --agent-identifiers=default-route=true
        - --keepalive-time=30s
        - --probe-interval=5s
        - --sync-interval=5s
EOF
```

**How this works**:
1. `{{.BootstrapNodeIP}}` is resolved when bootkube.sh.template → bootkube.sh (ignition generation)
2. `${KONNECTIVITY_IMAGE}` is resolved when bootkube.sh runs on bootstrap node
3. The manifest is written to `manifests/konnectivity-agent-daemonset.yaml`
4. `cluster-bootstrap` applies all manifests in `manifests/` to the cluster
5. The DaemonSet is created in the cluster and schedules pods on cluster nodes

**Disadvantages**:
- Manifest not visible until bootstrap runtime
- Mixed templating (Go + shell) is confusing
- Harder to test

#### Decision: Option B (Inline)

For the PoC, use **Option B (inline in bootkube.sh)** to reduce implementation complexity. Option A (Go asset) is deferred to post-PoC work.

---

### Task 2.5: Integrate Teardown with Bootstrap Completion

**Objective**: Remove the Konnectivity agent DaemonSet when bootstrap completes.

**Approach**: Add cleanup logic to `pkg/destroy/bootstrap/bootstrap.go` (per Task 1.3).

**Files to modify**:

| File | Purpose |
|------|---------|
| `pkg/destroy/bootstrap/bootstrap.go` | Add `deleteKonnectivityResources()` function |

**Implementation**:
```go
func deleteKonnectivityResources(ctx context.Context, dir string) error {
    // Load kubeconfig from asset directory
    // Create Kubernetes client
    // Delete DaemonSet kube-system/konnectivity-agent
    // Log warning on failure (don't fail the teardown)
}
```

**Timing**: Called after `provider.DestroyBootstrap()` returns successfully.

---

### Task 2.6: Ignition Config Validation

**Status**: ✅ Complete

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

## Open Questions

### Resolved

| Question | Answer | Source |
|----------|--------|--------|
| Which port should the Konnectivity server listen on? | UDS socket for KAS, TCP port 8091 for agents | Task 1.5 (HyperShift), Task 1.6 |
| Should the Konnectivity server be a static pod or systemd service? | Static pod (consistent with KAS, etcd) | Task 1.1 |
| What is the expected lifecycle overlap between bootstrap and cluster control plane? | Pod network is operational before bootstrap completes; production KAS instances are running | Task 1.2, Task 1.3 |

### Remaining Uncertainties

These should be investigated during implementation if they become blocking:

1. **Bootstrap node IP discovery** ✅ Resolved
   - **Issue**: The agent DaemonSet needs to know the bootstrap node's IP address to connect to the Konnectivity server.
   - **Solution**: Use `{{.BootstrapNodeIP}}` Go template variable, which is already available in `bootstrapTemplateData` (see [common.go:94](pkg/asset/ignition/bootstrap/common.go)). This is the same variable used by `kubelet.sh.template` for the `--node-ip` flag.

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
7. Complete bootstrap teardown and verify agent DaemonSet is removed:
   ```bash
   kubectl get daemonset -n kube-system konnectivity-agent
   # Should return "not found" after teardown
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

### 3. Agent-to-Server Authentication (mTLS)

**Current approach**: Agents connect to the Konnectivity server over unauthenticated TCP on port 8091.

**Security risks**:
- Any process on the infrastructure network could connect to the Konnectivity server
- Potential for unauthorized tunnel access during the bootstrap window
- Does not meet production security requirements

**Post-PoC implementation** (see Task 2.1 for deferred work):
- Generate a Konnectivity CA certificate as part of the installer's TLS asset pipeline
- Generate server certificate for the Konnectivity server (agent endpoint)
- Generate agent certificates or use a shared CA for agent authentication
- Configure mTLS on both server (`--cluster-cert`, `--cluster-key`, `--server-ca-cert`) and agents (`--ca-cert`, `--agent-cert`, `--agent-key`)
- Reference HyperShift's implementation in `control-plane-operator/controllers/hostedcontrolplane/pki/konnectivity.go`
