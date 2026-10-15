#!/usr/bin/env bash

# shellcheck disable=SC1091  # using path on bootstrap machine
. /usr/local/bin/bootstrap-service-record.sh

KUBECONFIG="${1}"

# This script still approves every pending CSR. Before doing so it logs the
# values that are contained in the CSR.

# subject_field prints the values a PEM encoded certificate request carries for
# one subject field, comma joined.
#   1 - field name, e.g. CN or O   2 - PEM encoded certificate request
subject_field() {
    printf '%s' "${2}" |
        openssl req -noout -subject -nameopt sep_multiline,utf8 2>/dev/null |
        sed -n "s/^ *${1}=//p" |
        paste -sd, -
}

# san_field prints the subject alternative names of one type carried by a PEM
# encoded certificate request, comma joined.
#   1 - SAN type as openssl prints it: DNS, "IP Address", email or URI
#   2 - PEM encoded certificate request
san_field() {
    printf '%s' "${2}" |
        openssl req -noout -text 2>/dev/null |
        sed -n '/X509v3 Subject Alternative Name/{n;p}' |
        tr ',' '\n' |
        sed -n "s/^ *${1}://p" |
        paste -sd, -
}

# log_value prints one labelled value, or <none> when it is empty.
log_value() {
    printf '  %-32s %s\n' "${1}" "${2:-<none>}"
}

# cluster_inputs renders the cluster side inputs. It is printed only when it
# changes, so the log shows when each of these starts being populated.
cluster_inputs() {
    echo "--- cluster inputs ---"

    local machines
    machines="$(oc --kubeconfig="$KUBECONFIG" get machines.machine.openshift.io --all-namespaces -o json 2>/dev/null | jq -r '
        .items[]?
        | "    \(.metadata.name) created=\(.metadata.creationTimestamp) nodeRef=\(.status.nodeRef.name // "<none>") addresses=\([.status.addresses[]? | "\(.type):\(.address)"] | join(",") | if . == "" then "<none>" else . end)"' 2>/dev/null)"
    if [ -z "${machines}" ]; then
        log_value "machines.machine.openshift.io:" "<none, or API not served>"
    else
        echo "  machines.machine.openshift.io:"
        printf '%s\n' "${machines}"
    fi

    local nodes
    nodes="$(oc --kubeconfig="$KUBECONFIG" get nodes -o json 2>/dev/null | jq -r '
        .items[]?
        | "    \(.metadata.name) addresses=\([.status.addresses[]? | "\(.type):\(.address)"] | join(","))"' 2>/dev/null)"
    if [ -z "${nodes}" ]; then
        log_value "nodes:" "<none>"
    else
        echo "  nodes:"
        printf '%s\n' "${nodes}"
    fi

    log_value "network .status.networkType:" \
        "$(oc --kubeconfig="$KUBECONFIG" get network.config.openshift.io cluster -o jsonpath='{.status.networkType}' 2>/dev/null)"
    log_value "egressips.k8s.ovn.org:" \
        "$(oc --kubeconfig="$KUBECONFIG" get egressips.k8s.ovn.org -o jsonpath='{range .items[*]}{.metadata.name}{" "}{end}' 2>/dev/null)"
    log_value "openshift-config-managed/csr-controller-ca:" \
        "$(oc --kubeconfig="$KUBECONFIG" get configmap csr-controller-ca -n openshift-config-managed -o jsonpath='{.metadata.name}' 2>/dev/null)"
}

