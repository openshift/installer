# Integration Guidelines

## Platform Structure

Most cloud platforms that support automated provisioning require code in these locations:

| Directory | Purpose |
|---|---|
| `pkg/types/<platform>/` | Type definitions (`platform.go`, `machinepool.go`, `metadata.go`, `doc.go`) |
| `pkg/types/<platform>/defaults/` | `SetPlatformDefaults()` function |
| `pkg/types/<platform>/validation/` | `ValidatePlatform()` and `ValidateMachinePool()` functions |
| `pkg/asset/installconfig/<platform>/` | Session/client setup, credential validation, cloud-specific install-config logic |
| `pkg/infrastructure/<platform>/clusterapi/` | CAPI provider implementing lifecycle hooks |
| `pkg/infrastructure/platform/platform.go` | Registration in `ProviderForPlatform()` switch |
| `pkg/destroy/<platform>/` | Cluster teardown logic with `register.go` |

Note: Some platforms (baremetal, external, none) have different structures and may not require all of these directories.

## Platform Name Constant

Each platform declares its name in `pkg/types/<platform>/doc.go`:

```go
package aws

const Name string = "aws"
```

This constant is used everywhere for switch dispatch (credentials checks, provisioning, destruction).

## Destroyer Registration

Destroyers use `init()` to self-register in a global map. Each platform creates `register.go`:

```go
func init() {
    providers.Registry["aws"] = New
}
```

The `New` function signature is `func(logrus.FieldLogger, *types.ClusterMetadata) (providers.Destroyer, error)`. The `Destroyer` interface requires only `Run() (*types.ClusterQuota, error)`.

## Cluster API Provider Hooks

CAPI providers implement `clusterapi.Provider` (mandatory) plus optional interfaces for lifecycle hooks. Use compile-time interface checks:

```go
var _ clusterapi.Provider           = (*Provider)(nil)
var _ clusterapi.PreProvider        = (*Provider)(nil)
var _ clusterapi.IgnitionProvider   = (*Provider)(nil)
var _ clusterapi.InfraReadyProvider = (*Provider)(nil)
var _ clusterapi.BootstrapDestroyer = (*Provider)(nil)
var _ clusterapi.PostDestroyer      = (*Provider)(nil)
```

Hook execution order: `PreProvision` -> infrastructure ready -> `InfraReady` -> `Ignition` -> machines provisioned -> `PostProvision` -> (later) `DestroyBootstrap` -> `PostDestroy`.

Optional `Timeouts` interface lets platforms override `NetworkTimeout()` and `ProvisionTimeout()`.

Register the provider in `pkg/infrastructure/platform/platform.go` (and `platform_altinfra.go`):

```go
case awstypes.Name:
    return clusterapi.InitializeProvider(&awscapi.Provider{}), nil
```

## Cluster Metadata

`ClusterMetadata` (in `pkg/types/clustermetadata.go`) carries per-platform metadata and is the input to destroyers. Add a new `*<platform>.Metadata` field to `ClusterPlatformMetadata` and extend the `Platform()` method.

## Resource Tagging Conventions

Tags are the primary mechanism for identifying cluster-owned resources during teardown.

**AWS**: Uses tag key `kubernetes.io/cluster/<infraID>` with values `"owned"` or `"shared"`. CAPI also uses `sigs.k8s.io/cluster-api-provider-aws/cluster/<infraID>`. The `resourcegroupstaggingapi` is used to discover all tagged resources by ARN, which are then dispatched by ARN service type. Shared resources have their tags removed; owned resources are deleted.

**GCP**: Uses labels `kubernetes-io-cluster-<infraID>: owned|shared` (format constant: `gcpconsts.ClusterIDLabelFmt`). CAPI adds `capg-cluster-<infraID>: owned`. Resources are also matched by name prefix `<infraID>-*`.

**Nutanix**: Uses Prism categories with key `kubernetes-io-cluster-<infraID>` and values `"owned"` or `"shared"`. Shared images have their category removed; owned resources are deleted.

**OpenStack**: Uses tag `openshiftClusterID=<infraID>` on Neutron resources. Some services (Swift) mangle case.

**vSphere**: Creates a tag category `openshift-<infraID>` with a tag named `<infraID>`, attached to VMs, resource pools, folders, and datastores.

