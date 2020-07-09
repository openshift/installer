#!/bin/bash

# Update iptables rules based on google cloud load balancer VIPS
#
# This is needed because the GCP L3 load balancer doesn't actually do DNAT;
# the destination IP address is still the VIP. Normally, there is an agent that
# adds the vip to the local routing table, tricking the kernel in to thinking
# it's a local IP and allowing processes doing an accept(0.0.0.0) to receive
# the packets. Clever.
#
# We don't do that. Instead, we DNAT with conntrack. This is so we don't break
# existing connections when the vip is removed. This is useful for draining
# connections - take ourselves out of the vip, but service existing conns.
#
# Additionally, clients can write a file to /run/gcp-routes/$IP.down to force
# a VIP as down. This is useful for graceful shutdown / upgrade.
#
# ~cdc~

set -e

# the list of load balancer IPs that are assigned to this node
declare -A vips

curler() {
  curl --silent -L -H "Metadata-Flavor: Google" "http://metadata.google.internal/computeMetadata/v1/instance/${1}"
}

CHAIN_NAME="gcp-vips"
RUN_DIR="/run/gcp-routes"

# Create a chan if it doesn't exist
ensure_chain() {
    local table="${1}"
    local chain="${2}"

    if ! iptables -w -t "${table}" -S "${chain}" &> /dev/null ; then
        iptables -w -t "${table}" -N "${chain}";
    fi;
}

ensure_rule() {
    local table="${1}"
    local chain="${2}"
    shift 2

    if ! iptables -w -t "${table}" -C "${chain}" "$@" &> /dev/null; then
        iptables -w -t "${table}" -A "${chain}" "$@"
    fi
}

# set the chain, ensure entry rules, ensure ESTABLISHED rule
initialize() {
    ensure_chain nat "${CHAIN_NAME}"
    ensure_chain nat "${CHAIN_NAME}-local"
    ensure_rule nat PREROUTING -m comment --comment 'gcp LB vip DNAT' -j ${CHAIN_NAME}
    ensure_rule nat OUTPUT -m comment --comment 'gcp LB vip DNAT for local clients' -j ${CHAIN_NAME}-local

    # Need this so that existing flows (with an entry in conntrack) continue to be
    # balanced, even if the DNAT entry is removed
    ensure_rule filter INPUT -m comment --comment 'gcp LB vip existing' -m addrtype ! --dst-type LOCAL -m state --state ESTABLISHED,RELATED -j ACCEPT

    mkdir -p "${RUN_DIR}"
}

remove_stale() {
    ## find extra iptables rules
    for ipt_vip in $(iptables -w -t nat -S "${CHAIN_NAME}" | awk '$4{print $4}' | awk -F/ '{print $1}'); do
        if [[ -z "${vips[${ipt_vip}]}" ]]; then
            echo removing stale vip "${ipt_vip}" for external clients
            iptables -w -t nat -D "${CHAIN_NAME}" --dst "${ipt_vip}" -j REDIRECT
        fi
    done
    for ipt_vip in $(iptables -w -t nat -S "${CHAIN_NAME}-local" | awk '$4{print $4}' | awk -F/ '{print $1}'); do
        if [[ -z "${vips[${ipt_vip}]}" ]] || [[ "${vips[${ipt_vip}]}" = down ]]; then
            echo removing stale vip "${ipt_vip}" for local clients
            iptables -w -t nat -D "${CHAIN_NAME}-local" --dst "${ipt_vip}" -j REDIRECT
        fi
    done
}

add_rules() {
    for vip in "${!vips[@]}"; do
        echo "ensuring rule for ${vip} for external clients"
        ensure_rule nat "${CHAIN_NAME}" --dst "${vip}" -j REDIRECT

        if [[ "${vips[${vip}]}" != down ]]; then
            echo "ensuring rule for ${vip} for internal clients"
            ensure_rule nat "${CHAIN_NAME}-local" --dst "${vip}" -j REDIRECT
        fi
    done
}

clear_rules() {
    iptables -t nat -F "${CHAIN_NAME}" || true
    iptables -t nat -F "${CHAIN_NAME}-local" || true
}

# out paramater: vips
list_lb_ips() {
    for k in "${!vips[@]}"; do
        unset vips["${k}"]
    done

    local net_path="network-interfaces/"
    for vif in $(curler ${net_path}); do
        local hw_addr; hw_addr=$(curler "${net_path}${vif}mac")
        local fwip_path; fwip_path="${net_path}${vif}forwarded-ips/"
        for level in $(curler "${fwip_path}"); do
            for fwip in $(curler "${fwip_path}${level}"); do
                if [[ -e "${RUN_DIR}/${fwip}.down" ]]; then
                    echo "${fwip} is manually marked as down, skipping for internal clients..."
                    vips[${fwip}]="down"
                else
                    echo "Processing route for NIC ${vif}${hw_addr} for ${fwip}"
                    vips[${fwip}]="${fwip}"
                fi
            done
        done
    done
}

sleep_or_watch() {
    if hash inotifywait &> /dev/null; then
        inotifywait -t 30 -r "${RUN_DIR}" &> /dev/null || true
    else
        # no inotify, need to manually poll
        for i in {0..5}; do
            for vip in "${!vips[@]}"; do
                if [[ "${vips[${vip}]}" != down ]] && [[ -e "${RUN_DIR}/${vip}.down" ]]; then
                    echo "new downfile detected"
                    break 2
                elif [[ "${vips[${vip}]}" = down ]] && ! [[ -e "${RUN_DIR}/${vip}.down" ]]; then
                    echo "downfile disappeared"
                    break 2
                fi
            done
            sleep 5
        done
    fi
}

case "$1" in
  start)
    initialize
    while :; do
      list_lb_ips
      remove_stale
      add_rules
      echo "done applying vip rules"
      sleep_or_watch
    done
    ;;
  cleanup)
    clear_rules
    ;;
  *)
    echo $"Usage: $0 {start|cleanup}"
    exit 1
esac
