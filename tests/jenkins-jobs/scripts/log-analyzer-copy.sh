#!/bin/bash
# Uploads smoke test logs or Jenkins console output to an AWS S3 bucket (to be consumed by LogStash).
# Usage:
#    log-analyzer-copy.sh [-d] smoke-test-logs testfile          # Upload all smoke test logs
#    log-analyzer-copy.sh [-d] jenkins-logs                      # Upload Jenkins console output
#
# Parameters:
#    -d                 Dry run (prepare the files but don't upload to the bucket)
# # Example:
#    log-analyzer-copy.sh smoke-test-logs tests/aws/some_test_spec.rb
#
# Required environment variables (apart from the standard ones provided by Jenkins)
# - LOGSTASH_BUCKET             # Target AWS S3 bucket name
# - LOG_ANALYZER_USER           # Jenkins User/Pass to retrieve the Jenkins console output
# - LOG_ANALYZER_PASSWORD

function flattenFields {
  # convert to "name=value,..."
  local k v sep=""
  for key in "${!fields[@]}" ; do
    k=$(sanitize "${key}")
    v=$(sanitize "${fields[${key}]}")
    printf "%s%s=%s" "${sep}" "${k}" "${v}"
    sep=","
  done
}

function sanitize {
  # Logstash uses the key=value filter on the filename, and ',' as the split char. This means that field names/values
  # may not contain '=', ',' or '/'.
  local p="$*"
  echo "${p//[=,\/]/_}"
}

function fatal {
  echo >&2 "$*"
  exit 1
}

set -xeo pipefail

if [[ -z ${LOGSTASH_BUCKET} ]]; then
  fatal "Required environment vars are not set."
fi

do_upload=1
while getopts ":d" arg; do
  case "${arg}" in
    d)
      do_upload=0
    ;;

    *)
      fatal "Invalid option"
    ;;
  esac
done
# Remove any script params that getopts has processed.
shift $((OPTIND-1))

action=$1
case "${action}" in
  "smoke-test-logs")
    testFilePath=$2
    if [[ -z ${testFilePath} ]]; then
      fatal "Missing testFilePath"
    fi
    ;;

  "jenkins-logs")
    if [[ -z ${LOG_ANALYZER_USER} || -z ${LOG_ANALYZER_PASSWORD} ]]; then
      fatal "Required environment vars are not set."
    fi
    ;;

  *)
    fatal "Unknown command."
    ;;
esac

jenkins_filename="jenkins.log"
logfile_location=$(realpath ./templogfiles)

mkdir -p "${logfile_location}"
cd "${logfile_location}"

declare -A fields
fields[job]="${JOB_NAME}"
fields[pull_request]="${CHANGE_ID}"
fields[build_number]="${BUILD_NUMBER}"
fields[author]="${CHANGE_AUTHOR}"
fields[source_branch]="${BRANCH_NAME}"
fields[target_branch]="${CHANGE_TARGET}"

if [[ -n ${BUILD_RESULT} ]]; then
  fields[build_result]="${BUILD_RESULT}"
fi

case "${action}" in
  "smoke-test-logs")
    # expected to be in the form "tests/<platform>/test_name_spec.rb"
    fields[platform]=$(echo -n "${testFilePath}" | cut -d'/' -f2)
    fields[specName]=$(echo -n "${testFilePath}" | cut -d'/' -f3 | sed -e 's/_spec\.rb//')

    # Checks whether there are any log files
    if ! compgen -G ../build/*/*.log >/dev/null ; then
      # Nothing to do
      exit 0
    fi
    # Note: only 1 sub-dir under build/ is expected. If the glob were to match >1 files with the same name in different
    # build/* dirs, then cp is smart enough to exit with an error.
    cp ../build/*/*.log .
    ;;

  "jenkins-logs")
    curl -u  "${LOG_ANALYZER_USER}:${LOG_ANALYZER_PASSWORD}" "${BUILD_URL}consoleText" >> "${jenkins_filename}"
    ;;
esac

for file in *.log ; do
  declare -i log_file_size
  log_file_size=$(stat --format=%s "${file}")
  if [[ ${log_file_size} -ge 1000000000 ]]; then
    echo "Log file is bigger than 1GB, dropping the file"
    rm "${file}"
  else
    fields[logfile]="${file}"
    newfile=$(flattenFields)
    if [[ ${#newfile} -gt 255 ]]; then
      echo "New file name is too long, skipping the file."
    else
      # Decorate the filename with useful metadata so Logstash can parse it.
      mv "${file}" "${newfile}"
    fi
  fi
done

if [[ ${do_upload} == 1 ]]; then
  aws s3 sync "${logfile_location}" "s3://${LOGSTASH_BUCKET}"
else
  echo "Dry run, skipping AWS upload"
fi
