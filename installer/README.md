# Tectonic Installer

The Tectonic Installer is an app for creating Tectonic clusters.

See [official installation documentation](https://coreos.com/tectonic/docs/latest/) if you'd like to use a published release.

Read on if you'd like to build and run the installer yourself.

## License

Get a [license](https://account.coreos.com) and follow the guides to create Tectonic clusters end to end.

## Build prerequisites

- [Go 1.8](https://golang.org/doc/install)
- The tectonic-installer repo must be located at `$GOPATH/src/github.com/coreos/tectonic-installer`

### Build / Run

All commands mentioned here must be run from the same working directory as this README file, `./installer/` from the root of this repo.

## Managing Dependencies

### Go

Dependencies are managed with [glide](https://glide.sh/), but committed directly to the repository.

If you don't have glide, install the latest release from [https://glide.sh/](https://glide.sh/). We require version 0.12 at a minimum.

The vendor directory is pruned using [glide-vc](https://github.com/sgotti/glide-vc). Follow the [installation instructions](https://github.com/sgotti/glide-vc#install) in the project's README.

To add a new dependency:

- Edit the `glide.yaml` file to add your dependency.
- Ensure you add a `version` field for the sha or tag you want to pin to.
- Revendor the dependencies:

```
rm glide.lock
glide install -v
glide-vc --use-lock-file --no-tests --only-code
git checkout vendor/golang.org/x/crypto/ssh/terminal/BUILD.bazel vendor/github.com/Sirupsen/logrus/BUILD.bazel
bazel run //:gazelle
```

If it worked correctly it should:
- Clone your new dep to the `/vendor` dir, and check out the ref you specified.
- Update `glide.lock` to include your new package, adds any transitive dependencies, and updates its hash.
- Regenerate BUILD.bazel files.

For the sake of your fellow reviewers, commit vendored code changes as a separate commit from any other changes.

#### Regenerate or Repair Vendored Code

Should you need to regenerate or repair the vendored code en-mass from their source repositories, you can run
the same steps above under "Revendor the dependencies".
