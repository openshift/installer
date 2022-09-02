#!/bin/bash

wait_for_assisted_service() {
    echo "Waiting for assisted-service to be ready"
    until $(curl --output /dev/null --silent --fail {{.ServiceBaseURL}}/api/assisted-install/v2/infra-envs); do
        printf '.'
        sleep 5
    done
}
