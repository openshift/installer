# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## 0.4.0 - 2018-11-22

### Added

- The creation targets have been moved below a new `create` subcommand
  (e.g. `openshift-install create cluster` instead of the old
  `openshift-install cluster`).  This makes them easier to distinguish
  from other `openshift-install` subcommands and also mirrors the
  approach taken by `destroy` in 0.3.0.
- A new `manifest-templates` target has been added to `create`,
  allowing users to edit templates and have descendant assets
  generated from their altered templates during [a staged
  install](docs/user/overview.md#multiple-invocations).
- [The ingress operator][ingress-operator] is no longer masked.  The
  old Tectonic ingress operator has been removed.
- The [the registry operator][registry-operator] has been added, and
  the kube-addon operator which used to provide a registry (among
  other things) has been removed.
- The [checkpointer operator][checkpointer-operator] is no longer
  masked.  It runs on the production cluster, but not on the bootstrap
  node.
- Cloud credentials are now pushed into a secret where they can be
  consumed by cluster-API operators and other tools.
- OpenStack now has `destroy` support.
- We log verbosely to `${INSTALL_DIR}/.openshift_install.log` for most
  operations, giving access to the logs for troubleshooting even if
  you neglected to run with `--log-level=debug`.
- We've grown [troubleshooting
  documentation](docs/user/troubleshooting.md).

### Changed

- The `create cluster` subcommand now waits for the
  `bootstrap-complete` event and automatically removes the bootstrap
  assets after receiving it.  This means that after `create cluster`
  returns successfully, the cluster has its production control plane
  and topology (although there may still be operators working through
  their initialization).  The `bootstrap-complete` event was new in
  0.3.0, and it is now pushed at the appropriate time (it was too
  early in 0.3.0).  The `destroy bootstrap` subcommand is still
  available, to allow users to manually trigger bootstrap deletion if
  the automatic removal fails for whatever reason.
- On AWS, bootstrap deletion now also removes the S3 bucket used for
  the bootstrap node's Igntition configuration.
