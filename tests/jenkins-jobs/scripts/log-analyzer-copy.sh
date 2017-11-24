#!/bin/bash -xe

OUTPUT=pn:"${CHANGE_ID}"-bn:"${BUILD_NUMBER}"-au:"${CHANGE_AUTHOR}".log

curl -u  "${LOG_ANALYZER_USER}":"${LOG_ANALYZER_PASSWORD}" "${BUILD_URL}""consoleText" >> "${OUTPUT}"
echo "Started" >> "${OUTPUT}" #Line used to indicate that this is the end of the logfile (used by Logstash)
header="s3://"
aws s3 cp "${OUTPUT}" "${header}""${LOGSTASH_BUCKET}"
rm "${OUTPUT}"
