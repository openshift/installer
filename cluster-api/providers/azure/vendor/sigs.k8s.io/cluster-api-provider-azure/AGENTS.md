## Overview

Cluster API Provider Azure (CAPZ) is a Kubernetes-native declarative infrastructure provider for managing Azure clusters. It implements the Cluster API (CAPI) specification for both self-managed (IaaS) and AKS-managed Kubernetes clusters on Azure.

**Key Concepts:**
- CAPI: Cluster API - the upstream Kubernetes project defining cluster lifecycle APIs
- CAPZ: Cluster API Provider Azure - this repository
- ASO: Azure Service Operator - used for declarative Azure resource management
- Reconciler: Controller-runtime pattern for managing Kubernetes custom resources

## Architecture

### Core Components

1. **API Definitions (`/api` and `/exp/api`)**
   - `v1beta1`: Stable API version for core resources
   - `v1alpha1`: Deprecated, use v1beta1
   - `/exp`: Experimental features (MachinePools, Managed clusters)
   - Key resources: AzureCluster, AzureMachine, AzureMachinePool, AzureManagedCluster, AzureManagedControlPlane

2. **Controllers (`/controllers` and `/exp/controllers`)**
   - Each controller reconciles a specific custom resource type
   - Controllers use Azure SDK and ASO to manage Azure infrastructure
   - Reconciliation pattern: observe state → determine actions → apply changes → update status
   - Key controllers: AzureClusterReconciler, AzureMachineReconciler, AzureMachinePoolReconciler

3. **Azure Services Layer (`/azure/services`)**
   - Service-specific Azure API clients organized by Azure resource type
   - Examples: `virtualnetworks`, `subnets`, `loadbalancers`, `virtualmachines`, `vmss`
   - Each service implements the `Reconciler` interface: `Reconcile()`, `Delete()`
   - Services handle Azure API calls, credential management, and error handling

4. **Scope Package (`/azure/scope`)**
   - Provides context and configuration for controllers
   - Scopes encapsulate cluster/machine specs, credentials, and Azure clients
   - Key scopes: ClusterScope, MachineScope, MachinePoolScope, ManagedControlPlaneScope

5. **Feature Gates (`/feature`)**
   - Controls experimental/optional functionality
   - Important gates: `MachinePool`, `ASOAPI`, `EdgeZone`, `ClusterResourceSet`

### Data Flow

```
User creates K8s resource → Controller watches → Reconciler triggered →
Scope created → Azure service methods called → Azure API interactions →
Status updated → Requeue if needed
```

### Two Deployment Models

1. **Self-Managed (IaaS)**: CAPZ creates VMs, networks, load balancers via Azure APIs
2. **Managed (AKS)**: CAPZ creates/manages AKS clusters and agent pools via ASO

## Development Commands

### Building and Testing

```bash
# Build the manager binary
make manager

# Run unit tests
make test

# Run unit tests with race detector
make test-cover

# Run linting
make lint

# Fix lint issues automatically
make lint-fix

# Generate code (deepcopy, CRDs, webhooks, mocks)
make generate

# Verify all generated files are up-to-date
make verify
```

### Running a Single Test

```bash
# Run specific test by name
KUBEBUILDER_ASSETS="$(make setup-envtest 2>&1 | grep -o '/.*')" go test -v -run TestFunctionName ./path/to/package

# Run all tests in a package
KUBEBUILDER_ASSETS="$(make setup-envtest 2>&1 | grep -o '/.*')" go test -v ./controllers/
```

### Docker Images

```bash
# Build controller image (defaults to dev tag)
make docker-build

# Build and push to registry
REGISTRY=myregistry.io make docker-build docker-push

# Build all architectures
make docker-build-all
```

### Local Development with Tilt

```bash
# Create Kind cluster and start Tilt
make kind-create tilt-up

# Use AKS as management cluster instead
make aks-create tilt-up

# Delete Kind cluster
make kind-reset
```

**tilt-settings.yaml** is required with Azure credentials (see docs/book/src/developers/development.md for details).

### E2E Testing

```bash
# Run E2E tests (requires Azure credentials in env)
./scripts/ci-e2e.sh

# Run conformance tests
./scripts/ci-conformance.sh
```

