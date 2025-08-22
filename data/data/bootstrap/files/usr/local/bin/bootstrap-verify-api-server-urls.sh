#!/usr/bin/env bash

# shellcheck disable=SC1091
. /usr/local/bin/bootstrap-service-record.sh

# This functions expects 2 arguments:
# 1. name of the URL
# 2. The value of the URL
function lookup_url() {
    unset IPS
    unset IP
    IPS=$(dig "${2}" +short)
    V6IPS=$(dig "${2}" AAAA +short)
    if [[ ! -z "${IPS}" ]] || [[ ! -z "${V6IPS}" ]]; then
        echo "Successfully resolved ${1} ${2}"
        return 0
    else
        echo "Unable to resolve ${1} ${2}"
        return 1
    fi
}

# This functions expects 2 arguments:
# 1. name of the URL
# 2. URL to validate
function validate_url() {
    if [[ $(curl --head -k --silent --fail --write-out "%{http_code}\\n" "${2}" -o /dev/null) == 200 ]]; then
        echo "Success while trying to reach ${1}'s https endpoint at ${2}"
        return 0
    else
        echo "Unable to reach ${1}'s https endpoint at ${2}"
        return 1
    fi
}
#
# This functions expects 2 arguments:
# 1. name of the URL
# 2. URL to validate
function validate_get_url() {
    if [[ $(curl -k --get --silent --fail --write-out "%{http_code}\\n" "${2}" -o /dev/null) == 200 ]]; then
        echo "Success while trying to reach ${1}'s https endpoint at ${2}"
        return 0
    else
        echo "Unable to reach ${1}'s https endpoint at ${2}"
        return 1
    fi
}

function resolve_url() {
    if [[ -z "${1}" ]] || [[ -z "${2}" ]]; then
        echo "Usage: resolve_url <API_URL or API_INT URL> <URL that needs to be verified>"
        return 1
    fi

    local URL_TYPE=${1}
    local SERVER_URL=${2}

    if [[ ${URL_TYPE} != API_URL ]] && [[ ${URL_TYPE} != API_INT_URL ]]; then
        echo "Usage: resolve_url <API_URL or API_INT URL> <URL that needs to be verified>"
        return 1
    fi

    echo "Checking if ${SERVER_URL} of type ${URL_TYPE} is resolvable"

    if [[ "${URL_TYPE}" = "API_URL" ]]; then
        local URL_STAGE_NAME="resolve-api-url"
    else 
        local URL_STAGE_NAME="resolve-api-int-url"
    fi

    echo "Starting stage ${URL_STAGE_NAME}"
    record_service_stage_start ${URL_STAGE_NAME}
    if lookup_url "$URL_TYPE" "$SERVER_URL"; then
        record_service_stage_success
        return 0
    else
        record_service_stage_failure
        # We do not want to stop bootkube service due to this failure.
        # So not returning failure at this point.
        return 1
    fi
}

function check_url() {
    if [[ -z "${1}" ]] || [[ -z "${2}" ]]; then
        echo "Usage: check_url <API_URL or API_INT URL> <URL that needs to be verified>"
        return 1
    fi

    local URL_TYPE=${1}
    local SERVER_URL=${2}

    if [[ ${URL_TYPE} != API_URL ]] && [[ ${URL_TYPE} != API_INT_URL ]]; then
        echo "Usage: check_url <API_URL or API_INT URL> <URL that needs to be verified>"
        return 1
    fi

    echo "Checking if ${SERVER_URL} of type ${URL_TYPE} reachable"

    if [[ "${URL_TYPE}" = "API_URL" ]]; then
        local URL_STAGE_NAME="check-api-url"
    else 
        local URL_STAGE_NAME="check-api-int-url"
    fi

    CURL_URL="https://${2}:6443/readyz"

    record_service_stage_start ${URL_STAGE_NAME}
    if validate_get_url "$URL_TYPE" "$CURL_URL"; then
        record_service_stage_success
        # Return the value from the validate_url- even on success
        return 0
    else
        echo "Unable to validate. ${CURL_URL} is currently unreachable."
        record_service_stage_failure
        # We do not want to stop bootkube service due to this failure.
        # So not returning failure at this point.
        return 1
    fi
}
