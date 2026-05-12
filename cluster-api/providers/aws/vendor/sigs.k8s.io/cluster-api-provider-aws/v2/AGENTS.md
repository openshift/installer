# Agent Development Guide

This guide is the primary instruction file for AI coding agents working in this
repository. Read [AI_POLICY.md](AI_POLICY.md) as well — it defines the
project-level policy that governs every AI-assisted contribution (human
oversight, quality, transparency, DCO sign-off).

## Rules & Constraints

### Process

- Before writing code, restate the goal, list the files you intend to touch, and
  wait for confirmation on anything ambiguous.
- Work iteratively. Keep each change as small as it can be while still making
  sense on its own.
- Do not open or submit pull requests automatically — leave that to the human.
- Commits must be signed off (`git commit -s`) to satisfy DCO. Follow the
  existing commit style in `git log`.

### Code

- Add unit tests for every change. Major new features also need new e2e tests.
- `make lint` and `make test` must be clean before handoff.
- Do not manually edit generated files. Regenerate with `make generate`.
  Generated files include `zz_generated.*.go`, CRDs under `config/crd/bases/`,
  RBAC under `config/rbac/`, and anything else emitted by `make generate`.
- All new API work targets `v1beta2`. `v1beta1` is deprecated — never add
  fields to it. When a `v1beta2` field is added or changed, update the
  corresponding conversion in `api/v1beta1/*_conversion.go` (and the equivalent
  under `bootstrap/eks/api/v1beta1/` or `controlplane/eks/api/v1beta1/`).
- Do not update the Go version, Kubernetes, or controller-runtime dependencies
  in `go.mod`.
- Follow the existing "scopes and services" pattern (see
  [pkg/cloud/scope/cluster.go](pkg/cloud/scope/cluster.go) and
  [pkg/cloud/services/ec2/instances.go](pkg/cloud/services/ec2/instances.go)
  as reference implementations).

## Project Context

This project implements a [Cluster API](https://github.com/kubernetes-sigs/cluster-api)
(CAPI) provider for provisioning Kubernetes clusters in AWS. It supports pure
EC2-based clusters (unmanaged), EKS clusters (managed), and ROSA (experimental,
under `/exp`). It implements the CAPI **infrastructure**, **bootstrap**, and
**control plane** provider contracts.

## Where Things Live

| Task | Where to go |
| --- | --- |
| Add / change a core CR field | [api/v1beta2/](api/v1beta2/) + `make generate` + update [api/v1beta1/*_conversion.go](api/v1beta1/) |
| Add / change an experimental CR field (MachinePool, ROSA, Fargate) | [exp/api/v1beta2/](exp/api/v1beta2/) + `make generate` |
| Add / change an EKS bootstrap CR field | [bootstrap/eks/api/v1beta2/](bootstrap/eks/api/v1beta2/) + `make generate` |
| Add / change an EKS control-plane CR field | [controlplane/eks/api/v1beta2/](controlplane/eks/api/v1beta2/) + `make generate` |
| Modify core reconciliation logic | [controllers/](controllers/) |
| Modify experimental reconciliation logic | [exp/controllers/](exp/controllers/) |
| Add a new AWS API call | [pkg/cloud/services/&lt;svc&gt;/](pkg/cloud/services/) (one package per AWS service) |
| Add cluster/machine context or an AWS client | [pkg/cloud/scope/](pkg/cloud/scope/) |
| Gate new functionality behind a feature flag | [feature/feature.go](feature/feature.go) (gates default to off) |
| Add an e2e test | [test/e2e/](test/e2e/) |

## Architecture

### Core Components

1. **Infrastructure Provider API Definitions** — [api/](api/) and [exp/api/](exp/api/).
   Key resources: `AWSCluster`, `AWSMachine`, `AWSClusterTemplate`,
   `AWSMachineTemplate`, and experimental `AWSMachinePool`, `AWSFargateProfile`.
2. **Infrastructure Provider Controllers** — [controllers/](controllers/) and
   [exp/controllers/](exp/controllers/). Key controllers: `AWSClusterReconciler`,
   `AWSMachineReconciler`, `AWSMachinePoolReconciler`.
3. **EKS Bootstrap Provider API Definitions** — [bootstrap/eks/api/](bootstrap/eks/api/).
   Key resources: `EKSConfig`, `EKSConfigTemplate`.
4. **EKS Bootstrap Provider Controllers** — [bootstrap/eks/controllers/](bootstrap/eks/controllers/).
   Key controller: `EKSConfigReconciler`.
5. **EKS Control Plane Provider API Definitions** — [controlplane/eks/api/](controlplane/eks/api/).
   Key resource: `AWSManagedControlPlane`. ROSA lives under [exp/](exp/).
6. **EKS Control Plane Provider Controllers** — [controlplane/eks/controllers/](controlplane/eks/controllers/).
   Key controller: `AWSManagedControlPlaneReconciler`.
7. **Services Layer** — [pkg/cloud/services/](pkg/cloud/services/). AWS
   service-specific clients organised by functional area (`ec2`, `s3`, `eks`,
   `elb`, …). Each service exposes `Reconcile*` and `Delete*` entry points
   called from the controllers.
8. **Scope Package** — [pkg/cloud/scope/](pkg/cloud/scope/). Scopes encapsulate
   cluster/machine specs, credentials, and AWS clients. Key scopes:
   `ClusterScope`, `MachineScope`, `ManagedControlPlaneScope`.
9. **Feature Gates** — [feature/feature.go](feature/feature.go). Controls
   experimental/optional functionality (`MachinePool`, `EKSFargate`,
   `ROSA`, …). All gates default to off.

All controllers use the same reconciliation pattern:
**observe state → determine actions → apply changes → update status**.

### Data Flow

```
User creates K8s resource → Controller watches → Reconciler triggered →
Scope created → Service methods called → AWS API interactions →
Status updated → Requeue if needed
```

## Commands

### Build

```bash
make binaries          # build all binaries
make docker-build      # build controller image for current arch
make docker-build-all  # build controller image for all supported archs
```

### Test & Lint

```bash
make test              # unit tests
make test-e2e          # full e2e suite (slow; requires AWS credentials)
make test-e2e-eks      # EKS-only e2e suite
make lint              # golangci-lint
```

### Generate

```bash
make generate          # regenerate CRDs, deepcopy, conversions, RBAC, mocks
```

## Pre-Handoff Checklist

Before declaring a task done, run and pass:

```bash
make generate && make lint && make test
```

Then sanity-check:

- No generated files edited by hand.
- New exported types / fields have doc comments.
- Unit tests cover the new behaviour; e2e tests added for major features.
- Commit(s) are signed off and follow the existing commit style.
