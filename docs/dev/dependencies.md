# Managing Dependencies

## Build Dependencies

The following dependencies must be installed on your system before you can build the installer.

### Fedora

```sh
sudo dnf install golang-bin gcc-c++
```

### CentOS, RHEL

```sh
sudo yum install golang-bin gcc-c++
```

## Go

We follow a hard flattening approach; i.e. direct and inherited dependencies are installed in the base `vendor/`.

Dependencies are managed with [glide](https://glide.sh/) but committed directly to the repository. If you don't have glide, install the latest release from [https://glide.sh/](https://glide.sh/). We require version 0.12 at a minimum.

The vendor directory is pruned using [glide-vc](https://github.com/sgotti/glide-vc). Follow the [installation instructions](https://github.com/sgotti/glide-vc#install) in the project's README.

To add a new dependency:
- Edit the `glide.yaml` file to add your dependency.
- Ensure you add a `version` field for the sha or tag you want to pin to.
- Revendor the dependencies:

```sh
rm glide.lock
glide install --strip-vendor
glide-vc --use-lock-file --no-tests --only-code
```

If it worked correctly it should:
- Clone your new dep to the `/vendor` dir and check out the ref you specified.
- Update `glide.lock` to include your new package, add any transitive dependencies and update its hash.

For the sake of your fellow reviewers, commit vendored code separately from any other changes.

## Tests

See [tests/README.md](../../tests/README.md).
