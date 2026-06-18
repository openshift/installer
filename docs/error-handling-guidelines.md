# Error Handling Guidelines

## Error Wrapping

Prefer **`fmt.Errorf("...: %w", err)`** (stdlib) for all new code. The `github.com/pkg/errors` package is deprecated and no longer maintained -- do not introduce new usages. Existing `errors.Wrap`/`errors.Wrapf` calls should be migrated to `fmt.Errorf` when touching those files.

- Always wrap with context describing **what failed**, not why: `"failed to create EC2 client: %w"`, not `"EC2 client error because credentials were bad"`.
- Use `errors.Is` and `errors.As` (stdlib) for error inspection. Avoid `errors.Cause()` from `pkg/errors`.

## Validation Errors (`field.ErrorList`)

All install-config and platform validation uses `k8s.io/apimachinery/pkg/util/validation/field`.

**Structure:** Validation functions typically return `field.ErrorList`. Some helper functions (e.g., `ValidateIPinMachineCIDR`) return `error` for use in interactive prompts or simple checks.

```go
func ValidatePlatform(p *aws.Platform, fldPath *field.Path) field.ErrorList {
    allErrs := field.ErrorList{}
    if p.Region == "" {
        allErrs = append(allErrs, field.Required(fldPath.Child("region"), "region must be specified"))
    }
    // Compose child validators
    allErrs = append(allErrs, validateSubnets(p.VPC.Subnets, fldPath.Child("vpc", "subnets"))...)
    return allErrs
}
```

**Rules:**
1. Use `field.Required`, `field.Invalid`, `field.Forbidden`, `field.NotSupported`, `field.Duplicate` -- choose the one that matches the problem.
2. Build `fldPath` with `.Child()` and `.Index()` so errors print as `platform.aws.vpc.subnets[0].id`.
3. Pure type validation lives in `pkg/types/<platform>/validation/`. Platform-API validation (calling cloud APIs) lives in `pkg/asset/installconfig/<platform>/validation.go`.
4. Convert to `error` at the boundary using `.ToAggregate()`:
```go
if err := validation.ValidateInstallConfig(cfg, false).ToAggregate(); err != nil {
    return fmt.Errorf("invalid %q file: %w", filename, err)
}
```

## `diagnostics.Err` -- Structured User-Facing Errors

`pkg/diagnostics/error.go` defines `Err` with `Source`, `Reason` (CamelCase), and `Message` fields. Use it when the error should be categorized for end-user reporting.

```go
return &diagnostics.Err{Reason: "MissingQuota", Message: msg}
```

Currently only used for quota check failures. `Reason` must be a single CamelCase word. `Message` is free-form text written for non-expert users.

## Asset Dependency Graph Error Propagation

The asset store (`pkg/asset/store/store.go`) propagates errors with dependency context automatically:

```text
failed to generate asset "Cluster" -> failed to fetch dependency of "Cluster" -> ...
```

- `Generate()` returning a non-nil error stops the entire graph.
- The store wraps each level: `errors.Wrapf(err, "failed to fetch dependency of %q", a.Name())` and `errors.Wrapf(err, "failed to generate asset %q", a.Name())`.
- Do NOT duplicate the asset name in your own error messages; the store adds it.

**Sentinel error constants** in `pkg/asset/asset.go`:
- `InstallConfigError` = `"failed to create install config"`
- `ClusterCreationError` = `"failed to create cluster"`
- `ControlPlaneCreationError` = `"failed to provision control-plane machines"`

These are used as wrapper messages at top-level boundaries (e.g., `errors.Wrap(err, asset.InstallConfigError)`).

## Cloud Provider Error Handling

### AWS
- Extract error codes via `errors.As(err, &apiErr)` and `apiErr.ErrorCode()`. Helper: `pkg/destroy/aws/errors.go:HandleErrorCode()`.
- Permission checks: `pkg/asset/installconfig/aws/awserrors.go` provides `IsUnauthorized(err)` and `IsHTTPForbidden(err)`.
- In destroy paths, when a resource is already gone (e.g., `InvalidDhcpOptionsID.NotFound`), treat it as success and return nil.

### Azure
- Use `errors.As(err, &dErr)` with `autorest.DetailedError` to check HTTP status codes.
- Dedicated helpers: `isNotFoundError()`, `isAuthError()`, `isResourceGroupBlockedError()` in `pkg/destroy/azure/azure.go`.
- `NotFound` during destroy is not an error -- return nil.

### GCP
- Use `errors.As(err, &ae)` with `googleapi.Error` and check `ae.Code` (e.g., 404, 409, 429).
- Permission checks: `pkg/quota/gcp/gcp.go:IsUnauthorized()`.

### Power VS / IBM Cloud
- `ibmError` struct with `Status` and `Message`. Use `isNoOp()` to detect 404s during destroy.
- Aggregate errors from multi-resource operations with `utilerrors.NewAggregate(errs)`.

### OpenStack
- Custom error type `networkextensions.Error` (a `string` type implementing `error`) designed for `errors.As` matching.

## Destroy Operations -- Retry and Error Tracking

Destroy flows run in a poll loop (`wait.PollImmediateUntil` / `wait.PollImmediateInfinite`) and follow these conventions:

1. **Continue on failure**: Individual resource deletion errors are logged but do NOT abort the loop. The loop only exits when all resources are gone or the context is canceled.
2. **`pendingItemTracker`**: GCP and Power VS track resources remaining to delete. The loop returns `done=true` only when the tracker is empty.
3. **`ErrorTracker`** (AWS, `pkg/destroy/aws/errortracker.go`): Rate-limits repeated warnings. First occurrence logs at DEBUG; after a suppression duration, escalates to WARN. Use `tracker.suppressWarning(id, err, logger)`.
4. **Staged deletion**: GCP uses ordered stages (stop instances -> delete resources -> delete networking). Errors in a stage skip subsequent stages for that iteration but the poll retries.

## General Conventions

- **Graceful degradation for permissions** (exception, not the norm): Quota and capability pre-flight checks are best-effort -- when they fail due to missing permissions, log a warning and skip the check rather than failing the install. This applies **only** to optional pre-flight checks, not to operations required for a successful install. Pattern: `if IsUnauthorized(err) { logrus.Warn(...); return nil }`.
- **`utilerrors.NewAggregate`**: Use to combine multiple independent errors (destroy loops, SSH key loading, multi-resource operations). Returns nil if the slice is empty.
- **Error messages**: Start with lowercase, do not end with punctuation. Use `"failed to <verb>"` not `"error <verbing>"`.
- **`errors.As` / `errors.Is`**: Prefer over type assertions for cloud SDK errors. The codebase uses these consistently across all platforms.
