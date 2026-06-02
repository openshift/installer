# Performance Guidelines

## Destroy Flow Concurrency

The installer uses three distinct concurrency models for cloud resource teardown, depending on the platform.

### OpenStack: Goroutine-per-resource-type with exponential backoff

Each resource type (servers, ports, networks, etc.) runs in its own goroutine. A shared `returnChannel` collects completion signals. Each goroutine retries internally using `wait.ExponentialBackoff` with `Duration: 15s, Factor: 1.3, Steps: 25`.

```go
// pkg/destroy/openstack/openstack.go
for name, function := range deleteFuncs {
    go deleteRunner(ctx, name, function, opts, o.Filter, o.Logger, returnChannel)
}
for i := 0; i < len(deleteFuncs); i++ {
    res := <-returnChannel
}
```

Routers are deleted sequentially *after* all other resources complete, because FIP tracking depends on them existing.

### OpenStack port deletion: Fixed worker pool (10 workers)

Port deletion fans out to exactly 10 worker goroutines reading from a shared channel. This bounds concurrent API calls to the networking service.

```go
// pkg/destroy/openstack/openstack.go
const workersNumber = 10
portsChannel := make(chan ports.Port, workersNumber)
```

### IBMCloud / PowerVS: Staged parallel teardown with WaitGroup

Resources are organized into ordered stages. Within each stage, functions run in parallel goroutines. The next stage only starts after all goroutines in the current stage complete. This enforces dependency ordering (e.g., instances before load balancers, load balancers before subnets).

```go
// pkg/destroy/ibmcloud/ibmcloud.go, pkg/destroy/powervs/powervs.go
for _, stage := range stagedFuncs {
    var wg sync.WaitGroup
    for _, f := range stage {
        wg.Add(1)
        go o.executeStageFunction(f, errCh, &wg)
    }
    // wait for stage completion before proceeding
}
```

PowerVS additionally enforces a per-stage timeout of `15 * time.Minute` via `time.After`.

### GCP: Sequential stages, no goroutines

GCP runs stages sequentially. Within each stage, resource functions run serially. If any function in a stage fails, later stages are skipped entirely and the outer `wait.PollImmediateInfinite` retries the full pass at 10-second intervals.

## Semaphore Pattern (OpenStack object storage)

`pkg/destroy/openstack/semaphore.go` implements a bounded-concurrency semaphore using a buffered channel and `sync.WaitGroup`. Used when bulk-deleting Swift objects (capped at 3 concurrent operations):

```go
queue := newSemaphore(3)
queue.Add(func() { /* delete batch */ })
queue.Wait()
```

Use this pattern (not raw goroutines) when you need to cap concurrency for API calls to a specific service.

## Static network config: `x/sync/semaphore`

`pkg/asset/agent/manifests/staticnetworkconfig/generator.go` uses `golang.org/x/sync/semaphore.Weighted` to limit concurrent `nmstatectl` subprocess invocations to 30 (configurable via `MAX_CONCURRENT_NMSTATECTL_GENERATIONS`). Use weighted semaphores when limiting external process concurrency.

## AWS SDK Retry Configuration

`pkg/asset/installconfig/aws/sessionv2.go` configures the AWS SDK v2 retryer with specific tuning:

- **MaxAttempts**: 25
- **Backoff**: Exponential with jitter, **300-second max** (matching SDK v1 behavior; the SDK v2 default of 20s was too aggressive)
- **Client-side rate limiter**: Explicitly **disabled** (`ratelimit.None`) because the default token bucket does not suit the destroy code's burst patterns

These values are set in `getDefaultConfigOptions()` and apply to all AWS API calls unless overridden.

## Thread-Safe Logging: LinePrinter

`pkg/lineprinter/lineprinter.go` wraps a `Print` function to provide a mutex-protected `io.WriteCloser`. It buffers bytes, emits complete newline-terminated lines during `Write()`, and flushes partial lines on `Close()`.

This is used when connecting subprocess stdout/stderr to logrus (e.g., SSH commands in `pkg/gather/ssh/ssh.go`). Always use `LinePrinter` (not raw logrus calls) when multiple goroutines write to the same log sink through a shared writer.

## Error Suppression in Destroy Loops

`errorTracker` (duplicated in `pkg/destroy/gcp/`, `pkg/destroy/powervs/`, `pkg/destroy/ibmcloud/`) suppresses repeated warnings. It logs an error at DEBUG level and only promotes to WARN every 5 minutes per unique identifier. This prevents log flooding during retry loops.

## Cluster API System Lifecycle

`pkg/clusterapi/system.go` manages a local control plane (etcd + kube-apiserver) with these concurrency primitives:

- `sync.Mutex` protects `Run()`, `Client()`, and `Teardown()` from concurrent access
- `sync.Once` ensures `Teardown()` is idempotent
- `sync.WaitGroup` tracks controller goroutine lifecycle
- Teardown has a 60-second hard timeout; controllers that do not stop in time are abandoned

## Metadata Caching with Mutex

`pkg/asset/installconfig/aws/metadata.go` uses `sync.Mutex` to protect lazy-loaded metadata fields (availability zones, regions, subnets). Each getter locks, checks the cache, populates on miss, and returns. This is the standard pattern for metadata objects accessed from multiple asset generators.

## Session/Credential Logging with sync.Once

Azure and OpenStack session setup (`pkg/asset/installconfig/azure/session.go`, `pkg/asset/installconfig/openstack/session.go`) use `map[string]*sync.Once` keyed by file path to ensure credential-source log messages appear exactly once per file, even when sessions are created multiple times.

## Azure Page Blob Upload: Thread-Grouped Parallelism

`pkg/infrastructure/azure/storage.go` uploads VHD images using up to 64 parallel goroutines per group, with groups processed sequentially. Each goroutine retries its page upload 3 times internally. Errors are collected via a buffered channel; only the first error is preserved.

## Context and Timeout Conventions

- **PowerVS**: Uses a factory function `contextWithTimeout()` returning `context.WithTimeout(context.Background(), 30*time.Minute)` for every API call. Each function creates and defers-cancels its own context.
- **GCP**: Uses `context.Background()` for most operations, with `context.WithTimeout(ctx, 10*time.Minute)` for expensive aggregated-list calls.
- **OpenStack**: Passes `context.TODO()` from the top-level `Run()` through all delete functions. No per-call timeouts; retry is handled by `wait.ExponentialBackoff`.

## Key Rules

1. **Never spawn unbounded goroutines for API calls.** Use the semaphore pattern, fixed worker pools, or staged WaitGroups.
2. **Order destroy stages by dependency where the platform supports it.** Platforms like GCP and IBMCloud use explicit stages to enforce ordering (instances before LBs, LBs before subnets, everything before routers). AWS uses tag-based discovery where resources are attempted in arbitrary order each pass -- dependencies resolve naturally over multiple iterations as blocking resources are deleted first.
3. **Retry with backoff, not tight loops.** Use `wait.ExponentialBackoff` or `wait.PollImmediateInfinite` from `k8s.io/apimachinery/pkg/util/wait`.
4. **Treat 404 as success in destroy paths.** Resources may already be deleted by a concurrent operation or prior attempt.
5. **Log transient errors at Debug, not Warn.** Use `errorTracker.suppressWarning` to avoid flooding logs during long retry sequences.
6. **Protect shared mutable state with Mutex.** Metadata objects, CAPI system state, and LinePrinter buffers all require synchronization.
7. **Always defer context cancellation** immediately after creation to prevent goroutine leaks.
