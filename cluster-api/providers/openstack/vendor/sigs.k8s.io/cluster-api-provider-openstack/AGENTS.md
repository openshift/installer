# AGENTS.md - Agent Guidelines for Cluster API Provider OpenStack

This document provides guidelines and useful commands for AI agents contributing to the Cluster API Provider OpenStack (CAPO) repository.

> **⚠️ IMPORTANT**: When making changes to Makefile targets, PR requirements, code generation workflows, verification steps, or any other information referenced in this document, **AGENTS.md must be updated accordingly** to keep it synchronized with the actual project state.

## Overview

Cluster API Provider OpenStack (CAPO) is a Kubernetes-native declarative infrastructure provider for managing Kubernetes clusters on OpenStack. It implements the Cluster API (CAPI) specification for self-managed Kubernetes clusters on OpenStack infrastructure.

**Key Concepts:**
- **CAPI**: Cluster API - the upstream Kubernetes project defining cluster lifecycle APIs
- **CAPO**: Cluster API Provider OpenStack - this repository
- **Reconciler**: Controller-runtime pattern for managing Kubernetes custom resources
- **Scope**: Context and configuration wrapper for controllers and services

## Key Requirements for Contributors

### Legal Requirements

- **CLA Required**: All contributors MUST sign the Kubernetes Contributor License Agreement (CLA)
- See: https://git.k8s.io/community/CLA.md

### Pull Request Labels

All code PRs MUST be labeled with one of:
- ⚠️ `:warning:` - major or breaking changes
- ✨ `:sparkles:` - feature additions
- 🐛 `:bug:` - patch and bugfixes
- 📖 `:book:` - documentation or proposals
- 🌱 `:seedling:` - minor or other

## Essential Make Targets

### Code Quality & Verification

> **⚠️ IMPORTANT**: The `make verify` targets compare the working tree against `HEAD` (the last commit). This means:
> - `make verify` will **always fail** if you have uncommitted changes, even if those changes are correct
> - You should **commit your changes first**, then run `make verify` to confirm everything is in order
> - The typical workflow is: make changes → `make generate` → `make modules` → `make test` → **commit** → `make verify`

```bash
# Run all verification checks (should pass before submitting PR)
# NOTE: This compares against HEAD, so commit your changes first!
make verify

# Verify boilerplate headers
make verify-boilerplate

# Verify go modules are up to date (runs `go mod tidy` and diffs against HEAD)
make verify-modules

# Verify generated code is up to date (runs `make generate` and diffs against HEAD)
make verify-gen

# Verify container images for vulnerabilities
make verify-container-images

# Check code for security vulnerabilities
make verify-govulncheck

# Run all security checks (images + code)
make verify-security
```

### Linting

```bash
# Lint codebase
make lint

# Lint and auto-fix issues
make lint-update

# Run faster linters only
make lint-fast
```

### Testing

```bash
# Run unit tests
make test

# Run unit tests for CAPO specifically
make test-capo

# Build e2e test binaries
make build-e2e-tests

# Run e2e tests (requires OpenStack environment)
make test-e2e

# Run conformance tests
make test-conformance

# Compile e2e tests (verify they build)
make compile-e2e
```

### Code Generation

```bash
# Generate ALL code (manifests, deepcopy, clients, mocks, docs)
make generate

# Generate Go code (mocks, etc.)
make generate-go

# Generate controller-gen code (deepcopy, etc.)
make generate-controller-gen

# Generate client code (clientsets, listers, informers)
make generate-codegen

# Generate CRD manifests
make generate-manifests

# Generate API documentation
make generate-api-docs

# Generate cluster templates
make templates
```

### Dependency Management

```bash
# Update go modules
make modules

# Check for API differences (useful before breaking changes)
make apidiff
```

### Building

```bash
# Build manager binaries
make managers

# Build all binaries
make binaries

# Build Docker image
make docker-build

# Build debug Docker image
make docker-build-debug

# Build for all architectures
make docker-build-all
```

### Cleanup

```bash
# Clean all build artifacts
make clean

# Clean binaries only
make clean-bin

# Clean temporary files
make clean-temporary

# Clean release artifacts
make clean-release
```

