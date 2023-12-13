# Cluster API

The installer uses Cluster API controllers through a local control plane powered by `kube-apiserver` and `etcd` running locally.
For more details, see the [enhancement][enhancement].

### Local control plane

The local control plane is setup using the previously available work done in Controller Runtime through [envtest](https://github.com/kubernetes-sigs/controller-runtime/tree/main/tools/setup-envtest). Envtest was born due to a necessity to run integration tests for controllers against a real API server, register webhooks (conversion, admission, validation), and managing the lifecycle of Custom Resource Definitions.

Over time, `envtest` matured in a way that now can be used to run controllers in a local environment, reducing or eliminating the need for a full Kubernetes cluster to run controllers.

At a high level, the local control plane is responsible for:
- Setting up certificates for the apiserver and etcd.
- Running (and cleaning up, on shutdown) the local control plane components.
- Installing any required component, like Custom Resource Definitions (CRDs)
    - For Cluster API core the CRDs are stored in `data/data/cluster-api/core-components.yaml`.
    - Infrastructure providers are expected to store their components in `data/data/cluster-api/<name>-infrastructure-components.yaml`
- Upon install, the local control plane takes care of modifying any webhook (conversion, admission, validation) to point to the `host:post` combination assigned.
    - Each controller manager will have its own `host:port` combination assigned.
    - Certificates are generated and injected in the server, and the client certs in the api-server webhook configuration.
- For each process that the local control plane manages, a health check (ping to `/healthz`) is required to pass similarly how, when running in a Deployment, a health probe is configured.

### Build and packaging

The Cluster API system is formed of a set of binaries. The core Cluster API manager, and the infrastructure providers are built using Go Modules in the `cluster-api` folder.

The binaries are built and packaged during the standard installer build process, `hack/build.sh`. Cluster API specific build flow is contained in the `hack/build-cluster-api.sh` script:
- Only enabled if the `OPENSHIFT_INSTALL_CLUSTER_API` environment variable is set.
- Builds (as needed) every binary listed as a Go Module in  the `cluster-api` folder.
- Downloads (as needed) the specified version of `envtest` to pacakge `kube-apiserver` and `etcd`.
- Produces a single `cluster-api.zip` file which is then copied in `pkg/clusterapi/mirror`.

To build an `openshift-install` binary with Cluster API bundled:
- Set `export OPENSHIFT_INSTALL_CLUSTER_API=y`
    - Optionally `export SKIP_TERRAFORM=y` if you don't need to use Terraform.
- Run `./hack/build.sh`, the binary is then produced in `bin/openshift-install`.

### Platform Onboarding

This is a general guide for onboarding a new platform to provision infrastructure with cluster api.

#### CAPI Cloud Provider Controller Dependencies

Providers are vendored to [`cluster-api/providers/<name>`][providerDir]. The README contains
detailed guidelines on how to vendor a new platform.

#### CAPI Cloud Provider CRDs

To add the CRDs for the CAPI cloud provider to the local control plane, add them to the
[`/data/data/cluster-api`][crdDir] dir. Typically these CRDs are published as an
`infrastructure-components.yaml` file as part of the release of the CAPI controller. The
yaml files in this directory will automatically be added to the Installer's local
control plane.

#### Run the CAPI Controller

In order for the Installer to run the controller, which was added in a previous step, add
the details for the controller in [this switch case in pkg/clusterapi/system.go][runController].
The controller will be unpacked and run based on the provider name; flags and environment variables
can be specified.

#### Create Cluster API Manifests

Next you need to add manifests to configure the infrastructure provisioned by the provider. The
cluster and any additional non-machine manifests should be added to platform-specific subpackages
of the `manifests` package. Introduce your platform in the [platform switch statement][clusterManifests]
and then add your implementation (see [AWS for an example][awsCluster].)

Machine manifests should be added similarly but in the [machines package][machineManifests]. If there is
existing code for your platform to generate master machines, it will serve as a good foundation for
generating these manifests. Note that when generating machine manifests, there will be a pair:
a general CAPI machine manifest, which handles the common bootstrapping pattern, and a platform-specific
machine manifest. [AWS machine manifests][awsMachines] offer an example.

NOTE: At the time of writing, the Installer will not produce manifests when you run
`openshift-install create manifests`. You can enable the Installer to write these to disk by adding
the `&clusterapi.Cluster{}` and the `&machines.ClusterAPI{}` targets in the [`manifests`][manifestsTarget],
similar to this [PR][manifestsPR]. This functionality has not merged yet, because it has been buggy
when running `create manifests` in conjunction with `create cluster`.

#### Implement Cluster API Provider Interface

Finally, add a [cluster-api provider interface][capiInterface] implementation for the platform. Initially
it should be fine to create an empty, simple provider which only implements the simple interface. Add
your new provider, behind a feature gate, to the [platform.go switch][platform] and add the exact
same code in [platform_altinfra.go][altinfra]. Make sure to use the `FeatureGateInstallClusterAPI`
feature gate:

```go
		if fg.Enabled(configv1.FeatureGateInstallClusterAPI) {
            //return interface implementation
		}
```

The `platform_altinfra.go` file uses a build tag to produce a separate altinfra image, which is
free of Terraform. We are also using this to expedite testing of cluster API while we implement
build dependencies needed in the main installer image.

#### Extend Cluster API Provider Interface

Once the "empty" provider is in place, that should be sufficient to allow the Installer to
apply the cluster API manifests and start provisioning resources! The Installer defines
additional hooks which you can implement in the provider to setup additional prerequisites
or provision resources out of band. The enhancement discusses this in more detail. Here is a
summary from the enhancement which describes the hooks, and how they run in relation to
the CAPI provisioning:

* PreProvision Hook - handle prequisites: e.g. IAM role creation, OS Disk upload
* Provision Infrastructure - create `cluster` manifest on local control plane
* InfraReady Hook - handle requirements in between cluster infrastructure and machine provisioning
* Ignition Hook - similar to InfraReady but specifically for generating (bootstrap) ignition
* Provision Machines - create bootstrap ignition and control-plane machines
* PostProvision Hook - post-provision requirements, such as DNS record creation

You can extend your platform implementation to implement these hooks as needed.

[enhancement]: https://github.com/openshift/enhancements/pull/1528 #TODO(padillon): update when merged
[providerDir]: https://github.com/openshift/installer/tree/master/cluster-api/providers
[crdDir]: https://github.com/openshift/installer/tree/master/data/data/cluster-api
[runController]: https://github.com/openshift/installer/blob/master/pkg/clusterapi/system.go#L120
[clusterManifests]: https://github.com/openshift/installer/blob/master/pkg/asset/manifests/clusterapi/cluster.go
[awsCluster]: https://github.com/openshift/installer/blob/master/pkg/asset/manifests/aws/cluster.go
[machineManifests]: https://github.com/openshift/installer/blob/master/pkg/asset/machines/clusterapi.go
[awsMachines]: https://github.com/openshift/installer/blob/master/pkg/asset/machines/aws/awsmachines.go
[manifestsTarget]: https://github.com/openshift/installer/blob/master/pkg/asset/targets/targets.go#L26
[manifestsPR]: https://github.com/openshift/installer/pull/7774/files
[capiInterface]: https://github.com/openshift/installer/blob/master/pkg/infrastructure/clusterapi/types.go#L15-L18
[platform]: https://github.com/openshift/installer/blob/master/pkg/infrastructure/platform/platform.go
[altinfra]: https://github.com/openshift/installer/blob/master/pkg/infrastructure/platform/platform_altinfra.go
