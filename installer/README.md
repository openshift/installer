# Tectonic Installer

The Tectonic Installer is an app for creating Tectonic clusters.

See [official installation documentation](https://coreos.com/tectonic/docs/latest/) if you'd like to use a published release.

Read on if you'd like to build and run the installer yourself.

## Usage

### Notable Flags and Environment Variables

| flag                   | env variable                    | example              |
|------------------------|---------------------------------|----------------------|
| -address               | INSTALLER_ADDRESS               | 0.0.0.0:8080         |
| -platforms             | INSTALLER_PLATFORMS             | bare-metal, aws      |
| -cookie-signing-secret | INSTALLER_COOKIE_SIGNING_SECRET | secret               |
| -disable-secure-cookie | INSTALLER_DISABLE_SECURE_COOKIE | false                |
| -open-browser          | INSTALLER_OPEN_BROWSER          | false                |
| -log-level             | INSTALLER_LOG_LEVEL             | debug, warn, error   |
| -version               | INSTALLER_VERSION               | NA                   |
| -help                  | INSTALLER_HELP                  | NA                   |

## License

Get a [license](https://account.coreos.com) and follow the guides to create Tectonic clusters end to end.

## Build prerequisites

- [Go 1.8](https://golang.org/doc/install)
- [Nodejs >=8.x](https://nodejs.org/en/download/)
- [Yarn >=0.24.x](https://yarnpkg.com/lang/en/docs/install/)
- The tectonic-installer repo must be located at `$GOPATH/src/github.com/coreos/tectonic-installer`

### Build / Run

All commands mentioned here must be run from the same working directory as this README file, `./installer/` from the root of this repo.

Build the static binary.

```
make build
```

Run the binary for your platform (linux, darwin).

```
./bin/linux/installer -help
```

Visit [http://127.0.0.1:4444](http://127.0.0.1:4444).

## Managing Dependencies

### Frontend

Dependencies are managed with yarn and browserify. Unlike go
dependencies, yarn dependencies are *not* vendored directly, because
`yarn install` will build native extensions that could break builds on
other platforms/operating systems. To add a dependency, run:

```
cd $GOPATH/src/github.com/coreos/tectonic-installer/installer/frontend
yarn add $MY_PACKAGE # for a runtime dependency
```

If you are adding a build dependency, run the following commands instead:

```
cd $GOPATH/src/github.com/coreos/tectonic-installer/installer/frontend
yarn add --dev $MY_BUILD_PACKAGE # for a development dependency
```

Both sets of commands will update the `package.json` and
`yarn.lock` files in the repository - those changes should
then be committed.

### Go

Dependencies are managed with [glide](https://glide.sh/), but committed directly to the repository.

If you don't have glide, install the latest release from [https://glide.sh/](https://glide.sh/). We require version 0.12 at a minimum.

To add a new dependency:

- Edit the `glide.yaml` file to add your dependency.
- Ensure you add a `version` field for the sha or tag you want to pin to.
- Revendor the dependencies:

```
make vendor
```

If it worked correctly it should:
- Clone your new dep to the `/vendor` dir, and check out the ref you specified.
- Update `glide.lock` to include your new package, adds any transitive dependencies, and updates its hash.

For the sake of your fellow reviewers, commit vendored code changes as a separate commit from any other changes.

#### Regenerate or Repair Vendored Code

Should you need to regenerate or repair the vendored code en-mass from their source repositories, you can run:

```
make vendor
```
