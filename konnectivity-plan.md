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

**Deliverables**:
- Document which assets create the bootstrap machine configuration
- Identify where the bootstrap KAS is configured
- Locate the bootstrap Ignition config generation code
- List the relevant source files and their responsibilities

**Search starting points**:
- `pkg/asset/` directory
- Bootstrap-related assets and ignition configs

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
