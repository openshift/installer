#!/bin/bash
while [ ! -f /opt/openshift/node-tuning-bootstrap.done ]
do
  sleep 0.2 # or less like 0.2
done
systemctl stop bootkube