# log_csr renders the values carried by one CSR, plus the Node and Machine that
# cluster-machine-approver would look up for it.
#   1 - name  2 - signer  3 - username  4 - groups  5 - usages
#   6 - creationTimestamp  7 - conditions  8 - PEM encoded request
log_csr() {
    local pem="${8}" common_name node_name

    echo "=== CSR ${1} ==="
    log_value ".spec.signerName:" "${2}"
    log_value ".spec.username:" "${3}"
    log_value ".spec.groups:" "${4}"
    log_value ".spec.usages:" "${5}"
    log_value ".metadata.creationTimestamp:" "${6}"
    log_value ".status.conditions:" "${7}"

    common_name="$(subject_field CN "${pem}")"
    log_value "x509 Subject.CommonName:" "${common_name}"
    log_value "x509 Subject.Organization:" "$(subject_field O "${pem}")"
    log_value "x509 DNSNames:" "$(san_field DNS "${pem}")"
    log_value "x509 IPAddresses:" "$(san_field 'IP Address' "${pem}")"
    log_value "x509 EmailAddresses:" "$(san_field email "${pem}")"
    log_value "x509 URIs:" "$(san_field URI "${pem}")"

    # The node name is taken from the common name for a client request and from
    # the requesting username for a serving request; both are system:node:<name>.
    node_name="${common_name#system:node:}"
    case "${3}" in
        system:node:*) node_name="${3#system:node:}" ;;
        *) ;;
    esac
    log_value "node ${node_name} exists:" \
        "$(oc --kubeconfig="$KUBECONFIG" get node "${node_name}" -o jsonpath='{.metadata.name}' 2>/dev/null)"
    log_value "machine with InternalDNS:" \
        "$(oc --kubeconfig="$KUBECONFIG" get machines.machine.openshift.io --all-namespaces -o json 2>/dev/null |
            jq -r --arg n "${node_name}" '[.items[]? | select(any(.status.addresses[]?; .type == "InternalDNS" and (.address | sub("\\.$"; "") | ascii_downcase) == ($n | ascii_downcase))) | .metadata.name] | join(",")' 2>/dev/null)"
    log_value "machine with nodeRef:" \
        "$(oc --kubeconfig="$KUBECONFIG" get machines.machine.openshift.io --all-namespaces -o json 2>/dev/null |
            jq -r --arg n "${node_name}" '[.items[]? | select(.status.nodeRef.name == $n) | .metadata.name] | join(",")' 2>/dev/null)"
}

echo "Approving all CSR requests until bootstrapping is complete..."
# Names already logged, so each CSR is described once rather than on every pass.
logged=" "
last_cluster_inputs=""
while [ ! -f /opt/openshift/.bootkube.done ]
do
    current_cluster_inputs="$(cluster_inputs)"
    if [ "${current_cluster_inputs}" != "${last_cluster_inputs}" ]; then
        printf '%s\n' "${current_cluster_inputs}"
        last_cluster_inputs="${current_cluster_inputs}"
    fi

    # A CSR is pending when nothing has yet added an Approved or Denied
    # condition to it.
    while IFS=$'\t' read -r name signer username groups usages created conditions request
    do
        [ -n "${name}" ] || continue
        case "${logged}" in
            *" ${name} "*) continue ;;
            *) logged="${logged}${name} " ;;
        esac
        [ "${request}" = "<none>" ] && request=""
        log_csr "${name}" "${signer}" "${username}" "${groups}" "${usages}" \
            "${created}" "${conditions}" "$(printf '%s' "${request}" | base64 --decode 2>/dev/null)"
    done <<< "$(oc --kubeconfig="$KUBECONFIG" get csr -o json 2>/dev/null | jq -r '
        # Tab is an IFS whitespace character, so an empty field would be
        # swallowed by read and shift every field after it. Never emit one.
        def or_none: if . == null or . == "" then "<none>" else . end;
        .items[]?
        | select(
            ((.status.conditions // []) | map(.type)) as $types
            | ($types | index("Approved")) == null and ($types | index("Denied")) == null
          )
        | [
            .metadata.name,
            (.spec.signerName | or_none),
            (.spec.username | or_none),
            ((.spec.groups // []) | sort | join(",") | or_none),
            ((.spec.usages // []) | sort | join(",") | or_none),
            (.metadata.creationTimestamp | or_none),
            ((.status.conditions // []) | map(.type) | join(",") | or_none),
            (.spec.request | or_none)
          ]
        | @tsv')"

    oc --kubeconfig="$KUBECONFIG" get csr --no-headers | grep Pending | \
        awk '{print $1}' | \
        xargs --no-run-if-empty oc --kubeconfig="$KUBECONFIG" adm certificate approve
    sleep 20
done
