#!/bin/sh

# download-terraform.sh downloads the terraform binary and
# the terraform providers that are sourced from remote registries.

TERRAFORM_VERSION=1.0.11

TARGET_OS_ARCHES=${1}
if [ -z "${TARGET_OS_ARCHES}" ]
then
  TARGET_OS_ARCHES="linux_amd64 darwin_amd64 darwin_arm64 freebsd_amd64"
fi

for os_arch in ${TARGET_OS_ARCHES}
do
  # Download the terraform binary.
  rm -rf "./terraform/terraform/${os_arch}"
  mkdir -p "./terraform/terraform/${os_arch}"
  terraform_dir=$(mktemp -d -t terraform-XXXXXX)
  curl "https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_${os_arch}.zip" -o "${terraform_dir}/terraform.zip" &&
    unzip "${terraform_dir}/terraform.zip" -d "./terraform/terraform/${os_arch}/"
  rm -rf "${terraform_dir}"

  # Download providers specified in terraform/versions.tf.
  podman run --rm \
    --volume "${PWD}:/go/src/github.com/openshift/installer:z" \
    --workdir /go/src/github.com/openshift/installer \
    docker.io/hashicorp/terraform:${TERRAFORM_VERSION} \
    -chdir=terraform providers mirror -platform="${os_arch}" ./mirror
done
