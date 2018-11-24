# Managing Dependencies

## Build Dependencies

The following dependencies must be installed on your system before you can build the installer.

### Fedora

```sh
sudo yum install golang-bin gcc-c++
```

### CentOS 7

```sh
sudo yum-config-manager --enable extras
sudo yum install -y centos-release-scl
sudo yum install -y go-toolset-7-golang-bin gcc-c++
# open a shell where Go 1.10 from SCL is in PATH
scl enable go-toolset-7 bash
```

### Red Hat Enterprise Linux 7

You need to install [go-toolset-1.10](https://access.redhat.com/documentation/en-us/red_hat_developer_tools/2018.4/html/using_go_toolset/).

```sh
sudo yum-config-manager --enable rhel-7-server-devtools-rpms
sudo yum-config-manager --enable rhel-server-rhscl-8-rpms
sudo yum install -y go-toolset-1.10-golang-bin gcc-c++
# open a shell where Go 1.10 from SCL is in PATH
scl enable go-toolset-1.10 bash
```

Note: If you have a `workstation` variant of RHEL, then reporitory name is `rhel-7-workstation-devtools-rpms` as documented in the link above.

### libvirt

If you need support for [libvirt destroy](libvirt/README.md#cleanup), you should also install `libvirt-devel`.

### Go

We follow a hard flattening approach; i.e. direct and inherited dependencies are installed in the base `vendor/`.

Dependencies are managed with [dep](https://golang.github.io/dep/) but committed directly to the repository. If you don't have dep, install the latest release from [Installation](https://golang.github.io/dep/docs/installation.html) link.

We require at least following version for dep:

```console
$ dep version
dep:
 version     : v0.5.0
 build date  : 2018-07-26
 git hash    : 224a564
 go version  : go1.10.3
```

To add a new dependency:

- Edit the `Gopkg.toml` file to add your dependency.
- Ensure you add a `version` field for the tag or the `revision` field for commit id you want to pin to.
- Revendor the dependencies:

```sh
dep ensure
```

This [guide](https://golang.github.io/dep/docs/daily-dep.html) a great source to learn more about using `dep` is .

For the sake of your fellow reviewers, commit vendored code separately from any other changes.

## Tests

See [tests/README.md](../../tests/README.md).
