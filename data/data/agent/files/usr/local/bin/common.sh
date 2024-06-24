#!/bin/bash

curl_assisted_service() {
    local endpoint=$1
    local method=${2:-GET}
    local data=$3
    BASE_URL="${SERVICE_BASE_URL}api/assisted-install/v2"

    case "${method}" in
        "HEAD")
            curl -s -S -I "${BASE_URL}${endpoint}" \
            -H "Authorization: ${AGENT_AUTH_TOKEN}" \
            --output /dev/null --silent --fail
            ;;
        "POST")
            curl -s -S -X -w "%{http_code}\\n" -o /dev/null "${BASE_URL}${endpoint}" \
            -H "Authorization: ${AGENT_AUTH_TOKEN}" \
            -d "${data}"
            ;;
        "GET")
            curl -s -S "${BASE_URL}${endpoint}" \
            -H "Authorization: ${AGENT_AUTH_TOKEN}"
            ;;
    esac
}