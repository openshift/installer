# Tectonic Installer

The Tectonic Installer is an app for creating Tectonic clusters.

## Usage

Run the binary for your platform (linux, darwin).

    ./bin/linux/installer -help

Run the Docker image.

    sudo docker run -p 4444:4444 quay.io/coreos/tectonic-installer -address=0.0.0.0:4444 -platforms=aws

Visit [http://127.0.0.1](http://127.0.0.1).

### Flags and Environment Variables

| flag            | env variable           | example              |
|-----------------|------------------------|----------------------|
| -address        | INSTALLER_ADDRESS      | 0.0.0.0:8080         |
| -platforms      | INSTALLER_PLATFORMS    | bare-metal, aws      |
| -cookie-signing-secret | INSTALLER_COOKIE_SIGNING_SECRET | secret |
| -disable-secure-cookie | INSTALLER_DISABLE_SECURE_COOKIE | false |
| -open-browser   | INSTALLER_OPEN_BROWSER | false |
| -log-level      | INSTALLER_LOG_LEVEL    | debug, warn, error   |
| -version        | INSTALLER_VERSION      | NA                   |
| -help           | INSTALLER_HELP         | NA                   |

## License

Get a [license](../docs-internal/get-license.md) and follow the guides to create Tectonic clusters end to end.

## Development

Follow the [developer](../docs-internal) guides to create Tectonic clusters end to end.

* [Bare-metal Usage](../docs-internal/usage-baremetal.md)
* [AWS Usage](../docs-internal/usage-aws.md)

Bare-metal can be developed locally with QEMU/KVM, rkt, and a running matchbox service. AWS can be developed with an AWS account and an internet connection.

### Static Binary

Build the static binary.

    make build

### Container Image

    make docker-image

### Run

Run the binary for your platform (linux, darwin).

    ./bin/linux/installer -help

Run the Docker image.

    sudo docker run -p 4444:4444 coreos/tectonic-installer -address=0.0.0.0:4444

### Devleopment Flags

Mock calls to `matchbox` or AWS to avoid creating real clusters with `-no-configure`.

    ./bin/linux/installer -no-configure

For non-UI development, use `curl` to POST example data to the installer.

    $ ./bin/linux/installer -open-browser=false
    $ curl -H "Content-Type: application/json" -X POST -d @examples/metal.json http://127.0.0.1:4444/cluster/create

    $ ./bin/linux/installer -open-browser=false
    $ curl -H "Content-Type: application/json" -X POST -d @examples/aws.json http://127.0.0.1:4444/cluster/create -o ~/Downloads/assets.zip

Use the `-asset-dir` flag to serve assets directly from a local directory of your choice. The front end asset build system has a watch on change mode that will build new assets as changes are made to source files, and deposit the build results in `./assets/frontend`. Using both of these capabilities together can make UI development a bit more convenient:

    ./bin/linux/installer -asset-dir ./assets/frontend &
    pushd frontend/ && yarn run dev

## Managing Dependencies

### Frontend

Dependencies are managed with yarn and browserify. Unlike go
dependencies, yarn dependencies are *not* vendored directly, because
`yarn install` will build native extensions that could break builds on
other platforms/operating systems. To add a dependency, run:

    cd $GOPATH/src/github.com/coreos/tectonic-installer/installer/frontend
    yarn add $MY_PACKAGE # for a runtime dependency

If you are adding a build dependency, run the following commands instead:

    cd $GOPATH/src/github.com/coreos/tectonic-installer/installer/frontend
    yarn add --dev $MY_BUILD_PACKAGE # for a development dependency

Both sets of commands will update the `package.json` and
`yarn.lock` files in the repository - those changes should
then be committed.

### Go

Dependencies are managed with [glide](https://glide.sh/), but committed directly to the repository.

If you don't have glide, install the latest release from [https://glide.sh/](https://glide.sh/). We require version 0.12 at a minimum.

To add a new dependency:

- Edit the `glide.yaml` file to add your dependency.
- Ensure you add a `version` field for the sha or tag you want to pin to.
- Run glide to update the vendor source directory (with --strip-vendor)

To run glide, use the following commands.

    cd $GOPATH/src/github.com/coreos/tectonic-installer/installer
    glide update --strip-vendor --skip-test

If it worked correctly it should:
- Clone your new dep to the `/vendor` dir, and check out the ref you specified.
- Update `glide.lock` to include your new package, adds any transitive dependencies, and updates its hash.

For the sake of your fellow reviewers, commit vendored code changes as a separate commit from any other changes.

#### Regenerate or Repair Vendored Code

Should you need to regenerate or repair the vendored code en-mass from their source repositories, you can run:

    glide install --strip-vendor --skip-test

