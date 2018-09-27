#!/bin/sh

case "$OSTYPE" in
  darwin*)  PLATFORM="darwin" ;;
  linux*)   PLATFORM="linux" ;;
  *)        PLATFORM="linux" ;;
esac

TERRAFORM_VERSION="0.11.8" &&
TERRAFORM_URL="https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_${PLATFORM}_amd64.zip" &&
cd "$(dirname "$0")/.." &&
mkdir -p bin &&
curl -L "${TERRAFORM_URL}" | funzip > bin/terraform &&
chmod +x bin/terraform
