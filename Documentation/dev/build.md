# Building Tectonic Installer

The Tectonic Installer leverages the [Bazel build system](https://bazel.build/) to build all artifacts, binaries, and documentation.

## Getting Started

Install Bazel > 0.10.x using the instructions [provided online](https://docs.bazel.build/versions/master/install.html).
*Note*: compiling Bazel from source requires Bazel.
*Note*: cross-compilation of Golang binaries, which the Tectonic Installer leverages, is [broken in Bazel 0.9.0](https://github.com/bazelbuild/rules_go/issues/1161).

Clone the Tectonic Installer git repository locally:

```sh
git clone git@github.com:coreos/tectonic-installer.git && cd tectonic-installer
```

## Quickstart

To build Tectonic for development or testing, build the `tarball` target with Bazel:

```sh
bazel build tarball
```

This will produce an archive named `tectonic.tar.gz` in the `bazel-bin` directory, containing all the assets necessary to bring up a Tectonic cluster, namely the:

* Tectonic Installer binary;
* Terraform modules;
* Terraform binary;
* Terraform provider binaries; and
* examples

To build a versioned Tectonic release tarball, set a `TECTONIC_VERSION` environment variable and build the `tarball` target:

```sh
export TECTONIC_VERSION=1.2.3-beta
bazel build tarball --action_env=TECTONIC_VERSION
```

This will also create an archive named `tectonic.tar.gz` in the `bazel-bin` directory, however when the archive is extracted, the base directory name will include the specified version number.

For more details on building a Tectonic release or other Tectonic assets as well as workarounds to some known issues, read on.

## Building A Release Tarball

To build a release tarball for the Tectonic Installer, issue the following command from the `tectonic-installer` root directory:

```sh
bazel build tarball
```

*Note*: Bazel < 0.9.0 is known to [fail to build tarballs when using Python 3](https://github.com/bazelbuild/bazel/issues/3816); to avoid this issue, force Python 2 by using:

```sh
bazel build --force_python=py2 --python_path=/usr/bin/python2 tarball
```

This will create a tarball named `tectonic.tar.gz` in the `bazel-bin` directory with the following directory structure:

```
tectonic
├── config.tf
├── examples
├── modules
├── platforms
└── tectonic-installer
    ├── darwin
    │   ├── installer
    │   ├── terraform
    │   └── terraform-provider-matchbox
    └── linux
        ├── installer
        ├── terraform
        └── terraform-provider-matchbox
```

In order to build a release tarball with the version string in the directory name within the tarball, export a `TECTONIC_VERSION` environment variable and then build the tarball while passing the variable to the build:

```sh
export TECTONIC_VERSION=1.2.3-beta
bazel build tarball --action_env=TECTONIC_VERSION
```

This will create a tarball named `tectonic.tar.gz` in the `bazel-bin` directory with the following directory structure:

```
tectonic_1.2.3-beta
├── config.tf
├── examples
├── modules
├── platforms
└── tectonic-installer
```

*Note*: the generated tarball will not include the version string in its own name since output names must be known ahead of time in Bazel. To include the version in the tarball name, copy or move the archive with the desired name in the destination.

## Building the Installer Binary

For cases where the entire Tectonic tarball is not needed and only the Tectonic Installer GUI is required, the installer binary can be built by itself.
To build the Tectonic Installer binary, issue the following command from the `tectonic-installer` root directory:

```sh
bazel build backend
```

This will produce a binary located at `bazel-bin/installer/cmd/installer/linux_amd64_pure_stripped/installer` when built on a Linux machine or `bazel-bin/installer/cmd/installer/darwin_amd64_pure_stripped/installer` if built on a Mac.

To build a cross-compiled binary for another platform, e.g. a Darwin binary on a Linux machine, specify the target platform explicitly:

```sh
bazel build backend --experimental_platforms=@io_bazel_rules_go//go/toolchain:darwin_amd64
```

## Building the Smoke Test Binary

We also use Bazel to build the smoke test binary. To do so, run:

```sh
bazel build tests/smoke
```

This operation will produce a binary located at `bazel-bin/tests/smoke/linux_amd64_stripped/smoke`, if on a Linux machine, or `bazel-bin/tests/smoke/darwin_amd64_stripped/smoke`, if on a Mac.
Follow the [smoke test instructions][smoke-test] to test a Tectonic cluster using this newly compiled binary.

[smoke-test]: ../../tests/smoke/README.md
