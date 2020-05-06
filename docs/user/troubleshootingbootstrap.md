# Troubleshooting Bootstrap Failures

Unfortunately, there will always be some cases where OpenShift fails to install properly. In these events, it is helpful to understand the likely failure modes as well as how to troubleshoot the failure.

## Gathering bootstrap failure logs

### Using the installer provisioned workflow

When users are using the installer to create the OpenShift cluster, the installer has all the information to automatically capture the logs from bootstrap host in case of failure.

#### Authenticating to bootstrap host for ipi

The installer will use the user's environment to discover the credentials to connect to the bootstrap host over SSH. One of the following methods is used by the installer,

1. Use the user's already setup `SSH_AGENT`. If the user has a ssh-agent setup, the installer will use it for SSH authentication.

2. Use the user's home directory, `~/.ssh` on Linux hosts, to load all the SSH private keys and use those for SSH authentication.
    a. The installer also configures the bootstrap host with a *generated* SSH key, and this private key will be used for SSH authentication if none of the user keys are trusted.
    The installer only configures the bootstrap host to trust the generated key, and therefore the log bundle will only contain the logs from the bootstrap host and not the control-plane hosts.

### Using the user provisioned workflow

When users are creating the infrastructure for the OpenShift cluster and the cluster fails to bootstrap, the users can use the `gather bootstrap` subcommand to gather the logs from the bootstrap host.

```console
$ openshift-install gather bootstrap --help
Gather debugging data for a failing-to-bootstrap control plane

Usage:
  openshift-install gather bootstrap [flags]

Flags:
      --bootstrap string     Hostname or IP of the bootstrap host
  -h, --help                 help for bootstrap
      --key stringArray      Path to SSH private keys that should be used for authentication. If no key was provided, SSH private keys from user's environment will be used
      --master stringArray   Hostnames or IPs of all control plane hosts
```

An example of a invocation for a cluster with three control-plane machines would be,

```sh
openshift-install gather bootstrap --bootstrap ${BOOTSTRAP_HOST_IP} --master ${CONTROL_PLANE_1_HOST_IP} --master ${CONTROL_PLANE_2_HOST_IP} --master ${CONTROL_PLANE_3_HOST_IP}
```

#### Authenticating to bootstrap host for upi

When explicitly using the `gather bootstrap` subcommand, user can either utilize the installer's discovery mechanism like detailed [above](#authenticating-with bootstrap host-for-ipi) or provide the keys using the `--key` flag.

An example of a invocation for a cluster with three control-plane machines would be,

```sh
openshift-install gather bootstrap --key ${KEY_1} --key ${KEY_2} --bootstrap ${BOOTSTRAP_HOST_IP} --master ${CONTROL_PLANE_1_HOST_IP} --master ${CONTROL_PLANE_2_HOST_IP} --master ${CONTROL_PLANE_3_HOST_IP}
```

## Understanding the bootstrap failure log bundle

Here's what a log bundle looks like,

```console
.
├── bootstrap
├── control-plane
├── failed-units.txt
├── rendered-assets
├── resources
└── unit-status

5 directories, 1 file
```

### file: failed-units.txt

The failed-units.txt contains a list of all the **failed** systemd units on the bootstrap host.

### directory: unit-status

