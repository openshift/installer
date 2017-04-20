#!/bin/bash -e

# check if cluster domain is set
if [ -z "${CLUSTER_DOMAIN}" ]; then
  echo "\$CLUSTER_DOMAIN must be set to check for DNS."
  exit 1
fi

# TODO: Replace with check on cluster/status
echo "Checking if records exist for ${CLUSTER_DOMAIN}"
for i in `seq 1 50`; do
  # hack to speed up NCACHE invalidation
  dig @8.8.8.8 -t NS +short ${CLUSTER_DOMAIN#*.} >/dev/null || true

  # get IPs for cluster domain
  ips=$(dig @8.8.8.8 -t A +short ${CLUSTER_DOMAIN})
  if [ "$(printf "${ips}" | wc -l)" != "0" ]; then
    echo "Found records!"
    printf "IPs:\n${ips}"
    echo
    break
  fi
  echo "no records found yet"
  sleep 10
done
