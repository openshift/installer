# AGENTS.md

This document is for AI agents working in the OpenShift Installer codebase. It supplements the existing docs -- read `CLAUDE.md` for build/test commands, `README.md` for user-facing overview, and `CONTRIBUTING.md` for contribution workflow.

## Guideline Documents

These files contain detailed, domain-specific guidance. Read the relevant one before making changes in that area:

- [`docs/security-guidelines.md`](docs/security-guidelines.md) -- TLS, credentials, FIPS compliance, secrets handling
- [`docs/error-handling-guidelines.md`](docs/error-handling-guidelines.md) -- Error wrapping, validation patterns, cloud error handling
- [`docs/testing-guidelines.md`](docs/testing-guidelines.md) -- Mocks, table-driven tests, integration test setup
- [`docs/integration-guidelines.md`](docs/integration-guidelines.md) -- Cloud platform structure, Cluster API, tagging, SDK usage
- [`docs/performance-guidelines.md`](docs/performance-guidelines.md) -- Concurrency, rate limiting, destroy flow optimization

## Code Style and Conventions

### Import Organization

Imports are sorted into four groups by `hack/go-fmt.sh` using the [gci tool](https://github.com/daixiang0/gci):

1. Standard library
2. Third-party packages
3. `github.com/openshift` packages (project-internal prefix)
4. Blank imports (for side effects)

Groups are separated by blank lines. Always run `hack/go-fmt.sh .` before submitting.

### Import Alias Conventions

Platform-specific imports frequently collide on package names. The codebase uses consistent alias prefixes:

- Type packages: `awstypes`, `azuretypes`, `gcptypes`, `vspheretypes`, etc.
- Install-config sub-assets: `icazure`, `icgcp`, `icibmcloud`, etc.
- Infrastructure CAPI providers: `awscapi`, `gcpcapi`, `vspherecapi`, etc.
- Infrastructure packages: `azureinfra`, `baremetalinfra`, etc.
- CAPI manifests: `capimanifests`

When adding a new import that conflicts with an existing package name, follow the alias pattern already used in that file or the patterns above.

### Package Documentation

Every package should have a `doc.go` file with a package comment. Platform-type packages also define their platform `Name` constant in `doc.go`:

```go
// Package aws contains AWS-specific structures for installer
// configuration and management.
// +k8s:deepcopy-gen=package
package aws

// Name is name for the AWS platform.
const Name string = "aws"
```

### Generated Code

Several files are generated and must not be edited by hand:

- `zz_generated.deepcopy.go` -- deep copy methods for types, generated via `+k8s:deepcopy-gen=package` directive
- `data/data/install.openshift.io_installconfigs.yaml` -- CRD definition, regenerated with `go generate ./pkg/types/installconfig.go`
- Mock files under `pkg/asset/mock/` -- regenerated with `hack/go-genmock.sh`

If you change an interface that has mocks or a type that has deepcopy, you must regenerate.

## Go Module and Vendoring

- The module is `github.com/openshift/installer`.
- All dependencies are vendored in `vendor/`. The `vendor/` directory is checked in.
- **Critical rule**: Always commit vendored dependency changes in a separate commit from functional code changes. This is enforced by convention and makes reviews tractable.
- When updating dependencies: `go get <pkg>`, then `go mod tidy`, then `go mod vendor`.
- Updating `github.com/openshift/api` is a special case that also requires regenerating the install-config CRD (see `CLAUDE.md`).

## Architecture: What Is Not Obvious

### The Asset DAG

The installer's core architecture is a directed acyclic graph (DAG) of "assets." Full design is in `docs/design/assetgeneration.md`, but here are the essentials for working in the code:

**Every piece of installer output is an Asset.** The `Asset` interface (`pkg/asset/asset.go`) has three methods:
- `Dependencies() []Asset` -- declares what this asset needs
- `Generate(ctx, Parents) error` -- produces the asset from its dependencies
- `Name() string` -- human-readable identifier

**`WritableAsset`** extends `Asset` with `Files()` and `Load()` -- it can be serialized to disk and read back. The installer chains targets: `install-config` -> `manifests` -> `ignition-configs` -> `cluster`. Each target consumes (and removes from disk) the previous target's files.

**Targets** are defined in `pkg/asset/targets/targets.go`. They group the writable assets for each CLI command (`create install-config`, `create manifests`, `create ignition-configs`, `create cluster`).

**The Store** (`pkg/asset/store/`) manages the DAG resolution. It does depth-first traversal, generating dependencies before dependents. Assets can be loaded from disk (user-provided overrides) or from an internal state file.

**When adding a new asset:**
1. Create a struct implementing `Asset` (or `WritableAsset` if it produces files)
2. Declare dependencies in `Dependencies()`
3. In `Generate()`, call `parents.Get(...)` to retrieve dependency state
4. If writable, add it to the appropriate target list in `pkg/asset/targets/targets.go`
5. Verify the interface is satisfied with `var _ asset.WritableAsset = (*YourAsset)(nil)`

### Platform Code is Spread Across Many Packages

Adding or modifying a platform feature typically requires touching multiple locations. For a platform named `<plat>`:

| Concern | Location |
|---|---|
| Type definitions (Platform, MachinePool structs) | `pkg/types/<plat>/` |
| Default values | `pkg/types/<plat>/defaults/` |
| Validation | `pkg/types/<plat>/validation/` |
| Platform name constant | `pkg/types/<plat>/doc.go` |
| Install-config sub-assets (metadata, creds checks) | `pkg/asset/installconfig/<plat>/` |
| Infrastructure provisioning (CAPI) | `pkg/infrastructure/<plat>/` |
| CAPI provider binaries | `cluster-api/providers/<plat>/` |
| Terraform variables | `pkg/asset/cluster/tfvars/` or `pkg/tfvars/<plat>/` |
| Cluster destroy logic | `pkg/destroy/<plat>/` |
| Destroy provider registration | `pkg/destroy/<plat>/register.go` (via `init()`) |
| Embedded data (manifests, templates) | `data/data/` |
| Platform wiring in provider switch | `pkg/infrastructure/platform/platform.go` |

### The Registry Pattern for Destroyers

Destroy providers use an `init()`-based registry. Each platform's `pkg/destroy/<plat>/register.go` registers a factory function into `providers.Registry` (a `map[string]NewFunc`). The platform string key matches the `Name` constant from `pkg/types/<plat>/doc.go`. If you add a new destroy provider, you need both the implementation and the `register.go` with the `init()` function.

### Infrastructure Providers and CAPI

Most platforms now provision infrastructure through Cluster API (CAPI). The wiring is in `pkg/infrastructure/platform/platform.go`, which maps platform names to `infrastructure.Provider` implementations. CAPI providers are initialized via `clusterapi.InitializeProvider()`.

The `platform.go` file uses build tags (`altinfra` / default) to support alternative infrastructure configurations. When adding a new CAPI-based platform, update both `platform.go` and `platform_altinfra.go`.

### Embedded Data

Static assets (bootstrap scripts, manifest templates, CAPI manifests) live in `data/data/`. In development builds, `data/assets.go` serves these from disk (honoring `OPENSHIFT_INSTALL_DATA` env var). In release builds, they are embedded into the binary via `data/assets_generate.go`.

### Feature Gates

The installer uses OpenShift feature gates (`configv1.FeatureGateName`) to conditionally enable functionality. The local interface is in `pkg/types/featuregates/`. Feature gates are passed through to platform provider selection and validation. When adding gated behavior, follow the existing pattern of accepting a `featuregates.FeatureGate` parameter and checking `fg.Enabled(...)`.

### Build Tags

The codebase uses build tags for conditional compilation:
- `release` / default -- controls whether data assets are embedded or read from disk
- `altinfra` -- alternative infrastructure provider wiring (used for CI/dev variants)

## Commit Message and PR Expectations

### Commit Format

```text
<subsystem>: <what changed>

<why this change was made>

Fixes #<issue-number>
```

- Subject line under 70 characters; body wrapped at 80
- The subsystem is typically the platform name (`aws`, `azure`, `vsphere`), installation method (`agent`, `ibi`), or component (`cluster-api`, `terraform`)
- For test-only changes, use `unit tests` or `integration tests` as the subsystem

### PR Conventions

- Each PR is reviewed by OWNERS (defined per-directory in `OWNERS` files, with aliases in `OWNERS_ALIASES`)
- Run all linters listed in `CONTRIBUTING.md` before submitting
- Vendored dependency updates go in their own commit, separate from functional changes
- All code under `cmd/`, `data/`, and `pkg/` must have unit tests

## Common Pitfalls

1. **Forgetting to update multiple platform locations.** A new install-config field for a platform typically needs type definition, defaults, validation, and possibly asset/infrastructure changes. See the platform locations table above.

2. **Editing generated files.** Files named `zz_generated.deepcopy.go` and mock files are regenerated by tooling. Edit the source interface or type, then regenerate.

3. **Breaking the asset DAG.** If you add a dependency to an asset but create a cycle, the installer will panic at runtime. The graph is traversed depth-first; dependencies must form a DAG.

4. **Missing `init()` registration.** New destroy providers or platform registrations that use the `init()` pattern will silently not work if the `register.go` file is missing or the package is not imported somewhere in the binary's import chain.

5. **Import alias drift.** The codebase has established alias patterns (e.g., `awstypes`, `icazure`). Using inconsistent aliases makes the code harder to navigate. Check existing files in the same package for conventions.

6. **Not running `hack/go-fmt.sh`** after adding imports. The four-group import ordering is enforced and will fail CI if not followed.

7. **Validation uses `field.ErrorList`.** Platform validation functions return `field.ErrorList` (from `k8s.io/apimachinery/pkg/util/validation/field`), not plain errors. Follow this pattern for all new validation code.

8. **Data directory changes require rebuild.** Changes to files under `data/data/` are picked up automatically in dev builds (from disk), but release builds require regeneration via `go generate`.