## User-Agent Strings

Every SDK client must set a user-agent. The format varies by operation phase:

**AWS** (defined in `pkg/asset/installconfig/aws/sessionv2.go`):
- Install: `OpenShift/4.x Installer/<version>`
- Gather: `OpenShift/4.x Gather/<version>`
- Destroy: `OpenShift/4.x Destroyer/<version>`

Applied per-client via `awsmiddleware.AddUserAgentKeyValue(agent, version.Raw)`.

**GCP**: `option.WithUserAgent(fmt.Sprintf("OpenShift/4.x Destroyer/%s", version.Raw))`

**OpenStack**: `openshift-installer/<version>` via `gophercloud.UserAgent.Prepend()`. Use `openstackdefaults.NewServiceClient()` to ensure consistent user-agent.

**IBM Cloud / PowerVS**: `fmt.Sprintf("OpenShift/4.x Destroyer/%s", version.Raw)` via `Service.SetUserAgent()`.

## AWS SDK Conventions

- **SDK version**: AWS SDK v2 (`github.com/aws/aws-sdk-go-v2`). Import alias convention: service-specific aliases like `ec2v2` for EC2.
- **Retry config** (in `sessionv2.go`): 25 max attempts, 300s exponential jitter backoff, client-side rate limiting disabled.
- **Per-service clients**: Create via helper functions in `pkg/asset/installconfig/aws/` (e.g. `NewEC2Client`, `NewIAMClient`). Each accepts `EndpointOptions` and service-specific options.
- **Custom endpoints**: Platforms support `ServiceEndpoints` for private/gov regions. The `EndpointOptions` struct routes through a custom resolver.
- **Error handling**: Use `HandleErrorCode()` from `pkg/destroy/aws/errors.go` which extracts `smithy.APIError` codes. Check for specific codes like `"NoSuchEntity"`, `"NoSuchHostedZone"` via `strings.Contains`.

## Destroy Loop Pattern

All destroyers follow the same poll-until-clean pattern:

```go
err = wait.PollImmediateInfinite(time.Second*10, o.destroyCluster)
```

or

```go
wait.PollImmediateUntil(time.Second*10, deleteFunc, ctx.Done())
```

**AWS** uses a two-phase destroy: terminate EC2 instances first (to prevent resource recreation), then poll-delete remaining resources by ARN.

**GCP** uses staged functions -- resources with dependencies are grouped into sequential stages; resources within a stage run in parallel.

**Error suppression**: AWS uses `ErrorTracker` that logs errors at DEBUG level, escalating to WARN only after 5 minutes of repeated failures for the same resource.

## Defaults and Validation Pattern

Each platform provides:

1. `pkg/types/<platform>/defaults/SetPlatformDefaults(p *<platform>.Platform)` -- called during install-config generation.
2. `pkg/types/<platform>/validation/ValidatePlatform(...)` -- returns `field.ErrorList`.
3. `pkg/types/<platform>/validation/ValidateMachinePool(...)` -- validates machine pool config.

Signatures vary slightly per platform (some receive `*types.InstallConfig`, `*types.Networking`, or `fldPath`).

## Credentials and Permissions Checks

Platform credentials are validated in `pkg/asset/installconfig/platformcredscheck.go`. Each platform case creates a session/client to verify connectivity. Permission pre-flight checks live in `platformpermscheck.go` (currently only AWS and GCP implement substantive checks).

## Adding a New Platform

Checklist:

1. Add type definitions in `pkg/types/<platform>/` with `doc.go` (Name constant), `platform.go`, `machinepool.go`, `metadata.go`
2. Add `defaults/` and `validation/` packages
3. Add `+k8s:deepcopy-gen=package` to `doc.go` and generate `zz_generated.deepcopy.go`
4. Add `*<platform>.Metadata` to `ClusterPlatformMetadata` and extend `Platform()`
5. Create `pkg/asset/installconfig/<platform>/` for sessions and cloud validation
6. Add platform cases to `platformcredscheck.go` and `platformpermscheck.go`
7. Create CAPI provider in `pkg/infrastructure/<platform>/clusterapi/`
8. Register in both `platform.go` and `platform_altinfra.go`
9. Create destroyer in `pkg/destroy/<platform>/` with `register.go`
10. Add OWNERS file in each new directory