## Important Development Patterns

### Adding New OpenStack Resources

1. Define API types in `/api/v1beta1` (or `/api/v1alpha1` for experimental features)
2. Run `make generate` to create deepcopy methods and update CRDs
3. Create controller in `/controllers`
4. Create service implementation in `/pkg/cloud/services/<category>`
5. Update or create scope in `/pkg/scope` if needed
6. Add webhooks in `/pkg/webhooks` for validation/defaulting
7. Add unit tests for controller and services
8. Update documentation
9. Generate cluster templates if applicable with `make templates`

### Testing Strategy

1. **Unit Tests**: Test individual functions/methods with mocks
2. **Integration Tests**: Test controller behavior with envtest
3. **E2E Tests**: Deploy real clusters on OpenStack, verify functionality
4. **Conformance Tests**: Run upstream Kubernetes conformance suite

## Pre-Submit Checklist for Agents

Before submitting a PR, ensure:

1. **Code is generated and up to date**:
   ```bash
   make generate
   ```

2. **Modules are tidy**:
   ```bash
   make modules
   ```

3. **Code passes linting**:
   ```bash
   make lint
   ```

4. **Tests pass**:
   ```bash
   make test
   ```

5. **All verification checks pass**:
   ```bash
   make verify
   ```

## Common Workflows

### Making Code Changes

1. Make your code changes
2. Run code generation: `make generate`
3. Update modules if needed: `make modules`
4. Run tests: `make test`
5. Lint the code: `make lint`
6. **Commit changes** with descriptive message
7. Verify everything: `make verify` (this compares against HEAD, so must be done after commit)

### Updating Dependencies

1. Update `go.mod` or `hack/tools/go.mod` as needed (e.g., `go get sigs.k8s.io/cluster-api@v1.x.x`)
2. Run: `make modules` to tidy dependencies
3. Run: `make generate` to regenerate code (dependency updates often change generated code)
4. Run: `make test` to ensure everything still works
5. **Commit all changes** (go.mod, go.sum, and any regenerated files)
6. Run: `make verify` to confirm everything is in order

## Common Issues

### Linting Errors

The project uses golangci-lint. If you get lint errors:
1. Run `make lint-update` first to auto-fix
2. Check `.golangci.yml` for enabled linters
3. Some issues require manual fixes (cognitive complexity, error handling, etc.)
4. Don't disable linters without good reason - fix the underlying issue

### Test Failures

- **envtest issues**: Ensure KUBEBUILDER_ASSETS is set correctly
- **Flaky E2E tests**: Transient infrastructure issues, failure to deploy devstack

### Generated File Drift

If `make verify` fails with generated file drift:
1. Run `make generate` to regenerate all files
2. Run `make modules` to ensure go.mod/go.sum are tidy
3. Review the changes to ensure they're expected
4. **Commit the generated files** - verify targets compare against HEAD!
5. Run `make verify` again - it should pass now
6. Never manually edit generated files

### Verify Fails Even After Running Generate/Modules

If `make verify-gen` or `make verify-modules` fails even after running `make generate` or `make modules`:
- Remember that these targets compare against `HEAD` (the last commit)
- If you have uncommitted changes, they will show as "drift" even if correct
- **Solution**: Commit your changes first, then run `make verify`
- This is by design - it ensures the committed code is complete and self-consistent

## Documentation

Primary documentation is in `/docs/book/src/` (mdBook format):
- Getting started guides
- Developer documentation
- Troubleshooting guides
- API reference
- Cluster template documentation

Build and serve docs locally:
```bash
make -C docs/book serve
```

## Quick Reference

| Task | Command |
|------|---------|
| Full verification before PR | `make verify && make test` |
| Generate all code | `make generate` |
| Update dependencies | `make modules` |
| Lint and fix | `make lint-update` |
| Run tests | `make test` |
| Build binary | `make managers` |
| Build Docker image | `make docker-build` |
| Clean everything | `make clean` |
| Check API compatibility | `make apidiff` |
| Generate templates | `make templates` |
| Build and serve docs | `make -C docs/book serve` |
