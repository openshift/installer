# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## [Unreleased]

### Changed

- On AWS and OpenStack, master ports 10251 (scheduler) and 10252
  (controller manager) have been opened to access from all machines.
  This allows Prometheus (which runs on the worker nodes) to scrape
  all machines for metrics.
- On AWS, the installer and subsequent cluster will now tag resources
  it creates with `openshiftClusterID`.  `tectonicClusterID` is
  deprecated.

## 0.6.0 - 2018-12-09

### Added

- We now push a `kubeadmin` user (with an internally-generated
  password) into the cluster for the new [bootstrap identity
  provider][bootstrap-identity-provider].  This gives users a way to
  access the web console, Prometheus, etc. without needing to
  configure a full-fledged identity provider or install `oc`.  The
  `create cluster` subcommand now blocks until the web-console route
  is available and then exits after printing instructions for using
  the new credentials.
- The installer binary now includes Terraform so there is no longer a
  need to separately download and install it.

### Changed

- The SSH public key configuration has moved a level up in the install
  config, now that the `admin` structure has been removed.
- `build.sh` now checks to make sure you have a new enough `go`,
  instead of erroring out partway through the build.
- We now resolve the update payload to a digest on the bootstrap node,
  so [the cluster-version-operator][cluster-version-operator] can
  figure out exactly which image we used.
- Creation logging has been overhauled to increase it's
  signal-to-noise while waiting for the Kubernetes API to come up.
- On AWS, the installer will now prompt you for an access key and
  secret if it cannot find your AWS credentials in the usual places.
- On AWS, the installer will look at `AWS_DEFAULT_REGION` and in other
  usual places when picking a default for the region prompt.  You
  still have to set `OPENSHIFT_INSTALL_AWS_REGION` if you want to skip
  the prompt entirely.
- On libvirt, we've bumped masters from 3 GiB of memory to 4 GiB to
  address out-of-memory issues we had been seeing at 3 GiB.
- Lots of doc and internal cleanup and minor fixes.

### Removed

- The old admin username and password inputs have been removed.  They
  weren't being used anyway, and their intended role has been replaced
  by the newly-added `kubeadmin` user and bootstrap identity provider.
- The old `openshift-web-console` namespace is gone.  The new console
  is in the `openshift-console` namespace.

## 0.5.0 - 2018-12-03

### Added

- We now push the ingress custom resource defintion and initial
  configuration, allowing the ingress operator to configure itself
  without referencing the deprecated `cluster-config-v1` resource.

### Changed

- Pull secret documentation now points to
  [try.openshift.com](https://try.openshift.com) for pull-secret
  acquisition, instead of pointing at `account.coreos.com`.  Users
  will need to update their pull secrets.
- If the automatic bootstrap teardown (which landed in 0.4.0) times
  out waiting for the `bootstrap-complete` event, the installer exits
  with a non-zero exit code.  We had ignored watcher timeouts in 0.4.0
  due to concerns about watcher robustness, but the current watcher
  code has been reliable in our continuous integration testing.
- The hard-coded `quay.io/coreos/bootkube` dependency has been
  replaced by the new [cluster-bootstrap][] image, which is referenced
  from the release image.
- The etcd service now uses [selectors][kube-selector] to determine
  the pods it exposes, and the explict etcd endpoints object is gone
  (replaced by the one Kubernetes maintains based on the selector).
- On AWS, both masters and worker have moved from t2.medium nodes
  m4.large nodes (more on AWS instance types
  [here][aws-instance-types]) to address CPU and memory constraints.
- On AWS, master volume size has been bumped from 30 GiB to 120 GiB to
  increase our baseline performance from on [gp2's sliding IOPS
  scale][aws-ebs-gp2-iops] from the 100 IOPS floor up to 360 IOPS.
  Volume information is not currently supported by [the cluster-API
  AWS provider's
  `AWSMachineProviderConfig`][cluster-api-provider-aws-012575c1-AWSMachineProviderConfig],
  so this change is currently limited to masters created by the
  installer.
- On Openstack, we now validate cloud, region, and image-name user
  input instead of blindly accepting entries.
- On libvirt, we pass Ignition information for masters and workers via
  secrets instead of passing a libvirt volume path.  This makes the
  libvirt approach consistent with how we already handle AWS and
  OpenStack.
- Lots of internal cleanup, especially around trimming dead code.

### Fixed

- The `.openshift_install.log` addition from 0.4.0 removed Terraform
  output from `--log-level=debug`.  We've fixed that in 0.5.0; now
  `.openshift_install.log` will always contain the full Terraform
  output, while standard error returns to containing the Terraform
  output if and only if `--log-level=debug` or higher.
- On AWS teardown, errors retrieving tags for S3 buckets and Route 53
  zones are no longer fatal.  This allows the teardown code to
  continue it's exponential backoff and try to remove the bucket or
  zone later.  It avoids some resource leaks we were seeing due to AWS
  rate limiting on those tag lookups as many simultaneous CI jobs
  searched for Route 53 zones with their cluster's tags.  We'll still
  hit those rate limits, but they no longer cause us to give up on
  reaping resources.
- On AWS, we've removed some unused data blocks, fixing occasional
  errors like:

        data.aws_route_table.worker.1: Your query returned no results.

- On OpenStack, similar retry-during-teardown changes were made for
  removing ports and for removing subnets from routers.
- On libvirt, Terraform no longer errors out when launching clusters
  configured for more than one master, fixing a bug from 0.4.0.

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

[aws-ebs-gp2-iops]: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/EBSVolumeTypes.html#EBSVolumeTypes_gp2
[aws-elb]: https://docs.aws.amazon.com/elasticloadbalancing/latest/classic/introduction.html
[aws-elb-latency]: https://github.com/openshift/installer/pull/594#issue-227786691
[aws-instance-types]: https://aws.amazon.com/ec2/instance-types/
[aws-nlb]: https://docs.aws.amazon.com/elasticloadbalancing/latest/network/introduction.html
[bootstrap-identity-provider]: https://github.com/openshift/origin/pull/21580
[checkpointer-operator]: https://github.com/openshift/pod-checkpointer-operator
[cluster-api-provider-aws]: https://github.com/openshift/cluster-api-provider-aws
[cluster-api-provider-aws-012575c1-AWSMachineProviderConfig]: https://github.com/openshift/cluster-api-provider-aws/blob/012575c1c8d758f81c979b0b2354950a2193ec1a/pkg/apis/awsproviderconfig/v1alpha1/awsmachineproviderconfig_types.go#L86-L139
[cluster-bootstrap]: https://github.com/openshift/cluster-bootstrap
[cluster-version-operator]: https://github.com/openshift/cluster-version-operator
[dot]: https://www.graphviz.org/doc/info/lang.html
[ingress-operator]: https://github.com/openshift/cluster-ingress-operator
[kube-apiserver-operator]: https://github.com/openshift/cluster-kube-apiserver-operator
[kube-controller-manager-operator]: https://github.com/openshift/cluster-kube-controller-manager-operator
[kube-selector]: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#label-selectors
[machine-api-operator]: https://github.com/openshift/machine-api-operator
[machine-config-operator]: https://github.com/openshift/machine-config-operator
[Prometheus]: https://github.com/prometheus/prometheus
[registry-operator]: https://github.com/openshift/cluster-image-registry-operator
[rhcos-pipeline]: https://releases-rhcos.svc.ci.openshift.org/storage/releases/maipo/builds.json
[service-serving-cert-signer]: https://github.com/openshift/service-serving-cert-signer
