#!/bin/bash
#
# Calls assisted-service REST-API to create cluster using the params generated from reading the ZTP manifests.
# The cluster's id is found in the create cluster response.
# The cluster-id is then used in the REST-API call to create the infra-env.
#

source common.sh

wait_for_assisted_service

CLUSTER_CREATE_OUT=$(curl -X POST -H "Content-Type: application/json" -d '{{.ClusterCreateParamsJSON}}' {{.ServiceBaseURL}}/api/assisted-install/v2/clusters)
echo "cluster create response: $CLUSTER_CREATE_OUT"
# pick cluster_id out from cluster create response
CLUSTER_ID=$(echo $CLUSTER_CREATE_OUT | jq .id)
# trim quotes from CLUSTER_ID
CLUSTER_ID=${CLUSTER_ID//\"}

INFRA_ENV_PARAMS_JSON='{{.InfraEnvCreateParamsJSON}}'
# replace "replace-cluster-id" with $CLUSTER_ID
INFRA_ENV_PARAMS_JSON=${INFRA_ENV_PARAMS_JSON//"replace-cluster-id"/$CLUSTER_ID}

curl -X POST -H "Content-Type: application/json" -d $INFRA_ENV_PARAMS_JSON {{.ServiceBaseURL}}/api/assisted-install/v2/infra-envs
