#!/bin/bash

curl_assisted_service() {
    local endpoint=$1
    local method=${2:-GET}
    local authz=${3:-}
    local additional_options=("${@:4}")  # Capture all arguments starting from the fourth one
    local baseURL="${SERVICE_BASE_URL}api/assisted-install/v2"

    case "${method}" in
        "POST")
            curl -s -S -X POST "${additional_options[@]}" "${baseURL}${endpoint}" \
            -H "Authorization: ${authz}" \
            -H "accept: application/json" \
            -H "Content-Type: application/json" \
            ;;
        "GET")
            curl -s -S -X GET "${additional_options[@]}" "${baseURL}${endpoint}" \
            -H "Authorization: ${authz}" \
            -H "Accept: application/json"
            ;;
    esac
}   
