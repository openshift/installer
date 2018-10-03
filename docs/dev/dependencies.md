# Managing Dependencies

## Build Dependencies

The following dependencies must be installed on your system before you can build the installer.

### Fedora

```sh
sudo dnf install golang-bin gcc-c++
```

If you need support for [libvirt destroy](libvirt-howto#cleanup), you should also install `libvirt-devel`.

### CentOS, RHEL

```sh
sudo yum install golang-bin gcc-c++
```

If you need support for [libvirt destroy](libvirt-howto#cleanup), you should also install `libvirt-devel`.

## Go

We follow a hard flattening approach; i.e. direct and inherited dependencies are installed in the base `vendor/`.

Dependencies are managed with [dep](https://golang.github.io/dep/) but committed directly to the repository. If you don't have dep, install the latest release from [Installation](https://golang.github.io/dep/docs/installation.html) link.

We require atleast following version for dep:

```
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
