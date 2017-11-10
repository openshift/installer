#!/bin/bash -xe

OUTPUT="${JOB_BASE_NAME}"-"${BUILD_NUMBER}".log

curl -u  "${LOG_ANALYZER_USER}":"${LOG_ANALYZER_PASSWORD}" "${BUILD_URL}""consoleText" >> "${OUTPUT}"
echo "Started" >> "${OUTPUT}" #Line used to indicate that this is the end of the logfile (used by Logstash)
header="s3://"
aws s3 cp "${OUTPUT}" "${header}""${LOGSTASH_BUCKET}"
rm "${OUTPUT}"
