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

package secretsmanager

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
MAX_RETRIES=10
RETRY_DELAY=10 # in seconds

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
  local id="${SECRET_PREFIX}-${1}"
  local out
  log::info "deleting secret from AWS Secrets Manager"
  set +o errexit
  set +o nounset
  set +o pipefail
  out=$(
    aws secretsmanager ${ENDPOINT} --region ${REGION} delete-secret --force-delete-without-recovery --secret-id "${id}" 2>&1
  )
  local delete_return=$?
  check_aws_command "SecretsManager::DeleteSecret" "${delete_return}" "${out}"
  if [ ${delete_return} -ne 0 ]; then
    log::error "Could not delete secret value"
    return 1
  fi
}

retry_delete_secret_value() {
  local retries=0
  while [ ${retries} -lt ${MAX_RETRIES} ]; do
    delete_secret_value "$1"
    local return_code=$?
    if [ ${return_code} -eq 0 ]; then
      return 0
    else
      ((retries++))
      log::info "Retrying in ${RETRY_DELAY} seconds..."
      sleep ${RETRY_DELAY}
    fi
  done
  return 1
}

get_secret_value() {
  local chunk=$1
  local id="${SECRET_PREFIX}-${chunk}"

  log::info "getting userdata from AWS Secrets Manager"
  log::info "getting secret value from AWS Secrets Manager"

  local data
  set +o errexit
  set +o nounset
  set +o pipefail
  data=$(
    set +e
    set +o pipefail
    aws secretsmanager ${ENDPOINT} --region ${REGION} get-secret-value --output text --query 'SecretBinary' --secret-id "${id}" 2>&1
  )
  local get_return=$?
  check_aws_command "SecretsManager::GetSecretValue" "${get_return}" "${data}"
  if [ ${get_return} -ne 0 ]; then
    log::error "could not get secret value"
    return 1
  fi
  set -o errexit
  set -o nounset
  set -o pipefail
  log::info "appending data to temporary file ${FILE}.gz"
  echo "${data}" | base64 -d >>${FILE}.gz
}

retry_get_secret_value() {
  local retries=0
  while [ ${retries} -lt ${MAX_RETRIES} ]; do
    get_secret_value "$1"
    local return_code=$?
    if [ ${return_code} -eq 0 ]; then
      return 0
    else
      ((retries++))
      log::info "Retrying in ${RETRY_DELAY} seconds..."
      sleep ${RETRY_DELAY}
    fi
  done
  return 1
}

log::info "aws.cluster.x-k8s.io encrypted cloud-init script $0 started"
log::info "secret prefix: ${SECRET_PREFIX}"
log::info "secret count: ${CHUNKS}"

if test -f "${FILE}"; then
  log::info "encrypted userdata already written to disk"
  log::success_exit
fi

for i in $(seq 0 "${FINAL_INDEX}"); do
  retry_get_secret_value "$i"
  return_code=$?
  if [ ${return_code} -ne 0 ]; then
    log::error "Failed to get secret value after ${MAX_RETRIES} attempts"
  fi
done

for i in $(seq 0 ${FINAL_INDEX}); do
  retry_delete_secret_value "$i"
  return_code=$?
  if [ ${return_code} -ne 0 ]; then
    log::error "Failed to delete secret value after ${MAX_RETRIES} attempts"
    log::error_exit "couldn't delete the secret value, exiting" 1
  fi
done

log::info "decompressing userdata to ${FILE}"
gunzip "${FILE}.gz"
GUNZIP_RETURN=$?
if [ ${GUNZIP_RETURN} -ne 0 ]; then
  log::error_exit "could not unzip data" 4
fi

log::info "restarting cloud-init"
systemctl restart cloud-init
log::success_exit
`
