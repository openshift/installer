#!/bin/bash
set -x

yum update -y 

echo "Subscribing to Red Hat"
subscription-manager register --username=${RHSM_USERNAME} --password=${RHSM_PASSWORD} --name=packer-rhel7-$(date +%Y%m%d)-${RANDOM}
subscription-manager refresh
subscription-manager attach --auto

echo "Getting OC client https://docs.openshift.com/container-platform/4.3/installing/install_config/installing-restricted-networks-preparations.html#cli-installing-cli_installing-restricted-networks-preparations"
wget -O oc-client.tar.gz https://mirror.openshift.com/pub/openshift-v4/clients/ocp/latest/openshift-client-linux.tar.gz
tar -zxvf oc-client.tar.gz
mv kubectl /usr/bin/
mv oc /usr/bin/
oc version
kubectl version

wget https://github.com/stedolan/jq/releases/download/jq-1.6/jq-linux64
mv jq-linux64 /usr/bin/jq
chmod +x /usr/bin/jq

