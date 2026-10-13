# Cluster API

The installer uses Cluster API controllers through a local control plane powered by `kube-apiserver` and `etcd` running locally.

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
- Builds (as needed) every binary listed as a Go Module in  the `cluster-api` folder.
- Downloads (as needed) the specified version of `envtest` to package `kube-apiserver` and `etcd`.
- Produces a single `cluster-api.zip` file which is then copied in `pkg/clusterapi/mirror`.

To build an `openshift-install` binary with Cluster API bundled:
- Optionally `export SKIP_TERRAFORM=y` if you don't need to use Terraform.
- Run `./hack/build.sh`, the binary is then produced in `bin/openshift-install`.
