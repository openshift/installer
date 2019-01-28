# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## 0.11.0 - 2019-01-27

### Added

- On AWS, the installer creates [DHCP options][aws-dhcp-options] for
  the VPC to support internal unqualified-hostname resolution.  This
  works around some limitations with `oc rsh` and Kubernetes node
  registration in the face of inappropriate default DHCP options.  And
  because [the AWS `domain-name` logic is
  region-specific][aws-dhcp-options], there is no single DHCP options
  configuration that provides internal unqualified-hostname resolution
  for multiple regions.

### Changed

- On AWS, the installer now prompts for missing credentials even if
  you supplied an `install-config.yaml`.  Previously, only the
  install-config wizard would prompt.
- On OpenStack, the developer-only internal DNS server which was
  removed in 0.10.0 has been restored, because the approach taken in
  0.10.0 broke etcd cluster formation for some users.
- Several doc and internal cleanups.

### Fixed

- `openshift-install` has improved error handling for various invalid
  command lines.  It now errors when additional positional arguments
  are passed to commands that do not take positional arguments
  (previously those commands silently ignored the presence of
  positional arguments).  And it logs an error and exits 1 when an
  invalid value is provided to --log-level (previously it exited 1 but
  did not write to the standard error stream).
- The slow-input issues for the install-config wizard have been fixed.
- On AWS, `destroy cluster` fixed a bug in the 0.10.1 refactor which
  could lead to leaked resources and a claim of successful deletion if
  a call to get tagged resources failed (for example, because the
  caller lacked the `tag:GetResources` permission).
- On AWS, a new explicit dependency in the Terraform modules prevents
  errors like:

        * module.vpc.aws_lb.api_external: 1 error occurred:
        * aws_lb.api_external: Error creating Application Load Balancer: InvalidSubnet: VPC vpc-0765c67bbc82a1b7d has no internet gateway
        status code: 400, request id: 5a...d5

- On libvirt, the installer no longer holds the OS image in memory
  after it has been written to disk.  Ideally it would stream the OS
  image to disk instead of ever holding it in memory, but this fix
  mitigates our current in-memory buffering.

## 0.10.1 - 2019-01-22

### Changed

- `create ignition-configs` now also writes `metadata.json` to the
  asset directory, which allows [Hive][] to more reliably destroy
  clusters.
- `destroy cluster` now removes `.openshift_install_state.json` on
  success, clearing the way for future `create cluster` runs in the
  same asset directory.
- On AWS, we now default to m4.xlarge masters.  The increased CPU
  reduces etcd latencies, which in turn helps with cluster stability.
- On AWS, the bootstrap machine has a new security-group allowing
  journald-gateway and kubelet access, for easier debugging when
  bootstrapping fails.
- Several doc and internal cleanups.

### Removed

- The SSH public key is no longer inserted in the pointer Ignition
  configurations, now that authorized public keys are [managed by the
  machine-config daemon][machine-config-daemon-ssh-keys].

### Fixed

- On AWS, the cluster-API provider now supports configuring machine
  volumes, so `rootVolume` settings in `install-config.yaml` will be
  respected.
- On AWS, the generated Terraform variables no longer clobber master
  instance type and root volume configuration set via
  `install-config.yaml`.  You can now use:

    ```yaml
    machines:
    - name: master
      platform:
        aws:
          type: m5.large
          rootVolume:
            iops: 3000
            size: 220
            type: io1
      replicas: 3
    - name: worker
      ...
    ```

    and similar to successfully customize your master machines.
- On AWS, `delete cluster` has been adjusted to use more efficient
  tag-based lookup and fix several bugs due to previously-missing
  pagination.  This should address some issues we had been seeing with
  leaking AWS resources despite `delete cluster` claiming success.

## 0.10.0 - 2019-01-15

### Added

- The installer pushes an Infrastructure object to
  infrastructures.config.openshift.io with platform information.
  Cluster components should use this instead of the deprecated
  `cluster-config-v1` resource.
- `openshift-install` has a new `completion` subcommand, to generation
  shell-completion code (currently only for Bash).
- On AWS, `destroy cluster` now also removed IAM users with the usual
  tags.  We don't create these users yet, but the removal sets the
  stage for the coming credentials operator.

### Changed

