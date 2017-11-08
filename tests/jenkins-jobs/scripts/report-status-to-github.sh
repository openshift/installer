#!/bin/bash -ex

STATUS="$1"
CONTEXT=${2/spec/smoke-tests}
COMMIT=$(git rev-parse origin/"${BRANCH_NAME}")

curl -f \
     -H 'Content-Type: application/json' \
     -u "$GITHUB_CREDENTIALS" \
     "https://api.github.com/repos/coreos/tectonic-installer/statuses/${COMMIT}" \
     -d "{\"state\": \"${STATUS}\", \"target_url\": \"${BUILD_URL}\", \"description\": \"\", \"context\": \"${CONTEXT}\"}"
