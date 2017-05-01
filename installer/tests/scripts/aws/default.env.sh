#!/bin/bash -e
# This script configures variables for a new cluster. Some elements can be set using environment variables, their behavior is documented below.

#
# configuration variables
#

# ClUSTER_NAME is used for the name of the CloudFormation stack and if not otherwise specified the domain.
export CLUSTER_NAME=${CLUSTER_NAME:-"tectonic-test-$(date +%s)"}

# WORKER_COUNT is the number of workers to be created for the cluster.
export WORKER_COUNT=${WORKER_COUNT:-2}

# AWS_REGION is the region the CloudFormation stack should be deployed in
export AWS_REGION=${AWS_REGION:-"us-west-2"}

# AWS_HOSTEDZONE is the ID of the hosted zone to be used with the cluster. The variable CLUSTER_DOMAIN below should be within the hosted zone.
export AWS_HOSTEDZONE=${AWS_HOSTEDZONE:-"ZP8I2D1YRBH2D"}

# CLUSTER_DOMAIN is the domain name that should have A records created for it in the hosted zone provided above. If not specified, it is based on the cluster name.
export CLUSTER_DOMAIN=${CLUSTER_DOMAIN:-"${CLUSTER_NAME}-k8s.staging.core-os.net"}
export TECTONIC_DOMAIN=${TECTONIC_DOMAIN:-"${CLUSTER_NAME}.staging.core-os.net"}
export ETCD_DOMAIN=${ETCD_DOMAIN:-"${CLUSTER_NAME}-etc.staging.core-os.net"}

# UPDATER_ENABLED indicates whether to install the updater.
export UPDATER_ENABLED=${UPDATER_ENABLED:-true}

# UPDATER_COREUPDATE_SERVER is the address of the coreupdate server.
export UPDATER_SERVER=${UPDATER_SERVER:-"https://tectonic.update.core-os.net"}

# UPDATER_CHANNEL is the update channel to pull from.
export UPDATER_CHANNEL=${UPDATER_CHANNEL:-"tectonic-1.5"}

# UPDATER_APPID is the app ID of Tectonic in the coreupdate server.
export UPDATER_APPID=${UPDATER_APPID:-"6bc7b986-4654-4a0f-94b3-84ce6feb1db4"}
