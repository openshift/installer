# Launch

## Goals

1. Validate credentials for cloud.
2. Create a functional cluster.
3. Validate the cluster is ready to be used by an user.

## Overview

The launch phase creates a cluster based on the assets generated in the prepare phase. The launch performs operations,

1. Pre-flight checks: Verify that credentials being passed to the installer are sufficient to bring up the cluster.
2. Launch all the platform specific resources required for the cluster. For example in AWS this includes, DNS in Route53, VPCs, ELBs, IAM roles, Security groups etc.
3. Bootstrap the cluster.
4. Wait for the cluster to be ready for use by an user.

## Detailed design

### Idempotency

Launch is **NOT** idempotent. Re-running the `installer launch` command should exit with a failure.

### Pre flight checks

### Platform Specific checks

1. AWS

    The launch phase performs following checks

    * Credentials are sufficient to create all the resources.
    * If VPC / Route53 zones were supplied, ensure they are valid.
    * TODO: add more

2. Libvirt

    The launch phase performs following checks

    * Ensure QEMU URI is reachable
    * OS image path is valid.
    * TODO: add more

### Launch platform-specific resources

Use terraform to create resources.

### Bootstrapping cluster

**[Various options are discussed here](https://docs.google.com/document/d/17sTJ1mdWtPTFkaHLENX2aeEYw1-o4aDS2NbWnZRgOlY/edit#heading=h.r9how0eg6txs)** *This link might be private*

The goal is to eliminate external coordination steps and keep the "special" steps constrained to a single throw-away node.

a. Launch a bootstrap node from `bootstrap.ign` that was generated in prepare step. The user-data of the bootstrap node either contains the `bootstrap.ign` or points to a remote location that contains the `bootstrap.ign` file.

b. Launch 3 master nodes with ign endpoint as `api.example.com`

c. Launch ELB or equivalent resource that fronts bootstrap node, and masters

d. ALIAS `api.example.com` → ELB

e. Start bootstrap MachineConfig Server to serve ignition for the masters. Also start the `etcd-signer-server` for serving TLS assets for the etcd members.

f. When etcd on masters has formed quorum, stop local MachineConfig Server and the `etcd-signer-server`.

g. Run bootkube to launch self-hosted control-plane.

    * Bootkube starts `boostrap-*` static pods on boostrap node to create a temporary control plane.
    * Then bootkube uses the temporary control plane to bootstrap the self-hosted control plane.
    * When bootkube validates that the self-hosted plane is ready, it shuts down.
    * When bootkube shuts down, it tears down bootstrap-apiserver (this triggers fail on ELB healthcheck, and bootstrap node is removed from backend pool).

h. Destroy bootstrap node and destroy the remote location if ign for bootstrap node was served from remote endpoint as it stores secrets.

### Bootstrapping etcd

Etcd is co-located with master nodes. Therefore the bootstrap MachineConfig Server serves the ignition file with the etcd static pods. But the etcd nodes need TLS assets to communicate with each other.

1. The boostrap node runs a [`etcd-signer-server`](https://github.com/coreos/kubecsr/tree/master/cmd/kube-etcd-signer-server) docker container which mimics the kube-apiserver's `CertificateSigningRequest` endpoint.
2. etcd systemd-service has a `PreStartHook` defined that runs a [`etcd-client`](https://github.com/coreos/kubecsr/tree/master/cmd/kube-client-agent). The `etcd-client` reaches out to the API server endpoint, currently being served by `etcd-signer-server`, for certificates.
3. After the `PreStartHook` succeeds, each etcd member has all the TLS assets for creating the etcd cluster.

### Verify cluster is ready

The launch step needs to exit only when the cluster is ready for use by an user.

TODO: need more info to decide when cluster is up.

a. Is installer done when control-plane is ready.
b. Installer is ready when all the second-level operators report `Done` condition.