#!/usr/bin/env bash

#
# For more information, see installer/data/data/bootstrap/baremetal/README.md
#

set -e

# Script to configure and run the Keepalived instance used to manage the DNS
# and API VIPs.

mkdir --parents /etc/keepalived

KEEPALIVED_IMAGE=quay.io/celebdor/keepalived:latest
if ! podman inspect "$KEEPALIVED_IMAGE" &>/dev/null; then
    echo "Pulling release image..."
    podman pull "$KEEPALIVED_IMAGE"
fi

API_DNS="$(sudo awk -F[/:] '/apiServerURL/ {print $5}' /opt/openshift/manifests/cluster-infrastructure-02-config.yml)"
CLUSTER_NAME="$(awk -F. '{print $2}' <<< "$API_DNS")"
API_VIP="$(dig +noall +answer "$API_DNS" | awk '{print $NF}')"
IFACE_CIDRS="$(ip addr show | grep -v "scope host" | grep -Po 'inet \K[\d.]+/[\d.]+' | xargs)"
SUBNET_CIDR="$(/usr/local/bin/get_vip_subnet_cidr "$API_VIP" "$IFACE_CIDRS")"
NET_MASK="$(echo "$SUBNET_CIDR" | cut -d "/" -f 2)"
INTERFACE="$(ip -o addr show to "$SUBNET_CIDR" | head -n 1 | awk '{print $2}')"
CLUSTER_DOMAIN="${API_DNS#*.}"
DNS_VIP="$(dig +noall +answer "ns1.${CLUSTER_DOMAIN}" | awk '{print $NF}')"

# Virtual Router IDs. They must be different and 8 bit in length
API_VRID=$(/usr/local/bin/fletcher8 "$CLUSTER_NAME-api")
DNS_VRID=$(/usr/local/bin/fletcher8 "$CLUSTER_NAME-dns")

export API_VIP
export CLUSTER_NAME
export INTERFACE
export DNS_VIP
export API_VRID
export DNS_VRID
export NET_MASK
envsubst < /etc/keepalived/keepalived.conf.tmpl | sudo tee /etc/keepalived/keepalived.conf

MATCHES="$(podman ps -a --format "{{.Names}}" | awk '/keepalived$/ {print $0}')"
if [[ ! -z "$MATCHES" ]]; then
    # Remove old pod, or you can get storage for container removed
    # errors on restart if the container exits unexpectedly
    podman rm -f keepalived
fi
# TODO(bnemec): Figure out how to run with less perms
podman create --rm \
        --name keepalived \
        --volume /etc/keepalived:/etc/keepalived:z \
        --network=host \
        --privileged \
        --cap-add=ALL \
        "${KEEPALIVED_IMAGE}" \
        /usr/sbin/keepalived -f /etc/keepalived/keepalived.conf \
            --dont-fork -D -l -P
