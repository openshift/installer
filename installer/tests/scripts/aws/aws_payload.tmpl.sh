#!/bin/bash
cat <<EOF
{
  "clusterKind": "tectonic-aws",
  "cluster": {
    "accessKeyID": "${AWS_ACCESS_KEY_ID}",
    "secretAccessKey": "${AWS_SECRET_ACCESS_KEY}",
    "cloudForm": {
      "channel": "stable",
      "elbScheme": "internet-facing",
      "controllerDomain": "${CLUSTER_DOMAIN}",
      "tectonicDomain": "${TECTONIC_DOMAIN}",
      "clusterName": "${CLUSTER_NAME}",
      "region": "${AWS_REGION}",
      "hostedZoneID": "${AWS_HOSTEDZONE}",
      "kmsKeyARN": "arn:aws:kms:us-west-2:477645798544:key/61f91d7c-db07-4b1f-a7e3-201d68d8db4b",
      "keyName": "tectonic-jenkins",
      "workerCount": ${WORKER_COUNT}
    },
    "tectonic": {
      "ingressKind": "NodePort",
      "identityAdminUser": "admin@example.com",
      "identityAdminPassword": "dGVjdG9uaWNUZXN0UGFzczExMDQyMDE2",
      "license": "${TECTONIC_LICENSE}",
      "dockercfg": "${TECTONIC_PULL_SECRET}",
      "updater": {
        "enabled": ${UPDATER_ENABLED},
        "server": "${UPDATER_SERVER}",
        "channel": "${UPDATER_CHANNEL}",
        "appID": "${UPDATER_APPID}"
      }
    }
  }
}
EOF
