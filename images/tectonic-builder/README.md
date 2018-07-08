# tectonic-builder

[![Container Repository on Quay](https://quay.io/repository/coreos/tectonic-builder/status "Container Repository on Quay")](https://quay.io/repository/coreos/tectonic-builder)

This container image contains the environment required to build and test the
[Tectonic Installer](../../installer) and aims at facilitating the implementation
of CI/CD pipelines. More particularly, this image is used in several Jenkins
jobs today for testing purposes.

Example:

```sh
docker build -t quay.io/coreos/tectonic-builder:v1.33 -f images/tectonic-builder/Dockerfile .
```
