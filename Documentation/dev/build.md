# Building Tectonic Installer

The Tectonic Installer leverages the [Bazel build system](https://bazel.build/) to build all artifacts, binaries, and documentation.

## Getting Started

Install Bazel > 0.11.x using the instructions [provided online](https://docs.bazel.build/versions/master/install.html).
*Note*: compiling Bazel from source requires Bazel.
*Note*: some Linux platforms may require installing _Static libraries for the GNU standard C++ library_ (on Fedora `dnf install libstdc++-static`)

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

To use the installer you can now do the following:

```sh
cd bazel-bin
tar -xvzf tectonic.tar.gz
cd tectonic
```

Then proceed using the installer as documented on [coreos.com](https://coreos.com/tectonic/docs/).

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
├── installer
├── modules
└── steps
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
├── installer
├── modules
└── steps
```

*Note*: the generated tarball will not include the version string in its own name since output names must be known ahead of time in Bazel. To include the version in the tarball name, copy or move the archive with the desired name in the destination.

## Cleaning

You can cleanup all generated files by running:
```sh
bazel clean
```

Additionally you can remove all toolchains (in addition to the generated files) with:
```sh
bazel clean --expunge
```
