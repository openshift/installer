#!/bin/bash

curl_assisted_service() {
    local endpoint=$1
    local method=${2:-GET}
    local additional_options=("${@:3}")  # Capture all arguments starting from the third one
    local baseURL="${SERVICE_BASE_URL}api/assisted-install/v2"

    headers=(
        -s -S
        -H "Authorization: ${USER_AUTH_TOKEN}"
        -H "accept: application/json"
    )

    [[ "$method" == "POST" || "$method" == "PATCH" ]] && headers+=(-H "Content-Type: application/json")

    curl "${headers[@]}" -X "${method}" "${additional_options[@]}" "${baseURL}${endpoint}"

}   
