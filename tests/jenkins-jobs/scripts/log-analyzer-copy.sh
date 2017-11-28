#!/bin/bash -xe

jenkins_filename="jenkins.log"
logfile_location=./templogfiles
action=$1

mkdir -p ${logfile_location}
cd ${logfile_location}

if [ "$action" = "smoke-test-logs" ]; then
  # Checks whether there are any log files
  if ls ../build/*/*.log 1> /dev/null 2>&1; then
    cp ../build/*/*.log .
  else
    exit 0
  fi
elif [ "$action" = "jenkins-logs" ]; then
  curl -u  "${LOG_ANALYZER_USER}":"${LOG_ANALYZER_PASSWORD}" "${BUILD_URL}""consoleText" >> "${jenkins_filename}"
fi

for file in *.log
  do
    mv "$file" "pn:${CHANGE_ID}-bn:${BUILD_NUMBER}-au:${CHANGE_AUTHOR}-fn:${file}"
  done

aws s3 sync ../templogfiles "s3://""${LOGSTASH_BUCKET}"
