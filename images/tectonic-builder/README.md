# tectonic-builder

[![Container Repository on Quay](https://quay.io/repository/coreos/tectonic-builder/status "Container Repository on Quay")](https://quay.io/repository/coreos/tectonic-builder)

This container image contains the environment required to build and test the
[Tectonic Installer](../installer) and aims at facilitating the implementation
of CI/CD pipelines. More particularly, this image is used in several Jenkins
jobs today for testing purposes.

## Upstream and CoreOS Terraform

The Tectonic CI is currently using a custom Terraform version for the default
pipeline (See https://github.com/coreos/tectonic-installer/pull/1247). As end
users of Tectonic installer use upstream Terraform we need to test with upstream
Terraform as well. We are building and publishing the Tectonic builder image
both with CoreOS Terraform as well as upstream Terraform. This is done via
docker `--build-arg` `TERRAFORM_URL`. CoreOS Terraform is used for normal PR and
branch tests, upstream Terraform is used once per day on master.

Example:
- CoreOS Terraform (default):
`docker build -t quay.io/coreos/tectonic-builder:v1.33 -f images/tectonic-builder/Dockerfile .`
- Upstream Terraform:
`docker build -t quay.io/coreos/tectonic-builder:v1.33-upstream-terraform --build-arg TERRAFORM_URL=https://releases.hashicorp.com/terraform/0.9.11/terraform_0.9.11_linux_amd64.zip -f images/tectonic-builder/Dockerfile .`
