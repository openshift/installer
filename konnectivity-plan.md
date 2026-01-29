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

- [ ] The bootstrap KAS will be able to access cluster-hosted webhooks via Konnectivity
- [ ] Non-bootstrap KAS instances will not be impacted and will continue routing normally

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
      proxyProtocol: "CONNECT"
      transport:
        type: "TCP"
        tcp:
          url: "http://127.0.0.1:<konnectivity-port>"
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
         proxyProtocol: "CONNECT"
         transport:
           type: "TCP"
           tcp:
             url: "http://127.0.0.1:{{.KonnectivityPort}}"
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

**Deliverables**:
- Document the network topology of the bootstrap environment
- Explain why pod network connectivity is unavailable from the bootstrap node
- Identify any existing documentation on bootstrap networking constraints

---

### Task 1.3: Investigate Bootstrap Teardown Mechanism

**Objective**: Find the existing mechanism for cleaning up bootstrap resources and determine how to integrate Konnectivity agent removal.

**Deliverables**:
- Document the bootstrap teardown/completion flow
- Identify where cluster resources are cleaned up when bootstrap completes
- Determine how to add the Konnectivity DaemonSet removal to this flow
- List relevant source files

---

### Task 1.4: Investigate Konnectivity in OpenShift Payload

**Objective**: Determine if Konnectivity components are already available in the OpenShift payload.

**Deliverables**:
- Confirm whether Konnectivity server and agent images are in the payload
- If present, document the image references and how to obtain them
- If not present, document alternative sources for Konnectivity binaries/images

---

### Task 1.5: Investigate HyperShift Konnectivity Deployment

**Objective**: Examine how HyperShift deploys Konnectivity to inform our implementation.

**Deliverables**:
- Document where HyperShift obtains Konnectivity components
- Extract relevant configuration patterns (server config, agent config, certificates)
- Identify any reusable manifests or configuration templates
- Note any authentication/certificate setup between agent and server

**Source**: `/home/mbooth/src/openshift/hypershift`

---

### Task 1.6: Research Konnectivity Authentication

**Objective**: Determine how the Konnectivity agent authenticates to the Konnectivity server.

**Deliverables**:
- Document the authentication mechanism (mTLS, tokens, etc.)
- Identify what certificates or credentials need to be generated
- Determine how to provision these credentials in the bootstrap environment
- Reference upstream Konnectivity documentation as needed

---

### Task 1.7: Validate Assumptions via Documentation

**Objective**: Verify the architectural assumptions through documentation research.

**Deliverables**:
- Confirm that EgressSelectorConfiguration can route webhook traffic through Konnectivity
- Confirm that non-bootstrap KAS instances (without EgressSelectorConfiguration) route normally
- Document any caveats or limitations found

**Sources**:
- Kubernetes EgressSelectorConfiguration documentation
- Konnectivity project documentation
- Existing OpenShift/HyperShift design docs

---

## Phase 2: Implementation

### Task 2.1: Generate Konnectivity Certificates

**Objective**: Create the necessary certificates for Konnectivity server and agent authentication.

**Deliverables**:
- Add certificate generation to the installer's asset pipeline
- Generate server certificate/key
- Generate agent certificate/key (or CA for agent authentication)
- Ensure certificates are included in appropriate Ignition configs

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