- Install configuration now includes a new `apiVersion` property which
  must be set to `v1beta1`.  Future changes to the install-config
  schema will result in new versions, allowing new installers to
  continue to support older install-config schema (and older
  installers to error out when presented with newer install-config
  schema).  Changes to the schema since 0.9.0:

    - `clusterID` has been removed.  This should be a new UUID for
      every cluster, so there is no longer an easy way for users to
      configure it.
    - Image configuration has been removed.  Almost all users should
      be fine with the installer-chosen RHCOS.  Users who need to
      override the RHCOS build (because they're testing new RHCOS
      releases) can set a new `OPENSHIFT_INSTALL_OS_IMAGE_OVERRIDE`
      environment variable.
    - Per-platform network properties have been consolidated into the
      platform-agnostic `machineCIDR` from which node IP addresses are
      assigned.
    - On libvirt, all machine-pool configuration has been removed, as
      no remaining properties were supported by the libvirt
      cluster-API provider.

- `install-config.yaml` read during [staged
  installs](docs/user/overview.md#multiple-invocations) will now have
  installer-defaults applied for missing properties.  This allows you
  to set only the properties you are interested in overriding, and
  allow the installer to manage the remaining properties.
- `create ignition-configs` now also writes the admin kubeconfig to
  the asset directory, to support bring-your-own-infrastructure use
  cases.
- The bootstrap node now [serves
  journals](docs/user/troubleshooting.md#troubleshooting-the-bootstrap-node)
  for easier troubleshooting.
- The validity for the initial kubelet TLS certificate has been
  increased from one hour to 24 hours, to give
  bring-your-own-infrastructure users longer to manually distribute
  the certificate before it expires.
- The key for the root certificate authority is no longer pushed into
  the cluster (not even to the bootstrap node).
- Machine(set)s generated by the installer now use `providerSpec`
  instead of the deprecated `providerConfig`.
- On AWS, the load balancers now use HTTPS health checks to reduce log
  noise like:

        http: TLS handshake error from 10.0.20.86:28372: EOF

- On AWS, IAM roles are now tagged with the usual resource tags
  (`openshiftClusterID`, etc.).  Some other resources have had their
  tags updated to match those conventions (e.g. the internal Route 53
  hosted zone was updated from `KubernetesCluster` to
  `kubernetes.io/cluster/{name}: owned`).
- The OpenStack platform has been removed from the install-config
  wizard while it remains experimental.  It is still available for
  users who supply their own `install-config.yaml`.
- On OpenStack, the service VP now respects any SSH key specified in
  the install configuration.
- On OpenStack, a developer-only internal DNS server has been removed,
  so users need to configure additional records for the existing
  external DNS zone.
- On OpenStack, Neutron trunk ports are now used for VM network
  interfaces if Neutron supports them to support future Kuryr
  integration.
- On OpenStack, masters and workers have been consolidated in a single
  subnet to simplify the deployment.
- On OpenStack, the Ignition security group now only allows internal
  connections, and no longer allows connections from outside the
  cluster network.
- On OpenStack, the machine(set) templates have been updated to set
  `cloudName` and some other properties.
- On libvirt, `destroy cluster` is now more robust in the face of
  domains which were already shutdown.
- Lots of doc and internal cleanup and minor fixes.

### Removed

- Support for `install-config.yml` (deprecated in 0.8.0) has been
  removed.

### Fixed

- On AWS, domain pagination for the wizard's base-domain select widget
  has been fixed.  Previously, it would continuously fetch the first
  page of hosted zones (for accounts with multiple pages of zones)
  until it hit an error like:

    ```
    ERROR list hosted zones: Throttling: Rate exceeded
            status code: 400, request id: ...
    ```

    before falling back to a free-form base-domain input.

## 0.9.0 - 2019-01-05

### Added

- There is a new `none` platform for bring-your-own infrastructure
  users who want to generate Ignition configurations.  The new
  platform is mostly undocumented; users will usually interact with it
  via [OpenShift Ansible][openshift-ansible].

### Changed

- On OpenStack, there's no longer a default flavor, because flavor
  names are not standardized.  Instead, there's a new prompt to choose
  among flavors supported by the target OpenStack implementation.
- On libvirt, we now use the host-passthrough CPU type, which should
  improve performance for some use-cases.
- Some doc and internal cleanup and minor fixes.

## 0.8.0 - 2018-12-23

### Added

- The installer binary now includes all required Terraform plugins, so
  there is no longer a need to separately download and install them.
  This will be most noticeable on libvirt, where users used to be
  required to install the libvirt plugin manually.  This avoids issues
  with mismatched plugin versions (which we saw sometimes on libvirt)
  and network-connectivity issues during `create cluster` invocations
  (which we saw sometimes on all platforms).
- The configured base domain is now pushed into the cluster's
  `config.openshift.io` as a DNS custom resource.

### Changed

- `install-config.yml` is now `install-config.yaml` to align with our
  usual YAML extension.  `install-config.yml` is deprecated, and
  support for it will be removed completely in the next release.
- On AWS, we now use a select widget for the base-domain wizard
  prompt, making it easier to choose an existing public zone.
- On AWS, Route 53 rate limits during `cluster destroy` are now less
  disruptive, reducing the AWS request load in busy accounts.
- On OpenStack, the HAProxy configuration no longer hard-codes the
  cluster name and base domain.
- On OpenStack, the 0.7.0 fix for:

        FATAL Expected HTTP response code [202 204] when accessing [DELETE https://osp-xxxxx:13696/v2.0/routers/52093478-dcf1-4bcc-9a2c-dbb1e42da880], but got 409 instead
        {"NeutronError": {"message": "Router 52093478-dcf1-4bcc-9a2c-dbb1e42da880 still has ports", "type": "RouterInUse", "detail": ""}}

    was incorrect and has been reverted.  We'll land a real fix for
    this issue in future work.
- On OpenStack, the service VM from 0.7.0 now has a floating IP
  address.
- All libvirt functionality is behind `TAGS=libvirt` now.  Previously
  installer builds with `TAGS=libvirt_destroy` included all libvirt
  functionality, while builds without that tag would include `create
  cluster` but not `destroy cluster` functionality.  With the change,
  all users using the installer with libvirt clusters will need to set
  the new build tag.
- Lots of doc and internal cleanup and minor fixes.

### Removed

- On AWS and OpenStack, the `tectonicClusterID` tag which was
  deprecated in 0.7.0 has been removed.

## 0.7.0 - 2018-12-14

### Added

- We now validate install-config when loading it during [staged
  installs](docs/user/overview.md#multiple-invocations).  Previously
  we only validated that input when it was entered into the wizard or
  via environment variables.  This also leads to some changes in which
  values are considered valid:

    - Base domains may no longer contain uppercase letters.
    - Cluster names may now be longer than 63 characters, although as
      specified in [RFC 1123][rfc-1123-s2.1] host software may not
      support names longer than 63 characters.
    - Pull secrets require more content (e.g. it must contain an
      `auths` property).  Previously we only required pull secrets to
      be valid JSON.
    - SSH public keys must be parsable with
      [`ParseAuthorizedKey`][ssh.ParseAuthorizedKey].  Previously we
      had our own logic that was not as well developed.

- We've added `images/installer/Dockerfile.ci.rhel7` for building
  installer images on a RHEL base.
- On AWS, we now create [an S3 endpoint][aws-s3-endpoint] for the VPC.
- We've added OpenStack documentation.

### Changed

- The pull-secret prompt now masks the input to avoid leaking it to
  shoulder surfers and similar.
- The pull-secret prompt's help now points to
  [cloud.openshift.com](https://cloud.openshift.com/clusters/install#pull-secret)
  instead of [try.openshift.com](https://try.openshift.com).  This
  should make it easier to find the pull secret without digging
  through a page of introductory content.
- The initial kubeconfig inserted on master nodes used to have `admin`
  privileges, but only for 30 minutes.  Now it has role bindings that
  allow it to create and receive automatic approval for certificate
  signing requests, but it does not have additional privileges beyond
  that.
- On AWS and OpenStack, master ports 10251 (scheduler) and 10252
  (controller manager) have been opened to access from all machines.
  This allows Prometheus (which runs on the worker nodes) to scrape
  all machines for metrics.
- On AWS and OpenStack, the installer and subsequent cluster will now
  tag resources it creates with `openshiftClusterID`.
  `tectonicClusterID` is deprecated.
- On OpenStack, only the OpenStack `clouds` entry is marshalled into
  the `openstack-creds` secret.  Previously we had injected the host's
  entire cloud configuration.
- On OpenStack, there is now a service VM to provide DNS and load
  balancing for the OpenShift cluster.  The service VM will eventually
  be removed, but for now its a convenient hack to get usable clusters
  on OpenStack.
- On libvirt, we now document host DNS configuration as required,
  because too many users were skipping that step and then reporting
  errors with Kubernetes API detection when the install-host failed to
  resolve the cluster domain name while waiting for the
  `bootstrap-complete` event.
- Lots of doc and internal cleanup and minor fixes.

### Fixed

- Fixed OpenShift manifest loading during [staged
  installs](docs/user/overview.md#multiple-invocations).  The
  installer had been ignoring changes to those files since 0.4.0.
- Fixed `you must pass a pointer as the target of a Write operation`
  errors introduced in 0.6.0 for the AWS access key ID prompt.
- When `create cluster` times out waiting for the Kubernetes API, we
  now exit immediately.  Previously we'd wait through another 30
  minutes of failed event-listener connections before failing this
  case.  We've also fixed similar timeout detection for the code that
  waits for the OpenShift console route.
- On OpenStack, we've fixed a bug in router deletion:

        FATAL Expected HTTP response code [202 204] when accessing [DELETE https://osp-xxxxx:13696/v2.0/routers/52093478-dcf1-4bcc-9a2c-dbb1e42da880], but got 409 instead
        {"NeutronError": {"message": "Router 52093478-dcf1-4bcc-9a2c-dbb1e42da880 still has ports", "type": "RouterInUse", "detail": ""}}

- On libvirt, we've fixed a bug introduced in 0.6.0 and are now back
  to removing the bootstrap node from round-robin DNS when we destroy
  the bootstrap resources.

### Removed

- The user-facing `OPENSHIFT_INSTALL_*` environment variables are
  gone.  Instead, users who want to skip the wizard are encouraged to
  [provide their own
  install-config](docs/user/overview.md#multiple-invocations).
- On AWS, the option to install a cluster into an existing VPC is
  gone.  Users who would have previously done this can use [VPC
  peering][aws-vpc-peering].

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

- We now push the ingress custom resource definition and initial
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
  the bootstrap node's Ignition configuration.
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

[aws-dhcp-options]: https://docs.aws.amazon.com/vpc/latest/userguide/VPC_DHCP_Options.html
[aws-ebs-gp2-iops]: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/EBSVolumeTypes.html#EBSVolumeTypes_gp2
[aws-elb]: https://docs.aws.amazon.com/elasticloadbalancing/latest/classic/introduction.html
[aws-elb-latency]: https://github.com/openshift/installer/pull/594#issue-227786691
[aws-instance-types]: https://aws.amazon.com/ec2/instance-types/
[aws-nlb]: https://docs.aws.amazon.com/elasticloadbalancing/latest/network/introduction.html
[aws-s3-endpoint]: https://docs.aws.amazon.com/vpc/latest/userguide/vpc-endpoints-s3.html
[aws-vpc-peering]: https://docs.aws.amazon.com/vpc/latest/peering/what-is-vpc-peering.html
[bootstrap-identity-provider]: https://github.com/openshift/origin/pull/21580
[checkpointer-operator]: https://github.com/openshift/pod-checkpointer-operator
[cluster-api-provider-aws]: https://github.com/openshift/cluster-api-provider-aws
[cluster-api-provider-aws-012575c1-AWSMachineProviderConfig]: https://github.com/openshift/cluster-api-provider-aws/blob/012575c1c8d758f81c979b0b2354950a2193ec1a/pkg/apis/awsproviderconfig/v1alpha1/awsmachineproviderconfig_types.go#L86-L139
[cluster-bootstrap]: https://github.com/openshift/cluster-bootstrap
[cluster-version-operator]: https://github.com/openshift/cluster-version-operator
[dot]: https://www.graphviz.org/doc/info/lang.html
[Hive]: https://github.com/openshift/hive/
[ingress-operator]: https://github.com/openshift/cluster-ingress-operator
[kube-apiserver-operator]: https://github.com/openshift/cluster-kube-apiserver-operator
[kube-controller-manager-operator]: https://github.com/openshift/cluster-kube-controller-manager-operator
[kube-selector]: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#label-selectors
[machine-api-operator]: https://github.com/openshift/machine-api-operator
[machine-config-operator]: https://github.com/openshift/machine-config-operator
[machine-config-daemon-ssh-keys]: https://github.com/openshift/machine-config-operator/blob/master/docs/Update-SSHKeys.md
[openshift-ansible]: https://github.com/openshift/openshift-ansible
[Prometheus]: https://github.com/prometheus/prometheus
[ssh.ParseAuthorizedKey]: https://godoc.org/golang.org/x/crypto/ssh#ParseAuthorizedKey
[registry-operator]: https://github.com/openshift/cluster-image-registry-operator
[rfc-1123-s2.1]: https://tools.ietf.org/html/rfc1123#section-2
[rhcos-pipeline]: https://releases-rhcos.svc.ci.openshift.org/storage/releases/maipo/builds.json
[service-serving-cert-signer]: https://github.com/openshift/service-serving-cert-signer
