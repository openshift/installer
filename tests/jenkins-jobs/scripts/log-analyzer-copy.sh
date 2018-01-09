#!/bin/bash


set -xeo pipefail

jenkins_filename="jenkins.log"
logfile_location=./templogfiles
action=$1
testFilePath=$2
additionalFields=""

mkdir -p ${logfile_location}
cd ${logfile_location}

if [ "$action" = "smoke-test-logs" ] || [ "$action" = "baremetal-smoke-test-logs" ]; then
  # Checks whether there are any log files
  if ls ../build/*/*.log 1> /dev/null 2>&1; then
    cp ../build/*/*.log .
  else
    exit 0
  fi
elif [ "$action" = "jenkins-logs" ]; then
  curl -u  "${LOG_ANALYZER_USER}":"${LOG_ANALYZER_PASSWORD}" "${BUILD_URL}""consoleText" >> "${jenkins_filename}"
fi

#Logstash uses the key=value filter and the ',' character as split field. This means that branches can neither have the '=' nor ',' in them.
SOURCE_BRANCH=${BRANCH_NAME//[=,]/_}
TARGET_BRANCH=${CHANGE_TARGET//[=,]/_}

#Check whether the BUILD_RESULT is not empty. If the variable is not empty (so success/failure/*), add the value to the additionalField variable.
if [ -n "${BUILD_RESULT}" ]; then
   additionalFields="${additionalFields}build_result=${BUILD_RESULT},"
fi

if [ -n "$testFilePath" ]; then
        platform=$(echo "${testFilePath}" | cut -d'/' -f2)
        specName=$(echo "${testFilePath}" | cut -d'/' -f3 | sed 's/_spec\.rb//')
        additionalFields="${additionalFields}platform=${platform},specName=${specName},"
fi

for file in *.log
  do
    mv "$file" "pull_request=${CHANGE_ID},build_number=${BUILD_NUMBER},author=${CHANGE_AUTHOR},source_branch=${SOURCE_BRANCH},target_branch=${TARGET_BRANCH},${additionalFields}${file}"
  done

aws s3 sync ../templogfiles "s3://""${LOGSTASH_BUCKET}"