## Important Patterns

### Controller Reconciliation

Controllers follow this pattern:
1. Fetch the resource being reconciled
2. Check deletion timestamp; run finalizer logic if deleting
3. Create/update scope with credentials and config
4. Call service Reconcile() or Delete() methods
5. Update resource status
6. Return result with requeue if needed

### Error Handling

- Transient errors: Return `ctrl.Result{RequeueAfter: timeout}` to retry
- Permanent errors: Log error, update status conditions, don't requeue
- Long-running operations: Use Azure async patterns with futures

### Adding New Azure Resources

1. Define API in `/api/v1beta1` or `/exp/api/v1beta1`
2. Run `make generate` to create deepcopy methods
3. Create controller in `/controllers`
4. Create service in `/azure/services/<resourcetype>`
5. Update scope if needed in `/azure/scope`
6. Add webhooks in `/api/v1beta1` for validation/defaulting
7. Generate CRDs with `make generate-manifests`
8. Add tests

### Credential Management

CAPZ supports multiple authentication methods:
- Service Principal (client ID/secret)
- Managed Identity (system or user-assigned)
- Workload Identity (recommended for AKS)

Credentials are cached in `azure.CredentialCache` and scoped to cluster identity.

## Code Generation

Several code generators are used:

```bash
make generate-go          # controller-gen for deepcopy, conversion-gen
make generate-manifests   # CRDs, RBAC, webhooks
make generate-flavors     # Cluster template flavors
make generate-addons      # CNI and other addons
```

Always run after changing:
- API types (add fields, change validation)
- RBAC annotations
- Webhook logic
- Flavor templates

## Testing Strategy

1. **Unit Tests**: Test individual functions/methods with mocks
2. **Integration Tests**: Test controller logic with fake Kubernetes client
3. **E2E Tests**: Deploy real clusters in Azure, verify functionality
4. **Conformance Tests**: Run upstream Kubernetes conformance suite

### Mocks

GoMock is used for Azure client mocks:
```bash
make generate-go  # Regenerates mocks in azure/services/*/mock_*/
```

## Key Files and Folders to Know

- `main.go`: Entry point, registers controllers and webhooks
- `Makefile`: All build/test/dev targets
- `Tiltfile`: Local development with Tilt
- `go.mod`: Go dependencies (uses Go 1.24+)
- `config/`: Kustomize configurations for CRDs, RBAC, webhooks, manager
- `templates/`: Cluster template flavors for different scenarios
- `test/e2e/`: E2E test suites and data files

## Common Issues

### Linting Errors

The project uses golangci-lint with strict settings. If you get lint errors:
1. Run `make lint-fix` first
2. Check `.golangci.yml` for enabled linters
3. Some issues require manual fixes (cognitive complexity, etc.)

### Test Failures

- **envtest issues**: Ensure KUBEBUILDER_ASSETS is set correctly
- **Race detector failures**: May indicate real concurrency bugs
- **Flaky E2E tests**: Azure API throttling or transient infrastructure issues

### Generated File Drift

If `make verify` fails:
1. Run `make generate`
2. Commit the generated files
3. Do not manually edit generated files

## Experimental Features

Features under `/exp` require feature gates to be enabled:
- Set in environment: `export EXP_MACHINE_POOL=true`
- Documented in `/feature/feature.go`

## Azure Service Operator (ASO) Integration

For managed clusters, CAPZ uses ASO to create Azure resources declaratively:
- ASO CRDs are vendored in `config/aso/`
- Controllers create ASO resources which ASO reconciles to Azure
- CAPZ watches ASO resources and updates CAPI status based on ASO status

## Performance Considerations

- Request coalescing: `/pkg/coalescing` deduplicates rapid reconcile requests
- Credential caching: Avoids repeated Azure auth calls
- Concurrency: Configured via flags like `--azurecluster-concurrency`
- Timeouts: Service reconcile timeout, Azure API call timeout (configurable)

## Documentation

Primary documentation is in `/docs/book/src/` (mdBook format):
- Getting started guides
- Developer documentation
- Troubleshooting guides
- API reference

Build docs: `make -C docs/book serve` (requires mdBook)
