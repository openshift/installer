#!/bin/bash -ex

STATUS="$1"
CONTEXT=${2/spec/smoke-tests}
CONTEXT=${CONTEXT/_spec.rb/}
COMMIT="$3"

curl -f \
     -H 'Content-Type: application/json' \
     -u "$GITHUB_CREDENTIALS" \
     "https://api.github.com/repos/coreos/tectonic-installer/statuses/${COMMIT}" \
     -d "{\"state\": \"${STATUS}\", \"target_url\": \"${BUILD_URL}\", \"description\": \"\", \"context\": \"${CONTEXT}\"}"
