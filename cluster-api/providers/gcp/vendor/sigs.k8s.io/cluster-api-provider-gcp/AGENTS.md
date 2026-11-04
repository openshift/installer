# Agent Development Guide

This guide describes the repository layout and development practices for
Cluster API Provider GCP (CAPG). Follow the existing code and test patterns in
the area you are changing. Also read and strictly follow the [AI_POLICY.md](AI_POLICY.md) for the
project's AI-assisted contribution expectations.
Exceptions apply for maintainers, when explicitly needed.

## Rules and Constraints

- Keep changes focused on the requested work. Do not create commits or pull
  requests unless explicitly asked.
- Use [Conventional Commits](https://www.conventionalcommits.org/) format for
  commit messages, PR titles, and release notes:
  `<type>[optional scope]: <description>` (for example, `fix: handle missing
  GCPCluster`).
- Add or update tests for behavior changes. Prefer tests alongside the code,
  following the patterns already used in that package.
- make lint and make test must be clean before handoff.
- Do not edit generated code or manifests by hand. Run `make generate` and
  include the resulting changes when generated output needs updating.
- Keep the Go version and dependencies in `go.mod` and `hack/tools/go.mod`
  unchanged unless the task requires a dependency or Go version update.
- New API fields need appropriate descriptions and validation. Follow the
  current API version and the project's existing conventions. Update generated
  code and manifests through the repository's generation targets.
- Do not add feature gates unless explicitly requested. When introducing
  gate-controlled functionality, make it opt-in (disabled by default) and
  document it. Follow the pattern in `feature/feature.go` and its setup in
  `main.go`.

## Version Upgrades/Updates/Bumps

For Go and Kubernetes/CAPI version bumps, follow the repository's upgrade skills and their linked
developer guides rather than updating versions ad hoc:

- Go version bumps: `.claude/skills/bump-go/SKILL.md`
- Kubernetes/CAPI version bumps: `.claude/skills/bump-k8s-capi/SKILL.md`

## Data Flow

```text
User creates CAPI and CAPG resources
  → Manager watches resources and enqueues the relevant controller
  → Reconciler loads the owning CAPI resources and creates a scope
  → Scope and service reconcilers compare desired state with GCP resources
  → Compute Engine or GKE APIs are called to create, update, or delete resources
  → Controller updates status, conditions, events, and finalizers
  → Reconcile again when cloud operations are still in progress
```

## Testing

Choose the smallest test layer that exercises the behavior:

- **Unit tests:** Use for package-level logic that does not need a Kubernetes
  API server or GCP. Follow nearby `_test.go` patterns; CAPG uses Go's testing
  package, Ginkgo, and Gomega.
- **Integration tests:** Use when testing controllers, webhooks, or Kubernetes
  API interactions. CAPG suites use `envtest`; `make test` runs the unit and
  integration tests. `envtest` does not exercise real GCP APIs.
- **End-to-end tests:** Add or update tests for user-visible provisioning
  lifecycle behavior that needs a real management/workload cluster or GCP API
  interactions. Follow the helpers, configured wait intervals, cleanup, and
  artifact collection patterns under `test/e2e/`.

Avoid fixed sleeps for asynchronous behavior; use the existing polling and
wait-interval patterns. E2E tests can create billable GCP resources and require
credentials, so run `make test-e2e` only when explicitly requested. Preserve
the normal cleanup and artifact collection behavior when adding or changing
E2E tests.

## Development and Validation

The repository's `Makefile` is the source of truth for available targets.
Common commands are:

```bash
make test       # Unit and integration tests
make lint       # Go lint checks
make generate   # Generate Go code and Kubernetes manifests
make verify     # Verify generated output, modules, conversions, and formatting
```

API or controller changes may require updating user documentation under
`docs/book/` and tests under the corresponding package. Keep generated CRDs,
RBAC, and other generated files in sync by running `make generate`.
