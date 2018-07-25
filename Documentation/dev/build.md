# Building OpenShift Installer

The OpenShift Installer leverages the [Bazel build system](https://bazel.build/) to build all artifacts, binaries, and documentation.

## Getting Started

Install Bazel > 0.11.x using the instructions [provided online](https://docs.bazel.build/versions/master/install.html).
*Note*: compiling Bazel from source requires Bazel.
*Note*: some Linux platforms may require installing _Static libraries for the GNU standard C++ library_ (on Fedora `dnf install libstdc++-static`)

Clone the OpenShift Installer git repository locally:

```sh
git clone git@github.com:openshift/installer.git && cd installer
```

## Quickstart

To build OpenShift for development or testing, build the `tarball` target with Bazel:

```sh
bazel build tarball
```

This will produce an archive named `openshift-installer-dev.tar.gz` in the `bazel-bin` directory, containing all the assets necessary to bring up a OpenShift cluster, namely the:

* OpenShift Installer binary;
* Terraform modules;
* Terraform binary;
* Terraform provider binaries; and
* examples

To use the installer you can now do the following:

```sh
cd bazel-bin
tar -xvzf openshift-installer-dev.tar.gz
cd openshift-installer-dev
```

Then proceed using the installer as documented on [coreos.com](https://coreos.com/tectonic/docs/).

For more details on building a OpenShift release or other OpenShift assets as well as workarounds to some known issues, read on.

## Building A Release Tarball

To build a release tarball for the OpenShift Installer, issue the following command from the `installer` root directory:

```sh
bazel build tarball
```

*Note*: Bazel < 0.9.0 is known to [fail to build tarballs when using Python 3](https://github.com/bazelbuild/bazel/issues/3816); to avoid this issue, force Python 2 by using:

```sh
bazel build --force_python=py2 --python_path=/usr/bin/python2 tarball
```

This will create a tarball named `openshift-installer-dev.tar.gz` in the `bazel-bin` directory with the following directory structure:

```
openshift-installer-dev
├── config.tf
├── examples
├── installer
├── modules
└── steps
```

In order to build a release tarball with the version string in the directory name within the tarball, export a `OPENSHIFT_VERSION` environment variable and then build the tarball while passing the variable to the build:

```sh
export OPENSHIFT_VERSION=1.2.3-beta
bazel build tarball --action_env=OPENSHIFT_VERSION
```

This will create a tarball named `openshift-install.tar.gz` in the `bazel-bin` directory with the following directory structure:

```
openshift-install_1.2.3-beta
├── config.tf
├── examples
├── installer
├── modules
└── steps
```

*Note*: the generated tarball will not include the version string in its own name since output names must be known ahead of time in Bazel. To include the version in the tarball name, copy or move the archive with the desired name in the destination.

## Building the Smoke Test Binary

We also use Bazel to build the smoke test binary. To do so, run:

```sh
bazel build tests/smoke
```

This operation will produce a binary located at `bazel-bin/tests/smoke/linux_amd64_stripped/smoke`, if on a Linux machine, or `bazel-bin/tests/smoke/darwin_amd64_stripped/smoke`, if on a Mac.
Follow the [smoke test instructions][smoke-test] to test a OpenShift cluster using this newly compiled binary.


## Cleaning

You can cleanup all generated files by running:
```sh
bazel clean
```

Additionally you can remove all toolchains (in addition to the generated files) with:
```sh
bazel clean --expunge
```

[smoke-test]: ../../tests/smoke/README.md
