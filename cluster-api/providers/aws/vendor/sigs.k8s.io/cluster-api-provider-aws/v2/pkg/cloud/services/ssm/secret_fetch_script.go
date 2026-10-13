/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ssm

//nolint:gosec
const secretFetchScript = `#cloud-boothook 
#!/bin/bash

# Copyright 2020 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# 	http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

umask 006

REGION="{{.Region}}"
if [ "{{.Endpoint}}" != "" ]; then
  ENDPOINT="--endpoint-url {{.Endpoint}}"
fi
SECRET_PREFIX="{{.SecretPrefix}}"
CHUNKS="{{.Chunks}}"
FILE="/etc/secret-userdata.txt"
FINAL_INDEX=$((CHUNKS - 1))

# Log an error and exit.
# Args:
#   $1 Message to log with the error
#   $2 The error code to return
log::error_exit() {
  local message="${1}"
  local code="${2}"

  log::error "${message}"
  log::error "aws.cluster.x-k8s.io encrypted cloud-init script $0 exiting with status ${code}"
  exit "${code}"
}

log::success_exit() {
  log::info "aws.cluster.x-k8s.io encrypted cloud-init script $0 finished"
  exit 0
}

# Log an error but keep going.
log::error() {
  local message="${1}"
  timestamp=$(date --iso-8601=seconds)
  echo "!!! [${timestamp}] ${1}" >&2
  shift
  for message; do
    echo "    ${message}" >&2
  done
}

# Print a status line.  Formatted to show up in a stream of output.
log::info() {
  timestamp=$(date --iso-8601=seconds)
  echo "+++ [${timestamp}] ${1}"
  shift
  for message; do
    echo "    ${message}"
  done
}

check_aws_command() {
  local command="${1}"
  local code="${2}"
  local out="${3}"
  local sanitised="${out//[$'\t\r\n']/}"
  case ${code} in
  "0")
    log::info "AWS CLI reported successful execution for ${command}"
    ;;
  "2")
    log::error "AWS CLI reported that it could not parse ${command}"
    log::error "${sanitised}"
    ;;
  "130")
    log::error "AWS CLI reported SIGINT signal during ${command}"
    log::error "${sanitised}"
    ;;
  "255")
    log::error "AWS CLI reported service error for ${command}"
    log::error "${sanitised}"
    ;;
  *)
    log::error "AWS CLI reported unknown error ${code} for ${command}"
    log::error "${sanitised}"
    ;;
  esac
}
delete_secret_value() {
  local id="${SECRET_PREFIX}/${1}"
  local out
  log::info "deleting secret from AWS SSM Parameter Store"
  set +o errexit
  set +o nounset
  set +o pipefail
  out=$(
    aws ssm ${ENDPOINT} --region ${REGION} delete-parameter --name "${id}" 2>&1
  )
  local delete_return=$?
  set -o errexit
  set -o nounset
  set -o pipefail
  check_aws_command "SSM::DeleteSecret" "${delete_return}" "${out}"
  if [ ${delete_return} -ne 0 ]; then
    log::error_exit "Could not delete secret value" 2
  fi
}

delete_secrets() {
  for i in $(seq 0 ${FINAL_INDEX}); do
    delete_secret_value "$i"
  done
}

get_secret_value() {
  local chunk=$1
  local id="${SECRET_PREFIX}/${chunk}"

  log::info "getting userdata from AWS SSM Parameter Store"
  log::info "getting secret value from AWS SSM Parameter Store"

  local data
  set +o errexit
  set +o nounset
  set +o pipefail
  data=$(
    set +e
    set +o pipefail
    aws ssm ${ENDPOINT} --region ${REGION} get-parameter --output text --query 'Parameter.Value' --with-decryption --name "${id}" 2>&1
  )
  local get_return=$?
  check_aws_command "SSM::GetSecretValue" "${get_return}" "${data}"
  set -o errexit
  set -o nounset
  set -o pipefail
  if [ ${get_return} -ne 0 ]; then
    log::error "could not get secret value, deleting secret"
    delete_secrets
    log::error_exit "could not get secret value, but secret was deleted" 1
  fi
  log::info "appending data to temporary file ${FILE}.b64.gz"
  echo "${data}" >>${FILE}.b64.gz
}

log::info "aws.cluster.x-k8s.io encrypted cloud-init script $0 started"
log::info "secret prefix: ${SECRET_PREFIX}"
log::info "secret count: ${CHUNKS}"

if test -f "${FILE}"; then
  log::info "encrypted userdata already written to disk"
  log::success_exit
fi

for i in $(seq 0 "${FINAL_INDEX}"); do
  get_secret_value "$i"
done

delete_secrets

log::info "decompressing userdata to ${FILE}"
cat "${FILE}.b64.gz" | base64 -d > "${FILE}.gz"
rm -f "${FILE}.b64.gz"
gunzip "${FILE}.gz"
GUNZIP_RETURN=$?
if [ ${GUNZIP_RETURN} -ne 0 ]; then
  log::error_exit "could not unzip data" 4
fi

log::info "restarting cloud-init"
systemctl restart cloud-init
log::success_exit
`
