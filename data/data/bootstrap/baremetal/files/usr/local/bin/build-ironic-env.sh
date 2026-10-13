#!/bin/bash

set -euo pipefail

# shellcheck disable=SC1091
source /usr/local/bin/release-image.sh

build_ironic_env() {
    # Retrieve the provisioning interface name based on mac address (case
    # insensitive).
    printf 'PROVISIONING_INTERFACE="%s"\n' "$(ip -j -o link | jq -r ".[] | select(.address != null) | select(.address | match (\"${PROVISIONING_MAC}\"; \"i\")).ifname")"

    # If the IP contains a colon, then it's an IPv6 address, and the HTTP
    # host needs surrounding with brackets
    if [[ "$IRONIC_IP" =~ : ]]; then
        printf 'IRONIC_BASE_URL="https://[%s]"\n' "${IRONIC_IP}"
    else
        printf 'IRONIC_BASE_URL="https://%s"\n' "${IRONIC_IP}"
    fi

    printf 'IRONIC_IMAGE="%s"\n' "$(image_for ironic)"
    printf 'IRONIC_AGENT_IMAGE="%s"\n' "$(image_for ironic-agent)"
    printf 'CUSTOMIZATION_IMAGE="%s"\n' "$(image_for machine-image-customization-controller)"
    printf 'MACHINE_OS_IMAGES_IMAGE="%s"\n' "$(image_for machine-os-images)"

    # set password for ironic basic auth
    # The ironic container contains httpd (and thus httpd-tools), so rely on it
    # to supply the htpasswd command
    printf 'IRONIC_HTPASSWD="%s"\n' "$(podman run -i --rm "$(image_for ironic)" htpasswd -niB "${IRONIC_USERNAME}" </opt/metal3/auth/ironic/password)"

    if [ "${EXTERNAL_IP_FAMILY:-}" = "ip4" ]; then
        dhcp_opt="dhcp"
    elif [ "${EXTERNAL_IP_FAMILY:-}" = "ip6" ]; then
        dhcp_opt="dhcp6"
    else
        dhcp_opt="dhcp,dhcp6"
    fi
    printf 'EXTERNAL_IP_OPTIONS="ip=%s"\n' "${dhcp_opt}"

    if [ "${PROVISIONING_NETWORK_TYPE}" != "Disabled" ]; then
        if [ "${PROV_IP_FAMILY:-}" = "ip4" ]; then
            dhcp_opt="dhcp"
        elif [ "${PROV_IP_FAMILY:-}" = "ip6" ]; then
            dhcp_opt="dhcp6"
        fi
    fi
    printf 'PROVISIONING_IP_OPTIONS="ip=%s"\n' "${dhcp_opt}"
    # Allow some extra time for the provisioning interface to come up
    # https://issues.redhat.com/browse/OCPBUGS-5151
    printf 'IRONIC_KERNEL_PARAMS="rd.net.timeout.carrier=30 ip=%s"\n' "${dhcp_opt}"
}

build_ironic_env | tee -a /etc/ironic.env
