#!/bin/bash

source common.sh

HOST=$(get_host)
echo Using hostname ${HOST} 1>&2

if [[ ${HOST} == {{.NodeZeroIP}} ]] ;then
   mkdir -p /etc/assisted-service && cd /etc/assisted-service && touch node0
   echo "This host is identified and set as node zero to run OpenShift Assisted Installer Service." > /etc/assisted-service/node0
   echo "Created file /etc/assisted-service/node0"
fi
