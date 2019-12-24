# Managing Dependencies

## Build Dependencies

The following dependencies must be installed on your system before you can build the installer.

### Fedora, CentOS, RHEL

```sh
sudo yum install golang-bin gcc-c++
```

If you need support for [libvirt destroy](libvirt/README.md#cleanup), you should also install `libvirt-devel`.

### Go

We follow a hard flattening approach; i.e. direct and inherited dependencies are installed in the base `vendor/`.

Dependencies are managed with [Go Modules](https://github.com/golang/go/wiki/Modules) but committed directly to the repository.

We require at least Go 1.13.

- Add or update a dependency with `go get <dependency>@<version>`.
- If you want to use a fork of a project or ensure that a dependency is not updated even when another dependency requires a newer version of it, manually add a [replace directive in the go.mod file](https://github.com/golang/go/wiki/Modules#when-should-i-use-the-replace-directive). 
- Run `go mod tidy` to tidy `go.mod` and update `go.sum`, then commit the changes.
- Run `go mod vendor` to re-vendor and then commit updated vendored code separately.

This [guide](https://github.com/golang/go/wiki/Modules#how-to-use-modules) is a great source to learn more about using `go mod`.

For the sake of your fellow reviewers, commit vendored code separately from any other changes.
