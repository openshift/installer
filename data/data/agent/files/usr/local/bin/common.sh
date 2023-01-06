#!/bin/bash

source /usr/local/share/assisted-service/assisted-service.env 

wait_for_assisted_service() {
    echo "Waiting for assisted-service to be ready"
    until $(curl --output /dev/null --silent --fail ${SERVICE_BASE_URL}/api/assisted-install/v2/infra-envs); do
        printf '.'
        sleep 5
    done
}
