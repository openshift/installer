#!/usr/bin/env bash
# This library provides a helper functions for recording when a service
# and its stages start and end.

# SERVICE_RECORDS_DIR is the directory under which service records will be stored.
SERVICE_RECORDS_DIR="${SERVICE_RECORDS_DIR:-/var/log/openshift/}"
# SYSTEMD_UNIT_NAME is the name of the systemd unit for the service
SYSTEMD_UNIT_NAME="$(ps -o unit= $$)"
# SERVICE_NAME is the name of the service
SERVICE_NAME="${SERVICE_NAME:-${SYSTEMD_UNIT_NAME%.service}}"
# STAGE_NAME is the name of the current stage in the service
STAGE_NAME=""


# add_service_record_entry adds a record entry to the service records file.
#   PHASE - phase being recorded; one of "start", "end", "stage start", or "stage end"
#   RESULT - result of the action
#   MESSAGE - message giving more detail about the result of the action
add_service_record_entry() {
  local FILENAME="${SERVICE_RECORDS_DIR}/${SERVICE_NAME}.json"
  mkdir --parents $(dirname "${FILENAME}")
  # Append the new entry to the existing array in the file.
  # If the file does not already exist, start with an empty array.
  # The new entry contains only the fields that have non-empty values, since the
  # stage, result, error line, and error message are all optional.
  ([ -f "${FILENAME}" ] && cat "${FILENAME}" || echo '[]') | \
      jq \
        --arg timestamp "$(date +"%Y-%m-%dT%H:%M:%SZ")" \
        --arg preCommand "${PRE_COMMAND-}" \
        --arg postCommand "${POST_COMMAND-}" \
        --arg stage "${STAGE_NAME-}" \
        --arg phase "${PHASE}" \
        --arg result "${RESULT-}" \
        --arg errorLine "${ERROR_LINE-}" \
        --arg errorMessage "${ERROR_MESSAGE-}" \
        '. += [
          {$timestamp,$preCommand,$postCommand,$stage,$phase,$result,$errorLine,$errorMessage} |
          reduce keys[] as $k (.; if .[$k] == "" then del(.[$k]) else . end)
        ]' \
      > "${FILENAME}.tmp" && \
    mv "${FILENAME}.tmp" "${FILENAME}"
}

# record_service_start() records the start of a service.
record_service_start() {
  local PHASE="start"

  add_service_record_entry
}

# record_service_end(result) records the end of a service.
#   ERROR_LINE - line where the error occurred, if there was an error
#   ERROR_MESSAGE - error message, if there was an error
record_service_end() {
  local PHASE="end"
  local RESULT=${1:?Must specify a result}

  add_service_record_entry
}

# record_service_stage_start(stage_name) records the start of a stage of a service.
record_service_stage_start() {
  if [ -n "${STAGE_NAME}" ]
  then
    echo "attempt to record the start of a stage without ending the previous one"
    exit 1
  fi

  local PHASE="stage start"
  STAGE_NAME=${1:?Must specify a stage name}

  add_service_record_entry
}

# record_service_stage_end(result) records the end  of a stage of a service.
#   ERROR_LINE - line where the error occurred, if there was an error
#   ERROR_MESSAGE - error message, if there was an error
record_service_stage_end() {
  if [ -z "${STAGE_NAME}" ]
  then
    echo "attempt to record the end of a stage without starting one"
    exit 1
  fi

  local PHASE="stage end"
  local RESULT=${1:?Must specify a result}

  add_service_record_entry

  STAGE_NAME=""
}

# record_service_stage_success records the successful end of a stage of a service.
record_service_stage_success() {
  record_service_stage_end "success"
}

record_service_stage_failure() {
  local ERROR_LINE
  local ERROR_MESSAGE
  get_error_info ERROR_LINE ERROR_MESSAGE
  record_service_stage_end "failure"
}

record_service_exit() {
  if [ "$1" -eq 0 ]
  then
    local RESULT="success"
  else
    local RESULT="failure"
    local ERROR_LINE
    local ERROR_MESSAGE
    get_error_info ERROR_LINE ERROR_MESSAGE
  fi

  if [ -n "${STAGE_NAME}" ]
  then
    record_service_stage_end "${RESULT}"
  fi

  record_service_end "${RESULT}"
}

get_error_info() {
  local -n error_line=$1
  local -n error_message=$2
  error_line="$(caller 1)"
  error_message="$(journalctl --unit="${SYSTEMD_UNIT_NAME}" --lines=3 --output=cat)"
}

record_service_start

trap 'record_service_exit $?' EXIT