- Asset state is preserved even while moving backwards through [a
  staged install](docs/user/overview.md#multiple-invocations).  For
  example:

    ```sh
    openshift-install --dir=example create ignition-configs
    openshift-install --dir=example create install-config
    ```

    now preserves the full state including the generated Ignition
    configuration.  In 0.3.0, the `install-config` call would have
    removed the Ignition configuration and other downstream assets
    from the stored state.
- Some asset state is removed by successful `destroy cluster` runs.
  This reduces the change of contaminating future cluster creation
  with assets left over from a previous cluster, but users are [still
  encouraged](README.md#cleanup) to remove state between clusters to
  avoid accidentally contaminating the subsequent cluster's state.
- etcd discovery now happens via `SRV` records.  On libvirt, this
  requires a new Terraform provider, so users with older providers
  should [install a newer
  version](docs/dev/libvirt-howto.md#install-the-terraform-provider).
  This also allows all masters to use a single Ignition file.
- On AWS, the API and service load balancers have been changed from
  [classic load balancers][aws-elb] to [network load
  balancers][aws-nlb].  This should avoid [some latency issues we were
  seeing with classic load balancers][aws-elb-latency], and network
  load balancers are cheaper.
- On AWS, master `Machine` entries now include load balancer
  references, ensuring that new masters created by [the AWS
  cluster-API provider][cluster-api-provider-aws] will be attached to
  the load balancers.
- On AWS and OpenStack, the default network CIDRs have changed to
  172.30.0.0/16 for services and 10.128.0.0/14 for the cluster, to be
  consistent with previous versions of OpenStack.
- The bootstrap kubelet is no longer part of the production cluster.
  This reduces complexity and keeps production pods off of the
  temporary bootstrap node.
- [The cluster-version operator][cluster-version-operator] now runs in
  a static pod on the bootstrap node until the production control
  plane comes up.  This breaks a cyclic dependency between the
  production API server and operators.
- The bootstrap control plane now waits for some core pods to come up
  before exiting.
- [The machine-API operator][machine-api-operator] now reads the
  install-config from the `cluster-config-v1` config-map, instead of
  from an operator-specific configuration.
- AWS AMIs and libvirt images are now pulled from the new [RHCOS
  pipeline][rhcos-pipeline].
- Updated the security contact information for CoreOS -> Red Hat.
- We push a `ClusterVersion` custom resource.  The old `CVOConfig` is
  still being pushed, but it is deprecated.
- OpenStack credentials are loaded from standard system paths.
- On AWS and OpenStack, ports 9000-9999 are now open for host network
  services.
- Lots of doc and internal cleanup and minor fixes.

### Fixed

- On AWS, `destroy cluster` is now more robust, removing resources with
  either the `tectonicClusterID` or `kubernetes.io/cluster/<name>:
  owned` tags.  It also removes pending instances as well (it used to
  only remove running instances).
- On libvirt, `destroy cluster` is now more precise, only removing
  resources which are prefixed by the cluster name.
- Bootstrap Ignition edits (via `create ignition-configs`) no longer
  suffer from a `worker.ign` dependency cycle, which had been
  clobbering manual `bootstrap.ign` changes.
- The state-purging implementation respects `--dir`, avoiding `remove
  ...: no such file or directory` errors during [staged
  installs](docs/user/overview.md#multiple-invocations).
- Cross-filesystem Terraform state recovery during `destroy bootstrap`
  no longer raises `invalid cross-device link`.
- Bootstrap binaries are now located under `/usr/local/bin`, avoiding
  SELinux violations on RHEL 8.

### Removed

- All the old Tectonic operators and the `tectonic-system` namespace
  have been removed.
- On libvirt, the image URI prompt has been removed.  You can still
  control this via the `OPENSHIFT_INSTALL_LIBVIRT_IMAGE` environment
  variable, but too many users were breaking their cluster by pointing
  the installer at an outdated RHCOS, so we removed the prompt to make
  that knob less obvious.
- On libvirt, we've removed `.gz` suffix handling for images.  The new
  RHCOS pipeline supports `Content-Encoding: gzip`, so the
  suffix-based hack is no longer necessary.
- The `destroy-cluster` command, which was deprecated in favor of
  `destroy cluster` in 0.3.0, has been removed.
- The creation target subcommands of `openshift-install` have been
  removed.  Use the target subcommands of `create` instead
  (e.g. `openshift-install create cluster` instead of
  `openshift-install cluster`).

## 0.3.0 - 2018-10-22

### Added

- Asset state is loaded from the install directory, allowing for a [staged
  install](docs/user/overview.md#multiple-invocations).
- A new `openshift-install destroy bootstrap` command destroys the
  bootstrap resources.  Ideally, this would be safe to run after the
  new `bootstrap-complete` event is pushed to the `kube-system`
  namespace, but there is currently a bug causing that event to be
  pushed too early.  For now, you're on your own figuring out when to
  call this command.

    For consistency, the old `destroy-cluster` has been deprecated in
    favor of `openshift-install destroy cluster`.

- The installer creates worker `MachineSet`s, instead of leaving that to
  [the machine-API operator][machine-api-operator].
- The installer creates master `Machine`s and tags masters to be
  picked up by the [AWS cluster-API
  provider][cluster-api-provider-aws].

### Changed

- The installer now respects the `AWS_PROFILE` environment variable
  when launching AWS clusters.
- Worker subnets are now created in the appropriate availability zone
  for AWS clusters.
- Use the released hyperkube and hypershift instead of hard-coded
  images.
- Lots of changes to keep up with the advancing release image, as
  OpenShift operators are added to control various cluster components.
- Lots of internal cleanup and minor fixes.

### Removed

- The Tectonic kube-core operator, which has been replaced by
  OpenShift operators.

## 0.2.0 - 2018-10-12

### Added

- Asset state is preserved between invocations, allowing for a staged
    install like:

    ```console
    $ openshift-install --dir=example install-config
    $ openshift-install --dir=example cluster
    ```

    which creates a cluster using the same data given in the
    install-config (including the same random cluster ID, etc.).
- [The kube-apiserver][kube-apiserver-operator] and
  [kube-controller-manager][kube-controller-manager-operator]
  operators are called to render additional cluster manifests.
- etcd is now available as a service in the `kube-system` namespace,
  and the new service is labeled so [Prometheus][] will scrape it.
- The `service-serving-cert-signer-signing-key` secret is now
  available in the `openshift-service-cert-signer` namespace, which
  gives [the service-serving cert signer][service-serving-cert-signer]
  the keys it needs to mint and manage certificates for Kubernetes
  services.
- The etcd-serving certificate is now passed through to [the
  kube-controller-manager operator][kube-controller-manager-operator].
- We disable some components which [the cluster-version
  operator][cluster-version-operator] would otherwise install but
  which conflict with the legacy tectonic-operators.
- The new `openshift-install graph` outputs the asset graph in [the
  DOT language][dot].
- `openshift-install version` now outputs the Terraform version as
  well as the installer version.

### Changed

- The [cluster-version operator][cluster-version-operator] is no
  longer run as a static pod.  Instead, we just wait until the control
  plane comes up and run it them.
- Terraform errors are logged to standard error even when
  `--log-level` is less than `debug`.
- Terraform is now invoked with `-no-color` and `-input=false`.
- The `cluster` target now includes both launching the cluster and
  populating `metadata.json`, regardless of whether the `terraform`
  invocation succeeds.  This allows `destroy-cluster` to cleanup
  cluster resources even when the `terraform` invocation fails.
- Reported errors now include more context, making them less
  enigmatic.
- Libvirt image caching is more efficient, caching unzipped images
  with a cache that grows by one unzipped image per RHCOS release in
  `$XDG_CACHE_HOME/openshift-install/libvirt/image`.  The previous
  implementation unzipped, when necessary, for every launched cluster,
  which was slow.  And the previous implementation added one unzipped
  image to `/tmp` per cluster launch, which consumed more disk space.
- Work continues on the OpenStack platform.
- Lots of internal cleanup, especially around asset generation.

### Removed

- The operatorstatus CRD.  Now [the cluster-version
  operator][cluster-version-operator] creates this on its own.
- The `machine-config-operator-images` config-map.  Now [the
  cluster-version operator][cluster-version-operator] pulls these from
  [the machine-config images][machine-config-operator].
- The `machine-api` app-version from the `tectonic-system` namespace.

## 0.1.0 - 2018-10-02

### Added

The `openshift-install` command.  This moves us to the new
install-config approach with [asset
generation](docs/design/assetgeneration.md) in Go instead of in
Terraform.  Terraform is still used to push the assets out to
resources on the backing platform (AWS, libvirt, or OpenStack), but
that push happens in a single Terraform invocation instead of in
multiple steps.  This makes installation faster, because more
resources can be created in parallel.  `openshift-install` also
dispenses with the distribution tarball; all required assets except
for a `terraform` binary are distributed in the `openshift-install`
binary.

The configuration and command-line interface are quite different, so
previous `tectonic` users are encouraged to start from scratch when
getting acquainted with `openshift-install`.  AWS users should look
[here](README.md#quick-start).  Libvirt users should look
[here](docs/dev/libvirt-howto.md).  The new `openshift-install` also
includes an interactive configuration generator, so you can launch the
installer and follow along as it guides you through the process.

### Removed

The `tectonic` command and tarball distribution are gone.  Please use
the new `openshift-install` command instead.

[aws-elb]: https://docs.aws.amazon.com/elasticloadbalancing/latest/classic/introduction.html
[aws-elb-latency]: https://github.com/openshift/installer/pull/594#issue-227786691
[aws-nlb]: https://docs.aws.amazon.com/elasticloadbalancing/latest/network/introduction.html
[checkpointer-operator]: https://github.com/openshift/pod-checkpointer-operator
[cluster-api-provider-aws]: https://github.com/openshift/cluster-api-provider-aws
[cluster-version-operator]: https://github.com/openshift/cluster-version-operator
[dot]: https://www.graphviz.org/doc/info/lang.html
[ingress-operator]: https://github.com/openshift/cluster-ingress-operator
[kube-apiserver-operator]: https://github.com/openshift/cluster-kube-apiserver-operator
[kube-controller-manager-operator]: https://github.com/openshift/cluster-kube-controller-manager-operator
[machine-api-operator]: https://github.com/openshift/machine-api-operator
[machine-config-operator]: https://github.com/openshift/machine-config-operator
[Prometheus]: https://github.com/prometheus/prometheus
[registry-operator]: https://github.com/openshift/cluster-image-registry-operator
[rhcos-pipeline]: https://releases-rhcos.svc.ci.openshift.org/storage/releases/maipo/builds.json
[service-serving-cert-signer]: https://github.com/openshift/service-serving-cert-signer
