#!/bin/sh

have() {
	command -v "${@}" >/dev/null 2>/dev/null
}

OS=linux
ARCH=amd64
FUNZIP="${FUNZIP:-funzip}"
if ! have "${FUNZIP}"
then
	if have gunzip
	then
		FUNZIP=gunzip
	else
		command -V "${FUNZIP}"
		exit 1
	fi
fi &&
if have go
then
	OS="$(go env GOOS)" &&
	ARCH="$(go env GOARCH)"
fi &&
TERRAFORM_VERSION="0.11.8" &&
TERRAFORM_URL="https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_${OS}_${ARCH}.zip" &&
echo "pulling ${TERRAFORM_URL}" >&2 &&
cd "$(dirname "$0")/.." &&
mkdir -p bin &&
curl -L "${TERRAFORM_URL}" | "${FUNZIP}" >bin/terraform &&
chmod +x bin/terraform