The unit-status directory contains the details of each failed systemd unit from [failed-units](#file-failed-units-txt),

### directory: bootstrap

The bootstrap directory consists of all the important logs and files from the bootstrap host. There are three subdirectories for the bootstrap host

```console
bootstrap
├── containers
├── journals
└── pods

3 directories, 0 files
```

#### directory: bootstrap/containers

The containers directory contains the descriptions and logs from all the containers created by the kubelet using CRI-O for the static pods.
This directory contains all the operators or their operands running on the bootstrap host in special bootstrap modes. For example the machine-config-server container, or the bootstrap-kube-controlplane pods etc.

For each container the directory has two files,

* `<human readable id>.log`, which contains the log of the container.
* `<human readable id>.inspect`, which containts the information about the container like the image, volume mounts, arguments etc.

#### directory: bootstrap/journals

The journals directory contains the logs for *important* systemd units. These units are,

* `release-image.log`, the release-image unit is responsible for pulling the Release Image to the bootstrap host.
* `crio-configure.log` and `crio.log`, these units are responsible for configuring the CRI-O on the bootstrap host and CRI-O daemon respectively.
* `kubelet.log`, the kubelet service is responsible for running the kubelet on the bootstrap host. The kubelet on the bootstrap host is responsible for running the static pods for etcd, bootstrap-kube-controlplane and various other operators in bootstrap mode.
* `approve-csr.log`, the approve-csr unit is responsible for allowing control-plane machines to join OpenShift cluster. This unit performs the job of in-cluster approver while the bootstrapping is in progress.
* `bootkube.log`, the bootkube service is the unit that performs the bootstrapping of OpenShift clusters using all the operators. This service is respnsible for running all the required steps to bootstrap the API and then wait for success.

There might also be other services that are important for some platforms like OpenStack, that will have logs in this directory.

#### directory: bootstrap/pods

The pods directory contains the information and logs from all the render commands for various operators run by the bootkube unit.

For each container the directory has two files,

* `<human readable id>.log`, which contains the log of the container.
* `<human readable id>.inspect`, which containts the information about the container like the image, volume mounts, arguments etc.

### directory: resources

The resources directory contains various Kubernetes objects that are present in the cluster. These resources are pulled using the bootstrap API running on the bootstrap host.

### directory: rendered-assets

The rendered-assets directory contains all the files and directories created by the bootkube unit using various render command for operators. This directory is a snapshot of the `/opt/openshift` directory on the bootstrap-host.

### directory: control-plane

The control-plane directory contains logs for each control-plane host. It contains a sub directory for each control-plane host, usually the IP address of the hosts.

```console
control-plane
├── 10.0.128.114
│   ├── containers
│   ├── failed-units.txt
│   ├── journals
│   └── unit-status
├── 10.0.142.138
│   ├── containers
│   ├── failed-units.txt
│   ├── journals
│   └── unit-status
└── 10.0.148.48
    ├── containers
    ├── failed-units.txt
    ├── journals
    └── unit-status

12 directories, 3 files
```

#### directory: control-plane/name/containers

The containers directory contains the descriptions and logs from all the containers created by the kubelet using CRI-O on the control-plane host. The files are same as [containers directory](#directory-bootstrap-containers) on bootstrap host.

#### directory: control-plane/name/journals

The journals directory contains the logs of **important** units on the control plane hosts. The list of such units is,

* `crio.log`
* `kubelet.log`
* `machine-config-daemon-host.log` and `pivot.log`, these files have logs for RHCOS pivot related actions on the control plane host.

## Common Failures

Here are some common failures that the users can troubleshoot using the bootstrap failure log bundle.

### Unable to pull the bootstrap failure logs

1. `Attempted to gather debug logs after installation failure: failed to create SSH client: failed to initialize the SSH agent: no keys found for SSH agent`
The installer tried to create a new SSH agent, but there were no keys found in user's home directory, usually `~/.ssh` on Linux. The user can use the `--key` flag to provide the private key for SSH to gather the bootstrap failure logs.

2. `failed to create SSH client: ssh: handshake failed: ssh: unable to authenticate, attempted methods [none publickey], no supported methods remain`
The keys provided to the installer from the `SSH_AGENT` or the keys loaded from user's home directory do not have permission to SSH to the bootstrap host. The user can use the `--key` flag to provide the private key for SSH to gather the bootstrap failure logs.

### Unable to pull Release Image

When the pull secret provided to the installer does not have correct permissions to pull the Release Image, the `bootstrap/journals/release-image.log` should contain the debugging logs.

For example,

```txt
-- Logs begin at Fri 2020-04-24 17:08:15 UTC, end at Fri 2020-04-24 17:33:16 UTC. --
Apr 24 17:08:46 ci-op-2cbvx-bootstrap.c.openshift-gce-devel-ci.internal systemd[1]: Starting Download the OpenShift Release Image...
Apr 24 17:08:46 ci-op-2cbvx-bootstrap.c.openshift-gce-devel-ci.internal release-image-download.sh[1688]: Pulling registry.svc.ci.openshift.org/ci-op-8dv01g3m/release@sha256:50b07a8b4529d8fd2ac6c23ecc311034a3b86cada41c948baaced8c6a46077bc...
Apr 24 17:08:49 ci-op-2cbvx-bootstrap.c.openshift-gce-devel-ci.internal podman[1698]: 2020-04-24 17:08:49.307961668 +0000 UTC m=+1.119158273 system refresh
Apr 24 17:08:49 ci-op-2cbvx-bootstrap.c.openshift-gce-devel-ci.internal release-image-download.sh[1688]: Error: error pulling image "registry.svc.ci.openshift.org/ci-op-8dv01g3m/release@sha256:50b07a8b4529d8fd2ac6c23ecc311034a3b86cada41c948baaced8c6a46077bc": unable to pull registry.svc.ci.openshift.org/ci-op-8dv01g3m/release@sha256:50b07a8b4529d8fd2ac6c23ecc311034a3b86cada41c948baaced8c6a46077bc: unable to pull image: Error initializing source docker://registry.svc.ci.openshift.org/ci-op-8dv01g3m/release@sha256:50b07a8b4529d8fd2ac6c23ecc311034a3b86cada41c948baaced8c6a46077bc: Error reading manifest sha256:50b07a8b4529d8fd2ac6c23ecc311034a3b86cada41c948baaced8c6a46077bc in registry.svc.ci.openshift.org/ci-op-8dv01g3m/release: unauthorized: authentication required
```

### Bootkube logs are empty

For cases where the bootkube logs are empty in `bootstrap/journals/bootkube.log` like,

```txt
-- Logs begin at Fri 2020-04-24 17:08:15 UTC, end at Fri 2020-04-24 17:33:16 UTC. --
-- No entries --
```

There is high likelihood that the Release Image cannot be downloaded and more details can be found using [release-image.log](#unable-to-pull-release-image)

## Control-plane logs missing from log bundle

When the control-plane logs are missing from the log bundle, for example,

```console
$ tree control-plane -L 2
control-plane
├── 10.0.0.4
├── 10.0.0.5
└── 10.0.0.6

3 directories, 0 files
```

The troubleshooting would require the logs of the installer gathering the log bundle, which are easily availble in `.openshift_install.log`.
